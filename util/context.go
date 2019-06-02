package util

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

// CreateDefaultContext will create a default context which will be cancelled if the program is terminated
// The terminal state will be automatically restored if it was altered in any way
func CreateDefaultContext() (context.Context, context.CancelFunc) {
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	terminalState, _ := terminal.GetState(int(syscall.Stdin))

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-signalChan
		cancel()

		// make sure we restore the terminal state
		if terminalState != nil {
			if err := terminal.Restore(int(syscall.Stdin), terminalState); err != nil {
				panic(err)
			}
		}
	}()

	return ctx, cancel
}
