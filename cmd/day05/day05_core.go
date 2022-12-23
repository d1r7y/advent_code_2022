package day05

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

type Crate struct {
	label string
}

func NewCrate(label string) Crate {
	return Crate{label: label}
}

type CrateStack []Crate

type CrateLocation struct {
	crate      Crate
	stackIndex int
}

type Warehouse struct {
	crates []CrateStack
}

func NewWarehouse() *Warehouse {
	return &Warehouse{crates: make([]CrateStack, 0)}
}

func (w *Warehouse) AddCrate(c Crate, si int) {
	missingStacks := si - len(w.crates)
	for i := 0; i < missingStacks; i++ {
		w.crates = append(w.crates, make(CrateStack, 0))
	}

	w.crates[si-1] = append(w.crates[si-1], c)
}

func (w *Warehouse) ApplyMovementOp(mo MovementOp) {
	if mo.startStackIndex > len(w.crates) {
		log.Fatalf("invalid start stack index %d vs %d\n", mo.startStackIndex, len(w.crates))
	}
	if mo.endStackIndex > len(w.crates) {
		log.Fatalf("invalid end stack index %d vs %d\n", mo.endStackIndex, len(w.crates))
	}
	if mo.crateCount <= 0 {
		log.Fatalf("invalid crate count %d\n", mo.crateCount)
	}

	for i := 0; i < mo.crateCount; i++ {
		// Get top of starting crate stack.
		c := w.crates[mo.startStackIndex-1][0]
		// Remove top of starting crate stack.
		w.crates[mo.startStackIndex-1] = w.crates[mo.startStackIndex-1][1:]
		// Add crate to top of new stack.
		w.crates[mo.endStackIndex-1] = append(CrateStack{c}, w.crates[mo.endStackIndex-1]...)
	}
}

func (w *Warehouse) ApplyMovementOp9001(mo MovementOp) {
	if mo.startStackIndex > len(w.crates) {
		log.Fatalf("invalid start stack index %d vs %d\n", mo.startStackIndex, len(w.crates))
	}
	if mo.endStackIndex > len(w.crates) {
		log.Fatalf("invalid end stack index %d vs %d\n", mo.endStackIndex, len(w.crates))
	}
	if mo.crateCount <= 0 {
		log.Fatalf("invalid crate count %d\n", mo.crateCount)
	}

	crane9001 := make(CrateStack, 0)

	for i := 0; i < mo.crateCount; i++ {
		// Get top of starting crate stack.
		c := w.crates[mo.startStackIndex-1][0]
		// Remove top of starting crate stack.
		w.crates[mo.startStackIndex-1] = w.crates[mo.startStackIndex-1][1:]
		// Hold the crate in the CrateMover 9001.  Crates are stored FIFO.
		crane9001 = append(crane9001, c)
	}

	// Add crate to top of new stack.
	w.crates[mo.endStackIndex-1] = append(crane9001, w.crates[mo.endStackIndex-1]...)
}

func (w *Warehouse) Describe() {
	var highestStackHeight = 0

	for _, stack := range w.crates {
		if len(stack) > highestStackHeight {
			highestStackHeight = len(stack)
		}
	}

	for i := highestStackHeight; i > 0; i-- {
		for _, stack := range w.crates {
			if len(stack) < i {
				fmt.Printf("    ")
			} else {
				fmt.Printf("[%s] ", stack[len(stack)-i].label)
			}
		}
		fmt.Print("\n")
	}

	for i := range w.crates {
		fmt.Printf(" %d  ", i+1)
	}
	fmt.Print("\n")
}

func NewCrateLocation(crate Crate, stackIndex int) CrateLocation {
	return CrateLocation{crate: crate, stackIndex: stackIndex}
}

type MovementOp struct {
	startStackIndex int
	endStackIndex   int
	crateCount      int
}

func NewMovementOp(si, ei, c int) MovementOp {
	return MovementOp{startStackIndex: si, endStackIndex: ei, crateCount: c}
}

func ParseMovementOp(str string) (MovementOp, error) {
	if str == "" {
		return MovementOp{}, errors.New("movementop: empty string")
	}

	var c int
	var s int
	var e int

	count, err := fmt.Sscanf(str, "move %d from %d to %d", &c, &s, &e)
	if err != nil {
		return MovementOp{}, err
	}
	if count != 3 {
		return MovementOp{}, errors.New("invalid string")
	}

	return NewMovementOp(s, e, c), nil
}

func ParseInitialCratesLine(str string) ([]CrateLocation, error) {
	crateLocations := make([]CrateLocation, 0)

	if str == "" {
		return crateLocations, errors.New("initialcratesline: empty string")
	}

	for i, c := range str {
		// Sanity check the line.
		if i%4 == 0 && (c != ' ' && c != '[') {
			return []CrateLocation{}, fmt.Errorf("invalid line.  unexpected char '%c' at index %d", c, i)
		}
		if i%4 == 1 && c != ' ' {
			stackIndex := (i / 4) + 1
			crate := NewCrate(fmt.Sprintf("%c", c))
			crateLocations = append(crateLocations, NewCrateLocation(crate, stackIndex))
		}
		if i%4 == 2 && (c != ' ' && c != ']') {
			return []CrateLocation{}, fmt.Errorf("invalid line.  unexpected char '%c' at index %d", c, i)
		}
		if i%4 == 3 && c != ' ' {
			return []CrateLocation{}, fmt.Errorf("invalid line.  unexpected char '%c' at index %d", c, i)
		}
	}

	return crateLocations, nil
}

func IsCrateStackLegendLine(line string) bool {
	if line[0] == ' ' && line[1] == '1' && line[2] == ' ' {
		return true
	}

	return false
}

func day05(fileContents string) error {
	var inInitialStackMode = true
	var inMovementOpsMode = false

	movementOps := make([]MovementOp, 0)

	// Part 1: After executing the movement operations for the initial crate stacks, what are the labels of the crates on
	// top of each stack?

	warehouse1 := NewWarehouse()

	for _, line := range strings.Split(string(fileContents), "\n") {
		if inInitialStackMode && line == "" {
			inInitialStackMode = false
			inMovementOpsMode = true
			continue
		}
		if inInitialStackMode {
			if IsCrateStackLegendLine(line) {
				continue
			}

			crateLocations, err := ParseInitialCratesLine(line)
			if err != nil {
				return err
			}

			for _, cl := range crateLocations {
				warehouse1.AddCrate(cl.crate, cl.stackIndex)
			}
		} else if inMovementOpsMode {
			mo, err := ParseMovementOp(line)
			if err != nil {
				return err
			}

			movementOps = append(movementOps, mo)
		}
	}

	for _, mo := range movementOps {
		warehouse1.ApplyMovementOp(mo)
	}

	warehouse1.Describe()

	// Part 2: If multiple crates are moved in a single movement op, their order is kept.  Now what are the labels of the crates on
	// top of each stack?

	inInitialStackMode = true

	warehouse2 := NewWarehouse()

	for _, line := range strings.Split(string(fileContents), "\n") {
		if inInitialStackMode && line == "" {
			break
		}
		if inInitialStackMode {
			if IsCrateStackLegendLine(line) {
				continue
			}

			crateLocations, err := ParseInitialCratesLine(line)
			if err != nil {
				return err
			}
			for _, cl := range crateLocations {
				warehouse2.AddCrate(cl.crate, cl.stackIndex)
			}
		}
	}

	for _, mo := range movementOps {
		warehouse2.ApplyMovementOp9001(mo)
	}

	warehouse2.Describe()
	return nil
}
