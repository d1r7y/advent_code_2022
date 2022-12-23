package day04

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

type SectionID int

func NewSectionID(i int) SectionID {
	return SectionID(i)
}

type SectionRange struct {
	startID SectionID
	endID   SectionID
}

func NewSectionRange(s, e SectionID) SectionRange {
	if int(s) > int(e) {
		panic(fmt.Sprintf("start range greater than end range (%d > %d", int(s), int(e)))
	}

	return SectionRange{startID: s, endID: e}
}

type CleaningPair struct {
	first  SectionRange
	second SectionRange
}

func NewCleaningPair(f, s SectionRange) CleaningPair {
	return CleaningPair{first: f, second: s}
}

func (cp CleaningPair) FullyContained() bool {
	// See if second is fully contained with first.
	if cp.first.startID <= cp.second.startID && cp.first.endID >= cp.second.endID {
		return true
	}

	// See if first is fully contained with second.
	if cp.second.startID <= cp.first.startID && cp.second.endID >= cp.first.endID {
		return true
	}

	return false
}

func (cp CleaningPair) Intersect() bool {
	if cp.first.endID < cp.second.startID {
		return false
	}

	if cp.first.startID > cp.second.endID {
		return false
	}

	return true
}

func ParseSectionRange(str string) (SectionRange, error) {
	if str == "" {
		return SectionRange{}, errors.New("empty string")
	}

	var s int
	var e int

	count, err := fmt.Sscanf(str, "%d-%d", &s, &e)
	if err != nil {
		return SectionRange{}, err
	}

	if count != 2 {
		return SectionRange{}, errors.New("invalid string")
	}

	return NewSectionRange(NewSectionID(s), NewSectionID(e)), nil
}

func ParseCleaningPair(str string) (CleaningPair, error) {
	if str == "" {
		return CleaningPair{}, errors.New("empty string")
	}

	cleaningElves := strings.Split(str, ",")
	if len(cleaningElves) != 2 {
		return CleaningPair{}, errors.New("invalid string")
	}

	f, err := ParseSectionRange(cleaningElves[0])
	if err != nil {
		return CleaningPair{}, err
	}

	s, err := ParseSectionRange(cleaningElves[1])
	if err != nil {
		return CleaningPair{}, err
	}

	return NewCleaningPair(f, s), nil
}

func ParseCleaningAssignments(text string) ([]CleaningPair, error) {
	assignments := make([]CleaningPair, 0)

	if text == "" {
		return assignments, nil
	}

	for _, line := range strings.Split(text, "\n") {
		cp, err := ParseCleaningPair(line)
		if err != nil {
			return []CleaningPair{}, err
		}

		assignments = append(assignments, cp)
	}

	return assignments, nil
}

func day04(fileContents string) error {
	// Part 1: In how many cleaning assignments does one SectionRange fully contain the other?
	assignments, err := ParseCleaningAssignments(string(fileContents))
	if err != nil {
		return err
	}

	totalFullyContained := 0
	for _, assignment := range assignments {
		if assignment.FullyContained() {
			totalFullyContained++
		}
	}

	log.Printf("%d assignments fully contained\n", totalFullyContained)

	// Part 2: In how many cleaning assignments is there any overlap between the SectionRanges?

	totalIntersect := 0
	for _, assignment := range assignments {
		if assignment.Intersect() {
			totalIntersect++
		}
	}

	log.Printf("%d assignments intersect\n", totalIntersect)
	return err
}
