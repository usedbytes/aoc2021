package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"strconv"
)

func doLines(filename string, do func(line string) error) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if err := do(line); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

type Point struct {
	X, Y int
}

// My first time using container/heap. It seems rather cumbersome
// compared to how sort works.

type QueueNode struct {
	Point
	Guess int

	index int
}

type NodeQueue struct {
	list []*QueueNode         // will be kept ordered by heap
	dict map[Point]*QueueNode // for quick lookup/update
}

func (nq *NodeQueue) Len() int {
	return len(nq.list)
}

func (nq NodeQueue) Less(i, j int) bool {
	return nq.list[i].Guess < nq.list[j].Guess
}

func (nq NodeQueue) Swap(i, j int) {
	if i == j || i < 0 || j < 0 {
		return
	}
	nq.list[i], nq.list[j] = nq.list[j], nq.list[i]
	nq.list[i].index = i
	nq.list[j].index = j
}

func (nq *NodeQueue) Push(x interface{}) {
	node := x.(*QueueNode)
	if existing, ok := nq.dict[node.Point]; ok {
		existing.Guess = node.Guess
		return
	}

	n := len(nq.list)
	node.index = n
	nq.list = append(nq.list, node)
}

func (nq *NodeQueue) Pop() interface{} {
	if len(nq.list) == 0 {
		return nil
	}

	old := nq.list
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	delete(nq.dict, item.Point)
	nq.list = old[0 : n-1]
	return item
}

func run() error {
	cavern := [][]int{}

	if err := doLines(os.Args[1], func(line string) error {
		row := make([]int, len(line))
		for i := range line {
			v, err := strconv.Atoi(line[i:i+1])
			if err != nil {
				return err
			}
			row[i] = v
		}
		cavern = append(cavern, row)

		return nil
	}); err != nil {
		return err
	}

	sizeX := len(cavern[0])
	sizeY := len(cavern)

	// Implements A* search

	// Map from Point to best cost to reach Point
	tree := make(map[Point]int)

	// Queue of nodes to explore next
	var queue NodeQueue
	heap.Init(&queue)

	// Cost function from Point to goal
	nodeGuess := func(p Point) int {
		return (sizeX - p.X) + (sizeY - p.Y)
	}

	dirs := [][2]int {
		{  0, -1 },
		{ -1,  0 },
		{  1,  0 },
		{  0,  1 },
	}

	// Seed the start node
	tree[Point{0, 0}] = 0
	heap.Push(&queue, &QueueNode{
		Point: Point{0, 0},
		Guess: nodeGuess(Point{0, 0}),
	})
	goal := Point{ sizeX - 1, sizeY - 1 }

	// Go!
	for i := heap.Pop(&queue); i != nil ; i = heap.Pop(&queue) {
		node := i.(*QueueNode)
		current := node.Point

		if current == goal {
			break
		}

		for _, d := range dirs {
			newX, newY := current.X + d[0], current.Y + d[1]
			if newX < 0 || newX >= sizeX || newY < 0 || newY >= sizeY {
				continue
			}

			next := Point{ newX, newY }
			costNext := tree[current] + cavern[next.Y][next.X]
			if v, ok := tree[next]; !ok || costNext < v {
				// This is the best route to "next" we've
				// found so far
				tree[next] = costNext

				// Explore around 'next' with this new-found
				// route
				heap.Push(&queue, &QueueNode{
					Point: next,
					Guess: costNext + nodeGuess(next),
				});
			}
		}
	}

	fmt.Println(tree[goal])

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
