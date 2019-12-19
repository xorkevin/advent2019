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
	ramSize     = 8192
)

type (
	Machine struct {
		pc       int
		mem      []int
		inp      chan int
		out      chan int
		outGauge int
		relBase  int
	}
)

func NewMachine(mem []int) *Machine {
	return &Machine{
		pc:       0,
		mem:      mem,
		inp:      make(chan int, 2),
		out:      make(chan int, 2),
		outGauge: 0,
		relBase:  0,
	}
}

const (
	modePos = iota
	modeImm
	modeRel
)

func paramMode(mode int) int {
	switch mode {
	case 0:
		return modePos
	case 1:
		return modeImm
	case 2:
		return modeRel
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

func (m *Machine) getMem(pos int) int {
	return m.mem[pos]
}

func (m *Machine) evalArg(mode int, arg int) int {
	switch mode {
	case modePos:
		return m.getMem(arg)
	case modeImm:
		return arg
	case modeRel:
		return m.getMem(arg + m.relBase)
	default:
		log.Fatal("Illegal arg mode")
		return 0
	}
}

func (m *Machine) getArg(mode, offset int) int {
	arg := m.getMem(m.pc + offset)
	return m.evalArg(mode, arg)
}

func (m *Machine) setMem(pos, val int) {
	m.mem[pos] = val
}

func (m *Machine) setArg(mode, offset, val int) {
	arg := m.getMem(m.pc + offset)
	switch mode {
	case modePos:
		m.setMem(arg, val)
	case modeImm:
		log.Fatal("Illegal mem write imm mode")
	case modeRel:
		m.setMem(arg+m.relBase, val)
	default:
		log.Fatal("Illegal mem write mode")
	}
}

func (m *Machine) stepPC(offset int) {
	m.pc += offset
}

func (m *Machine) Write(inp int) {
	m.inp <- inp
}

func (m *Machine) recvInput() int {
	return <-m.inp
}

func (m *Machine) sendOutput(out int) {
	m.outGauge = out
	m.out <- out
}

func (m *Machine) Read() (int, bool) {
	v, ok := <-m.out
	if !ok {
		return 0, false
	}
	return v, true
}

func (m *Machine) Exec() bool {
	op, a1, a2, a3 := decodeOp(m.getMem(m.pc))
	switch op {
	case 1:
		m.setArg(a3, 3, m.getArg(a1, 1)+m.getArg(a2, 2))
		m.stepPC(4)
	case 2:
		m.setArg(a3, 3, m.getArg(a1, 1)*m.getArg(a2, 2))
		m.stepPC(4)
	case 3:
		m.setArg(a1, 1, m.recvInput())
		m.stepPC(2)
	case 4:
		m.sendOutput(m.getArg(a1, 1))
		m.stepPC(2)
	case 5:
		if m.getArg(a1, 1) != 0 {
			m.pc = m.getArg(a2, 2)
		} else {
			m.stepPC(3)
		}
	case 6:
		if m.getArg(a1, 1) == 0 {
			m.pc = m.getArg(a2, 2)
		} else {
			m.stepPC(3)
		}
	case 7:
		if m.getArg(a1, 1) < m.getArg(a2, 2) {
			m.setArg(a3, 3, 1)
		} else {
			m.setArg(a3, 3, 0)
		}
		m.stepPC(4)
	case 8:
		if m.getArg(a1, 1) == m.getArg(a2, 2) {
			m.setArg(a3, 3, 1)
		} else {
			m.setArg(a3, 3, 0)
		}
		m.stepPC(4)
	case 9:
		m.relBase += m.getArg(a1, 1)
		m.stepPC(2)
	case 99:
		m.stepPC(1)
		return false
	default:
		log.Fatal("Illegal op code", m.pc, m.getMem(m.pc))
	}
	return true
}

func (m *Machine) Execute() {
	for m.Exec() {
	}
	close(m.out)
}

type (
	Point struct {
		x, y int
	}
)

func FindSquare(pos Point, size int, grid map[Point]struct{}) bool {
	if _, ok := grid[Point{pos.x, pos.y}]; !ok {
		return false
	}
	if _, ok := grid[Point{pos.x + size - 1, pos.y}]; !ok {
		return false
	}
	if _, ok := grid[Point{pos.x, pos.y + size - 1}]; !ok {
		return false
	}
	if _, ok := grid[Point{pos.x + size - 1, pos.y + size - 1}]; !ok {
		return false
	}
	return true
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
		count := 0
		for y := 0; y < 50; y++ {
			for x := 0; x < 50; x++ {
				mem := make([]int, ramSize)
				copy(mem, tokens)
				m := NewMachine(mem)
				go m.Execute()
				m.Write(x)
				m.Write(y)
				out, ok := m.Read()
				if !ok {
					log.Fatalln("Failed to read")
				}
				if out == 1 {
					count++
				}
			}
		}
		fmt.Println(count)
	}
	{
		ystart := 900
		yend := 1100
		xstart := 300
		xend := 500
		points := map[Point]struct{}{}
		for y := ystart; y < yend; y++ {
			for x := xstart; x < xend; x++ {
				mem := make([]int, ramSize)
				copy(mem, tokens)
				m := NewMachine(mem)
				go m.Execute()
				m.Write(x)
				m.Write(y)
				out, ok := m.Read()
				if !ok {
					log.Fatalln("Failed to read")
				}
				if out == 1 {
					points[Point{x, y}] = struct{}{}
					fmt.Print("#")
				} else {
					fmt.Print(".")
				}
			}
			fmt.Println()
		}
		for y := ystart; y < yend; y++ {
			for x := xstart; x < xend; x++ {
				if FindSquare(Point{x, y}, 100, points) {
					fmt.Println(x*10000 + y)
					return
				}
			}
		}
	}
}
