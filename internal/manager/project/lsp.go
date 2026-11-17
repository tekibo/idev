package project

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"

	"github.com/tekibo/idev/internal/platform/config"
	"github.com/tekibo/idev/internal/platform/execx"
)

// eofReader reports when the editor closes the LSP stdin.
type eofReader struct {
	reader io.Reader
	once   sync.Once
	done   chan struct{}
}

func newEOFReader(reader io.Reader) *eofReader {
	return &eofReader{reader: reader, done: make(chan struct{})}
}

func (r *eofReader) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)
	if err != nil {
		r.once.Do(func() { close(r.done) })
	}
	return n, err
}

// RunLSP runs the private SourceKit-LSP over stdio.
func (s *Service) RunLSP(args []string, stdin io.Reader, stdout, stderr io.Writer) (int, error) {
	path := s.Config.Path("bin", "sourcekit-lsp")
	if !config.IsExec(path) {
		return 127, fmt.Errorf("private sourcekit-lsp not installed; run idev setup --swift")
	}
	stream, dir := lspRoot(stdin)
	input := newEOFReader(stream)
	cmd := execx.NewCommand(execx.Command{
		Name:   path,
		Args:   args,
		Env:    s.Config.SwiftBuildEnv(),
		Dir:    dir,
		Stdin:  input,
		Stdout: stdout,
		Stderr: stderr,
	})
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		return 1, err
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	defer signal.Stop(stop)
	finished := make(chan struct{})
	var stopping atomic.Bool
	go func() {
		select {
		case <-input.done:
		case <-stop:
		case <-finished:
			return
		}
		stopping.Store(true)
		killChildren(cmd.Process.Pid)
		_ = cmd.Process.Signal(syscall.SIGTERM)
	}()
	err := cmd.Wait()
	close(finished)
	if cmd.Process != nil {
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)
	}
	if stopping.Load() {
		return 0, nil
	}
	return s.ExitCode(err), err
}

func killChildren(pid int) {
	for _, child := range childPIDs(pid) {
		killChildren(child)
		_ = syscall.Kill(child, syscall.SIGTERM)
	}
}

func childPIDs(pid int) []int {
	paths, err := filepath.Glob(fmt.Sprintf("/proc/%d/task/*/children", pid))
	if err != nil {
		return nil
	}
	seen := map[int]bool{}
	pids := []int{}
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		for _, field := range strings.Fields(string(data)) {
			child, err := strconv.Atoi(field)
			if err != nil || seen[child] {
				continue
			}
			seen[child] = true
			pids = append(pids, child)
		}
	}
	return pids
}
