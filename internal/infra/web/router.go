package web

import "net/http"

func NewRouter(uc WeatherUseCase) *http.ServeMux {
    mux := http.NewServeMux()
    handler := NewWeatherHandler(uc)

    // Rota raiz
    mux.HandleFunc("GET /", func(w http.ResponseWriter, _ *http.Request) {
       w.WriteHeader(http.StatusOK)
       _, _ = w.Write([]byte("Weather API is running! Use /weather/{cep}"))
    })

    // Rota de clima
    mux.HandleFunc("GET /weather/{cep}", handler.Get)

    // Rota de health check
    mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
       w.WriteHeader(http.StatusOK)
       _, _ = w.Write([]byte("ok"))
    })

    return mux
}