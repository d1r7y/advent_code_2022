package day14

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDescribe(t *testing.T) {
	type testCase struct {
		paths               string
		expectedDescription string
	}
	testCases := []testCase{
		{
			paths: `498,4 -> 498,6 -> 496,6
503,4 -> 502,4 -> 502,9 -> 494,9`,
			expectedDescription: `......+...
..........
..........
..........
....#...##
....#...#.
..###...#.
........#.
........#.
#########.`,
		},
	}

	for _, test := range testCases {
		cave := ParseCave(test.paths, true)
		assert.Equal(t, test.expectedDescription, cave.Describe())
	}
}

func TestDropSand(t *testing.T) {
	type testCase struct {
		paths               string
		sandCount           int
		expectedDescription string
	}
	testCases := []testCase{
		{
			paths: `498,4 -> 498,6 -> 496,6
503,4 -> 502,4 -> 502,9 -> 494,9`,
			sandCount: 1,
			expectedDescription: `......+...
..........
..........
..........
....#...##
....#...#.
..###...#.
........#.
......o.#.
#########.`,
		},
		{
			paths: `498,4 -> 498,6 -> 496,6
503,4 -> 502,4 -> 502,9 -> 494,9`,
			sandCount: 5,
			expectedDescription: `......+...
..........
..........
..........
....#...##
....#...#.
..###...#.
......o.#.
....oooo#.
#########.`,
		},
		{
			paths: `498,4 -> 498,6 -> 496,6
503,4 -> 502,4 -> 502,9 -> 494,9`,
			sandCount: 24,
			expectedDescription: `......+...
..........
......o...
.....ooo..
....#ooo##
...o#ooo#.
..###ooo#.
....oooo#.
.o.ooooo#.
#########.`,
		},
	}

	for _, test := range testCases {
		cave := ParseCave(test.paths, true)
		for i := 0; i < test.sandCount; i++ {
			assert.Equal(t, SandAtRest, cave.DropSand())
		}
		assert.Equal(t, test.expectedDescription, cave.Describe())
	}
}

func TestDropSandFalling(t *testing.T) {
	type testCase struct {
		paths               string
		sandCount           int
		expectedDescription string
	}
	testCases := []testCase{
		{
			paths: `498,4 -> 498,6 -> 496,6
503,4 -> 502,4 -> 502,9 -> 494,9`,
			sandCount: 24,
			expectedDescription: `......+...
..........
......o...
.....ooo..
....#ooo##
...o#ooo#.
..###ooo#.
....oooo#.
.o.ooooo#.
#########.`,
		},
	}

	for _, test := range testCases {
		cave := ParseCave(test.paths, true)
		for i := 0; i < test.sandCount; i++ {
			assert.Equal(t, SandAtRest, cave.DropSand())
		}
		assert.Equal(t, test.expectedDescription, cave.Describe())

		// Now if we drop one more sand, the it should fall to infinity.
		assert.Equal(t, SandFalling, cave.DropSand())
	}
}

func TestDropSandBlocked(t *testing.T) {
	type testCase struct {
		paths     string
		sandCount int
	}
	testCases := []testCase{
		{
			paths: `498,4 -> 498,6 -> 496,6
503,4 -> 502,4 -> 502,9 -> 494,9`,
			sandCount: 92,
		},
	}

	for _, test := range testCases {
		cave := ParseCave(test.paths, false)
		for i := 0; i < test.sandCount; i++ {
			assert.Equal(t, SandAtRest, cave.DropSand())
		}

		// Now if we drop one more sand, the it should block the source.
		assert.Equal(t, SandBlocked, cave.DropSand())
	}
}
