# Auth Package - Руководство по интеграции нового API

## Структура пакета

```
auth/
├── provider.go          # Общий интерфейс для всех OAuth2 провайдеров
├── storage.go           # Хранилище токенов (один файл для всех провайдеров)
├── config.go            # УСТАРЕЛО - можно удалить
├── oauth.go             # УСТАРЕЛО - можно удалить
├── token.go             # УСТАРЕЛО - можно удалить
├── wakatime/
│   └── provider.go      # Реализация для WakaTime API
└── googlefit/
    └── provider.go      # Реализация для Google Fit API (TODO)
```

## Как использовать (WakaTime)

### 1. Инициализация провайдера

```go
import (
    "DataLake/auth"
    "DataLake/auth/wakatime"
)

// Создаем провайдер WakaTime
wakatimeProvider := wakatime.NewProvider(
    "YOUR_CLIENT_ID",
    "YOUR_CLIENT_SECRET",
    "http://localhost:8080/callback/wakatime",
)

// Создаем хранилище токенов (один раз для всех провайдеров)
storage := auth.NewFileTokenStorage("tokens.json")
```

### 2. Авторизация пользователя

```go
// Получаем URL для авторизации
authURL := wakatimeProvider.GetAuthURL("random_state_string")

// Отправляем пользователя на authURL
// После авторизации пользователь будет перенаправлен на redirect_uri с кодом
```

### 3. Обмен кода на токены

```go
// В callback handler получаем код из URL параметров
code := r.URL.Query().Get("code")

// Обмениваем код на токены
token, err := wakatimeProvider.ExchangeToken(r.Context(), code)
if err != nil {
    // Обработка ошибки
}

// Сохраняем токены
err = storage.SaveToken("wakatime", token)
```

### 4. Загрузка сохраненных токенов

```go
// Загружаем токены из хранилища
token, err := storage.LoadToken("wakatime")
if err != nil {
    // Токен не найден или ошибка чтения
}

// Используем токен для запросов к API
req.Header.Set("Authorization", "Bearer " + token.AccessToken)
```

### 5. Обновление токенов

```go
// Если токен истек, обновляем его
newToken, err := wakatimeProvider.RefreshToken(ctx, token.RefreshToken)
if err != nil {
    // Обработка ошибки
}

// Сохраняем новый токен
storage.SaveToken("wakatime", newToken)
```

## Как добавить новый API (пошаговая инструкция)

### Шаг 1: Создать папку провайдера

```bash
mkdir -p auth/your_api_name
```

### Шаг 2: Создать файл `auth/your_api_name/provider.go`

```go
package your_api_name

import (
    "DataLake/auth"
    "DataLake/internal/logger"
    "context"
    "encoding/json"
    "io"
    "net/http"
    "net/url"
)

type Provider struct {
    clientID     string
    clientSecret string
    redirectURI  string
}

func NewProvider(clientID, clientSecret, redirectURI string) *Provider {
    return &Provider{
        clientID:     clientID,
        clientSecret: clientSecret,
        redirectURI:  redirectURI,
    }
}

// GetAuthURL возвращает URL для авторизации пользователя
func (p *Provider) GetAuthURL(state string) string {
    // Пример для стандартного OAuth2
    baseURL := "https://api.example.com/oauth/authorize"
    params := url.Values{}
    params.Set("client_id", p.clientID)
    params.Set("redirect_uri", p.redirectURI)
    params.Set("response_type", "code")
    params.Set("state", state)
    // Добавьте scope если нужно
    params.Set("scope", "read write")
    
    return baseURL + "?" + params.Encode()
}

// ExchangeToken обменивает код авторизации на токены
func (p *Provider) ExchangeToken(ctx context.Context, code string) (auth.TokenResponse, error) {
    log := logger.Get()
    
    // URL для обмена кода на токен
    tokenURL := "https://api.example.com/oauth/token"
    
    // Подготовка данных запроса
    data := url.Values{}
    data.Set("client_id", p.clientID)
    data.Set("client_secret", p.clientSecret)
    data.Set("redirect_uri", p.redirectURI)
    data.Set("grant_type", "authorization_code")
    data.Set("code", code)
    
    // Создание и выполнение запроса
    req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, 
        strings.NewReader(data.Encode()))
    if err != nil {
        log.Error().Err(err).Msg("failed to create request")
        return auth.TokenResponse{}, err
    }
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Error().Err(err).Msg("failed to execute request")
        return auth.TokenResponse{}, err
    }
    defer resp.Body.Close()
    
    // Парсинг ответа
    bodyBytes, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Error().Err(err).Msg("failed to read response body")
        return auth.TokenResponse{}, err
    }
    
    // В зависимости от API, ответ может быть в JSON или URL-encoded
    // Для JSON:
    var response struct {
        AccessToken  string `json:"access_token"`
        RefreshToken string `json:"refresh_token"`
        ExpiresIn    int    `json:"expires_in"`
        TokenType    string `json:"token_type"`
    }
    
    err = json.Unmarshal(bodyBytes, &response)
    if err != nil {
        log.Error().Err(err).Msg("failed to parse response")
        return auth.TokenResponse{}, err
    }
    
    token := auth.TokenResponse{
        AccessToken:  response.AccessToken,
        RefreshToken: response.RefreshToken,
        TokenType:    response.TokenType,
        // Конвертируйте expires_in в expires_at если нужно
        ExpiresAt:    "", 
    }
    
    log.Info().Msg("successfully exchanged code for token")
    return token, nil
}

// RefreshToken обновляет access token используя refresh token
func (p *Provider) RefreshToken(ctx context.Context, refreshToken string) (auth.TokenResponse, error) {
    log := logger.Get()
    
    tokenURL := "https://api.example.com/oauth/token"
    
    data := url.Values{}
    data.Set("client_id", p.clientID)
    data.Set("client_secret", p.clientSecret)
    data.Set("grant_type", "refresh_token")
    data.Set("refresh_token", refreshToken)
    
    // ... аналогично ExchangeToken
    
    return auth.TokenResponse{}, nil
}
```

