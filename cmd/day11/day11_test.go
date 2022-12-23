package day11

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewItem(t *testing.T) {
	item := NewItem(5)

	assert.Equal(t, big.NewInt(int64(5)), item.WorryLevel)
}

func TestParseItemList(t *testing.T) {
	type testCase struct {
		str           string
		expectedErr   bool
		expectedItems []Item
	}

	testCases := []testCase{
		{"", true, []Item{}},
		{"  Starting items: 79, 98", false, []Item{NewItem(79), NewItem(98)}},
		{"  Starting items: 74", false, []Item{NewItem(74)}},
	}

	for _, test := range testCases {
		items, err := ParseItemList(test.str)

		if test.expectedErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedItems, items)
		}
	}
}

func TestParseTest(t *testing.T) {
	type testCase struct {
		str            string
		expectedErr    bool
		item           Item
		expectedResult bool
	}

	testCases := []testCase{
		{"", true, NewItem(0), false},
		{"  Test: divisible by xxz", true, NewItem(0), false},
		{"  Tests: divisible by 100", true, NewItem(0), false},
		{"  Test: multiply by 100", true, NewItem(0), false},
		{"  Test: divisible by 23", false, NewItem(23), true},
		{"  Test: divisible by 23", false, NewItem(22), false},
	}

	for _, test := range testCases {
		testFunction, err := ParseTest(test.str)

		if test.expectedErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedResult, testFunction(test.item))
		}
	}
}

func TestParseTestResult(t *testing.T) {
	type testCase struct {
		str              string
		evaluator        string
		expectedErr      bool
		expectedMonkeyID int
	}

	testCases := []testCase{
		{"", "true", true, 0},
		{"    If zztrue: throw to monkey 2", "true", true, 0},
		{"    If true: throw to monkey bb", "true", true, 0},
		{"    If true: throw to hyena 2", "true", true, 0},
		{"    If true: throw to monkey 2", "true", false, 2},
		{"    If true: throw to monkey 5", "true", false, 5},
		{"", "false", true, 0},
		{"    If zzfalse: throw to monkey 2", "false", true, 0},
		{"    If false: throw to monkey bb", "false", true, 0},
		{"    If false: throw to hyena 2", "false", true, 0},
		{"    If false: throw to monkey 2", "false", false, 2},
		{"    If false: throw to monkey 5", "false", false, 5},
	}

	for _, test := range testCases {
		monkeyID, err := ParseTestResult(test.str, test.evaluator)

		if test.expectedErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedMonkeyID, monkeyID)
		}
	}
}

func TestParseOperation(t *testing.T) {
	type testCase struct {
		str          string
		expectedErr  bool
		item         Item
		expectedItem Item
	}

	testCases := []testCase{
		{"", true, NewItem(0), NewItem(0)},
		{"  Operation: new = old ? 19", true, NewItem(0), NewItem(0)},
		{"  Operation: new = old + bunko", true, NewItem(0), NewItem(0)},
		{"  Operation: new = old * 19", false, NewItem(10), NewItem(190)},
		{"  Operation: new = old + 6", false, NewItem(10), NewItem(16)},
		{"  Operation: new = old - 5", false, NewItem(10), NewItem(5)},
		{"  Operation: new = old / old", false, NewItem(10), NewItem(1)},
		{"  Operation: new = old / 2", false, NewItem(10), NewItem(5)},
	}

	for _, test := range testCases {
		operation, err := ParseOperation(test.str)

		if test.expectedErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedItem, operation(test.item))
		}
	}
}
