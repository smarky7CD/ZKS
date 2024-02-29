package zks

import (
	"encoding/binary"
	"math"

	"github.com/bwesterb/go-ristretto"
	mc "github.com/smarky7CD/go-dl-mercurial-commitments"
)

type TreeNode struct {
	soft bool
	c0   ristretto.Point
	c1   ristretto.Point
	r0   ristretto.Scalar
	r1   ristretto.Scalar
}

func NewNode(soft bool, c0 ristretto.Point, c1 ristretto.Point, r0 ristretto.Scalar, r1 ristretto.Scalar) *TreeNode {
	return &TreeNode{soft, c0, c1, r0, r1}
}

type Tree struct {
	root   TreeNode
	tree   map[uint64]map[uint64]*TreeNode
	levels uint64
}

func ComputeNearestPowerof2(n uint64) uint64 {
	return uint64(math.Ceil(math.Log2(float64(n))))
}

func ComputeLeaves(pp *PubVerPar, es *EnumSet, level uint64) map[uint64]*TreeNode {
	var leaves = make(map[uint64]*TreeNode)

	for x := uint64(0); x < uint64(math.Pow(float64(level), 2)); x++ {

		if es.In(x) {
			bx := make([]byte, 8)
			bl := make([]byte, 8)
			binary.PutUvarint(bx, x)
			binary.PutUvarint(bl, level)
			bx = append(bx, bl...)
			ra0, _ := pp.ps.ComputePrimaryPRF(bx, 64)
			ra1, _ := pp.ps.ComputePrimaryPRF(ra0, 64)
			var r0, r1 ristretto.Scalar
			r0.Derive(ra0)
			r1.Derive(ra1)
			c0, c1 := mc.HardCommit(&pp.h, bx, &r0, &r1)
			leaves[x] = NewNode(false, c0, c1, r0, r1)
		}

		if !es.In(x) && es.In(x^1) {
			bx := make([]byte, 8)
			bl := make([]byte, 8)
			binary.PutUvarint(bx, x)
			binary.PutUvarint(bl, level)
			bx = append(bx, bl...)
			ra0, _ := pp.ps.ComputePrimaryPRF(bx, 64)
			ra1, _ := pp.ps.ComputePrimaryPRF(ra0, 64)
			var r0, r1 ristretto.Scalar
			r0.Derive(ra0)
			r1.Derive(ra1)
			c0, c1 := mc.SoftCommit(&r0, &r1)
			leaves[x] = NewNode(true, c0, c1, r0, r1)
		}
	}

	return leaves
}

func ComputeLayer(pp *PubVerPar, level uint64, prev_layer_nodes map[uint64]*TreeNode) map[uint64]*TreeNode {
	var layer_nodes = make(map[uint64]*TreeNode)

	for i := uint64(0); i < uint64(math.Pow(float64(level), 2)); i++ {
		val0, ok0 := prev_layer_nodes[2*i]
		val1, ok1 := prev_layer_nodes[(2*i)+1]

		if ok0 && ok1 {
			bsigma := val0.c0.Bytes()
			bsigma = append(bsigma, val0.c1.Bytes()...)
			bsigma = append(bsigma, val1.c0.Bytes()...)
			bsigma = append(bsigma, val1.c1.Bytes()...)
			bx := make([]byte, 8)
			bl := make([]byte, 8)
			binary.PutUvarint(bx, i)
			binary.PutUvarint(bl, level)
			bx = append(bx, bl...)
			ra0, _ := pp.ps.ComputePrimaryPRF(bx, 64)
			ra1, _ := pp.ps.ComputePrimaryPRF(ra0, 64)
			var r0, r1 ristretto.Scalar
			r0.Derive(ra0)
			r1.Derive(ra1)
			c0, c1 := mc.HardCommit(&pp.h, bsigma, &r0, &r1)
			layer_nodes[i] = NewNode(false, c0, c1, r0, r1)
		}

		if ok0 != ok1 {
			bx := make([]byte, 8)
			bl := make([]byte, 8)
			binary.PutUvarint(bx, i)
			binary.PutUvarint(bl, level)
			bx = append(bx, bl...)
			ra0, _ := pp.ps.ComputePrimaryPRF(bx, 64)
			ra1, _ := pp.ps.ComputePrimaryPRF(ra0, 64)
			var r0, r1 ristretto.Scalar
			r0.Derive(ra0)
			r1.Derive(ra1)
			c0, c1 := mc.SoftCommit(&r0, &r1)
			layer_nodes[i] = NewNode(true, c0, c1, r0, r1)
		}

	}

	return layer_nodes

}

