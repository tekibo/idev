package cli

import (
	"encoding/json"
	"fmt"
	"os"
)

func agentsJSON(args []string) (bool, bool) {
	asJSON := false
	for _, arg := range args {
		switch arg {
		case "--json":
			asJSON = true
		case "-h", "--help":
			return false, false
		default:
			fmt.Fprintf(os.Stderr, "idev agents: unknown option %q\n", arg)
			return false, false
		}
	}
	return asJSON, true
}

func agentsPrint(value any) int {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev agents: %v\n", err)
		return 2
	}
	fmt.Println(string(data))
	return 0
}
