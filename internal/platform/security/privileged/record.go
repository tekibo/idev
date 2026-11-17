package privileged

import (
	"context"

	"github.com/tekibo/idev/internal/platform/security/changeset"
)

// Recorder records an applied plan.
type Recorder interface {
	Record(context.Context, Plan) error
}

// ChangeRecorder writes plans to the changeset store.
type ChangeRecorder struct {
	Store changeset.Store
}

// Record appends one applied plan.
func (r ChangeRecorder) Record(_ context.Context, plan Plan) error {
	return r.Store.Append(changeset.Change{
		Kind:     "system",
		Command:  plan.Display(),
		Packages: plan.Packages,
		Status:   "applied",
	})
}
