package cmd

import (
	"context"
	"notification-service/config"
	"notification-service/internal/adapter/message"
	"notification-service/internal/adapter/message/consumer"
	"notification-service/internal/adapter/repository"
	"notification-service/internal/core/service"
	"notification-service/utils"
	"notification-service/utils/logger"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/labstack/gommon/color"
	"github.com/spf13/cobra"
)

func startConsumerWorker() {
	customLogger := logger.NewLogger().Logger()

	ctx, cancel := context.WithCancel(context.Background())

	cfg := config.NewConfig()

	db, err := cfg.ConnectionPostgres(ctx)
	if err != nil {
		customLogger.Fatalf("[StartConsumerWorker-1] %v", err.Error())
		return
	}

	conn, err := cfg.NewRabbitMQ()
	if err != nil {
		customLogger.Fatalf("[StartConsumerWorker-2] %v", err.Error())
		return
	}

	notificationRepo := repository.NewNotificationRepository(customLogger)

	emailService := message.NewEmailMessage(cfg)
	notificationService := service.NewNotificationService(notificationRepo, db.DB, customLogger)

	var wg sync.WaitGroup

	consumerWorker := consumer.NewNotificationConsumerWorker(emailService, notificationRepo, notificationService, conn, db.DB, customLogger)

	wg.Go(func() {
		consumerWorker.StartCreateNotificationWorker(ctx, utils.NOTIF_EMAIL_VERIFICATION)
	})

	wg.Go(func() {
		consumerWorker.StartCreateNotificationWorker(ctx, utils.NOTIF_EMAIL_FORGOT_PASSWORD)
	})

	wg.Go(func() {
		consumerWorker.StartCreateNotificationWorker(ctx, utils.NOTIF_EMAIL_CREATE_CUSTOMER)
	})

	wg.Go(func() {
		consumerWorker.StartCreateNotificationWorker(ctx, utils.NOTIF_EMAIL_UPDATE_ORDER_STATUS)
	})

	wg.Go(func() {
		consumerWorker.StartCreateNotificationWorker(ctx, utils.NOTIF_PUSH)
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	cancel()

	wg.Wait()

	conn.Close()

	customLogger.Infof("[StartConsumerWorker-2] shutting down consumer worker...")
}

var workerConsumerCmd = &cobra.Command{
	Use:   "consumer-worker",
	Short: "Menjalankan worker untuk consume RabbitMQ",
	Run: func(cmd *cobra.Command, args []string) {
		color.Println(color.Green("Consumer worker is running..."))
		startConsumerWorker()
	},
}
