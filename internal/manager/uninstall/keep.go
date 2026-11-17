package uninstall

import "fmt"

// resolveKeep validates the requested keep list and applies dependencies.
func resolveKeep(components []Component, requested []string) (map[string]bool, error) {
	keep := map[string]bool{}
	installed := map[string]Component{}
	for _, component := range components {
		installed[component.ID] = component
	}
	for _, id := range requested {
		if id == "" {
			continue
		}
		component, ok := installed[id]
		if !ok {
			return nil, fmt.Errorf("uninstall: component %q is not installed", id)
		}
		if !component.Keepable {
			return nil, fmt.Errorf("uninstall: component %q cannot be kept", id)
		}
		keep[id] = true
	}
	if keep["xtool"] {
		keep["swift"] = true
	}
	return keep, nil
}
