//go:build linux

package terminal

import (
	"syscall"
	"unsafe"
)

// Get reads the terminal attributes of fd.
func Get(fd uintptr) (syscall.Termios, error) {
	var state syscall.Termios
	err := ioctl(fd, syscall.TCGETS, unsafe.Pointer(&state))
	return state, err
}

// Set writes terminal attributes to fd.
func Set(fd uintptr, state syscall.Termios) error {
	return ioctl(fd, syscall.TCSETS, unsafe.Pointer(&state))
}

// WithoutEcho returns the attributes with local echo disabled.
func WithoutEcho(state syscall.Termios) syscall.Termios {
	state.Lflag &^= syscall.ECHO
	return state
}
