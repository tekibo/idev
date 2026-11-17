package official

import "errors"

var (
	// ErrMissing reports an empty component.
	ErrMissing = errors.New("official binary component is missing")
	// ErrVersionMismatch reports a version outside the pin.
	ErrVersionMismatch = errors.New("official binary version mismatch")
)
