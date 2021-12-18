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

type Token int
const (
	TokenOpen Token = -1
	TokenClose      = -2
)

type SFN []Token

func ParseSFN(s string) SFN {
	var sfn SFN

	for i := 0; i < len(s); i++ {
		l := s[i]
		if l == '[' {
			sfn = append(sfn, TokenOpen)
		} else if l == ']' {
			sfn = append(sfn, TokenClose)
		} else if l >= '0' && l <= '9' {
			count := strings.IndexAny(s[i:], "[],")

			v, err := strconv.Atoi(s[i:i+count])
			if err != nil {
				panic(err)
			}

			sfn = append(sfn, Token(v))
			i += count - 1
		}
	}

	sfn.Reduce()
	return sfn
}

func (s SFN) String() string {
	res := ""
	for _, t := range s {
		switch t {
		case TokenOpen:
			res += "[ "
		case TokenClose:
			res += "] "
		default:
			res += fmt.Sprintf("%d ", int(t))
		}
	}
	return res
}

func assert(cond bool, msg string) {
	if !cond {
		panic(msg)
	}
}

func explode(old SFN, idx int) SFN {
	// Scan left for numbers
	left := old[idx+1]
	assert(left >= 0, "left not a number")
	for j := idx; j >= 0; j-- {
		if old[j] < 0 {
			continue
		}
		old[j] += left
		break
	}

	// Scan right for numbers
	right := old[idx+2]
	assert(right >= 0, "right not a number")
	for j := idx+3; j < len(old); j++ {
		if old[j] < 0 {
			continue
		}
		old[j] += right
		break
	}

	// We're swapping 4 tokens for 1, so net decrease of 3
	reduced := SFN(make([]Token, len(old)-3))

	copy(reduced, old[:idx])
	reduced[idx] = 0
	copy(reduced[idx+1:], old[idx+4:])

	return reduced
}

func split(old SFN, idx int) SFN {
	// We're swapping 1 token for 4, so net increase of 3
	reduced := SFN(make([]Token, len(old)+3))

	copy(reduced, old[:idx])
	reduced[idx] = TokenOpen
	reduced[idx+1] = old[idx] / 2
	reduced[idx+2] = (old[idx] + 1) / 2
	reduced[idx+3] = TokenClose
	copy(reduced[idx+4:], old[idx+1:])

	return reduced
}

func (s *SFN) Reduce() {
	old := *s
	for {
		var reduced SFN
		depth := 0

		// First check for explosions
		for i, tok := range old {
			if tok == TokenOpen {
				depth += 1

				if depth > 4 {
					reduced = explode(old, i)
					break
				}
			} else if tok == TokenClose {
				depth -= 1
			}
		}
		if reduced != nil {
			old = reduced
			continue
		}

		// If there were none, look for splits
		for i, tok := range old {
			if tok >= 10 {
				reduced = split(old, i)
				break
			}
		}
		if reduced != nil {
			old = reduced
			continue
		}

		// If still nothing happened, we're done
		break
	}

	*s = old
}

func Add(a, b SFN) SFN {
	if len(a) == 0 {
		return b
	}
	if len(b) == 0 {
		return a
	}

	res := SFN(make([]Token, len(a)+len(b)+2))
	res[0] = TokenOpen
	copy(res[1:], a)
	copy(res[len(a)+1:], b)
	res[len(res)-1]= TokenClose

	res.Reduce()
	return res
}

func (s SFN) magnitude(idx int) (int, int) {
	if s[idx] >= 0 {
		return int(s[idx]), 1
	}

	left, lused := s.magnitude(idx+1)
	right, rused := s.magnitude(idx+1+lused)

	return 3 * left + 2 * right, lused + rused + 2
}

func Magnitude(s SFN) int {
	v, used := s.magnitude(0)
	if len(s) != used {
		panic(fmt.Sprintf("used %d, want %d", used, len(s)))
	}
	return v
}

func run() error {
	var sfns []SFN
	var res SFN
	doLine := func(line string) error {
		sfn := ParseSFN(line)
		sfns = append(sfns, sfn)
		res = Add(res, sfn)
		return nil
	}
	if err := doLines(os.Args[1], doLine); err != nil {
		return err
	}

	fmt.Println("Part 1:", Magnitude(res))

	largestMag := 0
	for i := 0; i < len(sfns); i++ {
		for j := 0; j < len(sfns); j++ {
			if i == j {
				continue
			}

			a := Magnitude(Add(sfns[i], sfns[j]))
			if a > largestMag {
				largestMag = a
			}

			b := Magnitude(Add(sfns[j], sfns[i]))
			if b > largestMag {
				largestMag = b
			}
		}
	}

	fmt.Println("Part 2:", largestMag)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
