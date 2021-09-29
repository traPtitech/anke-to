package model

import (
	"context"
	"database/sql"
	"fmt"

	"gorm.io/gorm"
)

type ctxKey string

const (
	txKey ctxKey = "transaction"
)

// Transaction ITransactionの実装
type Transaction struct{}

func NewTransaction() *Transaction {
	return &Transaction{}
}

// Do トランザクション用の関数
func (*Transaction) Do(ctx context.Context, txOption *sql.TxOptions, f func(ctx context.Context) error) error {
	fc := func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, txKey, tx)

		err := f(ctx)
		if err != nil {
			return err
		}

		return nil
	}

	if txOption == nil {
		err := db.Transaction(fc)
		if err != nil {
			return fmt.Errorf("failed in transaction: %w", err)
		}
	} else {
		err := db.Transaction(fc, txOption)
		if err != nil {
			return fmt.Errorf("failed in transaction: %w", err)
		}
	}

	return nil
}

func getTx(ctx context.Context) (*gorm.DB, error) {
	iDB := ctx.Value(txKey)
	if iDB == nil {
		return db.Session(&gorm.Session{
			Context: ctx,
		}), nil
	}

	db, ok := iDB.(*gorm.DB)
	if !ok {
		return nil, ErrInvalidTx
	}

	return db.Session(&gorm.Session{
		Context: ctx,
	}), nil
}
