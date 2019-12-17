package main

import (
	"bufio"
	"bytes"
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

type (
	Bot struct {
		grid [][]byte
		w, h int
		x, y int
		dir  int
	}
)

func isBot(c byte) bool {
	return c == '^'
}

func NewBot(grid [][]byte) *Bot {
	bx := -1
	by := -1
	for y, i := range grid {
		for x, j := range i {
			if isBot(j) {
				bx = x
				by = y
				break
			}
		}
	}
	return &Bot{
		grid: grid,
		w:    len(grid[0]),
		h:    len(grid),
		x:    bx,
		y:    by,
		dir:  dirUp,
	}
}

func (b *Bot) inBounds(x, y int) bool {
	return y < b.h && y >= 0 && x >= 0 && x < b.w
}

func (b *Bot) isPath(x, y int) bool {
	if !b.inBounds(x, y) {
		return false
	}
	return b.grid[y][x] == '#'
}

func (b *Bot) isIntersection(x, y int) bool {
	if !b.isPath(x, y) {
		return false
	}

	return b.isPath(x, y-1) && b.isPath(x, y+1) && b.isPath(x-1, y) && b.isPath(x+1, y)
}

func (b *Bot) Sum() int {
	sum := 0
	for y, i := range b.grid {
		for x := range i {
			if b.isIntersection(x, y) {
				sum += x * y
			}
		}
	}
	return sum
}

func (b *Bot) getFLR() (bool, bool, bool) {
	lx := b.x
	ly := b.y
	rx := b.x
	ry := b.y
	fx := b.x
	fy := b.y

	switch b.dir {
	case dirUp:
		lx -= 1
		rx += 1
		fy -= 1
	case dirDown:
		lx += 1
		rx -= 1
		fy += 1
	case dirLeft:
		ly += 1
		ry -= 1
		fx -= 1
	case dirRight:
		ly -= 1
		ry += 1
		fx += 1
	default:
		log.Fatalln("Illegal direction")
	}

	return b.isPath(fx, fy), b.isPath(lx, ly), b.isPath(rx, ry)
}

func (b *Bot) forward() {
	switch b.dir {
	case dirUp:
		b.y -= 1
	case dirDown:
		b.y += 1
	case dirLeft:
		b.x -= 1
	case dirRight:
		b.x += 1
	default:
		log.Fatalln("Illegal direction")
	}
}

func (b *Bot) turnLeft() {
	switch b.dir {
	case dirUp:
		b.dir = dirLeft
	case dirDown:
		b.dir = dirRight
	case dirLeft:
		b.dir = dirDown
	case dirRight:
		b.dir = dirUp
	}
}

func (b *Bot) turnRight() {
	switch b.dir {
	case dirUp:
		b.dir = dirRight
	case dirDown:
		b.dir = dirLeft
	case dirLeft:
		b.dir = dirUp
	case dirRight:
		b.dir = dirDown
	}
}

func (b *Bot) FindDirections() string {
	instrs := bytes.Buffer{}
	for {
		f, l, r := b.getFLR()
		if !f && !l && !r {
			break
		}

		if f {
			b.forward()
			instrs.WriteByte('F')
		} else if l {
			b.turnLeft()
			instrs.WriteByte('L')
		} else {
			b.turnRight()
			instrs.WriteByte('R')
		}
	}

	run := 0
	s := strings.Builder{}
	for _, i := range instrs.Bytes() {
		if i == 'F' {
			run++
			continue
		}
		if run > 0 {
			s.WriteString(strconv.Itoa(run))
			s.WriteByte(',')
			run = 0
		}
		s.WriteByte(i)
		s.WriteByte(',')
	}
	if run > 0 {
		s.WriteString(strconv.Itoa(run))
		s.WriteByte(',')
		run = 0
	}
	return s.String()
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

	grid := [][]byte{}
	{
		mem := make([]int, ramSize)
		copy(mem, tokens)
		m := NewMachine(mem)
		go m.Execute()
		line := []byte{}
		for {
			out, ok := m.Read()
			if !ok {
				break
			}
			switch byte(out) {
			case '\n':
				if len(line) > 0 {
					grid = append(grid, line)
					line = []byte{}
				}
			default:
				line = append(line, byte(out))
			}
		}
	}

	b := NewBot(grid)
	fmt.Println(b.Sum())

	fmt.Println(b.FindDirections())

	instructions :=
		`A,A,B,C,C,A,C,B,C,B
L,4,L,4,L,6,R,10,L,6
L,12,L,6,R,10,L,6
R,8,R,10,L,6
n
`

	{
		mem := make([]int, ramSize)
		copy(mem, tokens)
		m := NewMachine(mem)
		mem[0] = 2
		go m.Execute()
		go func() {
			for _, c := range []byte(instructions) {
				m.Write(int(c))
			}
		}()
		for {
			out, ok := m.Read()
			if !ok {
				break
			}
			if out > 255 {
				fmt.Println(out)
			} else {
				fmt.Print(string(out))
			}
		}
	}
}
