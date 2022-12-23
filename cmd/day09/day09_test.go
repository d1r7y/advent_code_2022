package day09

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func nkmo(dir MovementDirection, count int) KnotMovementOp {
	return NewKnotMovementOp(dir, count)
}

func nkp(x, y int) KnotPosition {
	return NewKnotPosition(x, y)
}

func TestParseKnotMovementOp(t *testing.T) {
	type testCase struct {
		str                string
		expectedErr        bool
		expectedMovementOp KnotMovementOp
	}

	testCases := []testCase{
		{"", true, nkmo(UpDirection, 0)},
		{"L 2", false, nkmo(LeftDirection, 2)},
		{"U 20", false, nkmo(UpDirection, 20)},
		{"D 6", false, nkmo(DownDirection, 6)},
		{"R 1", false, nkmo(RightDirection, 1)},
		{"X 1", true, nkmo(UpDirection, 0)},
		{"55 1", true, nkmo(UpDirection, 0)},
		{"D J", true, nkmo(UpDirection, 0)},
		{"D U", true, nkmo(UpDirection, 0)},
		{"U 0", true, nkmo(UpDirection, 0)},
		{"L -11", true, nkmo(UpDirection, 0)},
	}

	for _, test := range testCases {
		kmo, err := ParseKnotMovementOp(test.str)

		if test.expectedErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedMovementOp, kmo)
		}
	}
}

func TestParseKnotMovementOps(t *testing.T) {
	type testCase struct {
		lines               []string
		expectedErr         bool
		expectedMovementOps KnotMovementOpList
	}

	testCases := []testCase{
		{[]string{""}, true, KnotMovementOpList{}},
		{[]string{"", "", ""}, true, KnotMovementOpList{}},
		{[]string{"U 1", "D 2", "L 3"}, false, KnotMovementOpList{
			nkmo(UpDirection, 1),
			nkmo(DownDirection, 2),
			nkmo(LeftDirection, 3),
		}},
		{[]string{"U 1", "D 2", "", "L 3"}, true, KnotMovementOpList{}},
		{[]string{"Z 1", "D 2", "L 3"}, true, KnotMovementOpList{}},
		{[]string{"U 1", "D -2", "L 3"}, true, KnotMovementOpList{}},
	}

	for _, test := range testCases {
		list, err := ParseKnotMovementOps(test.lines)

		if test.expectedErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedMovementOps, list)
		}
	}
}

func TestGetMovementAmount(t *testing.T) {
	type testCase struct {
		direction MovementDirection
		expectedX int
		expectedY int
	}

	testCases := []testCase{
		{UpDirection, 0, 1},
		{DownDirection, 0, -1},
		{LeftDirection, -1, 0},
		{RightDirection, 1, 0},
	}

	for _, test := range testCases {
		x, y := GetMovementAmount(test.direction)

		assert.Equal(t, test.expectedX, x)
		assert.Equal(t, test.expectedY, y)
	}
}

func withinTolerance(a, b, e float64) bool {
	if a == b {
		return true
	}

	d := math.Abs(a - b)
	if b == 0 {
		return d < e
	}

	return (d / math.Abs(b)) < e
}

func TestDistance(t *testing.T) {
	type testCase struct {
		knot1            KnotPosition
		knot2            KnotPosition
		expectedDistance float64
	}

	testCases := []testCase{
		{nkp(0, 0), nkp(0, 0), 0.0},
		{nkp(1, 1), nkp(0, 0), math.Sqrt(2.0)},
		{nkp(-1, -1), nkp(0, 0), math.Sqrt(2.0)},
		{nkp(2, 2), nkp(0, 0), 2.0 * math.Sqrt(2.0)},
	}

	for _, test := range testCases {
		distance := Distance(test.knot1, test.knot2)

		assert.True(t, withinTolerance(test.expectedDistance, distance, 1e-12))
	}
}

