package day17

import (
	"errors"
	"fmt"
	"sort"
)

type Column []bool

type Tower struct {
	Width  int
	Rows   []byte
	Shapes []Shape
}

func NewTower(width int) *Tower {
	return &Tower{
		Width:  width,
		Rows:   make([]byte, 0),
		Shapes: make([]Shape, 0),
	}
}

func (t *Tower) GetHeights() []int64 {
	heights := make([]int64, t.Width)

	for i, row := range t.Rows {
		mask := byte(0x80)

		for j := 0; j < t.Width; j++ {
			if (row & mask) != 0 {
				heights[j] = int64(i + 1)
			}

			mask >>= 1
		}
	}

	return heights
}

func (t *Tower) AddEmptyRows(count int) {
	for i := 0; i < count; i++ {
		t.Rows = append(t.Rows, 0)
	}
}

func (t *Tower) CanShapeMoveToPosition(shape Shape, p Point) bool {
	bitmap := shape.GetBitmap()
	s := shape.GetSize()

	// Do we need to add rows?  Room normally handles this, but unit tests call
	// CanShapeMoveToPosition() directly...
	if p.Y+int64(s.Height) > int64(len(t.Rows)) {
		t.AddEmptyRows(int(p.Y + int64(s.Height) - int64(len(t.Rows))))
	}

	// Check against the bounds of the room.
	if p.X < 0 || int(p.X)+s.Width > t.Width {
		return false
	}

	if p.Y < 0 {
		return false
	}

	// See if shape will intersect with any existing locked shapes.
	for y := 0; y < s.Height; y++ {
		if (t.Rows[p.Y+int64(y)] & (bitmap.Rows[y] >> byte(p.X))) != 0 {
			return false
		}
	}

	return true
}

func (t *Tower) LockShape(shape Shape) {
	// Shape has stopped moving.  Using its position and bitmap, update the individual columns in
	// the tower.

	bitmap := shape.GetBitmap()
	p := shape.GetPosition()
	s := shape.GetSize()

	// Do we need to add rows?  Room normally handles this, but unit tests call LockShape() directly...
	if p.Y+int64(s.Height) > int64(len(t.Rows)) {
		t.AddEmptyRows(int(p.Y + int64(s.Height) - int64(len(t.Rows))))
	}

	for y := 0; y < s.Height; y++ {
		t.Rows[p.Y+int64(y)] |= bitmap.Rows[y] >> byte(p.X)
	}
}

type Room struct {
	NextShape        int
	NextJetDirection int
	JetDirections    JetDirectionList
	Tower            *Tower
}

func NewRoom(width int, jetDirections JetDirectionList) *Room {
	return &Room{JetDirections: jetDirections, Tower: NewTower(width)}
}

func (r *Room) MakeNextShape() Shape {
	shapes := []func() Shape{NewHorizontalLineShape, NewCrossShape, NewAngleShape, NewVerticalLineShape, NewSquareShape}

	newShape := shapes[r.NextShape]()

	r.NextShape = (r.NextShape + 1) % len(shapes)

	return newShape
}

func (r *Room) GetNextJetDirection() JetDirection {
	direction := r.JetDirections[r.NextJetDirection]

	r.NextJetDirection = (r.NextJetDirection + 1) % len(r.JetDirections)

	return direction
}

type Int64Slice []int64

func (x Int64Slice) Len() int           { return len(x) }
func (x Int64Slice) Less(i, j int) bool { return x[i] < x[j] }
func (x Int64Slice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

func (r *Room) GetTowerHeight() int64 {
	heights := r.Tower.GetHeights()
	sort.Sort(sort.Reverse(Int64Slice(heights)))

	return heights[0]
}

func (r *Room) DropShape() {
	shape := r.MakeNextShape()

	// Position shape 2 units away from the left edge and 3 units above the highest shape in the room
	// or the floor.
	towerHeight := r.GetTowerHeight()
	shapePosition := NewPoint(2, towerHeight+3)

	// Add new rows to tower to account for location of newly created shape.
	r.Tower.AddEmptyRows(3 + shape.GetSize().Height)

	shape.SetPosition(shapePosition)

	for {
		// Handle jet direction.
		direction := r.GetNextJetDirection()
		newPosition := shape.GetPosition()

		if direction == Left {
			newPosition.X--
		} else {
			newPosition.X++
		}

		if r.Tower.CanShapeMoveToPosition(shape, newPosition) {
			shape.SetPosition(newPosition)
		}

		// Handle downward movement.
		newPosition = shape.GetPosition()

		newPosition.Y--

		if r.Tower.CanShapeMoveToPosition(shape, newPosition) {
			shape.SetPosition(newPosition)
		} else {
			r.Tower.LockShape(shape)
			break
		}
	}
}

type JetDirection int
type JetDirectionList []JetDirection

const (
	Left JetDirection = iota
	Right
)

func ParseJetDirections(line string) (JetDirectionList, error) {
	directionList := make(JetDirectionList, 0)

	for _, d := range line {
		var direction JetDirection
		switch d {
		case '>':
			direction = Right
		case '<':
			direction = Left
		default:
			return JetDirectionList{}, errors.New("invalid jet direction")
		}

		directionList = append(directionList, direction)
	}

	return directionList, nil
}

func day17(fileContents string) error {
	jetDirections, err := ParseJetDirections(fileContents)
	if err != nil {
		return err
	}

	// Part 1: After dropping 2022 rocks (shapes) which were buffeted by the jets, how tall will the tower of rocks be?
	room := NewRoom(7, jetDirections)

	for i := 0; i < 2022; i++ {
		room.DropShape()
	}

	fmt.Printf("Tower height: %d\n", room.GetTowerHeight())

	// Part 2: Elephants still don't believe you.  They want you to drop 1,000,000,000,000 rocks.  Now how tall will the tower of rocks be?
	room = NewRoom(7, jetDirections)

	for i := 0; i < 1000000000000; i++ {
		if i%10000 == 0 {
			fmt.Println("Shape", i)
		}
		room.DropShape()
	}

	fmt.Printf("Tower height: %d\n", room.GetTowerHeight())

	return nil
}
