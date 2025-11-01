#!/bin/bash

echo "🛑 Останавливаем старые контейнеры..."
docker-compose -f docker-compose.monitoring.yml down

echo "🗑️  Удаляем старый образ..."
docker rmi data_lake-data_lake 2>/dev/null || true

echo "🏗️  Собираем Docker образ..."
docker-compose -f docker-compose.monitoring.yml build data_lake

echo "🚀 Запускаем PostgreSQL..."
docker-compose -f docker-compose.monitoring.yml up -d postgres

echo "⏳ Ждем запуска PostgreSQL (8 сек)..."
sleep 8

echo "🚀 Запускаем все остальные сервисы..."
docker-compose -f docker-compose.monitoring.yml up -d

echo "⏳ Ждем запуска всех сервисов..."
sleep 8

echo ""
echo "📊 Статус контейнеров:"
docker-compose -f docker-compose.monitoring.yml ps

echo ""
echo "✅ Готово! Сервисы доступны:"
echo "   🐘  PostgreSQL:  localhost:5432"
echo "   🚀  App:         http://localhost:8080"
echo "   📊 Metrics:     http://localhost:8080/metrics"
echo "   📈 Grafana:     http://localhost:3000 (admin/admin)"
echo "   📉 Prometheus:  http://localhost:9090"
echo "   📝 Loki:        http://localhost:3100"
echo ""
echo "📋 Просмотр логов приложения:"
echo "   docker-compose -f docker-compose.monitoring.yml logs -f data_lake"
echo ""
echo "🔍 Проверка подключения к БД:"
echo "   docker-compose -f docker-compose.monitoring.yml exec postgres psql -U postgres -d datalake -c '\\dt'"
