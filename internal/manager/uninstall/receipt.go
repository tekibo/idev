package uninstall

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/tekibo/idev/internal/platform/config"
)

// Receipt is the applied uninstall record.
type Receipt struct {
	Time           string          `json:"time"`
	Home           string          `json:"home"`
	Path           string          `json:"path"`
	Removed        []string        `json:"removed"`
	Kept           []string        `json:"kept"`
	KeptComponents []string        `json:"kept_components"`
	Remaining      []string        `json:"remaining"`
	System         []SystemPackage `json:"system"`
	Configs        []Config        `json:"configs"`
	ProjectData    []string        `json:"project_data"`
	Failure        string          `json:"failure,omitempty"`
	SystemErr      string          `json:"system_error,omitempty"`
}

// keptComponents returns the sorted kept component IDs.
func keptComponents(components []Component) []string {
	ids := []string{}
	for _, component := range components {
		if component.Keep {
			ids = append(ids, component.ID)
		}
	}
	sort.Strings(ids)
	return ids
}

// receiptPath returns the receipt location outside the private home.
func receiptPath(now time.Time) string {
	name := "uninstall-" + now.Format("20060102T150405Z") + ".json"
	return filepath.Join(config.UserHome(), ".local", "state", "idev", name)
}

// writeReceipt writes the receipt after the removal attempt.
func writeReceipt(receipt Receipt) error {
	if err := os.MkdirAll(filepath.Dir(receipt.Path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(receipt, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(receipt.Path, data, 0644)
}
