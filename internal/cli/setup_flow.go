package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/tekibo/idev/internal/manager"
)

func takeFlowSelection(args []string) (manager.FlowSelection, []string, bool, error) {
	selection := manager.FlowSelection{Expo: true}
	rest := []string{}
	expoSet := false
	for index := 0; index < len(args); index++ {
		switch args[index] {
		case "--flow":
			if index+1 >= len(args) {
				return selection, nil, false, fmt.Errorf("--flow needs a name")
			}
			selection.Flows = append(selection.Flows, args[index+1])
			index++
		case "--no-expo":
			selection.Expo = false
			expoSet = true
		default:
			rest = append(rest, args[index])
		}
	}
	return selection, rest, expoSet, nil
}

func setupSaved(force, expoSet bool, selection manager.FlowSelection, dryRun bool) int {
	manager := manager.New()
	saved, err := manager.LoadFlowSelection()
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev setup: %v\n", err)
		return 1
	}
	if expoSet {
		saved.Expo = selection.Expo
	}
	if len(saved.Flows) == 0 {
		fmt.Fprintln(os.Stderr, "idev setup: no saved flows; run idev setup --flow <name>")
		return 2
	}
	return setupFlow(saved, force, dryRun)
}

func setupFlow(selection manager.FlowSelection, force, dryRun bool) int {
	manager := manager.New()
	plan, err := manager.SetupFlowPlan(selection)
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev setup: %v\n", err)
		return 1
	}
	fmt.Printf("Flows: %s\n", strings.Join(plan.Flows, ", "))
	fmt.Printf("Expo: %v\n", plan.Expo)
	fmt.Printf("Payloads: %s\n", strings.Join(plan.Payloads, ", "))
	fmt.Printf("Components: %s\n", strings.Join(plan.Components, ", "))
	if dryRun {
		return 0
	}
	if err := manager.SaveFlowSelection(selection); err != nil {
		fmt.Fprintf(os.Stderr, "idev setup: %v\n", err)
		return 1
	}
	if err := manager.ApplySetupFlow(context.Background(), plan, force, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "idev setup: %v\n", err)
		return 1
	}
	reportSDK(manager)
	return 0
}
