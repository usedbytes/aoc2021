package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

var properDigits []string = []string{
	"abcefg",
	"cf",
	"acdeg",
	"acdfg",
	"bcdf",
	"abdfg",
	"abdefg",
	"acf",
	"abcdefg",
	"abcdfg",
};

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

func countCommonLetters(a, b string) int {
	count := 0
	for _, l := range a {
		if strings.ContainsRune(b, l) {
			count++
		}
	}

	return count
}

func orderString(s string) string {
	b := []byte(s)
	sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
	return string(b)
}

func run() error {
	// Counts the number of occurrences of specific string lengths
	part1 := make([]int, 8)

	// Accumulates all the displays
	part2 := 0

	if err := doLines(os.Args[1], func(line string) error {
		// Parse the input
		chunks := strings.Split(line, " | ")
		signals := strings.Split(chunks[0], " ")
		digits := strings.Split(chunks[1], " ")

		for i, d := range digits {
			d = orderString(strings.TrimSpace(d))
			digits[i] = d
			part1[len(d)]++
		}

		for i, s := range signals {
			signals[i] = orderString(strings.TrimSpace(s))
		}

		// Build a map from number to possible candidate strings for
		// that number, based on the number of segments in the string
		candidateMap := make(map[int][]string)
		for _, s := range signals {
			for i, d := range properDigits {
				if len(s) == len(d) {
					candidateMap[i] = append(candidateMap[i], s)
				}
			}
		}

		// Map from string to number once known
		assignedMap := make(map[string]int)

		// Assign anything we know for sure (i.e. only has one candidate)
		// Should be 1, 4, 7 and 8 after the 1st pass
		for n, candidates := range candidateMap {
			if len(candidates) == 1 {
				assignedMap[candidates[0]] = n
				delete(candidateMap, n)
			}
		}

		// While we still have values to find
		remaining := len(candidateMap)
		for remaining > 0 {
			for n, candidates := range candidateMap {
				// If 'assigned' and 'candidate' don't share the same
				// number of common letters as 'm' and 'n' share segments,
				// then 'candidate' is *not* a valid possibility for 'n'
				// and can be eliminated.
				//
				// Note: The candidate *must* be consistent with *all*
				// already-assigned values (in assignedMap)
				//
				// We'll build a new list of candidates which
				// weren't eliminated.
				newCandidates := []string{}
				for _, candidate := range candidates {
					eliminated := false
					for assigned, m := range assignedMap {
						// We only check that the number of segments matches
						// This is hopefully enough to remove all ambiguity.
						// We never actually determine which letters represent
						// which segments.
						expected := countCommonLetters(properDigits[n], properDigits[m])
						got := countCommonLetters(assigned, candidate)
						if expected != got {
							eliminated = true
							break
						}
					}
					if !eliminated {
						// This candidate was consistent with *all* currently-assigned
						// numbers, so it gets to stay in the list
						newCandidates = append(newCandidates, candidate)
					}
				}

				if len(newCandidates) == 1 {
					// Hooray, we found one!
					assignedMap[newCandidates[0]] = n
					delete(candidateMap, n)
				} else {
					// Add the (maybe shorter) candidate list back
					// to the map for the next iteration
					candidateMap[n] = newCandidates
				}
			}

			if len(candidateMap) == remaining {
				return fmt.Errorf("couldn't assign any numbers, will loop forever")
			}
			remaining = len(candidateMap)
		}

		// Phew! We found all the digits, so what's being displayed?
		display := 0
		for _, d := range digits {
			display *= 10
			display += assignedMap[d]
		}

		part2 += display

		return nil
	}); err != nil {
		return err
	}

	fmt.Println("Part1:", part1[2] + part1[4] + part1[3] + part1[7])
	fmt.Println("Part2:", part2)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