### Шаг 3: Использование нового провайдера

```go
import "DataLake/auth/your_api_name"

// Создаем провайдер
provider := your_api_name.NewProvider(clientID, clientSecret, redirectURI)

// Используем через интерфейс auth.Provider
var p auth.Provider = provider

// Получаем URL для авторизации
authURL := p.GetAuthURL("state")

// Обмениваем код на токены
token, err := p.ExchangeToken(ctx, code)

// Сохраняем токены
storage.SaveToken("your_api_name", token)
```

### Шаг 4: Добавить callback handler

```go
func handleYourAPICallback(w http.ResponseWriter, r *http.Request) {
    code := r.URL.Query().Get("code")
    state := r.URL.Query().Get("state")
    
    // Проверка state для защиты от CSRF
    
    provider := your_api_name.NewProvider(clientID, clientSecret, redirectURI)
    token, err := provider.ExchangeToken(r.Context(), code)
    if err != nil {
        http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
        return
    }
    
    storage := auth.NewFileTokenStorage("tokens.json")
    err = storage.SaveToken("your_api_name", token)
    if err != nil {
        http.Error(w, "Failed to save token", http.StatusInternalServerError)
        return
    }
    
    fmt.Fprintf(w, "Authorization successful!")
}
```

## Формат хранения токенов

Все токены хранятся в одном JSON файле (`tokens.json`):

```json
{
  "wakatime": {
    "access_token": "waka_xxx",
    "refresh_token": "waka_yyy",
    "expires_at": "2024-12-31T23:59:59Z",
    "token_type": "Bearer",
    "scope": "read_stats",
    "uid": "user123"
  },
  "googlefit": {
    "access_token": "google_xxx",
    "refresh_token": "google_yyy",
    "expires_at": "2024-12-31T23:59:59Z",
    "token_type": "Bearer"
  },
  "your_api_name": {
    "access_token": "your_xxx",
    "refresh_token": "your_yyy",
    "expires_at": "2024-12-31T23:59:59Z"
  }
}
```

## Интерфейс Provider

Все провайдеры должны реализовывать интерфейс `auth.Provider`:

```go
type Provider interface {
    // GetAuthURL возвращает URL для авторизации пользователя
    GetAuthURL(state string) string
    
    // ExchangeToken обменивает код авторизации на токены
    ExchangeToken(ctx context.Context, code string) (TokenResponse, error)
    
    // RefreshToken обновляет access token используя refresh token
    RefreshToken(ctx context.Context, refreshToken string) (TokenResponse, error)
}
```

## Полезные ссылки

