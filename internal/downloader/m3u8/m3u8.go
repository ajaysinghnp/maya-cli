package m3u8

import "fmt"

func Download(opts Options) error {
	opts.Log.Info("Starting M3U8 download: " + opts.URL)
	opts.Log.Debug("Output file: " + opts.Output)
	opts.Log.Debug(fmt.Sprintf("Temp dir: %s | Resume: %v | Concurrency: %d", opts.TempDir, opts.Resume, opts.Concurrent))

	// TODO:
	// 1. Fetch playlist
	// 2. Parse segments
	// 3. Download segments concurrently
	// 4. Merge segments into opts.Output
	// 5. Handle resume logic via tempDir
	return nil
}
