package repository

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"ivpn.net/email/api/internal/model"
)

const (
	baseDir = "/var/log/bounce"
)

func (d *Database) GetBouncesByUser(ctx context.Context, userID string) ([]model.Bounce, error) {
	var bounces []model.Bounce
	err := d.Client.Where("user_id = ?", userID).Find(&bounces).Error
	return bounces, err
}

func (d *Database) GetBounce(ctx context.Context, bounceID string, userId string) (model.Bounce, error) {
	var bounce model.Bounce
	err := d.Client.Where("id = ? AND user_id = ?", bounceID, userId).First(&bounce).Error
	return bounce, err
}

func (d *Database) PostBounce(ctx context.Context, bounce model.Bounce) error {
	return d.Client.Create(&bounce).Error
}

func (d *Database) DeleteBounceByUserID(ctx context.Context, userID string) error {
	return d.Client.Where("user_id = ?", userID).Delete(&model.Bounce{}).Error
}

func (d *Database) SaveBounceToFile(ctx context.Context, filename string, data []byte) error {
	filePath := baseDir + "/" + filename + ".eml"

	if err := os.WriteFile(filePath, data, 0600); err != nil {
		log.Println("error writing bounce file:", err)
		return err
	}

	return nil
}

func (d *Database) GetBounceFile(ctx context.Context, filename string) ([]byte, error) {
	data, err := readFile(filename, ".eml")
	if err != nil {
		return nil, err
	}

	return data, nil
}

func readFile(filename, ext string) ([]byte, error) {
	filePath := filename + ext
	cleanPath := filepath.Clean(filePath)
	fullPath := filepath.Join(baseDir, cleanPath)

	// Ensure fullPath is still within baseDir
	if !strings.HasPrefix(fullPath, baseDir+string(os.PathSeparator)) {
		return nil, fmt.Errorf("invalid file path: %s", filePath)
	}

	f, err := os.Open(fullPath)
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
