package app

import (
	"github.com/djmarrerajr/common-lib/utils"
)

func WithLoggerAtLevel(level string) Option {
	return func(a *application) {
		a.AppContext.Logger = utils.NewLogger(level).WithCtx(a.AppContext.RootCtx)
	}
}
