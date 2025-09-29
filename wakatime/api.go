package wakatime

import (
	"DataLake/auth"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func FetchSummaries() Summary {
	token, err := auth.LoadTokens()
	if err != nil {
		log.Fatalf("failed to load tokens: %v", err)
	}

	url := "https://wakatime.com/api/v1/users/current/stats/last_7_days"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("unexpected status code: %d", resp.StatusCode)
	}

	var respData SummaryResponse
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		log.Fatalf("failed to decode JSON: %v", err)
	}

	summary := respData.Data

	fmt.Println(summary)
	return summary
}
