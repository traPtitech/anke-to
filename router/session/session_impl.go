package session

import (
	"fmt"
	"github.com/srinathgs/mysqlstore"
	"github.com/traPtitech/anke-to/model"
)

type Store struct {
	store *mysqlstore.MySQLStore
}

func NewStore(sess model.Session) (*Store,error) {
	store,err := sess.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	return &Store{store: store},nil
}
