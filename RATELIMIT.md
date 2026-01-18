# Rate Limiting Guide

## Overview

Both Anime and Manga APIs have rate limiting to prevent abuse.

## Rate Limits

### Default (IP-based)

- **45 requests per minute** per IP address
- Applies to all users without API key

### Premium (API Key-based)

- **450 requests per minute** (10x more)
- Requires valid API key

## Usage

### Without API Key

Just make requests normally:

```javascript
const response = await fetch("http://localhost:3001/api/v1/trending");
```

### With API Key

Include `X-API-Key` header:

```javascript
const response = await fetch("http://localhost:3001/api/v1/trending", {
  headers: {
    "X-API-Key": "your-api-key-here",
  },
});
```

## Response When Rate Limit Exceeded

```json
{
  "success": false,
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Too many requests. Please try again later."
  }
}
```

HTTP Status: `429 Too Many Requests`

## API Key Management

### Current Implementation

API keys are defined in `internal/api/middleware/ratelimit.go`:

```go
validKeys := map[string]bool{
    "demo-key-12345": true,
}
```

### Production Setup

For production, implement proper API key management:

#### Option 1: Environment Variables

```bash
# .env file
API_KEYS=key1,key2,key3
```

```go
keys := strings.Split(os.Getenv("API_KEYS"), ",")
```

#### Option 2: Database

```go
func IsValidAPIKey(key string) bool {
    // Query database
    var count int
    db.QueryRow("SELECT COUNT(*) FROM api_keys WHERE key = ? AND active = true", key).Scan(&count)
    return count > 0
}
```

#### Option 3: Redis

```go
func IsValidAPIKey(key string) bool {
    val, err := redisClient.Get(ctx, "apikey:"+key).Result()
    return err == nil && val == "1"
}
```

## Best Practices

### For API Users

1. **Cache responses** - Don't request same data repeatedly
2. **Use API key** - Get 10x more requests
3. **Handle 429 errors** - Implement exponential backoff
4. **Respect limits** - Don't try to bypass rate limiting

### For API Maintainers

1. **Monitor usage** - Track which IPs/keys use most requests
2. **Adjust limits** - Increase for trusted users
3. **Rotate keys** - Expire old API keys periodically
4. **Ban abusers** - Block IPs that try to bypass limits

## Example: Handling Rate Limits in Frontend

```typescript
async function fetchWithRetry(url: string, options = {}, retries = 3) {
  for (let i = 0; i < retries; i++) {
    const response = await fetch(url, options);

    if (response.status === 429) {
      // Rate limited, wait and retry
      const waitTime = Math.pow(2, i) * 1000; // Exponential backoff
      await new Promise((resolve) => setTimeout(resolve, waitTime));
      continue;
    }

    return response;
  }

  throw new Error("Rate limit exceeded after retries");
}

// Usage
const data = await fetchWithRetry("http://localhost:3001/api/v1/trending", {
  headers: { "X-API-Key": "your-key" },
});
```

## Testing Rate Limits

### Manual Test

```bash
# Send 50 requests rapidly
for i in {1..50}; do
  curl http://localhost:3001/api/v1/trending
done
```

After 45 requests, you should see 429 errors.

### With API Key

```bash
# Send 100 requests with API key
for i in {1..100}; do
  curl -H "X-API-Key: demo-key-12345" http://localhost:3001/api/v1/trending
done
```

Should handle 450 before limiting.

## Security Considerations

1. **Don't hardcode keys** - Use environment variables
2. **Use HTTPS** - Protect API keys in transit
3. **Rotate keys regularly** - Expire old keys
4. **Monitor for abuse** - Log suspicious patterns
5. **Rate limit by endpoint** - Different limits for expensive operations

## Future Enhancements

- [ ] Per-endpoint rate limits (e.g., stream resolution stricter)
- [ ] Usage tracking/analytics
- [ ] API key generation endpoint
- [ ] Subscription tiers (free/pro/enterprise)
- [ ] Rate limit headers (X-RateLimit-Remaining)
