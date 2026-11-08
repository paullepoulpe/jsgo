// Port of com/google/javascript/jscomp/JSDocInfoPrinter.java @ closure-compiler v20261104

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
	"regexp"
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// JSDocInfoPrinter prints a JSDocInfo, used for preserving type annotations
// in ES6 transpilation.
type JSDocInfoPrinter struct {
	useOriginalName bool
	printDesc       bool
}

// Java: JSDocInfoPrinter(boolean)
func NewJSDocInfoPrinter(useOriginalName bool) *JSDocInfoPrinter {
	return NewJSDocInfoPrinterWithPrintDesc(useOriginalName, false)
}

// NewJSDocInfoPrinterWithPrintDesc: useOriginalName is whether to use the
// original name field when printing types; printDesc is whether to print
// block, param, and return descriptions.
//
// Java: JSDocInfoPrinter(boolean, boolean)
func NewJSDocInfoPrinterWithPrintDesc(useOriginalName bool, printDesc bool) *JSDocInfoPrinter {
	return &JSDocInfoPrinter{useOriginalName: useOriginalName, printDesc: printDesc}
}

// The regular expression of print(JSDocInfo). Java's \s is
// [ \t\n\x0B\f\r]; Go's \s lacks \x0B, so the class is spelled out. The
// pattern uses neither backreferences nor lookaround, and Go's leftmost-first
// matching picks the same matches as Java's backtracking for it.
var blockDescriptionLineStartPattern = regexp.MustCompile("\n[ \t\n\x0B\f\r]*\\*[ \t\n\x0B\f\r]*")

// jsDocWriter streams JSDoc parts directly into a single strings.Builder while
// applying the line-prefix formatting (" *\n" / " * ") on the fly.
type jsDocWriter struct {
	sb             strings.Builder
	multiline      bool
	afterBlockDesc bool
	pendingNL      bool
}

func (w *jsDocWriter) startPart() {
	if w.afterBlockDesc {
		w.writeByte('\n')
		w.afterBlockDesc = false
	}
	if w.multiline {
		w.writeByte('\n')
	} else {
		w.writeByte(' ')
	}
}

func (w *jsDocWriter) writeByte(b byte) {
	if b == '\n' {
		if w.pendingNL {
			w.sb.WriteString("\n *")
		}
		w.pendingNL = true
		return
	}
	if w.pendingNL {
		if b != ' ' {
			w.sb.WriteString("\n * ")
		} else {
			w.sb.WriteString("\n *")
		}
		w.pendingNL = false
	}
	w.sb.WriteByte(b)
}

func (w *jsDocWriter) writeString(s string) {
	if len(s) == 0 {
		return
	}
	if w.pendingNL {
		if s[0] != ' ' && s[0] != '\n' {
			w.sb.WriteString("\n * ")
		} else {
			w.sb.WriteString("\n *")
		}
		w.pendingNL = false
	}
	for {
		nl := strings.IndexByte(s, '\n')
		if nl < 0 {
			w.sb.WriteString(s)
			return
		}
		w.sb.WriteString(s[:nl])
		if nl+1 == len(s) {
			w.pendingNL = true
			return
		}
		if s[nl+1] != ' ' && s[nl+1] != '\n' {
			w.sb.WriteString("\n * ")
		} else {
			w.sb.WriteString("\n *")
		}
		s = s[nl+1:]
	}
}

func (w *jsDocWriter) finish() string {
	if w.pendingNL {
		w.sb.WriteString("\n *")
		w.pendingNL = false
	}
	if w.multiline {
		w.sb.WriteString("\n */\n")
	} else {
		w.sb.WriteString(" */ ")
	}
	return w.sb.String()
}

func jsDocTrimRightSpace(s string) string {
	end := len(s)
	for end > 0 && s[end-1] <= ' ' {
		end--
	}
	return s[:end]
}

