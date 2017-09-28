package env

import (
	"time"
)

type tournamentDatastore interface {
	CreateTournament(tournament Tournament) (*Tournament, error)
	GetTournament(selector string) (*Tournament, error)
	UpdateTournament(tournament Tournament) error
	DeleteTournament(selector string) error
}

type Tournament struct {
	Selectable
	id uint

	Name     string
	Managers map[string]*Manager
	Teams    map[string]*Team
	Games    map[int]*Game
}

type Manager struct {
	*User
	Rank
}

type gameDatastore interface {
}

type Game struct {
	Selectable

	BracketPosition uint
	HomeTeam        *Team
	AwayTeam        *Team
	Time            time.Time
}
