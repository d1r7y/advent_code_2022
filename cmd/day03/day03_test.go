package day03

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestItemPriority(t *testing.T) {
	for i := 97; i <= 122; i++ {
		item := NewItem(byte(i))

		assert.Equal(t, item.Priority().Value(), i-96)
	}

	for i := 65; i <= 90; i++ {
		item := NewItem(byte(i))

		assert.Equal(t, item.Priority().Value(), i-38)
	}
}

func TestGetCommonItem(t *testing.T) {
	r := NewRucksack("abcdEFGa")

	assert.Equal(t, getCommonItem(r), NewItem(byte('a')))
}

func getCompartmentString(c Compartment) string {
	b := []byte{}

	for _, i := range c.items {
		b = append(b, i.Value())
	}

	return string(b)
}

func TestParseRucksack(t *testing.T) {
	tests := []string{
		"vJrwpWtwJgWrhcsFMMfFFhFp",
		"jqHRNqRjqzjGDLGLrsFMfFZSrLrFZsSL",
		"PmmdzqPrVvPwwTWBwg",
	}

	for _, line := range tests {
		r := NewRucksack(line)

		assert.Equal(t, getCompartmentString(r.compartment1)+getCompartmentString(r.compartment2), line)
	}
}

func TestParseGroup(t *testing.T) {
	lines := []string{
		"vJrwpWtwJgWrhcsFMMfFFhFp",
		"jqHRNqRjqzjGDLGLrsFMfFZSrLrFZsSL",
		"PmmdzqPrVvPwwTWBwg",
	}

	r1 := NewRucksack(lines[0])
	r2 := NewRucksack(lines[1])
	r3 := NewRucksack(lines[2])

	g := NewGroup(r1, r2, r3)

	for i, r := range g.rucksacks {
		assert.Equal(t, getCompartmentString(r.compartment1)+getCompartmentString(r.compartment2), lines[i])
	}
}

func TestGetGroupBadge(t *testing.T) {
	type getGroupBadgeTest struct {
		lines         [3]string
		expectedBadge Item
	}
	tests := []getGroupBadgeTest{
		{
			lines: [3]string{
				"vJrwpWtwJgWrhcsFMMfFFhFp",
				"jqHRNqRjqzjGDLGLrsFMfFZSrLrFZsSL",
				"PmmdzqPrVvPwwTWBwg",
			},
			expectedBadge: NewItem('r'),
		},
		{
			lines: [3]string{
				"wMqvLMZHhHMvwLHjbvcjnnSBnvTQFn",
				"ttgJtRGJQctTZtZT",
				"CrZsJsPPZsGzwwsLwLmpwMDw",
			},
			expectedBadge: NewItem('Z'),
		},
	}

	for _, test := range tests {
		r1 := NewRucksack(test.lines[0])
		r2 := NewRucksack(test.lines[1])
		r3 := NewRucksack(test.lines[2])

		g := NewGroup(r1, r2, r3)

		assert.Equal(t, getGroupBadge(g), test.expectedBadge)
	}
}
