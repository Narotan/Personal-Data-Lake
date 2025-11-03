package main

import (
	"os"

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

	fullURL := wakatimeProvider.GetAuthURL("state")
	log.Info().Str("auth_url", fullURL).Msg("oauth authorization url generated")

	// Создаем userID для scheduler
	userID := uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000001"))

	// Запускаем scheduler в отдельной goroutine
	sched := scheduler.NewScheduler(store, &log, userID)
	go sched.Start()

	// Запускаем сервер
	srv := server.NewServer(store)

	if err := srv.Run(); err != nil {
		log.Fatal().Err(err).Msg("server failed")
	}
}
