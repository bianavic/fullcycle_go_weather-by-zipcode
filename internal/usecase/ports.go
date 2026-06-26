package usecase

import (
	"context"

	"github.com/bianavic/fullcycle_go_weather-by-zipcode/internal/domain"
)

type CityFinder interface {
	Find(ctx context.Context, zipcode domain.Zipcode) (city string, err error)
}

type WeatherFinder interface {
	CurrentTempC(ctx context.Context, city string) (float64, error)
}
