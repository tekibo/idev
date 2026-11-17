// Package system probes system dependencies and builds the exact package plan.
package system

// Check is one system dependency.
type Check struct {
	ID        string `json:"id"`
	Class     string `json:"class"`
	Status    string `json:"status"`
	FixableBy string `json:"fixable_by"`
	Detail    string `json:"detail"`
	Hint      string `json:"hint"`
	Package   string `json:"package,omitempty"`
	Critical  bool   `json:"critical"`
}

// Status values.
const (
	StatusOK      = "ok"
	StatusMissing = "missing"
	StatusWarn    = "warn"
)

// Class values.
const ClassSystem = "system"

// Fix owners.
const (
	FixUser = "user"
	FixIdev = "idev"
)

// OK reports a passing check.
func (c Check) OK() bool {
	return c.Status == StatusOK
}
