package day02

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetShapeScore(t *testing.T) {
	assert.Equal(t, getShapeScore(Rock), 1)
	assert.Equal(t, getShapeScore(Paper), 2)
	assert.Equal(t, getShapeScore(Scissor), 3)
}

func TestGetResultScore(t *testing.T) {
	assert.Equal(t, getResultScore(Win), 6)
	assert.Equal(t, getResultScore(Lose), 0)
	assert.Equal(t, getResultScore(Draw), 3)
}

func TestGetShapeToWin(t *testing.T) {
	assert.Equal(t, getShapeToWin(Rock), Paper)
	assert.Equal(t, getShapeToWin(Paper), Scissor)
	assert.Equal(t, getShapeToWin(Scissor), Rock)
}

func TestGetShapeToLose(t *testing.T) {
	assert.Equal(t, getShapeToLose(Rock), Scissor)
	assert.Equal(t, getShapeToLose(Paper), Rock)
	assert.Equal(t, getShapeToLose(Scissor), Paper)
}

func TestGetShapeForResult(t *testing.T) {
	type getShapeForResultTest struct {
		result        Result
		shape         Shape
		expectedShape Shape
	}

	tests := []getShapeForResultTest{
		{Win, Rock, Paper},
		{Lose, Rock, Scissor},
		{Draw, Rock, Rock},

		{Win, Paper, Scissor},
		{Lose, Paper, Rock},
		{Draw, Paper, Paper},

		{Win, Scissor, Rock},
		{Lose, Scissor, Paper},
		{Draw, Scissor, Scissor},
	}

	for _, test := range tests {
		assert.Equal(t, getShapeForResult(test.result, test.shape), test.expectedShape)
	}
}

func TestParseFirstShape(t *testing.T) {
	type parseFirstShapeTest struct {
		shapeStr      string
		expectedErr   bool
		expectedShape Shape
	}

	tests := []parseFirstShapeTest{
		{"A", false, Rock},
		{"B", false, Paper},
		{"C", false, Scissor},
		{"X", true, Rock},
		{"Y", true, Rock},
		{"Z", true, Rock},

		{"AA", true, Rock},
		{"122345", true, Rock},
		{"", true, Rock},
	}

	for _, test := range tests {
		s, err := parseFirstShape(test.shapeStr)
		if test.expectedErr {
			assert.Error(t, err)
		} else {
			assert.Equal(t, s, test.expectedShape)
		}
	}
}

func TestParseSecondShape(t *testing.T) {
	type parseFirstShapeTest struct {
		shapeStr      string
		expectedErr   bool
		expectedShape Shape
	}

	tests := []parseFirstShapeTest{
		{"A", true, Rock},
		{"B", true, Rock},
		{"C", true, Rock},
		{"X", false, Rock},
		{"Y", false, Paper},
		{"Z", false, Scissor},

		{"AA", true, Rock},
		{"122345", true, Rock},
		{"", true, Rock},
	}

	for _, test := range tests {
		s, err := parseSecondShape(test.shapeStr)
		if test.expectedErr {
			assert.Error(t, err)
		} else {
			assert.Equal(t, s, test.expectedShape)
		}
	}
}

func TestParseResult(t *testing.T) {
	type parseResultTest struct {
		resultStr      string
		expectedErr    bool
		expectedResult Result
	}

	tests := []parseResultTest{
		{"A", true, Lose},
		{"B", true, Lose},
		{"C", true, Lose},
		{"X", false, Lose},
		{"Y", false, Draw},
		{"Z", false, Win},

		{"AA", true, Lose},
		{"122345", true, Lose},
		{"", true, Lose},
	}

	for _, test := range tests {
		s, err := parseResult(test.resultStr)
		if test.expectedErr {
			assert.Error(t, err)
		} else {
			assert.Equal(t, s, test.expectedResult)
		}
	}
}

