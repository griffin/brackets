package env

import (
	"database/sql"
	"fmt"
	"github.com/go-redis/redis"
	"html/template"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq" // Postgres driver
)

// Env represents the environment that is needed to
// respond to http queries
type Env struct {
	Db       datastore
	Template *template.Template
	Log      *log.Logger
}

type db struct {
	*sql.DB
	*log.Logger
	*redis.Client

	randLk *sync.Mutex
	rand   *rand.Rand
}

type datastore interface {
	userDatastore
	sessionDatastore
	teamDatastore
	//postDatastore
	tournamentDatastore
	gameDatastore
}

// SQLOptions lets you define the parameters to connect to
// a sql database
type SQLOptions struct {
	User     string
	Password string
	Host     string
	Port     int
	Database string
}

// RedisOptions lets you define the parameters to connect to
// a redis database
type RedisOptions struct {
	Password string
	Host     string
	Port     int
}

// String returns a string representation of a sql options
func (s SQLOptions) String() string {
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", s.User, s.Password, s.Host, s.Port, s.Database)
}

// Selectable allows a struct to be selectable in a database
type Selectable struct {
	selector string
}

// Selector returns the structs selector
func (s Selectable) Selector() string {
	return s.selector
}

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// GenerateSelector generates a selector of length n using the env's random
// Safe for concurrent use.
func (d *db) GenerateSelector(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, d.randInt63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = d.randInt63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// randInt63 gets a random int
// Safe for concurrent use.
func (d *db) randInt63() (n int64) {
	d.randLk.Lock()
	n = d.rand.Int63()
	d.randLk.Unlock()
	return
}

// New enviornment which allows for database calls
func New() *Env {
	logger := log.New(os.Stdout, "log: ", log.Lshortfile)
	return &Env{nil, template.New("web"), logger}
}

// ConnectDb connects the env to the sql database with the sqlOpt and the redis
// database with redisOpt
func (env *Env) ConnectDb(sqlOpt SQLOptions, redisOpt RedisOptions) {
	loggerDb := log.New(os.Stdout, "db: ", log.Lshortfile)

	lk := &sync.Mutex{}
	ra := rand.New(rand.NewSource(time.Now().UnixNano()))

	d, err := sql.Open("postgres", sqlOpt.String())
	if err != nil {
		loggerDb.Fatal(err)
	}

	redisConv := &redis.Options{
		Addr:     fmt.Sprintf("%v:%v", redisOpt.Host, redisOpt.Port),
		Password: "",
		DB:       0,
	}

	r := redis.NewClient(redisConv)

	_, err = r.Ping().Result()
	if err != nil {
		loggerDb.Fatal(err)
	}

	loggerDb.Printf("Connected to database")
	env.Db = &db{d, loggerDb, nil, lk, ra}
}
