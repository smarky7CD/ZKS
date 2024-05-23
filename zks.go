package zks

import (
	"github.com/bwesterb/go-ristretto"
	"github.com/google/tink/go/keyset"
	"github.com/google/tink/go/prf"
	mc "github.com/smarky7CD/go-dl-mercurial-commitments"
)

// h is the randomly selected point on the EC used for the commitment scheme
// ps is the randomly selected PRF
type PubVerPar struct {
	h  ristretto.Point
	ps prf.Set
}

// A ZKS representation is the tree and the underlying EnumSet.
type Repr struct {
	tree Tree
	set  EnumSet
}

// A commitment is two points on the EC.
type Com struct {
	c0 ristretto.Point
	c1 ristretto.Point
}

// An answer contains the boolean set-membership reply and information used in the proof.
type Answer struct {
	answer  bool
	levels  uint64
	xcoms   map[uint64]*Com
	sibcoms map[uint64]*Com
	opens   map[uint64]*Open
	teases  map[uint64]*Tease
}

// Generate h (value used for commitments) and *ps (the PRF).
func Gen() *PubVerPar {
	h := mc.GeneratePublicParameters()
	kh, _ := keyset.NewHandle(prf.HMACSHA256PRFKeyTemplate())
	ps, _ := prf.NewPRFSet(kh)
	return &PubVerPar{h, *ps}
}

// Input: public parameters (h,ps) and an EnumSet.
// Return: ZKS representation and a commitment to it.
func Rep(pp *PubVerPar, es *EnumSet) (*Repr, Com) {
	tree := NewTree(pp, es)
	return &Repr{*tree, *es}, Com{tree.root.c0, tree.root.c1}
}

// Input: The public parameters (h,ps), a ZKS representation, and an element x.
// Return: Answer struct containing set-membership response and a proof.
func Qry(pp *PubVerPar, repr *Repr, x uint64) *Answer {
	return repr.tree.Path(pp, x, repr.set.In(x))
}

// Input: The public parameters (h,ps), a ZKS commitment, an element x that was queried, and the answer/proof struct.
// Return: True if answer verifies, false otherwise.
func Vfy(pp *PubVerPar, com Com, x uint64, answer *Answer) bool {
	return VerifyPath(pp, com, x, answer)
}
