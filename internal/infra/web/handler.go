package web

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bianavic/fullcycle_go_weather-by-zipcode/internal/domain"
	"github.com/bianavic/fullcycle_go_weather-by-zipcode/internal/usecase"
)

type WeatherUseCase interface {
	Get(ctx context.Context, rawCEP string) (usecase.WeatherByZipcodeOutput, error)
}

type WeatherHandler struct {
	usecase WeatherUseCase
}

func NewWeatherHandler(uc WeatherUseCase) *WeatherHandler {
	return &WeatherHandler{usecase: uc}
}

func (h *WeatherHandler) Get(w http.ResponseWriter, r *http.Request) {
	cep := r.PathValue("cep")

	out, err := h.usecase.Get(r.Context(), cep)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidZipcode):
			writeText(w, http.StatusUnprocessableEntity, "invalid zipcode")
		case errors.Is(err, domain.ErrZipcodeNotFound):
			writeText(w, http.StatusNotFound, "can not find zipcode")
		default:
			writeText(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(out)
}

func writeText(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(body))
}