// Java: JSDocInfoPrinter.print(JSDocInfo)
func (p *JSDocInfoPrinter) Print(info *rhino.JSDocInfo) string {
	authors := info.GetAuthors()
	description, hasDescription := info.GetDescriptionString()
	refs := info.GetReferences()
	extendedInterfaces := info.GetExtendedInterfaces()
	implementedInterfaces := info.GetImplementedInterfaces()
	paramCount := info.GetParameterCount()
	throws := info.GetThrowsAnnotations()
	hasThrows := len(throws) != 0 && throws[0] != ""
	templateEntries := info.GetTemplateTypeEntries()
	typeTransformations := info.GetTypeTransformationEntries()
	suppressions := info.GetSuppressionsAndTheirDescription()

	var cleanedBlockDesc string
	hasBlockDesc := false
	if p.printDesc {
		if blockDesc, ok := info.GetBlockDescriptionString(); ok {
			cleaned := blockDescriptionLineStartPattern.ReplaceAllLiteralString(blockDesc, "\n")
			if cleaned != "" {
				cleanedBlockDesc = jutil.StringTrim(cleaned)
				hasBlockDesc = true
			}
		}
	}

	multiline := hasBlockDesc ||
		authors != nil ||
		hasDescription ||
		refs != nil ||
		info.HasBaseType() ||
		len(extendedInterfaces) > 0 ||
		len(implementedInterfaces) > 0 ||
		info.HasThisType() ||
		paramCount > 0 ||
		info.HasReturnType() ||
		hasThrows ||
		len(templateEntries) > 0 ||
		len(typeTransformations) > 0 ||
		len(suppressions) != 0 ||
		info.IsDeprecated() ||
		info.IsPolymer() ||
		info.IsPolymerBehavior() ||
		info.IsMixinFunction() ||
		info.IsMixinClass() ||
		info.IsCustomElement()

	var w jsDocWriter
	w.multiline = multiline
	if multiline {
		w.sb.Grow(96)
	} else {
		w.sb.Grow(32)
	}
	w.sb.WriteString("/**")

	if hasBlockDesc {
		w.startPart()
		w.writeString(cleanedBlockDesc)
		w.afterBlockDesc = true
	}

	if info.IsExterns() {
		w.startPart()
		w.writeString("@externs")
	}
	if info.IsTypeSummary() {
		w.startPart()
		w.writeString("@typeSummary")
	}

	if info.IsExport() {
		w.startPart()
		w.writeString("@export")
	} else if vis := info.GetVisibility(); vis != rhino.JSDocInfo_Visibility_INHERITED {
		// Go-only: language-adapter — getVisibility() is never null (it
		// defaults to INHERITED), so `getVisibility() != null` always holds.
		w.startPart()
		switch vis {
		case rhino.JSDocInfo_Visibility_PRIVATE:
			w.writeString("@private")
		case rhino.JSDocInfo_Visibility_PROTECTED:
			w.writeString("@protected")
		case rhino.JSDocInfo_Visibility_PACKAGE:
			w.writeString("@package")
		case rhino.JSDocInfo_Visibility_PUBLIC:
			w.writeString("@public")
		default:
			w.writeByte('@')
			w.writeString(jutil.StringToLowerCaseASCII(vis.String()))
		}
	}

	for _, name := range authors {
		w.startPart()
		w.writeString("@author ")
		w.writeString(name)
	}

	if info.IsAbstract() {
		w.startPart()
		w.writeString("@abstract")
	}

	if info.HasLendsName() {
		w.startPart()
		p.writeAnnotationWithTypeNode(&w, "lends", info.GetLendsName().GetRoot(), "", false)
	}

	if info.HasConstAnnotation() && !info.IsDefine() {
		w.startPart()
		w.writeString("@const")
	}

	if info.IsFinal() {
		w.startPart()
		w.writeString("@final")
	}

	if hasDescription {
		w.startPart()
		w.writeString("@desc ")
		w.writeString(description)
	}

	for _, desc := range refs {
		w.startPart()
		w.writeString("@see ")
		w.writeString(desc)
	}

	if info.IsWizaction() {
		w.startPart()
		w.writeString("@wizaction")
	}

	if info.IsNoSideEffects() {
		w.startPart()
		w.writeString("@nosideeffects")
	}

	if info.IsNoCompile() {
		w.startPart()
		w.writeString("@nocompile")
	}

	if info.IsNoInline() {
		w.startPart()
		w.writeString("@noinline")
	}

	if info.IsRequireInlining() {
		w.startPart()
		w.writeString("@requireInlining")
	}

	if info.IsEncourageInlining() {
		w.startPart()
		w.writeString("@encourageInlining")
	}

	if info.IsProvideAlreadyProvided() {
		w.startPart()
		w.writeString("@provideAlreadyProvided")
	}

	if info.IsIdGenerator() {
		w.startPart()
		w.writeString("@idGenerator {unique}")
	}

	if info.IsConsistentIdGenerator() {
		w.startPart()
		w.writeString("@idGenerator {consistent}")
	}

	if info.IsStableIdGenerator() {
		w.startPart()
		w.writeString("@idGenerator {stable}")
	}

	if info.IsXidGenerator() {
		w.startPart()
		w.writeString("@idGenerator {xid}")
	}

	if info.IsMappedIdGenerator() {
		w.startPart()
		w.writeString("@idGenerator {mapped}")
	}

	if info.MakesDicts() {
		w.startPart()
		w.writeString("@dict")
	}

	if info.MakesStructs() {
		w.startPart()
		w.writeString("@struct")
	}

	if info.MakesUnrestricted() {
		w.startPart()
		w.writeString("@unrestricted ")
	}

	if info.IsConstructor() {
		w.startPart()
		w.writeString("@constructor")
	}

	if info.IsInterface() && !info.UsesImplicitMatch() {
		w.startPart()
		w.writeString("@interface")
	}

	if info.IsInterface() && info.UsesImplicitMatch() {
		w.startPart()
		w.writeString("@record")
	}

	if info.IsClosureUnawareCode() {
		w.startPart()
		w.writeString("@closureUnaware")
	}

	if info.HasBaseType() {
		w.startPart()
		p.writeAnnotationWithTypeNode(&w, "extends", jsDocStripBang(info.GetBaseType().GetRoot()), "", false)
	}

	for _, type_ := range extendedInterfaces {
		w.startPart()
		p.writeAnnotationWithTypeNode(&w, "extends", jsDocStripBang(type_.GetRoot()), "", false)
	}

	for _, type_ := range implementedInterfaces {
		w.startPart()
		p.writeAnnotationWithTypeNode(&w, "implements", jsDocStripBang(type_.GetRoot()), "", false)
	}

	if info.HasThisType() {
		w.startPart()
		p.writeAnnotationWithTypeNode(&w, "this", jsDocStripBang(info.GetThisType().GetRoot()), "", false)
	}

	for i := 0; i < paramCount; i++ {
		name, _ := info.GetParameterNameValueAt(i)
		w.startPart()
		w.writeString("@param ")
		p.writeParamType(&w, info, name)
	}

	if info.HasReturnType() {
		retDesc, hasRetDesc := info.GetReturnDescriptionString()
		w.startPart()
		p.writeAnnotationWithTypeNode(&w, "return", info.GetReturnType().GetRoot(), retDesc, hasRetDesc)
	}

	if hasThrows {
		w.startPart()
		w.writeString("@throws ")
		w.writeString(throws[0])
	}

	for _, entry := range templateEntries {
		name, boundExpr := entry.Name, entry.Type
		boundRoot := boundExpr.GetRoot()
		w.startPart()
		if boundRoot.GetToken() == rhino.Token_QMARK && !boundRoot.HasChildren() {
			// If the bound of the expression is `?` (as it is after
			// parsing an unbounded template) don't specify a bound.
			// TODO(b/140187077): This case becomes redundant when fixed.
			// It also only covers explicit `?` bounds, typedefs will
			// remain explicit.
			w.writeString("@template ")
			w.writeString(name)
		} else {
			p.writeAnnotationWithTypeNode(&w, "template", boundRoot, name, true)
		}
	}

	for _, tr := range typeTransformations {
		tranformationDefinition := PrintCode(CodePrintConfig{Root: tr.TypeNode})
		w.startPart()
		w.writeString("@template ")
		w.writeString(tr.Name)
		w.writeString(" := ")
		w.writeString(tranformationDefinition)
		w.writeString(" =:")
	}

	if info.IsOverride() {
		w.startPart()
		w.writeString("@override")
	}

	if info.HasType() && !info.IsDefine() {
		w.startPart()
		if info.IsInlineType() {
			p.appendTypeNode(&w, info.GetType().GetRoot())
		} else {
			p.writeAnnotationWithTypeNode(&w, "type", info.GetType().GetRoot(), "", false)
		}
	}

	if info.IsDefine() {
		w.startPart()
		p.writeAnnotationWithTypeNode(&w, "define", info.GetType().GetRoot(), "", false)
	}

	if info.HasTypedefType() {
		w.startPart()
		p.writeAnnotationWithTypeNode(&w, "typedef", info.GetTypedefType().GetRoot(), "", false)
	}

	if info.HasEnumParameterType() {
		w.startPart()
		p.writeAnnotationWithTypeNode(&w, "enum", info.GetEnumParameterType().GetRoot(), "", false)
	}

	if info.IsImplicitCast() {
		w.startPart()
		w.writeString("@implicitCast")
	}

	if info.IsNoCollapse() {
		w.startPart()
		w.writeString("@nocollapse")
	}

	var suppressBuf [8]string
	for _, sup := range suppressions {
		var warnings []string
		if len(sup.Warnings) <= len(suppressBuf) {
			warnings = suppressBuf[:len(sup.Warnings)]
			copy(warnings, sup.Warnings)
		} else {
			warnings = slices.Clone(sup.Warnings)
		}
		// Even the warnings inside a suppress annotation are printed in
		// natural order for consistency
		slices.SortStableFunc(warnings, jutil.StringCompareTo)
		w.startPart()
		w.writeString("@suppress {")
		for i, warn := range warnings {
			if i > 0 {
				w.writeByte(',')
			}
			w.writeString(warn)
		}
		w.writeByte('}')
		if sup.Description != "" {
			w.writeByte(' ')
			w.writeString(sup.Description)
		}
	}

	if info.IsDeprecated() {
		w.startPart()
		if deprecationReason, ok := info.GetDeprecationReasonString(); ok {
			w.writeString("@deprecated ")
			w.writeString(deprecationReason)
		} else {
			w.writeString("@deprecated")
		}
	}

	if info.IsPolymer() {
		w.startPart()
		w.writeString("@polymer")
	}
	if info.IsPolymerBehavior() {
		w.startPart()
		w.writeString("@polymerBehavior")
	}
	if info.IsMixinFunction() {
		w.startPart()
		w.writeString("@mixinFunction")
	}
	if info.IsMixinClass() {
		w.startPart()
		w.writeString("@mixinClass")
	}
	if info.IsCustomElement() {
		w.startPart()
		w.writeString("@customElement")
	}

	if id, ok := info.GetClosurePrimitiveIdString(); ok {
		w.startPart()
		w.writeString("@closurePrimitive {")
		w.writeString(id)
		w.writeByte('}')
	}

	if info.IsNgInject() {
		w.startPart()
		w.writeString("@ngInject")
	}

	if info.IsLegacyNamedTsExport() {
		w.startPart()
		w.writeString("@legacyNamedTsExport")
	}

	for _, tsType := range info.GetTsTypes() {
		w.startPart()
		w.writeString("@tsType ")
		w.writeString(tsType)
	}

	return w.finish()
}

