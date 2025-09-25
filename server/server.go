package server

import (
	"DataLake/auth"
	"log"
	"net/http"
)

func RunServer(cfg auth.Config) {
	http.HandleFunc("/callback", MakeCallbackHandler(cfg))
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
