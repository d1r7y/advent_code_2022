package day12

import (
	"fmt"
	"log"
	"math"
	"sort"
	"strings"
)

type Direction byte

const (
	Up Direction = iota
	Down
	Left
	Right
)

func getDirectionString(d Direction) string {
	directionMap := map[Direction]string{
		Up:    "Up",
		Down:  "Down",
		Left:  "Left",
		Right: "Right",
	}

	return directionMap[d]
}

type Columns []int

func NewColumns() Columns {
	return make(Columns, 0)
}

type Position struct {
	X int
	Y int
}

type Size struct {
	W int
	H int
}

type World struct {
	Start      Position
	End        Position
	Dimensions Size
	Rows       []Columns
}

func NewWorld() *World {
	return &World{Rows: make([]Columns, 0)}
}

func (w *World) GetDimensions() Size {
	return w.Dimensions
}

func (w *World) GetPositionHeight(p Position) int {
	return w.Rows[p.Y][p.X]
}

func (w *World) GetAllPositionsAtHeight(height int) []Position {
	positions := make([]Position, 0)

	for i := range w.Rows {
		for j, h := range w.Rows[i] {
			if height == h {
				positions = append(positions, Position{X: j, Y: i})
			}
		}
	}

	return positions
}

type SolutionState struct {
	World    *World
	Position Position
	Visited  []bool
	Moves    []Direction
}

func NewSolutionState(world *World) *SolutionState {
	dimensions := world.GetDimensions()
	visited := make([]bool, dimensions.H*dimensions.W)

	return &SolutionState{World: world, Visited: visited, Moves: make([]Direction, 0)}
}

func (s *SolutionState) SetPosition(p Position) {
	s.Position = p

	// Mark new position as visited.
	dimensions := s.World.GetDimensions()
	s.Visited[p.Y*dimensions.W+p.X] = true
}

func (s *SolutionState) MoveDescription() string {
	description := ""

	for _, d := range s.Moves {
		description += getDirectionString(d) + " "
	}

	return description
}

func (s *SolutionState) AtEnd() bool {
	return s.Position == s.World.End
}

func (s *SolutionState) WasPositionVisited(proposed Direction) bool {
	dimensions := s.World.GetDimensions()

	proposedPosition := s.Position

	switch proposed {
	case Up:
		proposedPosition.Y--
	case Left:
		proposedPosition.X--
	case Right:
		proposedPosition.X++
	case Down:
		proposedPosition.Y++
	}

	// Make sure we don't exceed the bounds of the world.
	if proposedPosition.X < 0 || proposedPosition.X == dimensions.W {
		return false
	}

	if proposedPosition.Y < 0 || proposedPosition.Y == dimensions.H {
		return false
	}

	return s.Visited[proposedPosition.Y*dimensions.W+proposedPosition.X]
}

func (s *SolutionState) IsMoveLegal(proposed Direction) bool {
	// Have we already visited this position?
	if s.WasPositionVisited(proposed) {
		return false
	}

	// Next, check the proposed direction and the current position against the bounds of the world.
	if s.DoesMoveExceedBounds(proposed) {
		return false
	}

	// Now see if this would backtrack.
	if s.DoesMoveBacktrack(proposed) {
		return false
	}

	// Now check if the height of the proposed destination position is not too high.
	if s.IsProposedDestinationHeightInvalid(proposed) {
		return false
	}

	return true
}

func (s *SolutionState) GetLegalMoves() []Direction {
	moves := make([]Direction, 0)

	candidates := []Direction{Up, Down, Left, Right}

	for _, c := range candidates {
		if s.IsMoveLegal(c) {
			moves = append(moves, c)
		}
	}

	return moves
}

func (s *SolutionState) DoesMoveExceedBounds(proposed Direction) bool {
	dimensions := s.World.GetDimensions()
	switch proposed {
	case Up:
		return s.Position.Y == 0
	case Down:
		return s.Position.Y == dimensions.H-1
	case Left:
		return s.Position.X == 0
	case Right:
		return s.Position.X == dimensions.W-1
	}

	log.Panic("Invalid previous direction")
	return false
}

