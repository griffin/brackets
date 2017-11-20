package env

import (
	"errors"
)

const (
	createTeam = "INSERT INTO teams (selector, name, tournament_id) VALUES ($1, $2, $3)"
	getTeam    = "SELECT id, tournamentID, name FROM teams WHERE selector=$1"
	updateTeam = "UPDATE teams SET name=$1 WHERE team_id=$2"
	deleteTeam = "DELETE FROM teams WHERE team_id=$1"

	insertPlayer = "INSERT INTO players (user_id, team_id, rank) VALUES ($1, $2, $3)"
	updatePlayer = "UPDATE players SET rank=$1 WHERE team_id=$2 AND user_id=$3"
	deletePlayer = "DELETE FROM players WHERE team_id=$2 AND user_id=$3" //FIX
	selectPlayer = "SELECT rank WHERE team_id=$2 AND user_id=$3"

	selectPlayers    = "SELECT users.selector, users.id, users.first_name, users.last_name, users.gender, users.dob, users.email, players.rank FROM users JOIN players WHERE players.team_id=$1"
	deleteAllPlayers = "DELETE FROM players WHERE team_id=$1"
)

type Rank int

const (
	Delete    Rank = -1
	Owner     Rank = 0
	Moderator Rank = 1
	Member    Rank = 3
)

type teamDatastore interface {
	CreateTeam(team Team) (*Team, error)
	GetTeam(selector string) (*Team, error)
	UpdateTeam(team Team) error
	DeleteTeam(selector Team) error
}

type Team struct {
	Selectable

	ID           uint
	tournamentID uint

	Name    string
	Players []*Player
}

type Player struct {
	*User
	Rank
}

func (d *db) CreateTeam(team Team) (*Team, error) {
	selector := d.GenerateSelector(selectorLen)
	tx, err := d.DB.Begin()
	tx.Exec(createTeam, selector, team.Name, team.tournamentID)
	for _, e := range team.Players {
		tx.Exec(insertPlayer, e.ID, team.ID, e.Rank)
	}
	err = tx.Commit()
	if err != nil {
		return nil, errors.New("failed to create team")
	}

	return &team, nil
}

func (d *db) GetTeam(selector string) (*Team, error) {
	var team Team
	team.sel = selector

	tx, err := d.DB.Begin()
	if err != nil {
		return nil, errors.New("Couldn't get team")
	}

	tx.QueryRow(getTeam, team.Selector()).Scan(team.ID, team.tournamentID, team.Name)
	rows, err := tx.Query(selectPlayers, team.ID)

	err = tx.Commit()
	if err != nil {
		return nil, errors.New("failed to get team")
	}

	for rows.Next() {
		pl := &Player{}
		rows.Scan(pl.sel, pl.ID, pl.FirstName, pl.LastName, pl.Gender, pl.DateOfBirth, pl.Email, pl.Rank)
		team.Players = append(team.Players, pl)
	}

	return &team, nil
}

func (d *db) UpdateTeam(team Team) error {
	tx, err := d.DB.Begin()
	tx.Exec(updateTeam, team.Name, team.Selector())
	for _, e := range team.Players {
		if e.Rank > 0 { // INSERT if new
			tx.Exec(updatePlayer, e.Rank, team.ID, e.ID)
		} else {
			tx.Exec(deletePlayer, team.ID, e.ID)
		}
	}
	err = tx.Commit()
	if err != nil {
		return errors.New("failed to update team")
	}

	return nil
}

func (d *db) DeleteTeam(team Team) error {
	tx, err := d.DB.Begin()
	if err != nil {
		return errors.New("Couldn't delete team")
	}
	tx.Exec(deleteTeam, team.ID)
	tx.Exec(deleteAllPlayers, team.ID)
	tx.Commit()

	return nil
}
