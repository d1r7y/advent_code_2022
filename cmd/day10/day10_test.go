package day10

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCPU(t *testing.T) {
	cpu := NewCPU()

	assert.Equal(t, CPU{X: 1, Cycle: 1, CRTPosition: 0, CRTScreen: NewNullOutput(), XSampleBuffer: make([]Sample, 0)}, *cpu)
}

func TestParseInstruction(t *testing.T) {
	type testCase struct {
		str                 string
		expectedErr         bool
		expectedInstruction Instruction
	}

	testCases := []testCase{
		{"", true, nil},
		{"noop", false, NewNoop()},
		{"addx 0", false, NewAddx(0)},
		{"addx 100", false, NewAddx(100)},
		{"addx -100", false, NewAddx(-100)},
		{"addy 0", true, nil},
		{"addx a", true, nil},
		{"100 addx", true, nil},
	}

	for _, test := range testCases {
		instruction, err := ParseInstruction(test.str)

		if test.expectedErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedInstruction, instruction)
		}
	}
}

func TestValidateInstruction(t *testing.T) {
	type testCase struct {
		instruction        Instruction
		expectedDescribe   string
		expectedCycleCount int
	}

	testCases := []testCase{
		{NewNoop(), "noop", 1},
		{NewAddx(0), "addx 0", 2},
		{NewAddx(-10), "addx -10", 2},
	}

	for _, test := range testCases {
		assert.Equal(t, test.expectedDescribe, test.instruction.Describe())
		assert.Equal(t, test.expectedCycleCount, test.instruction.CycleCount())
	}
}

func TestRunInstructions(t *testing.T) {
	type testCase struct {
		instructions  []Instruction
		expectedX     int
		expectedCycle int
	}

	testCases := []testCase{
		{[]Instruction{NewNoop()}, 1, 2},
		{[]Instruction{NewAddx(0)}, 1, 3},
		{[]Instruction{NewAddx(-10)}, -9, 3},
		{[]Instruction{NewAddx(5)}, 6, 3},
		{[]Instruction{NewAddx(5), NewNoop(), NewAddx(-2)}, 4, 6},
	}

	for _, test := range testCases {
		c := NewCPU()

		for _, i := range test.instructions {
			c.RunInstruction(i)
		}

		assert.Equal(t, test.expectedX, c.X)
		assert.Equal(t, test.expectedCycle, c.Cycle)
	}
}

func TestSampleX(t *testing.T) {
	type testCase struct {
		program                     string
		sampleCycles                []int
		expectedSamples             []Sample
		expectedTotalSignalStrength int
	}

	testCases := []testCase{
		{
			`addx 15
addx -11
addx 6
addx -3
addx 5
addx -1
addx -8
addx 13
addx 4
noop
addx -1
addx 5
addx -1
addx 5
addx -1
addx 5
addx -1
addx 5
addx -1
addx -35
addx 1
addx 24
addx -19
addx 1
addx 16
addx -11
noop
noop
addx 21
addx -15
noop
noop
addx -3
addx 9
addx 1
addx -3
addx 8
addx 1
addx 5
noop
noop
noop
noop
noop
addx -36
noop
addx 1
addx 7
noop
noop
noop
addx 2
addx 6
noop
noop
noop
noop
noop
addx 1
noop
noop
addx 7
addx 1
noop
addx -13
addx 13
addx 7
noop
addx 1
addx -33
noop
noop
noop
addx 2
noop
noop
noop
addx 8
noop
addx -1
addx 2
addx 1
noop
addx 17
addx -9
addx 1
addx 1
addx -3
addx 11
noop
noop
addx 1
noop
addx 1
noop
noop
addx -13
addx -19
addx 1
addx 3
addx 26
addx -30
addx 12
addx -1
addx 3
addx 1
noop
noop
noop
addx -9
addx 18
addx 1
addx 2
noop
noop
addx 9
noop
noop
noop
addx -1
addx 2
addx -37
addx 1
addx 3
noop
addx 15
addx -21
addx 22
addx -6
addx 1
noop
addx 2
addx 1
noop
addx -10
noop
noop
addx 20
addx 1
addx 2
addx 2
addx -6
addx -11
noop
noop
noop`,
			[]int{20, 60, 100, 140, 180, 220},
			[]Sample{NewSample(20, 21), NewSample(60, 19), NewSample(100, 18), NewSample(140, 21), NewSample(180, 16), NewSample(220, 18)},
			13140,
		},
	}

	for _, test := range testCases {
		c := NewCPU()

		c.SetXSampleCycles(test.sampleCycles)

		instructions, err := ParseInstructions(strings.Split(test.program, "\n"))
		assert.NoError(t, err)

		for _, i := range instructions {
			c.RunInstruction(i)
		}

		samples := c.GetXSampleBuffer()

		assert.Equal(t, test.expectedSamples, samples)

		totalSignalStrength := 0

		for _, sample := range samples {
			totalSignalStrength += sample.Cycle * sample.Value
		}

		assert.Equal(t, test.expectedTotalSignalStrength, totalSignalStrength)
	}
}

