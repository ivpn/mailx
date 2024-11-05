package repository

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"ivpn.net/email/api/config"
	"ivpn.net/email/api/internal/model"
)

type Database struct {
	Client *gorm.DB
}

func NewDB(cfg config.DBConfig) (*Database, error) {
	db, err := connect(cfg)
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

func (d *Database) Close() error {
	db, err := d.Client.DB()
	if err != nil {
		return err
	}

	return db.Close()
}

func connect(cfg config.DBConfig) (*gorm.DB, error) {
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	dsn_a := cfg.User + ":" + cfg.Password + "@tcp(:" + cfg.Port + ")/" + cfg.Name + "?charset=utf8mb4&parseTime=True&loc=Local"
	dsn_b := cfg.User + ":" + cfg.Password + "@tcp(:" + cfg.Port + ")/" + cfg.Name + "?charset=utf8mb4&parseTime=True&loc=Local"
	dsn_c := cfg.User + ":" + cfg.Password + "@tcp(:" + cfg.Port + ")/" + cfg.Name + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn_a), config)
	if err != nil {
		return nil, err
	}

	// DBResolver adds multiple databases support to GORM
	// https://github.com/go-gorm/dbresolver
	db.Use(dbresolver.Register(dbresolver.Config{
		Sources: []gorm.Dialector{mysql.Open(dsn_b), mysql.Open(dsn_c)},
		Policy:  dbresolver.RandomPolicy{},
	}))

	log.Println("DB connection OK")

	return db, nil
}

func migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&model.User{},
		&model.Subscription{},
		&model.Recipient{},
		&model.Alias{},
		&model.Message{},
		&model.Settings{},
		&model.Session{},
		&model.Credential{},
	)
	if err != nil {
		return err
	}

	log.Println("DB migration OK")

	return nil
}
