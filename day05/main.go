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
	Machine struct {
		pc  int
		mem []int
		inp int
		out int
	}
)

func NewMachine(mem []int, inp int) *Machine {
	return &Machine{
		pc:  0,
		mem: mem,
		inp: inp,
		out: 0,
	}
}

const (
	modePos = iota
	modeImm
)

func paramMode(mode int) int {
	switch mode {
	case 0:
		return modePos
	case 1:
		return modeImm
	default:
		log.Fatal("Illegal param mode", mode)
		return 0
	}
}

func decodeOp(code int) (int, int, int, int) {
	op := code % 100
	mode1 := paramMode(code / 100 % 10)
	mode2 := paramMode(code / 1000 % 10)
	mode3 := paramMode(code / 10000 % 10)
	return op, mode1, mode2, mode3
}

func (m *Machine) evalArg(mode int, arg int) int {
	switch mode {
	case modePos:
		return m.mem[arg]
	case modeImm:
		return arg
	default:
		log.Fatal("Illegal arg mode")
		return 0
	}
}

func (m *Machine) Exec() (bool, bool) {
	out := false
	op, a1, a2, _ := decodeOp(m.mem[m.pc])
	switch op {
	case 1:
		arg1 := m.mem[m.pc+1]
		arg2 := m.mem[m.pc+2]
		dest := m.mem[m.pc+3]
		m.mem[dest] = m.evalArg(a1, arg1) + m.evalArg(a2, arg2)
		m.pc += 4
	case 2:
		arg1 := m.mem[m.pc+1]
		arg2 := m.mem[m.pc+2]
		dest := m.mem[m.pc+3]
		m.mem[dest] = m.evalArg(a1, arg1) * m.evalArg(a2, arg2)
		m.pc += 4
	case 3:
		arg1 := m.mem[m.pc+1]
		m.mem[arg1] = m.inp
		m.pc += 2
	case 4:
		arg1 := m.mem[m.pc+1]
		m.out = m.evalArg(a1, arg1)
		m.pc += 2
		out = true
	case 5:
		arg1 := m.mem[m.pc+1]
		arg2 := m.mem[m.pc+2]
		if m.evalArg(a1, arg1) != 0 {
			m.pc = m.evalArg(a2, arg2)
		} else {
			m.pc += 3
		}
	case 6:
		arg1 := m.mem[m.pc+1]
		arg2 := m.mem[m.pc+2]
		if m.evalArg(a1, arg1) == 0 {
			m.pc = m.evalArg(a2, arg2)
		} else {
			m.pc += 3
		}
	case 7:
		arg1 := m.mem[m.pc+1]
		arg2 := m.mem[m.pc+2]
		arg3 := m.mem[m.pc+3]
		if m.evalArg(a1, arg1) < m.evalArg(a2, arg2) {
			m.mem[arg3] = 1
		} else {
			m.mem[arg3] = 0
		}
		m.pc += 4
	case 8:
		arg1 := m.mem[m.pc+1]
		arg2 := m.mem[m.pc+2]
		arg3 := m.mem[m.pc+3]
		if m.evalArg(a1, arg1) == m.evalArg(a2, arg2) {
			m.mem[arg3] = 1
		} else {
			m.mem[arg3] = 0
		}
		m.pc += 4
	case 99:
		m.pc += 1
		return false, false
	default:
		log.Fatal("Illegal op code", m.pc, m.mem[m.pc])
	}
	return out, true
}

func (m *Machine) Execute() {
	for {
		out, ok := m.Exec()
		if !ok {
			break
		}
		if out {
			fmt.Println(m.out)
		}
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
		m := NewMachine(mem, 1)
		m.Execute()
	}
	{
		mem := make([]int, len(tokens))
		copy(mem, tokens)
		m := NewMachine(mem, 5)
		m.Execute()
	}
}
