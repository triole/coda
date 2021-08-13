package main

import (
	"io/ioutil"
	"log"

	"github.com/pelletier/go-toml"
)

type tConfig struct {
	FileTypes []tFileType `toml:"ft"`
	Settings  tSettings   `toml:"settings"`
}

type tFileType struct {
	Name  string `toml:"name"`
	Regex string
	Cmds  []string
}

type tSettings struct {
	IgnoreList []string
}

func readConfig(filename string) (config tConfig) {
	if filename != "" {
		var err error
		raw, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatalf("Error reading general config %q, %q", filename, err)
		}
		err = toml.Unmarshal(raw, &config)
		if err != nil {
			log.Fatalf("Error unmarshal %q, %q", filename, err)
		}
	}
	return
}
