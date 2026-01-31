package resolver

import (
	"github.com/ajaysinghnp/maya-cli/internal/logger/iface"
	"github.com/ajaysinghnp/maya-cli/internal/metadata"
)

func Resolve(source SourceType, url string, log iface.Logger) (*metadata.Metadata, error) {
	log.Info("Resolving metadata for URL: " + url)

	switch source {
	case SourceM3U8:
		log.Info("M3U8 detected → manual metadata input required")
		return metadata.PromptUser(log)

	case SourceMoviesBazar:
		log.Info("MoviesBazar detected → scraping metadata")
		return fetchMoviesBazarMetadata(url, log)

	case SourceYouTube:
		log.Info("YouTube detected → scraping metadata")
		return fetchYouTubeMetadata(url, log)

	default:
		log.Warn("Unknown source → manual metadata input required")
		return metadata.PromptUser(log)
	}
}
