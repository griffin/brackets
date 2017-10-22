package env

import (
	"time"
	"github.com/go-redis/redis"
	"strings"
	"crypto/sha256"
	"bytes"
	"strconv"
)

const (
	createSession = "INSERT INTO sessions (validator, selector, user_id, exp) VALUES (?, ?, ?, ?)"
	selectSession = "SELECT validator, user_id, exp FROM sessions WHERE selector=?"
	invalidateSession = "DELETE FROM sessions WHERE selector=?"
	invalidateAllSession = "DELETE FROM sessions WHERE user_id=?"
)

type sessionDatastore interface {
	CreateSession(userID int)
	CheckSession(token string) (int, error)
	InvalidateSession(token string)
	InvalidateAllSessions(userID int)
}

func (d *db) CreateSession(userID int) {

}

func (d *db) InvalidateSession(token string) {

}

func (d *db) InvalidateAllSessions(userID int) {

}

func (d *db) CheckSession(token string) (int, error) {
	
	split := strings.Split(token, ":")
	selector := split[0]
	validator := sha256.Sum256([]byte(split[1]))

	var valQuery string
	var exp int64
	var userID int

	//Check Redis cache

	token, err := d.Get(selector).Result()
	if err == redis.Nil {
		goto query
	}

	split = strings.Split(token, ":")
	userID, err = strconv.Atoi(split[0])
	valQuery = split[1]

	if err != nil {
		d.Println("error parsing redis token")
		goto query

	}

	goto check

	query: // Skip to the SQL query

	err = d.QueryRow(selectSession, selector).Scan(&valQuery, &userID, &exp)
	if err != nil {
		d.Println("no validator found")
		return -1, nil
	}

	check:

	if bytes.Equal(validator[:], []byte(valQuery)[:]) { // TODO also check exp
		d.Println("validator != valQ")
		return -1, nil
	}

	if(exp > 0){ //TODO come up with when to update cache
		go d.updateCache(string(validator[:]))
	}

	return 0, nil

}

func (d *db) updateCache(selector string){

	var validator string
	var userID int
	var exp int64

	err := d.QueryRow(selectSession, selector).Scan(&validator, &userID, &exp)


	err = d.Set(selector, validator + ":" + strconv.Itoa(userID), time.Hour * 3).Err()
	if err != nil {

	}
}
