package cli

import (
	"fmt"

	"github.com/tekibo/idev/internal/manager"
)

func printUpdateHints(app *manager.App) {
	if plan := app.CleanPlan(); plan.Bytes > 0 {
		fmt.Printf("cache: %s; run idev update clean\n", cleanSize(plan.Bytes))
	}
}
