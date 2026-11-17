// Package project owns the project inspection and tooling operations.
package project

import (
	"github.com/tekibo/idev/internal/manager/core"
)

// Service owns the project operations.
type Service struct {
	*core.Services
}

// New builds the project service.
func New(services *core.Services) *Service {
	return &Service{Services: services}
}
