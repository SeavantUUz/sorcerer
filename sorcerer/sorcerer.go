package sorcerer

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sorcerer/services"
	"syscall"
)

func main() {
	stop := make(chan bool)
	c := make(chan os.Signal)
	logger := logrus.New()
	publishService, err := services.NewPublishService(logger)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	syncService, err := services.NewSyncService(logger)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	if err := syncService.Start(); err != nil {
		_, _ = fmt.Printf("Fatal: %v \n", err)
		os.Exit(1)
	}
	if err := publishService.Start(); err != nil {
		_, _ = fmt.Printf("Fatal: %v \n", err)
		os.Exit(1)
	}
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				_ = syncService.Stop()
				_ = publishService.Stop()
				stop <- true
				fmt.Println("success exit", s)
			default:
				fmt.Println("terminate exit", s)
			}
		}
	}()
	<-stop
}
