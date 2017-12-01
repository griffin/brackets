package env

import (
	"errors"
	"time"
	"bytes"
	"strconv"
)

const (
	getGame    = "SELECT (id, away_id, home_id, time) FROM games WHERE selector=$1"
	getGames   = "SELECT (selector, id, away_id, home_id, time) FROM games WHERE away_id=$1 OR home_id=$1"
	createGame = "INSERT INTO games (selector, away_id, home_id, time, location) VALUES ($1, $2, $3, $4, $5)"
	deleteGame = "DELETE FROM games WHERE id=$1"
	updateGame = "UPDATE games SET away_id=$1, home_id=$2, time=$3 WHERE id=$4"


	getAllGames = "SELECT games.id, games.selector, games.time, games.location, t1.id, t1.name, " +
					"t1.selector, t2.id, t2.name, t2.selector " +
					"FROM games JOIN teams t1 ON games.away_id=t1.id JOIN teams t2 ON games.home_id=t2.id WHERE t1.id IN ($1) OR t2.id IN ($1) " +
					"ORDER BY games.time ASC"

)

type gameDatastore interface {
	CreateGame(game Game) (*Game, error)
	GetGame(selector string) (*Game, error)
	GetUpcomingGames(team []Team) ([]*Game, error)
	UpdateGame(team Team) error
	DeleteGame(game Game) error
}

type Game struct {
	Selector

	ID uint

	HomeTeam *Team
	AwayTeam *Team
	Location string
	Time     time.Time
}

func (d *db) CreateGame(game Game) (*Game, error) {
	game.sel = d.GenerateSelector(selectorLen)
	_, err := d.DB.Exec(createGame, game.sel, game.AwayTeam.ID, game.HomeTeam.ID, game.Time, game.Location)
	if err != nil {
		d.Logger.Println(err)
		return nil, errors.New("couldn't create game")
	}

	return &game, nil
}

func (d *db) GetGame(selector string) (*Game, error) {
	return nil, nil
}

func (d *db) GetUpcomingGames(team []Team) ([]*Game, error) {

	var buf bytes.Buffer
	
	for i, t := range team {
		buf.WriteString(strconv.Itoa(int(t.ID)))
		
		if i+1 != len(team) {
			buf.WriteString(",")
		}
	}

	rows, err := d.Query(getAllGames, buf.String()[:len(buf.String())])
	if err != nil {
		return nil, err
	}

	var rt []*Game
	for rows.Next() {
		var away Team
		var home Team
		var game Game

		rows.Scan(&game.ID, &game.sel, &game.Time, &game.Location, &away.ID, &away.Name, &away.sel, &home.ID, &home.Name, &home.sel)

		game.HomeTeam = &home
		game.AwayTeam = &away
		rt = append(rt, &game)

	}

	return rt, nil
}

func (d *db) UpdateGame(team Team) error {
	return nil
}

func (d *db) DeleteGame(game Game) error {
	return nil
}
