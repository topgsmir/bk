package manage

import (
	"context"

	"os"
	"os/signal"

	"syscall"
)

// connTestContext ends on Ctrl+C, so a test can be stopped and still clean up.
func connTestContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
}
