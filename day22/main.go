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
	Deck struct {
		cards []int
		alt   []int
		size  int
	}
)

func NewDeck(num int) *Deck {
	cards := make([]int, 0, num)
	for i := 0; i < num; i++ {
		cards = append(cards, i)
	}
	return &Deck{
		cards: cards,
		alt:   make([]int, num),
		size:  num,
	}
}

func (d *Deck) DealNew() {
	l := d.size
	for i := 0; i < l/2; i++ {
		d.cards[i], d.cards[l-i-1] = d.cards[l-i-1], d.cards[i]
	}
}

func (d *Deck) Cut(n int) {
	n = (n + d.size) % d.size
	copy(d.alt, d.cards[n:])
	copy(d.alt[d.size-n:], d.cards[:n])
	d.cards, d.alt = d.alt, d.cards
}

func (d *Deck) DealIncr(incr int) {
	k := 0
	l := d.size
	for _, i := range d.cards {
		d.alt[k] = i
		k = (k + incr) % l
	}
	d.cards, d.alt = d.alt, d.cards
}

func (d *Deck) FindCard(c int) int {
	for n, i := range d.cards {
		if i == c {
			return n
		}
	}
	return -1
}

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

	deck := NewDeck(10007)

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

	fmt.Println(deck.FindCard(2019))
}
