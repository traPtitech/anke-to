package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jinzhu/gorm"
)

type ctxKey string

const (
	txKey ctxKey = "transaction"
)

// Transaction ITransactionの実装
type Transaction struct{}

// Do トランザクション用の関数
func (*Transaction) Do(ctx context.Context, txOption *sql.TxOptions, f func(ctx context.Context) error) error {
	tx := db.BeginTx(ctx, txOption)

	ctx = context.WithValue(ctx, txKey, tx)

	err := f(ctx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed in transaction: %w", err)
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

func getTx(ctx context.Context) (*gorm.DB, error) {
	iDB := ctx.Value(txKey)
	if iDB == nil {
		return db, nil
	}

	db, ok := iDB.(*gorm.DB)
	if !ok {
		return nil, ErrInvalidTx
	}

	return db, nil
}
