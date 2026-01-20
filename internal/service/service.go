package service

import (
	"fmt"
	"io"
	"komiku-scraper/scraper/cache"
	"komiku-scraper/scraper/komiku"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// KomikuService handles data fetching logic
type KomikuService struct {
	Client *komiku.KomikuClient
	Cache  *cache.Cache
}

func NewKomikuService(client *komiku.KomikuClient, c *cache.Cache) *KomikuService {
	return &KomikuService{Client: client, Cache: c}
}

func (s *KomikuService) FetchAndParseList(url string) ([]komiku.Manga, error) {
	cacheKey := fmt.Sprintf(cache.KomikuSearchKey, url)
	if val, found := s.Cache.Get(cacheKey); found {
		log.Printf("[Komiku] Cache HIT for list: %s", url)
		return val.([]komiku.Manga), nil
	}

	log.Printf("[Komiku] Fetching manga list from: %s", url)
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := s.Client.Do(req)
	if err != nil {
		log.Printf("[Komiku] Error fetching list: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("[Komiku] List response status: %d", resp.StatusCode)

	// Read body for debug logging
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[Komiku] Error reading body: %v", err)
		return nil, err
	}

	// Debug: Log first 2000 chars of HTML to see actual content
	htmlPreview := string(bodyBytes)
	if len(htmlPreview) > 2000 {
		htmlPreview = htmlPreview[:2000]
	}
	log.Printf("[Komiku] HTML Preview (first 2000 chars): %s", htmlPreview)

	// Create reader from bytes
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(bodyBytes)))
	if err != nil {
		log.Printf("[Komiku] Error parsing list HTML: %v", err)
		return nil, err
	}

	result, err := komiku.ParseMangaList(doc)
	if err == nil {
		log.Printf("[Komiku] Successfully parsed %d manga from list", len(result))
		s.Cache.Set(cacheKey, result, cache.SearchTTL)
	} else {
		log.Printf("[Komiku] Parse error: %v", err)
	}
	return result, err
}

func (s *KomikuService) FetchAndParseDetail(url string) (*komiku.MangaDetail, error) {
	cacheKey := fmt.Sprintf(cache.KomikuDetailKey, url)
	if val, found := s.Cache.Get(cacheKey); found {
		log.Printf("[Komiku] Cache HIT for detail: %s", url)
		return val.(*komiku.MangaDetail), nil
	}

	log.Printf("[Komiku] Fetching manga detail from: %s", url)
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := s.Client.Do(req)
	if err != nil {
		log.Printf("[Komiku] Error fetching detail: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("[Komiku] Detail response status: %d", resp.StatusCode)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("[Komiku] Error parsing detail HTML: %v", err)
		return nil, err
	}

	result, err := komiku.ParseMangaDetail(doc)
	if err == nil && result != nil {
		log.Printf("[Komiku] Successfully parsed manga: %s (%d chapters)", result.Title, len(result.Chapters))
		s.Cache.Set(cacheKey, result, cache.DetailTTL)
	}
	return result, err
}

func (s *KomikuService) FetchHomeData() (*komiku.HomeData, error) {
	if val, found := s.Cache.Get(cache.KomikuHomeKey); found {
		log.Printf("[Komiku] Cache HIT for home data")
		return val.(*komiku.HomeData), nil
	}

	req, _ := http.NewRequest("GET", "https://komiku.org/", nil)
	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	result, err := komiku.ParseHomeData(doc)
	if err == nil {
		s.Cache.Set(cache.KomikuHomeKey, result, cache.HomeTTL)
	}
	return result, err
}

func (s *KomikuService) FetchChapterImages(url string) ([]komiku.ChapterImage, error) {
	cacheKey := fmt.Sprintf(cache.KomikuChapterKey, url)
	if val, found := s.Cache.Get(cacheKey); found {
		log.Printf("[Komiku] Cache HIT for chapter: %s", url)
		return val.([]komiku.ChapterImage), nil
	}

	log.Printf("[Komiku] Fetching chapter images from: %s", url)
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := s.Client.Do(req)
	if err != nil {
		log.Printf("[Komiku] Error fetching chapter: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("[Komiku] Chapter response status: %d", resp.StatusCode)

	// Chapter parser takes string body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[Komiku] Error reading chapter body: %v", err)
		return nil, err
	}

	log.Printf("[Komiku] Chapter response size: %d bytes", len(bodyBytes))

	result, err := komiku.ParseChapterImages(string(bodyBytes))
	if err == nil {
		log.Printf("[Komiku] Successfully parsed %d images from chapter", len(result))
		s.Cache.Set(cacheKey, result, cache.ChapterTTL)
	} else {
		log.Printf("[Komiku] Error parsing chapter images: %v", err)
	}
	return result, err
}

func (s *KomikuService) FetchRecommendations(url string) ([]komiku.Manga, error) {
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return komiku.ParseRecommendations(doc)
}

func (s *KomikuService) FetchGenreList() ([]komiku.Genre, error) {
	req, _ := http.NewRequest("GET", "https://komiku.org/", nil)
	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return komiku.ParseGenreList(doc)
}
