// Port of com/google/javascript/jscomp/JSDocInfoPrinter.java @ closure-compiler v20261015

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

// The regular expressions of print(JSDocInfo). Java's \s is
// [ \t\n\x0B\f\r]; Go's \s lacks \x0B, so the class is spelled out. Neither
// pattern uses backreferences or lookaround, and Go's leftmost-first
// matching picks the same matches as Java's backtracking for them.
var (
	jsDocInfoPrinter_BLOCK_DESCRIPTION_LINE_START = regexp.MustCompile("\n[ \t\n\x0B\f\r]*\\*[ \t\n\x0B\f\r]*")
	jsDocInfoPrinter_LINE_START_WITHOUT_SPACE     = regexp.MustCompile("\n \\*([^ \n])")
)

// Java: JSDocInfoPrinter.print(JSDocInfo)
func (p *JSDocInfoPrinter) Print(info *rhino.JSDocInfo) string {
	multiline := false

	var parts []string

	parts = append(parts, "/**")

	if info.IsExterns() {
		parts = append(parts, "@externs")
	}
	if info.IsTypeSummary() {
		parts = append(parts, "@typeSummary")
	}

	if info.IsExport() {
		parts = append(parts, "@export")
	} else if info.GetVisibility() != rhino.JSDocInfo_Visibility_INHERITED {
		// Go-only: language-adapter — getVisibility() is never null (it
		// defaults to INHERITED), so `getVisibility() != null` always holds.
		parts = append(parts, "@"+jutil.StringToLowerCaseASCII(info.GetVisibility().String()))
	}

	if info.GetAuthors() != nil {
		multiline = true
		for name := range info.GetAuthors().All() {
			parts = append(parts, "@author "+name)
		}
	}

	if info.IsAbstract() {
		parts = append(parts, "@abstract")
	}

	if info.HasLendsName() {
		parts = append(parts, p.BuildAnnotationWithTypeNode("lends", info.GetLendsName().GetRoot()))
	}

	if info.HasConstAnnotation() && !info.IsDefine() {
		parts = append(parts, "@const")
	}

	if info.IsFinal() {
		parts = append(parts, "@final")
	}

	description := info.GetDescription()
	if description != nil {
		multiline = true
		parts = append(parts, "@desc "+*description)
	}

	if info.GetReferences() != nil {
		multiline = true
		for desc := range info.GetReferences().All() {
			parts = append(parts, "@see "+desc)
		}
	}

	if info.IsWizaction() {
		parts = append(parts, "@wizaction")
	}

	if info.IsNoSideEffects() {
		parts = append(parts, "@nosideeffects")
	}

	if info.IsNoCompile() {
		parts = append(parts, "@nocompile")
	}

	if info.IsNoInline() {
		parts = append(parts, "@noinline")
	}

	if info.IsRequireInlining() {
		parts = append(parts, "@requireInlining")
	}

	if info.IsEncourageInlining() {
		parts = append(parts, "@encourageInlining")
	}

	if info.IsProvideAlreadyProvided() {
		parts = append(parts, "@provideAlreadyProvided")
	}

	if info.IsIdGenerator() {
		parts = append(parts, "@idGenerator {unique}")
	}

	if info.IsConsistentIdGenerator() {
		parts = append(parts, "@idGenerator {consistent}")
	}

	if info.IsStableIdGenerator() {
		parts = append(parts, "@idGenerator {stable}")
	}

	if info.IsXidGenerator() {
		parts = append(parts, "@idGenerator {xid}")
	}

	if info.IsMappedIdGenerator() {
		parts = append(parts, "@idGenerator {mapped}")
	}

	if info.MakesDicts() {
		parts = append(parts, "@dict")
	}

	if info.MakesStructs() {
		parts = append(parts, "@struct")
	}

	if info.MakesUnrestricted() {
		parts = append(parts, "@unrestricted ")
	}

	if info.IsConstructor() {
		parts = append(parts, "@constructor")
	}

	if info.IsInterface() && !info.UsesImplicitMatch() {
		parts = append(parts, "@interface")
	}

	if info.IsInterface() && info.UsesImplicitMatch() {
		parts = append(parts, "@record")
	}

	if info.IsClosureUnawareCode() {
		parts = append(parts, "@closureUnaware")
	}

	if info.HasBaseType() {
		multiline = true
		typeNode := p.StripBang(info.GetBaseType().GetRoot())
		parts = append(parts, p.BuildAnnotationWithTypeNode("extends", typeNode))
	}

	for _, type_ := range info.GetExtendedInterfaces() {
		multiline = true
		typeNode := p.StripBang(type_.GetRoot())
		parts = append(parts, p.BuildAnnotationWithTypeNode("extends", typeNode))
	}

	for _, type_ := range info.GetImplementedInterfaces() {
		multiline = true
		typeNode := p.StripBang(type_.GetRoot())
		parts = append(parts, p.BuildAnnotationWithTypeNode("implements", typeNode))
	}

	if info.HasThisType() {
		multiline = true
		typeNode := p.StripBang(info.GetThisType().GetRoot())
		parts = append(parts, p.BuildAnnotationWithTypeNode("this", typeNode))
	}

	if paramCount := info.GetParameterCount(); paramCount > 0 {
		multiline = true
		for i := 0; i < paramCount; i++ {
			name, _ := info.GetParameterNameValueAt(i)
			parts = append(parts, "@param "+p.BuildParamType(info, name))
		}
	}

	if info.HasReturnType() {
		multiline = true
		parts = append(parts,
			p.BuildAnnotationWithTypeAndDescription("return", info.GetReturnType(), info.GetReturnDescription()))
	}

	if len(info.GetThrowsAnnotations()) != 0 && info.GetThrowsAnnotations()[0] != "" {
		multiline = true
		parts = append(parts, "@throws "+info.GetThrowsAnnotations()[0])
	}

	templates := info.GetTemplateTypes()
	if !templates.IsEmpty() {
		multiline = true
		for name, boundExpr := range templates.All() {
			boundRoot := boundExpr.GetRoot()
			if boundRoot.GetToken() == rhino.Token_QMARK && !boundRoot.HasChildren() {
				// If the bound of the expression is `?` (as it is after
				// parsing an unbounded template) don't specify a bound.
				// TODO(b/140187077): This case becomes redundant when fixed.
				// It also only covers explicit `?` bounds, typedefs will
				// remain explicit.
				parts = append(parts, "@template "+name)
			} else {
				parts = append(parts, p.BuildAnnotationWithTypeAndDescription("template", boundExpr, &name))
			}
		}
	}

	typeTransformations := info.GetTypeTransformations()
	if !typeTransformations.IsEmpty() {
		multiline = true
		for name, value := range typeTransformations.All() {
			tranformationDefinition := NewCodePrinter_Builder(value).Build()
			parts = append(parts, "@template "+name+" := "+tranformationDefinition+" =:")
		}
	}

	if info.IsOverride() {
		parts = append(parts, "@override")
	}

	if info.HasType() && !info.IsDefine() {
		if info.IsInlineType() {
			parts = append(parts, p.TypeNode(info.GetType().GetRoot()))
		} else {
			parts = append(parts, p.BuildAnnotationWithType("type", info.GetType()))
		}
	}

	if info.IsDefine() {
		parts = append(parts, p.BuildAnnotationWithType("define", info.GetType()))
	}

	if info.HasTypedefType() {
		parts = append(parts, p.BuildAnnotationWithType("typedef", info.GetTypedefType()))
	}

	if info.HasEnumParameterType() {
		parts = append(parts, p.BuildAnnotationWithType("enum", info.GetEnumParameterType()))
	}

	if info.IsImplicitCast() {
		parts = append(parts, "@implicitCast")
	}

	if info.IsNoCollapse() {
		parts = append(parts, "@nocollapse")
	}

	suppressions := info.GetSuppressionsAndTheirDescription()
	if !suppressions.IsEmpty() {
		// With ImmutableMap, the iteration order will be same as insertion
		// order (i.e. parse order).
		for key, text := range suppressions.All() {
			warnings := key.ToSlice()
			// Even the warnings inside a suppress annotation are printed in
			// natural order for consistency
			slices.SortStableFunc(warnings, jutil.StringCompareTo)
			var sb strings.Builder
			sb.WriteString("@suppress {")
			sb.WriteString(strings.Join(warnings, ","))
			sb.WriteString("}")
			if text != "" {
				sb.WriteString(" ")
				sb.WriteString(text)
			}
			parts = append(parts, sb.String())
		}
		multiline = true
	}

	if info.IsDeprecated() {
		deprecationReason := info.GetDeprecationReason()
		if deprecationReason != nil {
			parts = append(parts, "@deprecated "+*deprecationReason)
		} else {
			parts = append(parts, "@deprecated")
		}
		multiline = true
	}

	if info.IsPolymer() {
		multiline = true
		parts = append(parts, "@polymer")
	}
	if info.IsPolymerBehavior() {
		multiline = true
		parts = append(parts, "@polymerBehavior")
	}
	if info.IsMixinFunction() {
		multiline = true
		parts = append(parts, "@mixinFunction")
	}
	if info.IsMixinClass() {
		multiline = true
		parts = append(parts, "@mixinClass")
	}
	if info.IsCustomElement() {
		multiline = true
		parts = append(parts, "@customElement")
	}

	if info.GetClosurePrimitiveId() != nil {
		parts = append(parts, "@closurePrimitive {"+*info.GetClosurePrimitiveId()+"}")
	}

	if info.IsNgInject() {
		parts = append(parts, "@ngInject")
	}

	for _, tsType := range info.GetTsTypes() {
		parts = append(parts, "@tsType "+tsType)
	}

	if p.printDesc && info.GetBlockDescription() != nil {
		cleaned := jsDocInfoPrinter_BLOCK_DESCRIPTION_LINE_START.ReplaceAllLiteralString(*info.GetBlockDescription(), "\n")
		if cleaned != "" {
			multiline = true
			cleaned = jutil.StringTrim(cleaned)
			if len(parts) > 1 {
				// If there is more than one part - the opening "/**" - then
				// add blank line between the description and everything else.
				cleaned += "\n"
			}
			parts = slices.Insert(parts, 1, cleaned)
		}
	}

	var sb strings.Builder
	if multiline {
		sb.WriteString(strings.Join(parts, "\n"))
	} else {
		sb.WriteString(strings.Join(parts, " "))
		sb.WriteString(" */")
	}
	// Ensure all lines start with " *", and then ensure all non blank lines
	// have a space after the *.
	s := jsDocInfoPrinter_LINE_START_WITHOUT_SPACE.ReplaceAllString(
		strings.ReplaceAll(sb.String(), "\n", "\n *"), "\n * $1")
	if multiline {
		s += "\n */\n"
	} else {
		s += " "
	}
	return s
}

