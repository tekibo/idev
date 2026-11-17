package privileged

import "errors"

// ErrPasswordRequired reports that sudo needs a password and no provider exists.
var ErrPasswordRequired = errors.New("sudo password required")

// ErrEmptyCommand reports a plan without a command.
var ErrEmptyCommand = errors.New("privileged command is empty")

// ErrPasswordRejected reports that sudo asked for the password again.
var ErrPasswordRejected = errors.New("sudo password rejected")

// ErrPromptTimeout reports that sudo did not show a password prompt in time.
var ErrPromptTimeout = errors.New("sudo password prompt timed out")
