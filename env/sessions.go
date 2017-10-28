package env

import (
	"fmt"
	"time"
	"github.com/go-redis/redis"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"crypto/sha256"
	"bytes"
	"errors"
	"encoding/json"
)

const (
	createSession = "INSERT INTO sessions (validator, selector, user_id, exp) VALUES ($1, $2, $3, $4)"
	selectSession = "SELECT users.id, users.email, users.first_name, users.last_name, users.gender, users.dob, sessions.validator, sessions.exp FROM sessions JOIN users WHERE users.selector=$1 "
	validateUser = "SELECT id, selector, validator, first_name, last_name, gender, dob FROM users WHERE email=$1"
	invalidateSession = "DELETE FROM sessions WHERE selector=$1"
	invalidateAllSession = "DELETE FROM sessions WHERE user_id=$1"
	getAllSessionsForUser = "SELECT selector FROM sessions WHERE user_id=$1"

	selectorLen = 12
	validatorLen = 20
)

type sessionDatastore interface {
	CreateSession(username, password string) (*User, error)
	CheckSession(token string) (*User, error)
	InvalidateSession(usr User) error
	InvalidateAllSessions(usr User) error
}

func (d *db) CreateSession(email, password string) (*User, error) {

	usr := &User{}
	var validator string

	//TODO check login times to prevent spam

	err := d.DB.QueryRow(validateUser, email).Scan(&usr.ID, &usr.selector, &validator, &usr.FirstName, &usr.LastName, &usr.Gender, &usr.DateOfBirth)
	if err != nil {
		return nil, errors.New("couldn't find user")
	}

	err = bcrypt.CompareHashAndPassword([]byte(validator), []byte(password))
	if err != nil {
		//TODO send goroutine to add to spam
		return nil, errors.New("incorrect password")
	}

	session := d.insertSession(*usr)

	return usr, nil
}

func (d *db) insertSession(user User) string {
	validator := d.GenerateSelector(validatorLen)
	selector := d.GenerateSelector(selectorLen)
	exp := time.Now().Add(time.Hour * 2).UnixNano() //TODO

	hashValidator := sha256.Sum256([]byte(validator))

	_, err := d.DB.Exec(createSession, string(hashValidator[:]), selector, user.ID, exp)
	if err != nil {
		d.Panicf("Could not insert session: %v", err)
	}

	return fmt.Sprintf("%v:%v", selector, validator)
}

func (d *db) InvalidateSession(token string) error {
	split := strings.Split(token, ":")
	selector := split[0]

	_, err := d.DB.Exec(invalidateSession, selector)
	if err != nil {
		return errors.New("couldn't invalidate sessions")
	}

	err = d.Del(selector).Err()

	return nil
}

func (d *db) InvalidateAllSessions(usr User) error {
	var v []string

	rows, err := d.DB.Query(getAllSessionsForUser, usr.ID)
	for rows.Next() {
        var validator string
		rows.Scan(&validator);
		v = append(v, validator)
	}

	err = d.Del(v...).Err()

	_, err = d.DB.Exec(invalidateAllSession, usr.ID)
	if err != nil {
		return errors.New("couldn't invalidate all sessions")
	}
	
	return nil
}

func (d *db) CheckSession(token string) (*User, error) {
	split := strings.Split(token, ":")
	selector := split[0]
	validator := sha256.Sum256([]byte(split[1]))
	var valQuery string
	var exp int64

	var usr User

	//Check Redis cache

	jsonUser, err := d.Get(selector).Result()
	if err == redis.Nil {
		goto query
	}

	err = json.Unmarshal([]byte(jsonUser), usr)
	if err != nil {
		d.Println("error parsing redis token")
		goto query

	}

	goto check

	query: // Skip to the SQL query

	err = d.QueryRow(selectSession, selector).Scan(&usr.ID, &usr.Email, &usr.FirstName, &usr.LastName, usr.Gender, usr.DateOfBirth, &valQuery, &exp)
	if err != nil {
		return nil, errors.New("no validator found")
	}

	check:

	if bytes.Equal(validator[:], []byte(valQuery)) { // TODO also check exp
		return nil, errors.New("validator != valQ")
	}

	if(time.Now().UnixNano() < exp){ //TODO come up with when to update cache
		go d.updateCache(usr, token)
	}

	return &usr, nil

}

func (d *db) updateCache(usr User, token string) error {
	split := strings.Split(token, ":")
	selector := split[0]
	var validator string
	var userID int
	var exp int64

	json, err := json.Marshal(usr)
	if err != nil {
		return errors.New("falied to marshal user")
	}

	err = d.Set(validator, string(json), time.Duration(exp - time.Now().UnixNano())).Err()
	if err != nil {
		return errors.New("couldn't update cache")
	}

	return nil
}
