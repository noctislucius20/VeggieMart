package cmd

import (
	"context"
	"os"
	"os/signal"
	"product-service/config"
	"product-service/internal/adapter/message/consumer"
	"sync"
	"syscall"

	"github.com/labstack/gommon/color"
	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
)

func startConsumerWorker() {
	cfg := config.NewConfig()

	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())

	consumerWorker := consumer.NewProductConsumerWorker(cfg)

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

	log.Infof("[StartConsumerWorker-1] shutting down consumer worker...")
}

var workerConsumerCmd = &cobra.Command{
	Use:   "consumer-worker",
	Short: "Menjalankan worker untuk consume RabbitMQ",
	Run: func(cmd *cobra.Command, args []string) {
		color.Println(color.Green("Consumer worker is running..."))
		startConsumerWorker()
	},
}
