package viacep

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/bianavic/fullcycle_go_weather-by-zipcode/internal/domain"
)

type Client struct {
	http    *http.Client
	baseURL string
}

func New(httpClient *http.Client, baseURL string) *Client {
	return &Client{http: httpClient, baseURL: baseURL}
}

func (c *Client) Find(ctx context.Context, zip domain.Zipcode) (string, error) {
	url := strings.TrimRight(c.baseURL, "/") + "/" + zip.String() + "/json/"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("viacep request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("viacep do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		return "", domain.ErrZipcodeNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("viacep unexpected status %d", resp.StatusCode)
	}

	var payload struct {
		Erro       json.RawMessage `json:"erro"`
		Localidade string          `json:"localidade"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", fmt.Errorf("viacep decode: %w", err)
	}
	if isViaCEPNotFound(payload.Erro) {
		return "", domain.ErrZipcodeNotFound
	}
	return payload.Localidade, nil
}

func isViaCEPNotFound(raw json.RawMessage) bool {
	s := strings.TrimSpace(string(raw))
	return s == "true" || s == `"true"`
}
