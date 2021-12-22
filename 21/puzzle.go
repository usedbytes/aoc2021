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

func Part2Approach1(initialPositions [2]int) int64 {
	// Calculate all the possible scores for the Dirac dice
	// and the number of ways to reach that score
	dieScores := make([]int, 10)
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			for k := 1; k <= 3; k++ {
				dieScores[i + j + k]++
			}
		}
	}

	// We can actually feasibly covert the *whole* game space,
	// that is all possible combinations of P1 scores, P2 scores
	// and board positions.
	// Store that in one massive array indexed by:
	// [P1 score][P2 score][P1 position][P2 position]
	// Note that to avoid any off-by-1 confusion, I store 11 positions,
	// position 0 is never used.
	states := [32][32][11][11]int64{}

	// Initial state - one way to reach it
	states[0][0][initialPositions[0]][initialPositions[1]] = 1

	scores := [2]int{}
	pos := [2]int{}

	for scores[0] = 0; scores[0] < 21; scores[0]++ {
		for scores[1] = 0; scores[1] < 21; scores[1]++ {
			for pos[0] = 1; pos[0] <= 10; pos[0]++ {
				for pos[1] = 1; pos[1] <= 10; pos[1]++ {
					n := states[scores[0]][scores[1]][pos[0]][pos[1]]
					if n == 0 {
						// No way to reach here
						continue
					}

					for move1, ntimes1 := range dieScores {
						newPos1 := ((pos[0] - 1 + move1) % 10) + 1
						newScore1 := scores[0] + newPos1

						// Does player 2 get a go?
						if newScore1 < 21 {
							// If so, calculate all the outcomes
							for move2, ntimes2 := range dieScores {
								newPos2 := ((pos[1] - 1 + move2) % 10) + 1
								newScore2 := scores[1] + newPos2
								states[newScore1][newScore2][newPos1][newPos2] += n * int64(ntimes1) * int64(ntimes2)
							}
						} else {
							// Otherwise, only take into account
							// player 1's outcomes
							states[newScore1][scores[1]][newPos1][pos[1]] += n * int64(ntimes1)
						}
					}
				}
			}
		}
	}

	// Now we've exhaustively computed the whole game, figure out how
	// many times P1 won first (that is, P1's score was >= 21, while P2 was
	// <= 20
	p1Wins := int64(0)
	for p1Score := 21; p1Score < len(states); p1Score++ {
		for p2Score := 21; p2Score >= 0; p2Score-- {
			for p1Pos := 0; p1Pos < len(states[0][0]); p1Pos++ {
				for p2Pos := 0; p2Pos < len(states[0][0][0]); p2Pos++ {
					p1Wins += states[p1Score][p2Score][p1Pos][p2Pos]
				}
			}
		}
	}

	// And the same for P2
	p2Wins := int64(0)
	for p2Score := 21; p2Score < len(states[0]); p2Score++ {
		for p1Score := 21; p1Score >= 0; p1Score-- {
			for p1Pos := 0; p1Pos < len(states[0][0]); p1Pos++ {
				for p2Pos := 0; p2Pos < len(states[0][0][0]); p2Pos++ {
					p2Wins += states[p1Score][p2Score][p1Pos][p2Pos]
				}
			}
		}
	}

	if p1Wins > p2Wins {
		return p1Wins
	} else {
		return p2Wins
	}
}

func run() error {
	var initialPositions [2]int
	if err := doLines(os.Args[1], func(line string) error {
		var p, pos int
		_, err := fmt.Sscanf(line, "Player %d starting position: %d", &p, &pos)
		if err != nil {
			return err
		}

		initialPositions[p-1] = pos

		return nil
	}); err != nil {
		return err
	}

	nrolls := 0
	roll := func() int {
		nrolls++
		ret := ((nrolls - 1) % 100) + 1;
		return ret
	}

	var scores [2]int
	positions := [2]int{
		initialPositions[0],
		initialPositions[1],
	}
	winner := 0
	for scores[0] < 1000 && scores[1] < 1000 {
		for p := 0; p < len(positions); p++ {
			for i := 0; i < 3; i++ {
				move := roll()
				positions[p] = (((positions[p] - 1) + move) % 10) + 1
			}
			scores[p] += positions[p]

			if scores[p] >= 1000 {
				winner = p
				break
			}
		}
	}

	fmt.Println("Part 1:", scores[(winner + 1) % 2] * nrolls)

	part2 := Part2Approach1(initialPositions)
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
