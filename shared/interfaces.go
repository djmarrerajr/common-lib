package shared

import (
	"context"
	"net/http"

	"github.com/djmarrerajr/common-lib/services"
)

// nolint: unused
type Servable interface {
	services.Serviceable

	DefineRoute(string, http.HandlerFunc, ...string)
	DefineRequestHandler(string, RequestHandlerFunc, any, ...string)
}

type RequestHandlerFunc func(context.Context, *ApplicationContext, any) any
