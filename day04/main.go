package main

import (
	"fmt"
)

const (
	puzzleInputMin = 231832
	puzzleInputMax = 767346
)

func passToDigits(pass int) [6]int {
	return [6]int{
		pass / 100000 % 10,
		pass / 10000 % 10,
		pass / 1000 % 10,
		pass / 100 % 10,
		pass / 10 % 10,
		pass % 10,
	}
}

func isValidPass(pass int) bool {
	s := passToDigits(pass)

	hasDouble := false

	k := s[0]
	for _, i := range s[1:] {
		if i == k {
			hasDouble = true
		}
		if i < k {
			return false
		}
		k = i
	}

	return hasDouble
}

func isValidPass2(pass int) bool {
	if !isValidPass(pass) {
		return false
	}

	s := passToDigits(pass)

	hasRun2 := false

	run := 1
	k := s[0]
	for _, i := range s[1:] {
		if i == k {
			run++
		} else {
			if run == 2 {
				hasRun2 = true
			}
			run = 1
		}
		k = i
	}

	return hasRun2 || run == 2
}

func main() {
	count := 0
	count2 := 0
	for i := puzzleInputMin; i <= puzzleInputMax; i++ {
		if isValidPass(i) {
			count++
		}
		if isValidPass2(i) {
			count2++
		}
	}
	fmt.Println(count)
	fmt.Println(count2)
}
