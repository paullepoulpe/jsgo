// Port of com/google/javascript/refactoring/ErrorToFixMapper.java @ closure-compiler v20261021

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

package refactoring

import (
	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ErrorToFixMapper maps a JSError to a list of SuggestedFixes, if possible.
// TODO(tbreisacher): Move this into the compiler itself (i.e. into the
// jscomp package). This will make it easier for people adding new warnings
// to also add fixes for them.
//
// Java: public final class ErrorToFixMapper
type ErrorToFixMapper struct {
	compiler              jscomp.AbstractCompiler
	metadataForEachScript *jutil.LinkedHashMap[rhino.Node, *ScriptMetadata]
}

var errorToFixMapper_DID_YOU_MEAN = jutil.Pattern_Compile(".*Did you mean (.*)\\?")
var errorToFixMapper_EARLY_REF = jutil.Pattern_Compile("Variable referenced before declaration: (.*)")

// errorToFixMapper_MISSING_REQUIRE is
// Pattern.compile("'([^']+)' references a .*", Pattern.DOTALL).
//
// Go-only: language-adapter — jutil.Pattern has no flags: under DOTALL `.`
// matches every char, which is `.` (every char but the line terminators
// \n, \r, U+0085, U+2028 and U+2029) or one of those terminators.
var errorToFixMapper_MISSING_REQUIRE = jutil.Pattern_Compile(
	"'([^']+)' references a (?:.|[\n\r\u0085  ])*")
var errorToFixMapper_FULLY_QUALIFIED_NAME = jutil.Pattern_Compile("Reference to fully qualified import name '([^']+)'.*")
var errorToFixMapper_USE_SHORT_NAME = jutil.Pattern_Compile(".*Please use the short name '(.*)' instead.")

// Java: ErrorToFixMapper(AbstractCompiler)
func NewErrorToFixMapper(compiler jscomp.AbstractCompiler) *ErrorToFixMapper {
	m := &ErrorToFixMapper{
		metadataForEachScript: jutil.NewLinkedHashMap[rhino.Node, *ScriptMetadata](),
	}
	m.compiler = compiler
	return m
}

// Java: ErrorToFixMapper.getFixesForJsError(JSError)
func (m *ErrorToFixMapper) GetFixesForJsError(error *jscomp.JSError) []*SuggestedFix {
	fix := m.GetFixForJsError(error)
	if fix != nil {
		return []*SuggestedFix{fix}
	}
	switch error.Type().Key() {
	case "JSC_IMPLICITLY_NONNULL_JSDOC",
		"JSC_IMPLICITLY_NULLABLE_JSDOC",
		"JSC_MISSING_NULLABILITY_MODIFIER_JSDOC",
		"JSC_NULL_MISSING_NULLABILITY_MODIFIER_JSDOC":
		return m.GetFixesForImplicitNullabilityErrors(error)
	default:
		return []*SuggestedFix{}
	}
}

// GetFixForJsError creates a SuggestedFix for the given error. Note that
// some errors have multiple fixes so GetFixesForJsError should often be used
// instead of this. It may return nil.
//
// Java: ErrorToFixMapper.getFixForJsError(JSError)
func (m *ErrorToFixMapper) GetFixForJsError(error *jscomp.JSError) *SuggestedFix {
	switch error.Type().Key() {
	case "JSC_REDECLARED_VARIABLE":
		return m.GetFixForRedeclaration(error)
	case "JSC_REFERENCE_BEFORE_DECLARE":
		return m.GetFixForEarlyReference(error)
	case "JSC_MISSING_SEMICOLON":
		return m.GetFixForMissingSemicolon(error)
	case "JSC_REQUIRES_NOT_SORTED":
		return m.GetFixForUnsortedRequires(error)
	case "JSC_PROVIDES_NOT_SORTED":
		return m.GetFixForUnsortedProvides(error)
	case "JSC_DEBUGGER_STATEMENT_PRESENT":
		return NewSuggestedFix_Builder().
			AttachMatchedNodeInfo(error.Node(), m.compiler).
			SetDescription("Remove debugger statement").
			Delete(error.Node()).
			Build()
	case "JSC_USELESS_EMPTY_STATEMENT":
		return m.RemoveEmptyStatement(error)
	case "JSC_INEXISTENT_PROPERTY_WITH_SUGGESTION",
		"JSC_STRICT_INEXISTENT_PROPERTY_WITH_SUGGESTION":
		return m.GetFixForInexistentProperty(error)
	case "JSC_MISSING_CALL_TO_SUPER":
		return m.GetFixForMissingSuper(error)
	case "JSC_INVALID_SUPER_CALL_WITH_SUGGESTION":
		return m.GetFixForInvalidSuper(error)
	case "JSC_MISSING_REQUIRE", "JSC_MISSING_REQUIRE_IN_PROVIDES_FILE":
		return m.GetFixForMissingRequire(error, SuggestedFix_ImportType_REQUIRE)
	case "JSC_MISSING_REQUIRE_TYPE", "JSC_MISSING_REQUIRE_TYPE_IN_PROVIDES_FILE":
		return m.GetFixForMissingRequire(error, SuggestedFix_ImportType_REQUIRE_TYPE)
	case "JSC_EXTRA_REQUIRE_WARNING":
		return m.GetFixForExtraRequire(error)
	case "JSC_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME",
		"JSC_REFERENCE_TO_FULLY_QUALIFIED_IMPORT_NAME":
		return m.GetFixForReferenceToShortImportByLongName(error)
	case "JSC_REDUNDANT_NULLABILITY_MODIFIER_JSDOC":
		return m.GetFixForRedundantNullabilityModifierJsDoc(error)
	case "JSC_MISSING_CONST_ON_CONSTANT_CASE":
		return m.GetFixForConstantCaseErrors(error)
	default:
		return nil
	}
}

// GetFixForRedeclaration may return nil.
//
// Java: ErrorToFixMapper.getFixForRedeclaration(JSError)
func (m *ErrorToFixMapper) GetFixForRedeclaration(error *jscomp.JSError) *SuggestedFix {
	name := error.Node()
	jutil.CheckStateMsg1(name.IsName(), "%s", name)
	parent := name.GetParent()
	if !jscomp.NodeUtil_IsNameDeclaration(parent) {
		return nil
	}

	fix := NewSuggestedFix_Builder().AttachMatchedNodeInfo(name, m.compiler)

	if !name.HasChildren() {
		var nodeToDelete rhino.Node
		if parent.HasOneChild() {
			nodeToDelete = parent
		} else {
			nodeToDelete = error.Node()
		}
		return fix.SetDescription("Remove redundant declaration").Delete(nodeToDelete).Build()
	}

	fix.SetDescription("Convert redundant declaration to assignment")

	assign := rhino.IR_ExprResult(
		rhino.IR_Assign(name.CloneNode(), name.GetFirstChild().CloneTree()))
	if parent.HasOneChild() {
		return fix.Replace(parent, assign, m.compiler).Build()
	}

	// Split the var statement into an assignment and up to two var
	// statements.
	// var a = 0,
	//     b = 1,  // This is the one we're removing.
	//     c = 2;
	//
	// becomes
	//
	// var a = 0;  // This is the "added" var statement.
	// b = 1;
	// var c = 2;  // This is the original var statement.
	var childrenOfAddedVarStatement []rhino.Node
	for n := parent.GetFirstChild(); !n.IsNil(); n = n.GetNext() {
		if n == name {
			break
		}
		childrenOfAddedVarStatement = append(childrenOfAddedVarStatement, n)
	}

	if len(childrenOfAddedVarStatement) > 0 {
		var_ := rhino.NewNode(parent.GetToken())
		for _, n := range childrenOfAddedVarStatement {
			var_.AddChildToBack(n.CloneTree())
		}
		// Use a sortKey of "1" to make sure this is applied before the
		// statement below.
		fix.InsertBeforeNodeWithSortKey(parent, var_, m.compiler, "1")
	}

	if !name.GetNext().IsNil() {
		// Keep the original var statement, just remove the names that will
		// be put in the added one.
		for _, n := range childrenOfAddedVarStatement {
			fix.Delete(n)
		}
		fix.Delete(name)
		// Use a sortKey of "2" to make sure this is applied after the
		// statement above.
		fix.InsertBeforeNodeWithSortKey(parent, assign, m.compiler, "2")
	} else {
		// Remove the original var statement.
		fix.Replace(parent, assign, m.compiler)
	}

	return fix.Build()
}

// GetFixForEarlyReference: this fix is not ideal. It trades one warning
// (JSC_REFERENCE_BEFORE_DECLARE) for another (JSC_REDECLARED_VARIABLE). But
// after running the fixer once, you can then run it again and
// GetFixForRedeclaration will take care of the JSC_REDECLARED_VARIABLE
// warning. It may return nil.
//
// Java: ErrorToFixMapper.getFixForEarlyReference(JSError)
func (m *ErrorToFixMapper) GetFixForEarlyReference(error *jscomp.JSError) *SuggestedFix {
	matcher := errorToFixMapper_EARLY_REF.Matcher(error.Description())
	if matcher.Matches() {
		name := *matcher.Group(1)
		stmt := jscomp.NodeUtil_GetEnclosingStatement(error.Node())
		return NewSuggestedFix_Builder().
			AttachMatchedNodeInfo(error.Node(), m.compiler).
			SetDescription("Insert var declaration statement").
			InsertBefore(stmt, "var "+name+";\n").
			Build()
	}
	return nil
}

// Java: ErrorToFixMapper.getFixForReferenceToShortImportByLongName(JSError)
func (m *ErrorToFixMapper) GetFixForReferenceToShortImportByLongName(error *jscomp.JSError) *SuggestedFix {
	fix := NewSuggestedFix_Builder().AttachMatchedNodeInfo(error.Node(), m.compiler)
	metadata := NewNodeMetadata(m.compiler)
	match := NewMatch(error.Node(), metadata)

	script := m.compiler.GetScriptNode(*error.SourceName())
	scriptMetadata := m.GetMetadataForScript(script)

	fullNameMatcher := errorToFixMapper_FULLY_QUALIFIED_NAME.Matcher(error.Description())
	jutil.CheckStateMsg1(fullNameMatcher.Matches(), "%s", error.Description())
	fullName := *fullNameMatcher.Group(1)

	shortNameMatcher := errorToFixMapper_USE_SHORT_NAME.Matcher(error.Description())
	var shortName *string
	if shortNameMatcher.Matches() {
		shortName = shortNameMatcher.Group(1)
	} else {
		fix.AddGoogRequire(match, fullName, scriptMetadata)
		shortName = scriptMetadata.GetAlias(fullName)
	}

	var oldName string
	if error.Node().IsQualifiedName() {
		oldName = *error.Node().GetQualifiedName()
	} else {
		oldName = error.Node().GetString()
	}

	return fix.Replace(
		error.Node(),
		jscomp.NodeUtil_NewQName(m.compiler, jutil.StringReplace(oldName, fullName, *shortName)),
		m.compiler).
		Build()
}

// Java: ErrorToFixMapper.getFixesForImplicitNullabilityErrors(JSError)
func (m *ErrorToFixMapper) GetFixesForImplicitNullabilityErrors(error *jscomp.JSError) []*SuggestedFix {

	if error.Node().GetSourceFileName() == nil {
		// If we don't have a source location we can't generate a valid fix.
		return []*SuggestedFix{}
	}

	qmark :=
		NewSuggestedFix_Builder().
			AttachMatchedNodeInfo(error.Node(), m.compiler).
			InsertBefore(error.Node(), "?").
			SetDescription("Make nullability explicit").
			Build()
	bang :=
		NewSuggestedFix_Builder().
			AttachMatchedNodeInfo(error.Node(), m.compiler).
			InsertBefore(error.Node(), "!").
			SetDescription("Make type non-nullable").
			Build()
	switch error.Type().Key() {
	case "JSC_NULL_MISSING_NULLABILITY_MODIFIER_JSDOC":
		// When initializer was null, we can be confident about nullability
		return []*SuggestedFix{qmark}
	case "JSC_MISSING_NULLABILITY_MODIFIER_JSDOC":
		// Otherwise, the linter should assume ! is preferred over ?.
		return []*SuggestedFix{bang, qmark}
	case "JSC_IMPLICITLY_NULLABLE_JSDOC":
		// The type-based check prefers ? over ! since it only warns for names
		// that are nullable.
		return []*SuggestedFix{qmark, bang}
	case "JSC_IMPLICITLY_NONNULL_JSDOC":
		// Using type-information, we can also be confident about ! for enums
		// and typedefs.
		return []*SuggestedFix{bang}
	default:
		panic(jutil.NewIllegalArgumentExceptionMsg("Unexpected JSError Type: " + error.Type().Key()))
	}
}

// Java: ErrorToFixMapper.removeEmptyStatement(JSError)
func (m *ErrorToFixMapper) RemoveEmptyStatement(error *jscomp.JSError) *SuggestedFix {
	return NewSuggestedFix_Builder().
		AttachMatchedNodeInfo(error.Node(), m.compiler).
		SetDescription("Remove empty statement").
		DeleteWithoutRemovingWhitespace(error.Node()).
		Build()
}

// Java: ErrorToFixMapper.getFixForMissingSemicolon(JSError)
func (m *ErrorToFixMapper) GetFixForMissingSemicolon(error *jscomp.JSError) *SuggestedFix {
	return NewSuggestedFix_Builder().
		AttachMatchedNodeInfo(error.Node(), m.compiler).
		InsertAfter(error.Node(), ";").
		Build()
}

// Java: ErrorToFixMapper.getFixForMissingSuper(JSError)
func (m *ErrorToFixMapper) GetFixForMissingSuper(error *jscomp.JSError) *SuggestedFix {
	constructorFunction := error.Node().GetFirstChild()
	body := jscomp.NodeUtil_GetFunctionBody(constructorFunction)
	return NewSuggestedFix_Builder().
		AttachMatchedNodeInfo(error.Node(), m.compiler).
		AddChildToFront(body, "super();").
		Build()
}

// GetFixForInvalidSuper may return nil.
//
// Java: ErrorToFixMapper.getFixForInvalidSuper(JSError)
func (m *ErrorToFixMapper) GetFixForInvalidSuper(error *jscomp.JSError) *SuggestedFix {
	matcher := errorToFixMapper_DID_YOU_MEAN.Matcher(error.Description())
	if matcher.Matches() {
		superDotSuggestion := *jutil.CheckNotNull(matcher.Group(1))
		return NewSuggestedFix_Builder().
			AttachMatchedNodeInfo(error.Node(), m.compiler).
			SetDescription("Call '"+superDotSuggestion+"' instead").
			Replace(error.Node(), jscomp.NodeUtil_NewQName(m.compiler, superDotSuggestion), m.compiler).
			Build()
	}
	return nil
}

// GetFixForInexistentProperty may return nil.
//
// Java: ErrorToFixMapper.getFixForInexistentProperty(JSError)
func (m *ErrorToFixMapper) GetFixForInexistentProperty(error *jscomp.JSError) *SuggestedFix {
	matcher := errorToFixMapper_DID_YOU_MEAN.Matcher(error.Description())
	if matcher.Matches() {
		suggestedPropName := *matcher.Group(1)
		return NewSuggestedFix_Builder().
			AttachMatchedNodeInfo(error.Node(), m.compiler).
			SetDescription("Change property name to '"+suggestedPropName+"'").
			Rename(error.Node(), suggestedPropName).
			Build()
	}
	return nil
}

// Java: ErrorToFixMapper.getFixForMissingRequire(JSError, ImportType)
func (m *ErrorToFixMapper) GetFixForMissingRequire(error *jscomp.JSError, importType SuggestedFix_ImportType) *SuggestedFix {
	regexMatcher := errorToFixMapper_MISSING_REQUIRE.Matcher(error.Description())
	jutil.CheckStateMsg1(regexMatcher.Matches(), "Unexpected error description: %s", error.Description())
	namespaceToRequire := *regexMatcher.Group(1)

	var qName string
	if error.Node().IsQualifiedName() {
		qName = *error.Node().GetQualifiedName()
	} else {
		qName = error.Node().GetString()
	}

	jutil.CheckStateMsg2(
		jutil.StringStartsWith(qName, namespaceToRequire),
		"Expected error location %s to start with namespace <%s>",
		error.Node(),
		namespaceToRequire)

	script := m.compiler.GetScriptNode(*error.SourceName())
	scriptMetadata := m.GetMetadataForScript(script)

	metadata := NewNodeMetadata(m.compiler)
	match := NewMatch(error.Node(), metadata)
	fix :=
		NewSuggestedFix_Builder().
			AttachMatchedNodeInfo(error.Node(), m.compiler).
			AddImport(match, namespaceToRequire, importType, scriptMetadata)

	alias := scriptMetadata.GetAlias(namespaceToRequire)
	if alias != nil {
		fix.Replace(
			error.Node(),
			jscomp.NodeUtil_NewQName(m.compiler, jutil.StringReplace(qName, namespaceToRequire, *alias)),
			m.compiler)
	}

	return fix.Build()
}

// Java: ErrorToFixMapper.getFixForExtraRequire(JSError)
func (m *ErrorToFixMapper) GetFixForExtraRequire(error *jscomp.JSError) *SuggestedFix {
	node := error.Node()

	fix := NewSuggestedFix_Builder().AttachMatchedNodeInfo(node, m.compiler)
	destructuring := !jscomp.NodeUtil_GetEnclosingType(node, rhino.Token_OBJECT_PATTERN).IsNil()
	if destructuring {
		fix.SetDescription("Delete unused symbol")
		if node.IsStringKey() {
			fix.Delete(node)
		} else {
			jutil.CheckStateMsg1(node.GetParent().IsStringKey(), "%s", node.GetParent())
			fix.Delete(node.GetParent())
		}
	} else {
		fix.SetDescription("Delete extra require")
		fix.DeleteWithoutRemovingWhitespaceBefore(jscomp.NodeUtil_GetEnclosingStatement(node))
	}
	return fix.Build()
}

// GetFixForUnsortedRequires may return nil.
//
// Java: ErrorToFixMapper.getFixForUnsortedRequires(JSError)
func (m *ErrorToFixMapper) GetFixForUnsortedRequires(error *jscomp.JSError) *SuggestedFix {
	// TODO(tjgq): Encode enough information in the error to avoid the need
	// to run a traversal in order to produce the fix.
	script := jscomp.NodeUtil_GetEnclosingScript(error.Node())
	callback := jscomp.NewCheckRequiresSorted(jscomp.CheckRequiresSorted_Mode_COLLECT_ONLY)
	jscomp.NodeTraversal_Traverse(m.compiler, script, callback)

	if !callback.NeedsFix() {
		return nil
	}

	return NewSuggestedFix_Builder().
		AttachMatchedNodeInfo(callback.GetFirstNode(), m.compiler).
		ReplaceRange(callback.GetFirstNode(), callback.GetLastNode(), *callback.GetReplacement()).
		Build()
}

// GetFixForUnsortedProvides may return nil.
//
// Java: ErrorToFixMapper.getFixForUnsortedProvides(JSError)
func (m *ErrorToFixMapper) GetFixForUnsortedProvides(error *jscomp.JSError) *SuggestedFix {
	// TODO(tjgq): Encode enough information in the error to avoid the need
	// to run a traversal in order to produce the fix.
	script := jscomp.NodeUtil_GetEnclosingScript(error.Node())
	callback := jscomp.NewCheckProvidesSorted(jscomp.CheckProvidesSorted_Mode_COLLECT_ONLY)
	jscomp.NodeTraversal_Traverse(m.compiler, script, callback)

	if !callback.NeedsFix() {
		return nil
	}

	return NewSuggestedFix_Builder().
		AttachMatchedNodeInfo(callback.GetFirstNode(), m.compiler).
		ReplaceRange(callback.GetFirstNode(), callback.GetLastNode(), *callback.GetReplacement()).
		Build()
}

// Java: ErrorToFixMapper.getFixForRedundantNullabilityModifierJsDoc(JSError)
func (m *ErrorToFixMapper) GetFixForRedundantNullabilityModifierJsDoc(error *jscomp.JSError) *SuggestedFix {
	return NewSuggestedFix_Builder().
		AttachMatchedNodeInfo(error.Node(), m.compiler).
		ReplaceText(error.Node(), 1, "").
		Build()
}

// GetFixForConstantCaseErrors suggests a fix for a constant case error.
//
// If the variable is in a let clause, suggest adding a const. If the
// variable is not, suggest adding a @const annotation. Don't try to adjust
// the JSDoc, because that can produce invalid output. It may return nil.
//
// Java: ErrorToFixMapper.getFixForConstantCaseErrors(JSError)
func (m *ErrorToFixMapper) GetFixForConstantCaseErrors(error *jscomp.JSError) *SuggestedFix {
	n := error.Node()
	parent := n.GetParent()
	if !n.IsName() || !jscomp.NodeUtil_IsNameDeclaration(parent) {
		return nil
	}

	if parent.IsLet() {
		// Convert to a `const` declaration
		return NewSuggestedFix_Builder().
			SetDescription("Make explicitly constant").
			AttachMatchedNodeInfo(parent, m.compiler).
			ReplaceText(parent, 3, "const").
			Build()
	}
	// Don't convert a `var` to a `const` to avoid breaking variable scoping.
	jutil.CheckStateMsg1(parent.IsVar(), "%s", parent)
	info := parent.GetJSDocInfo()
	if info == nil {
		return NewSuggestedFix_Builder().
			SetDescription("Make explicitly constant").
			AttachMatchedNodeInfo(parent, m.compiler).
			AddOrReplaceJsDoc(parent, "/** @const */\n").
			Build()
	}
	return nil
}

// Java: ErrorToFixMapper.getMetadataForScript(Node)
func (m *ErrorToFixMapper) GetMetadataForScript(script rhino.Node) *ScriptMetadata {
	jutil.CheckArgument(script.IsScript())
	return m.metadataForEachScript.ComputeIfAbsent(
		script, func(s rhino.Node) *ScriptMetadata { return ScriptMetadata_Create(s, m.compiler) })
}
