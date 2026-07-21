package app

import (
	"context"
	"os"
	"os/signal"
	"payment-service/config"
	"payment-service/internal/adapter/handler"
	httpclient "payment-service/internal/adapter/http_client"
	"payment-service/internal/adapter/repository"
	"payment-service/internal/adapter/repository/cache"
	"payment-service/internal/core/service"

	"payment-service/utils/logger"
	"payment-service/utils/validator"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10/translations/en"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RunServer() {
	serviceCtx, serviceCancel := context.WithCancel(context.Background())

	customValidator := validator.NewValidator()
	customLogger := logger.NewLogger()

	cfg := config.NewConfig()
	db, err := cfg.ConnectionPostgres(serviceCtx)
	if err != nil {
		customLogger.Logger().Fatalf("[RunServer-1] %v", err)
		return
	}

	redisClient, err := cfg.NewRedisClient(serviceCtx)
	if err != nil {
		customLogger.Logger().Fatalf("[RunServer-2] %v", err)
		return
	}

	e := echo.New()

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:        true,
		LogMethod:     true,
		LogStatus:     true,
		LogValuesFunc: customLogger.RequestLogger,
	}))

	en.RegisterDefaultTranslations(customValidator.Validator, customValidator.Translator)

	e.Validator = customValidator
	e.Logger = customLogger.Logger()

	httpClient := httpclient.NewHttpClient(cfg)

	txManager := repository.NewGormTransactionManager(db.DB)

	midtransClient := httpclient.NewMidtransclient(cfg, customLogger.Logger())

	paymentRepo := repository.NewPaymentRepository(db.DB, customLogger.Logger())
	outboxRepo := repository.NewOutboxEventRepository(db.DB, customLogger.Logger())

	cachePayment := cache.NewPaymentCache(redisClient, paymentRepo, customLogger.Logger())

	jwtService := service.NewJwtService(cfg)
	httpService := service.NewHttpService(cfg, httpClient)
	paymentService := service.NewPaymentService(paymentRepo, outboxRepo, cachePayment, cfg, httpService, midtransClient, txManager, customLogger.Logger())

	handler.NewPaymentHandler(paymentService, e, cfg, jwtService, redisClient)

	e.GET("/api/check", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	go func() {
		if cfg.App.AppPort == "" {
			cfg.App.AppPort = os.Getenv("APP_PORT")
		}

		err = e.Start(":" + cfg.App.AppPort)
		if err != nil {
			customLogger.Logger().Fatalf("[RunServer-4] %v", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	serviceCancel()

	customLogger.Logger().Infof("[RunServer-5] shutting down server on 5 seconds...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	e.Shutdown(ctx)
}
