package day12

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseWorld(t *testing.T) {
	type testCase struct {
		str           string
		expectedWorld *World
	}
	testCases := []testCase{
		{
			str: `Sabc
hijk
lmno
pqrE
`,
			expectedWorld: &World{
				Start:      Position{0, 0},
				End:        Position{3, 3},
				Dimensions: Size{4, 4},
				Rows: []Columns{
					{0, 0, 1, 2},
					{7, 8, 9, 10},
					{11, 12, 13, 14},
					{15, 16, 17, 25},
				},
			},
		},
	}

	for _, test := range testCases {
		w := ParseWorld(test.str)
		assert.Equal(t, test.expectedWorld, w)
	}
}

func TestDoesMoveBacktrack(t *testing.T) {
	type testCase struct {
		previous       Direction
		proposed       Direction
		expectedResult bool
	}

	testCases := []testCase{
		{Up, Up, false},
		{Up, Right, false},
		{Up, Left, false},
		{Up, Down, true},

		{Right, Right, false},
		{Right, Up, false},
		{Right, Left, true},
		{Right, Down, false},

		{Left, Left, false},
		{Left, Up, false},
		{Left, Right, true},
		{Left, Down, false},

		{Down, Down, false},
		{Down, Up, true},
		{Down, Right, false},
		{Down, Left, false},
	}

	str := `Sabc
hijk
lmno
pqrE
`

	for _, test := range testCases {
		w := ParseWorld(str)

		s := NewSolutionState(w)
		s.Moves = append(s.Moves, test.previous)

		assert.Equal(t, test.expectedResult, s.DoesMoveBacktrack(test.proposed))
	}
}

func TestDoesMoveExceedBounds(t *testing.T) {
	type testCase struct {
		proposed       Direction
		dimensions     Size
		position       Position
		expectedResult bool
	}

	testCases := []testCase{
		{Up, Size{4, 4}, Position{0, 0}, true},
		{Down, Size{4, 4}, Position{0, 0}, false},
		{Left, Size{4, 4}, Position{0, 0}, true},
		{Right, Size{4, 4}, Position{0, 0}, false},

		{Up, Size{4, 4}, Position{1, 1}, false},
		{Down, Size{4, 4}, Position{1, 1}, false},
		{Left, Size{4, 4}, Position{1, 1}, false},
		{Right, Size{4, 4}, Position{1, 1}, false},

		{Up, Size{4, 4}, Position{3, 3}, false},
		{Down, Size{4, 4}, Position{3, 3}, true},
		{Left, Size{4, 4}, Position{3, 3}, false},
		{Right, Size{4, 4}, Position{3, 3}, true},

		{Up, Size{4, 4}, Position{3, 0}, true},
		{Down, Size{4, 4}, Position{0, 3}, true},
		{Left, Size{4, 4}, Position{0, 3}, true},
		{Right, Size{4, 4}, Position{3, 0}, true},
	}

	for _, test := range testCases {
		w := NewWorld()
		w.Dimensions = test.dimensions

		s := NewSolutionState(w)
		s.SetPosition(test.position)

		assert.Equal(t, test.expectedResult, s.DoesMoveExceedBounds(test.proposed))
	}
}

func TestIsProposedDestinationHeightInvalid(t *testing.T) {
	type testCase struct {
		str            string
		position       Position
		proposed       Direction
		expectedResult bool
	}
	testCases := []testCase{
		{
			str: `Sabc
hijk
lmno
pqrE
`,
			position:       Position{0, 0},
			proposed:       Down,
			expectedResult: true,
		},
		{
			str: `Sabc
hijk
lmno
pqrE
`,
			position:       Position{0, 0},
			proposed:       Right,
			expectedResult: false,
		},
		{
			str: `Sabc
hijk
lmno
pqrE
`,
			position:       Position{3, 3},
			proposed:       Up,
			expectedResult: false,
		},
		{
			str: `Sabc
hijk
lmnz
pqrE
`,
			position:       Position{3, 3},
			proposed:       Up,
			expectedResult: false,
		},
	}

	for _, test := range testCases {
		w := ParseWorld(test.str)

		s := NewSolutionState(w)
		s.SetPosition(test.position)

		assert.Equal(t, test.expectedResult, s.IsProposedDestinationHeightInvalid(test.proposed))
	}
}