func NewTree(pp *PubVerPar, es *EnumSet) *Tree {
	levels := ComputeNearestPowerof2(es.max)
	var tree = make(map[uint64]map[uint64]*TreeNode)

	// compute the leaves of the tree
	leaves := ComputeLeaves(pp, es, levels)
	tree[levels] = leaves

	// build the tree in a bottom up fashion
	prev_layer_nodes := leaves
	for i := int(levels) - 1; i >= 0; i-- {
		layer_nodes := ComputeLayer(pp, uint64(i), prev_layer_nodes)
		tree[uint64(i)] = layer_nodes
		prev_layer_nodes = layer_nodes
	}

	// check for nil root
	if len(tree[0]) == 0 {
		bx := make([]byte, 8)
		bl := make([]byte, 8)
		binary.PutUvarint(bx, 0)
		binary.PutUvarint(bl, 0)
		bx = append(bx, bl...)
		ra0, _ := pp.ps.ComputePrimaryPRF(bx, 64)
		ra1, _ := pp.ps.ComputePrimaryPRF(ra0, 64)
		var r0, r1 ristretto.Scalar
		r0.Derive(ra0)
		r1.Derive(ra1)
		c0, c1 := mc.SoftCommit(&r0, &r1)
		tree[0][0] = NewNode(true, c0, c1, r0, r1)
	}

	return &Tree{*tree[0][0], tree, levels}
}

type Open struct {
	r0 ristretto.Scalar
	r1 ristretto.Scalar
}

type Tease = ristretto.Scalar

func MemberPath(tree *Tree, pp *PubVerPar, x uint64) *Answer {
	var opens = make(map[uint64]*Open)
	var xcoms = make(map[uint64]*Com)
	var sibcoms = make(map[uint64]*Com)
	var teases = make(map[uint64]*Tease)
	for i := uint64(0); i < tree.levels; i++ {
		j := tree.levels - i
		xi := x >> i
		opens[j] = &Open{tree.tree[j][xi].r0, tree.tree[j][xi].r1}
		if i >= 1 {
			xcoms[j] = &Com{tree.tree[j][xi].c0, tree.tree[j][xi].c1}
			sibcoms[j] = &Com{tree.tree[j][xi^1].c0, tree.tree[j][xi^1].c1}
		}
	}
	return &Answer{true, tree.levels, xcoms, sibcoms, opens, teases}
}

func NonMemberPath(tree *Tree, pp *PubVerPar, x uint64) *Answer {
	// refine tree if necessary
	for i := uint64(0); i < tree.levels; i++ {
		j := tree.levels - i
		xi := x >> i
		_, ok := tree.tree[j][xi]
		if !ok {
			if i == 0 {
				bx := make([]byte, 8)
				bl := make([]byte, 8)
				binary.PutUvarint(bx, xi)
				binary.PutUvarint(bl, j)
				bx = append(bx, bl...)
				ra0, _ := pp.ps.ComputePrimaryPRF(bx, 64)
				ra1, _ := pp.ps.ComputePrimaryPRF(ra0, 64)
				var r0, r1 ristretto.Scalar
				r0.Derive(ra0)
				r1.Derive(ra1)
				c0, c1 := mc.HardCommit(&pp.h, []byte("bot"), &r0, &r1)
				tree.tree[j][xi] = NewNode(false, c0, c1, r0, r1)
			}
			bxip := make([]byte, 8)
			bl := make([]byte, 8)
			binary.PutUvarint(bxip, xi^1)
			binary.PutUvarint(bl, j)
			bxip = append(bxip, bl...)
			ra0, _ := pp.ps.ComputePrimaryPRF(bxip, 64)
			ra1, _ := pp.ps.ComputePrimaryPRF(ra0, 64)
			var r0, r1 ristretto.Scalar
			r0.Derive(ra0)
			r1.Derive(ra1)
			c0, c1 := mc.SoftCommit(&r0, &r1)
			tree.tree[j][xi] = NewNode(true, c0, c1, r0, r1)
		}
	}

	// build answer
	var opens = make(map[uint64]*Open)
	var xcoms = make(map[uint64]*Com)
	var sibcoms = make(map[uint64]*Com)
	var teases = make(map[uint64]*Tease)
	for i := uint64(0); i < tree.levels; i++ {
		j := tree.levels - i
		xi := x >> i
		val := tree.tree[j][xi]
		var r ristretto.Scalar
		if val.soft {
			bxi := make([]byte, 8)
			bl := make([]byte, 8)
			binary.PutUvarint(bxi, xi)
			binary.PutUvarint(bl, j)
			bxi = append(bxi, bl...)
			r = mc.SoftTease(bxi, &val.r0, &val.r1)
		} else {
			r = tree.tree[j][xi].r0
		}
		teases[j] = &r
		if i >= 1 {
			xcoms[j] = &Com{tree.tree[j][xi].c0, tree.tree[j][xi].c1}
			sibcoms[j] = &Com{tree.tree[j][xi^1].c0, tree.tree[j][xi^1].c1}
		}
	}
	return &Answer{false, tree.levels, xcoms, sibcoms, opens, teases}
}

func (tree *Tree) Path(pp *PubVerPar, x uint64, a bool) *Answer {
	if a {
		return MemberPath(tree, pp, x)
	} else {
		return NonMemberPath(tree, pp, x)
	}
}

