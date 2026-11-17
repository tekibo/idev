package help

import "strings"

// PreDispatch prints focused help before a command runs.
// Passthrough commands keep their own --help handling.
func PreDispatch(args []string, asJSON bool) (bool, int) {
	node := Root()
	path := []*Command{}
	index := 0
	for index < len(args) {
		if strings.HasPrefix(args[index], "-") {
			break
		}
		next := childByName(node, args[index])
		if next == nil {
			break
		}
		node = next
		path = append(path, next)
		index++
	}
	if len(path) == 0 {
		return false, 0
	}
	for _, part := range path {
		if part.Passthrough {
			return false, 0
		}
	}
	if !hasHelp(args[index:]) {
		return false, 0
	}
	if asJSON {
		return true, printJSON(node)
	}
	names := make([]string, 0, len(path))
	for _, part := range path {
		names = append(names, part.Name)
	}
	printText(node, names)
	return true, 0
}

func findPath(root *Command, path []string) (*Command, bool) {
	node := root
	for _, name := range path {
		next := childByName(node, name)
		if next == nil {
			return node, false
		}
		node = next
	}
	return node, true
}

func childByName(node *Command, name string) *Command {
	for _, child := range node.Children {
		if child.Name == name {
			return child
		}
		for _, alias := range child.Aliases {
			if alias == name {
				return child
			}
		}
	}
	return nil
}

func hasHelp(args []string) bool {
	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			return true
		}
	}
	return false
}
