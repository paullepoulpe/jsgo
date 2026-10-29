// Port of com/google/javascript/jscomp/ConformanceRules.java @ closure-compiler v20261021

/*
 * Copyright 2014 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jscomp

import (
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Java: ConformanceRules.GOOG_DOM
var conformanceRules_GOOG_DOM = rhino.QualifiedName_Of("goog.dom")

// Java: ConformanceRules.GOOG_DOM_TAGNAME
var conformanceRules_GOOG_DOM_TAGNAME = rhino.QualifiedName_Of("goog.dom.TagName")

// ConformanceRules_BanCreateElement bans document.createElement and similar
// methods with string literal parameter specified in value, e.g. value:
// 'script'.
type ConformanceRules_BanCreateElement struct {
	ConformanceRules_AbstractRule

	bannedTags    map[string]struct{}
	domHelperType rhino.JSType
	documentType  rhino.JSType
}

// Java: ConformanceRules.BanCreateElement(AbstractCompiler, Requirement)
func newConformanceRules_BanCreateElement(compiler AbstractCompiler, requirement *Requirement) (*ConformanceRules_BanCreateElement, error) {
	r := &ConformanceRules_BanCreateElement{}
	if err := r.initAbstractRule(compiler, requirement, false); err != nil {
		return nil, err
	}
	values := requirement.GetValueList()
	if len(values) == 0 {
		return nil, errInvalidRequirementSpec("Specify one or more values.")
	}
	bannedTags := make(map[string]struct{}, len(values))
	for _, value := range values {
		bannedTags[jutil.StringToLowerCaseASCII(value)] = struct{}{}
	}
	r.bannedTags = bannedTags
	r.domHelperType = compiler.GetTypeRegistry().GetGlobalType("goog.dom.DomHelper")
	r.documentType = compiler.GetTypeRegistry().GetGlobalType("Document")
	return r, nil
}

func (r *ConformanceRules_BanCreateElement) Check(
	t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	if result := r.CheckConformance(t, n); result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
		r.Report(n, result, behavior)
	}
}

// Java: ConformanceRules.BanCreateElement.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BanCreateElement) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if call, ok := ast.AsCall(n); ok {
		if firstArg, ok := call.FirstArg(); ok {
			if tag, ok := ast.AsStringLit(firstArg.Node()); ok {
				if _, banned := r.bannedTags[jutil.StringToLowerCaseASCII(tag.Value())]; banned {
					return r.checkCreateElement(call)
				}
			}
		}
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// Java: ConformanceRules.BanCreateElement.checkCreateElement(Node)
func (r *ConformanceRules_BanCreateElement) checkCreateElement(call ast.Call) *ConformanceRules_ConformanceResult {
	target, ok := ast.AsGetProp(call.Callee().Node())
	if !ok {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}
	functionName := target.Property()
	if functionName != "createElement" && functionName != "createDom" {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	srcObj := target.Object().Node()
	if conformanceRules_GOOG_DOM.Matches(srcObj) {
		return ConformanceRules_ConformanceResult_VIOLATION
	}
	type_ := srcObj.GetJSType()
	if type_ == nil || ConformanceRules_ConformanceUtil_IsLooseType(type_) {
		if r.reportLooseTypeViolations {
			return conformanceRules_ConformanceResult_POSSIBLE_VIOLATION_DUE_TO_LOOSE_TYPES
		}
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}
	if (r.domHelperType != nil && r.domHelperType.IsSubtypeOf(type_)) ||
		(r.documentType != nil && r.documentType.IsSubtypeOf(type_)) {
		return ConformanceRules_ConformanceResult_VIOLATION
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// ConformanceRules_BanCreateDom bans goog.dom.createDom and
// goog.dom.DomHelper#createDom with parameters specified in value in the
// format tagname.attribute, e.g. value: 'iframe.src'.
type ConformanceRules_BanCreateDom struct {
	ConformanceRules_AbstractRule

	bannedTagAttrs [][2]string
	domHelperType  rhino.JSType
	classNameTypes rhino.JSType
}

// Java: ConformanceRules.BanCreateDom(AbstractCompiler, Requirement)
func newConformanceRules_BanCreateDom(compiler AbstractCompiler, requirement *Requirement) (*ConformanceRules_BanCreateDom, error) {
	r := &ConformanceRules_BanCreateDom{}
	if err := r.initAbstractRule(compiler, requirement, false); err != nil {
		return nil, err
	}
	values := requirement.GetValueList()
	bannedTagAttrs := make([][2]string, 0, len(values))
	for _, value := range values {
		tagAttr := conformanceRules_splitToList(conformanceRules_ON_DOT, value)
		if len(tagAttr) != 2 || tagAttr[0] == "" || tagAttr[1] == "" {
			return nil, errInvalidRequirementSpec("Values must be in the format tagname.attribute.")
		}
		lowercasedTag := jutil.StringToLowerCaseRoot(tagAttr[0])
		bannedTagAttrs = append(bannedTagAttrs, [2]string{lowercasedTag, tagAttr[1]})
	}
	if len(bannedTagAttrs) == 0 {
		return nil, errInvalidRequirementSpec("Specify one or more values.")
	}
	r.bannedTagAttrs = bannedTagAttrs
	registry := compiler.GetTypeRegistry()
	r.domHelperType = registry.GetGlobalType("goog.dom.DomHelper")
	r.classNameTypes = registry.CreateUnionType(
		registry.GetNativeType(rhino.JSTypeNative_STRING_TYPE),
		registry.GetNativeType(rhino.JSTypeNative_ARRAY_TYPE),
		registry.GetNativeType(rhino.JSTypeNative_NULL_TYPE),
		registry.GetNativeType(rhino.JSTypeNative_VOID_TYPE))
	return r, nil
}

func (r *ConformanceRules_BanCreateDom) Check(
	t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	if result := r.CheckConformance(t, n); result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
		r.Report(n, result, behavior)
	}
}

// Java: ConformanceRules.BanCreateDom.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BanCreateDom) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	call, ok := r.asCreateDomCall(n)
	if !ok {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}
	firstArg, hasFirst := call.FirstArg()
	if !hasFirst {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}
	attrsNode := firstArg.Node().GetNext()
	if attrsNode.IsNil() {
		// goog.dom.createDom('iframe') is fine.
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	tagNames, tagNamesIsNull := r.GetTagNames(firstArg.Node())
	attrsType := attrsNode.GetJSType()
	if attrsType == nil {
		// Type information is not available; don't run this check.
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	// String or array attribute sets the class.
	isClassName := !attrsType.IsUnknownType() && attrsType.IsSubtypeOf(r.classNameTypes)

	if attrsNode.IsNull() || attrsType.IsVoidType() {
		// goog.dom.createDom('iframe', null) is fine.
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	hasExtraChildren := !attrsNode.GetNext().IsNil()

	for _, tagAttr := range r.bannedTagAttrs {
		if !tagNamesIsNull && !slices.Contains(tagNames, tagAttr[0]) && tagAttr[0] != "*" {
			continue
		}
		var violation *ConformanceRules_ConformanceResult
		if !tagNamesIsNull || tagAttr[0] == "*" {
			violation = ConformanceRules_ConformanceResult_VIOLATION
		} else if r.reportLooseTypeViolations {
			violation = ConformanceRules_ConformanceResult_POSSIBLE_VIOLATION
		} else {
			violation = ConformanceRules_ConformanceResult_CONFORMANCE
		}
		if isClassName {
			if tagAttr[1] != "class" {
				continue
			}
			return violation
		}
		if tagAttr[1] == "textContent" &&
			hasExtraChildren &&
			violation != ConformanceRules_ConformanceResult_CONFORMANCE {
			return violation
		}
		objLit, isObjLit := ast.AsObjectLit(attrsNode)
		if !isObjLit {
			// Attrs is not an object literal and tagName matches or is unknown.
			if r.reportLooseTypeViolations {
				return ConformanceRules_ConformanceResult_POSSIBLE_VIOLATION
			}
			return ConformanceRules_ConformanceResult_CONFORMANCE
		}
		prop := NodeUtil_GetFirstPropMatchingKey(attrsNode, tagAttr[1])
		if !prop.IsNil() {
			if NodeUtil_IsSomeCompileTimeConstStringValue(prop) {
				// Ignore string literal values.
				continue
			}
			return violation
		}

		for attr := objLit.Node().GetFirstChild(); !attr.IsNil(); attr = attr.GetNext() {
			if attr.IsComputedProp() {
				// We don't know if the computed property matches 'src' or not
				if r.reportLooseTypeViolations {
					return conformanceRules_ConformanceResult_POSSIBLE_VIOLATION_DUE_TO_LOOSE_TYPES
				}
				return ConformanceRules_ConformanceResult_CONFORMANCE
			}
		}
	}

	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// Java: ConformanceRules.BanCreateDom.getTagNames(Node)
func (r *ConformanceRules_BanCreateDom) GetTagNames(tag rhino.Node) (tagNames []string, isNull bool) {
	if s, ok := ast.AsStringLit(tag); ok {
		return []string{jutil.StringToLowerCaseRoot(s.Value())}, false
	} else if gp, ok := ast.AsGetProp(tag); ok && conformanceRules_GOOG_DOM_TAGNAME.Matches(gp.Object().Node()) {
		return []string{jutil.StringToLowerCaseRoot(gp.Property())}, false
	}
	type_ := tag.GetJSType()
	if type_ == nil || !type_.IsTemplatizedType() {
		return nil, true
	}
	typeAsObj := type_.ToMaybeObjectType()
	if *typeAsObj.GetRawType().GetDisplayName() == "goog.dom.TagName" {
		tagType := conformanceRules_getOnlyElement(typeAsObj.GetTemplateTypes())
		var names []string
		if displayName := tagType.GetDisplayName(); displayName != nil {
			names = conformanceRules_elementTagNames[*displayName]
		}
		return names, false
	}
	return nil, true
}

func conformanceRules_getOnlyElement(types []rhino.JSType) rhino.JSType {
	if len(types) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	if len(types) == 1 {
		return types[0]
	}
	var sb strings.Builder
	sb.WriteString("expected one element but was: <")
	sb.WriteString(types[0].String())
	i := 1
	for ; i < 5 && i < len(types); i++ {
		sb.WriteString(", ")
		sb.WriteString(types[i].String())
	}
	if i < len(types) {
		sb.WriteString(", ...")
	}
	sb.WriteByte('>')
	panic(jutil.NewIllegalArgumentExceptionMsg(sb.String()))
}

// Java: ConformanceRules.BanCreateDom.ELEMENT_TAG_NAMES (an ImmutableMultimap)
var conformanceRules_elementTagNames = map[string][]string{
	"HTMLAnchorElement":       {"a"},
	"HTMLAppletElement":       {"applet"},
	"HTMLAreaElement":         {"area"},
	"HTMLAudioElement":        {"audio"},
	"HTMLBRElement":           {"br"},
	"HTMLBaseElement":         {"base"},
	"HTMLBaseFontElement":     {"basefont"},
	"HTMLBodyElement":         {"body"},
	"HTMLButtonElement":       {"button"},
	"HTMLCanvasElement":       {"canvas"},
	"HTMLDListElement":        {"dl"},
	"HTMLDataListElement":     {"datalist"},
	"HTMLDetailsElement":      {"details"},
	"HTMLDialogElement":       {"dialog"},
	"HTMLDirectoryElement":    {"dir"},
	"HTMLDivElement":          {"div"},
	"HTMLEmbedElement":        {"embed"},
	"HTMLFieldSetElement":     {"fieldset"},
	"HTMLFontElement":         {"font"},
	"HTMLFormElement":         {"form"},
	"HTMLFrameElement":        {"frame"},
	"HTMLFrameSetElement":     {"frameset"},
	"HTMLHRElement":           {"hr"},
	"HTMLHeadElement":         {"head"},
	"HTMLHeadingElement":      {"h1", "h2", "h3", "h4", "h5", "h6"},
	"HTMLHtmlElement":         {"html"},
	"HTMLIFrameElement":       {"iframe"},
	"HTMLImageElement":        {"img"},
	"HTMLInputElement":        {"input"},
	"HTMLIsIndexElement":      {"isindex"},
	"HTMLLIElement":           {"li"},
	"HTMLLabelElement":        {"label"},
	"HTMLLegendElement":       {"legend"},
	"HTMLLinkElement":         {"link"},
	"HTMLMapElement":          {"map"},
	"HTMLMenuElement":         {"menu"},
	"HTMLMetaElement":         {"meta"},
	"HTMLMeterElement":        {"meter"},
	"HTMLModElement":          {"del", "ins"},
	"HTMLOListElement":        {"ol"},
	"HTMLObjectElement":       {"object"},
	"HTMLOptGroupElement":     {"optgroup"},
	"HTMLOptionElement":       {"option"},
	"HTMLOutputElement":       {"output"},
	"HTMLParagraphElement":    {"p"},
	"HTMLParamElement":        {"param"},
	"HTMLPreElement":          {"pre"},
	"HTMLProgressElement":     {"progress"},
	"HTMLQuoteElement":        {"blockquote", "q"},
	"HTMLScriptElement":       {"script"},
	"HTMLSelectElement":       {"select"},
	"HTMLSourceElement":       {"source"},
	"HTMLSpanElement":         {"span"},
	"HTMLStyleElement":        {"style"},
	"HTMLTableCaptionElement": {"caption"},
	"HTMLTableCellElement":    {"td", "th"},
	"HTMLTableColElement":     {"col", "colgroup"},
	"HTMLTableElement":        {"table"},
	"HTMLTableRowElement":     {"tr"},
	"HTMLTableSectionElement": {"tbody", "tfoot", "thead"},
	"HTMLTemplateElement":     {"template"},
	"HTMLTextAreaElement":     {"textarea"},
	"HTMLTitleElement":        {"title"},
	"HTMLTrackElement":        {"track"},
	"HTMLUListElement":        {"ul"},
	"HTMLVideoElement":        {"video"},
}

// Java: ConformanceRules.BanCreateDom.isCreateDomCall(Node)
func (r *ConformanceRules_BanCreateDom) asCreateDomCall(n rhino.Node) (ast.Call, bool) {
	call, ok := ast.AsCall(n)
	if !ok {
		return ast.Call{}, false
	}
	target, ok := ast.AsGetProp(call.Callee().Node())
	if !ok || target.Property() != "createDom" {
		return ast.Call{}, false
	}

	srcObj := target.Object().Node()
	if conformanceRules_GOOG_DOM.Matches(srcObj) {
		return call, true
	}
	type_ := srcObj.GetJSType()
	if type_ != nil && type_.Equals(r.domHelperType) {
		return call, true
	}
	return ast.Call{}, false
}

// ConformanceRules_SecuritySensitiveAttributes checks nodes for conformance
// with banning the setting of attributes that are on the blocklist.
type ConformanceRules_SecuritySensitiveAttributes struct {
	bannedAttrs             map[string]struct{}
	globalNamespaceSupplier func() *GlobalNamespace
}

// Java: ConformanceRules.SecuritySensitiveAttributes.ALL_BANNED_ATTRS
var conformanceRules_allBannedAttrs = [...]string{
	"href",
	"rel",
	"src",
	"srcdoc",
	"action",
	"formaction",
	"sandbox",
	"icon",
	"codebase",
	"data",
}

// Java: ConformanceRules.SecuritySensitiveAttributes()
func NewConformanceRules_SecuritySensitiveAttributes() *ConformanceRules_SecuritySensitiveAttributes {
	return NewConformanceRules_SecuritySensitiveAttributesWithBannedAttrsAndSupplier(
		conformanceRules_allBannedAttrs[:], func() *GlobalNamespace { return nil })
}

// Java: ConformanceRules.SecuritySensitiveAttributes(Collection<String>, Supplier<GlobalNamespace>)
func NewConformanceRules_SecuritySensitiveAttributesWithBannedAttrsAndSupplier(
	bannedAttrs []string, globalNamespaceSupplier func() *GlobalNamespace) *ConformanceRules_SecuritySensitiveAttributes {
	m := make(map[string]struct{}, len(bannedAttrs))
	for _, a := range bannedAttrs {
		m[jutil.StringToLowerCaseRoot(a)] = struct{}{}
	}
	return &ConformanceRules_SecuritySensitiveAttributes{
		bannedAttrs:             m,
		globalNamespaceSupplier: globalNamespaceSupplier,
	}
}

// Java: ConformanceRules.SecuritySensitiveAttributes.contains(String)
func (s *ConformanceRules_SecuritySensitiveAttributes) Contains(attributeName string) bool {
	_, ok := s.bannedAttrs[attributeName]
	return ok
}

// Java: ConformanceRules.SecuritySensitiveAttributes.checkConformanceForAttributeName(NodeTraversal, Node)
func (s *ConformanceRules_SecuritySensitiveAttributes) CheckConformanceForAttributeName(
	traversal *NodeTraversal, attrName rhino.Node) *ConformanceRules_ConformanceResult {
	literalName :=
		ConformanceRules_ConformanceUtil_InferStringValue(traversal.GetScope(), attrName, s.globalNamespaceSupplier)
	if literalName == nil {
		// xid() obfuscates attribute names, thus never clashing with security-sensitive attributes.
		if ConformanceRules_ConformanceUtil_IsXid(attrName.GetJSType()) {
			return ConformanceRules_ConformanceResult_CONFORMANCE
		}
		return ConformanceRules_ConformanceResult_VIOLATION
	}

	lowerName := jutil.StringToLowerCaseRoot(*literalName)
	if s.Contains(lowerName) || ConformanceRules_ConformanceUtil_IsEventHandlerAttrName(lowerName) {
		return ConformanceRules_ConformanceResult_VIOLATION
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// Java: ConformanceRules.SecuritySensitiveAttributes.checkConformanceForAttributeNameWithHighConfidence(NodeTraversal, Node)
func (s *ConformanceRules_SecuritySensitiveAttributes) CheckConformanceForAttributeNameWithHighConfidence(
	traversal *NodeTraversal, attrName rhino.Node) *ConformanceRules_ConformanceResult {
	literalName :=
		ConformanceRules_ConformanceUtil_InferStringValue(traversal.GetScope(), attrName, s.globalNamespaceSupplier)
	if literalName == nil {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}
	if s.Contains(jutil.StringToLowerCaseRoot(*literalName)) {
		return ConformanceRules_ConformanceResult_VIOLATION
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// ConformanceRules_BanElementSetAttribute bans Element#setAttribute with
// attribute names specified in value or any dynamic string.
type ConformanceRules_BanElementSetAttribute struct {
	ConformanceRules_AbstractRule

	performGlobalNamespaceAnalysis bool
	globalNamespace                *GlobalNamespace

	elementType                      rhino.JSType
	securitySensitiveAttributes      *ConformanceRules_SecuritySensitiveAttributes
	defaultDecisionForUncertainCases *ConformanceRules_ConformanceResult
}

const (
	conformanceRules_BanElementSetAttribute_ELEMENT_TYPE_NAME               = "Element"
	conformanceRules_BanElementSetAttribute_SET_ATTRIBUTE                   = "setAttribute"
	conformanceRules_BanElementSetAttribute_SET_ATTRIBUTE_NS                = "setAttributeNS"
	conformanceRules_BanElementSetAttribute_SET_ATTRIBUTE_NODE              = "setAttributeNode"
	conformanceRules_BanElementSetAttribute_SET_ATTRIBUTE_NODE_NS           = "setAttributeNodeNS"
	conformanceRules_BanElementSetAttribute_GLOBAL_NAMESPACE_ANALYSIS_LIMIT = 10000
)

func conformanceRules_isBannedSetAttributeProp(name string) bool {
	switch name {
	case conformanceRules_BanElementSetAttribute_SET_ATTRIBUTE,
		conformanceRules_BanElementSetAttribute_SET_ATTRIBUTE_NS,
		conformanceRules_BanElementSetAttribute_SET_ATTRIBUTE_NODE,
		conformanceRules_BanElementSetAttribute_SET_ATTRIBUTE_NODE_NS:
		return true
	default:
		return false
	}
}

// Java: ConformanceRules.BanElementSetAttribute(AbstractCompiler, Requirement)
func newConformanceRules_BanElementSetAttribute(compiler AbstractCompiler, requirement *Requirement) (*ConformanceRules_BanElementSetAttribute, error) {
	r := &ConformanceRules_BanElementSetAttribute{
		performGlobalNamespaceAnalysis: true,
	}
	if err := r.initAbstractRule(compiler, requirement, false); err != nil {
		return nil, err
	}
	r.securitySensitiveAttributes =
		NewConformanceRules_SecuritySensitiveAttributesWithBannedAttrsAndSupplier(requirement.GetValueList(), r.GetGlobalNamespace)
	if requirement.GetReportLooseTypeViolations() {
		r.defaultDecisionForUncertainCases = ConformanceRules_ConformanceResult_VIOLATION
	} else {
		r.defaultDecisionForUncertainCases = ConformanceRules_ConformanceResult_CONFORMANCE
	}
	r.elementType = compiler.GetTypeRegistry().GetGlobalType(conformanceRules_BanElementSetAttribute_ELEMENT_TYPE_NAME)
	return r, nil
}

// Java: ConformanceRules.BanElementSetAttribute.getGlobalNamespace()
func (r *ConformanceRules_BanElementSetAttribute) GetGlobalNamespace() *GlobalNamespace {
	if r.performGlobalNamespaceAnalysis && r.globalNamespace == nil {
		if NodeUtil_CountAstSizeUpToLimit(r.compiler.GetJsRoot(), conformanceRules_BanElementSetAttribute_GLOBAL_NAMESPACE_ANALYSIS_LIMIT) >=
			conformanceRules_BanElementSetAttribute_GLOBAL_NAMESPACE_ANALYSIS_LIMIT {
			r.performGlobalNamespaceAnalysis = false
		} else {
			r.globalNamespace = NewGlobalNamespace(r.compiler, r.compiler.GetJsRoot())
		}
	}
	return r.globalNamespace
}

func (r *ConformanceRules_BanElementSetAttribute) Check(
	t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	if result := r.CheckConformance(t, n); result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
		r.Report(n, result, behavior)
	}
}

// Java: ConformanceRules.BanElementSetAttribute.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BanElementSetAttribute) CheckConformance(traversal *NodeTraversal, node rhino.Node) *ConformanceRules_ConformanceResult {
	if call, ok := ast.AsCall(node); ok {
		return r.checkConformanceOnPropertyCall(traversal, call)
	} else if ge, ok := ast.AsGetElem(node); ok && NodeUtil_IsLValue(node) {
		return r.checkConformanceOnGetElement(traversal, ge)
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// Java: ConformanceRules.BanElementSetAttribute.checkConformanceOnPropertyCall(NodeTraversal, Node)
func (r *ConformanceRules_BanElementSetAttribute) checkConformanceOnPropertyCall(
	traversal *NodeTraversal, call ast.Call) *ConformanceRules_ConformanceResult {
	calledProperty, ok := r.getBannedPropertyName(call)
	if !ok {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	if calledProperty == conformanceRules_BanElementSetAttribute_SET_ATTRIBUTE {
		firstArg, hasFirst := call.FirstArg()
		if !hasFirst || firstArg.Node().GetNext().IsNil() {
			return ConformanceRules_ConformanceResult_CONFORMANCE
		}

		if r.requirement.GetReportLooseTypeViolations() {
			return r.securitySensitiveAttributes.CheckConformanceForAttributeName(
				traversal, firstArg.Node())
		}
		return r.securitySensitiveAttributes.CheckConformanceForAttributeNameWithHighConfidence(
			traversal, firstArg.Node())
	}

	if calledProperty == conformanceRules_BanElementSetAttribute_SET_ATTRIBUTE_NS {
		firstArg, hasFirst := call.FirstArg()
		if !hasFirst {
			return ConformanceRules_ConformanceResult_CONFORMANCE
		}
		secondArg := firstArg.Node().GetNext()
		if secondArg.IsNil() || secondArg.GetNext().IsNil() {
			return ConformanceRules_ConformanceResult_CONFORMANCE
		}

		if firstArg.Node().IsNull() {
			// A null namespace makes this equivalent to a regular setAttribute call.
			if r.requirement.GetReportLooseTypeViolations() {
				return r.securitySensitiveAttributes.CheckConformanceForAttributeName(
					traversal, secondArg)
			}
			return r.securitySensitiveAttributes.CheckConformanceForAttributeNameWithHighConfidence(
				traversal, secondArg)
		}
	}

	// We haven't defined a security contract on setAttributeNode and setAttributeNodeNS yet, so
	// flag them as risky if report_loose_type_violations is set.
	return r.defaultDecisionForUncertainCases
}

// Java: ConformanceRules.BanElementSetAttribute.getBannedPropertyName(Node)
func (r *ConformanceRules_BanElementSetAttribute) getBannedPropertyName(call ast.Call) (string, bool) {
	target, ok := ast.AsGetProp(call.Callee().Node())
	if !ok {
		return "", false
	}
	propertyName := target.Property()
	if !conformanceRules_isBannedSetAttributeProp(propertyName) {
		return "", false
	}
	type_ := target.Object().Node().GetJSType()
	if type_ == nil {
		return "", false
	}
	if r.elementType != nil && type_.RestrictByNotNullOrUndefined().IsSubtypeOf(r.elementType) {
		return propertyName, true
	}
	return "", false
}

// Java: ConformanceRules.BanElementSetAttribute.checkConformanceOnGetElement(NodeTraversal, Node)
func (r *ConformanceRules_BanElementSetAttribute) checkConformanceOnGetElement(
	traversal *NodeTraversal, ge ast.GetElem) *ConformanceRules_ConformanceResult {
	key := ge.Index().Node()
	keyName :=
		ConformanceRules_ConformanceUtil_InferStringValue(traversal.GetScope(), key, r.GetGlobalNamespace)
	if keyName != nil {
		if r.securitySensitiveAttributes.Contains(jutil.StringToLowerCaseRoot(*keyName)) ||
			(r.requirement.GetReportLooseTypeViolations() &&
				(*keyName == "innerHTML" || *keyName == "outerHTML")) {
			if r.hasElementType(ge) {
				return ConformanceRules_ConformanceResult_VIOLATION
			}
		}
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	if !r.requirement.GetReportLooseTypeViolations() {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	keyType := key.GetJSType()
	if keyType == nil || ConformanceRules_ConformanceUtil_IsXid(keyType) {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}
	if r.hasElementType(ge) {
		if keyType.IsString() {
			return ConformanceRules_ConformanceResult_VIOLATION
		} else if keyType.IsUnionType() {
			for _, alternate := range keyType.ToMaybeUnionType().GetAlternates() {
				if alternate.IsString() {
					return ConformanceRules_ConformanceResult_VIOLATION
				}
			}
		}
	}

	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// Java: ConformanceRules.BanElementSetAttribute.hasElementType(Node)
func (r *ConformanceRules_BanElementSetAttribute) hasElementType(ge ast.GetElem) bool {
	objType := ge.Object().Node().GetJSType()
	if objType == nil || r.elementType == nil {
		return false
	}
	objType = objType.RestrictByNotNullOrUndefined()
	if objType.IsUnknownType() || objType.IsEmptyType() || !objType.IsSubtypeOf(r.elementType) {
		return false
	}
	return true
}

// ConformanceRules_BanSettingAttributes can be used to ban any function like
// Element#setAttribute which takes two parameters, an attribute name and a
// value.
type ConformanceRules_BanSettingAttributes struct {
	ConformanceRules_AbstractRule

	restrictions []ConformanceRules_BanSettingAttributes_AttributeSettingRestriction
}

// Java: ConformanceRules.BanSettingAttributes.AttributeSettingRestriction
type ConformanceRules_BanSettingAttributes_AttributeSettingRestriction struct {
	type_    rhino.JSType
	property string
}

// Java: ConformanceRules.BanSettingAttributes.BANNED_ATTRS
var conformanceRules_BanSettingAttributes_BANNED_ATTRS = NewConformanceRules_SecuritySensitiveAttributes()

// Java: ConformanceRules.BanSettingAttributes(AbstractCompiler, Requirement)
func newConformanceRules_BanSettingAttributes(compiler AbstractCompiler, requirement *Requirement) (*ConformanceRules_BanSettingAttributes, error) {
	r := &ConformanceRules_BanSettingAttributes{}
	if err := r.initAbstractRule(compiler, requirement, false); err != nil {
		return nil, err
	}

	values := requirement.GetValueList()
	if len(values) == 0 {
		return nil, errInvalidRequirementSpec("missing value")
	}

	registry := compiler.GetTypeRegistry()
	restrictions := make([]ConformanceRules_BanSettingAttributes_AttributeSettingRestriction, 0, len(values))
	for _, value := range values {
		typeName, property, ok := conformanceRules_splitClassAndProperty(value)
		if !ok {
			return nil, errInvalidRequirementSpec("bad prop value")
		}
		restrictions = append(restrictions, ConformanceRules_BanSettingAttributes_AttributeSettingRestriction{
			type_:    registry.GetGlobalType(typeName),
			property: property,
		})
	}

	r.restrictions = restrictions
	return r, nil
}

func (r *ConformanceRules_BanSettingAttributes) Check(
	t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	if result := r.CheckConformance(t, n); result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
		r.Report(n, result, behavior)
	}
}

// Java: ConformanceRules.BanSettingAttributes.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BanSettingAttributes) CheckConformance(traversal *NodeTraversal, node rhino.Node) *ConformanceRules_ConformanceResult {
	if call, ok := ast.AsCall(node); ok {
		return r.checkConformanceOnPropertyCall(traversal, call)
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// Java: ConformanceRules.BanSettingAttributes.checkConformanceOnPropertyCall(NodeTraversal, Node)
func (r *ConformanceRules_BanSettingAttributes) checkConformanceOnPropertyCall(traversal *NodeTraversal, call ast.Call) *ConformanceRules_ConformanceResult {
	if !r.matchesBannedProperty(call) {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	firstArg, hasFirst := call.FirstArg()
	if !hasFirst || firstArg.Node().GetNext().IsNil() {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	return conformanceRules_BanSettingAttributes_BANNED_ATTRS.CheckConformanceForAttributeName(traversal, firstArg.Node())
}

// Java: ConformanceRules.BanSettingAttributes.getBannedPropertyName(Node)
func (r *ConformanceRules_BanSettingAttributes) matchesBannedProperty(call ast.Call) bool {
	target, ok := ast.AsGetProp(call.Callee().Node())
	if !ok {
		return false
	}
	type_ := target.Object().Node().GetJSType()
	if type_ == nil {
		return false
	}
	propertyName := target.Property()
	for i := range r.restrictions {
		restricted := &r.restrictions[i]
		if propertyName == restricted.property &&
			restricted.type_ != nil &&
			type_.IsSubtypeOf(restricted.type_) {
			return true
		}
	}
	return false
}

// ConformanceRules_BanExecCommand bans Document#execCommand with command
// names specified in value or any dynamic string.
type ConformanceRules_BanExecCommand struct {
	ConformanceRules_AbstractRule

	bannedAttrs []string
}

// Java: ConformanceRules.BanExecCommand(AbstractCompiler, Requirement)
func newConformanceRules_BanExecCommand(compiler AbstractCompiler, requirement *Requirement) (*ConformanceRules_BanExecCommand, error) {
	r := &ConformanceRules_BanExecCommand{}
	if err := r.initAbstractRule(compiler, requirement, false); err != nil {
		return nil, err
	}
	values := requirement.GetValueList()
	bannedAttrs := make([]string, 0, len(values))
	for _, v := range values {
		bannedAttrs = append(bannedAttrs, jutil.StringToLowerCaseRoot(v))
	}
	r.bannedAttrs = bannedAttrs
	return r, nil
}

func (r *ConformanceRules_BanExecCommand) Check(
	t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	if result := r.CheckConformance(t, n); result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
		r.Report(n, result, behavior)
	}
}

// Java: ConformanceRules.BanExecCommand.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BanExecCommand) CheckConformance(traversal *NodeTraversal, node rhino.Node) *ConformanceRules_ConformanceResult {
	call, ok := r.asBannedExecCommandCall(node)
	if !ok {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	firstArg, hasFirst := call.FirstArg()
	if !hasFirst {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	attr, ok := ast.AsStringLit(firstArg.Node())
	if !ok {
		return ConformanceRules_ConformanceResult_VIOLATION
	}

	if slices.Contains(r.bannedAttrs, jutil.StringToLowerCaseRoot(attr.Value())) {
		return ConformanceRules_ConformanceResult_VIOLATION
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// Java: ConformanceRules.BanExecCommand.isBannedProperty(Node)
func (r *ConformanceRules_BanExecCommand) asBannedExecCommandCall(node rhino.Node) (ast.Call, bool) {
	call, ok := ast.AsCall(node)
	if !ok {
		return ast.Call{}, false
	}
	target, ok := ast.AsGetProp(call.Callee().Node())
	if !ok || target.Property() != "execCommand" {
		return ast.Call{}, false
	}
	type_ := target.Object().Node().GetJSType()
	if type_ == nil {
		return ast.Call{}, false
	}
	documentType := r.compiler.GetTypeRegistry().GetGlobalType("Document")
	if documentType == nil || !type_.RestrictByNotNullOrUndefined().IsSubtypeOf(documentType) {
		return ast.Call{}, false
	}
	return call, true
}

// ConformanceRules_BanStaticThis checks that this is not being referenced
// directly within a static member function.
type ConformanceRules_BanStaticThis struct {
	ConformanceRules_AbstractRule
}

// Java: ConformanceRules.BanStaticThis(AbstractCompiler, Requirement)
func newConformanceRules_BanStaticThis(compiler AbstractCompiler, requirement *Requirement) (*ConformanceRules_BanStaticThis, error) {
	r := &ConformanceRules_BanStaticThis{}
	if err := r.initTypeRestrictionRule(compiler, requirement, false); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *ConformanceRules_BanStaticThis) Check(
	t *NodeTraversal, n rhino.Node, behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) {
	if result := r.CheckConformance(t, n); result.level != ConformanceRules_ConformanceLevel_CONFORMANCE {
		r.Report(n, result, behavior)
	}
}

// Java: ConformanceRules.BanStaticThis.checkConformance(NodeTraversal, Node)
func (r *ConformanceRules_BanStaticThis) CheckConformance(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if !n.IsThis() {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	enclosingFunction := NodeUtil_GetEnclosingNonArrowFunction(n)
	if !enclosingFunction.IsNil() && ConformanceRules_BanStaticThis_IsStaticMethod(enclosingFunction) {
		return ConformanceRules_ConformanceResult_VIOLATION
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// Java: ConformanceRules.BanStaticThis.isStaticMethod(Node)
func ConformanceRules_BanStaticThis_IsStaticMethod(n rhino.Node) bool {
	jutil.CheckArgument(n.IsFunction())
	mfd, ok := ast.AsMemberFunctionDef(n.GetParent())
	return ok && mfd.IsStatic()
}
