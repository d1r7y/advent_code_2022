package day13

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePacketElements(t *testing.T) {
	type testCase struct {
		line                  string
		expectedPacketElement *PacketElement
	}
	testCases := []testCase{
		{
			line: "[1,1,3,1,1]",
			expectedPacketElement: &PacketElement{
				line: "[1,1,3,1,1]",
				List: PacketElementList{
					{Number: true, Value: 1, List: PacketElementList{}},
					{Number: true, Value: 1, List: PacketElementList{}},
					{Number: true, Value: 3, List: PacketElementList{}},
					{Number: true, Value: 1, List: PacketElementList{}},
					{Number: true, Value: 1, List: PacketElementList{}},
				},
			},
		},
		{
			line: "[10,1,3,1,10]",
			expectedPacketElement: &PacketElement{
				line: "[10,1,3,1,10]",
				List: PacketElementList{
					{Number: true, Value: 10, List: PacketElementList{}},
					{Number: true, Value: 1, List: PacketElementList{}},
					{Number: true, Value: 3, List: PacketElementList{}},
					{Number: true, Value: 1, List: PacketElementList{}},
					{Number: true, Value: 10, List: PacketElementList{}},
				},
			},
		},
		{
			line: "[5]",
			expectedPacketElement: &PacketElement{
				line: "[5]",
				List: PacketElementList{
					{Number: true, Value: 5, List: PacketElementList{}},
				},
			},
		},
		{
			line: "[5,[1]]",
			expectedPacketElement: &PacketElement{
				line: "[5,[1]]",
				List: PacketElementList{
					{Number: true, Value: 5, List: PacketElementList{}},
					{List: PacketElementList{
						{Number: true, Value: 1, List: PacketElementList{}},
					}},
				},
			},
		},
	}

	for _, test := range testCases {
		pe := ParsePacketElements(test.line)
		assert.Equal(t, test.expectedPacketElement, pe)
	}
}

func TestComparePacketElements(t *testing.T) {
	type testCase struct {
		str1                     string
		str2                     string
		expectedComparisonResult ComparisonResult
	}
	testCases := []testCase{
		{
			str1:                     "[1,1,3,1,1]",
			str2:                     "[1,1,3,1,1]",
			expectedComparisonResult: EqualResult,
		},
		{
			str1:                     "[1,1,3,1,1]",
			str2:                     "[1,1,3,1]",
			expectedComparisonResult: IncorrectOrder,
		},
		{
			str1:                     "[1,1,3,1]",
			str2:                     "[1,1,3,1,1]",
			expectedComparisonResult: CorrectOrder,
		},
		{
			str1:                     "[]",
			str2:                     "[1,1,3,1,1]",
			expectedComparisonResult: CorrectOrder,
		},
		{
			str1:                     "[1,1,3,1,1]",
			str2:                     "[]",
			expectedComparisonResult: IncorrectOrder,
		},
		{
			str1:                     "[[[]]]",
			str2:                     "[[]]",
			expectedComparisonResult: IncorrectOrder,
		},
		{
			str1:                     "[[]]",
			str2:                     "[[[]]]",
			expectedComparisonResult: CorrectOrder,
		},
		{
			str1:                     "[9]",
			str2:                     "[[8,7,6]]",
			expectedComparisonResult: IncorrectOrder,
		},
		{
			str1:                     "[[4,4],4,4]",
			str2:                     "[[4,4],4,4,4]",
			expectedComparisonResult: CorrectOrder,
		},
		{
			str1:                     "[1,[2,[3,[4,[5,6,7]]]],8,9]",
			str2:                     "[1,[2,[3,[4,[5,6,0]]]],8,9]",
			expectedComparisonResult: IncorrectOrder,
		},
		{
			str1:                     "[]",
			str2:                     "[[[[10],6,[1,8]],[],7]]",
			expectedComparisonResult: CorrectOrder,
		},
		{
			str1:                     "[[0,[4,1,[10,0,5],[]],4,[9]],[[4],[]],[7],[]]",
			str2:                     "[[0,[10,3,[5,7,3],[5],[]],[[9,3,7,9,7],5,[1,7,8,7],[]],[[4],0,[],2],[[7,2,10,7,0],10,5,[]]],[5,10,[5,[2,3],0,[3,10,7,10]],5],[[8,2,10],[4,[]],0,2,8],[4,[9,[2,4,10]],[[2,4,7],5],[[2]]],[[1,[0,9,3],7,[]],[[],[],2],3]]",
			expectedComparisonResult: CorrectOrder,
		},
		{
			str1:                     "[[],[8,[],[[],[9,2,2],7],[[9,7,6],10]],[[0,[8,6,6,9,5],10],10,[3]],[[1,[10,4,1,5],2,[8,8,8,4,7]],8,4,1,[1,[3,2,6,7,5],3]],[6]]",
			str2:                     "[[],[[[3,5],[6,3,5],[2,9]],[[6,10,4,2],4],9],[[[3,6,6],4],[[4],[5,0,5],2,2],[9,[8,7,5,3]]]]",
			expectedComparisonResult: IncorrectOrder,
		},
	}

	for _, test := range testCases {
		p1 := ParsePacketElements(test.str1)
		p2 := ParsePacketElements(test.str2)
		assert.Equal(t, test.expectedComparisonResult, ComparePacketElements(p1.List, p2.List), test.str1)
	}
}

func TestPairCorrectOrderIndices(t *testing.T) {
	type testCase struct {
		str             string
		expectedIndices []int
	}
	testCases := []testCase{
		{
			str: `[1,1,3,1,1]
[1,1,5,1,1]

[[1],[2,3,4]]
[[1],4]

[9]
[[8,7,6]]

[[4,4],4,4]
[[4,4],4,4,4]

[7,7,7,7]
[7,7,7]

[]
[3]

[[[]]]
[[]]

[1,[2,[3,[4,[5,6,7]]]],8,9]
[1,[2,[3,[4,[5,6,0]]]],8,9]`,
			expectedIndices: []int{1, 2, 4, 6},
		},
	}

	for _, test := range testCases {
		pairs := ParsePairs(test.str)

		assert.Equal(t, test.expectedIndices, PairCorrectOrderIndices(pairs))
	}
}

func TestPacketSorting(t *testing.T) {
	type testCase struct {
		str                  string
		expectedSortedOutput string
	}
	testCases := []testCase{
		{
			str: `[1,1,3,1,1]
[[1],[2,3,4]]
[[1],4]
[1,1,5,1,1]`,
			expectedSortedOutput: `[1,1,3,1,1]
[1,1,5,1,1]
[[1],[2,3,4]]
[[1],4]`,
		},
		{
			str: `[1,1,3,1,1]
[[1],4]
[1,1,5,1,1]`,
			expectedSortedOutput: `[1,1,3,1,1]
[1,1,5,1,1]
[[1],4]`,
		},
	}

	for _, test := range testCases {
		list := ParsePackets(test.str)

		sort.Sort(list)

		output := ""

		for i, pes := range list {
			if i != 0 {
				output += "\n"
			}

			output += pes.line
		}

		assert.Equal(t, test.expectedSortedOutput, output)
	}
}
