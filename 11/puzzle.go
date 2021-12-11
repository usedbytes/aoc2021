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

type Point struct {
	X, Y int
}

func run() error {
	cavern := make(map[Point]int)

	y := 0
	if err := doLines(os.Args[1], func(line string) error {
		for i, l := range line {
			pt := Point{ X: i, Y: y }

			cavern[pt] = int(l - '0')
		}

		y++
		return nil
	}); err != nil {
		return err
	}

	nsteps := 100
	flashes := 0

	adjacent_dirs := [][2]int{
		{ -1, -1, },
		{  0, -1, },
		{  1, -1, },
		{ -1,  0, },
		// Skip self
		{  1,  0, },
		{ -1,  1, },
		{  0,  1, },
		{  1,  1, },
	}

	finished := false
	step := 0
	for !finished {
		flashed := make(map[Point]bool)

		// First all octopodes increment
		for k, v := range cavern {
			cavern[k] = v + 1
		}

		// Then process flashes until there are no more
		for {
			done := true
			for k, v := range cavern {
				if v <= 9 || flashed[k] {
					// Just save an indent level
					continue
				}

				// There was a flash, so we aren't done yet
				done = false
				flashed[k] = true

				for _, d := range adjacent_dirs {
					adjacent := Point{
						X: k.X + d[0],
						Y: k.Y + d[1],
					}

					if v, exists := cavern[adjacent]; exists {
						// Just increment, any flash will be
						// processed in a later iteration
						cavern[adjacent] = v + 1
					}
				}
			}

			if done {
				break
			}
		}

		// Then all that flashed go to 0
		for k, _ := range flashed {
			cavern[k] = 0
		}

		flashes += len(flashed)
		step++

		if (step == nsteps) {
			fmt.Println("Part 1:", flashes)
		}

		if len(flashed) == len(cavern) {
			fmt.Println("Part 2:", step)
			finished = true
		}
	}

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
