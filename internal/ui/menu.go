package ui

import (
	"bufio"
	"fmt"
	"komiku-scraper/internal/service"
	"komiku-scraper/scraper/winbu"
	"log"
	"os"
	"strings"
)

// StartMenu starts the interactive CLI
func StartMenu(komikuSvc *service.KomikuService, winbuSvc *service.WinbuService) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n=== AUTO SCRAPER BOT ===")
		fmt.Println("Pilih Provider:")
		fmt.Println("1. Komiku.org (Manga/Komik)")
		fmt.Println("2. Winbu.net (Anime/Streaming)")
		fmt.Println("0. Keluar")

		fmt.Print("Pilihan: ")
		if scanner.Scan() {
			choice := scanner.Text()
			switch choice {
			case "1":
				menuKomiku(komikuSvc, scanner)
			case "2":
				menuWinbu(winbuSvc, scanner)
			case "0":
				fmt.Println("Bye!")
				return
			default:
				fmt.Println("Pilihan salah")
			}
		}
	}
}

func menuKomiku(svc *service.KomikuService, scanner *bufio.Scanner) {
	for {
		fmt.Println("\n--- KOMIKU PROVIDER ---")
		fmt.Println("Features:")
		fmt.Println("1. Search Manga (Auto Selector)")
		fmt.Println("2. Detail Manga (By URL)")
		fmt.Println("3. Manga Trending")
		fmt.Println("4. Manga Populer")
		fmt.Println("5. Read Chapter (Extract Images)")
		fmt.Println("6. Recommendations (From Chapter Page)")
		fmt.Println("7. List Genre")
		fmt.Println("0. Keluar")

		fmt.Print("\nPilih Menu (1-7, 0 Exit): ")
		if scanner.Scan() {
			choice := scanner.Text()

			switch choice {
			case "0":
				fmt.Println("Bye!")
				return

			case "1":
				fmt.Print("Masukkan Kata Kunci: ")
				if scanner.Scan() {
					keyword := scanner.Text()
					searchURL := fmt.Sprintf("https://api.komiku.org/?post_type=manga&s=%s", strings.ReplaceAll(keyword, " ", "+"))

					res, err := svc.FetchAndParseList(searchURL)
					if err != nil {
						log.Println("Error:", err)
						continue
					}
					fmt.Printf("\nHasil Pencarian '%s': %d manga\n", keyword, len(res))
					for i, m := range res {
						fmt.Printf("%d. %s (%s)\n   Url: %s\n", i+1, m.Title, m.Type, m.Endpoint)
						if i >= 9 {
							break
						}
					}

					if len(res) > 0 {
						fmt.Print("\nPilih nomor untuk detail (0 batal): ")
						if scanner.Scan() {
							var sel int
							fmt.Sscanf(scanner.Text(), "%d", &sel)
							if sel > 0 && sel <= len(res) {
								handleDetail(svc, scanner, res[sel-1].Endpoint)
							}
						}
					}
				}

			case "2":
				fmt.Print("Masukkan URL Manga (contoh: /manga/komik-one-piece-indo/): ")
				if scanner.Scan() {
					handleDetail(svc, scanner, scanner.Text())
				}

			case "3": // Manga Trending
				fmt.Println("Mengambil Data Trending...")
				homeData, err := svc.FetchHomeData()
				if err != nil {
					log.Println("Error:", err)
					continue
				}

				fmt.Println("\n🔥 Trending / Peringkat:")
				for i, m := range homeData.Trending {
					fmt.Printf("#%d %s\n", i+1, m.Title)
				}

				// Add selection prompt
				if len(homeData.Trending) > 0 {
					fmt.Print("\nPilih nomor manga trending untuk detail (0 batal): ")
					if scanner.Scan() {
						var sel int
						fmt.Sscanf(scanner.Text(), "%d", &sel)
						if sel > 0 && sel <= len(homeData.Trending) {
							handleDetail(svc, scanner, homeData.Trending[sel-1].Endpoint)
						}
					}
				}

			case "4": // Manga Populer
				fmt.Println("Mengambil Data Populer...")
				homeData, err := svc.FetchHomeData()
				if err != nil {
					log.Println("Error:", err)
					continue
				}

				fmt.Println("\n⭐ Populer:")
				for i, m := range homeData.Popular {
					fmt.Printf("%d. %s (%s)\n", i+1, m.Title, m.Type)
				}

				// Add selection prompt
				if len(homeData.Popular) > 0 {
					fmt.Print("\nPilih nomor manga populer untuk detail (0 batal): ")
					if scanner.Scan() {
						var sel int
						fmt.Sscanf(scanner.Text(), "%d", &sel)
						if sel > 0 && sel <= len(homeData.Popular) {
							handleDetail(svc, scanner, homeData.Popular[sel-1].Endpoint)
						}
					}
				}

			case "5": // Read Chapter
				fmt.Print("Masukkan URL Chapter: ")
				if scanner.Scan() {
					url := scanner.Text()
					if !strings.HasPrefix(url, "http") {
						url = "https://komiku.org" + url
					}

					images, err := svc.FetchChapterImages(url)
					if err != nil {
						log.Println("Error:", err)
						continue
					}
					fmt.Printf("\nDitemukan %d gambar:\n", len(images))
					OpenInBrowser("Read Chapter", images)
				}

			case "6": // Recommendations
				fmt.Print("Masukkan URL Chapter untuk Rekomendasi: ")
				if scanner.Scan() {
					url := scanner.Text()
					if !strings.HasPrefix(url, "http") {
						url = "https://komiku.org" + url
					}

					recs, err := svc.FetchRecommendations(url)
					if err != nil {
						log.Println("Error:", err)
						continue
					}
					fmt.Println("\nRekomendasi Manga:")
					for i, m := range recs {
						fmt.Printf("- %s\n", m.Title)
						if i >= 5 {
							break
						}
					}
				}

			case "7": // List Genre
				fmt.Println("Mengambil Daftar Genre...")
				genres, err := svc.FetchGenreList()
				if err != nil {
					log.Println("Error:", err)
					continue
				}
				fmt.Printf("\nDitemukan %d Genre:\n", len(genres))
				cutoff := 15
				for i, g := range genres {
					fmt.Printf("- %s [%s]\n", g.Name, g.Endpoint)
					if i >= cutoff {
						fmt.Printf("... dan %d lainnya\n", len(genres)-cutoff)
						break
					}
				}
			}
		}
	}
}

