package cmd

import (
	"fmt"

	"github.com/ajaysinghnp/maya-cli/internal/downloader"
	"github.com/ajaysinghnp/maya-cli/internal/logger"
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download <url>",
	Short: "Download movies or series from a given URL",
	Long: `Download is a flexible command to fetch movies or TV series
from supported sources. It handles:

  - Metadata extraction in Jellyfin-friendly format
  - Episode organization for series
  - Resumable downloads using temporary files
  - M3U8 playlist handling
  - Parallel downloads to speed up series downloads

Examples:
  # Download a single movie
  maya download https://example.com/movie123

  # Download a series
  maya download https://example.com/series456
`,
	Args: cobra.MinimumNArgs(1), // requires at least one argument (the URL)
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]

		verbose, _ := cmd.Flags().GetBool("verbose")
		output, _ := cmd.Flags().GetString("output")
		resume, _ := cmd.Flags().GetBool("resume")
		concurrency, _ := cmd.Flags().GetInt("concurrency")
		log := logger.New(verbose, "")
		defer log.Close()

		if verbose {
			log.Debug("Starting download for URL: " + url)
			log.Debug("Output: " + output)
			log.Debug(fmt.Sprintf("Resume: %v", resume))
			log.Debug(fmt.Sprintf("Concurrency: %d", concurrency))
		}

		// Here you can call your downloader package logic
		// downloader.StartDownload(url, output, resume, concurrency)
		log.Info("Analyzing URL: " + url)

		// Create downloader
		dl := downloader.New(log)
		err := dl.StartDownload(url, output, resume, concurrency)
		if err != nil {
			log.Error(fmt.Sprintf("Download failed: %v", err))
			return
		}

		log.Info("Download completed successfully!")
	},
}

func init() {
	// Attach the download command to root
	rootCmd.AddCommand(downloadCmd)

	// Flags specific to the download command
	downloadCmd.Flags().BoolP("resume", "r", true, "Resume an interrupted download if cached files exist")
	downloadCmd.Flags().StringP("output", "o", "", "Specify output directory or filename (default: auto-generated)")
	downloadCmd.Flags().IntP("concurrency", "c", 5, "Number of simultaneous downloads for series episodes")
}
