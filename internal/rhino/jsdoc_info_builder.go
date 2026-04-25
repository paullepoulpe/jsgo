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

// RecordThrowsAnnotation mirrors Builder.recordThrowsAnnotation(String). Returns
// true if the annotation was recorded. Drops the annotation when a singleton
// type tag is already present (same guard as Java). If the annotation text is
// empty the call still records an empty-string sentinel so PureFunctionIdentifier
// can detect the @throws presence.
func (b *JSDocInfoBuilder) RecordThrowsAnnotation(annotation string) bool {
	b.populated = true
	if b.hasAnySingletonTypeTags() {
		return false
	}
	b.info.AddThrowsAnnotation(annotation)
	return true
}

// RecordClosurePrimitiveId mirrors Builder.recordClosurePrimitiveId(String).
func (b *JSDocInfoBuilder) RecordClosurePrimitiveId(id string) bool {
	if !b.info.SetClosurePrimitiveId(id) {
		return false
	}
	b.populated = true
	return true
}

// RecordModifies mirrors Builder.recordModifies(Set<String>). Returns false
// when any singleton side-effect tag is already present, otherwise records
// the modifies set.
func (b *JSDocInfoBuilder) RecordModifies(targets map[string]bool) bool {
	if b.hasAnySingletonSideEffectTags() {
		return false
	}
	if !b.info.SetModifies(targets) {
		return false
	}
	b.populated = true
	return true
}

// RecordSuppressions merges a set of suppressions into the existing set,
// mirroring Builder.recordSuppressions(Set<String>). Unlike RecordSuppression
// (single name), this accepts multiple names at once.
func (b *JSDocInfoBuilder) RecordSuppressions(names map[string]bool) {
	for n := range names {
		b.info.AddSuppression(n)
	}
	b.populated = true
}

// RecordTypeTransformation mirrors Builder.recordTypeTransformation(String, Node).
// In jsgo the AST node is not available yet, so we store the TTL template name
// and return true to signal acceptance. Returns false when a regular @template
// of the same name already exists.
func (b *JSDocInfoBuilder) RecordTypeTransformation(name string) bool {
	for _, existing := range b.info.templateNames {
		if existing == name {
			return false
		}
	}
	b.info.AddTTLTemplateName(name)
	b.populated = true
	return true
}

// RecordNgInject mirrors Builder.recordNgInject(boolean).
func (b *JSDocInfoBuilder) RecordNgInject(val bool) bool {
	if b.info.CheckBit(BitNgInject) {
		return false
	}
	b.info.SetBit(BitNgInject, val)
	b.populated = true
	return true
}

// RecordWizaction mirrors Builder.recordWizaction().
func (b *JSDocInfoBuilder) RecordWizaction() bool {
	if b.info.CheckBit(BitWizAction) {
		return false
	}
	b.info.SetBit(BitWizAction, true)
	b.populated = true
	return true
}

// RecordPolymerBehavior mirrors Builder.recordPolymerBehavior().
func (b *JSDocInfoBuilder) RecordPolymerBehavior() bool {
	if b.info.CheckBit(BitPolymerBehavior) {
		return false
	}
	b.info.SetBit(BitPolymerBehavior, true)
	b.populated = true
	return true
}

// RecordCustomElement mirrors Builder.recordCustomElement().
func (b *JSDocInfoBuilder) RecordCustomElement() bool {
	if b.info.CheckBit(BitCustomElement) {
		return false
	}
	b.info.SetBit(BitCustomElement, true)
	b.populated = true
	return true
}

// RecordIdGenerator mirrors Builder.recordIdGenerator(). Plain @idGenerator.
func (b *JSDocInfoBuilder) RecordIdGenerator() bool {
	if !b.info.SetIdGenerator(IdGeneratorUnique) {
		return false
	}
	b.populated = true
	return true
}

// RecordConsistentIdGenerator mirrors Builder.recordConsistentIdGenerator().
func (b *JSDocInfoBuilder) RecordConsistentIdGenerator() bool {
	if !b.info.SetIdGenerator(IdGeneratorConsistent) {
		return false
	}
	b.populated = true
	return true
}

