package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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

type Point struct {
	X, Y int
}

func makeHorizontalFold(y int) func (Point) Point {
	return func(in Point) Point {
		if in.Y == y {
			panic("dots shouldn't lie on folds")
		}

		if in.Y < y {
			return in
		}

		return Point{in.X, y - (in.Y - y)}
	}
}

func makeVerticalFold(x int) func (Point) Point {
	return func(in Point) Point {
		if in.X == x {
			panic("dots shouldn't lie on folds")
		}

		if in.X < x {
			return in
		}

		return Point{x - (in.X - x), in.Y}
	}
}

func wrapXformFunc(current, next func(Point) Point) func(Point) Point {
	return func (in Point) Point {
		return next(current(in))
	}
}

type Instruction struct {
	Direction string
	Line int
}

func run() error {
	coords := []Point{}
	instructions := []Instruction{}

	reachedInstructions := false
	if err := doLines(os.Args[1], func(line string) error {
		if len(line) == 0 {
			reachedInstructions = true
			return nil
		}

		if !reachedInstructions {
			split := strings.Split(line, ",")
			x, err := strconv.Atoi(split[0])
			if err != nil {
				return err
			}
			y, err := strconv.Atoi(split[1])
			if err != nil {
				return err
			}

			coords = append(coords, Point{x, y})
		} else {
			split := strings.Split(line, "=")
			if !strings.HasPrefix(split[0], "fold along ") {
				return fmt.Errorf("unknown instruction %s", line)
			}

			v, err := strconv.Atoi(split[1])
			if err != nil {
				return err
			}

			dir := ""
			if split[0][len(split[0])-1] == 'x' {
				dir = "vertical"
			} else if split[0][len(split[0])-1] == 'y' {
				dir = "horizontal"
			} else {
				return fmt.Errorf("unknown dir %c", split[0][len(split[0])-1])
			}

			i := Instruction{
				Direction: dir,
				Line: v,
			}

			instructions = append(instructions, i)
		}

		return nil
	}); err != nil {
		return err
	}

	f := func(in Point) Point {
		return in
	}

	if instructions[0].Direction == "horizontal" {
		f = wrapXformFunc(f, makeHorizontalFold(instructions[0].Line))
	} else {
		f = wrapXformFunc(f, makeVerticalFold(instructions[0].Line))
	}

	paper := make(map[Point]bool)
	for _, c := range coords {
		paper[f(c)] = true
	}

	fmt.Println("Part 1:", len(paper))

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
