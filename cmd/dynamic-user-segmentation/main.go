package main

import (
	"context"
	"flag"
	zap_logger "github.com/romandnk/dynamic-user-segmentation-service/internal/logger/zap"
	http_server "github.com/romandnk/dynamic-user-segmentation-service/internal/server/http"
	v1 "github.com/romandnk/dynamic-user-segmentation-service/internal/server/http/v1"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/service"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/storage/postgres"
	"go.uber.org/zap"
	"log"
	"os/signal"
	"syscall"
	"time"
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

	// initialize services
	services := service.NewService(postgresStorage)

	// initialize http handler
	handler := v1.NewHandler(services, logg)

	// initialize http server
	server := http_server.NewServer(config.Server, handler.InitRoutes())

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		logg.Info("stopping dynamic user segmentation service...")

		// stopping http server in three seconds
		if err := server.Stop(ctx); err != nil {
			logg.Error("error stopping dynamic user segmentation service", zap.String("error", err.Error()))
			return
		}

		logg.Info("dynamic-user-segmentation is stopped")
	}()

	logg.Info("starting dynamic user segmentation service...")

	// starting http server
	if err := server.Start(); err != nil {
		logg.Error("error dynamic user segmentation service", zap.String("error", err.Error()))
		return
	}
}
