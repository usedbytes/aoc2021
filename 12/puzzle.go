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

func contains(visited []string, place string) bool {
	for _, s := range visited {
		if s == place {
			return true
		}
	}

	return false
}

func isLower(a string) bool {
	return strings.ToLower(a) == a
}

func part1CanVisit(visited []string, node string) bool {
	if contains(visited, node) && isLower(node) {
		return false
	}

	return true
}

func part2CanVisit(visited []string, node string) bool {
	counts := make(map[string]int)

	if !isLower(node) {
		return true
	}

	if node == "start" {
		return false
	}

	if node == "end" {
		return true
	}

	for _, n := range visited {
		if isLower(n) {
			counts[n] = counts[n] + 1
		}
	}

	if counts[node] > 1 {
		return false
	}

	visitedBefore := contains(visited, node)
	for _, v := range counts {
		if v > 1 && visitedBefore {
			return false
		}
	}

	return true
}

func explore(system map[string][]string, from string, visited []string, canVisit func([]string, string) bool) [][]string {
	conns := system[from]
	visited = append(visited, from)

	routes := make([][]string, 0)

	for _, l := range conns {
		if l == "end" {
			newVisited := make([]string, len(visited), len(visited)+1)
			copy(newVisited, visited)

			routes = append(routes, append(newVisited, l))
			continue
		}

		if !canVisit(visited, l) {
			continue
		}

		routes = append(routes, explore(system, l, visited, canVisit)...)
	}

	return routes
}

func run() error {
	system := make(map[string][]string)

	if err := doLines(os.Args[1], func(line string) error {
		parts := strings.Split(line, "-")

		conns := system[parts[0]]
		system[parts[0]] = append(conns, parts[1])

		conns = system[parts[1]]
		system[parts[1]] = append(conns, parts[0])

		return nil
	}); err != nil {
		return err
	}

	routes := explore(system, "start", []string{}, part1CanVisit)
	fmt.Println("Part 1:", len(routes))

	routes = explore(system, "start", []string{}, part2CanVisit)
	fmt.Println("Part 2:", len(routes))

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
