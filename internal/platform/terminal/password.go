//go:build linux

// Package tty provides the minimal terminal password reader used by the CLI.
// It is a scoped local implementation.
package terminal

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// ReadPassword writes a prompt and reads one line with echo disabled.
func ReadPassword(prompt string, in *os.File, out io.Writer) ([]byte, error) {
	state, err := Get(in.Fd())
	if err != nil {
		return nil, err
	}
	if err := Set(in.Fd(), WithoutEcho(state)); err != nil {
		return nil, err
	}
	defer func() { _ = Set(in.Fd(), state) }()

	if _, err := io.WriteString(out, prompt); err != nil {
		return nil, err
	}
	line, err := bufio.NewReader(in).ReadString('\n')
	if err != nil {
		return nil, err
	}
	if _, err := io.WriteString(out, "\n"); err != nil {
		return nil, err
	}
	return []byte(strings.TrimRight(line, "\r\n")), nil
}
