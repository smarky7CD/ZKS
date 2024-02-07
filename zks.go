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
	tree TreeNode
	set  EnumSet
}

type Com struct {
	com0 ristretto.Point
	com1 ristretto.Point
}

func Gen() (ristretto.Point, prf.Set) {
	h := mc.GeneratePublicParameters()
	kh, _ := keyset.NewHandle(prf.HMACSHA256PRFKeyTemplate())
	ps, _ := prf.NewPRFSet(kh)
	return h, *ps
}
