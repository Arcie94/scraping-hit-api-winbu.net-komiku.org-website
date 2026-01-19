package winbu

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ParseEpisodePage(doc *goquery.Document) (*EpisodePageData, error) {
	data := &EpisodePageData{}

	// Title
	data.Title = strings.TrimSpace(doc.Find("h1.titless").Text())
	if data.Title == "" {
		data.Title = strings.TrimSpace(doc.Find("title").Text())
	}

	// Stream Options (The Servers)
	// Based on CSS .east_player_option
	doc.Find(".east_player_option").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())

		// Attempt to split Name and Quality
		// Common formats: "Server 720p", "Server - 1080p", "Server [360p]"
		var name, quality string

		resolutions := []string{"1080p", "1080", "720p", "720", "480p", "480", "360p", "360"}
		foundRes := false

		for _, res := range resolutions {
			if strings.Contains(strings.ToLower(text), res) {
				quality = res
				if !strings.HasSuffix(quality, "p") {
					quality += "p"
				}
				name = strings.TrimSpace(strings.NewReplacer(res, "", "-", "", "[", "", "]", "", "(", "", ")", "").Replace(text))
				foundRes = true
				break
			}
		}

		// Check title attribute if text parsing failed
		if !foundRes {
			titleAttr := s.AttrOr("title", "")
			for _, res := range resolutions {
				if strings.Contains(strings.ToLower(titleAttr), res) {
					quality = res
					if !strings.HasSuffix(quality, "p") {
						quality += "p"
					}
					// Only use title parsing for quality, keep original name or use title?
					// Usually name is better.
					foundRes = true
					break
				}
			}
		}

		if !foundRes {
			// Fallback checks
			if strings.Contains(strings.ToLower(text), "hd") {
				quality = "HD"
				name = strings.TrimSpace(strings.ReplaceAll(text, "HD", ""))
			} else if strings.Contains(strings.ToLower(text), "sd") {
				quality = "SD"
				name = strings.TrimSpace(strings.ReplaceAll(text, "SD", ""))
			} else {
				// No quality found in text, check attributes?
				name = text
				quality = "" // Leave empty if not found, to avoid ugly "Unknown"
			}
		}

		opt := StreamOption{
			Name:    name,
			Server:  name,
			Quality: quality,
			PostID:  s.AttrOr("data-post", ""),
			Nume:    s.AttrOr("data-nume", ""),
			Type:    s.AttrOr("data-type", ""),
		}
		// Only add if we have at least ID and Nume
		if opt.PostID != "" && opt.Nume != "" {
			data.StreamOptions = append(data.StreamOptions, opt)
		}
	})

	// Fallback: Check if there's only one default player without options list
	if len(data.StreamOptions) == 0 {
		// Try to find hidden inputs or script variables if needed
		// For now, check if there is a meta tag or similar
	}

	// Navigation
	// CSS .naveps .nvsc a
	doc.Find(".naveps .nvsc a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}
		text := strings.ToLower(s.Text())
		if strings.Contains(text, "next") || strings.Contains(text, "selanjutnya") {
			data.NextEpisodeEndpoint = href
		} else if strings.Contains(text, "prev") || strings.Contains(text, "sebelumnya") {
			data.PrevEpisodeEndpoint = href
		}
	})

	// Fallback for navigation
	if data.NextEpisodeEndpoint == "" {
		data.NextEpisodeEndpoint = doc.Find(".fr a").AttrOr("href", "")
	}
	if data.PrevEpisodeEndpoint == "" {
		data.PrevEpisodeEndpoint = doc.Find(".fl a").AttrOr("href", "")
	}

	// Download Links Parsing
	// Attempt to find download links in common containers
	// Strategy 1: .download-eps
	doc.Find(".download-eps a").Each(func(i int, s *goquery.Selection) {
		link := DownloadLink{
			Server:  strings.TrimSpace(s.Text()),
			URL:     s.AttrOr("href", ""),
			Quality: "Unknown", // Often mixed in text
		}

		// Try to extract quality from parent or text
		if strings.Contains(link.Server, "360") {
			link.Quality = "360p"
		}
		if strings.Contains(link.Server, "480") {
			link.Quality = "480p"
		}
		if strings.Contains(link.Server, "720") {
			link.Quality = "720p"
		}
		if strings.Contains(link.Server, "1080") {
			link.Quality = "1080p"
		}

		if link.URL != "" && !strings.HasPrefix(link.URL, "javascript") {
			data.DownloadLinks = append(data.DownloadLinks, link)
		}
	})

	// Strategy 2: #download container
	if len(data.DownloadLinks) == 0 {
		doc.Find("#download a").Each(func(i int, s *goquery.Selection) {
			link := DownloadLink{
				Server:  strings.TrimSpace(s.Text()),
				URL:     s.AttrOr("href", ""),
				Quality: "Unknown",
			}
			if link.URL != "" {
				data.DownloadLinks = append(data.DownloadLinks, link)
			}
		})
	}

	return data, nil
}
