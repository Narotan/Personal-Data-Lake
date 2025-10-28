package wakatime

import (
	"DataLake/auth"
	internal_db "DataLake/internal/db"
	wakatime_db "DataLake/internal/db/wakatime"
	"DataLake/internal/logger"
	"DataLake/internal/metrics"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	uuid "github.com/satori/go.uuid"
)

// FetchSummaries получает данные по всем дням за последние 7 дней
func FetchSummaries() ([]DailySummary, error) {
	log := logger.Get()
	start := time.Now()

	metrics.WakatimeFetchTotal.Inc()

	storage := auth.NewFileTokenStorage("tokens.json")
	token, err := storage.LoadToken("wakatime")
	if err != nil {
		metrics.WakatimeFetchErrors.Inc()
		log.Error().Err(err).Msg("failed to load tokens")
		return nil, fmt.Errorf("failed to load tokens: %w", err)
	}

	end := time.Now().UTC()
	startDate := end.AddDate(0, 0, -7)

	url := fmt.Sprintf(
		"https://wakatime.com/api/v1/users/current/summaries?start=%s&end=%s",
		startDate.Format("2006-01-02"),
		end.Format("2006-01-02"),
	)

	log.Info().
		Str("start_date", startDate.Format("2006-01-02")).
		Str("end_date", end.Format("2006-01-02")).
		Msg("fetching wakatime summaries")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		metrics.WakatimeFetchErrors.Inc()
		log.Error().Err(err).Msg("failed to create request")
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		metrics.WakatimeFetchErrors.Inc()
		log.Error().Err(err).Msg("request failed")
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		metrics.WakatimeFetchErrors.Inc()
		log.Error().Int("status_code", resp.StatusCode).Msg("unexpected status code")
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var respData SummariesResponse
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		metrics.WakatimeFetchErrors.Inc()
		log.Error().Err(err).Msg("failed to decode JSON")
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	metrics.WakatimeFetchDuration.Observe(time.Since(start).Seconds())
	log.Info().
		Int("days_fetched", len(respData.Data)).
		Dur("duration", time.Since(start)).
		Msg("successfully fetched wakatime summaries")

	return respData.Data, nil
}

func SaveSummaries(store *internal_db.Store, dailySummaries []DailySummary, userID uuid.UUID) error {
	log := logger.Get()
	start := time.Now()

	ctx := context.Background()

	var uuidBytes [16]byte
	copy(uuidBytes[:], userID.Bytes())

	log.Info().
		Int("days_count", len(dailySummaries)).
		Str("user_id", userID.String()).
		Msg("saving wakatime summaries")

	for _, daySummary := range dailySummaries {
		err := store.ExecTx(ctx, func(q *wakatime_db.Queries) error {
			dayDate, err := time.Parse("2006-01-02", daySummary.Range.Date)
			if err != nil {
				log.Error().
					Err(err).
					Str("date", daySummary.Range.Date).
					Msg("failed to parse date")
				return fmt.Errorf("failed to parse date %q: %w", daySummary.Range.Date, err)
			}

			existingDay, err := q.GetDayByDate(ctx, wakatime_db.GetDayByDateParams{
				UserID: pgtype.UUID{Bytes: uuidBytes, Valid: true},
				Date:   pgtype.Date{Time: dayDate, Valid: true},
			})

			var day wakatime_db.WakatimeDay
			if err == nil {
				day, err = q.UpdateDay(ctx, wakatime_db.UpdateDayParams{
					ID:           existingDay.ID,
					TotalSeconds: daySummary.GrandTotal.TotalSeconds,
					Text:         pgtype.Text{String: daySummary.GrandTotal.Text, Valid: true},
				})
				if err != nil {
					metrics.DatabaseOperationsTotal.WithLabelValues("update_day", "error").Inc()
					log.Error().Err(err).Str("date", daySummary.Range.Date).Msg("failed to update day")
					return fmt.Errorf("failed to update day %s: %w", daySummary.Range.Date, err)
				}
				metrics.DatabaseOperationsTotal.WithLabelValues("update_day", "success").Inc()
				log.Debug().Str("date", daySummary.Range.Date).Msg("updated existing day")
			} else {
				day, err = q.CreateDay(ctx, wakatime_db.CreateDayParams{
					UserID:       pgtype.UUID{Bytes: uuidBytes, Valid: true},
					Date:         pgtype.Date{Time: dayDate, Valid: true},
					TotalSeconds: daySummary.GrandTotal.TotalSeconds,
					Text:         pgtype.Text{String: daySummary.GrandTotal.Text, Valid: true},
				})
				if err != nil {
					metrics.DatabaseOperationsTotal.WithLabelValues("create_day", "error").Inc()
					log.Error().Err(err).Str("date", daySummary.Range.Date).Msg("failed to create day")
					return fmt.Errorf("failed to create day %s: %w", daySummary.Range.Date, err)
				}
				metrics.DatabaseOperationsTotal.WithLabelValues("create_day", "success").Inc()
				log.Debug().Str("date", daySummary.Range.Date).Msg("created new day")
			}

			dayID := day.ID

			if err := q.DeleteProjectsByDay(ctx, dayID); err != nil {
				return fmt.Errorf("failed to delete old projects: %w", err)
			}
			if err := q.DeleteLanguagesByDay(ctx, dayID); err != nil {
				return fmt.Errorf("failed to delete old languages: %w", err)
			}
			if err := q.DeleteEditorsByDay(ctx, dayID); err != nil {
				return fmt.Errorf("failed to delete old editors: %w", err)
			}
			if err := q.DeleteOSByDay(ctx, dayID); err != nil {
				return fmt.Errorf("failed to delete old OS: %w", err)
			}
			if err := q.DeleteDependenciesByDay(ctx, dayID); err != nil {
				return fmt.Errorf("failed to delete old dependencies: %w", err)
			}
			if err := q.DeleteMachinesByDay(ctx, dayID); err != nil {
				return fmt.Errorf("failed to delete old machines: %w", err)
			}

			// Сохраняем проекты
			for _, p := range daySummary.Projects {
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

			// Сохраняем языки
			for _, lang := range daySummary.Languages {
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

			// Сохраняем редакторы
			for _, e := range daySummary.Editors {
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

			// Сохраняем OS
			for _, os := range daySummary.OS {
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

			// Сохраняем зависимости
			for _, d := range daySummary.Dependencies {
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

			for _, m := range daySummary.Machines {
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

			log.Debug().
				Str("date", daySummary.Range.Date).
				Int32("day_id", day.ID).
				Str("total", daySummary.GrandTotal.Text).
				Msg("saved day summary")
			return nil
		})

		if err != nil {
			metrics.DatabaseOperationsTotal.WithLabelValues("save_summary", "error").Inc()
			log.Error().Err(err).Msg("failed to save summary")
			return err
		}
		metrics.DatabaseOperationsTotal.WithLabelValues("save_summary", "success").Inc()
	}

	metrics.DatabaseOperationDuration.WithLabelValues("save_summaries").Observe(time.Since(start).Seconds())
	log.Info().
		Int("days_saved", len(dailySummaries)).
		Dur("duration", time.Since(start)).
		Msg("successfully saved all wakatime summaries")

	return nil
}
