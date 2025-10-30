package googlefit

import (
	"strconv"
	"time"
)

// Value represents a data value in Google Fit
type Value struct {
	MapVal []interface{} `json:"mapVal,omitempty"`
	IntVal int           `json:"intVal,omitempty"`
	FpVal  float64       `json:"fpVal,omitempty"`
}

// Point represents a data point in a dataset
type Point struct {
	StartTimeNanos     string  `json:"startTimeNanos"`
	EndTimeNanos       string  `json:"endTimeNanos"`
	OriginDataSourceId string  `json:"originDataSourceId,omitempty"`
	DataTypeName       string  `json:"dataTypeName"`
	Value              []Value `json:"value"`
}

// StartTime returns the start time as time.Time
func (p *Point) StartTime() (time.Time, error) {
	nanos, err := strconv.ParseInt(p.StartTimeNanos, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(0, nanos), nil
}

// EndTime returns the end time as time.Time
func (p *Point) EndTime() (time.Time, error) {
	nanos, err := strconv.ParseInt(p.EndTimeNanos, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(0, nanos), nil
}

// Dataset represents a collection of data points for a specific data source
type Dataset struct {
	DataSourceId string  `json:"dataSourceId"`
	Point        []Point `json:"point"`
}

// Bucket represents a time bucket containing multiple datasets
type Bucket struct {
	StartTimeMillis string    `json:"startTimeMillis"`
	EndTimeMillis   string    `json:"endTimeMillis"`
	Dataset         []Dataset `json:"dataset"`
}

// StartTime returns the start time as time.Time
func (b *Bucket) StartTime() (time.Time, error) {
	millis, err := strconv.ParseInt(b.StartTimeMillis, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(0, millis*1000000), nil
}

// EndTime returns the end time as time.Time
func (b *Bucket) EndTime() (time.Time, error) {
	millis, err := strconv.ParseInt(b.EndTimeMillis, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(0, millis*1000000), nil
}

// AggregatedDataResponse represents the response from Google Fit API for aggregated data
type AggregatedDataResponse struct {
	Bucket []Bucket `json:"bucket"`
}

// DailyStats represents aggregated daily statistics
type DailyStats struct {
	Date     time.Time
	Steps    int
	Distance float64
}

// ExtractDailyStats extracts daily statistics from the aggregated response
func (r *AggregatedDataResponse) ExtractDailyStats() ([]DailyStats, error) {
	var stats []DailyStats

	for _, bucket := range r.Bucket {
		stat := DailyStats{}

		startTime, err := bucket.StartTime()
		if err != nil {
			return nil, err
		}
		stat.Date = startTime

		for _, dataset := range bucket.Dataset {
			for _, point := range dataset.Point {
				if len(point.Value) == 0 {
					continue
				}

				switch dataset.DataSourceId {
				case "derived:com.google.step_count.delta:com.google.android.gms:aggregated":
					stat.Steps = point.Value[0].IntVal
				case "derived:com.google.distance.delta:com.google.android.gms:aggregated":
					stat.Distance = point.Value[0].FpVal
				}
			}
		}

		stats = append(stats, stat)
	}

	return stats, nil
}

// DataType constants for Google Fit
const (
	DataTypeStepCount = "com.google.step_count.delta"
	DataTypeDistance  = "com.google.distance.delta"
)
