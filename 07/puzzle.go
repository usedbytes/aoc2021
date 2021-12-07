package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

func sumOfCosts(vals []int, x int, cost func(a, b int) int) int {
	res := 0
	for _, v := range vals {
		res += cost(v, x)
	}

	return res
}

func run() error {
	filename := os.Args[1]

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	rd := bufio.NewReader(f)

	line, err := rd.ReadString('\n')
	if err != nil {
		return err
	}

	line = strings.TrimSpace(line)

	startStrings := strings.Split(line, ",")
	positions := make([]int, len(startStrings))
	for i, s := range startStrings {
		v, err := strconv.Atoi(s)
		if err != nil {
			return err
		}

		positions[i] = v
	}

	maxPos := 0
	minPos := 0x7fffffff
	for _, p := range positions {
		if p < minPos {
			minPos = p
		}
		if p > maxPos {
			maxPos = p
		}
	}

	// Part 1 - Absolute difference
	part1Cost := func(a, b int) int {
		return abs(a - b)
	}

	// Part 2 - Integral
	costMap := make(map[int]int)
	costMap[0] = 0
	var recurseAdd func(int) int
	recurseAdd = func(a int) int {
		v, ok := costMap[a]
		if ok {
			return v
		}

		costMap[a] = a + recurseAdd(a - 1)
		return costMap[a]
	}

	part2Cost := func(a, b int) int {
		diff := abs(a - b)
		return recurseAdd(diff)
	}


	part1Min := 0x7fffffff
	part2Min := 0x7fffffff
	for p := minPos; p <= maxPos; p++ {
		diff := sumOfCosts(positions, p, part1Cost)
		if diff < part1Min {
			part1Min = diff
		}

		diff = sumOfCosts(positions, p, part2Cost)
		if diff < part2Min {
			part2Min = diff
		}
	}
	fmt.Println("Part 1:", part1Min)
	fmt.Println("Part 2:", part2Min)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
