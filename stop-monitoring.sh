#!/bin/bash

echo "🛑 Останавливаем все сервисы..."
docker-compose -f docker-compose.monitoring.yml down

echo ""
echo "✅ Все сервисы остановлены!"
echo ""
echo "💾 Данные сохранены в Docker volumes:"
echo "   - postgres_data"
echo "   - prometheus_data"
echo "   - loki_data"
echo "   - grafana_data"
echo ""
echo "🗑️  Для полного удаления (включая данные):"
echo "   docker-compose -f docker-compose.monitoring.yml down -v"

