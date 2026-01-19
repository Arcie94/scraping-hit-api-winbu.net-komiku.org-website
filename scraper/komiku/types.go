package komiku

type Manga struct {
	Title       string
	Endpoint    string
	Thumb       string
	Type        string
	Score       string
	Description string
}

type MangaDetail struct {
	Title       string
	Thumb       string
	Synopsis    string
	Description string // For UI compatibility
	Status      string
	Authors     []string
	Genres      []string
	Chapters    []ChapterLink
	Metadata    map[string]string
}

type ChapterLink struct {
	Title        string
	Endpoint     string
	Number       string
	DateUploaded string
	ViewCount    string
}

type HomeData struct {
	Trending []Manga
	Popular  []Manga
	Latest   []Manga
}

type ChapterImage struct {
	URL    string
	Number int
}

type Genre struct {
	Name     string
	Endpoint string
}
