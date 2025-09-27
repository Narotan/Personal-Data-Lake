package server

import (
	"DataLake/auth"
	"fmt"
	"net/http"
)

func MakeCallbackHandler(cfg auth.Config) http.HandlerFunc {
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
			fmt.Println("Failed to save token:", err)
		} else {
			fmt.Println("Token saved to token.json")
		}

	}
}