// RecordXidGenerator mirrors Builder.recordXidGenerator().
func (b *JSDocInfoBuilder) RecordXidGenerator() bool {
	if !b.info.SetIdGenerator(IdGeneratorXid) {
		return false
	}
	b.populated = true
	return true
}

// RecordOriginalVisibility mirrors the originalVisibility setter pattern.
// Returns false if already recorded.
func (b *JSDocInfoBuilder) RecordOriginalVisibility(v Visibility) bool {
	if !b.info.SetOriginalVisibility(v) {
		return false
	}
	b.populated = true
	return true
}

// RecordOriginalCommentString mirrors Builder.recordOriginalCommentString(String).
func (b *JSDocInfoBuilder) RecordOriginalCommentString(comment string) {
	if b.info.IsDocumentationIncluded() {
		b.info.SetOriginalCommentString(comment)
		b.populated = true
	}
}

// RecordOriginalCommentPosition mirrors Builder.recordOriginalCommentPosition(int).
func (b *JSDocInfoBuilder) RecordOriginalCommentPosition(pos int) {
	if b.info.IsDocumentationIncluded() {
		b.info.SetOriginalCommentPosition(pos)
		b.populated = true
	}
}

// RecordDescription mirrors Builder.recordDescription(String).
func (b *JSDocInfoBuilder) RecordDescription(desc string) bool {
	if b.info.description != "" {
		return false
	}
	b.info.description = desc
	b.populated = true
	return true
}

// RecordMeaning mirrors Builder.recordMeaning(String).
func (b *JSDocInfoBuilder) RecordMeaning(meaning string) bool {
	if b.info.meaning != "" {
		return false
	}
	b.info.meaning = meaning
	b.populated = true
	return true
}

// RecordAlternateMessageId mirrors Builder.recordAlternateMessageId(String).
func (b *JSDocInfoBuilder) RecordAlternateMessageId(id string) bool {
	if b.info.alternateMessageId != "" {
		return false
	}
	b.info.alternateMessageId = id
	b.populated = true
	return true
}

// RecordAuthor mirrors Builder.recordAuthor(String).
func (b *JSDocInfoBuilder) RecordAuthor(author string) bool {
	if b.info.IsDocumentationIncluded() {
		b.info.AddAuthor(author)
	}
	b.populated = true
	return true
}

// RecordReference mirrors Builder.recordReference(String).
func (b *JSDocInfoBuilder) RecordReference(ref string) bool {
	if b.info.IsDocumentationIncluded() {
		b.info.AddReference(ref)
	}
	b.populated = true
	return true
}

// RecordReturnDescription mirrors Builder.recordReturnDescription(String).
func (b *JSDocInfoBuilder) RecordReturnDescription(desc string) bool {
	if !b.info.IsDocumentationIncluded() {
		return true
	}
	if b.info.returnDescription != "" {
		return false
	}
	b.info.returnDescription = desc
	b.populated = true
	return true
}

// RecordParameterDescription mirrors Builder.recordParameterDescription(String, String).
func (b *JSDocInfoBuilder) RecordParameterDescription(name, desc string) bool {
	if !b.info.IsDocumentationIncluded() {
		return true
	}
	if b.info.HasDescriptionForParameter(name) {
		return false
	}
	b.info.SetDescriptionForParameter(name, desc)
	b.populated = true
	return true
}

// RecordLicense mirrors Builder.recordLicense(String).
func (b *JSDocInfoBuilder) RecordLicense(license string) bool {
	b.info.license = license
	b.populated = true
	return true
}

// ShouldParseDocumentation mirrors Builder.shouldParseDocumentation().
func (b *JSDocInfoBuilder) ShouldParseDocumentation() bool {
	return b.info.IsDocumentationIncluded()
}

// ParseDocumentation mirrors Builder.parseDocumentation() — enables
// documentation collection mode.
func (b *JSDocInfoBuilder) ParseDocumentation() *JSDocInfoBuilder {
	b.info.ParseDocumentation()
	return b
}

// IsDescriptionRecorded mirrors Builder.isDescriptionRecorded().
func (b *JSDocInfoBuilder) IsDescriptionRecorded() bool {
	return b.info.description != ""
}

