# Docker Deployment Guide

## Quick Start

### Build & Run dengan Docker Compose

```bash
# Build images dan run containers
docker-compose up -d

# Check logs
docker-compose logs -f

# Stop containers
docker-compose down
```

---

## Manual Docker Commands

### Build Images

```bash
# Anime API
docker build -f Dockerfile.anime -t anime-api:latest .

# Manga API
docker build -f Dockerfile.manga -t manga-api:latest .
```

### Run Containers

```bash
# Anime API
docker run -d \
  --name anime-api \
  -p 3001:3001 \
  --restart unless-stopped \
  anime-api:latest

# Manga API
docker run -d \
  --name manga-api \
  -p 3002:3002 \
  --restart unless-stopped \
  manga-api:latest
```

---

## Integration dengan Website Docker

### Option 1: Shared Network

Jika website Astro Anda di Docker network yang sama:

```yaml
# Tambahkan ke docker-compose.yml website Anda
networks:
  default:
    external:
      name: scraper-network
```

Atau tambahkan API ke network website:

```bash
docker network connect website-network anime-api
docker network connect website-network manga-api
```

Akses dari Astro:

```typescript
// Gunakan container name sebagai hostname
const response = await fetch("http://anime-api:3001/api/v1/trending");
```

### Option 2: Nginx Reverse Proxy

Tambahkan di Nginx config (bisa container atau host):

```nginx
# Anime API
location /api/anime/ {
    proxy_pass http://anime-api:3001/api/v1/;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
}

# Manga API
location /api/manga/ {
    proxy_pass http://manga-api:3002/api/v1/;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
}
```

---

## Production Configuration

### Update CORS

Edit kedua file API (`cmd/anime-api/main.go` dan `cmd/manga-api/main.go`):

```go
app.Use(cors.New(cors.Config{
    AllowOrigins: "https://yourdomain.com", // Ganti dengan domain Anda
    AllowMethods: "GET,POST",
    AllowHeaders: "Origin, Content-Type, Accept",
}))
```

### Environment Variables (Optional)

Tambahkan di `docker-compose.yml`:

```yaml
services:
  anime-api:
    environment:
      - ALLOWED_ORIGIN=https://yourdomain.com
      - PORT=3001
```

---

## Deployment ke Proxmox

### 1. Transfer Files

```bash
# Zip project
tar -czf scraper.tar.gz .

# SCP ke Proxmox
scp scraper.tar.gz user@proxmox-ip:/path/to/deploy/

# Extract di server
ssh user@proxmox-ip
cd /path/to/deploy
tar -xzf scraper.tar.gz
```

### 2. Build & Run

```bash
docker-compose up -d --build
```

### 3. Verify

```bash
# Check containers running
docker ps

# Test endpoints
curl http://localhost:3001/
curl http://localhost:3002/

# Check logs
docker-compose logs -f anime-api
docker-compose logs -f manga-api
```

---

## Monitoring & Maintenance

### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f anime-api

# Last 100 lines
docker-compose logs --tail=100 anime-api
```

### Restart Services

```bash
# Restart all
docker-compose restart

# Restart specific
docker-compose restart anime-api
```

### Update & Redeploy

```bash
# Pull latest code
git pull

# Rebuild & restart
docker-compose up -d --build
```

### Clean Up

```bash
# Stop and remove containers
docker-compose down

# Remove images
docker rmi anime-api:latest manga-api:latest

# Clean build cache
docker builder prune
```

---

## Health Checks

Both APIs include health checks:

- Interval: 30s
- Timeout: 10s
- Start period: 40s

Check health status:

```bash
docker ps
# Look for (healthy) status
```

---

## Troubleshooting

### Container won't start

```bash
# Check logs
docker logs anime-api

# Inspect container
docker inspect anime-api
```

### Port conflicts

```bash
# Change ports in docker-compose.yml
ports:
  - "8001:3001"  # Host:Container
```

### Network issues

```bash
# List networks
docker network ls

# Inspect network
docker network inspect scraper-network
```

### Can't access from Astro

1. Check containers are on same network
2. Use container name as hostname
3. Check CORS configuration
4. Verify firewall rules

---

## Size Optimization

Current images use multi-stage builds:

- Builder stage: Full Go toolchain
- Runtime stage: Alpine (~5MB)
- Final image: ~20-30MB

Check image sizes:

```bash
docker images | grep -E 'anime-api|manga-api'
```

---

## Security Best Practices

1. **Don't expose ports publicly** - Use Nginx reverse proxy
2. **Update CORS** - Set specific domain, not `*`
3. **Add rate limiting** - Consider adding middleware
4. **Use secrets** - For API keys (if needed later)
5. **Regular updates** - Keep base images updated

```bash
# Update base images
docker-compose pull
docker-compose up -d --build
```
