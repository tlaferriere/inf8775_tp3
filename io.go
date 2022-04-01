package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func ReadProb(path string) [] {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)
	var text []string
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}
	_ = f.Close()

	var blocks []block
	for _, ln := range text {
		var blockDims [3]int
		for i, s := range strings.Split(ln, " ") {
			num, err := strconv.Atoi(s)
			if err != nil {
				panic(err)
			}
			blockDims[i] = num
		}
		blocks = append(blocks, block{blockDims[0], blockDims[1], blockDims[2]})
	}
	return blocks
}

var (
	algoName        = flag.String("a", "glouton", "Algorithme")
	exemplairesPath = flag.String("e", "b100_1.txt", "Path vers exemplaire")
	printTower      = flag.Bool("p", false, "Afficher les blocs utilisés")
	printTime       = flag.Bool("t", false, "Afficher le temps pris en ms")
	graphTime       = flag.Bool("g", false, "Générer les graphiques de temps avec le serveur de graphiques")
	graphServerAddr = flag.String("addr", "localhost:50051", "Addresse du serveur de gaphiques")
)

func cliMode(exemplairesPath *string, algoName *string, printTower *bool, printTime *bool) {
	blocks := readBlocks(*exemplairesPath)
	var tower []block
	start := time.Now()
	if *algoName == "glouton" {
		_, tower, _ = glouton(blocks)
	} else if *algoName == "progdyn" {
		_, tower = progdyn(blocks)
	} else if *algoName == "tabou" {
		_, tower = tabou(blocks)
	}
	elapsed := time.Since(start)
	if *printTower {
		for _, b := range tower {
			fmt.Printf("%v %v %v\n", b.height, b.length, b.width)
		}
	}
	if *printTime {
		fmt.Println(elapsed.Seconds() * 1000)
	}
}
