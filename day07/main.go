package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	puzzleInput = "input.txt"
)

type (
	Machine struct {
		pc       int
		mem      []int
		inp      chan int
		out      chan int
		outGauge int
	}
)

func NewMachine(mem []int) *Machine {
	return &Machine{
		pc:       0,
		mem:      mem,
		inp:      make(chan int, 2),
		out:      make(chan int, 2),
		outGauge: 0,
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

func (m *Machine) GetInp() int {
	return <-m.inp
}

func (m *Machine) Write(inp int) {
	m.inp <- inp
}

func (m *Machine) Output(out int) {
	m.outGauge = out
	m.out <- out
}

func (m *Machine) Read() int {
	return <-m.out
}

func (m *Machine) Exec() bool {
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
		m.mem[arg1] = m.GetInp()
		m.pc += 2
	case 4:
		arg1 := m.mem[m.pc+1]
		m.Output(m.evalArg(a1, arg1))
		m.pc += 2
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

func Perm(a []int, f func([]int)) {
	perm(a, f, 0)
}

// Permute the values at index i to len(a)-1.
func perm(a []int, f func([]int), i int) {
	if i > len(a) {
		f(a)
		return
	}
	perm(a, f, i+1)
	for j := i + 1; j < len(a); j++ {
		a[i], a[j] = a[j], a[i]
		perm(a, f, i+1)
		a[i], a[j] = a[j], a[i]
	}
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
			nums := strings.Split(scanner.Text(), ",")
			for _, i := range nums {
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
		maxOut := 0
		Perm([]int{0, 1, 2, 3, 4}, func(phases []int) {
			out := 0
			for _, phase := range phases {
				mem := make([]int, len(tokens))
				copy(mem, tokens)
				m := NewMachine(mem)
				m.Write(phase)
				m.Write(out)
				m.Execute()
				out = m.Read()
			}
			if out > maxOut {
				maxOut = out
			}
		})
		fmt.Println(maxOut)
	}

	{
		maxOut := 0
		Perm([]int{5, 6, 7, 8, 9}, func(phases []int) {
			m := make([]*Machine, 0, len(phases))
			{
				var prev *Machine
				for _, phase := range phases {
					mem := make([]int, len(tokens))
					copy(mem, tokens)
					k := NewMachine(mem)
					if prev != nil {
						k.inp = prev.out
						k.Write(phase)
					}
					prev = k
					m = append(m, k)
				}
				m[0].inp = prev.out
				m[0].Write(phases[0])
			}

			wg := sync.WaitGroup{}
			for _, i := range m {
				k := i
				wg.Add(1)
				go func() {
					defer wg.Done()
					k.Execute()
				}()
			}

			m[0].Write(0)
			wg.Wait()
			if k := m[len(m)-1].outGauge; k > maxOut {
				maxOut = k
			}
		})
		fmt.Println(maxOut)
	}
}
