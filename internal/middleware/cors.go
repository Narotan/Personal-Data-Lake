package middleware

import (
	"net/http"
	"os"
	"strings"
)

// CORS middleware для обработки Cross-Origin Resource Sharing
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получаем список разрешенных источников из переменной окружения
		allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
		if allowedOrigins == "" {
			// По умолчанию разрешаем только localhost для разработки
			allowedOrigins = "http://localhost:8000,http://localhost"
		}

		origin := r.Header.Get("Origin")

		// Проверяем, есть ли Origin запроса в списке разрешенных
		if origin != "" {
			origins := strings.Split(allowedOrigins, ",")
			for _, allowedOrigin := range origins {
				if strings.TrimSpace(allowedOrigin) == origin {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}

		// Разрешаем методы
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Разрешаем заголовки
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-Key")

		// Разрешаем credentials (cookies, authorization headers)
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Время кэширования preflight запроса (в секундах)
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Обрабатываем preflight запросы
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
