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

func costAbsDiff(a, b int) int {
	return abs(a - b)
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

	minDiff := 0x7fffffff

	for p := minPos; p <= maxPos; p++ {
		diff := sumOfCosts(positions, p, costAbsDiff)
		if diff < minDiff {
			minDiff = diff
		}
	}
	fmt.Println(minDiff)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
