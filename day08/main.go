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
	imgWidth    = 25
	imgHeight   = 6
)

type (
	Layer [][]int
	Img   []Layer
)

func NewLayer(w, h int, nums []int) Layer {
	layer := make(Layer, 0, h)
	for i := 0; i < h; i++ {
		layer = append(layer, make([]int, w))
	}
	for n, i := range nums[:w*h] {
		layer[n/w][n%w] = i
	}
	return layer
}

func (l Layer) CountNum(num int) int {
	count := 0
	for _, i := range l {
		for _, j := range i {
			if j == num {
				count++
			}
		}
	}
	return count
}

func (l Layer) Print() {
	for _, i := range l {
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

func NewImg() Img {
	return []Layer{}
}

func (m *Img) AddLayer(l Layer) {
	*m = append(*m, l)
}

func (m Img) Render() Layer {
	if len(m) == 0 {
		return nil
	}
	h := len(m[0])
	if h == 0 {
		return nil
	}
	w := len(m[0][0])

	nums := make([]int, w*h)
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			for k := 0; k < len(m); k++ {
				if l := m[k][i][j]; l != 2 {
					nums[i*w+j] = l
					break
				}
			}
		}
	}

	return NewLayer(w, h, nums)
}

func main() {
	var nums []int
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
			line := strings.Split(scanner.Text(), "")
			nums = make([]int, 0, len(line))
			for _, i := range line {
				num, err := strconv.Atoi(i)
				if err != nil {
					log.Fatal(err)
				}
				nums = append(nums, num)
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	img := NewImg()
	for len(nums) > 0 {
		img.AddLayer(NewLayer(imgWidth, imgHeight, nums))
		nums = nums[imgWidth*imgHeight:]
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

	img.Render().Print()
}
