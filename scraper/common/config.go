package common

const (
	// ChromeAndroidUserAgent is the User-Agent string for Chrome on Android
	ChromeAndroidUserAgent = "Mozilla/5.0 (Linux; Android 13; SM-S908B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.6099.230 Mobile Safari/537.36"

	// DefaultTimeout is the default HTTP client timeout
	DefaultTimeout = 60 // Keeping 60s timeout as it is safer
)

// Base URLs for scraper targets
const (
	WinbuBaseURL  = "https://winbu.net"
	KomikuBaseURL = "https://komiku.id"
)
