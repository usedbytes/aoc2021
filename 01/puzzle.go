package main

import (
	"bufio"
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

func run() error {
	var err error

	lastSum := -1
	deeper := 0

	windowSize := 1
	if len(os.Args) > 2 {
		windowSize, err = strconv.Atoi(os.Args[2])
		if err != nil {
			return err
		}
	}

	window := make([]int, windowSize)
	idx := 0

	if err := doLines(os.Args[1], func(line string) error {
		this, err := strconv.Atoi(line)
		if err != nil {
			return err
		}

		window[idx % windowSize] = this
		idx++
		if idx < windowSize {
			return nil
		}

		thisSum := 0
		for _, v := range window {
			thisSum += v
		}

		if lastSum > 0 && thisSum > lastSum {
			deeper++
		}
		lastSum = thisSum

		return nil
	}); err != nil {
		return err
	}

	fmt.Println(deeper)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
