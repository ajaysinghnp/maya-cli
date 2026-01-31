package metadata

type MediaType string

const (
	Movie  MediaType = "movie"
	Series MediaType = "series"
)

// IDs struct to hold IMDB/TMDB as alternatives
type IDs struct {
	IMDB string `json:"imdbId,omitempty"`
	TMDB string `json:"tmdb,omitempty"`
}

// Watch/source info
type Source struct {
	Label    string `json:"label"`              // e.g., "Hindi (No Ads)"
	URL      string `json:"source"`             // the m3u8 or download link
	LabelTag string `json:"labelTag,omitempty"` // optional
}

// Main metadata struct
type Metadata struct {
	// Core
	Title    string    `json:"title"`
	Year     int       `json:"releaseYear"`
	Type     MediaType `json:"type"`
	Category string    `json:"category"`

	// Dates
	ReleaseDate string `json:"fullReleaseDate"`

	// Content
	Language string   `json:"language"`
	Genres   []string `json:"genre"`
	Cast     []string `json:"castDetails"`

	// IDs
	IDs IDs `json:"ids"` // keep your IDs struct as-is, filled manually if needed

	// Media
	Thumbnail string   `json:"thumbnail"`
	Sources   []Source `json:"watchLink"` // renamed from WatchLink

	// TV-only (optional)
	Season       int    `json:"season,omitempty"`
	Episode      int    `json:"episode,omitempty"`
	EpisodeTitle string `json:"episodeTitle,omitempty"`
	SeasonDir    string `json:"seasonDir,omitempty"`

	// Filesystem (runtime)
	RootDir   string `json:"-"`
	MediaFile string `json:"-"`

	HlsSourceDomain string `json:"hlsSourceDomain"`
}
