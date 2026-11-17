// Package setup owns the setup plans, payloads, and system packages.
package setup

import (
	"github.com/tekibo/idev/internal/manager/components"
	"github.com/tekibo/idev/internal/manager/core"
)

// FlowSelection is the setup group selection.
type FlowSelection struct {
	Flows []string
	Expo  bool
}

// Service owns the setup operations.
type Service struct {
	*core.Services
	Components *components.Service
}

// New builds the setup service.
func New(services *core.Services, components *components.Service) *Service {
	return &Service{Services: services, Components: components}
}
