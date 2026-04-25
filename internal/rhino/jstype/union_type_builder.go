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

	// Dedup: check against already-collected alternates. Mirrors Java's
	// Builder.addAlternate (UnionType.java:697) suffix loop.
	for idx := 0; idx < len(b.alternates); idx++ {
		removeCurrent := false
		current := b.alternates[idx]

		// Unknown types: leave in the union and use equality dedup.
		// Mirrors Java's `alternate.isUnknownType() || current.isUnknownType()
		// || alternate.hasAnyTemplateTypes() || current.hasAnyTemplateTypes()` branch.
		if isUnknown(current) || isUnknown(alternate) {
			if alternate.TypeName() == current.TypeName() {
				return b // alternate is unnecessary
			}
		} else if isTemplatizedJSType(alternate) || isTemplatizedJSType(current) {
			// Templatized handling — Java cases 1-9 (UnionType.java:773-832).
			//   Object is considered a subtype of Object<string>, but the
			//   union should keep "Object", not "Object<string>".
			altT, _ := alternate.(*TemplatizedType)
			curT, _ := current.(*TemplatizedType)
			if curT == nil {
				// Cases 1-3: alternate is templatized, current is not.
				if IsStrictSubtypeOf(alternate, current) {
					// Case 1, 2: alternate is unnecessary.
					return b
				}
				// Case 3: leave current, add alternate.
			} else if altT == nil {
				// Cases 4-6: current is templatized, alternate is not.
				if IsStrictSubtypeOf(current, alternate) {
					// Case 4, 5: alternate makes current obsolete.
					removeCurrent = true
				}
				// Case 6: leave current, add alternate.
			} else {
				// Both are templatized. Cases 7-9.
				if templatizedSameRaw(curT, altT) {
					if alternate.TypeName() == current.TypeName() {
						// Case 8: identical templatized — skip.
						return b
					}
					// Case 7a: replace with a `?`-specialised alternate.
					// Java passes ImmutableList.of() to specialise on `?`;
					// the Go side approximates by setting Args to nil.
					alternate = NewTemplatizedType(curT.Base, nil)
					removeCurrent = true
					// Case 7b (ReadonlyArray) — we don't have a special-case
					// readonly type in jsgo yet, so we treat all same-raw
					// templatized merges via the case-7a path.
				}
				// Case 9 / case 7b leave current, add alternate.
			}
		} else if IsSubtype(alternate, current) {
			// Alternate is unnecessary. Java also calls
			// mayRegisterDroppedProperties(alternate, current) here for
			// record types — Go has no central registry for that yet, so
			// the call is omitted (the registry would track which
			// properties got dropped during union construction so the
			// strict-mode "don't drop properties in union" diagnostic can
			// re-attach them. testdontdroppropertiesinunion* fixtures rely
			// on this; closing them needs a registry-side hook).
			return b
		} else if IsSubtype(current, alternate) {
			// Alternate makes current obsolete.
			removeCurrent = true
		}

		if removeCurrent {
			b.alternates = append(b.alternates[:idx], b.alternates[idx+1:]...)
			idx--
		}
	}

	b.alternates = append(b.alternates, alternate)
	return b
}

// isTemplatizedJSType reports whether t is a TemplatizedType. Mirrors
// Java's JSType.isTemplatizedType() — used by Builder.addAlternate for the
// case 1-9 branching.
func isTemplatizedJSType(t JSType) bool {
	_, ok := t.(*TemplatizedType)
	return ok
}

// templatizedSameRaw reports whether two TemplatizedTypes share the same
// Base ObjectType. Mirrors TemplatizedType.wrapsSameRawType
// (TemplatizedType.java:163).
func templatizedSameRaw(a, b *TemplatizedType) bool {
	return a != nil && b != nil && a.Base != nil && a.Base == b.Base
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

