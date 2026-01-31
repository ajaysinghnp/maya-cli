package metadata

import (
	"fmt"
	"os"
	"path/filepath"
)

func WriteMovieNFO(dir string, m *Metadata) error {
	content := fmt.Sprintf(
		`<movie>
  <title>%s</title>
  <year>%s</year>
  <imdbid>%s</imdbid>
  <tmdbid>%s</tmdbid>
</movie>`,
		m.Title,
		m.Year,
		m.IDs.IMDB,
		m.IDs.TMDB,
	)

	return os.WriteFile(
		filepath.Join(dir, "movie.nfo"),
		[]byte(content),
		0644,
	)
}

func WriteEpisodeNFO(dir string, m *Metadata) error {
	content := fmt.Sprintf(
		`<episodedetails>
  <title>%s</title>
  <season>%d</season>
  <episode>%d</episode>
</episodedetails>`,
		m.EpisodeTitle,
		m.Season,
		m.Episode,
	)

	return os.WriteFile(
		filepath.Join(dir, fmt.Sprintf("episode-%02d.nfo", m.Episode)),
		[]byte(content),
		0644,
	)
}