// IsModsRecorded mirrors Builder.isModsRecorded(). The mods namespace is not
// tracked in this port; always returns false.
func (b *JSDocInfoBuilder) IsModsRecorded() bool {
	if b == nil {
		return false
	}
	return b.info.HasMods()
}

// IsPopulatedWithFileOverview mirrors Builder.isPopulatedWithFileOverview().
func (b *JSDocInfoBuilder) IsPopulatedWithFileOverview() bool {
	if !b.populated {
		return false
	}
	return b.info.HasFileOverview() ||
		b.info.IsExterns() ||
		b.info.IsNoCompile() ||
		b.info.IsNoCoverage() ||
		b.info.IsTypeSummary()
}

// RecordFileOverview mirrors Builder.recordFileOverview(String). Records
// the @fileoverview bit and optionally the description text.
func (b *JSDocInfoBuilder) RecordFileOverview(description string) bool {
	if b.info.HasFileOverview() {
		return false
	}
	b.info.SetBit(BitFileoverview, true)
	if b.info.IsDocumentationIncluded() && description != "" {
		b.info.SetFileOverview(description)
	}
	b.populated = true
	return true
}

// RecordNoCollapse mirrors Builder.recordNoCollapse().
func (b *JSDocInfoBuilder) RecordNoCollapse() bool {
	if b.info.CheckBit(BitNoCollapse) {
		return false
	}
	b.info.SetBit(BitNoCollapse, true)
	b.populated = true
	return true
}

// RecordNoInline mirrors Builder.recordNoInline().
func (b *JSDocInfoBuilder) RecordNoInline() bool {
	if b.info.CheckBit(BitNoInline) {
		return false
	}
	b.info.SetBit(BitNoInline, true)
	b.populated = true
	return true
}

// RecordExport mirrors Builder.recordExport().
func (b *JSDocInfoBuilder) RecordExport() bool {
	if b.info.CheckBit(BitExport) {
		return false
	}
	b.info.SetBit(BitExport, true)
	b.populated = true
	return true
}

// RemoveExport mirrors Builder.removeExport().
func (b *JSDocInfoBuilder) RemoveExport() bool {
	if !b.info.CheckBit(BitExport) {
		return false
	}
	b.info.SetBit(BitExport, false)
	b.populated = true
	return true
}

// RecordImplicitCast mirrors Builder.recordImplicitCast().
func (b *JSDocInfoBuilder) RecordImplicitCast() bool {
	if b.info.CheckBit(BitImplicitCast) {
		return false
	}
	b.info.SetBit(BitImplicitCast, true)
	b.populated = true
	return true
}

// RecordExterns mirrors Builder.recordExterns().
func (b *JSDocInfoBuilder) RecordExterns() bool {
	if b.info.CheckBit(BitExterns) {
		return false
	}
	b.info.SetBit(BitExterns, true)
	b.populated = true
	return true
}

// RecordNoCoverage mirrors Builder.recordNoCoverage().
func (b *JSDocInfoBuilder) RecordNoCoverage() bool {
	if b.info.CheckBit(BitNoCoverage) {
		return false
	}
	b.info.SetBit(BitNoCoverage, true)
	b.populated = true
	return true
}

// RecordTypeSummary mirrors Builder.recordTypeSummary().
func (b *JSDocInfoBuilder) RecordTypeSummary() bool {
	if b.info.CheckBit(BitTypeSummary) {
		return false
	}
	b.info.SetBit(BitTypeSummary, true)
	b.populated = true
	return true
}

// RecordPureOrBreakMyCode mirrors Builder.recordPureOrBreakMyCode().
func (b *JSDocInfoBuilder) RecordPureOrBreakMyCode() bool {
	if b.info.CheckBit(BitPureOrBreakMyCode) {
		return false
	}
	b.info.SetBit(BitPureOrBreakMyCode, true)
	b.populated = true
	return true
}

// RecordCollapsibleOrBreakMyCode mirrors Builder.recordCollapsibleOrBreakMyCode().
func (b *JSDocInfoBuilder) RecordCollapsibleOrBreakMyCode() bool {
	if b.info.CheckBit(BitCollapsibleOrBreak) {
		return false
	}
	b.info.SetBit(BitCollapsibleOrBreak, true)
	b.populated = true
	return true
}

