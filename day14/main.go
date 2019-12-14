package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	puzzleInput = "input.txt"
)

type (
	Quantity struct {
		chem  string
		count int
	}

	Reaction struct {
		out Quantity
		inp []Quantity
	}

	Reactions struct {
		reactions map[string]Reaction
	}

	Cache map[string]int
)

func (c Cache) Add(chem string, amt int) {
	if v, ok := c[chem]; ok {
		c[chem] = v + amt
	} else {
		c[chem] = amt
	}
}

func (c Cache) Rm(chem string, amt int) int {
	if v, ok := c[chem]; ok {
		if v < amt {
			amt = v
		}
		c[chem] = v - amt
		return amt
	}
	return 0
}

func NewReactions() *Reactions {
	return &Reactions{
		reactions: map[string]Reaction{},
	}
}

func (r *Reactions) Add(re Reaction) {
	r.reactions[re.out.chem] = re
}

func (r *Reactions) Craft(chem string, amt int, cache Cache) int {
	if chem == "ORE" {
		cache.Add("ORE", amt)
		return amt
	}

	re, ok := r.reactions[chem]
	if !ok {
		log.Fatalf("Illegal chem inp: %s\n", chem)
	}

	outMul := amt / re.out.count
	totalOut := outMul * re.out.count
	if totalOut < amt {
		outMul++
		totalOut += re.out.count
	}

	oreCount := 0
	for _, i := range re.inp {
		total := i.count * outMul
		obtained := cache.Rm(i.chem, total)
		if obtained < total {
			oreCount += r.Craft(i.chem, total-obtained, cache)
			cache.Rm(i.chem, total-obtained)
		}
	}

	cache.Add(chem, totalOut)
	return oreCount
}

func parseQuantity(line string) Quantity {
	l := strings.Split(line, " ")
	num, err := strconv.Atoi(l[0])
	if err != nil {
		log.Fatalln(err)
	}
	return Quantity{
		chem:  l[1],
		count: num,
	}
}

func parseReaction(line string) Reaction {
	l := strings.Split(line, " => ")
	inpS := strings.Split(l[0], ", ")
	q := make([]Quantity, 0, len(inpS))
	for _, i := range inpS {
		q = append(q, parseQuantity(i))
	}
	return Reaction{
		out: parseQuantity(l[1]),
		inp: q,
	}
}

const (
	trillion = 1000000000000
)

func main() {
	r := NewReactions()
	{
		file, err := os.Open(puzzleInput)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := file.Close(); err != nil {
				log.Fatal(err)
			}
		}()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			r.Add(parseReaction(scanner.Text()))
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println(r.Craft("FUEL", 1, Cache{}))
	fuel := 1766154
	k := r.Craft("FUEL", fuel, Cache{})
	fmt.Println(k < trillion)
	fmt.Println("remaining:", trillion-k)
	fmt.Println("fuel:", fuel)
}
