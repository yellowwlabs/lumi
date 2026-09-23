package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/charmbracelet/log"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"lumi.yellowlabs.space/pkg/metrics"
)

func LoggingMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			path := r.URL.Path
			method := r.Method

			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = generateRequestID()
				w.Header().Set("X-Request-ID", requestID)
			}
			r = r.WithContext(context.WithValue(r.Context(), "request_id", requestID))

			ww := chimiddleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)

			latency := time.Since(start)
			status := ww.Status()
			if status == 0 {
				status = http.StatusOK
			}

			metrics.HTTPRequestDuration.WithLabelValues(method, path, strconv.Itoa(status)).Observe(latency.Seconds())
			metrics.HTTPRequestTotal.WithLabelValues(method, path, strconv.Itoa(status)).Inc()
			log.Info("request completed",
				"request_id", requestID,
				"method", method,
				"path", path,
				"status", status,
				"latency", latency,
				"user_agent", r.UserAgent(),
				"ip", r.RemoteAddr,
			)
		})
	}
}

func generateRequestID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}
