package env

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // Postgres driver
)

type db struct {
	*sql.DB
	*log.Logger
}

type Env struct {
	db  datastore
	log *log.Logger
}

type datastore interface {
	players
}

func New() *Env {
	/*logger := log.New(os.Stdout, "log: ", log.Lshortfile)

	loggerDb := log.New(os.Stdout, "db: ", log.Lshortfile)
	sql, err := sql.Open("postgres", "")
	if err != nil {
		loggerDb.Fatal(err)
	}

	loggerDb.Printf("Connected to database")
	a := &db{sql, loggerDb}*/
	v := Env{nil, nil}
	return &v //&Env{a, logger}
}

func DbString(database, server, username, password string) string {
	return fmt.Sprintf("postgres://%v:%v@%v/%v", username, password, server, database)
}

/*
func (env *Env) ConnectDb(dbString string) {
	loggerDb := log.New(os.Stdout, "db: ", log.Lshortfile)
	sql, err := sql.Open("postgres", dbString)
	if err != nil {
		loggerDb.Fatal(err)
	}

	loggerDb.Printf("Connected to database")
	env.db = &db{sql, loggerDb}
}
*/
