package main

import (
	"testing"
)

func TestIsFile(t *testing.T) {
	assertIsFile("/tmp", false, t)
	assertIsFile("/does_not_exist", false, t)
	assertIsFile("main.go", true, t)
	assertIsFile("does_not_exist.go", false, t)
}

func assertIsFile(filePath string, exp bool, t *testing.T) {
	if isFile(filePath) != exp {
		t.Errorf("isFile check failed: %q did not meet %v", filePath, exp)
	}
}
