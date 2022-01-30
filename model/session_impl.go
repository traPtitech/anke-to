package model

import (
	"fmt"
	"github.com/srinathgs/mysqlstore"
	"os"
)

type Session struct {
}

func (s *Session) Get() (*mysqlstore.MySQLStore, error) {
	_db, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB :%w", err)
	}

	store, err := mysqlstore.NewMySQLStoreFromConnection(_db, "sessions", "/", 60*60*24*14, []byte(os.Getenv("SESSION_SECRET")))
	if err != nil {
		return nil, fmt.Errorf("failed to create session store:%w", err)
	}

	return store, nil
}

func NewSession() *Session {
	return &Session{}
}
