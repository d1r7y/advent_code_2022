package day21

import (
	"errors"
	"fmt"
	"strings"
)

type OperationFn func(a, b int) int

type TreeNodeLookupMap map[string]*TreeNode

type TreeRoot struct {
	Lookup TreeNodeLookupMap
}

func NewTreeRoot(lookup TreeNodeLookupMap) *TreeRoot {
	return &TreeRoot{Lookup: lookup}
}

func (t *TreeRoot) FindNode(nodeName string) *TreeNode {
	return t.Lookup[nodeName]
}

func (t *TreeRoot) Evaluate(nodeName string) int {
	node := t.FindNode(nodeName)

	if node.Leaf {
		return node.Value
	}

	// Parent node, evaluate the left and right children.
	node.LeftValue = t.Evaluate(node.Left)
	node.RightValue = t.Evaluate(node.Right)

	// If either of the children are poisoned, we can't evaluate.  Return a sentinel for debugging.
	if node.LeftPoisoned || node.RightPoisoned {
		return -1
	}

	return node.Op(node.LeftValue, node.RightValue)
}

func (t *TreeRoot) Solve(nodeName string, result int) int {
	// Follow the poisoned path, evaluating each node until
	// we get to the poisoned node.

	node := t.FindNode(nodeName)

	// Handle root special.
	if nodeName == RootMonkeyName {
		if node.LeftPoisoned {
			result = node.RightValue
			node = t.FindNode(node.Left)
		} else {
			result = node.LeftValue
			node = t.FindNode(node.Right)
		}
	}

	if node.Poisoned {
		return result
	}

	// There should only be one poisoned child...
	if node.LeftPoisoned {
		newResult := NodeSolve(node, node.RightValue, result)
		return t.Solve(node.Left, newResult)
	} else {
		newResult := NodeSolve(node, node.LeftValue, result)
		return t.Solve(node.Right, newResult)
	}
}

func (t *TreeRoot) PromotePoison(nodeName string) bool {
	node := t.Lookup[nodeName]

	if !node.Leaf {
		if t.PromotePoison(node.Left) {
			node.LeftPoisoned = true
		}

		if t.PromotePoison(node.Right) {
			node.RightPoisoned = true
		}

		return node.LeftPoisoned || node.RightPoisoned
	}

	return node.Poisoned
}

type TreeNode struct {
	Poisoned bool

	Name string

	// Each node is either a leaf and has a value or has children and an operation.
	Leaf bool

	Value int

	Op            OperationFn
	InvertOpLeft  OperationFn
	InvertOpRight OperationFn
	Left          string
	LeftPoisoned  bool
	LeftValue     int
	Right         string
	RightPoisoned bool
	RightValue    int
}

const RootMonkeyName = "root"

var AddOp = func(a, b int) int { return a + b }
var InvertAddLeft = func(a, b int) int { return b - a }
var InvertAddRight = func(a, b int) int { return b - a }

var SubtractOp = func(a, b int) int { return a - b }
var InvertSubtractLeft = func(a, b int) int { return b + a }
var InvertSubtractRight = func(a, b int) int { return a - b }

var MultiplyOp = func(a, b int) int { return a * b }
var InvertMultiplyLeft = func(a, b int) int { return b / a }
var InvertMultiplyRight = func(a, b int) int { return b / a }

var DivideOp = func(a, b int) int { return a / b }
var InvertDivideLeft = func(a, b int) int { return b * a }
var InvertDivideRight = func(a, b int) int { return a / b }

func NewLeafNode(name string, value int) *TreeNode {
	return &TreeNode{Name: name, Leaf: true, Value: value}
}

func NodeSolve(node *TreeNode, a, result int) int {
	if node.LeftPoisoned {
		return node.InvertOpLeft(a, result)
	}
	return node.InvertOpRight(a, result)
}

func NewParentNode(name string, operation byte) *TreeNode {
	n := &TreeNode{Name: name, Leaf: false}

	switch operation {
	case '+':
		n.Op = AddOp
		n.InvertOpLeft = InvertAddLeft
		n.InvertOpRight = InvertAddRight
	case '-':
		n.Op = SubtractOp
		n.InvertOpLeft = InvertSubtractLeft
		n.InvertOpRight = InvertSubtractRight
	case '*':
		n.Op = MultiplyOp
		n.InvertOpLeft = InvertMultiplyLeft
		n.InvertOpRight = InvertMultiplyRight
	case '/':
		n.Op = DivideOp
		n.InvertOpLeft = InvertDivideLeft
		n.InvertOpRight = InvertDivideRight
	default:
		panic("invalid operation")
	}

	return n
}

