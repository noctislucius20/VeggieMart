package cmd

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"user-service/config"
	"user-service/internal/adapter/message/publisher"
	"user-service/internal/adapter/repository"
	"user-service/utils/logger"

	"github.com/labstack/gommon/color"
	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
)

func startPublisherWorker() {
	customLogger := logger.NewLogger().Logger()

	conn, err := config.NewConfig().NewRabbitMQ()
	if err != nil {
		customLogger.Fatalf("[PublisherWorker-1] %v", err.Error())
	}

	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())

	cfg := config.NewConfig()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	db, err := cfg.ConnectionPostgres(ctx)
	if err != nil {
		log.Fatalf("[PublisherWorker-2] %v", err.Error())
	}

	outboxRepo := repository.NewOutboxEventRepository(db.DB, customLogger)

	wg.Go(func() {
		publisher.NewStartPublisherWorker(db.DB, conn, outboxRepo, customLogger).StartPublisherWorker(ctx)
	})

	<-quit

	cancel()

	wg.Wait()

	conn.Close()

	customLogger.Infof("[PublisherWorker-3] shutting down publisher worker...")
}

var workerPublisherCmd = &cobra.Command{
	Use:   "publisher-worker",
	Short: "Menjalankan worker untuk publish RabbitMQ",
	Run: func(cmd *cobra.Command, args []string) {
		color.Println(color.Green("Publisher worker is running..."))
		startPublisherWorker()
	},
}
