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

	fmt.Println(c, " ^ ", b)
	fmt.Println("xs", disx)
	fmt.Println("ys", disy)
	fmt.Println("zs", disz)

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

	/*
	on := []Cuboid{}
	off := []Cuboid{}

	if err := doLines(os.Args[1], func(line string) error {
		fmt.Println(line)

		var s string
		var x1, x2, y1, y2, z1, z2 int
		_ ,err := fmt.Sscanf(line, "%s x=%d..%d,y=%d..%d,z=%d..%d", &s, &x1, &x2, &y1, &y2, &z1, &z2)
		if err != nil {
			return err
		}

		fmt.Println(s, x1, x2, y1, y2, z1, z2)

		c := Cuboid{
			X: MakeRange(x1, x2),
			Y: MakeRange(y1, y2),
			Z: MakeRange(z1, z2),
		}

		if s == "on" {
			on = append(on, c)
		} else {
			off = append(off, c)
		}

		return nil
	}); err != nil {
		return err
	}

	for i, o := range on {
		for j, m := range on {
			if j <= i {
				continue
			}

			ndjs = append(ndjs, o.Disjoint(m))
		}
	}
	*/

	c := Cuboid{
		X: MakeRange(0, 2),
		Y: MakeRange(0, 2),
		Z: MakeRange(0, 2),
	}

	d := Cuboid{
		X: MakeRange(1, 1),
		Y: MakeRange(1, 1),
		Z: MakeRange(1, 1),
	}

	fmt.Println(c.Z.Disjoint(d.Z))
	fmt.Println(d.Z.Disjoint(c.Z))
	fmt.Println(c.Intersect(d))
	fmt.Println(c.Disjoint(d), len(c.Disjoint(d)))

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
