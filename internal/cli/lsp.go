package cli

import (
	"fmt"
	"os"

	"github.com/tekibo/idev/internal/manager"
)

// LSP runs the private SourceKit-LSP.
func LSP(args []string) int {
	manager := manager.New()
	code, err := manager.RunLSP(args, os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev lsp: %v\n", err)
		if code == 0 {
			return 1
		}
	}
	return code
}
