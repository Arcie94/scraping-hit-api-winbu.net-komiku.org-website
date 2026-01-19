package downloader

import (
	"fmt"
	"io"
	"komiku-scraper/scraper/komiku"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// DownloadChapter downloads all images from a chapter
func (d *Downloader) DownloadChapter(mangaTitle, chapterTitle string, images []komiku.ChapterImage) error {
	safeMangaTitle := SanitizeFilename(mangaTitle)
	safeChapterTitle := SanitizeFilename(chapterTitle)

	// Create structure: Downloads/Manga/Title/Chapter
	saveDir := filepath.Join(d.BaseDir, "Manga", safeMangaTitle, safeChapterTitle)
	if err := EnsureDir(saveDir); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	fmt.Printf("\nDownloading to: %s\n", saveDir)
	fmt.Printf("Total Images: %d\n", len(images))

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 5) // Limit to 5 concurrent downloads
	errorsChan := make(chan error, len(images))

	completed := 0
	total := len(images)

	for i, img := range images {
		wg.Add(1)
		go func(idx int, imgUrl string) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire token
			defer func() { <-semaphore }() // Release token

			// Generate filename: 01.jpg, 02.jpg, ...
			ext := ".jpg" // default
			// You could detect extension from URL or Content-Type if needed
			filename := fmt.Sprintf("%03d%s", idx+1, ext)
			filePath := filepath.Join(saveDir, filename)

			// Retry logic
			var err error
			for attempt := 0; attempt < 3; attempt++ {
				err = d.downloadFile(imgUrl, filePath)
				if err == nil {
					break
				}
				time.Sleep(1 * time.Second)
			}

			if err != nil {
				log.Printf("Failed to download image %d: %v", idx+1, err)
				errorsChan <- err
			} else {
				completed++
				// Update progress inline
				fmt.Printf("\rProgress: %d/%d images [%.0f%%]", completed, total, float64(completed)/float64(total)*100)
			}
		}(i, img.URL)
	}

	wg.Wait()
	close(errorsChan)

	// Check for errors
	errCount := 0
	for range errorsChan {
		errCount++
	}

	fmt.Println() // New line after progress
	if errCount > 0 {
		return fmt.Errorf("finished with %d errors", errCount)
	}

	fmt.Println("✅ Download Complete!")
	return nil
}

func (d *Downloader) downloadFile(url, filepath string) error {
	// Check if file already exists
	if _, err := os.Stat(filepath); err == nil {
		return nil // Skip existing
	}

	resp, err := d.GetRequest(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("status code %d", resp.StatusCode)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
