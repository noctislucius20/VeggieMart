package config

import (
	"context"
	"fmt"

	"user-service/database/seeds"
	"user-service/internal/core/domain/model"
	customLogger "user-service/utils/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Postgres struct {
	DB *gorm.DB
}

func (cfg Config) ConnectionPostgres(ctx context.Context) (*Postgres, error) {
	log := customLogger.NewLogger().Logger()

	dbConnString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?connect_timeout=%d",
		cfg.Psql.User,
		cfg.Psql.Password,
		cfg.Psql.Host,
		cfg.Psql.Port,
		cfg.Psql.DBName,
		cfg.Psql.DBConnectTimeout,
	)

	db, err := gorm.Open(postgres.Open(dbConnString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Errorf("[ConnectionPostgres-1] failed to connect with database %s", cfg.Psql.Host)
		return nil, err
	}

	go func() {
		sqlDB, err := db.DB()
		if err != nil {
			log.Errorf("[ConnectionPostgres-2] failed to connect with database")
			return
		}

		db.AutoMigrate(
			&model.User{},
			&model.Role{},
			&model.UserRole{},
			&model.VerificationToken{},
			&model.OutboxEvent{},
		)

		seeds.SeedRole(db)
		seeds.SeedAdmin(db)

		sqlDB.SetMaxOpenConns(cfg.Psql.DBMaxOpen)
		sqlDB.SetMaxIdleConns(cfg.Psql.DBMaxIdle)

		<-ctx.Done()

		sqlDB.Close()
		log.Infof("[ConnectionPostgres-3] database connection closed")
	}()

	return &Postgres{DB: db}, nil
}
