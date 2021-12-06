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

func run() error {
	filename := os.Args[1]

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	rd := bufio.NewReader(f)

	line, err := rd.ReadString('\n')
	if err != nil {
		return err
	}

	maxTime := 9 // At most we need 9 days
	fish := make([]int, maxTime)

	line = strings.TrimSpace(line)

	startStrings := strings.Split(line, ",")
	for _, s := range startStrings {
		v, err := strconv.Atoi(s)
		if err != nil {
			return err
		}

		fish[v]++
	}

	days := 80
	if len(os.Args) > 2 {
		v, err := strconv.Atoi(os.Args[2])
		if err != nil {
			return err
		}

		days = v
	}

	for day := 0; day < days; day++ {
		modDay := day % maxTime

		today := fish[modDay]
		fish[modDay] = 0

		// All of today's fish will spawn again in 7 days
		fish[(modDay + 7) % maxTime] += today

		// And they also spawn new fish, which will spawn in 9 days
		fish[(modDay + 9) % maxTime] += today
	}

	// How many do we have?
	total := 0
	for _, v := range fish {
		total += v
	}
	fmt.Println("Total fish:", total)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
