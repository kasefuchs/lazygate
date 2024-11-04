package utils

import (
	"context"
	"time"
)

// Tick runs a function at interval until the context is cancelled.
func Tick(ctx context.Context, interval time.Duration, fn func()) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			fn()
		case <-ctx.Done():
			return
		}
	}
}
