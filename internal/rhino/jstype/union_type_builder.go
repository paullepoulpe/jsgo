// Port of UnionType.Builder (inner class) from
// com/google/javascript/rhino/jstype/UnionType.java @ closure-compiler v20260415
//
// UnionTypeBuilder assembles a new UnionType from individual alternates,
// enforcing the invariants described in Java's Builder:
//   - nested unions are flattened (union-of-union collapse)
//   - duplicates and subtypes are deduplicated
//   - if the union contains AllType or an UnknownType the result collapses
//     to that wildcard (plus Void when containsVoidType)
//   - empty builds return NoType (the bottom)
//   - unions larger than maxUnionSize collapse to Unknown

package jstype

// maxUnionSize mirrors UnionType.MAX_UNION_SIZE in Java.
const maxUnionSize = 30

// UnionTypeBuilder is the state machine that constructs a UnionType.
// Port of the static final class UnionType.Builder.
// Create via NewUnionTypeBuilder; call AddAlternate for each member;
// call Build to obtain the final JSType.
type UnionTypeBuilder struct {
	alternates       []JSType
	containsVoidType bool
	isAllType        bool
	isNativeUnknown  bool
}

// NewUnionTypeBuilder returns a fresh builder. Mirrors Builder(JSTypeRegistry).
func NewUnionTypeBuilder() *UnionTypeBuilder {
	return &UnionTypeBuilder{}
}

// AddAlternates adds all types from the slice. Convenience wrapper.
func (b *UnionTypeBuilder) AddAlternates(types []JSType) *UnionTypeBuilder {
	for _, t := range types {
		b.AddAlternate(t)
	}
	return b
}

// AddAlternate adds a single type to the union under construction.
// Mirrors Builder.addAlternate(JSType) — flattens nested unions,
// tracks wildcard flags, deduplicates using subtype relationships.
func (b *UnionTypeBuilder) AddAlternate(alternate JSType) *UnionTypeBuilder {
	if alternate == nil {
		return b
	}

	// Flatten union-of-union. Cycle detection: we won't recurse infinitely
	// because AddAlternate is never called with the builder's own output.
	if u, ok := alternate.(*UnionType); ok {
		for _, m := range u.members {
			b.AddAlternate(m)
		}
		return b
	}

	if len(b.alternates) > maxUnionSize {
		return b
	}

	// bottom type contributes nothing
	if isNoType(alternate) {
		return b
	}

	// Track wildcard flags.
	if isAll(alternate) {
		b.isAllType = true
		return b
	}
	if isVoidType(alternate) {
		b.containsVoidType = true
		// still add void as an explicit member below (Java keeps void even with wildcards)
	}
	if _, ok := alternate.(unknownType); ok {
		b.isNativeUnknown = true
		return b
	}
	if b.isAllType || b.isNativeUnknown {
		// collect only void when a wildcard dominates
		if !isVoidType(alternate) {
			return b
		}
	}

	// Dedup: check against already-collected alternates.
	for idx := 0; idx < len(b.alternates); idx++ {
		current := b.alternates[idx]
		if isUnknown(current) {
			if alternate.TypeName() == current.TypeName() {
				return b // same checked-unknown name, skip
			}
			continue
		}
		if isUnknown(alternate) {
			if alternate.TypeName() == current.TypeName() {
				return b
			}
			continue
		}
		if IsSubtype(alternate, current) {
			// alternate is redundant (current is already wider)
			return b
		}
		if IsSubtype(current, alternate) {
			// current is made redundant by the wider alternate — remove it
			b.alternates = append(b.alternates[:idx], b.alternates[idx+1:]...)
			idx--
		}
	}

	b.alternates = append(b.alternates, alternate)
	return b
}

// Build finalises the union and returns the collapsed JSType.
// Mirrors Builder.build() in Java. The builder must not be reused.
func (b *UnionTypeBuilder) Build() JSType {
	// Wildcard dominates.
	if b.isAllType {
		if b.containsVoidType {
			return NewUnion(All, Void)
		}
		return All
	}
	if b.isNativeUnknown {
		// In the Go port there is no separate checkedUnknownType; both
		// checked and unchecked unknowns are represented as Unknown.
		if b.containsVoidType {
			return NewUnion(Unknown, Void)
		}
		return Unknown
	}

	switch len(b.alternates) {
	case 0:
		return None
	case 1:
		return b.alternates[0]
	default:
		if len(b.alternates) > maxUnionSize {
			return Unknown
		}
		// Use NewUnion to get deterministic ordering and dedup.
		return NewUnion(b.alternates...)
	}
}

