package manager

import (
	"context"

	"github.com/tekibo/idev/internal/manager/agents"
	"github.com/tekibo/idev/internal/manager/next"
	"github.com/tekibo/idev/internal/manager/xcode"
)

// AgentInstall reports one skill install.
type AgentInstall = agents.Result

// AgentStatus reports the installed skill and links.
type AgentStatus = agents.Status

// AgentUninstall reports one skill removal.
type AgentUninstall = agents.UninstallResult

// InstallAgents installs the idev skill and instruction snippet.
func (a *App) InstallAgents() (AgentInstall, error) {
	return agents.Install(a.Config)
}

// UninstallAgents removes the idev skill and its links.
func (a *App) UninstallAgents() (AgentUninstall, error) {
	return agents.Uninstall(a.Config)
}

// AgentStatus reads the current agent install state.
func (a *App) AgentStatus() AgentStatus {
	return agents.Inspect(a.Config)
}

// NextPlan is the agent-facing setup state.
type NextPlan = next.Plan

// NextStep is one setup state line.
type NextStep = next.Step

// NextPlan reads the setup state and the next command.
func (a *App) NextPlan(project string) NextPlan {
	return next.Build(a.Config, project)
}

// XcodeFiles lists the idev-owned xip folder.
func (a *App) XcodeFiles() (xcode.List, error) {
	return xcode.Files(a.Config)
}

// XcodeSelect records one xip as current.
func (a *App) XcodeSelect(path string) (xcode.List, error) {
	return xcode.Select(a.Config, path)
}

// XcodeRemove deletes one xip and selects the next file.
func (a *App) XcodeRemove(path string) (xcode.List, error) {
	return xcode.Remove(a.Config, path)
}

// XcodeImport moves one chosen xip into the idev-owned folder.
func (a *App) XcodeImport(ctx context.Context, source, mode string, log func(string),
	progress func(done, total int64)) (xcode.List, error) {
	return xcode.Import(ctx, a.Config, source, mode, log, progress)
}
