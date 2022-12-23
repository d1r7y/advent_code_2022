package day03

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

type Item byte

func NewItem(v byte) Item {
	return Item(v)
}

func (i Item) Priority() Priority {
	// a-z
	if i >= 97 && i <= 122 {
		return NewPriority(int(i) - 97 + 1)
	}

	// A-Z
	if i >= 65 && i <= 90 {
		return NewPriority(int(i) - 65 + 27)
	}

	return NewPriority(0)
}

func (i Item) Value() byte {
	return byte(i)
}

type Priority byte

func NewPriority(v int) Priority {
	return Priority(v)
}

func (p Priority) Value() int {
	return int(p)
}

type Compartment struct {
	items []Item
}

type Rucksack struct {
	compartment1 Compartment
	compartment2 Compartment
}

type Group struct {
	rucksacks [3]Rucksack
}

func getGroupBadge(g Group) Item {
	type itemMap map[Item]int

	rucksacksItemMap := make([]itemMap, 3)

	for index, r := range g.rucksacks {
		rucksacksItemMap[index] = make(itemMap)

		for _, item := range r.compartment1.items {
			rucksacksItemMap[index][item] = 1
		}

		for _, item := range r.compartment2.items {
			rucksacksItemMap[index][item] = 1
		}
	}

	for k := range rucksacksItemMap[0] {
		if _, ok := rucksacksItemMap[1][k]; !ok {
			continue
		}

		if _, ok := rucksacksItemMap[2][k]; !ok {
			continue
		}

		return k
	}

	return NewItem(0)
}

func NewGroup(r1, r2, r3 Rucksack) Group {
	return Group{rucksacks: [3]Rucksack{r1, r2, r3}}
}

func NewRucksack(itemsString string) Rucksack {
	itemCount := len(itemsString)

	c1Items := make([]Item, 0)
	c2Items := make([]Item, 0)

	for _, i := range itemsString[:itemCount/2] {
		c1Items = append(c1Items, NewItem(byte(i)))
	}

	for _, i := range itemsString[itemCount/2:] {
		c2Items = append(c2Items, NewItem(byte(i)))
	}

	c1 := Compartment{items: c1Items}
	c2 := Compartment{items: c2Items}

	return Rucksack{compartment1: c1, compartment2: c2}
}

func getCommonItem(r Rucksack) Item {
	if len(r.compartment1.items) != len(r.compartment2.items) {
		panic("mismatched compartment sizes")
	}

	for _, item1 := range r.compartment1.items {
		for _, item2 := range r.compartment2.items {
			if item1 == item2 {
				return item1
			}
		}
	}

	panic("no common element")
}

func ParseRucksack(line string) (Rucksack, error) {
	if line == "" {
		return Rucksack{}, errors.New("empty line")
	}

	itemCount := len(line)

	// Make sure the line has an even number of items.
	if itemCount%2 != 0 {
		return Rucksack{}, fmt.Errorf("non-even number of items '%s'", line)
	}

	return NewRucksack(line), nil
}

func ParseRucksacks(text string) ([]Rucksack, error) {
	rucksacks := make([]Rucksack, 0)

	if text == "" {
		return rucksacks, nil
	}

	for _, line := range strings.Split(text, "\n") {
		if line == "" {
			continue
		}

		rucksack, err := ParseRucksack(line)
		if err != nil {
			return []Rucksack{}, err
		}

		rucksacks = append(rucksacks, rucksack)
	}

	return rucksacks, nil
}

func ParseGroups(text string) ([]Group, error) {
	groups := make([]Group, 0)

	if text == "" {
		return groups, nil
	}

	lines := strings.Split(text, "\n")

	for i := 0; i < len(lines); i += 3 {
		r1, err := ParseRucksack(lines[i])
		if err != nil {
			return []Group{}, err
		}
		r2, err := ParseRucksack(lines[i+1])
		if err != nil {
			return []Group{}, err
		}
		r3, err := ParseRucksack(lines[i+2])
		if err != nil {
			return []Group{}, err
		}

		g := NewGroup(r1, r2, r3)

		groups = append(groups, g)
	}

	return groups, nil
}

func day03(fileContents string) error {
	// Part 1: What is the total priority of all the common elements in each rucksack?
	rucksacks, err := ParseRucksacks(string(fileContents))
	if err != nil {
		return err
	}

	totalPriority := 0

	for _, r := range rucksacks {
		commonItem := getCommonItem(r)
		totalPriority += commonItem.Priority().Value()
	}

	log.Printf("Total priority: %d\n", totalPriority)

	// Part 2: What is the total priority of all the badges for a given elf group?
	groups, err := ParseGroups(string(fileContents))
	if err != nil {
		return err
	}

	totalPriority = 0

	for _, g := range groups {
		badget := getGroupBadge(g)
		totalPriority += badget.Priority().Value()
	}

	log.Printf("Total badge priority: %d\n", totalPriority)
	return nil
}
