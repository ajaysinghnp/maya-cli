package metadata

import (
	"fmt"
	"path/filepath"
)

func (m *Metadata) BuildPaths(base string, ext string, log Logger) {
	log.Info("Building media paths")

	idPart := ""
	if m.IDs.TMDB != "" {
		idPart = " [tmdb-" + m.IDs.TMDB + "]"
		log.Debug("Using TMDB ID: " + m.IDs.TMDB)
	} else if m.IDs.IMDB != "" {
		idPart = " [imdb-" + m.IDs.IMDB + "]"
		log.Debug("Using IMDB ID: " + m.IDs.IMDB)
	}

	// Base defaults to current directory
	if base == "" {
		base = "."
		log.Debug("No output directory specified, using current directory")
	} else {
		log.Debug("Using output directory: " + base)
	}

	if m.Type == Movie {
		m.RootDir = filepath.Join(
			base,
			fmt.Sprintf("%s (%s)%s", m.Title, m.Year, idPart),
		)

		m.MediaFile = filepath.Join(
			m.RootDir,
			fmt.Sprintf("%s (%s).%s", m.Title, m.Year, ext),
		)

		log.Info("Detected movie content")
		log.Info("Movie directory: " + m.RootDir)
		log.Info("Movie file: " + m.MediaFile)
		return
	}

	// Series
	log.Info("Detected series content")

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

	log.Info("Series root directory: " + m.RootDir)
	log.Info("Season directory: " + m.SeasonDir)
	log.Info("Episode file: " + m.MediaFile)
}
