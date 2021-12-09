package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"sort"
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

func searchAround(heightMap [][]int, pt Point, basin map[Point]bool) {
	if _, ok := basin[pt]; ok {
		// Already hit
		return
	}

	if get(heightMap, pt.X, pt.Y, 10) >= 9 {
		// Edge of the basin
		basin[pt] = false
		return
	}

	basin[pt] = true

	dirs := [][2]int {
		{  0, -1 },
		{  0,  1 },
		{ -1,  0 },
		{  1,  0 },
	}

	for _, d := range dirs {
		searchAround(heightMap, Point{pt.X + d[0], pt.Y + d[1]}, basin)
	}
}

func basinSize(basin map[Point]bool) int {
	size := 0
	for _, v := range basin {
		if v {
			size++
		}
	}

	return size
}

type Point struct {
	X, Y int
}

type LowPoint struct {
	Point
	Height int
	BasinSize int
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

	lowPoints := []*LowPoint{}

	for y := 0; y < len(heightMap); y++ {
		for x := 0; x < len(heightMap[0]); x++ {
			here := get(heightMap, x, y, 10)

			if get(heightMap, x - 1, y, 10) > here &&
			   get(heightMap, x + 1, y, 10) > here &&
			   get(heightMap, x, y - 1, 10) > here &&
			   get(heightMap, x, y + 1, 10) > here {
				pt := LowPoint{
					Point: Point {
						X: x,
						Y: y,
					},
					Height: here,
				}
				lowPoints = append(lowPoints, &pt)
			}
		}
	}

	part1 := 0

	basinSizes := []int{}
	for _, p := range lowPoints {
		part1 += (p.Height + 1)

		basin := make(map[Point]bool)
		searchAround(heightMap, p.Point, basin)
		basinSizes = append(basinSizes, basinSize(basin))
	}


	fmt.Println("Part 1:", part1)

	sort.Ints(basinSizes)
	part2 := 1
	for i := len(basinSizes) - 1; i >= len(basinSizes) - 3; i-- {
		part2 *= basinSizes[i]
	}
	fmt.Println("Part 2:", part2)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
