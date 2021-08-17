package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

var (
	coda = readConfig("../conf/coda.toml")
)

func TestDetect(t *testing.T) {
	files := getFiles("../testdata")
	for _, filename := range files {
		expectedType := find(`[^/]+/[^/]+$`, filename)
		expectedType = find(`^.*?/`, expectedType)
		expectedType = expectedType[0 : len(expectedType)-1]
		assertDetect(filename, expectedType, t)
	}
}

func assertDetect(filename string, expectedName string, t *testing.T) {
	ft := coda.detect(filename)
	if expectedName != ft.Name {
		t.Errorf(
			"Assertion detect failed: %q, %q != %q",
			filename, expectedName, ft.Name,
		)
	}
}

func getFiles(p string) (files []string) {
	root, err := filepath.Abs(p)
	if err != nil {
		fmt.Printf("Can not make absolute file path: %s\n", err)
		os.Exit(1)
	}
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			fmt.Printf("Can not walk over files: %s\n", err)
			return nil
		}

		if !info.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
	sort.Strings(files)
	return
}
