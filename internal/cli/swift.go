package cli

import (
	"fmt"
	"os"
)

// Swift runs the xtool SwiftPM flow with the build environment.
func Swift(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: idev swift <build|run|new|generate-xcode-project|build-server> [args...]")
		return 2
	}
	switch args[0] {
	case "build", "run", "generate-xcode-project", "build-server":
		return runXtool(append([]string{"dev", args[0]}, args[1:]...))
	case "new":
		return runXtool(append([]string{"new"}, args[1:]...))
	default:
		fmt.Fprintf(os.Stderr, "idev swift: unknown command %q\n", args[0])
		return 2
	}
}
