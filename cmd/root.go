package cmd

import (
	"os"

	"github.com/ajaysinghnp/maya-cli/internal/logger"
	"github.com/spf13/cobra"
)

var log *logger.Logger

var rootCmd = &cobra.Command{
	Use:   "maya",
	Short: "Maya - A Modular Multimedia CLI Tool",
	Long: `Maya is a modular command-line tool designed to simplify 
media management and downloading workflows.

With Maya, you can:
  - Download movies and series from supported sources
  - Extract and save metadata in Jellyfin-friendly format
  - Resume interrupted downloads seamlessly
  - Extend functionality easily by adding new subcommands

Usage Examples:
  maya download <url>           # Download a movie or series from a URL
  maya other-tool --option xyz   # Run other future tools
`,
	Version: Version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		PrintBanner()
		// Initialize logger per run, read verbose flag
		verbose, _ := cmd.Flags().GetBool("verbose")
		log = logger.New(verbose, "") // empty string = no file logging by default
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Use Success to greet the user
		log.Success("Welcome to Maya! Use `maya --help` to see available commands.")
	},
}

// Execute runs the root command and all subcommands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		PrintBanner()
		if log != nil {
			log.Error(err.Error())
		} else {
			// fallback if logger not initialized
			os.Stderr.WriteString(err.Error() + "\n")
		}
		os.Exit(1)
	}
}

func init() {
	// Persistent flags available to all commands
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
}
