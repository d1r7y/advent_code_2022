package day18

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCube(t *testing.T) {
	type testCase struct {
		line         string
		expectedCube *Cube
	}

	testCases := []testCase{
		{
			line:         "2,2,2",
			expectedCube: &Cube{Position: Point{2, 2, 2}, FacesExposed: 0},
		},
		{
			line:         "2,1,2",
			expectedCube: &Cube{Position: Point{2, 1, 2}, FacesExposed: 0},
		},
		{
			line:         "0,0,0",
			expectedCube: &Cube{Position: Point{0, 0, 0}, FacesExposed: 0},
		},
		{
			line:         "-1,-2,-3",
			expectedCube: &Cube{Position: Point{-1, -2, -3}, FacesExposed: 0},
		},
	}

	for _, test := range testCases {
		assert.Equal(t, test.expectedCube, ParseCube(test.line))
	}
}

func TestGetSurfaceArea(t *testing.T) {
	str := `2,2,2
1,2,2
3,2,2
2,1,2
2,3,2
2,2,1
2,2,3
2,2,4
2,2,6
1,2,5
3,2,5
2,1,5
2,3,5`

	g := ParseCubes(str)

	assert.Equal(t, 64, g.GetSurfaceArea())
}

func TestGetExteriorSurfaceArea(t *testing.T) {
	str := `2,2,2
1,2,2
3,2,2
2,1,2
2,3,2
2,2,1
2,2,3
2,2,4
2,2,6
1,2,5
3,2,5
2,1,5
2,3,5`

	g := ParseCubes(str)

	assert.Equal(t, 58, g.GetExternalSurfaceArea())
}
