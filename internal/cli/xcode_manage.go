package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tekibo/idev/internal/manager"
)

func xcodeSelect(manager *manager.App, args []string) int {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Usage: idev xcode select <path>")
		return 2
	}
	list, err := manager.XcodeSelect(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev xcode: %v\n", err)
		return 1
	}
	return printXcodeList(list, false)
}

func xcodeRemove(manager *manager.App, args []string) int {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Usage: idev xcode remove <path>")
		return 2
	}
	list, err := manager.XcodeRemove(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev xcode: %v\n", err)
		return 1
	}
	return printXcodeList(list, false)
}

func printXcodeList(list manager.XcodeXIPList, asJSON bool) int {
	if asJSON {
		data, err := json.MarshalIndent(list, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "idev xcode: %v\n", err)
			return 2
		}
		fmt.Println(string(data))
		return 0
	}
	for _, file := range list.Files {
		marker := " "
		if file.Current {
			marker = "*"
		}
		fmt.Printf("%s %s\n", marker, file.Path)
	}
	return 0
}
