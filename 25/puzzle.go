package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/pprof"
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

func moveEast(bed, newBed map[[2]int]byte, coord, size [2]int) bool {
	v, ok := bed[coord]
	if ok && v != '>' {
		newBed[coord] = v
		return false
	}

	newPos := [2]int{(coord[0]+1) % size[0], coord[1]}

	//fmt.Println("Move from", coord, "to", newPos, string(v), string(bed[newPos]))

	if _, occupied := bed[newPos]; !occupied {
		newBed[newPos] = v
		return true
	}

	newBed[coord] = v

	return false
}

func moveSouth(bed, newBed map[[2]int]byte, coord, size [2]int) bool {
	v, ok := bed[coord]
	if ok && v != 'v' {
		newBed[coord] = v
		return false
	}

	newPos := [2]int{coord[0], (coord[1]+1) % size[1]}

	if _, occupied := bed[newPos]; !occupied {
		newBed[newPos] = v
		return true
	}

	newBed[coord] = v

	return false
}

func printBed(bed map[[2]int]byte, size [2]int) {
	for y := 0; y < size[1]; y++ {
		line := ""
		for x := 0; x < size[0]; x++ {
			if v, ok := bed[[2]int{x, y}]; ok {
				line += string(v)
			} else {
				line += "."
			}
		}
		fmt.Println(line)
	}
}

func run() error {
	bed := make(map[[2]int]byte)

	var size [2]int
	y := 0
	if err := doLines(os.Args[1], func(line string) error {
		for x, c := range line {
			if c == '.' {
				continue
			}

			bed[[2]int{x, y}] = byte(c)
		}

		size[0] = len(line)

		y++
		return nil
	}); err != nil {
		return err
	}

	size[1] = y

	step := 0
	canMove := true
	for step = 0 ; canMove; step++ {
		newBed := make(map[[2]int]byte)
		canMove = false

		for k, _ := range bed {
			move := moveEast(bed, newBed, k, size)
			canMove = move || canMove
		}

		bed = newBed
		newBed = make(map[[2]int]byte)

		for k, _ := range bed {
			move := moveSouth(bed, newBed, k, size)
			canMove = move || canMove
		}

		bed = newBed
	}

	fmt.Println("Part 1:", step)

	return nil
}

func main() {
	profileEnv := os.Getenv("PROFILE")
	if profileEnv != "" {
		f, err := os.Create(profileEnv)
		if err != nil {
			fmt.Println("ERROR:", err)
			os.Exit(1)
		}

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
