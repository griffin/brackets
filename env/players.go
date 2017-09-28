package env

type Player struct {
	userId    int
	Email     string
	FirstName string
	LastName  string
	validator string
}

type forgot_tokens struct {
	id        int
	userId    int
	validator string
	exp       int64
}

type players interface {
	GetPlayer(int) *Player
}

func (db *db) GetPlayer(id int) *Player {
	db.Ping()
	return &Player{id}
}
