package help

import (
	"fmt"
	"strings"
)

var groupOrder = []string{
	"Project",
	"Build & Run",
	"Devices",
	"Setup",
	"System",
	"Maintenance",
	"Tools",
}

func printText(node *Command, path []string) {
	if node.Name == "idev" {
		printRoot(node)
		return
	}
	label := "idev " + strings.Join(path, " ")
	fmt.Printf("%s — %s\n\n", label, node.Summary)
	if node.Usage != "" {
		fmt.Printf("Usage:\n  %s\n\n", node.Usage)
	}
	if node.Description != "" {
		fmt.Printf("%s\n\n", strings.TrimSpace(node.Description))
	}
	if node.Passthrough {
		fmt.Println("The manager passes arguments to the underlying CLI.")
		fmt.Println()
	}
	if len(node.Children) > 0 {
		fmt.Println("Commands:")
		for _, child := range node.Children {
			fmt.Printf("  %-18s %s\n", child.Name, child.Summary)
		}
		fmt.Println()
	}
	writeOptions(node.Options)
	writeEnv(node.Env)
	writeExamples(node.Examples)
	if len(node.JSON) > 0 {
		fmt.Println("JSON fields:")
		for _, field := range node.JSON {
			fmt.Printf("  %s\n", field)
		}
		fmt.Println()
	}
	if len(node.Children) > 0 {
		fmt.Printf("Run 'idev help %s <command>' for details.\n", strings.Join(path, " "))
	}
}

func printRoot(node *Command) {
	fmt.Printf("%s\n\n", node.Summary)
	fmt.Printf("Usage:\n  %s\n\n", node.Usage)
	for _, group := range groupOrder {
		children := groupChildren(node, group)
		if len(children) == 0 {
			continue
		}
		fmt.Printf("%s:\n", group)
		for _, child := range children {
			fmt.Printf("  %-18s %s\n", child.Name, child.Summary)
		}
		fmt.Println()
	}
	writeOptions(node.Options)
	writeEnv(node.Env)
	fmt.Println("Run 'idev help <command>' for focused help.")
}

func groupChildren(node *Command, group string) []*Command {
	children := []*Command{}
	for _, child := range node.Children {
		if child.Group == group {
			children = append(children, child)
		}
	}
	return children
}
