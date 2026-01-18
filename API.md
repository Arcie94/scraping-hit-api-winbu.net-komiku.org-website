# API Documentation

## Overview

2 separate API servers untuk anime dan manga scraping dengan caching & image proxy.

##⚡ NEW FEATURES

### Image Proxy

Bypass CORS & hotlink protection:

```http
GET /api/v1/proxy/image?url=<base64_encoded_url>&size=small

Sizes:
- small: 150px width
- medium: 300px width
- large: 600px width
- (no size): original
```

Example (JavaScript):

```javascript
const imageUrl = "https://example.com/image.jpg";
const encoded = btoa(imageUrl);
const proxyUrl = `http://localhost:3001/api/v1/proxy/image?url=${encoded}&size=medium`;
```

### Response Caching

Frequently accessed endpoints are cached:

- Homepage data: 5 minutes
- Genre lists: 10 minutes
- Reduces server load & improves response time

---

## Anime API (Port 3001)

### Base URL

```
http://localhost:3001/api/v1
```

### Endpoints

#### Image Proxy

```http
GET /proxy/image?url=<base64_url>&size=small|medium|large
```

#### Search Anime

```http
GET /search?q=naruto
```

#### Top Series (Cached 5min)

```http
GET /top-series
```

#### Top Movies (Cached 5min)

```http
GET /top-movies
```

#### Latest Movies (Cached 5min)

```http
GET /latest-movies
```

#### Latest Anime (Cached 5min)

```http
GET /latest-anime
```

#### Drama (Cached 5min)

```http
GET /drama
```

#### Genres (Cached 10min)

```http
GET /genres
```

#### Anime Detail

```http
GET /anime/:endpoint
```

#### Episode Data

```http
GET /episode/:endpoint
```

#### Resolve Stream

```http
POST /stream/resolve
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

#### Image Proxy

```http
GET /proxy/image?url=<base64_url>&size=small|medium|large
```

#### Search Manga

```http
GET /search?q=one+piece
```

#### Trending (Cached 5min)

```http
GET /trending
```

#### Popular (Cached 5min)

```http
GET /popular
```

#### Genres (Cached 10min)

```http
GET /genres
```

#### Manga Detail

```http
GET /manga/:endpoint
```

#### Chapter Images

```http
GET /chapter/:endpoint
```

#### Recommendations

```http
GET /recommendations/:endpoint
```

---

## Usage Examples

### Astro/React - Using Image Proxy

```typescript
// Utility function to get proxied image URL
function getProxiedImageUrl(originalUrl: string, size?: 'small' | 'medium' | 'large') {
  const encoded = btoa(originalUrl);
  const sizeParam = size ? `&size=${size}` : '';
  return `http://your-server:3001/api/v1/proxy/image?url=${encoded}${sizeParam}`;
}

// In component
const manga = await fetch('http://your-server:3002/api/v1/trending').then(r => r.json());

// Use proxied images
<img src={getProxiedImageUrl(manga.data[0].image, 'medium')} alt={manga.data[0].title} />
```

### Performance Benefits

- **Cached responses**: 5-10 min cache reduces redundant scraping
- **Smaller images**: Thumbnails load faster
- **Single origin**: No CORS issues

---

## Response Format

### Success

```json
{
  "success": true,
  "data": [...],
  "meta": {
    "total": 100
  }
}
```

### Error

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Description"
  }
}
```

---

## Deployment (see DOCKER.md)

### Quick Start

```bash
docker-compose up -d
```

### Production

- Update CORS origins in both API files
- Use Nginx reverse proxy
- Enable HTTPS

---

## Performance Notes

### Cache TTL

- Homepage endpoints: 5 minutes
- Genre lists: 10 minutes
- Detail/chapter: No cache (always fresh)

### Image Proxy

- Caches images: 24 hours (browser)
- Resize quality: 85% JPEG
- Supports PNG & JPEG

---

## Changelog

### v1.1 (Current)

- ✅ Image proxy with resize
- ✅ Response caching
- ✅ Better error messages

### v1.0

- Initial release
- Basic scraping endpoints
