package model

import "github.com/srinathgs/mysqlstore"

type ISession interface {
	Get() (*mysqlstore.MySQLStore,error)
}
