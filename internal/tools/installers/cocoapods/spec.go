// Package cocoapods installs the private CocoaPods gem tree.
package cocoapods

import "github.com/tekibo/idev/internal/platform/config"

// Spec is the private gem tree install plan.
type Spec struct {
	Version string
	Erb     string
	Root    string
	Bin     string
	Link    string
	Cache   string
}

// NewSpec returns the private gem tree paths.
func NewSpec(cfg config.Config) Spec {
	return Spec{
		Version: cfg.Lock.Idev.Cocoapods,
		Erb:     cfg.Lock.Idev.Erb,
		Root:    cfg.GemRoot(),
		Bin:     cfg.PodBin(),
		Link:    cfg.PodLink(),
		Cache:   cfg.GemCache(),
	}
}