func jsDocOnlyChild(n rhino.Node) (rhino.Node, bool) {
	for c := range n.Children() {
		return c, true
	}
	return rhino.NilNode, false
}

func jsDocColonChildren(n rhino.Node) (string, rhino.Node) {
	var key string
	var val rhino.Node
	first := true
	for c := range n.Children() {
		if first {
			if sk, ok := ast.AsStringKey(c); ok {
				key = sk.Key()
			} else {
				key = c.GetString()
			}
			first = false
		} else {
			val = c
		}
	}
	return key, val
}

// Java: JSDocInfoPrinter.stripBang(Node)
func jsDocStripBang(typeNode rhino.Node) rhino.Node {
	if typeNode.GetToken() == rhino.Token_BANG {
		child, _ := jsDocOnlyChild(typeNode)
		return child
	}
	return typeNode
}

func (p *JSDocInfoPrinter) writeAnnotationWithTypeNode(
	w *jsDocWriter, annotation string, type_ rhino.Node, description string, hasDesc bool,
) {
	w.writeByte('@')
	w.writeString(annotation)
	w.writeString(" {")
	p.appendTypeNode(w, type_)
	w.writeByte('}')
	if hasDesc {
		w.writeByte(' ')
		w.writeString(description)
	}
}

// Java: JSDocInfoPrinter.buildParamType(JSDocInfo, String)
func (p *JSDocInfoPrinter) writeParamType(w *jsDocWriter, info *rhino.JSDocInfo, name string) {
	type_ := info.GetParameterType(name)
	if type_ == nil {
		w.writeString(name)
		return
	}
	desc := ""
	if p.printDesc {
		if d, ok := info.GetDescriptionForParameterString(name); ok {
			// Don't add a leading space; the parser retained it.
			desc = d
		}
	}
	w.writeByte('{')
	p.appendTypeNode(w, type_.GetRoot())
	trimmedDesc := jsDocTrimRightSpace(desc)
	if trimmedDesc != "" {
		w.writeString("} ")
		w.writeString(name)
		w.writeString(trimmedDesc)
		return
	}
	trimmedName := jsDocTrimRightSpace(name)
	if trimmedName != "" {
		w.writeString("} ")
		w.writeString(trimmedName)
		return
	}
	w.writeByte('}')
}

