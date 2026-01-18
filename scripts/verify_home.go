package main

import (
	"fmt"
	"komiku-scraper/internal/service"
	"komiku-scraper/scraper/winbu"
	"log"
)

func main() {
	client := winbu.NewWinbuClient()
	svc := service.NewWinbuService(client)

	fmt.Println("Fetching Homepage Data...")
	data, err := svc.FetchHomeData()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n=== TOP 10 SERIES (%d) ===\n", len(data.TopSeries))
	for i, a := range data.TopSeries {
		fmt.Printf("%d. %s [%s] (%s)\n   %s | %s\n", i+1, a.Title, a.Rating, a.Status, a.Type, a.Endpoint)
	}

	fmt.Printf("\n=== TOP 10 MOVIES (%d) ===\n", len(data.TopMovies))
	for i, a := range data.TopMovies {
		fmt.Printf("%d. %s [%s] (%s)\n   %s | %s\n", i+1, a.Title, a.Rating, a.Status, a.Type, a.Endpoint)
	}

	fmt.Printf("\n=== LATEST MOVIES (%d) ===\n", len(data.LatestMovies))
	for i, a := range data.LatestMovies {
		if i >= 5 {
			break
		}
		fmt.Printf("%d. %s [%s]\n   %s | %s\n", i+1, a.Title, a.Rating, a.Type, a.Endpoint)
	}

	fmt.Printf("\n=== LATEST ANIME (%d) ===\n", len(data.LatestAnime))
	for i, a := range data.LatestAnime {
		if i >= 5 {
			break
		}
		fmt.Printf("%d. %s [%s]\n   %s | %s\n", i+1, a.Title, a.Status, a.Type, a.Endpoint)
	}

	fmt.Printf("\n=== INTERNATIONAL SERIES (%d) ===\n", len(data.InternationalSeries))
	for i, a := range data.InternationalSeries {
		if i >= 5 {
			break
		}
		fmt.Printf("%d. %s [%s]\n   %s | %s\n", i+1, a.Title, a.Status, a.Type, a.Endpoint)
	}

	fmt.Printf("\n=== GENRES (%d) ===\n", len(data.Genres))
	for i, g := range data.Genres {
		if i >= 10 {
			break
		}
		fmt.Printf("- %s: %s\n", g.Name, g.Endpoint)
	}
}
