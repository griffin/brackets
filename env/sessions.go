package env

import (
	"github.com/go-redis/redis"
)

type sessionDatastore interface {
}

type session_tokens struct {
	Selectable

	id        uint
	validator string
	userID    uint
	exp       int64
}
