package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/pprof"
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

type ALUState struct {
	X, Y, Z, W int
}

func (s *ALUState) getDestination(name string) *int {
	switch name {
	case "x":
		return &s.X
	case "y":
		return &s.Y
	case "z":
		return &s.Z
	case "w":
		return &s.W
	}

	panic("getDestination" + name)
}

func (s *ALUState) getSource(name string) int {
	switch name {
	case "x":
		return s.X
	case "y":
		return s.Y
	case "z":
		return s.Z
	case "w":
		return s.W
	}

	var val int
	_, err := fmt.Sscanf(name, "%d", &val)
	if err != nil {
		panic(err)
	}

	return val
}

// Returns if input was consumed
func (s *ALUState) Execute(insn string, input int) bool {
	ret := false
	parts := strings.Split(insn, " ")
	switch parts[0] {
	case "inp":
		dst := s.getDestination(parts[1])
		*dst = input
		ret = true
	case "add":
		a := s.getDestination(parts[1])
		b := s.getSource(parts[2])
		*a = *a + b
	case "mul":
		a := s.getDestination(parts[1])
		b := s.getSource(parts[2])
		*a = *a * b
	case "div":
		a := s.getDestination(parts[1])
		b := s.getSource(parts[2])
		*a = *a / b
	case "mod":
		a := s.getDestination(parts[1])
		b := s.getSource(parts[2])
		*a = *a % b
	case "eql":
		a := s.getDestination(parts[1])
		b := s.getSource(parts[2])
		if *a == b {
			*a = 1
		} else {
			*a = 0
		}
	}

	return ret
}

type Operator string
const (
	OpLiteral Operator = ""
	OpVar              = "inp"
	OpRes              = "out"
	OpAdd		   = "+"
	OpMul              = "*"
	OpDiv              = "/"
	OpMod              = "%"
	OpEquals           = "=="
)

// The zero-value Expression is a '0' literal
type Expression struct {
	A, B *Expression
	Op   Operator
	Val  int
	Min, Max int
}

