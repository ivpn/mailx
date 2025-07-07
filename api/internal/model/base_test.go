package model

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestBeforeCreate(t *testing.T) {
	// Initialize a new in-memory SQLite database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	// Define a test model that embeds BaseModel
	type TestModel struct {
		BaseModel
		Name string
	}

	// Auto migrate the test model
	err = db.AutoMigrate(&TestModel{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	// Create a new instance of the test model
	testModel := TestModel{Name: "Test Name"}

	// Create the record in the database
	err = db.Create(&testModel).Error
	if err != nil {
		t.Fatalf("failed to create record: %v", err)
	}

	// Check if the ID was set and is a valid UUID
	_, err = uuid.Parse(testModel.ID)
	assert.NoError(t, err, "ID should be a valid UUID")
	assert.NotEmpty(t, testModel.ID, "ID should not be empty")
}
