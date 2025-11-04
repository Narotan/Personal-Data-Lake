package v1

import (
	handlers_api_v1 "DataLake/api/v1/handlers"
	internal_db "DataLake/internal/db"
	"DataLake/internal/middleware"
	"net/http"

	"github.com/rs/zerolog"
)

func NewRouter(store *internal_db.Store, logger *zerolog.Logger) http.Handler {
	mux := http.NewServeMux()

	wakaTimeHandler := handlers_api_v1.NewWakatimeHandler(store, logger)
	googleFitHandler := handlers_api_v1.NewGoogleFitHandler(store, logger)
	googleCalendar := handlers_api_v1.NewGoogleCalendarHandler(store, logger)
	activityWatchHandler := handlers_api_v1.NewActivityWatchHandler(store, logger)

	// wakatime endpoints
	mux.Handle("/wakatime/stats", middleware.APIKeyAuth(http.HandlerFunc(wakaTimeHandler.GetStats)))

	// googlefit endpoints
	mux.Handle("/googlefit/stats", middleware.APIKeyAuth(http.HandlerFunc(googleFitHandler.GetStats)))
	// googlecalendar endpoints
	mux.Handle("/googlecalendar/events", middleware.APIKeyAuth(http.HandlerFunc(googleCalendar.GetEvents)))

	// activitywatch endpoints
	mux.Handle("/activitywatch/events", middleware.APIKeyAuth(http.HandlerFunc(activityWatchHandler.HandleEvents)))
	mux.Handle("/activitywatch/stats", middleware.APIKeyAuth(http.HandlerFunc(activityWatchHandler.GetStats)))

	return middleware.Logging(mux)
}