func TestShouldDrawSprite(t *testing.T) {
	type testCase struct {
		spritePosition int
		beamPosition   int
		expectedResult bool
	}

	testCases := []testCase{
		{5, 5, true},
		{5, 4, true},
		{5, 6, true},
		{5, 7, false},
		{5, 1, false},
		{5, 10, false},
		{5, 3, false},
		{1, 0, true},
	}

	for _, test := range testCases {
		assert.Equal(t, test.expectedResult, ShouldDrawSprite(test.spritePosition, test.beamPosition))
	}
}

func TestOutput(t *testing.T) {
	type testCase struct {
		program string
		output  string
	}

	testCases := []testCase{
		{
			`addx 15
addx -11
addx 6
addx -3
addx 5
addx -1
addx -8
addx 13
addx 4
noop
addx -1
addx 5
addx -1
addx 5
addx -1
addx 5
addx -1
addx 5
addx -1
addx -35
addx 1
addx 24
addx -19
addx 1
addx 16
addx -11
noop
noop
addx 21
addx -15
noop
noop
addx -3
addx 9
addx 1
addx -3
addx 8
addx 1
addx 5
noop
noop
noop
noop
noop
addx -36
noop
addx 1
addx 7
noop
noop
noop
addx 2
addx 6
noop
noop
noop
noop
noop
addx 1
noop
noop
addx 7
addx 1
noop
addx -13
addx 13
addx 7
noop
addx 1
addx -33
noop
noop
noop
addx 2
noop
noop
noop
addx 8
noop
addx -1
addx 2
addx 1
noop
addx 17
addx -9
addx 1
addx 1
addx -3
addx 11
noop
noop
addx 1
noop
addx 1
noop
noop
addx -13
addx -19
addx 1
addx 3
addx 26
addx -30
addx 12
addx -1
addx 3
addx 1
noop
noop
noop
addx -9
addx 18
addx 1
addx 2
noop
noop
addx 9
noop
noop
noop
addx -1
addx 2
addx -37
addx 1
addx 3
noop
addx 15
addx -21
addx 22
addx -6
addx 1
noop
addx 2
addx 1
noop
addx -10
noop
noop
addx 20
addx 1
addx 2
addx 2
addx -6
addx -11
noop
noop
noop`,
			`##..##..##..##..##..##..##..##..##..##..
###...###...###...###...###...###...###.
####....####....####....####....####....
#####.....#####.....#####.....#####.....
######......######......######......####
#######.......#######.......#######.....`,
		},
	}

	for _, test := range testCases {
		c := NewCPU()
		o := NewBufferedOutput()

		c.SetOutput(o)

		instructions, err := ParseInstructions(strings.Split(test.program, "\n"))
		assert.NoError(t, err)

		for _, i := range instructions {
			c.RunInstruction(i)
		}

		screen := ""

		for i, line := range o.Screen {
			screen += string(line)

			if i != len(o.Screen)-1 {
				screen += "\n"
			}
		}

		assert.Equal(t, test.output, screen)
	}
}
