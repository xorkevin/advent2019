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
	Layer [][]int
)

func NewLayer() Layer {
	layer := make(Layer, 0, 6)
	for i := 0; i < 6; i++ {
		layer = append(layer, make([]int, 25))
	}
	return layer
}

func (l *Layer) CountNum(num int) int {
	count := 0
	for _, i := range *l {
		for _, j := range i {
			if j == num {
				count++
			}
		}
	}
	return count
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

	img := []Layer{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "")
		nums := make([]int, 0, len(line))
		for _, i := range line {
			num, err := strconv.Atoi(i)
			if err != nil {
				log.Fatal(err)
			}
			nums = append(nums, num)
		}
		for len(nums) > 0 {
			layer := NewLayer()
			for n, i := range nums[:25*6] {
				layer[n/25][n%25] = i
			}
			img = append(img, layer)
			nums = nums[25*6:]
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	layer := 0
	minZeros := 25
	for n, i := range img {
		count := i.CountNum(0)
		if count < minZeros {
			layer = n
			minZeros = count
		}
	}

	count1 := img[layer].CountNum(1)
	count2 := img[layer].CountNum(2)

	fmt.Println(count1 * count2)

	finalLayer := NewLayer()
	for i := 0; i < 6; i++ {
		for j := 0; j < 25; j++ {
			for k := 0; k < len(img); k++ {
				if l := img[k][i][j]; l != 2 {
					finalLayer[i][j] = l
					break
				}
			}
		}
	}

	for _, i := range finalLayer {
		for _, j := range i {
			if j == 1 {
				fmt.Print("#")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}
