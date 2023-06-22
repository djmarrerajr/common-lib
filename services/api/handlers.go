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
	"github.com/djmarrerajr/common-lib/observability/traces"
	"github.com/djmarrerajr/common-lib/shared"
	"github.com/djmarrerajr/common-lib/utils"
)

// this needs to be thought out and refactored...
type ErrorResponse struct {
	RequestId   string         `json:"requestId" xml:"requestId"`
	Type        errs.ErrorType `json:"type"  xml:"type"`
	Code        int            `json:"code"  xml:"code"`
	Description string         `json:"error" xml:"error"`
}

// ContextualHandler wraps the underlying 'business logic' handler function
// and provides a way in which we can inject request specific and application
// wide context in to each request
type ContextualHandler struct {
	*shared.ApplicationContext

	CustomHandlerFunc func(context.Context, *shared.ApplicationContext, any) any
	any
}

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

func (h ContextualHandler) returnErrorResponse(w http.ResponseWriter, reqCtx context.Context, ctype string, err error) {
	var buff []byte

	h.Logger.WithCtx(reqCtx).Error("error processing request", err)

	reqID, _ := utils.GetFieldValueFromContext[string](reqCtx, "requestID")

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

func (h ContextualHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var resp any
	var buff []byte
	var err error

	// reqID := r.Header.Get(HeaderRequestId)
	// if reqID == "" {
	// 	reqID = uuid.NewString()
	// }

	reqCtx := utils.AddMapToContext(r.Context(), utils.GetFieldMapFromContext(h.RootCtx))
	reqCtx = utils.AddMapToContext(reqCtx, utils.FieldMap{
		// "requestID":  reqID,
		"requestURL": r.URL.Path,
	})

	span, childCtx := traces.StartChildSpan(reqCtx, runtime.FuncForPC(reflect.ValueOf(h.CustomHandlerFunc).Pointer()).Name())
	defer traces.FinishChildSpan(span)

	ctype := r.Header.Get("Content-Type")

	w.Header().Add("Content-Type", ctype)

	data, err := h.unmarshalRequest(ctype, r.Body)
	if err != nil {
		h.returnErrorResponse(w, reqCtx, ctype, errs.WithType(err, errs.ErrTypeUnmarshal))
		return
	}

	if data != nil {
		err = h.ApplicationContext.Validator.Struct(data)
		if err != nil {
			h.returnErrorResponse(w, reqCtx, ctype, errs.WithType(err, errs.ErrTypeValidation))
			return
		}
	}

	resp = h.CustomHandlerFunc(childCtx, h.ApplicationContext, data)

	buff, err = h.marshalRequest(ctype, resp)
	if err != nil {
		h.returnErrorResponse(w, reqCtx, ctype, errs.WithType(err, errs.ErrTypeMarshal))
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
//	... this should be replaced as it will cause the application to
//		always 'appear' healthy but is provided as a default
func defaultHealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
