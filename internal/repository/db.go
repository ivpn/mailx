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

func New() (*Database, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	db, err := connect(cfg.DB)
	if err != nil {
		return nil, err
	}

	err = migrate(db)
	if err != nil {
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

	log.Println("DB connection OK")

	return db, nil
}

func migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&model.Recipient{}, &model.Alias{}, &model.User{})
	if err != nil {
		return err
	}

	log.Println("DB migration OK")

	return nil
}
