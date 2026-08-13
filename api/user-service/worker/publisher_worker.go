package worker

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"user-service/config"
	"user-service/internal/adapter/message/publisher"
	"user-service/internal/adapter/repository"
	"user-service/utils/logger"
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
	}

	outboxRepo := repository.NewOutboxEventRepository(db.DB, customLogger)

	txManager := repository.NewGormTransactionManager(db.DB)

	wg.Go(func() {
		publisher.NewStartPublisherWorker(conn, txManager, outboxRepo, customLogger).StartPublisherWorker(ctx)
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	cancel()

	wg.Wait()

	conn.Close()

	customLogger.Infof("[PublisherWorker] shutting down publisher worker...")
}
