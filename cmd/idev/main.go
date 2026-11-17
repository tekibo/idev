// Command idev runs the CLI.
package main

import (
	"os"

	"github.com/tekibo/idev/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:]))
}
