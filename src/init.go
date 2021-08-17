package main

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/pelletier/go-toml"
)

type tCoda struct {
	FileTypes     []tFileType `toml:"ft"`
	Settings      tSettings   `toml:"settings"`
	FileConfig    string
	FileToProcess string
	VarMap        map[string]string
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
		raw, err := ioutil.ReadFile(fileConfig)
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

func makeVarMap(filename string) (varmap map[string]string) {
	varmap = make(map[string]string)
	varmap["filename"] = filename
	varmap["shortname"] = find(`[^/]+$`, filename)
	if strings.Contains(filename, ".") == true {
		arr := strings.Split(filename, ".")
		varmap["ext"] = arr[len(arr)-1]
	}
	varmap["filename_no_ext"] = strings.Replace(
		filename, "."+varmap["ext"], "", -1,
	)
	varmap["shortname_no_ext"] = strings.Replace(
		varmap["shortname"], "."+varmap["ext"], "", -1,
	)
	return
}
