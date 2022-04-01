package main

import (
	"flag"
)

func main() {
	flag.Parse()

	if *graphTime {
		graphMode(*graphServerAddr)
	} else {
		cliMode(exemplairesPath, algoName, printTower, printTime)
	}
}
