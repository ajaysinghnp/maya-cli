package logger

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type Logger struct {
	verbose   bool
	file      *os.File
	writeFile bool
	mu        sync.Mutex
}

// ANSI background colors for label
var levelColors = map[string]string{
	"INFO":  "\033[1;37;44m", // white text on blue bg
	"WARN":  "\033[1;30;43m", // black text on yellow bg
	"ERROR": "\033[1;37;41m", // white text on red bg
	"DEBUG": "\033[1;37;45m", // white text on magenta bg
}

var resetColor = "\033[0m"

// New creates a new logger instance
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

func (l *Logger) Close() {
	if l.writeFile && l.file != nil {
		l.file.Close()
	}
}

func (l *Logger) logMessage(level, msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	color, ok := levelColors[level]
	if !ok {
		color = ""
	}

	// Colored label + timestamp + message
	formatted := fmt.Sprintf("%s[%s]%s %s %s\n", color, level, resetColor, timestamp, msg)

	// Terminal
	fmt.Print(formatted)

	// File (without ANSI colors)
	if l.writeFile && l.file != nil {
		l.file.WriteString(fmt.Sprintf("[%s] %s %s\n", level, timestamp, msg))
	}
}

func (l *Logger) Info(msg string)  { l.logMessage("INFO", msg) }
func (l *Logger) Warn(msg string)  { l.logMessage("WARN", msg) }
func (l *Logger) Error(msg string) { l.logMessage("ERROR", msg) }
func (l *Logger) Debug(msg string) {
	if l.verbose {
		l.logMessage("DEBUG", msg)
	}
}
