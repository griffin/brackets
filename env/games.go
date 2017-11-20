package env

import (
	"errors"
	"time"
)

const (
	getGame    = "SELECT (id, away_id, home_id, time) FROM games WHERE selector=$1"
	getGames   = "SELECT (selector, id, away_id, home_id, time) FROM games WHERE away_id=$1 OR home_id=$1"
	createGame = "INSERT INTO games (selector, away_id, home_id, time) VALUES ($1, $2, $3, $4)"
	deleteGame = "DELETE FROM games WHERE id=$1"
	updateGame = "UPDATE games SET away_id=$1, home_id=$2, time=$3 WHERE id=$4"
)

type gameDatastore interface {
	CreateGame(game Game) (*Game, error)
	GetGame(selector string) (*Game, error)
	GetGames(team Team) ([]*Game, error)
	UpdateGame(team Team) error
	DeleteGame(game Game) error
}

type Game struct {
	Selectable

	ID uint

	HomeTeam *Team
	AwayTeam *Team
	Time     time.Time
}

func (d *db) CreateGame(game Game) (*Game, error) {
	game.sel = d.GenerateSelector(selectorLen)
	res, err := d.DB.Exec(createGame, game.sel, game.AwayTeam.ID, game.HomeTeam.ID, game.Time)
	if err != nil {
		return nil, errors.New("couldn't create game")
	}

	id, err := res.LastInsertId()
	game.ID = uint(id)

	return &game, nil
}

func (d *db) GetGame(selector string) (*Game, error) {
	return nil, nil
}

func (d *db) GetGames(team Team) ([]*Game, error) {
	return nil, nil
}

func (d *db) UpdateGame(team Team) error {
	return nil
}

func (d *db) DeleteGame(game Game) error {
	return nil
}
