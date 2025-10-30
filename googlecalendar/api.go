package googlecalendar

import (
	"DataLake/auth"
	internal_db "DataLake/internal/db"
	googlecalendar_db "DataLake/internal/db/googlecalendar"
	"DataLake/internal/logger"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	uuid "github.com/satori/go.uuid"
)

const (
	calendarAPIBaseURL = "https://www.googleapis.com/calendar/v3"
)

// FetchCalendars получает список календарей пользователя
func FetchCalendars() (*CalendarListResponse, error) {
	log := logger.Get()

	storage := auth.NewFileTokenStorage("tokens.json")
	token, err := storage.LoadToken("googlecalendar")
	if err != nil {
		log.Error().Err(err).Msg("failed to load tokens")
		return nil, fmt.Errorf("failed to load tokens: %w", err)
	}

	apiURL := fmt.Sprintf("%s/users/me/calendarList", calendarAPIBaseURL)

	log.Info().Msg("fetching google calendar list")

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to create request")
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute request")
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Error().
			Int("status_code", resp.StatusCode).
			Str("response", string(body)).
			Msg("unexpected status code")
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response CalendarListResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Error().Err(err).Msg("failed to decode JSON")
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	log.Info().Int("calendar_count", len(response.Items)).Msg("successfully fetched calendars")

	return &response, nil
}

// FetchEvents получает события из календаря за указанный период
func FetchEvents(calendarID string, startTime, endTime time.Time) (*EventsResponse, error) {
	log := logger.Get()

	storage := auth.NewFileTokenStorage("tokens.json")
	token, err := storage.LoadToken("googlecalendar")
	if err != nil {
		log.Error().Err(err).Msg("failed to load tokens")
		return nil, fmt.Errorf("failed to load tokens: %w", err)
	}

	baseURL := fmt.Sprintf("%s/calendars/%s/events", calendarAPIBaseURL, url.PathEscape(calendarID))

	u, err := url.Parse(baseURL)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse URL")
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	q := u.Query()
	q.Set("timeMin", startTime.Format(time.RFC3339))
	q.Set("timeMax", endTime.Format(time.RFC3339))
	q.Set("singleEvents", "true")
	q.Set("orderBy", "startTime")
	q.Set("maxResults", "2500")
	u.RawQuery = q.Encode()

	log.Info().
		Str("calendar_id", calendarID).
		Str("start_time", startTime.Format("2006-01-02")).
		Str("end_time", endTime.Format("2006-01-02")).
		Msg("fetching google calendar events")

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to create request")
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute request")
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Error().
			Int("status_code", resp.StatusCode).
			Str("response", string(body)).
			Msg("unexpected status code")
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response EventsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Error().Err(err).Msg("failed to decode JSON")
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	log.Info().
		Str("calendar_id", calendarID).
		Int("event_count", len(response.Items)).
		Msg("successfully fetched events")

	return &response, nil
}

// FetchAndStoreEvents получает события и сохраняет их в базу данных
func FetchAndStoreEvents(store *internal_db.Store, days int, userID uuid.UUID) error {
	log := logger.Get()
	ctx := context.Background()

	// Получаем список календарей
	calendars, err := FetchCalendars()
	if err != nil {
		return fmt.Errorf("failed to fetch calendars: %w", err)
	}

	// Временной диапазон
	endTime := time.Now().UTC()
	startTime := endTime.AddDate(0, 0, -days)

	var uuidBytes [16]byte
	copy(uuidBytes[:], userID.Bytes())

	totalEvents := 0

	// Проходим по каждому календарю
	for _, calendar := range calendars.Items {
		// Получаем события для календаря
		events, err := FetchEvents(calendar.ID, startTime, endTime)
		if err != nil {
			log.Error().
				Err(err).
				Str("calendar_id", calendar.ID).
				Msg("failed to fetch events for calendar")
			continue
		}

		// Сохраняем события в БД
		for _, event := range events.Items {
			if event.Status == "cancelled" {
				continue
			}

			eventStartTime, err := event.Start.ParseEventTime()
			if err != nil {
				log.Error().Err(err).Str("event_id", event.ID).Msg("failed to parse start time")
				continue
			}

			eventEndTime, err := event.End.ParseEventTime()
			if err != nil {
				log.Error().Err(err).Str("event_id", event.ID).Msg("failed to parse end time")
				continue
			}

			duration, _ := event.GetDuration()

			err = store.ExecTxGoogleCalendar(ctx, func(q *googlecalendar_db.Queries) error {
				_, err := q.UpsertEvent(ctx, googlecalendar_db.UpsertEventParams{
					UserID:      pgtype.UUID{Bytes: uuidBytes, Valid: true},
					EventID:     event.ID,
					CalendarID:  calendar.ID,
					Summary:     pgtype.Text{String: event.Summary, Valid: event.Summary != ""},
					Description: pgtype.Text{String: event.Description, Valid: event.Description != ""},
					Location:    pgtype.Text{String: event.Location, Valid: event.Location != ""},
					StartTime:   pgtype.Timestamptz{Time: eventStartTime, Valid: true},
					EndTime:     pgtype.Timestamptz{Time: eventEndTime, Valid: true},
					IsAllDay:    pgtype.Bool{Bool: event.IsAllDayEvent(), Valid: true},
					Duration:    pgtype.Int4{Int32: int32(duration.Minutes()), Valid: true},
					Status:      pgtype.Text{String: event.Status, Valid: true},
				})

				if err != nil {
					log.Error().
						Err(err).
						Str("event_id", event.ID).
						Msg("failed to upsert event")
					return fmt.Errorf("failed to upsert event: %w", err)
				}

				return nil
			})

			if err != nil {
				log.Error().Err(err).Str("event_id", event.ID).Msg("transaction failed")
				continue
			}

			totalEvents++
		}

		log.Info().
			Str("calendar_id", calendar.ID).
			Str("calendar_name", calendar.Summary).
			Int("events_count", len(events.Items)).
			Msg("processed calendar events")
	}

	log.Info().
		Int("total_events", totalEvents).
		Int("days", days).
		Msg("successfully stored calendar events")

	return nil
}
