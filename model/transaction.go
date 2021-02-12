package model

import (
	"context"
	"database/sql"
)

// ITransaction Transaction処理のinterface
type ITransaction interface {
	Do(context.Context, *sql.TxOptions, func(context.Context) error) error
}
