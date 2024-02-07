package zks

import (
	"github.com/bwesterb/go-ristretto"
)

type TreeNode struct {
	c0    ristretto.Point
	c1    ristretto.Point
	left  *TreeNode
	right *TreeNode
}

func NewNode(c0 ristretto.Point, c1 ristretto.Point, left *TreeNode, right *TreeNode) *TreeNode {
	return &TreeNode{c0, c1, left, right}
}
