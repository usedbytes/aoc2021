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
		if (x1 != x2) && (y1 != y2) {
			return nil
		}

		// TODO: This doesn't work if lines aren't horizontal/vertical
		for y := min(y1, y2); y <= max(y1, y2); y++ {
			for x := min(x1, x2); x <= max(x1, x2); x++ {
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
