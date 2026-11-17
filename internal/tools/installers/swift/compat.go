package swift

import (
	"os"
	"path/filepath"

	"github.com/tekibo/idev/internal/platform/config"
)

// InstallCompat creates the private ncurses compatibility links.
func InstallCompat(cfg config.Config) error {
	links := map[string]string{
		"libncurses.so.6": "libncursesw.so.6",
		"libpanel.so.6":   "libpanelw.so.6",
		"libform.so.6":    "libformw.so.6",
		"libtinfo.so.6":   "libncursesw.so.6",
	}
	for _, directory := range []string{cfg.Path("lib", "compat"), cfg.SwiftLib()} {
		if err := os.MkdirAll(directory, 0755); err != nil {
			return err
		}
		for name, target := range links {
			path := filepath.Join(directory, name)
			_ = os.Remove(path)
			if err := os.Symlink(filepath.Join("/usr/lib", target), path); err != nil {
				return err
			}
		}
	}
	return nil
}
