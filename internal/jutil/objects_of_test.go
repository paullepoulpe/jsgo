package jutil

import "testing"

type objectsOfKey struct{ id int }

func (k objectsOfKey) Equals(o any) bool { x, ok := o.(objectsOfKey); return ok && x.id == k.id }
func (k objectsOfKey) HashCode() int32   { return int32(k.id * 31) }

func checkObjectsOf[T any](t *testing.T, vals []T) {
	t.Helper()
	for _, a := range vals {
		if got, want := ObjectHashCodeOf(a), ObjectHashCode(a); got != want {
			t.Errorf("ObjectHashCodeOf(%v) = %d, ObjectHashCode %d", a, got, want)
		}
		for _, b := range vals {
			if got, want := ObjectsEqualsOf(a, b), ObjectsEquals(a, b); got != want {
				t.Errorf("ObjectsEqualsOf(%v, %v) = %v, ObjectsEquals %v", a, b, got, want)
			}
		}
	}
}

func TestObjectsEqualsOfAndHashCodeOfMatchObjects(t *testing.T) {
	checkObjectsOf(t, []string{"", "a", "b", "ab", "é", "\U0001F600", string([]byte{0xed, 0xa0, 0x80})})
	checkObjectsOf(t, []int32{0, 1, -1, 1 << 30})
	checkObjectsOf(t, []int{0, 1, -1, 1 << 40})
	checkObjectsOf(t, []any{nil, "a", int32(1), 1, int64(1), 1.0, objectsOfKey{1}, objectsOfKey{2}, (*int)(nil), true})
	checkObjectsOf(t, []objectsOfKey{{1}, {2}, {1}})
	type named string
	checkObjectsOf(t, []named{"a", "b"})
}

func objectsOfStringOps[T any](a, b T) bool {
	return ObjectsEqualsOf(a, b) && ObjectHashCodeOf(a) == ObjectHashCodeOf(b)
}

func TestObjectsEqualsOfDoesNotBoxStrings(t *testing.T) {
	a, b := "some key", string([]byte("some key"))
	if n := testing.AllocsPerRun(100, func() {
		if !objectsOfStringOps(a, b) {
			t.Fatal("not equal")
		}
	}); n != 0 {
		t.Errorf("%v allocations per call", n)
	}
}
