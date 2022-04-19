package main

import (
	"flag"
	"fmt"
)

func main() {

	var (
		// startAlgoName   = flag.String("a", "start", "Algorithme de départ")
		// improveAlgoName = flag.String("b", "heur_glouton", "Algorithme d'amélioration")
		exemplairesPath = flag.String("e", "N100_K3_0", "Path vers exemplaire")
		printSol        = flag.Bool("p", false, "Afficher la meilleure solution trouvée")
	)
	flag.Parse()
	t, nAtoms, h, edges := readProb(*exemplairesPath)
	e, sol, graphMap := start(t, nAtoms, h, edges)
	if *printSol {
		for _, a := range sol {
			fmt.Printf("%v ", a)
		}
		fmt.Printf("\n")
	} else {
		print(e)
	}

	return

	improvedEnergy := make(chan int)
	improvedSolution := make(chan []uint)
	go improve(e, sol, t, nAtoms, h, graphMap, *printSol, improvedEnergy, improvedSolution)
	if *printSol {
		for sol := range improvedSolution {
			for _, a := range sol {
				fmt.Printf("%v ", a)
			}
			fmt.Printf("\n")
		}
	} else {
		for e := range improvedEnergy {
			print(e)
		}
	}
}
