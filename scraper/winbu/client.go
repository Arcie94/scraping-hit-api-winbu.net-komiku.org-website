package winbu

import (
	"komiku-scraper/scraper/common"
	"net/http"
)

// WinbuClient wraps the common BaseClient with Winbu-specific functionality
type WinbuClient struct {
	*common.BaseClient
}

// NewWinbuClient creates a new Winbu scraper client
func NewWinbuClient() *WinbuClient {
	return &WinbuClient{
		BaseClient: common.NewBaseClient("Winbu"),
	}
}

// Do executes an HTTP request with Winbu-specific headers
func (c *WinbuClient) Do(req *http.Request) (*http.Response, error) {
	// Add Winbu-specific Referer header if not set
	c.SetCustomHeader(req, "Referer", common.WinbuBaseURL+"/")

	// Use BaseClient's Do method which handles User-Agent and logging
	return c.BaseClient.Do(req)
}
