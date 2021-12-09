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

func get(heightMap [][]int, x, y, oob int) int {
	if y < 0 || y >= len(heightMap) {
		return oob
	}

	if x < 0 || x >= len(heightMap[0]) {
		return oob
	}

	return heightMap[y][x]
}

func run() error {

	var heightMap [][]int

	if err := doLines(os.Args[1], func(line string) error {
		heightRow := []int{}
		for _, l := range line {
			n, err := strconv.Atoi(string(l))
			if err != nil {
				return err
			}
			heightRow = append(heightRow, n)
		}
		heightMap = append(heightMap, heightRow)

		return nil
	}); err != nil {
		return err
	}

	lowPoints := []int{}

	for y := 0; y < len(heightMap); y++ {
		for x := 0; x < len(heightMap[0]); x++ {
			here := get(heightMap, x, y, 10)

			if get(heightMap, x - 1, y, 10) > here &&
			   get(heightMap, x + 1, y, 10) > here &&
			   get(heightMap, x, y - 1, 10) > here &&
			   get(heightMap, x, y + 1, 10) > here {
				lowPoints = append(lowPoints, here)
			}
		}
	}

	part1 := 0
	for _, p := range lowPoints {
		part1 += (p + 1)
	}

	fmt.Println(part1)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
