package repository

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	dsn := cfg.User + ":" + cfg.Password + "@tcp(" + cfg.Host + ":" + cfg.Port + ")/" + cfg.Name + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	log.Println("DB connection successful")

	return db, nil
}

func migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&model.Recipient{}, &model.Alias{})
	if err != nil {
		return err
	}

	log.Println("DB migration successful")

	return nil
}
