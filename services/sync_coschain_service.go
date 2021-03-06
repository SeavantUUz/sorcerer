package services

import (
	"context"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"sorcerer/constants"
	grpcpb "sorcerer/rpc"
	"sorcerer/structure"
	"sorcerer/util"
	"sync"
	"sync/atomic"
	"time"
)

type SyncProgress struct {
	ID          uint64 `gorm:"primary_key;auto_increment"`
	BlockHeight uint64
	FinishAt    time.Time
}

type SyncService struct {
	sync.Mutex
	log      *logrus.Logger
	db       *gorm.DB
	jobTimer *time.Timer
	stop     int32
	working  int32
	workStop *sync.Cond
	client   grpcpb.ApiServiceClient
	trxsCh   chan<- *structure.Transaction
}

func NewSyncService(log *logrus.Logger, ch chan<- *structure.Transaction) (*SyncService, error) {
	s := &SyncService{log: log}
	s.db, s.stop, s.working = nil, 0, 0
	s.trxsCh = ch
	s.workStop = sync.NewCond(&s.Mutex)
	return s, nil
}

func (s *SyncService) Start() error {
	err := s.initDatabase()
	if err != nil {
		return err
	}
	err = s.initRpc()
	if err != nil {
		return err
	}
	s.scheduleNextJob()
	return nil
}

func (s *SyncService) Stop() error {
	s.waitWorkDone()
	if s.db != nil {
		_ = s.db.Close()
	}
	s.db, s.stop, s.working = nil, 0, 0
	s.log.Infoln("Sync Service Exit Success")
	return nil
}

func (s *SyncService) initDatabase() error {
	if db, err := gorm.Open("mysql", constants.DSN); err != nil {
		return err
	} else {
		s.db = db
	}
	if !s.db.HasTable(&SyncProgress{}) {
		if err := s.db.CreateTable(&SyncProgress{}).Error; err != nil {
			_ = s.db.Close()
			return err
		}
	}
	return nil
}

func (s *SyncService) initRpc() error {
	conn, err := util.Dial(constants.Node)
	if err != nil {
		return err
	}
	client := grpcpb.NewApiServiceClient(conn)
	s.client = client
	return nil
}

func (s *SyncService) scheduleNextJob() {
	s.jobTimer = time.AfterFunc(1*time.Second, s.work)
}

func (s *SyncService) work() {
	var (
		userBreak = false
		err       error
	)
	atomic.StoreInt32(&s.working, 1)
	chainInfo, err := s.chainInfo()
	if err != nil {
		s.log.Error(err)
	}
	headBlockNum := chainInfo.State.Dgpo.HeadBlockNumber
	progress := &SyncProgress{}
	// localstorage is empty, write current info into storage
	if s.db.Where(&SyncProgress{}).First(progress).RecordNotFound() {
		progress.BlockHeight = headBlockNum
		progress.FinishAt = time.Now()
		tx := s.db.Begin()
		if err = tx.Save(progress).Error; err == nil {
			tx.Commit()
		} else {
			s.log.Error(err)
			tx.Rollback()
		}
	}
	if atomic.LoadInt32(&s.stop) != 0 {
		userBreak = true
	}
	if err == nil {
		minBlockNum, maxBlockNum := progress.BlockHeight+1, headBlockNum
		for blockNum := minBlockNum; blockNum <= maxBlockNum; blockNum++ {
			if atomic.LoadInt32(&s.stop) != 0 {
				userBreak = true
				break
			}
			blockInfo, err := s.blockInfo(blockNum)
			if err != nil {
				s.log.Error(err)
				break
			}
			tx := s.db.Begin()
			trxs := util.ExtractTransactions(blockInfo.Block)
			for _, trx := range trxs {
				s.trxsCh <- trx
			}
			progress.BlockHeight = blockNum
			progress.FinishAt = time.Now()
			if err := tx.Save(progress).Error; err == nil {
				tx.Commit()
			} else {
				s.log.Error(err)
				tx.Rollback()
				break
			}
		}
	}
	s.Lock()
	atomic.StoreInt32(&s.working, 0)
	if !userBreak {
		s.scheduleNextJob()
	}
	s.workStop.Signal()
	s.Unlock()
}

func (s *SyncService) chainInfo() (*grpcpb.GetChainStateResponse, error) {
	chainState, err := s.client.GetChainState(context.Background(), &grpcpb.NonParamsRequest{})
	return chainState, err
}

func (s *SyncService) blockInfo(blockHeight uint64) (*grpcpb.GetSignedBlockResponse, error) {
	request := grpcpb.GetSignedBlockRequest{Start: blockHeight}
	blockInfo, err := s.client.GetSignedBlock(context.Background(), &request)
	return blockInfo, err
}

func (s *SyncService) waitWorkDone() {
	s.Lock()
	if s.jobTimer != nil {
		s.jobTimer.Stop()
	}
	atomic.StoreInt32(&s.stop, 1)
	for atomic.LoadInt32(&s.working) != 0 {
		s.workStop.Wait()
	}
	s.Unlock()
}