// RecordNoCompile mirrors Builder.recordNoCompile().
func (b *JSDocInfoBuilder) RecordNoCompile() bool {
	if b.info.CheckBit(BitNoCompile) {
		return false
	}
	b.info.SetBit(BitNoCompile, true)
	b.populated = true
	return true
}

// RecordNoDts mirrors Builder.recordNoDts().
func (b *JSDocInfoBuilder) RecordNoDts() bool {
	if b.info.CheckBit(BitNoDts) {
		return false
	}
	b.info.SetBit(BitNoDts, true)
	b.populated = true
	return true
}

// RecordEncourageInlining mirrors Builder.recordEncourageInlining().
func (b *JSDocInfoBuilder) RecordEncourageInlining() bool {
	if b.info.CheckBit(BitEncourageInlining) {
		return false
	}
	b.info.SetBit(BitEncourageInlining, true)
	b.populated = true
	return true
}

// RecordRequireInlining mirrors Builder.recordRequireInlining().
func (b *JSDocInfoBuilder) RecordRequireInlining() bool {
	if b.info.CheckBit(BitRequireInlining) {
		return false
	}
	b.info.SetBit(BitRequireInlining, true)
	b.populated = true
	return true
}

// RecordProvideGoog mirrors Builder.recordProvideGoog().
func (b *JSDocInfoBuilder) RecordProvideGoog() bool {
	if b.info.CheckBit(BitProvideGoog) {
		return false
	}
	b.info.SetBit(BitProvideGoog, true)
	b.populated = true
	return true
}

// RecordProvideAlreadyProvided mirrors Builder.recordProvideAlreadyProvided().
func (b *JSDocInfoBuilder) RecordProvideAlreadyProvided() bool {
	if b.info.CheckBit(BitProvideAlreadyProvided) {
		return false
	}
	b.info.SetBit(BitProvideAlreadyProvided, true)
	b.populated = true
	return true
}

// RecordWizcallback mirrors Builder.recordWizcallback().
func (b *JSDocInfoBuilder) RecordWizcallback() bool {
	if b.info.CheckBit(BitWizCallback) {
		return false
	}
	b.info.SetBit(BitWizCallback, true)
	b.populated = true
	return true
}

// RecordMixinClass mirrors Builder.recordMixinClass().
func (b *JSDocInfoBuilder) RecordMixinClass() bool {
	if b.info.CheckBit(BitMixinClass) {
		return false
	}
	b.info.SetBit(BitMixinClass, true)
	b.populated = true
	return true
}

// RecordMixinFunction mirrors Builder.recordMixinFunction().
func (b *JSDocInfoBuilder) RecordMixinFunction() bool {
	if b.info.CheckBit(BitMixinFunction) {
		return false
	}
	b.info.SetBit(BitMixinFunction, true)
	b.populated = true
	return true
}

// RecordPolymer mirrors Builder.recordPolymer().
func (b *JSDocInfoBuilder) RecordPolymer() bool {
	if b.info.CheckBit(BitPolymer) {
		return false
	}
	b.info.SetBit(BitPolymer, true)
	b.populated = true
	return true
}

// RecordUnrestricted mirrors Builder.recordUnrestricted().
func (b *JSDocInfoBuilder) RecordUnrestricted() bool {
	if b.info.CheckBit(BitUnrestricted) || b.hasAnySingletonTypeTags() {
		return false
	}
	if b.info.isInterface || b.info.isDict || b.info.isStruct {
		return false
	}
	b.info.SetBit(BitUnrestricted, true)
	b.populated = true
	return true
}

// IsUnrestrictedRecorded mirrors Builder.isUnrestrictedRecorded().
func (b *JSDocInfoBuilder) IsUnrestrictedRecorded() bool {
	return b.info.CheckBit(BitUnrestricted)
}

