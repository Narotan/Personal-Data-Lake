package middleware

import (
	"context"
	"crypto/subtle"
	"net/http"
	"os"
)

type contextKey string

const UserIDKey contextKey = "userID"

// APIKeyAuth проверяет наличие валидного X-API-Key в заголовках
// и добавляет userID из .env в контекст запроса.
func APIKeyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := os.Getenv("API_KEY")
		if apiKey == "" {
			next.ServeHTTP(w, r)
			return
		}

		providedKey := r.Header.Get("X-API-Key")
		if subtle.ConstantTimeCompare([]byte(providedKey), []byte(apiKey)) != 1 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID := os.Getenv("API_USER_ID")
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserID извлекает userID из контекста.
func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}
