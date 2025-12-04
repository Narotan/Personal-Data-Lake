package models_api_v1

// ProjectStat представляет статистику по одному проекту.
type ProjectStat struct {
	Name         string  `json:"name"`
	TotalSeconds float64 `json:"total_seconds"`
}

// LanguageStat представляет статистику по одному языку.
type LanguageStat struct {
	Name         string  `json:"name"`
	TotalSeconds float64 `json:"total_seconds"`
	Percent      float64 `json:"percent"`
}

// DailyStat представляет полную дневную статистику.
type DailyStat struct {
	Date         string         `json:"date"`
	TotalSeconds float64        `json:"total_seconds"`
	Text         string         `json:"text"`
	Projects     []ProjectStat  `json:"projects"`
	Languages    []LanguageStat `json:"languages"`
}

// AggregatedLanguageStat для агрегированной статистики языков
type AggregatedLanguageStat struct {
	Name         string  `json:"name"`
	TotalSeconds float64 `json:"total_seconds"`
	Percent      float64 `json:"percent"`
}

// AggregatedProjectStat для агрегированной статистики проектов
type AggregatedProjectStat struct {
	Name         string  `json:"name"`
	TotalSeconds float64 `json:"total_seconds"`
}

type DailyFitStat struct {
	Date     string  `json:"date"`
	Steps    int     `json:"steps"`
	Distance float64 `json:"distance"` // Калории убраны
}

type CalendarEvent struct {
	ID          string `json:"id"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
}
