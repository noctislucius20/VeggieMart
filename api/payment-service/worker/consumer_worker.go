package worker

import (
	"context"
	"os"
	"os/signal"
	"payment-service/config"
	"payment-service/internal/adapter/message/consumer"
	"payment-service/internal/adapter/repository"
	"payment-service/internal/adapter/repository/cache"
	"payment-service/utils/logger"
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

	db, err := cfg.ConnectionPostgres(ctx)
	if err != nil {
		customLogger.Fatalf("[StartConsumerWorker-2] %v", err)
		return
	}

	redisClient, err := cfg.NewRedisClient(ctx)
	if err != nil {
		customLogger.Fatalf("[StartConsumerWorker-3] %v", err)
		return
	}

	paymentRepo := repository.NewPaymentRepository(db.DB, customLogger)
	txManager := repository.NewGormTransactionManager(db.DB)
	cachePayment := cache.NewPaymentCache(redisClient, paymentRepo, customLogger)
	consumerWorker := consumer.NewPaymentConsumerWorker(conn, db.DB, cfg, cachePayment, paymentRepo, txManager, customLogger)

	wg.Go(func() {
		consumerWorker.StartCreatePaymentWorker(ctx)
	})

	wg.Go(func() {
		consumerWorker.StartUpdatePaymentWorker(ctx)
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	cancel()

	wg.Wait()

	conn.Close()

	customLogger.Infof("[StartConsumerWorker-4] shutting down consumer worker...")
}
