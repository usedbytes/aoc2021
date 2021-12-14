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
func add(a *[26]int, b [26]int) {
	for i, v := range b {
		a[i] += v
	}
}

func expandPair(rules map[string]byte, pair string, n int) [26]int {
	res := [26]int{}

	if n == 0 {
		return res
	}

	insert := rules[pair]
	res[insert - 'A']++

	add(&res, expandPair(rules, string([]byte{pair[0], insert}), n - 1))
	add(&res, expandPair(rules, string([]byte{insert, pair[1]}), n - 1))

	return res
}

// Returns count of letters inserted
func expand(rules map[string]byte, template string, n int) [26]int {
	res := [26]int{}
	for _, l := range template {
		res[l - 'A']++
	}

	for i := 0; i < len(template) - 1; i++ {
		pair := template[i:i+2]

		add(&res, expandPair(rules, pair, n))
	}

	return res
}

func min(l []int) int {
	min := 0x7fffffff

	for _, v := range l {
		if v != 0 && v < min {
			min = v
		}
	}

	return min
}

func max(l []int) int {
	max := 0

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

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
