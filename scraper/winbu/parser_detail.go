package winbu

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ParseAnimeDetail(doc *goquery.Document) (*AnimeDetail, error) {
	detail := &AnimeDetail{
		Metadata: make(map[string]string),
	}

	container := doc.Find(".movies-list.movies-list-full .t-item")

	// Title
	detail.Title = strings.TrimSpace(container.Find(".mli-info .judul").Text())
	if detail.Title == "" {
		detail.Title = strings.TrimSpace(doc.Find("h1.titless").Text())
	}

	// Thumb
	detail.Thumb = container.Find(".ml-mask .mli-thumb-box img").AttrOr("src", "")

	// Synopsis
	detail.Synopsis = strings.TrimSpace(container.Find(".ml-mask .mli-desc").Text())

	// Score
	detail.Score = strings.TrimSpace(container.Find(".ml-mask .mli-mvi span[itemprop='ratingValue']").Text())

	// Genres
	container.Find(".ml-mask .mli-mvi a[itemprop='genre']").Each(func(i int, s *goquery.Selection) {
		detail.Genres = append(detail.Genres, strings.TrimSpace(s.Text()))
	})

	// Episodes
	doc.Find(".tvseason .les-content a").Each(func(i int, s *goquery.Selection) {
		url, exists := s.Attr("href")
		if exists {
			detail.Episodes = append(detail.Episodes, Episode{
				Title:    strings.TrimSpace(s.Text()),
				Endpoint: url,
			})
		}
	})

	// Other Metadata
	container.Find(".mli-mvi").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if strings.Contains(text, "Status :") {
			detail.Metadata["Status"] = strings.TrimSpace(strings.Replace(text, "Status :", "", 1))
		} else if strings.Contains(text, "Duration :") {
			detail.Metadata["Duration"] = strings.TrimSpace(strings.Replace(text, "Duration :", "", 1))
		}
	})

	return detail, nil
}
