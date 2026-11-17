package swift

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tekibo/idev/internal/platform/config"
)

// The Swift 6.4.2 tarball ships a plain ld64.lld symlink. Xcode passes
// -fatal_warnings, which turns LLD's iOS-platform warning into an error. The
// Docker flow wrapped the linker; native installs must apply the same repair
// without replacing Swift or the Linux lld entry point.
const linkerWrapper = `#!/bin/sh
here=$(dirname "$(readlink -f "$0")")
root=$(cd "$here/../../../.." && pwd)
exec "$root/react-native/tools/ld64.lld" "$@"
`

// InstallLinkers writes the Darwin linker wrapper into the private toolchain.
func InstallLinkers(cfg config.Config) error {
	tool := cfg.Path("react-native", "tools", "ld64.lld")
	if !config.IsExec(tool) {
		return fmt.Errorf("React Native linker wrapper missing at %s; run idev setup", tool)
	}
	target := filepath.Join(cfg.SwiftBin(), "ld64.lld")
	if current, err := os.ReadFile(target); err == nil && string(current) == linkerWrapper {
		return nil
	}
	if err := os.WriteFile(target, []byte(linkerWrapper), 0755); err != nil {
		return fmt.Errorf("write %s: %w", target, err)
	}
	return os.Chmod(target, 0755)
}
