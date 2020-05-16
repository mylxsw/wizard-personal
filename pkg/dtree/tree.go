package dtree

import (
	"github.com/spf13/afero"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileInfo is a struct created from os.FileInfo interface for serialization.
type FileInfo struct {
	Name    string      `json:"name"`
	Size    int64       `json:"size"`
	Mode    os.FileMode `json:"mode"`
	ModTime time.Time   `json:"mod_time"`
	IsDir   bool        `json:"is_dir"`
}

// Helper function to create a local FileInfo struct from os.FileInfo interface.
func fileInfoFromInterface(v os.FileInfo) *FileInfo {
	return &FileInfo{v.Name(), v.Size(), v.Mode(), v.ModTime(), v.IsDir()}
}

// Node represents a node in a directory tree.
type Node struct {
	FullPath string    `json:"path"`
	Info     *FileInfo `json:"info"`
	Children []*Node   `json:"children"`
	Parent   *Node     `json:"-"`
}

func (n Node) Len() int {
	return len(n.Children)
}

func (n Node) Less(i, j int) bool {
	if n.Children[i].Info.IsDir && !n.Children[j].Info.IsDir {
		return true
	}

	if !n.Children[i].Info.IsDir && n.Children[j].Info.IsDir {
		return false
	}

	return strings.Compare(n.Children[i].Info.Name, n.Children[j].Info.Name) < 0
}

func (n Node) Swap(i, j int) {
	n.Children[i], n.Children[j] = n.Children[j], n.Children[i]
}

// Create directory hierarchy.
func NewTree(fs afero.Fs, root string, filter func(node *Node) bool) (result *Node, err error) {
	parents := make(map[string]*Node)
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		node := &Node{
			FullPath: filepath.ToSlash(path),
			Info:     fileInfoFromInterface(info),
			Children: make([]*Node, 0),
		}

		if filter(node) {
			parents[path] = node
		}

		return nil
	}
	if err = afero.Walk(fs, root, walkFunc); err != nil {
		return
	}
	for path, node := range parents {
		parentPath := filepath.Dir(path)
		parent, exists := parents[parentPath]
		if !exists { // If a parent does not exist, this is the root.
			if result == nil {
				result = &Node{
					FullPath: "",
					Info:     &FileInfo{
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
			result.Children = append(result.Children, node)
		} else {
			node.Parent = parent
			parent.Children = append(parent.Children, node)
		}
	}
	return
}
