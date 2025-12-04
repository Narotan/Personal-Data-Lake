#!/bin/bash

API_KEY="ac9dce6189d8d3983779004612684f9e86e5033b161deb38273c72892b6039d2"
BASE_URL="http://localhost:8080"

echo "Starting traffic generation to $BASE_URL..."

for i in {1..20}; do
    echo "Iteration $i"

    # 1. Health check (200 OK)
    curl -s -o /dev/null -w "Health: %{http_code}\n" "$BASE_URL/health"

    # 2. Root Page (Valid)
    curl -s -o /dev/null -w "Root (Valid): %{http_code}\n" \
        -H "X-API-Key: $API_KEY" \
        "$BASE_URL/"

    # 3. Root Page (Unauthorized)
    curl -s -o /dev/null -w "Root (Unauthorized): %{http_code}\n" \
        "$BASE_URL/"

    # 4. Metrics (200 OK)
    curl -s -o /dev/null -w "Metrics: %{http_code}\n" "$BASE_URL/metrics"

    # 3. WakaTime Stats (Valid)
    curl -s -o /dev/null -w "WakaTime (Valid): %{http_code}\n" \
        -H "X-API-Key: $API_KEY" \
        "$BASE_URL/api/v1/wakatime/stats"

    # 4. Google Fit Stats (Valid)
    curl -s -o /dev/null -w "GoogleFit (Valid): %{http_code}\n" \
        -H "X-API-Key: $API_KEY" \
        "$BASE_URL/api/v1/googlefit/stats"

    # 5. Unauthorized (401)
    curl -s -o /dev/null -w "Unauthorized: %{http_code}\n" \
        "$BASE_URL/api/v1/wakatime/stats"

    # 6. Bad Request (400)
    curl -s -o /dev/null -w "Bad Request: %{http_code}\n" \
        -H "X-API-Key: $API_KEY" \
        "$BASE_URL/api/v1/wakatime/stats?start_date=invalid-date"

    sleep 1
done

echo "Traffic generation complete."
