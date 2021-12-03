package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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

func countOnes(lines []string) []int {
	counts := make([]int, len(lines[0]))

	for _, line := range lines {
		for i, c := range line {
			if c == '1' {
				counts[i]++
			}
		}
	}

	return counts
}

func filter(in []string, prefix string) []string {
	var out []string
	for _, s := range in {
		if strings.HasPrefix(s, prefix) {
			out = append(out, s)
		}
	}

	return out
}

func decode(lines []string, rule func(current string, ones, zeroes int) string, filter func(lines []string, prefix string) []string) string {
	result := ""

	counts := countOnes(lines)
	for i := 0; i < len(lines[0]); i++ {
		count := counts[i]
		zeroes := len(lines) - count
		ones := count

		result = rule(result, ones, zeroes)

		if filter != nil {
			lines = filter(lines, result)
			if len(lines) <= 1 {
				return lines[0]
			}
			counts = countOnes(lines)
		}
	}

	return result
}

func run() error {
	var lines []string

	if err := doLines(os.Args[1], func(line string) error {
		lines = append(lines, line)
		return nil
	}); err != nil {
		return err
	}

	mostCommonRule := func(current string, ones, zeroes int) string {
		if ones >= zeroes {
			return current + "1"
		} else {
			return current + "0"
		}
	}

	leastCommonRule := func(current string, ones, zeroes int) string {
		if ones >= zeroes {
			return current + "0"
		} else {
			return current + "1"
		}
	}

	gammaRate := decode(lines, mostCommonRule, nil)

	i64, err := strconv.ParseUint(gammaRate, 2, 32)
	if err != nil {
		return err
	}
	gamma := int(i64)
	epsilon := (^gamma) & ((1 << len(gammaRate)) - 1)
	fmt.Println("Power consumption:", gamma * epsilon)

	oxygen := decode(lines, mostCommonRule, filter)
	i64, err = strconv.ParseUint(oxygen, 2, 32)
	if err != nil {
		return err
	}
	oxygenRating := int(i64)

	co2 := decode(lines, leastCommonRule, filter)
	i64, err = strconv.ParseUint(co2, 2, 32)
	if err != nil {
		return err
	}
	co2Rating := int(i64)

	fmt.Println("Life Support Rating:", oxygenRating * co2Rating)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
