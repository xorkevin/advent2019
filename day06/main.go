package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	puzzleInput = "input.txt"
)

type (
	OrbitsMap struct {
		orbits map[string]map[string]struct{}
	}
)

func (m *OrbitsMap) Add(a, b string) {
	_, ok := m.orbits[a]
	if !ok {
		m.orbits[a] = map[string]struct{}{}
	}
	m.orbits[a][b] = struct{}{}
}

func (m *OrbitsMap) Traverse(a string, depth int) int {
	c := depth
	for k := range m.orbits[a] {
		c += m.Traverse(k, depth+1)
	}
	return c
}

func (m *OrbitsMap) GetPath(target, a string) []string {
	if _, ok := m.orbits[a]; !ok {
		return nil
	}

	if _, ok := m.orbits[a][target]; ok {
		return []string{target}
	}

	for i := range m.orbits[a] {
		k := m.GetPath(target, i)
		if k != nil {
			return append(k, a)
		}
	}
	return nil
}

func main() {
	file, err := os.Open(puzzleInput)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	orbits := OrbitsMap{
		orbits: map[string]map[string]struct{}{},
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		k := strings.Split(scanner.Text(), ")")
		orbits.Add(k[0], k[1])
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(orbits.Traverse("COM", 0))
	x := orbits.GetPath("YOU", "COM")
	y := orbits.GetPath("SAN", "COM")

	for i := 0; i < len(x)/2; i++ {
		x[i], x[len(x)-i-1] = x[len(x)-i-1], x[i]
	}
	for i := 0; i < len(y)/2; i++ {
		y[i], y[len(y)-i-1] = y[len(y)-i-1], y[i]
	}

	for i := 0; i < len(x) && i < len(y); i++ {
		if x[i] != y[i] {
			fmt.Println(len(x) - i + len(y) - i - 2)
			break
		}
	}
}
