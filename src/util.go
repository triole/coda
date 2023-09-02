package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

func pprint(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	println(string(b))
}

func getFirstLineOfFile(filename string) (l string) {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("error reading file %q\n", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		l = scanner.Text()
		break
	}
	return
}

func makeAbs(filename string) string {
	filename, err := filepath.Abs(filename)
	if err != nil {
		fmt.Printf("can not assemble absolute filename %q\n", err)
		os.Exit(1)
	}
	return filename
}

func isFile(filePath string) bool {
	stat, err := os.Stat(makeAbs(filePath))
	if !os.IsNotExist(err) && !stat.IsDir() {
		return true
	}
	return false
}

func getHome() string {
	usr, err := user.Current()
	if err != nil {
		fmt.Printf("unable to retrieve user's home folder")
	}
	return usr.HomeDir
}
