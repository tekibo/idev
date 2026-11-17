package pymobile

import "errors"

var (
	// ErrPythonMissing reports that python3 is absent.
	ErrPythonMissing = errors.New("python3 not found; run idev system plan")
	// ErrPipMissing reports that the venv pip is absent.
	ErrPipMissing = errors.New("venv pip not found")
	// ErrVersionMismatch reports a pymobiledevice3 version outside the pin.
	ErrVersionMismatch = errors.New("pymobiledevice3 version mismatch")
)
