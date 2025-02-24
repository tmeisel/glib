package exec

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	ctxPkg "github.com/tmeisel/glib/ctx"
)

func SigtermHandler(ctx context.Context, cancelFn context.CancelFunc) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT)

	sig := <-sigs

	if logger := ctxPkg.GetLogger(ctx); logger != nil {
		logger.Infof(ctx, "received %s, shutting down gracefully", sig.String())
	}

	cancelFn()
}

func Deferred(fn func() error) {
	if err := fn(); err != nil {
		log.Printf("deferred fn: %v", err)
	}
}
