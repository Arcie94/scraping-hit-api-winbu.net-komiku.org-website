package common

import (
	"log"
	"net/http"
	"time"
)

// BaseClient provides common HTTP client functionality for all scrapers
type BaseClient struct {
	Client      *http.Client
	ServiceName string // e.g., "Winbu" or "Komiku"
}

// NewBaseClient creates a new BaseClient with default configuration
func NewBaseClient(serviceName string) *BaseClient {
	return &BaseClient{
		Client: &http.Client{
			Timeout: DefaultTimeout * time.Second,
		},
		ServiceName: serviceName,
	}
}

// Do executes an HTTP request with common headers and logging
func (c *BaseClient) Do(req *http.Request) (*http.Response, error) {
	// Set User-Agent header
	req.Header.Set("User-Agent", ChromeAndroidUserAgent)

	// Log the request
	log.Printf("[%s] Fetching: %s", c.ServiceName, req.URL.String())

	// Execute request
	resp, err := c.Client.Do(req)
	if err != nil {
		log.Printf("[%s] Request error: %v", c.ServiceName, err)
		return nil, err
	}

	return resp, nil
}

// SetCustomHeader allows scrapers to add custom headers before making request
func (c *BaseClient) SetCustomHeader(req *http.Request, key, value string) {
	if req.Header.Get(key) == "" {
		req.Header.Set(key, value)
	}
}
