package app

import (
	"context"
	"os"
	"os/signal"
	"payment-service/config"
	"payment-service/internal/adapter/handler"
	httpclient "payment-service/internal/adapter/http_client"
	"payment-service/internal/adapter/repository"
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
	dbCtx, dbCancel := context.WithCancel(context.Background())

	customValidator := validator.NewValidator()
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

	en.RegisterDefaultTranslations(customValidator.Validator, customValidator.Translator)

	e.Validator = customValidator
	e.Logger = customLogger.Logger()

	httpClient := httpclient.NewHttpClient(cfg)
	midtransClient := httpclient.NewMidtransclient(cfg, customLogger.Logger())

	paymentRepo := repository.NewPaymentRepository(customLogger.Logger())
	outboxRepo := repository.NewOutboxEventRepository(customLogger.Logger())

	httpService := service.NewHttpService(cfg, httpClient)
	paymentService := service.NewPaymentService(paymentRepo, outboxRepo, cfg, httpService, midtransClient, db.DB, customLogger.Logger())

	handler.NewPaymentHandler(paymentService, e, cfg)

	e.GET("/api/check", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	go func() {
		if cfg.App.AppPort == "" {
			cfg.App.AppPort = os.Getenv("APP_PORT")
		}

		err = e.Start(":" + cfg.App.AppPort)
		if err != nil {
			customLogger.Logger().Fatalf("[RunServer-4] %v", err.Error())
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	dbCancel()

	customLogger.Logger().Infof("[RunServer-5] shutting down server on 5 seconds...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	e.Shutdown(ctx)
}
