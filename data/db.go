package data

import (
	"os"

	"github.com/upper/bond"

	"upper.io/db.v2"
	"upper.io/db.v2/postgresql"
)

type Database struct {
	bond.Session
}

type DBConf struct {
	DebugQueries bool `toml:"debug_queries"`
}

// String implements db.ConnectionURL
func (cf *DBConf) String() string {
	return os.Getenv("DATABASE_URL")
}

func NewDBSession(conf *DBConf) (*Database, error) {
	if conf.DebugQueries {
		db.Conf.SetLogging(true)
	}

	var err error
	db := &Database{}
	db.Session, err = bond.Open(postgresql.Adapter, conf)
	if err != nil {
		return nil, err
	}

	return db, nil
}
