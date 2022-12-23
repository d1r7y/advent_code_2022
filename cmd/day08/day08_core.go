package day08

import (
	"errors"
	"log"
	"strings"
)

type TreeRow []byte

type Forest struct {
	Trees []TreeRow
}

func NewForest() *Forest {
	return &Forest{Trees: make([]TreeRow, 0)}
}

func (f *Forest) scenicScoreForTree(x, y int) int {
	row := f.Trees[y]

	// Need to check four cardinal directions.
	height := row[x]

	// Up
	upViewDistance := 0
Up:
	for up := y - 1; up >= 0; up-- {
		upRow := f.Trees[up]
		upViewDistance++
		if upRow[x] >= height {
			break Up
		}
	}

	// Down
	downViewDistance := 0
Down:
	for down := y + 1; down < len(f.Trees); down++ {
		downRow := f.Trees[down]
		downViewDistance++
		if downRow[x] >= height {
			break Down
		}
	}

	// Left
	leftViewDistance := 0
Left:
	for left := x - 1; left >= 0; left-- {
		leftViewDistance++
		if row[left] >= height {
			break Left
		}
	}

	// Right
	rightViewDistance := 0
Right:
	for right := x + 1; right < len(row); right++ {
		rightViewDistance++
		if row[right] >= height {
			break Right
		}
	}

	return upViewDistance * downViewDistance * leftViewDistance * rightViewDistance
}

func (f *Forest) BestScenicScore() int {
	bestScenicScore := 0

	for i, row := range f.Trees {
		for j := range row {
			scenicScore := f.scenicScoreForTree(j, i)
			if scenicScore > bestScenicScore {
				bestScenicScore = scenicScore
			}
		}
	}
	return bestScenicScore
}

func (f *Forest) NumberVisibleTrees() int {
	// The edges are always visible.

	// The beginning and end of each row is visible.
	numVisible := 2 * len(f.Trees)

	// The entire top and bottom rows are visible.  Make sure not to count the beginning and
	// end of the top and bottom rows twice: they are accounted for above.
	numVisible += 2 * (len(f.Trees[0]) - 2)

	// Now look at the inner trees.
	for i := 1; i < len(f.Trees)-1; i++ {
		row := f.Trees[i]
	NextTree:
		for j := 1; j < len(row)-1; j++ {
			// Need to check four cardinal directions.
			height := row[j]

			treeVisible := true

			// Up
		Up:
			for up := i - 1; up >= 0; up-- {
				upRow := f.Trees[up]
				if upRow[j] >= height {
					treeVisible = false
					break Up
				}
			}

			// If it's visible, then we can short circuit the remaining checks.
			if treeVisible {
				numVisible++
				continue NextTree
			}

			treeVisible = true

		Down:
			// Down
			for down := i + 1; down < len(f.Trees); down++ {
				downRow := f.Trees[down]
				if downRow[j] >= height {
					treeVisible = false
					break Down
				}
			}

			// If it's visible, then we can short circuit the remaining checks.
			if treeVisible {
				numVisible++
				continue NextTree
			}

			treeVisible = true

			// Left
		Left:
			for left := j - 1; left >= 0; left-- {
				if row[left] >= height {
					treeVisible = false
					break Left
				}
			}

			// If it's visible, then we can short circuit the remaining checks.
			if treeVisible {
				numVisible++
				continue NextTree
			}

			treeVisible = true

			// Right
		Right:
			for right := j + 1; right < len(row); right++ {
				if row[right] >= height {
					treeVisible = false
					break Right
				}
			}

			// If it's visible, then we can short circuit the remaining checks.
			if treeVisible {
				numVisible++
				continue NextTree
			}
		}
	}

	return numVisible
}

func ParseTreeRow(str string) (TreeRow, error) {
	if str == "" {
		return TreeRow{}, errors.New("empty line")
	}

	row := make(TreeRow, 0)

	for _, char := range str {
		if char < '0' || char > '9' {
			return TreeRow{}, errors.New("invalid line")
		}

		height := byte(char) - byte('0')
		row = append(row, height)
	}

	return row, nil
}

func ParseForest(strs []string) (*Forest, error) {
	f := NewForest()

	for _, line := range strs {
		row, err := ParseTreeRow(line)
		if err != nil {
			return nil, err
		}

		f.Trees = append(f.Trees, row)
	}

	return f, nil
}

func day08(fileContents string) error {
	// Scan the forest in.
	f, err := ParseForest(strings.Split(string(fileContents), "\n"))
	if err != nil {
		return err
	}

	// Part 1: How many trees are visible from outside the forest?
	log.Printf("Number of visible trees: %d\n", f.NumberVisibleTrees())

	// Part 2: What is the highest possible scenic score possible for any tree?
	log.Printf("Best possible scenic score: %d\n", f.BestScenicScore())

	return nil
}
