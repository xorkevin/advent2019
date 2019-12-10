package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strings"
)

const (
	puzzleInput = "input.txt"
)

func sign(a int) int {
	if a < 0 {
		return -1
	}
	return 1
}

func abs(a int) int {
	if a < 0 {
		return a * -1
	}
	return a
}

func gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

type (
	Angle struct {
		dx, dy int
	}

	AngleList []Angle

	Point struct {
		x, y int
	}

	PointList []Point

	Grid [][]int
)

func NewAngle(x2, y2, x1, y1 int) Angle {
	dx := x2 - x1
	dy := y2 - y1
	if dx == 0 {
		return Angle{
			dx: 0,
			dy: sign(dy),
		}
	}
	if dy == 0 {
		return Angle{
			dx: sign(dx),
			dy: 0,
		}
	}
	g := gcd(abs(dx), abs(dy))
	return Angle{
		dx: dx / g,
		dy: dy / g,
	}
}

func (a Angle) Arctan() float64 {
	k := math.Atan2(float64(a.dx), float64(-a.dy))
	if k < 0 {
		k += math.Pi * 2
	}
	return k
}

func (s AngleList) Len() int {
	return len(s)
}
func (s AngleList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s AngleList) Less(i, j int) bool {
	return s[i].Arctan() < s[j].Arctan()
}

func (p Point) Dist() int {
	return abs(p.x) + abs(p.y)
}

func (s PointList) Len() int {
	return len(s)
}
func (s PointList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s PointList) Less(i, j int) bool {
	return s[i].Dist() < s[j].Dist()
}

func (g Grid) Visible(px, py int) int {
	count := 0
	angles := map[Angle]struct{}{}
	for y, i := range g {
		for x, j := range i {
			if j != 1 || (x == px && y == py) {
				continue
			}
			k := NewAngle(x, y, px, py)
			if _, ok := angles[k]; !ok {
				angles[k] = struct{}{}
				count++
			}
		}
	}
	return count
}

func (g Grid) OrderAngles(px, py, count int) Point {
	points := map[Angle]PointList{}

	for y, i := range g {
		for x, j := range i {
			if j != 1 || (x == px && y == py) {
				continue
			}
			k := NewAngle(x, y, px, py)
			if _, ok := points[k]; !ok {
				points[k] = PointList{}
			}
			s := append(points[k], Point{x - px, y - py})
			sort.Sort(s)
			points[k] = s
		}
	}

	angles := make(AngleList, len(points))
	for k := range points {
		angles = append(angles, k)
	}

	sort.Sort(angles)

	k := Point{}
	i := 0
	n := 0
	for i < count {
		a := angles[n]
		p := points[a]
		if len(p) == 0 {
			n = (n + 1) % len(angles)
			continue
		}
		k = points[a][0]
		k = Point{k.x + px, k.y + py}
		points[a] = points[a][1:]
		i++
		n = (n + 1) % len(angles)
	}
	return k
}

func main() {
	grid := Grid{}
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
			chars := strings.Split(scanner.Text(), "")
			row := make([]int, 0, len(chars))
			for _, i := range chars {
				switch i {
				case ".":
					row = append(row, 0)
				case "#":
					row = append(row, 1)
				}
			}
			grid = append(grid, row)
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	max := 0
	maxX := 0
	maxY := 0
	for y, i := range grid {
		for x, j := range i {
			if j == 1 {
				k := grid.Visible(x, y)
				if k > max {
					max = k
					maxX = x
					maxY = y
				}
			}
		}
	}

	fmt.Println(max)

	k := grid.OrderAngles(maxX, maxY, 200)
	fmt.Println(k.x*100 + k.y)
}