func menuWinbu(svc *service.WinbuService, scanner *bufio.Scanner) {
	for {
		fmt.Println("\n--- WINBU PROVIDER ---")
		fmt.Println("1. Search Anime (Kata Kunci)")
		fmt.Println("2. Top 10 Anime")
		fmt.Println("3. Top 10 Film")
		fmt.Println("4. Film Terbaru")
		fmt.Println("5. Anime/Donghua Terbaru")
		fmt.Println("6. Drama")
		fmt.Println("7. List Genre")
		fmt.Println("0. Kembali")

		fmt.Print("Pilihan: ")
		if scanner.Scan() {
			switch scanner.Text() {
			case "1":
				fmt.Print("Masukkan Kata Kunci: ")
				if scanner.Scan() {
					keyword := scanner.Text()
					res, err := svc.FetchSearch(strings.ReplaceAll(keyword, " ", "+"))
					if err != nil {
						fmt.Println("Error:", err)
						continue
					}
					fmt.Printf("\nHasil %d anime:\n", len(res))
					for i, a := range res {
						fmt.Printf("%d. %s\n   %s | %s\n", i+1, a.Title, a.Status, a.Type)
					}

					if len(res) > 0 {
						fmt.Print("\nPilih nomor untuk detail (0 batal): ")
						if scanner.Scan() {
							var sel int
							fmt.Sscanf(scanner.Text(), "%d", &sel)
							if sel > 0 && sel <= len(res) {
								handleDetailWinbu(svc, scanner, res[sel-1].Endpoint)
							}
						}
					}
				}
			case "2": // Top 10 Anime
				fmt.Println("Mengambil Top 10 Anime...")
				data, err := svc.FetchHomeData()
				if err != nil {
					log.Println("Error:", err)
					continue
				}
				if len(data.TopSeries) > 0 {
					fmt.Println("\n🏆 Top 10 Anime:")
					for i, a := range data.TopSeries {
						fmt.Printf("%d. %s [%s] (%s)\n", i+1, a.Title, a.Rating, a.Status)
					}

					fmt.Print("\nPilih nomor untuk detail (0 batal): ")
					if scanner.Scan() {
						var sel int
						fmt.Sscanf(scanner.Text(), "%d", &sel)
						if sel > 0 && sel <= len(data.TopSeries) {
							handleDetailWinbu(svc, scanner, data.TopSeries[sel-1].Endpoint)
						}
					}
				} else {
					fmt.Println("Tidak ada data ditemukan.")
				}

			case "3": // Top 10 Film (NEW)
				fmt.Println("Mengambil Top 10 Film...")
				data, err := svc.FetchHomeData()
				if err != nil {
					log.Println("Error:", err)
					continue
				}
				if len(data.TopMovies) > 0 {
					fmt.Println("\n🏆 Top 10 Film:")
					for i, a := range data.TopMovies {
						fmt.Printf("%d. %s [%s] (%s)\n", i+1, a.Title, a.Rating, a.Status)
					}

					fmt.Print("\nPilih nomor untuk detail (0 batal): ")
					if scanner.Scan() {
						var sel int
						fmt.Sscanf(scanner.Text(), "%d", &sel)
						if sel > 0 && sel <= len(data.TopMovies) {
							handleDetailWinbu(svc, scanner, data.TopMovies[sel-1].Endpoint)
						}
					}
				} else {
					fmt.Println("Tidak ada data Top Film ditemukan atau filter parser perlu diperiksa.")
				}

			case "4": // Latest Movies
				fmt.Println("Mengambil Film Terbaru...")
				data, err := svc.FetchHomeData()
				if err != nil {
					log.Println("Error:", err)
					continue
				}
				if len(data.LatestMovies) > 0 {
					fmt.Println("\n🎬 Film Terbaru:")
					for i, a := range data.LatestMovies {
						fmt.Printf("%d. %s [%s]\n", i+1, a.Title, a.Rating)
						if i >= 19 {
							break
						} // Limit display
					}

					fmt.Print("\nPilih nomor untuk detail (0 batal): ")
					if scanner.Scan() {
						var sel int
						fmt.Sscanf(scanner.Text(), "%d", &sel)
						if sel > 0 && sel <= len(data.LatestMovies) {
							handleDetailWinbu(svc, scanner, data.LatestMovies[sel-1].Endpoint)
						}
					}
				}

			case "5": // Latest Anime
				fmt.Println("Mengambil Anime/Donghua Terbaru...")
				data, err := svc.FetchHomeData()
				if err != nil {
					log.Println("Error:", err)
					continue
				}
				if len(data.LatestAnime) > 0 {
					fmt.Println("\n🆕 Anime/Donghua Terbaru:")
					for i, a := range data.LatestAnime {
						fmt.Printf("%d. %s [%s]\n", i+1, a.Title, a.Status) // Status used for episode
						if i >= 19 {
							break
						} // Limit display
					}

					fmt.Print("\nPilih nomor untuk detail (0 batal): ")
					if scanner.Scan() {
						var sel int
						fmt.Sscanf(scanner.Text(), "%d", &sel)
						if sel > 0 && sel <= len(data.LatestAnime) {
							handleDetailWinbu(svc, scanner, data.LatestAnime[sel-1].Endpoint)
						}
					}
				}

			case "6": // International Series
				fmt.Println("Mengambil Drama Terbaru...")
				data, err := svc.FetchHomeData()
				if err != nil {
					log.Println("Error:", err)
					continue
				}
				if len(data.InternationalSeries) > 0 {
					fmt.Println("\n📺 Drama:")
					for i, a := range data.InternationalSeries {
						fmt.Printf("%d. %s [%s]\n", i+1, a.Title, a.Status)
						if i >= 19 {
							break
						}
					}

					fmt.Print("\nPilih nomor untuk detail (0 batal): ")
					if scanner.Scan() {
						var sel int
						fmt.Sscanf(scanner.Text(), "%d", &sel)
						if sel > 0 && sel <= len(data.InternationalSeries) {
							handleDetailWinbu(svc, scanner, data.InternationalSeries[sel-1].Endpoint)
						}
					}
				}

			case "7": // List Genre
				fmt.Println("Mengambil Daftar Genre...")
				data, err := svc.FetchHomeData()
				if err != nil {
					log.Println("Error:", err)
					continue
				}
				if len(data.Genres) > 0 {
					fmt.Println("\n📂 Genre Tersedia:")
					for i, g := range data.Genres {
						fmt.Printf("%d. %s\n", i+1, g.Name)
					}
					// Note: Selecting genre is not implemented yet as it requires fetching category/search result parsing
					fmt.Println("\n(Pilih genre belum diimplementasikan, fitur hanya menampilkan list saat ini)")
				}

			case "0":
				return
			default:
				fmt.Println("Pilihan tidak valid")
			}
		}
	}
}

