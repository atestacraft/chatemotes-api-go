package database

import (
	"log"

	db "github.com/sonyarouje/simdb"
)

func New() *db.Driver {
	driver, err := db.New("db")
	if err != nil {
		log.Fatal(err)
	}

	return driver
}
