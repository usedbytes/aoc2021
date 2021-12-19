package main

import (
	"bufio"
	"fmt"
	"os"
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

var RotFuncs []func(Point) Point

func calcRotFuncs() {
	var transforms = [3][]func(Point) Point {
		{
			// Rotate 90 CCW around X
			func(p Point) Point { return Point{  p.X,  p.Y,  p.Z } },
			func(p Point) Point { return Point{  p.X, -p.Z,  p.Y } },
			func(p Point) Point { return Point{  p.X, -p.Y, -p.Z } },
			func(p Point) Point { return Point{  p.X,  p.Z, -p.Y } },
		},

		{
			// Rotate 90 CCW around Y
			func(p Point) Point { return Point{  p.X,  p.Y,  p.Z } },
			func(p Point) Point { return Point{ -p.Z,  p.Y,  p.X } },
			func(p Point) Point { return Point{ -p.X,  p.Y, -p.Z } },
			func(p Point) Point { return Point{  p.Z,  p.Y, -p.X } },
		},

		{
			// Rotate 90 CCW around Z
			func(p Point) Point { return Point{  p.X,  p.Y,  p.Z } },
			func(p Point) Point { return Point{ -p.Y,  p.X,  p.Z } },
			func(p Point) Point { return Point{ -p.X, -p.Y,  p.Z } },
			func(p Point) Point { return Point{  p.Y, -p.X,  p.Z } },
		},
	}

	makeRotFunc := func(rx, ry, rz int) func(Point) Point {
		return func (p Point) Point {
			return transforms[0][rx](transforms[1][ry](transforms[2][rz](p)))
		}
	}

	oris := make(map[Point]bool)

	p := Point{ 1, 2, 3 }
	for rx := 0; rx < 4; rx++ {
		for ry := 0; ry < 4; ry++ {
			for rz := 0; rz < 4; rz++ {
				f := makeRotFunc(rx, ry, rz)

				q := f(p)

				if _, ok := oris[q]; !ok {
					oris[q] = true
					RotFuncs = append(RotFuncs, f)
				}
			}
		}
	}
}

type Point struct {
	X, Y, Z int
}

func (p Point) Sub(b Point) Point {
	return Point{p.X - b.X, p.Y - b.Y, p.Z - b.Z}
}

func (p Point) Add(b Point) Point {
	return Point{p.X + b.X, p.Y + b.Y, p.Z + b.Z}
}

type Scanner struct {
	Idx int
	Beacons []Point
	Coordinates Point
}

func (s *Scanner) ConvertToAbs(f func(Point) Point) {
	s.Coordinates = f(Point{0, 0, 0})
	for i, b := range s.Beacons {
		s.Beacons[i] = f(b)
	}
}

// Match 's' to 't' transformed by 'r'. If they match, the transformation
// function from 't' coordinates to 's' coordinates is returned.
func match(s, t *Scanner, r int) (int, func(Point) Point) {
	f := RotFuncs[r]

	maxMatch := 0
	for _, b := range s.Beacons {
		for _, c := range t.Beacons {
			// For each pair of 'b' and 'c', we assume that they
			// are the same beacon, and therefore in the same place
			// in absolute coordinates.
			// We find  the transformation from 's' to 't' which would
			// make that so.
			//
			// Then for each other beacon in 't', we transform it
			// by the same transformation, and see how many beacons
			// in 's' it matches.
			// If more than 12 match, then we say it's good, and
			// return the transformation between 's' and 't'
			c = f(c)
			scannerTtoS := b.Sub(c)
			abs := func(p Point) Point {
				p = f(p)
				return p.Add(scannerTtoS)
			}

			absT := make(map[Point]bool)
			for _, d := range t.Beacons {
				absT[abs(d)] = true
			}

			match := 0
			for _, e := range s.Beacons {
				if _, ok := absT[e]; ok {
					match++
				}
			}
			if match > maxMatch {
				maxMatch = match
			}
			if match >= 12 {
				return match, abs
			}
		}
	}

	return maxMatch, nil
}

func run() error {
	var scanners []*Scanner
	var scanner *Scanner

	if err := doLines(os.Args[1], func(line string) error {
		if len(line) == 0 {
			return nil
		}

		if strings.HasPrefix(line, "---") {
			var v int
			_, err := fmt.Sscanf(line, "--- scanner %d ---", &v)
			if err != nil {
				return err
			}

			scanner = &Scanner{
				Idx: v,
			}
			scanners = append(scanners, scanner)
			return nil
		}

		var x, y, z int
		_, err := fmt.Sscanf(line, "%d,%d,%d", &x, &y, &z)
		if err != nil {
			return err
		}
		scanner.Beacons = append(scanner.Beacons, Point{x, y, z})

		return nil
	}); err != nil {
		return err
	}

	// Track which scanners we've found
	foundScanners := make(map[int]bool)

	// We will assume scanners[0] is at (0, 0, 0) with rotation '0'
	foundScanners[0] = true

	// This only works if there's no ambiguity. If there are two scanners
	// which look like they overlap, but they actually don't, then this
	// falls over.
	// I assume the input is crafted such that this never happens.
	// I think to handle that, after finding a scanner you could check if
	// it's consistent with all the other one's you've found already and
	// keep track of permutations you've tried.
	for len(foundScanners) < len(scanners) {
		for i, _ := range scanners {
			// We already know where scanner[i] is, don't look for it again
			if _, ok := foundScanners[i]; ok {
				continue
			}

			// Compare scanners[i] with all of the ones we already
			// found, trying to find a match.
			// Note: This might not succeed if scanners[i] doesn't
			// overlap any of foundScanners yet.
			for k, _ := range foundScanners {
				s := scanners[k]
				t := scanners[i]

				// Try and find s relative to t
				for r := 0; r < len(RotFuncs); r++ {
					n, f := match(s, t, r)
					if n >= 12 {
						// Convert 't' to absolute coordinates, so we can
						// use it as a future reference
						t.ConvertToAbs(f)
						foundScanners[i] = true
						break
					}
				}

				// Found a match, stop looking
				if _, ok := foundScanners[i]; ok {
					fmt.Println("Scanner", i, "at", t.Coordinates)
					break
				}
			}
		}
	}

	// Now just build a map of all the beacon absolute coordinates
	beacons := make(map[Point]bool)
	for _, s := range scanners {
		for _, b := range s.Beacons {
			beacons[b] = true
		}
	}

	fmt.Println("Part 1:", len(beacons))

	return nil
}

func init() {
	calcRotFuncs()
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
