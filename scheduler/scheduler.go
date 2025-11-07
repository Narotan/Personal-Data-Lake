package scheduler

import (
	"DataLake/googlecalendar"
	"DataLake/googlefit"
	internal_db "DataLake/internal/db"
	wakatime_api "DataLake/wakatime"
	"time"

	"github.com/rs/zerolog"
	uuid "github.com/satori/go.uuid"
)

type Scheduler struct {
	store  *internal_db.Store
	logger *zerolog.Logger
	userID uuid.UUID
}

func NewScheduler(store *internal_db.Store, logger *zerolog.Logger, userID uuid.UUID) *Scheduler {
	return &Scheduler{
		store:  store,
		logger: logger,
		userID: userID,
	}
}

func (s *Scheduler) Start() {
	s.logger.Info().Msg("Scheduler started - будет собирать данные каждую минуту (режим тестирования)")

	s.runCollectors()

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.runCollectors()
	}
}

func (s *Scheduler) runCollectors() {
	s.logger.Info().Msg("Запуск сбора данных из всех API")

	// Собираем данные из WakaTime
	s.collectWakatimeData()

	// Собираем данные из Google Fit
	s.collectGoogleFitData()

	// Собираем данные из Google Calendar
	s.collectGoogleCalendarData()

	s.logger.Info().Msg("Сбор данных завершен")
}

func (s *Scheduler) collectWakatimeData() {
	s.logger.Info().Msg("Сбор данных WakaTime...")

	dailySummaries, err := wakatime_api.FetchSummaries()
	if err != nil {
		s.logger.Error().Err(err).Msg("ошибка при получении данных WakaTime")
		return
	}

	if err := wakatime_api.SaveSummaries(s.store, dailySummaries, s.userID); err != nil {
		s.logger.Error().Err(err).Msg("ошибка при сохранении данных WakaTime")
	} else {
		s.logger.Info().Int("count", len(dailySummaries)).Msg("данные WakaTime успешно сохранены")
	}
}

func (s *Scheduler) collectGoogleFitData() {
	s.logger.Info().Msg("Сбор данных Google Fit...")

	response, err := googlefit.FetchSummaries(7)
	if err != nil {
		s.logger.Error().Err(err).Msg("ошибка при получении данных Google Fit")
		return
	}

	if err := googlefit.SaveSummaries(s.store, response, s.userID); err != nil {
		s.logger.Error().Err(err).Msg("ошибка при сохранении данных Google Fit")
	} else {
		s.logger.Info().Msg("данные Google Fit успешно сохранены")
	}
}

func (s *Scheduler) collectGoogleCalendarData() {
	s.logger.Info().Msg("Сбор данных Google Calendar...")

	if err := googlecalendar.FetchAndStoreEvents(s.store, 30, s.userID); err != nil {
		s.logger.Error().Err(err).Msg("ошибка при получении данных Google Calendar")
	} else {
		s.logger.Info().Msg("данные Google Calendar успешно сохранены")
	}
}
