package database

import (
	"log"

	"github.com/jinzhu/gorm"
)

type Database struct {
	Client *gorm.DB
}

func New() (*Database, error) {
	db, err := connect()
	if err != nil {
		log.Printf("an error occured initializing DB: %s", err.Error())
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

func connect() (*gorm.DB, error) {
	dsn := ""
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	log.Println("DB connection successful")

	return db, nil
}
