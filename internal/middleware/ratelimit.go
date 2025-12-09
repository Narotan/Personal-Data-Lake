package middleware

import (
	"DataLake/internal/logger"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RateLimiter хранит лимитеры для каждого IP
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

// NewRateLimiter создает новый rate limiter
func NewRateLimiter() *RateLimiter {
	// Читаем настройки из переменных окружения
	requestsPerSecond := 10 // по умолчанию
	burstSize := 20         // по умолчанию

	if rps := os.Getenv("RATE_LIMIT_RPS"); rps != "" {
		if val, err := strconv.Atoi(rps); err == nil && val > 0 {
			requestsPerSecond = val
		}
	}

	if burst := os.Getenv("RATE_LIMIT_BURST"); burst != "" {
		if val, err := strconv.Atoi(burst); err == nil && val > 0 {
			burstSize = val
		}
	}

	rl := &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     rate.Limit(requestsPerSecond),
		burst:    burstSize,
	}

	// Очищаем старые лимитеры каждые 5 минут
	go rl.cleanupLimiters()

	return rl
}

// getLimiter возвращает лимитер для IP адреса
func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[ip] = limiter
	}

	return limiter
}

// cleanupLimiters периодически очищает неиспользуемые лимитеры
func (rl *RateLimiter) cleanupLimiters() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		// Очищаем все лимитеры (можно улучшить, отслеживая время последнего использования)
		rl.limiters = make(map[string]*rate.Limiter)
		rl.mu.Unlock()
	}
}

// Middleware для rate limiting
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получаем IP адрес клиента
		ip := r.RemoteAddr

		// Проверяем заголовок X-Forwarded-For для случая работы за прокси
		if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
			ip = forwardedFor
		}

		limiter := rl.getLimiter(ip)

		if !limiter.Allow() {
			log := logger.Get()
			log.Warn().
				Str("ip", ip).
				Str("path", r.URL.Path).
				Msg("rate limit exceeded")

			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"error": "rate limit exceeded", "message": "too many requests, please slow down"}`))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// GlobalRateLimiter - глобальный экземпляр rate limiter
var GlobalRateLimiter *RateLimiter

func init() {
	GlobalRateLimiter = NewRateLimiter()
}

// RateLimit - удобная функция для добавления rate limiting
func RateLimit(next http.Handler) http.Handler {
	return GlobalRateLimiter.Middleware(next)
}
