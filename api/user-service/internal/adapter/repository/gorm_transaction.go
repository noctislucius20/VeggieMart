package repository

import (
	"context"
	"user-service/internal/core/service/transaction"

	"gorm.io/gorm"
)

type txKey struct{}

type gormTransactionManager struct {
	db *gorm.DB
}

func (d *gormTransactionManager) WithinTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, txKey{}, tx)
		return fn(txCtx)
	})
}

func NewGormTransactionManager(db *gorm.DB) transaction.TransactionManager {
	return &gormTransactionManager{
		db: db,
	}
}
