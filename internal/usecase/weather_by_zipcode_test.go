package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bianavic/fullcycle_go_weather-by-zipcode/internal/domain"
	"github.com/bianavic/fullcycle_go_weather-by-zipcode/internal/usecase"
)

type fakeCities struct {
	city   string
	err    error
	called bool
}

func (f *fakeCities) Find(_ context.Context, _ domain.Zipcode) (string, error) {
	f.called = true
	return f.city, f.err
}

type fakeWeather struct {
	tempC  float64
	err    error
	called bool
}

func (f *fakeWeather) CurrentTempC(_ context.Context, _ string) (float64, error) {
	f.called = true
	return f.tempC, f.err
}

func TestWeatherByZipcodeExecute(t *testing.T) {
	t.Run("returns DTO with C/F/K on success", func(t *testing.T) {
		cities := &fakeCities{city: "São Paulo"}
		wx := &fakeWeather{tempC: 28.5}
		uc := usecase.NewWeatherByZipcode(cities, wx)

		got, err := uc.Get(context.Background(), "01001000")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := usecase.WeatherByZipcodeOutput{TempC: 28.5, TempF: 83.3, TempK: 301.5}
		if got != want {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})

	t.Run("returns ErrInvalidZipcode without calling adapters", func(t *testing.T) {
		cities := &fakeCities{}
		wx := &fakeWeather{}
		uc := usecase.NewWeatherByZipcode(cities, wx)

		_, err := uc.Get(context.Background(), "abc")
		if !errors.Is(err, domain.ErrInvalidZipcode) {
			t.Fatalf("err = %v, want %v", err, domain.ErrInvalidZipcode)
		}
		if cities.called || wx.called {
			t.Errorf("adapters should not be called on invalid input (cities=%v, wx=%v)", cities.called, wx.called)
		}
	})

	t.Run("propagates ErrZipcodeNotFound from city finder", func(t *testing.T) {
		cities := &fakeCities{err: domain.ErrZipcodeNotFound}
		wx := &fakeWeather{}
		uc := usecase.NewWeatherByZipcode(cities, wx)

		_, err := uc.Get(context.Background(), "00000000")
		if !errors.Is(err, domain.ErrZipcodeNotFound) {
			t.Fatalf("err = %v, want %v", err, domain.ErrZipcodeNotFound)
		}
		if wx.called {
			t.Errorf("weather finder should not be called when city lookup fails")
		}
	})

	t.Run("wraps unexpected city finder error", func(t *testing.T) {
		boom := errors.New("network down")
		cities := &fakeCities{err: boom}
		wx := &fakeWeather{}
		uc := usecase.NewWeatherByZipcode(cities, wx)

		_, err := uc.Get(context.Background(), "01001000")
		if !errors.Is(err, boom) {
			t.Fatalf("err = %v, want wrapped %v", err, boom)
		}
		if errors.Is(err, domain.ErrZipcodeNotFound) || errors.Is(err, domain.ErrInvalidZipcode) {
			t.Errorf("unexpected domain error reported: %v", err)
		}
	})

	t.Run("wraps weather finder error", func(t *testing.T) {
		boom := errors.New("weather api 500")
		cities := &fakeCities{city: "São Paulo"}
		wx := &fakeWeather{err: boom}
		uc := usecase.NewWeatherByZipcode(cities, wx)

		_, err := uc.Get(context.Background(), "01001000")
		if !errors.Is(err, boom) {
			t.Fatalf("err = %v, want wrapped %v", err, boom)
		}
	})
}
