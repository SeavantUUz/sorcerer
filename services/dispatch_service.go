package services

import (
	"errors"
	"github.com/sirupsen/logrus"
	"sorcerer/structure"
	"sorcerer/watcher"
	"sync"
	"sync/atomic"
)

type PublishService struct {
	sync.Mutex
	log *logrus.Logger
	//publisher     *philote.Client
	working       int32
	stop          int32
	workCron      *sync.Cond
	trxsChan      <-chan []*structure.Transaction
	eventWatchers []watcher.Watcher
}

func NewPublishService(log *logrus.Logger, ch <-chan []*structure.Transaction) (*PublishService, error) {
	service := &PublishService{}
	service.log = log
	service.working = 0
	service.stop = 0
	service.trxsChan = ch
	if contractApplyWatcher, err := watcher.NewContractApplyWatcher(log); err == nil {
		service.eventWatchers = append(service.eventWatchers, contractApplyWatcher)
	} else {
		msg := "contract apply watcher initialize failed"
		log.Error(msg)
		return nil, errors.New(msg)
	}
	return service, nil
}

func (s *PublishService) Start() error {
	go func() {
		for !(atomic.LoadInt32(&s.stop) > 0) {
			select {
			case trxs := <-s.trxsChan:
				atomic.StoreInt32(&s.working, 1)
				for _, trx := range trxs {
					for _, eventWatcher := range s.eventWatchers {
						if eventWatcher.Agree(trx) {

							eventWatcher.Publish(trx)
						}
					}

				}
				s.Lock()
				atomic.StoreInt32(&s.working, 0)
				s.workCron.Signal()
				s.Unlock()
			}
		}
	}()
	return nil
}

func (s *PublishService) Stop() error {
	s.Lock()
	atomic.StoreInt32(&s.stop, 1)
	for atomic.LoadInt32(&s.working) != 0 {
		s.workCron.Wait()
	}
	s.Unlock()
	return nil
}
