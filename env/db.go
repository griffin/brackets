package env

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // Postgres driver
)

type db struct {
	*sql.DB
	*log.Logger
}

type Env struct {
	Db  datastore
	Log *log.Logger
}

type datastore interface {
	userDatastore
	sessionDatastore
	teamDatastore
	tournamentDatastore
	gameDatastore
}

type Selectable struct {
	selector string
}

func (s Selectable) Selector() string {
	return s.selector
}

func New() *Env {
	logger := log.New(os.Stdout, "log: ", log.Lshortfile)
	return &Env{nil, logger}
}

func DbString(database, host, username, password string) string {
	return fmt.Sprintf("postgres://%v:%v@%v/%v", username, password, host, database)
}

func (env *Env) ConnectDb(dbString string) {
	loggerDb := log.New(os.Stdout, "db: ", log.Lshortfile)
	sql, err := sql.Open("postgres", dbString)
	if err != nil {
		loggerDb.Fatal(err)
	}

	loggerDb.Printf("Connected to database")
	env.Db = &db{sql, loggerDb}
}
