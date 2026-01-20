package komiku

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ParseMangaList(doc *goquery.Document) ([]Manga, error) {
	var mangas []Manga

	// DEBUG: Log total elements found
	bgeCount := doc.Find(".bge").Length()
	ls2Count := doc.Find("article.ls2").Length()
	log.Printf("[Parser DEBUG] Found .bge elements: %d, article.ls2 elements: %d", bgeCount, ls2Count)

	// Selector untuk halaman search/daftar komik (.bge adalah wrapper umum untuk item list)
	doc.Find(".bge, article.ls2").Each(func(i int, s *goquery.Selection) {
		// Coba ambil dari struktur .bge dulu
		title := strings.TrimSpace(s.Find(".kan h3").Text())
		endpoint, _ := s.Find(".kan a").Attr("href")
		thumb, _ := s.Find(".bgei img").Attr("src")

		// DEBUG: Log setiap item yang ditemukan
		log.Printf("[Parser DEBUG] Item %d - .bge structure: title='%s', endpoint='%s', thumb='%s'", i, title, endpoint, thumb)

		// Jika tidak ketemu, coba struktur ls2
		if title == "" {
			title = strings.TrimSpace(s.Find(".ls2j h3 a").Text())
			endpoint, _ = s.Find(".ls2j h3 a").Attr("href")
			thumb, _ = s.Find(".ls2v img").Attr("src")
			// Coba data-src untuk lazy load
			if thumb == "" {
				thumb, _ = s.Find(".ls2v img").Attr("data-src")
			}
			log.Printf("[Parser DEBUG] Item %d - .ls2 structure: title='%s', endpoint='%s', thumb='%s'", i, title, endpoint, thumb)
		}

		// Bersihkan thumbnail dari parameter query string dan lazy placeholder
		if strings.Contains(thumb, "?") {
			parts := strings.Split(thumb, "?")
			thumb = parts[0]
		}
		if strings.Contains(thumb, "lazy.jpg") {
			thumb = ""
		}

		if title != "" && endpoint != "" {
			log.Printf("[Parser DEBUG] ✓ Adding manga: %s", title)
			mangas = append(mangas, Manga{
				Title:    title,
				Endpoint: endpoint,
				Thumb:    thumb,
			})
		} else {
			log.Printf("[Parser DEBUG] ✗ Skipping item %d - missing title or endpoint", i)
		}
	})

	log.Printf("[Parser DEBUG] Total manga parsed: %d", len(mangas))
	return mangas, nil
}

func ParseMangaDetail(doc *goquery.Document) (*MangaDetail, error) {
	var detail MangaDetail

	// Detail Utama
	detail.Title = strings.TrimSpace(doc.Find("#Judul h1").Text())
	detail.Thumb, _ = doc.Find(".ims img").Attr("src")
	if strings.Contains(detail.Thumb, "?") {
		detail.Thumb = strings.Split(detail.Thumb, "?")[0]
	}
	detail.Synopsis = strings.TrimSpace(doc.Find(".desc").Text())
	detail.Description = detail.Synopsis // UI Compatibility

	// Metadata Table
	doc.Find(".inftable tr").Each(func(i int, s *goquery.Selection) {
		label := strings.ToLower(strings.TrimSpace(s.Find("td").First().Text()))
		value := strings.TrimSpace(s.Find("td").Last().Text())

		if strings.Contains(label, "pengarang") {
			detail.Authors = []string{value}
		}
		if strings.Contains(label, "status") {
			detail.Status = value
		}
	})

	// Genre
	doc.Find(".genre li a").Each(func(i int, s *goquery.Selection) {
		detail.Genres = append(detail.Genres, strings.TrimSpace(s.Text()))
	})

	// Chapter List
	doc.Find("table#Daftar_Chapter tr").Each(func(i int, s *goquery.Selection) {
		// Skip header rows by checking for th
		if s.Find("th").Length() > 0 {
			return
		}

		titleEl := s.Find("td.judulseries a")
		title := strings.TrimSpace(titleEl.Text())
		endpoint, _ := titleEl.Attr("href")
		date := strings.TrimSpace(s.Find("td.tanggalseries").Text())

		// Bersihkan whitespace berlebih di tanggal
		date = strings.Join(strings.Fields(date), " ")

		if title != "" && endpoint != "" {
			detail.Chapters = append(detail.Chapters, ChapterLink{
				Title:        title,
				Endpoint:     endpoint,
				DateUploaded: date,
			})
		}
	})

	return &detail, nil
}

