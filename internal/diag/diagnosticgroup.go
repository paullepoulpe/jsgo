// Port of com/google/javascript/jscomp/DiagnosticGroup.java
// @ closure-compiler v20260415

package diag

import "sync"

// DiagnosticGroup groups a set of related DiagnosticTypes so they can be
// toggled as one unit via CompilerOptions. Matches Java's DiagnosticGroup
// final class.
type DiagnosticGroup struct {
	// name is the human-readable name used on the command line
	// (e.g. "checkTypes"). May be empty for anonymous groups.
	name string
	// types is the set of DiagnosticTypes belonging to this group, keyed by Key.
	types map[string]*DiagnosticType
}

// NewDiagnosticGroup creates a named group covering the given DiagnosticTypes.
// Mirrors Java's DiagnosticGroup(String, DiagnosticType...) constructor.
func NewDiagnosticGroup(name string, types ...*DiagnosticType) *DiagnosticGroup {
	g := &DiagnosticGroup{name: name, types: make(map[string]*DiagnosticType, len(types))}
	for _, t := range types {
		if t != nil {
			g.types[t.Key] = t
		}
	}
	return g
}

// NewDiagnosticGroupAnonymous creates a group with no name. Mirrors
// Java's DiagnosticGroup(DiagnosticType...) constructor.
func NewDiagnosticGroupAnonymous(types ...*DiagnosticType) *DiagnosticGroup {
	return NewDiagnosticGroup("", types...)
}

// NewDiagnosticGroupComposite creates a named composite group from sub-groups.
// Mirrors Java's DiagnosticGroup(String, DiagnosticGroup...) constructor.
func NewDiagnosticGroupComposite(name string, groups ...*DiagnosticGroup) *DiagnosticGroup {
	g := &DiagnosticGroup{name: name, types: make(map[string]*DiagnosticType)}
	for _, sub := range groups {
		if sub == nil {
			continue
		}
		for k, t := range sub.types {
			g.types[k] = t
		}
	}
	return g
}

// singletons caches DiagnosticGroup instances created by ForType.
var (
	singletonsMu sync.Mutex
	singletons   = make(map[string]*DiagnosticGroup)
)

// ForType returns (or creates) a singleton DiagnosticGroup that matches only
// the given DiagnosticType. Mirrors Java's DiagnosticGroup.forType.
func ForType(t *DiagnosticType) *DiagnosticGroup {
	singletonsMu.Lock()
	defer singletonsMu.Unlock()
	if g, ok := singletons[t.Key]; ok {
		return g
	}
	g := &DiagnosticGroup{types: map[string]*DiagnosticType{t.Key: t}}
	singletons[t.Key] = g
	return g
}

// MatchesType reports whether the given DiagnosticType's key belongs to this group.
func (g *DiagnosticGroup) MatchesType(t *DiagnosticType) bool {
	if g == nil || t == nil {
		return false
	}
	_, ok := g.types[t.Key]
	return ok
}

// MatchesKey reports whether the given key string belongs to this group.
func (g *DiagnosticGroup) MatchesKey(key string) bool {
	if g == nil {
		return false
	}
	_, ok := g.types[key]
	return ok
}

// GetName returns the human-readable name of the group. May be empty.
func (g *DiagnosticGroup) GetName() string {
	if g == nil {
		return ""
	}
	return g.name
}

// GetTypes returns all DiagnosticTypes in this group.
func (g *DiagnosticGroup) GetTypes() []*DiagnosticType {
	if g == nil {
		return nil
	}
	out := make([]*DiagnosticType, 0, len(g.types))
	for _, t := range g.types {
		out = append(out, t)
	}
	return out
}

func (g *DiagnosticGroup) String() string {
	if g == nil || g.name == "" {
		return "DiagnosticGroup<anonymous>"
	}
	return "DiagnosticGroup<" + g.name + ">"
}
