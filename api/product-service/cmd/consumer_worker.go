package cmd

import (
	"context"
	"os"
	"os/signal"
	"product-service/config"
	"product-service/internal/adapter/message/consumer"
	"product-service/utils/logger"
	"sync"
	"syscall"

	"github.com/labstack/gommon/color"
	"github.com/spf13/cobra"
)

func startConsumerWorker() {
	var (
		customLogger = logger.NewLogger().Logger()
		cfg          = config.NewConfig()
		wg           sync.WaitGroup
		ctx, cancel  = context.WithCancel(context.Background())
	)

	conn, err := cfg.NewRabbitMQ()
	if err != nil {
		customLogger.Fatalf("[StartConsumerWorker-1] %v", err)
		return
	}

	esClient, err := cfg.NewElasticsearchClient()
	if err != nil {
		customLogger.Fatalf("[StartConsumerWorker-2] %v", err)
		return
	}

	consumerWorker := consumer.NewProductConsumerWorker(conn, esClient, cfg, customLogger)

	wg.Go(func() {
		consumerWorker.StartCreateProductWorker(ctx)
	})

	wg.Go(func() {
		consumerWorker.StartUpdateProductWorker(ctx)
	})

	wg.Go(func() {
		consumerWorker.StartDeleteProductWorker(ctx)
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	cancel()

	wg.Wait()

	customLogger.Infof("[StartConsumerWorker-3] shutting down consumer worker...")
}

var workerConsumerCmd = &cobra.Command{
	Use:   "consumer-worker",
	Short: "Menjalankan worker untuk consume RabbitMQ",
	Run: func(cmd *cobra.Command, args []string) {
		color.Println(color.Green("Consumer worker is running..."))
		startConsumerWorker()
	},
}
