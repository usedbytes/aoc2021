package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

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

	// Acceleration is -1, so the projectile reaches its peak height
	// at step 'u', where 'u' is the starting velocity.
	// From the rules, it will spend 2 'steps' at its peak, one step
	// with v = 0, and one with v = -1, meaning it returns to y = 0
	// at step (2 * u) + 1, with a velocity of -(u + 1).
	//
	// At this point, we need its velocity to be such that on
	// _the very next step_, (2 * u) + 2, it lands in the target area.
	//  * If it lands short of the target area on step (2 * u) + 2,
	//    then it could have been going faster, and reached a higher
	//    altitude.
	//  * If it overshoots the target area, then it was going too
	//    fast.
	// To maximise the speed, we need it to land _just inside_ the target
	// area, nearly overshooting. That means we need the velocity on
	// step ((2 * u) + 1) (when it reaches y = 0) to be the "lowest"
	// y position of the target region:
	//
	//  -(u + 1) = y1
	//
	// So the starting velocity is -(y1 + 1)
	//
	// We can simulate the trajectory to find the max altitude.
	// For some reason, in s = 0.5(u)t, we need to use (t + 1) to get
	// the right answer. I'm not sure why that is, maybe to do with
	// the rules of motion we are given?

	initialVelocity := -(y1 + 1)

	y := 0
	for v := initialVelocity; v > 0; v-- {
		y = y + v
	}

	fmt.Println("Part 1:", y)

	hits := make(map[[2]int]bool)

	// For Part 2, we know the velocity has to be between
	// y1 and +initialVelocity, so just scan
	// that whole range.
	//
	// For each, determine which steps have 'y' in the target
	// range.
	for uy := y1; uy <= initialVelocity; uy += 1 {
		y = 0
		for vy, ty := uy, 0; ; vy, ty = vy - 1, ty + 1 {
			y = y + vy

			if y > y2 {
				// Not there yet
				continue
			}

			if y < y1 {
				// Overshot
				break
			}

			// I can't figure out how to solve for 'x', so just brute
			// force it :-)
			for ux := 0; ux <= x2; ux++ {
				x := 0
				vx := ux
				for tx := 0; tx <= ty; tx++ {
					x += vx
					if vx > 0 {
						vx--
					}
				}

				if x >= x1 && x <= x2 {
					hits[[2]int{ux, uy}] = true
				}
			}
		}
	}

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
