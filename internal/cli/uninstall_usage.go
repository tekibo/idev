package cli

import (
	"fmt"
	"os"
	"strings"
)

// uninstallUsage prints one usage error.
func uninstallUsage(message string) int {
	fmt.Fprintf(os.Stderr, "idev uninstall: %s\n", message)
	fmt.Fprintln(os.Stderr, "Usage: idev uninstall [--dry-run|--latest] [--json] [--keep <id>]... [--confirm \"remove idev\"]")
	return 2
}

// splitKeep parses one comma-separated keep value.
func splitKeep(value string) []string {
	keep := []string{}
	for _, part := range strings.Split(value, ",") {
		part = strings.TrimSpace(part)
		if part != "" {
			keep = append(keep, part)
		}
	}
	return keep
}
