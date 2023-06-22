package api

import (
	"net/http"

	"github.com/djmarrerajr/common-lib/errs"
)

type metricsResponseWriter struct {
	writer    http.ResponseWriter
	code      int
	ErrorType errs.ErrorType
}

func (m *metricsResponseWriter) Header() http.Header {
	return m.writer.Header()
}

func (m *metricsResponseWriter) WriteHeader(statusCode int) {
	m.code = statusCode
	m.writer.WriteHeader(statusCode)
}

func (m *metricsResponseWriter) Write(data []byte) (int, error) {
	return m.writer.Write(data)
}
