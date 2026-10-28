package main

import (
	"os"
	"sync"
	"text/template"

	yaml "gopkg.in/yaml.v3"
)

type tCoda struct {
	FileTypes     []tFileType `yaml:"filetypes"`
	Settings      tSettings   `yaml:"settings"`
	FileConfig    string
	FileToProcess string
	VarMap        tVarMap
	tmplCache     map[string]*template.Template
	cacheMu       sync.RWMutex
}

type tSettings struct {
	IgnoreList []string
}

func initCoda(fileConfig, fileToProcess string) (coda tCoda) {
	coda.FileConfig = fileConfig
	coda.FileToProcess = fileToProcess
	coda.tmplCache = make(map[string]*template.Template, 8)
	if coda.FileConfig != "" {
		var err error
		raw, err := os.ReadFile(coda.FileConfig)
		if err != nil {
			log.Fatal("error reading config", map[string]any{"file": coda.FileConfig, "error": err})
		}
		err = yaml.Unmarshal(raw, &coda)
		if err != nil {
			log.Fatal("unmarshal error", map[string]any{"file": coda.FileConfig, "error": err})
		}
	}
	coda.VarMap = makeVarMap(fileToProcess)
	compileRegexes(&coda.FileTypes)
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
