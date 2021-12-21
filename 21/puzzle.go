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

func run() error {
	var positions [2]int
	if err := doLines(os.Args[1], func(line string) error {
		var p, pos int
		_, err := fmt.Sscanf(line, "Player %d starting position: %d", &p, &pos)
		if err != nil {
			return err
		}

		positions[p-1] = pos

		return nil
	}); err != nil {
		return err
	}

	nrolls := 0
	roll := func() int {
		nrolls++
		ret := ((nrolls - 1) % 100) + 1;
		return ret
	}

	var scores [2]int
	winner := 0
	for scores[0] < 1000 && scores[1] < 1000 {
		for p := 0; p < len(positions); p++ {
			for i := 0; i < 3; i++ {
				move := roll()
				positions[p] = (((positions[p] - 1) + move) % 10) + 1
			}
			scores[p] += positions[p]

			if scores[p] >= 1000 {
				winner = p
				break
			}
		}
	}

	fmt.Println("Part 1:", scores[(winner + 1) % 2] * nrolls)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
