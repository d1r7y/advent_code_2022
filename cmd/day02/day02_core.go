package day02

import (
	"fmt"
	"log"
	"strings"
)

type Shape int

const (
	Rock Shape = iota
	Paper
	Scissor
)

type Result int

const (
	Win Result = iota
	Lose
	Draw
)

func parseFirstShape(str string) (Shape, error) {
	shapes := map[string]Shape{
		"A": Rock,
		"B": Paper,
		"C": Scissor,
	}

	if shape, ok := shapes[str]; ok {
		return shape, nil
	}

	return Rock, fmt.Errorf("invalid shape '%s'", str)
}

func parseSecondShape(str string) (Shape, error) {
	shapes := map[string]Shape{
		"X": Rock,
		"Y": Paper,
		"Z": Scissor,
	}

	if shape, ok := shapes[str]; ok {
		return shape, nil
	}

	return Rock, fmt.Errorf("invalid shape '%s'", str)
}

func parseResult(str string) (Result, error) {
	results := map[string]Result{
		"X": Lose,
		"Y": Draw,
		"Z": Win,
	}

	if result, ok := results[str]; ok {
		return result, nil
	}

	return Lose, fmt.Errorf("invalid result '%s'", str)
}

func getShapeScore(s Shape) int {
	scores := map[Shape]int{
		Rock:    1,
		Paper:   2,
		Scissor: 3,
	}

	return scores[s]
}

func getResultScore(r Result) int {
	scores := map[Result]int{
		Lose: 0,
		Draw: 3,
		Win:  6,
	}

	return scores[r]
}

func getShapeToWin(s Shape) Shape {
	switch s {
	case Rock:
		return Paper
	case Paper:
		return Scissor
	case Scissor:
		return Rock
	default:
		panic(fmt.Errorf("unknown shape: '%v'", s))
	}
}

func getShapeToLose(s Shape) Shape {
	switch s {
	case Rock:
		return Scissor
	case Paper:
		return Rock
	case Scissor:
		return Paper
	default:
		panic(fmt.Errorf("unknown shape: '%v'", s))
	}
}

func getShapeForResult(r Result, s Shape) Shape {
	switch r {
	case Draw:
		return s
	case Win:
		return getShapeToWin(s)
	case Lose:
		return getShapeToLose(s)
	default:
		panic(fmt.Errorf("unknown result: '%v'", r))
	}
}

type Round struct {
	shape1 Shape
	shape2 Shape
}

func NewRound(s1, s2 Shape) Round {
	return Round{shape1: s1, shape2: s2}
}

// Result What's the result of this round?  shape1 is the shape played by the other player, shape2 is your shape.
func (r Round) Result() Result {
	// Check if the same shapes were played.
	if r.shape1 == r.shape2 {
		return Draw
	}

	switch r.shape2 {
	case Rock:
		if r.shape1 == Scissor {
			return Win
		}
		return Lose
	case Paper:
		if r.shape1 == Rock {
			return Win
		}
		return Lose
	case Scissor:
		if r.shape1 == Paper {
			return Win
		}
		return Lose
	default:
		panic(fmt.Errorf("unknown shape: '%v'", r.shape2))
	}
}

// Result What's the score of this round?
func (r Round) Score() int {
	score := getShapeScore(r.shape2)
	score += getResultScore(r.Result())
	return score
}

func ParseRounds(text string, partOne bool) ([]Round, error) {

	rounds := make([]Round, 0)

	if text == "" {
		return rounds, nil
	}

	for _, line := range strings.Split(text, "\n") {
		if line == "" {
			continue
		}

		var shape1Str string
		var shape2Str string
		var resultStr string
		var count int
		var err error

		if partOne {
			count, err = fmt.Sscanln(line, &shape1Str, &shape2Str)
		} else {
			count, err = fmt.Sscanln(line, &shape1Str, &resultStr)
		}

		if err != nil {
			return []Round{}, err
		}

		if count != 2 {
			return []Round{}, fmt.Errorf("unexpected line in input '%s'", line)
		}

		var shape2 Shape

		shape1, err := parseFirstShape(shape1Str)
		if err != nil {
			return []Round{}, err
		}

		if partOne {
			shape2, err = parseSecondShape(shape2Str)
			if err != nil {
				return []Round{}, err
			}
		} else {
			result, err := parseResult(resultStr)
			if err != nil {
				return []Round{}, err
			}

			shape2 = getShapeForResult(result, shape1)
		}

		round := NewRound(shape1, shape2)
		rounds = append(rounds, round)
	}

	return rounds, nil
}

func day02(fileContents string) error {
	// Part 1: What is the total score if you followed the strategy?
	roundsPartOne, err := ParseRounds(string(fileContents), true)
	if err != nil {
		return err
	}

	totalScore := 0

	for _, r := range roundsPartOne {
		totalScore += r.Score()
	}

	log.Printf("Total score: %d\n", totalScore)

	// Part 2: What is the total score if you followed the strategy, where the second item in each round
	// is the result?
	roundsPartTwo, err := ParseRounds(string(fileContents), false)
	if err != nil {
		return err
	}

	totalScore = 0

	for _, r := range roundsPartTwo {
		totalScore += r.Score()
	}

	log.Printf("Total score: %d\n", totalScore)
	return nil
}
