package swiftbuild

import "errors"

var (
	// ErrSwiftMissing reports that the private Swift toolchain is absent.
	ErrSwiftMissing = errors.New("private Swift toolchain not installed; run idev setup --swift")
	// ErrEngineMissing reports that the Expo engine source is absent.
	ErrEngineMissing = errors.New("Expo engine source not found; run idev setup")
	// ErrSourceMissing reports an empty Swift Build component.
	ErrSourceMissing = errors.New("Swift Build source component is missing")
	// ErrBuildMissing reports a missing build product.
	ErrBuildMissing = errors.New("Swift Build output is missing")
)
