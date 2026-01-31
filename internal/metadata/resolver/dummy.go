package resolver

import (
	"github.com/ajaysinghnp/maya-cli/internal/logger/iface"
	"github.com/ajaysinghnp/maya-cli/internal/metadata"
)

// dummyMetadata generates placeholder metadata for testing
func dummyMetadata(url string, log iface.Logger) (*metadata.Metadata, error) {
	log.Debug("Generating dummy metadata for URL: " + url)
	return &metadata.Metadata{
		Title: "Dummy Title",
		Year:  2026,
		Type:  metadata.Movie, // or metadata.Series if you want
		IDs: metadata.IDs{
			IMDB: "tt0000000",
			TMDB: "000000",
		},
		Episode:      1,
		EpisodeTitle: "Dummy Episode",
		Season:       1,
	}, nil
}