func ParseHomeData(doc *goquery.Document) (*HomeData, error) {
	var data HomeData

	// Popular Manga - dari section #Komik_Hot_Manga
	doc.Find("#Komik_Hot_Manga article.ls2").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find(".ls2j h3 a").Text())
		endpoint, _ := s.Find(".ls2j h3 a").Attr("href")
		thumb, _ := s.Find(".ls2v img").Attr("src")

		// Coba data-src untuk lazy load
		if thumb == "" || strings.Contains(thumb, "lazy.jpg") {
			thumb, _ = s.Find(".ls2v img").Attr("data-src")
		}

		// Bersihkan
		if strings.Contains(thumb, "?") {
			thumb = strings.Split(thumb, "?")[0]
		}

		if title != "" && endpoint != "" {
			data.Popular = append(data.Popular, Manga{
				Title: title, Endpoint: endpoint, Thumb: thumb,
			})
		}
	})

	// Latest Manga - dari section #Terbaru (menggunakan struktur yang berbeda)
	// Catatan: Perlu mencari struktur ls4 atau daftar terbaru
	doc.Find("#Terbaru .ls4").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find(".ls4j h4 a").Text())
		endpoint, _ := s.Find(".ls4j h4 a").Attr("href")
		thumb, _ := s.Find(".ls4v img").Attr("src")

		if thumb == "" || strings.Contains(thumb, "lazy.jpg") {
			thumb, _ = s.Find(".ls4v img").Attr("data-src")
		}

		if strings.Contains(thumb, "?") {
			thumb = strings.Split(thumb, "?")[0]
		}

		if title != "" && endpoint != "" {
			data.Latest = append(data.Latest, Manga{
				Title: title, Endpoint: endpoint, Thumb: thumb,
			})
		}
	})

	// Ensure Trending is initialized to avoid nil
	data.Trending = []Manga{}

	return &data, nil
}

func ParseChapterImages(htmlContent string) ([]ChapterImage, error) {
	var images []ChapterImage
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	doc.Find("#Baca_Komik img").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		// Handling lazy load jika src adalah placeholder/kosong tapi ada data-src
		if src == "" || strings.Contains(src, "lazy.jpg") {
			src, _ = s.Attr("data-src")
		}

		if src != "" && !strings.Contains(src, "lazy.jpg") {
			images = append(images, ChapterImage{
				URL: src,
			})
		}
	})

	return images, nil
}

func ParseRecommendations(doc *goquery.Document) ([]Manga, error) {
	var recommendations []Manga

	// Selector rekomendasi dari sidebar atau bagian "Mirip"
	// Di chapter ada di #Terbaru .ls8
	doc.Find("#Terbaru .ls8, .ls8").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find(".ls8j h3 a").Text())
		if title == "" {
			title = strings.TrimSpace(s.Find("h3 a").Text())
		}
		endpoint, _ := s.Find("a").Attr("href")
		thumb, _ := s.Find("img").Attr("src")

		// Handle lazy load
		if thumb == "" || strings.Contains(thumb, "lazy.jpg") {
			thumb, _ = s.Find("img").Attr("data-src")
		}

		if strings.Contains(thumb, "?") {
			thumb = strings.Split(thumb, "?")[0]
		}

		if title != "" && endpoint != "" {
			recommendations = append(recommendations, Manga{
				Title: title, Endpoint: endpoint, Thumb: thumb,
			})
		}
	})

	return recommendations, nil
}

func ParseGenreList(doc *goquery.Document) ([]Genre, error) {
	var genres []Genre

	// Coba ambil dari menu navigasi atau list genre di sidebar jika ada.
	// Struktur umum menu genre: .genre li a atau ul.nav li a[href*="/genre/"]

	doc.Find("ul.genre li a, a[href*='/genre/']").Each(func(i int, s *goquery.Selection) {
		name := strings.TrimSpace(s.Text())
		endpoint, _ := s.Attr("href")

		// Filter nama yang valid dan bukan link "Genre" parent saja
		if name != "" && endpoint != "" && strings.Contains(endpoint, "/genre/") {
			genres = append(genres, Genre{
				Name:     name,
				Endpoint: endpoint,
			})
		}
	})

	return genres, nil
}
