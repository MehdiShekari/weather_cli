# Weather CLI – Go

A beautiful, production‑ready command‑line weather tool with **smart caching** (Redis or file‑based), coloured terminal output, and support for both city names and geographic coordinates.

## ✨ Features

- Fetch current weather by **city name** or **latitude/longitude**
- **Automatic caching** – responses are stored for 10 minutes (TTL configurable)
- Two cache backends:
  - **File cache** (JSON, persistent across runs) – no extra services required
  - **Redis cache** – faster, ideal for development or shared environments
- `-refresh` flag to **force a fresh API call** (bypass cache)
- Beautiful terminal output with:
  - Emoji icons based on weather conditions
  - Coloured temperature, labels, and separators
  - Clean, human‑readable layout
- Graceful fallback: if Redis is unavailable, the CLI automatically switches to file cache
- Environment variable support (`.env` file or system env) for the API key

## 🚀 Quick Start

### 1. Clone the repository

```bash
git clone https://github.com/MehdiShekari/weather_cli.git
cd weather_cli
```

### 2. Get an API key

Sign up at [OpenWeatherMap](https://openweathermap.org/api) (free tier works perfectly) and copy your API key.

### 3. Set your API key

Create a `.env` file in the project root:

```env
OPENWEATHER_API_KEY=your_api_key_here
# Optional – Redis cache (see section below)
REDIS_ADDR=localhost:6379
```

Alternatively, export the variable in your shell:

```bash
export OPENWEATHER_API_KEY=your_api_key_here
```

### 4. Build and run

```bash
go mod tidy
go build -o weather_cli .
```

**On Windows** (the executable will be `weather_cli.exe`):

```cmd
go build -o weather_cli.exe .
```

### 5. First weather query

```bash
./weather_cli London
```

Or on Windows:

```cmd
weather_cli.exe Paris
```

## 📖 Usage

### Basic commands

| Command | Description |
|---------|-------------|
| `weather_cli <city>` | Weather for a city (e.g., `Tehran`, `New York`) |
| `weather_cli -city <city>` | Same as above |
| `weather_cli -lat <lat> -lon <lon>` | Weather for coordinates (e.g., `-lat 35.6892 -lon 51.3890`) |
| `weather_cli <city> -refresh` | Ignore cache, fetch fresh data |

### Examples

```bash
# By city name
weather_cli Tokyo

# With explicit flag
weather_cli -city "Mexico City"

# By coordinates (Mount Everest)
weather_cli -lat 27.9881 -lon 86.9250

# Force refresh (bypass cache)
weather_cli Paris -refresh
```

## 💾 Caching behaviour

- **Cache TTL**: 10 minutes (hardcoded, easy to change in `weather.go`)
- **Cache key**: lowercased city name or `"lat,lon"` string
- **First request** – fetches from API, stores result
- **Subsequent requests within TTL** – returns cached data (shows `📦 Using cached weather data`)
- **With `-refresh`** – ignores cache, always calls API, then updates cache

### Cache backends

#### File cache (default)

No setup required. Cache is stored in your OS user cache directory:
- Windows: `%LOCALAPPDATA%\weather_cli\cache.json`
- Linux/macOS: `~/.cache/weather_cli/cache.json`

#### Redis cache (optional)

Set the environment variable `REDIS_ADDR` (in `.env` or system). Example:

```env
REDIS_ADDR=localhost:6379
```

**Running Redis on Windows** (for development):

Download a Windows build from the community port:  
[Redis‑8.8.0‑Windows‑x64‑cygwin.zip](https://github.com/redis-windows/redis-windows/releases/download/8.8.0/Redis-8.8.0-Windows-x64-cygwin.zip)

1. Extract to `C:\Redis`
2. Run `redis-server.exe` (keep the terminal open)
3. Set `REDIS_ADDR=localhost:6379`
4. Run the CLI – you’ll see `✅ Using Redis cache`

**With Docker** (any OS):

```bash
docker run -d -p 6379:6379 --name redis redis:alpine
```

If Redis is unreachable, the CLI automatically falls back to file cache (warning shown).

## 🖥️ Example output

```
📁 Using file cache (set REDIS_ADDR to enable Redis)
🌐 Fetching fresh data from API...

 🌍 Weather in Tehran, IR

 ☀️  Clear Sky

 24.5°C (feels like 23.8°C)
────────────────────────────────
💧 Humidity    : 38%
🌬️ Wind        : 2.1 m/s
🔽 Pressure    : 1012 hPa
────────────────────────────────
🕒 Cache TTL    : 10 minutes
```

## 🧪 Testing the `-refresh` flag

```bash
# First run – API call (cache miss)
./weather_cli Rome
# Output: 🌐 Fetching fresh data from API...

# Second run – cache hit
./weather_cli Rome
# Output: 📦 Using cached weather data

# Force fresh data
./weather_cli Rome -refresh
# Output: 🔄 Refresh mode: ignoring cache
#         🌐 Fetching fresh data from API...
```

## 📁 Project structure

```
weather_cli/
├── main.go         # CLI entry point, flag parsing, cache initialisation
├── cache.go        # Cache interface + FileCache + RedisCache
├── weather.go      # OpenWeatherMap API client
├── display.go      # Beautiful coloured output
├── go.mod          # Dependencies (fatih/color, go-redis, godotenv)
└── .env            # Your API key (not committed)
```

## 🔧 Dependencies

- [fatih/color](https://github.com/fatih/color) – coloured terminal output
- [go-redis/redis](https://github.com/go-redis/redis) – Redis client
- [joho/godotenv](https://github.com/joho/godotenv) – `.env` file loader

Install all with `go mod tidy`.

## 📜 License

MIT – free to use and modify.

---

**Enjoy the weather!** ☀️🌧️❄️  
Made with ❤️ by [Mehdi Shekari](https://github.com/MehdiShekari)

---
