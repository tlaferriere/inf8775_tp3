package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func readProb(path string) (uint, []uint, [][]int, [][2]uint) {
	// Read problem description
	// Input:
	//   path: path to file containing problem
	// Output:
	//   t: number of nodes
	//   nAtoms: number of each atom. shape = k
	//   h: energy of each atom combination. shape = k x k
	//   edges: available edges of target graph. shape = A x 2
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)
	var text [4][]string
	part := 0
	for scanner.Scan() {
		ln := scanner.Text()
		if ln != "" {
			text[part] = append(text[part], ln)
		} else {
			part++
		}
	}
	_ = f.Close()

	var metadata [3]int
	for i, s := range strings.Split(text[0][0], " ") {
		num, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		metadata[i] = num
	}
	nAtoms := make([]uint, metadata[1])
	nAtomsLn := strings.Split(text[1][0], " ")
	if len(nAtomsLn) != metadata[1] {
		log.Fatalf("Wrong number of atom types: %v types found, expected %v types", len(nAtomsLn), metadata[1])
	}
	h := make([][]int, metadata[1])
	for i, s := range nAtomsLn {
		n, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		nAtoms[i] = uint(n)
		h[i] = make([]int, metadata[1])
		hLn := strings.Split(text[2][i], " ")
		if len(hLn) != metadata[1] {
			log.Fatalf("Wrong number of atom energies: %v energies found, expected %v energies", len(hLn), metadata[1])
		}
		for j, s := range hLn {
			e, err := strconv.Atoi(s)
			if err != nil {
				panic(err)
			}
			h[i][j] = e
		}
	}
	edges := make([][2]uint, metadata[2])
	if len(text[3]) != metadata[2] {
		log.Fatalf("Wrong number of edges: %v edges found, expected %v edges", len(text[3]), metadata[2])
	}
	for i, ln := range text[3] {
		s := strings.Split(ln, " ")
		a, err := strconv.Atoi(s[0])
		if err != nil {
			panic(err)
		}
		b, err := strconv.Atoi(s[1])
		if err != nil {
			panic(err)
		}
		edges[i] = [2]uint{uint(a), uint(b)}
	}
	return uint(metadata[0]), nAtoms, h, edges
}
