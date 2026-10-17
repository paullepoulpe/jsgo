package jutil

import "strconv"

// Integer is a boxed java.lang.Integer seen as a java.lang.Number (for
// example the value of JsonObject.addProperty(String, Number)). A Go int32
// stands for an Integer elsewhere; this named type is the Number.
//
// Written from the Javadoc of java.lang.Integer and java.lang.Number: the
// conversions are the JLS widening primitive conversions.
type Integer int32

// Java: Integer.intValue()
func (n Integer) IntValue() int32 { return int32(n) }

// Java: Integer.longValue()
func (n Integer) LongValue() int64 { return int64(n) }

// Java: Integer.doubleValue()
func (n Integer) DoubleValue() float64 { return float64(n) }

// Java: Integer.toString()
func (n Integer) String() string { return strconv.FormatInt(int64(n), 10) }

// Java: Integer.hashCode()
func (n Integer) HashCode() int32 { return int32(n) }

// Java: Integer.equals(Object)
func (n Integer) Equals(obj any) bool {
	other, ok := obj.(Integer)
	return ok && other == n
}
