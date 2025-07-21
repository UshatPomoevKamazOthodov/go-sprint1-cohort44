package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func Logger() func(http.Handler) http.Handler {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Сохраняем время начала
			start := time.Now()

			// Оборачиваем ResponseWriter, чтобы получить статус и размер
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// Передаём управление следующему обработчику
			next.ServeHTTP(ww, r)

			// Логируем после выполнения
			latency := time.Since(start)

			logger.Info("Request completed",
				zap.String("method", r.Method),
				zap.String("uri", r.RequestURI),
				zap.Int("status", ww.Status()),
				zap.Int64("response_size", int64(ww.BytesWritten())),
				zap.Duration("latency", latency),
			)
		})
	}
}
