package worker

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
)

func StartPublisherWorker() {
	var (
		customLogger = logger.NewLogger().Logger()
		cfg          = config.NewConfig()
		wg           sync.WaitGroup
		ctx, cancel  = context.WithCancel(context.Background())
	)

	conn, err := config.NewConfig().NewRabbitMQ()
	if err != nil {
		customLogger.Fatalf("[PublisherWorker] %v", err)
	}

	db, err := cfg.ConnectionPostgres(ctx)
	if err != nil {
		customLogger.Fatalf("[PublisherWorker] %v", err)
		return
	}

	outboxRepo := repository.NewOutboxEventRepository(db.DB, customLogger)

	txManager := repository.NewGormTransactionManager(db.DB)

	wg.Go(func() {
		publisher.NewStartPublisherWorker(conn, outboxRepo, txManager, customLogger).StartPublisherWorker(ctx)
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	cancel()

	wg.Wait()

	conn.Close()

	customLogger.Infof("[PublisherWorker] shutting down publisher worker...")
}
