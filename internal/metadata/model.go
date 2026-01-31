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

	Title string `json:"title"`
	Year  string `json:"year,omitempty"`
	IDs   IDs    `json:"ids,omitempty"`

	// Series-specific
	Season       int    `json:"season,omitempty"`
	Episode      int    `json:"episode,omitempty"`
	EpisodeTitle string `json:"episode_title,omitempty"`

	// Computed paths
	RootDir   string `json:"-"`
	MediaDir  string `json:"-"`
	SeasonDir string `json:"-"`
	MediaFile string `json:"-"`

	// New field: store extracted M3U8 URLs for this media
	M3U8URLs []string `json:"-"`
}
