package main

import (
	"bufio"
	"fmt"
	"os"
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

type Coord struct {
	X, Y int
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func dir(a1, a2 int) int {
	if a1 == a2 {
		return 0
	} else if a1 < a2 {
		return 1
	} else {
		return -1
	}
}

func run() error {

	chart := make(map[Coord]int)
	numTwo := 0

	if err := doLines(os.Args[1], func(line string) error {
		var x1, y1, x2, y2 int

		n, err := fmt.Sscanf(line, "%d,%d -> %d,%d", &x1, &y1, &x2, &y2)
		if n != 4 || err != nil {
			return err
		}

		// Only considering horizontal/vertical for Part 1
		if len(os.Args) < 3 && (x1 != x2) && (y1 != y2) {
			return nil
		}

		dirX := dir(x1, x2)
		dirY := dir(y1, y2)

		for x, y := x1, y1; (x != (x2 + dirX)) || (y != (y2 + dirY)); x, y = x + dirX, y + dirY {
			c := Coord{ x, y }
			if v, ok := chart[c]; ok {
				// Already have at least one line at this Coord
				chart[c] = v + 1

				// Count specifically two or more for Part 1
				if v == 1 {
					numTwo++
				}
			} else {
				chart[c] = 1
			}
		}

		return nil
	}); err != nil {
		return err
	}

	fmt.Println(numTwo)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
