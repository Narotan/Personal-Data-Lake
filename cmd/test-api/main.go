package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ActivityEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Duration  float64   `json:"duration"`
	App       string    `json:"app"`
	Title     string    `json:"title"`
	BucketID  string    `json:"bucket_id"`
}

func main() {
	// Тестовые данные
	events := []ActivityEvent{
		{
			Timestamp: time.Now().Add(-5 * time.Minute),
			Duration:  300.0,
			App:       "test-app",
			Title:     "Test Event",
			BucketID:  "test-bucket",
		},
	}

	// Проверка 1: ActivityWatch API
	fmt.Println("=== Test 1: ActivityWatch API ===")
	resp1, err := http.Get("http://localhost:5600/api/0/buckets/")
	if err != nil {
		fmt.Printf("❌ ActivityWatch не доступен: %v\n", err)
	} else {
		fmt.Printf("✅ ActivityWatch доступен (status: %d)\n", resp1.StatusCode)
		resp1.Body.Close()
	}

	// Проверка 2: Data Lake API без ключа
	fmt.Println("\n=== Test 2: Data Lake API без ключа ===")
	data, _ := json.Marshal(events)
	resp2, err := http.Post("http://localhost:8080/api/v1/activitywatch/events", "application/json", bytes.NewReader(data))
	if err != nil {
		fmt.Printf("❌ Ошибка подключения: %v\n", err)
	} else {
		body2, _ := io.ReadAll(resp2.Body)
		fmt.Printf("Статус: %d\nОтвет: %s\n", resp2.StatusCode, string(body2))
		resp2.Body.Close()
	}

	// Проверка 3: Data Lake API с ключом
	fmt.Println("\n=== Test 3: Data Lake API с API ключом ===")
	req, _ := http.NewRequest("POST", "http://localhost:8080/api/v1/activitywatch/events", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "ac9dce6189d8d3983779004612684f9e86e5033b161deb38273c72892b6039d2")

	client := &http.Client{}
	resp3, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Ошибка подключения: %v\n", err)
	} else {
		body3, _ := io.ReadAll(resp3.Body)
		fmt.Printf("Статус: %d\nОтвет: %s\n", resp3.StatusCode, string(body3))
		resp3.Body.Close()

		if resp3.StatusCode == 200 {
			fmt.Println("✅ Данные успешно отправлены!")
		} else {
			fmt.Println("❌ Ошибка отправки данных")
		}
	}
}
