package xcode

import (
	"context"
	"fmt"
	"io"
	"os"
)

func copyFile(ctx context.Context, source, destination string, total int64, progress func(done, total int64)) error {
	input, err := os.Open(source)
	if err != nil {
		return err
	}
	defer func() { _ = input.Close() }()
	output, err := os.OpenFile(destination, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	buffer := make([]byte, 1<<20)
	var done int64
	for {
		if err := ctx.Err(); err != nil {
			_ = output.Close()
			_ = os.Remove(destination)
			return err
		}
		count, readErr := input.Read(buffer)
		if count > 0 {
			if _, writeErr := output.Write(buffer[:count]); writeErr != nil {
				_ = output.Close()
				_ = os.Remove(destination)
				return writeErr
			}
			done += int64(count)
			if progress != nil {
				progress(done, total)
			}
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			_ = output.Close()
			_ = os.Remove(destination)
			return readErr
		}
	}
	if err := output.Close(); err != nil {
		_ = os.Remove(destination)
		return err
	}
	if done != total {
		_ = os.Remove(destination)
		return fmt.Errorf("copied %d bytes, expected %d", done, total)
	}
	return nil
}
