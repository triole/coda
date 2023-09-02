package main

import (
	"fmt"
)

func main() {
	parseArgs()

	fileConfig := makeAbs(CLI.Config)
	fileToProcess := makeAbs(CLI.Filename)

	coda := initCoda(fileConfig, fileToProcess)
	ft := coda.detect()

	if CLI.DryRun {
		fmt.Printf("\n\n%s\n\n", "Used VarMap")
		pprint(coda.VarMap)
	}
	coda.execute(ft.Cmds)
}
