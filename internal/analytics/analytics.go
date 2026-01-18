package analytics

import (
	"context"
	"database/sql"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/redis/go-redis/v9"
)

// Analytics handles usage tracking with Redis + SQLite
type Analytics struct {
	redis *redis.Client
	db    *sql.DB
	ctx   context.Context
	mu    sync.Mutex
}

// Event represents an analytics event
type Event struct {
	Timestamp    time.Time
	Endpoint     string
	Method       string
	StatusCode   int
	ResponseTime int64 // milliseconds
	IP           string
	APIKey       string
	UserAgent    string
}

// NewAnalytics creates a new analytics instance
func NewAnalytics() *Analytics {
	// Redis client
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("REDIS_PORT")
	if port == "" {
		port = "6379"
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: host + ":" + port,
		DB:   1, // Use DB 1 for analytics
	})

	// SQLite database
	dbPath := os.Getenv("ANALYTICS_DB_PATH")
	if dbPath == "" {
		dbPath = "./analytics.db"
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Printf("[Analytics] Failed to open SQLite: %v", err)
		return nil
	}

	// Create table if not exists
	schema := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME NOT NULL,
		endpoint TEXT NOT NULL,
		method TEXT NOT NULL,
		status_code INTEGER NOT NULL,
		response_time_ms INTEGER NOT NULL,
		ip TEXT,
		api_key TEXT,
		user_agent TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_timestamp ON events(timestamp);
	CREATE INDEX IF NOT EXISTS idx_endpoint ON events(endpoint);
	`

	if _, err := db.Exec(schema); err != nil {
		log.Printf("[Analytics] Failed to create schema: %v", err)
		return nil
	}

	log.Println("[Analytics] Initialized successfully")

	return &Analytics{
		redis: redisClient,
		db:    db,
		ctx:   context.Background(),
	}
}

// Track records an analytics event
func (a *Analytics) Track(event Event) {
	if a == nil {
		return
	}

	// Update Redis counters (real-time)
	go func() {
		a.redis.Incr(a.ctx, "analytics:total_requests")
		a.redis.Incr(a.ctx, "analytics:endpoint:"+event.Endpoint)
		a.redis.HIncrBy(a.ctx, "analytics:ips", event.IP, 1)

		if event.APIKey != "" {
			a.redis.Incr(a.ctx, "analytics:apikey:"+event.APIKey)
		}

		if event.StatusCode >= 400 {
			a.redis.Incr(a.ctx, "analytics:errors")
		}
	}()

	// Store in SQLite (detailed logs)
	go func() {
		a.mu.Lock()
		defer a.mu.Unlock()

		_, err := a.db.Exec(`
			INSERT INTO events (timestamp, endpoint, method, status_code, response_time_ms, ip, api_key, user_agent)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`, event.Timestamp, event.Endpoint, event.Method, event.StatusCode, event.ResponseTime, event.IP, event.APIKey, event.UserAgent)

		if err != nil {
			log.Printf("[Analytics] Failed to insert event: %v", err)
		}
	}()
}

// GetSummary returns analytics summary
func (a *Analytics) GetSummary() map[string]interface{} {
	if a == nil {
		return map[string]interface{}{"error": "analytics not initialized"}
	}

	total, _ := a.redis.Get(a.ctx, "analytics:total_requests").Int64()
	errors, _ := a.redis.Get(a.ctx, "analytics:errors").Int64()

	return map[string]interface{}{
		"total_requests": total,
		"total_errors":   errors,
		"error_rate":     float64(errors) / float64(total) * 100,
	}
}

// GetPopularEndpoints returns most hit endpoints
func (a *Analytics) GetPopularEndpoints(limit int) ([]map[string]interface{}, error) {
	if a == nil {
		return nil, nil
	}

	rows, err := a.db.Query(`
		SELECT endpoint, COUNT(*) as hits
		FROM events
		WHERE timestamp > datetime('now', '-7 days')
		GROUP BY endpoint
		ORDER BY hits DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]map[string]interface{}, 0)
	for rows.Next() {
		var endpoint string
		var hits int64
		if err := rows.Scan(&endpoint, &hits); err != nil {
			continue
		}
		results = append(results, map[string]interface{}{
			"endpoint": endpoint,
			"hits":     hits,
		})
	}

	return results, nil
}

// Close closes database connections
func (a *Analytics) Close() {
	if a != nil {
		if a.redis != nil {
			a.redis.Close()
		}
		if a.db != nil {
			a.db.Close()
		}
	}
}
