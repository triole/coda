package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alecthomas/kong"
)

var (
	// BUILDTAGS are injected ld flags during build
	BUILDTAGS      string
	appName        = "coda"
	appDescription = "coda brings the music back to coding"
	appMainversion = "0.1"
)

var CLI struct {
	Filename    string `help:"file to process, positional arg required" arg optional`
	Config      string `help:"configuration file" short:c default:${config}`
	PrintVars   bool   `help:"print available vars" short:p`
	Debug       bool   `help:"debug mode" short:d`
	VersionFlag bool   `help:"display version" short:V`
}

func parseArgs() {
	curdir, _ := os.Getwd()
	ctx := kong.Parse(&CLI,
		kong.Name(appName),
		kong.Description(appDescription),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact:      true,
			Summary:      true,
			NoAppSummary: true,
			FlagsLast:    false,
		}),
		kong.Vars{
			"curdir": curdir,
			"config": path.Join(getBindir(), appName+".toml"),
		},
	)
	_ = ctx.Run()

	if CLI.VersionFlag == true {
		printBuildTags(BUILDTAGS)
		os.Exit(0)
	}
	if CLI.PrintVars == true {
		printAvailableVars()
		os.Exit(0)
	}
	if CLI.Filename == "" {
		fmt.Printf("%s\n", "Error: Positional arg expected. Please pass file name.")
		os.Exit(1)
	}
	// ctx.FatalIfErrorf(err)
}

func printBuildTags(buildtags string) {
	regexp, _ := regexp.Compile(`({|}|,)`)
	s := regexp.ReplaceAllString(buildtags, "\n")
	s = strings.Replace(s, "_subversion: ", "version: "+appMainversion+".", -1)
	arr := strings.Split(s, "\n")
	fmt.Printf("\n%s, %s\n", appName, appDescription)
	for _, el := range arr {
		if el != "" {
			fmt.Printf("%s\n", strings.TrimSpace(el))
		}
	}
}

func alnum(s string) string {
	s = strings.ToLower(s)
	re := regexp.MustCompile("[^a-z0-9_-]")
	return re.ReplaceAllString(s, "-")
}

func getBindir() (s string) {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	s = filepath.Dir(ex)
	return
}
