package worker

import (
	"context"
	"order-service/config"
	"order-service/internal/adapter/message/consumer"
	"order-service/internal/adapter/repository"
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
		customLogger.Fatalf("[StartConsumerWorker] %v", err)
		return
	}

	db, err := cfg.ConnectionPostgres(ctx)
	if err != nil {
		customLogger.Fatalf("[StartConsumerWorker] %v", err)
		return
	}

	esClient, err := cfg.NewElasticsearchClient()
	if err != nil {
		customLogger.Fatalf("[StartConsumerWorker] %v", err)
		return
	}

	orderRepo := repository.NewOrderRepository(db.DB, customLogger)
	txManager := repository.NewGormTransactionManager(db.DB)

	consumerWorker := consumer.NewOrderConsumerWorker(conn, esClient, orderRepo, txManager, cfg, customLogger)

	wg.Go(func() {
		consumerWorker.StartCreateOrderWorker(ctx)
	})

	wg.Go(func() {
		consumerWorker.StartOrderPaymentSuccessWorker(ctx)
	})

	wg.Go(func() {
		consumerWorker.StartElasticUpdateStatusOrderWorker(ctx)
	})

	wg.Go(func() {
		consumerWorker.StartDbUpdateStatusOrderWorker(ctx)
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	cancel()

	wg.Wait()

	conn.Close()

	customLogger.Infof("[StartConsumerWorker] shutting down consumer worker...")
}
