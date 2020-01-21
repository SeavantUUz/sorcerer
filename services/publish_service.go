package services

import (
	"github.com/pote/philote-go"
	"github.com/sirupsen/logrus"
	"sorcerer/constants"
	"sorcerer/structure"
	"sorcerer/watcher"
	"sync"
	"sync/atomic"
)

type PublishService struct {
	sync.Mutex
	log           *logrus.Logger
	publisher     *philote.Client
	working       int32
	stop          int32
	workCron      *sync.Cond
	trxsChan      <-chan []*structure.Transaction
	eventWatchers []watcher.Watcher
}

func NewPublishService(log *logrus.Logger, ch <-chan []*structure.Transaction) (*PublishService, error) {
	c, err := philote.NewClient("ws://localhost:6380", constants.TokenContractRW)
	if err != nil {
		return nil, err
	}
	service := &PublishService{log: log, publisher: c}
	service.working = 0
	service.stop = 0
	service.trxsChan = ch
	service.eventWatchers = []watcher.Watcher{watcher.NewContractApplyWatcher()}
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
							data := eventWatcher.Do(trx)
							_ = s.publisher.Publish(&philote.Message{
								Channel: eventWatcher.Channel(),
								Data:    data,
							})
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
