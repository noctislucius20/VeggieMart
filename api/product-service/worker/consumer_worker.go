package worker

import (
	"context"
	"os"
	"os/signal"
	"product-service/config"
	"product-service/internal/adapter/message/consumer"
	"product-service/utils/logger"
	"sync"
	"syscall"
)

func StartConsumerWorker() {
	var (
		customLogger = logger.NewLogger().Logger()
		cfg          = config.NewConfig()
		wg           sync.WaitGroup
		ctx, cancel  = context.WithCancel(context.Background())
	)

	conn, err := cfg.NewRabbitMQ()
	if err != nil {
		customLogger.Fatalf("[StartConsumerWorker] %v", err)
		return
	}

	esClient, err := cfg.NewElasticsearchClient()
	if err != nil {
		customLogger.Fatalf("[StartConsumerWorker] %v", err)
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

	customLogger.Infof("[StartConsumerWorker] shutting down consumer worker...")
}
