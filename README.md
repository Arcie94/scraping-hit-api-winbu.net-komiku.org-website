# Komiku & Winbu Scraper

Scraper sederhana untuk membaca manga dari **Komiku.org** dan streaming anime dari **Winbu.net** dengan CLI interaktif.

## 🌟 Features

### Komiku (Manga)

- 🔍 **Search Manga** - Cari manga berdasarkan kata kunci
- 📖 **Detail Manga** - Lihat info lengkap + daftar chapter
- 🔥 **Manga Trending** - Lihat manga yang sedang trending (bisa pilih)
- ⭐ **Manga Populer** - Lihat manga populer (bisa pilih)
- 📷 **Read Chapter** - Ekstrak semua gambar chapter untuk dibaca
- 💡 **Recommendations** - Dapatkan rekomendasi dari halaman chapter
- 🏷️ **List Genre** - Browse genre yang tersedia

### Winbu (Anime)

- 🔍 **Search Anime** - Cari anime berdasarkan kata kunci
- 🏆 **Top 10 Anime** - Top 10 series anime
- 🎬 **Top 10 Film** - Top 10 film anime
- 🆕 **Film Terbaru** - Film anime terbaru
- 📺 **Anime/Donghua Terbaru** - Update terbaru anime & donghua
- 🌏 **Drama** - Series dari Jepang/Korea/China/Barat
- 🏷️ **List Genre** - Browse genre anime
- **Stream Video** - Resolusi video otomatis dengan 6 fallback strategies

## 📦 Installation

### Prerequisites

- Go 1.20+
- Internet connection

### Setup

```bash
# Clone repository
git clone <repo-url>
cd "260118 Sniffing & Hit API Manga & Anime"

# Install dependencies
go mod tidy

# Build
go build -o scraper.exe

# Run
./scraper.exe
```

## 🚀 Usage

### Quick Start

```bash
./scraper.exe
```

Menu utama akan muncul:

```
=== MAIN MENU ===
1. KOMIKU (Manga Scraper)
2. WINBU (Anime Scraper)
0. Exit
```

### Contoh Flow - Membaca Manga

1. Pilih `1` (KOMIKU)
2. Pilih `3` (Manga Trending)
3. Pilih nomor manga yang ingin dibaca
4. Lihat daftar chapter
5. Pilih chapter untuk ekstrak gambar
6. Gambar akan dibuka di browser

### Contoh Flow - Streaming Anime

1. Pilih `2` (WINBU)
2. Pilih `1` (Search Anime)
3. Ketik nama anime (contoh: "one piece")
4. Pilih anime dari hasil pencarian
5. Pilih episode
6. Pilih server streaming
7. URL video akan ditampilkan

## 🛠️ Development

### Project Structure

```
├── cmd/
│   └── main.go              # Entry point
├── internal/
│   ├── service/             # Business logic layer
│   │   ├── service.go       # Komiku service
│   │   └── winbu.go         # Winbu service
│   └── ui/
│       └── menu.go          # CLI menu
├── scraper/
│   ├── komiku/              # Komiku parsers
│   │   ├── client.go
│   │   ├── parser.go
│   │   └── types.go
│   └── winbu/               # Winbu parsers
│       ├── client.go
│       ├── parser_home.go
│       ├── parser_detail.go
│       ├── parser_stream.go
│       └── types.go
├── scripts/
│   └── verify_home.go       # Test script
└── go.mod
```

### Running Tests

```bash
# Verify homepage parsing
go run scripts/verify_home.go
```

### Logging

All operations are logged with prefixes:

- `[Komiku]` - Manga operations
- `[Winbu]` - Anime operations

Logs include:

- HTTP response status
- Response size
- Parse results
- Strategy used (for stream resolution)

## 🔧 Technical Details

### Winbu Stream Resolution

Uses 6 different strategies to find video iframe:

1. Direct `iframe`
2. `div iframe`
3. `iframe[class]`
4. `iframe[id]`
5. `iframe[src]`
6. `iframe[data-src]` (lazy-loading)

### Komiku Image Extraction

- Extracts images from `div#Baca_Komik img`
- Opens all images in browser for easy reading
- Preserves page order

## 📝 Notes

- **Rate Limiting**: Please use responsibly, don't spam requests
- **Legal**: For educational purposes only
- **Maintenance**: Selectors may break if websites update their HTML structure

## 🐛 Troubleshooting

### "No iframe src found in response"

- Try different server (MEGA, HyDRAX usually work better)
- Check logs for response preview
- Server might be down temporarily

### "Tidak ada gambar ditemukan"

- Chapter page structure might have changed
- Check if chapter exists on website
- Verify URL is correct

### Build Errors

```bash
# Clean and rebuild
go clean
go mod tidy
go build
```

## 📄 License

Educational use only. Respect the original content creators and website owners.

## 🙏 Credits

- [Komiku.org](https://komiku.org) - Manga source
- [Winbu.net](https://winbu.net) - Anime source
- [goquery](https://github.com/PuerkitoBio/goquery) - HTML parsing library
