package googlecalendar

import (
	"time"
)

// CalendarListResponse представляет список календарей пользователя
type CalendarListResponse struct {
	Kind          string     `json:"kind"`
	Etag          string     `json:"etag"`
	NextPageToken string     `json:"nextPageToken,omitempty"`
	Items         []Calendar `json:"items"`
}

// Calendar представляет календарь
type Calendar struct {
	Kind            string `json:"kind"`
	Etag            string `json:"etag"`
	ID              string `json:"id"`
	Summary         string `json:"summary"`
	Description     string `json:"description,omitempty"`
	TimeZone        string `json:"timeZone"`
	ColorID         string `json:"colorId,omitempty"`
	BackgroundColor string `json:"backgroundColor,omitempty"`
	ForegroundColor string `json:"foregroundColor,omitempty"`
	Selected        bool   `json:"selected,omitempty"`
	AccessRole      string `json:"accessRole"`
	Primary         bool   `json:"primary,omitempty"`
}

// EventsResponse представляет список событий
type EventsResponse struct {
	Kind          string  `json:"kind"`
	Etag          string  `json:"etag"`
	Summary       string  `json:"summary"`
	Updated       string  `json:"updated"`
	TimeZone      string  `json:"timeZone"`
	AccessRole    string  `json:"accessRole"`
	NextPageToken string  `json:"nextPageToken,omitempty"`
	Items         []Event `json:"items"`
}

// Event представляет событие календаря
type Event struct {
	Kind             string         `json:"kind"`
	Etag             string         `json:"etag"`
	ID               string         `json:"id"`
	Status           string         `json:"status"`
	HtmlLink         string         `json:"htmlLink"`
	Created          string         `json:"created"`
	Updated          string         `json:"updated"`
	Summary          string         `json:"summary"`
	Description      string         `json:"description,omitempty"`
	Location         string         `json:"location,omitempty"`
	ColorID          string         `json:"colorId,omitempty"`
	Creator          *Person        `json:"creator,omitempty"`
	Organizer        *Person        `json:"organizer,omitempty"`
	Start            *EventDateTime `json:"start"`
	End              *EventDateTime `json:"end"`
	Transparency     string         `json:"transparency,omitempty"`
	Visibility       string         `json:"visibility,omitempty"`
	ICalUID          string         `json:"iCalUID"`
	Sequence         int            `json:"sequence,omitempty"`
	Attendees        []Attendee     `json:"attendees,omitempty"`
	Recurrence       []string       `json:"recurrence,omitempty"`
	RecurringEventID string         `json:"recurringEventId,omitempty"`
	EventType        string         `json:"eventType,omitempty"`
}

// EventDateTime представляет дату/время события
type EventDateTime struct {
	Date     string `json:"date,omitempty"`     // Для событий на весь день (формат: yyyy-mm-dd)
	DateTime string `json:"dateTime,omitempty"` // Для событий с конкретным временем (RFC3339)
	TimeZone string `json:"timeZone,omitempty"`
}

// Person представляет участника события
type Person struct {
	Email       string `json:"email,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Self        bool   `json:"self,omitempty"`
}

// Attendee представляет участника события
type Attendee struct {
	Email          string `json:"email"`
	DisplayName    string `json:"displayName,omitempty"`
	Organizer      bool   `json:"organizer,omitempty"`
	Self           bool   `json:"self,omitempty"`
	Resource       bool   `json:"resource,omitempty"`
	Optional       bool   `json:"optional,omitempty"`
	ResponseStatus string `json:"responseStatus,omitempty"`
	Comment        string `json:"comment,omitempty"`
}

func (e *EventDateTime) ParseEventTime() (time.Time, error) {
	if e.DateTime != "" {
		return time.Parse(time.RFC3339, e.DateTime)
	}
	if e.Date != "" {
		return time.Parse("2006-01-02", e.Date)
	}
	return time.Time{}, nil
}

// GetDuration возвращает длительность события
func (e *Event) GetDuration() (time.Duration, error) {
	if e.Start == nil || e.End == nil {
		return 0, nil
	}

	startTime, err := e.Start.ParseEventTime()
	if err != nil {
		return 0, err
	}

	endTime, err := e.End.ParseEventTime()
	if err != nil {
		return 0, err
	}

	return endTime.Sub(startTime), nil
}

// IsAllDayEvent проверяет, является ли событие событием на весь день
func (e *Event) IsAllDayEvent() bool {
	return e.Start != nil && e.Start.Date != ""
}
