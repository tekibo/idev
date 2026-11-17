// Package manager composes the manager services.
package manager

import (
	"github.com/tekibo/idev/internal/manager/auth"
	"github.com/tekibo/idev/internal/manager/components"
	"github.com/tekibo/idev/internal/manager/core"
	"github.com/tekibo/idev/internal/manager/devices"
	"github.com/tekibo/idev/internal/manager/project"
	"github.com/tekibo/idev/internal/manager/setup"
	"github.com/tekibo/idev/internal/manager/uninstall"
	"github.com/tekibo/idev/internal/manager/update"
	"github.com/tekibo/idev/internal/platform/redact"
	"github.com/tekibo/idev/internal/platform/terminal"
)

// ComponentsService is the component installer service.
type ComponentsService = components.Service

// SetupService is the setup service.
type SetupService = setup.Service

// ProjectService is the project service.
type ProjectService = project.Service

// DevicesService is the device service.
type DevicesService = devices.Service

// AuthService is the auth service.
type AuthService = auth.Service

// UpdateService is the update service.
type UpdateService = update.Service

// UninstallService is the uninstall service.
type UninstallService = uninstall.Service

// App composes the CLI services.
type App struct {
	*core.Services
	*ComponentsService
	*SetupService
	*ProjectService
	*DevicesService
	*AuthService
	*UpdateService
	*UninstallService
}

// New builds the default app.
func New() *App {
	services := core.New()
	componentsService := components.New(services)
	setupService := setup.New(services, componentsService)
	return &App{
		Services:          services,
		ComponentsService: componentsService,
		SetupService:      setupService,
		ProjectService:    project.New(services),
		DevicesService:    devices.New(services),
		AuthService:       auth.New(services),
		UpdateService:     update.New(services, componentsService),
		UninstallService:  uninstall.New(services),
	}
}

// RedactText removes secret values from terminal text.
func RedactText(text string) string {
	return redact.RedactText(text)
}

// IsTerminal reports whether the file descriptor is a terminal.
func (a *App) IsTerminal(fd uintptr) bool {
	_, err := terminal.Get(fd)
	return err == nil
}
