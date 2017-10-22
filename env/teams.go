package env


const (
	createTeam = "INSERT INTO teams () VALUES ()"
	getTeam = "SELECT id, tournamentID, name FROM"
	updateTeam = ""
	deleteTeam = ""
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

	id           uint
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
