package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/djmarrerajr/common-lib/shared"
	"github.com/djmarrerajr/common-lib/utils"
)

var _ Runnable = new(application)

type application struct {
	name    string // the name of the application        (used for obvservability)
	version string // the version of the application     (used for obvservability)
	commit  string // the commit id of the current build (used for obvservability)

	env            utils.Environ               // map of environment values
	signalHandlers map[os.Signal]signalHandler // map of signal handlers

	AppContext *shared.ApplicationContext // application wide resources
}

func (a *application) Run() (err error) {
	sigChan := make(chan os.Signal, 1)
	defer close(sigChan)

	signal.Notify(sigChan)

	ctx, cancel := context.WithCancel(a.AppContext.RootCtx)

	grp, gCtx := errgroup.WithContext(ctx)
	if a.AppContext.Database != nil {
		err = a.AppContext.Database.Start(gCtx, grp)
		if err != nil {
			a.AppContext.Logger.WithCtx(ctx).Fatalf("unable to start application: %v", err)
		}
	}
	if a.AppContext.Server != nil {
		err = a.AppContext.Server.Start(gCtx, grp)
		if err != nil {
			a.AppContext.Logger.WithCtx(ctx).Fatalf("unable to start application: %v", err)
		}
	}

	a.AppContext.Logger.WithCtx(ctx).Infof("application startup complete")
	run := true
	for run {
		sig := <-sigChan

		switch {
		case sig == os.Interrupt:
			run = false
			// cancel()
		case sig == syscall.SIGUSR1:
			a.toggleDebug()
		default:
			if fn, exists := a.signalHandlers[sig]; exists {
				fn(a.AppContext.Logger)
			}
		}
	}

	a.Shutdown(cancel)

	<-gCtx.Done()

	return grp.Wait()
}

func (a *application) Shutdown(cancel context.CancelFunc) {
	if a.AppContext.Server != nil {
		err := a.AppContext.Server.Stop()
		if err != nil {
			a.AppContext.Logger.WithCtx(a.AppContext.RootCtx).Fatalf("unable to shutdown server: %v", err)
		}
	}

	if a.AppContext.Database != nil {
		err := a.AppContext.Database.Stop()
		if err != nil {
			a.AppContext.Logger.WithCtx(a.AppContext.RootCtx).Fatalf("unable to shutdown database: %v", err)
		}
	}

	cancel()
}

func (a *application) toggleDebug() {
	a.AppContext.Logger.ToggleDebug()
}
