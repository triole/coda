package main

import (
	"fmt"
	"regexp"
	"time"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

type tSettings struct {
	Command    []string
	Interval   time.Duration
	Folder     string
	Regex      *regexp.Regexp
	Spectate   bool
	KeepOutput bool
	Logging    *logrus.Logger
	LogInit    bool
}

func main() {
	parseArgs()
	settings := tSettings{
		Command:    CLI.Command,
		Folder:     CLI.Folder,
		Interval:   time.Duration(CLI.Interval) * time.Second,
		Regex:      regexp.MustCompile(CLI.Regex),
		Spectate:   CLI.Spectate,
		KeepOutput: CLI.KeepOutput,
		LogInit:    false,
	}
	if len(settings.Command) < 1 {
		settings.Spectate = true
	}

	mode := fmt.Sprintf("command on change: %q", settings.Command)
	if settings.Spectate == true {
		mode = "just spectate"
	}

	color.Green("\nWatch folder %q, %s", settings.Folder, mode)
}