// Java: JSDocInfoPrinter.stripBang(Node)
func (p *JSDocInfoPrinter) StripBang(typeNode rhino.Node) rhino.Node {
	if typeNode.GetToken() == rhino.Token_BANG {
		typeNode = typeNode.GetFirstChild()
	}
	return typeNode
}

// Java: JSDocInfoPrinter.buildAnnotationWithType(String, JSTypeExpression)
func (p *JSDocInfoPrinter) BuildAnnotationWithType(annotation string, type_ *rhino.JSTypeExpression) string {
	return p.BuildAnnotationWithTypeAndDescription(annotation, type_, nil)
}

// Java: JSDocInfoPrinter.buildAnnotationWithType(String, JSTypeExpression, String)
func (p *JSDocInfoPrinter) BuildAnnotationWithTypeAndDescription(
	annotation string, type_ *rhino.JSTypeExpression, description *string) string {
	return p.BuildAnnotationWithTypeNodeAndDescription(annotation, type_.GetRoot(), description)
}

// Java: JSDocInfoPrinter.buildAnnotationWithType(String, Node)
func (p *JSDocInfoPrinter) BuildAnnotationWithTypeNode(annotation string, type_ rhino.Node) string {
	return p.BuildAnnotationWithTypeNodeAndDescription(annotation, type_, nil)
}

