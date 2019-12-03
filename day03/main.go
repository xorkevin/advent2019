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
	Tuple struct {
		x, y int
	}
)

func Abs(a, b int) int {
	k := a - b
	if k < 0 {
		return -k
	}
	return k
}

func Dist(t Tuple, t2 Tuple) int {
	return Abs(t.x, t2.x) + Abs(t.y, t2.y)
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

	wire1 := map[Tuple]int{}

	dist := -1
	dist2 := -1

	first := true
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		x := 0
		y := 0
		n := 0
		for _, step := range strings.Split(line, ",") {
			num, err := strconv.Atoi(step[1:])
			if err != nil {
				log.Fatal(err)
			}
			for i := 0; i < num; i++ {
				switch step[0] {
				case 'U':
					y -= 1
				case 'D':
					y += 1
				case 'R':
					x += 1
				case 'L':
					x -= 1
				}
				n++
				if first {
					wire1[Tuple{x, y}] = n
				} else {
					if w1, ok := wire1[Tuple{x, y}]; ok {
						if k := Dist(Tuple{x, y}, Tuple{0, 0}); dist < 0 || k < dist {
							dist = k
						}
						if dist2 < 0 || w1+n < dist2 {
							dist2 = w1 + n
						}
					}
				}
			}
		}
		if first {
			first = false
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(dist)
	fmt.Println(dist2)
}
