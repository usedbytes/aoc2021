package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

type PacketType int

type Packet struct {
	Version int
	Type	int

	Literal int

	Children []*Packet
}

// Returns value, consumed
func decodeLiteral(s string) (int, int) {
	idx := 0
	val := 0

	for {
		chunk, err := strconv.ParseUint(s[idx+1:idx+5], 2, 4)
		if err != nil {
			panic(err)
		}

		val *= 16
		val += int(chunk)

		if s[idx] == '0' {
			break
		}

		idx += 5
	}

	// Account for the last chunk
	idx += 5

	return val, idx
}

func decodeOperatorChildren(s string) ([]*Packet, int) {
	idx := 0
	ltype := s[0]
	idx += 1

	children := []*Packet{}
	switch ltype {
	case '0':
		sublength, err := strconv.ParseUint(s[idx:idx+15], 2, 15)
		if err != nil {
			panic(err)
		}
		idx += 15
		n := 0
		for n < int(sublength) {
			p, m := decodePacket(s[idx:])
			children = append(children, p)
			n += m
			idx += m
		}
	case '1':
		subpkts, err := strconv.ParseUint(s[idx:idx+11], 2, 11)
		if err != nil {
			panic(err)
		}
		idx += 11
		for n := 0; n < int(subpkts); n++ {
			p, m := decodePacket(s[idx:])
			children = append(children, p)
			idx += m
		}
	}

	return children, idx
}

func decodePacket(s string) (*Packet, int) {
	idx := 0
	ver, err := strconv.ParseUint(s[idx:idx+3], 2, 3)
	if err != nil {
		panic(err)
	}
	idx += 3

	t, err := strconv.ParseUint(s[idx:idx+3], 2, 3)
	if err != nil {
		panic(err)
	}
	idx += 3

	p := &Packet{
		Version: int(ver),
		Type: int(t),
	}

	switch t {
	case 4:
		// Literal
		v, c := decodeLiteral(s[idx:])
		idx += c
		p.Literal = v
	default:
		// Operator
		ps, c := decodeOperatorChildren(s[idx:])
		p.Children = ps
		idx += c
	}

	return p, idx
}

func sumVersions(p *Packet) int {
	v := p.Version
	for _, c := range p.Children {
		v += sumVersions(c)
	}

	return v
}

func run() error {
	f, err := os.Open(os.Args[1])
	if err != nil {
		return err
	}
	defer f.Close()

	bs, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	// Simpler than using strconv
	// Keeping the data as a big binary string makes bit slicing
	// easier, and it's not very big.
	lut := map[byte]string {
		byte('0'): "0000",
		byte('1'): "0001",
		byte('2'): "0010",
		byte('3'): "0011",
		byte('4'): "0100",
		byte('5'): "0101",
		byte('6'): "0110",
		byte('7'): "0111",
		byte('8'): "1000",
		byte('9'): "1001",
		byte('A'): "1010",
		byte('B'): "1011",
		byte('C'): "1100",
		byte('D'): "1101",
		byte('E'): "1110",
		byte('F'): "1111",
	}

	s := ""
	for _, b := range bs {
		s += lut[b]
	}

	p, _ := decodePacket(s)

	fmt.Println("Part 1:", sumVersions(p))

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
