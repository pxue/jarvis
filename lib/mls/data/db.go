package data

import (
	"github.com/upper/bond"

	"upper.io/db.v2/postgresql"
)

type Database struct {
	bond.Session

	Listing ListingStore
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

	db.Listing = ListingStore{db.Store(&Listing{})}

	return db, nil
}
