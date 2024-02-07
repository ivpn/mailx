package database

import (
	"log"

	"github.com/jinzhu/gorm"
)

type Database struct {
	Client *gorm.DB
}

func New() (*Database, error) {
	db, err := Connect()
	if err != nil {
		log.Printf("an error occured connecting DB: %s", err.Error())
		return nil, err
	}

	err = Migrate(db)
	if err != nil {
		log.Printf("an error occured migrating DB: %s", err.Error())
		return nil, err
	}

	return &Database{
		Client: db,
	}, nil
}
