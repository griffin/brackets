package env

import (
	"errors"
)

const (
	createTournament = "INSERT INTO tournaments (id, selector, name) VALUES ($1, $2)"
	getTournament = "SELECT selector, name FROM tournament WHERE selector=$1"
	updateTournament = "UPDATE tournaments SET name=$1 WHERE id=$2"
	deleteTournament = "DELETE FROM tournaments WHERE id=$1"

	insertOrganizer = "INSERT INTO organizers (user_id, tournament_id, rank) VALUES ($1, $2, $3)"
	deleteOrganizer = "DELETE FROM organizers WHERE tournament_id=$2 AND user_id=$3" //FIX
	selectOrganizer = "SELECT rank WHERE tournament_id=$2 AND user_id=$3"
	updateOrganizer = "UPDATE organizers SET rank=$1 WHERE tournament_id=$2 AND user_id=$3"

	selectOrganizers = "SELECT users.selector, users.id, users.first_name, users.last_name, users.gender, users.dob, users.email, organizers.rank FROM users JOIN organizers WHERE organizers.tournament_id=$1"
	deleteAllOrganizers = "DELETE FROM organizers WHERE tournament_id=$1"
)


type tournamentDatastore interface {
	CreateTournament(tour Tournament) (*Tournament, error)
	GetTournament(selector string) (*Tournament, error)
	UpdateTournament(tour Tournament) error
	DeleteTournament(tour Tournament) error
}

type Tournament struct {
	Selectable
	ID uint

	Name     string
	Organizers []*Organizer
	Teams    []*Team
}

type Organizer struct {
	*User
	Rank
}

func (org *Organizer) Delete(){
	org.Rank = Delete // TODO
}

// CreateTournament creates a new tournament using the struct provided
// and returns a pointer to a new struct
func (d *db) CreateTournament(tour Tournament) (*Tournament, error) {
	selector := d.GenerateSelector(selectorLen)
	tx, err := d.DB.Begin()
	res, err := tx.Exec(createTournament, selector, tour.Name);
	for _, e := range tour.Organizers {
		res, err := tx.Exec(insertOrganizer, e.ID, tour.ID, e.Rank)
	}
	err = tx.Commit()
	if err != nil {
		return nil, errors.New("failed to creat tournament")
	}


	return &tour, nil
}

func (d *db) GetTournament(selector string) (*Tournament, error) {

	var tour Tournament

	tx, err := d.DB.Begin()
	tx.QueryRow(getTournament, selector).Scan(tour.ID, tour.selector, tour.Name)
	rows, err := tx.Query(selectOrganizers, tour.ID)
	err = tx.Commit()
	if err != nil {
		return nil, errors.New("failed to get tournament")
	}

	for rows.Next() {
		org := &Organizer{}
		rows.Scan(org.selector, org.ID, org.FirstName, org.LastName, org.Gender, org.DateOfBirth, org.Email, org.Rank)
		tour.Organizers = append(tour.Organizers, org)
	}


	return &tour, nil
}

func (d *db) UpdateTournament(tour Tournament) error {
	tx, err := d.DB.Begin()
	res, err := tx.Exec(updateTournament, tour.ID, tour.Name);
	for _, e := range tour.Organizers {
		if e.Rank > 0 { // INSERT if new
			res, err := tx.Exec(updateOrganizer, e.Rank, tour.ID, e.ID)
		} else {
			res, err := tx.Exec(deleteOrganizer, tour.ID, e.ID)
		}
	}
	err = tx.Commit()
	if err != nil {
		return errors.New("failed to update tournament")
	}

	return nil
}

func (d *db) DeleteTournament(tour Tournament) error {
	tx, err := d.DB.Begin()
	res, err := tx.Exec(deleteTournament, tour.ID)
	res, err = tx.Exec(deleteAllOrganizers, tour.ID)
	err = tx.Commit()

	return nil
}

func (d *db) GetOrganizer(selector string, usr User) (*Organizer, error) {
	//FINISH, used to validate permissions
	return nil, nil
}