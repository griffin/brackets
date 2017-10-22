package env

import (
	"time"
)

type gameDatastore interface {
}

type Game struct {
	Selectable

	BracketPosition uint
	HomeTeam        *Team
	AwayTeam        *Team
	Time            time.Time
}