func (s *SolutionState) DoesMoveBacktrack(proposed Direction) bool {
	if len(s.Moves) == 0 {
		return false
	}

	previous := s.Moves[len(s.Moves)-1]
	switch previous {
	case Up:
		return proposed == Down
	case Left:
		return proposed == Right
	case Right:
		return proposed == Left
	case Down:
		return proposed == Up
	}

	log.Panic("Invalid previous direction")
	return false
}

func (s *SolutionState) IsProposedDestinationHeightInvalid(proposed Direction) bool {
	proposedPosition := s.Position

	switch proposed {
	case Up:
		proposedPosition.Y--
	case Left:
		proposedPosition.X--
	case Right:
		proposedPosition.X++
	case Down:
		proposedPosition.Y++
	}

	currentHeight := s.World.GetPositionHeight(s.Position)
	proposedHeight := s.World.GetPositionHeight(proposedPosition)

	return proposedHeight > currentHeight+1
}

func (s *SolutionState) Move(d Direction) *SolutionState {
	p := s.Position

	switch d {
	case Up:
		p.Y--
	case Left:
		p.X--
	case Right:
		p.X++
	case Down:
		p.Y++
	}

	ns := NewSolutionState(s.World)
	ns.Moves = make([]Direction, len(s.Moves))
	copy(ns.Moves, s.Moves)

	// This bit is a little tricky: we're sharing the Visited slice across all SolutionState copies.
	ns.Visited = s.Visited

	ns.Moves = append(ns.Moves, d)
	ns.SetPosition(p)

	return ns
}

func ParseWorld(fileContents string) *World {
	world := NewWorld()

	for row, line := range strings.Split(fileContents, "\n") {
		if line == "" {
			continue
		}

		columns := NewColumns()
		world.Dimensions.W = 0

		for column, c := range line {
			if c == 'S' {
				world.Start.X = column
				world.Start.Y = row

				// The start is at the lowest elevation 'a'
				columns = append(columns, int('a')-'a')
			} else if c == 'E' {
				world.End.X = column
				world.End.Y = row

				// The destination is at the highest elevation 'z'
				columns = append(columns, int('z')-'a')
			} else {
				if c < 'a' || c > 'z' {
					log.Fatal(fmt.Errorf("unknown character '%c' at %d,%d", c, column, row))
				}
				columns = append(columns, int(c)-'a')
			}
			world.Dimensions.W++
		}

		if len(columns) > 0 {
			world.Rows = append(world.Rows, columns)
			world.Dimensions.H++
		}
	}

	return world
}

func FindMinimumMovement(world *World) int {
	return FindMinimumMovementFromPosition(world, world.Start)
}

func FindMinimumMovementFromPosition(world *World, p Position) int {
	solutionState := NewSolutionState(world)
	solutionState.SetPosition(p)

	candidates := make([]*SolutionState, 0)
	candidates = append(candidates, solutionState)

	for i := 0; i < len(candidates); i++ {
		s := candidates[i]

		if s.AtEnd() {
			return len(s.Moves)
		}

		for _, m := range s.GetLegalMoves() {
			candidates = append(candidates, s.Move(m))
		}
	}

	return math.MaxInt
}

func day12(fileContents string) error {
	// Part 1: What is the fewest number of steps to go from the starting position to the
	// ending position.
	world := ParseWorld(fileContents)

	fmt.Printf("Minimum moves %d\n", FindMinimumMovement(world))

	// Part 2: Let's plan a more scenic route to the destination.  What is the fewest steps
	// required to move starting from any square with elevation a to the location that should
	// get the best signal?

	movesCount := make([]int, 0)

	for _, p := range world.GetAllPositionsAtHeight(int('a') - 'a') {
		moves := FindMinimumMovementFromPosition(world, p)
		movesCount = append(movesCount, moves)
	}

	sort.Ints(movesCount)

	fmt.Printf("Minimum moves from scenic positions %d\n", movesCount[0])

	return nil
}
