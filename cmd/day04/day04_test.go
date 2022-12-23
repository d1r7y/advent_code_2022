package day04

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSectionID(t *testing.T) {
	id := NewSectionID(5)

	assert.Equal(t, int(id), 5)
}

func TestSectionRange(t *testing.T) {
	sr := NewSectionRange(NewSectionID(5), NewSectionID(10))
	assert.Equal(t, int(sr.startID), 5)
	assert.Equal(t, int(sr.endID), 10)
}

func TestCleaningPair(t *testing.T) {
	sr1 := NewSectionRange(NewSectionID(5), NewSectionID(10))
	sr2 := NewSectionRange(NewSectionID(15), NewSectionID(23))
	cp := NewCleaningPair(sr1, sr2)

	assert.Equal(t, int(cp.first.startID), 5)
	assert.Equal(t, int(cp.first.endID), 10)
	assert.Equal(t, int(cp.second.startID), 15)
	assert.Equal(t, int(cp.second.endID), 23)
}

func TestParseSectionRange_Valid(t *testing.T) {
	type parseSectionRangeTest struct {
		str                  string
		expectedSectionRange SectionRange
	}

	tests := []parseSectionRangeTest{
		{"1-20", NewSectionRange(NewSectionID(1), NewSectionID(20))},
		{"4-4", NewSectionRange(NewSectionID(4), NewSectionID(4))},
	}

	for _, test := range tests {
		sr, err := ParseSectionRange(test.str)
		assert.NoError(t, err)
		assert.Equal(t, test.expectedSectionRange, sr)
	}
}

func TestParseSectionRange_Invalid(t *testing.T) {
	type parseSectionRangeTest struct {
		str string
	}

	tests := []parseSectionRangeTest{
		{""},
		{"ax?"},
		{"1*5"},
		{"1,2,3"},
	}

	for _, test := range tests {
		_, err := ParseSectionRange(test.str)
		assert.Error(t, err)
	}
}

func TestParseCleaningPair_Valid(t *testing.T) {
	type parseCleaningPairTest struct {
		str                  string
		expectedCleaningPair CleaningPair
	}

	tests := []parseCleaningPairTest{
		{"1-20,25-30", NewCleaningPair(NewSectionRange(NewSectionID(1), NewSectionID(20)), NewSectionRange(NewSectionID(25), NewSectionID(30)))},
		{"4-4,8-8", NewCleaningPair(NewSectionRange(NewSectionID(4), NewSectionID(4)), NewSectionRange(NewSectionID(8), NewSectionID(8)))},
	}

	for _, test := range tests {
		cp, err := ParseCleaningPair(test.str)
		assert.NoError(t, err)
		assert.Equal(t, test.expectedCleaningPair, cp)
	}
}

func TestParseCleaningPair_Invalid(t *testing.T) {
	type parseSectionRangeTest struct {
		str string
	}

	tests := []parseSectionRangeTest{
		{""},
		{"ax?,asj"},
		{"1*5,,"},
	}

	for _, test := range tests {
		_, err := ParseCleaningPair(test.str)
		assert.Error(t, err)
	}
}

func TestCleaningPairFullyContained(t *testing.T) {
	type fullyContainedTest struct {
		str                    string
		expectedFullyContained bool
	}

	tests := []fullyContainedTest{
		{"1-20,25-30", false},
		{"4-4,8-8", false},

		{"4-8,8-8", true},
		{"8-8,4-8", true},
		{"4-8,8-10", false},
		{"8-10,4-8", false},

		{"4-10,5-8", true},
	}

	for _, test := range tests {
		cp, err := ParseCleaningPair(test.str)
		assert.NoError(t, err)
		assert.Equal(t, test.expectedFullyContained, cp.FullyContained())
	}
}

func TestCleaningPairIntersect(t *testing.T) {
	type intersectTest struct {
		str               string
		expectedIntersect bool
	}

	tests := []intersectTest{
		{"1-20,25-30", false},

		{"4-4,8-8", false},

		{"4-8,8-8", true},
		{"8-8,4-8", true},

		{"4-8,8-10", true},
		{"8-10,4-8", true},

		{"4-10,5-8", true},

		{"4-10,11-16", false},
		{"11-16,4-10", false},
	}

	for _, test := range tests {
		cp, err := ParseCleaningPair(test.str)
		assert.NoError(t, err)
		assert.Equal(t, test.expectedIntersect, cp.Intersect())
	}
}
