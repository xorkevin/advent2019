package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

const (
	puzzleInput = "input.txt"
)

func main() {
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

		k := 0
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			num, err := strconv.Atoi(scanner.Text())
			if err != nil {
				log.Fatal(err)
			}
			k += num/3 - 2
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		fmt.Println(k)
	}
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

		k := 0
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			num, err := strconv.Atoi(scanner.Text())
			if err != nil {
				log.Fatal(err)
			}
			j := num/3 - 2
			for j > 0 {
				k += j
				j = j/3 - 2
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		fmt.Println(k)
	}
}
