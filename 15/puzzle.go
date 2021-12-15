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

func doAStar(start, goal Point, cavern func(Point) int) int {
	// Map from Point to best cost to reach Point
	tree := make(map[Point]int)

	// Queue of nodes to explore next
	var queue NodeQueue
	heap.Init(&queue)

	// Cost function from Point to goal
	nodeGuess := func(p Point) int {
		return (goal.X - p.X) + (goal.Y - p.Y)
	}

	dirs := [][2]int {
		{  0, -1 },
		{ -1,  0 },
		{  1,  0 },
		{  0,  1 },
	}

	// Seed the start node
	tree[start] = 0
	heap.Push(&queue, &QueueNode{
		Point: start,
		Guess: nodeGuess(start),
	})

	// Go!
	for i := heap.Pop(&queue); i != nil ; i = heap.Pop(&queue) {
		node := i.(*QueueNode)
		current := node.Point

		if current == goal {
			break
		}

		for _, d := range dirs {
			next := Point{ current.X + d[0], current.Y + d[1] }
			v := cavern(next)
			if v < 0 {
				continue
			}

			costNext := tree[current] + v
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

	return tree[goal]
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

	start := Point{0, 0}

	part1Cavern := func(p Point) int {
		if p.X < 0 || p.X >= len(cavern[0]) || p.Y < 0 || p.Y >= len(cavern) {
			return -1
		}

		return cavern[p.Y][p.X]
	}
	part1Goal := Point{ len(cavern[0]) - 1, len(cavern) - 1 }
	part1 := doAStar(start, part1Goal, part1Cavern)
	fmt.Println("Part 1:", part1)

	part2Cavern := func(p Point) int {
		xTile, yTile := p.X / len(cavern[0]), p.Y / len(cavern)

		if p.X < 0 || p.Y < 0 {
			return -1
		}

		if xTile >= 5 || yTile >= 5 {
			return -1
		}

		xOffs, yOffs := p.X % len(cavern[0]), p.Y % len(cavern)

		// Annoying number system is 1-9
		// Subtract 1, wrap to 0-8, add 1
		v := cavern[yOffs][xOffs] - 1
		v = (v + xTile + yTile) % 9
		v += 1

		return v
	}
	part2Goal := Point{ len(cavern[0]) * 5 - 1, len(cavern) * 5 - 1 }
	part2 := doAStar(start, part2Goal, part2Cavern)
	fmt.Println("Part 2:", part2)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
