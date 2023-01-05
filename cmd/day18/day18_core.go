package day18

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"sync"
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

func getNeighborString(n Neighbor) string {
	neighborMap := map[Neighbor]string{
		Top:    "Top",
		Bottom: "Bottom",
		Left:   "Left",
		Right:  "Right",
		Front:  "Front",
		Back:   "Back",
	}

	return neighborMap[n]
}

func Backtracking(n1, n2 Neighbor) bool {
	if n1 == Top && n2 == Bottom {
		return true
	}
	if n1 == Bottom && n2 == Top {
		return true
	}
	if n1 == Left && n2 == Right {
		return true
	}
	if n1 == Right && n2 == Left {
		return true
	}
	if n1 == Front && n2 == Back {
		return true
	}
	if n1 == Back && n2 == Front {
		return true
	}

	return false
}

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
	Position             Point
	FacesExposed         int
	ExternalFacesExposed int
}

type Plane []*Cube

type Tristate int

const (
	// Sentinel Tristate = iota
	False Tristate = iota
	True
	Pending
)

type Grid struct {
	Bounds  Bounds
	Offsets Offsets
	Min     Point
	Max     Point

	Space                       []Plane
	Cubes                       []*Cube
	EmptyCubeExternalPathCache  map[Point]bool
	EmptyCubeExternalPathCache2 map[Point]Tristate
	VisitedCube                 map[Point]Tristate

	mutex           sync.Mutex
	PendingChannels map[Point]chan bool
}

func (g *Grid) GetPendingChannel(p Point) chan bool {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if channel, ok := g.PendingChannels[p]; ok {
		return channel
	}
	// Channel doesn't exist, create it.
	channel := make(chan bool)

	g.PendingChannels[p] = channel
	return channel
}

func (g *Grid) GetCube(p Point) *Cube {
	return g.Space[p.Z+g.Offsets.Z][g.Bounds.W*(p.Y+g.Offsets.Y)+(p.X+g.Offsets.X)]
}

func (g *Grid) GetNeighbor(point Point, neighbor Neighbor) (*Cube, bool) {
	nc, _, valid := g.GetNeighborExtended(point, neighbor)
	return nc, valid
}

func (g *Grid) EdgePoint(p Point) bool {
	if p.Z == g.Max.Z {
		return true
	}
	if p.Z == g.Min.Z {
		return true
	}
	if p.Y == g.Max.Y {
		return true
	}
	if p.Y == g.Min.Y {
		return true
	}
	if p.X == g.Max.X {
		return true
	}
	if p.X == g.Min.X {
		return true
	}

	return false
}

func (g *Grid) CanReachEdge(emptyPoint Point) bool {

	g.VisitedCube[emptyPoint] = Pending
	defer func() {
		fmt.Printf("Marking %d,%d,%d as visited\n", emptyPoint.X, emptyPoint.Y, emptyPoint.Z)
		g.VisitedCube[emptyPoint] = True
	}()

	// Check our cache to see if this empty space has access to an edge.
	if access, ok := g.EmptyCubeExternalPathCache[emptyPoint]; ok {
		fmt.Printf("Cached %t for %d,%d,%d\n", access, emptyPoint.X, emptyPoint.Y, emptyPoint.Z)
		return access
	}

	for _, n := range []Neighbor{Top, Bottom, Left, Right, Front, Back} {
		if neighbor, point, ok := g.GetNeighborExtended(emptyPoint, n); ok {
			fmt.Printf("%d,%d,%d [%s] recurse\n", point.X, point.Y, point.Z, getNeighborString(n))
			if neighbor != nil {
				continue
			}
			if g.VisitedCube[point] == True {
				if access, ok := g.EmptyCubeExternalPathCache[point]; ok {
					if access {
						fmt.Printf("Access %t for %d,%d,%d\n", access, emptyPoint.X, emptyPoint.Y, emptyPoint.Z)
						g.EmptyCubeExternalPathCache[emptyPoint] = true
						return true
					}
				} else {
					fmt.Printf("unexpected: visited says yes, but cache has no entry for %d,%d,%d\n", point.X, point.Y, point.Z)
					os.Exit(0)
				}
			} else if g.VisitedCube[point] == False {
				if g.CanReachEdge(point) {
					fmt.Printf("Recursive result %d,%d,%d\n", emptyPoint.X, emptyPoint.Y, emptyPoint.Z)
					g.EmptyCubeExternalPathCache[emptyPoint] = true
					return true
				}
			} else if g.VisitedCube[point] == Pending {
				fmt.Printf("Pending for %d,%d,%d\n", point.X, point.Y, point.Z)
			}
		} else {
			fmt.Printf("Edged %t for %d,%d,%d\n", true, emptyPoint.X, emptyPoint.Y, emptyPoint.Z)
			g.EmptyCubeExternalPathCache[emptyPoint] = true
			return true
		}
	}

	fmt.Printf("Can't reach for %d,%d,%d\n", emptyPoint.X, emptyPoint.Y, emptyPoint.Z)

	// Can't reach an edge.
	g.EmptyCubeExternalPathCache[emptyPoint] = false

	return false
}

