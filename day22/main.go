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
	Card struct {
		a, b int
		size int
	}
)

func NewCard(size int) *Card {
	return &Card{
		a:    1,
		b:    0,
		size: size,
	}
}

func (c *Card) Reverse() {
	c.a = (c.a*-1 + c.size) % c.size
	c.b = (c.b*-1 - 1 + c.size) % c.size
}

func (c *Card) Cut(n int) {
	c.b = (c.b - n + c.size) % c.size
}

func (c *Card) DealIncr(n int) {
	c.a = (c.a * n) % c.size
	c.b = (c.b * n) % c.size
}

func (c *Card) FindCard(n int) int {
	return (n*c.a + c.b) % c.size
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

	deck := NewCard(10007)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if line[0] == "deal" {
			if line[1] == "into" {
				deck.Reverse()
			} else if line[1] == "with" {
				num, err := strconv.Atoi(line[3])
				if err != nil {
					log.Fatalln(err)
				}
				deck.DealIncr(num)
			} else {
				log.Fatalln("invalid input")
			}
		} else if line[0] == "cut" {
			num, err := strconv.Atoi(line[1])
			if err != nil {
				log.Fatalln(err)
			}
			deck.Cut(num)
		} else {
			log.Fatalln("invalid input")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(deck.FindCard(2019))
}
