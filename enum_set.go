package zks

type EnumSet struct {
	set map[uint]bool
	max uint
}

func NewEnumSet(set map[uint]bool, max uint) *EnumSet {
	return &EnumSet{set, max}
}
