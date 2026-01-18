package service

import (
	"io"
	"komiku-scraper/scraper/komiku"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// KomikuService handles data fetching logic
type KomikuService struct {
	Client *komiku.KomikuClient
}

func NewKomikuService(client *komiku.KomikuClient) *KomikuService {
	return &KomikuService{Client: client}
}

func (s *KomikuService) FetchAndParseList(url string) ([]komiku.Manga, error) {
	log.Printf("[Komiku] Fetching manga list from: %s", url)
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := s.Client.Do(req)
	if err != nil {
		log.Printf("[Komiku] Error fetching list: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("[Komiku] List response status: %d", resp.StatusCode)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("[Komiku] Error parsing list HTML: %v", err)
		return nil, err
	}

	result, err := komiku.ParseMangaList(doc)
	if err == nil {
		log.Printf("[Komiku] Successfully parsed %d manga from list", len(result))
	}
	return result, err
}

func (s *KomikuService) FetchAndParseDetail(url string) (*komiku.MangaDetail, error) {
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
	}
	return result, err
}

func (s *KomikuService) FetchHomeData() (*komiku.HomeData, error) {
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

	return komiku.ParseHomeData(doc)
}

func (s *KomikuService) FetchChapterImages(url string) ([]komiku.ChapterImage, error) {
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
