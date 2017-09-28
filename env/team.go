package env

import (
	"time"
)

type Rank int

const (
	Owner Rank = iota
	Moderator
	Member
)

type teamDatastore interface {
	CreateTeam(team Team) (*Team, error)
	GetTeam(selector string) (*Team, error)
	UpdateTeam(team Team) error
	DeleteTeam(selector string) error
}

type Team struct {
	Selectable

	id           uint
	tournamentID uint

	Name    string
	Players map[string]*Player
}

type Player struct {
	*User
	Rank
}

type Post struct {
	Selectable

	id     uint
	Title  string
	Author *User
	Path   string
	Posted time.Time
}
