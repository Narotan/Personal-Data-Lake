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

