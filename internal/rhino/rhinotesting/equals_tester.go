package rhinotesting

// Go-only: language-adapter — Guava's com.google.common.testing.EqualsTester,
// used by the ported tests.

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// EqualsTester checks equals/hashCode: items in the same group are equal
// with equal hash codes, items in different groups are not equal, and no
// item equals null or an instance of an unrelated class.
type EqualsTester struct {
	groups [][]any
}

// NewEqualsTester returns an empty tester.
func NewEqualsTester() *EqualsTester {
	return &EqualsTester{}
}

// AddEqualityGroup adds a group of mutually equal items.
func (e *EqualsTester) AddEqualityGroup(items ...any) *EqualsTester {
	e.groups = append(e.groups, items)
	return e
}

type equalsTesterNotAnInstance struct{}

// TestEquals runs the checks.
func (e *EqualsTester) TestEquals(t testing.TB) *EqualsTester {
	t.Helper()
	for gi, group := range e.groups {
		for _, item := range group {
			if jutil.ObjectsEquals(item, nil) {
				t.Fatalf("%s must not be Object#equals to null", Show(item))
			}
			if jutil.ObjectsEquals(item, &equalsTesterNotAnInstance{}) {
				t.Fatalf("%s must not be Object#equals to an arbitrary object of another class", Show(item))
			}
			for gj, other := range e.groups {
				for _, related := range other {
					eq := jutil.ObjectsEquals(item, related)
					if gi == gj {
						if !eq {
							t.Fatalf("%s [group %d] must be Object#equals to %s [group %d]", Show(item), gi+1, Show(related), gj+1)
						}
						if jutil.ObjectHashCode(item) != jutil.ObjectHashCode(related) {
							t.Fatalf("the Object#hashCode of %s must be equal to the Object#hashCode of %s", Show(item), Show(related))
						}
					} else if eq {
						t.Fatalf("%s [group %d] must not be Object#equals to %s [group %d]", Show(item), gi+1, Show(related), gj+1)
					}
				}
			}
		}
	}
	return e
}
