package env

import (
	_ "github.com/go-redis/redis"
	"strings"
)

const (
	selectSession = ""
)

type sessionDatastore interface {
}

func (d *db) CheckSession(token string) error {
	split := strings.Split(token, ":")
	selector := split[0]
	var selQ string
	validator := split[1]
	var valQ string

	//TODO redis check

	err := d.QueryRow(selectSession).Scan(&selQ, &valQ)
	if err != nil {

	}

	if selector != selQ || validator != valQ {
		return nil //TODO error
	}

	//TODO redis cache

	return nil

}
