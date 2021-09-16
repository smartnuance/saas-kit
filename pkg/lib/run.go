package lib

import (
	"context"
	"os"
	"os/signal"
)

func RunInterruptible(f func(ctx context.Context) error) error {
	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)

	// Wait for interrupt/kill signal to gracefully shutdown the server with a timeout
	quit := make(chan os.Signal, 2)
	signal.Notify(quit, os.Interrupt, os.Kill)

	defer func() {
		signal.Stop(quit)
		cancel()
	}()
	go func() {
		select {
		case <-quit:
			cancel()
		case <-ctx.Done():
		}
	}()

	return f(ctx)
}
