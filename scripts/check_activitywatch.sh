#!/bin/bash
# Проверка и установка ActivityWatch

echo "🔍 Проверка ActivityWatch..."
echo ""

# Проверка запущен ли ActivityWatch
if curl -s http://localhost:5600/api/0/buckets/ > /dev/null 2>&1; then
    echo "✅ ActivityWatch запущен и доступен"
    
    # Проверка наличия buckets
    BUCKETS=$(curl -s http://localhost:5600/api/0/buckets/ | jq 'length' 2>/dev/null || echo "0")
    
    if [ "$BUCKETS" -gt 0 ]; then
        echo "✅ Найдено $BUCKETS bucket(s)"
        echo ""
        echo "Доступные buckets:"
        curl -s http://localhost:5600/api/0/buckets/ | jq -r 'keys[]' 2>/dev/null || echo "  (не удалось получить список)"
    else
        echo "⚠️  ActivityWatch запущен, но ещё не собрал данные"
        echo "   Поработайте на компьютере несколько минут и попробуйте снова"
    fi
    
    echo ""
    echo "Web интерфейс: http://localhost:5600"
    exit 0
fi

echo "❌ ActivityWatch НЕ запущен"
echo ""
echo "═══════════════════════════════════════════════════════════"
echo "📥 Инструкция по установке ActivityWatch"
echo "═══════════════════════════════════════════════════════════"
echo ""

# Определение ОС
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    echo "Обнаружена ОС: Linux"
    echo ""
    echo "Установка через пакетный менеджер:"
    echo ""
    
    if command -v yay &> /dev/null; then
        echo "  yay -S activitywatch-bin"
    elif command -v pacman &> /dev/null; then
        echo "  # Для Arch Linux:"
        echo "  yay -S activitywatch-bin"
    elif command -v apt &> /dev/null; then
        echo "  # Для Ubuntu/Debian:"
        echo "  # Скачайте .deb с https://activitywatch.net/downloads/"
        echo "  wget https://github.com/ActivityWatch/activitywatch/releases/latest/download/activitywatch-v0.XX.X-linux-x86_64.deb"
        echo "  sudo dpkg -i activitywatch-*.deb"
    elif command -v dnf &> /dev/null; then
        echo "  # Для Fedora:"
        echo "  # Скачайте .rpm с https://activitywatch.net/downloads/"
    fi
    
    echo ""
    echo "Или скачайте напрямую:"
    echo "  https://activitywatch.net/downloads/"
    echo ""
    echo "После установки запустите:"
    echo "  aw-qt"
    
elif [[ "$OSTYPE" == "darwin"* ]]; then
    echo "Обнаружена ОС: macOS"
    echo ""
    echo "Установка через Homebrew:"
    echo "  brew install --cask activitywatch"
    echo ""
    echo "Или скачайте .dmg:"
    echo "  https://activitywatch.net/downloads/"
    echo ""
    echo "После установки запустите из Applications"
    
elif [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" ]]; then
    echo "Обнаружена ОС: Windows"
    echo ""
    echo "Скачайте установщик:"
    echo "  https://activitywatch.net/downloads/"
    echo ""
    echo "После установки запустите из меню Пуск"
    
else
    echo "Неизвестная ОС: $OSTYPE"
    echo ""
    echo "Скачайте ActivityWatch для вашей системы:"
    echo "  https://activitywatch.net/downloads/"
fi

echo ""
echo "═══════════════════════════════════════════════════════════"
echo ""
echo "После установки и запуска ActivityWatch:"
echo "  1. Дождитесь несколько минут (ActivityWatch собирает данные)"
echo "  2. Запустите этот скрипт снова для проверки"
echo "  3. Запустите сбор данных: make run-aw"
echo ""
