package main

import (
	"log"
	"net/http"
	"time"

	"github.com/bianavic/fullcycle_go_weather-by-zipcode/internal/config"
	"github.com/bianavic/fullcycle_go_weather-by-zipcode/internal/infra/client/viacep"
	"github.com/bianavic/fullcycle_go_weather-by-zipcode/internal/infra/client/weatherapi"
	"github.com/bianavic/fullcycle_go_weather-by-zipcode/internal/infra/web"
	"github.com/bianavic/fullcycle_go_weather-by-zipcode/internal/usecase"
)

const httpClientTimeout = 10 * time.Second

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	addr := ":" + cfg.ServerPort
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, newRouter(cfg)); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func newRouter(cfg config.Config) http.Handler {
	httpClient := &http.Client{Timeout: httpClientTimeout}
	cities := viacep.New(httpClient, cfg.ViaCEPAPIURL)
	weather := weatherapi.New(httpClient, cfg.WeatherAPIURL, cfg.WeatherAPIKey)
	uc := usecase.NewWeatherByZipcode(cities, weather)
	return web.NewRouter(uc)
}
