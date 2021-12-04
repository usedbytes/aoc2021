package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Board struct {
	Numbers   [5][5]int
	ColMarked [5]int
	RowMarked [5]int
}

func NewBoard(rd *bufio.Reader) (*Board, error) {
	var b Board

	i := 0
	for i < 5 {
		line, err := rd.ReadString('\n')
		if err != nil {
			return nil, err
		}

		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		}

		strs := strings.Split(line, " ")
		j := 0
		for _, s := range strs {
			if len(s) == 0 {
				continue
			}

			n, err := strconv.Atoi(s)
			if err != nil {
				return nil, err
			}

			b.Numbers[i][j] = n
			j++
		}

		i++
	}

	return &b, nil
}

func (b *Board) Score() int {
	score := 0
	for _, row := range b.Numbers {
		for _, num := range row {
			if num >= 0 {
				score += num
			}
		}
	}

	return score
}

func (b *Board) Mark(n int) (bool, int) {
	for i, row := range b.Numbers {
		for j, number := range row {
			if number != n {
				continue
			}

			b.RowMarked[i]++
			b.ColMarked[j]++
			b.Numbers[i][j] = -1

			if b.RowMarked[i] == 5 || b.ColMarked[j] == 5 {
				return true, b.Score()
			}
		}
	}

	return false, 0
}

func run() error {
	filename := os.Args[1]

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	rd := bufio.NewReader(f)

	// Decode moves
	movesLine, err := rd.ReadString('\n')
	if err != nil {
		return err
	}
	movesLine = strings.TrimSpace(movesLine)
	moveStrings := strings.Split(movesLine, ",")
	moves := make([]int, len(moveStrings))

	for i, m := range moveStrings {
		n, err := strconv.Atoi(m)
		if err != nil {
			return err
		}

		moves[i] = n
	}

	// Decode Boards
	var boards []*Board
	var b *Board
	for {
		b, err = NewBoard(rd)
		if err != nil {
			break
		}
		boards = append(boards, b)
	}

	if err != io.EOF {
		return err
	}

	// Play the game
	numWins := 0
	for _, m := range moves {
		for i, b := range boards {
			if b == nil {
				continue
			}

			win, score := b.Mark(m)
			if win {
				boards[i] = nil
				numWins++
				if numWins == 1 {
					fmt.Println("Part 1:", score * m)
				} else if numWins == len(boards) {
					fmt.Println("Part 2:", score * m)
					return nil
				}
			}
		}
	}

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