// Java: JSDocInfoPrinter.appendTypeNode(StringBuilder, Node)
func (p *JSDocInfoPrinter) appendTypeNode(w *jsDocWriter, typeNode rhino.Node) {
	if p.useOriginalName {
		if orig, ok := typeNode.GetOriginalNameString(); ok {
			w.writeString(orig)
			return
		}
	}
	switch typeNode.GetToken() {
	case rhino.Token_BANG:
		w.writeByte('!')
		child, _ := jsDocOnlyChild(typeNode)
		p.appendTypeNode(w, child)
	case rhino.Token_EQUALS:
		child, _ := jsDocOnlyChild(typeNode)
		p.appendTypeNode(w, child)
		w.writeByte('=')
	case rhino.Token_PIPE:
		w.writeByte('(')
		first := true
		for child := range typeNode.Children() {
			if !first {
				w.writeByte('|')
			}
			first = false
			p.appendTypeNode(w, child)
		}
		w.writeByte(')')
	case rhino.Token_ITER_REST:
		w.writeString("...")
		if target := ast.MustIterRest(typeNode).Target().Node(); !target.IsNil() && !target.IsEmpty() {
			p.appendTypeNode(w, target)
		}
	case rhino.Token_STAR:
		w.writeByte('*')
	case rhino.Token_QMARK:
		w.writeByte('?')
		if child, ok := jsDocOnlyChild(typeNode); ok {
			p.appendTypeNode(w, child)
		}
	case rhino.Token_FUNCTION:
		p.appendFunctionNode(w, typeNode)
	case rhino.Token_LC:
		w.writeByte('{')
		if lb, ok := jsDocOnlyChild(typeNode); ok {
			first := true
			for field := range lb.Children() {
				if !first {
					w.writeByte(',')
				}
				first = false
				if sk, isSk := ast.AsStringKey(field); isSk && !field.HasChildren() {
					w.writeString(sk.Key())
				} else if field.HasChildren() {
					key, val := jsDocColonChildren(field)
					w.writeString(key)
					w.writeByte(':')
					p.appendTypeNode(w, val)
				} else {
					w.writeString(field.GetString())
				}
			}
		}
		w.writeByte('}')
	case rhino.Token_VOID:
		w.writeString("void")
	case rhino.Token_TYPEOF:
		w.writeString("typeof ")
		p.appendTypeNode(w, ast.MustUnary(typeNode).Operand().Node())
	case rhino.Token_BLOCK:
		w.writeByte('<')
		first := true
		for type_ := range typeNode.Children() {
			if !first {
				w.writeByte(',')
			}
			first = false
			p.appendTypeNode(w, type_)
		}
		w.writeByte('>')
	case rhino.Token_STRINGLIT:
		w.writeString(ast.MustStringLit(typeNode).Value())
		if child, ok := jsDocOnlyChild(typeNode); ok {
			p.appendTypeNode(w, child)
		}
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Unexpected typeNode: " + typeNode.String()))
	}
}

