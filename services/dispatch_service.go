package services

import (
	"errors"
	"github.com/sirupsen/logrus"
	"sorcerer/structure"
	"sorcerer/watcher"
	"sync"
	"sync/atomic"
	"time"
)

type DispatchService struct {
	sync.Mutex
	log           *logrus.Logger
	working       int32
	stop          int32
	ticker        *time.Ticker
	workCron      *sync.Cond
	trxsChan      <-chan *structure.Transaction
	eventWatchers []watcher.Watcher
}

func NewPublishService(log *logrus.Logger, ch <-chan *structure.Transaction) (*DispatchService, error) {
	service := &DispatchService{log: log, working: 0, stop: 0, trxsChan: ch, ticker: time.NewTicker(500 * time.Millisecond)}
	service.workCron = sync.NewCond(&service.Mutex)
	if contractApplyWatcher, err := watcher.NewContractApplyWatcher(log); err == nil {
		service.eventWatchers = append(service.eventWatchers, contractApplyWatcher)
	} else {
		msg := "contract apply watcher initialize failed"
		log.Error(msg)
		return nil, errors.New(msg)
	}
	return service, nil
}

func (s *DispatchService) Start() error {
	for {
		select {
		case trx := <-s.trxsChan:
			atomic.StoreInt32(&s.working, 1)
			for _, eventWatcher := range s.eventWatchers {
				if eventWatcher.Agree(trx) {
					eventWatcher.Publish(trx)
				}
			}
			s.Lock()
			atomic.StoreInt32(&s.working, 0)
			s.workCron.Signal()
			s.Unlock()
		case <-s.ticker.C:
			if !(atomic.LoadInt32(&s.stop) == 0) {
				s.log.Infoln("Dispatch Service Exit Successful")
				return nil
			}
		}
	}
}

func (s *DispatchService) Stop() error {
	s.Lock()
	atomic.StoreInt32(&s.stop, 1)
	for atomic.LoadInt32(&s.working) != 0 {
		s.workCron.Wait()
	}
	s.Unlock()
	return nil
}
