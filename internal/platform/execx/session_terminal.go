package execx

import (
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"unsafe"
)

// attachedTerminal owns the controlling terminal while an attached child is
// running. Docker ran expo as the foreground process; without this handoff a
// child in its own process group becomes a background job and its first tty
// read stops it with SIGTTIN.
type attachedTerminal struct {
	fd     uintptr
	parent int
	active bool
}

func attachTerminal(cmd *exec.Cmd, in io.Reader) *attachedTerminal {
	term := &attachedTerminal{}
	file, ok := in.(*os.File)
	if !ok {
		setGroup(cmd)
		return term
	}
	fd := file.Fd()
	foreground, ok := terminalPgrp(fd)
	if !ok || foreground != syscall.Getpgrp() {
		setGroup(cmd)
		return term
	}
	term.fd = fd
	term.parent = foreground
	term.active = true
	cmd.SysProcAttr = &syscall.SysProcAttr{Foreground: true, Ctty: int(fd)}
	return term
}

func (t *attachedTerminal) restore() error {
	if !t.active {
		return nil
	}
	signal.Ignore(syscall.SIGTTOU)
	defer signal.Reset(syscall.SIGTTOU)
	return terminalSetPgrp(t.fd, t.parent)
}

func terminalPgrp(fd uintptr) (int, bool) {
	var pgrp int32
	_, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, fd, syscall.TIOCGPGRP,
		uintptr(unsafe.Pointer(&pgrp)), 0, 0, 0)
	if errno != 0 {
		return 0, false
	}
	return int(pgrp), true
}

func terminalSetPgrp(fd uintptr, pgrp int) error {
	value := int32(pgrp)
	_, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, fd, syscall.TIOCSPGRP,
		uintptr(unsafe.Pointer(&value)), 0, 0, 0)
	if errno != 0 {
		return errno
	}
	return nil
}
