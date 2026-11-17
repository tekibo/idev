package cocoapods

import "errors"

var (
	// ErrGemMissing reports that the system Ruby and gem command are absent.
	ErrGemMissing = errors.New("system Ruby and gem not found; run idev system plan")
	// ErrVersionMismatch reports a CocoaPods version outside the pin.
	ErrVersionMismatch = errors.New("CocoaPods version mismatch")
)
