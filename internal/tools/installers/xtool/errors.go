package xtool

import "errors"

var (
	// ErrSwiftMissing reports that the private Swift toolchain is absent.
	ErrSwiftMissing = errors.New("private Swift toolchain not installed; run idev setup --swift")
	// ErrPinMismatch reports an xtool version outside the pin.
	ErrPinMismatch = errors.New("xtool version does not match the pin")
)
