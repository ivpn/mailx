package database

import (
	"log"

	"github.com/jinzhu/gorm"
	"ivpn.net/email-service/internal/model"
)

func Migrate(db *gorm.DB) error {
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
