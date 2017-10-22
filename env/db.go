package env

import (
	"github.com/go-redis/redis"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // Postgres driver
)

type db struct {
	*sql.DB
	*log.Logger
	*redis.Client
}

type datastore interface {
	userDatastore
	sessionDatastore
	teamDatastore
	postDatastore
	tournamentDatastore
	gameDatastore
}

type SqlOptions struct {
	User string
	Password string
	Host string
	Port int
	Database string
}

type RedisOptions struct {
	Password string
	Host string
	Port int
}

func (s SqlOptions) String() string {
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v", s.User, s.Password, s.Host, s.Port, s.Database)
}

type Selectable struct {
	selector string
}

func (s Selectable) Selector() string {
	return s.selector
}

type Env struct {
	Db  datastore
	Log *log.Logger
}

func New() *Env {
	logger := log.New(os.Stdout, "log: ", log.Lshortfile)
	return &Env{nil, logger}
}

func (env *Env) ConnectDb(sqlOpt SqlOptions, redisOptions RedisOptions) {
	loggerDb := log.New(os.Stdout, "db: ", log.Lshortfile)

	d, err := sql.Open("postgres", sqlOpt.String())
	if err != nil {
		loggerDb.Fatal(err)
	}

	redisConv := &redis.Options {
		Addr: fmt.Sprintf("%v:%v", redisOptions.Host, redisOptions.Port),
		Password: redisOptions.Password,
		DB: 0,
	}

	r := redis.NewClient(redisConv)

	_, err = r.Ping().Result()
	if err != nil {
		loggerDb.Fatal(err)
	}

	loggerDb.Printf("Connected to database")
	env.Db = &db{d, loggerDb, r}
}
