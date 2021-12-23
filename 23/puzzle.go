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

type Cave [5][11]byte

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

func (c Cave) At(X, Y int) byte {
	return c[Y][X]
}

func AllowedDestinations(c Cave, p Position) []Position {
	var poss []Position

	currentPos := p
	color := c.At(p.X, p.Y)

	if !strings.Contains("ABCD", string(color)) {
		return poss
	}

	// Already in a room, so see if there are any hallway positions to
	// move to
	if currentPos.Y > 0 {
		currentRoom := PositionToRoomIdx(currentPos)
		targetRoom := TargetRoomIdx(color)

		if currentRoom == targetRoom {
			// Check if everyone else in the room is already the target color
			homogenous := true
			for y := currentPos.Y + 1; y < len(c); y++ {
				at := c.At(currentPos.X, y)
				if at == '#' {
					break
				}

				if at != color {
					homogenous = false
					break
				}
			}

			if homogenous {
				// Don't want to move
				return nil
			}
		}

		// Can we get out of the room?
		for y := currentPos.Y - 1; y > 0; y-- {
			if c.At(currentPos.X, y) != '.' {
				// Can't move
				return nil
			}
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

				if c.At(newX, 0) != '.' {
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

	frontOfRoom := c.At(roomX, 1)
	roomFull := (frontOfRoom != '.')

	if roomFull {
		// Can't get in
		return poss
	}

	for y := 1; y < len(c); y++ {
		at := c.At(roomX, y)
		if at == '#' {
			break
		}

		if at != '.' && at != color {
			// Will refuse to get in
			return poss
		}
	}

	// Check if the hallway is clear all the way to the room
	dir := roomX - currentPos.X
	dir /= abs(dir)
	for i := currentPos.X + dir; i != roomX; i += dir {
		// Can't get past another amphipod
		if c.At(i, 0) != '.' {
			return poss
		}
	}

	// We can make it to the room! Take the lowest position available
	for y := len(c) - 1; y > 0; y-- {
		at := c.At(roomX, y)
		if at == '#' {
			continue
		}

		if at == '.' {
			poss = append(poss, Position{
				X: roomX,
				Y: y,
			})
			break
		}
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
	colorIdx := int(c.At(from.X, from.Y) - 'A')

	r := c
	r[from.Y][from.X], r[to.Y][to.X] = r[to.Y][to.X], r[from.Y][from.X]

	return r, distance * pricePerMove[colorIdx]
}

func (c Cave) IsSolved() bool {
	for room := 0; room < 4; room++ {
		x := RoomXPosition(room)
		for y := 1; y < len(c); y++ {
			at := c.At(x, y)
			if at == '#' {
				break
			}

			if at != 'A' + byte(room) {
				return false
			}
		}
	}

	return true
}

func solve(c Cave, dp map[Cave]int) int {
	if c.IsSolved() {
		return 0
	}

	if v, ok := dp[c]; ok {
		return v
	}

	minCost := -1
	pods := FindPods(c)
	for _, p := range pods {
		moves := AllowedDestinations(c, p)
		for _, m := range moves {
			d, moveCost := c.Move(p, m)

			newCost := solve(d, dp)
			if newCost >= 0 {
				newCost += moveCost
				if minCost < 0 || (newCost < minCost) {
					minCost = newCost
				}
			}
		}
	}

	dp[c] = minCost

	return minCost
}

func run() error {
	var cave Cave

	part2 := (len(os.Args) > 2)

	y := -1
	if err := doLines(os.Args[1], func(line string) error {
		if y < 0 || y > len(cave) - 1 {
			y++
			return nil
		}

		for x, c := range line[1:] {
			if x > len(cave[0]) - 1 {
				break
			}
			cave[y][x] = byte(c)
		}

		y++

		if part2 && y == 2 {
			copy(cave[2][:], []byte(" #D#C#B#A#"))
			copy(cave[3][:], []byte(" #D#B#A#C#"))
			y += 2
		}

		return nil
	}); err != nil {
		return err
	}

	cave.Print()

	dp := make(map[Cave]int)
	totalCost := solve(cave, dp)

	fmt.Println(totalCost)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
