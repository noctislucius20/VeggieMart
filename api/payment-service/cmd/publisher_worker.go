package cmd

import (
	"context"
	"os"
	"os/signal"
	"payment-service/config"
	"payment-service/internal/adapter/message/publisher"
	"payment-service/internal/adapter/repository"
	"payment-service/utils/logger"
	"sync"
	"syscall"

	"github.com/labstack/gommon/color"
	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
)

func startPublisherWorker() {
	var wg sync.WaitGroup
	defer wg.Wait()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.NewConfig()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	db, err := cfg.ConnectionPostgres(ctx)
	if err != nil {
		log.Fatalf("[PublisherWorker-1] %v", err.Error())
		return
	}

	customLogger := logger.NewLogger().Logger()

	outboxRepo := repository.NewOutboxEventRepository(customLogger)

	wg.Go(func() {
		publisher.NewStartPublisherWorker(db.DB, outboxRepo, customLogger).StartPublisherWorker(ctx)
	})

	<-quit
	customLogger.Infof("[PublisherWorker-2] shutting down publisher worker...")
}

var workerPublisherCmd = &cobra.Command{
	Use:   "publisher-worker",
	Short: "Menjalankan worker untuk publish RabbitMQ",
	Run: func(cmd *cobra.Command, args []string) {
		color.Println(color.Green("Publisher worker is running..."))
		startPublisherWorker()
	},
}
