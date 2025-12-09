# 📊 Personal Data Lake

> Self-hosted платформа для сбора и анализа персональной продуктивности с современным веб-интерфейсом


[![Go](https://img.shields.io/badge/Go-1.25-00ADD8?style=flat&logo=go)](https://golang.org)
[![React](https://img.shields.io/badge/React-18+-61DAFB?style=flat&logo=react)](https://reactjs.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791?style=flat&logo=postgresql)](https://postgresql.org)
[![Docker](https://img.shields.io/badge/Docker-ready-2496ED?style=flat&logo=docker)](https://docker.com)
[![Grafana](https://img.shields.io/badge/Grafana-9+-F46800?style=flat&logo=grafana)](https://grafana.com/)
[![Prometheus](https://img.shields.io/badge/Prometheus-2+-E6522C?style=flat&logo=prometheus)](https://prometheus.io/)
[![Loki](https://img.shields.io/badge/Loki-2+-00B2A9?style=flat&logo=loki)](https://grafana.com/oss/loki/)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](docs/LICENSE)


---

## 🎯 О проекте

**Personal Data Lake** - это self-hosted решение для автоматического сбора и визуализации данных о персональной продуктивности. Приложение собирает данные из различных источников, хранит их локально в PostgreSQL и предоставляет современный интерфейс для просмотра и анализа.

### ✨ Функционал

- 📊 **Интерактивный Dashboard** - графики, KPI
- 🔄 **Автоматический сбор** - данные собираются каждые 10 минут (можно легко настроить)
- 🎨 **Адаптивная визуализация** - месячная агрегация для больших периодов
- 🔐 **Полная приватность** - все данные хранятся локально
- 🚀 **Быстрый старт** - запуск одной командой через Docker
- 📈 **Мониторинг** - Prometheus + Grafana из коробки

### 📦 Источники данных

| Сервис | Что собирается |
|--------|----------------|
| 💻 **WakaTime** | Время кодирования, языки, проекты 
| 🏃 **Google Fit** | Шаги, дистанция, физическая активность 
| 📅 **Google Calendar** | События, встречи, расписание 
| 🖥️ **ActivityWatch** | Активность на компьютере, приложения
---

## ⚡ Быстрый старт

```bash
# 1. Клонировать репозиторий
git clone https://github.com/Narotan/Personal-Data-Lake.git
cd Personal-Data-Lake

# 2. Настроить окружение (создаст .env и сгенерирует ENCRYPTION_KEY)
./setup.sh

# 3. Добавить OAuth секреты в .env
nano .env

# 4. Запустить всё одной командой
make start

# 5. Запустить frontend (в новом терминале)
cd web
npm install
npm run dev
```

**📖 Подробная инструкция:** [docs/QUICKSTART.md](docs/QUICKSTART.md)

## 📸 Демонстрация

### Dashboard - Главная страница

![Dashboard Overview](docs/screenshots/dashboard-overview.png)
*Интерактивный dashboard с KPI метриками, графиками продуктивности и активности*

### Графики и визуализации

![Productivity Charts](docs/screenshots/productivity-charts.png)
*Детальная визуализация времени кодирования и топ проектов*

### Мониторинг и метрики

![Grafana Dashboard](docs/screenshots/grafana-dashboard.png)
*Мониторинг системы через Grafana с метриками Prometheus*

---

## 📊 Web Dashboard

### Возможности интерфейса

- **KPI Карточки** - мгновенная статистика за выбранный период
- **Графики продуктивности** - визуализация времени кодирования
- **График активности** - отслеживание шагов с целями
- **Топ языков** - круговая диаграмма с процентами
- **Топ проектов** - прогресс-бары по времени работы
- **Топ приложений** - анализ использования ПК
- **Timeline расписания** - события календаря в удобном виде

### Периоды анализа

- Today / Yesterday
- Last 7 / 30 Days
- Current Month
- **Last Year** (автоматическая агрегация по месяцам)
- **All Time** (автоматическая агрегация по месяцам)
- Custom Range (произвольный период)

### Умная агрегация

- **≤ 90 дней** → дневные данные для детального анализа
- **> 90 дней** → месячная агрегация для обзора трендов

---

## 🛠 Технологии

### Frontend
```
React 18 + TypeScript + Vite
TailwindCSS + Framer Motion
Recharts + date-fns
```

### Backend
```
Go 1.23 + net/http
SQLC (type-safe SQL)
zerolog + Prometheus
```

### Infrastructure
```
PostgreSQL 15
Docker + Docker Compose
Grafana + Loki + Promtail
```

---

## 🏗 Архитектура проекта

### 📁 Структура Backend

```
.
├── cmd/                    # Точки входа приложений
│   ├── main.go            # Основной сервер
│   ├── aw-client/         # Клиент ActivityWatch
│   └── test-api/          # Тестовый клиент API
│
├── api/v1/                # REST API
│   ├── router.go          # Маршрутизация API
│   ├── handlers/          # HTTP обработчики
│   │   ├── activitywatch.go
│   │   ├── googlecalendar.go
│   │   ├── googlefit.go
│   │   └── wakatime.go
│   └── models/            # API модели
│
├── server/                # HTTP сервер
│   ├── server.go          # Конфигурация сервера
│   ├── routes.go          # Общие маршруты
│   └── handlers/          # Обработчики (auth, callbacks)
│
├── auth/                  # OAuth 2.0 провайдеры
│   ├── provider.go        # Базовый провайдер
│   ├── token_manager.go   # Управление токенами
│   ├── storage.go         # Хранение токенов (tokens.json)
│   ├── googlecalendar/
│   ├── googlefit/
│   └── wakatime/
│
├── db/                    # База данных
│   ├── db.go              # Подключение к PostgreSQL
│   ├── migrations/        # Миграции схемы
│   ├── schema/            # DDL определения таблиц
│   ├── queries/           # SQL запросы для SQLC
│   └── views/             # SQL представления
│
├── internal/db/           # Генерируемый SQLC код
│   ├── store.go           # Общий интерфейс хранилища
│   ├── activitywatch/     # Queries для ActivityWatch
│   ├── googlecalendar/    # Queries для Google Calendar
│   ├── googlefit/         # Queries для Google Fit
│   └── wakatime/          # Queries для WakaTime
│
├── scheduler/             # Cron задачи
│   └── scheduler.go       # Периодический сбор данных
│
├── internal/
│   ├── logger/            # Структурированное логирование
│   ├── metrics/           # Prometheus метрики
│   └── middleware/        # HTTP middleware
│
├── wakatime/              # Интеграция WakaTime
├── googlefit/             # Интеграция Google Fit
├── googlecalendar/        # Интеграция Google Calendar
│
└── monitoring/            # Конфигурация мониторинга
    ├── prometheus.yml
    ├── loki-config.yml
    └── grafana/
```

#### 🔄 Поток данных Backend

1. **Scheduler** → запускает сбор каждые 10 минут
2. **API Clients** (wakatime, googlefit, googlecalendar) → получают данные из внешних API
3. **SQLC Stores** → сохраняют в PostgreSQL с type-safety
4. **REST API** → обслуживает запросы фронтенда
5. **Middleware** → логирование, метрики, авторизация

### 📁 Структура Frontend

```
web/
├── src/
│   ├── main.tsx           # Точка входа
│   ├── App.tsx            # Корневой компонент
│   │
│   ├── components/        # React компоненты
│   │   ├── Dashboard.tsx          # Главная страница
│   │   ├── KPICard.tsx            # KPI карточки
│   │   ├── ProductivityChart.tsx  # График времени кодирования
│   │   ├── ActivityChart.tsx      # График шагов
│   │   ├── LanguagesPieChart.tsx  # Круговая диаграмма языков
│   │   ├── ProjectsBarChart.tsx   # Топ проектов
│   │   ├── ApplicationsChart.tsx  # Топ приложений
│   │   ├── CalendarTimeline.tsx   # Timeline событий
│   │   └── DateRangePicker.tsx    # Выбор периода
│   │
│   ├── hooks/             # Custom React hooks
│   │   └── useDashboardData.ts   # Загрузка данных с API
│   │
│   ├── lib/               # Утилиты
│   │   └── api.ts                 # API клиент
│   │
│   └── index.css          # Глобальные стили
│
├── index.html             # HTML шаблон
├── vite.config.ts         # Vite конфигурация
├── tailwind.config.js     # TailwindCSS конфигурация
├── tsconfig.json          # TypeScript конфигурация
└── package.json           # Зависимости
```

#### 🎨 Архитектура Frontend

- **React 18** с TypeScript для type-safety
- **Vite** для быстрой разработки и сборки
- **TailwindCSS** для стилизации
- **Recharts** для интерактивных графиков
- **Framer Motion** для плавных анимаций
- **date-fns** для работы с датами

#### 📊 Компонентная структура

```
App.tsx
└── Dashboard.tsx
    ├── DateRangePicker       # Выбор периода
    ├── KPICard (x4)          # Метрики
    ├── ProductivityChart     # WakaTime график
    ├── ActivityChart         # Google Fit шаги
    ├── LanguagesPieChart     # Языки программирования
    ├── ProjectsBarChart      # Проекты
    ├── ApplicationsChart     # Приложения ПК
    └── CalendarTimeline      # События календаря
```

---

## 🔒 Безопасность

Personal Data Lake реализует современные практики информационной безопасности:

- 🔐 **AES-256-GCM шифрование** - все OAuth токены зашифрованы перед сохранением
- 🛡️ **Rate Limiting** - защита от DDoS и брутфорса (настраиваемый лимит)
- 🌐 **CORS** - контроль доступа с разрешенных доменов
- 👤 **Непривилегированный пользователь** - Docker контейнер работает без root
- 🔑 **Переменные окружения** - все секреты в ENV, без хардкода
- 📊 **Аудит** - структурированное логирование всех операций

### Быстрая настройка безопасности

```bash
# 1. Сгенерировать ключ шифрования (32 символа для AES-256)
openssl rand -base64 32 | cut -c1-32

# 2. Добавить в .env
ENCRYPTION_KEY=ваш-32-символьный-ключ

# 3. Настроить CORS для продакшена
ALLOWED_ORIGINS=https://yourdomain.com

# 4. Настроить Rate Limiting
RATE_LIMIT_RPS=10
RATE_LIMIT_BURST=20
```

---

## 🎮 Использование

### Основные команды

```bash
make help          # Показать все команды
make start         # Запустить backend
make logs          # Логи приложения
make stop          # Остановить всё
make restart       # Перезапустить
```

### ActivityWatch клиент

⚠️ **Важно:** ActivityWatch должен поработать 5-10 минут перед сбором данных.  
См. [ACTIVITYWATCH_NO_DATA.md](./ACTIVITYWATCH_NO_DATA.md) если получаете "No events collected".

```bash
make check-aw              # Проверить статус ActivityWatch
make build-aw              # Собрать клиент
make run-aw                # Запустить (собрать данные за час)
make check-db-aw           # Проверить данные в БД
make install-aw-service    # Установить systemd сервис (Linux)
```

### Frontend разработка

```bash
cd web
npm run dev        # Dev сервер
npm run build      # Production build
npm run preview    # Preview build
```

---

## 📄 Лицензия

Проект распространяется под лицензией MIT. Подробности в файле [LICENSE](docs/LICENSE).

---
