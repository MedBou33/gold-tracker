package store

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFirstRun(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	s := NewFileStore(path)

	_, found, err := s.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found {
		t.Error("found = true, want false on first run")
	}
}

func TestSaveThenLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	s := NewFileStore(path)

	err := s.Save(State{LastPrice: 3500.5})
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}

	got, found, err := s.Load()
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if !found {
		t.Fatal("found = false, want true after save")
	}
	if got.LastPrice != 3500.5 {
		t.Errorf("LastPrice = %v, want 3500.5", got.LastPrice)
	}
}

func TestLoadBrokenFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	err := os.WriteFile(path, []byte("not json"), 0o600)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	s := NewFileStore(path)
	_, _, err = s.Load()
	if err == nil {
		t.Error("expected an error for a broken file, got nil")
	}
}
