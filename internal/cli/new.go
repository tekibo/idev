package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tekibo/idev/internal/manager"
)

// NewProject runs xtool new and applies the requested bundle ID.
func NewProject(args []string) int {
	bundleID := ""
	rest := []string{}
	skipSetup := false
	for index := 0; index < len(args); index++ {
		switch args[index] {
		case "--bundle-id":
			index++
			if index >= len(args) {
				fmt.Fprintln(os.Stderr, "idev new: --bundle-id needs a value")
				return 2
			}
			bundleID = args[index]
		case "--skip-setup":
			skipSetup = true
			rest = append(rest, args[index])
		default:
			rest = append(rest, args[index])
		}
	}
	if len(rest) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: idev new <name> [--bundle-id <id>]")
		return 2
	}
	if bundleID != "" {
		if err := manager.ValidBundleID(bundleID); err != nil {
			fmt.Fprintf(os.Stderr, "idev new: %v\n", err)
			return 2
		}
	}
	if !skipSetup {
		rest = append(rest, "--skip-setup")
	}
	code := runXtool(append([]string{"new"}, rest...))
	if code != 0 || bundleID == "" {
		return code
	}
	name := ""
	for _, arg := range rest {
		if !strings.HasPrefix(arg, "-") {
			name = arg
			break
		}
	}
	if name == "" {
		fmt.Fprintln(os.Stderr, "idev new: the project name is required")
		return 2
	}
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev new: %v\n", err)
		return 1
	}
	if err := manager.New().SetBundleID(filepath.Join(cwd, name), bundleID); err != nil {
		fmt.Fprintf(os.Stderr, "idev new: %v\n", err)
		return 1
	}
	return 0
}
