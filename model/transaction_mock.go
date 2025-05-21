package model

import (
	"context"
	"fmt"
)

type MockTransaction struct{}

func (m *MockTransaction) Do(ctx context.Context, f func(ctx context.Context) error) error {
	err := f(ctx)
	if err != nil {
		return fmt.Errorf("failed in transaction: %w", err)
	}

	return nil
}
