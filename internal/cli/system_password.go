package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
)

// stdinPassword reads one sudo password from the JSON session stdin.
type stdinPassword struct{}

func (stdinPassword) Password(ctx context.Context) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	fmt.Fprint(os.Stderr, "sudo password: ")
	line, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return nil, err
	}
	return []byte(strings.TrimRight(line, "\r\n")), nil
}
