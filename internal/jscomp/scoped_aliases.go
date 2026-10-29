// Port of com/google/javascript/jscomp/ScopedAliases.java @ closure-compiler v20261021

/*
 * Copyright 2010 The Closure Compiler Authors.
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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

var scopedAliasesJscompID = strpool.Intern("$jscomp")

// ScopedAliases processes aliases in goog.scope blocks.
//
//	goog.scope(function() {
//	  var dom = goog.dom;
//	  var DIV = dom.TagName.DIV;
//
//	  dom.createElement(DIV);
//	});
//
// should become
//
//	goog.dom.createElement(goog.dom.TagName.DIV);
//
// The advantage of using goog.scope is that the compiler will *guarantee*
// the anonymous function will be inlined, even if it can't prove that it's
// semantically correct to do so. For example, consider this case:
//
//	goog.scope(function() {
//	  goog.getBar = function () { return alias; };
//	  ...
//	  var alias = foo.bar;
//	})
//
// In theory, the compiler can't inline 'alias' unless it can prove that
// goog.getBar is called only after 'alias' is defined. In practice, the
// compiler will inline 'alias' anyway, at the risk of 'fixing' bad code.
type ScopedAliases struct {
	compiler                AbstractCompiler
	preprocessorSymbolTable *PreprocessorSymbolTable
	moduleMetadataMap       *ModuleMetadataMap
	uniqueIdInput           *CompilerInput

	closureNamespaces        map[string]struct{}
	invalidModuleGetHandling ScopedAliases_InvalidModuleGetHandling
}

var _ CompilerPass = (*ScopedAliases)(nil)

// ScopedAliases_SCOPING_METHOD_NAME is the name used to denote an scoped
// function block used for aliasing.
const ScopedAliases_SCOPING_METHOD_NAME = "goog.scope"

const scopedAliases_SCOPED_ALIASES_PREFIX = "$jscomp$scope$"

// LINT.IfChange
const scopedAliases_MISSING_ALIASES_PREFIX = "$$jscomp$missingAlias$"

// LINT.ThenChange(//depot/google3/third_party/java_src/clutz/src/main/java/com/google/javascript/clutz/imports/ImportRenameMapBuilder.java)

// Errors

// Java: ScopedAliases.GOOG_SCOPE_MUST_BE_ALONE
var ScopedAliases_GOOG_SCOPE_MUST_BE_ALONE = DiagnosticType_Error(
	"JSC_GOOG_SCOPE_MUST_BE_ALONE",
	"The call to goog.scope must be alone in a single statement.")

// Java: ScopedAliases.GOOG_SCOPE_MUST_BE_IN_GLOBAL_SCOPE
var ScopedAliases_GOOG_SCOPE_MUST_BE_IN_GLOBAL_SCOPE = DiagnosticType_Error(
	"JSC_GOOG_SCOPE_MUST_BE_IN_GLOBAL_SCOPE",
	"The call to goog.scope must be in the global scope.")

// Java: ScopedAliases.GOOG_SCOPE_HAS_BAD_PARAMETERS
var ScopedAliases_GOOG_SCOPE_HAS_BAD_PARAMETERS = DiagnosticType_Error(
	"JSC_GOOG_SCOPE_HAS_BAD_PARAMETERS",
	"The call to goog.scope must take only a single parameter.  It must"+
		" be an anonymous function that itself takes no parameters.")

// Java: ScopedAliases.GOOG_SCOPE_REFERENCES_THIS
var ScopedAliases_GOOG_SCOPE_REFERENCES_THIS = DiagnosticType_Error(
	"JSC_GOOG_SCOPE_REFERENCES_THIS",
	"The body of a goog.scope function cannot reference 'this'.")

// Java: ScopedAliases.GOOG_SCOPE_USES_RETURN
var ScopedAliases_GOOG_SCOPE_USES_RETURN = DiagnosticType_Error(
	"JSC_GOOG_SCOPE_USES_RETURN", "The body of a goog.scope function cannot use 'return'.")

// Java: ScopedAliases.GOOG_SCOPE_USES_THROW
var ScopedAliases_GOOG_SCOPE_USES_THROW = DiagnosticType_Error(
	"JSC_GOOG_SCOPE_USES_THROW", "The body of a goog.scope function cannot use 'throw'.")

// Java: ScopedAliases.GOOG_SCOPE_ALIAS_REDEFINED
var ScopedAliases_GOOG_SCOPE_ALIAS_REDEFINED = DiagnosticType_Error(
	"JSC_GOOG_SCOPE_ALIAS_REDEFINED", "The alias {0} is assigned a value more than once.")

// Java: ScopedAliases.GOOG_SCOPE_ALIAS_CYCLE
var ScopedAliases_GOOG_SCOPE_ALIAS_CYCLE = DiagnosticType_Error("JSC_GOOG_SCOPE_ALIAS_CYCLE", "The aliases {0} has a cycle.")

// Java: ScopedAliases.GOOG_SCOPE_NON_ALIAS_LOCAL
var ScopedAliases_GOOG_SCOPE_NON_ALIAS_LOCAL = DiagnosticType_Error(
	"JSC_GOOG_SCOPE_NON_ALIAS_LOCAL",
	"The local variable {0} is in a goog.scope and is not an alias.")

// Java: ScopedAliases.GOOG_SCOPE_INVALID_VARIABLE
var ScopedAliases_GOOG_SCOPE_INVALID_VARIABLE = DiagnosticType_Error(
	"JSC_GOOG_SCOPE_INVALID_VARIABLE", "The variable {0} cannot be declared in this scope")

var scopedAliases_GOOG_MODULE_GET = rhino.QualifiedName_Of("goog.module.get")

// ScopedAliases_InvalidModuleGetHandling is what to do with goog.module.get
// calls importing an inexistent Closure namespace.
type ScopedAliases_InvalidModuleGetHandling int32

const (
	ScopedAliases_InvalidModuleGetHandling_PRESERVE ScopedAliases_InvalidModuleGetHandling = iota
	ScopedAliases_InvalidModuleGetHandling_GIVE_UNIQUE_NAME
)

// Java: ScopedAliases.InvalidModuleGetHandling.name()
func (h ScopedAliases_InvalidModuleGetHandling) String() string {
	return [...]string{"PRESERVE", "GIVE_UNIQUE_NAME"}[h]
}

// Java: ScopedAliases.InvalidModuleGetHandling.shouldGiveUniqueName()
func (h ScopedAliases_InvalidModuleGetHandling) ShouldGiveUniqueName() bool {
	return h == ScopedAliases_InvalidModuleGetHandling_GIVE_UNIQUE_NAME
}

// Java: ScopedAliases(AbstractCompiler, PreprocessorSymbolTable, Set<String>, InvalidModuleGetHandling)
func newScopedAliases(
	compiler AbstractCompiler,
	preprocessorSymbolTable *PreprocessorSymbolTable,
	closureNamespaces map[string]struct{},
	invalidModuleGetHandling ScopedAliases_InvalidModuleGetHandling) *ScopedAliases {
	return &ScopedAliases{
		compiler:                 compiler,
		preprocessorSymbolTable:  preprocessorSymbolTable,
		moduleMetadataMap:        compiler.GetModuleMetadataMap(),
		closureNamespaces:        closureNamespaces,
		invalidModuleGetHandling: invalidModuleGetHandling,
	}
}

// Java: ScopedAliases.builder(AbstractCompiler)
func ScopedAliases_Builder_(compiler AbstractCompiler) *ScopedAliases_Builder {
	return newScopedAliases_Builder(compiler)
}

// ScopedAliases_Builder is ScopedAliases.Builder.
type ScopedAliases_Builder struct {
	compiler                 AbstractCompiler
	preprocessorSymbolTable  *PreprocessorSymbolTable
	moduleMetadataMap        *ModuleMetadataMap
	invalidModuleGetHandling ScopedAliases_InvalidModuleGetHandling
}

// Java: ScopedAliases.Builder(AbstractCompiler)
func newScopedAliases_Builder(compiler AbstractCompiler) *ScopedAliases_Builder {
	return &ScopedAliases_Builder{
		compiler:                 compiler,
		invalidModuleGetHandling: ScopedAliases_InvalidModuleGetHandling_PRESERVE,
	}
}

// Java: ScopedAliases.Builder.setPreprocessorSymbolTable(PreprocessorSymbolTable)
func (b *ScopedAliases_Builder) SetPreprocessorSymbolTable(preprocessorSymbolTable *PreprocessorSymbolTable) *ScopedAliases_Builder {
	b.preprocessorSymbolTable = preprocessorSymbolTable
	return b
}

// Java: ScopedAliases.Builder.setModuleMetadataMap(ModuleMetadataMap)
func (b *ScopedAliases_Builder) SetModuleMetadataMap(moduleMetadataMap *ModuleMetadataMap) *ScopedAliases_Builder {
	b.moduleMetadataMap = moduleMetadataMap
	return b
}

// SetInvalidModuleGetHandling configures whether to delete or preserve
// invalid goog.module get calls that are top-level aliases in a goog.scope.
//
// Java: ScopedAliases.Builder.setInvalidModuleGetHandling(InvalidModuleGetHandling)
func (b *ScopedAliases_Builder) SetInvalidModuleGetHandling(invalidModuleGetHandling ScopedAliases_InvalidModuleGetHandling) *ScopedAliases_Builder {
	b.invalidModuleGetHandling = invalidModuleGetHandling
	return b
}

// Java: ScopedAliases.Builder.build()
func (b *ScopedAliases_Builder) Build() *ScopedAliases {
	var closureNamespaces map[string]struct{}
	if b.moduleMetadataMap != nil {
		entries := b.moduleMetadataMap.ModulesByGoogNamespaceMap().Entries()
		if len(entries) > 0 {
			closureNamespaces = make(map[string]struct{}, len(entries))
			for _, e := range entries {
				closureNamespaces[e.Key] = struct{}{}
			}
		}
	}
	return newScopedAliases(
		b.compiler,
		b.preprocessorSymbolTable,
		closureNamespaces,
		b.invalidModuleGetHandling)
}

// Java: ScopedAliases.process(Node, Node)
func (s *ScopedAliases) Process(externs rhino.Node, root rhino.Node) {
	traversal := s.NewTraversal()
	NodeTraversal_Traverse(s.compiler, root, traversal)

	if !traversal.HasErrors() {
		// Apply the aliases.
		aliasWorkQueue := traversal.aliasUsages
		for len(aliasWorkQueue) != 0 {
			n := 0
			for _, u := range aliasWorkQueue {
				if u.referencesOtherAlias(traversal.deletedAliasVars) {
					aliasWorkQueue[n] = u
					n++
				} else {
					u.applyAlias(s.compiler)
				}
			}

			// Prevent an infinite loop.
			if n == len(aliasWorkQueue) {
				cycleVar := aliasWorkQueue[0].aliasVar
				s.compiler.Report(
					JSError_MakeWithNode(cycleVar.GetNode(), ScopedAliases_GOOG_SCOPE_ALIAS_CYCLE, cycleVar.GetName()))
				break
			}
			clear(aliasWorkQueue[n:])
			aliasWorkQueue = aliasWorkQueue[:n]
		}

		// Remove the alias definitions.
		for _, aliasDefinition := range traversal.aliasDefinitionsToDelete {
			s.compiler.ReportChangeToEnclosingScope(aliasDefinition)
			if NodeUtil_IsNameDeclaration(aliasDefinition.GetParent()) &&
				aliasDefinition.GetParent().HasOneChild() {
				aliasDefinition.GetParent().Detach()
			} else {
				aliasDefinition.Detach()
			}
		}

		// Collapse the scopes.
		for _, scopeCall := range traversal.scopeCalls {
			expressionWithScopeCall := scopeCall.GetParent()
			scopeClosureBlock := ast.MustFunction(ast.MustCall(scopeCall).MustOnlyArg().Node()).Body().Node()
			scopeClosureBlock.Detach()
			expressionWithScopeCall.ReplaceWith(scopeClosureBlock)
			NodeUtil_MarkFunctionsDeleted(expressionWithScopeCall, s.compiler)
			s.compiler.ReportChangeToEnclosingScope(scopeClosureBlock)
			NodeUtil_TryMergeBlock(scopeClosureBlock, false)
		}
	}
}

// aliasUsage collapses ScopedAliases.AliasUsage, AliasedNode, and AliasedTypeNode.
type aliasUsage struct {
	aliasVar       *Var
	aliasReference rhino.Node
	isTypeNode     bool
}

func scopedAliases_aliasedRootID(rhs rhino.Node, intern bool) (rootID uint32, isQualified bool) {
	cur := rhs
	for cur.IsGetProp() {
		isQualified = true
		cur = cur.GetFirstChild()
	}
	switch cur.GetToken() {
	case rhino.Token_NAME:
		return cur.GetStringID(), isQualified
	case rhino.Token_CALL:
		jutil.CheckStateMsg1(NodeUtil_IsCallToQualifiedName(cur, scopedAliases_GOOG_MODULE_GET), "%s", cur)
		arg, ok := ast.MustCall(cur).OnlyArg()
		jutil.CheckStateMsg1(ok, "%s", cur)
		ns := arg.Node().GetString()
		if dotIndex := strings.IndexByte(ns, '.'); dotIndex != -1 {
			root := ns[:dotIndex]
			if intern {
				return uint32(strpool.Intern(root)), true
			}
			return uint32(strpool.Lookup(root)), true
		}
		return arg.Node().GetStringID(), isQualified
	default:
		panic(jutil.NewRuntimeException("Invalid alias RHS:" + rhs.String()))
	}
}

func (u aliasUsage) referencesOtherAlias(deletedAliasVars map[*Var]struct{}) bool {
	rootID, _ := scopedAliases_aliasedRootID(u.aliasVar.GetInitialValue(), false)
	if rootID == 0 {
		return false
	}
	otherAliasVar := u.aliasVar.scope.lookupOwnVarByID(rootID)
	if otherAliasVar == nil {
		return false
	}
	_, deleted := deletedAliasVars[otherAliasVar]
	return !deleted
}

func (u aliasUsage) applyAlias(compiler AbstractCompiler) {
	if u.isTypeNode {
		aliasDefinition := u.aliasVar.GetInitialValue()
		aliasName := u.aliasVar.GetName()
		typeName := u.aliasReference.GetString()
		if strings.HasPrefix(typeName, "$jscomp$scope$") {
			// Already visited.
			return
		}
		aliasExpanded := scopedAliases_getAliasedNamespace(aliasDefinition)
		jutil.CheckStateMsg2(
			strings.HasPrefix(typeName, aliasName), "%s must start with %s", typeName, aliasName)
		u.aliasReference.SetString(aliasExpanded + typeName[len(aliasName):])
		return
	}

	aliasDefinition := u.aliasVar.GetInitialValue()
	replacement := aliasDefinition.CloneTree()
	replacement.SrcrefTree(u.aliasReference)
	// Given alias "var Bar = foo.Bar;" here we replace a usage of Bar with
	// foo.Bar. foo is generated and never visible to user. Because of that we
	// should mark all new nodes as non-indexable leaving only Bar indexable.
	// Given that replacement is GETPROP node, prefix is first child. It's also
	// possible that replacement is single-part namespace. Like
	// goog.provide('Foo') in that case replacement won't have children.
	if replacement.HasChildren() {
		replacement.GetFirstChild().MakeNonIndexableRecursive()
	}
	u.aliasReference.ReplaceWith(replacement)
	compiler.ReportChangeToEnclosingScope(replacement)
}

// Java: ScopedAliases.isValidAliasRhs(Node)
func scopedAliases_isValidAliasRhs(rhs rhino.Node) bool {
	switch rhs.GetToken() {
	case rhino.Token_GETPROP:
		return scopedAliases_isValidAliasRhs(ast.MustGetProp(rhs).Object().Node())
	case rhino.Token_NAME:
		return true
	case rhino.Token_CALL:
		return NodeUtil_IsCallToQualifiedName(rhs, scopedAliases_GOOG_MODULE_GET)
	default:
		return false
	}
}

// Java: ScopedAliases.isAliasDefinition(Node)
func scopedAliases_isAliasDefinition(nameNode rhino.Node) bool {
	initExpr, ok := ast.MustName(nameNode).Init()
	return ok && scopedAliases_isValidAliasRhs(initExpr.Node())
}

// Java: ScopedAliases.getAliasedNamespace(Node)
func scopedAliases_getAliasedNamespace(rhs rhino.Node) string {
	switch rhs.GetToken() {
	case rhino.Token_GETPROP:
		gp := ast.MustGetProp(rhs)
		return scopedAliases_getAliasedNamespace(gp.Object().Node()) + "." + gp.Property()
	case rhino.Token_NAME:
		return ast.MustName(rhs).Ident()
	case rhino.Token_CALL:
		jutil.CheckStateMsg1(NodeUtil_IsCallToQualifiedName(rhs, scopedAliases_GOOG_MODULE_GET), "%s", rhs)
		arg, ok := ast.MustCall(rhs).OnlyArg()
		jutil.CheckStateMsg1(ok, "%s", rhs)
		return arg.Node().GetString()
	default:
		panic(jutil.NewRuntimeException("Invalid alias RHS:" + rhs.String()))
	}
}

// ScopedAliases_Traversal is ScopedAliases.Traversal (an inner class: outer
// is the enclosing ScopedAliases).
type ScopedAliases_Traversal struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *ScopedAliases

	// The job of this class is to collect these three data sets.
	aliasDefinitionsToDelete []rhino.Node
	scopeCalls               []rhino.Node
	aliasUsages              []aliasUsage

	// This map is temporary and cleared for each scope.
	aliasByID map[uint32]*Var

	// Also temporary and cleared for each scope.
	injectedDecls map[rhino.Node]struct{}

	// Persists across scopes.
	deletedAliasVars map[*Var]struct{}

	// Suppose you create an alias.
	// var x = goog.x;
	// As a side-effect, this means you can shadow the namespace 'goog'
	// in inner scopes. When we inline the namespaces, we have to rename
	// these shadows.
	//
	// Fortunately, we already have a name uniquifier that runs during tree
	// normalization (before optimizations). We run it here on a limited
	// set of variables, but only as a last resort (because this will screw
	// up warning messages downstream).
	forbiddenLocals map[strpool.ID]struct{}

	hasNamespaceShadows bool
	hasErrors           bool

	// The body of the function that is passed to goog.scope.
	// Set when the traversal enters the body, and set back to null when it
	// exits.
	scopeFunctionBody rhino.Node
}

var _ NodeTraversal_ScopedCallback = (*ScopedAliases_Traversal)(nil)

// Java: ScopedAliases.Traversal()
func (s *ScopedAliases) NewTraversal() *ScopedAliases_Traversal {
	return &ScopedAliases_Traversal{outer: s}
}

// Java: ScopedAliases.Traversal.hasErrors()
func (v *ScopedAliases_Traversal) HasErrors() bool {
	return v.hasErrors
}

// InGoogScopeBody returns true if this NodeTraversal is currently within a
// goog.scope function body.
//
// Java: ScopedAliases.Traversal.inGoogScopeBody()
func (v *ScopedAliases_Traversal) InGoogScopeBody() bool {
	return !v.scopeFunctionBody.IsNil()
}

// IsGoogScopeFunctionBody returns true if n is the goog.scope function body.
//
// Java: ScopedAliases.Traversal.isGoogScopeFunctionBody(Node)
func (v *ScopedAliases_Traversal) IsGoogScopeFunctionBody(n rhino.Node) bool {
	return v.InGoogScopeBody() && n == v.scopeFunctionBody
}

// Java: ScopedAliases.Traversal.isCallToScopeMethod(Node)
func (v *ScopedAliases_Traversal) IsCallToScopeMethod(n rhino.Node) bool {
	call, ok := ast.AsCall(n)
	return ok && call.Callee().Node().MatchesQualifiedName(ScopedAliases_SCOPING_METHOD_NAME)
}

// FindScopeMethodCall: scopeRoot is the Node which is the root of the
// current scope. Returns the goog.scope() CALL node containing the
// scopeRoot, or null if scopeRoot is not in a goog.scope() call.
//
// Java: ScopedAliases.Traversal.findScopeMethodCall(Node)
func (v *ScopedAliases_Traversal) FindScopeMethodCall(scopeRoot rhino.Node) rhino.Node {
	n := scopeRoot.GetGrandparent()
	if v.IsCallToScopeMethod(n) {
		return n
	}
	return rhino.NilNode
}

// Java: ScopedAliases.Traversal.enterScope(NodeTraversal)
func (v *ScopedAliases_Traversal) EnterScope(t *NodeTraversal) {
	if t.InGlobalHoistScope() {
		return
	}
	scopeMethodCall := v.FindScopeMethodCall(t.GetScopeRoot())
	if !scopeMethodCall.IsNil() {
		if v.forbiddenLocals == nil {
			v.forbiddenLocals = map[strpool.ID]struct{}{scopedAliasesJscompID: {}}
		}
		v.FindAliases(t.GetScope())
		if lastArg, ok := ast.MustCall(scopeMethodCall).LastArg(); ok {
			if fn, ok := ast.AsFunction(lastArg.Node()); ok {
				v.scopeFunctionBody = fn.Body().Node()
			} else {
				v.scopeFunctionBody = lastArg.Node().GetLastChild()
			}
		}
	}
}

// Java: ScopedAliases.Traversal.exitScope(NodeTraversal)
func (v *ScopedAliases_Traversal) ExitScope(t *NodeTraversal) {
	if v.IsGoogScopeFunctionBody(t.GetScopeRoot()) {
		v.scopeFunctionBody = rhino.NilNode
		v.RenameNamespaceShadows(t)
		clear(v.injectedDecls)
		clear(v.aliasByID)
		clear(v.forbiddenLocals)
		v.hasNamespaceShadows = false
	} else if v.InGoogScopeBody() {
		// Called on inner scopes within a goog.scope, including both block
		// scopes and function scopes.
		v.FindNamespaceShadows(t)
		v.ReportInvalidVariables(t)
	}
}

// Java: ScopedAliases.Traversal.reportInvalidVariables(NodeTraversal)
func (v *ScopedAliases_Traversal) ReportInvalidVariables(t *NodeTraversal) {
	scopeRoot := t.GetScopeRoot()
	enclosingFunctionBody := ast.MustFunction(t.GetEnclosingFunction()).Body().Node()
	if v.IsGoogScopeFunctionBody(enclosingFunctionBody) &&
		scopeRoot.IsBlock() &&
		!scopeRoot.GetParent().IsFunction() {
		for var_ := range t.GetScope().GetVarIterable() {
			parent := var_.GetNameNode().GetParent()
			if NodeUtil_IsFunctionDeclaration(parent) {
				// Disallow block-scoped function declarations that leak into the
				// goog.scope function body. Technically they shouldn't leak in ES6
				// but the browsers don't agree on that yet.
				v.Report(var_.GetNode(), ScopedAliases_GOOG_SCOPE_INVALID_VARIABLE, var_.GetName())
			}
		}
	}
}

// Java: ScopedAliases.Traversal.report(Node, DiagnosticType, String...)
func (v *ScopedAliases_Traversal) Report(n rhino.Node, error_ *DiagnosticType, arguments ...string) {
	v.outer.compiler.Report(JSError_MakeWithNode(n, error_, arguments...))
	v.hasErrors = true
}

// Java: ScopedAliases.Traversal.findAliases(Scope)
func (v *ScopedAliases_Traversal) FindAliases(scope *Scope) {
	compiler := v.outer.compiler
	for var_ := range scope.GetVarIterable() {
		n := var_.GetNode()
		parent := n.GetParent()
		// We use isBlock to avoid variables declared in loop headers.
		isVar := NodeUtil_IsNameDeclaration(parent) && parent.GetParent().IsBlock()
		isFunctionDecl := NodeUtil_IsFunctionDeclaration(parent)
		if isVar && scopedAliases_isAliasDefinition(n) {
			v.RecordAlias(var_)
		} else if var_.IsBleedingFunction() {
			// Bleeding functions already get a BAD_PARAMETERS error, so just
			// do nothing.
		} else if parent.IsParamList() {
			// Parameters of the scope function also get a BAD_PARAMETERS
			// error.
		} else if isVar || isFunctionDecl || NodeUtil_IsClassDeclaration(parent) {
			isHoisted := NodeUtil_IsHoistedFunctionDeclaration(parent)
			grandparent := parent.GetParent()
			value := var_.GetInitialValue()
			var varNode rhino.Node

			// Pull out inline type declaration if present.
			if n.GetJSDocInfo() != nil {
				builder := rhino.JSDocInfo_Builder_MaybeCopyFrom(parent.GetJSDocInfo())
				if isFunctionDecl { // Fix inline return type.
					builder.RecordReturnType(n.GetJSDocInfo().GetType())
				} else {
					builder.RecordType(n.GetJSDocInfo().GetType())
				}
				parent.SetJSDocInfo(builder.Build())
				n.SetJSDocInfo(nil)
			}
			// Grab the docinfo before we do any AST manipulation.
			varDocInfo := var_.GetJSDocInfo()

			name := n.GetString()
			uniqueId := compiler.GetUniqueIdSupplier().GetUniqueId(v.outer.uniqueIdInput)

			globalName := "$jscomp$scope$" + uniqueId + "$" + name

			// First, we need to free up the function expression (EXPR)
			// to be used in another expression.
			if isFunctionDecl || NodeUtil_IsClassDeclaration(parent) {
				// Replace "function NAME() { ... }" with "var NAME;".
				// Replace "class NAME { ... }" with "var NAME;".
				var newName rhino.Node
				if isFunctionDecl {
					newName = rhino.IR_Name("")
				} else {
					newName = rhino.IR_Empty()
				}
				newName.Srcref(n)
				n.ReplaceWith(newName)
				compiler.ReportChangeToEnclosingScope(newName)

				varNode = rhino.IR_Var(n).Srcref(n)
				parent.ReplaceWith(varNode)
			} else {
				if !value.IsNil() {
					// If this is a VAR, we can just detach the expression and
					// the tree will still be valid.
					value.Detach()
				}
				varNode = parent
			}

			// Add var declarations for name
			// var $jscomp$scope$uniqueId$name = EXPR; or var $jscomp$scope$uniqueId$name;
			newDecl :=
				NodeUtil_NewVarNode(globalName, value).
					SetJSDocInfo(varDocInfo).
					SrcrefTreeIfMissing(n)

			newDecl.Srcref(n)
			newDecl.SetOriginalName(&name)
			if isHoisted {
				grandparent.AddChildToFront(newDecl)
			} else {
				newDecl.InsertBefore(varNode)
			}
			if !value.IsNil() {
				compiler.ReportChangeToEnclosingScope(newDecl)
			}
			if v.injectedDecls == nil {
				v.injectedDecls = make(map[rhino.Node]struct{})
			}
			v.injectedDecls[newDecl] = struct{}{}

			// Rewrite "var name = EXPR;" to "var name = $jscomp$scope$uniqueId$name;"
			var_.GetNameNode().AddChildToFront(NodeUtil_NewNameWithOriginalName(compiler, globalName, n, name))
			v.RecordAlias(var_)
		} else {
			// Do not allow other kinds of local symbols, like catch params.
			v.Report(n, ScopedAliases_GOOG_SCOPE_NON_ALIAS_LOCAL, n.GetString())
		}
	}
}

// Java: ScopedAliases.Traversal.recordAlias(Var)
func (v *ScopedAliases_Traversal) RecordAlias(aliasVar *Var) {
	initialValue := aliasVar.GetInitialValue()
	nameNode := aliasVar.GetNameNode()
	v.aliasDefinitionsToDelete = append(v.aliasDefinitionsToDelete, nameNode)

	if v.aliasByID == nil {
		v.aliasByID = make(map[uint32]*Var)
	}

	if v.outer.invalidModuleGetHandling.ShouldGiveUniqueName() &&
		v.ContainsInvalidGoogModuleGet(initialValue) {
		// Create a new alias with a unique name, to replace the existing one.
		// Needed for integration with Clutz only: by giving the alias a unique
		// name, Clutz may then map back from the unique name to the original
		// `const alias = goog.module.get(` call.
		name := aliasVar.GetName()
		fileName := initialValue.GetSourceFileName()
		var metadata *ModuleMetadataMap_ModuleMetadata
		if fileName != nil {
			metadata = v.outer.moduleMetadataMap.ModuleByPath(*fileName)
		}
		// LINT.IfChange
		suffix := "unknownScript"
		if metadata != nil {
			suffix = metadata.Path().ToModuleName()
		}
		newName := name + scopedAliases_MISSING_ALIASES_PREFIX + suffix
		// LINT.ThenChange(//depot/google3/third_party/java_src/clutz/src/main/java/com/google/javascript/clutz/imports/ImportRenameMapBuilder.java)
		initialValue.ReplaceWith(rhino.IR_Name(newName))
		if v.deletedAliasVars == nil {
			v.deletedAliasVars = make(map[*Var]struct{})
		}
		v.deletedAliasVars[aliasVar] = struct{}{}
		v.aliasByID[nameNode.GetStringID()] = aliasVar
		return
	}

	v.aliasByID[nameNode.GetStringID()] = aliasVar

	if rootID, isQualified := scopedAliases_aliasedRootID(initialValue, true); isQualified && rootID != 0 {
		if _, ok := v.aliasByID[rootID]; !ok {
			v.forbiddenLocals[strpool.ID(rootID)] = struct{}{}
		}
	}
}

// ContainsInvalidGoogModuleGet returns whether the rhs contains any
// goog.module.get calls to inexistent namespaces.
//
// Java: ScopedAliases.Traversal.containsInvalidGoogModuleGet(Node)
func (v *ScopedAliases_Traversal) ContainsInvalidGoogModuleGet(expression rhino.Node) bool {
	switch expression.GetToken() {
	case rhino.Token_NAME:
		return false
	case rhino.Token_GETPROP:
		return v.ContainsInvalidGoogModuleGet(ast.MustGetProp(expression).Object().Node())
	case rhino.Token_CALL:
		namespace := ast.MustCall(expression).MustFirstArg().Node().GetString()
		_, ok := v.outer.closureNamespaces[namespace]
		return !ok
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Unrecognized alias rhs " + expression.String()))
	}
}

// FindNamespaceShadows finds out if there are any local shadows of
// namespaces.
//
// Java: ScopedAliases.Traversal.findNamespaceShadows(NodeTraversal)
func (v *ScopedAliases_Traversal) FindNamespaceShadows(t *NodeTraversal) {
	if v.hasNamespaceShadows {
		return
	}

	scope := t.GetScope()
	for var_ := range scope.GetVarIterable() {
		if _, ok := v.forbiddenLocals[strpool.ID(var_.nameID)]; ok {
			v.hasNamespaceShadows = true
			return
		}
	}
}

// RenameNamespaceShadows renames any local shadows of namespaces. This should
// be a very rare occurrence, so only do this traversal if we know that we
// need it.
//
// Java: ScopedAliases.Traversal.renameNamespaceShadows(NodeTraversal)
func (v *ScopedAliases_Traversal) RenameNamespaceShadows(t *NodeTraversal) {
	jutil.CheckStateMsg1(NodeUtil_IsFunctionBlock(t.GetScopeRoot()), "%s", t.GetScopeRoot())

	if v.hasNamespaceShadows {
		renamer := &MakeDeclaredNamesUnique_TargettedRenamer{
			delegate:   NewMakeDeclaredNamesUnique_ContextualRenamer(),
			targetIDs:  v.forbiddenLocals,
			cachedSize: len(v.forbiddenLocals),
		}
		// order-independent: root ContextualRenamer.addDeclaredNameID sets each independent nameUsage[id] = 1
		for id := range v.forbiddenLocals {
			renamer.addDeclaredNameID(id, false)
		}
		uniquifier :=
			MakeDeclaredNamesUnique_Builder_().WithRenamer(renamer).Build()
		NodeTraversal_TraverseScopeRoots(
			v.outer.compiler, []rhino.Node{t.GetScopeRoot()}, uniquifier, true)
	}
}

var scopedAliasesRenameBleedingTokens = makeVisitTokenSet(rhino.Token_NAME)

// scopedAliases_renameBleedingFunctionNameCallback is the anonymous
// AbstractPostOrderCallback of Traversal.renameBleedingFunctionName.
type scopedAliases_renameBleedingFunctionNameCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler
	name     string
	nameID   uint32
	suffix   string
	fnName   rhino.Node
}

var (
	_ nodeTraversal_visitTokens    = (*scopedAliases_renameBleedingFunctionNameCallback)(nil)
	_ nodeTraversal_onlyNameLeaves = (*scopedAliases_renameBleedingFunctionNameCallback)(nil)
)

func (*scopedAliases_renameBleedingFunctionNameCallback) visitTokens() *ast.TokenSet {
	return &scopedAliasesRenameBleedingTokens
}

func (*scopedAliases_renameBleedingFunctionNameCallback) onlyNameLeaves() {}

// Java: ScopedAliases.Traversal.renameBleedingFunctionName(NodeTraversal, Node).new AbstractPostOrderCallback() {...}.visit(NodeTraversal, Node, Node)
func (c *scopedAliases_renameBleedingFunctionNameCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if name, ok := ast.AsName(n); ok &&
		n.GetStringID() == c.nameID &&
		t.GetScope().GetVarByID(c.nameID).GetNode() == c.fnName {
		name.SetIdent(c.name + "$jscomp$scopedAliases$" + c.suffix)
		c.compiler.ReportChangeToEnclosingScope(n)
	}
}

// Java: ScopedAliases.Traversal.renameBleedingFunctionName(NodeTraversal, Node)
func (v *ScopedAliases_Traversal) RenameBleedingFunctionName(t *NodeTraversal, fnName rhino.Node) {
	compiler := v.outer.compiler
	name := fnName.GetString()
	suffix := compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())
	cb := &scopedAliases_renameBleedingFunctionNameCallback{
		compiler: compiler, name: name, nameID: fnName.GetStringID(), suffix: suffix, fnName: fnName}
	NodeTraversal_Builder_().
		SetCompiler(compiler).
		SetCallback(cb).
		SetScopeCreator(t.GetScopeCreator()).
		TraverseAtScope(t.GetScope())
	fnName.SetString(name + "$jscomp$scopedAliases$" + suffix)
}

// Java: ScopedAliases.Traversal.validateScopeCall(NodeTraversal, Node, Node)
func (v *ScopedAliases_Traversal) ValidateScopeCall(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	call := ast.MustCall(n)
	if v.outer.preprocessorSymbolTable != nil {
		v.outer.preprocessorSymbolTable.AddReference(call.Callee().Node())
	}
	if !parent.IsExprResult() {
		v.Report(n, ScopedAliases_GOOG_SCOPE_MUST_BE_ALONE)
	}
	if !t.GetEnclosingFunction().IsNil() {
		v.Report(n, ScopedAliases_GOOG_SCOPE_MUST_BE_IN_GLOBAL_SCOPE)
	}
	arg, ok := call.OnlyArg()
	if !ok {
		// The goog.scope call should have exactly 1 parameter.
		v.Report(n, ScopedAliases_GOOG_SCOPE_HAS_BAD_PARAMETERS)
	} else {
		anonymousFnNode := arg.Node()
		fn, isFn := ast.AsFunction(anonymousFnNode)
		if !isFn || NodeUtil_GetName(anonymousFnNode) != nil || fn.Params().HasParams() {
			v.Report(anonymousFnNode, ScopedAliases_GOOG_SCOPE_HAS_BAD_PARAMETERS)
		} else {
			v.scopeCalls = append(v.scopeCalls, n)
		}
	}
}

// Java: ScopedAliases.Traversal.visit(NodeTraversal, Node, Node)
func (v *ScopedAliases_Traversal) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	compiler := v.outer.compiler
	// grab the source file path to generate the unique id
	v.outer.uniqueIdInput = t.GetInput()
	if v.IsCallToScopeMethod(n) {
		v.ValidateScopeCall(t, n, n.GetParent())
	}

	if !v.InGoogScopeBody() {
		return
	}

	type_ := n.GetToken()
	var aliasVar *Var
	if type_ == rhino.Token_NAME {
		id := n.GetStringID()
		if candidate := v.aliasByID[id]; candidate != nil {
			if t.GetScope().GetVarByID(id) == candidate {
				aliasVar = candidate
				// For nodes that are referencing the aliased type, set the original
				// name so it can be accessed later in tools such as the CodePrinter
				// or refactoring tools.
				if compiler.GetOptions().PreservesDetailedSourceInfo() {
					name := n.GetString()
					n.SetOriginalName(&name)
				}
			}
		}
	}

	if v.IsGoogScopeFunctionBody(ast.MustFunction(t.GetEnclosingFunction()).Body().Node()) {
		if aliasVar != nil && NodeUtil_IsLValue(n) {
			if aliasVar.GetNode() == n {
				// Return early, to ensure that we don't record this as an alias
				// usage.
				return
			}
			v.Report(n, ScopedAliases_GOOG_SCOPE_ALIAS_REDEFINED, n.GetString())
		}

		switch type_ {
		case rhino.Token_RETURN:
			v.Report(n, ScopedAliases_GOOG_SCOPE_USES_RETURN)
		case rhino.Token_THIS:
			v.Report(n, ScopedAliases_GOOG_SCOPE_REFERENCES_THIS)
		case rhino.Token_THROW:
			v.Report(n, ScopedAliases_GOOG_SCOPE_USES_THROW)
		default:
		}
	}

	// If this is a bleeding function expression, like
	// var x = function y() { ... }
	// then old versions of IE declare "y" in the current scope. We don't
	// want the scope unboxing to add "y" to the global scope, so we
	// need to rename it.
	if NodeUtil_IsBleedingFunctionName(n) {
		v.RenameBleedingFunctionName(t, n)
	}

	// Check if this name points to an alias.
	if aliasVar != nil {
		// Note, to support the transitive case, it's important we don't
		// clone aliasedNode here.
		v.aliasUsages = append(v.aliasUsages, aliasUsage{aliasVar: aliasVar, aliasReference: n})
	}

	// When we inject declarations, we duplicate jsdoc. Make sure
	// we only process that jsdoc once.
	if info := n.GetJSDocInfo(); info != nil {
		if _, injected := v.injectedDecls[n]; !injected {
			for it := info.TypeExprIter(); it.Next(); {
				v.FixTypeNode(it.Node())
			}
		}
	}
}

// Java: ScopedAliases.Traversal.fixTypeNode(Node)
func (v *ScopedAliases_Traversal) FixTypeNode(typeNode rhino.Node) {
	if strLit, ok := ast.AsStringLit(typeNode); ok {
		name := strLit.Value()
		var baseID uint32
		if endIndex := strings.IndexByte(name, '.'); endIndex != -1 {
			baseID = uint32(strpool.Lookup(name[:endIndex]))
		} else {
			baseID = typeNode.GetStringID()
		}
		if aliasVar := v.aliasByID[baseID]; aliasVar != nil {
			v.aliasUsages = append(v.aliasUsages, aliasUsage{
				aliasVar:       aliasVar,
				aliasReference: typeNode,
				isTypeNode:     true,
			})
		}
		// For nodes that are referencing the aliased type, set the original
		// name so it can be accessed later in tools such as the CodePrinter or
		// refactoring tools.
		if v.outer.compiler.GetOptions().PreservesDetailedSourceInfo() {
			typeNode.SetOriginalName(&name)
		}
	}
	for child := typeNode.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		v.FixTypeNode(child)
	}
}

// Java: ScopedAliases.isScopedAliases(String)
func ScopedAliases_IsScopedAliases(name string) bool {
	return strings.HasPrefix(name, scopedAliases_SCOPED_ALIASES_PREFIX)
}
