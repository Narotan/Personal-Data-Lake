package wakatime

import (
	"strings"
	"time"
)

type DateOnly time.Time

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*d = DateOnly(t)
	return nil
}

type DaySummary struct {
	Date         DateOnly `json:"date"`
	TotalSeconds float64  `json:"total_seconds"`
	Text         string   `json:"text"`
}
type Project struct {
	Name         string  `json:"name"`
	TotalSeconds float64 `json:"total_seconds"`
	Percent      float64 `json:"percent"`
	Text         string  `json:"text"`
}

type Language struct {
	Name         string  `json:"name"`
	TotalSeconds float64 `json:"total_seconds"`
	Percent      float64 `json:"percent"`
	Text         string  `json:"text"`
}

type Editor struct {
	Name         string  `json:"name"`
	TotalSeconds float64 `json:"total_seconds"`
	Percent      float64 `json:"percent"`
	Text         string  `json:"text"`
}

type OperatingSystem struct {
	Name         string  `json:"name"`
	TotalSeconds float64 `json:"total_seconds"`
	Percent      float64 `json:"percent"`
	Text         string  `json:"text"`
}

type Dependency struct {
	Name         string  `json:"name"`
	TotalSeconds float64 `json:"total_seconds"`
	Percent      float64 `json:"percent"`
	Text         string  `json:"text"`
}

type Machine struct {
	Name         string  `json:"name"`
	TotalSeconds float64 `json:"total_seconds"`
	Percent      float64 `json:"percent"`
	Text         string  `json:"text"`
}

type Summary struct {
	Start        time.Time         `json:"start"`
	End          time.Time         `json:"end"`
	Range        string            `json:"range"`
	TotalSeconds float64           `json:"total_seconds"`
	DailyAverage float64           `json:"daily_average"`
	BestDay      DaySummary        `json:"best_day"`
	Projects     []Project         `json:"projects"`
	Languages    []Language        `json:"languages"`
	Editors      []Editor          `json:"editors"`
	OS           []OperatingSystem `json:"operating_systems"`
	Dependencies []Dependency      `json:"dependencies"`
	Machines     []Machine         `json:"machines"`
}

type SummaryResponse struct {
	Data Summary `json:"data"`
}

type DailySummary struct {
	GrandTotal struct {
		TotalSeconds float64 `json:"total_seconds"`
		Text         string  `json:"text"`
	} `json:"grand_total"`
	Projects     []Project         `json:"projects"`
	Languages    []Language        `json:"languages"`
	Editors      []Editor          `json:"editors"`
	OS           []OperatingSystem `json:"operating_systems"`
	Dependencies []Dependency      `json:"dependencies"`
	Machines     []Machine         `json:"machines"`
	Range        struct {
		Date string `json:"date"` // YYYY-MM-DD
		Text string `json:"text"`
	} `json:"range"`
}

type SummariesResponse struct {
	Data  []DailySummary `json:"data"`
	Start string         `json:"start"`
	End   string         `json:"end"`
}
