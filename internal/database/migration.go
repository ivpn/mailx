package database

import (
	"log"

	"github.com/jinzhu/gorm"
	"ivpn.net/email-service/internal/service"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&service.Recipient{}).Error
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&service.Alias{}).Error
	if err != nil {
		return err
	}

	log.Println("DB migration successful")

	return nil
}
