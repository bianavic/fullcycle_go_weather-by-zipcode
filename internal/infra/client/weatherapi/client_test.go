package weatherapi_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bianavic/fullcycle_go_weather-by-zipcode/internal/infra/client/weatherapi"
)

func TestClientCurrentTempC(t *testing.T) {
	t.Run("returns temp_c on success and sends key + url-encoded city", func(t *testing.T) {
		var gotQuery string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotQuery = r.URL.RawQuery
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"location":{"name":"Sao Paulo"},"current":{"temp_c":28.5}}`))
		}))
		defer server.Close()

		client := weatherapi.New(server.Client(), server.URL, "secret-key")
		tempC, err := client.CurrentTempC(context.Background(), "São Paulo")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tempC != 28.5 {
			t.Errorf("tempC = %v, want 28.5", tempC)
		}
		// Query must include both key and url-encoded city
		want := "key=secret-key&q=S%C3%A3o+Paulo"
		if gotQuery != want {
			t.Errorf("query = %q, want %q", gotQuery, want)
		}
	})

	t.Run("returns wrapped error on non-200 status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client := weatherapi.New(server.Client(), server.URL, "secret-key")
		_, err := client.CurrentTempC(context.Background(), "São Paulo")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("returns wrapped error on malformed JSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(`{not json`))
		}))
		defer server.Close()

		client := weatherapi.New(server.Client(), server.URL, "secret-key")
		_, err := client.CurrentTempC(context.Background(), "São Paulo")
		if err == nil {
			t.Fatal("expected decode error, got nil")
		}
	})

	t.Run("propagates context cancellation", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(`{"current":{"temp_c":28.5}}`))
		}))
		defer server.Close()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		client := weatherapi.New(server.Client(), server.URL, "secret-key")
		_, err := client.CurrentTempC(ctx, "São Paulo")
		if err == nil {
			t.Fatal("expected context error, got nil")
		}
	})
}
