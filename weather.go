package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "time"
)

const (
    openWeatherMapURL = "https://api.openweathermap.org/data/2.5/weather"
    cacheTTL          = 10 * time.Minute
)

type WeatherData struct {
    City        string  `json:"city"`
    Country     string  `json:"country"`
    TempC       float64 `json:"temp_c"`
    FeelsLikeC  float64 `json:"feels_like_c"`
    Humidity    int     `json:"humidity"`
    Pressure    int     `json:"pressure"`
    WindSpeed   float64 `json:"wind_speed"`
    Conditions  string  `json:"conditions"`
    IconCode    string  `json:"icon_code"`
    Description string  `json:"description"`
}

type openWeatherResponse struct {
    Weather []struct {
        Main        string `json:"main"`
        Description string `json:"description"`
        Icon        string `json:"icon"`
    } `json:"weather"`
    Main struct {
        Temp     float64 `json:"temp"`
        FeelsLike float64 `json:"feels_like"`
        Humidity int     `json:"humidity"`
        Pressure int     `json:"pressure"`
    } `json:"main"`
    Wind struct {
        Speed float64 `json:"speed"`
    } `json:"wind"`
    Name    string `json:"name"`
    Sys     struct {
        Country string `json:"country"`
    } `json:"sys"`
}

func FetchWeather(city string, lat, lon float64, apiKey string) (*WeatherData, error) {
    params := url.Values{}
    params.Set("appid", apiKey)
    params.Set("units", "metric")
    if city != "" {
        params.Set("q", city)
    } else if lat != 0 || lon != 0 {
        params.Set("lat", fmt.Sprintf("%f", lat))
        params.Set("lon", fmt.Sprintf("%f", lon))
    } else {
        return nil, fmt.Errorf("either city or lat/lon must be provided")
    }

    reqURL := fmt.Sprintf("%s?%s", openWeatherMapURL, params.Encode())
    client := http.Client{Timeout: 10 * time.Second}
    resp, err := client.Get(reqURL)
    if err != nil {
        return nil, fmt.Errorf("HTTP request failed: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
    }

    var owm openWeatherResponse
    if err := json.NewDecoder(resp.Body).Decode(&owm); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }

    // extract city name (if not provided by lat/lon)
    finalCity := owm.Name
    if finalCity == "" && city != "" {
        finalCity = city
    } else if finalCity == "" {
        finalCity = fmt.Sprintf("%.2f,%.2f", lat, lon)
    }

    weather := &WeatherData{
        City:        finalCity,
        Country:     owm.Sys.Country,
        TempC:       owm.Main.Temp,
        FeelsLikeC:  owm.Main.FeelsLike,
        Humidity:    owm.Main.Humidity,
        Pressure:    owm.Main.Pressure,
        WindSpeed:   owm.Wind.Speed,
        Conditions:  owm.Weather[0].Main,
        IconCode:    owm.Weather[0].Icon,
        Description: owm.Weather[0].Description,
    }
    return weather, nil
}