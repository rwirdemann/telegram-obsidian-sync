package state

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// State holds sync progress across restarts.
type State struct {
	Offset int `json:"offset"`
}

// Load reads state from dir/state.json. Returns an empty State
// if the file does not exist.
func Load(dir string) (*State, error) {
	path := filepath.Join(dir, "state.json")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &State{}, nil
	}
	if err != nil {
		return nil, err
	}
	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// Save writes state to dir/state.json.
func Save(dir string, s *State) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	path := filepath.Join(dir, "state.json")
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
