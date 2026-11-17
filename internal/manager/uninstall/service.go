// Package uninstall owns the uninstall report and apply operations.
package uninstall

import (
	"github.com/tekibo/idev/internal/manager/core"
)

// Service owns the uninstall operations.
type Service struct {
	*core.Services
}

// New builds the uninstall service.
func New(services *core.Services) *Service {
	return &Service{Services: services}
}
