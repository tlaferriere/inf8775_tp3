package main

import (
	"sort"
)

// Départ start
// On remplit le tableau de tous les sommets
// Return: (e uint, nodes []uint)
//  - e: energie totale
//	- nodes: le tableau de tous les sommets
func start(t uint, nAtoms []uint, h [][]int, edges [][2]uint) (int, []uint, [][]uint) {
	possibleEdges := make(chan [][2]uint)
	graphMapChan := make(chan [][]uint)
	go makeGraphMap(edges, t, graphMapChan)
	go sortEdges(nAtoms, h, possibleEdges)

	prioNodes := make(chan []uint)
	mapOut := make(chan [][]uint)
	go sortNodes(t, graphMapChan, mapOut, prioNodes)

	// Placer les arêtes les moins énergisantes en premier dans le graphe sur les noeuds les plus connectés en premier
	occupied := make([]bool, t)
	nodes := make([]uint, t)
	prioIdx := 0 // Pour résumer l'itération dans la liste des noeuds prioritaires
	graphMap := <-mapOut
	prio := <-prioNodes
	for _, e := range <-possibleEdges {
		var limitingAtom, nLimitingAtoms, otherAtom, nOtherAtoms uint
		if nAtoms[e[0]] < nAtoms[e[1]] {
			limitingAtom = e[0]
			nLimitingAtoms = nAtoms[e[0]]
			otherAtom = e[1]
			nOtherAtoms = nAtoms[e[1]]
		} else {
			limitingAtom = e[1]
			nLimitingAtoms = nAtoms[e[1]]
			otherAtom = e[0]
			nOtherAtoms = nAtoms[e[0]]
		}
		sameAtom := limitingAtom == otherAtom
		if nLimitingAtoms > 0 && nOtherAtoms > 0 { // S'assurer que les deux atomes sont encore présents

			for i, n := range prio[prioIdx:] { // Reprendre l'itération dans la liste des noeuds prioritaires
				if !occupied[n] {
					if sameAtom { // Si on place le même atome, on s'assure d'en avoir assez
						if nLimitingAtoms < 2 {
							prioIdx += i // Mise à jour de l'index de priorité pour reprendre l'itération
							break
						}
					}
					// Placer l'atome limitant en premier
					nodes[n] = limitingAtom
					occupied[n] = true
					nLimitingAtoms--
					for _, o := range graphMap[n] {
						// Remplir les connections de l'atome limitant avec l'autre atome
						if !occupied[o] {
							nodes[o] = otherAtom
							occupied[o] = true
							if sameAtom {
								nLimitingAtoms--
								if nLimitingAtoms == 0 {
									break
								}
							} else {
								nOtherAtoms--
								if nOtherAtoms == 0 {
									break
								}
							}
						}
					}
					// Mise à jour du nombre d'atomes restant
					nAtoms[limitingAtom] = nLimitingAtoms
					if !sameAtom {
						nAtoms[otherAtom] = nOtherAtoms
					}
					if nOtherAtoms == 0 || nLimitingAtoms == 0 {
						prioIdx += i + 1 // Mise à jour de l'index de priorité pour reprendre l'itération
						break
					}
				}
			}
		}
	}
	for i, o := range occupied {
		// Remplir les trous restants
		if !o {
			for a := 0; a < len(nAtoms); a++ {
				if nAtoms[a] > 0 {
					nodes[i] = uint(a)
					nAtoms[a]--
					break
				}
			}
		}
	}

	energy := 0
	for _, e := range edges {
		energy += h[nodes[e[0]]][nodes[e[1]]]
	}
	return energy, nodes, graphMap
}

func sortNodes(t uint, mapIn chan [][]uint, mapOut chan [][]uint, nodes chan []uint) {
	prioNodes := make([]uint, t)
	for i := uint(0); i < t; i++ {
		prioNodes[i] = i
	}
	d := <-mapIn
	sort.Slice(prioNodes, func(i, j int) bool {
		return len(d[i]) > len(d[j])
	})
	mapOut <- d
	nodes <- prioNodes
}

func sortEdges(nAtoms []uint, h [][]int, sortedEdges chan [][2]uint) {
	// Création de la liste des arêtes possibles entre les atomes
	var possibleEdges [][2]uint
	for i := 0; i < len(nAtoms); i++ {
		for j := i; j < len(nAtoms); j++ {
			possibleEdges = append(possibleEdges, [2]uint{uint(i), uint(j)})
		}
	}

	// Tri en ordre croissant d'énergie
	// puis en ordre décroissant de nombre d'atomes
	sort.SliceStable(possibleEdges, func(i, j int) bool {
		if h[possibleEdges[i][0]][possibleEdges[i][1]] == h[possibleEdges[j][0]][possibleEdges[j][1]] {
			var minI, minJ uint
			if nAtoms[possibleEdges[i][0]] < nAtoms[possibleEdges[i][1]] {
				minI = possibleEdges[i][0]
			} else {
				minI = possibleEdges[i][1]
			}
			if nAtoms[possibleEdges[j][0]] < nAtoms[possibleEdges[j][1]] {
				minJ = possibleEdges[j][0]
			} else {
				minJ = possibleEdges[j][1]
			}
			return minI > minJ
		}

		return h[possibleEdges[i][0]][possibleEdges[i][1]] < h[possibleEdges[j][0]][possibleEdges[j][1]]
	})

	sortedEdges <- possibleEdges
	close(sortedEdges)
}

type nodeEdges struct {
	id    uint
	edges []bool
}

type nodeDegree struct {
	id     uint
	degree uint
}

// Count all true values in a boolean array
func edgeCounter(edges chan nodeEdges, degrees chan nodeDegree) {
	for n := range edges {
		var nEdges uint
		for _, e := range n.edges {
			if e {
				nEdges++
			}
		}
		degrees <- nodeDegree{n.id, nEdges}
	}
}

// Retourne une liste des noeuds connectés à chaque noeud
func makeGraphMap(edges [][2]uint, t uint, eMap chan [][]uint) {
	mapping := make([][]uint, t)
	for i := uint(0); i < t; i++ {
		mapping[i] = make([]uint, 1)
	}
	for _, e := range edges {
		mapping[e[0]] = append(mapping[e[0]], e[1])
		mapping[e[1]] = append(mapping[e[1]], e[0])
	}
	eMap <- mapping
}

func improve(e int, sol []uint, t uint, nAtoms []uint, h [][]int, graphMap [][]uint, printSol bool, improvedEnergy chan int, improvedSol chan []uint) (uint, []uint) {
	if printSol {
		close(improvedEnergy)
	} else {
		close(improvedSol)
	}
	return 0, make([]uint, t)
}