// Java: JSDocInfoPrinter.buildAnnotationWithType(String, Node, String)
func (p *JSDocInfoPrinter) BuildAnnotationWithTypeNodeAndDescription(
	annotation string, type_ rhino.Node, description *string) string {
	var sb strings.Builder
	sb.WriteString("@")
	sb.WriteString(annotation)
	sb.WriteString(" {")
	p.AppendTypeNode(&sb, type_)
	sb.WriteString("}")
	if description != nil {
		sb.WriteString(" ")
		sb.WriteString(*description)
	}
	return sb.String()
}

// Java: JSDocInfoPrinter.buildParamType(JSDocInfo, String)
func (p *JSDocInfoPrinter) BuildParamType(info *rhino.JSDocInfo, name string) string {
	type_ := info.GetParameterType(name)
	if type_ != nil {
		desc := ""
		if p.printDesc && info.GetDescriptionForParameter(name) != nil {
			// Don't add a leading space; the parser retained it.
			desc = *info.GetDescriptionForParameter(name)
		}
		// Go-only: language-adapter — Java's local `p` would shadow the
		// receiver here.
		p_ := "{" + p.TypeNode(type_.GetRoot()) + "} " + name + desc
		return jutil.StringTrim(p_)
	} else {
		return name
	}
}

// Java: JSDocInfoPrinter.typeNode(Node)
func (p *JSDocInfoPrinter) TypeNode(typeNode rhino.Node) string {
	var sb strings.Builder
	p.AppendTypeNode(&sb, typeNode)
	return sb.String()
}

