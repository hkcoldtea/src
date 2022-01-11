package util

import (
	"fmt"
	"sync"
)

type Node struct {
	Url      string
	Children map[string]*Node
	mu  sync.RWMutex
}

func (n *Node) String() string {
	return fmt.Sprintf("Node{Url=%s, with %d Children}", n.Url, len(n.Children))
}

func (n *Node) Insert(url string, parts []string, index int) {
	n.mu.RLock()
	defer n.mu.RUnlock()
	if len(parts) == index {
		n.Url = url
		return
	}
	part := parts[index]
	child := n.match(part)
	if child == nil {
		child = &Node{
			Children: make(map[string]*Node),
		}
		n.Children[part] = child
	}
	child.Insert(url, parts, index + 1)
}

func (n *Node) match(part string) *Node {
	return n.Children[part]
}

func (n *Node) Search(parts []string, index int) string {
	if index == len(parts) {
		return n.Url
	}
	part := parts[index]
	child := n.Children[part]
	if child == nil {
		for k, v := range n.Children {
			if k[0] == ':' {
				searchRes := v.Search(parts, index + 1)
				if searchRes != "" {
					return searchRes
				}
			} else {
				if k[0] == '*' {
					return v.Url
				}
			}
		}
		return ""
	}
	return child.Search(parts, index + 1)
}
