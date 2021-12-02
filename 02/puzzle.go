package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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

type State struct {
	Horizontal int
	Depth      int
	Aim        int
}

func (s *State) Mutate(cmd string, arg int) {
	switch cmd {
	case "forward":
		s.Horizontal += arg
	case "up":
		s.Depth -= arg
	case "down":
		s.Depth += arg
	}
}

func (s *State) MutatePart2(cmd string, arg int) {
	switch cmd {
	case "forward":
		s.Horizontal += arg
		s.Depth += arg * s.Aim
	case "up":
		s.Aim -= arg
	case "down":
		s.Aim += arg
	}
}

func run() error {

	var s State

	if err := doLines(os.Args[1], func(line string) error {
		parts := strings.Split(line, " ")

		arg, err := strconv.Atoi(parts[1])
		if err != nil {
			return err
		}

		if len(os.Args) > 2 {
			s.MutatePart2(parts[0], arg)
		} else {
			s.Mutate(parts[0], arg)
		}

		return nil
	}); err != nil {
		return err
	}

	fmt.Println(s.Horizontal * s.Depth)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
