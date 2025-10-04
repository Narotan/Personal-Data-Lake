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
		userID := "00000000-0000-0000-0000-000000000001"

		summary := wakatime.FetchSummaries(store)
		if err := wakatime.SaveSummaries(store, summary, userID); err != nil {
			log.Printf("error saving summaries: %v", err)
		} else {
			log.Printf("successfully fetched and saved summary: %+v\n", summary)
		}
	}()

	logger := log.New(os.Stdout, "[DataLake] ", log.LstdFlags)
	srv := server.NewServer(cfg, store, logger)

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
