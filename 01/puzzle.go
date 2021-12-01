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
	last := -1
	deeper := 0

	if err := doLines(os.Args[1], func(line string) error {
		this, err := strconv.Atoi(line)
		if err != nil {
			return err
		}

		if last > 0 && this > last {
			deeper++
		}
		last = this

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
