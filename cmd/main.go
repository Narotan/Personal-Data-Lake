package main

import (
	"os"

	"DataLake/auth"
	"DataLake/db"
	wakatime_db "DataLake/internal/db/wakatime"
	"DataLake/internal/logger"
	"DataLake/internal/metrics"
	"DataLake/server"
	"DataLake/wakatime"

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

	if err := godotenv.Load(); err != nil {
		log.Fatal().Err(err).Msg("error loading .env file")
	}

	cfg := auth.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		RedirectURI:  os.Getenv("REDIRECT_URI"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
	}

	err := db.Connect()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	store := wakatime_db.NewStore(db.Pool)

	fullURL := auth.BuildAuthRequest(cfg)
	log.Info().Str("auth_url", fullURL).Msg("oauth authorization url generated")

	go func() {
		userID := uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000001"))

		dailySummaries, err := wakatime.FetchSummaries()
		if err != nil {
			log.Error().Err(err).Msg("error fetching summaries")
			return
		}

		if err := wakatime.SaveSummaries(store, dailySummaries, userID); err != nil {
			log.Error().Err(err).Msg("error saving summaries")
		} else {
			log.Info().Int("count", len(dailySummaries)).Msg("successfully saved summaries")
		}
	}()

	srv := server.NewServer(cfg, store)

	if err := srv.Run(); err != nil {
		log.Fatal().Err(err).Msg("server failed")
	}
}
