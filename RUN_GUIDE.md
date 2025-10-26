# 🚀 Запуск Data Lake приложения

## Быстрый старт

### 1. Запуск всех сервисов

```bash
./start-monitoring.sh
```

Этот скрипт запустит:
- 🐘 **PostgreSQL** - база данных (порт 5432)
- 🚀 **Data Lake App** - Go приложение (порт 8080)
- 📉 **Prometheus** - сбор метрик (порт 9090)
- 📈 **Grafana** - дашборды (порт 3000)
- 📝 **Loki** - логи (порт 3100)
- 📋 **Promtail** - сборщик логов

### 2. Авторизация в WakaTime

После запуска в логах увидите URL для авторизации:

```bash
# Смотрим логи
docker-compose -f docker-compose.monitoring.yml logs data-lake | grep "oauth authorization url"
```

Скопируйте URL и откройте в браузере → авторизуйтесь → токен сохранится автоматически.

### 3. Доступ к сервисам

| Сервис | URL | Логин/Пароль |
|--------|-----|--------------|
| **Приложение** | http://localhost:8080 | - |
| **Метрики** | http://localhost:8080/metrics | - |
| **Grafana** | http://localhost:3000 | admin/admin |
| **Prometheus** | http://localhost:9090 | - |
| **Loki** | http://localhost:3100 | - |

### 4. Остановка сервисов

```bash
./stop-monitoring.sh
```

## 📋 Полезные команды

### Просмотр логов

```bash
# Все сервисы
docker-compose -f docker-compose.monitoring.yml logs -f

# Только приложение
docker-compose -f docker-compose.monitoring.yml logs -f data-lake

# Последние 50 строк
docker-compose -f docker-compose.monitoring.yml logs --tail=50 data-lake
```

### Статус контейнеров

```bash
docker-compose -f docker-compose.monitoring.yml ps
```

### Перезапуск приложения

```bash
docker-compose -f docker-compose.monitoring.yml restart data-lake
```

### Пересборка после изменений в коде

```bash
./start-monitoring.sh
```

Скрипт автоматически:
1. Остановит старые контейнеры
2. Пересоберет Docker образ
3. Запустит обновленное приложение

### Подключение к БД

```bash
docker-compose -f docker-compose.monitoring.yml exec postgres psql -U postgres -d datalake
```

### Проверка таблиц

```bash
docker-compose -f docker-compose.monitoring.yml exec postgres psql -U postgres -d datalake -c '\dt'
```

## 🔧 Переменные окружения

Создайте файл `.env` в корне проекта:

```env
# WakaTime OAuth
CLIENT_ID=your_wakatime_client_id
CLIENT_SECRET=your_wakatime_client_secret
REDIRECT_URI=http://localhost:8080/callback

# Google Fit OAuth (для будущего)
GOOGLEFIT_CLIENT_ID=your_google_client_id
GOOGLEFIT_CLIENT_SECRET=your_google_client_secret
GOOGLEFIT_REDIRECT_URI=http://localhost:8080/callback/googlefit
```

## 📊 Мониторинг

### Grafana дашборды

1. Откройте http://localhost:3000
2. Логин: `admin`, Пароль: `admin`
3. Перейдите в **Dashboards** → **Data Lake Dashboard**

### Prometheus метрики

1. Откройте http://localhost:9090
2. Попробуйте запросы:
   - `wakatime_fetch_total` - количество запросов к WakaTime
   - `wakatime_fetch_errors` - количество ошибок
   - `http_requests_total` - общее количество HTTP запросов

### Loki логи

1. Откройте Grafana
2. Перейдите в **Explore**
3. Выберите **Loki** как источник данных
4. Используйте запрос: `{container_name="data-lake"}`

## 🐛 Отладка

### Приложение не стартует

```bash
# Проверьте логи
docker-compose -f docker-compose.monitoring.yml logs data-lake

# Проверьте .env файл
cat .env

# Пересоберите образ
docker-compose -f docker-compose.monitoring.yml build data-lake
./start-monitoring.sh
```

### База данных не подключается

```bash
# Проверьте что PostgreSQL запущен
docker-compose -f docker-compose.monitoring.yml ps postgres

# Проверьте логи PostgreSQL
docker-compose -f docker-compose.monitoring.yml logs postgres

# Перезапустите PostgreSQL
docker-compose -f docker-compose.monitoring.yml restart postgres
```

### Токены не сохраняются

Токены хранятся внутри контейнера в файле `/app/tokens.json`. 

Для постоянного хранения добавьте volume в `docker-compose.monitoring.yml`:

```yaml
data-lake:
  volumes:
    - ./tokens.json:/app/tokens.json
```

## 🔄 Workflow разработки

1. **Внесите изменения в код**
2. **Запустите**: `./start-monitoring.sh`
3. **Проверьте логи**: `docker-compose -f docker-compose.monitoring.yml logs -f data-lake`
4. **Проверьте метрики**: http://localhost:8080/metrics
5. **Проверьте дашборды**: http://localhost:3000

## 📚 Дополнительно

- **Архитектура OAuth2**: см. `auth/README.md`
- **Миграция**: см. `MIGRATION.md`
- **Добавление нового API**: см. `auth/README.md` → "Как добавить новый API"

---

**Готово!** Ваше приложение запущено и собирает данные из WakaTime! 🎉

