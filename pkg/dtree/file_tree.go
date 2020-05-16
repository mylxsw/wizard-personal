package dtree

import (
	"sort"
	"strings"
	"time"

	"github.com/spf13/afero"
)

func BuildVueTreeForExt(fs afero.Fs, root string, ext ...string) (*VueTree, error) {
	tree, err := BuildTreeForExt(fs, root, ext...)
	if err != nil {
		return nil, err
	}

	return transformToVueTree(tree, root), nil
}

func BuildTreeForExt(fs afero.Fs, root string, ext ...string) (*Node, error) {
	tree, err := NewTree(fs, root, func(node *Node) bool {
		if node.FullPath == root {
			return false
		}

		relativePath := strings.TrimPrefix(node.FullPath, root)
		for _, p := range strings.Split(relativePath, "/") {
			if strings.HasPrefix(p, ".") {
				return false
			}
		}

		if node.Info.IsDir {
			return true
		}

		for _, ex := range ext {
			if strings.HasSuffix(node.Info.Name, ex) {
				return true
			}
		}

		return false
	})
	if err != nil {
		return nil, err
	}

	if tree == nil {
		tree = &Node{
			FullPath: "",
			Info: &FileInfo{
				Name:    "",
				Size:    0,
				Mode:    0,
				ModTime: time.Time{},
				IsDir:   true,
			},
			Children: make([]*Node, 0),
			Parent:   nil,
		}
	}

	sortTree(tree)
	return filterTree(tree), nil
}

func sortTree(tree *Node) {
	if tree == nil {
		return
	}

	sort.Sort(tree)
	for _, c := range tree.Children {
		sortTree(c)
	}
}

type VueTree struct {
	Title    string    `json:"title"`
	FullPath string    `json:"full_path"`
	Expand   bool      `json:"expand"`
	IsDir    bool      `json:"is_dir"`
	Children []VueTree `json:"children"`
}

func transformToVueTree(tree *Node, root string) *VueTree {
	if tree == nil {
		return nil
	}

	newTree := VueTree{
		Title:    tree.Info.Name,
		FullPath: strings.TrimPrefix(tree.FullPath, root),
		Expand:   false,
		IsDir:    tree.Info.IsDir,
		Children: make([]VueTree, 0),
	}

	if tree.Info.IsDir && len(tree.Children) == 0 {
		return &newTree
	}

	if tree.Children != nil && len(tree.Children) > 0 {
		for _, c := range tree.Children {
			if c.Info.IsDir && len(c.Children) == 0 {
				continue
			}

			res := transformToVueTree(c, root)
			if res == nil || (res.IsDir && len(res.Children) == 0) {
				continue
			}

			newTree.Children = append(newTree.Children, *res)
		}
	}

	return &newTree
}

func filterTree(tree *Node) *Node {
	if tree == nil {
		return nil
	}

	newTree := Node{
		FullPath: tree.FullPath,
		Info:     tree.Info,
		Children: make([]*Node, 0),
		Parent:   tree.Parent,
	}

	if tree.Info.IsDir && len(tree.Children) == 0 {
		return &newTree
	}

	if tree.Children != nil && len(tree.Children) > 0 {
		for _, c := range tree.Children {
			if c.Info.IsDir && len(c.Children) == 0 {
				continue
			}

			res := filterTree(c)
			if res == nil || (res.Info.IsDir && len(res.Children) == 0) {
				continue
			}

			newTree.Children = append(newTree.Children, res)
		}
	}

	return &newTree
}
