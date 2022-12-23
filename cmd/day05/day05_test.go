package day05

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMovementOp(t *testing.T) {
	mo := NewMovementOp(0, 1, 3)

	assert.Equal(t, mo.startStackIndex, 0)
	assert.Equal(t, mo.endStackIndex, 1)
	assert.Equal(t, mo.crateCount, 3)
}

func TestParseMovementOp(t *testing.T) {
	type parseMovementOpTest struct {
		str                string
		expectedErr        bool
		expectedMovementOp MovementOp
	}

	tests := []parseMovementOpTest{
		{"", true, NewMovementOp(0, 0, 0)},
		{"move 5 from 2 to 3", false, NewMovementOp(2, 3, 5)},
		{"move 5 from ! to z", true, NewMovementOp(0, 0, 0)},
	}

	for _, test := range tests {
		mo, err := ParseMovementOp(test.str)
		if test.expectedErr {
			assert.Error(t, err)
		} else {
			assert.Equal(t, test.expectedMovementOp, mo)
		}
	}
}

func TestParseInitialCratesLine(t *testing.T) {
	type parseInitialCratesLineTest struct {
		str                    string
		expectedErr            bool
		expectedCrateLocations []CrateLocation
	}

	tests := []parseInitialCratesLineTest{
		{"", true, []CrateLocation{}},
		{"    [C] [R] [Z]     [R]     [H] [Z]", false, []CrateLocation{
			NewCrateLocation(NewCrate("C"), 2),
			NewCrateLocation(NewCrate("R"), 3),
			NewCrateLocation(NewCrate("Z"), 4),
			NewCrateLocation(NewCrate("R"), 6),
			NewCrateLocation(NewCrate("H"), 8),
			NewCrateLocation(NewCrate("Z"), 9),
		}},
		{"[T] [R] [B] [C] [L] [P] [F] [L] [H]", false, []CrateLocation{
			NewCrateLocation(NewCrate("T"), 1),
			NewCrateLocation(NewCrate("R"), 2),
			NewCrateLocation(NewCrate("B"), 3),
			NewCrateLocation(NewCrate("C"), 4),
			NewCrateLocation(NewCrate("L"), 5),
			NewCrateLocation(NewCrate("P"), 6),
			NewCrateLocation(NewCrate("F"), 7),
			NewCrateLocation(NewCrate("L"), 8),
			NewCrateLocation(NewCrate("H"), 9),
		}},
		{"(T] [R] [B] [C] [L] [P] [F] [L] [H]", true, []CrateLocation{}},
		{"[T]_[R] [B] [C] [L] [P] [F] [L] [H]", true, []CrateLocation{}},
		{"[T] [RR] [B] [C] [L] [P] [F] [L] [H]", true, []CrateLocation{}},
	}

	for _, test := range tests {
		cl, err := ParseInitialCratesLine(test.str)
		if test.expectedErr {
			assert.Error(t, err)
		} else {
			assert.Equal(t, test.expectedCrateLocations, cl)
		}
	}
}

func createLabeledCrates(labels []string) CrateStack {
	stack := make(CrateStack, 0)
	for _, label := range labels {
		crate := NewCrate(label)
		stack = append(stack, crate)
	}

	return stack
}

func TestWarehouseAddCrate(t *testing.T) {
	// initialCrates := []string{
	// 	"        [H]         [S]         [D]",
	// 	"    [S] [C]         [C]     [Q] [L]",
	// 	"    [C] [R] [Z]     [R]     [H] [Z]",
	// 	"    [G] [N] [H] [S] [B]     [R] [F]",
	// 	"[D] [T] [Q] [F] [Q] [Z]     [Z] [N]",
	// 	"[Z] [W] [F] [N] [F] [W] [J] [V] [G]",
	// 	"[T] [R] [B] [C] [L] [P] [F] [L] [H]",
	// 	"[H] [Q] [P] [L] [G] [V] [Z] [D] [B]",
	// }

	type addCrateTest struct {
		crates         []CrateLocation
		expectedCrates []CrateStack
	}

	tests := []addCrateTest{
		{
			[]CrateLocation{NewCrateLocation(NewCrate("C"), 2), NewCrateLocation(NewCrate("D"), 5)},
			[]CrateStack{make(CrateStack, 0), createLabeledCrates([]string{"C"}), make(CrateStack, 0), make(CrateStack, 0), createLabeledCrates([]string{"D"})},
		},
		{
			[]CrateLocation{NewCrateLocation(NewCrate("C"), 2), NewCrateLocation(NewCrate("D"), 2)},
			[]CrateStack{make(CrateStack, 0), createLabeledCrates([]string{"C", "D"})},
		},
	}

	for _, test := range tests {
		w := NewWarehouse()
		for _, cl := range test.crates {
			w.AddCrate(cl.crate, cl.stackIndex)
		}

		assert.Equal(t, test.expectedCrates, w.crates)
	}
}

