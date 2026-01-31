package metadata

type MediaType string

const (
	Movie  MediaType = "movie"
	Series MediaType = "series"
)

type IDs struct {
	IMDB string `json:"imdb,omitempty"`
	TMDB string `json:"tmdb,omitempty"`
}

type Metadata struct {
	Type MediaType `json:"type"`

	Title       string   `json:"title"`
	Year        string   `json:"year,omitempty"`
	ReleaseDate string   `json:"release_date,omitempty"`
	Language    string   `json:"language,omitempty"`
	Cast        []string `json:"cast,omitempty"`
	Genres      []string `json:"genres,omitempty"`
	IDs         IDs      `json:"ids,omitempty"`

	// Series-specific
	Season       int    `json:"season,omitempty"`
	Episode      int    `json:"episode,omitempty"`
	EpisodeTitle string `json:"episode_title,omitempty"`

	// M3U8 URLs
	M3U8URLs []string `json:"m3u8_urls,omitempty"`

	// Computed paths
	RootDir   string `json:"-"`
	MediaDir  string `json:"-"`
	SeasonDir string `json:"-"`
	MediaFile string `json:"-"`
}
