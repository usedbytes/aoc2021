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
	var oneCounts []int
	numBits := 0
	numLines := 0

	if err := doLines(os.Args[1], func(line string) error {
		if numBits == 0 {
			numBits = len(line)
			oneCounts = make([]int, numBits)
		}

		for i, c := range line {
			if c == '1' {
				oneCounts[i]++
			}
		}

		numLines++

		return nil
	}); err != nil {
		return err
	}

	output := 0
	for i, count := range oneCounts {
		if count > numLines / 2 {
			output |= (1 << (numBits - i - 1))
		}
	}

	gamma := output
	epsilon := (^output) & ((1 << numBits) - 1)

	fmt.Println(gamma * epsilon)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
