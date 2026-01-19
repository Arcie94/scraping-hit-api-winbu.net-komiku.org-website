package winbu

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ParseHome(doc *goquery.Document) (*HomeData, error) {
	data := &HomeData{}

	// Iterate over each movie list wrapper to identify the section content
	doc.Find(".movies-list-wrap").Each(func(i int, s *goquery.Selection) {
		title := strings.ToLower(strings.TrimSpace(s.Find(".list-title h2").Text()))

		var targetList *[]Anime

		if strings.Contains(title, "top 10 series") {
			targetList = &data.TopSeries
		} else if strings.Contains(title, "top 10 film") {
			targetList = &data.TopMovies
		} else if strings.Contains(title, "anime donghua terbaru") || strings.Contains(title, "anime terbaru") {
			targetList = &data.LatestAnime
		} else if strings.Contains(title, "film terbaru") {
			targetList = &data.LatestMovies
		} else if strings.Contains(title, "jepang korea china barat") {
			targetList = &data.InternationalSeries
		}

		if targetList != nil {
			s.Find(".ml-item").Each(func(_ int, item *goquery.Selection) {
				*targetList = append(*targetList, extractAnimeFromItem(item))
			})
		}
	})

	// Genres
	// Updated selector based on winbu_home.html structure
	doc.Find("#List-Anime .list-group-item a").Each(func(i int, s *goquery.Selection) {
		name := strings.TrimSpace(s.Text())
		endpoint := s.AttrOr("href", "")
		// Filter out non-genre links if necessary, though most seem to be genres or lists
		if name != "" && endpoint != "" && !strings.Contains(strings.ToLower(name), "daftar anime") {
			data.Genres = append(data.Genres, Genre{
				Name:     name,
				Endpoint: endpoint,
			})
		}
	})

	return data, nil
}

func ParseSearch(doc *goquery.Document) ([]Anime, error) {
	var results []Anime
	// Updated selector: search results use .a-item instead of .ml-item
	doc.Find(".a-item").Each(func(i int, s *goquery.Selection) {
		anime := extractAnimeFromItem(s)
		results = append(results, anime)
	})
	return results, nil
}

func extractAnimeFromItem(s *goquery.Selection) Anime {
	anime := Anime{
		Title:    strings.TrimSpace(s.Find(".mli-info").Text()), // Fallback
		Endpoint: s.Find("a").First().AttrOr("href", ""),
		Thumb:    s.Find("img").AttrOr("data-original", ""),
	}

	// Specific title selector
	if title := strings.TrimSpace(s.Find(".mli-info h2").Text()); title != "" {
		anime.Title = title
	} else if title := strings.TrimSpace(s.Find(".mli-info .judul").Text()); title != "" {
		anime.Title = title
	}

	// For search results (.a-item), try ml-mask title attribute
	if anime.Title == "" {
		if titleAttr, exists := s.Find("a.ml-mask").Attr("title"); exists {
			anime.Title = strings.TrimSpace(titleAttr)
		}
	}

	if anime.Thumb == "" {
		anime.Thumb = s.Find("img").AttrOr("src", "")
	}

	// Extract hidden info if available
	infoHidden := s.Find(".info-hidden")
	if infoHidden.Length() > 0 {
		anime.Rating = infoHidden.AttrOr("data-rating", "")
		// Use episode count as status/extra info if available
		if ep := infoHidden.AttrOr("data-episode", ""); ep != "" && ep != "0" {
			anime.Status = "Ep " + ep
		}
	}

	// Fallback for Rating from visual element (star icon)
	if anime.Rating == "" || anime.Rating == "0" {
		// Look for .mli-mvi that has fa-star
		s.Find(".mli-mvi").Each(func(i int, sel *goquery.Selection) {
			if sel.Find(".fa-star").Length() > 0 {
				anime.Rating = strings.TrimSpace(sel.Text())
			}
		})
	}

	// Fallback for Status from visual element (e.g. Top 10)
	if anime.Status == "" {
		if top := strings.TrimSpace(s.Find(".mli-topten b").Text()); top != "" {
			anime.Status = "Rank " + top
		}
	}

	// Clean up Rating (remove newlines/spaces)
	anime.Rating = strings.TrimSpace(anime.Rating)

	return anime
}
