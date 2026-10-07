package main

import (
	"regexp"
)

type tFileType struct {
	Name                string     `yaml:"name"`
	Shebang             string     `yaml:"shebang"`
	Regex               string     `yaml:"regex"`
	RegexIgnore         string     `yaml:"regex_ignore"`
	Cmds                [][]string `yaml:"cmds"`
	WriteStdoutTo       string     `yaml:"write_stdout_to"`
	compiledRegex       *regexp.Regexp
	compiledRegexIgnore *regexp.Regexp
}

// pre compile regexes improving performance
func compileRegexes(filetypes *[]tFileType) {
	for i := range *filetypes {
		ft := &(*filetypes)[i]
		if ft.Regex != "" {
			if rx, err := regexp.Compile(ft.Regex); err == nil {
				ft.compiledRegex = rx
			}
		}
		if ft.RegexIgnore != "" {
			if rx, err := regexp.Compile(ft.RegexIgnore); err == nil {
				ft.compiledRegexIgnore = rx
			}
		}
	}
}

func (coda tCoda) detect() (ft tFileType) {
	compileRegexes(&coda.FileTypes)
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
	if filetype.compiledRegexIgnore != nil {
		if filetype.compiledRegexIgnore.MatchString(filename) {
			return tFileType{}
		}
	}
	if filetype.compiledRegex != nil && filetype.compiledRegex.MatchString(filename) {
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
