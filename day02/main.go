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
	puzzleInput  = "input.txt"
	puzzleInput2 = 19690720
)

type (
	Machine struct {
		pc  int
		mem []int
	}
)

func NewMachine(mem []int) *Machine {
	return &Machine{
		pc:  0,
		mem: mem,
	}
}

func (m *Machine) Exec() bool {
	switch m.mem[m.pc] {
	case 1:
		arg1 := m.mem[m.pc+1]
		arg2 := m.mem[m.pc+2]
		dest := m.mem[m.pc+3]
		m.mem[dest] = m.mem[arg1] + m.mem[arg2]
		m.pc += 4
	case 2:
		arg1 := m.mem[m.pc+1]
		arg2 := m.mem[m.pc+2]
		dest := m.mem[m.pc+3]
		m.mem[dest] = m.mem[arg1] * m.mem[arg2]
		m.pc += 4
	case 99:
		m.pc += 1
		return false
	default:
		log.Fatal("Illegal op code", m.pc, m.mem[m.pc])
	}
	return true
}

func (m *Machine) Execute() {
	for m.Exec() {
	}
}

func (m *Machine) MemAt(offset int) int {
	return m.mem[offset]
}

func (m *Machine) MemSet(offset, val int) {
	m.mem[offset] = val
}

func main() {
	tokens := []int{}
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
			line := scanner.Text()
			for _, i := range strings.Split(line, ",") {
				num, err := strconv.Atoi(i)
				if err != nil {
					log.Fatal(err)
				}
				tokens = append(tokens, num)
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
	{
		mem := make([]int, len(tokens))
		copy(mem, tokens)
		m := NewMachine(mem)
		m.MemSet(1, 12)
		m.MemSet(2, 2)
		m.Execute()
		fmt.Println(m.MemAt(0))
	}
	{
	outer:
		for i := 0; i < 100; i++ {
			for j := 0; j < 100; j++ {
				mem := make([]int, len(tokens))
				copy(mem, tokens)
				m := NewMachine(mem)
				m.MemSet(1, i)
				m.MemSet(2, j)
				m.Execute()
				if m.MemAt(0) == puzzleInput2 {
					fmt.Println(i*100 + j)
					break outer
				}
			}
		}
	}
}
