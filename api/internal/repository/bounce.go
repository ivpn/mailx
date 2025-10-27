package repository

import (
	"context"
	"errors"
	"log"
	"os"
	"path/filepath"
	"regexp"

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
	data, err := readBounceEmail(basePath, filename, ".eml")
	if err != nil {
		return nil, err
	}

	return data, nil
}

func readBounceEmail(basePath, filename string, ext string) ([]byte, error) {
	uuidPattern := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)
	if !uuidPattern.MatchString(filename) {
		return nil, errors.New("invalid filename format (must be UUID)")
	}

	filePath := filepath.Join(basePath, filename+ext)

	return os.ReadFile(filePath)
}
