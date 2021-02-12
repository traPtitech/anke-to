package model

import "context"

// ITransaction Transaction処理のinterface
type ITransaction interface {
	Do(context.Context, func(context.Context) error) error
}
