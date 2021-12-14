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

// Mutates a
func add(a *[26]uint64, b [26]uint64) {
	for i, v := range b {
		a[i] += v
	}
}

type Tag struct {
	Pair string
	N int
}

func expandPair(rules map[string]byte, dp map[Tag][26]uint64, pair string, n int) [26]uint64 {
	tag := Tag{ pair, n }
	if res, ok := dp[tag]; ok {
		return res
	}

	res := [26]uint64{}

	if n == 0 {
		return res
	}

	insert := rules[pair]
	res[insert - 'A']++

	add(&res, expandPair(rules, dp, string([]byte{pair[0], insert}), n - 1))
	add(&res, expandPair(rules, dp, string([]byte{insert, pair[1]}), n - 1))

	dp[tag] = res

	return res
}

// Returns count of letters inserted
func expand(rules map[string]byte, template string, n int) [26]uint64 {
	res := [26]uint64{}
	for _, l := range template {
		res[l - 'A']++
	}

	dp := make(map[Tag][26]uint64)

	for i := 0; i < len(template) - 1; i++ {
		pair := template[i:i+2]

		add(&res, expandPair(rules, dp, pair, n))
	}

	return res
}

func min(l []uint64) uint64 {
	min := ^uint64(0)

	for _, v := range l {
		if v != 0 && v < min {
			min = v
		}
	}

	return min
}

func max(l []uint64) uint64 {
	max := uint64(0)

	for _, v := range l {
		if v > max {
			max = v
		}
	}

	return max
}

func run() error {
	template := ""
	rules := make(map[string]byte)

	if err := doLines(os.Args[1], func(line string) error {
		if len(line) == 0 {
			return nil
		}

		if len(template) == 0 {
			template = line
			return nil
		}

		parts := strings.Split(line, " -> ")
		rules[parts[0]] = parts[1][0]
		return nil
	}); err != nil {
		return err
	}

	res := expand(rules, template, 10)
	fmt.Println("Part 1:", max(res[:])-min(res[:]))

	res = expand(rules, template, 40)
	fmt.Println("Part 2:", max(res[:])-min(res[:]))

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
