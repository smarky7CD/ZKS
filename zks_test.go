package zks

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCorrectness(t *testing.T) {

	values := map[uint64]bool{
		0:  true,
		1:  false,
		2:  false,
		3:  true,
		4:  true,
		5:  true,
		6:  false,
		7:  false,
		8:  false,
		9:  false,
		10: true,
		11: true,
		12: true,
		13: true,
		14: true,
		15: true,
	}

	set := NewEnumSet(values, 16)

	pp := Gen()

	repr, com := Rep(pp, set)

	for i := uint64(0); i < 16; i++ {
		val := values[i]
		/*
			if !val {
				a := Qry(pp, repr, i)
				v := Vfy(pp, com, i, a)
				fmt.Println(i)
				assert.True(t, v, "v should be true.")
			}
		*/
		if val {
			a := Qry(pp, repr, i)
			v := Vfy(pp, com, i, a)
			fmt.Println(i)
			assert.True(t, v, "v should be true.")

		}
	}

}