func TestWarehouseApplyMovementOp(t *testing.T) {
	initialCrates := []string{
		"        [H]         [S]         [D]",
		"    [S] [C]         [C]     [Q] [L]",
		"    [C] [R] [Z]     [R]     [H] [Z]",
		"    [G] [N] [H] [S] [B]     [R] [F]",
		"[D] [T] [Q] [F] [Q] [Z]     [Z] [N]",
		"[Z] [W] [F] [N] [F] [W] [J] [V] [G]",
		"[T] [R] [B] [C] [L] [P] [F] [L] [H]",
		"[H] [Q] [P] [L] [G] [V] [Z] [D] [B]",
	}

	type applyMovementOpTest struct {
		movementOps    []MovementOp
		expectedCrates []CrateStack
	}

	tests := []applyMovementOpTest{
		{
			[]MovementOp{NewMovementOp(1, 2, 1)},
			[]CrateStack{
				createLabeledCrates([]string{"Z", "T", "H"}),
				createLabeledCrates([]string{"D", "S", "C", "G", "T", "W", "R", "Q"}),
				createLabeledCrates([]string{"H", "C", "R", "N", "Q", "F", "B", "P"}),
				createLabeledCrates([]string{"Z", "H", "F", "N", "C", "L"}),
				createLabeledCrates([]string{"S", "Q", "F", "L", "G"}),
				createLabeledCrates([]string{"S", "C", "R", "B", "Z", "W", "P", "V"}),
				createLabeledCrates([]string{"J", "F", "Z"}),
				createLabeledCrates([]string{"Q", "H", "R", "Z", "V", "L", "D"}),
				createLabeledCrates([]string{"D", "L", "Z", "F", "N", "G", "H", "B"})},
		},
		{
			[]MovementOp{NewMovementOp(1, 2, 1), NewMovementOp(2, 3, 2)},
			[]CrateStack{
				createLabeledCrates([]string{"Z", "T", "H"}),
				createLabeledCrates([]string{"C", "G", "T", "W", "R", "Q"}),
				createLabeledCrates([]string{"S", "D", "H", "C", "R", "N", "Q", "F", "B", "P"}),
				createLabeledCrates([]string{"Z", "H", "F", "N", "C", "L"}),
				createLabeledCrates([]string{"S", "Q", "F", "L", "G"}),
				createLabeledCrates([]string{"S", "C", "R", "B", "Z", "W", "P", "V"}),
				createLabeledCrates([]string{"J", "F", "Z"}),
				createLabeledCrates([]string{"Q", "H", "R", "Z", "V", "L", "D"}),
				createLabeledCrates([]string{"D", "L", "Z", "F", "N", "G", "H", "B"})},
		},
	}

	for _, test := range tests {
		w := NewWarehouse()

		for _, line := range initialCrates {
			cls, err := ParseInitialCratesLine(line)
			assert.NoError(t, err)
			for _, cl := range cls {
				w.AddCrate(cl.crate, cl.stackIndex)
			}
		}

		for _, mo := range test.movementOps {
			w.ApplyMovementOp(mo)
		}

		assert.Equal(t, test.expectedCrates, w.crates)
	}
}

func TestWarehouseApplyMovementOp9001(t *testing.T) {
	initialCrates := []string{
		"        [H]         [S]         [D]",
		"    [S] [C]         [C]     [Q] [L]",
		"    [C] [R] [Z]     [R]     [H] [Z]",
		"    [G] [N] [H] [S] [B]     [R] [F]",
		"[D] [T] [Q] [F] [Q] [Z]     [Z] [N]",
		"[Z] [W] [F] [N] [F] [W] [J] [V] [G]",
		"[T] [R] [B] [C] [L] [P] [F] [L] [H]",
		"[H] [Q] [P] [L] [G] [V] [Z] [D] [B]",
	}

	type applyMovementOpTest struct {
		movementOps    []MovementOp
		expectedCrates []CrateStack
	}

	tests := []applyMovementOpTest{
		{
			[]MovementOp{NewMovementOp(1, 2, 1)},
			[]CrateStack{
				createLabeledCrates([]string{"Z", "T", "H"}),
				createLabeledCrates([]string{"D", "S", "C", "G", "T", "W", "R", "Q"}),
				createLabeledCrates([]string{"H", "C", "R", "N", "Q", "F", "B", "P"}),
				createLabeledCrates([]string{"Z", "H", "F", "N", "C", "L"}),
				createLabeledCrates([]string{"S", "Q", "F", "L", "G"}),
				createLabeledCrates([]string{"S", "C", "R", "B", "Z", "W", "P", "V"}),
				createLabeledCrates([]string{"J", "F", "Z"}),
				createLabeledCrates([]string{"Q", "H", "R", "Z", "V", "L", "D"}),
				createLabeledCrates([]string{"D", "L", "Z", "F", "N", "G", "H", "B"})},
		},
		{
			[]MovementOp{NewMovementOp(1, 2, 1), NewMovementOp(2, 3, 2)},
			[]CrateStack{
				createLabeledCrates([]string{"Z", "T", "H"}),
				createLabeledCrates([]string{"C", "G", "T", "W", "R", "Q"}),
				createLabeledCrates([]string{"D", "S", "H", "C", "R", "N", "Q", "F", "B", "P"}),
				createLabeledCrates([]string{"Z", "H", "F", "N", "C", "L"}),
				createLabeledCrates([]string{"S", "Q", "F", "L", "G"}),
				createLabeledCrates([]string{"S", "C", "R", "B", "Z", "W", "P", "V"}),
				createLabeledCrates([]string{"J", "F", "Z"}),
				createLabeledCrates([]string{"Q", "H", "R", "Z", "V", "L", "D"}),
				createLabeledCrates([]string{"D", "L", "Z", "F", "N", "G", "H", "B"})},
		},
	}

	for _, test := range tests {
		w := NewWarehouse()

		for _, line := range initialCrates {
			cls, err := ParseInitialCratesLine(line)
			assert.NoError(t, err)
			for _, cl := range cls {
				w.AddCrate(cl.crate, cl.stackIndex)
			}
		}

		for _, mo := range test.movementOps {
			w.ApplyMovementOp9001(mo)
		}

		assert.Equal(t, test.expectedCrates, w.crates)
	}
}
