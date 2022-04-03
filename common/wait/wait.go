package wait

import (
	"context"
	"time"
)

func WaitForCondition(timeout time.Duration, condition func() bool) bool {
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, timeout)
	for !condition() {
		select {
		case <-ctx.Done():
			return false
		default:
			time.Sleep(200 * time.Millisecond)
		}
	}

	return true
}