// Java: JSDocInfoPrinter.appendTypeNode(StringBuilder, Node)
func (p *JSDocInfoPrinter) AppendTypeNode(sb *strings.Builder, typeNode rhino.Node) {
	if p.useOriginalName && typeNode.GetOriginalName() != nil {
		sb.WriteString(*typeNode.GetOriginalName())
		return
	}
	switch typeNode.GetToken() {
	case rhino.Token_BANG:
		sb.WriteString("!")
		p.AppendTypeNode(sb, typeNode.GetFirstChild())
	case rhino.Token_EQUALS:
		p.AppendTypeNode(sb, typeNode.GetFirstChild())
		sb.WriteString("=")
	case rhino.Token_PIPE:
		sb.WriteString("(")
		lastChild := typeNode.GetLastChild()
		for child := typeNode.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
			p.AppendTypeNode(sb, child)
			if child != lastChild {
				sb.WriteString("|")
			}
		}
		sb.WriteString(")")
	case rhino.Token_ITER_REST:
		sb.WriteString("...")
		if typeNode.HasChildren() && !typeNode.GetFirstChild().IsEmpty() {
			p.AppendTypeNode(sb, typeNode.GetFirstChild())
		}
	case rhino.Token_STAR:
		sb.WriteString("*")
	case rhino.Token_QMARK:
		sb.WriteString("?")
		if typeNode.HasChildren() {
			p.AppendTypeNode(sb, typeNode.GetFirstChild())
		}
	case rhino.Token_FUNCTION:
		p.AppendFunctionNode(sb, typeNode)
	case rhino.Token_LC:
		sb.WriteString("{")
		lb := typeNode.GetFirstChild()
		lastColon := lb.GetLastChild()
		for colon := lb.GetFirstChild(); !colon.IsNil(); colon = colon.GetNext() {
			if colon.HasChildren() {
				sb.WriteString(colon.GetFirstChild().GetString())
				sb.WriteString(":")
				p.AppendTypeNode(sb, colon.GetLastChild())
			} else {
				sb.WriteString(colon.GetString())
			}
			if colon != lastColon {
				sb.WriteString(",")
			}
		}
		sb.WriteString("}")
	case rhino.Token_VOID:
		sb.WriteString("void")
	case rhino.Token_TYPEOF:
		sb.WriteString("typeof ")
		p.AppendTypeNode(sb, typeNode.GetFirstChild())
	case rhino.Token_BLOCK:
		sb.WriteString("<")
		last := typeNode.GetLastChild()
		for type_ := typeNode.GetFirstChild(); !type_.IsNil(); type_ = type_.GetNext() {
			p.AppendTypeNode(sb, type_)
			if type_ != last {
				sb.WriteString(",")
			}
		}
		sb.WriteString(">")
	case rhino.Token_STRINGLIT:
		sb.WriteString(typeNode.GetString())
		if typeNode.HasChildren() {
			p.AppendTypeNode(sb, typeNode.GetOnlyChild())
		}
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Unexpected typeNode: " + typeNode.String()))
	}
}

// Java: JSDocInfoPrinter.appendFunctionNode(StringBuilder, Node)
func (p *JSDocInfoPrinter) AppendFunctionNode(sb *strings.Builder, function rhino.Node) {
	hasNewOrThis := false
	sb.WriteString("function(")
	first := function.GetFirstChild()
	if first.IsNew() {
		sb.WriteString("new:")
		p.AppendTypeNode(sb, first.GetFirstChild())
		hasNewOrThis = true
	} else if first.IsThis() {
		sb.WriteString("this:")
		p.AppendTypeNode(sb, first.GetFirstChild())
		hasNewOrThis = true
	} else if first.IsEmpty() {
		sb.WriteString(")")
		return
	} else if !first.IsParamList() {
		sb.WriteString("):")
		p.AppendTypeNode(sb, first)
		return
	}
	var paramList rhino.Node
	if first.IsParamList() {
		paramList = first
	} else if first.GetNext().IsParamList() {
		paramList = first.GetNext()
	}
	if !paramList.IsNil() {
		firstParam := true
		for param := paramList.GetFirstChild(); !param.IsNil(); param = param.GetNext() {
			if !firstParam || hasNewOrThis {
				sb.WriteString(",")
			}
			p.AppendTypeNode(sb, param)
			firstParam = false
		}
	}
	sb.WriteString(")")

	returnType := function.GetLastChild()
	if !returnType.IsEmpty() {
		sb.WriteString(":")
		p.AppendTypeNode(sb, returnType)
	}
}
