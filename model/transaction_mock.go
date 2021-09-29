package model

import (
	"context"
	"database/sql"
	"fmt"
)

type MockTransaction struct{}

func (m *MockTransaction) Do(ctx context.Context, txOption *sql.TxOptions, f func(ctx context.Context) error) error {
	err := f(ctx)
	if err != nil {
		return fmt.Errorf("failed in transaction: %w", err)
	}

	return nil
}
