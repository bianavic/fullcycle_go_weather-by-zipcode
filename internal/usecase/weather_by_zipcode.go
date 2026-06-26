package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/bianavic/fullcycle_go_weather-by-zipcode/internal/domain"
)

type WeatherByZipcode struct {
	cities  CityFinder
	weather WeatherFinder
}

func NewWeatherByZipcode(cities CityFinder, weather WeatherFinder) *WeatherByZipcode {
	return &WeatherByZipcode{cities: cities, weather: weather}
}

func (uc *WeatherByZipcode) Get(ctx context.Context, rawCEP string) (WeatherByZipcodeOutput, error) {
	zip, err := domain.NewZipcode(rawCEP)
	if err != nil {
		return WeatherByZipcodeOutput{}, err
	}

	city, err := uc.cities.Find(ctx, zip)
	if err != nil {
		if errors.Is(err, domain.ErrZipcodeNotFound) {
			return WeatherByZipcodeOutput{}, err
		}
		return WeatherByZipcodeOutput{}, fmt.Errorf("find city: %w", err)
	}

	celsius, err := uc.weather.CurrentTempC(ctx, city)
	if err != nil {
		return WeatherByZipcodeOutput{}, fmt.Errorf("fetch weather: %w", err)
	}

	temp := domain.FromCelsius(celsius)
	return WeatherByZipcodeOutput{
		TempC: temp.Celsius(),
		TempF: temp.Fahrenheit(),
		TempK: temp.Kelvin(),
	}, nil
}
