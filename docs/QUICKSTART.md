# 🚀 Быстрый старт Data Lake

> Полная установка за 10-15 минут

## 📋 Требования

- **Docker** и **Docker Compose**
- **Make** (обычно уже установлен)
- Аккаунты в сервисах (если хочешь собирать данные):
  - WakaTime 
  - Google Account (для Fit и Calendar)

---

## ⚡ Установка

### Шаг 1: Клонировать репозиторий

```bash
git clone https://github.com/Narotan/Personal-Data-Lake.git
cd data-lake
```

### Шаг 2: Создать файл конфигурации

```bash
make setup
```

Это создаст файл `.env` из шаблона `.env.example`.

---

## 🔐 Шаг 3: Настройка .env (ВАЖНО!)

Открой `.env` в редакторе:

```bash
nano .env  # или code .env, vim .env
```

### 3.1. Обязательные параметры

Сгенерируй ключи и добавь в `.env`:

```bash
# Генерация API Key
openssl rand -hex 32
# Скопируй результат в .env как API_KEY=<результат>

# Генерация User ID
uuidgen
# Скопируй результат в .env как API_USER_ID=<результат>
```

Пример `.env`:
```bash
API_KEY=a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6
API_USER_ID=123e4567-e89b-12d3-a456-426614174000
```

## 🔗 Шаг 4: Настройка OAuth (для сбора данных)

### 4.1. WakaTime (статистика кодирования)

1. **Зайди на** https://wakatime.com/apps
2. **Войди** в свой аккаунт WakaTime (или создай новый)
3. **Создай новое приложение:**
   - Нажми "Create an App"
   - App Name: `Data Lake` (любое имя)
   - Redirect URI: `http://localhost:8080/callback`
4. **Скопируй credentials:**
   - `Client ID` → в .env как `CLIENT_ID=<твой_client_id>`
   - `Client Secret` → в .env как `CLIENT_SECRET=<твой_client_secret>`
5. **В .env добавь:**
   ```bash
   REDIRECT_URI=http://localhost:8080/callback
   ```

### 4.2. Google Cloud (Fit + Calendar)

#### Шаг 1: Создать проект

1. **Зайди на** https://console.cloud.google.com/
2. **Войди** в свой Google аккаунт
3. **Создай новый проект:**
   - Нажми на выпадающий список проектов (вверху)
   - "New Project"
   - Project name: `Data Lake` (любое имя)
   - Нажми "Create"

#### Шаг 2: Включить API

1. **Перейди в** "APIs & Services" → "Library"
2. **Найди и включи:**
   - **Google Fit API** (найди через поиск, нажми "Enable")
   - **Google Calendar API** (найди через поиск, нажми "Enable")

#### Шаг 3: Настроить OAuth Consent Screen

1. **Перейди в** "APIs & Services" → "OAuth consent screen"
2. **Выбери:** "External" → "Create"
3. **Заполни:**
   - App name: `Data Lake`
   - User support email: твой email
   - Developer contact: твой email
4. **Нажми:** "Save and Continue"
5. **Scopes:** Пропусти (нажми "Save and Continue")
6. **Test users:** Добавь свой email → "Save and Continue"
7. **Нажми:** "Back to Dashboard"

#### Шаг 4: Создать OAuth Credentials

1. **Перейди в** "APIs & Services" → "Credentials"
2. **Нажми:** "Create Credentials" → "OAuth 2.0 Client ID"
3. **Выбери:**
   - Application type: "Web application"
   - Name: `Data Lake Web Client`
4. **Authorized redirect URIs - добавь ОБА:**
   - `http://localhost:8080/oauth2callback`
   - `http://localhost:8080/oauth2callback/calendar`
5. **Нажми:** "Create"
6. **Скопируй credentials:**
   - `Client ID` → в .env как `GOOGLE_CLIENT_ID=<твой_google_client_id>`
   - `Client Secret` → в .env как `GOOGLE_CLIENT_SECRET=<твой_google_client_secret>`

### 4.3. Финальный .env

Твой `.env` должен выглядеть так:

```bash
# Обязательные
API_KEY=a1b2c3d4e5f6...
API_USER_ID=123e4567-e89b...

# WakaTime
CLIENT_ID=waka_твой_client_id
CLIENT_SECRET=waka_твой_secret
REDIRECT_URI=http://localhost:8080/callback

# Google
GOOGLE_CLIENT_ID=твой_google_client_id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=GOCSPX-твой_google_secret

# База данных (по умолчанию)
DSN=postgres://postgres:postgres@postgres:5432/datalake?sslmode=disable

# Опционально
ENVIRONMENT=production
ENABLE_SCHEDULER=true
```

