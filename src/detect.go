package main

import (
	"regexp"
)

type tFileType struct {
	Name          string     `yaml:"name"`
	Shebang       string     `yaml:"shebang"`
	Regex         string     `yaml:"regex"`
	RegexIgnore   string     `yaml:"regex_ignore"`
	Cmds          [][]string `yaml:"cmds"`
	WriteStdoutTo string     `yaml:"write_stdout_to"`
}

func (coda tCoda) detect() (ft tFileType) {
	for _, filetype := range coda.FileTypes {
		ft = coda.detectByRegex(coda.FileToProcess, filetype)
		if ft.Name != "" {
			return
		}
	}

	for _, filetype := range coda.FileTypes {
		ft = coda.detectByShebang(coda.FileToProcess, filetype)
		if ft.Name != "" {
			return
		}
	}
	return
}

func (coda tCoda) detectByRegex(filename string, filetype tFileType) (ft tFileType) {
	rx := regexp.MustCompile(filetype.Regex)
	if filetype.RegexIgnore != "" {
		ign := regexp.MustCompile(filetype.RegexIgnore)
		if ign.MatchString(filename) {
			return tFileType{}
		}
	}
	if rx.MatchString(filename) {
		ft = filetype
	}
	return
}

func (coda tCoda) detectByShebang(filename string, filetype tFileType) (ft tFileType) {
	shebang := getFirstLineOfFile(filename)
	if shebang == filetype.Shebang {
		ft = filetype
	}
	return
}
