package env

const (
	createTournament = "INSERT INTO tournaments (selector, name) VALUES ()"
	getTournament = ""
	updateTournament = ""
	deleteTournament = ""
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
	Organizers []*Organizer
	Teams    []*Team
	Games    []*Game
}

type Organizer struct {
	*User
	Rank
}

// CreateTournament creates a new tournament using the struct provided
// and returns a pointer to a new struct
func (d *db) CreateTournament(tournament Tournament) (*Tournament, error) {
	return nil, nil
}

func (d *db) GetTournament(selector string) (*Tournament, error) {
	return nil, nil
}

func (d *db) UpdateTournament(tournament Tournament) error {
	return nil
}

func (d *db) DeleteTournament(selector string) error {
	return nil
}

func (d *db) GetOrganizer(selector string) (*Organizer, error) {
	return nil, nil
}

func (d *db) UpdateOrganizer(selector string) (*Organizer, error) {
	return nil, nil
}

func (d *db) DeleteOrganizer(selector string) error {
	return nil
}