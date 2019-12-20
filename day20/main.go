package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"os"
)

const (
	puzzleInput = "input.txt"
)

type (
	Point2 struct {
		x, y int
	}

	Point struct {
		x, y, z int
	}

	Maze struct {
		grid      [][]byte
		ringWidth int
		size      int
		tele      map[Point2]Point2
		start     Point
		end       Point
	}
)

func isPath(c byte) bool {
	return c == '.'
}

func isWall(c byte) bool {
	return c == '#'
}

func isMaze(c byte) bool {
	return isPath(c) || isWall(c)
}

func isTeleMark(c byte) bool {
	return c >= 'A' && c <= 'Z'
}

func NewMaze(grid [][]byte) *Maze {
	l := len(grid)
	mid := l / 2
	ringWidth := -1
	for i := 2; i < mid; i++ {
		if !isMaze(grid[mid][i]) {
			ringWidth = i - 2
			break
		}
	}
	if ringWidth < 0 {
		log.Fatalln("Invalid map")
	}

	tele := map[Point2]Point2{}
	teleEnd := map[string]Point2{}
	for i := 0; i < l; i++ {
		if isTeleMark(grid[0][i]) {
			p := Point2{i, 2}
			mark := string([]byte{grid[0][i], grid[1][i]})
			if v, ok := teleEnd[mark]; ok {
				tele[p] = v
				tele[v] = p
				delete(teleEnd, mark)
			} else {
				teleEnd[mark] = p
			}
		}
		if isTeleMark(grid[2+ringWidth][i]) {
			p := Point2{i, 1 + ringWidth}
			mark := string([]byte{grid[2+ringWidth][i], grid[3+ringWidth][i]})
			if v, ok := teleEnd[mark]; ok {
				tele[p] = v
				tele[v] = p
				delete(teleEnd, mark)
			} else {
				teleEnd[mark] = p
			}
		}
		if isTeleMark(grid[l-2][i]) {
			p := Point2{i, l - 3}
			mark := string([]byte{grid[l-2][i], grid[l-1][i]})
			if v, ok := teleEnd[mark]; ok {
				tele[p] = v
				tele[v] = p
				delete(teleEnd, mark)
			} else {
				teleEnd[mark] = p
			}
		}
		if isTeleMark(grid[l-ringWidth-4][i]) {
			p := Point2{i, l - ringWidth - 2}
			mark := string([]byte{grid[l-ringWidth-4][i], grid[l-ringWidth-3][i]})
			if v, ok := teleEnd[mark]; ok {
				tele[p] = v
				tele[v] = p
				delete(teleEnd, mark)
			} else {
				teleEnd[mark] = p
			}
		}
		if isTeleMark(grid[i][0]) {
			p := Point2{2, i}
			mark := string([]byte{grid[i][0], grid[i][1]})
			if v, ok := teleEnd[mark]; ok {
				tele[p] = v
				tele[v] = p
				delete(teleEnd, mark)
			} else {
				teleEnd[mark] = p
			}
		}
		if isTeleMark(grid[i][2+ringWidth]) {
			p := Point2{1 + ringWidth, i}
			mark := string([]byte{grid[i][2+ringWidth], grid[i][3+ringWidth]})
			if v, ok := teleEnd[mark]; ok {
				tele[p] = v
				tele[v] = p
				delete(teleEnd, mark)
			} else {
				teleEnd[mark] = p
			}
		}
		if isTeleMark(grid[i][l-2]) {
			p := Point2{l - 3, i}
			mark := string([]byte{grid[i][l-2], grid[i][l-1]})
			if v, ok := teleEnd[mark]; ok {
				tele[p] = v
				tele[v] = p
				delete(teleEnd, mark)
			} else {
				teleEnd[mark] = p
			}
		}
		if isTeleMark(grid[i][l-ringWidth-4]) {
			p := Point2{l - ringWidth - 2, i}
			mark := string([]byte{grid[i][l-ringWidth-4], grid[i][l-ringWidth-3]})
			if v, ok := teleEnd[mark]; ok {
				tele[p] = v
				tele[v] = p
				delete(teleEnd, mark)
			} else {
				teleEnd[mark] = p
			}
		}
	}

	start, ok := teleEnd["AA"]
	if !ok {
		log.Fatalln("No maze start")
	}
	end, ok := teleEnd["ZZ"]
	if !ok {
		log.Fatalln("No maze end")
	}

	return &Maze{
		grid:      grid,
		ringWidth: ringWidth,
		size:      l,
		tele:      tele,
		start:     Point{start.x, start.y, 0},
		end:       Point{end.x, end.y, 0},
	}
}

