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

const (
	dirUp = iota
	dirDown
	dirLeft
	dirRight
)

const (
	colorBlack = 0
	colorWhite = 1
)

const (
	turnLeft  = 0
	turnRight = 1
)

type (
	Robot struct {
		x, y  int
		dir   int
		board map[Point]int
		maxXL, maxXR,
		maxYT, maxYB int
	}

	Point struct {
		x, y int
	}
)

func NewRobot() *Robot {
	return &Robot{
		x:     0,
		y:     0,
		dir:   dirUp,
		board: map[Point]int{},
		maxXL: 0,
		maxXR: 0,
		maxYT: 0,
		maxYB: 0,
	}
}

func (r *Robot) getPos() Point {
	return Point{x: r.x, y: r.y}
}

func (r *Robot) getPaint() int {
	val, ok := r.board[r.getPos()]
	if ok {
		return val
	}
	return colorBlack
}

func (r *Robot) paint(color int) {
	if color != colorBlack && color != colorWhite {
		log.Fatal("Invalid color")
	}
	r.board[r.getPos()] = color
}

func (r *Robot) turn(turn int) {
	switch turn {
	case turnLeft:
		switch r.dir {
		case dirUp:
			r.dir = dirLeft
		case dirDown:
			r.dir = dirRight
		case dirLeft:
			r.dir = dirDown
		case dirRight:
			r.dir = dirUp
		}
	case turnRight:
		switch r.dir {
		case dirUp:
			r.dir = dirRight
		case dirDown:
			r.dir = dirLeft
		case dirLeft:
			r.dir = dirUp
		case dirRight:
			r.dir = dirDown
		}
	default:
		log.Fatal("Invalid turn")
	}
}

func (r *Robot) forward() {
	switch r.dir {
	case dirUp:
		r.y -= 1
		if r.y < r.maxYT {
			r.maxYT = r.y
		}
	case dirDown:
		r.y += 1
		if r.y > r.maxYB {
			r.maxYB = r.y
		}
	case dirLeft:
		r.x -= 1
		if r.x < r.maxXL {
			r.maxXL = r.x
		}
	case dirRight:
		r.x += 1
		if r.x > r.maxXR {
			r.maxXR = r.x
		}
	}
}

func (r *Robot) Print() {
	width := r.maxXR - r.maxXL + 1
	height := r.maxYB - r.maxYT + 1
	offsetY := r.maxYT
	offsetX := r.maxXL
	grid := make([][]int, 0, height)
	for i := 0; i < height; i++ {
		grid = append(grid, make([]int, width))
	}
	for k, v := range r.board {
		grid[k.y-offsetY][k.x-offsetX] = v
	}
	for _, i := range grid {
		for _, j := range i {
			switch j {
			case colorBlack:
				fmt.Print(" ")
			case colorWhite:
				fmt.Print("#")
			}
		}
		fmt.Println()
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
		r := NewRobot()
		mem := make([]int, ramSize)
		copy(mem, tokens)
		m := NewMachine(mem)
		go m.Execute()
		for {
			curColor := r.getPaint()
			m.Write(curColor)
			nextColor, ok := m.Read()
			if !ok {
				break
			}
			turn, ok := m.Read()
			if !ok {
				log.Fatal("Must read a turn")
			}
			r.paint(nextColor)
			r.turn(turn)
			r.forward()
		}
		fmt.Println(len(r.board))
	}
	{
		r := NewRobot()
		r.paint(colorWhite)
		mem := make([]int, ramSize)
		copy(mem, tokens)
		m := NewMachine(mem)
		go m.Execute()
		for {
			curColor := r.getPaint()
			m.Write(curColor)
			nextColor, ok := m.Read()
			if !ok {
				break
			}
			turn, ok := m.Read()
			if !ok {
				log.Fatal("Must read a turn")
			}
			r.paint(nextColor)
			r.turn(turn)
			r.forward()
		}
		r.Print()
	}
}
