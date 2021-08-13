package main

import (
	"gmm/util"
)

func main() {
	parseArgs()
	conf := readConfig(CLI.Config)

	util.Pprint(conf)
}
