package main

import (
	"strings"
)

func detectByExtension(filename string) (r string) {
	var ext string
	extArr := strings.Split(filename, ".")
	if len(extArr) > 1 {
		ext = extArr[len(extArr)-1]
	}
	switch strings.ToLower(ext) {
	case "go":
		r = "golang"
	case "js":
		r = "javascript"
	case "json":
		r = "json"
	case "md":
		r = "markdown"
	case "py":
		r = "python"
	case "rs":
		r = "rust"
	}
	return
}