func handleDetail(svc *service.KomikuService, scanner *bufio.Scanner, slug string) {
	if !strings.HasPrefix(slug, "http") {
		slug = "https://komiku.org" + slug
	}

	detail, err := svc.FetchAndParseDetail(slug)
	if err != nil {
		log.Println("Error fetching detail:", err)
		return
	}

	fmt.Printf("\n--- %s ---\n", detail.Title)
	fmt.Printf("Status: %s\n", detail.Status)
	fmt.Printf("Genres: %v\n", detail.Genres)
	if len(detail.Description) > 100 {
		fmt.Printf("Sinopsis: %s...\n", detail.Description[:100])
	} else {
		fmt.Printf("Sinopsis: %s\n", detail.Description)
	}
	fmt.Printf("Total Chapter: %d\n", len(detail.Chapters))

	if len(detail.Chapters) > 0 {
		fmt.Printf("Chapter Terakhir: %s (%s)\n", detail.Chapters[0].Title, detail.Chapters[0].DateUploaded)

		fmt.Println("\n5 Chapter Teratas:")
		limit := 5
		if len(detail.Chapters) < 5 {
			limit = len(detail.Chapters)
		}

		for i := 0; i < limit; i++ {
			c := detail.Chapters[i]
			fmt.Printf("%d. %s (%s)\n", i+1, c.Title, c.ViewCount)
		}

		fmt.Print("\nPilih nomor chapter untuk membaca (0 kembali): ")
		if scanner.Scan() {
			var sel int
			fmt.Sscanf(scanner.Text(), "%d", &sel)
			if sel > 0 && sel <= limit {
				targetChapter := detail.Chapters[sel-1]
				fmt.Printf("Membaca %s...\n", targetChapter.Title)

				images, err := svc.FetchChapterImages("https://komiku.org" + targetChapter.Endpoint)
				if err != nil {
					log.Println("Error:", err)
					return
				}

				fmt.Printf("\nDitemukan %d gambar:\n", len(images))
				OpenInBrowser(targetChapter.Title, images)
			}
		}
	}
}

