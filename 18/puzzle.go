package main

import (
	"bufio"
	"bytes"
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

type SFN []string

func sfnScanner(s string) *bufio.Scanner {
	scanner := bufio.NewScanner(strings.NewReader(s))
	scanner.Split(func (data []byte, atEOF bool) (int, []byte, error) {
		skip := 0
		for {
			if len(data) == 0 {
				return 0, nil, nil
			}

			if data[0] == ',' {
				skip++
				data = data[1:]
				continue
			}

			if data[0] == '[' {
				return skip+1, data[:1], nil
			}

			if data[0] == ']' {
				return skip+1, data[:1], nil
			}

			end := bytes.IndexAny(data, "[],")
			if end == -1 {
				return 0, nil, nil
			} else {
				return skip+end, data[:end], nil
			}

			return 0, nil, fmt.Errorf("unexpected condition: %s", data)
		}
	});

	return scanner
}

func ParseSFN(s string) SFN {
	var sfn SFN
	scanner := sfnScanner(s)
	for scanner.Scan() {
		tok := scanner.Text()
		sfn = append(sfn, tok)
	}

	sfn.Reduce()
	return sfn
}

func (s SFN) String() string {
	return strings.Join([]string(s), " ")
}

func explode(old SFN, idx int) SFN {
	left, err := strconv.Atoi(old[idx+1])
	if err != nil {
		panic(err)
	}

	right, err := strconv.Atoi(old[idx+2])
	if err != nil {
		panic(err)
	}

	var j int
	for j = idx; j >= 0; j-- {
		t := old[j]
		if t == "[" || t == "]" {
			continue
		}

		v, err := strconv.Atoi(t)
		if err != nil {
			panic(err)
		}

		v += left

		old[j] = fmt.Sprintf("%d", v)
		break
	}

	for j := idx+3; j < len(old); j++ {
		t := old[j]
		if t == "[" || t == "]" {
			continue
		}

		v, err := strconv.Atoi(t)
		if err != nil {
			panic(err)
		}

		v += right

		old[j] = fmt.Sprintf("%d", v)
		break
	}

	// We're swapping 4 tokens for 1, so net decrease of 3
	reduced := SFN(make([]string, len(old)-3))

	copy(reduced, old[:idx])
	reduced[idx] = "0"
	copy(reduced[idx+1:], old[idx+4:])

	return reduced
}

func split(old SFN, idx, val int) SFN {
	// We're swapping 1 token for 4, so net increase of 3
	reduced := SFN(make([]string, len(old)+3))

	copy(reduced, old[:idx])
	reduced[idx] = "["
	reduced[idx+1] = fmt.Sprintf("%d", val / 2)
	reduced[idx+2] = fmt.Sprintf("%d", (val + 1) / 2)
	reduced[idx+3] = "]"
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
			if tok == "[" {
				depth += 1

				if depth > 4 {
					reduced = explode(old, i)
					break
				}
			} else if tok == "]" {
				depth -= 1
			}
		}
		if reduced != nil {
			old = reduced
			continue
		}

		// If there were none, look for splits
		for i, tok := range old {
			if tok == "[" || tok == "]" {
				continue
			} else {
				v, err := strconv.Atoi(tok)
				if err != nil {
					panic(err)
				}

				if v >= 10 {
					reduced = split(old, i, v)
					break
				}
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

	res := SFN(make([]string, len(a)+len(b)+2))
	res[0] = "["
	copy(res[1:], a)
	copy(res[len(a)+1:], b)
	res[len(res)-1]= "]"

	res.Reduce()
	return res
}

func (s SFN) magnitude(idx int) (int, int) {
	if s[idx] != "[" && s[idx] != "]" {
		v, err := strconv.Atoi(s[idx])
		if err != nil {
			panic(err)
		}
		return v, 1
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
