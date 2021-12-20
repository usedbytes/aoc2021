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

type Point struct {
	X, Y int
}

type Image struct {
	Min, Max Point
	Pixels map[Point]bool
	Oob bool
}

func (i *Image) Get(p Point) bool {
	if p.X < i.Min.X || p.Y < i.Min.Y ||
	   p.X > i.Max.X || p.Y > i.Max.Y {
		return i.Oob
	}

	return i.Pixels[p]
}

func sig(img *Image, p Point) int {
	bitPositions := [][2]int{
		{  1,  1 },
		{  0,  1 },
		{ -1,  1 },
		{  1,  0 },
		{  0,  0 },
		{ -1,  0 },
		{  1, -1 },
		{  0, -1 },
		{ -1, -1 },
	}

	res := 0
	for i, d := range bitPositions {
		np := Point{p.X + d[0], p.Y + d[1]}
		if img.Get(np) {
			res |= (1 << i)
		}
	}

	return res
}

func run() error {
	var algorithm []byte

	pixels := make(map[Point]bool)

	var y int
	var maxX, maxY int
	if err := doLines(os.Args[1], func(line string) error {
		if algorithm == nil {
			if len(line) != 512 {
				panic("algorithm should be 512 bytes")
			}
			algorithm = []byte(line)
			return nil
		}

		if len(line) == 0 {
			return nil
		}

		maxX = len(line)
		for x := 0; x < len(line); x++ {
			if line[x] == '#' {
				pixels[Point{x, y}] = true
			}
		}
		y++
		return nil
	}); err != nil {
		return err
	}
	maxY = y

	img := &Image{
		Min: Point{0, 0},
		Max: Point{maxX - 1, maxY - 1},
		Pixels: pixels,
		Oob: false, // Initially, out-of-bounds areas of the input are '.' (unlit)
	}

	part1 := 0
	part2 := 0

	for i := 0; i < 50; i++ {
		newPixels := make(map[Point]bool)
		for y := img.Min.Y - 1; y <= img.Max.Y + 1; y++ {
			for x := img.Min.X - 1; x <= img.Max.X + 1; x++ {
				idx := sig(img, Point{x, y})
				b := algorithm[idx]
				if b == '#' {
					newPixels[Point{x, y}] = true
				}
			}
		}

		// Just use some random far out-of-bounds point to get Oob
		// Note: This is the trick! algorithm[0] is '#', which means
		// after the first image enhancement pass, all of the "infinite"
		// outer edges of the image become '#', so the second pass needs
		// to take that into account.
		oob := algorithm[sig(img, Point{-10000, -10000})]
		oobb := false
		if oob == '#' {
			oobb = true
		}

		newImg := &Image{
			Min: Point{ img.Min.X - 1, img.Min.Y - 1},
			Max: Point{ img.Max.X + 1, img.Max.Y + 1},
			Pixels: newPixels,
			Oob: oobb,
		}
		img = newImg

		if i == 1 {
			part1 = len(img.Pixels)
		}
	}

	part2 = len(img.Pixels)

	if len(os.Args) > 2 {
		for y := img.Min.Y - 1; y <= img.Max.Y + 1; y++ {
			line := ""
			for x := img.Min.X - 1; x <= img.Max.X + 1; x++ {
				if img.Get(Point{x, y}) {
					line += "#"
				} else {
					line += "."
				}
			}
			fmt.Println(line)
		}
	}

	fmt.Println("Part 1:", part1)
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
