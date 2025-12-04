package main

import (
	"fmt"
	"os"

	googlecalendarauth "DataLake/auth/googlecalendar"
	googlefitauth "DataLake/auth/googlefit"
	wakatimeauth "DataLake/auth/wakatime"
	"DataLake/db"
	internal_db "DataLake/internal/db"
	"DataLake/internal/logger"
	"DataLake/internal/metrics"
	"DataLake/scheduler"
	"DataLake/server"

	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
)

// printAuthorizationBanner выводит красивый баннер с ссылками на авторизацию
func printAuthorizationBanner(wakatimeURL, googleFitURL, googleCalendarURL string) {
	banner := `
╔══════════════════════════════════════════════════════════════════════════════╗
║                     🔐 PERSONAL DATA LAKE - OAUTH SETUP                      ║
╚══════════════════════════════════════════════════════════════════════════════╝

📊 To start collecting your personal data, please authorize the following services:

┌─────────────────────────────────────────────────────────────────────────────┐
│ 💜 WakaTime (Coding Activity)                                               │
├─────────────────────────────────────────────────────────────────────────────┤
│ %s
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 💚 Google Fit (Health & Activity)                                           │
├─────────────────────────────────────────────────────────────────────────────┤
│ %s
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 📅 Google Calendar (Events & Meetings)                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│ %s
└─────────────────────────────────────────────────────────────────────────────┘

💡 Tip: Copy each URL and paste it into your browser to complete authorization.
      After authorization, the system will automatically start collecting data.

╔══════════════════════════════════════════════════════════════════════════════╗
`
	fmt.Printf(banner, wakatimeURL, googleFitURL, googleCalendarURL)
}

func main() {
	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "development"
	}

	logger.Init(environment)
	metrics.Init()

	log := logger.Get()

	// Load .env only in development
	if environment == "development" {
		if err := godotenv.Load(); err != nil {
			log.Warn().Err(err).Msg("no .env file found, using environment variables")
		}
	}

	wakatimeProvider := wakatimeauth.NewProvider(
		os.Getenv("CLIENT_ID"),
		os.Getenv("CLIENT_SECRET"),
		os.Getenv("REDIRECT_URI"),
	)

	err := db.Connect()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	store := internal_db.NewStore(db.Pool)

	// Инициализация всех провайдеров OAuth
	googleFitProvider := googlefitauth.NewProviderFromEnv()
	googleCalendarProvider := googlecalendarauth.NewProviderFromEnv()

	// Красивый вывод авторизационных ссылок
	printAuthorizationBanner(
		wakatimeProvider.GetAuthURL("wakatime"),
		googleFitProvider.GetAuthURL("googlefit"),
		googleCalendarProvider.GetAuthURL("googlecalendar"),
	)

	// Создаем userID для scheduler
	userIDStr := os.Getenv("API_USER_ID")
	if userIDStr == "" {
		log.Fatal().Msg("API_USER_ID environment variable not set")
	}
	userID, err := uuid.FromString(userIDStr)
	if err != nil {
		log.Fatal().Err(err).Msg("Invalid API_USER_ID format")
	}

	// Запускаем scheduler в отдельной goroutine, если включено
	if os.Getenv("ENABLE_SCHEDULER") == "true" {
		sched := scheduler.NewScheduler(store, &log, userID)
		go sched.Start()
		log.Info().Msg("Scheduler enabled and started")
	} else {
		log.Info().Msg("Scheduler is disabled")
	}

	// Запускаем сервер
	srv := server.NewServer(store)

	if err := srv.Run(); err != nil {
		log.Fatal().Err(err).Msg("server failed")
	}
}
