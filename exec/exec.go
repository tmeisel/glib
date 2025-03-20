package exec

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	ctxPkg "github.com/tmeisel/glib/ctx"
)

// SigHandler waits for syscall.SIG* and calls cancelFn, if received. If finally is
// a function, it will be called last before it returns. Any error
func SigHandler(ctx context.Context, cancelFn context.CancelFunc, finally func() error) error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT)

	sig := <-sigs

	if logger := ctxPkg.GetLogger(ctx); logger != nil {
		logger.Infof(ctx, "received %s, shutting down gracefully", sig.String())
	}

	cancelFn()

	if finally == nil {
		return nil
	}

	return finally()
}

func Deferred(fn func() error) {
	if err := fn(); err != nil {
		log.Printf("deferred fn: %v", err)
	}
}
