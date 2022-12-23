package day17

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewItem(t *testing.T) {
	assert.True(t, true)
}

func TestMakeNewShape(t *testing.T) {
	expectedTypes := []reflect.Type{
		reflect.TypeOf((*HorizontalLineShape)(nil)).Elem(),
		reflect.TypeOf((*CrossShape)(nil)).Elem(),
		reflect.TypeOf((*AngleShape)(nil)).Elem(),
		reflect.TypeOf((*VerticalLineShape)(nil)).Elem(),
		reflect.TypeOf((*SquareShape)(nil)).Elem(),
	}

	room := NewRoom(7, nil)

	// Allocate 20 shapes and verify they always come out in the correct order.
	for i := 0; i < 20; i++ {
		shape := room.MakeNextShape()
		assert.Equal(t, expectedTypes[i%len(expectedTypes)], reflect.TypeOf(shape).Elem())
	}
}

func TestParseJetDirections(t *testing.T) {
	type testCase struct {
		str                   string
		expectedErr           bool
		expectedJetDirections JetDirectionList
	}

	testCases := []testCase{
		{"", false, JetDirectionList{}},
		{"Q", true, JetDirectionList{}},
		{">Q", true, JetDirectionList{}},
		{">", false, JetDirectionList{Right}},
		{">>", false, JetDirectionList{Right, Right}},
		{"<Q", true, JetDirectionList{}},
		{"<", false, JetDirectionList{Left}},
		{"<<", false, JetDirectionList{Left, Left}},
		{"<>", false, JetDirectionList{Left, Right}},
		{"<><><>", false, JetDirectionList{Left, Right, Left, Right, Left, Right}},
	}

	for _, test := range testCases {
		directions, err := ParseJetDirections(test.str)

		if test.expectedErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)

			assert.Equal(t, test.expectedJetDirections, directions)
		}
	}
}

func TestBoundsIntersect(t *testing.T) {
	type testCase struct {
		p1             Point
		s              Size
		p2             Point
		expectedResult bool
	}

	testCases := []testCase{
		{NewPoint(0, 0), NewSize(2, 2), NewPoint(1, 1), true},
		{NewPoint(0, 0), NewSize(2, 2), NewPoint(0, 0), true},
		{NewPoint(0, 0), NewSize(2, 2), NewPoint(2, 2), true},
		{NewPoint(-1, -1), NewSize(2, 2), NewPoint(0, 0), true},
		{NewPoint(-1, -1), NewSize(2, 2), NewPoint(-1, -1), true},
		{NewPoint(-1, -1), NewSize(2, 2), NewPoint(-2, -2), false},
		{NewPoint(-1, -1), NewSize(2, 2), NewPoint(-2, -1), false},
	}

	for _, test := range testCases {
		intersects := BoundsIntersect(test.p1, test.s, test.p2)
		assert.Equal(t, test.expectedResult, intersects)
	}
}

func TestNewBitmap(t *testing.T) {
	type testCase struct {
		str            string
		expectedErr    bool
		expectedBitmap Bitmap
	}

	testCases := []testCase{
		{
			"1234\n1234\n1234", true, Bitmap{},
		},
		{
			"####\n####\n####", false, Bitmap{Size: NewSize(4, 3), Rows: []byte{0xF0, 0xF0, 0xF0}},
		},
		{
			"....\n....\n....", false, Bitmap{Size: NewSize(4, 3), Rows: []byte{0x00, 0x00, 0x00}},
		},
	}

	for _, test := range testCases {
		b, err := NewBitmap(test.str)

		if test.expectedErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)

			assert.Equal(t, test.expectedBitmap, b)
		}
	}
}

func TestBitmap_Describe(t *testing.T) {
	type testCase struct {
		str                 string
		expectedDescription string
	}

	testCases := []testCase{
		{
			"####\n####\n####", "####\n####\n####",
		},
		{
			"....\n....\n....", "....\n....\n....",
		},
	}

	for _, test := range testCases {
		b, err := NewBitmap(test.str)
		assert.NoError(t, err)

		description := b.Describe()
		assert.Equal(t, test.expectedDescription, description)
	}
}

