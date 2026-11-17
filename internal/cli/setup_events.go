package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/tekibo/idev/internal/manager"
)

// Setup installs the selected setup groups and components.
func Setup(args []string) int {
	if !JSON {
		return setupBody(args)
	}
	manager := manager.New()
	reader, writer, err := os.Pipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev setup: %v\n", err)
		return 1
	}
	oldOut, oldErr := os.Stdout, os.Stderr
	sink := manager.LogWriter(oldOut, "idev setup", "stdout")
	os.Stdout, os.Stderr = writer, writer
	copied := make(chan struct{})
	go func() {
		_, _ = io.Copy(sink, reader)
		close(copied)
	}()
	code := setupBody(args)
	_ = writer.Close()
	os.Stdout, os.Stderr = oldOut, oldErr
	<-copied
	sink.Flush()
	sink.Done(code)
	return code
}
