package repository

import (
	"context"
	"log"
	"os"

	"ivpn.net/email/api/internal/model"
)

// Configuration
const (
	basePath = "/var/log/bounce"
)

func (d *Database) GetBouncesByUser(ctx context.Context, userID string) ([]model.Bounce, error) {
	var bounces []model.Bounce
	err := d.Client.Where("user_id = ?", userID).Find(&bounces).Error
	return bounces, err
}

func (d *Database) PostBounce(ctx context.Context, bounce model.Bounce) error {
	return d.Client.Create(&bounce).Error
}

func (d *Database) DeleteBounceByUserID(ctx context.Context, userID string) error {
	return d.Client.Where("user_id = ?", userID).Delete(&model.Bounce{}).Error
}

func (d *Database) SaveBounceToFile(ctx context.Context, filename string, data []byte) error {
	filePath := basePath + "/" + filename + ".eml"

	if err := os.WriteFile(filePath, data, 0600); err != nil {
		log.Println("error writing bounce file:", err)
		return err
	}

	return nil
}

func (d *Database) GetBounceFile(ctx context.Context, filename string) ([]byte, error) {
	filePath := basePath + "/" + filename + ".eml"
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Println("error reading bounce file:", err)
		return nil, err
	}

	return data, nil
}
