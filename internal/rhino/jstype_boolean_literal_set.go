// Port of com/google/javascript/rhino/jstype/BooleanLiteralSet.java @ closure-compiler v20260415

package rhino

import (
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// BooleanLiteralSet is a set in the domain {true,false}. There are four
// possible sets: {}, {true}, {false}, {true,false}.
type BooleanLiteralSet int32

const (
	BooleanLiteralSet_EMPTY BooleanLiteralSet = iota
	BooleanLiteralSet_TRUE
	BooleanLiteralSet_FALSE
	BooleanLiteralSet_BOTH
)

var booleanLiteralSet_names = [...]string{"EMPTY", "TRUE", "FALSE", "BOTH"}

// BooleanLiteralSet_values is BooleanLiteralSet.values().
func BooleanLiteralSet_values() []BooleanLiteralSet {
	return []BooleanLiteralSet{
		BooleanLiteralSet_EMPTY, BooleanLiteralSet_TRUE, BooleanLiteralSet_FALSE, BooleanLiteralSet_BOTH,
	}
}

// String is Enum.name().
func (s BooleanLiteralSet) String() string { return booleanLiteralSet_names[s] }

// Ordinal is Enum.ordinal().
func (s BooleanLiteralSet) Ordinal() int { return int(s) }

// Java: BooleanLiteralSet.fromOrdinal(int)
func (s BooleanLiteralSet) FromOrdinal(ordinal int) BooleanLiteralSet {
	switch ordinal {
	case 0:
		return BooleanLiteralSet_EMPTY
	case 1:
		return BooleanLiteralSet_TRUE
	case 2:
		return BooleanLiteralSet_FALSE
	case 3:
		return BooleanLiteralSet_BOTH
	}
	panic(jutil.NewIllegalArgumentExceptionMsg("Ordinal: " + jutil.ToJavaString(ordinal)))
}

// Intersection computes the intersection of this set and that.
//
// Java: BooleanLiteralSet.intersection(BooleanLiteralSet)
func (s BooleanLiteralSet) Intersection(that BooleanLiteralSet) BooleanLiteralSet {
	return s.FromOrdinal(s.Ordinal() & that.Ordinal())
}

// Union computes the union of this set and that.
//
// Java: BooleanLiteralSet.union(BooleanLiteralSet)
func (s BooleanLiteralSet) Union(that BooleanLiteralSet) BooleanLiteralSet {
	return s.FromOrdinal(s.Ordinal() | that.Ordinal())
}

// Contains returns whether this contains the given literal value.
//
// Java: BooleanLiteralSet.contains(boolean)
func (s BooleanLiteralSet) Contains(literalValue bool) bool {
	switch s.Ordinal() {
	case 0:
		return false
	case 1:
		return literalValue
	case 2:
		return !literalValue
	case 3:
		return true
	}
	panic(jutil.NewIndexOutOfBoundsException("Ordinal: " + jutil.ToJavaString(s.Ordinal())))
}

// BooleanLiteralSet_Get returns the singleton set {literalValue}.
//
// Java: BooleanLiteralSet.get(boolean)
func BooleanLiteralSet_Get(literalValue bool) BooleanLiteralSet {
	if literalValue {
		return BooleanLiteralSet_TRUE
	}
	return BooleanLiteralSet_FALSE
}

// ToTri converts to a Tri.
//
// Java: BooleanLiteralSet.toTri()
func (s BooleanLiteralSet) ToTri() base.Tri {
	if s == BooleanLiteralSet_TRUE {
		return base.Tri_TRUE
	} else if s == BooleanLiteralSet_FALSE {
		return base.Tri_FALSE
	}
	return base.Tri_UNKNOWN
}
