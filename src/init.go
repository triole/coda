package main

import (
	"log"
	"os"

	"github.com/pelletier/go-toml"
)

type tCoda struct {
	FileTypes     []tFileType `toml:"ft"`
	Settings      tSettings   `toml:"settings"`
	FileConfig    string
	FileToProcess string
	VarMap        tVarMap
}

type tFileType struct {
	Name    string `toml:"name"`
	Shebang string `toml:"shebang"`
	Regex   string
	Cmds    [][]string
}

type tSettings struct {
	IgnoreList []string
}

func initCoda(fileConfig, fileToProcess string) (coda tCoda) {
	coda.FileConfig = fileConfig
	coda.FileToProcess = fileToProcess
	if fileConfig != "" {
		var err error
		raw, err := os.ReadFile(fileConfig)
		if err != nil {
			log.Fatalf("Error reading config %q, %q", fileConfig, err)
		}
		err = toml.Unmarshal(raw, &coda)
		if err != nil {
			log.Fatalf("Error unmarshal %q, %q", fileConfig, err)
		}
	}
	coda.VarMap = makeVarMap(fileToProcess)
	return
}
