package repository

import (
	"log"
	"time"

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

	// Determine the main host (first in the Hosts array)
	host_main := ""
	if len(cfg.Hosts) > 0 && cfg.Hosts[0] != "" {
		host_main = cfg.Hosts[0]
	}

	dsn_main := cfg.User + ":" + cfg.Password + "@tcp(" + host_main + ":" + cfg.Port + ")/" + cfg.Name + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn_main), config)
	if err != nil {
		return nil, err
	}

	// DBResolver adds multiple databases support to GORM
	// https://github.com/go-gorm/dbresolver
	if len(cfg.Hosts) > 1 { // Check if we have additional hosts besides the main one
		replicas := make([]gorm.Dialector, 0)

		// Start from index 1 since index 0 is already used as the primary
		for i := 1; i < len(cfg.Hosts); i++ {
			host := cfg.Hosts[i]
			if host == "" {
				continue
			}

			replicas = append(replicas, mysql.Open(cfg.User+":"+cfg.Password+"@tcp("+host+":"+cfg.Port+")/"+cfg.Name+"?charset=utf8mb4&parseTime=True&loc=Local"))
		}

		if len(replicas) > 0 {
			err = db.Use(dbresolver.Register(dbresolver.Config{
				Sources:  replicas,
				Replicas: replicas,
				Policy:   dbresolver.RandomPolicy{},
			}).
				SetMaxIdleConns(100).
				SetMaxOpenConns(200).
				SetConnMaxIdleTime(time.Hour).
				SetConnMaxLifetime(24 * time.Hour))

			if err != nil {
				return nil, err
			}
		}
	}

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
		&model.Log{},
	)
	if err != nil {
		return err
	}

	log.Println("DB migration OK")

	return nil
}
