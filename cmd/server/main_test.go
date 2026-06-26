package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bianavic/fullcycle_go_weather-by-zipcode/internal/config"
)

func TestNewRouter(t *testing.T) {
	cfg := config.Config{
		ServerPort:    "8080",
		WeatherAPIKey: "test-key",
		WeatherAPIURL: "http://example.invalid/current.json",
		ViaCEPAPIURL:  "http://example.invalid/ws",
	}
	router := newRouter(cfg)

	t.Run("GET /health returns 200 ok", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rr.Code)
		}
		if body := strings.TrimSpace(rr.Body.String()); body != "ok" {
			t.Errorf("body = %q, want \"ok\"", body)
		}
	})

	t.Run("invalid zipcode returns 422 without hitting external APIs", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/weather/abc", nil)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnprocessableEntity {
			t.Fatalf("status = %d, want 422", rr.Code)
		}
		if body := strings.TrimSpace(rr.Body.String()); body != "invalid zipcode" {
			t.Errorf("body = %q, want \"invalid zipcode\"", body)
		}
	})
}
