# API Documentation v1.1

## Base URLs

- **Anime API**: `http://localhost:3001`
- **Manga API**: `http://localhost:3002`

---

## NEW in v1.1

### Health Check

Monitor API and scraper status:

```http
GET /health
```

Response:

```json
{
  "status": "healthy",
  "timestamp": "2026-01-18T22:30:00Z",
  "uptime": 86400.5,
  "scraper": {
    "name": "Winbu",
    "status": "ok",
    "url": "https://winbu.net"
  },
  "requests_served": 15432
}
```

### Batch Endpoints

Fetch multiple items in one request:

#### Anime Batch

```http
POST /api/v1/batch/anime
Content-Type: application/json

{
  "endpoints": [
    "/anime/naruto",
    "/anime/one-piece"
  ]
}
```

Response:

```json
{
  "success": true,
  "data": [{...}, {...}],
  "meta": {
    "total": 2,
    "requested": 2,
    "failed": 0
  }
}
```

#### Manga Batch

```http
POST /api/v1/batch/manga
Content-Type: application/json

{
  "endpoints": [
    "/manga/one-piece",
    "/manga/naruto"
  ]
}
```

**Limits:**

- Maximum 10 items per batch
- Failed items return in `errors` array
- Partial success supported

---

## Image Proxy

Bypass CORS & resize images:

```http
GET /api/v1/proxy/image?url=<base64_url>&size=small|medium|large
```

Sizes:

- `small`: 150px width
- `medium`: 300px width
- `large`: 600px width

Example:

```javascript
const imageUrl = "https://example.com/image.jpg";
const encoded = btoa(imageUrl);
const proxyUrl = `http://localhost:3001/api/v1/proxy/image?url=${encoded}&size=medium`;
```

---

## Anime API Endpoints

### Health & Info

```http
GET /           # Simple info
GET /health     # Enhanced health check
```

### Image & Batch

```http
GET /api/v1/proxy/image?url=<base64>&size=<size>
POST /api/v1/batch/anime
```

### Search & Lists

```http
GET /api/v1/search?q=naruto
GET /api/v1/top-series        # Cached 5min
GET /api/v1/top-movies        # Cached 5min
GET /api/v1/latest-movies     # Cached 5min
GET /api/v1/latest-anime      # Cached 5min
GET /api/v1/drama             # Cached 5min
GET /api/v1/genres            # Cached 10min
```

### Details & Stream

```http
GET /api/v1/anime/:endpoint
GET /api/v1/episode/:endpoint
POST /api/v1/stream/resolve
```

---

## Manga API Endpoints

### Health & Info

```http
GET /           # Simple info
GET /health     # Enhanced health check
```

### Image & Batch

```http
GET /api/v1/proxy/image?url=<base64>&size=<size>
POST /api/v1/batch/manga
```

### Search & Lists

```http
GET /api/v1/search?q=one+piece
GET /api/v1/trending         # Cached 5min
GET /api/v1/popular          # Cached 5min
GET /api/v1/genres           # Cached 10min
```

### Details & Chapters

```http
GET /api/v1/manga/:endpoint
GET /api/v1/chapter/:endpoint
GET /api/v1/recommendations/:endpoint
```

---

## Features

### Caching

- Homepage data: 5 minutes
- Genre lists: 10 minutes
- Detail/chapter: No cache (always fresh)

### Rate Limiting

- Default: 45 req/min per IP
- With API key: 450 req/min
- Header: `X-API-Key: your-key`

### Error Responses

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable message"
  }
}
```

Common codes:

- `RATE_LIMIT_EXCEEDED` - Too many requests
- `INVALID_QUERY` - Missing/invalid parameters
- `FETCH_FAILED` - Scraping error
- `BATCH_TOO_LARGE` - Batch size > 10

---

## Usage Examples

### JavaScript/TypeScript

```typescript
// Health check
const health = await fetch("http://localhost:3001/health").then((r) =>
  r.json(),
);
console.log(`Uptime: ${health.uptime}s, Requests: ${health.requests_served}`);

// Batch fetch
const batch = await fetch("http://localhost:3001/api/v1/batch/anime", {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({
    endpoints: ["/anime/naruto", "/anime/one-piece"],
  }),
}).then((r) => r.json());

// Image proxy
function getProxiedImage(url, size = "medium") {
  const encoded = btoa(url);
  return `http://localhost:3001/api/v1/proxy/image?url=${encoded}&size=${size}`;
}
```

### Astro Component

```astro
---
const health = await fetch('http://anime-api:3001/health').then(r => r.json());
const trending = await fetch('http://manga-api:3002/api/v1/trending').then(r => r.json());
---

<div>
  <p>API Uptime: {Math.floor(health.uptime / 3600)}h</p>
  <p>Requests Served: {health.requests_served.toLocaleString()}</p>

  {trending.data.map(manga => (
    <img src={getProxiedImage(manga.image, 'small')} alt={manga.title} />
  ))}
</div>
```

---

## Deployment

See [DOCKER.md](DOCKER.md) for Docker deployment guide.

See [RATELIMIT.md](RATELIMIT.md) for rate limiting details.

---

## Changelog

### v1.1 (Current)

- ✅ Health check with scraper status
- ✅ Batch endpoints (anime & manga)
- ✅ Request counter
- ✅ Image proxy with resize
- ✅ Response caching (5-10min)
- ✅ Rate limiting (45/450 req/min)

### v1.0

- Initial release
- Basic scraping endpoints
