package config

import (
	"errors"
	"os"
)

type Config struct {
	ServerPort    string
	WeatherAPIKey string
	WeatherAPIURL string
	ViaCEPAPIURL  string
}

var ErrMissingWeatherAPIKey = errors.New("WEATHER_API_KEY is required")

func LoadConfig() (Config, error) {
	cfg := Config{
		ServerPort:    getEnv("SERVER_PORT", "8080"),
		WeatherAPIKey: getEnv("WEATHER_API_KEY", ""),
		WeatherAPIURL: getEnv("WEATHER_API_URL", "https://api.weatherapi.com/v1/current.json"),
		ViaCEPAPIURL:  getEnv("VIACEP_API_URL", "https://viacep.com.br/ws"),
	}

	if cfg.WeatherAPIKey == "" {
		return Config{}, ErrMissingWeatherAPIKey
	}

	if port := os.Getenv("PORT"); port != "" {
		cfg.ServerPort = port
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