func VerifyOpen(pp *PubVerPar, com Com, x uint64, answer *Answer) bool {
	// verify all internal tree nodes
	for i := uint64(1); i < answer.levels; i++ {
		c := answer.xcoms[i]
		pi := answer.opens[i]
		vx := answer.xcoms[i+1]
		vs := answer.sibcoms[i+1]
		vxleft := ((x & (1 << ((answer.levels - i + 1) - 1))) >> ((answer.levels - i + 1) - 1))
		if vxleft == 0 {
			bsigma := vx.c0.Bytes()
			bsigma = append(bsigma, vx.c1.Bytes()...)
			bsigma = append(bsigma, vs.c0.Bytes()...)
			bsigma = append(bsigma, vs.c1.Bytes()...)
			if !mc.VerOpen(&pp.h, &c.c0, &c.c1, bsigma, &pi.r0, &pi.r1) {
				return false
			}
		} else {
			bsigma := vs.c0.Bytes()
			bsigma = append(bsigma, vs.c1.Bytes()...)
			bsigma = append(bsigma, vx.c0.Bytes()...)
			bsigma = append(bsigma, vx.c1.Bytes()...)
			if !mc.VerOpen(&pp.h, &c.c0, &c.c1, bsigma, &pi.r0, &pi.r1) {
				return false
			}
		}
	}

	// check root commit
	pi := answer.opens[0]
	vx := answer.xcoms[1]
	vs := answer.sibcoms[1]
	vxleft := ((x & (1 << ((answer.levels - 1) - 1))) >> ((answer.levels - 1) - 1))
	if vxleft == 0 {
		bsigma := vx.c0.Bytes()
		bsigma = append(bsigma, vx.c1.Bytes()...)
		bsigma = append(bsigma, vs.c0.Bytes()...)
		bsigma = append(bsigma, vs.c1.Bytes()...)
		if !mc.VerOpen(&pp.h, &com.c0, &com.c1, bsigma, &pi.r0, &pi.r1) {
			return false
		}
	} else {
		bsigma := vs.c0.Bytes()
		bsigma = append(bsigma, vs.c1.Bytes()...)
		bsigma = append(bsigma, vx.c0.Bytes()...)
		bsigma = append(bsigma, vx.c1.Bytes()...)
		if !mc.VerOpen(&pp.h, &com.c0, &com.c1, bsigma, &pi.r0, &pi.r1) {
			return false
		}
	}

	return true
}

func VerifyTease(com Com, x uint64, answer *Answer) bool {
	// verify all internal tree nodes
	for i := uint64(1); i < answer.levels; i++ {
		c := answer.xcoms[i]
		tau := answer.teases[i]
		vx := answer.xcoms[i+1]
		vs := answer.sibcoms[i+1]
		vxleft := ((x & (1 << ((answer.levels - i + 1) - 1))) >> ((answer.levels - i + 1) - 1))
		if vxleft == 0 {
			bsigma := vx.c0.Bytes()
			bsigma = append(bsigma, vx.c1.Bytes()...)
			bsigma = append(bsigma, vs.c0.Bytes()...)
			bsigma = append(bsigma, vs.c1.Bytes()...)
			if !mc.VerTease(&c.c0, &c.c1, bsigma, tau) {
				return false
			}
		} else {
			bsigma := vs.c0.Bytes()
			bsigma = append(bsigma, vs.c1.Bytes()...)
			bsigma = append(bsigma, vx.c0.Bytes()...)
			bsigma = append(bsigma, vx.c1.Bytes()...)
			if !mc.VerTease(&c.c0, &c.c1, bsigma, tau) {
				return false
			}
		}
	}

	// check root commit
	tau := answer.teases[1]
	vx := answer.xcoms[1]
	vs := answer.sibcoms[1]
	vxleft := ((x & (1 << ((answer.levels - 1) - 1))) >> ((answer.levels - 1) - 1))
	if vxleft == 0 {
		bsigma := vx.c0.Bytes()
		bsigma = append(bsigma, vx.c1.Bytes()...)
		bsigma = append(bsigma, vs.c0.Bytes()...)
		bsigma = append(bsigma, vs.c1.Bytes()...)
		if !mc.VerTease(&com.c0, &com.c1, bsigma, tau) {
			return false
		}
	} else {
		bsigma := vs.c0.Bytes()
		bsigma = append(bsigma, vs.c1.Bytes()...)
		bsigma = append(bsigma, vx.c0.Bytes()...)
		bsigma = append(bsigma, vx.c1.Bytes()...)
		if !mc.VerTease(&com.c0, &com.c1, bsigma, tau) {
			return false
		}
	}

	// check x commit
	cx := answer.xcoms[answer.levels]
	taux := answer.teases[answer.levels]
	return mc.VerTease(&cx.c0, &cx.c1, []byte("bot"), taux)
}

func VerifyPath(pp *PubVerPar, com Com, x uint64, answer *Answer) bool {
	if answer.answer {
		return VerifyOpen(pp, com, x, answer)
	} else {
		return VerifyTease(com, x, answer)
	}
}
