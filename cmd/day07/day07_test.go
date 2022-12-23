package day07

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFilesystemTree(t *testing.T) {
	tree := NewFilesystemTree()

	assert.Equal(t, tree.Root.GetName(), "/")
}

func TestNewFile(t *testing.T) {
	file := NewFile(nil, "child1")
	assert.Equal(t, "child1", file.Name)
	assert.Equal(t, "child1", file.GetName())
}

func TestNewDirectory(t *testing.T) {
	dir := NewDirectory(nil, "child1")
	assert.Equal(t, "child1", dir.Name)
	assert.Equal(t, "child1", dir.GetName())
}

func TestAddFileChild(t *testing.T) {
	parent := NewDirectory(nil, "parent")
	file := NewFile(nil, "child1")

	err := parent.AddChildren([]FilesystemNode{file})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(parent.Children))

	children, err := parent.GetChildren()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(children))
	assert.Equal(t, file.GetName(), children[0].GetName())
}

func TestAddDirChild(t *testing.T) {
	parent := NewDirectory(nil, "parent")
	dir := NewDirectory(nil, "child1")

	err := parent.AddChildren([]FilesystemNode{dir})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(parent.Children))

	children, err := parent.GetChildren()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(children))
	assert.Equal(t, dir.GetName(), children[0].GetName())
}

func TestAddChildError(t *testing.T) {
	parent := NewFile(nil, "parent")

	file := NewFile(nil, "child1")

	err := parent.AddChildren([]FilesystemNode{file})
	assert.Error(t, err)
	_, err = parent.GetChildren()
	assert.Error(t, err)

	dir := NewDirectory(nil, "child1")

	err = parent.AddChildren([]FilesystemNode{dir})
	assert.Error(t, err)
	_, err = parent.GetChildren()
	assert.Error(t, err)
}

func TestParseLsOutputLine(t *testing.T) {
	type testCase struct {
		str          string
		expectedErr  bool
		expectedType NodeType
		expectedName string
		expectedSize int64
	}

	tests := []testCase{
		{"dir lqwntmdg", false, DirectoryType, "lqwntmdg", 0},
		{"264381 tmwzlzn", false, FileType, "tmwzlzn", 264381},
		{"264381 12345", false, FileType, "12345", 264381},
		{"xyz abc", true, FileType, "", 0},
	}

	for _, test := range tests {
		node, err := ParseLsOutputLine(test.str)
		if test.expectedErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedType, node.GetType())
			assert.Equal(t, test.expectedName, node.GetName())
			assert.Equal(t, test.expectedSize, node.GetSize())
		}
	}
}
