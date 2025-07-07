package utils

import (
	"testing"
	"time"
)

func TestBackground(t *testing.T) {
	done := make(chan bool)

	// Test if the function runs in the background
	Background(func() {
		time.Sleep(1 * time.Second)
		done <- true
	})

	select {
	case <-done:
		// Test passed
	case <-time.After(2 * time.Second):
		t.Error("Background function did not complete in time")
	}

	// Test if the function recovers from panic
	defer func() {
		if r := recover(); r != nil {
			t.Error("Background function did not recover from panic")
		}
	}()

	Background(func() {
		panic("test panic")
	})
}
