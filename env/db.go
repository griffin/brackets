package env

import (
	"sync"
	"github.com/go-redis/redis"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
	"math/rand"

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

// SQLOptions lets you define the parameters to connect to
// a sql database
type SQLOptions struct {
	User string
	Password string
	Host string
	Port int
	Database string
}

// RedisOptions lets you define the parameters to connect to
// a redis database
type RedisOptions struct {
	Password string
	Host string
	Port int
}

func (s SQLOptions) String() string {
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v", s.User, s.Password, s.Host, s.Port, s.Database)
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
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    letterIdxBits = 6                    // 6 bits to represent a letter index
    letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
    letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// GenerateSelector generates a selector of length n using the env's random
// Safe for concurrent use.
func (env *Env) GenerateSelector(n int) string {
    b := make([]byte, n)
    // A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
    for i, cache, remain := n-1, env.randInt63(), letterIdxMax; i >= 0; {
        if remain == 0 {
            cache, remain = env.randInt63(), letterIdxMax
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

// Env represents the environment that is needed to
// respond to http queries
type Env struct {
	Db  datastore
	Log *log.Logger

	randLk  *sync.Mutex
	rand *rand.Rand
}

func (env *Env) randInt63() (n int64) {
	env.randLk.Lock()
	n = rand.Int63()
	env.randLk.Unlock()
	return
}

// New enviornment which allows for database calls
func New() *Env {
	logger := log.New(os.Stdout, "log: ", log.Lshortfile)
	
	lk := &sync.Mutex{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &Env{nil, logger, lk, r}
}

// ConnectDb connects the env to the sql database with the sqlOpt and the redis
// database with redisOpt
func (env *Env) ConnectDb(sqlOpt SQLOptions, redisOpt RedisOptions) {
	loggerDb := log.New(os.Stdout, "db: ", log.Lshortfile)

	d, err := sql.Open("postgres", sqlOpt.String())
	if err != nil {
		loggerDb.Fatal(err)
	}

	redisConv := &redis.Options {
		Addr: fmt.Sprintf("%v:%v", redisOpt.Host, redisOpt.Port),
		Password: redisOpt.Password,
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
