package env

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // Postgres driver
)

type Db struct {
	User     string
	Password string
	Database string
	sql      *sql.DB
}

type Env struct {
	Db     *Datastore
	logger log.Logger
}

type Datastore interface {
	playersDatastore
}

func (db *Db) Init() {
	sql, err := sql.Open("postgres", "user=pqgotest dbname=pqgotest sslmode=verify-full")
	if err != nil {
		log.Fatal(err)
	}

	db.sql = sql
}
