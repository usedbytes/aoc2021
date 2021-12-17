package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func run() error {

	f, err := os.Open(os.Args[1])
	if err != nil {
		return err
	}
	defer f.Close()

	bs, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	var x1, x2, y1, y2 int
	fmt.Sscanf(string(bs), "target area: x=%d..%d, y=%d..%d", &x1, &x2, &y1, &y2)

	// The Y velocity can't be any more than the furthest edge of the
	// region, otherwise we'll overshoot on the first step after
	// leaving (or arriving back at) y = 0. So this sets our Y bounds.
	maxYVel := max(abs(y1), abs(y2))

	maxYPos := 0
	hits := make(map[[2]int]bool)

	// Just brute force exhaustive search the whole range
	for uy := -maxYVel; uy <= maxYVel; uy += 1 {
		y := 0
		peak := 0
		for vy, ty := uy, 0; ; vy, ty = vy - 1, ty + 1 {
			y = y + vy

			if y > peak {
				peak = y
			}

			if y > y2 {
				// Not there yet
				continue
			}

			if y < y1 {
				// Overshot
				break
			}

			// This is a valid Y velocity, if it has the highest
			// peak, save it.
			// Note: This is expected to happen for uy == maxYVel
			if peak > maxYPos {
				maxYPos = peak
			}

			// Now look for X velocities which will land in the
			// on the same time step.
			// The same as for Y, we can't overshoot on our first
			// step, which sets the upper bound. The problem
			// description says we can only fire "forward", so 0
			// is our lower bound.
			// Note: We could save these results somewhere instead
			// of re-searching each time.
			for ux := 0; ux <= max(x1, x2); ux++ {
				x := 0
				vx := ux
				// Simulate trajectory up to 'ty'
				for tx := 0; tx <= ty; tx++ {
					x += vx
					if vx > 0 {
						vx--
					}
				}

				// And then check if we're inside the target
				if x >= x1 && x <= x2 {
					hits[[2]int{ux, uy}] = true
				}
			}
		}
	}


	fmt.Println("Part 1:", maxYPos)
	fmt.Println("Part 2:", len(hits))

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
