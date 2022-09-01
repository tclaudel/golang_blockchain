package http

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type middlewareRequestID struct {
	logger *zap.Logger
}

func NewMiddlewareRequestID(logger *zap.Logger) mux.MiddlewareFunc {
	middleware := &middlewareRequestID{
		logger: logger,
	}

	return middleware.serveMiddleware
}

func (m middlewareRequestID) serveMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := m.handle(w, r)
		if err != nil {
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func (m middlewareRequestID) handle(w http.ResponseWriter, r *http.Request) error {
	requestID := r.Header.Get(HeaderXRequestID)
	if requestID == "" {
		requestID = uuid.New().String()
		r.Header.Set(HeaderXRequestID, requestID)
	}

	return nil
}

func GetRequestID(r *http.Request) (string, error) {
	requestID := r.Header.Get(HeaderXRequestID)
	if requestID == "" {
		return "", fmt.Errorf("no %s header provided", HeaderXRequestID)
	}

	return requestID, nil
}