---

## 🚀 Шаг 5: Запустить проект

```bash
make start
```

Подожди ~30 секунд пока все сервисы запустятся.

### Проверь что всё работает:


# Проверь статус
```bash
docker ps
```


## 🔗 Шаг 6: OAuth авторизация

**После запуска проекта** нужно авторизоваться в каждом сервисе.

### Посмотри логи приложения:

```bash
make logs
```


В логах ты увидишь OAuth URL (примерно так):

```
INF WakaTime authorization URL: https://wakatime.com/oauth/authorize?client_id=...
INF Google Fit authorization URL: https://accounts.google.com/o/oauth2/auth?client_id=...
INF Google Calendar authorization URL: https://accounts.google.com/o/oauth2/auth?client_id=...
```

### Авторизуйся:

1. **Скопируй WakaTime URL** → открой в браузере → "Authorize" → вернёт на localhost
2. **Скопируй Google Fit URL** → открой в браузере → выбери аккаунт → "Allow" → вернёт на localhost
3. **Скопируй Google Calendar URL** → открой в браузере → выбери аккаунт → "Allow" → вернёт на localhost

После успешной авторизации токены сохранятся в `tokens.json` и данные начнут собираться автоматически каждые 30 минут!

---

## ✅ Проверка

### Доступные сервисы:

- **API:** http://localhost:8080
- **Grafana:** http://localhost:3000 (логин: `admin`, пароль: `admin`)
- **Prometheus:** http://localhost:9090

---

## 🖥️ Шаг 7: ActivityWatch (опционально)

**ActivityWatch** - это отдельное приложение для отслеживания активности на компьютере. Оно работает ОТДЕЛЬНО от Data Lake.

### Что нужно сделать:

#### 1. Установить ActivityWatch

**ActivityWatch НЕ включен в Docker!** Это отдельное приложение.

**Установка:**
- **Официальный сайт:** https://activitywatch.net/downloads/

Или через пакетный менеджер:
```bash
# Arch Linux
yay -S activitywatch-bin

# macOS
brew install --cask activitywatch

# Windows
# Скачай с официального сайта
```

#### 2. Запустить ActivityWatch

После установки запусти ActivityWatch:
- **Linux:** `aw-qt` или найди в меню приложений
- **macOS/Windows:** Запусти из Applications/Programs

ActivityWatch будет работать в фоне и собирать данные локально.

**По умолчанию работает на:** http://localhost:5600

#### 3. Настроить Data Lake клиент

Data Lake имеет специальный клиент который отправляет данные из ActivityWatch в базу данных.

**Собрать клиент:**
```bash
# В папке проекта
bash scripts/build_aw_client.sh
```

Это создаст бинарник `bin/aw-client`.

**Тестовый запуск:**
```bash
# Собрать данные за последние 5 минут
./bin/aw-client -minutes 5

# С кастомными параметрами
./bin/aw-client -aw-host http://localhost:5600 -server http://localhost:8080 -minutes 10
```

#### 4. Автоматический сбор (SystemD - только Linux)

Для автоматического сбора данных каждые 5 минут:

```bash
# Установить systemd сервис
bash scripts/install_service.sh

# Проверить статус
systemctl --user status aw-client@$(whoami).timer

# Посмотреть логи
journalctl --user -u aw-client@$(whoami).service -f
```

**Удалить сервис:**
```bash
bash scripts/uninstall_service.sh
```

### Как это работает:

```
ActivityWatch (localhost:5600)
      ↓ (собирает данные о приложениях)
aw-client (каждые 5 минут)
      ↓ (отправляет данные)
Data Lake API (localhost:8080)
      ↓ (сохраняет)
PostgreSQL
```

**Важно:**
- ActivityWatch должен быть ЗАПУЩЕН для сбора данных
- `aw-client` работает ОТДЕЛЬНО от основного приложения (`make start`)
- Данные отправляются через API с использованием твоего `API_KEY`
- SystemD сервис работает только на Linux (для Windows/macOS используй Task Scheduler/cron)

---

## 🎉 Готово!

Теперь у тебя запущен **Data Lake** который:
- ✅ Собирает данные из WakaTime, Google Fit, Google Calendar
- ✅ Хранит всё в PostgreSQL
- ✅ Предоставляет REST API
- ✅ Показывает метрики в Grafana

**Следующий шаг:** Изучи [API Documentation](API.md)
---

## 📚 Дополнительно

- **API Документация:** [docs/API.md](API.md)
- **Полный README:** [README.md](README.md)


