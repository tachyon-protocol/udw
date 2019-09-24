package udwStrings

import (
	"fmt"
	"strings"
)

func GetShareCommonPrefix(s1 string, s2 string) string {
	pos := 0
	for {
		if pos >= len(s1) {
			break
		}
		if pos >= len(s2) {
			break
		}
		if s1[pos] != s2[pos] {
			break
		}
		pos++
	}

	return s1[:pos]
}

func NewRootPrefixNode() *prefixNode {
	return &prefixNode{
		Children: map[string]*prefixNode{},
	}
}

type prefixNode struct {
	Prefix   string
	Children map[string]*prefixNode
	Level    int
}

func (node *prefixNode) Print(infoCb func(prefix string) string) {
	var info string
	if infoCb != nil {
		info = infoCb(node.Prefix)
	}
	fmt.Println(strings.Repeat(" | ", node.Level), node.Prefix, info)
	for _, child := range node.Children {
		child.Print(infoCb)
	}
}

func (node *prefixNode) AddChild(child string) {
	if !strings.HasPrefix(child, node.Prefix) {
		return
	}
	if len(child) < len(node.Prefix)+1 {
		return
	}
	nextPrefix := child[:len(node.Prefix)+1]
	nextNode, ok := node.Children[nextPrefix]
	if !ok {
		nextNode = &prefixNode{
			Prefix:   nextPrefix,
			Children: map[string]*prefixNode{},
		}
		node.Children[nextPrefix] = nextNode
	}
	nextNode.AddChild(child)
}

func (node *prefixNode) MergeSingleChild() {
	if len(node.Children) == 1 {
		var grandChildren map[string]*prefixNode
		for _, child := range node.Children {
			node.Prefix = child.Prefix
			grandChildren = child.Children
		}
		node.Children = grandChildren
		node.MergeSingleChild()
	}
	for _, child := range node.Children {
		child.MergeSingleChild()
	}
}

func (node *prefixNode) CalculateLevel(start int) {
	node.Level = start + 1
	for _, child := range node.Children {
		child.CalculateLevel(node.Level)
	}
}
