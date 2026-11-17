package swift

import (
	"os"
	"path/filepath"

	"github.com/tekibo/idev/internal/platform/config"
)

// rollback removes the Swift root and its manager links after a failed install.
func rollback(cfg config.Config, root string) {
	_ = os.RemoveAll(root)
	for _, name := range []string{"swift", "swiftc", "sourcekit-lsp"} {
		_ = os.Remove(filepath.Join(cfg.Path("bin"), name))
	}
	for _, name := range []string{"libncurses.so.6", "libpanel.so.6", "libform.so.6", "libtinfo.so.6"} {
		_ = os.Remove(filepath.Join(cfg.Path("lib", "compat"), name))
	}
}
