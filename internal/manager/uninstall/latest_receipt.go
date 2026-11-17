package uninstall

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/tekibo/idev/internal/platform/config"
)

// LatestReceipt reads the newest uninstall receipt.
func LatestReceipt() (Receipt, error) {
	directory := filepath.Join(config.UserHome(), ".local", "state", "idev")
	entries, err := os.ReadDir(directory)
	if err != nil {
		return Receipt{}, err
	}
	var names []string
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasPrefix(name, "uninstall-") || !strings.HasSuffix(name, ".json") {
			continue
		}
		names = append(names, name)
	}
	if len(names) == 0 {
		return Receipt{}, fmt.Errorf("uninstall: no receipt found")
	}
	sort.Strings(names)
	data, err := os.ReadFile(filepath.Join(directory, names[len(names)-1]))
	if err != nil {
		return Receipt{}, err
	}
	var receipt Receipt
	if err := json.Unmarshal(data, &receipt); err != nil {
		return Receipt{}, err
	}
	return receipt, nil
}
