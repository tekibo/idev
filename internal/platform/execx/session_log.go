package execx

import "io"

// LineWriter emits one event per output line.
type LineWriter struct {
	line *lineBuffer
}

// LogWriter returns a line-oriented event writer.
func LogWriter(out io.Writer, command, stream string) *LineWriter {
	return &LineWriter{line: newLineBuffer(stream, newEmitter(out, command))}
}

func (w *LineWriter) Write(p []byte) (int, error) {
	w.line.write(p)
	return len(p), nil
}

// Flush emits the final partial line.
func (w *LineWriter) Flush() {
	w.line.flush()
}

// Done emits the completion event.
func (w *LineWriter) Done(code int) {
	state := "succeeded"
	if code != 0 {
		state = "failed"
	}
	w.line.emit.send(Event{Type: "done", State: state, Code: code})
}
