package main

import (
	"fmt"

	"github.com/jedib0t/go-pretty/table"
)

func main() {
	parseArgs()

	fileConfig := makeAbs(CLI.Config)
	fileToProcess := makeAbs(CLI.Filename)

	coda := initCoda(fileConfig, fileToProcess)
	ft := coda.detect()

	if CLI.DryRun {
		fmt.Printf("\nvariables and their mapped values\n\n")
		t := newTable()
		t.AppendHeader(table.Row{
			"variable", "value", "description",
		})
		for _, val := range orderedIterator(coda.VarMap) {
			t.AppendRow(
				[]interface{}{
					val,
					coda.VarMap[val].Variable,
					coda.VarMap[val].Desc,
				},
			)
		}
		t.Render()
	}
	coda.execute(ft.Cmds)
}
