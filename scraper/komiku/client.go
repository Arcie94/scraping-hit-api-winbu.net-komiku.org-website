package komiku

import (
	"komiku-scraper/scraper/common"
)

// KomikuClient wraps the common BaseClient with Komiku-specific functionality
type KomikuClient struct {
	*common.BaseClient
}

// NewKomikuClient creates a new Komiku scraper client
func NewKomikuClient() *KomikuClient {
	return &KomikuClient{
		BaseClient: common.NewBaseClient("Komiku"),
	}
}

// No custom Do method needed - Komiku uses BaseClient's Do directly
// This demonstrates how simple a client can be when using shared infrastructure
