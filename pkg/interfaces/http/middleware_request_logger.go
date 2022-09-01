package http

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type contextRequestLoggerType string

var contextRequestLoggerKey contextRequestLoggerType = "request-logger"

type middlewareRequestLogger struct {
	logger *zap.Logger
}

func NewMiddlewareRequestLogger(logger *zap.Logger) mux.MiddlewareFunc {
	middleware := &middlewareRequestLogger{
		logger: logger,
	}

	return middleware.serveMiddleware
}

func (m middlewareRequestLogger) serveMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w, r, err := m.handle(w, r)
		if err != nil {
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func (m middlewareRequestLogger) handle(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *http.Request, error) {
	requestLoggerFields := []zap.Field{
		zap.String("method", r.Method),
		zap.String("request_uri", r.RequestURI),
		zap.String("request_id", r.Header.Get(HeaderXRequestID)),
		zap.String("remote_ip", r.RemoteAddr),
		zap.String("forwarded_for", r.Header.Get(HeaderXForwardedFor)),
		zap.String("host", r.Host),
		zap.String("user_agent", r.UserAgent()),
	}

	ctx := context.WithValue(r.Context(), contextRequestLoggerKey, m.logger.With(requestLoggerFields...))
	r = r.Clone(ctx)

	return w, r, nil
}

func GetRequestLogger(r *http.Request) *zap.Logger {
	logger, ok := r.Context().Value(contextRequestLoggerKey).(*zap.Logger)
	if !ok {
		panic("no logger set")
	}

	return logger
}
