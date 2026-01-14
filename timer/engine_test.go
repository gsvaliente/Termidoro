package timer

import (
	"testing"
	"time"
)

func TestEngine(t *testing.T) {
	engine := NewEngine()

	// Test AddSession
	engine.AddSession(10 * time.Minute)
	if len(engine.Sessions) != 1 {
		t.Errorf("Expected 1 session, got %d", len(engine.Sessions))
	}
	if engine.Sessions[0].Duration != 10*time.Minute {
		t.Errorf("Expected duration 10m, got %v", engine.Sessions[0].Duration)
	}

	// Test CompleteSession
	engine.CompleteSession(0)
	if !engine.Sessions[0].Completed {
		t.Error("Expected session to be completed")
	}
	if engine.TotalTime != 10*time.Minute {
		t.Errorf("Expected total time to be 10m, got %v", engine.TotalTime)
	}

	// Test CancelSession
	engine.AddSession(5 * time.Minute)
	engine.CancelSession(1)
	if !engine.Sessions[1].WasCancelled {
		t.Error("Expected session to be cancelled")
	}
	if engine.TotalTime != 10*time.Minute {
		t.Errorf("Expected total time to remain 10m, got %v", engine.TotalTime)
	}
}
