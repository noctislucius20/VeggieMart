package app

import (
	"context"
	"notification-service/config"
	"notification-service/internal/adapter/handler"
	"notification-service/internal/adapter/repository"
	"notification-service/internal/core/service"
	"notification-service/utils/logger"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RunServer() {
	dbCtx, dbCancel := context.WithCancel(context.Background())

	customLogger := logger.NewLogger()

	cfg := config.NewConfig()

	db, err := cfg.ConnectionPostgres(dbCtx)
	if err != nil {
		customLogger.Logger().Fatalf("[RunServer-1] %v", err.Error())
		return
	}

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:        true,
		LogMethod:     true,
		LogStatus:     true,
		LogValuesFunc: customLogger.RequestLogger,
	}))

	notificationRepo := repository.NewNotificationRepository(customLogger.Logger())

	notificationService := service.NewNotificationService(notificationRepo, db.DB, customLogger.Logger())

	handler.NewNotificationHandler(notificationService, e, cfg)

	go func() {
		if cfg.App.AppPort == "" {
			cfg.App.AppPort = os.Getenv("APP_PORT")
		}

		err := e.Start(":" + cfg.App.AppPort)
		if err != nil {
			customLogger.Logger().Fatalf("[RunServer-2] %v", err.Error())
			return
		}
	}()

	var wg sync.WaitGroup

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	customLogger.Logger().Infof("[RunServer-3] shutting down server on 5 seconds...")

	dbCancel()

	wg.Wait()

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	e.Shutdown(ctx)

	cancel()
}
