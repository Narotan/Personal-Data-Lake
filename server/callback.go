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

		fmt.Println("Received code:", code)

		token, err := auth.ExchangeCodeForToken(cfg, code)
		if err != nil {
			http.Error(w, "Failed to exchange code for token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Printf("AccessToken: %s\nRefreshToken: %s\nTokenType: %s\nUID: %s\nExpiresAt: %s\nScope:%s\n",
			token.AccessToken, token.RefreshToken, token.TokenType, token.UID, token.ExpiresAt, token.Scope)
	}
}
