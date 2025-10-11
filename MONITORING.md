# Data Lake - Monitoring & Logging

## Логирование с Zerolog

Все логи выводятся в формате JSON через zerolog для удобной интеграции с Loki.

В development режиме логи выводятся в читаемом формате с цветами.
В production режиме логи выводятся в JSON формате.

Установка режима через переменную окружения:
```bash
export ENVIRONMENT=production  # или development
```

## Метрики Prometheus

Приложение экспортирует метрики на endpoint `/metrics`.

### Доступные метрики:

**HTTP метрики:**
- `http_requests_total` - общее количество HTTP запросов (labels: method, endpoint, status)
- `http_request_duration_seconds` - время обработки HTTP запросов

**WakaTime метрики:**
- `wakatime_fetch_total` - общее количество запросов к WakaTime API
- `wakatime_fetch_errors_total` - количество ошибок при запросах
- `wakatime_fetch_duration_seconds` - время выполнения запросов

**Database метрики:**
- `database_operations_total` - общее количество операций с БД (labels: operation, status)
- `database_operation_duration_seconds` - время выполнения операций с БД

**OAuth метрики:**
- `oauth_token_refresh_total` - количество обновлений токенов
- `oauth_token_refresh_errors_total` - ошибки при обновлении токенов

## Запуск мониторинга

### 1. Запустить Prometheus, Loki и Grafana:
```bash
docker-compose -f docker-compose.monitoring.yml up -d
```

### 2. Запустить приложение:
```bash
go run cmd/main.go
```

### 3. Доступ к сервисам:
- **Grafana**: http://localhost:3000 (admin/admin)
- **Prometheus**: http://localhost:9090
- **Loki**: http://localhost:3100
- **Metrics endpoint**: http://localhost:8080/metrics

## Настройка Grafana

1. Откройте Grafana: http://localhost:3000
2. Войдите с учетными данными: admin/admin
3. Datasources уже настроены автоматически (Prometheus и Loki)
4. Создайте дашборды для визуализации метрик и логов

### Примеры запросов в Grafana:

**Prometheus queries:**
```promql
# Количество запросов в секунду
rate(http_requests_total[5m])

# P95 латентность HTTP запросов
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))

# Ошибки WakaTime
rate(wakatime_fetch_errors_total[5m])
```

**Loki queries:**
```logql
# Все логи приложения
{job="data-lake"}

# Только ошибки
{job="data-lake"} |= "error"

# Логи OAuth
{job="data-lake"} | json | service="data-lake" | caller=~".*oauth.*"
```

## Экспорт логов в файл

Для интеграции с Promtail можно настроить экспорт логов в файл.

Добавьте в `internal/logger/logger.go`:
```go
logFile, err := os.OpenFile("logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
output = io.MultiWriter(os.Stdout, logFile)
```

## Остановка мониторинга

```bash
docker-compose -f docker-compose.monitoring.yml down
```

Для полного удаления данных:
```bash
docker-compose -f docker-compose.monitoring.yml down -v
```
