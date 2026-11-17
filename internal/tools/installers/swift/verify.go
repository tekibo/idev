package swift

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/tekibo/idev/internal/platform/config"
	"github.com/tekibo/idev/internal/platform/execx"
)

// Result reports one install.
type Result struct {
	Version string
	Root    string
	Cache   string
	Already bool
}

// Verify runs the private Swift and checks SourceKit-LSP presence.
func Verify(ctx context.Context, cfg config.Config, runner execx.Runner) error {
	swift := cfg.Path("bin", "swift")
	output, err := runner.Output(execx.Command{
		Name: swift,
		Args: []string{"--version"},
		Env:  cfg.SwiftEnv(),
	})
	if err != nil {
		return fmt.Errorf("swift --version: %w", err)
	}
	if !strings.Contains(string(output), cfg.Sources.Swift.Version) {
		return fmt.Errorf("swift version mismatch: %s", strings.TrimSpace(string(output)))
	}
	lsp := cfg.Path("bin", "sourcekit-lsp")
	if !config.IsExec(lsp) {
		return fmt.Errorf("sourcekit-lsp missing at %s", lsp)
	}
	linker := filepath.Join(cfg.SwiftBin(), "ld64.lld")
	if !config.IsExec(linker) {
		return fmt.Errorf("ld64.lld missing at %s", linker)
	}
	return nil
}
