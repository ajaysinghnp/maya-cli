package logger

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ajaysinghnp/maya-cli/internal/logger/iface"
	"github.com/ajaysinghnp/maya-cli/internal/metadata"
)

type Logger struct {
	verbose   bool
	file      *os.File
	writeFile bool
	mu        sync.Mutex
}

// ANSI color codes
var (
	reset   = "\033[0m"
	bold    = "\033[1m"
	black   = "\033[30m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	white   = "\033[37m"

	bgRed     = "\033[41m"
	bgGreen   = "\033[42m"
	bgYellow  = "\033[43m"
	bgBlue    = "\033[44m"
	bgMagenta = "\033[45m"
	bgCyan    = "\033[46m"
	bgWhite   = "\033[47m"
)

// New creates a new logger instance
// If filePath is empty, logger will only print to terminal
func New(verbose bool, filePath string) *Logger {
	var file *os.File
	writeFile := false

	if filePath != "" {
		var err error
		file, err = os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Failed to open log file: %v\n", err)
		} else {
			writeFile = true
		}
	}

	return &Logger{
		verbose:   verbose,
		file:      file,
		writeFile: writeFile,
	}
}

// Close closes the log file if any
func (l *Logger) Close() {
	if l.writeFile && l.file != nil {
		l.file.Close()
	}
}

func (l *Logger) logMessage(level, labelColor, textColor, msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	formatted := fmt.Sprintf("%s[%s]%s %s%s%s\n",
		labelColor, level, reset,
		textColor, timestamp+" "+msg, reset,
	)

	// Always log to terminal
	fmt.Print(formatted)

	// Only log to file if enabled
	if l.writeFile {
		l.file.WriteString(fmt.Sprintf("[%s] %s %s\n", level, timestamp, msg))
	}
}

func (l *Logger) Info(msg string) {
	l.logMessage("INFO", bgCyan+white, cyan, msg)
}

func (l *Logger) Debug(msg string) {
	if l.verbose {
		l.logMessage("DEBUG", bgBlue+white, blue, msg)
	}
}

func (l *Logger) Warn(msg string) {
	l.logMessage("WARN", bgYellow+black, yellow, msg)
}

func (l *Logger) Error(msg string) {
	l.logMessage("ERROR", bgRed+white, red, msg)
}

func (l *Logger) Success(msg string) {
	l.logMessage("SUCCESS", bgGreen+black, green, msg)
}

// LogMetadata prints a Metadata struct in a clean, colorful format.
func LogMetadata(m *metadata.Metadata) {
	fmt.Println(bold + cyan + "Parsed Media Info:" + reset)
	fmt.Println(cyan + "------------------" + reset)

	if m.Title != "" {
		fmt.Printf("%sTitle:%s        %s\n", yellow, reset, m.Title)
	}
	if m.Year != 0 {
		fmt.Printf("%sYear:%s         %d\n", yellow, reset, m.Year)
	}
	if m.Type != "" {
		fmt.Printf("%sType:%s         %s\n", yellow, reset, m.Type)
	}
	if m.Category != "" {
		fmt.Printf("%sCategory:%s     %s\n", yellow, reset, m.Category)
	}
	if m.ReleaseDate != "" {
		fmt.Printf("%sReleaseDate:%s  %s\n", yellow, reset, m.ReleaseDate)
	}
	if m.Language != "" {
		fmt.Printf("%sLanguage:%s     %s\n", yellow, reset, m.Language)
	}
	if len(m.Genres) > 0 {
		fmt.Printf("%sGenres:%s       %s\n", yellow, reset, strings.Join(m.Genres, ", "))
	}
	if len(m.Cast) > 0 {
		fmt.Printf("%sCast:%s         %s\n", yellow, reset, strings.Join(m.Cast, ", "))
	}
	if m.IDs.IMDB != "" || m.IDs.TMDB != "" {
		fmt.Println(bold + green + "IDs:" + reset)
		if m.IDs.IMDB != "" {
			fmt.Printf("  %sIMDB:%s %s\n", magenta, reset, m.IDs.IMDB)
		}
		if m.IDs.TMDB != "" {
			fmt.Printf("  %sTMDB:%s %s\n", magenta, reset, m.IDs.TMDB)
		}
	}
	if m.Thumbnail != "" {
		fmt.Printf("%sThumbnail:%s    %s\n", yellow, reset, m.Thumbnail)
	}
	if len(m.Sources) > 0 {
		fmt.Println(bold + green + "Sources:" + reset)
		for _, s := range m.Sources {
			fmt.Printf("  %s- %s:%s %s", cyan, s.Label, reset, s.URL)
			if s.LabelTag != "" {
				fmt.Printf(" (%s%s%s)", magenta, s.LabelTag, reset)
			}
			fmt.Println()
		}
	}
	if m.HlsSourceDomain != "" {
		fmt.Printf("%sHLS Domain:%s   %s\n", yellow, reset, m.HlsSourceDomain)
	}

	fmt.Println(cyan + "------------------" + reset)
}

// Implement the iface.Logger interface
var _ iface.Logger = (*Logger)(nil)
