// Session events are the machine-readable command stream.
package execx

import (
	"encoding/json"
	"io"
	"sync"
	"time"
)

// Event is one machine-readable command event.
type Event struct {
	Time     time.Time `json:"time"`
	Command  string    `json:"command"`
	Type     string    `json:"type"`
	Stream   string    `json:"stream,omitempty"`
	Line     string    `json:"line,omitempty"`
	Prompt   string    `json:"prompt,omitempty"`
	PromptID string    `json:"prompt_id,omitempty"`
	State    string    `json:"state,omitempty"`
	Error    string    `json:"error,omitempty"`
	Code     int       `json:"code,omitempty"`
}

type emitter struct {
	mu      sync.Mutex
	enc     *json.Encoder
	command string
}

func newEmitter(out io.Writer, command string) *emitter {
	return &emitter{enc: json.NewEncoder(out), command: command}
}

func (e *emitter) send(event Event) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if event.Time.IsZero() {
		event.Time = time.Now().UTC()
	}
	if event.Command == "" {
		event.Command = e.command
	}
	_ = e.enc.Encode(event)
}
