package data

import (
	"github.com/upper/bond"
	"upper.io/db.v2/postgresql"
)

type Database struct {
	bond.Session

	Profile ProfileStore
}

func NewDBSession() (*Database, error) {
	var (
		err     error
		connURL postgresql.ConnectionURL
	)
	connURL, err = postgresql.ParseURL("postgres://localhost/jarvis")
	if err != nil {
		return nil, err
	}

	db := &Database{}
	db.Session, err = bond.Open(postgresql.Adapter, connURL)
	if err != nil {
		return nil, err
	}

	db.Profile = ProfileStore{db.Store(&Profile{})}

	return db, nil
}