func CreateTree(poisonName string, fileContents string) (*TreeRoot, error) {
	lookupMap := make(TreeNodeLookupMap)
	var rootNode *TreeNode

	for _, line := range strings.Split(fileContents, "\n") {

		var n string
		var operation byte
		var src1MonkeyName string
		var src2MonkeyName string

		// Try line of "monkeyName: monkeyName OP monkeyName" form first.
		count, err := fmt.Sscanf(line, "%s %s %c %s", &n, &src1MonkeyName, &operation, &src2MonkeyName)
		if err == nil && count == 4 {
			monkeyName := strings.TrimSuffix(n, ":")

			node := NewParentNode(monkeyName, operation)

			node.Left = src1MonkeyName
			node.Right = src2MonkeyName

			// Remember root node.
			if node.Name == RootMonkeyName {
				rootNode = node
			}

			// Add the node to the lookup table.
			lookupMap[monkeyName] = node

		} else {
			// Try line of "monkeyName: number" form.
			var number int

			count, err = fmt.Sscanf(line, "%s %d", &n, &number)
			if err != nil {
				return nil, err
			}
			if count != 2 {
				return nil, errors.New("invalid monkey line")
			}

			monkeyName := strings.TrimSuffix(n, ":")

			node := NewLeafNode(monkeyName, number)

			if monkeyName == poisonName {
				node.Poisoned = true

				// Let's stick a sentinel value here for debugging...
				node.Value = -1
			}

			// Add the node to the lookup table.
			lookupMap[monkeyName] = node
		}
	}

	treeRoot := NewTreeRoot(lookupMap)

	// Now walk through the tree, bottom up, bringing up the poison flag.
	treeRoot.PromotePoison(RootMonkeyName)

	if rootNode.LeftPoisoned && rootNode.RightPoisoned {
		return nil, errors.New("both sides of the root tree are poisoned")
	}

	return treeRoot, nil
}

type MonkeyChannel chan int

func GetOrCreateMonkeyChannel(monkeyChannelMap map[string]MonkeyChannel, monkeyName string) MonkeyChannel {
	mc, ok := monkeyChannelMap[monkeyName]
	if !ok {
		// Need to create the channel.
		mc = make(MonkeyChannel, 1)
		monkeyChannelMap[monkeyName] = mc
	}

	return mc
}

func CreateChannels(fileContents string) (MonkeyChannel, error) {
	monkeyChannelMap := make(map[string]MonkeyChannel)

	for _, line := range strings.Split(fileContents, "\n") {

		var monkeyName string
		var operation byte
		var src1MonkeyName string
		var src2MonkeyName string

		// Try line of "monkeyName: monkeyName OP monkeyName" form first.
		count, err := fmt.Sscanf(line, "%s %s %c %s", &monkeyName, &src1MonkeyName, &operation, &src2MonkeyName)
		if err == nil && count == 4 {
			// Has this monkey channel already been created?
			monkeyChannel := GetOrCreateMonkeyChannel(monkeyChannelMap, strings.TrimSuffix(monkeyName, ":"))

			// Get the channels for our dependencies.
			src1MonkeyChannel := GetOrCreateMonkeyChannel(monkeyChannelMap, src1MonkeyName)
			src2MonkeyChannel := GetOrCreateMonkeyChannel(monkeyChannelMap, src2MonkeyName)

			// Now spin up a go routine to wait for our dependencies.
			switch operation {
			case '+':
				go func() {
					num1 := <-src1MonkeyChannel
					num2 := <-src2MonkeyChannel
					monkeyChannel <- num1 + num2
				}()
			case '-':
				go func() {
					num1 := <-src1MonkeyChannel
					num2 := <-src2MonkeyChannel
					monkeyChannel <- num1 - num2
				}()
			case '*':
				go func() {
					num1 := <-src1MonkeyChannel
					num2 := <-src2MonkeyChannel
					monkeyChannel <- num1 * num2
				}()
			case '/':
				go func() {
					num1 := <-src1MonkeyChannel
					num2 := <-src2MonkeyChannel
					monkeyChannel <- num1 / num2
				}()
			default:
				return nil, errors.New("invalid operation")
			}
		} else {
			// Try line of "monkeyName: number" form.
			var number int

			count, err = fmt.Sscanf(line, "%s %d", &monkeyName, &number)
			if err != nil {
				return nil, err
			}
			if count != 2 {
				return nil, errors.New("invalid monkey line")
			}

			// Has this monkey channel already been created?
			monkeyChannel := GetOrCreateMonkeyChannel(monkeyChannelMap, strings.TrimSuffix(monkeyName, ":"))

			// Push the number into the channel
			monkeyChannel <- number
		}
	}

	rootChannel, ok := monkeyChannelMap[RootMonkeyName]
	if !ok {
		return nil, errors.New("missing root channel")
	}

	return rootChannel, nil
}

func day21(fileContents string) error {

	// Part 1: Monkeys yell numbers.  Other monkeys listen for specific other monkeys and do math on the numbers they here.
	// root is the alpha monkey.  What number will it yell?
	rootChannel, err := CreateChannels(fileContents)
	if err != nil {
		return err
	}

	// Now wait for the output of RootMonkeyName channel.
	fmt.Printf("Root will yell: %d\n", <-rootChannel)

	// Part 2: Confusion!  root monkey isn't doing math on its two dependent numbers: it's equality.  Both numbers need to be the same.
	// And humn monkey isn't a monkey, it's you!  So what number do you have to yell such that root's two dependent numbers are equal?
	t, err := CreateTree("humn", fileContents)
	if err != nil {
		return err
	}

	t.Evaluate(RootMonkeyName)

	result := t.Solve(RootMonkeyName, 0)
	fmt.Printf("You should yell: %d\n", result)

	return nil
}
