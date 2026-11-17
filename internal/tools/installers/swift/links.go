package swift

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tekibo/idev/internal/platform/config"
)

// InstallLinks exposes the private tools under $IDEV_HOME/bin.
func InstallLinks(cfg config.Config) error {
	destination := cfg.Path("bin")
	if err := os.MkdirAll(destination, 0755); err != nil {
		return err
	}
	for _, name := range []string{"swift", "swiftc", "sourcekit-lsp"} {
		source := filepath.Join(cfg.SwiftBin(), name)
		if !config.IsExec(source) && !config.IsFile(source) {
			return fmt.Errorf("private Swift tool missing: %s", source)
		}
		target := filepath.Join(destination, name)
		_ = os.Remove(target)
		if err := os.Symlink(source, target); err != nil {
			return err
		}
	}
	return nil
}
