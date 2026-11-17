package xtool

import "errors"

// ErrMissing reports that xtool was not found.
var ErrMissing = errors.New("xtool not found")

// ErrSecretArg reports a secret value passed in argv.
var ErrSecretArg = errors.New("secret values must be entered at the prompt; do not pass them in argv")