func TestMustMoveKnot(t *testing.T) {
	type testCase struct {
		knot1            KnotPosition
		knot2            KnotPosition
		expectedMustMove bool
	}

	testCases := []testCase{
		{nkp(0, 0), nkp(0, 0), false},
		{nkp(1, 1), nkp(0, 0), false},
		{nkp(-1, -1), nkp(0, 0), false},
		{nkp(2, 2), nkp(0, 0), true},
		{nkp(100, 50), nkp(11, 99), true},
	}

	for _, test := range testCases {
		assert.Equal(t, test.expectedMustMove, MustMoveKnot(test.knot1, test.knot2))
	}
}

func TestGetNewTailPosition(t *testing.T) {
	type testCase struct {
		head            KnotPosition
		tail            KnotPosition
		expectedNewTail KnotPosition
	}

	testCases := []testCase{
		{nkp(3, 1), nkp(1, 1), nkp(2, 1)},
		{nkp(1, 1), nkp(1, 3), nkp(1, 2)},
		{nkp(2, 3), nkp(1, 1), nkp(2, 2)},
		{nkp(3, 2), nkp(1, 1), nkp(2, 2)},
	}

	for _, test := range testCases {
		assert.Equal(t, test.expectedNewTail, GetNewKnotPosition(test.head, test.tail))
	}
}

func TestApplyMovementOp(t *testing.T) {
	type testCase struct {
		head            KnotPosition
		tail            KnotPosition
		movementOp      KnotMovementOp
		expectedNewHead KnotPosition
		expectedNewTail KnotPosition
	}

	testCases := []testCase{
		{nkp(0, 0), nkp(0, 0), nkmo(RightDirection, 4), nkp(4, 0), nkp(3, 0)},
		{nkp(2, 2), nkp(1, 1), nkmo(UpDirection, 1), nkp(2, 3), nkp(2, 2)},
	}

	for _, test := range testCases {
		w := NewWorld(2)
		w.Head = test.head
		w.RemainingKnots[0] = test.tail

		w.ApplyMovementOp(test.movementOp)

		assert.Equal(t, test.expectedNewHead, w.Head)
		assert.Equal(t, test.expectedNewTail, w.RemainingKnots[0])
	}
}

func TestManyKnotsApplyMovementOp(t *testing.T) {
	type testCase struct {
		knots         []KnotPosition
		movementOp    KnotMovementOp
		expectedKnots []KnotPosition
	}

	testCases := []testCase{
		{
			[]KnotPosition{nkp(0, 0), nkp(0, 0), nkp(0, 0), nkp(0, 0), nkp(0, 0), nkp(0, 0), nkp(0, 0), nkp(0, 0), nkp(0, 0), nkp(0, 0)},
			nkmo(RightDirection, 4),
			[]KnotPosition{nkp(4, 0), nkp(3, 0), nkp(2, 0), nkp(1, 0), nkp(0, 0), nkp(0, 0), nkp(0, 0), nkp(0, 0), nkp(0, 0), nkp(0, 0)},
		},
		{
			[]KnotPosition{nkp(4, 0), nkp(3, 0), nkp(2, 0), nkp(1, 0), nkp(0, 0), nkp(0, 0), nkp(0, 0), nkp(0, 0), nkp(0, 0), nkp(0, 0)},
			nkmo(UpDirection, 4),
			[]KnotPosition{nkp(4, 4), nkp(4, 3), nkp(4, 2), nkp(3, 2), nkp(2, 2), nkp(1, 1), nkp(0, 0), nkp(0, 0), nkp(0, 0), nkp(0, 0)},
		},
	}

	for _, test := range testCases {
		knotCount := len(test.knots)
		w := NewWorld(knotCount)

		w.Head = test.knots[0]
		w.RemainingKnots = test.knots[1:]

		w.ApplyMovementOp(test.movementOp)

		finalKnots := append([]KnotPosition{w.Head}, w.RemainingKnots...)
		assert.Equal(t, test.expectedKnots, finalKnots)
	}
}
