package day08

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTreeRow(t *testing.T) {
	type testCase struct {
		str             string
		expectedErr     bool
		expectedTreeRow TreeRow
	}

	testCases := []testCase{
		{"30373", false, TreeRow{3, 0, 3, 7, 3}},
		{"25512", false, TreeRow{2, 5, 5, 1, 2}},
		{"65332", false, TreeRow{6, 5, 3, 3, 2}},
		{"33549", false, TreeRow{3, 3, 5, 4, 9}},
		{"35390", false, TreeRow{3, 5, 3, 9, 0}},
		{"3a390", true, TreeRow{}},
		{"", true, TreeRow{}},
	}

	for _, test := range testCases {
		row, err := ParseTreeRow(test.str)

		if test.expectedErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedTreeRow, row)
		}
	}
}

func TestParseForest(t *testing.T) {
	strs := []string{
		"30373",
		"25512",
		"65332",
		"33549",
		"35390",
	}
	rows := []TreeRow{
		{3, 0, 3, 7, 3},
		{2, 5, 5, 1, 2},
		{6, 5, 3, 3, 2},
		{3, 3, 5, 4, 9},
		{3, 5, 3, 9, 0},
	}

	f, err := ParseForest(strs)
	assert.NoError(t, err)
	for i, r := range f.Trees {
		assert.Equal(t, rows[i], r)
	}
}

func TestNumberVisibleTrees(t *testing.T) {
	strs := []string{
		"30373",
		"25512",
		"65332",
		"33549",
		"35390",
	}
	f, err := ParseForest(strs)
	assert.NoError(t, err)
	assert.Equal(t, 21, f.NumberVisibleTrees())
}

func TestScenicScoreForTree(t *testing.T) {
	strs := []string{
		"30373",
		"25512",
		"65332",
		"33549",
		"35390",
	}
	f, err := ParseForest(strs)
	assert.NoError(t, err)

	assert.Equal(t, 4, f.scenicScoreForTree(2, 1))
	assert.Equal(t, 8, f.scenicScoreForTree(2, 3))
}
