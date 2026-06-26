package domain_test

import (
	"math"
	"testing"

	"github.com/bianavic/fullcycle_go_weather-by-zipcode/internal/domain"
)

func TestTemperatureFromCelsius(t *testing.T) {
	tests := []struct {
		name       string
		celsius    float64
		fahrenheit float64
		kelvin     float64
	}{
		{name: "challenge sample 28.5C", celsius: 28.5, fahrenheit: 83.3, kelvin: 301.5},
		{name: "freezing point", celsius: 0, fahrenheit: 32, kelvin: 273},
		{name: "boiling point", celsius: 100, fahrenheit: 212, kelvin: 373},
		{name: "negative temperature", celsius: -40, fahrenheit: -40, kelvin: 233},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			temp := domain.FromCelsius(tc.celsius)

			if temp.Celsius() != tc.celsius {
				t.Errorf("Celsius() = %v, want %v", temp.Celsius(), tc.celsius)
			}
			if !approxEqual(temp.Fahrenheit(), tc.fahrenheit) {
				t.Errorf("Fahrenheit() = %v, want %v", temp.Fahrenheit(), tc.fahrenheit)
			}
			if !approxEqual(temp.Kelvin(), tc.kelvin) {
				t.Errorf("Kelvin() = %v, want %v", temp.Kelvin(), tc.kelvin)
			}
		})
	}
}

func approxEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-9
}
