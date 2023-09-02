package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/alecthomas/kong"
)

var (
	// BUILDTAGS are injected ld flags during build
	BUILDTAGS      string
	appName        = "coda"
	appDescription = "brings the music back to coding"
	appMainversion = "0.1"
)

var CLI struct {
	Filename    string `help:"file to process, required positional arg" arg:"" optional:""`
	Config      string `help:"configuration file" short:"c" default:"${config}"`
	PrintVars   bool   `help:"print available vars" short:"p"`
	DryRun      bool   `help:"dry run, just print don't do" short:"n"`
	VersionFlag bool   `help:"display version" short:"V"`
}

func parseArgs() {
	curdir, _ := os.Getwd()
	homeFolder := getHome()
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
			"config": returnFirstExistingFile(
				[]string{
					path.Join(getBindir(), appName+".toml"),
					path.Join(homeFolder, ".conf", "coda", "conf.yaml"),
					path.Join(homeFolder, ".conf", "coda", "conf.toml"),
					path.Join(homeFolder, ".config", "coda", "conf.yaml"),
					path.Join(homeFolder, ".config", "coda", "conf.toml"),
				},
			),
		},
	)
	_ = ctx.Run()

	if CLI.VersionFlag {
		printBuildTags(BUILDTAGS)
		os.Exit(0)
	}
	if CLI.PrintVars {
		printAvailableVars()
		os.Exit(0)
	}
	if CLI.Filename == "" {
		fmt.Printf("%s\n", "Error: Positional arg expected. Please pass file name.")
		os.Exit(1)
	}
	// ctx.FatalIfErrorf(err)
}

type tPrinter []tPrinterEl
type tPrinterEl struct {
	Key string
	Val string
}

func printBuildTags(buildtags string) {
	regexp, _ := regexp.Compile(`({|}|,)`)
	s := regexp.ReplaceAllString(buildtags, "\n")
	s = strings.Replace(s, "_subversion: ", "version: "+appMainversion+".", -1)
	fmt.Printf("\n%s\n%s\n\n", appName, appDescription)
	arr := strings.Split(s, "\n")
	var pr tPrinter
	var maxlen int
	for _, line := range arr {
		if strings.Contains(line, ":") {
			l := strings.Split(line, ":")
			if len(l[0]) > maxlen {
				maxlen = len(l[0])
			}
			pr = append(pr, tPrinterEl{l[0], strings.Join(l[1:], ":")})
		}
	}
	for _, el := range pr {
		fmt.Printf("%"+strconv.Itoa(maxlen)+"s\t%s\n", el.Key, el.Val)
	}
	fmt.Printf("\n")
}

// func alnum(s string) string {
// 	s = strings.ToLower(s)
// 	re := regexp.MustCompile("[^a-z0-9_-]")
// 	return re.ReplaceAllString(s, "-")
// }

func getBindir() (s string) {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	s = filepath.Dir(ex)
	return
}
