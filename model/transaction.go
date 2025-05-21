//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package model

import (
	"context"
	"database/sql"
)

// ITransaction Transaction処理のinterface
type ITransaction interface {
	Do(context.Context, *sql.TxOptions, func(context.Context) error) error
}