func TestShape_GetBitmap(t *testing.T) {
	type testCase struct {
		newShape            func() Shape
		expectedDescription string
	}

	testCases := []testCase{
		{NewHorizontalLineShape, "####"},
		{NewCrossShape, ".#.\n###\n.#."},
		{NewAngleShape, "..#\n..#\n###"},
		{NewVerticalLineShape, "#\n#\n#\n#"},
		{NewSquareShape, "##\n##"},
	}

	for _, test := range testCases {
		shape := test.newShape()

		description := shape.GetBitmap().Describe()
		assert.Equal(t, test.expectedDescription, description)
	}
}

func TestTower_LockShape(t *testing.T) {
	type testCase struct {
		newShape     func() Shape
		location     Point
		expectedRows []byte
	}

	testCases := []testCase{
		{NewHorizontalLineShape, NewPoint(0, 0), []byte{0xF0}},
		{NewHorizontalLineShape, NewPoint(1, 0), []byte{0x78}},
		{NewHorizontalLineShape, NewPoint(2, 0), []byte{0x3C}},
		{NewHorizontalLineShape, NewPoint(3, 0), []byte{0x1E}},

		{NewHorizontalLineShape, NewPoint(0, 1), []byte{0x00, 0xF0}},
		{NewHorizontalLineShape, NewPoint(1, 1), []byte{0x00, 0x78}},
		{NewHorizontalLineShape, NewPoint(2, 1), []byte{0x00, 0x3C}},
		{NewHorizontalLineShape, NewPoint(3, 1), []byte{0x00, 0x1E}},

		{NewCrossShape, NewPoint(0, 0), []byte{0x40, 0xE0, 0x40}},
		{NewCrossShape, NewPoint(1, 0), []byte{0x20, 0x70, 0x20}},
		{NewCrossShape, NewPoint(2, 0), []byte{0x10, 0x38, 0x10}},
		{NewCrossShape, NewPoint(3, 0), []byte{0x08, 0x1C, 0x08}},
		{NewCrossShape, NewPoint(4, 0), []byte{0x04, 0x0E, 0x04}},

		{NewCrossShape, NewPoint(0, 1), []byte{0x00, 0x40, 0xE0, 0x40}},
		{NewCrossShape, NewPoint(1, 1), []byte{0x00, 0x20, 0x70, 0x20}},
		{NewCrossShape, NewPoint(2, 1), []byte{0x00, 0x10, 0x38, 0x10}},
		{NewCrossShape, NewPoint(3, 1), []byte{0x00, 0x08, 0x1C, 0x08}},
		{NewCrossShape, NewPoint(4, 1), []byte{0x00, 0x04, 0x0E, 0x04}},

		{NewAngleShape, NewPoint(0, 0), []byte{0xE0, 0x20, 0x20}},
		{NewAngleShape, NewPoint(1, 0), []byte{0x70, 0x10, 0x10}},
		{NewAngleShape, NewPoint(2, 0), []byte{0x38, 0x08, 0x08}},
		{NewAngleShape, NewPoint(3, 0), []byte{0x1C, 0x04, 0x04}},
		{NewAngleShape, NewPoint(4, 0), []byte{0x0E, 0x02, 0x02}},

		{NewAngleShape, NewPoint(0, 1), []byte{0x00, 0xE0, 0x20, 0x20}},
		{NewAngleShape, NewPoint(1, 1), []byte{0x00, 0x70, 0x10, 0x10}},
		{NewAngleShape, NewPoint(2, 1), []byte{0x00, 0x38, 0x08, 0x08}},
		{NewAngleShape, NewPoint(3, 1), []byte{0x00, 0x1C, 0x04, 0x04}},
		{NewAngleShape, NewPoint(4, 1), []byte{0x00, 0x0E, 0x02, 0x02}},

		{NewVerticalLineShape, NewPoint(0, 0), []byte{0x80, 0x80, 0x80, 0x80}},
		{NewVerticalLineShape, NewPoint(1, 0), []byte{0x40, 0x40, 0x40, 0x40}},
		{NewVerticalLineShape, NewPoint(2, 0), []byte{0x20, 0x20, 0x20, 0x20}},
		{NewVerticalLineShape, NewPoint(3, 0), []byte{0x10, 0x10, 0x10, 0x10}},
		{NewVerticalLineShape, NewPoint(4, 0), []byte{0x08, 0x08, 0x08, 0x08}},
		{NewVerticalLineShape, NewPoint(5, 0), []byte{0x04, 0x04, 0x04, 0x04}},
		{NewVerticalLineShape, NewPoint(6, 0), []byte{0x02, 0x02, 0x02, 0x02}},

		{NewVerticalLineShape, NewPoint(0, 1), []byte{0x00, 0x80, 0x80, 0x80, 0x80}},
		{NewVerticalLineShape, NewPoint(1, 1), []byte{0x00, 0x40, 0x40, 0x40, 0x40}},
		{NewVerticalLineShape, NewPoint(2, 1), []byte{0x00, 0x20, 0x20, 0x20, 0x20}},
		{NewVerticalLineShape, NewPoint(3, 1), []byte{0x00, 0x10, 0x10, 0x10, 0x10}},
		{NewVerticalLineShape, NewPoint(4, 1), []byte{0x00, 0x08, 0x08, 0x08, 0x08}},
		{NewVerticalLineShape, NewPoint(5, 1), []byte{0x00, 0x04, 0x04, 0x04, 0x04}},
		{NewVerticalLineShape, NewPoint(6, 1), []byte{0x00, 0x02, 0x02, 0x02, 0x02}},

		{NewSquareShape, NewPoint(0, 0), []byte{0xC0, 0xC0}},
		{NewSquareShape, NewPoint(1, 0), []byte{0x60, 0x60}},
		{NewSquareShape, NewPoint(2, 0), []byte{0x30, 0x30}},
		{NewSquareShape, NewPoint(3, 0), []byte{0x18, 0x18}},
		{NewSquareShape, NewPoint(4, 0), []byte{0x0C, 0x0C}},
		{NewSquareShape, NewPoint(5, 0), []byte{0x06, 0x06}},

		{NewSquareShape, NewPoint(0, 1), []byte{0x00, 0xC0, 0xC0}},
		{NewSquareShape, NewPoint(1, 1), []byte{0x00, 0x60, 0x60}},
		{NewSquareShape, NewPoint(2, 1), []byte{0x00, 0x30, 0x30}},
		{NewSquareShape, NewPoint(3, 1), []byte{0x00, 0x18, 0x18}},
		{NewSquareShape, NewPoint(4, 1), []byte{0x00, 0x0C, 0x0C}},
		{NewSquareShape, NewPoint(5, 1), []byte{0x00, 0x06, 0x06}},
	}

	for _, test := range testCases {
		tower := NewTower(7)

		shape := test.newShape()
		shape.SetPosition(test.location)

		tower.LockShape(shape)
		assert.Equal(t, test.expectedRows, tower.Rows)
	}
}

