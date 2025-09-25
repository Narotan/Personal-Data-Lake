package main

import (
	"fmt"
	"log"
	"os"

	"DataLake/auth"
	"DataLake/server"

	"github.com/joho/godotenv"
)

type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	fmt.Println(os.Getenv("CLIENT_ID"))
	fmt.Println(os.Getenv("REDIRECT_URI"))
	fmt.Println(os.Getenv("CLIENT_SECRET"))

	cfg := auth.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		RedirectURI:  os.Getenv("REDIRECT_URI"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
	}

	fullURL := auth.BuildAuthRequest(cfg)

	fmt.Println(fullURL)
	server.RunServer(cfg)
}
