package weatherapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Client struct {
	http    *http.Client
	baseURL string
	apiKey  string
}

func New(httpClient *http.Client, baseURL, apiKey string) *Client {
	return &Client{http: httpClient, baseURL: baseURL, apiKey: apiKey}
}

func (c *Client) CurrentTempC(ctx context.Context, city string) (float64, error) {
	query := url.Values{}
	query.Set("key", c.apiKey)
	query.Set("q", city)
	endpoint := c.baseURL + "?" + query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return 0, fmt.Errorf("weatherapi request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return 0, fmt.Errorf("weatherapi do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("weatherapi unexpected status %d", resp.StatusCode)
	}

	var payload struct {
		Current struct {
			TempC float64 `json:"temp_c"`
		} `json:"current"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return 0, fmt.Errorf("weatherapi decode: %w", err)
	}
	return payload.Current.TempC, nil
}
