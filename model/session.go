//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package model

import "github.com/srinathgs/mysqlstore"

type ISession interface {
	Get() (*mysqlstore.MySQLStore,error)
}
