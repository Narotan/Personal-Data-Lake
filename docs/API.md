# API Documentation

**Base URL:** `http://localhost:8080/api/v1`

Все endpoints требуют аутентификации через заголовок `X-API-Key`.

---

## Содержание

- [Аутентификация](#аутентификация)
- [WakaTime Endpoints](#wakatime)
- [Google Fit Endpoints](#google-fit)
- [Google Calendar Endpoints](#google-calendar)
- [ActivityWatch Endpoints](#activitywatch)
- [Обработка ошибок](#error-responses)
- [Примеры использования](#examples)

---

## Аутентификация

Все API запросы должны включать API ключ в заголовке:

```
X-API-Key: your_api_key_from_env
```

**Пример:**
```bash
curl -H "X-API-Key: your_api_key" http://localhost:8080/api/v1/wakatime/stats
```

---

## WakaTime

### Получение статистики кодирования

**GET** `/wakatime/stats`

Получение статистики кодирования из WakaTime за указанный период.

**Query Parameters:**
- `start_date` (required): Start date in format `YYYY-MM-DD`
- `end_date` (required): End date in format `YYYY-MM-DD`

**Example Request:**
```bash
curl -H "X-API-Key: your_api_key" \
  "http://localhost:8080/api/v1/wakatime/stats?start_date=2024-11-01&end_date=2024-11-07"
```

**Example Response:**
```json
{
  "total_seconds": 86400,
  "daily_average": 12342,
  "projects": [
    {
      "name": "data-lake",
      "total_seconds": 43200,
      "percent": 50.0
    }
  ],
  "languages": [
    {
      "name": "Go",
      "total_seconds": 60000,
      "percent": 69.4
    }
  ],
  "editors": [
    {
      "name": "GoLand",
      "total_seconds": 86400,
      "percent": 100.0
    }
  ],
  "operating_systems": [
    {
      "name": "Linux",
      "total_seconds": 86400,
      "percent": 100.0
    }
  ]
}
```

---

---

## Google Fit

### Получение статистики физической активности

**GET** `/googlefit/stats`

Получение данных о физической активности из Google Fit за указанный период.

**Query Parameters:**
- `start_date` (required): Start date in format `YYYY-MM-DD`
- `end_date` (required): End date in format `YYYY-MM-DD`

**Example Request:**
```bash
curl -H "X-API-Key: your_api_key" \
  "http://localhost:8080/api/v1/googlefit/stats?start_date=2024-11-01&end_date=2024-11-07"
```

**Example Response:**
```json
{
  "summary": {
    "total_steps": 70000,
    "total_calories": 14000,
    "total_distance_meters": 56000,
    "avg_steps_per_day": 10000,
    "avg_calories_per_day": 2000
  },
  "daily_data": [
    {
      "date": "2024-11-01",
      "steps": 10500,
      "calories": 2100,
      "distance_meters": 8400
    }
  ]
}
```

---

---

## Google Calendar

### Получение событий календаря

**GET** `/googlecalendar/events`

Получение событий из Google Calendar за указанный период.

**Query Parameters:**
- `start_date` (optional): Start date in format `YYYY-MM-DD` (default: 7 days ago)
- `end_date` (optional): End date in format `YYYY-MM-DD` (default: today)

**Example Request:**
```bash
curl -H "X-API-Key: your_api_key" \
  "http://localhost:8080/api/v1/googlecalendar/events?start_date=2024-11-01&end_date=2024-11-07"
```

**Example Response:**
```json
[
  {
    "id": "event123",
    "summary": "Team Meeting",
    "description": "Weekly sync",
    "start_time": "2024-11-01T10:00:00Z",
    "end_time": "2024-11-01T11:00:00Z"
  },
  {
    "id": "event124",
    "summary": "Code Review",
    "description": "",
    "start_time": "2024-11-01T14:00:00Z",
    "end_time": "2024-11-01T15:00:00Z"
  }
]
```

---

---

## ActivityWatch

### Отправка событий активности

**POST** `/activitywatch/events`

Отправка событий отслеживания активности в ActivityWatch.

**Request Body:**
```json
{
  "bucket": "aw-watcher-window",
  "events": [
    {
      "timestamp": "2024-11-01T10:00:00Z",
      "duration": 3600,
      "data": {
        "app": "GoLand",
        "title": "main.go - data-lake"
      }
    }
  ]
}
```

**Example Request:**
```bash
curl -X POST \
  -H "X-API-Key: your_api_key" \
  -H "Content-Type: application/json" \
  -d '{
    "bucket": "aw-watcher-window",
    "events": [{
      "timestamp": "2024-11-01T10:00:00Z",
      "duration": 3600,
      "data": {
        "app": "GoLand",
        "title": "main.go"
      }
    }]
  }' \
  http://localhost:8080/api/v1/activitywatch/events
```

**Response:**
```json
{
  "success": true,
  "events_processed": 1
}
```

### Получение статистики активности

**GET** `/activitywatch/stats`

Получение статистики активности из ActivityWatch за указанный период.

**Query Parameters:**
- `start_date` (optional): Start date in format `YYYY-MM-DD` (default: 7 days ago)
- `end_date` (optional): End date in format `YYYY-MM-DD` (default: today)

**Example Request:**
```bash
curl -H "X-API-Key: your_api_key" \
  "http://localhost:8080/api/v1/activitywatch/stats?start_date=2024-11-01&end_date=2024-11-07"
```

**Example Response:**
```json
{
  "total_duration_seconds": 86400,
  "by_app": [
    {
      "app": "GoLand",
      "duration_seconds": 43200,
      "percent": 50.0
    },
    {
      "app": "Chrome",
      "duration_seconds": 21600,
      "percent": 25.0
    }
  ],
  "by_category": [
    {
      "category": "Development",
      "duration_seconds": 60000,
      "percent": 69.4
    }
  ]
}
```

---

---

## Error Responses

Все endpoints могут возвращать следующие ошибки:

### 401 Unauthorized
Отсутствует или недействителен API ключ.

```json
{
  "error": "Unauthorized"
}
```

### 400 Bad Request
Некорректные параметры запроса.

```json
{
  "error": "Invalid start_date format. Use YYYY-MM-DD"
}
```

### 500 Internal Server Error
Внутренняя ошибка сервера.

```json
{
  "error": "Internal Server Error"
}
```

---

## Ограничения

**Rate Limiting:** Не реализовано. Приложение предназначено для одного пользователя.

**Хранение данных:** Все данные хранятся в PostgreSQL без ограничений по времени. При необходимости устаревшие данные можно удалить вручную.

---

## Examples

### Получение еженедельной сводки кодирования

```bash
#!/bin/bash

API_KEY="your_api_key"
START_DATE=$(date -d "7 days ago" +%Y-%m-%d)
END_DATE=$(date +%Y-%m-%d)

curl -H "X-API-Key: $API_KEY" \
  "http://localhost:8080/api/v1/wakatime/stats?start_date=$START_DATE&end_date=$END_DATE" \
  | jq '.languages'
```

### Получение данных о фитнесе за сегодня

```bash
#!/bin/bash

API_KEY="your_api_key"
TODAY=$(date +%Y-%m-%d)

curl -H "X-API-Key: $API_KEY" \
  "http://localhost:8080/api/v1/googlefit/stats?start_date=$TODAY&end_date=$TODAY" \
  | jq '.daily_data[0]'
```

### Получение событий календаря за текущую неделю

```bash
#!/bin/bash

API_KEY="your_api_key"
WEEK_START=$(date -d "monday" +%Y-%m-%d)
WEEK_END=$(date -d "sunday" +%Y-%m-%d)

curl -H "X-API-Key: $API_KEY" \
  "http://localhost:8080/api/v1/googlecalendar/events?start_date=$WEEK_START&end_date=$WEEK_END" \
  | jq 'length'
```

---

## Пример использования Python

```python
import requests
from datetime import datetime, timedelta

API_KEY = "your_api_key"
BASE_URL = "http://localhost:8080/api/v1"

headers = {
    "X-API-Key": API_KEY
}

def get_wakatime_stats(start_date: str, end_date: str):
    response = requests.get(
        f"{BASE_URL}/wakatime/stats",
        params={
            "start_date": start_date,
            "end_date": end_date
        },
        headers=headers
    )
    response.raise_for_status()
    return response.json()

# Get last 7 days
end_date = datetime.now()
start_date = end_date - timedelta(days=7)

stats = get_wakatime_stats(
    start_date.strftime("%Y-%m-%d"),
    end_date.strftime("%Y-%m-%d")
)

print(f"Total coding time: {stats['total_seconds'] / 3600:.2f} hours")
print(f"Daily average: {stats['daily_average'] / 3600:.2f} hours")
```

