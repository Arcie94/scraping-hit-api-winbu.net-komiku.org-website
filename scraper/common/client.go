package common

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/proxy"
)

// BaseClient provides common HTTP client functionality for all scrapers
type BaseClient struct {
	Client      *http.Client
	ServiceName string // e.g., "Winbu" or "Komiku"
}

// NewBaseClient creates a new BaseClient with default configuration
func NewBaseClient(serviceName string) *BaseClient {
	// Create TLS config that skips certificate verification
	// This is needed because some sites (like Komiku) have expired certificates
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // Skip SSL verification
		},
	}

	// Add SOCKS5 proxy support if enabled
	if ProxyEnabled {
		proxyURL, err := url.Parse(ProxyURL)
		if err == nil {
			dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
			if err == nil {
				transport.Dial = dialer.Dial
				log.Printf("[%s] Using proxy: %s", serviceName, ProxyURL)
			} else {
				log.Printf("[%s] Proxy setup failed: %v", serviceName, err)
			}
		} else {
			log.Printf("[%s] Invalid proxy URL: %v", serviceName, err)
		}
	}

	return &BaseClient{
		Client: &http.Client{
			Timeout:   DefaultTimeout * time.Second,
			Transport: transport,
		},
		ServiceName: serviceName,
	}
}

// Do executes an HTTP request with common headers and logging
func (c *BaseClient) Do(req *http.Request) (*http.Response, error) {
	// Set complete browser-like headers
	req.Header.Set("User-Agent", ChromeAndroidUserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "id-ID,id;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Cache-Control", "max-age=0")

	// Set Referer to make it look like navigation from homepage
	if req.Header.Get("Referer") == "" {
		req.Header.Set("Referer", req.URL.Scheme+"://"+req.URL.Host+"/")
	}

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
