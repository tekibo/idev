package download

import (
	"fmt"
	"io"
)

const progressStep = 64 << 20

// tracked reports progress while a response body is copied.
type tracked struct {
	reader   io.Reader
	done     int64
	total    int64
	step     int64
	next     int64
	report   Progress
	activity chan<- struct{}
}

func (t *tracked) Read(buffer []byte) (int, error) {
	count, err := t.reader.Read(buffer)
	t.done += int64(count)
	if t.activity != nil && count > 0 {
		select {
		case t.activity <- struct{}{}:
		default:
		}
	}
	if t.report != nil && t.next == 0 {
		t.next = t.done + t.step
	}
	if t.report != nil && t.done >= t.next {
		t.next = t.done + t.step
		t.report(t.done, t.total)
	}
	return count, err
}

// Log returns a progress reporter that writes whole-MiB status lines.
func Log(w io.Writer, label string) Progress {
	if w == nil {
		return nil
	}
	return func(done, total int64) {
		if total > 0 {
			fmt.Fprintf(w, "%s: %s of %s\n", label, size(done), size(total))
			return
		}
		fmt.Fprintf(w, "%s: %s\n", label, size(done))
	}
}

func size(bytes int64) string {
	const unit = 1024
	if bytes < unit*unit {
		return fmt.Sprintf("%d KiB", bytes/unit)
	}
	return fmt.Sprintf("%.0f MiB", float64(bytes)/(unit*unit))
}
