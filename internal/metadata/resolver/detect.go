package resolver

import "strings"

type SourceType int

const (
	SourceUnknown SourceType = iota
	SourceM3U8
	SourceMoviesBazar
	SourceYouTube
	SourceIMDB
	SourceTMDB
)

// DetectSource analyzes a URL and returns the source type
func DetectSource(url string) SourceType {
	switch {
	case strings.Contains(url, ".m3u8"):
		return SourceM3U8
	case strings.Contains(url, "moviesbazar"):
		return SourceMoviesBazar
	case strings.Contains(url, "youtube.com") || strings.Contains(url, "youtu.be"):
		return SourceYouTube
	default:
		return SourceUnknown
	}
}
