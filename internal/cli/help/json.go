package help

import (
	"encoding/json"
	"fmt"
	"os"
)

type jsonCommand struct {
	Name        string         `json:"name"`
	Aliases     []string       `json:"aliases,omitempty"`
	Group       string         `json:"group,omitempty"`
	Summary     string         `json:"summary"`
	Usage       string         `json:"usage,omitempty"`
	Description string         `json:"description,omitempty"`
	Passthrough bool           `json:"passthrough,omitempty"`
	Options     []Option       `json:"options,omitempty"`
	Env         []Env          `json:"environment,omitempty"`
	Examples    []Example      `json:"examples,omitempty"`
	JSON        []string       `json:"json_fields,omitempty"`
	Children    []*jsonCommand `json:"commands,omitempty"`
}

func printJSON(node *Command) int {
	data, err := json.MarshalIndent(toJSON(node), "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev help: %v\n", err)
		return 2
	}
	fmt.Println(string(data))
	return 0
}

func toJSON(node *Command) *jsonCommand {
	out := &jsonCommand{
		Name:        node.Name,
		Aliases:     node.Aliases,
		Group:       node.Group,
		Summary:     node.Summary,
		Usage:       node.Usage,
		Description: node.Description,
		Passthrough: node.Passthrough,
		Options:     node.Options,
		Env:         node.Env,
		Examples:    node.Examples,
		JSON:        node.JSON,
	}
	for _, child := range node.Children {
		out.Children = append(out.Children, toJSON(child))
	}
	return out
}
