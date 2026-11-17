package privileged

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"time"
)

// watchPrompts reports every occurrence of the sudo prompt and forwards output
// after the first prompt. A closed channel means the PTY ended.
func watchPrompts(reader io.Reader, prompt []byte, out io.Writer) (<-chan struct{}, chan struct{}) {
	events := make(chan struct{}, 4)
	done := make(chan struct{})
	go func() {
		defer close(events)
		buffer := make([]byte, 0, 512)
		chunk := make([]byte, 256)
		seen := false
		for {
			n, err := reader.Read(chunk)
			if n > 0 {
				data := chunk[:n]
				if !seen {
					buffer = append(buffer, data...)
					index := bytes.Index(buffer, prompt)
					if index >= 0 {
						seen = true
						rest := append([]byte{}, buffer[index+len(prompt):]...)
						buffer = buffer[:0]
						select {
						case events <- struct{}{}:
						case <-done:
							return
						}
						writeOutput(out, rest)
					}
				} else {
					buffer = append(buffer, data...)
					for {
						index := bytes.Index(buffer, prompt)
						if index < 0 {
							break
						}
						buffer = buffer[index+len(prompt):]
						select {
						case events <- struct{}{}:
						case <-done:
							return
						}
					}
					if keep := len(prompt) - 1; keep > 0 && len(buffer) > keep {
						buffer = buffer[len(buffer)-keep:]
					}
					writeOutput(out, data)
				}
			}
			if err != nil {
				return
			}
		}
	}()
	return events, done
}

func writeOutput(out io.Writer, data []byte) {
	if out == nil || len(data) == 0 {
		return
	}
	_, _ = out.Write(data)
}

func waitForPrompt(ctx context.Context, timeout time.Duration, prompts <-chan struct{}, wait <-chan error) (bool, error) {
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case _, ok := <-prompts:
		if !ok {
			return false, authResult(<-wait)
		}
		return true, nil
	case err := <-wait:
		return false, authResult(err)
	case <-timer.C:
		return false, ErrPromptTimeout
	case <-ctx.Done():
		return false, ctx.Err()
	}
}

func waitForResult(ctx context.Context, prompts <-chan struct{}, wait <-chan error) error {
	select {
	case _, ok := <-prompts:
		if !ok {
			return authResult(<-wait)
		}
		return ErrPasswordRejected
	case err := <-wait:
		var exit *exec.ExitError
		if errors.As(err, &exit) && exit.ExitCode() == 125 {
			return authResult(err)
		}
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func authResult(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("sudo authentication failed: %w", err)
}
