package env

import (
	"time"
)

type Gender int

const (
	Male Gender = iota
	Female
	Other
	PreferNotToSay
)

type User struct {
	Selectable

	userID    uint
	validator string

	Email     string
	FirstName string
	LastName  string

	Gender      Gender
	DateOfBirth time.Time
}

type userDatastore interface {
	CreateUser(user User) (*User, error)
	GetUser(selector string) (*User, error)
	UpdateUser(user User) error
	DeleteUser(selector string) error
}

func (db *db) GetUser(email string) *User {
	db.Ping()
	return &User{}
}

func (g Gender) String() string {
	switch g {
	case 0:
		return "Male"
	case 1:
		return "Female"
	case 2:
		return "Other"
	}
	return "Prefer not to say"
}
