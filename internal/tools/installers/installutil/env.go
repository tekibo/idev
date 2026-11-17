package installutil

import (
	"strings"

	"github.com/tekibo/idev/internal/platform/config"
)

// SwiftEnv returns the private Swift environment without DEVELOPER_DIR.
// Extra keys replace inherited values and are appended.
func SwiftEnv(cfg config.Config, extra map[string]string) []string {
	env := cfg.SwiftEnv()
	filtered := env[:0]
	for _, entry := range env {
		key, _, _ := strings.Cut(entry, "=")
		if key == "DEVELOPER_DIR" {
			continue
		}
		if _, replace := extra[key]; replace {
			continue
		}
		filtered = append(filtered, entry)
	}
	for key, value := range extra {
		filtered = append(filtered, key+"="+value)
	}
	return filtered
}
