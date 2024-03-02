package main

import (
	"advertiser/channel_owner/internal/dep_container"
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	container, err := dep_container.New()
	if err != nil {
		panic(fmt.Sprintf("error initializing DI container: %+v", err))
	}

	go container.RunChannelListener()
	go container.RunNotificationService()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	<-exit

	if err = container.Delete(); err != nil {
		zap.S().Error("error deleting DI container", zap.Error(err))
	}
}
