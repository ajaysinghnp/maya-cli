package downloader

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ajaysinghnp/maya-cli/internal/downloader/m3u8"
	moviebazar "github.com/ajaysinghnp/maya-cli/internal/downloader/movie-bazar"
	"github.com/ajaysinghnp/maya-cli/internal/logger"
)

type Downloader struct {
	log *logger.Logger
}

func New(log *logger.Logger) *Downloader {
	return &Downloader{log: log}
}

func downloadM3U8(url string, concurrent int, log *logger.Logger) error {
	return m3u8.Download(url, concurrent, log)
}

func extractAndDownloadFromWebpage(url string, concurrent int, log *logger.Logger) error {
	m3u8Url, err := moviebazar.ExtractM3U8(url, log)
	if err != nil {
		return err
	}
	if m3u8Url == "" {
		return errors.New("No M3U8 found on webpage")
	}
	return downloadM3U8(m3u8Url, concurrent, log)
}

// StartDownload decides which module to use based on URL
func (d *Downloader) StartDownload(url string, output string, resume bool, concurrent int) error {
	d.log.Info("Getting metadata for URL: " + url)
	d.log.Info(fmt.Sprintf("Output: %s | Resume: %v | Concurrency: %d", output, resume, concurrent))

	// Simplest URL inspection
	if strings.HasSuffix(url, ".m3u8") {
		d.log.Info("Detected direct M3U8 link.")
		return downloadM3U8(url, concurrent, d.log)
	}

	if strings.Contains(url, "youtube.com") || strings.Contains(url, "youtu.be") {
		d.log.Info("Detected YouTube URL.")
		// placeholder for youtube downloader
		return errors.New("YouTube downloader not implemented yet")
	}

	if strings.HasPrefix(url, "moviesbazar") {
		d.log.Info("Detected moviesbazar URL, attempting to extract M3U8.")
		return extractAndDownloadFromWebpage(url, concurrent, d.log)
	}

	return errors.New("Unsupported URL format")
}
