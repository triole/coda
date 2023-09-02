package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
	"gopkg.in/yaml.v3"
)

type tCoda struct {
	FileTypes     []tFileType `toml|yaml:"filetypes"`
	Settings      tSettings   `toml|yaml:"settings"`
	FileConfig    string
	FileToProcess string
	VarMap        tVarMap
}

type tFileType struct {
	Name    string     `toml|yaml:"name"`
	Shebang string     `toml|yaml:"shebang"`
	Regex   string     `toml|yaml:"regex"`
	Cmds    [][]string `toml|yaml:"cmds"`
}

type tSettings struct {
	IgnoreList []string
}

func initCoda(fileConfig, fileToProcess string) (coda tCoda) {
	coda.FileConfig = fileConfig
	coda.FileToProcess = fileToProcess
	if coda.FileConfig != "" {
		var err error
		raw, err := os.ReadFile(coda.FileConfig)
		if err != nil {
			log.Fatalf("Error reading config %q, %q", coda.FileConfig, err)
		}
		ext := filepath.Ext(coda.FileConfig)
		if ext == ".toml" {
			err = toml.Unmarshal(raw, &coda)
		}
		if ext == ".yaml" {
			err = yaml.Unmarshal(raw, &coda)
		}
		if err != nil {
			log.Fatalf("unmarshal error %q, %q", coda.FileConfig, err)
		}
	}
	coda.VarMap = makeVarMap(fileToProcess)
	return
}

func returnFirstExistingFile(arr []string) (s string) {
	for _, el := range arr {
		if isFile(el) {
			s = el
			break
		}
	}
	return
}
