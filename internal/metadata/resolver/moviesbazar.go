package resolver

import (
	"github.com/ajaysinghnp/maya-cli/internal/logger/iface"
	"github.com/ajaysinghnp/maya-cli/internal/metadata"
)

// fetchMoviesBazarMetadata simulates fetching metadata from MoviesBazar
func fetchMoviesBazarMetadata(url string, log iface.Logger) (*metadata.Metadata, error) {
	log.Info("Fetching metadata from MoviesBazar for URL: " + url)
	return dummyMetadata(url, log)
}
