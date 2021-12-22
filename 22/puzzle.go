package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
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
	On bool
}

func (c Cuboid) Intersect(b Cuboid) Cuboid {
	return Cuboid{
		X: c.X.Intersect(b.X),
		Y: c.Y.Intersect(b.Y),
		Z: c.Z.Intersect(b.Z),
	}
}

func (c Cuboid) String() string {
	return fmt.Sprintf("{%v %v %v %v %d}", c.X, c.Y, c.Z, c.On, c.Count())
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
				nc := Cuboid{ x, y, z, c.On }
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
	if c.X.Length < 0 || c.Y.Length < 0 || c.Z.Length < 0 {
		return 0
	}
	return c.X.Length * c.Y.Length * c.Z.Length
}

// Return the number of eventual 'on' cells contributed by *only* 'c'
// i.e. that aren't masked or turned off by a another entry in 'through'
func propagate(c Cuboid, through []Cuboid) int64 {
	if !c.On {
		// Can't ever turn anything on
		return 0
	}

	// If we made it to the end, then return what's left
	if len(through) == 0 {
		return int64(c.Count())
	}

	// Anything that intersects with 'next' will be handled by a later
	// stage.
	// If that's an "off" then those cells just get dropped, if it's an
	// 'on' then that command's own 'propagate' will handle it

	// Therefore, we just propagate the Disjoint parts of 'c'
	next := through[0]
	djs := c.Disjoint(next)
	count := int64(0)
	for _, d := range djs {
		count += propagate(d, through[1:])
	}

	return count
}

func run() error {
	part1Range := Cuboid{
		X: MakeRange( -50, 50 ),
		Y: MakeRange( -50, 50 ),
		Z: MakeRange( -50, 50 ),
	}

	// p1Cmds is filtered by the range
	var p1Cmds, p2Cmds []Cuboid

	if err := doLines(os.Args[1], func(line string) error {
		fmt.Println(line)
		var s string
		var x1, x2, y1, y2, z1, z2 int
		_ ,err := fmt.Sscanf(line, "%s x=%d..%d,y=%d..%d,z=%d..%d", &s, &x1, &x2, &y1, &y2, &z1, &z2)
		if err != nil {
			return err
		}

		c := Cuboid{
			X: MakeRange(x1, x2),
			Y: MakeRange(y1, y2),
			Z: MakeRange(z1, z2),
		}

		if s == "on" {
			c.On = true
		}

		if c.Intersect(part1Range).Count() > 0 {
			p1Cmds = append(p1Cmds, c)
		}

		p2Cmds = append(p2Cmds, c)

		return nil
	}); err != nil {
		return err
	}

	part1 := int64(0)
	for i, cmd := range p1Cmds {
		this := propagate(cmd, p1Cmds[i+1:])
		part1 += this
	}

	// Throw more cores at the problem... This clearly isn't the "right"
	// solution, it takes ~15 minutes on my M1 Mac
	var wg sync.WaitGroup
	counts := make(chan int64)

	part2 := int64(0)
	for i, cmd := range p2Cmds {
		wg.Add(1)
		go func(c Cuboid, i int) {
			this := propagate(c, p2Cmds[i+1:])
			fmt.Printf("Cmd %d/%d adds %d\n", i, len(p2Cmds), this)
			counts <- this
			wg.Done()
		}(cmd, i)
	}

	go func() {
		wg.Wait()
		close(counts)
	}()

	for c := range counts {
		part2 += c
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
