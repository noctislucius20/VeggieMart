package worker

import (
	"context"
	"order-service/config"
	"order-service/internal/adapter/message/consumer"
	"order-service/utils/logger"
	"os"
	"os/signal"
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
		customLogger.Fatalf("[StartConsumerWorker-1] %v", err)
		return
	}

	esClient, err := cfg.NewElasticsearchClient()
	if err != nil {
		customLogger.Fatalf("[StartConsumerWorker-2] %v", err)
		return
	}

	consumerWorker := consumer.NewOrderConsumerWorker(conn, esClient, cfg, customLogger)

	wg.Go(func() {
		consumerWorker.StartCreateOrderWorker(ctx)
	})

	wg.Go(func() {
		consumerWorker.StartOrderPaymentSuccessWorker(ctx)
	})

	wg.Go(func() {
		consumerWorker.StartUpdateStatusOrderWorker(ctx)
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	cancel()

	wg.Wait()

	conn.Close()

	customLogger.Infof("[StartConsumerWorker-3] shutting down consumer worker...")
}
