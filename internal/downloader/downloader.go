package downloader

import (
	"errors"
	"path/filepath"

	"github.com/ajaysinghnp/maya-cli/internal/downloader/m3u8"
	moviebazar "github.com/ajaysinghnp/maya-cli/internal/downloader/movie-bazar"
	"github.com/ajaysinghnp/maya-cli/internal/logger/iface"
	"github.com/ajaysinghnp/maya-cli/internal/metadata/resolver"
)

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
	// Detect source
	source := resolver.DetectSource(url)

	// 1️⃣ Resolve metadata
	meta, err := resolver.Resolve(source, url, d.log)
	if err != nil {
		return err
	}
	d.log.Success("Metadata resolved successfully!")

	// 2️⃣ Build paths (single source of truth)
	meta.BuildPaths(output, "mp4", d.log)
	d.log.Info("Paths prepared for download.")

	// 3️⃣ Prepare temp path
	tempDir := filepath.Join(meta.RootDir, ".temp")
	tempFile := filepath.Join(tempDir, filepath.Base(meta.MediaFile))
	d.log.Debug("Temporary file path: " + tempFile)
	d.log.Info("Final output file: " + meta.MediaFile)

	// 4️⃣ Dispatch by source type
	switch source {
	case resolver.SourceM3U8:
		d.log.Info("Detected direct M3U8 link.")
		return m3u8.Download(m3u8.Options{
			URL:        url,
			Output:     meta.MediaFile,
			TempDir:    tempDir,
			Resume:     resume,
			Concurrent: concurrent,
			Log:        d.log,
		})

	case resolver.SourceMoviesBazar:
		d.log.Info("Detected MoviesBazar webpage, extracting M3U8...")
		m3u8URL, err := moviebazar.ExtractM3U8(url, d.log)
		if err != nil {
			d.log.Error("Failed to extract M3U8: " + err.Error())
			return err
		}
		if m3u8URL == "" {
			d.log.Warn("No M3U8 found on webpage")
			return errors.New("no m3u8 found on webpage")
		}
		d.log.Success("M3U8 URL extracted successfully!")
		return m3u8.Download(m3u8.Options{
			URL:        m3u8URL,
			Output:     meta.MediaFile,
			TempDir:    tempDir,
			Resume:     resume,
			Concurrent: concurrent,
			Log:        d.log,
		})

	case resolver.SourceYouTube:
		d.log.Info("Detected YouTube URL")
		return errors.New("YouTube downloader not implemented yet")

	default:
		d.log.Warn("Unsupported URL format")
		return errors.New("unsupported URL format")
	}
}
