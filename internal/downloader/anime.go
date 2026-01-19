package downloader

import (
	"fmt"
	"komiku-scraper/scraper/winbu"
	"os"
	"path/filepath"
	"time"
)

// SaveAnimeInfo saves episode stream and download info to a text file
func (d *Downloader) SaveAnimeInfo(animeTitle, episodeTitle string, data *winbu.EpisodePageData, streamURL string) error {
	safeAnimeTitle := SanitizeFilename(animeTitle)
	safeEpisodeTitle := SanitizeFilename(episodeTitle)

	// Create structure: Downloads/Anime/Title/Episode
	saveDir := filepath.Join(d.BaseDir, "Anime", safeAnimeTitle, safeEpisodeTitle)
	if err := EnsureDir(saveDir); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	filename := filepath.Join(saveDir, "info.txt")
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write Info
	fmt.Fprintf(file, "Title: %s\n", animeTitle)
	fmt.Fprintf(file, "Episode: %s\n", episodeTitle)
	fmt.Fprintf(file, "Date: %s\n", time.Now().Format(time.RFC1123))
	fmt.Fprintf(file, "\n--- STREAM URL ---\n")
	if streamURL != "" {
		fmt.Fprintf(file, "%s\n", streamURL)
		fmt.Fprintf(file, "(Use this URL in IDM, XDM, or VLC to play/download)\n")
	} else {
		fmt.Fprintf(file, "No direct stream URL found.\n")
	}

	fmt.Fprintf(file, "\n--- DOWNLOAD LINKS ---\n")
	if len(data.DownloadLinks) > 0 {
		for _, link := range data.DownloadLinks {
			fmt.Fprintf(file, "[%s] %s: %s\n", link.Quality, link.Server, link.URL)
		}
	} else {
		fmt.Fprintf(file, "No direct download links parsed found on page.\n")
	}

	fmt.Printf("\n✅ Info saved to: %s\n", filename)
	return nil
}
