package day07

import (
	"errors"
	"fmt"
	"math"
	"path/filepath"
	"strings"
)

type NodeType int

const (
	DirectoryType NodeType = iota
	FileType
)

type WalkFunction func(n FilesystemNode) error

type FilesystemNode interface {
	GetType() NodeType

	GetName() string
	SetName(name string)

	Find(path string) (FilesystemNode, error)

	GetSize() int64
	SetSize(size int64) error

	GetChildren() ([]FilesystemNode, error)
	AddChildren(children []FilesystemNode) error

	GetParent() FilesystemNode
	SetParent(parent FilesystemNode) error

	Walk(depth bool, wf WalkFunction) error
}

type FilesystemTree struct {
	Root FilesystemNode
}

func NewFilesystemTree() *FilesystemTree {
	return &FilesystemTree{Root: NewRootDirectory()}
}

func (t *FilesystemTree) Find(path string) (FilesystemNode, error) {
	if !filepath.IsAbs(path) {
		return nil, fmt.Errorf("non-absolute path '%s'", path)
	}

	if path == "/" {
		return t.Root, nil
	}

	// Remove leading slash
	rel := path[1:]

	return t.Root.Find(rel)
}

type File struct {
	Name   string
	Size   int64
	Parent FilesystemNode
}

func NewFile(parent FilesystemNode, name string) *File {
	return &File{Parent: parent, Name: name}
}

func (f *File) GetType() NodeType {
	return FileType
}

func (f *File) GetName() string {
	return f.Name
}

func (f *File) SetName(name string) {
	f.Name = name
}

func (f *File) Find(path string) (FilesystemNode, error) {
	return nil, errors.New("files don't have children")
}

func (f *File) GetSize() int64 {
	return f.Size
}

func (f *File) SetSize(size int64) error {
	f.Size = size
	return nil
}

func (f *File) GetChildren() ([]FilesystemNode, error) {
	return []FilesystemNode{}, errors.New("files don't have children")
}

func (f *File) AddChildren(children []FilesystemNode) error {
	return errors.New("files don't have children")
}

func (f *File) GetParent() FilesystemNode {
	return f.Parent
}

func (f *File) SetParent(parent FilesystemNode) error {
	if parent.GetType() != DirectoryType {
		return errors.New("parent not directory")
	}
	if f == parent {
		return errors.New("parent/child loop")
	}

	f.Parent = parent

	return nil
}

func (f *File) Walk(depth bool, wf WalkFunction) error {
	return wf(f)
}

type Directory struct {
	Name     string
	Children []FilesystemNode
	Parent   FilesystemNode
}

func NewDirectory(parent FilesystemNode, name string) *Directory {
	return &Directory{Parent: parent, Name: name}
}

func NewRootDirectory() *Directory {
	root := &Directory{Name: "/"}
	root.Parent = root

	return root
}

func (d *Directory) GetType() NodeType {
	return DirectoryType
}

func (d *Directory) GetName() string {
	return d.Name
}

func (d *Directory) SetName(name string) {
	d.Name = name
}

func (d *Directory) Find(path string) (FilesystemNode, error) {
	elements := strings.Split(path, "/")

	if path == "" {
		return d, nil
	}

	if path == ".." {
		return d.Parent, nil
	}

	for _, child := range d.Children {
		if child.GetName() == elements[0] {
			return child.Find(filepath.Join(elements[1:]...))
		}
	}

	return nil, fmt.Errorf("no such child '%s'", elements[0])
}

func (d *Directory) GetSize() int64 {
	var totalSize int64

	for _, child := range d.Children {
		totalSize += child.GetSize()
	}

	return totalSize
}

func (d *Directory) SetSize(size int64) error {
	return errors.New("can't change directory size")
}

func (d *Directory) GetChildren() ([]FilesystemNode, error) {
	return d.Children, nil
}

func (d *Directory) AddChildren(children []FilesystemNode) error {
	d.Children = append(d.Children, children...)
	return nil
}

func (d *Directory) GetParent() FilesystemNode {
	return d.Parent
}

