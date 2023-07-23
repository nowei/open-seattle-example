package logger

import (
	"net/http"
	"time"

	"github.com/urfave/negroni"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func init() {
	Logger, _ = zap.NewProduction()
}

func GetLogger() *zap.Logger {
	return Logger
}

// Adapted from https://gist.github.com/ndrewnee/6187a01427b9203b9f11ca5864b8a60d for chi
// so that we can have structured logging using zap as a middleware for chi
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lrw := negroni.NewResponseWriter(w)
		next.ServeHTTP(lrw, r)
		statusCode := lrw.Status()

		fields := []zapcore.Field{
			zap.Int("status", statusCode),
			zap.String("latency", time.Since(start).String()),
			zap.String("method", r.Method),
			zap.String("uri", r.RequestURI),
			zap.String("host", r.Host),
			zap.String("remote_ip", r.RemoteAddr),
		}

		switch {
		case statusCode >= 500:
			Logger.Error("Server error", fields...)
		case statusCode >= 400:
			Logger.Warn("Client error", fields...)
		case statusCode >= 300:
			Logger.Info("Redirection", fields...)
		default:
			Logger.Info("Success", fields...)
		}
	})
}
