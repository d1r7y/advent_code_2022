package day14

import (
	"fmt"
	"log"
	"math"
	"strings"
)

type Cell byte

const (
	Edge     Cell = '#'
	Air      Cell = '.'
	Sand     Cell = 'o'
	Source   Cell = '+'
	Infinity Cell = '^'
)

func getCellString(cell Cell) string {
	cellMap := map[Cell]string{
		Edge:     "#",
		Air:      ".",
		Sand:     "o",
		Source:   "+",
		Infinity: "^",
	}

	return cellMap[cell]
}

type Point struct {
	X int
	Y int
}

type Size struct {
	W int
	H int
}

type Bounds struct {
	Origin Point
	Size   Size
}

type Column []Cell

type Cave struct {
	Bounds Bounds

	SandSource Point
	Columns    []Column
}

func NewCave(bounds Bounds, sandSource Point) *Cave {
	c := &Cave{Bounds: bounds}

	columns := make([]Column, bounds.Size.W+1)

	for i := range columns {
		column := make(Column, bounds.Size.H+1)
		for j := range column {
			column[j] = Air
		}

		columns[i] = column
	}

	c.Columns = columns

	// Add in the sand source.
	c.SandSource = sandSource

	c.SetCell(c.SandSource, Source)

	return c
}

func (c *Cave) SetCell(p Point, cell Cell) {
	c.Columns[p.X-c.Bounds.Origin.X][p.Y-c.Bounds.Origin.Y] = cell
}

func (c *Cave) GetCell(p Point) Cell {
	if p.X < c.Bounds.Origin.X {
		return Infinity
	}

	if p.X >= c.Bounds.Origin.X+c.Bounds.Size.W {
		return Infinity
	}

	if p.Y < c.Bounds.Origin.Y {
		return Infinity
	}

	if p.Y >= c.Bounds.Origin.Y+c.Bounds.Size.H {
		return Infinity
	}

	return c.Columns[p.X-c.Bounds.Origin.X][p.Y-c.Bounds.Origin.Y]
}

func (c *Cave) AddEdge(p1 Point, p2 Point) {
	if p1.X == p2.X {
		// Vertical edge
		// Find the starting point.
		var startY int
		var endY int

		if p1.Y > p2.Y {
			startY = p2.Y
			endY = p1.Y
		} else {
			startY = p1.Y
			endY = p2.Y
		}

		for i := startY; i <= endY; i++ {
			c.SetCell(Point{p1.X, i}, Edge)
		}
	} else if p1.Y == p2.Y {
		// Horizontal edge
		// Find the starting point.
		var startX int
		var endX int

		if p1.X > p2.X {
			startX = p2.X
			endX = p1.X
		} else {
			startX = p1.X
			endX = p2.X
		}

		for i := startX; i <= endX; i++ {
			c.SetCell(Point{i, p1.Y}, Edge)
		}
	} else {
		log.Panic("Non vertical or horizontal line")
	}
}

func (c *Cave) Describe() string {
	description := ""
	for i := 0; i < c.Bounds.Size.H; i++ {
		for j := 0; j < c.Bounds.Size.W; j++ {
			p := Point{j + c.Bounds.Origin.X, i + c.Bounds.Origin.Y}
			cell := c.GetCell(p)
			description += getCellString(cell)
		}

		if i != c.Bounds.Size.H-1 {
			description += "\n"
		}
	}
	return description
}

type DropResult int

const (
	SandAtRest DropResult = iota
	SandFalling
	SandBlocked
)

func (c *Cave) DropSand() DropResult {
	result := SandFalling

	position := c.SandSource

	for result != SandAtRest {
		// Check if it can drop down.
		nextPosition := position
		nextPosition.Y++

		cell := c.GetCell(nextPosition)

		if cell == Infinity {
			return SandFalling
		}

		if cell == Air {
			position = nextPosition
			continue
		}

		// Nope.  Try down and left.
		nextPosition = position
		nextPosition.Y++
		nextPosition.X--

		cell = c.GetCell(nextPosition)

		if cell == Infinity {
			return SandFalling
		}

		if cell == Air {
			position = nextPosition
			continue
		}

		// Try down and right.
		nextPosition = position
		nextPosition.Y++
		nextPosition.X++

		cell = c.GetCell(nextPosition)

		if cell == Infinity {
			return SandFalling
		}

		if cell == Air {
			position = nextPosition
			continue
		}

		// Nope.  Sand is at rest.
		c.SetCell(position, Sand)
		if position != c.SandSource {
			result = SandAtRest
		} else {
			return SandBlocked
		}
	}

	return result
}

func ParseCave(fileContents string, infiniteAbyss bool) *Cave {
	var maxDepth = 0
	var minWidth = math.MaxInt
	var maxWidth = -1

	// Determine the size of the cave from the paths of the rock structures.
	for _, line := range strings.Split(fileContents, "\n") {
		points := ParsePath(line)
		for _, p := range points {
			if p.Y > maxDepth {
				maxDepth = p.Y
			}
			if p.X < minWidth {
				minWidth = p.X
			}
			if p.X > maxWidth {
				maxWidth = p.X
			}
		}
	}

	sandSource := Point{500, 0}

	var caveBounds Bounds

	// Probably a smarter way to handle this, but let's just make the cave quite wide to handle the "infinite"
	// bottom row...
	if infiniteAbyss {
		caveBounds = Bounds{Origin: Point{minWidth, 0}, Size: Size{(maxWidth - minWidth) + 1, maxDepth + 1}}
	} else {
		caveBounds = Bounds{Origin: Point{0, 0}, Size: Size{1001, maxDepth + 3}}
	}

	cave := NewCave(caveBounds, sandSource)

	// Now add the points
	for _, line := range strings.Split(fileContents, "\n") {
		points := ParsePath(line)
		for i := 0; i < len(points)-1; i++ {
			p1 := points[i]
			p2 := points[i+1]

			cave.AddEdge(p1, p2)
		}
	}

	// If there's no infinite abyss below us, add the final line.
	if !infiniteAbyss {
		p1 := Point{0, maxDepth + 2}
		p2 := Point{1000, maxDepth + 2}
		cave.AddEdge(p1, p2)
	}

	return cave
}

func ParsePath(line string) []Point {
	points := make([]Point, 0)

	for _, ps := range strings.Split(line, " -> ") {
		var x, y int

		count, err := fmt.Sscanf(ps, "%d,%d", &x, &y)
		if err != nil {
			log.Panic(err)
		}

		if count != 2 {
			log.Panic("invalid point")
		}

		points = append(points, Point{X: x, Y: y})
	}

	return points
}

func day14(fileContents string) error {
	// Part 1: How many units of sand come to rest before sand starts flowing into the abyss below?
	cave := ParseCave(fileContents, true)

	fmt.Println(cave.Describe())

	sandCount := 0

	for {
		result := cave.DropSand()
		if result == SandFalling {
			break
		}
		sandCount++
	}

	fmt.Printf("%d sand units come to rest before the others start flowing into the abyss.\n", sandCount)

	fmt.Println(cave.Describe())

	// Part 2: You misread the scan.  There isn't an infinite void.  You're standing on the floor.  It's
	// an infinite horizontal line with a Y coordinate +2 of the highest Y coordinate of any point in your
	// scan.  How much sand can drop until it blocks the source?

	cave2 := ParseCave(fileContents, false)

	sandCount2 := 1

	for {
		result := cave2.DropSand()
		if result == SandBlocked {
			break
		}
		sandCount2++
	}

	fmt.Printf("%d sand units come to rest before the source is blocked.\n", sandCount2)

	return nil
}