type (
	Item struct {
		value Point
		g, f  int
		index int
	}

	PriorityQueue []*Item

	OpenSet struct {
		pq     PriorityQueue
		valSet map[Point]struct{}
	}

	ClosedSet map[Point]struct{}
)

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].f < pq[j].f
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func NewOpenSet() *OpenSet {
	return &OpenSet{
		pq:     PriorityQueue{},
		valSet: map[Point]struct{}{},
	}
}

func (os *OpenSet) Empty() bool {
	return os.pq.Len() == 0
}

func (os *OpenSet) Has(val Point) bool {
	_, ok := os.valSet[val]
	return ok
}

func (os *OpenSet) Push(value Point, g, f int) {
	item := &Item{
		value: value,
		g:     g,
		f:     f,
	}
	heap.Push(&os.pq, item)
	os.valSet[value] = struct{}{}
}

func (os *OpenSet) Pop() (Point, int, int) {
	item := heap.Pop(&os.pq).(*Item)
	delete(os.valSet, item.value)
	return item.value, item.g, item.f
}

func NewClosedSet() ClosedSet {
	return ClosedSet{}
}

func (cs ClosedSet) Has(val Point) bool {
	_, ok := cs[val]
	return ok
}

func (cs ClosedSet) Push(val Point) {
	cs[val] = struct{}{}
}

func (m *Maze) isPath(pos Point) bool {
	return isPath(m.grid[pos.y][pos.x])
}

func (m *Maze) neighbors(pos Point) []Point {
	points := make([]Point, 0, 5)
	if k := (Point{pos.x, pos.y - 1, 0}); m.isPath(k) {
		points = append(points, k)
	}
	if k := (Point{pos.x, pos.y + 1, 0}); m.isPath(k) {
		points = append(points, k)
	}
	if k := (Point{pos.x - 1, pos.y, 0}); m.isPath(k) {
		points = append(points, k)
	}
	if k := (Point{pos.x + 1, pos.y, 0}); m.isPath(k) {
		points = append(points, k)
	}
	if v, ok := m.tele[Point2{pos.x, pos.y}]; ok {
		points = append(points, Point{v.x, v.y, 0})
	}
	return points
}

func (m *Maze) BFS() int {
	openSet := NewOpenSet()
	openSet.Push(m.start, 0, 0)
	closedSet := NewClosedSet()
	for !openSet.Empty() {
		cur, curg, _ := openSet.Pop()
		closedSet.Push(cur)
		if cur == m.end {
			return curg
		}

		for _, neighbor := range m.neighbors(cur) {
			if closedSet.Has(neighbor) || openSet.Has(neighbor) {
				continue
			}
			openSet.Push(neighbor, curg+1, curg+1)
		}
	}
	return -1
}

func (m *Maze) isOuter(pos Point) bool {
	if pos.x == 2 || pos.y == 2 {
		return true
	}
	if pos.x == m.size-3 || pos.y == m.size-3 {
		return true
	}
	return false
}

func (m *Maze) neighbors2(pos Point) []Point {
	points := make([]Point, 0, 5)
	if k := (Point{pos.x, pos.y - 1, pos.z}); m.isPath(k) {
		points = append(points, k)
	}
	if k := (Point{pos.x, pos.y + 1, pos.z}); m.isPath(k) {
		points = append(points, k)
	}
	if k := (Point{pos.x - 1, pos.y, pos.z}); m.isPath(k) {
		points = append(points, k)
	}
	if k := (Point{pos.x + 1, pos.y, pos.z}); m.isPath(k) {
		points = append(points, k)
	}
	if v, ok := m.tele[Point2{pos.x, pos.y}]; ok {
		if m.isOuter(pos) {
			if pos.z > 0 {
				k := Point{v.x, v.y, pos.z - 1}
				points = append(points, k)
			}
		} else {
			k := Point{v.x, v.y, pos.z + 1}
			points = append(points, k)
		}
	}
	return points
}

func (m *Maze) BFS2() int {
	openSet := NewOpenSet()
	openSet.Push(m.start, 0, 0)
	closedSet := NewClosedSet()
	for !openSet.Empty() {
		cur, curg, _ := openSet.Pop()
		closedSet.Push(cur)
		if cur == m.end {
			return curg
		}

		for _, neighbor := range m.neighbors2(cur) {
			if closedSet.Has(neighbor) || openSet.Has(neighbor) {
				continue
			}
			openSet.Push(neighbor, curg+1, curg+1)
		}
	}
	return -1
}

func (m *Maze) Print() {
	for _, i := range m.grid {
		fmt.Println(string(i))
	}
}

func main() {
	grid := [][]byte{}
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
			grid = append(grid, []byte(scanner.Text()))
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	m := NewMaze(grid)
	fmt.Println(m.BFS())
	fmt.Println(m.BFS2())
}
