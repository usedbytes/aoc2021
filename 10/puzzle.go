package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

var pairs map[rune]rune = map[rune]rune{
	'(': ')',
	'[': ']',
	'{': '}',
	'<': '>',
}

var syntaxErrorScores map[rune]int = map[rune]int{
	')': 3,
	']': 57,
	'}': 1197,
	'>': 25137,
}

var completionScores map[rune]int = map[rune]int{
	')': 1,
	']': 2,
	'}': 3,
	'>': 4,
}

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

func run() error {

	part1 := 0
	part2 := []int{}

	if err := doLines(os.Args[1], func(line string) error {
		fmt.Println(line)

		var err error
		var illegal rune
		stack := []rune{}
		for _, l := range line {
			if closing, ok := pairs[l]; ok {
				stack = append(stack, closing)
			} else if l != stack[len(stack)-1] {
				err = fmt.Errorf("expected %c, got %c", stack[len(stack)-1], l)
				illegal = l
				break
			} else {
				stack = stack[:len(stack)-1]
			}
		}

		if err != nil {
			fmt.Println("corrupt:", err)
			part1 += syntaxErrorScores[illegal]
		} else if len(stack) != 0 {
			fmt.Println("incomplete")
			complete := ""
			score := 0
			for i, _ := range stack {
				l := stack[len(stack)-(i+1)]
				score *= 5
				score += completionScores[l]

				complete += string(l)
			}
			fmt.Println("complete with", complete, score)
			part2 = append(part2, score)
		}

		return nil
	}); err != nil {
		return err
	}

	fmt.Println("Part 1:", part1)

	sort.Ints(part2)
	fmt.Println("Part 2:", part2[len(part2) / 2])

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
