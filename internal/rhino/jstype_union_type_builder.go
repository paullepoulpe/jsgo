// Port of com/google/javascript/rhino/jstype/UnionType.java @ closure-compiler v20261021

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Bob Jervis
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino

import (
	"iter"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// UnionType_Builder implements type unioning logic, since UnionTypes only
// actually need to perform unioning operations when being (re)built.
//
// Builders exist in two forms. One for assembing a new union and one for
// updating an existing union. Only the former is exposed.
//
// Most users of this class should prefer JSTypeRegistry#createUnionType
// instead.
type UnionType_Builder struct {
	rebuildTarget *UnionType
	registry      *JSTypeRegistry

	inlineAlternates [32]JSType
	inlineLen        int
	alternates       []JSType
	finalAlternates  []JSType

	// If a union has ? or *, we do not care about any other types, except for
	// undefined (for optional properties).
	containsVoidType              bool
	isAllType                     bool
	isNativeUnknownType           bool
	areAllUnknownsChecked         bool
	alternatesResolvedBeforeBuild bool
	hasBuilt                      bool

	// Every UnionType may have at most one structural function in it.
	//
	// NOTE(nicksantos): I've read some literature that says that
	// type-inferenced languages are fundamentally incompatible with union
	// types. I refuse to believe this. But they do make the type lattice much
	// more complicated.
	//
	// For this reason, when we deal with function types, we actually merge
	// some nodes on the lattice, and treat them as fundamentally equivalent.
	// For example, we treat
	//
	//	function(): string | function(): number
	//
	// as equivalent to
	//
	//	function(): (string|number)
	//
	// and normalize the first type into the second type.
	//
	// To perform this normalization, we've modified Builder to disallow
	// multiple structural functions in a union. We always delegate to
	// FunctionType::getLeastSupertype, which either merges the functions into
	// one structural function, or just bails out and uses the top function
	// type.
	functionTypePosition int
}

func (b *UnionType_Builder) size() int {
	if b.alternates != nil {
		return len(b.alternates)
	}
	return b.inlineLen
}

func (b *UnionType_Builder) get(i int) JSType {
	if b.alternates != nil {
		return b.alternates[i]
	}
	return b.inlineAlternates[i]
}

func (b *UnionType_Builder) set(i int, t JSType) {
	if b.alternates != nil {
		b.alternates[i] = t
	} else {
		b.inlineAlternates[i] = t
	}
}

func (b *UnionType_Builder) removeAt(i int) {
	if b.alternates == nil {
		last := b.inlineLen - 1
		copy(b.inlineAlternates[i:last], b.inlineAlternates[i+1:b.inlineLen])
		b.inlineAlternates[last] = nil
		b.inlineLen = last
		return
	}
	last := len(b.alternates) - 1
	copy(b.alternates[i:], b.alternates[i+1:])
	b.alternates[last] = nil
	b.alternates = b.alternates[:last]
}

func (b *UnionType_Builder) add(t JSType) {
	if b.alternates != nil {
		b.alternates = append(b.alternates, t)
		return
	}
	if b.inlineLen < len(b.inlineAlternates) {
		b.inlineAlternates[b.inlineLen] = t
		b.inlineLen++
		return
	}
	b.alternates = make([]JSType, len(b.inlineAlternates)+1, 64)
	copy(b.alternates, b.inlineAlternates[:])
	b.alternates[len(b.inlineAlternates)] = t
	clear(b.inlineAlternates[:])
	b.inlineLen = 0
}

func initUnionType_Builder(b *UnionType_Builder, registry *JSTypeRegistry) {
	b.rebuildTarget = nil
	b.registry = registry
	b.inlineLen = 0
	b.alternates = nil
	b.finalAlternates = nil
	b.containsVoidType = false
	b.isAllType = false
	b.isNativeUnknownType = false
	b.areAllUnknownsChecked = true
	b.alternatesResolvedBeforeBuild = true
	b.hasBuilt = false
	b.functionTypePosition = -1
}

// NewUnionType_Builder creates a builder for a new union.
//
// Java: UnionType.Builder(JSTypeRegistry)
func NewUnionType_Builder(registry *JSTypeRegistry) *UnionType_Builder {
	return &UnionType_Builder{
		registry:                      registry,
		areAllUnknownsChecked:         true,
		alternatesResolvedBeforeBuild: true,
		functionTypePosition:          -1,
	}
}

// NewUnionType_BuilderOfRebuildTarget creates a re-builder for an existing
// union.
//
// Java: UnionType.Builder(UnionType)
func NewUnionType_BuilderOfRebuildTarget(rebuildTarget *UnionType) *UnionType_Builder {
	b := NewUnionType_Builder(rebuildTarget.GetRegistry())
	b.rebuildTarget = rebuildTarget
	return b
}

// Java: UnionType.Builder.isSubtype(JSType, JSType)
func UnionType_Builder_IsSubtype(rightType JSType, leftType JSType) bool {
	return rightType.IsSubtypeWithoutStructuralTyping(leftType)
}

// Java: UnionType.Builder.addAlternates(Collection<? extends JSType>)
func (b *UnionType_Builder) AddAlternatesOfCollection(c iter.Seq[JSType]) *UnionType_Builder {
	for type_ := range c {
		b.AddAlternate(type_)
	}
	return b
}

// AddAlternates: a specific override that avoid creating an iterator. This
// version is currently used when adding a union as an alternate.
//
// Java: UnionType.Builder.addAlternates(List<? extends JSType>)
func (b *UnionType_Builder) AddAlternates(list []JSType) *UnionType_Builder {
	for i := 0; i < len(list); i++ {
		b.AddAlternate(list[i])
	}
	return b
}

// AddAlternate adds an alternate to the union type under construction.
//
// Returns this for easy chaining.
//
// Java: UnionType.Builder.addAlternate(JSType)
func (b *UnionType_Builder) AddAlternate(alternate JSType) *UnionType_Builder {
	b.CheckHasNotBuilt()

	altBase := jsTypeBase(alternate)
	altKind := altBase.objFlags & objKindMask
	if jsTypeIsUnion(alternate, altKind) {
		b.AddAlternates(alternate.ToMaybeUnionType().GetAlternates())
		return b
	}

	if b.size() > UnionType_MAX_UNION_SIZE {
		return b
	}

	// Defer removing duplicate elements until all alternates in the union are
	// resolved. JSType operations like equality and subtyping are not reliable
	// pre-resolution.
	if altBase.resolveResult == nil {
		b.alternatesResolvedBeforeBuild = false

		for i := 0; i < b.size(); i++ {
			if jsTypeBase(b.get(i)) == altBase {
				return b
			}
		}
		b.add(alternate)
		return b
	}

	// build() returns the bottom type by default, so we can just bail out
	// early here.
	if altKind == objKindNo || (((uint32(1)<<altKind)&jsKindProxyOrFallbackMask) != 0 && alternate.IsNoType()) {
		return b
	}

	if !b.isAllType && (altKind == jsKindAll || (((uint32(1)<<altKind)&jsKindProxyOrFallbackMask) != 0 && alternate.IsAllType())) {
		b.isAllType = true
	}
	if !b.containsVoidType && (altKind == jsKindVoid || (altKind == objKindFallback && alternate.IsVoidType())) {
		b.containsVoidType = true
	}

	if altKind == objKindUnknown {
		b.isNativeUnknownType = true
		b.areAllUnknownsChecked = b.areAllUnknownsChecked && alternate.(*UnknownType).isChecked
	}

	if b.isAllType || b.isNativeUnknownType {
		return b
	}

	// Function types are special, because they have their own bizarre
	// sub-lattice. See the comments on FunctionType#supAndInf helper and above
	// at functionTypePosition.
	altIsFunc := jsTypeIsFunction(alternate, altKind)
	if altIsFunc && b.functionTypePosition != -1 {
		// See the comments on functionTypePosition above.
		other := b.get(b.functionTypePosition).ToMaybeFunctionType()
		supremum := alternate.ToMaybeFunctionType().SupAndInfHelper(other, true)
		b.set(b.functionTypePosition, supremum)
		return b
	}

	altUnknown := jsTypeIsUnknown(alternate, altBase, altKind)
	altHasTemplate := !altUnknown && jsTypeHasAnyTemplateTypes(altBase)
	altTemplatized := !altUnknown && !altHasTemplate && jsTypeIsTemplatized(alternate, altKind)

	// Look through the alternates we've got so far, and check if any of them
	// are duplicates of one another.
	for index := 0; index < b.size(); index++ {
		removeCurrent := false
		current := b.get(index)
		curBase := jsTypeBase(current)
		if curBase.resolveResult == nil { // Defer de-duplicating unresolved alternates.
			continue
		}
		if curBase == altBase {
			return b
		}
		curKind := curBase.objFlags & objKindMask

		// Unknown types may just be names that haven't been resolved yet. So
		// keep these in the union, and just use equality checking for simple
		// de-duping.
		if altUnknown ||
			jsTypeIsUnknown(current, curBase, curKind) ||
			altHasTemplate ||
			jsTypeHasAnyTemplateTypes(curBase) {
			if equalityChecker_AreIdentityEqualWithBases(alternate, altBase, current, curBase) {
				// Alternate is unnecessary.
				return b
			}
		} else if curTemplatized := jsTypeIsTemplatized(current, curKind); altTemplatized || curTemplatized {
			// Because "Foo" and "Foo<?>" are roughly equivalent templatized types,
			// special care is needed when building the union. For example:
			//   Object is consider a subtype of Object<string>
			// but we want to leave "Object" not "Object<string>" when building the
			// subtype.
			//

			// Cases:
			// 1) alternate:Array<string> and current:Object ==> Object
			// 2) alternate:Array<string> and current:Array ==> Array
			// 3) alternate:Object<string> and
			//    current:Array ==> Array|Object<string>
			// 4) alternate:Object and current:Array<string> ==> Object
			// 5) alternate:Array and current:Array<string> ==> Array
			// 6) alternate:Array and
			//    current:Object<string> ==> Array|Object<string>
			// 7a) alternate:Object<string> and
			//     current:Object<number> ==> Object<?>
			// 7b) alternate:ReadonlyArray<string> and
			//     current:ReadonlyArray<number> ==> ReadonlyArray<string>|ReadonlyArray<number>
			// 8) alternate:Array<string> and
			//    current:Array<string> ==> Array<string>
			// 9) alternate:Array<string> and
			//    current:Object<string> ==> Object<string>|Array<string>

			if !curTemplatized {
				if UnionType_Builder_IsSubtype(alternate, current) {
					// case 1, 2
					return b
				}
				// case 3: leave current, add alternate
			} else if !altTemplatized {
				if UnionType_Builder_IsSubtype(current, alternate) {
					// case 4, 5
					removeCurrent = true
				}
				// case 6: leave current, add alternate
			} else {
				templatizedAlternate := alternate.ToMaybeTemplatizedType()
				templatizedCurrent := current.ToMaybeTemplatizedType()

				if templatizedCurrent.WrapsSameRawType(templatizedAlternate) {

					if equalityChecker_AreIdentityEqualWithBases(current, curBase, alternate, altBase) {
						// case 8
						return b
					} else if !templatizedCurrent.GetReferencedType().IsReadonlyArrayType() {
						// case 7a: replace with a merged alternate specialized on `?`.
						rawType := templatizedCurrent.GetReferencedObjTypeInternal()
						// Providing no type-parameter values specializes `rawType` on `?`
						// by default.
						alternate = b.registry.CreateTemplatizedTypeOfList(rawType, nil)
						altBase = jsTypeBase(alternate)
						altKind = altBase.objFlags & objKindMask
						altIsFunc = jsTypeIsFunction(alternate, altKind)
						altUnknown = jsTypeIsUnknown(alternate, altBase, altKind)
						altHasTemplate = !altUnknown && jsTypeHasAnyTemplateTypes(altBase)
						altTemplatized = !altUnknown && !altHasTemplate && jsTypeIsTemplatized(alternate, altKind)
						removeCurrent = true
					}
				}
				// case 9 & 7b: leave current, add alternate
			}
			// Otherwise leave both templatized types.
		} else {
			bothBits := (uint32(1) << altKind) | (uint32(1) << curKind)
			if (bothBits&jsKindSimpleDisjointValueMask) == 0 ||
				(bothBits&jsKindCanRelateToValueMask) != 0 {
				if UnionType_Builder_IsSubtype(alternate, current) {
					// Alternate is unnecessary.
					b.MayRegisterDroppedProperties(alternate, current)
					return b
				} else if UnionType_Builder_IsSubtype(current, alternate) {
					// Alternate makes current obsolete
					b.MayRegisterDroppedProperties(current, alternate)
					removeCurrent = true
				}
			}
		}

		if removeCurrent {
			b.removeAt(index)

			if index == b.functionTypePosition {
				b.functionTypePosition = -1
			} else if index < b.functionTypePosition {
				b.functionTypePosition--
			}

			index--
		}
	}

	if altIsFunc {
		// See the comments on functionTypePosition above.
		jutil.CheckState(b.functionTypePosition == -1)
		b.functionTypePosition = b.size()
	}

	b.add(alternate)
	return b
}

// Java: UnionType.Builder.mayRegisterDroppedProperties(JSType, JSType)
func (b *UnionType_Builder) MayRegisterDroppedProperties(subtype JSType, supertype JSType) {
	subKind := jsTypeBase(subtype).objFlags & objKindMask
	supKind := jsTypeBase(supertype).objFlags & objKindMask
	if ((uint32(1)<<subKind)&jsKindMaybeRecordMask) == 0 ||
		((uint32(1)<<supKind)&jsKindMaybeRecordMask) == 0 {
		return
	}
	if subRec := subtype.ToMaybeRecordType(); subRec != nil {
		if supRec := supertype.ToMaybeRecordType(); supRec != nil && subRec != supRec {
			b.registry.RegisterDroppedPropertiesInUnion(subRec, supRec)
		}
	}
}

// Build returns a type, not necessarily a UnionType, that represents the
// union of the inputs.
//
// The Builder cannot be used again once this method is called.
//
// Java: UnionType.Builder.build()
func (b *UnionType_Builder) Build() JSType {
	if b.rebuildTarget == nil {
		b.CheckHasNotBuilt()
		b.hasBuilt = true

		wildcard := b.GetNativeWildcardType()
		if wildcard != nil {
			if !b.containsVoidType {
				return wildcard
			}
			b.inlineAlternates[0] = wildcard
			b.inlineAlternates[1] = b.registry.GetNativeType(JSTypeNative_VOID_TYPE)
			b.inlineLen = 2
			b.alternates = nil
			b.finalAlternates = nil
			return NewUnionType(b)
		}
		n := b.size()
		if n == 0 {
			return b.registry.GetNativeType(JSTypeNative_NO_TYPE)
		}
		if n == 1 {
			return b.get(0)
		}
		if n > UnionType_MAX_UNION_SIZE {
			return b.registry.GetNativeType(JSTypeNative_UNKNOWN_TYPE)
		}
		if b.alternates != nil {
			b.finalAlternates = b.alternates
		}
		return NewUnionType(b)
	}

	b.BuildInternal()
	b.rebuildTarget.FillFromBuilder(b)
	if b.rebuildTarget.heapAlternates == nil && b.rebuildTarget.inlineLen == 1 {
		return b.rebuildTarget.inlineAlternates[0]
	}
	return b.rebuildTarget
}

// BuildInternal creates the final set of alternates for either a new union
// or a union being rebuilt.
//
// Java: UnionType.Builder.buildInternal()
func (b *UnionType_Builder) BuildInternal() *UnionType_Builder {
	b.CheckHasNotBuilt()
	b.hasBuilt = true

	wildcard := b.GetNativeWildcardType()
	if wildcard != nil {
		if b.containsVoidType {
			b.inlineAlternates[0] = wildcard
			b.inlineAlternates[1] = b.registry.GetNativeType(JSTypeNative_VOID_TYPE)
			b.inlineLen = 2
		} else {
			b.inlineAlternates[0] = wildcard
			b.inlineLen = 1
		}
		b.alternates = nil
		b.finalAlternates = nil
	} else if b.size() == 0 {
		// To simplify the typesystem, empty union types are forbidden. Using a
		// single `bottom` makes it essentially a proxy instead.
		b.inlineAlternates[0] = b.registry.GetNativeType(JSTypeNative_NO_TYPE)
		b.inlineLen = 1
		b.alternates = nil
		b.finalAlternates = nil
	} else if b.size() > UnionType_MAX_UNION_SIZE {
		b.inlineAlternates[0] = b.registry.GetNativeType(JSTypeNative_UNKNOWN_TYPE)
		b.inlineLen = 1
		b.alternates = nil
		b.finalAlternates = nil
	} else if b.alternates != nil {
		b.finalAlternates = b.alternates
	}

	return b
}

// GetNativeWildcardType returns ALL_TYPE, UNKNOWN_TYPE,
// CHECKED_UNKNOWN_TYPE, or null as specified by the flags.
//
// Java: UnionType.Builder.getNativeWildcardType()
func (b *UnionType_Builder) GetNativeWildcardType() JSType {
	if b.isAllType {
		return b.registry.GetNativeType(JSTypeNative_ALL_TYPE)
	} else if b.isNativeUnknownType {
		if b.areAllUnknownsChecked {
			return b.registry.GetNativeType(JSTypeNative_CHECKED_UNKNOWN_TYPE)
		} else {
			return b.registry.GetNativeType(JSTypeNative_UNKNOWN_TYPE)
		}
	}
	return nil
}

// Java: UnionType.Builder.checkHasNotBuilt()
func (b *UnionType_Builder) CheckHasNotBuilt() {
	if b.hasBuilt {
		panic(jutil.NewIllegalStateExceptionMsg("Cannot reuse a `UnionType.Builder` that has already built."))
	}
}
