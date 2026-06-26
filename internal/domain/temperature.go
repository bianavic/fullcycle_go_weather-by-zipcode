package domain

type Temperature struct {
	celsius float64
}

func FromCelsius(c float64) Temperature {
	return Temperature{celsius: c}
}

func (t Temperature) Celsius() float64 {
	return t.celsius
}

func (t Temperature) Fahrenheit() float64 {
	return t.celsius*1.8 + 32
}

func (t Temperature) Kelvin() float64 {
	return t.celsius + 273
}
