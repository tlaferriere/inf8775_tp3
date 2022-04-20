package main

import (
	"flag"
	"fmt"
	"runtime"
)

func main() {

	// Enable parallelism
	runtime.GOMAXPROCS(runtime.NumCPU())

	var (
		// startAlgoName   = flag.String("a", "start", "Algorithme de départ")
		// improveAlgoName = flag.String("b", "heur_glouton", "Algorithme d'amélioration")
		exemplairesPath = flag.String("e", "N100_K3_0", "Path vers exemplaire")
		printSol        = flag.Bool("p", false, "Afficher la meilleure solution trouvée")
	)
	flag.Parse()
	t, nAtoms, h, edges := readProb(*exemplairesPath)

	sol, graphMap := start(t, nAtoms, h, edges)
	if *printSol {
		for _, a := range sol.nodes {
			fmt.Printf("%v ", a)
		}
		fmt.Printf("\n")
	} else {
		println(sol.energy)
	}

	improvedSol := make(chan solution)
	go improve(sol, t, h, edges, graphMap, improvedSol)
	if *printSol {
		for sol := range improvedSol {
			for _, a := range sol.nodes {
				fmt.Printf("%v ", a)
			}
			fmt.Printf("\n")
		}
	} else {
		for sol := range improvedSol {
			println(sol.energy)
		}
	}
}
