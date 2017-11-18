package env

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"
)

func GenUsers() {
	e := New()

	sql := SQLOptions{
		"postgres",
		"mysecretpassword",
		"localhost",
		5432,
		"brackets",
	}

	redis := RedisOptions{}

	e.ConnectDb(sql, redis)

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