func (d *Directory) SetParent(parent FilesystemNode) error {
	if parent.GetType() != DirectoryType {
		return errors.New("parent not directory")
	}
	if d == parent {
		return errors.New("parent/child loop")
	}

	d.Parent = parent

	return nil
}

func (d *Directory) Walk(depth bool, wf WalkFunction) error {
	if !depth {
		if err := wf(d); err != nil {
			return err
		}
	}

	for _, child := range d.Children {
		if err := child.Walk(depth, wf); err != nil {
			return err
		}
	}

	if depth {
		if err := wf(d); err != nil {
			return err
		}
	}

	return nil
}

func ParseLsOutputLine(str string) (FilesystemNode, error) {
	// See if it's a directory first.
	var name string

	count, err := fmt.Sscanf(str, "dir %s", &name)
	if err == nil && count == 1 {
		return NewDirectory(nil, name), nil
	}

	// Now see if it's a file.
	var size int64
	count, err = fmt.Sscanf(str, "%d %s", &size, &name)
	if err == nil && count == 2 {
		file := NewFile(nil, name)
		file.SetSize(size)
		return file, nil
	}

	return nil, errors.New("invalid line")
}

func IsCommand(str string) bool {
	return strings.HasPrefix(str, "$ ")
}

func day07(fileContents string) error {
	fs := NewFilesystemTree()

	var cwd = fs.Root

	lsOutputMode := false

NextLine:
	for _, line := range strings.Split(string(fileContents), "\n") {
		if lsOutputMode {
			if IsCommand(line) {
				lsOutputMode = false
			} else {
				node, err := ParseLsOutputLine(line)
				if err != nil {
					return err
				}
				node.SetParent(cwd)
				cwd.AddChildren([]FilesystemNode{node})
				continue NextLine
			}
		}

		if IsCommand(line) {
			if strings.HasPrefix(line, "$ ls") {
				lsOutputMode = true
				continue NextLine
			} else if strings.HasPrefix(line, "$ cd") {
				var name string
				count, err := fmt.Sscanf(line, "$ cd %s", &name)
				if err != nil {
					return err
				}
				if count != 1 {
					return errors.New("invalid line")
				}

				if filepath.IsAbs(name) {
					node, err := fs.Find(name)
					if err != nil {
						return err
					}
					cwd = node
				} else {
					node, err := cwd.Find(name)
					if err != nil {
						return err
					}
					cwd = node
				}
			} else {
				return fmt.Errorf("unknown command '%s'", line)
			}

		}
	}

	// Part 1: Find all of the directories with a total size of at most 100000.  What is the sum of the total sizes
	// of those directories?

	var totalSize int64

	fs.Root.Walk(true, func(n FilesystemNode) error {
		if n.GetType() == DirectoryType && n.GetSize() <= 100000 {
			totalSize += n.GetSize()
		}

		return nil
	})

	fmt.Printf("Sum of the total sizes of directories whose size is at most 100000: %d\n", totalSize)

	const DiskSize = int64(70000000)
	const UpdateSize = int64(30000000)

	// Part 2: Disk is DiskSize.  Update needs UpdateSize bytes free.  Find the smallest directory we can delete that will
	// allow us to do the update.

	availableSpace := DiskSize - fs.Root.GetSize()
	if availableSpace < UpdateSize {
		amountToDelete := UpdateSize - availableSpace
		fmt.Printf("Amount to delete: %d\n", amountToDelete)

		minimumSize := int64(math.MaxInt64)

		fs.Root.Walk(true, func(n FilesystemNode) error {
			if n.GetType() == DirectoryType {
				dirSize := n.GetSize()

				// Is this directory bigger than what we need to delete, but smaller than the smallest we've seen up to now?
				// Then remember its size.
				if dirSize >= amountToDelete && dirSize < minimumSize {
					minimumSize = dirSize
				}
			}

			return nil
		})

		fmt.Printf("Smallest directory size which can satisfy our update: %d\n", minimumSize)
	}

	return nil
}
