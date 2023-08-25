package main

import (
	"context"
	"flag"
	zap_logger "github.com/romandnk/dynamic-user-segmentation-service/internal/logger/zap"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/storage/postgres"
	"go.uber.org/zap"
	"log"
	"os/signal"
	"syscall"
)

var configFile string

func main() {
	flag.StringVar(&configFile, "config", "./configs/dynamic-user-segmentation.yaml", "Path to configuration file")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	// initialize config
	config, err := NewConfig(configFile)
	if err != nil {
		if err != nil {
			log.Fatalf("error initializing config: %s", err.Error())
		}
	}

	// initialize zap logger
	logg, err := zap_logger.NewZapLogger(config.ZapLogger)
	if err != nil {
		log.Fatalf("error initializing zap logger: %s", err.Error())
	}

	logg.Info("using zap logger")

	// initialize postgres storage
	postgresStorage := postgres.NewStoragePostgres()
	err = postgresStorage.Connect(ctx, config.Postgres)
	defer postgresStorage.Close()
	if err != nil {
		logg.Error("error connecting postgres db", zap.String("error", err.Error()))
		return
	}

	logg.Info("using postgres storage")
}
