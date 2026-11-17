package download

import (
	"context"
	"errors"
	"time"
)

const idleTimeout = 60 * time.Second

var errIdle = errors.New("download: idle timeout")

// watchIdle cancels the attempt when no bytes arrive for the idle timeout.
func watchIdle(ctx context.Context, activity <-chan struct{}) (context.Context, context.CancelFunc) {
	attempt, cancel := context.WithCancel(ctx)
	go func() {
		timer := time.NewTimer(idleTimeout)
		defer timer.Stop()
		for {
			select {
			case <-activity:
				if !timer.Stop() {
					select {
					case <-timer.C:
					default:
					}
				}
				timer.Reset(idleTimeout)
			case <-timer.C:
				cancel()
				return
			case <-attempt.Done():
				return
			}
		}
	}()
	return attempt, cancel
}
