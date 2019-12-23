package main

import (
	"bufio"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
)

const (
	puzzleInput = "input.txt"
)

type (
	Card struct {
		a, b int64
		size int64
	}
)

func NewCard(size int64) *Card {
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

func (c *Card) Cut(n int64) {
	c.b = (c.b - n + c.size) % c.size
}

func (c *Card) DealIncr(n int64) {
	c.a = (c.a * n) % c.size
	c.b = (c.b * n) % c.size
}

func (c *Card) FindCard(n int64) int64 {
	return (n*c.a + c.b) % c.size
}

// y = a x + b
// x = y*a' - b*a'
func (c *Card) Inverse() *Card {
	a := new(big.Int).ModInverse(big.NewInt(c.a), big.NewInt(c.size)).Int64()
	b := c.b * a % c.size
	return &Card{
		a:    a,
		b:    b,
		size: c.size,
	}
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
				deck.DealIncr(int64(num))
			} else {
				log.Fatalln("invalid input")
			}
		} else if line[0] == "cut" {
			num, err := strconv.Atoi(line[1])
			if err != nil {
				log.Fatalln(err)
			}
			deck.Cut(int64(num))
		} else {
			log.Fatalln("invalid input")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(deck.FindCard(2019))
}
