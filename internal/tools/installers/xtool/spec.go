// Package xtool builds the pinned xtool from source.
package xtool

import (
	"github.com/tekibo/idev/internal/platform/config"
)

// Spec is the pinned xtool install plan.
type Spec struct {
	Pin     string
	Repo    string
	Archive string
	SHA256  string
	Root    string
	Bin     string
	Source  string
	Link    string
}

// NewSpec returns the pinned xtool paths.
func NewSpec(cfg config.Config) Spec {
	return Spec{
		Pin:    cfg.Lock.AppleDev.XtoolSHA,
		Repo:   cfg.Lock.AppleDev.XtoolRepo,
		Root:   cfg.XtoolRoot(),
		Bin:    cfg.XtoolBin(),
		Source: cfg.XtoolSource(),
		Link:   cfg.XtoolLink(),
	}
}
