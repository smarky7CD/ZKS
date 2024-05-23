package zks

// An enumerated set.
//
// Values in the set are integers and the boolean value they map to indicates set membership.
// True indicates membership and false indicates non-membership.
// The max value is the maximum value of the set.
// Any value not explicitly in the set or greater than the maximum value will return false.
type EnumSet struct {
	set map[uint64]bool
	max uint64
}

// Takes in (a possibly empty) integer to bool map and a maximum value. Returns a new EnumSet structure.
func NewEnumSet(set map[uint64]bool, max uint64) *EnumSet {
	return &EnumSet{set, max}
}

// Add a value x to the EnumSet if it is less than the maximum value.
func (es *EnumSet) Add(x uint64) {
	if x < es.max {
		es.set[x] = true
	}
}

// Remove a value x from the EnumSet if its less than the maximum value.
func (es *EnumSet) Remove(x uint64) {
	if x < es.max {
		es.set[x] = false
	}
}

// Respond true if x is in the Enum Set.
// Respond false if x is not in the Enum Set.
func (es *EnumSet) In(x uint64) bool {
	if x < es.max {
		v, ok := es.set[x]
		if ok {
			return v
		}
	}
	return false
}
