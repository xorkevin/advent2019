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
		getInp   func() int
	}
)

func NewMachine(mem []int, getInp func() int) *Machine {
	return &Machine{
		pc:       0,
		mem:      mem,
		inp:      make(chan int, 2),
		out:      make(chan int, 2),
		outGauge: 0,
		relBase:  0,
		getInp:   getInp,
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
	m.Write(m.getInp())
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

const (
	tileEmpty  = 0
	tileWall   = 1
	tileBlock  = 2
	tilePaddle = 3
	tileBall   = 4
)

type (
	Point struct {
		x, y int
	}

	Board struct {
		grid map[Point]int
	}
)

func NewBoard() *Board {
	return &Board{
		grid: map[Point]int{},
	}
}

func (b *Board) Emplace(p Point, tile int) {
	b.grid[p] = tile
}

func (b *Board) BlockCount() int {
	count := 0
	for _, v := range b.grid {
		if v == tileBlock {
			count++
		}
	}
	return count
}

func (b *Board) Print() {
	maxX := 0
	maxY := 0
	for k, _ := range b.grid {
		if k.x > maxX {
			maxX = k.x
		}
		if k.y > maxY {
			maxY = k.y
		}
	}

	for i := 0; i <= maxY; i++ {
		for j := 0; j <= maxX; j++ {
			v, ok := b.grid[Point{j, i}]
			if !ok {
				fmt.Fprint(os.Stderr, " ")
				continue
			}
			switch v {
			case tileEmpty:
				fmt.Fprint(os.Stderr, " ")
			case tileWall:
				fmt.Fprint(os.Stderr, "#")
			case tileBlock:
				fmt.Fprint(os.Stderr, "+")
			case tilePaddle:
				fmt.Fprint(os.Stderr, "_")
			case tileBall:
				fmt.Fprint(os.Stderr, "O")
			}
		}
		fmt.Fprintln(os.Stderr)
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
		b := NewBoard()
		mem := make([]int, ramSize)
		copy(mem, tokens)
		m := NewMachine(mem, func() int { return 0 })
		go m.Execute()
		for {
			x, ok := m.Read()
			if !ok {
				break
			}
			y, ok := m.Read()
			if !ok {
				log.Fatalln("Failed to read y")
			}
			tile, ok := m.Read()
			if !ok {
				log.Fatalln("Failed to read tile")
			}
			b.Emplace(Point{x, y}, tile)
		}
		fmt.Println(b.BlockCount())
	}
	{
		b := NewBoard()
		mem := make([]int, ramSize)
		copy(mem, tokens)
		mem[0] = 2
		m := NewMachine(mem, func() int { return 0 })
		go m.Execute()
		score := 0
		for {
			//b.Print()
			//fmt.Fprintf(os.Stderr, "Score: %d\n", score)
			x, ok := m.Read()
			if !ok {
				break
			}
			y, ok := m.Read()
			if !ok {
				log.Fatalln("Failed to read y")
			}
			tile, ok := m.Read()
			if !ok {
				log.Fatalln("Failed to read tile")
			}
			if x == -1 && y == 0 {
				score = tile
			} else {
				b.Emplace(Point{x, y}, tile)
			}
		}
		fmt.Println(score)
	}
}