func (e Expression) String() string {
	if e.A == nil {
		return fmt.Sprintf("%s%d", e.Op, e.Val)
	} else {
		//return fmt.Sprintf("(%s %s %s {%d..%d})", e.A, e.Op, e.B, e.Min, e.Max)
		return fmt.Sprintf("(%s %s %s)", e.A, e.Op, e.B)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// An ugly ugly set of hand-crafted optimisation based on the patterns in the
// input.
// We track a "Min" and "Max" for each expression, so that we can eliminate
// 'eql' expressions which will always evaluate to 0.
func (e *Expression) Simplify() {
	if e.A.Op == OpLiteral && e.B.Op == OpLiteral {
		switch e.Op {
		case OpAdd:
			e.Val = e.A.Val + e.B.Val
		case OpMul:
			e.Val = e.A.Val * e.B.Val
		case OpDiv:
			e.Val = e.A.Val / e.B.Val
		case OpMod:
			e.Val = e.A.Val % e.B.Val
		case OpEquals:
			if e.A.Val == e.B.Val {
				e.Val = 1
			} else {
				e.Val = 0
			}
		}
		e.A = nil
		e.B = nil
		e.Op = OpLiteral
		e.Min, e.Max = e.Val, e.Val
	} else if e.Op == OpEquals && min(e.A.Max, e.B.Max) < max(e.A.Min, e.B.Min) {
		e.Val = 0
		e.A = nil
		e.B = nil
		e.Op = OpLiteral
		e.Min, e.Max = 0, 0
	} else if e.Op == OpEquals && e.A.Op == OpVar && e.B.Op == OpLiteral && (e.B.Val > 9 || e.B.Val < 1) {
		e.Val = 0
		e.A = nil
		e.B = nil
		e.Op = OpLiteral
		e.Min, e.Max = 0, 0
	} else if e.Op == OpEquals && e.A.Op == OpLiteral && e.B.Op == OpVar && (e.A.Val > 9 || e.A.Val < 1) {
		e.Val = 0
		e.A = nil
		e.B = nil
		e.Op = OpLiteral
		e.Min, e.Max = 0, 0
	} else if e.Op == OpMul && e.B.Op == OpLiteral && e.B.Val == 0 {
		e.Val = 0
		e.A = nil
		e.B = nil
		e.Op = OpLiteral
		e.Min, e.Max = 0, 0
	} else if e.Op == OpMul && e.A.Op == OpLiteral && e.A.Val == 0 {
		e.Val = 0
		e.A = nil
		e.B = nil
		e.Op = OpLiteral
		e.Min, e.Max = 0, 0
	} else if e.Op == OpMul && e.A.Op == OpLiteral && e.A.Val == 1 {
		e.Val = e.B.Val
		e.Op = e.B.Op
		e.Min, e.Max = e.B.Min, e.B.Max
		e.A, e.B = e.B.A, e.B.B
	} else if e.Op == OpMul && e.B.Op == OpLiteral && e.B.Val == 1 {
		e.Val = e.A.Val
		e.Op = e.A.Op
		e.Min, e.Max = e.A.Min, e.A.Max
		e.A, e.B = e.A.A, e.A.B
	} else if e.Op == OpAdd && e.A.Op == OpLiteral && e.A.Val == 0 {
		e.Val = e.B.Val
		e.Op = e.B.Op
		e.Min, e.Max = e.B.Min, e.B.Max
		e.A, e.B = e.B.A, e.B.B
	} else if e.Op == OpAdd && e.B.Op == OpLiteral && e.B.Val == 0 {
		e.Val = e.A.Val
		e.Op = e.A.Op
		e.Min, e.Max = e.A.Min, e.A.Max
		e.A, e.B = e.A.A, e.A.B
	} else if e.Op == OpDiv && e.B.Op == OpLiteral && e.B.Val == 1 {
		e.Val = e.A.Val
		e.Op = e.A.Op
		e.Min, e.Max = e.A.Min, e.A.Max
		e.A, e.B = e.A.A, e.A.B
	} else {
		switch e.Op {
		case OpVar:
			e.Min = 1
			e.Max = 9
		case OpAdd:
			e.Min = e.A.Min + e.B.Min
			e.Max = e.A.Max + e.B.Max
		case OpMul:
			e.Min = e.A.Min * e.B.Min
			e.Max = e.A.Max * e.B.Max
		case OpDiv:
			e.Min = e.A.Min / max(1, e.B.Max)
			e.Max = e.A.Max / max(1, e.B.Min)
		case OpMod:
			e.Min = 0
			e.Max = e.B.Val
			if e.B.Op != OpLiteral {
				panic("can't determine range for non-literal mod")
			}
		case OpEquals:
			e.Min = 0
			e.Max = 1
		}
	}
}

type SymbolicALUState struct {
	W, X, Y, Z *Expression
	InpCount   int
}

func NewSymbolicAlu() *SymbolicALUState {
	return &SymbolicALUState{
		W: &Expression{},
		X: &Expression{},
		Y: &Expression{},
		Z: &Expression{},
	}
}

func (s SymbolicALUState) String() string {
	return fmt.Sprintf("W: %s, X: %s, Y: %s, Z: %s", s.W, s.X, s.Y, s.Z)

}

func (s *SymbolicALUState) getDestination(name string) **Expression {
	switch name {
	case "x":
		return &s.X
	case "y":
		return &s.Y
	case "z":
		return &s.Z
	case "w":
		return &s.W
	}

	panic("getDestination" + name)
}

func (s *SymbolicALUState) getSource(name string) *Expression {
	switch name {
	case "x":
		return s.X
	case "y":
		return s.Y
	case "z":
		return s.Z
	case "w":
		return s.W
	}

	var val int
	_, err := fmt.Sscanf(name, "%d", &val)
	if err != nil {
		panic(err)
	}

	return &Expression{
		Val: val,
		Min: val,
		Max: val,
	}
}

// Returns if input was consumed
func (s *SymbolicALUState) Execute(insn string, input int) bool {
	ret := false
	parts := strings.Split(insn, " ")
	switch parts[0] {
	case "inp":
		dst := s.getDestination(parts[1])
		expr := &Expression{
			Op: OpVar,
			Val: s.InpCount,
			Min: 1,
			Max: 9,
		}
		s.InpCount++
		*dst = expr
		ret = true
	case "add":
		a := s.getDestination(parts[1])
		b := s.getSource(parts[2])
		expr := &Expression{
			A: *a,
			B: b,
			Op: OpAdd,
		}
		expr.Simplify()
		*a = expr
	case "mul":
		a := s.getDestination(parts[1])
		b := s.getSource(parts[2])
		expr := &Expression{
			A: *a,
			B: b,
			Op: OpMul,
		}
		expr.Simplify()
		*a = expr
	case "div":
		a := s.getDestination(parts[1])
		b := s.getSource(parts[2])
		expr := &Expression{
			A: *a,
			B: b,
			Op: OpDiv,
		}
		expr.Simplify()
		*a = expr
	case "mod":
		a := s.getDestination(parts[1])
		b := s.getSource(parts[2])
		expr := &Expression{
			A: *a,
			B: b,
			Op: OpMod,
		}
		expr.Simplify()
		*a = expr
	case "eql":
		a := s.getDestination(parts[1])
		b := s.getSource(parts[2])
		expr := &Expression{
			A: *a,
			B: b,
			Op: OpEquals,
		}
		expr.Simplify()
		*a = expr
	}

	return ret
}

type ALU interface {
	Execute(string, int) bool
}

func RunProgram(alu ALU, program []string, input []int) {
	i := 0
	for _, insn := range program {
		var v int
		if i < len(input) {
			v = input[i]
		}

		consumed := alu.Execute(insn, v)
		if consumed {
			if i >= len(input) {
				panic("input underflow")
			}

			i++
		}
	}
}

func parseInput(in string) []int {
	out := make([]int, len(in))
	for i, c := range in {
		out[i] = int(c - '0')
	}

	return out
}

func run() error {

	program := []string{}

	if err := doLines(os.Args[1], func(line string) error {
		program = append(program, line)
		return nil
	}); err != nil {
		return err
	}

	digits := [][]string{}
	digit := []string{}
	for _, i := range program {
		if strings.HasPrefix(i, "inp") {
			if len(digit) > 0 {
				digits = append(digits, digit)
				digit = make([]string, 0)
			}
		}
		digit = append(digit, i)
	}
	digits = append(digits, digit)

	// 'in' actually doesn't matter for the symbolic evaluation,
	// but I need it to keep the interface happy.
	in := []int{ 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1 }
	salu := NewSymbolicAlu()
	for i, p := range digits[:] {
		RunProgram(salu, p, in[i:])

		// It's quite easy to find/prove that only 'Z' is important
		// between stages.
		// To keep the output expressions manageable, we replace each
		// stage's output with a special "output" expression
		// Now, just simplify the "outX" expressions by hand 😅
		fmt.Printf("out%d: %s\n", i, salu.Z)
		salu.Z = &Expression{
			Op: OpRes,
			Val: i,
		}
	}

	// Maximum for my input
	inputVals := []int{ 9, 5, 2, 9, 9, 8, 9, 7, 9, 9, 9, 8, 9, 7 }
	input := "95299897999897"
	in = parseInput(input)
	var alu ALUState
	for i := range inputVals {
		RunProgram(&alu, digits[i], in[i:i+1])
	}
	fmt.Println("Part 1:", input, "->", alu.Z == 0)

	// Minimum for my input
	inputVals = []int{ 3, 1, 1, 1, 1, 1, 2, 1, 3, 8, 2, 1, 5, 1 }
	input = "31111121382151"
	in = parseInput(input)
	var alu2 ALUState
	for i := range inputVals {
		RunProgram(&alu2, digits[i], in[i:i+1])
	}
	fmt.Println("Part 2:", input, "->", alu2.Z == 0)

	return nil
}

func main() {
	profileEnv := os.Getenv("PROFILE")
	if profileEnv != "" {
		f, err := os.Create(profileEnv)
		if err != nil {
			fmt.Println("ERROR:", err)
			os.Exit(1)
		}

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
