package main

import (
	"log"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type tCoda struct {
	FileTypes     []tFileType `yaml:"filetypes"`
	Settings      tSettings   `yaml:"settings"`
	FileConfig    string
	FileToProcess string
	VarMap        tVarMap
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
			log.Fatalf("[coda] error reading config %q, %q", coda.FileConfig, err)
		}
		err = yaml.Unmarshal(raw, &coda)
		if err != nil {
			log.Fatalf("[coda] unmarshal error %q, %q", coda.FileConfig, err)
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
