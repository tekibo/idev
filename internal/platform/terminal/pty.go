//go:build linux

// Package terminal provides the minimal Linux pseudo-terminal, termios,
// terminal-text, and password primitives used by the CLI and privileged runner.
package terminal

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"unsafe"
)

// PTY owns the master and slave sides of one pseudo-terminal.
type PTY struct {
	Master *os.File
	Slave  *os.File
}

// Open creates a pseudo-terminal pair.
func Open() (*PTY, error) {
	master, err := os.OpenFile("/dev/ptmx", os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		return nil, err
	}
	result := &PTY{Master: master}
	if err := result.initialize(); err != nil {
		_ = result.Close()
		return nil, err
	}
	return result, nil
}

func (p *PTY) initialize() error {
	var number uint32
	if err := ioctl(p.Master.Fd(), syscall.TIOCGPTN, unsafe.Pointer(&number)); err != nil {
		return err
	}
	var unlock int32
	if err := ioctl(p.Master.Fd(), syscall.TIOCSPTLCK, unsafe.Pointer(&unlock)); err != nil {
		return err
	}
	slave, err := os.OpenFile(fmt.Sprintf("/dev/pts/%d", number), os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		return err
	}
	p.Slave = slave
	return nil
}

// DisableEcho turns off terminal echo on the slave side.
func (p *PTY) DisableEcho() error {
	if p.Slave == nil {
		return errors.New("pty: slave is closed")
	}
	state, err := Get(p.Slave.Fd())
	if err != nil {
		return err
	}
	return Set(p.Slave.Fd(), WithoutEcho(state))
}

// Start attaches the slave side to the command and starts it in a new session.
func (p *PTY) Start(cmd *exec.Cmd) error {
	if p.Slave == nil {
		return errors.New("pty: slave is closed")
	}
	cmd.Stdin = p.Slave
	cmd.Stdout = p.Slave
	cmd.Stderr = p.Slave
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true, Setctty: true, Ctty: 0}
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := p.Slave.Close(); err != nil {
		return err
	}
	p.Slave = nil
	return nil
}

// Close closes both sides.
func (p *PTY) Close() error {
	var errs []error
	if p.Slave != nil {
		errs = append(errs, p.Slave.Close())
		p.Slave = nil
	}
	if p.Master != nil {
		errs = append(errs, p.Master.Close())
		p.Master = nil
	}
	return errors.Join(errs...)
}

func ioctl(fd uintptr, request uintptr, arg unsafe.Pointer) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, request, uintptr(arg))
	if errno != 0 {
		return errno
	}
	return nil
}
