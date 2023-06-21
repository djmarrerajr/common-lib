package app

import (
	"os"

	"github.com/djmarrerajr/common-lib/utils"
)

type signalHandler func(utils.Logger)

// WithSignalHandler allows for the definition of a handler that will be
// invoked should the application recieve the signal in question
func WithSignalHandler(sig os.Signal, fn signalHandler) Option {
	return func(a *application) {
		a.signalHandlers[sig] = fn

	}
}
