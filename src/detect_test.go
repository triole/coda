package main

import "testing"

func TestDetectByExtension(t *testing.T) {
	assertDetectByExtension("main.go", "golang", t)
	assertDetectByExtension("main.py", "python", t)
	assertDetectByExtension("main.md", "markdown", t)
	assertDetectByExtension("main.rs", "rust", t)
}

func assertDetectByExtension(filename string, exprectedResult string, t *testing.T) {
	res := detectByExtension(filename)
	if res != exprectedResult {
		t.Errorf(
			"Assertion detect by extension failed: %q, %q != %q",
			filename, res, exprectedResult,
		)
	}
}
