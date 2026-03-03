package cmd

import (
	"context"
	"os"
	"os/signal"
	"product-service/config"
	"product-service/internal/adapter/message/publisher"
	"product-service/internal/adapter/repository"
	"product-service/utils/logger"
	"sync"
	"syscall"

	"github.com/labstack/gommon/color"
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
		customLogger.Fatalf("[PublisherWorker-2] %v", err)
		return
	}

	outboxRepo := repository.NewOutboxEventRepository(db.DB, customLogger)

	wg.Go(func() {
		publisher.NewStartPublisherWorker(db.DB, conn, outboxRepo, customLogger).StartPublisherWorker(ctx)
	})

	<-quit

	cancel()

	wg.Wait()

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
