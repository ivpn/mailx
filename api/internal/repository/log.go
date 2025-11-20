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
	logEmlBaseDir = "/var/log/eml"
)

func (d *Database) GetLogs(ctx context.Context, userID string) ([]model.Log, error) {
	var logs []model.Log
	err := d.Client.Where("user_id = ?", userID).Order("created_at desc").Find(&logs).Error
	return logs, err
}

func (d *Database) GetLog(ctx context.Context, logID string, userId string) (model.Log, error) {
	var log model.Log
	err := d.Client.Where("id = ? AND user_id = ?", logID, userId).First(&log).Error
	return log, err
}

func (d *Database) PostLog(ctx context.Context, log model.Log) error {
	return d.Client.Create(&log).Error
}

func (d *Database) DeleteLogs(ctx context.Context, userID string) error {
	return d.Client.Where("user_id = ?", userID).Delete(&model.Log{}).Error
}

func (d *Database) SaveLogToFile(ctx context.Context, filename string, data []byte) error {
	filePath := logEmlBaseDir + "/" + filename + ".eml"

	// Ensure the directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0750); err != nil {
		log.Println("error creating bounce directory:", err)
		return err
	}

	// Write the file
	if err := os.WriteFile(filePath, data, 0600); err != nil {
		log.Println("error writing bounce file:", err)
		return err
	}

	return nil
}

func (d *Database) GetLogFile(ctx context.Context, filename string) ([]byte, error) {
	data, err := readFile(filename, ".eml")
	if err != nil {
		return nil, err
	}

	return data, nil
}

func readFile(filename, ext string) ([]byte, error) {
	filePath := filename + ext
	cleanPath := filepath.Clean(filePath)
	fullPath := filepath.Join(logEmlBaseDir, cleanPath)

	// Ensure fullPath is still within baseDir
	if !strings.HasPrefix(fullPath, logEmlBaseDir+string(os.PathSeparator)) {
		return nil, fmt.Errorf("invalid file path: %s", filePath)
	}

	file, err := os.Open(fullPath)
	if err != nil {
		log.Println("failed to open file:", err)
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Println("failed to read file:", err)
		return nil, err
	}

	return data, nil
}
