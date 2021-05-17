package model

import (
	"fmt"
	"os"

	"github.com/srinathgs/mysqlstore"
)

type Session struct{}

func NewSession() *Session {
	return &Session{}
}

func (*Session) Get() (*mysqlstore.MySQLStore, error) {
	store, err := mysqlstore.NewMySQLStoreFromConnection(db.DB(), "sessions", "/", 60*60*24*14, []byte(os.Getenv("SESSION_SECRET")))
	if err != nil {
		return nil, fmt.Errorf("failed to create session store: %w", err)
	}

	return store, nil
}
