package day06

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPacketMarkerStart(t *testing.T) {
	type getPacketMarkerStartTestCase struct {
		str                 string
		expectedValidOffset int
	}

	testCases := []getPacketMarkerStartTestCase{
		{"mjqjpqmgbljsphdztnvjfqwrcgsmlb", 7},
		{"bvwbjplbgvbhsrlpgdmjqwftvncz", 5},
		{"nppdvjthqldpwncqszvftbrmjlhg", 6},
		{"nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg", 10},
		{"zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw", 11},
	}

	for _, test := range testCases {
		ds := Datastream(test.str)

		assert.Equal(t, test.expectedValidOffset, ds.GetPacketMarkerStart())
	}
}

func TestGetMessageMarkerStart(t *testing.T) {
	type getMessageMarkerStartTestCase struct {
		str                 string
		expectedValidOffset int
	}

	testCases := []getMessageMarkerStartTestCase{
		{"mjqjpqmgbljsphdztnvjfqwrcgsmlb", 19},
		{"bvwbjplbgvbhsrlpgdmjqwftvncz", 23},
		{"nppdvjthqldpwncqszvftbrmjlhg", 23},
		{"nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg", 29},
		{"zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw", 26},
	}

	for _, test := range testCases {
		ds := Datastream(test.str)

		assert.Equal(t, test.expectedValidOffset, ds.GetMessageMarkerStart())
	}
}
