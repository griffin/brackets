package env

import (
	"encoding/csv"
	"fmt"
	"io"
	//	"math/rand"
	"os"
	"time"
)

func (e *Env) GenUsers() {

	file, _ := os.Open("test.csv")
	r := csv.NewReader(file)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		fmt.Println(record)

		usr := User{
			Email:       record[2],
			FirstName:   record[0],
			LastName:    record[1],
			Gender:      0,
			DateOfBirth: time.Now(),
		}

		_, err = e.Db.CreateUser(usr, record[3])
		fmt.Println(err)
	}
}

func (e *Env) GenTour() {

	file, _ := os.Open("tour.csv")
	r := csv.NewReader(file)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		fmt.Println(record)

		tour := Tournament{
			Name: record[0],
		}

		_, err = e.Db.CreateTournament(tour)
		fmt.Println(err)
	}
}