func TestTower_GetHeights(t *testing.T) {
	type testCase struct {
		newShape        func() Shape
		location        Point
		expectedHeights []int64
	}

	testCases := []testCase{
		{NewHorizontalLineShape, NewPoint(0, 0), []int64{1, 1, 1, 1, 0, 0, 0}},
		{NewHorizontalLineShape, NewPoint(1, 0), []int64{0, 1, 1, 1, 1, 0, 0}},
		{NewHorizontalLineShape, NewPoint(2, 0), []int64{0, 0, 1, 1, 1, 1, 0}},
		{NewHorizontalLineShape, NewPoint(3, 0), []int64{0, 0, 0, 1, 1, 1, 1}},

		{NewHorizontalLineShape, NewPoint(0, 1), []int64{2, 2, 2, 2, 0, 0, 0}},
		{NewHorizontalLineShape, NewPoint(1, 1), []int64{0, 2, 2, 2, 2, 0, 0}},
		{NewHorizontalLineShape, NewPoint(2, 1), []int64{0, 0, 2, 2, 2, 2, 0}},
		{NewHorizontalLineShape, NewPoint(3, 1), []int64{0, 0, 0, 2, 2, 2, 2}},

		{NewCrossShape, NewPoint(0, 0), []int64{2, 3, 2, 0, 0, 0, 0}},
		{NewCrossShape, NewPoint(1, 0), []int64{0, 2, 3, 2, 0, 0, 0}},
		{NewCrossShape, NewPoint(2, 0), []int64{0, 0, 2, 3, 2, 0, 0}},
		{NewCrossShape, NewPoint(3, 0), []int64{0, 0, 0, 2, 3, 2, 0}},
		{NewCrossShape, NewPoint(4, 0), []int64{0, 0, 0, 0, 2, 3, 2}},

		{NewCrossShape, NewPoint(0, 1), []int64{3, 4, 3, 0, 0, 0, 0}},
		{NewCrossShape, NewPoint(1, 1), []int64{0, 3, 4, 3, 0, 0, 0}},
		{NewCrossShape, NewPoint(2, 1), []int64{0, 0, 3, 4, 3, 0, 0}},
		{NewCrossShape, NewPoint(3, 1), []int64{0, 0, 0, 3, 4, 3, 0}},
		{NewCrossShape, NewPoint(4, 1), []int64{0, 0, 0, 0, 3, 4, 3}},

		{NewAngleShape, NewPoint(0, 0), []int64{1, 1, 3, 0, 0, 0, 0}},
		{NewAngleShape, NewPoint(1, 0), []int64{0, 1, 1, 3, 0, 0, 0}},
		{NewAngleShape, NewPoint(2, 0), []int64{0, 0, 1, 1, 3, 0, 0}},
		{NewAngleShape, NewPoint(3, 0), []int64{0, 0, 0, 1, 1, 3, 0}},
		{NewAngleShape, NewPoint(4, 0), []int64{0, 0, 0, 0, 1, 1, 3}},

		{NewAngleShape, NewPoint(0, 1), []int64{2, 2, 4, 0, 0, 0, 0}},
		{NewAngleShape, NewPoint(1, 1), []int64{0, 2, 2, 4, 0, 0, 0}},
		{NewAngleShape, NewPoint(2, 1), []int64{0, 0, 2, 2, 4, 0, 0}},
		{NewAngleShape, NewPoint(3, 1), []int64{0, 0, 0, 2, 2, 4, 0}},
		{NewAngleShape, NewPoint(4, 1), []int64{0, 0, 0, 0, 2, 2, 4}},

		{NewVerticalLineShape, NewPoint(0, 0), []int64{4, 0, 0, 0, 0, 0, 0}},
		{NewVerticalLineShape, NewPoint(1, 0), []int64{0, 4, 0, 0, 0, 0, 0}},
		{NewVerticalLineShape, NewPoint(2, 0), []int64{0, 0, 4, 0, 0, 0, 0}},
		{NewVerticalLineShape, NewPoint(3, 0), []int64{0, 0, 0, 4, 0, 0, 0}},
		{NewVerticalLineShape, NewPoint(4, 0), []int64{0, 0, 0, 0, 4, 0, 0}},
		{NewVerticalLineShape, NewPoint(5, 0), []int64{0, 0, 0, 0, 0, 4, 0}},
		{NewVerticalLineShape, NewPoint(6, 0), []int64{0, 0, 0, 0, 0, 0, 4}},

		{NewVerticalLineShape, NewPoint(0, 1), []int64{5, 0, 0, 0, 0, 0, 0}},
		{NewVerticalLineShape, NewPoint(1, 1), []int64{0, 5, 0, 0, 0, 0, 0}},
		{NewVerticalLineShape, NewPoint(2, 1), []int64{0, 0, 5, 0, 0, 0, 0}},
		{NewVerticalLineShape, NewPoint(3, 1), []int64{0, 0, 0, 5, 0, 0, 0}},
		{NewVerticalLineShape, NewPoint(4, 1), []int64{0, 0, 0, 0, 5, 0, 0}},
		{NewVerticalLineShape, NewPoint(5, 1), []int64{0, 0, 0, 0, 0, 5, 0}},
		{NewVerticalLineShape, NewPoint(6, 1), []int64{0, 0, 0, 0, 0, 0, 5}},

		{NewSquareShape, NewPoint(0, 0), []int64{2, 2, 0, 0, 0, 0, 0}},
		{NewSquareShape, NewPoint(1, 0), []int64{0, 2, 2, 0, 0, 0, 0}},
		{NewSquareShape, NewPoint(2, 0), []int64{0, 0, 2, 2, 0, 0, 0}},
		{NewSquareShape, NewPoint(3, 0), []int64{0, 0, 0, 2, 2, 0, 0}},
		{NewSquareShape, NewPoint(4, 0), []int64{0, 0, 0, 0, 2, 2, 0}},
		{NewSquareShape, NewPoint(5, 0), []int64{0, 0, 0, 0, 0, 2, 2}},

		{NewSquareShape, NewPoint(0, 1), []int64{3, 3, 0, 0, 0, 0, 0}},
		{NewSquareShape, NewPoint(1, 1), []int64{0, 3, 3, 0, 0, 0, 0}},
		{NewSquareShape, NewPoint(2, 1), []int64{0, 0, 3, 3, 0, 0, 0}},
		{NewSquareShape, NewPoint(3, 1), []int64{0, 0, 0, 3, 3, 0, 0}},
		{NewSquareShape, NewPoint(4, 1), []int64{0, 0, 0, 0, 3, 3, 0}},
		{NewSquareShape, NewPoint(5, 1), []int64{0, 0, 0, 0, 0, 3, 3}},
	}

	for _, test := range testCases {
		tower := NewTower(7)

		shape := test.newShape()
		shape.SetPosition(test.location)

		tower.LockShape(shape)
		assert.Equal(t, test.expectedHeights, tower.GetHeights())
	}
}

