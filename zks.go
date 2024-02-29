package zks

import (
	"github.com/bwesterb/go-ristretto"
	"github.com/google/tink/go/keyset"
	"github.com/google/tink/go/prf"
	mc "github.com/smarky7CD/go-dl-mercurial-commitments"
)

type PubVerPar struct {
	h  ristretto.Point
	ps prf.Set
}

type Repr struct {
	tree Tree
	set  EnumSet
}

type Com struct {
	c0 ristretto.Point
	c1 ristretto.Point
}

type Answer struct {
	answer  bool
	levels  uint64
	xcoms   map[uint64]*Com
	sibcoms map[uint64]*Com
	opens   map[uint64]*Open
	teases  map[uint64]*Tease
}

func Gen() *PubVerPar {
	h := mc.GeneratePublicParameters()
	kh, _ := keyset.NewHandle(prf.HMACSHA256PRFKeyTemplate())
	ps, _ := prf.NewPRFSet(kh)
	return &PubVerPar{h, *ps}
}

func Rep(pp *PubVerPar, es *EnumSet) (*Repr, Com) {
	tree := NewTree(pp, es)
	return &Repr{*tree, *es}, Com{tree.root.c0, tree.root.c1}
}

func Qry(pp *PubVerPar, repr *Repr, x uint64) *Answer {
	return repr.tree.Path(pp, x, repr.set.In(x))
}

func Vfy(pp *PubVerPar, com Com, x uint64, answer *Answer) bool {
	return VerifyPath(pp, com, x, answer)
}
