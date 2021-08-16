package main

import (
	"regexp"
)

func (coda tCoda) detect(filename string) (b bool, r [][]string) {
	for _, filetype := range coda.FileTypes {
		b, r = coda.detectByRegex(filename, filetype)
		if b == true {
			break
		}
	}
	return
}

func (coda tCoda) detectByRegex(filename string, filetype tFileType) (b bool, r [][]string) {
	rx := regexp.MustCompile(filetype.Regex)
	b = false
	if rx.MatchString(filename) == true {
		b = true
		r = filetype.Cmds
	}
	return
}
