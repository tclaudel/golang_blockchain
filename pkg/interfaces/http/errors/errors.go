package errors

import (
	"encoding/json"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/tclaudel/golang_blockchain/pkg/interfaces/http/rest"
	"go.uber.org/zap"
)

type HttpCode struct {
	code int
}

const (
	ErrNilRequestIDCode = iota
	ErrMarshalingJSONCode

	ErrInvalidPublicKeyCode
	ErrInvalidSignatureCode

	ErrAppendingTransactionCode
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

	ErrInvalidPublicKey = HttpError{
		code:    HttpCode{code: ErrInvalidPublicKeyCode},
		message: "invalid public key",
	}

	ErrInvalidSignature = HttpError{
		code:    HttpCode{code: ErrInvalidSignatureCode},
		message: "invalid signature",
	}

	ErrAppendingTransaction = HttpError{
		code:    HttpCode{code: ErrAppendingTransactionCode},
		message: http.StatusText(http.StatusInternalServerError),
	}
)

func (e HttpError) Write(logger *zap.Logger, w http.ResponseWriter, httpCode int) {
	logger.Error(e.message, zap.Int("code", e.code.code), zap.String("repository", "http"))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	data, err := json.Marshal(e)
	spew.Dump(data, err)
	w.Write(data)
}

func (e HttpError) MarshalJSON() ([]byte, error) {
	return json.Marshal(rest.ErrorResponse{
		ErrCode: e.code.code,
		Message: e.message,
	})
}
