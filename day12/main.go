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
	Vec struct {
		x, y, z int
	}

	Moon struct {
		pos Vec
		vel Vec
	}
)

func NewMoon(pos Vec) *Moon {
	return &Moon{
		pos: pos,
		vel: Vec{0, 0, 0},
	}
}

func gravity2(x, ox int) int {
	if x < ox {
		return +1
	}
	if x > ox {
		return -1
	}
	return 0
}

func (m *Moon) Gravity(other *Moon) {
	dx := gravity2(m.pos.x, other.pos.x)
	dy := gravity2(m.pos.y, other.pos.y)
	dz := gravity2(m.pos.z, other.pos.z)

	m.vel = Vec{
		x: m.vel.x + dx,
		y: m.vel.y + dy,
		z: m.vel.z + dz,
	}
}

func (m *Moon) Velocity() {
	m.pos = Vec{
		x: m.pos.x + m.vel.x,
		y: m.pos.y + m.vel.y,
		z: m.pos.z + m.vel.z,
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (m *Moon) Energy() int {
	e1 := abs(m.pos.x) + abs(m.pos.y) + abs(m.pos.z)
	e2 := abs(m.vel.x) + abs(m.vel.y) + abs(m.vel.z)
	return e1 * e2
}

func (m *Moon) GetX() string {
	return strconv.Itoa(m.pos.x) + "," + strconv.Itoa(m.vel.x)
}
func (m *Moon) GetY() string {
	return strconv.Itoa(m.pos.y) + "," + strconv.Itoa(m.vel.y)
}
func (m *Moon) GetZ() string {
	return strconv.Itoa(m.pos.z) + "," + strconv.Itoa(m.vel.z)
}

type (
	GravSystem struct {
		moons []*Moon
	}
)

func NewGravSystem(moons []*Moon) *GravSystem {
	return &GravSystem{
		moons: moons,
	}
}

func (g *GravSystem) stepGrav() {
	for n, i := range g.moons {
		for n2, j := range g.moons {
			if n == n2 {
				continue
			}
			i.Gravity(j)
		}
	}
}

func (g *GravSystem) stepVel() {
	for _, i := range g.moons {
		i.Velocity()
	}
}

func (g *GravSystem) Step() {
	g.stepGrav()
	g.stepVel()
}

func (g *GravSystem) Energy() int {
	count := 0
	for _, i := range g.moons {
		count += i.Energy()
	}
	return count
}

func (g *GravSystem) GetX() string {
	s := strings.Builder{}
	for _, i := range g.moons {
		s.WriteString(i.GetX())
		s.WriteString("\n")
	}
	return s.String()
}
func (g *GravSystem) GetY() string {
	s := strings.Builder{}
	for _, i := range g.moons {
		s.WriteString(i.GetY())
		s.WriteString("\n")
	}
	return s.String()
}
func (g *GravSystem) GetZ() string {
	s := strings.Builder{}
	for _, i := range g.moons {
		s.WriteString(i.GetZ())
		s.WriteString("\n")
	}
	return s.String()
}

func parseLine2(text string) int {
	l := strings.Split(text, "=")
	num, err := strconv.Atoi(l[1])
	if err != nil {
		log.Fatal(err)
	}
	return num
}

func parseLine(text string) Vec {
	line := strings.Split(strings.Trim(text, "<>"), ", ")
	x := parseLine2(line[0])
	y := parseLine2(line[1])
	z := parseLine2(line[2])
	return Vec{x, y, z}
}

func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func LCM(a, b int, integers ...int) int {
	result := a * b / GCD(a, b)

	for i := 0; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}

	return result
}

func main() {
	{
		moons := []*Moon{}
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
				moons = append(moons, NewMoon(parseLine(scanner.Text())))
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
		}

		gravsys := NewGravSystem(moons)
		for i := 0; i < 1000; i++ {
			gravsys.Step()
		}
		fmt.Println(gravsys.Energy())
	}
	{
		moons := []*Moon{}
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
				moons = append(moons, NewMoon(parseLine(scanner.Text())))
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
		}

		prevX := map[string]struct{}{}
		prevY := map[string]struct{}{}
		prevZ := map[string]struct{}{}
		xmatch := -1
		ymatch := -1
		zmatch := -1

		gravsys := NewGravSystem(moons)
		prevX[gravsys.GetX()] = struct{}{}
		prevY[gravsys.GetY()] = struct{}{}
		prevZ[gravsys.GetZ()] = struct{}{}
		for i := 0; xmatch < 0 || ymatch < 0 || zmatch < 0; i++ {
			gravsys.Step()
			if xmatch < 0 {
				x := gravsys.GetX()
				if _, ok := prevX[x]; ok {
					xmatch = i + 1
				} else {
					prevX[x] = struct{}{}
				}
			}
			if ymatch < 0 {
				y := gravsys.GetY()
				if _, ok := prevY[y]; ok {
					ymatch = i + 1
				} else {
					prevY[y] = struct{}{}
				}
			}
			if zmatch < 0 {
				z := gravsys.GetZ()
				if _, ok := prevZ[z]; ok {
					zmatch = i + 1
				} else {
					prevZ[z] = struct{}{}
				}
			}
		}

		fmt.Println(LCM(xmatch, ymatch, zmatch))
	}
}
