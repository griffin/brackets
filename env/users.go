package env

import (
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	createUser = "INSERT INTO users (selector, validator, first_name, last_name, gender, dob, email) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	getUser    = "SELECT id, first_name, last_name, gender, dob, email FROM users WHERE selector=$1"
	updateUser = "UPDATE users SET first_name=$1, last_name=$2, gender=$3, dob=$4, email=$5 WHERE id=$6"
	deleteUser = "DELETE FROM users WHERE id=$1"

	selectUsers  = "SELECT id, selector, first_name, last_name, gender, dob FROM users ORDER BY id ASC LIMIT $1 OFFSET $2"
	getUserCount = "SELECT COUNT(*) FROM users"
)

type Gender int8

const (
	Male           Gender = 0
	Female         Gender = 1
	Other          Gender = 2
	PreferNotToSay Gender = 3
)

type User struct {
	Selector

	ID uint

	Email     string
	FirstName string
	LastName  string

	Gender      Gender
	DateOfBirth time.Time
}

type userDatastore interface {
	CreateUser(usr User, password string) (*User, error)
	GetUser(selector string) (*User, error)
	UpdateUser(usr User) error
	DeleteUser(usr User) error

	GetUsers(amount, page int) ([]*User, int, error)
}

func (d *db) CreateUser(usr User, password string) (*User, error) {
	validator, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	usr.sel = d.GenerateSelector(selectorLen)
	if err != nil {
		return nil, err
	}

	res, err := d.DB.Exec(createUser, usr.Selector.String(), string(validator), usr.FirstName, usr.LastName, usr.Gender, usr.DateOfBirth, usr.Email)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()

	usr.ID = uint(id)

	return &usr, nil
}

func (d *db) GetUser(selector string) (*User, error) {
	var usr User

	err := d.DB.QueryRow(getUser, selector).Scan(&usr.ID, &usr.FirstName, &usr.LastName, &usr.Gender, &usr.DateOfBirth, &usr.Email)
	if err != nil {
		return nil, errors.New("Couldn't find user")
	}

	return &usr, nil
}

func (d *db) UpdateUser(usr User) error {

	_, err := d.DB.Exec(updateUser, usr.FirstName, usr.LastName, usr.Gender, usr.DateOfBirth, usr.Email, usr.ID)
	if err != nil {
		return errors.New("update user failed")
	}

	return nil
}

func (d *db) DeleteUser(usr User) error {

	_, err := d.DB.Exec(deleteUser, usr.ID)
	if err != nil {
		return errors.New("delete user failed")
	}

	return nil
}

func (d *db) GetUsers(amount, page int) ([]*User, int, error) {
	rows, err := d.DB.Query(selectUsers, amount, amount*page)
	if err != nil {
		return nil, 0, errors.New("could not get next page")
	}
	var rt []*User

	for rows.Next() {
		var usr User
		err = rows.Scan(&usr.ID, &usr.sel, &usr.FirstName, &usr.LastName, &usr.Gender, &usr.DateOfBirth)
		rt = append(rt, &usr)
	}

	count := 0
	err = d.DB.QueryRow(getUserCount).Scan(&count)

	return rt, count - (amount * (page + 1)), nil
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

func ToGender(st string) Gender {
	switch strings.ToLower(st) {
	case "male":
	case "m":
		return Male
	case "female":
	case "f":
		return Female
	case "other":
	case "o":
		return Other
	}
	return PreferNotToSay
}

func Age(birthday time.Time) int {
	now := time.Now()
	years := now.Year() - birthday.Year()
	if now.YearDay() < birthday.YearDay() {
		years--
	}
	return years
}

func HttpString(date time.Time) string {
	return date.Format("2006-01-02")
}
