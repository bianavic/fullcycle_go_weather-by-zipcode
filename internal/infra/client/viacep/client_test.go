package viacep_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bianavic/fullcycle_go_weather-by-zipcode/internal/domain"
	"github.com/bianavic/fullcycle_go_weather-by-zipcode/internal/infra/client/viacep"
)

func TestClientFind(t *testing.T) {
	zip := domain.Zipcode("01001000")

	t.Run("returns city on success", func(t *testing.T) {
		var gotPath string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"cep":"01001-000","localidade":"São Paulo","uf":"SP"}`))
		}))
		defer server.Close()

		client := viacep.New(server.Client(), server.URL)
		city, err := client.Find(context.Background(), zip)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if city != "São Paulo" {
			t.Errorf("city = %q, want São Paulo", city)
		}
		if gotPath != "/01001000/json/" {
			t.Errorf("path = %q, want /01001000/json/", gotPath)
		}
	})

	t.Run("returns ErrZipcodeNotFound when payload has erro:true", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(`{"erro":true}`))
		}))
		defer server.Close()

		client := viacep.New(server.Client(), server.URL)
		_, err := client.Find(context.Background(), zip)
		if !errors.Is(err, domain.ErrZipcodeNotFound) {
			t.Errorf("err = %v, want ErrZipcodeNotFound", err)
		}
	})

	t.Run("returns ErrZipcodeNotFound on HTTP 400", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}))
		defer server.Close()

		client := viacep.New(server.Client(), server.URL)
		_, err := client.Find(context.Background(), zip)
		if !errors.Is(err, domain.ErrZipcodeNotFound) {
			t.Errorf("err = %v, want ErrZipcodeNotFound", err)
		}
	})

	t.Run("returns wrapped error on non-200 server status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client := viacep.New(server.Client(), server.URL)
		_, err := client.Find(context.Background(), zip)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if errors.Is(err, domain.ErrZipcodeNotFound) {
			t.Errorf("err should not match ErrZipcodeNotFound, got %v", err)
		}
	})

	t.Run("returns wrapped error on malformed JSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(`{not json`))
		}))
		defer server.Close()

		client := viacep.New(server.Client(), server.URL)
		_, err := client.Find(context.Background(), zip)
		if err == nil {
			t.Fatal("expected decode error, got nil")
		}
	})

	t.Run("propagates context cancellation", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(`{"localidade":"São Paulo"}`))
		}))
		defer server.Close()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		client := viacep.New(server.Client(), server.URL)
		_, err := client.Find(ctx, zip)
		if !errors.Is(err, context.Canceled) {
			t.Errorf("err = %v, want context.Canceled", err)
		}
	})
}
