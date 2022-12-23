package day17

import (
	"errors"
	"log"
	"strings"
)

type Bitmap struct {
	Size Size
	Rows []byte
}

func NewBitmap(str string) (Bitmap, error) {
	lines := strings.Split(str, "\n")

	size := NewSize(0, 0)
	rows := make([]byte, 0)
	maxWidth := 0

	// Store bitmap upside down.  The first row is the bottom of the shape, last row is the top.
	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]

		row := byte(0)
		mask := byte(0x80)

		for i, c := range line {
			switch c {
			case '#':
				row |= mask
			case '.':
				// Row is 0, so no need to clear it out.
			default:
				return Bitmap{}, errors.New("unexpected character in bitmap")
			}

			if i+1 > maxWidth {
				maxWidth = i + 1
			}

			mask >>= 1
		}

		rows = append(rows, row)
		size.Height++
	}

	size.Width = maxWidth

	return Bitmap{Size: size, Rows: rows}, nil
}

func (b Bitmap) Describe() string {
	// The bitmap is stored upside down, so make sure to invert it
	// here.
	str := ""

	for i := len(b.Rows) - 1; i >= 0; i-- {
		row := b.Rows[i]

		if i < len(b.Rows)-1 {
			str += "\n"
		}

		mask := byte(0x80)
		for i := 0; i < b.Size.Width; i++ {
			if (row & mask) != 0 {
				str += "#"
			} else {
				str += "."
			}

			mask >>= 1
		}
	}

	return str
}

type Shape interface {
	GetSize() Size

	GetPosition() Point
	SetPosition(p Point)

	GetBitmap() Bitmap
}

type Point struct {
	X int64
	Y int64
}

func NewPoint(x, y int64) Point {
	return Point{X: x, Y: y}
}

type Size struct {
	Width  int
	Height int
}

func NewSize(w, h int) Size {
	return Size{Width: w, Height: h}
}

// BoundsIntersect Does point p2 intersect a rectangle whose bottom left corner is p1 and has size s?
func BoundsIntersect(p1 Point, s Size, p2 Point) bool {
	if p2.X < p1.X {
		return false
	}
	if p2.Y < p1.Y {
		return false
	}
	if p2.X > p1.X+int64(s.Width) {
		return false
	}
	if p2.Y > p1.Y+int64(s.Height) {
		return false
	}

	return true
}

type SquareShape struct {
	P Point
	B Bitmap
}

func NewSquareShape() Shape {
	b, err := NewBitmap("##\n##")
	if err != nil {
		log.Fatal(err)
	}

	return &SquareShape{B: b}
}

func (s *SquareShape) GetSize() Size {
	return NewSize(2, 2)
}

func (s *SquareShape) GetPosition() Point {
	return s.P
}

func (s *SquareShape) SetPosition(p Point) {
	s.P = p
}

func (s *SquareShape) GetBitmap() Bitmap {
	return s.B
}

type HorizontalLineShape struct {
	P Point
	B Bitmap
}

func NewHorizontalLineShape() Shape {
	b, err := NewBitmap("####")
	if err != nil {
		log.Fatal(err)
	}

	return &HorizontalLineShape{B: b}
}

func (s *HorizontalLineShape) GetSize() Size {
	return NewSize(4, 1)
}

func (s *HorizontalLineShape) GetPosition() Point {
	return s.P
}

func (s *HorizontalLineShape) SetPosition(p Point) {
	s.P = p
}

func (s *HorizontalLineShape) GetBitmap() Bitmap {
	return s.B
}

type VerticalLineShape struct {
	P Point
	B Bitmap
}

func NewVerticalLineShape() Shape {
	b, err := NewBitmap("#\n#\n#\n#")
	if err != nil {
		log.Fatal(err)
	}

	return &VerticalLineShape{B: b}
}

func (s *VerticalLineShape) GetSize() Size {
	return NewSize(1, 4)
}

func (s *VerticalLineShape) GetPosition() Point {
	return s.P
}

func (s *VerticalLineShape) SetPosition(p Point) {
	s.P = p
}

func (s *VerticalLineShape) GetBitmap() Bitmap {
	return s.B
}

type AngleShape struct {
	P Point
	B Bitmap
}

func NewAngleShape() Shape {
	b, err := NewBitmap("..#\n..#\n###")
	if err != nil {
		log.Fatal(err)
	}

	return &AngleShape{B: b}
}

func (s *AngleShape) GetSize() Size {
	return NewSize(3, 3)
}

func (s *AngleShape) GetPosition() Point {
	return s.P
}

func (s *AngleShape) SetPosition(p Point) {
	s.P = p
}

func (s *AngleShape) GetBitmap() Bitmap {
	return s.B
}

type CrossShape struct {
	P Point
	B Bitmap
}

func NewCrossShape() Shape {
	b, err := NewBitmap(".#.\n###\n.#.")
	if err != nil {
		log.Fatal(err)
	}

	return &CrossShape{B: b}
}

func (s *CrossShape) GetSize() Size {
	return NewSize(3, 3)
}

func (s *CrossShape) GetPosition() Point {
	return s.P
}

func (s *CrossShape) SetPosition(p Point) {
	s.P = p
}

func (s *CrossShape) GetBitmap() Bitmap {
	return s.B
}
