package main

import (
	"regexp"
)

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
