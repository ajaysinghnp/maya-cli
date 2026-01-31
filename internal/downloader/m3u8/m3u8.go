package m3u8

import (
	"fmt"
	"io"
	"net/http"
)

func Download(opts Options) error {
	log := opts.Log
	log.Info("Starting M3U8 download: " + opts.URL)
	log.Debug("Output file: " + opts.Output)
	log.Info(fmt.Sprintf("Temp dir: %s | Resume: %v | Concurrency: %d", opts.TempDir, opts.Resume, opts.Concurrent))

	log.Info("Requesting M3U8 playlist...")
	resp, err := http.Get(opts.URL)
	if err != nil {
		return fmt.Errorf("failed to fetch playlist: %w", err)
	}
	defer resp.Body.Close()

	log.Success("M3U8 playlist requested successfully!")

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to fetch playlist, status: %d", resp.StatusCode)
	}

	log.Success("M3U8 playlist fetched successfully!")

	playlistData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read playlist: %w", err)
	}

	log.Debug("Playlist data length: " + fmt.Sprint(len(playlistData)))

	// TODO:
	// 1. Fetch playlist
	// 2. Parse segments
	// 3. Download segments concurrently
	// 4. Merge segments into opts.Output
	// 5. Handle resume logic via tempDir
	return nil
}
