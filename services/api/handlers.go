package api

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"reflect"

	"github.com/djmarrerajr/common-lib/errs"
	"github.com/djmarrerajr/common-lib/shared"
	"github.com/djmarrerajr/common-lib/utils"
	"github.com/google/uuid"
)

// this needs to be thought out and refactored...
type ErrorResponse struct {
	Code        int    `json:"code"  xml:"code"`
	Description string `json:"error" xml:"error"`
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

func (h ContextualHandler) returnErrorResponse(w http.ResponseWriter, ctype string, err error) {
	var buff []byte

	h.Logger.WithCtx(h.RootCtx).Error("error processing request", err)

	resp := ErrorResponse{
		Code:        http.StatusBadRequest,
		Description: err.Error(),
	}

	buff, err = h.marshalRequest(ctype, resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	// nolint: errcheck
	w.Write(buff)
}

func (h ContextualHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var resp any
	var buff []byte
	var err error

	ctype := r.Header.Get("Content-Type")

	w.Header().Add("Content-Type", ctype)

	data, err := h.unmarshalRequest(ctype, r.Body)
	if err != nil {
		h.returnErrorResponse(w, ctype, errs.WithType(err, errs.ErrTypeUnmarshal))
		return
	}

	if data != nil {
		err = h.ApplicationContext.Validator.Struct(data)
		if err != nil {
			h.returnErrorResponse(w, ctype, errs.WithType(err, errs.ErrTypeValidation))
			return
		}
	}

	reqID := r.Header.Get(HeaderRequestId)
	if reqID == "" {
		reqID = uuid.NewString()
	}

	reqCtx := utils.AddMapToContext(h.RootCtx, utils.FieldMap{
		"requestID":  reqID,
		"requestURL": r.URL.Path,
	})

	resp = h.CustomHandlerFunc(reqCtx, h.ApplicationContext, data)

	buff, err = h.marshalRequest(ctype, resp)
	if err != nil {
		h.returnErrorResponse(w, ctype, errs.WithType(err, errs.ErrTypeMarshal))
		return
	}

	w.WriteHeader(http.StatusOK)

	// nolint: errcheck
	w.Write(buff)
}

// defaultHealthCheckHandler will respond with a simple HTTP-200
//
//	... this should be replaced as it will cause the application to
//		always 'appear' healthy but is provided as a default
func defaultHealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
