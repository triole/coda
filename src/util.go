package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

func pprint(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	println(string(b))
}

func getFirstLineOfFile(filename string) (l string) {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		l = scanner.Text()
		break
	}
	return
}
