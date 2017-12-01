package env

import (
	"errors"
)

const (
	createTournament = "INSERT INTO tournaments (selector, name) VALUES ($1, $2) RETURNING id"
	getTournament    = "SELECT id, name FROM tournaments WHERE selector=$1"
	updateTournament = "UPDATE tournaments SET name=$1 WHERE id=$2"
	deleteTournament = "DELETE FROM tournaments WHERE id=$1"

	insertOrganizer = "INSERT INTO organizers (user_id, tournament_id, rank) VALUES ($1, $2, $3)"
	deleteOrganizer = "DELETE FROM organizers WHERE tournament_id=$2 AND user_id=$3" //FIX
	selectOrganizer = "SELECT rank FROM organizers WHERE tournament_id=$2 AND user_id=$3"
	updateOrganizer = "UPDATE organizers SET rank=$1 WHERE tournament_id=$2 AND user_id=$3"

	selectOrganizers    = "SELECT users.selector, users.id, users.first_name, users.last_name, users.gender, users.dob, users.email, organizers.rank FROM users JOIN organizers WHERE organizers.tournament_id=$1"
	deleteAllOrganizers = "DELETE FROM organizers WHERE tournament_id=$1"

	selectAllTeams = "SELECT id, selector, name FROM teams WHERE tournament_id=$1"

	selectTournaments  = "SELECT id, selector, name FROM tournaments ORDER BY id ASC LIMIT $1 OFFSET $2"
	getTournamentCount = "SELECT COUNT(*) FROM tournaments"
)

type tournamentDatastore interface {
	CreateTournament(tour Tournament) (*Tournament, error)
	GetTournament(selector string, full bool) (*Tournament, error)
	UpdateTournament(tour Tournament) error
	DeleteTournament(tour Tournament) error

	GetTournaments(amount, page int) ([]*Tournament, int, error)
}

type Tournament struct {
	Selector
	ID uint

	Name       string
	Owner      *Organizer
	Organizers []*Organizer
	Teams      []*Team
}

type Organizer struct {
	*User
	Rank
}

// CreateTournament creates a new tournament using the struct provided
// and returns a pointer to a new struct
func (d *db) CreateTournament(tour Tournament) (*Tournament, error) {
	tour.sel = d.GenerateSelector(selectorLen)

	err := d.QueryRow(createTournament, tour.sel, tour.Name).Scan(&tour.ID)
	if err != nil {
		d.Logger.Panicln(err)
		return nil, errors.New("failed to create tournament")
	}

	return &tour, nil
}

func (d *db) GetTournament(selector string, full bool) (*Tournament, error) {
	var tour Tournament
	tour.sel = selector

	err := d.QueryRow(getTournament, selector).Scan(&tour.ID, &tour.Name)
	if err != nil {
		d.Logger.Println(err)
		return nil, errors.New("failed to get tournament")
	}

	if full {
		rows, err := d.Query(selectOrganizers, tour.ID)
		if err == nil {
			for rows.Next() {
				var org Organizer
				rows.Scan(&org.sel, &org.ID, &org.FirstName, &org.LastName, &org.Gender, &org.DateOfBirth, &org.Email, &org.Rank)
				tour.Organizers = append(tour.Organizers, &org)
			}
		}

		rows, err = d.Query(selectAllTeams, tour.ID)
		if err == nil {
			for rows.Next() {
				var team Team
				team.TournamentID = tour.ID
				rows.Scan(&team.ID, &team.sel, &team.Name)
				tour.Teams = append(tour.Teams, &team)
			}
		}

	}

	return &tour, nil
}

func (d *db) UpdateTournament(tour Tournament) error {
	tx, err := d.DB.Begin()
	tx.Exec(updateTournament, tour.ID, tour.Name)
	for _, e := range tour.Organizers {
		if e.Rank > 0 { // INSERT if new
			tx.Exec(updateOrganizer, e.Rank, tour.ID, e.ID)
		} else {
			tx.Exec(deleteOrganizer, tour.ID, e.ID)
		}
	}
	err = tx.Commit()
	if err != nil {
		return errors.New("failed to update tournament")
	}

	return nil
}

func (d *db) DeleteTournament(tour Tournament) error {
	tx, _ := d.DB.Begin()
	tx.Exec(deleteTournament, tour.ID)
	tx.Exec(deleteAllOrganizers, tour.ID)
	tx.Commit()

	return nil
}

func (d *db) GetTournaments(amount, page int) ([]*Tournament, int, error) {
	rows, err := d.DB.Query(selectTournaments, amount, amount*page)
	if err != nil {
		return nil, 0, errors.New("could not get next page")
	}
	var rt []*Tournament

	for rows.Next() {
		var tour Tournament
		err = rows.Scan(&tour.ID, &tour.sel, &tour.Name)
		rt = append(rt, &tour)
	}

	count := 0
	err = d.DB.QueryRow(getTournamentCount).Scan(&count)

	return rt, count - (amount * (page + 1)), nil
}
