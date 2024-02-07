package database

import (
	"log"

	"github.com/jinzhu/gorm"
	"ivpn.net/email-service/config"
	"ivpn.net/email-service/internal/model"
)

type Database struct {
	Client *gorm.DB
}

func New(cfg config.DBConfig) (*Database, error) {
	db, err := connect(cfg)
	if err != nil {
		log.Printf("an error occured connecting DB: %s", err.Error())
		return nil, err
	}

	err = migrate(db)
	if err != nil {
		log.Printf("an error occured migrating DB: %s", err.Error())
		return nil, err
	}

	return &Database{
		Client: db,
	}, nil
}

func connect(cfg config.DBConfig) (*gorm.DB, error) {
	dsn := cfg.User + ":" + cfg.Password + "@tcp(" + cfg.Host + ":" + cfg.Port + ")/" + cfg.Name + "?charset=utf8&parseTime=True&loc=Europe%2FBerlin"

	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	log.Println("DB connection successful")

	return db, nil
}

func migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&model.Recipient{},
		&model.Alias{},
	).Error
	if err != nil {
		return err
	}

	log.Println("DB migration successful")

	return nil
}
