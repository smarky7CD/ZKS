package zks

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func randRange(min, max int) int {
	return rand.Intn(max-min) + min
}

func TestCorrectnessRand(t *testing.T) {

	max_value := uint64(randRange(8, 4096))
	mem_prob := rand.Float64()

	var values = make(map[uint64]bool)

	for i := uint64(0); i < max_value; i++ {
		r := rand.Float64()
		if r <= mem_prob {
			values[i] = true
		} else {
			values[i] = false
		}
	}

	fmt.Println("ZKS of universe size ", max_value, " with membership probability ", mem_prob, ".")

	set := NewEnumSet(values, max_value)

	pp := Gen()

	repr, com := Rep(pp, set)

	for i := uint64(0); i < max_value; i++ {
		a := Qry(pp, repr, i)
		v := Vfy(pp, com, i, a)
		fmt.Println(i)
		assert.True(t, v, "v should be true.")
	}

}

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
		a := Qry(pp, repr, i)
		v := Vfy(pp, com, i, a)
		fmt.Println(i)
		assert.True(t, v, "v should be true.")
	}
}

func TestCorrectness2(t *testing.T) {

	values := map[uint64]bool{
		0:  false,
		1:  false,
		2:  false,
		3:  false,
		4:  false,
		5:  false,
		6:  false,
		7:  false,
		8:  true,
		9:  false,
		10: false,
		11: false,
		12: false,
		13: false,
		14: false,
		15: true,
		16: false,
		17: true,
		18: true,
		19: false,
		20: false,
		21: false,
		22: false,
		23: false,
		24: false,
		25: false,
		26: false,
		27: false,
		28: false,
		29: false,
		30: false,
		31: true,
	}

	set := NewEnumSet(values, 32)

	pp := Gen()

	repr, com := Rep(pp, set)

	for i := uint64(0); i < 32; i++ {
		a := Qry(pp, repr, i)
		v := Vfy(pp, com, i, a)
		fmt.Println(i)
		assert.True(t, v, "v should be true.")

	}
}

func TestCorrectness3(t *testing.T) {

	values := map[uint64]bool{
		0:   false,
		1:   false,
		2:   false,
		3:   true,
		4:   false,
		5:   false,
		6:   false,
		7:   false,
		8:   false,
		9:   false,
		10:  true,
		11:  false,
		12:  false,
		13:  true,
		14:  false,
		15:  false,
		16:  false,
		17:  false,
		18:  false,
		19:  false,
		20:  false,
		21:  false,
		22:  false,
		23:  false,
		24:  true,
		25:  false,
		26:  false,
		27:  true,
		28:  true,
		29:  false,
		30:  false,
		31:  false,
		32:  true,
		33:  false,
		34:  false,
		35:  false,
		36:  false,
		37:  false,
		38:  false,
		39:  true,
		40:  false,
		41:  false,
		42:  false,
		43:  false,
		44:  true,
		45:  false,
		46:  false,
		47:  false,
		48:  false,
		49:  false,
		50:  true,
		51:  true,
		52:  false,
		53:  false,
		54:  false,
		55:  false,
		56:  false,
		57:  false,
		58:  false,
		59:  false,
		60:  false,
		61:  false,
		62:  false,
		63:  false,
		64:  false,
		65:  false,
		66:  false,
		67:  false,
		68:  true,
		69:  false,
		70:  false,
		71:  true,
		72:  false,
		73:  true,
		74:  false,
		75:  false,
		76:  true,
		77:  false,
		78:  false,
		79:  false,
		80:  true,
		81:  false,
		82:  false,
		83:  false,
		84:  true,
		85:  true,
		86:  false,
		87:  false,
		88:  false,
		89:  false,
		90:  true,
		91:  false,
		92:  true,
		93:  false,
		94:  false,
		95:  true,
		96:  false,
		97:  false,
		98:  false,
		99:  true,
		100: false,
		101: false,
		102: false,
		103: true,
		104: false,
		105: true,
		106: false,
		107: false,
		108: false,
		109: false,
		110: false,
		111: false,
		112: false,
		113: false,
		114: true,
		115: true,
		116: true,
		117: true,
		118: false,
		119: false,
		120: false,
		121: false,
		122: false,
		123: false,
		124: false,
		125: true,
		126: true,
		127: false,
		128: true,
		129: true,
		130: false,
		131: false,
		132: false,
		133: true,
		134: false,
		135: false,
		136: false,
		137: true,
		138: false,
		139: false,
		140: false,
		141: true,
		142: false,
		143: false,
		144: false,
		145: false,
		146: false,
		147: false,
		148: true,
		149: false,
		150: false,
		151: false,
		152: true,
		153: false,
		154: true,
		155: false,
		156: false,
		157: true,
		158: false,
		159: false,
		160: true,
		161: true,
		162: false,
		163: true,
		164: false,
		165: false,
		166: false,
		167: false,
		168: true,
		169: false,
		170: false,
		171: false,
		172: false,
		173: false,
		174: true,
		175: false,
		176: false,
		177: false,
		178: false,
		179: true,
		180: false,
		181: true,
		182: false,
		183: true,
		184: false,
		185: false,
		186: false,
		187: false,
		188: false,
		189: false,
		190: true,
		191: false,
		192: true,
		193: false,
		194: false,
		195: false,
		196: false,
		197: false,
		198: false,
		199: false,
		200: false,
		201: false,
		202: false,
		203: true,
		204: true,
		205: false,
		206: false,
		207: false,
		208: false,
		209: false,
		210: false,
		211: false,
		212: false,
		213: false,
		214: false,
		215: false,
		216: true,
		217: true,
		218: false,
		219: false,
		220: false,
		221: true,
		222: false,
		223: true,
		224: true,
		225: true,
		226: true,
		227: false,
		228: true,
		229: false,
		230: true,
		231: true,
		232: false,
		233: false,
		234: false,
		235: false,
		236: false,
		237: false,
		238: true,
		239: false,
		240: true,
		241: false,
		242: true,
		243: false,
		244: false,
		245: false,
		246: true,
		247: false,
		248: true,
		249: false,
		250: true,
		251: false,
		252: true,
		253: false,
		254: false,
		255: false,
	}

	set := NewEnumSet(values, 256)

	pp := Gen()

	repr, com := Rep(pp, set)

	for i := uint64(0); i < 256; i++ {
		a := Qry(pp, repr, i)
		v := Vfy(pp, com, i, a)
		fmt.Println(i)
		assert.True(t, v, "v should be true.")
	}
}