func TestGetLegalMoves(t *testing.T) {
	type testCase struct {
		str           string
		position      Position
		expectedMoves []Direction
	}
	testCases := []testCase{
		{
			str: `Sabc
hijk
lmno
pqrE
`,
			position:      Position{0, 0},
			expectedMoves: []Direction{Right},
		},
		{
			str: `Sabc
hijk
lmno
pqrE
`,
			position:      Position{3, 3},
			expectedMoves: []Direction{Up, Left},
		},
		{
			str: `Sabc
hijk
lmno
pqrE
`,
			position:      Position{2, 2},
			expectedMoves: []Direction{Up, Left, Right},
		},
		{
			str: `Sabc
hijk
lmzo
pqrE
`,
			position:      Position{2, 2},
			expectedMoves: []Direction{Up, Down, Left, Right},
		},
	}

	for _, test := range testCases {
		w := ParseWorld(test.str)

		s := NewSolutionState(w)
		s.SetPosition(test.position)

		assert.Equal(t, test.expectedMoves, s.GetLegalMoves())
	}
}

func TestWasPositionVisited(t *testing.T) {
	type testCase struct {
		historicalPositions []Position
		currentPosition     Position
		direction           Direction
		expectedResult      bool
	}
	testCases := []testCase{
		{
			historicalPositions: []Position{{3, 3}, {3, 2}, {2, 2}, {2, 3}},
			currentPosition:     Position{1, 2},
			direction:           Right,
			expectedResult:      true,
		},
		{
			historicalPositions: []Position{{3, 3}, {3, 2}, {2, 2}, {2, 3}},
			currentPosition:     Position{1, 1},
			direction:           Right,
			expectedResult:      false,
		},
	}

	str := `Sabc
hijk
lmno
pqrE
`

	for _, test := range testCases {
		w := ParseWorld(str)

		s := NewSolutionState(w)
		for _, p := range test.historicalPositions {
			s.SetPosition(p)
		}

		s.SetPosition(test.currentPosition)

		assert.Equal(t, test.expectedResult, s.WasPositionVisited(test.direction))
	}
}

func TestFindMinimumMovement1(t *testing.T) {
	str := `Sabqponm
abcryxxl
accszExk
acctuvwj
abdefghi`

	w := ParseWorld(str)

	assert.Equal(t, 31, FindMinimumMovement(w))
}

