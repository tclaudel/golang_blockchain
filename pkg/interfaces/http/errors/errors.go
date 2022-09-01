package errors

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type HttpCode struct {
	code int
}

const (
	ErrNilRequestIDCode = iota
	ErrMarshalingJSONCode
)

type HttpError struct {
	code    HttpCode
	message string
}

var (
	ErrNilRequestID = HttpError{
		code:    HttpCode{code: ErrNilRequestIDCode},
		message: "nil request ID",
	}
	ErrMarshalingJSON = HttpError{
		code:    HttpCode{code: ErrMarshalingJSONCode},
		message: "error marshaling JSON",
	}
)

func (e HttpError) Write(logger *zap.Logger, w http.ResponseWriter, httpCode int) {
	logger.Error(e.message, zap.Int("code", e.code.code), zap.String("repository", "http"))
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(httpCode)
	data, _ := json.Marshal(e)
	w.Write(data)
}

func (e *HttpError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}{
		Code:    e.code.code,
		Message: e.message,
	})
}