- [OAuth 2.0 RFC](https://datatracker.ietf.org/doc/html/rfc6749)
- [WakaTime API Documentation](https://wakatime.com/developers)
- [Google Fit API Documentation](https://developers.google.com/fit)

## Миграция со старого кода

### Было (старый способ):
```go
config := auth.Config{
    ClientID:     "xxx",
    ClientSecret: "yyy",
    RedirectURI:  "zzz",
}
authURL := auth.BuildAuthRequest(config)
tokenResp, _ := auth.ExchangeCodeForToken(config, code)
tokens := auth.Tokens{
    AccessToken:  tokenResp.AccessToken,
    RefreshToken: tokenResp.RefreshToken,
    ExpiresAt:    tokenResp.ExpiresAt,
}
auth.SaveTokens(tokens)
```

### Стало (новый способ):
```go
provider := wakatime.NewProvider("xxx", "yyy", "zzz")
storage := auth.NewFileTokenStorage("tokens.json")

authURL := provider.GetAuthURL("state")
token, _ := provider.ExchangeToken(ctx, code)
storage.SaveToken("wakatime", token)
```

## Примечания

- **Безопасность**: Всегда используйте `state` параметр для защиты от CSRF атак
- **Контекст**: Передавайте `context.Context` для управления таймаутами и отменой операций
- **Логирование**: Все провайдеры используют `internal/logger` для логирования
- **Ошибки**: Всегда проверяйте ошибки при работе с токенами

## Сравнение реализаций провайдеров

### WakaTime vs Google Fit

| Аспект | WakaTime | Google Fit |
|--------|----------|------------|
| **Формат ответа** | URL-encoded (`application/x-www-form-urlencoded`) | JSON (`application/json`) |
| **Парсинг токенов** | `url.ParseQuery()` | `json.Unmarshal()` |
| **ExpiresAt** | Готовая строка в ответе | Вычисляется из `expires_in` (секунды) |
| **Scopes** | Не требуются в URL | Обязательны (Fitness API permissions) |
| **Refresh token** | Всегда возвращается | Требует `access_type=offline` и `prompt=consent` |
| **HTTP Client** | Без таймаута | С таймаутом 10 сек |
| **Логирование** | Структурированное (zerolog) | Структурированное (zerolog) |
| **Обработка ошибок** | Logger + возврат ошибки | Logger + wrapped errors (`fmt.Errorf`) |

### Ключевые различия в GetAuthURL

**WakaTime** - минималистичный:
```go
params.Set("client_id", p.clientID)
params.Set("redirect_uri", p.redirectURI)
params.Set("response_type", "code")
params.Set("state", state)
```

**Google Fit** - расширенный с обязательными параметрами:
```go
q.Set("client_id", p.clientID)
q.Set("redirect_uri", p.redirectURI)
q.Set("response_type", "code")
q.Set("scope", strings.Join(defaultScopes, " "))     // Обязательно!
q.Set("access_type", "offline")                      // Для refresh token
q.Set("prompt", "consent")                            // Гарантия refresh token
q.Set("include_granted_scopes", "true")              // Инкрементальная авторизация
q.Set("state", state)
```

### Ключевые различия в ExchangeToken

**WakaTime** - парсинг URL-encoded:
```go
bodyBytes, err := io.ReadAll(resp.Body)
values, err := url.ParseQuery(string(bodyBytes))

token := auth.TokenResponse{
    AccessToken:  values.Get("access_token"),
    RefreshToken: values.Get("refresh_token"),
    ExpiresAt:    values.Get("expires_at"),  // Уже в формате RFC3339
    Scope:        values.Get("scope"),
    UID:          values.Get("uid"),
}
```

**Google Fit** - парсинг JSON + вычисление ExpiresAt:
```go
body, err := io.ReadAll(resp.Body)

var tokenResp tokenResponse
json.Unmarshal(body, &tokenResp)

return auth.TokenResponse{
    AccessToken:  tokenResp.AccessToken,
    RefreshToken: tokenResp.RefreshToken,
    ExpiresAt:    time.Now().Add(
        time.Duration(tokenResp.ExpiresIn) * time.Second
    ).Format(time.RFC3339),  // Конвертация из секунд в дату
}
```

### Рекомендации при добавлении нового провайдера

1. **Изучите документацию API**:
   - Формат ответа токенов (JSON или URL-encoded)
   - Обязательные параметры (scopes, access_type и т.д.)
   - Особенности получения refresh token

2. **Используйте единый стиль**:
   - Добавляйте логирование через `logger.Get()`
   - Используйте wrapped errors с `fmt.Errorf(..., %w, err)`
   - Добавляйте таймауты на HTTP клиенты
   - Проверяйте `resp.StatusCode`

3. **Адаптируйте под специфику API**:
   - Если API возвращает JSON → используйте подход Google Fit
   - Если API возвращает URL-encoded → используйте подход WakaTime
   - Если нужны scopes → добавьте их как в Google Fit

### Пример гибридной реализации

Если ваш API возвращает JSON, но требует специфичные параметры авторизации:

```go
// GetAuthURL с дополнительными параметрами
func (p *Provider) GetAuthURL(state string) string {
    u, _ := url.Parse("https://api.example.com/oauth/authorize")
    
    q := url.Values{}
    q.Set("client_id", p.clientID)
    q.Set("redirect_uri", p.redirectURI)
    q.Set("response_type", "code")
    q.Set("scope", "read write delete")  // Ваши scopes
    q.Set("state", state)
    
    u.RawQuery = q.Encode()
    return u.String()
}

// ExchangeToken с JSON парсингом
func (p *Provider) ExchangeToken(ctx context.Context, code string) (auth.TokenResponse, error) {
    log := logger.Get()
    
    // ... создание запроса ...
    
    body, err := io.ReadAll(resp.Body)
    
    var response struct {
        AccessToken  string `json:"access_token"`
        RefreshToken string `json:"refresh_token"`
        ExpiresIn    int    `json:"expires_in"`
    }
    
    if err := json.Unmarshal(body, &response); err != nil {
        log.Error().Err(err).Msg("failed to parse response")
        return auth.TokenResponse{}, fmt.Errorf("parse error: %w", err)
    }
    
    log.Info().Msg("successfully exchanged code for token")
    
    return auth.TokenResponse{
        AccessToken:  response.AccessToken,
        RefreshToken: response.RefreshToken,
        ExpiresAt:    time.Now().Add(time.Duration(response.ExpiresIn) * time.Second).Format(time.RFC3339),
    }, nil
}
```