func TestRoundResult(t *testing.T) {
	type roundResultTest struct {
		shape1 Shape
		shape2 Shape
		result Result
	}

	tests := []roundResultTest{
		{Rock, Rock, Draw},
		{Paper, Paper, Draw},
		{Scissor, Scissor, Draw},

		{Rock, Paper, Win},
		{Paper, Rock, Lose},
		{Rock, Scissor, Lose},
		{Scissor, Rock, Win},
		{Paper, Scissor, Win},
		{Scissor, Paper, Lose},
	}

	for _, test := range tests {
		r := NewRound(test.shape1, test.shape2)
		assert.Equal(t, r.Result(), test.result)
	}
}

func TestRoundScore(t *testing.T) {
	type roundScoreTest struct {
		shape1 Shape
		shape2 Shape
		score  int
	}

	tests := []roundScoreTest{
		{Rock, Rock, 4},
		{Paper, Paper, 5},
		{Scissor, Scissor, 6},

		{Rock, Paper, 8},
		{Paper, Rock, 1},
		{Rock, Scissor, 3},
		{Scissor, Rock, 7},
		{Paper, Scissor, 9},
		{Scissor, Paper, 2},
	}

	for _, test := range tests {
		r := NewRound(test.shape1, test.shape2)
		assert.Equal(t, r.Score(), test.score)
	}
}

func TestParseRounds_EmptyString(t *testing.T) {
	rounds, err := ParseRounds("", true)
	assert.NoError(t, err)
	assert.Len(t, rounds, 0)

	rounds, err = ParseRounds("", false)
	assert.NoError(t, err)
	assert.Len(t, rounds, 0)
}

func TestParseRounds_Newlines(t *testing.T) {
	rounds, err := ParseRounds("\n\n\n\n\n", true)
	assert.NoError(t, err)
	assert.Len(t, rounds, 0)

	rounds, err = ParseRounds("\n\n\n\n\n", false)
	assert.NoError(t, err)
	assert.Len(t, rounds, 0)
}

func TestParseRounds_ValidRound(t *testing.T) {
	type parseRoundsTest struct {
		str            string
		partOne        bool
		expectedShape1 Shape
		expectedShape2 Shape
	}

	tests := []parseRoundsTest{
		{"A Y", true, Rock, Paper},

		{"A X", false, Rock, Scissor},
		{"A Y", false, Rock, Rock},
		{"A Z", false, Rock, Paper},

		{"B X", false, Paper, Rock},
		{"B Y", false, Paper, Paper},
		{"B Z", false, Paper, Scissor},

		{"C X", false, Scissor, Paper},
		{"C Y", false, Scissor, Scissor},
		{"C Z", false, Scissor, Rock},
	}

	for _, test := range tests {
		rounds, err := ParseRounds(test.str, test.partOne)
		assert.NoError(t, err)
		assert.Len(t, rounds, 1)
		assert.Equal(t, rounds[0].shape1, test.expectedShape1)
		assert.Equal(t, rounds[0].shape2, test.expectedShape2)
	}
}

func TestParseRounds_InvalidRound_MultipleShapes(t *testing.T) {
	rounds, err := ParseRounds("A B C", true)
	assert.Error(t, err)
	assert.Len(t, rounds, 0)

	rounds, err = ParseRounds("A B C", false)
	assert.Error(t, err)
	assert.Len(t, rounds, 0)
}

func TestParseRounds_InvalidRound_InvalidShape(t *testing.T) {
	rounds, err := ParseRounds("A T", true)
	assert.Error(t, err)
	assert.Len(t, rounds, 0)

	rounds, err = ParseRounds("A T", false)
	assert.Error(t, err)
	assert.Len(t, rounds, 0)

	rounds, err = ParseRounds("T A", true)
	assert.Error(t, err)
	assert.Len(t, rounds, 0)

	rounds, err = ParseRounds("T A", false)
	assert.Error(t, err)
	assert.Len(t, rounds, 0)

	rounds, err = ParseRounds("A B", false)
	assert.Error(t, err)
	assert.Len(t, rounds, 0)
}
