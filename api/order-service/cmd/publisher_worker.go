package cmd

import (
	"context"
	"order-service/config"
	"order-service/internal/adapter/message/publisher"
	"order-service/internal/adapter/repository"
	"order-service/internal/core/service/transaction"
	"order-service/utils/logger"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/labstack/gommon/color"
	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
)

func startPublisherWorker() {
	var (
		customLogger = logger.NewLogger().Logger()
		cfg          = config.NewConfig()
		wg           sync.WaitGroup
		txManager    transaction.TransactionManager
	)

	conn, err := config.NewConfig().NewRabbitMQ()
	if err != nil {
		customLogger.Fatalf("[PublisherWorker-1] %v", err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	db, err := cfg.ConnectionPostgres(ctx)
	if err != nil {
		log.Fatalf("[PublisherWorker-2] %v", err.Error())
	}

	outboxRepo := repository.NewOutboxEventRepository(db.DB, customLogger)

	wg.Go(func() {
		publisher.NewStartPublisherWorker(conn, outboxRepo, txManager, customLogger).StartPublisherWorker(ctx)
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
