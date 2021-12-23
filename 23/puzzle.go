package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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

type Cave [3][11]byte

type Position struct {
	X, Y int
}

func TargetRoomIdx(b byte) int {
	return int(b - byte('A'))
}

func RoomXPosition(i int) int {
	return 2 + (i * 2)
}

func PositionToRoomIdx(p Position) int {
	return (p.X - 2) / 2
}

func IsDoorPosition(i int) bool {
	if i > 0 && i < 10 && (i & 1) == 0 {
		return true
	}

	return false
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func (c Cave) At(p Position) byte {
	return c[p.Y][p.X]
}

func AllowedDestinations(c Cave, p Position) []Position {
	var poss []Position

	currentPos := p
	color := c.At(p)

	if !strings.Contains("ABCD", string(color)) {
		return poss
	}

	// Already in a room, so see if there are any hallway positions to
	// move to
	if currentPos.Y > 0 {
		currentRoom := PositionToRoomIdx(currentPos)
		targetRoom := TargetRoomIdx(color)

		// At the back of the room?
		if currentPos.Y == 2 {
			if currentRoom == targetRoom {
				// Don't want to move
				return poss
			}

			if c[1][currentPos.X] != '.' {
				// Can't move
				return poss
			}
		}

		if currentRoom == targetRoom && c.At(Position{currentPos.X, 2}) == color {
			// Don't want to move
			return poss
		}

		// All the hallway positions
		dirs := []int{ -1, 1 }
		for _, d := range dirs {
			p := currentPos
			for i := 0; i < 11; i++ {
				newX := p.X + (i * d)

				if newX < 0 || newX >= 10 {
					// Reached end of hallway
					continue
				}

				if IsDoorPosition(newX) {
					// Can't stop in front of a door
					continue
				}

				if c.At(Position{newX, 0}) != '.' {
					// Can't pass another pod
					break
				}

				// Could move into the hallway at newX
				poss = append(poss, Position{
					X: newX,
					Y: 0,
				})
			}
		}
	}

	// Any room positions to move to?
	room := TargetRoomIdx(color)
	roomX := RoomXPosition(room)

	frontOfRoom := c.At(Position{roomX, 1})
	backOfRoom := c.At(Position{roomX, 2})
	roomFull := (frontOfRoom != '.')
	roomEmpty := (backOfRoom == '.')

	if roomFull {
		// Can't get in
		return poss
	}

	if !roomEmpty && backOfRoom != color {
		// Will refuse to get in
		return poss
	}

	// Check if the hallway is clear all the way to the room
	dir := roomX - currentPos.X
	dir /= abs(dir)
	for i := currentPos.X + dir; i != roomX; i += dir {
		// Can't get past another amphipod
		if c.At(Position{i, 0}) != '.' {
			return poss
		}
	}

	if roomEmpty {
		// Take the bottom slot
		poss = append(poss, Position{
			X: roomX,
			Y: 2,
		})
	} else {
		// Take the top slot
		poss = append(poss, Position{
			X: roomX,
			Y: 1,
		})
	}

	return poss
}

func (c Cave) Print() {
	fmt.Println("---")
	for _, row := range c {
		fmt.Printf("%s\n", row)
	}
	fmt.Println("---")
}

func FindPods(c Cave) []Position {
	res := make([]Position, 0, 8)
	for y, row := range c {
		for x, b := range row {
			if strings.Contains("ABCD", string(b)) {
				res = append(res, Position{x, y})
			}
		}
	}

	return res
}

func Distance(from, to Position) int {
	return abs(to.X - from.X) + from.Y + to.Y
}

func (c Cave) Move(from, to Position) (Cave, int) {
	distance := Distance(from, to)
	pricePerMove := []int{1, 10, 100, 1000}
	colorIdx := int(c.At(from) - 'A')

	r := c
	r[from.Y][from.X], r[to.Y][to.X] = r[to.Y][to.X], r[from.Y][from.X]

	return r, distance * pricePerMove[colorIdx]
}

func (c Cave) IsSolved() bool {
	return c.At(Position{2, 1}) == 'A' &&
		c.At(Position{2, 2}) == 'A' &&
		c.At(Position{4, 1}) == 'B' &&
		c.At(Position{4, 2}) == 'B' &&
		c.At(Position{6, 1}) == 'C' &&
		c.At(Position{6, 2}) == 'C' &&
		c.At(Position{8, 1}) == 'D' &&
		c.At(Position{8, 2}) == 'D'
}

func solve(c Cave, cost int, path [][2]Position) (int, [][2]Position) {
	if c.IsSolved() {
		return cost, path
	}

	var minPath [][2]Position

	minCost := -1
	pods := FindPods(c)
	for _, p := range pods {
		moves := AllowedDestinations(c, p)
		for _, m := range moves {
			d, price := c.Move(p, m)
			newCost, newPath := solve(d, cost + price, append(path, [2]Position{p, m}))
			if minCost < 0 || (newCost > 0 && newCost < minCost) {
				minPath = make([][2]Position, len(newPath))
				copy(minPath, newPath)
				minCost = newCost
			}
		}
	}

	return minCost, minPath
}

func run() error {
	var cave Cave

	lineNo := 0
	if err := doLines(os.Args[1], func(line string) error {
		lineNo++
		if lineNo < 2 || lineNo > 4 {
			return nil
		}

		y := lineNo - 2
		for x, c := range line[1:] {
			if x > len(cave[0]) - 1 {
				break
			}
			cave[y][x] = byte(c)
		}

		return nil
	}); err != nil {
		return err
	}

	cave.Print()

	totalCost, moves := solve(cave, 0, nil)

	cost := 0
	dave := cave
	dave.Print()
	for i, m := range moves {
		fmt.Println("Move", i, m)
		dave, cost = dave.Move(m[0], m[1])
		fmt.Println("Cost", cost)
		dave.Print()
	}

	fmt.Println(totalCost)
	fmt.Println(moves)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
