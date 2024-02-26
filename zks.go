package zks

import (
	"encoding/binary"

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
	tree *Tree
	set  EnumSet
}

type Com struct {
	com0 ristretto.Point
	com1 ristretto.Point
}

func Gen() *PubVerPar {
	h := mc.GeneratePublicParameters()
	kh, _ := keyset.NewHandle(prf.HMACSHA256PRFKeyTemplate())
	ps, _ := prf.NewPRFSet(kh)
	return &PubVerPar{h, *ps}
}

func Rep(pp *PubVerPar, es *EnumSet) *Repr {

	// inital commits to D represented as an enumerated set
	escommits := make(map[uint64]Com)
	for x := uint64(0); x < es.max; x++ {
		if es.In(x) {
			bx := make([]byte, 8)
			binary.PutUvarint(bx, x)
			ra0, _ := pp.ps.ComputePrimaryPRF(bx, 32)
			ra1, _ := pp.ps.ComputePrimaryPRF(ra0, 32)
			var r0, r1 ristretto.Scalar
			r0.Derive(ra0)
			r1.Derive(ra1)
			c0, c1 := mc.HardCommit(&pp.h, bx, &r0, &r1)
			escommits[x] = Com{c0, c1}
		}
		if !es.In(x) && es.In(x^1) {
			bx := make([]byte, 8)
			binary.PutUvarint(bx, x)
			ra0, _ := pp.ps.ComputePrimaryPRF(bx, 32)
			ra1, _ := pp.ps.ComputePrimaryPRF(ra0, 32)
			var r0, r1 ristretto.Scalar
			r0.Derive(ra0)
			r1.Derive(ra1)
			c0, c1 := mc.SoftCommit(&r0, &r1)
			escommits[x] = Com{c0, c1}
		}

	}

	tree := NewTree(pp, escommits, es.max)

	return &Repr{tree, *es}
}
