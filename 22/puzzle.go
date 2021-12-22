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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type Range struct {
	Start  int
	Length int
}

func MakeRange(start, end int) Range {
	return Range{ start, end - start + 1 }
}

func (r Range) Intersect(b Range) Range {
	nmin := max(r.Start, b.Start)
	nmax := min(r.Start + r.Length, b.Start + b.Length)

	return Range{ nmin, max(0, nmax - nmin) }
}

func (r Range) disjoint(b Range) []Range {
	minMin := min(r.Start, b.Start)
	maxMax := max(r.Start + r.Length, b.Start + b.Length)

	minMax := min(r.Start + r.Length, b.Start + b.Length)
	maxMin := max(r.Start, b.Start)

	r1Min := minMin
	r1Max := min(minMax, maxMin)

	r2Min := max(maxMin, minMax)
	r2Max := maxMax

	res := []Range{}
	if r1Max - r1Min > 0 {
		res = append(res, Range{ r1Min, r1Max - r1Min})
	}

	if r2Max - r2Min > 0 {
		res = append(res, Range{ r2Min, r2Max - r2Min})
	}

	return res
}

func (r Range) Disjoint(b Range) []Range {
	djs := r.disjoint(b)

	res := []Range{}
	for _, d := range djs {
		if r.Intersect(d).Length > 0 {
			res = append(res, d)
		}
	}
	return res
}

func (r Range) String() string {
	return fmt.Sprintf("%d..%d", r.Start, r.Start + r.Length - 1)
}

type Cuboid struct {
	X, Y, Z Range
}

func (c Cuboid) Intersect(b Cuboid) Cuboid {
	return Cuboid{
		X: c.X.Intersect(b.X),
		Y: c.Y.Intersect(b.Y),
		Z: c.Z.Intersect(b.Z),
	}
}

func (c Cuboid) Disjoint(b Cuboid) []Cuboid {
	disx := c.X.Disjoint(b.X)
	disy := c.Y.Disjoint(b.Y)
	disz := c.Z.Disjoint(b.Z)

	intx := c.X.Intersect(b.X)
	inty := c.Y.Intersect(b.Y)
	intz := c.Z.Intersect(b.Z)

	res := []Cuboid{}

	for _, x := range append(disx, intx) {
		for _, y := range append(disy, inty) {
			for _, z := range append(disz, intz) {
				nc := Cuboid{ x, y, z }
				ic := c.Intersect(nc)
				ib := b.Intersect(nc)

				icc := ic.Count()
				ibc := ib.Count()

				if (icc > 0) && !(icc > 0 && ibc > 0) {
					res = append(res, nc)
				}
			}
		}
	}

	return res
}

func (c Cuboid) Count() int {
	return c.X.Length * c.Y.Length * c.Z.Length
}

func run() error {
	//initMin, initMax := -50, 50

	on := []Cuboid{}

	rangeCube := Cuboid{
		X: MakeRange(-50, 50),
		Y: MakeRange(-50, 50),
		Z: MakeRange(-50, 50),
	}

	if err := doLines(os.Args[1], func(line string) error {
		fmt.Println(line)
		var s string
		var x1, x2, y1, y2, z1, z2 int
		_ ,err := fmt.Sscanf(line, "%s x=%d..%d,y=%d..%d,z=%d..%d", &s, &x1, &x2, &y1, &y2, &z1, &z2)
		if err != nil {
			return err
		}

		//fmt.Println(s, x1, x2, y1, y2, z1, z2)

		c := Cuboid{
			X: MakeRange(x1, x2),
			Y: MakeRange(y1, y2),
			Z: MakeRange(z1, z2),
		}

		if c.Intersect(rangeCube).Count() == 0 {
			fmt.Println(line, "out of range")
		}

		intersects := false
		for _, o := range on {
			cio := c.Intersect(o)
			//fmt.Println(c, " N ", o, cio, cio.Count())
			if cio.Count() > 0 {
				intersects = true
				break
			}
		}

		if !intersects {
			if s == "on" {
				fmt.Println(c, "adds", c.Count())
				on = append(on, c)
			}
			// Nothing to do if "off" doesn't intersect anything
		} else {
			// On only contains unique bits, so
			// Find the intersection between 'c' and each 'on'
			newOn := []Cuboid{}
			for _, ic := range on {
				cic := c.Intersect(ic)
				if cic.Count() == 0 {
					// Doesn't change anything in this chunk
					fmt.Println("keeping", ic)
					newOn = append(newOn, ic)
					continue
				}

				//fmt.Println("Blah")

				if s == "on" {
					fmt.Println(ic.Count(), "intersects", c.Count())
					// We keep 'ic'
					newOn = append(newOn, ic)
					// The disjoint parts of 'c' turn on
					cdc := c.Disjoint(ic)
					newOn = append(newOn, cdc...)
					//fmt.Println("fragments")
					for _, dc := range cdc {
						fmt.Println(dc.Count())
					}
				} else {
					// The disjoint parts if 'ic' *stay* on
					cdc := ic.Disjoint(c)
					newOn = append(newOn, cdc...)
					/*
					for _, dc := range cdc {
						fmt.Println(dc, "stays on", dc.Count())
					}
					*/

					//goesOff := c.Intersect(ic)
					//fmt.Println(goesOff, "turns off", goesOff.Count())
				}
			}

			on = newOn
		}

		fmt.Println("len(on)", len(on))

		if len(on) > 5000 {
			return fmt.Errorf("getting out of hand")
		}

		return nil
	}); err != nil {
		return err
	}

	count := 0
	for _, c := range on {
		fmt.Println(c, "turns on", c.Count())
		count += c.Count()
	}

	fmt.Println("count", count)

	/*

	// Explode the "on" list
	modified := true
	for modified {
		modified = false
		for i, ic := range on {
			for _, jc := range on[i+1:] {
				// Do these cubes intersect? If not, there's
				// nothing to do
				iji := ic.Intersect(jc)
				if iji.Count() == 0 {
					continue
				}

				// Discard the overlapping part of 'ic', and
				// keep the non-overlapping part(s)
				ijd := ic.Disjoint(jc)
				if len(ijd) > 0 {
					on[i] = ijd[0]
					on = append(on, ijd[1:]...)
				} else {
					copy(on[i:], on[i+1:])
					on = on[:len(on)-1]
				}
				modified = true
				break
			}

			if modified {
				break
			}
		}
	}

	// Explode the "off" list
	modified = true
	for modified {
		modified = false
		for i, ic := range off {
			for _, jc := range off[i+1:] {
				iji := ic.Intersect(jc)
				if iji.Count() == 0 {
					continue
				}

				ijd := ic.Disjoint(jc)
				if len(ijd) > 0 {
					off[i] = ijd[0]
					off = append(off, ijd[1:]...)
				} else {
					copy(off[i:], off[i+1:])
					off = off[:len(off)-1]
				}
				modified = true
				break
			}

			if modified {
				break
			}
		}
	}

	count := 0
	for _, c := range on {
		fmt.Println(c, "turns on", c.Count())
		count += c.Count()
	}

	fmt.Println("count", count)

	for _, c := range off {
		for _, d := range on {
			i := c.Intersect(d)
			fmt.Println(c, d, "turns off", i.Count())
			count -= i.Count()
		}
	}

	fmt.Println("On:", count)
	*/

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
