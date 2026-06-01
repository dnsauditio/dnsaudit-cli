package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const BaseURL = "https://dnsaudit.io/api"

type Client struct {
	APIKey     string
	HTTPClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		APIKey: apiKey,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type APIError struct {
	ErrorMsg   string `json:"error"`
	Message    string `json:"message"`
	Code       string `json:"code,omitempty"`
	RetryAfter int    `json:"retryAfter,omitempty"`
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("%s: %s", e.ErrorMsg, e.Message)
	}
	return e.ErrorMsg
}

// DoRequest handles the HTTP request and 429 logic
func (c *Client) DoRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("X-API-Key", c.APIKey)
	req.Header.Set("User-Agent", "dnsaudit-cli/1.0")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 429 {
		var apiErr APIError
		if err := json.Unmarshal(body, &apiErr); err == nil && apiErr.RetryAfter > 0 {
			// Burst limit hit. Let's auto-sleep and retry ONCE.
			fmt.Printf("[-] Burst rate limit hit. Sleeping for %d seconds...\n", apiErr.RetryAfter)
			time.Sleep(time.Duration(apiErr.RetryAfter) * time.Second)
			
			// Retry the request once
			retryResp, err := c.HTTPClient.Do(req)
			if err != nil {
				return nil, err
			}
			defer retryResp.Body.Close()
			body, err = io.ReadAll(retryResp.Body)
			if err != nil {
				return nil, err
			}
			resp = retryResp
		}
	}

	if resp.StatusCode >= 400 {
		var apiErr APIError
		if err := json.Unmarshal(body, &apiErr); err == nil {
			return nil, &apiErr
		}
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

func (c *Client) Scan(domain string) ([]byte, error) {
	reqURL := fmt.Sprintf("%s/v1/scan?domain=%s", BaseURL, url.QueryEscape(domain))
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	return c.DoRequest(req)
}

func (c *Client) ExportJSON(domain string) ([]byte, error) {
	reqURL := fmt.Sprintf("%s/export/json/%s", BaseURL, url.QueryEscape(domain))
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	return c.DoRequest(req)
}

func (c *Client) ExportPDF(domain string) ([]byte, error) {
	reqURL := fmt.Sprintf("%s/export/pdf/%s?format=detailed", BaseURL, url.QueryEscape(domain))
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	return c.DoRequest(req)
}

func (c *Client) History(limit int) ([]byte, error) {
	reqURL := fmt.Sprintf("%s/v1/scan-history?limit=%d", BaseURL, limit)
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	return c.DoRequest(req)
}
