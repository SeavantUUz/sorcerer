package main

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sorcerer/services"
	"sorcerer/structure"
	"sync"
	"syscall"
)

func main() {
	stop := make(chan bool)
	c := make(chan os.Signal)
	logger := logrus.New()
	ch := make(chan *structure.Transaction)
	publishService, err := services.NewPublishService(logger, ch)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	wg := sync.WaitGroup{}
	syncService, err := services.NewSyncService(logger, ch)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	go func() {
		wg.Add(1)
		defer wg.Done()
		if err := syncService.Start(); err != nil {
			logger.Errorf("Fatal: %v \n", err)
		}
	}()
	go func() {
		wg.Add(1)
		defer wg.Done()
		if err := publishService.Start(); err != nil {
			logger.Errorf("Fatal: %v \n", err)
		}
	}()
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				_ = syncService.Stop()
				_ = publishService.Stop()
				stop <- true
			default:
				logger.Infoln("unknown terminate exit")
			}
		}
	}()
	<-stop
	wg.Wait()
	logger.Infoln("exit successful")
}
