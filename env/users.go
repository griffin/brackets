package env

import (
	"time"
	"golang.org/x/crypto/bcrypt"
)

const (
	createUser = "INSERT INTO users (selector, validator, first_name, last_name, gender, dob, email) VALUES (?, ?, ?, ?, ?, ?, ?)"
	getUser = "SELECT id, first_name, last_name, gender, dob, email FROM users"
	updateUser = "UPDATE first_name, last_name, gender, dob, email IN users WHERE id=?"
	deleteUser = "DELETE FROM users WHERE selector=?"
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

	ID    uint

	Email     string
	FirstName string
	LastName  string

	Gender      Gender
	DateOfBirth time.Time
}

type userDatastore interface {
	CreateUser(user User, password string) (*User, error)
	GetUser(selector string) (*User, error)
	UpdateUser(user User) error
	DeleteUser(selector string) error
}

func (d *db) CreateUser(user User, password string) (*User, error) {
	validator, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	selector := d.GenerateSelector(selectorLen)
	if err != nil {
		return nil, err
	}

	res, err := d.DB.Exec(createUser, selector, string(validator), user.FirstName, user.LastName, user.Gender, user.DateOfBirth, user.Email)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()

	user.selector = selector
	user.ID = uint(id)

	return &user, nil
}

func (d *db) GetUser(selector string) (*User, error) {

	return &User{}, nil
}

func (d *db) UpdateUser(user User) error {
	
	return nil
}

func (d *db) DeleteUser(selector string) error {
		
	return nil
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
