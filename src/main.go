package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	parseArgs()

	fileConfig := makeAbs(CLI.Config)
	fileToProcess := makeAbs(CLI.Filename)

	coda := initCoda(fileConfig, fileToProcess)

	ft := coda.detect()

	if CLI.Debug == true {
		fmt.Printf("\n\n%s\n\n", "Used VarMap")
		pprint(coda.VarMap)
	}
	coda.execute(ft.Cmds)
}

func makeAbs(filename string) string {
	filename, err := filepath.Abs(filename)
	if err != nil {
		fmt.Printf("Can not assemble absolute filename: %s\n", err)
		os.Exit(1)
	}
	return filename
}
