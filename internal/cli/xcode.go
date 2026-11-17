package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/tekibo/idev/internal/manager"
)

// Xcode manages the idev-owned Xcode.xip folder.
func Xcode(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: idev xcode <list|import|select|remove> [args]")
		return 2
	}
	manager := manager.New()
	switch args[0] {
	case "list":
		return xcodeList(manager, args[1:])
	case "import":
		return xcodeImport(manager, args[1:])
	case "select":
		return xcodeSelect(manager, args[1:])
	case "remove":
		return xcodeRemove(manager, args[1:])
	case "-h", "--help":
		fmt.Println("Usage: idev xcode <list|import|select|remove> [args]")
		return 0
	default:
		fmt.Fprintf(os.Stderr, "idev xcode: unknown command %q\n", args[0])
		return 2
	}
}

func xcodeList(manager *manager.App, args []string) int {
	asJSON := false
	for _, arg := range args {
		if arg == "--json" {
			asJSON = true
		} else {
			fmt.Fprintf(os.Stderr, "idev xcode list: unknown option %q\n", arg)
			return 2
		}
	}
	list, err := manager.XcodeFiles()
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev xcode: %v\n", err)
		return 1
	}
	return printXcodeList(list, asJSON)
}

func xcodeImport(manager *manager.App, args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: idev xcode import <path> [--replace]")
		return 2
	}
	source, mode := args[0], "add"
	asJSON := false
	for _, arg := range args[1:] {
		switch arg {
		case "--replace":
			mode = "replace"
		case "--json":
			asJSON = true
		default:
			fmt.Fprintf(os.Stderr, "idev xcode import: unknown option %q\n", arg)
			return 2
		}
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	list, err := manager.XcodeImport(ctx, source, mode, func(line string) {
		fmt.Println(line)
	}, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev xcode: %v\n", err)
		return 1
	}
	return printXcodeList(list, asJSON)
}
