package database

import (
	"log"

	"github.com/jinzhu/gorm"
)

type Database struct {
	Client *gorm.DB
}

func New() (*Database, error) {
	db, err := connectDB()
	if err != nil {
		log.Printf("an error occured initializing DB: %s", err.Error())
		return nil, err
	}

	return &Database{
		Client: db,
	}, nil
}

func connectDB() (*gorm.DB, error) {
	dsn := ""
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
