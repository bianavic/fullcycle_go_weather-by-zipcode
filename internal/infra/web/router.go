package web

import "net/http"

func NewRouter(uc WeatherUseCase) *http.ServeMux {
	mux := http.NewServeMux()
	handler := NewWeatherHandler(uc)

	mux.HandleFunc("GET /weather/{cep}", handler.Get)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	return mux
}
