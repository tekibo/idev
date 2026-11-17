package execx

import (
	"bytes"
	"os"
	"strings"
	"sync"

	"github.com/tekibo/idev/internal/platform/redact"
)

type lineBuffer struct {
	stream   string
	data     []byte
	lastText string
	emit     *emitter
}

func newLineBuffer(stream string, emit *emitter) *lineBuffer {
	return &lineBuffer{stream: stream, emit: emit}
}

func (b *lineBuffer) write(p []byte) {
	b.data = append(b.data, p...)
	for {
		at := bytes.IndexByte(b.data, '\n')
		if at < 0 {
			break
		}
		line := strings.TrimRight(string(b.data[:at]), "\r")
		b.data = b.data[at+1:]
		b.lastText = ""
		if line != "" {
			b.emit.send(Event{Type: "log", Stream: b.stream, Line: redact.RedactText(line)})
		}
	}
	if len(b.data) > 65536 {
		b.data = append([]byte{}, b.data[len(b.data)-65536:]...)
	}
	id, text := promptFor(string(b.data))
	if id == "" || text == b.lastText {
		return
	}
	b.lastText = text
	b.emit.send(Event{Type: "prompt", Stream: b.stream, PromptID: id, Prompt: redact.RedactText(text)})
	b.data = b.data[:0]
}

func (b *lineBuffer) flush() {
	text := strings.TrimSpace(string(b.data))
	b.data = nil
	if text != "" {
		b.emit.send(Event{Type: "log", Stream: b.stream, Line: redact.RedactText(text)})
	}
}

func stream(r *os.File, name string, emit *emitter, wg *sync.WaitGroup) {
	defer wg.Done()
	defer r.Close()
	line := newLineBuffer(name, emit)
	buf := make([]byte, 4096)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			line.write(buf[:n])
		}
		if err != nil {
			break
		}
	}
	line.flush()
}
