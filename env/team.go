package env

import (
	"time"
)

type Team struct {
	teamId  uint
	Name    string
	Players map[string]string
}

type joined struct {
	userId uint
	teamId uint
	rank   string
}

type Post struct {
	Id     uint
	Title  string
	Author *Player
	Path   string
	Posted time.Time
}
