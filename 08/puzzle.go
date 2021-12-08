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

func run() error {
	// Counts the number of occurrences of specific string lengths
	count := make([]int, 8)

	if err := doLines(os.Args[1], func(line string) error {
		chunks := strings.Split(line, "|")
		digits := strings.Split(chunks[1], " ")

		for _, d := range digits {
			d = strings.TrimSpace(d)
			count[len(d)]++
		}

		return nil
	}); err != nil {
		return err
	}

	fmt.Println(count[2] + count[4] + count[3] + count[7])

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
