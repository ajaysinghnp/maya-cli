package logger

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/ajaysinghnp/maya-cli/internal/logger/iface"
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

// Implement the iface.Logger interface
var _ iface.Logger = (*Logger)(nil)
