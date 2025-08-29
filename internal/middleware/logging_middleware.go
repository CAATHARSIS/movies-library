package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type responseWriteWrapper struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func (w *responseWriteWrapper) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriteWrapper) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.bytesWritten += n
	return n, err
}

func (w *responseWriteWrapper) Status() int {
	if w.statusCode == 0 {
		return http.StatusOK
	}
	return w.statusCode
}

func (w *responseWriteWrapper) BytesWritten() int {
	return w.bytesWritten
}

func NewLoggingMiddleware(log *slog.Logger) func(next http.Handler) http.Handler {
	log.With(
		slog.String("component", "middleware/logger"),
	)

	log.Info("logger middleware is enabled")
	
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				// slog.String("request_id", GetReqID(r.Context())),
			)

			ww := &responseWriteWrapper{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}		
			t := time.Now()
			defer func() {
				entry.Info(
					"request completed",
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.String("duration", time.Since(t).String()),
				)
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
