package execx

import (
	"io"
	"os"
	"syscall"
	"unsafe"
)

// PromptTerminal restores canonical line input for a command that prompts on
// the controlling terminal, then restores the caller's mode. Metro leaves the
// tty in raw mode; xtool's yes/no prompt needs ICANON and ECHO.
func PromptTerminal(in io.Reader) (func() error, error) {
	file, ok := in.(*os.File)
	if !ok {
		return func() error { return nil }, nil
	}
	fd := file.Fd()
	var original syscall.Termios
	if err := ioctlTermios(fd, syscall.TCGETS, &original); err != nil {
		if err == syscall.ENOTTY {
			return func() error { return nil }, nil
		}
		return nil, err
	}
	sane := original
	sane.Iflag |= syscall.ICRNL | syscall.IXON
	sane.Lflag |= syscall.ICANON | syscall.ECHO | syscall.ECHOE | syscall.ECHOK |
		syscall.ECHOCTL | syscall.ECHOKE | syscall.ISIG | syscall.IEXTEN
	sane.Cc[syscall.VMIN] = 1
	sane.Cc[syscall.VTIME] = 0
	if err := ioctlTermios(fd, syscall.TCSETS, &sane); err != nil {
		return nil, err
	}
	return func() error { return ioctlTermios(fd, syscall.TCSETS, &original) }, nil
}

func ioctlTermios(fd uintptr, request uintptr, mode *syscall.Termios) error {
	_, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, fd, request,
		uintptr(unsafe.Pointer(mode)), 0, 0, 0)
	if errno != 0 {
		return errno
	}
	return nil
}
