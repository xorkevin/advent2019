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
		card int
		size int
	}
)

func NewCard(card, size int) *Card {
	return &Card{
		card: card,
		size: size,
	}
}

func (c *Card) DealNew() {
	c.card = c.size - c.card - 1
}

func (c *Card) Cut(n int) {
	n = (n + c.size) % c.size
	c.card = (c.card - n + c.size) % c.size
}

func (c *Card) DealIncr(n int) {
	c.card = (c.card * n) % c.size
}

func (c *Card) FindCard() int {
	return c.card
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

	deck := NewCard(2019, 10007)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if line[0] == "deal" {
			if line[1] == "into" {
				deck.DealNew()
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

	fmt.Println(deck.FindCard())
}
