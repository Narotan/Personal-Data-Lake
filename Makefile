.PHONY: help setup start stop restart logs clean build-aw run-aw install-aw-service swagger

.DEFAULT_GOAL := help

help: ## Показать доступные команды
	@echo 'Использование: make [команда]'
	@echo ''
	@echo 'Доступные команды:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup: ## Первоначальная настройка проекта
	@echo "Настройка Personal Data Lake..."
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "Создан .env из .env.example"; \
		echo ""; \
		echo "ВАЖНО: Отредактируйте .env и укажите:"; \
		echo "   1. API_KEY (openssl rand -hex 32)"; \
		echo "   2. API_USER_ID (uuidgen)"; \
		echo "   3. OAuth credentials для сервисов"; \
	else \
		echo "Файл .env уже существует"; \
	fi
	@if [ ! -f web/.env ]; then \
		if [ -f .env ]; then \
			API_KEY=$$(grep "^API_KEY=" .env 2>/dev/null | cut -d'=' -f2- | tr -d '"' | tr -d "'"); \
			if [ -n "$$API_KEY" ]; then \
				echo "VITE_API_KEY=$$API_KEY" > web/.env; \
				echo "Создан web/.env с API ключом"; \
			fi; \
		fi; \
	else \
		echo "Файл web/.env уже существует"; \
	fi
	@echo ""
	@echo "Следующий шаг: make start"

start: ## Запустить все сервисы
	@echo "Запуск Personal Data Lake..."
	@docker-compose up -d --build
	@echo ""
	@echo "Все сервисы запущены"
	@echo ""
	@echo "Доступные сервисы:"
	@echo "  Веб-интерфейс:  http://localhost/"
	@echo "  API:            http://localhost/api/"
	@echo "  Grafana:        http://localhost/grafana/ (admin/admin)"
	@echo ""
	@echo "Команды:"
	@echo "  make logs        - Логи приложения"
	@echo "  make stop        - Остановить всё"
	@echo "  make restart     - Перезапустить"

stop: ## Остановить все сервисы
	@echo "Остановка..."
	@docker-compose down
	@echo "Все сервисы остановлены"

restart: ## Перезапустить все сервисы
	@echo "Перезапуск..."
	@docker-compose down
	@sleep 1
	@docker-compose up -d --build
	@echo "Сервисы перезапущены"

logs: ## Показать все логи
	@docker-compose logs -f

clean: ## Удалить все данные (БД, контейнеры, volumes)
	@echo "Это удалит ВСЕ данные. Продолжить? [y/N]" && read ans && [ $${ans:-N} = y ]
	@docker-compose down -v
	@rm -rf bin/
	@echo "Все данные удалены"

build-aw: ## Собрать ActivityWatch клиент
	@echo "Сборка aw-client..."
	@mkdir -p bin
	@go build -o bin/aw-client ./cmd/aw-client
	@echo "Готово: ./bin/aw-client"

run-aw: build-aw ## Запустить ActivityWatch клиент (сбор данных)
	@echo "Запуск aw-client..."
	@if [ ! -f .env ]; then echo "Файл .env не найден. Выполните: make setup" && exit 1; fi
	@API_KEY=$$(grep "^API_KEY=" .env 2>/dev/null | cut -d'=' -f2- | tr -d '"' | tr -d "'"); \
	if [ -z "$$API_KEY" ]; then echo "API_KEY не найден в .env" && exit 1; fi; \
	./bin/aw-client -minutes 60 -api-key "$$API_KEY"
	@echo "Данные собраны"

install-aw-service: build-aw ## Установить aw-client как системный сервис (systemd/launchd)
	@echo "Установка aw-client как системный сервис..."
	@./scripts/install_service.sh

