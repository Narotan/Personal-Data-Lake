package main

import (
	"log"
	"os"

	"DataLake/auth"
	"DataLake/db"
	wakatime_db "DataLake/internal/db/wakatime"
	"DataLake/server"
	"DataLake/wakatime"

	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	cfg := auth.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		RedirectURI:  os.Getenv("REDIRECT_URI"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
	}

	err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	store := wakatime_db.NewStore(db.Pool)

	fullURL := auth.BuildAuthRequest(cfg)
	log.Println("Auth URL:", fullURL)

	go func() {
		userID := uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000001"))

		// FetchSummaries теперь не принимает store и возвращает массив дней + ошибку
		dailySummaries, err := wakatime.FetchSummaries()
		if err != nil {
			log.Printf("error fetching summaries: %v", err)
			return
		}

		// SaveSummaries теперь принимает массив дней
		if err := wakatime.SaveSummaries(store, dailySummaries, userID); err != nil {
			log.Printf("error saving summaries: %v", err)
		} else {
			log.Printf("successfully fetched and saved %d days\n", len(dailySummaries))
		}
	}()

	logger := log.New(os.Stdout, "[DataLake] ", log.LstdFlags)
	srv := server.NewServer(cfg, store, logger)

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
