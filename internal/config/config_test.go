package config_test

import (
	"errors"
	"testing"

	"github.com/bianavic/fullcycle_go_weather-by-zipcode/internal/config"
)

func TestLoadConfig(t *testing.T) {
	t.Run("returns config when WEATHER_API_KEY is set", func(t *testing.T) {
		t.Setenv("WEATHER_API_KEY", "test-key")
		t.Setenv("SERVER_PORT", "9090")
		t.Setenv("WEATHER_API_URL", "https://weather.example/v1/current.json")
		t.Setenv("VIACEP_API_URL", "https://viacep.example/ws")

		cfg, err := config.LoadConfig()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := config.Config{
			ServerPort:    "9090",
			WeatherAPIKey: "test-key",
			WeatherAPIURL: "https://weather.example/v1/current.json",
			ViaCEPAPIURL:  "https://viacep.example/ws",
		}
		if cfg != want {
			t.Errorf("got %+v, want %+v", cfg, want)
		}
	})

	t.Run("falls back to defaults for non-secret values", func(t *testing.T) {
		t.Setenv("WEATHER_API_KEY", "test-key")
		t.Setenv("SERVER_PORT", "")
		t.Setenv("WEATHER_API_URL", "")
		t.Setenv("VIACEP_API_URL", "")

		cfg, err := config.LoadConfig()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.ServerPort != "8080" {
			t.Errorf("ServerPort = %q, want 8080", cfg.ServerPort)
		}
		if cfg.WeatherAPIURL == "" || cfg.ViaCEPAPIURL == "" {
			t.Errorf("expected default URLs, got %+v", cfg)
		}
	})

	t.Run("PORT overrides SERVER_PORT for Cloud Run", func(t *testing.T) {
		t.Setenv("WEATHER_API_KEY", "test-key")
		t.Setenv("SERVER_PORT", "8080")
		t.Setenv("PORT", "8081")

		cfg, err := config.LoadConfig()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.ServerPort != "8081" {
			t.Errorf("ServerPort = %q, want 8081", cfg.ServerPort)
		}
	})

	t.Run("fails when WEATHER_API_KEY is missing", func(t *testing.T) {
		t.Setenv("WEATHER_API_KEY", "")

		_, err := config.LoadConfig()
		if !errors.Is(err, config.ErrMissingWeatherAPIKey) {
			t.Errorf("got %v, want ErrMissingWeatherAPIKey", err)
		}
	})
}
