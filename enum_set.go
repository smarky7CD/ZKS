package zks

type EnumSet struct {
	set map[uint64]bool
	max uint64
}

func NewEnumSet(set map[uint64]bool, max uint64) *EnumSet {
	return &EnumSet{set, max}
}

func (es *EnumSet) Add(x uint64) {
	if x < es.max {
		es.set[x] = true
	}
}

func (es *EnumSet) Remove(x uint64) {
	if x < es.max {
		es.set[x] = false
	}
}

func (es *EnumSet) In(x uint64) bool {
	if x < es.max {
		v := es.set[x]
		return v
	}
	return false
}
