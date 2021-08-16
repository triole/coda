package main

import (
	"testing"
)

var (
	coda = readConfig("../conf/coda.toml")
)

func TestDetect(t *testing.T) {
	assertDetect("main.go", "go", t)
	assertDetect("main.py", "autoflake", t)
	assertDetect("main.md", "mdtoc", t)
	assertDetect("main.rs", "rustfmt", t)
}

func assertDetect(filename string, expectedCmd string, t *testing.T) {
	_, cmds := coda.detect(filename)
	if expectedCmd != cmds[0][0] {
		t.Errorf(
			"Assertion detect by extension failed: %q, %q",
			filename, cmds,
		)
	}
}
