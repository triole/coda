package main

import "fmt"

func main() {
	parseArgs()
	coda := readConfig(CLI.Config)
	_, cmds := coda.detect(CLI.Filename)

	if CLI.Debug == true {
		fmt.Printf("\n%s\n", "Var map")
		pprint(coda.VarMap)
		fmt.Printf("\n%s\n", "Commands")
	}

	coda.execute(cmds)
}
