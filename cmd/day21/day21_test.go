package day21

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateChannels(t *testing.T) {
	src := `root: pppw + sjmn
dbpl: 5
cczh: sllz + lgvd
zczc: 2
ptdq: humn - dvpt
dvpt: 3
lfqf: 4
humn: 5
ljgn: 2
sjmn: drzm * dbpl
sllz: 4
pppw: cczh / lfqf
lgvd: ljgn * ptdq
drzm: hmdt - zczc
hmdt: 32`

	rootChannel, err := CreateChannels(src)
	assert.NoError(t, err)
	if err != nil {
		assert.FailNow(t, "CreateChannels error", err)
	}

	number := <-rootChannel
	assert.Equal(t, 152, number)
}

func TestCreateTree(t *testing.T) {
	src := `root: pppw + sjmn
dbpl: 5
cczh: sllz + lgvd
zczc: 2
ptdq: humn - dvpt
dvpt: 3
lfqf: 4
humn: 5
ljgn: 2
sjmn: drzm * dbpl
sllz: 4
pppw: cczh / lfqf
lgvd: ljgn * ptdq
drzm: hmdt - zczc
hmdt: 32`

	_, err := CreateTree("humn", src)
	assert.NoError(t, err)
}

func TestEvaluate(t *testing.T) {
	src := `root: pppw + sjmn
dbpl: 5
cczh: sllz + lgvd
zczc: 2
ptdq: humn - dvpt
dvpt: 3
lfqf: 4
humn: 5
ljgn: 2
sjmn: drzm * dbpl
sllz: 4
pppw: cczh / lfqf
lgvd: ljgn * ptdq
drzm: hmdt - zczc
hmdt: 32`

	treeRoot, err := CreateTree("humn", src)
	assert.NoError(t, err)

	treeRoot.Evaluate(RootMonkeyName)

	rootNode := treeRoot.FindNode(RootMonkeyName)

	if rootNode.LeftPoisoned {
		assert.Equal(t, 150, rootNode.RightValue)
	} else {
		assert.Equal(t, 150, rootNode.LeftValue)
	}
}

func TestSolve(t *testing.T) {
	src := `root: pppw + sjmn
dbpl: 5
cczh: sllz + lgvd
zczc: 2
ptdq: humn - dvpt
dvpt: 3
lfqf: 4
humn: 5
ljgn: 2
sjmn: drzm * dbpl
sllz: 4
pppw: cczh / lfqf
lgvd: ljgn * ptdq
drzm: hmdt - zczc
hmdt: 32`

	treeRoot, err := CreateTree("humn", src)
	assert.NoError(t, err)

	treeRoot.Evaluate(RootMonkeyName)

	result := treeRoot.Solve(RootMonkeyName, 0)

	assert.Equal(t, 301, result)
}

func TestNewLeafNode(t *testing.T) {
	type testCase struct {
		name         string
		value        int
		expectedNode TreeNode
	}
	testCases := []testCase{
		{"root", 100, TreeNode{Name: "root", Value: 100, Leaf: true}},
		{"argb", 100, TreeNode{Name: "argb", Value: 100, Leaf: true}},
	}

	for _, test := range testCases {
		assert.Equal(t, test.expectedNode, *NewLeafNode(test.name, test.value))
	}
}

func TestNewParentNode(t *testing.T) {
	type testCase struct {
		name          string
		value         int
		expectedName  string
		expectedValue int
	}
	testCases := []testCase{
		{"root", 100, "root", 0},
		{"argb", 100, "argb", 0},
	}

	for _, test := range testCases {
		n := NewParentNode(test.name, '+')
		assert.Equal(t, test.expectedName, n.Name)
		assert.Equal(t, test.expectedValue, n.Value)
	}
}

func TestNodeSolve(t *testing.T) {
	type testCase struct {
		knownValue           int
		knownResult          int
		leftSolve            bool
		op                   byte
		expectedUnknownValue int
	}

	testCases := []testCase{
		// x + a = b
		{5, 100, true, '+', 95},
		// x - a = b
		{5, 100, true, '-', 105},
		// x / a = b
		{5, 100, true, '/', 500},
		// x * a = b
		{5, 100, true, '*', 20},
		// a + x = b
		{5, 100, false, '+', 95},
		// a - x = b
		{500, 100, false, '-', 400},
		// a / x = b
		{20, 5, false, '/', 4},
		// a * x = b
		{5, 100, false, '*', 20},
	}
	for _, test := range testCases {
		n := NewParentNode("test", test.op)
		if test.leftSolve {
			n.LeftPoisoned = true
		} else {
			n.RightPoisoned = true
		}
		assert.Equal(t, test.expectedUnknownValue, NodeSolve(n, test.knownValue, test.knownResult))
	}
}