func (g *Grid) GetNeighborExtended(point Point, neighbor Neighbor) (*Cube, Point, bool) {
	switch neighbor {
	case Top:
		if point.Z == g.Max.Z {
			return nil, Point{}, false
		}
		position := Point{point.X, point.Y, point.Z + 1}
		return g.GetCube(position), position, true
	case Bottom:
		if point.Z == g.Min.Z {
			return nil, Point{}, false
		}
		position := Point{point.X, point.Y, point.Z - 1}
		return g.GetCube(position), position, true
	case Left:
		if point.X == g.Min.X {
			return nil, Point{}, false
		}
		position := Point{point.X - 1, point.Y, point.Z}
		return g.GetCube(position), position, true
	case Right:
		if point.X == g.Max.X {
			return nil, Point{}, false
		}
		position := Point{point.X + 1, point.Y, point.Z}
		return g.GetCube(position), position, true
	case Front:
		if point.Y == g.Max.Y {
			return nil, Point{}, false
		}
		position := Point{point.X, point.Y + 1, point.Z}
		return g.GetCube(position), position, true
	case Back:
		if point.Y == g.Min.Y {
			return nil, Point{}, false
		}
		position := Point{point.X, point.Y - 1, point.Z}
		return g.GetCube(position), position, true
	}

	log.Panic("unknown neighbor")
	return nil, Point{}, false
}

func (g *Grid) AddCube(cube *Cube) {
	g.Space[cube.Position.Z+g.Offsets.Z][g.Bounds.W*(cube.Position.Y+g.Offsets.Y)+(cube.Position.X+g.Offsets.X)] = cube
	g.Cubes = append(g.Cubes, cube)
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

	g.EmptyCubeExternalPathCache = make(map[Point]bool)
	g.EmptyCubeExternalPathCache2 = make(map[Point]Tristate)
	g.VisitedCube = make(map[Point]Tristate)
	g.PendingChannels = make(map[Point]chan bool)

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

	for i := range g.Space {
		g.Space[i] = make(Plane, g.Bounds.W*g.Bounds.D)
	}

	for _, line := range strings.Split(fileContents, "\n") {
		g.AddCube(ParseCube(line))
	}

	for _, cube := range g.Cubes {
		fmt.Printf("Cube: %d,%d,%d\n", cube.Position.X, cube.Position.Y, cube.Position.Z)
		g.VisitedCube[cube.Position] = Pending

		for _, n := range []Neighbor{Top, Bottom, Left, Right, Front, Back} {
			if neighbor, point, ok := g.GetNeighborExtended(cube.Position, n); ok {
				if neighbor != nil {
					continue
				}

				cube.FacesExposed++
				fmt.Printf("Scan: %d,%d,%d [%s]\n", point.X, point.Y, point.Z, getNeighborString(n))

				if g.CanReachEdge(point) {
					cube.ExternalFacesExposed++
				}
			} else {
				// If a side of a cube is on the edge, then it is an external face.
				cube.ExternalFacesExposed++
				cube.FacesExposed++
			}
		}
		fmt.Printf("Cube: %d,%d,%d faces %d external %d\n", cube.Position.X, cube.Position.Y, cube.Position.Z, cube.FacesExposed, cube.ExternalFacesExposed)
		g.VisitedCube[cube.Position] = True
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
