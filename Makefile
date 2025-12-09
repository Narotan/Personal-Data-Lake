.PHONY: help setup start stop restart logs clean

.DEFAULT_GOAL := help

help: ## Показать доступные команды
	@echo 'Использование: make [команда]'
	@echo ''
	@echo 'Доступные команды:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup: ## Первоначальная настройка (создать .env)
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "✅ Создан файл .env"; \
		echo ""; \
		echo "⚠️  ВАЖНО: Отредактируй .env и заполни:"; \
		echo "   1. API_KEY (генерируй: openssl rand -hex 32)"; \
		echo "   2. API_USER_ID (генерируй: uuidgen)"; \
		echo "   3. OAuth credentials (если нужны)"; \
		echo ""; \
		echo "Потом запускай: make start"; \
	else \
		echo "⚠️  .env уже существует"; \
	fi

start: ## Запустить всё (база + приложение + мониторинг)
	@echo "🚀 Запускаю Data Lake..."
	@docker-compose up -d --build
	@echo ""
	@echo "✅ Запущено!"
	@echo ""
	@echo "Доступные сервисы:"
	@echo "  • API:        http://localhost:8080"
	@echo "  • Grafana:    http://localhost:3000 (admin/admin)"
	@echo "  • Prometheus: http://localhost:9090"
	@echo ""
	@echo "Логи: make logs"

stop: ## Остановить всё
	@echo "🛑 Останавливаю..."
	@docker-compose down
	@echo "✅ Остановлено"

restart: ## Перезапустить
	@echo "🔄 Перезапускаю..."
	@docker-compose down
	@sleep 1
	@$(MAKE) start

prune: ## Удалить все остановленные контейнеры проекта
	@echo "🧹 Очистка остановленных контейнеров..."
	@docker ps -a | grep datalake | awk '{print $$1}' | xargs -r docker rm -f
	@echo "✅ Очищено"

logs: ## Показать логи приложения
	@docker-compose logs -f app

logs-all: ## Показать логи всех сервисов
	@docker-compose logs -f

clean: ## Удалить всё (включая данные)
	@echo "⚠️  Это удалит все контейнеры и данные. Продолжить? [y/N]" && read ans && [ $${ans:-N} = y ]
	@docker-compose down -v
	@echo "✅ Всё удалено"

check-aw: ## Проверить ActivityWatch (установлен ли и запущен)
	@./scripts/check_activitywatch.sh

build-aw: ## Собрать ActivityWatch клиент
	@echo "🔨 Сборка aw-client..."
	@./scripts/build_aw_client.sh
	@echo "✅ Готово: ./bin/aw-client"

run-aw: build-aw ## Запустить ActivityWatch клиент (собрать данные за последний час)
	@echo "🚀 Запуск aw-client..."
	@echo "Проверка ActivityWatch..."
	@./scripts/check_activitywatch.sh > /dev/null 2>&1 || (echo "❌ ActivityWatch не запущен. Запустите: make check-aw" && exit 1)
	@./bin/aw-client -minutes 60 -api-key "$$(grep API_KEY .env | cut -d'=' -f2 | tr -d '\"')"
	@echo ""
	@echo "💡 Проверьте данные в БД: make check-db-aw"

check-db-aw: ## Проверить данные ActivityWatch в БД
	@echo "📊 Проверка данных в базе..."
	@docker-compose exec -T postgres psql -U postgres -d datalake -c \
		"SELECT COUNT(*) as total_events, \
		MIN(timestamp) as first_event, \
		MAX(timestamp) as last_event, \
		COUNT(DISTINCT app) as unique_apps \
		FROM activity_events;" 2>/dev/null || echo "❌ Ошибка подключения к БД"
	@echo ""
	@echo "Последние 5 событий:"
	@docker-compose exec -T postgres psql -U postgres -d datalake -c \
		"SELECT timestamp, app, LEFT(title, 50) as title, duration \
		FROM activity_events \
		ORDER BY timestamp DESC \
		LIMIT 5;" 2>/dev/null || echo "❌ Ошибка подключения к БД"

install-aw-service: ## Установить aw-client как systemd сервис
	@echo "📦 Установка systemd сервиса..."
	@sudo ./scripts/install_service.sh
	@echo "✅ Сервис установлен"
	@echo "Проверка: sudo systemctl status aw-client@$$USER.timer"

