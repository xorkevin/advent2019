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
	Point struct {
		x, y int
	}

	Maze struct {
		grid    [][]byte
		w, h    int
		enter   Point
		keys    []byte
		keyPos  map[byte]Point
		doors   []byte
		doorPos map[byte]Point
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
	enter := Point{-1, -1}
	keys := []byte{}
	keyPos := map[byte]Point{}
	doors := []byte{}
	doorPos := map[byte]Point{}

	for y, i := range grid {
		for x, j := range i {
			if isEntrance(j) {
				enter = Point{x, y}
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

func (m *Maze) CalcPath(start Point, goal Point, keys map[byte]struct{}) int {
	openSet := NewOpenSet()
	openSet.Push(start, 0, Manhattan(start, goal))
	closedSet := NewClosedSet()

	for !openSet.Empty() {
		cur, curg, _ := openSet.Pop()
		closedSet.Push(cur)
		if cur == goal {
			return curg
		}

		for _, neighbor := range m.neighbors(cur, keys) {
			if closedSet.Has(neighbor) || openSet.Has(neighbor) {
				continue
			}
			openSet.Push(neighbor, curg+1, curg+1+Manhattan(neighbor, goal))
		}
	}

	return -1
}

func (m *Maze) Reachable(start Point, keys map[byte]struct{}) []byte {
	reachable := []byte{}
	openSet := NewOpenSet()
	openSet.Push(start, 0, 0)
	closedSet := NewClosedSet()
	for !openSet.Empty() {
		cur, curg, _ := openSet.Pop()
		closedSet.Push(cur)
		k := m.grid[cur.y][cur.x]
		if isKey(k) {
			if _, ok := keys[k]; !ok {
				reachable = append(reachable, k)
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

func Perm(a []byte, f func([]byte)) {
	perm(a, f, 0)
}

// Permute the values at index i to len(a)-1.
func perm(a []byte, f func([]byte), i int) {
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

func (m *Maze) CalcMultiPath(start Point, order []byte) int {
	dist := 0
	prev := start
	keys := map[byte]struct{}{}
	for _, key := range order {
		goal := m.keyPos[key]
		partial := m.CalcPath(prev, goal, keys)
		if partial < 0 {
			return -1
		}
		dist += partial
		prev = goal
		keys[key] = struct{}{}
	}
	return dist
}

//func (m *Maze) Salesman(start Point, allkeys []byte) int {
//	minPath := -1
//	Perm(allkeys, func(order []byte) {
//		dist := m.CalcMultiPath(start, order)
//		if dist < 0 {
//			return
//		}
//		if minPath < 0 || dist < minPath {
//			minPath = dist
//		}
//	})
//	return minPath
//}

func (m *Maze) Salesman(start Point, dist int, keys map[byte]struct{}, curpath []byte) int {
	if len(keys) >= len(m.keyPos) {
		//fmt.Println(string(curpath), dist)
		return dist
	}

	minPath := -1
	reachable := m.Reachable(start, keys)
	for _, i := range reachable {
		goal := m.keyPos[i]
		partial := m.CalcPath(start, goal, keys)
		keys[i] = struct{}{}
		nextDist := m.Salesman(goal, dist+partial, keys, append(curpath, i))
		delete(keys, i)
		if nextDist < 0 {
			continue
		}
		if minPath < 0 || nextDist < minPath {
			minPath = nextDist
		}
	}
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

	maze := NewMaze(grid)
	for _, i := range maze.grid {
		fmt.Println(string(i))
	}
	fmt.Println(maze.Salesman(maze.enter, 0, map[byte]struct{}{}, []byte{}))
}
