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

func Abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

type (
	Vector []int
)

func (v Vector) Ith(i int, repeat int) int {
	return v[(i/repeat)%len(v)]
}

func (v Vector) PhaseInner(pattern Vector, offset, repeat int) int {
	sum := 0
	for n, i := range v {
		sum += i * pattern.Ith(n+offset, repeat)
	}
	return Abs(sum) % 10
}

func SliceToNum(nums []int) int {
	if len(nums) == 0 {
		return -1
	}
	sum := 0
	for _, i := range nums {
		sum *= 10
		sum += i
	}
	return sum
}

func (v Vector) Phase(pattern Vector) Vector {
	next := make(Vector, 0, len(v))
	for n := range v {
		next = append(next, v.PhaseInner(pattern, 1, n+1))
	}
	return next
}

func (v Vector) Phase2(pattern Vector, offset int) Vector {
	l := len(v)
	partial := 0
	for i := l - 1; i >= offset; i-- {
		partial = (partial + v[i]) % 10
		v[i] = partial
	}
	return v
}

func main() {
	orignums := Vector{}
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
			numbers := strings.Split(scanner.Text(), "")
			for _, i := range numbers {
				num, err := strconv.Atoi(i)
				if err != nil {
					log.Fatal(err)
				}
				orignums = append(orignums, num)
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	{
		nums := orignums
		pattern := []int{0, 1, 0, -1}
		for i := 0; i < 100; i++ {
			nums = nums.Phase(pattern)
		}
		fmt.Println(SliceToNum(nums[0:8]))
	}
	{
		nums := make(Vector, 0, len(orignums)*10000)
		for i := 0; i < 10000; i++ {
			nums = append(nums, orignums...)
		}
		pattern := []int{0, 1, 0, -1}
		offset := SliceToNum(nums[0:7])
		for i := 0; i < 100; i++ {
			nums = nums.Phase2(pattern, offset)
		}
		fmt.Println(SliceToNum(nums[offset : offset+8]))
	}
}
