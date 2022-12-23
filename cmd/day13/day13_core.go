package day13

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type PacketElement struct {
	line   string
	Number bool

	Value int
	List  PacketElementList
}

type PacketElementList []*PacketElement

func (p PacketElementList) Len() int {
	return len(p)
}

func (p PacketElementList) Less(i, j int) bool {
	return ComparePacketElements(p[i].List, p[j].List) == CorrectOrder
}

func (p PacketElementList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func NewPacketElement() *PacketElement {
	return &PacketElement{List: make(PacketElementList, 0)}
}

type Pair struct {
	p1    *PacketElement
	line1 string

	p2    *PacketElement
	line2 string
}

func ParseNumber(str string, characterIndex int) (int, int) {
	number := 0
	valueStr := ""
	for i := characterIndex; i < len(str); i++ {
		c := str[i]
		switch c {
		case ']', ',':
			return number, i - 1
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			valueStr += string(c)
			value, err := strconv.Atoi(valueStr)
			if err != nil {
				return 0, i
			}
			number = value
		}
	}

	panic("unexpected end of string")
}

func ParsePacketElement(line string, characterIndex int, parentElement *PacketElement) int {
	for i := characterIndex; i < len(line); i++ {
		c := line[i]
		switch c {
		case '[':
			pe := NewPacketElement()
			updatedIndex := ParsePacketElement(line, i+1, pe)

			parentElement.List = append(parentElement.List, pe)

			i = updatedIndex
		case ']':
			return i
		case ',':
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			value, updatedIndex := ParseNumber(line, i)

			pe := NewPacketElement()
			pe.Number = true
			pe.Value = value

			parentElement.List = append(parentElement.List, pe)
			i = updatedIndex
		}
	}

	return 0
}

func ParsePairs(fileContents string) []Pair {
	lines := strings.Split(fileContents, "\n")

	pairs := make([]Pair, 0)

	for i := 0; i < len(lines); i++ {
		line1 := lines[i]
		pes1 := ParsePacketElements(line1)
		i++

		line2 := lines[i]
		pes2 := ParsePacketElements(line2)
		i++ // Skip over blank line

		pairs = append(pairs, Pair{p1: pes1, line1: line1, p2: pes2, line2: line2})
	}

	return pairs
}

func ParsePackets(fileContents string) PacketElementList {
	list := make(PacketElementList, 0)

	for _, line := range strings.Split(fileContents, "\n") {
		if line == "" {
			continue
		}

		pes := ParsePacketElements(line)
		list = append(list, pes)
	}

	return list
}

func PairCorrectOrderIndices(pairs []Pair) []int {
	correctIndices := make([]int, 0)

	for i, pair := range pairs {
		if ComparePacketElements(pair.p1.List, pair.p2.List) == CorrectOrder {
			correctIndices = append(correctIndices, i+1)
		}
	}

	return correctIndices
}

func ParsePacketElements(line string) *PacketElement {
	pe := NewPacketElement()

	pe.line = line
	ParsePacketElement(line, 1, pe)

	return pe
}

type ComparisonResult int

const (
	CorrectOrder ComparisonResult = iota
	IncorrectOrder
	EqualResult
)

func ComparePacketElements(p1, p2 PacketElementList) ComparisonResult {
	for i := 0; i < len(p1); i++ {
		// See if we've reached the end of p2.
		if i == len(p2) {
			return IncorrectOrder
		}
		e1 := p1[i]
		e2 := p2[i]

		if e1.Number && e2.Number {
			// Number/number
			if e1.Value < e2.Value {
				return CorrectOrder
			} else if e1.Value > e2.Value {
				return IncorrectOrder
			}
		} else if e1.Number && !e2.Number {
			// Number/list
			tempE := NewPacketElement()
			tempE.Number = false
			tempE.List = append(tempE.List, e1)
			result := ComparePacketElements(tempE.List, e2.List)
			if result != EqualResult {
				return result
			}
		} else if !e1.Number && e2.Number {
			// List/number
			tempE := NewPacketElement()
			tempE.Number = false
			tempE.List = append(tempE.List, e2)
			result := ComparePacketElements(e1.List, tempE.List)
			if result != EqualResult {
				return result
			}
		} else {
			// List/list
			result := ComparePacketElements(e1.List, e2.List)
			if result != EqualResult {
				return result
			}
		}
	}

	// If len(p1) < len(p2), then it's in the correct order.
	// If len(p1) == len(p2), then we need to keep checking.
	// If len(p1) > len(p2), then it's in the wrong order.  That's taken care of above.
	if len(p1) == len(p2) {
		return EqualResult
	}

	return CorrectOrder
}

func FindPacketIndex(line string, list PacketElementList) int {
	for i, pes := range list {
		if pes.line == line {
			return i + 1
		}
	}

	return 0
}

func day13(fileContents string) error {
	// Part 1: What is the sum of the packet pairs that are in the correct order?
	pairs := ParsePairs(fileContents)

	indexSum := 0

	for _, index := range PairCorrectOrderIndices(pairs) {
		indexSum += index
	}

	fmt.Printf("Index sum %d\n", indexSum)

	// Part 2: Break apart the pairs into individual packets.  Insert [[2]] and [[6]].  Sort the packets.
	// The decoder key is the indices of [[2]] and [[6]] multiplied together.  What's the decoder key?

	list := ParsePackets(fileContents)
	list = append(list, ParsePackets("[[2]]")...)
	list = append(list, ParsePackets("[[6]]")...)

	sort.Sort(list)

	two := FindPacketIndex("[[2]]", list)
	six := FindPacketIndex("[[6]]", list)

	fmt.Printf("Decoder key %d\n", two*six)
	return nil
}
