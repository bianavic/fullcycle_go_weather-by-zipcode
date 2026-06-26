package web_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bianavic/fullcycle_go_weather-by-zipcode/internal/domain"
	"github.com/bianavic/fullcycle_go_weather-by-zipcode/internal/infra/web"
	"github.com/bianavic/fullcycle_go_weather-by-zipcode/internal/usecase"
)

type fakeUseCase struct {
	out    usecase.WeatherByZipcodeOutput
	err    error
	gotCEP string
}

func (f *fakeUseCase) Get(_ context.Context, rawCEP string) (usecase.WeatherByZipcodeOutput, error) {
	f.gotCEP = rawCEP
	return f.out, f.err
}

func newRequest(t *testing.T, target string) (*httptest.ResponseRecorder, *http.Request) {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, target, nil)
	return httptest.NewRecorder(), req
}

func readBody(t *testing.T, rr *httptest.ResponseRecorder) string {
	t.Helper()
	b, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	return strings.TrimSpace(string(b))
}

func TestRouter(t *testing.T) {
	t.Run("success returns 200 with DTO JSON and forwards cep", func(t *testing.T) {
		uc := &fakeUseCase{out: usecase.WeatherByZipcodeOutput{TempC: 28.5, TempF: 83.3, TempK: 301.5}}
		router := web.NewRouter(uc)

		rr, req := newRequest(t, "/weather/01001000")
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rr.Code)
		}
		if ct := rr.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
			t.Errorf("content-type = %q, want application/json", ct)
		}
		var got usecase.WeatherByZipcodeOutput
		if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
			t.Fatalf("unmarshal body: %v", err)
		}
		want := usecase.WeatherByZipcodeOutput{TempC: 28.5, TempF: 83.3, TempK: 301.5}
		if got != want {
			t.Errorf("body = %+v, want %+v", got, want)
		}
		if uc.gotCEP != "01001000" {
			t.Errorf("usecase received cep = %q, want 01001000", uc.gotCEP)
		}
	})

	t.Run("ErrInvalidZipcode returns 422 with invalid zipcode body", func(t *testing.T) {
		uc := &fakeUseCase{err: domain.ErrInvalidZipcode}
		router := web.NewRouter(uc)

		rr, req := newRequest(t, "/weather/abc")
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnprocessableEntity {
			t.Errorf("status = %d, want 422", rr.Code)
		}
		if body := readBody(t, rr); body != "invalid zipcode" {
			t.Errorf("body = %q, want \"invalid zipcode\"", body)
		}
	})

	t.Run("ErrZipcodeNotFound returns 404 with can not find zipcode body", func(t *testing.T) {
		uc := &fakeUseCase{err: domain.ErrZipcodeNotFound}
		router := web.NewRouter(uc)

		rr, req := newRequest(t, "/weather/00000000")
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("status = %d, want 404", rr.Code)
		}
		if body := readBody(t, rr); body != "can not find zipcode" {
			t.Errorf("body = %q, want \"can not find zipcode\"", body)
		}
	})

	t.Run("unexpected error returns 500", func(t *testing.T) {
		uc := &fakeUseCase{err: errors.New("boom")}
		router := web.NewRouter(uc)

		rr, req := newRequest(t, "/weather/01001000")
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Errorf("status = %d, want 500", rr.Code)
		}
	})

	t.Run("GET /health returns 200", func(t *testing.T) {
		router := web.NewRouter(&fakeUseCase{})

		rr, req := newRequest(t, "/health")
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("status = %d, want 200", rr.Code)
		}
	})
}
