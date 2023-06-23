package api

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"reflect"
	"runtime"

	"github.com/djmarrerajr/common-lib/errs"
	"github.com/djmarrerajr/common-lib/observability/tracing"
	"github.com/djmarrerajr/common-lib/shared"
	"github.com/djmarrerajr/common-lib/utils"
)

// TODO: this needs to be validated and refactored...
type ErrorResponse struct {
	RequestId   string         `json:"requestId" xml:"requestId"`
	Type        errs.ErrorType `json:"type"  xml:"type"`
	Code        int            `json:"code"  xml:"code"`
	Description string         `json:"error" xml:"error"`
}

// ContextualHandler wraps the underlying domain handler function and
// provides a way in which we can inject request specific context
type ContextualHandler struct {
	*shared.ApplicationContext

	CustomHandlerFunc shared.RequestHandlerFunc
	any
}

// ServeHTTP is central to the operation of our API, it will:
//
//	... retrieve the content-type header value
//	...	create a span that can be used to trace the request
//	... transform the incoming body in to a domain object
//	... optionally validate the domain object
//	... invoke the domain logic
//	... transform the domain response
//	... return the response to the API client
func (h ContextualHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error

	// get the incoming content-type as it drives our behavior...
	ctype := r.Header.Get(HeaderContentType)

	// set our outgoing content-type to match...
	w.Header().Add(HeaderContentType, ctype)

	// combine our app-wide content with the context of this request...
	reqCtx := utils.AddMapToContext(r.Context(), utils.GetFieldMapFromContext(h.RootCtx))

	// start our outer (parent) span for this request...
	span, spanCtx := tracing.StartChildSpan(reqCtx, runtime.FuncForPC(reflect.ValueOf(h.CustomHandlerFunc).Pointer()).Name())
	defer tracing.FinishChildSpan(span)

	// turn our request body in to something more useful...
	data, err := h.unmarshalRequest(ctype, r.Body)
	if err != nil {
		h.returnErrorResponse(w, reqCtx, ctype, errs.WithType(err, errs.ErrTypeUnmarshal))
		return
	}

	// if we *have* any data validate it - if we *have* a validator
	if data != nil {
		if h.ApplicationContext.Validator != nil {
			err = h.ApplicationContext.Validator.Struct(data)
			if err != nil {
				h.returnErrorResponse(w, reqCtx, ctype, errs.WithType(err, errs.ErrTypeValidation))
				return
			}
		}
	}

	var resp any
	var status int
	var buff []byte

	// invoke our business logic/handler...
	resp, status = h.CustomHandlerFunc(spanCtx, h.ApplicationContext, data)

	// turn our response in to something more interesting...
	buff, err = h.marshalRequest(ctype, resp)
	if err != nil {
		h.returnErrorResponse(w, reqCtx, ctype, errs.WithType(err, errs.ErrTypeMarshal))
		return
	}

	// ensure we return a valid status...
	if status == 0 {
		status = http.StatusOK
	}

	w.WriteHeader(status)

	// send our response...
	_, err = w.Write(buff)
	if err != nil {
		h.Logger.WithCtx(reqCtx).Error("error writing response", err)
	}
}

// unmarshalRequest will, based on the incoming content-type, transform the incoming
// request body in to a pointer object that can be cast to the correct type by the receiver
func (h ContextualHandler) unmarshalRequest(ctype string, body io.ReadCloser) (any, error) {
	if h.any != nil {
		data := reflect.New(reflect.TypeOf(h.any)).Interface()

		switch ctype {
		case "application/json":
			err := json.NewDecoder(body).Decode(data)
			if err != nil {
				return nil, err
			}
		case "text/xml", "application/xml":
			err := xml.NewDecoder(body).Decode(data)
			if err != nil {
				return nil, err
			}
		}
		return data, nil
	}

	return nil, nil
}

// marshalRequest will, based on the incoming content-type, transform the outgoing
// domain object to an http response
func (h ContextualHandler) marshalRequest(ctype string, body any) ([]byte, error) {
	var buff []byte
	var err error

	if body != nil {
		switch ctype {
		case "application/json":
			buff, err = json.Marshal(body)
			if err != nil {
				return nil, err
			}
		case "text/xml", "application/xml":
			buff, err = xml.Marshal(body)
			if err != nil {
				return nil, err
			}
		default:
			switch body := body.(type) {
			case string:
				buff = []byte(body)
			case []byte:
				buff = body
			default:
				if reflect.TypeOf(body).Kind() == reflect.Struct {
					buff, err = json.Marshal(body)
					if err != nil {
						return nil, err
					}
				} else {
					return nil, errs.New(errs.ErrTypeMarshal, "unsupport type")
				}
			}
		}

		return buff, nil
	}

	return nil, nil
}

// returnErrorResponse will, as the name states, return a standardized error to the API caller
func (h ContextualHandler) returnErrorResponse(w http.ResponseWriter, reqCtx context.Context, ctype string, err error) {
	var buff []byte

	h.Logger.WithCtx(reqCtx).Error("error processing request", err)

	reqID, _ := utils.GetFieldValueFromContext[string](reqCtx, shared.RequestIdContextKey)

	w.(*metricsResponseWriter).errorType = errs.GetType(err)

	resp := ErrorResponse{
		RequestId:   reqID,
		Type:        w.(*metricsResponseWriter).errorType,
		Code:        http.StatusBadRequest,
		Description: err.Error(),
	}

	buff, err = h.marshalRequest(ctype, resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(buff)
	if err != nil {
		h.Logger.WithCtx(reqCtx).Error("error writing response", err)
	}
}

// defaultHealthCheckHandler will respond with a simple HTTP-200
//
//	... this should be probably replaced with something with more
//		insight in to the applications health as it will cause the
//		application to always 'appear' healthy
func defaultHealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
