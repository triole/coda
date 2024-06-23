package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestDetect(t *testing.T) {
	files := getFiles("../testdata")
	for _, filename := range files {
		expectedName := find(`[^/]+/[^/]+$`, filename)
		expectedName = find(`^.*?/`, expectedName)
		expectedName = expectedName[0 : len(expectedName)-1]
		coda := initCoda("../testdata/yaml/conf.yaml", filename)
		ft := coda.detect()
		if expectedName != ft.Name {
			t.Errorf(
				"assertion detect failed: %q, %q != %q",
				filename, expectedName, ft.Name,
			)
		}
	}
}

func getFiles(p string) (files []string) {
	root, err := filepath.Abs(p)
	if err != nil {
		fmt.Printf("can not make absolute file path: %s\n", err)
		os.Exit(1)
	}
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			fmt.Printf("can not walk over files: %s\n", err)
			return nil
		}

		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		log.Fatal("[coda] error: ", err)
	}
	sort.Strings(files)
	return
}