func handleDetailWinbu(svc *service.WinbuService, scanner *bufio.Scanner, slug string) {
	fmt.Println("Mengambil Detail Anime...")
	detail, err := svc.FetchAndParseDetail(slug)
	if err != nil {
		log.Println("Error fetching detail:", err)
		return
	}

	fmt.Printf("\n--- %s ---\n", detail.Title)
	fmt.Printf("Genres: %v\n", detail.Genres)
	if len(detail.Synopsis) > 150 {
		fmt.Printf("Sinopsis: %s...\n", detail.Synopsis[:150])
	} else {
		fmt.Printf("Sinopsis: %s\n", detail.Synopsis)
	}

	fmt.Printf("Total Episode: %d\n", len(detail.Episodes))

	if len(detail.Episodes) > 0 {
		var targetEpisode *winbu.Episode

		if len(detail.Episodes) == 1 {
			// Single Movie Logic
			fmt.Printf("\nFilm Tunggal Terdeteksi: %s\n", detail.Episodes[0].Title)
			fmt.Print("Mulai Videonya? (y/n, default y): ")
			if scanner.Scan() {
				ans := strings.ToLower(scanner.Text())
				if ans == "" || ans == "y" || ans == "yes" {
					targetEpisode = &detail.Episodes[0]
				}
			}
		} else {
			// Series Logic
			fmt.Println("\n5 Episode Teratas:")
			limit := 5
			if len(detail.Episodes) < 5 {
				limit = len(detail.Episodes)
			}

			for i := 0; i < limit; i++ {
				ep := detail.Episodes[i]
				fmt.Printf("%d. %s\n", i+1, ep.Title)
			}

			fmt.Print("\nPilih nomor episode untuk nonton (0 kembali): ")
			if scanner.Scan() {
				var sel int
				fmt.Sscanf(scanner.Text(), "%d", &sel)
				if sel > 0 && sel <= limit {
					targetEpisode = &detail.Episodes[sel-1]
				}
			}
		}

		if targetEpisode != nil {
			fmt.Printf("Mengambil data episode %s...\n", targetEpisode.Title)
			epData, err := svc.FetchEpisode(targetEpisode.Endpoint)
			if err != nil {
				log.Println("Error fetching episode:", err)
				return
			}

			fmt.Printf("\n--- %s ---\n", epData.Title)

			// Show Stream Options
			if len(epData.StreamOptions) > 0 {
				fmt.Println("\nServer Stream:")
				for i, opt := range epData.StreamOptions {
					if opt.Quality != "" {
						fmt.Printf("%d. %s [%s]\n", i+1, opt.Server, opt.Quality)
					} else {
						fmt.Printf("%d. %s\n", i+1, opt.Server)
					}
				}
			}

			// Show Download Links
			if len(epData.DownloadLinks) > 0 {
				fmt.Println("\nLink Download:")
				for _, link := range epData.DownloadLinks {
					fmt.Printf("- %s [%s]: %s\n", link.Server, link.Quality, link.URL)
				}
			}

			if len(epData.StreamOptions) > 0 {
				fmt.Print("\nPilih server untuk streaming (0 untuk batal): ")
				if scanner.Scan() {
					var sSel int
					fmt.Sscanf(scanner.Text(), "%d", &sSel)
					if sSel > 0 && sSel <= len(epData.StreamOptions) {
						targetOpt := epData.StreamOptions[sSel-1]
						fmt.Println("Mengambil URL video...")
						vidURL, err := svc.ResolveStream(targetOpt)
						if err != nil {
							log.Println("Error resolving stream:", err)
						} else {
							fmt.Printf("\nVIDEO URL: %s\n", vidURL)

							fmt.Print("Buka player di browser? (y/n, default n): ")
							if scanner.Scan() {
								ans := strings.ToLower(scanner.Text())
								if ans == "y" || ans == "yes" {
									OpenURL(vidURL)
								} else {
									fmt.Println("Silahkan copy URL di atas ke player favorit Anda.")
								}
							}
						}
					}
				}
			} else {
				fmt.Println("Tidak ada opsi streaming tersedia.")
			}
		}
	}
}