// Java: JSDocInfoPrinter.appendFunctionNode(StringBuilder, Node)
func (p *JSDocInfoPrinter) appendFunctionNode(w *jsDocWriter, function rhino.Node) {
	hasNewOrThis := false
	seenFirst := false
	closedParen := false
	w.writeString("function(")
	for child := range function.Children() {
		if !seenFirst {
			seenFirst = true
			if newNode, ok := ast.AsNew(child); ok {
				w.writeString("new:")
				p.appendTypeNode(w, newNode.Callee().Node())
				hasNewOrThis = true
				continue
			}
			if child.IsThis() {
				w.writeString("this:")
				ctxType, _ := jsDocOnlyChild(child)
				p.appendTypeNode(w, ctxType)
				hasNewOrThis = true
				continue
			}
		}
		if !closedParen {
			if pl, ok := ast.AsParamList(child); ok {
				firstParam := true
				for param := range pl.Params() {
					if !firstParam || hasNewOrThis {
						w.writeByte(',')
					}
					p.appendTypeNode(w, param.Node())
					firstParam = false
				}
				w.writeByte(')')
				closedParen = true
				continue
			}
			w.writeByte(')')
			closedParen = true
		}
		if !child.IsEmpty() {
			w.writeByte(':')
			p.appendTypeNode(w, child)
		}
		return
	}
	if !closedParen {
		w.writeByte(')')
	}
}
