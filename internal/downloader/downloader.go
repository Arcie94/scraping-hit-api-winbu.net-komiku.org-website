package downloader

import (
	"crypto/tls"
	"komiku-scraper/scraper/common"
	"net/http"
	"os"
	"strings"
	"time"
)

// Downloader manages file downloads
type Downloader struct {
	Client  *http.Client
	BaseDir string
}

// New creates a new Downloader instance
func New() *Downloader {
	// Create base Downloads directory
	baseDir := "Downloads"
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		os.Mkdir(baseDir, 0755)
	}

	return &Downloader{
		BaseDir: baseDir,
		Client: &http.Client{
			Timeout: 2 * time.Minute, // Longer timeout for large files
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
}

// SanitizeFilename encodes string to safe filename
func SanitizeFilename(name string) string {
	replacer := strings.NewReplacer(
		"/", "-",
		"\\", "-",
		":", "-",
		"*", "",
		"?", "",
		"\"", "",
		"<", "",
		">", "",
		"|", "",
	)
	return strings.TrimSpace(replacer.Replace(name))
}

// EnsureDir creates directory if not exists
func EnsureDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

// GetRequest creates a request with user agent
func (d *Downloader) GetRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", common.ChromeAndroidUserAgent)
	return d.Client.Do(req)
}
