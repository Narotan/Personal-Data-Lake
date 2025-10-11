package wakatime

import (
	"DataLake/auth"
	wakatime_db "DataLake/internal/db/wakatime"
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
	token, err := auth.LoadTokens()
	if err != nil {
		return nil, fmt.Errorf("failed to load tokens: %w", err)
	}

	// Вычисляем даты для последних 7 дней
	end := time.Now().UTC()
	start := end.AddDate(0, 0, -7)

	url := fmt.Sprintf(
		"https://wakatime.com/api/v1/users/current/summaries?start=%s&end=%s",
		start.Format("2006-01-02"),
		end.Format("2006-01-02"),
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var respData SummariesResponse
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return respData.Data, nil
}

func SaveSummaries(store *wakatime_db.Store, dailySummaries []DailySummary, userID uuid.UUID) error {
	ctx := context.Background()

	// Конвертируем UUID один раз
	var uuidBytes [16]byte
	copy(uuidBytes[:], userID.Bytes())

	// Проходим по каждому дню из массива
	for _, daySummary := range dailySummaries {
		err := store.ExecTx(ctx, func(q *wakatime_db.Queries) error {
			// Парсим дату из строки формата "2024-01-15"
			dayDate, err := time.Parse("2006-01-02", daySummary.Range.Date)
			if err != nil {
				return fmt.Errorf("failed to parse date %q: %w", daySummary.Range.Date, err)
			}

			// Попробуем найти существующий день
			existingDay, err := q.GetDayByDate(ctx, wakatime_db.GetDayByDateParams{
				UserID: pgtype.UUID{Bytes: uuidBytes, Valid: true},
				Date:   pgtype.Date{Time: dayDate, Valid: true},
			})

			var day wakatime_db.WakatimeDay
			if err == nil {
				// День существует, обновляем
				day, err = q.UpdateDay(ctx, wakatime_db.UpdateDayParams{
					ID:           existingDay.ID,
					TotalSeconds: daySummary.GrandTotal.TotalSeconds,
					Text:         pgtype.Text{String: daySummary.GrandTotal.Text, Valid: true},
				})
				if err != nil {
					return fmt.Errorf("failed to update day %s: %w", daySummary.Range.Date, err)
				}
			} else {
				// День не существует, создаем
				day, err = q.CreateDay(ctx, wakatime_db.CreateDayParams{
					UserID:       pgtype.UUID{Bytes: uuidBytes, Valid: true},
					Date:         pgtype.Date{Time: dayDate, Valid: true},
					TotalSeconds: daySummary.GrandTotal.TotalSeconds,
					Text:         pgtype.Text{String: daySummary.GrandTotal.Text, Valid: true},
				})
				if err != nil {
					return fmt.Errorf("failed to create day %s: %w", daySummary.Range.Date, err)
				}
			}

			dayID := day.ID

			// Удаляем старые связанные записи для этого дня
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

			// Сохраняем машины
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

			fmt.Printf("Saved day %s (day_id: %d, total: %s)\n",
				daySummary.Range.Date, day.ID, daySummary.GrandTotal.Text)
			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}
