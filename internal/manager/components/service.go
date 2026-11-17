// Package components installs the private toolchain components.
package components

import (
	"github.com/tekibo/idev/internal/manager/core"
)

// Service owns the component installers.
type Service struct {
	*core.Services
}

// New builds the component service.
func New(services *core.Services) *Service {
	return &Service{Services: services}
}
