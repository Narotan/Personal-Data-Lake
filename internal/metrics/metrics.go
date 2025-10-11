package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	WakatimeFetchTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "wakatime_fetch_total",
			Help: "Total number of WakaTime data fetches",
		},
	)

	WakatimeFetchErrors = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "wakatime_fetch_errors_total",
			Help: "Total number of WakaTime fetch errors",
		},
	)

	WakatimeFetchDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "wakatime_fetch_duration_seconds",
			Help:    "WakaTime fetch duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)

	DatabaseOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "database_operations_total",
			Help: "Total number of database operations",
		},
		[]string{"operation", "status"},
	)

	DatabaseOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "database_operation_duration_seconds",
			Help:    "Database operation duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	OAuthTokenRefreshTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "oauth_token_refresh_total",
			Help: "Total number of OAuth token refreshes",
		},
	)

	OAuthTokenRefreshErrors = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "oauth_token_refresh_errors_total",
			Help: "Total number of OAuth token refresh errors",
		},
	)
)

func Init() {
}
