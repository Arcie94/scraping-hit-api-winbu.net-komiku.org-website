# API Documentation

## Overview

2 separate API servers untuk anime dan manga scraping.

---

## Anime API (Port 3001)

### Base URL

```
http://localhost:3001/api/v1
```

### Endpoints

#### Search Anime

```http
GET /search?q=naruto
```

Response:

```json
{
  "success": true,
  "data": [
    {
      "title": "Naruto",
      "endpoint": "/anime/naruto",
      "image": "...",
      "rating": "8.5",
      "status": "Completed",
      "type": "TV"
    }
  ],
  "meta": { "total": 10 }
}
```

#### Top Series

```http
GET /top-series
```

#### Top Movies

```http
GET /top-movies
```

#### Latest Movies

```http
GET /latest-movies
```

#### Latest Anime

```http
GET /latest-anime
```

#### Drama (International)

```http
GET /drama
```

#### Genres

```http
GET /genres
```

#### Anime Detail

```http
GET /anime/:endpoint
# Example: GET /anime/anime/naruto
```

#### Episode Data

```http
GET /episode/:endpoint
# Example: GET /episode/episode/naruto-episode-1
```

#### Resolve Stream

```http
POST /stream/resolve
Content-Type: application/json

{
  "postId": "123",
  "nume": "1",
  "type": "anime"
}
```

---

## Manga API (Port 3002)

### Base URL

```
http://localhost:3002/api/v1
```

### Endpoints

#### Search Manga

```http
GET /search?q=one+piece
```

#### Trending

```http
GET /trending
```

#### Popular

```http
GET /popular
```

#### Genres

```http
GET /genres
```

#### Manga Detail

```http
GET /manga/:endpoint
# Example: GET /manga/manga/one-piece
```

#### Chapter Images

```http
GET /chapter/:endpoint
# Example: GET /chapter/ch/one-piece-chapter-1
```

#### Recommendations

```http
GET /recommendations/:endpoint
```

---

## Running the Servers

### Development

```bash
# Terminal 1 - Anime API
go run cmd/anime-api/main.go

# Terminal 2 - Manga API
go run cmd/manga-api/main.go
```

### Production

```bash
# Build
go build -o anime-api ./cmd/anime-api
go build -o manga-api ./cmd/manga-api

# Run
./anime-api &
./manga-api &
```

---

## Integration with Astro

### Example: Fetch Trending Manga

```typescript
// In your Astro component
const response = await fetch("http://your-server:3002/api/v1/trending");
const data = await response.json();

if (data.success) {
  const mangaList = data.data;
  // Render manga list
}
```

### Example: Search Anime

```typescript
const q = "naruto";
const response = await fetch(`http://your-server:3001/api/v1/search?q=${q}`);
const data = await response.json();
```

---

## Deployment on Proxmox

### Option 1: Systemd Services

Create `/etc/systemd/system/anime-api.service`:

```ini
[Unit]
Description=Anime API Server
After=network.target

[Service]
Type=simple
User=your-user
WorkingDirectory=/path/to/scraper
ExecStart=/path/to/scraper/anime-api
Restart=always

[Install]
WantedBy=multi-user.target
```

Create `/etc/systemd/system/manga-api.service`:

```ini
[Unit]
Description=Manga API Server
After=network.target

[Service]
Type=simple
User=your-user
WorkingDirectory=/path/to/scraper
ExecStart=/path/to/scraper/manga-api
Restart=always

[Install]
WantedBy=multi-user.target
```

Enable & start:

```bash
sudo systemctl enable anime-api manga-api
sudo systemctl start anime-api manga-api
```

### Option 2: Nginx Reverse Proxy

Add to your Nginx config:

```nginx
# Anime API
location /api/anime/ {
    proxy_pass http://localhost:3001/api/v1/;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
}

# Manga API
location /api/manga/ {
    proxy_pass http://localhost:3002/api/v1/;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
}
```

Then access via:

- `https://yourdomain.com/api/anime/search?q=naruto`
- `https://yourdomain.com/api/manga/trending`

---

## CORS Configuration

Update `AllowOrigins` in both API servers for production:

```go
app.Use(cors.New(cors.Config{
    AllowOrigins: "https://yourdomain.com",
    AllowMethods: "GET,POST",
    AllowHeaders: "Origin, Content-Type, Accept",
}))
```

---

## Error Handling

All errors return:

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable message"
  }
}
```

Common error codes:

- `INVALID_QUERY` - Query parameter missing
- `SEARCH_FAILED` - Search error
- `FETCH_FAILED` - Failed to fetch data
- `NOT_FOUND` - Resource not found