func TestFindMinimumMovement2(t *testing.T) {
	str := `abccccccccccccccccccaaaaaaaaacccccccccccccccccccccccccccccccccccccaaaa
abcccccccccccccccaaaaaaaaaaacccccccccccccccccccccccccccccccccccccaaaaa
abcaaccaacccccccccaaaaaaaaaacccccccccccccccccccccaaacccccccccccccaaaaa
abcaaaaaaccccccccaaaaaaaaaaaaacccccccccccccccccccaacccccccccccccaaaaaa
abcaaaaaacccaaacccccaaaaaaaaaaaccccccccccccccccccaaaccccccccccccccccaa
abaaaaaaacccaaaaccccaaaaaacaaaacccccccccccaaaacjjjacccccccccccccccccca
abaaaaaaaaccaaaaccccaaaaaaccccccaccccccccccaajjjjjkkcccccccccccccccccc
abaaaaaaaaccaaacccccccaaaccccccaaccccccccccajjjjjjkkkaaacccaaaccaccccc
abccaaacccccccccccccccaaccccaaaaaaaacccccccjjjjoookkkkaacccaaaaaaccccc
abcccaacccccccccccccccccccccaaaaaaaaccccccjjjjoooookkkkcccccaaaaaccccc
abcccccccaacccccccccccccccccccaaaacccccccijjjoooooookkkkccaaaaaaaccccc
abccaaccaaaccccccccccccccccccaaaaacccccciijjooouuuoppkkkkkaaaaaaaacccc
abccaaaaaaaccccccccccaaaaacccaacaaaccciiiiiooouuuuupppkkklllaaaaaacccc
abccaaaaaacccccccccccaaaaacccacccaaciiiiiiqooouuuuuupppkllllllacaccccc
abcccaaaaaaaacccccccaaaaaaccccaacaiiiiiqqqqoouuuxuuupppppplllllccccccc
abccaaaaaaaaaccaaaccaaaaaaccccaaaaiiiiqqqqqqttuxxxuuuppppppplllccccccc
abccaaaaaaaacccaaaaaaaaaaacccaaaahiiiqqqttttttuxxxxuuuvvpppplllccccccc
abcaaaaaaacccaaaaaaaaaaacccccaaaahhhqqqqtttttttxxxxuuvvvvvqqlllccccccc
abcccccaaaccaaaaaaaaaccccccccacaahhhqqqttttxxxxxxxyyyyyvvvqqlllccccccc
abcccccaaaccaaaaaaaacccccccccccaahhhqqqtttxxxxxxxyyyyyyvvqqqlllccccccc
SbcccccccccccaaaaaaaaaccccccccccchhhqqqtttxxxxEzzzyyyyvvvqqqmmlccccccc
abcccccccccccaaaaaaaacccaacccccccchhhppptttxxxxyyyyyvvvvqqqmmmcccccccc
abccccccccccaaaaaaaaaaccaacccccccchhhpppptttsxxyyyyyvvvqqqmmmccccccccc
abcaacccccccaaaaaaacaaaaaaccccccccchhhppppsswwyyyyyyyvvqqmmmmccccccccc
abaaaacccccccaccaaaccaaaaaaacccccccchhhpppsswwyywwyyyvvqqmmmddcccccccc
abaaaaccccccccccaaaccaaaaaaacccccccchhhpppsswwwwwwwwwvvqqqmmdddccccccc
abaaaacccccccccaaaccaaaaaaccccccccccgggpppsswwwwrrwwwwvrqqmmdddccccccc
abccccccaaaaaccaaaacaaaaaaccccccaacccggpppssswwsrrrwwwvrrqmmdddacccccc
abccccccaaaaaccaaaacccccaaccccaaaaaacggpppssssssrrrrrrrrrnmmdddaaccccc
abcccccaaaaaaccaaaccccccccccccaaaaaacggppossssssoorrrrrrrnnmdddacccccc
abcccccaaaaaaccccccccaaaaccccccaaaaacgggoooossoooonnnrrnnnnmddaaaacccc
abccccccaaaaaccccccccaaaacccccaaaaaccgggoooooooooonnnnnnnnndddaaaacccc
abccccccaaaccccccccccaaaacccccaaaaacccgggoooooooffennnnnnnedddaaaacccc
abcccccccccccccccccccaaacccccccaacccccggggffffffffeeeeeeeeeedaaacccccc
abccccccccccccccccccaaacccccaccaaccccccggfffffffffeeeeeeeeeecaaacccccc
abccccccccccccccccccaaaacccaaaaaaaaaccccfffffffaaaaaeeeeeecccccccccccc
abccccccccaacaaccccaaaaaacaaaaaaaaaaccccccccccaaaccaaaaccccccccccccccc
abccccccccaaaaacccaaaaaaaaaaacaaaaccccccccccccaaaccccaaccccccccccaaaca
abcccccccaaaaaccccaaaaaaaaaaacaaaaacccccccccccaaaccccccccccccccccaaaaa
abcccccccaaaaaacccaaaaaaaaaacaaaaaacccccccccccaaccccccccccccccccccaaaa
abcccccccccaaaaccaaaaaaaaaaaaaaccaaccccccccccccccccccccccccccccccaaaaa`

	w := ParseWorld(str)

	assert.Equal(t, 352, FindMinimumMovement(w))
}
