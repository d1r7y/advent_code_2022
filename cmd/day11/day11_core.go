package day11

import (
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"strings"
)

type Item struct {
	WorryLevel *big.Int
}

func NewItem(worryLevel int) Item {
	return Item{WorryLevel: big.NewInt(int64(worryLevel))}
}

func (i Item) Describe() string {
	return fmt.Sprintf("%d", i.WorryLevel)
}

type Jungle struct {
	Monkeys                  []*Monkey
	UndamagedWorryAdjustment *big.Int
}

func NewJungle(monkeys []*Monkey) *Jungle {
	j := &Jungle{Monkeys: monkeys, UndamagedWorryAdjustment: big.NewInt(int64(3))}

	for _, m := range j.Monkeys {
		m.SetHome(j)
	}

	return j
}

func (j *Jungle) Evaluate() {
	for _, m := range j.Monkeys {
		for _, ir := range m.EvaluateItems() {
			j.Monkeys[ir.MonkeyID].AddItem(ir.Item)
		}
	}
}

func (j *Jungle) Describe() string {
	str := ""

	for id, m := range j.Monkeys {
		str += fmt.Sprintf("Monkey %d: %s\n", id, m.Describe())
	}

	return str
}

func (j *Jungle) GetMonkeyInspectionCounts() []int {
	activity := make([]int, 0)
	for _, m := range j.Monkeys {
		activity = append(activity, m.GetInspectionCount())
	}

	return activity
}

func (j *Jungle) SetUndamagedWorryLevelAdjustment(adjustment int) {
	j.UndamagedWorryAdjustment = big.NewInt(int64(adjustment))
}

func (j *Jungle) ApplyUndamagedWorryAdjustment(item Item) Item {
	item.WorryLevel.Div(item.WorryLevel, j.UndamagedWorryAdjustment)
	return item
}

type OperationFn func(item Item) Item
type TestFn func(item Item) bool

type Monkey struct {
	Home            *Jungle
	Items           []Item
	InspectionCount int
	Operation       OperationFn
	Test            TestFn
	TrueResult      int
	FalseResult     int
}

func NewMonkey() *Monkey {
	monkey := &Monkey{Items: make([]Item, 0)}

	monkey.Operation = func(item Item) Item {
		panic("unset Operation")
	}

	monkey.Test = func(item Item) bool {
		panic("unset Test")
	}

	return monkey
}

func (m *Monkey) SetHome(jungle *Jungle) {
	m.Home = jungle
}

func (m *Monkey) Describe() string {
	strs := make([]string, 0)
	for _, item := range m.Items {
		strs = append(strs, item.Describe())
	}

	return strings.Join(strs, ", ")
}

func (m *Monkey) AddItem(item Item) {
	m.Items = append(m.Items, item)
}

func (m *Monkey) GetInspectionCount() int {
	return m.InspectionCount
}

type ItemDestination struct {
	Item     Item
	MonkeyID int
}

func NewItemDestination(item Item, monkeyID int) ItemDestination {
	return ItemDestination{Item: item, MonkeyID: monkeyID}
}

func (m *Monkey) EvaluateItems() []ItemDestination {
	ids := make([]ItemDestination, 0)

	for _, item := range m.Items {
		m.Items = m.Items[1:]

		m.InspectionCount++

		updatedItem := m.Home.ApplyUndamagedWorryAdjustment(m.Operation(item))

		var monkeyID int
		if m.Test(updatedItem) {
			monkeyID = m.TrueResult
		} else {
			monkeyID = m.FalseResult
		}

		id := NewItemDestination(updatedItem, monkeyID)

		ids = append(ids, id)
	}

	return ids
}

func (m *Monkey) SetOperation(f OperationFn) {
	m.Operation = f
}

func (m *Monkey) SetTest(f TestFn) {
	m.Test = f
}

func (m *Monkey) SetTrueResult(tr int) {
	m.TrueResult = tr
}

func (m *Monkey) SetFalseResult(fr int) {
	m.FalseResult = fr
}

func ParseTest(str string) (TestFn, error) {
	var constant int

	c, err := fmt.Sscanf(str, "  Test: divisible by %d", &constant)
	if err != nil {
		return nil, err
	}
	if c != 1 {
		return nil, errors.New("invalid test definition")
	}

	return func(item Item) bool {
		result := new(big.Int)
		return result.Mod(item.WorryLevel, big.NewInt(int64(constant))).Cmp(big.NewInt(0)) == 0
		// return (item.WorryLevel % int64(constant)) == 0
	}, nil
}

func ParseTestResult(str string, evaluator string) (int, error) {
	var monkeyID int

	formatStr := fmt.Sprintf("    If %s: throw to monkey", evaluator)
	c, err := fmt.Sscanf(str, formatStr+" %d", &monkeyID)
	if err != nil {
		return 0, err
	}
	if c != 1 {
		return 0, errors.New("invalid test result definition")
	}

	return monkeyID, nil
}

