package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	puzzleInput = "input.txt"
)

type (
	Point struct {
		x, y int
	}

	Maze struct {
		grid    [][]byte
		w, h    int
		enter   []Point
		keys    []byte
		keyPos  map[byte]Point
		doors   []byte
		doorPos map[byte]Point
		reCache map[string]int
	}
)

func isEntrance(c byte) bool {
	return c == '@'
}

func isKey(c byte) bool {
	return c >= 'a' && c <= 'z'
}

func isDoor(c byte) bool {
	return c >= 'A' && c <= 'Z'
}

func isWall(c byte) bool {
	return c == '#'
}

func isPath(c byte) bool {
	return c == '.' || isEntrance(c) || isKey(c)
}

func NewMaze(grid [][]byte) *Maze {
	enter := []Point{}
	keys := []byte{}
	keyPos := map[byte]Point{}
	doors := []byte{}
	doorPos := map[byte]Point{}

	for y, i := range grid {
		for x, j := range i {
			if isEntrance(j) {
				enter = append(enter, Point{x, y})
			} else if isKey(j) {
				keys = append(keys, j)
				keyPos[j] = Point{x, y}
			} else if isDoor(j) {
				doors = append(doors, j)
				doorPos[j] = Point{x, y}
			}
		}
	}

	return &Maze{
		grid:    grid,
		w:       len(grid[0]),
		h:       len(grid),
		enter:   enter,
		keys:    keys,
		keyPos:  keyPos,
		doors:   doors,
		doorPos: doorPos,
		reCache: map[string]int{},
	}
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func Manhattan(a, b Point) int {
	return abs(a.x-b.x) + abs(a.y-b.y)
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

func (m *Maze) isPath(pos Point, keys map[byte]struct{}) bool {
	c := m.grid[pos.y][pos.x]
	if isPath(c) {
		return true
	}
	if !isDoor(c) {
		return false
	}
	_, ok := keys[c-'A'+'a']
	return ok
}

func (m *Maze) neighbors(pos Point, keys map[byte]struct{}) []Point {
	points := make([]Point, 0, 4)
	if k := (Point{pos.x, pos.y - 1}); m.isPath(k, keys) {
		points = append(points, k)
	}
	if k := (Point{pos.x, pos.y + 1}); m.isPath(k, keys) {
		points = append(points, k)
	}
	if k := (Point{pos.x - 1, pos.y}); m.isPath(k, keys) {
		points = append(points, k)
	}
	if k := (Point{pos.x + 1, pos.y}); m.isPath(k, keys) {
		points = append(points, k)
	}
	return points
}

func (m *Maze) Reachable(start Point, keys map[byte]struct{}) map[byte]int {
	reachable := map[byte]int{}
	openSet := NewOpenSet()
	openSet.Push(start, 0, 0)
	closedSet := NewClosedSet()
	for !openSet.Empty() {
		cur, curg, _ := openSet.Pop()
		closedSet.Push(cur)
		k := m.grid[cur.y][cur.x]
		if isKey(k) {
			if _, ok := keys[k]; !ok {
				reachable[k] = curg
				continue
			}
		}

		for _, neighbor := range m.neighbors(cur, keys) {
			if closedSet.Has(neighbor) || openSet.Has(neighbor) {
				continue
			}
			openSet.Push(neighbor, curg+1, curg+1)
		}
	}
	return reachable
}

func (m *Maze) ReachableN(start []Point, keys map[byte]struct{}) []map[byte]int {
	reachable := make([]map[byte]int, 0, len(start))
	for _, i := range start {
		reachable = append(reachable, m.Reachable(i, keys))
	}
	return reachable
}

func ToState(pos []Point, keys map[byte]struct{}) string {
	s := strings.Builder{}
	for _, i := range pos {
		s.WriteString(strconv.Itoa(i.x))
		s.WriteByte(',')
		s.WriteString(strconv.Itoa(i.y))
		s.WriteByte(';')
	}
	keySlice := make([]byte, 0, len(keys))
	for k := range keys {
		keySlice = append(keySlice, k)
	}
	sort.Slice(keySlice, func(i, j int) bool { return keySlice[i] < keySlice[j] })
	s.Write(keySlice)
	return s.String()
}

func (m *Maze) Salesman(start []Point, keys map[byte]struct{}) int {
	if len(keys) >= len(m.keyPos) {
		return 0
	}

	stateID := ToState(start, keys)
	if val, ok := m.reCache[stateID]; ok {
		return val
	}

	minPath := -1
	for n, reachable := range m.ReachableN(start, keys) {
		for k, i := range reachable {
			goal := m.keyPos[k]
			keys[k] = struct{}{}
			next := make([]Point, len(start))
			copy(next, start)
			next[n] = goal
			partial := i + m.Salesman(next, keys)
			delete(keys, k)
			if partial < 0 {
				continue
			}
			if minPath < 0 || partial < minPath {
				minPath = partial
			}
		}
	}

	m.reCache[stateID] = minPath
	return minPath
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

	{
		maze := NewMaze(grid)
		fmt.Println(maze.Salesman(maze.enter, map[byte]struct{}{}))
	}
	{
		px := -1
		py := -1
		for y, i := range grid {
			for x, j := range i {
				if j == '@' {
					px = x
					py = y
					break
				}
			}
		}
		grid[py][px] = '#'
		grid[py-1][px] = '#'
		grid[py+1][px] = '#'
		grid[py][px-1] = '#'
		grid[py][px+1] = '#'
		grid[py-1][px-1] = '@'
		grid[py-1][px+1] = '@'
		grid[py+1][px-1] = '@'
		grid[py+1][px+1] = '@'
		maze := NewMaze(grid)
		fmt.Println(maze.Salesman(maze.enter, map[byte]struct{}{}))
	}
}
