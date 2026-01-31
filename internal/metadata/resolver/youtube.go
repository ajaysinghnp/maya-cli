package resolver

import (
	"github.com/ajaysinghnp/maya-cli/internal/logger/iface"
	"github.com/ajaysinghnp/maya-cli/internal/metadata"
)

// fetchYouTubeMetadata simulates fetching metadata for YouTube
func fetchYouTubeMetadata(url string, log iface.Logger) (*metadata.Metadata, error) {
	log.Info("Fetching metadata from YouTube for URL: " + url)
	return dummyMetadata(url, log)
}
