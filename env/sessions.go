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
	"time"
	"encoding/json"
)

const (
	createSession = "INSERT INTO sessions (validator, selector, user_id, exp) VALUES ($1, $2, $3, $4)"
	selectSession = "SELECT users.email, users.first_name, users.last_name, users.gender, users.dob, sessions.validator, sessions.exp FROM sessions JOIN users WHERE users.selector=$1 "
	validateUser = "SELECT id, selector, validator, first_name, last_name, gender, dob FROM users WHERE email=$1"
	invalidateSession = "DELETE FROM sessions WHERE selector=$1"
	invalidateAllSession = "DELETE FROM sessions WHERE user_id=$1"


	selectorLen = 12
	validatorLen = 20
)

type sessionDatastore interface {
	CreateSession(username, password string) (*User, error)
	CheckSession(token string) (*User, error)
	InvalidateSession(user User) error
	InvalidateAllSessions(user User) error
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
	exp := time.Now().Add(time.Hour * 2) //TODO

	hashValidator := sha256.Sum256([]byte(validator))

	_, err := d.DB.Exec(createSession, string(hashValidator[:]), selector, user.ID, exp)
	if err != nil {
		d.Panicf("Could not insert session: %v", err)
	}

	return fmt.Sprintf("%v:%v", selector, validator)
}

func (d *db) InvalidateSession(token string) {

}

func (d *db) InvalidateAllSessions(userID int) {

}

func (d *db) CheckSession(token string) (*User, error) {
	
	split := strings.Split(token, ":")
	selector := split[0]
	validator := sha256.Sum256([]byte(split[1]))

	var usr *User

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

	err = d.QueryRow(selectSession, selector).Scan(usr.valQuery, &userID, &exp)
	if err != nil {
		d.Println("no validator found")
		return nil, nil
	}

	check:

	if bytes.Equal(validator[:], []byte(valQuery)) { // TODO also check exp
		d.Println("validator != valQ")
		return nil, nil
	}

	if(exp > 0){ //TODO come up with when to update cache
		go d.updateCache(string(validator[:]))
	}

	return 0, usr

}

func (d *db) updateCache(selector string){

	var validator string
	var userID int
	var exp int64

	err := d.QueryRow(selectSession, selector).Scan(&validator, &userID, &exp)

}
