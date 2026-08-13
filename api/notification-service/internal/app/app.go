package app

import (
	"context"
	"notification-service/config"
	"notification-service/internal/adapter/handler"
	"notification-service/internal/adapter/repository"
	"notification-service/internal/adapter/repository/cache"
	"notification-service/internal/core/service"
	"notification-service/utils/logger"
	"notification-service/utils/ws"

	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RunServer() {
	serviceCtx, serviceCancel := context.WithCancel(context.Background())

	customLogger := logger.NewLogger()

	cfg := config.NewConfig()

	db, err := cfg.ConnectionPostgres(serviceCtx)
	if err != nil {
		customLogger.Logger().Fatalf("[RunServer] %v", err)
		return
	}

	redisClient, err := cfg.NewRedisClient(serviceCtx)
	if err != nil {
		customLogger.Logger().Fatalf("[RunServer] %v", err)
		return
	}

	e := echo.New()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:        true,
		LogMethod:     true,
		LogStatus:     true,
		LogValuesFunc: customLogger.RequestLogger,
	}))

	txManager := repository.NewGormTransactionManager(db.DB)

	notificationRepo := repository.NewNotificationRepository(db.DB, customLogger.Logger())

	cacheNotification := cache.NewNotificationCache(redisClient, customLogger.Logger())

	jwtService := service.NewJwtService(cfg)
	notificationService := service.NewNotificationService(notificationRepo, cacheNotification, txManager, db.DB, customLogger.Logger())

	handler.NewNotificationHandler(e, cfg, notificationService, jwtService, redisClient)
	handler.NewWebSocketHandler(e, cfg, jwtService, redisClient, customLogger.Logger())

	go ws.SubscribeWebSocketChannel(serviceCtx, redisClient, notificationService, customLogger.Logger())
	go ws.SubscribeAdminWebSocketChannel(serviceCtx, redisClient, customLogger.Logger())

	e.GET("/api/check", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	go func() {
		if cfg.App.AppPort == "" {
			cfg.App.AppPort = os.Getenv("APP_PORT")
		}

		err := e.Start(":" + cfg.App.AppPort)
		if err != nil {
			customLogger.Logger().Fatalf("[RunServer] %v", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	customLogger.Logger().Infof("[RunServer] shutting down server on 5 seconds...")

	serviceCancel()

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	e.Shutdown(ctx)
}
