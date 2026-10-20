package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type State struct {
	LastPrice float64 `json:"last_price"`
}

type FileStore struct {
	path string
}

func NewFileStore(path string) *FileStore {
	return &FileStore{path: path}
}

func (s *FileStore) Load() (State, bool, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return State{}, false, nil
	}
	if err != nil {
		return State{}, false, fmt.Errorf("read state file: %w", err)
	}

	var st State
	if err := json.Unmarshal(data, &st); err != nil {
		return State{}, false, fmt.Errorf("parse state file: %w", err)
	}
	return st, true, nil
}

func (s *FileStore) Save(st State) error {
	data, err := json.Marshal(st)
	if err != nil {
		return fmt.Errorf("encode state: %w", err)
	}
	if err := os.WriteFile(s.path, data, 0o600); err != nil {
		return fmt.Errorf("write state file: %w", err)
	}
	return nil
}
