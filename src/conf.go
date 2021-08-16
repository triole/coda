package main

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/pelletier/go-toml"
)

type tCoda struct {
	FileTypes []tFileType `toml:"ft"`
	Settings  tSettings   `toml:"settings"`
	VarMap    map[string]string
}

type tFileType struct {
	Name  string `toml:"name"`
	Regex string
	Cmds  [][]string
}

type tSettings struct {
	IgnoreList []string
}

func readConfig(filename string) (coda tCoda) {
	if filename != "" {
		var err error
		raw, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatalf("Error reading general config %q, %q", filename, err)
		}
		err = toml.Unmarshal(raw, &coda)
		if err != nil {
			log.Fatalf("Error unmarshal %q, %q", filename, err)
		}
	}
	coda.VarMap = makeVarMap(CLI.Filename)
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
