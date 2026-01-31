package metadata

import (
	"fmt"
	"path/filepath"

	"github.com/ajaysinghnp/maya-cli/internal/logger/iface"
)

func (m *Metadata) BuildPaths(base string, ext string, log iface.Logger) {
	log.Info("Building media paths")

	idPart := ""
	if m.IDs.TMDB != "" {
		idPart = " [tmdb-" + m.IDs.TMDB + "]"
		log.Success("Using TMDB ID: " + m.IDs.TMDB)
	} else if m.IDs.IMDB != "" {
		idPart = " [imdb-" + m.IDs.IMDB + "]"
		log.Success("Using IMDB ID: " + m.IDs.IMDB)
	}

	// Base defaults to current directory
	if base == "" {
		base = "."
		log.Info("No output directory specified, using current directory")
	} else {
		log.Info("Using output directory: " + base)
	}

	if m.Type == Movie {
		m.RootDir = filepath.Join(
			base,
			fmt.Sprintf("%s (%d)%s", m.Title, m.Year, idPart),
		)

		m.MediaFile = filepath.Join(
			m.RootDir,
			fmt.Sprintf("%s (%d).%s", m.Title, m.Year, ext),
		)

		log.Success(fmt.Sprintf("Detected content: %s", m.Type))
		log.Success("Movie directory: " + m.RootDir)
		log.Success("Movie file: " + m.MediaFile)
		return
	}

	// Series
	log.Success("Detected series content")

	m.RootDir = filepath.Join(
		base,
		fmt.Sprintf("%s (%s)%s", m.Title, m.Year, idPart),
	)

	m.SeasonDir = filepath.Join(
		m.RootDir,
		fmt.Sprintf("Season %02d (%s)", m.Season, m.Year),
	)

	m.MediaFile = filepath.Join(
		m.SeasonDir,
		fmt.Sprintf(
			"%s - S%02dE%02d - %s.%s",
			m.Title,
			m.Season,
			m.Episode,
			m.EpisodeTitle,
			ext,
		),
	)

	log.Success("Series root directory: " + m.RootDir)
	log.Success("Season directory: " + m.SeasonDir)
	log.Success("Episode file: " + m.MediaFile)
}
