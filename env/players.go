package env

type Player struct {
	Id int
}

type playersDatastore interface {
	GetPlayer(int) *Player
}

func (db *Db) GetPlayer(id int) *Player {
	db.sql.Ping()
	return &Player{id}
}
