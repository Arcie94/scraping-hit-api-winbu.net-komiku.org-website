package winbu

type Anime struct {
	Title    string
	Endpoint string
	Thumb    string
	Type     string // "Movie" or "Series"
	Rating   string
	Status   string
}

type AnimeDetail struct {
	Title    string
	Thumb    string
	Synopsis string
	Score    string
	Genres   []string
	Episodes []Episode
	Metadata map[string]string // Status, Type, Released, etc.
}

type Episode struct {
	Title    string
	Endpoint string
}

type EpisodePageData struct {
	Title               string
	EpisodeNumber       string
	StreamOptions       []StreamOption
	NextEpisodeEndpoint string
	PrevEpisodeEndpoint string
	AllEpisodes         []Episode
	DownloadLinks       []DownloadLink
}

type DownloadLink struct {
	Server  string
	URL     string
	Quality string
}

type StreamOption struct {
	Name    string
	Server  string
	Quality string
	PostID  string
	Nume    string
	Type    string
}

type HomeData struct {
	TopSeries           []Anime
	TopMovies           []Anime
	LatestMovies        []Anime
	LatestAnime         []Anime
	InternationalSeries []Anime
	Genres              []Genre
}

type Genre struct {
	Name     string
	Endpoint string
}
