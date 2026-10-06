package cron

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestEvery_RunsRepeatedlyAfterInterval(t *testing.T) {
	var runs atomic.Int32
	stop := make(chan struct{})
	defer close(stop)

	every(20*time.Millisecond, "test", func() { runs.Add(1) }, stop)

	time.Sleep(5 * time.Millisecond)
	if n := runs.Load(); n != 0 {
		t.Fatalf("job ran %d times before the first interval elapsed", n)
	}

	deadline := time.Now().Add(time.Second)
	for runs.Load() < 2 && time.Now().Before(deadline) {
		time.Sleep(5 * time.Millisecond)
	}
	if n := runs.Load(); n < 2 {
		t.Fatalf("expected at least 2 runs, got %d", n)
	}
}

func TestEvery_RecoversPanicAndKeepsRunning(t *testing.T) {
	var runs atomic.Int32
	stop := make(chan struct{})
	defer close(stop)

	every(10*time.Millisecond, "panics", func() {
		runs.Add(1)
		panic("boom")
	}, stop)

	deadline := time.Now().Add(time.Second)
	for runs.Load() < 2 && time.Now().Before(deadline) {
		time.Sleep(5 * time.Millisecond)
	}
	if n := runs.Load(); n < 2 {
		t.Fatalf("expected job to keep running after a panic, got %d runs", n)
	}
}

func TestEvery_StopsWhenStopClosed(t *testing.T) {
	var runs atomic.Int32
	stop := make(chan struct{})

	every(10*time.Millisecond, "stops", func() { runs.Add(1) }, stop)
	close(stop)

	time.Sleep(50 * time.Millisecond)
	if n := runs.Load(); n > 1 {
		t.Fatalf("expected job to stop, got %d runs", n)
	}
}
