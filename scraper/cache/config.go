package cache

import "time"

// TTL constants for different data types
const (
	// HomeTTL for home/popular pages (data changes infrequently)
	HomeTTL = 15 * time.Minute

	// SearchTTL for search results
	SearchTTL = 30 * time.Minute

	// DetailTTL for anime/manga detail pages
	DetailTTL = 1 * time.Hour

	// ChapterTTL for chapter/episode lists
	ChapterTTL = 2 * time.Hour

	// StreamTTL for stream URLs (can expire quickly)
	StreamTTL = 5 * time.Minute
)

// CacheKey formats for consistent key generation
const (
	// Winbu cache key formats
	WinbuHomeKey    = "winbu:home"
	WinbuSearchKey  = "winbu:search:%s"  // winbu:search:naruto
	WinbuDetailKey  = "winbu:detail:%s"  // winbu:detail:/anime/one-piece
	WinbuEpisodeKey = "winbu:episode:%s" // winbu:episode:/anime/one-piece/episode-1

	// Komiku cache key formats
	KomikuHomeKey    = "komiku:home"
	KomikuPopularKey = "komiku:popular"
	KomikuSearchKey  = "komiku:search:%s"  // komiku:search:dandadan
	KomikuDetailKey  = "komiku:detail:%s"  // komiku:detail:/manga/dandadan
	KomikuChapterKey = "komiku:chapter:%s" // komiku:chapter:/manga/dandadan/chapter-223
)
