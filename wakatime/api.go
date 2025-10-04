package wakatime

import (
	"DataLake/auth"
	wakatime_db "DataLake/internal/db/wakatime"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func FetchSummaries(store *wakatime_db.Store) Summary {
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

	return summary
}

func SaveSummaries(store *wakatime_db.Store, summary Summary, userID string) error {
	ctx := context.Background()

	return store.ExecTx(ctx, func(q *wakatime_db.Queries) error {
		bestDay, err := q.CreateDay(ctx, wakatime_db.CreateDayParams{
			UserID:       pgtype.UUID{Bytes: [16]byte{}, Valid: true}, // нужно конвертировать userID → UUID
			Date:         pgtype.Date{Time: time.Time(summary.BestDay.Date), Valid: true},
			TotalSeconds: summary.BestDay.TotalSeconds,
			Text:         pgtype.Text{String: summary.BestDay.Text, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to create best day: %w", err)
		}

		summaryRow, err := q.CreateSummary(ctx, wakatime_db.CreateSummaryParams{
			UserID:       pgtype.UUID{Bytes: [16]byte{}, Valid: true},
			StartTime:    pgtype.Timestamptz{Time: summary.Start, Valid: true},
			EndTime:      pgtype.Timestamptz{Time: summary.End, Valid: true},
			Range:        pgtype.Text{String: summary.Range, Valid: true},
			TotalSeconds: pgtype.Float8{Float64: summary.TotalSeconds, Valid: true},
			DailyAverage: pgtype.Float8{Float64: summary.DailyAverage, Valid: true},
			BestDayID:    pgtype.Int4{Int32: bestDay.ID, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to create summary: %w", err)
		}

		dayID := bestDay.ID

		for _, p := range summary.Projects {
			_, err := q.CreateProject(ctx, wakatime_db.CreateProjectParams{
				DayID:        dayID,
				Name:         p.Name,
				TotalSeconds: p.TotalSeconds,
				Percent:      pgtype.Float8{Float64: p.Percent, Valid: true},
				Text:         pgtype.Text{String: p.Text, Valid: true},
			})
			if err != nil {
				return fmt.Errorf("failed to insert project %q: %w", p.Name, err)
			}
		}

		for _, lang := range summary.Languages {
			_, err := q.CreateLanguage(ctx, wakatime_db.CreateLanguageParams{
				DayID:        dayID,
				Name:         lang.Name,
				TotalSeconds: lang.TotalSeconds,
				Percent:      pgtype.Float8{Float64: lang.Percent, Valid: true},
				Text:         pgtype.Text{String: lang.Text, Valid: true},
			})
			if err != nil {
				return fmt.Errorf("failed to insert language %q: %w", lang.Name, err)
			}
		}

		for _, e := range summary.Editors {
			_, err := q.CreateEditor(ctx, wakatime_db.CreateEditorParams{
				DayID:        dayID,
				Name:         e.Name,
				TotalSeconds: e.TotalSeconds,
				Percent:      pgtype.Float8{Float64: e.Percent, Valid: true},
				Text:         pgtype.Text{String: e.Text, Valid: true},
			})
			if err != nil {
				return fmt.Errorf("failed to insert editor %q: %w", e.Name, err)
			}
		}

		for _, os := range summary.OS {
			_, err := q.CreateOS(ctx, wakatime_db.CreateOSParams{
				DayID:        dayID,
				Name:         os.Name,
				TotalSeconds: os.TotalSeconds,
				Percent:      pgtype.Float8{Float64: os.Percent, Valid: true},
				Text:         pgtype.Text{String: os.Text, Valid: true},
			})
			if err != nil {
				return fmt.Errorf("failed to insert OS %q: %w", os.Name, err)
			}
		}

		for _, d := range summary.Dependencies {
			_, err := q.CreateDependency(ctx, wakatime_db.CreateDependencyParams{
				DayID:        dayID,
				Name:         d.Name,
				TotalSeconds: d.TotalSeconds,
				Percent:      pgtype.Float8{Float64: d.Percent, Valid: true},
				Text:         pgtype.Text{String: d.Text, Valid: true},
			})
			if err != nil {
				return fmt.Errorf("failed to insert dependency %q: %w", d.Name, err)
			}
		}

		for _, m := range summary.Machines {
			_, err := q.CreateMachine(ctx, wakatime_db.CreateMachineParams{
				DayID:        dayID,
				Name:         m.Name,
				TotalSeconds: m.TotalSeconds,
				Percent:      pgtype.Float8{Float64: m.Percent, Valid: true},
				Text:         pgtype.Text{String: m.Text, Valid: true},
			})
			if err != nil {
				return fmt.Errorf("failed to insert machine %q: %w", m.Name, err)
			}
		}

		fmt.Printf("saved summary %+v (summary id: %d)\n", summaryRow, summaryRow.ID)
		return nil
	})
}
