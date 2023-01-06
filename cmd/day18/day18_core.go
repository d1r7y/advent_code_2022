package day18

import (
	"fmt"
	"log"
	"math"
	"strings"
)

type Neighbor int

const (
	Top Neighbor = iota
	Bottom
	Left
	Right
	Front
	Back
)

type Point struct {
	X, Y, Z int
}

type Bounds struct {
	W, H, D int
}

type Offsets struct {
	X, Y, Z int
}

type Cube struct {
	Empty                bool
	ExternalAccess       bool
	Position             Point
	FacesExposed         int
	ExternalFacesExposed int
}

func NewEmptyCube(p Point) *Cube {
	return &Cube{Empty: true, Position: p}
}

type Plane []*Cube

type Grid struct {
	Bounds  Bounds
	Offsets Offsets
	Min     Point
	Max     Point

	Space []Plane
	Cubes []*Cube
}

func (g *Grid) GetCube(p Point) *Cube {
	return g.Space[p.Z+g.Offsets.Z][g.Bounds.W*(p.Y+g.Offsets.Y)+(p.X+g.Offsets.X)]
}

func (g *Grid) GetNeighbor(point Point, neighbor Neighbor) (*Cube, bool) {
	switch neighbor {
	case Top:
		if point.Z == g.Max.Z {
			return nil, false
		}
		position := Point{point.X, point.Y, point.Z + 1}
		return g.GetCube(position), true
	case Bottom:
		if point.Z == g.Min.Z {
			return nil, false
		}
		position := Point{point.X, point.Y, point.Z - 1}
		return g.GetCube(position), true
	case Left:
		if point.X == g.Min.X {
			return nil, false
		}
		position := Point{point.X - 1, point.Y, point.Z}
		return g.GetCube(position), true
	case Right:
		if point.X == g.Max.X {
			return nil, false
		}
		position := Point{point.X + 1, point.Y, point.Z}
		return g.GetCube(position), true
	case Front:
		if point.Y == g.Max.Y {
			return nil, false
		}
		position := Point{point.X, point.Y + 1, point.Z}
		return g.GetCube(position), true
	case Back:
		if point.Y == g.Min.Y {
			return nil, false
		}
		position := Point{point.X, point.Y - 1, point.Z}
		return g.GetCube(position), true
	}

	log.Panic("unknown neighbor")
	return nil, false
}

func (g *Grid) AddCube(cube *Cube) {
	g.Space[cube.Position.Z+g.Offsets.Z][g.Bounds.W*(cube.Position.Y+g.Offsets.Y)+(cube.Position.X+g.Offsets.X)] = cube
	g.Cubes = append(g.Cubes, cube)
}

func (g *Grid) AddEmptyCube(p Point) {
	cube := NewEmptyCube(p)
	g.Space[cube.Position.Z+g.Offsets.Z][g.Bounds.W*(cube.Position.Y+g.Offsets.Y)+(cube.Position.X+g.Offsets.X)] = cube
}

func (g *Grid) GetSurfaceArea() int {
	exposedFaces := 0

	for _, cube := range g.Cubes {
		exposedFaces += cube.FacesExposed
	}

	return exposedFaces
}

func (g *Grid) GetExternalSurfaceArea() int {
	exposedFaces := 0

	for _, cube := range g.Cubes {
		exposedFaces += cube.ExternalFacesExposed
	}

	return exposedFaces
}

func (g *Grid) FillExternalAccess() {
	candidates := make([]Point, 0)
	visited := make(map[Point]bool)

	// Set ExternalAccess on all empty cubes on the border.
	for z := g.Min.Z; z <= g.Max.Z; z++ {
		for x := g.Min.X; x <= g.Max.X; x++ {
			p := Point{x, g.Min.Y, z}
			cube := g.GetCube(p)
			if cube.Empty {
				cube.ExternalAccess = true
				candidates = append(candidates, p)
				visited[p] = true
			}

			p = Point{x, g.Max.Y, z}
			cube = g.GetCube(p)
			if cube.Empty {
				cube.ExternalAccess = true
				candidates = append(candidates, p)
				visited[p] = true
			}
		}

		for y := g.Min.Y + 1; y <= g.Max.Y-1; y++ {
			p := Point{g.Min.X, y, z}
			cube := g.GetCube(p)
			if cube.Empty {
				cube.ExternalAccess = true
				candidates = append(candidates, p)
				visited[p] = true
			}

			p = Point{g.Max.X, y, z}
			cube = g.GetCube(p)
			if cube.Empty {
				cube.ExternalAccess = true
				candidates = append(candidates, p)
				visited[p] = true
			}
		}
	}

	// Now, starting with the positions in candidates, flood fill ExternalAccess to all accessible empty cubes.
	for i := 0; i < len(candidates); i++ {
		cp := candidates[i]

		for _, n := range []Neighbor{Top, Bottom, Left, Right, Front, Back} {
			if neighbor, ok := g.GetNeighbor(cp, n); ok {
				if neighbor.Empty {
					if !visited[neighbor.Position] {
						neighbor.ExternalAccess = true
						candidates = append(candidates, neighbor.Position)
						visited[neighbor.Position] = true
					}
				}
			}
		}
	}
}

