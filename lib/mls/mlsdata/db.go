package mlsdata

import (
	"context"
	"net/http"
	"os"

	"github.com/goware/lg"
	"github.com/upper/bond"

	"upper.io/db.v2/postgresql"
)

var DefaultDatabase *Database

type Database struct {
	bond.Session

	Listing ListingStore
}

func DatabaseCtx(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if DefaultDatabase == nil {
			if _, err := NewDBSession(); err != nil {
				lg.Alert("cannot open database connection: %+v", err)
				return
			}
		}
		context.WithValue(ctx, "database", DefaultDatabase)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func NewDBSession() (*Database, error) {
	var (
		err     error
		connURL postgresql.ConnectionURL
	)
	connURL, err = postgresql.ParseURL(os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	db := &Database{}
	db.Session, err = bond.Open(postgresql.Adapter, connURL)
	if err != nil {
		return nil, err
	}

	db.Listing = ListingStore{db.Store(&Listing{})}

	DefaultDatabase = db
	return db, nil
}
