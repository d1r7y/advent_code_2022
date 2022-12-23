package day09

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

func GetMovementAmount(dir MovementDirection) (int, int) {
	switch dir {
	case UpDirection:
		return 0, 1
	case DownDirection:
		return 0, -1
	case RightDirection:
		return 1, 0
	case LeftDirection:
		return -1, 0
	default:
		panic("unexpected direction")
	}
}

func Distance(Head KnotPosition, Tail KnotPosition) float64 {
	return math.Sqrt(math.Pow(float64(Head.X-Tail.X), 2.0) + math.Pow(float64(Head.Y-Tail.Y), 2.0))
}

func MustMoveKnot(head KnotPosition, knot KnotPosition) bool {
	if head == knot {
		return false
	}

	if Distance(head, knot) > math.Sqrt2 {
		return true
	}

	return false
}

func GetNewKnotPosition(head KnotPosition, knot KnotPosition) KnotPosition {
	// If head is ever two steps directly up, down, left, or right from
	// the tail, knot must move one step in that direction.

	if head.X == knot.X {
		if head.Y > knot.Y {
			knot.Y++
		} else {
			knot.Y--
		}

		return knot
	}
	if head.Y == knot.Y {
		if head.X > knot.X {
			knot.X++
		} else {
			knot.X--
		}

		return knot
	}

	// Otherwise, if the head and knot aren't in the same row or column,
	// knot always moves one step diagonally to keep up.
	if head.X > knot.X {
		knot.X++
	} else {
		knot.X--
	}

	if head.Y > knot.Y {
		knot.Y++
	} else {
		knot.Y--
	}

	return knot
}

type KnotPosition struct {
	X int
	Y int
}

func NewKnotPosition(x, y int) KnotPosition {
	return KnotPosition{X: x, Y: y}
}

type World struct {
	MaxX int
	MaxY int
	MinX int
	MinY int

	Head           KnotPosition
	RemainingKnots []KnotPosition

	TailPositions map[KnotPosition]bool
}

func NewWorld(totalKnots int) *World {
	w := &World{TailPositions: make(map[KnotPosition]bool), RemainingKnots: make([]KnotPosition, totalKnots-1)}

	// Save away the initial Tail position
	w.TailPositions[w.GetTail()] = true

	return w
}

func (w *World) GetTail() KnotPosition {
	return w.RemainingKnots[len(w.RemainingKnots)-1]
}

func (w *World) TrackMinMaxPosition(kp KnotPosition) {
	if kp.X > w.MaxX {
		w.MaxX = kp.X
	}

	if kp.Y > w.MaxY {
		w.MaxY = kp.Y
	}

	if kp.X < w.MinX {
		w.MinX = kp.X
	}

	if kp.Y < w.MinY {
		w.MinY = kp.Y
	}
}

func (w *World) ApplyMovementOp(mo KnotMovementOp) {
	for i := mo.Count; i > 0; i-- {
		x, y := GetMovementAmount(mo.Direction)

		w.Head.X += x
		w.Head.Y += y

		w.TrackMinMaxPosition(w.Head)

		previousKP := w.Head

		for j, kp := range w.RemainingKnots {
			if !MustMoveKnot(previousKP, kp) {
				previousKP = kp
				continue
			}

			// Now adjust tail in response to head's movement.
			newKP := GetNewKnotPosition(previousKP, kp)
			w.RemainingKnots[j] = newKP

			w.TrackMinMaxPosition(newKP)

			// If this is the tail knot, remember its position.
			if j == len(w.RemainingKnots)-1 {
				w.TailPositions[newKP] = true
			}

			previousKP = newKP
		}
	}
}

func (w *World) GetTailPositions() []KnotPosition {
	keys := make([]KnotPosition, 0)

	for key := range w.TailPositions {
		keys = append(keys, key)
	}

	return keys
}

type MovementDirection int

const (
	UpDirection MovementDirection = iota
	DownDirection
	LeftDirection
	RightDirection
)

type KnotMovementOp struct {
	Direction MovementDirection
	Count     int
}

func NewKnotMovementOp(direction MovementDirection, count int) KnotMovementOp {
	return KnotMovementOp{Direction: direction, Count: count}
}

type KnotMovementOpList []KnotMovementOp

func ParseKnotMovementOp(line string) (KnotMovementOp, error) {
	var dir byte
	var count int

	c, err := fmt.Sscanf(line, "%c %d", &dir, &count)
	if err != nil {
		return KnotMovementOp{}, err
	}
	if c != 2 {
		return KnotMovementOp{}, errors.New("invalid line")
	}
	if count <= 0 {
		return KnotMovementOp{}, errors.New("invalid movement count")
	}

	var direction MovementDirection

	switch dir {
	case 'U':
		direction = UpDirection
	case 'D':
		direction = DownDirection
	case 'L':
		direction = LeftDirection
	case 'R':
		direction = RightDirection
	default:
		return KnotMovementOp{}, errors.New("invalid direction")
	}

	kmo := NewKnotMovementOp(direction, count)
	return kmo, nil
}

func ParseKnotMovementOps(lines []string) (KnotMovementOpList, error) {
	list := make(KnotMovementOpList, 0)

	for _, line := range lines {
		kmo, err := ParseKnotMovementOp(line)
		if err != nil {
			return KnotMovementOpList{}, err
		}

		list = append(list, kmo)
	}

	return list, nil
}

func day09(fileContents string) error {
	// Scan the head knot movement operations in.
	headMovementOperations, err := ParseKnotMovementOps(strings.Split(string(fileContents), "\n"))
	if err != nil {
		return err
	}

	// Part 1: After running through all the head knot movement operations, how many positions does the tail knot visit
	// at least once?
	w2 := NewWorld(2)

	for _, movementOp := range headMovementOperations {
		w2.ApplyMovementOp(movementOp)
	}

	fmt.Printf("Dynamic board size (%d,%d,%d,%d)\n", w2.MinX, w2.MinY, w2.MaxX, w2.MaxY)

	fmt.Printf("Tail knot visited %d positions\n", len(w2.GetTailPositions()))

	// Part 2: What if there are 10 knots?  How many positions does the final tail knot visit at least once?
	w10 := NewWorld(10)

	for _, movementOp := range headMovementOperations {
		w10.ApplyMovementOp(movementOp)
	}

	fmt.Printf("Dynamic board size (%d,%d,%d,%d)\n", w10.MinX, w10.MinY, w10.MaxX, w10.MaxY)

	fmt.Printf("Tail knot visited %d positions\n", len(w10.GetTailPositions()))

	return nil
}
