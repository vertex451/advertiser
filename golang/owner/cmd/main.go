package main

import (
	"advertiser/owner/internal/dep_container"
	"advertiser/shared/config/config"
	"advertiser/shared/pkg/service/repo"
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Fixes ERROR: duplicate key value violates unique constraint "pg_class_relname_nsp_index" (SQLSTATE 23505)
	repo.New(config.Load())
	// end of fix

	container, err := dep_container.New()
	if err != nil {
		panic(fmt.Sprintf("error initializing DI container: %+v", err))
	}

	go container.MonitorChannels()
	go container.RunNotificationService()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	<-exit

	if err = container.Delete(); err != nil {
		zap.S().Error("error deleting DI container", zap.Error(err))
	}
}