func TestTower_CanShapeMoveToPosition(t *testing.T) {
	type testCase struct {
		shape1         func() Shape
		location1      Point
		shape2         func() Shape
		location2      Point
		expectedResult bool
	}

	testCases := []testCase{
		{NewHorizontalLineShape, NewPoint(0, 0), NewHorizontalLineShape, NewPoint(0, 0), false},
		{NewHorizontalLineShape, NewPoint(0, 0), NewHorizontalLineShape, NewPoint(0, 1), true},
		{NewHorizontalLineShape, NewPoint(0, 0), NewHorizontalLineShape, NewPoint(3, 1), true},

		{NewAngleShape, NewPoint(0, 0), NewVerticalLineShape, NewPoint(1, 1), true},
		{NewAngleShape, NewPoint(0, 0), NewVerticalLineShape, NewPoint(2, 1), false},

		// Test bounds movements.
		{NewAngleShape, NewPoint(0, 0), NewVerticalLineShape, NewPoint(-1, 1), false},
		{NewHorizontalLineShape, NewPoint(0, 0), NewHorizontalLineShape, NewPoint(4, 1), false},
		{NewHorizontalLineShape, NewPoint(0, 0), NewVerticalLineShape, NewPoint(4, 0), true},
	}

	for _, test := range testCases {
		tower := NewTower(7)

		shape1 := test.shape1()
		shape1.SetPosition(test.location1)

		tower.LockShape(shape1)

		shape2 := test.shape2()

		assert.Equal(t, test.expectedResult, tower.CanShapeMoveToPosition(shape2, test.location2))
	}
}

func TestTower_DropShape(t *testing.T) {
	type testCase struct {
		jetDirections  string
		iterations     int64
		expectedHeight int64
	}

	testCases := []testCase{
		{">>><<><>><<<>><>>><<<>>><<<><<<>><>><<>>", 2022, 3068},
	}

	for _, test := range testCases {
		jetDirections, err := ParseJetDirections(test.jetDirections)
		assert.NoError(t, err)

		room := NewRoom(7, jetDirections)

		for i := int64(0); i < test.iterations; i++ {
			room.DropShape()
		}

		assert.Equal(t, test.expectedHeight, room.GetTowerHeight())
	}
}
