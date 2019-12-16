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
	puzzleInput = "input2.txt"
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
	return sum
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
		next = append(next, Abs(v.PhaseInner(pattern, 1, n+1)%10))
	}
	return next
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
		nums2 := orignums
		pattern := []int{0, 1, 0, -1}
		for i := 0; i < 100; i++ {
			nums2 = nums2.Phase(pattern)
		}
		fmt.Println(SliceToNum(nums2[0:8]))
	}
	//{
	//	nums2 := make(Vector, 0, len(orignums)*10000)
	//	for i := 0; i < 10000; i++ {
	//		nums2 = append(nums2, orignums...)
	//	}
	//	pattern := []int{0, 1, 0, -1}
	//	offset := SliceToNum(nums2[0:7])
	//	fmt.Println(offset)
	//	for i := 0; i < 100; i++ {
	//		nums2 = nums2.Phase(pattern)
	//		fmt.Println("done", i+1)
	//	}

	//	fmt.Println(SliceToNum(nums2[offset : offset+8]))
	//}
}
