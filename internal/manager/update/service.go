// Package update owns the manager and component update operations.
package update

import (
	"github.com/tekibo/idev/internal/manager/components"
	"github.com/tekibo/idev/internal/manager/core"
)

// Service owns the update operations.
type Service struct {
	*core.Services
	Components *components.Service
}

// New builds the update service.
func New(services *core.Services, components *components.Service) *Service {
	return &Service{Services: services, Components: components}
}
