package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	city    string
	lat     float64
	lon     float64
	refresh bool
)

func init() {
	flag.StringVar(&city, "city", "", "City name (e.g., 'Tehran')")
	flag.Float64Var(&lat, "lat", 0, "Latitude (if using coordinates)")
	flag.Float64Var(&lon, "lon", 0, "Longitude (if using coordinates)")
	flag.BoolVar(&refresh, "refresh", false, "Bypass cache and force new API call")
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [city]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nExamples:")
		fmt.Fprintln(os.Stderr, "  weather_cli Tehran")
		fmt.Fprintln(os.Stderr, "  weather_cli -city London")
		fmt.Fprintln(os.Stderr, "  weather_cli -lat 35.6892 -lon 51.3890")
		fmt.Fprintln(os.Stderr, "  weather_cli Tokyo -refresh")
	}
	flag.Parse()

	// Load .env file if present (ignore error if missing)
	_ = godotenv.Load()

	// positional argument if no -city flag
	if city == "" && flag.NArg() > 0 {
		city = flag.Arg(0)
	}

	if city == "" && (lat == 0 && lon == 0) {
		fmt.Println("❌ Error: Please provide a city name or coordinates (--lat, --lon)")
		flag.Usage()
		os.Exit(1)
	}

	// API key from environment (falls back to .env)
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		fmt.Println("❌ Error: OPENWEATHER_API_KEY not found in environment or .env file")
		fmt.Println("   Create a .env file with: OPENWEATHER_API_KEY=your_key_here")
		fmt.Println("   Or export it: export OPENWEATHER_API_KEY=your_key_here")
		os.Exit(1)
	}

	// Determine cache backend
	var cache Cache
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr != "" {
		rc, err := NewRedisCache(redisAddr)
		if err != nil {
			fmt.Printf("⚠️  Redis connection failed (%v), falling back to file cache\n", err)
			cache, _ = NewFileCache()
		} else {
			cache = rc
			defer cache.Close()
			fmt.Println("✅ Using Redis cache")
		}
	} else {
		fc, err := NewFileCache()
		if err != nil {
			fmt.Printf("❌ Failed to initialise file cache: %v\n", err)
			os.Exit(1)
		}
		cache = fc
		defer cache.Close()
		fmt.Println("📁 Using file cache (set REDIS_ADDR to enable Redis)")
	}

	// Build cache key
	cacheKey := strings.ToLower(city)
	if cacheKey == "" {
		cacheKey = fmt.Sprintf("%f,%f", lat, lon)
	}

	var weather *WeatherData
	var err error

	if !refresh {
		weather, err = cache.Get(cacheKey)
		if err == nil {
			fmt.Println("📦 Using cached weather data")
		} else {
			fmt.Println("🌐 Fetching fresh data from API...")
			weather, err = FetchWeather(city, lat, lon, apiKey)
			if err == nil {
				_ = cache.Set(cacheKey, weather)
			}
		}
	} else {
		fmt.Println("🔄 Refresh mode: ignoring cache")
		weather, err = FetchWeather(city, lat, lon, apiKey)
		if err == nil {
			_ = cache.Set(cacheKey, weather)
		}
	}

	if err != nil {
		fmt.Printf("❌ Failed to get weather: %v\n", err)
		os.Exit(1)
	}

	DisplayWeather(weather)
}
