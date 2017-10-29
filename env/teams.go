package env


const (
	createTeam = "INSERT INTO teams (selector, name, tournament_id) VALUES ($1, $2, $3)"
	getTeam = "SELECT id, tournamentID, name FROM teams WHERE selector"
	updateTeam = ""
	deleteTeam = ""
)

type Rank int

const (
	Delete Rank = -1
	Owner Rank = 0
	Moderator Rank = 1
	Member Rank = 3
)

type teamDatastore interface {
	CreateTeam(team Team) (*Team, error)
	GetTeam(selector string) (*Team, error)
	UpdateTeam(team Team) error
	DeleteTeam(selector string) error
}

func (d *db) CreateTeam(team Team) (*Team, error){

	


	return nil, nil
}

func (d *db) GetTeam(selector string) (*Team, error){
	return nil, nil
}

func (d *db) UpdateTeam(team Team) error{
	return nil
}

func (d *db) DeleteTeam(selector string) error {
	return nil
}



type Team struct {
	Selectable

	ID           uint
	tournamentID uint

	Name    string
	Players []*Player
}

type Manager struct {
	*User
	Rank
}

type Player struct {
	*User
	Rank
}
