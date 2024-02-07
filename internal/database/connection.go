package database

import (
	"log"

	"github.com/jinzhu/gorm"
)

func Connect() (*gorm.DB, error) {
	dsn := ""
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	log.Println("DB connection successful")

	return db, nil
}
