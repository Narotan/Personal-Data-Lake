package handlers

import (
	"DataLake/auth"
	"fmt"
	"log"
	"net/http"
)

func HandleCallback(cfg auth.Config, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Missing code", http.StatusBadRequest)
			return
		}

		token, err := auth.ExchangeCodeForToken(cfg, code)
		if err != nil {
			http.Error(w, "Failed to exchange code for token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Printf("AccessToken: %s\n", token.AccessToken)

		err = auth.SaveTokens(auth.Tokens{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			ExpiresAt:    token.ExpiresAt,
		})

		if err != nil {
			logger.Println("Failed to save token:", err)
		} else {
			logger.Println("Token saved to token.json")
		}
	}
}
