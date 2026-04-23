// Port of com/google/javascript/rhino/JSDocInfo.java @ closure-compiler v20260415 (Builder subset)
//
// This file ports only the nested `JSDocInfo.Builder` class. In Java it is
// defined inside JSDocInfo.java (not a top-level file); the Go port lifts it
// into a sibling file so jsdoc.go stays focused on the value type.
//
// Scope: the recorder subset whose target fields already exist on
// JSDocInfo (see jsdoc.go). This is a narrow port — the Java Builder carries
// many bits (POLYMER, MIXIN, NG_INJECT, LICENSE, FILEOVERVIEW, ...) that
// jsgo's type checker does not consume yet. New recorders should be added
// here as call sites appear, not preemptively.
//
// Idempotency semantics mirror Java:
//   - record*  methods return true when the value was freshly recorded,
//     false when a conflicting value was already present (duplicates ignored).
//   - Build()  returns the populated *JSDocInfo.
//
// Method names match Java verbatim (exported in Go form), e.g.
// `recordBlockDescription` → `RecordBlockDescription`.

package rhino

// JSDocInfoBuilder mirrors JSDocInfo.Builder (nested static class) in
// closure-compiler. It accumulates annotation state onto a single
// *JSDocInfo and returns it from Build().
type JSDocInfoBuilder struct {
	info      *JSDocInfo
	populated bool // mirrors Builder.populated — set by any recorder that accepted a value

	// blockDescription mirrors prop BLOCK_DESCRIPTION. JSDocInfo in jsgo
	// does not carry this field yet; we store it on the builder only so
	// RecordBlockDescription still matches Java's idempotent-returning
	// contract without requiring a jsdoc.go diff.
	blockDescription    string
	hasBlockDescription bool
}

// NewJSDocInfoBuilder mirrors `JSDocInfo.builder()` in Java, which returns a
// fresh Builder backed by an empty JSDocInfo.
func NewJSDocInfoBuilder() *JSDocInfoBuilder {
	return &JSDocInfoBuilder{info: &JSDocInfo{}}
}

// Build mirrors Builder.build() — returns the populated *JSDocInfo, or nil
// if nothing was recorded (Java returns null in that case). Once called,
// the builder should not be reused; callers that want reuse should call
// BuildAndReset (not ported until a call site needs it).
func (b *JSDocInfoBuilder) Build() *JSDocInfo {
	if b == nil || !b.populated {
		return nil
	}
	return b.info
}

// RecordBlockDescription mirrors Builder.recordBlockDescription(String). The
// Java version stores into the BLOCK_DESCRIPTION prop and returns true once
// populated. jsgo stores on the builder only (jsdoc.go has no matching
// field yet) — the return semantics are preserved so callers behave the
// same.
func (b *JSDocInfoBuilder) RecordBlockDescription(description string) bool {
	b.populated = true
	if b.hasBlockDescription {
		return false
	}
	b.blockDescription = description
	b.hasBlockDescription = true
	return true
}

// BlockDescription exposes the stored BLOCK_DESCRIPTION value for the few
// callers that need to inspect it before Build().
func (b *JSDocInfoBuilder) BlockDescription() string { return b.blockDescription }

// RecordVisibility mirrors Builder.recordVisibility(Visibility). Returns
// false (and does not overwrite) when visibility was already recorded.
func (b *JSDocInfoBuilder) RecordVisibility(v Visibility) bool {
	if !b.info.RecordVisibility(v) {
		return false
	}
	b.populated = true
	return true
}

// OverwriteVisibility mirrors Builder.overwriteVisibility(Visibility). Always
// sets the visibility, overriding any previous value.
func (b *JSDocInfoBuilder) OverwriteVisibility(v Visibility) {
	b.info.visibility = v
	b.info.hasVisibility = true
	b.populated = true
}

// RecordType mirrors Builder.recordType(JSTypeExpression). Returns false
// when `type` is nil or when a @type (or type-related singleton tag) has
// already been recorded.
func (b *JSDocInfoBuilder) RecordType(t *JSTypeExpression) bool {
	if t == nil || b.hasAnyTypeRelatedTags() {
		return false
	}
	b.info.typeExpr = t
	b.populated = true
	return true
}

// RecordReturnType mirrors Builder.recordReturnType(JSTypeExpression).
func (b *JSDocInfoBuilder) RecordReturnType(t *JSTypeExpression) bool {
	if t == nil || b.info.returnExpr != nil || b.hasAnySingletonTypeTags() {
		return false
	}
	b.info.returnExpr = t
	b.populated = true
	return true
}

