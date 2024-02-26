package zks

import (
	"math"

	"github.com/bwesterb/go-ristretto"
	mc "github.com/smarky7CD/go-dl-mercurial-commitments"
)

type TreeNode struct {
	c0     ristretto.Point
	c1     ristretto.Point
	parent *TreeNode
	left   *TreeNode
	right  *TreeNode
}

func NewNode(c0 ristretto.Point, c1 ristretto.Point, parent *TreeNode, left *TreeNode, right *TreeNode) *TreeNode {
	return &TreeNode{c0, c1, parent, left, right}
}

type Tree struct {
	root   *TreeNode
	leaves map[uint64]*TreeNode
	max    uint64
	np2    uint64
}

func ComputeNearestPowerof2(n uint64) uint64 {
	return uint64(math.Ceil(math.Log2(float64(n))))
}

func NewTree(pp *PubVerPar, escommits map[uint64]Com, max uint64) *Tree {
	np2 := ComputeNearestPowerof2(max)

	// handle empty tree
	if len(escommits) == 0 {
		ra0, _ := pp.ps.ComputePrimaryPRF([]byte("epsilon"), 32)
		ra1, _ := pp.ps.ComputePrimaryPRF(ra0, 32)
		var r0, r1 ristretto.Scalar
		r0.Derive(ra0)
		r1.Derive(ra1)
		c0, c1 := mc.SoftCommit(&r0, &r1)
		root := NewNode(c0, c1, nil, nil, nil)
		return &Tree{root, nil, 0, np2}
	}

	// build leaves
	leaves := make(map[uint64]*TreeNode)
	for i := uint64(0); i < uint64(math.Pow(float64(np2), 2)); i++ {
		if val, ok := escommits[i]; ok {
			leaves[i] = NewNode(val.com0, val.com1, nil, nil, nil)
		}
	}

	return NewTreeHelper(leaves, np2, leaves, max, np2, pp)
}

func NewTreeHelper(curnodes map[uint64]*TreeNode, level uint64, leaves map[uint64]*TreeNode, max uint64, np2 uint64, pp *PubVerPar) *Tree {

	if level == 0 {
		// base case -- do not have to worry about empty tree as already handled
		val0, ok0 := curnodes[0]
		val1, ok1 := curnodes[1]
		// 2 non-nil nodes
		if ok0 && ok1 {
			bsigma := val0.c0.Bytes()
			bsigma = append(bsigma, val0.c1.Bytes()...)
			bsigma = append(bsigma, val1.c0.Bytes()...)
			bsigma = append(bsigma, val1.c1.Bytes()...)
			rsigma0, _ := pp.ps.ComputePrimaryPRF(bsigma, 32)
			rsigma1, _ := pp.ps.ComputePrimaryPRF(rsigma0, 32)
			var r0, r1 ristretto.Scalar
			r0.Derive(rsigma0)
			r1.Derive(rsigma1)
			c0, c1 := mc.HardCommit(&pp.h, bsigma, &r0, &r1)
			root := NewNode(c0, c1, nil, curnodes[0], curnodes[1])
			curnodes[0].parent = root
			curnodes[1].parent = root
			return &Tree{root, leaves, max, np2}
			// only 1 non-nil node
		} else if ok0 && !ok1 {
			bsigma := val0.c0.Bytes()
			bsigma = append(bsigma, val0.c1.Bytes()...)
			rsigma0, _ := pp.ps.ComputePrimaryPRF(bsigma, 32)
			rsigma1, _ := pp.ps.ComputePrimaryPRF(rsigma0, 32)
			var r0, r1 ristretto.Scalar
			r0.Derive(rsigma0)
			r1.Derive(rsigma1)
			c0, c1 := mc.SoftCommit(&r0, &r1)
			root := NewNode(c0, c1, nil, curnodes[0], nil)
			curnodes[0].parent = root
			return &Tree{root, leaves, max, np2}

		} else if !ok0 && ok1 {
			bsigma := val1.c0.Bytes()
			bsigma = append(bsigma, val1.c1.Bytes()...)
			rsigma0, _ := pp.ps.ComputePrimaryPRF(bsigma, 32)
			rsigma1, _ := pp.ps.ComputePrimaryPRF(rsigma0, 32)
			var r0, r1 ristretto.Scalar
			r0.Derive(rsigma0)
			r1.Derive(rsigma1)
			c0, c1 := mc.SoftCommit(&r0, &r1)
			root := NewNode(c0, c1, nil, nil, curnodes[1])
			curnodes[1].parent = root
			return &Tree{root, leaves, max, np2}
		} else {
			// something bad has happened
			var p0, p1 ristretto.Point
			badroot := NewNode(*p0.SetBase(), *p1.SetBase(), nil, nil, nil)
			return &Tree{badroot, nil, 0, 0}
		}
	}

	nextnodes := make(map[uint64]*TreeNode)
	for i := uint64(0); i < uint64(math.Pow(float64(level), 2)); i = i + 2 {
		j := i / 2
		val0, ok0 := curnodes[i]
		val1, ok1 := curnodes[1+1]
		// 2 non-nil nodes
		if ok0 && ok1 {
			bsigma := val0.c0.Bytes()
			bsigma = append(bsigma, val0.c1.Bytes()...)
			bsigma = append(bsigma, val1.c0.Bytes()...)
			bsigma = append(bsigma, val1.c1.Bytes()...)
			rsigma0, _ := pp.ps.ComputePrimaryPRF(bsigma, 32)
			rsigma1, _ := pp.ps.ComputePrimaryPRF(rsigma0, 32)
			var r0, r1 ristretto.Scalar
			r0.Derive(rsigma0)
			r1.Derive(rsigma1)
			c0, c1 := mc.HardCommit(&pp.h, bsigma, &r0, &r1)
			nextnodes[j] = NewNode(c0, c1, nil, curnodes[i], curnodes[i+1])
			curnodes[i].parent = nextnodes[j]
			curnodes[i+1].parent = nextnodes[j]
		}
		// only 1 non-nil node
		if ok0 && !ok1 {
			bsigma := val0.c0.Bytes()
			bsigma = append(bsigma, val0.c1.Bytes()...)
			rsigma0, _ := pp.ps.ComputePrimaryPRF(bsigma, 32)
			rsigma1, _ := pp.ps.ComputePrimaryPRF(rsigma0, 32)
			var r0, r1 ristretto.Scalar
			r0.Derive(rsigma0)
			r1.Derive(rsigma1)
			c0, c1 := mc.SoftCommit(&r0, &r1)
			nextnodes[j] = NewNode(c0, c1, nil, curnodes[i], nil)
			curnodes[i].parent = nextnodes[j]
		}

		if !ok0 && ok1 {
			bsigma := val1.c0.Bytes()
			bsigma = append(bsigma, val1.c1.Bytes()...)
			rsigma0, _ := pp.ps.ComputePrimaryPRF(bsigma, 32)
			rsigma1, _ := pp.ps.ComputePrimaryPRF(rsigma0, 32)
			var r0, r1 ristretto.Scalar
			r0.Derive(rsigma0)
			r1.Derive(rsigma1)
			c0, c1 := mc.SoftCommit(&r0, &r1)
			nextnodes[j] = NewNode(c0, c1, nil, nil, curnodes[i+1])
			curnodes[i+1].parent = nextnodes[j]
		}
	}

	// update leaves only if this is the first bottom up call
	//  implicitly done in the base case (|D| == 1)
	if level == np2 {
		return NewTreeHelper(nextnodes, level-1, curnodes, max, np2, pp)
	}

	return NewTreeHelper(nextnodes, level-1, leaves, max, np2, pp)
}
