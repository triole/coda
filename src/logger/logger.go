// Package logger provides a minimal, zero-dependency logging implementation
// using Go's standard library. Supports log levels, basic file output,
// and timestamp formatting.

package logger

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

// Level represents log severity
type Level int

const (
	TraceLevel Level = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// Logger wraps Go's standard log package with level filtering
type Logger struct {
	mu     sync.Mutex
	level  Level
	out    io.Writer
	file   *os.File // file handle for cleanup, nil if stdout/stderr
	colors bool
	closed bool // track if already closed
	format string
}

// NewLogger creates a logger with the given level and output.
// Returns a logger that must be closed with Close() to properly release file handles.
func NewLogger(levelStr string, logFile string, noColors bool) *Logger {
	// Parse log level
	level := InfoLevel
	switch strings.ToLower(levelStr) {
	case "trace":
		level = TraceLevel
	case "debug":
		level = DebugLevel
	case "info":
		level = InfoLevel
	case "warn":
		level = WarnLevel
	case "error":
		level = ErrorLevel
	case "fatal":
		level = FatalLevel
	}

	// Set output destination
	output := os.Stdout
	file := (*os.File)(nil)
	if logFile != "" && logFile != "stdout" && logFile != "/dev/stdout" {
		f, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "can't open log file: %s\n", err)
			os.Exit(1)
		}
		// Ensure file is closed if creation fails later
		output = f
		file = f // store file handle for later cleanup
	}

	format := "%-14s %s %-32s %s\n"
	if noColors {
		format = "%-5s %s %-32s %s\n"
	}

	return &Logger{
		level:  level,
		out:    output,
		file:   file,
		colors: !noColors,
		closed: false,
		format: format,
	}
}

// formatTime returns timestamp string
func formatTime() string {
	return time.Now().Format("[2006-01-02 15:04:05.000 MST]")
}

func alphaIterator(fields map[string]any) (iterator []string) {
	iterator = make([]string, len(fields))
	i := 0
	for k := range fields {
		iterator[i] = k
		i++
	}
	sort.Strings(iterator)
	return
}

// formatFields converts key-value pairs to string
func formatFields(fields map[string]any) string {
	if len(fields) == 0 {
		return ""
	}
	parts := make([]string, 0, len(fields))
	for _, it := range alphaIterator(fields) {
		parts = append(parts, fmt.Sprintf("%s=%v", it, fields[it]))
	}
	return strings.Join(parts, " ")
}

// colorPrefix returns ANSI color codes if colors enabled
func (l *Logger) colorPrefix(level string) string {
	if !l.colors {
		return level
	}

	colors := map[string]string{
		"TRACE": "\033[37mTRACE\033[0m", // grey
		"DEBUG": "\033[36mDEBUG\033[0m", // cyan
		"INFO":  "\033[32mINFO\033[0m",  // green
		"WARN":  "\033[33mWARN\033[0m",  // yellow
		"ERROR": "\033[31mERROR\033[0m", // red
		"FATAL": "\033[31mFATAL\033[0m", // red
	}

	if c, ok := colors[level]; ok {
		return c
	}
	return "[" + level + "]"
}

// Trace logs a message at trace level
func (l *Logger) Trace(msg string, fields map[string]any) {
	if l.level <= TraceLevel {
		l.mu.Lock()
		defer l.mu.Unlock()
		fmt.Fprintf(
			l.out, l.format,
			l.colorPrefix("TRACE"), formatTime(), msg, formatFields(fields),
		)
	}
}

// Debug logs a message at debug level
func (l *Logger) Debug(msg string, fields map[string]any) {
	if l.level <= DebugLevel {
		l.mu.Lock()
		defer l.mu.Unlock()
		fmt.Fprintf(
			l.out, l.format,
			l.colorPrefix("DEBUG"), formatTime(), msg, formatFields(fields),
		)
	}
}

// Info logs a message at info level
func (l *Logger) Info(msg string, fields map[string]any) {
	if l.level <= InfoLevel {
		l.mu.Lock()
		defer l.mu.Unlock()
		fmt.Fprintf(
			l.out, l.format,
			l.colorPrefix("INFO"), formatTime(), msg, formatFields(fields),
		)
	}
}

// Warn logs a message at warn level
func (l *Logger) Warn(msg string, fields map[string]any) {
	if l.level <= WarnLevel {
		l.mu.Lock()
		defer l.mu.Unlock()
		fmt.Fprintf(
			l.out, l.format,
			l.colorPrefix("WARN"), formatTime(), msg, formatFields(fields),
		)
	}
}

// Error logs a message at error level
func (l *Logger) Error(msg string, fields map[string]any) {
	if l.level <= ErrorLevel {
		l.mu.Lock()
		defer l.mu.Unlock()
		fmt.Fprintf(
			l.out, l.format,
			l.colorPrefix("ERROR"), formatTime(), msg, formatFields(fields),
		)
	}
}

// closeInternal handles the actual file closing without locking.
// Should only be called when the mutex is already held.
func (l *Logger) closeInternal() error {
	if l.closed {
		return nil // already closed
	}
	l.closed = true

	if l.file != nil {
		_ = l.file.Sync() // flush to disk, ignore error as we're closing
		err := l.file.Close()
		l.file = nil // nil out reference to prevent double-close
		return err
	}
	return nil
}

// Close closes the file handle if a log file was opened.
// This should be called before program exit to ensure proper cleanup.
// Multiple calls to Close are safe and idempotent.
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.closeInternal()
}

// IfErrFatal logs an error message and exits if err is non-nil
func (l *Logger) IfErrFatal(msg string, err error, fields map[string]any) {
	if err != nil {
		if fields == nil {
			fields = make(map[string]any)
		}
		fields["error"] = err
		l.Fatal(msg, fields)
	}
}

// Fatal logs a message at fatal level, flushes, closes file handles, and exits.
func (l *Logger) Fatal(msg string, fields map[string]any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.level <= FatalLevel {
		fmt.Fprintf(
			l.out, l.format,
			l.colorPrefix("FATAL"), formatTime(), msg, formatFields(fields),
		)
		// Flush output before exit
		if f, ok := l.out.(interface{ Sync() error }); ok {
			_ = f.Sync() // ignore error
		}
	}

	// Always close file handle before exit to prevent file handle leak
	// This must be unconditional to ensure cleanup even when log level is higher than Fatal
	_ = l.closeInternal()
	os.Exit(1)
}
