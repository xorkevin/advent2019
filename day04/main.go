package main

import (
	"fmt"
	"strconv"
)

const (
	puzzleInputMin = 231832
	puzzleInputMax = 767346
)

func isValidPass(pass int) bool {
	s := strconv.Itoa(pass)

	hasDouble := false

	k := byte(s[0])
	for _, i := range s[1:] {
		if byte(i) == k {
			hasDouble = true
		}
		if byte(i) < k {
			return false
		}
		k = byte(i)
	}

	return hasDouble
}

func isValidPass2(pass int) bool {
	if !isValidPass(pass) {
		return false
	}

	s := strconv.Itoa(pass)

	hasRun2 := false

	run := 1
	k := byte(s[0])
	for _, i := range s[1:] {
		if byte(i) == k {
			run++
		} else {
			if run == 2 {
				hasRun2 = true
			}
			run = 1
		}
		k = byte(i)
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