func ParseCube(line string) *Cube {
	var x, y, z int

	count, err := fmt.Sscanf(line, "%d,%d,%d", &x, &y, &z)
	if err != nil {
		log.Panic("Invalid cube line.")
	}

	if count != 3 {
		log.Panic("Invalid cube line.")
	}

	return &Cube{Position: Point{X: x, Y: y, Z: z}, FacesExposed: 0}
}

func ParseCubes(fileContents string) *Grid {
	// Need to make three passes: first to get the bounds of the grid, next to allocate and store the
	// cubes, third to calculate the exposed faces.
	g := &Grid{}

	g.Min.X = math.MaxInt
	g.Max.X = math.MinInt
	g.Min.Y = math.MaxInt
	g.Max.Y = math.MinInt
	g.Min.Z = math.MaxInt
	g.Max.Z = math.MinInt

	for _, line := range strings.Split(fileContents, "\n") {
		cube := ParseCube(line)

		if cube.Position.X < g.Min.X {
			g.Min.X = cube.Position.X
		}
		if cube.Position.X > g.Max.X {
			g.Max.X = cube.Position.X
		}

		if cube.Position.Y < g.Min.Y {
			g.Min.Y = cube.Position.Y
		}
		if cube.Position.Y > g.Max.Y {
			g.Max.Y = cube.Position.Y
		}

		if cube.Position.Z < g.Min.Z {
			g.Min.Z = cube.Position.Z
		}
		if cube.Position.Z > g.Max.Z {
			g.Max.Z = cube.Position.Z
		}
	}

	g.Bounds.H = g.Max.Z - g.Min.Z + 1
	g.Bounds.W = g.Max.X - g.Min.X + 1
	g.Bounds.D = g.Max.Y - g.Min.Y + 1

	g.Offsets.X = -g.Min.X
	g.Offsets.Y = -g.Min.Y
	g.Offsets.Z = -g.Min.Z

	g.Space = make([]Plane, g.Bounds.H)

	// Prefill the space with empty cubes.
	z := g.Min.Z
	for i := range g.Space {
		g.Space[i] = make(Plane, g.Bounds.W*g.Bounds.D)
		for y := g.Min.Y; y <= g.Max.Y; y++ {
			for x := g.Min.X; x <= g.Max.X; x++ {
				g.AddEmptyCube(Point{x, y, z})
			}
		}
		z++
	}

	// Add the real cubes.
	for _, line := range strings.Split(fileContents, "\n") {
		g.AddCube(ParseCube(line))
	}

	// Set ExternalAccess on all empty cubes on the border.  Flood fill ExternalAccess to all reachable empty cubes.
	g.FillExternalAccess()

	for _, cube := range g.Cubes {
		for _, n := range []Neighbor{Top, Bottom, Left, Right, Front, Back} {
			if neighbor, ok := g.GetNeighbor(cube.Position, n); ok {
				if !neighbor.Empty {
					continue
				}

				cube.FacesExposed++
				if neighbor.ExternalAccess {
					cube.ExternalFacesExposed++
				}
			} else {
				// If a side of a cube is on the edge, then it is an external face.
				cube.ExternalFacesExposed++
				cube.FacesExposed++
			}
		}
	}

	return g
}

func day18(fileContents string) error {
	g := ParseCubes(fileContents)

	fmt.Printf("Bounds: %dx%dx%d\n", g.Bounds.W, g.Bounds.D, g.Bounds.H)
	fmt.Printf("Cubes read: %d\n", len(g.Cubes))
	fmt.Printf("Empty cubes: %d\n", g.Bounds.D*g.Bounds.H*g.Bounds.W-len(g.Cubes))

	// Part 1: After reading in the scanner report, what is the surface area of the lava droplet?
	fmt.Printf("Lava droplets surface area: %d units\n", g.GetSurfaceArea())

	// Part 2: Ignore the surfaces that are trapped within the droplets.  What is the exterior
	// surface area of the lava droplet?
	fmt.Printf("Lava droplets external surface area: %d units\n", g.GetExternalSurfaceArea())

	return nil
}
