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
	// Загружаем .env
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Создаём конфиг для OAuth
	cfg := auth.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		RedirectURI:  os.Getenv("REDIRECT_URI"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
	}

	// Подключение к базе данных
	if err := db.Connect(); err != nil {
		log.Fatal(err)
	}
	defer db.Pool.Close()
	queries := wakatime_db.New(db.Pool) // используем глобальную переменную db.Pool

	// Генерируем ссылку для авторизации
	fullURL := auth.BuildAuthRequest(cfg)
	log.Println("Auth URL:", fullURL)

	go func() {
		summary := wakatime.FetchSummaries()
		log.Printf("Fetched summary: %+v\n", summary)
	}()

	logger := log.New(os.Stdout, "[DataLake] ", log.LstdFlags)

	srv := server.NewServer(cfg, queries, logger)

	// заупуск сервер
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