// RecordParameter mirrors Builder.recordParameter(String, JSTypeExpression).
// Returns false when a parameter of the same name was already recorded or
// when a singleton type tag (@type, @typedef, @enum) is present.
func (b *JSDocInfoBuilder) RecordParameter(name string, t *JSTypeExpression) bool {
	if b.hasAnySingletonTypeTags() {
		return false
	}
	if _, seen := b.info.params[name]; seen {
		return false
	}
	b.info.AddParam(name, t)
	b.populated = true
	return true
}

// RecordTemplateTypeName mirrors Builder.recordTemplateTypeName(String).
func (b *JSDocInfoBuilder) RecordTemplateTypeName(name string) bool {
	for _, existing := range b.info.templateNames {
		if existing == name {
			return false
		}
	}
	if b.info.typedefExpr != nil {
		return false
	}
	b.info.AddTemplateName(name)
	b.populated = true
	return true
}

// RecordTypedef mirrors Builder.recordTypedef(JSTypeExpression).
func (b *JSDocInfoBuilder) RecordTypedef(t *JSTypeExpression) bool {
	if t == nil || b.hasAnyTypeRelatedTags() || len(b.info.templateNames) > 0 {
		return false
	}
	b.info.typedefExpr = t
	b.populated = true
	return true
}

// RecordEnumParameterType mirrors Builder.recordEnumParameterType(JSTypeExpression).
func (b *JSDocInfoBuilder) RecordEnumParameterType(t *JSTypeExpression) bool {
	if t == nil || b.hasAnyTypeRelatedTags() {
		return false
	}
	b.info.SetEnumParameterType(t)
	b.populated = true
	return true
}

// RecordBaseType mirrors Builder.recordBaseType(JSTypeExpression).
func (b *JSDocInfoBuilder) RecordBaseType(t *JSTypeExpression) bool {
	if t == nil || b.hasAnySingletonTypeTags() {
		return false
	}
	b.info.SetBaseType(t)
	b.populated = true
	return true
}

// RecordImplementedInterface mirrors Builder.recordImplementedInterface.
func (b *JSDocInfoBuilder) RecordImplementedInterface(t *JSTypeExpression) bool {
	if t == nil {
		return false
	}
	b.info.AddImplementation(t)
	b.populated = true
	return true
}

// RecordExtendedInterface mirrors Builder.recordExtendedInterface — for
// interface-extends-interface relationships. jsgo stores both constructor
// and interface extends on the same slice.
func (b *JSDocInfoBuilder) RecordExtendedInterface(t *JSTypeExpression) bool {
	if t == nil {
		return false
	}
	b.info.SetBaseType(t)
	b.populated = true
	return true
}

// RecordConstructor mirrors Builder.recordConstructor().
func (b *JSDocInfoBuilder) RecordConstructor() bool {
	if b.hasAnySingletonTypeTags() || b.isConstructorOrInterface() {
		return false
	}
	b.info.isConstructor = true
	b.populated = true
	return true
}

// RecordInterface mirrors Builder.recordInterface().
func (b *JSDocInfoBuilder) RecordInterface() bool {
	if b.hasAnySingletonTypeTags() || b.info.isConstructor || b.info.isAbstract {
		return false
	}
	b.info.isInterface = true
	b.populated = true
	return true
}

// RecordImplicitMatch mirrors Builder.recordImplicitMatch() — the @record tag.
// Sets both the RECORD and INTERFACE bits.
func (b *JSDocInfoBuilder) RecordImplicitMatch() bool {
	if b.hasAnySingletonTypeTags() || b.isConstructorOrInterface() {
		return false
	}
	b.info.isRecord = true
	b.info.isInterface = true
	b.populated = true
	return true
}

// RecordAbstract mirrors Builder.recordAbstract().
func (b *JSDocInfoBuilder) RecordAbstract() bool {
	if b.hasAnySingletonTypeTags() || b.info.isInterface || b.info.isFinal ||
		(b.info.hasVisibility && b.info.visibility == VisibilityPrivate) {
		return false
	}
	b.info.isAbstract = true
	b.populated = true
	return true
}

// RecordStruct mirrors Builder.recordStruct().
func (b *JSDocInfoBuilder) RecordStruct() bool {
	if b.hasAnySingletonTypeTags() || b.info.isDict {
		return false
	}
	b.info.isStruct = true
	b.populated = true
	return true
}

// RecordDict mirrors Builder.recordDict().
func (b *JSDocInfoBuilder) RecordDict() bool {
	if b.hasAnySingletonTypeTags() || b.info.isStruct {
		return false
	}
	b.info.isDict = true
	b.populated = true
	return true
}

// RecordOverride mirrors Builder.recordOverride().
func (b *JSDocInfoBuilder) RecordOverride() bool {
	if b.info.isOverride {
		return false
	}
	b.info.isOverride = true
	b.populated = true
	return true
}

