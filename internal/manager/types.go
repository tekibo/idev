package manager

import (
	"github.com/tekibo/idev/internal/manager/auth"
	"github.com/tekibo/idev/internal/manager/doctor"
	"github.com/tekibo/idev/internal/manager/setup"
	"github.com/tekibo/idev/internal/manager/uninstall"
	"github.com/tekibo/idev/internal/manager/xcode"
	"github.com/tekibo/idev/internal/platform/project"
	"github.com/tekibo/idev/internal/platform/system"
)

// UninstallServiceReport is the terminal uninstall report.
type UninstallServiceReport = uninstall.Report

// AuthStatus is the parsed Apple account status.
type AuthStatus = auth.AuthStatus

// XcodeXIP is one xip in the idev-owned folder.
type XcodeXIP = xcode.XIP

// XcodeXIPList is the xip folder state.
type XcodeXIPList = xcode.List

// SetupPlan is the resolved setup plan.
type SetupPlan = setup.Plan

// SystemPlan is the system transaction view.
type SystemPlan = setup.SystemPlan

// FlowSelection is the setup group selection.
type FlowSelection = setup.FlowSelection

// UninstallReceipt is one uninstall receipt.
type UninstallReceipt = uninstall.Receipt

// UninstallReport is the uninstall report.
type UninstallReport = uninstall.Report

// DoctorCheck is one diagnostic result.
type DoctorCheck = system.Check

// DoctorReport is the diagnostic snapshot.
type DoctorReport = doctor.Report

// Project is one detected project.
type Project = project.Info

// Workspace is one opened directory and its detected apps.
type Workspace = project.Workspace

// ValidBundleID checks the reverse-DNS form used by Apple.
func ValidBundleID(id string) error {
	return project.ValidBundleID(id)
}

// Kind is one supported project kind.
type Kind = project.Kind

// Node workflows.
const (
	WorkflowExpoManaged   = project.WorkflowExpoManaged
	WorkflowExpoPrebuild  = project.WorkflowExpoPrebuild
	WorkflowBareRN        = project.WorkflowBareRN
	WorkflowBareRNMissing = project.WorkflowBareRNMissing
)

// Project kinds.
const (
	KindSwiftUI      = project.SwiftUI
	KindFlutter      = project.Flutter
	KindExpo         = project.Expo
	KindReactNative  = project.ReactNative
	KindXcode        = project.Xcode
	KindSwiftPackage = project.SwiftPackage
	KindUnknown      = project.Unknown
)
