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

// WaitUntil runs a function at interval until it returns true or the context times out.
func WaitUntil(pctx context.Context, interval time.Duration, condition func(ctx context.Context) bool) bool {
	ctx, cancel := context.WithCancel(pctx)
	defer cancel()

	ok := false
	Tick(ctx, interval, func() {
		if ok = condition(ctx); ok {
			cancel()
		}
	})

	return ok
}
