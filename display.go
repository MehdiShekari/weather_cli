package main

import (
    "fmt"
    "strings"

    "github.com/fatih/color"
)

// weatherIconMap maps OpenWeatherMap icon codes to emoji/UTF8 symbols
func weatherIcon(iconCode string) string {
    switch {
    case strings.Contains(iconCode, "01"):
        return "☀️"
    case strings.Contains(iconCode, "02"):
        return "⛅"
    case strings.Contains(iconCode, "03"):
        return "☁️"
    case strings.Contains(iconCode, "04"):
        return "☁️"
    case strings.Contains(iconCode, "09"):
        return "🌧️"
    case strings.Contains(iconCode, "10"):
        return "🌦️"
    case strings.Contains(iconCode, "11"):
        return "⛈️"
    case strings.Contains(iconCode, "13"):
        return "❄️"
    case strings.Contains(iconCode, "50"):
        return "🌫️"
    default:
        return "🌈"
    }
}

// DisplayWeather prints a coloured, well‑formatted weather report
func DisplayWeather(w *WeatherData) {
    // Colours
    titleColor := color.New(color.FgCyan, color.Bold)
    tempColor := color.New(color.FgYellow, color.Bold)
    labelColor := color.New(color.FgHiBlack)
    valueColor := color.New(color.FgWhite)
    sepColor := color.New(color.FgHiBlack)

    // Header
    titleColor.Printf("\n 🌍 Weather in %s", w.City)
    if w.Country != "" {
        fmt.Printf(", %s", w.Country)
    }
    fmt.Println()

    // Icon + main condition
    icon := weatherIcon(w.IconCode)
    fmt.Printf("\n %s  %s\n", icon, strings.Title(w.Description))

    // Temperature
    tempColor.Printf("\n %.1f°C", w.TempC)
    valueColor.Printf(" (feels like %.1f°C)\n", w.FeelsLikeC)

    sepColor.Println(strings.Repeat("─", 32))

    // Details
    labelColor.Printf("💧 Humidity    : ")
    valueColor.Printf("%d%%\n", w.Humidity)

    labelColor.Printf("🌬️ Wind        : ")
    valueColor.Printf("%.1f m/s\n", w.WindSpeed)

    labelColor.Printf("🔽 Pressure    : ")
    valueColor.Printf("%d hPa\n", w.Pressure)

    sepColor.Println(strings.Repeat("─", 32))

    // Footer
    labelColor.Printf("🕒 Cache TTL    : ")
    valueColor.Printf("%.0f minutes\n", cacheTTL.Minutes())
    fmt.Println()
}