func ParseOperation(str string) (OperationFn, error) {
	var operand1Str string
	var operand2Str string
	var operatorStr string

	c, err := fmt.Sscanf(str, "  Operation: new = %s %s %s", &operand1Str, &operatorStr, &operand2Str)
	if err != nil {
		return nil, err
	}
	if c != 3 {
		return nil, errors.New("invalid operation definition")
	}

	type sourceFn func(currentWorryLevel *big.Int) *big.Int
	type operatorFn func(a, b *big.Int) *big.Int

	var operand1Src sourceFn
	var operand2Src sourceFn
	var operator operatorFn

	if operand1Str == "old" {
		operand1Src = func(currentWorryLevel *big.Int) *big.Int { return currentWorryLevel }
	} else {
		// We'll assume it's a number
		constant, err := strconv.Atoi(operand1Str)
		if err != nil {
			return nil, errors.New("invalid constant")
		}
		operand1Src = func(currentWorryLevel *big.Int) *big.Int { return big.NewInt(int64(constant)) }
	}

	if operand2Str == "old" {
		operand2Src = func(currentWorryLevel *big.Int) *big.Int { return currentWorryLevel }
	} else {
		// We'll assume it's a number
		constant, err := strconv.Atoi(operand2Str)
		if err != nil {
			return nil, errors.New("invalid constant")
		}
		operand2Src = func(currentWorryLevel *big.Int) *big.Int { return big.NewInt(int64(constant)) }
	}

	switch operatorStr {
	case "+":
		operator = func(a, b *big.Int) *big.Int {
			result := new(big.Int)
			return result.Add(a, b)
		}
	case "-":
		operator = func(a, b *big.Int) *big.Int {
			result := new(big.Int)
			return result.Sub(a, b)
		}
	case "*":
		operator = func(a, b *big.Int) *big.Int {
			result := new(big.Int)
			return result.Mul(a, b)
		}
	case "/":
		operator = func(a, b *big.Int) *big.Int {
			result := new(big.Int)
			return result.Div(a, b)
		}
	default:
		return nil, errors.New("unknown operator")
	}

	return func(item Item) Item {
		item.WorryLevel = operator(operand1Src(item.WorryLevel), operand2Src(item.WorryLevel))

		return item
	}, nil
}

func ParseItemList(str string) ([]Item, error) {
	items := make([]Item, 0)

	itemsPrefixString := "  Starting items: "

	if !strings.HasPrefix(str, itemsPrefixString) {
		return []Item{}, errors.New("invalid starting items definition")
	}

	// Now we need to scan the list...
	for _, itemWorryLevelStr := range strings.Split(strings.TrimPrefix(str, itemsPrefixString), ",") {
		itemWorryLevel, err := strconv.Atoi(strings.TrimSpace(itemWorryLevelStr))
		if err != nil {
			return []Item{}, err
		}

		item := NewItem(itemWorryLevel)
		items = append(items, item)
	}

	return items, nil
}

func ParseNotes(lines []string) ([]*Monkey, error) {
	monkeys := make([]*Monkey, 0)
	for i := 0; i < len(lines); i++ {
		// First line should be the monkey definition.
		var monkeyID int

		c, err := fmt.Sscanf(lines[i], "Monkey %d:", &monkeyID)
		if err != nil {
			return []*Monkey{}, err
		}
		if c != 1 {
			return []*Monkey{}, errors.New("invalid monkey definition")
		}

		monkey := NewMonkey()

		i++

		// Next parse the starting items definition.
		itemList, err := ParseItemList(lines[i])
		if err != nil {
			return []*Monkey{}, err
		}
		for _, item := range itemList {
			monkey.AddItem(item)
		}

		i++

		// Parse the operation.
		operation, err := ParseOperation(lines[i])
		if err != nil {
			return []*Monkey{}, err
		}

		monkey.SetOperation(operation)

		i++

		// Parse the test.
		test, err := ParseTest(lines[i])
		if err != nil {
			return []*Monkey{}, err
		}

		monkey.SetTest(test)

		i++

		// Parse the true test result
		thrownMonkeyID, err := ParseTestResult(lines[i], "true")
		if err != nil {
			return []*Monkey{}, err
		}

		monkey.SetTrueResult(thrownMonkeyID)

		i++

		// Parse the false test result
		thrownMonkeyID, err = ParseTestResult(lines[i], "false")
		if err != nil {
			return []*Monkey{}, err
		}

		monkey.SetFalseResult(thrownMonkeyID)

		// This next increment skips over the blank line.
		i++

		monkeys = append(monkeys, monkey)
	}

	return monkeys, nil
}

func day11(fileContents string) error {
	// Part 1: After evaluating all the monkey shines, the monkey level business is the activity of the top two monkeys multiplied together.  What is it?

	// Scan monkey notes.
	monkeys, err := ParseNotes(strings.Split(string(fileContents), "\n"))
	if err != nil {
		return err
	}

	j := NewJungle(monkeys)
	for i := 0; i < 20; i++ {
		j.Evaluate()
	}

	inspectionCounts := j.GetMonkeyInspectionCounts()
	sort.Sort(sort.Reverse(sort.IntSlice(inspectionCounts)))

	fmt.Printf("Level of monkey business: %d\n", inspectionCounts[0]*inspectionCounts[1])

	// Part 2: Now you're so worried that your relief that the items are undamaged don't lower your worry level by 3.  Now you need to run 10,000.
	// What is the new monkey level business?

	// Scan monkey notes.
	monkeys, err = ParseNotes(strings.Split(string(fileContents), "\n"))
	if err != nil {
		return err
	}

	j = NewJungle(monkeys)
	j.SetUndamagedWorryLevelAdjustment(1)

	for i := 1; i <= 10000; i++ {
		j.Evaluate()
		// fmt.Println(j.Describe())

		fmt.Printf("== After round %d ==\n", i)
		for mid, count := range j.GetMonkeyInspectionCounts() {
			fmt.Printf("Monkey %d inspected items %d times.\n", mid, count)
		}
	}

	inspectionCounts = j.GetMonkeyInspectionCounts()
	sort.Sort(sort.Reverse(sort.IntSlice(inspectionCounts)))

	fmt.Printf("Level of monkey business: %d\n", inspectionCounts[0]*inspectionCounts[1])

	return nil
}