// RecordThisType mirrors Builder.recordThisType with a full type-expression.
// Wraps SetThisType if the type is non-nil and no singleton type is set.
func (b *JSDocInfoBuilder) RecordThisTypeExpr(t *JSTypeExpression) bool {
	if t == nil || b.info.hasThisType || b.hasAnySingletonTypeTags() {
		return false
	}
	b.info.SetThisType(t)
	b.populated = true
	return true
}

// IsNgInline mirrors Builder.isNoInline() — returns true when @noinline set.
func (b *JSDocInfoBuilder) IsNoInline() bool {
	return b.info.CheckBit(BitNoInline)
}

// IsRequireInlining mirrors Builder.isRequireInlining().
func (b *JSDocInfoBuilder) IsRequireInlining() bool {
	return b.info.CheckBit(BitRequireInlining)
}

// IsEncourageInlining mirrors Builder.isEncourageInlining().
func (b *JSDocInfoBuilder) IsEncourageInlining() bool {
	return b.info.CheckBit(BitEncourageInlining)
}

// IsPolymerBehaviorRecorded mirrors Builder.isPolymerBehaviorRecorded().
func (b *JSDocInfoBuilder) IsPolymerBehaviorRecorded() bool {
	return b.info.CheckBit(BitPolymerBehavior)
}

// IsPolymerRecorded mirrors Builder.isPolymerRecorded().
func (b *JSDocInfoBuilder) IsPolymerRecorded() bool {
	return b.info.CheckBit(BitPolymer)
}

// IsCustomElementRecorded mirrors Builder.isCustomElementRecorded().
func (b *JSDocInfoBuilder) IsCustomElementRecorded() bool {
	return b.info.CheckBit(BitCustomElement)
}

// IsMixinClassRecorded mirrors Builder.isMixinClassRecorded().
func (b *JSDocInfoBuilder) IsMixinClassRecorded() bool {
	return b.info.CheckBit(BitMixinClass)
}

// IsMixinFunctionRecorded mirrors Builder.isMixinFunctionRecorded().
func (b *JSDocInfoBuilder) IsMixinFunctionRecorded() bool {
	return b.info.CheckBit(BitMixinFunction)
}

// IsWizcallbackRecorded mirrors Builder.isWizcallbackRecorded().
func (b *JSDocInfoBuilder) IsWizcallbackRecorded() bool {
	return b.info.CheckBit(BitWizCallback)
}

// IsSassGeneratedCssTsRecorded mirrors Builder.isSassGeneratedCssTsRecorded().
func (b *JSDocInfoBuilder) IsSassGeneratedCssTsRecorded() bool {
	return b.info.CheckBit(BitSassGeneratedCssTs)
}

// RecordSassGeneratedCssTs mirrors Builder.recordSassGeneratedCssTs().
func (b *JSDocInfoBuilder) RecordSassGeneratedCssTs() {
	b.info.SetBit(BitSassGeneratedCssTs, true)
	b.populated = true
}

// IsDeprecationReasonRecorded mirrors Builder.isDeprecationReasonRecorded().
func (b *JSDocInfoBuilder) IsDeprecationReasonRecorded() bool {
	return b.info.deprecationReason != ""
}

// RecordStableIdGenerator mirrors Builder.recordStableIdGenerator().
func (b *JSDocInfoBuilder) RecordStableIdGenerator() bool {
	if !b.info.SetIdGenerator(IdGeneratorStable) {
		return false
	}
	b.populated = true
	return true
}

// RecordMappedIdGenerator mirrors Builder.recordMappedIdGenerator().
func (b *JSDocInfoBuilder) RecordMappedIdGenerator() bool {
	if !b.info.SetIdGenerator(IdGeneratorMapped) {
		return false
	}
	b.populated = true
	return true
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

// hasAnySingletonSideEffectTags mirrors Builder.hasAnySingletonSideEffectTags()
// — prevents @modifies from coexisting with @nosideeffects.
func (b *JSDocInfoBuilder) hasAnySingletonSideEffectTags() bool {
	return b.info.isNoSideEffects
}

// isConstructorOrInterface mirrors Builder.isConstructorOrInterface().
func (b *JSDocInfoBuilder) isConstructorOrInterface() bool {
	return b.info.isConstructor || b.info.isInterface
}
