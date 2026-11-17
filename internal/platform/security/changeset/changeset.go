// Package changeset records applied manager changes.
package changeset

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Change is one applied change.
type Change struct {
	ID       string    `json:"id"`
	Kind     string    `json:"kind"`
	Command  string    `json:"command"`
	Packages []string  `json:"packages,omitempty"`
	Time     time.Time `json:"time"`
	Status   string    `json:"status"`
}

// Store writes change records under a directory.
type Store struct {
	Dir string
}

// New returns a store under $IDEV_HOME/state/changesets.
func New(home string) Store {
	return Store{Dir: filepath.Join(home, "state", "changesets")}
}

// Append writes one change atomically.
func (s Store) Append(change Change) error {
	if change.Time.IsZero() {
		change.Time = time.Now().UTC()
	}
	if change.Status == "" {
		change.Status = "applied"
	}
	if err := os.MkdirAll(s.Dir, 0755); err != nil {
		return err
	}
	file, err := os.CreateTemp(s.Dir, "change-*.tmp")
	if err != nil {
		return err
	}
	name := file.Name()
	if err := json.NewEncoder(file).Encode(change); err != nil {
		_ = file.Close()
		_ = os.Remove(name)
		return err
	}
	if err := file.Close(); err != nil {
		_ = os.Remove(name)
		return err
	}
	return os.Rename(name, strings.TrimSuffix(name, ".tmp")+".json")
}

// List reads every recorded change in name order.
func (s Store) List() ([]Change, error) {
	entries, err := os.ReadDir(s.Dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var changes []Change
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(s.Dir, entry.Name()))
		if err != nil {
			return nil, err
		}
		var change Change
		if err := json.Unmarshal(data, &change); err != nil {
			return nil, err
		}
		changes = append(changes, change)
	}
	return changes, nil
}