// RecordDeprecated mirrors Builder.recordDeprecated().
func (b *JSDocInfoBuilder) RecordDeprecated() bool {
	if b.info.isDeprecated {
		return false
	}
	b.info.isDeprecated = true
	b.populated = true
	return true
}

// RecordDeprecationReason mirrors Builder.recordDeprecationReason(String).
func (b *JSDocInfoBuilder) RecordDeprecationReason(reason string) bool {
	if b.info.deprecationReason != "" {
		return false
	}
	b.info.deprecationReason = reason
	b.populated = true
	return true
}

// RecordConstancy mirrors Builder.recordConstancy() — the @const flag.
func (b *JSDocInfoBuilder) RecordConstancy() bool {
	if b.info.isConst {
		return false
	}
	b.info.isConst = true
	b.populated = true
	return true
}

// RecordFinality mirrors Builder.recordFinality() — the @final flag.
func (b *JSDocInfoBuilder) RecordFinality() bool {
	if b.info.isFinal {
		return false
	}
	b.info.isFinal = true
	b.populated = true
	return true
}

// RecordNoSideEffects mirrors Builder.recordNoSideEffects().
func (b *JSDocInfoBuilder) RecordNoSideEffects() bool {
	if b.info.isNoSideEffects {
		return false
	}
	b.info.isNoSideEffects = true
	b.populated = true
	return true
}

// RecordThisType mirrors Builder.recordThisType(JSTypeExpression). jsgo
// currently tracks only the presence of @this (not the resolved type), so
// we accept any non-nil expression and flip hasThisType.
func (b *JSDocInfoBuilder) RecordThisType(t *JSTypeExpression) bool {
	if t == nil || b.info.hasThisType || b.hasAnySingletonTypeTags() {
		return false
	}
	b.info.hasThisType = true
	b.populated = true
	return true
}

// RecordLends mirrors Builder.recordLends(JSTypeExpression).
func (b *JSDocInfoBuilder) RecordLends(name *JSTypeExpression) bool {
	if name == nil || b.hasAnyTypeRelatedTags() || b.info.lendsName != "" {
		return false
	}
	b.info.lendsName = name.TypeName
	b.populated = true
	return true
}

// RecordSuppression mirrors Builder.recordSuppression(String).
func (b *JSDocInfoBuilder) RecordSuppression(name string) {
	b.info.AddSuppression(name)
	b.populated = true
}

// HasParameter mirrors Builder.hasParameter(String).
func (b *JSDocInfoBuilder) HasParameter(name string) bool {
	if b == nil || b.info == nil {
		return false
	}
	_, ok := b.info.params[name]
	return ok
}

// IsConstructorRecorded mirrors Builder.isConstructorRecorded().
func (b *JSDocInfoBuilder) IsConstructorRecorded() bool {
	return b != nil && b.info != nil && b.info.isConstructor
}

// IsInterfaceRecorded mirrors Builder.isInterfaceRecorded().
func (b *JSDocInfoBuilder) IsInterfaceRecorded() bool {
	return b != nil && b.info != nil && b.info.isInterface
}

// IsStructRecorded mirrors Builder.isStructRecorded().
func (b *JSDocInfoBuilder) IsStructRecorded() bool {
	return b != nil && b.info != nil && b.info.isStruct
}

// IsDictRecorded mirrors Builder.isDictRecorded().
func (b *JSDocInfoBuilder) IsDictRecorded() bool {
	return b != nil && b.info != nil && b.info.isDict
}

// --- Compatibility gates mirroring Java private helpers ----------------

// hasAnySingletonTypeTags mirrors Builder.hasAnySingletonTypeTags() — the
// set of tags that cannot coexist with @param/@return/@this/@extends.
// Java includes @type, @typedef, and @enum here.
func (b *JSDocInfoBuilder) hasAnySingletonTypeTags() bool {
	return b.info.typeExpr != nil ||
		b.info.typedefExpr != nil ||
		b.info.enumMarked
}

// hasAnyTypeRelatedTags mirrors Builder.hasAnyTypeRelatedTags() — stricter
// than hasAnySingletonTypeTags; also forbids @param/@return and class-shape
// tags alongside @type.
func (b *JSDocInfoBuilder) hasAnyTypeRelatedTags() bool {
	return b.hasAnySingletonTypeTags() ||
		b.info.returnExpr != nil ||
		len(b.info.params) > 0 ||
		b.info.isConstructor ||
		b.info.isInterface ||
		b.info.extendsExpr != nil ||
		len(b.info.implementsAll) > 0 ||
		b.info.hasThisType
}

// isConstructorOrInterface mirrors Builder.isConstructorOrInterface().
func (b *JSDocInfoBuilder) isConstructorOrInterface() bool {
	return b.info.isConstructor || b.info.isInterface
}
