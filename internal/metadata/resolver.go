package metadata

func Resolve(url string, log Logger) (*Metadata, error) {

	log.Info("Resolving metadata for URL: " + url)

	// Future:
	// if youtube → scrape
	// if moviesbazar → scrape
	// if imdb/tmdb url → fetch API

	// For now: m3u8 or unknown
	log.Warn("No metadata source detected, switching to interactive mode")
	return PromptUser(log)
}
