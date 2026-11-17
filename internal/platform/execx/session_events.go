package execx

import (
	"context"
	"io"
	"os"
	"sync"
)

func runEvents(ctx context.Context, command Command, in io.Reader, out io.Writer) error {
	preview := Command{Name: command.Name, Args: command.Args}.Preview()
	emit := newEmitter(out, preview)
	inputR, inputW, err := os.Pipe()
	if err != nil {
		return err
	}
	outputR, outputW, err := os.Pipe()
	if err != nil {
		closeFiles(inputR, inputW)
		return err
	}
	errorR, errorW, err := os.Pipe()
	if err != nil {
		closeFiles(inputR, inputW, outputR, outputW)
		return err
	}
	cmd := NewCommandContext(ctx, Command{
		Name: command.Name, Args: command.Args, Env: command.Env, Dir: command.Dir,
		Stdin: inputR, Stdout: outputW, Stderr: errorW,
	})
	setGroup(cmd)
	emit.send(Event{Type: "start"})
	if err := cmd.Start(); err != nil {
		closeFiles(inputR, inputW, outputR, outputW, errorR, errorW)
		return err
	}
	closeFiles(inputR, outputW, errorW)
	go copyInput(in, inputW)
	var wg sync.WaitGroup
	wg.Add(2)
	go stream(outputR, "stdout", emit, &wg)
	go stream(errorR, "stderr", emit, &wg)
	err = waitGroup(ctx, cmd)
	wg.Wait()
	emit.send(doneEvent(err))
	return err
}

func copyInput(in io.Reader, w *os.File) {
	_, _ = io.Copy(w, in)
	_ = w.Close()
}
