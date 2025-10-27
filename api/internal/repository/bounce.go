package repository

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"

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
	data, err := readFile(basePath, filename, ".eml")
	if err != nil {
		return nil, err
	}

	return data, nil
}

func readFile(basePath, filename, ext string) ([]byte, error) {
	filePath := filepath.Join(basePath, filename+ext)

	f, err := os.Open(filePath)
	if err != nil {
		log.Println("failed to open file:", err)
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		log.Println("failed to read file:", err)
		return nil, err
	}

	return data, nil
}
