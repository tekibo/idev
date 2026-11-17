// Package auth owns the Apple account operations.
package auth

import (
	"github.com/tekibo/idev/internal/manager/core"
)

// Service owns the auth operations.
type Service struct {
	*core.Services
}

// New builds the auth service.
func New(services *core.Services) *Service {
	return &Service{Services: services}
}
