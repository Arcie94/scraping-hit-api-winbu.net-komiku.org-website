package service

import (
	"fmt"
	"io"
	"komiku-scraper/scraper/cache"
	"komiku-scraper/scraper/winbu"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type WinbuService struct {
	Client *winbu.WinbuClient
	Cache  *cache.Cache
}

func NewWinbuService(client *winbu.WinbuClient, c *cache.Cache) *WinbuService {
	return &WinbuService{Client: client, Cache: c}
}

func (s *WinbuService) FetchSearch(keyword string) ([]winbu.Anime, error) {
	cacheKey := fmt.Sprintf(cache.WinbuSearchKey, keyword)
	if val, found := s.Cache.Get(cacheKey); found {
		log.Printf("[Winbu] Cache HIT for search: %s", keyword)
		return val.([]winbu.Anime), nil
	}

	// Winbu search URL: https://winbu.net/?s=keyword
	// Replace space with +
	req, _ := http.NewRequest("GET", "https://winbu.net/?s="+keyword, nil)
	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	result, err := winbu.ParseSearch(doc)
	if err == nil {
		s.Cache.Set(cacheKey, result, cache.SearchTTL)
	}
	return result, err
}

func (s *WinbuService) FetchAndParseDetail(url string) (*winbu.AnimeDetail, error) {
	cacheKey := fmt.Sprintf(cache.WinbuDetailKey, url)
	if val, found := s.Cache.Get(cacheKey); found {
		log.Printf("[Winbu] Cache HIT for detail: %s", url)
		return val.(*winbu.AnimeDetail), nil
	}

	if !strings.HasPrefix(url, "http") {
		url = "https://winbu.net" + url
	}

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

	result, err := winbu.ParseAnimeDetail(doc)
	if err != nil {
		return nil, err
	}

	// If no episodes found (e.g. Movies), use the current page as the episode
	if len(result.Episodes) == 0 {
		result.Episodes = append(result.Episodes, winbu.Episode{
			Title:    "Full Movie / Watch",
			Endpoint: url,
		})
	}

	s.Cache.Set(cacheKey, result, cache.DetailTTL)
	return result, nil
}

func (s *WinbuService) FetchEpisode(url string) (*winbu.EpisodePageData, error) {
	cacheKey := fmt.Sprintf(cache.WinbuEpisodeKey, url)
	if val, found := s.Cache.Get(cacheKey); found {
		log.Printf("[Winbu] Cache HIT for episode: %s", url)
		return val.(*winbu.EpisodePageData), nil
	}

	if !strings.HasPrefix(url, "http") {
		url = "https://winbu.net" + url
	}

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

	result, err := winbu.ParseEpisodePage(doc)
	if err == nil {
		s.Cache.Set(cacheKey, result, cache.ChapterTTL)
	}
	return result, err
}

// FetchHomeData loads homepage data for top series, latest movies, latest anime, and genres
func (s *WinbuService) FetchHomeData() (*winbu.HomeData, error) {
	if val, found := s.Cache.Get(cache.WinbuHomeKey); found {
		log.Printf("[Winbu] Cache HIT for home data")
		return val.(*winbu.HomeData), nil
	}

	req, _ := http.NewRequest("GET", "https://winbu.net", nil)
	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	result, err := winbu.ParseHome(doc)
	if err == nil {
		s.Cache.Set(cache.WinbuHomeKey, result, cache.HomeTTL)
	}
	return result, err
}

func (s *WinbuService) ResolveStream(opt winbu.StreamOption) (string, error) {
	data := url.Values{}
	data.Set("action", "player_ajax")
	data.Set("post", opt.PostID)
	data.Set("nume", opt.Nume)
	data.Set("type", opt.Type)

	req, err := http.NewRequest("POST", "https://winbu.net/wp-admin/admin-ajax.php", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "https://winbu.net/")

	resp, err := s.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Log response details for debugging
	log.Printf("Stream Response Status: %d", resp.StatusCode)
	log.Printf("Stream Response Length: %d bytes", len(bodyBytes))
	log.Printf("Stream Response Preview: %s", string(bodyBytes[:min(200, len(bodyBytes))]))

	// Parse response to get iframe src using GoQuery for robustness
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(bodyBytes)))
	if err != nil {
		// Fallback to simple string parsing if goquery fails on fragment
		content := string(bodyBytes)
		if start := strings.Index(content, "src=\""); start != -1 {
			start += 5
			if end := strings.Index(content[start:], "\""); end != -1 {
				return content[start : start+end], nil
			}
		}
		return "", fmt.Errorf("could not parse response: %v", err)
	}

	// Try multiple selectors to find iframe src
	var src string
	var exists bool

	// Strategy 1: Direct iframe
	src, exists = doc.Find("iframe").Attr("src")
	if exists {
		log.Printf("Found iframe using selector: iframe")
		return src, nil
	}

	// Strategy 2: Iframe inside div
	src, exists = doc.Find("div iframe").Attr("src")
	if exists {
		log.Printf("Found iframe using selector: div iframe")
		return src, nil
	}

	// Strategy 3: Iframe with class
	src, exists = doc.Find("iframe[class]").Attr("src")
	if exists {
		log.Printf("Found iframe using selector: iframe[class]")
		return src, nil
	}

	// Strategy 4: Iframe with id
	src, exists = doc.Find("iframe[id]").Attr("src")
	if exists {
		log.Printf("Found iframe using selector: iframe[id]")
		return src, nil
	}

	// Strategy 5: Any iframe with src attribute
	doc.Find("iframe[src]").Each(func(i int, s *goquery.Selection) {
		if !exists {
			src, exists = s.Attr("src")
			if exists {
				log.Printf("Found iframe using selector: iframe[src]")
			}
		}
	})
	if exists {
		return src, nil
	}

	// Strategy 6: Check for data-src attribute
	src, exists = doc.Find("iframe[data-src]").Attr("data-src")
	if exists {
		log.Printf("Found iframe using data-src attribute")
		return src, nil
	}

	// Log response for debugging
	log.Printf("No iframe found. Response body preview (first 500 chars): %s", string(bodyBytes[:min(500, len(bodyBytes))]))

	return "", fmt.Errorf("no iframe src found in response after trying all strategies")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
