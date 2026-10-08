package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func TestDetect(t *testing.T) {
	files := getFiles("../testdata", t)
	for _, filename := range files {
		parentDir := filepath.Dir(filename)
		expectedName := strings.TrimSuffix(filepath.Base(parentDir), "/")
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

func getFiles(p string, t *testing.T) (files []string) {
	root, err := filepath.Abs(p)
	if err != nil {
		t.Fatalf("can not make absolute file path: %v", err)
	}
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			t.Errorf("can not walk over files: %v", err)
			return nil
		}

		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		t.Fatalf("an error occurred: %v", err)
	}
	sort.Strings(files)
	return
}
