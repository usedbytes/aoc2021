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

func run() error {
	var lines []string

	if err := doLines(os.Args[1], func(line string) error {
		lines = append(lines, line)
		return nil
	}); err != nil {
		return err
	}

	gammaCounts := countOnes(lines)
	gammaRate := ""

	for _, count := range gammaCounts {
		zeroes := len(lines) - count
		ones := count
		if ones >= zeroes {
			gammaRate += "1"
		} else {
			gammaRate += "0"
		}
	}

	i64, err := strconv.ParseUint(gammaRate, 2, 32)
	if err != nil {
		return err
	}
	gamma := int(i64)
	epsilon := (^gamma) & ((1 << len(gammaRate)) - 1)
	fmt.Println("Power consumption:", gamma * epsilon)

	oxygenPrefix := ""
	oxygenList := make([]string, len(lines))
	copy(oxygenList, lines)

	for len(oxygenList) > 1 {
		for i := 0; i < len(lines[0]); i++ {
			oxygenCounts := countOnes(oxygenList)
			count := oxygenCounts[i]
			zeroes := len(oxygenList) - count
			ones := count

			if ones >= zeroes {
				oxygenPrefix += "1"
			} else {
				oxygenPrefix += "0"
			}

			oxygenList = filter(oxygenList, oxygenPrefix)
			if len(oxygenList) <= 1 {
				break
			}
		}
	}

	i64, err = strconv.ParseUint(oxygenList[0], 2, 32)
	if err != nil {
		return err
	}
	oxygenRating := int(i64)

	co2Prefix := ""
	co2List := make([]string, len(lines))
	copy(co2List, lines)

	for len(co2List) > 1 {
		for i := 0; i < len(lines[0]); i++ {
			co2Counts := countOnes(co2List)
			count := co2Counts[i]
			zeroes := len(co2List) - count
			ones := count

			if ones >= zeroes {
				co2Prefix += "0"
			} else {
				co2Prefix += "1"
			}

			co2List = filter(co2List, co2Prefix)
			if len(co2List) <= 1 {
				break
			}
		}
	}

	i64, err = strconv.ParseUint(co2List[0], 2, 32)
	if err != nil {
		return err
	}
	co2Rating := int(i64)

	fmt.Println(oxygenList)
	fmt.Println(co2List)

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
