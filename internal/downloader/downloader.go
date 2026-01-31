package downloader

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/ajaysinghnp/maya-cli/internal/downloader/m3u8"
	moviebazar "github.com/ajaysinghnp/maya-cli/internal/downloader/movie-bazar"
	"github.com/ajaysinghnp/maya-cli/internal/logger/iface"
	"github.com/ajaysinghnp/maya-cli/internal/metadata"
)

type Logger interface {
	Info(string)
	Debug(string)
	Warn(string)
	Error(string)
}

type Downloader struct {
	log iface.Logger
}

func New(log iface.Logger) *Downloader {
	return &Downloader{log: log}
}

// StartDownload decides which module to use based on URL
func (d *Downloader) StartDownload(
	url string,
	output string,
	resume bool,
	concurrent int,
) error {
	// 1️⃣ Resolve metadata
	meta, err := metadata.Resolve(url, d.log)
	if err != nil {
		return err
	}

	// 2️⃣ Build paths (single source of truth)
	meta.BuildPaths(output, "mp4", d.log)

	// 3️⃣ Prepare temp path
	tempDir := filepath.Join(meta.RootDir, ".temp")
	tempFile := filepath.Join(tempDir, filepath.Base(meta.MediaFile))

	d.log.Info("Final output file: " + meta.MediaFile)
	d.log.Debug("Temporary file: " + tempFile)

	// 4️⃣ Dispatch by URL type
	switch {
	case strings.HasSuffix(url, ".m3u8"):
		d.log.Info("Detected direct M3U8 link")
		tempDir := meta.RootDir + ".temp" // temp dir logic
		return m3u8.Download(m3u8.Options{
			URL:        url,
			Output:     meta.MediaFile,
			TempDir:    tempDir,
			Resume:     resume,
			Concurrent: concurrent,
			Log:        d.log,
		})

	case strings.Contains(url, "youtube.com"),
		strings.Contains(url, "youtu.be"):
		d.log.Info("Detected YouTube URL")
		return errors.New("YouTube downloader not implemented yet")

	case strings.Contains(url, "moviesbazar"):
		d.log.Info("Detected moviesbazar webpage")
		m3u8URL, err := moviebazar.ExtractM3U8(url, d.log)
		if err != nil {
			return err
		}
		if m3u8URL == "" {
			return errors.New("no m3u8 found on webpage")
		}

		return m3u8.Download(m3u8.Options{
			URL:        m3u8URL,
			Output:     meta.MediaFile,
			TempDir:    tempDir,
			Resume:     resume,
			Concurrent: concurrent,
			Log:        d.log,
		})
	}

	return errors.New("unsupported URL format")
}
