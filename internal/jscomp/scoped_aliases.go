// Port of com/google/javascript/jscomp/ScopedAliases.java @ closure-compiler v20261028

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

	// Traversal state collected during Process.
	aliasDefinitionsToDelete []rhino.Node
	scopeCalls               []rhino.Node
	aliasUsages              []aliasUsage

	// Temporary per-scope maps cleared when exiting a goog.scope body.
	aliasByID     map[uint32]*Var
	injectedDecls map[rhino.Node]struct{}

	// Persists across scopes within a single Process run.
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

	matchingScripts map[rhino.Node]struct{}

	// The body of the function that is passed to goog.scope.
	// Set when the traversal enters the body, and set back to NilNode when it exits.
	scopeFunctionBody rhino.Node
	uniquifier        MakeDeclaredNamesUnique
}

var (
	_ CompilerPass                 = (*ScopedAliases)(nil)
	_ NodeTraversal_ScopedCallback = (*ScopedAliases)(nil)
)

// ScopedAliases_SCOPING_METHOD_NAME is the name used to denote an scoped
// function block used for aliasing.
const ScopedAliases_SCOPING_METHOD_NAME = "goog.scope"

const scopedAliasesPrefix = "$jscomp$scope$"

// LINT.IfChange
const missingAliasesPrefix = "$$jscomp$missingAlias$"

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

var (
	googModuleGetQName = rhino.QualifiedName_Of("goog.module.get")
	googScopeQName     = rhino.QualifiedName_Of(ScopedAliases_SCOPING_METHOD_NAME)
)

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
	moduleMetadataMap *ModuleMetadataMap,
	closureNamespaces map[string]struct{},
	invalidModuleGetHandling ScopedAliases_InvalidModuleGetHandling) *ScopedAliases {
	if moduleMetadataMap == nil {
		moduleMetadataMap = compiler.GetModuleMetadataMap()
	}
	return &ScopedAliases{
		compiler:                 compiler,
		preprocessorSymbolTable:  preprocessorSymbolTable,
		moduleMetadataMap:        moduleMetadataMap,
		closureNamespaces:        closureNamespaces,
		invalidModuleGetHandling: invalidModuleGetHandling,
	}
}

// ScopedAliasesOptions configures a ScopedAliases pass.
type ScopedAliasesOptions struct {
	PreprocessorSymbolTable  *PreprocessorSymbolTable
	ModuleMetadataMap        *ModuleMetadataMap
	InvalidModuleGetHandling ScopedAliases_InvalidModuleGetHandling
}

// NewScopedAliases constructs a ScopedAliases pass.
func NewScopedAliases(compiler AbstractCompiler, opts ScopedAliasesOptions) *ScopedAliases {
	var closureNamespaces map[string]struct{}
	if opts.ModuleMetadataMap != nil {
		entries := opts.ModuleMetadataMap.GetModulesByGoogNamespace().Entries()
		if len(entries) > 0 {
			closureNamespaces = make(map[string]struct{}, len(entries))
			for _, e := range entries {
				closureNamespaces[e.Key] = struct{}{}
			}
		}
	}
	return newScopedAliases(
		compiler,
		opts.PreprocessorSymbolTable,
		compiler.GetModuleMetadataMap(),
		closureNamespaces,
		opts.InvalidModuleGetHandling)
}

func (s *ScopedAliases) resetTraversalState() {
	clear(s.aliasDefinitionsToDelete)
	s.aliasDefinitionsToDelete = s.aliasDefinitionsToDelete[:0]
	clear(s.scopeCalls)
	s.scopeCalls = s.scopeCalls[:0]
	clear(s.aliasUsages)
	s.aliasUsages = s.aliasUsages[:0]
	clear(s.aliasByID)
	clear(s.injectedDecls)
	clear(s.deletedAliasVars)
	clear(s.forbiddenLocals)
	s.hasNamespaceShadows = false
	s.hasErrors = false
	s.matchingScripts = nil
	s.scopeFunctionBody = rhino.NilNode
	s.uniqueIdInput = nil
}

// Java: ScopedAliases.process(Node, Node)
func (s *ScopedAliases) Process(externs rhino.Node, root rhino.Node) {
	s.resetTraversalState()
	defer s.resetTraversalState()

	if scripts, ok := collectRootScripts(root); ok {
		matchingScripts, prefiltered := filterMatchingScriptsSlice(s.compiler, scripts, subtreeContainsGoogScopeCall)
		if prefiltered {
			if len(matchingScripts) == 0 {
				return
			}
			s.matchingScripts = matchingScripts
		}
	}

	NodeTraversal_Traverse(s.compiler, root, s)

	if !s.hasErrors {
		// Apply the aliases.
		aliasWorkQueue := s.aliasUsages
		for len(aliasWorkQueue) != 0 {
			n := 0
			for _, u := range aliasWorkQueue {
				if u.referencesOtherAlias(s.deletedAliasVars) {
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
		for _, aliasDefinition := range s.aliasDefinitionsToDelete {
			s.compiler.ReportChangeToEnclosingScope(aliasDefinition)
			if decl, ok := ast.AsNameDeclaration(aliasDefinition.GetParent()); ok && decl.NumDeclarators() == 1 {
				decl.Detach()
			} else {
				aliasDefinition.Detach()
			}
		}

		// Collapse the scopes.
		for _, scopeCall := range s.scopeCalls {
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

func aliasedRootID(rhs rhino.Node, intern bool) (rootID uint32, isQualified bool) {
	cur := rhs
	for gp, ok := ast.AsGetProp(cur); ok; gp, ok = ast.AsGetProp(cur) {
		isQualified = true
		cur = gp.Object().Node()
	}
	switch cur.GetToken() {
	case rhino.Token_NAME:
		return cur.GetStringID(), isQualified
	case rhino.Token_CALL:
		jutil.CheckStateMsg1(NodeUtil_IsCallToQualifiedName(cur, googModuleGetQName), "%s", cur)
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
	rootID, _ := aliasedRootID(u.aliasVar.GetInitialValue(), false)
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
		if strings.HasPrefix(typeName, scopedAliasesPrefix) {
			// Already visited.
			return
		}
		aliasExpanded := getAliasedNamespace(aliasDefinition)
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
func isValidAliasRhs(rhs rhino.Node) bool {
	switch rhs.GetToken() {
	case rhino.Token_GETPROP:
		return isValidAliasRhs(ast.MustGetProp(rhs).Object().Node())
	case rhino.Token_NAME:
		return true
	case rhino.Token_CALL:
		return NodeUtil_IsCallToQualifiedName(rhs, googModuleGetQName)
	default:
		return false
	}
}

// Java: ScopedAliases.isAliasDefinition(Node)
func isAliasDefinition(nameNode rhino.Node) bool {
	initExpr, ok := ast.MustName(nameNode).Init()
	return ok && isValidAliasRhs(initExpr.Node())
}

// Java: ScopedAliases.getAliasedNamespace(Node)
func getAliasedNamespace(rhs rhino.Node) string {
	switch rhs.GetToken() {
	case rhino.Token_GETPROP:
		gp := ast.MustGetProp(rhs)
		return getAliasedNamespace(gp.Object().Node()) + "." + gp.Property()
	case rhino.Token_NAME:
		return ast.MustName(rhs).Ident()
	case rhino.Token_CALL:
		jutil.CheckStateMsg1(NodeUtil_IsCallToQualifiedName(rhs, googModuleGetQName), "%s", rhs)
		arg, ok := ast.MustCall(rhs).OnlyArg()
		jutil.CheckStateMsg1(ok, "%s", rhs)
		return arg.Node().GetString()
	default:
		panic(jutil.NewRuntimeException("Invalid alias RHS:" + rhs.String()))
	}
}

// inGoogScopeBody returns true if this NodeTraversal is currently within a
// goog.scope function body.
//
// Java: ScopedAliases.Traversal.inGoogScopeBody()
func (s *ScopedAliases) inGoogScopeBody() bool {
	return !s.scopeFunctionBody.IsNil()
}

// isGoogScopeFunctionBody returns true if n is the goog.scope function body.
//
// Java: ScopedAliases.Traversal.isGoogScopeFunctionBody(Node)
func (s *ScopedAliases) isGoogScopeFunctionBody(n rhino.Node) bool {
	return s.inGoogScopeBody() && n == s.scopeFunctionBody
}

func (*ScopedAliases) shouldTraverseOnlyScriptAndFunction() bool { return true }

func subtreeContainsGoogScopeCall(root rhino.Node) bool {
	for cur := root; !cur.IsNil(); {
		if cur.IsCall() && googScopeQName.Matches(cur.GetFirstChild()) {
			return true
		}
		if child := cur.GetFirstChild(); !child.IsNil() {
			cur = child
			continue
		}
		for cur != root {
			if next := cur.GetNext(); !next.IsNil() {
				cur = next
				break
			}
			cur = cur.GetParent()
		}
		if cur == root {
			break
		}
	}
	return false
}

// Java: ScopedAliases.Traversal.shouldTraverse(NodeTraversal, Node, Node)
func (s *ScopedAliases) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if n.IsScript() && !s.inGoogScopeBody() {
		if s.matchingScripts != nil {
			_, ok := s.matchingScripts[n]
			return ok
		}
		return subtreeContainsGoogScopeCall(n)
	}
	return true
}

// Java: ScopedAliases.Traversal.isCallToScopeMethod(Node)
func isCallToScopeMethod(n rhino.Node) bool {
	call, ok := ast.AsCall(n)
	return ok && googScopeQName.Matches(call.Callee().Node())
}

// findScopeMethodCall: scopeRoot is the Node which is the root of the
// current scope. Returns the goog.scope() CALL node containing the
// scopeRoot, or NilNode if scopeRoot is not in a goog.scope() call.
//
// Java: ScopedAliases.Traversal.findScopeMethodCall(Node)
func findScopeMethodCall(scopeRoot rhino.Node) rhino.Node {
	n := scopeRoot.GetGrandparent()
	if isCallToScopeMethod(n) {
		return n
	}
	return rhino.NilNode
}

// Java: ScopedAliases.Traversal.enterScope(NodeTraversal)
func (s *ScopedAliases) EnterScope(t *NodeTraversal) {
	if t.InGlobalHoistScope() {
		return
	}
	scopeMethodCall := findScopeMethodCall(t.GetScopeRoot())
	if !scopeMethodCall.IsNil() {
		if s.forbiddenLocals == nil {
			s.forbiddenLocals = map[strpool.ID]struct{}{scopedAliasesJscompID: {}}
		} else {
			s.forbiddenLocals[scopedAliasesJscompID] = struct{}{}
		}
		s.uniqueIdInput = t.GetInput()
		s.findAliases(t.GetScope())
		if lastArg, ok := ast.MustCall(scopeMethodCall).LastArg(); ok {
			if fn, ok := ast.AsFunction(lastArg.Node()); ok {
				s.scopeFunctionBody = fn.Body().Node()
			} else {
				s.scopeFunctionBody = lastArg.Node().GetLastChild()
			}
		}
	}
}

// Java: ScopedAliases.Traversal.exitScope(NodeTraversal)
func (s *ScopedAliases) ExitScope(t *NodeTraversal) {
	if s.isGoogScopeFunctionBody(t.GetScopeRoot()) {
		s.scopeFunctionBody = rhino.NilNode
		s.renameNamespaceShadows(t)
		clear(s.injectedDecls)
		clear(s.aliasByID)
		clear(s.forbiddenLocals)
		s.hasNamespaceShadows = false
	} else if s.inGoogScopeBody() {
		// Called on inner scopes within a goog.scope, including both block
		// scopes and function scopes.
		s.findNamespaceShadows(t)
		s.reportInvalidVariables(t)
	}
}

// Java: ScopedAliases.Traversal.reportInvalidVariables(NodeTraversal)
func (s *ScopedAliases) reportInvalidVariables(t *NodeTraversal) {
	scopeRoot := t.GetScopeRoot()
	if enclosingFn, ok := ast.AsFunction(t.GetEnclosingFunction()); ok &&
		s.isGoogScopeFunctionBody(enclosingFn.Body().Node()) &&
		scopeRoot.IsBlock() &&
		!scopeRoot.GetParent().IsFunction() {
		for var_ := range t.GetScope().GetVarIterable() {
			parent := var_.GetNameNode().GetParent()
			if NodeUtil_IsFunctionDeclaration(parent) {
				// Disallow block-scoped function declarations that leak into the
				// goog.scope function body. Technically they shouldn't leak in ES6
				// but the browsers don't agree on that yet.
				s.report(var_.GetNode(), ScopedAliases_GOOG_SCOPE_INVALID_VARIABLE, var_.GetName())
			}
		}
	}
}

// Java: ScopedAliases.Traversal.report(Node, DiagnosticType, String...)
func (s *ScopedAliases) report(n rhino.Node, error_ *DiagnosticType, arguments ...string) {
	s.compiler.Report(JSError_MakeWithNode(n, error_, arguments...))
	s.hasErrors = true
}

// Java: ScopedAliases.Traversal.findAliases(Scope)
func (s *ScopedAliases) findAliases(scope *Scope) {
	compiler := s.compiler
	for var_ := range scope.GetVarIterable() {
		n := var_.GetNode()
		parent := n.GetParent()
		// We use isBlock to avoid variables declared in loop headers.
		isVar := NodeUtil_IsNameDeclaration(parent) && parent.GetParent().IsBlock()
		isFunctionDecl := NodeUtil_IsFunctionDeclaration(parent)
		if isVar && isAliasDefinition(n) {
			s.recordAlias(var_)
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
			uniqueId := compiler.GetUniqueIdSupplier().GetUniqueId(s.uniqueIdInput)

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
			newDecl.SetOriginalNameString(name)
			if isHoisted {
				grandparent.AddChildToFront(newDecl)
			} else {
				newDecl.InsertBefore(varNode)
			}
			if !value.IsNil() {
				compiler.ReportChangeToEnclosingScope(newDecl)
			}
			if s.injectedDecls == nil {
				s.injectedDecls = make(map[rhino.Node]struct{})
			}
			s.injectedDecls[newDecl] = struct{}{}

			// Rewrite "var name = EXPR;" to "var name = $jscomp$scope$uniqueId$name;"
			var_.GetNameNode().AddChildToFront(NodeUtil_NewNameWithOriginalName(compiler, globalName, n, name))
			s.recordAlias(var_)
		} else {
			// Do not allow other kinds of local symbols, like catch params.
			s.report(n, ScopedAliases_GOOG_SCOPE_NON_ALIAS_LOCAL, n.GetString())
		}
	}
}

// Java: ScopedAliases.Traversal.recordAlias(Var)
func (s *ScopedAliases) recordAlias(aliasVar *Var) {
	initialValue := aliasVar.GetInitialValue()
	nameNode := aliasVar.GetNameNode()
	s.aliasDefinitionsToDelete = append(s.aliasDefinitionsToDelete, nameNode)

	if s.aliasByID == nil {
		s.aliasByID = make(map[uint32]*Var)
	}

	if s.invalidModuleGetHandling.ShouldGiveUniqueName() &&
		s.containsInvalidGoogModuleGet(initialValue) {
		// Create a new alias with a unique name, to replace the existing one.
		// Needed for integration with Clutz only: by giving the alias a unique
		// name, Clutz may then map back from the unique name to the original
		// `const alias = goog.module.get(` call.
		name := aliasVar.GetName()
		var metadata *ModuleMetadataMap_ModuleMetadata
		if s.moduleMetadataMap != nil {
			if sf := initialValue.GetStaticSourceFile(); sf != nil {
				metadata = s.moduleMetadataMap.ModuleByPath(sf.GetName())
			}
		}
		// LINT.IfChange
		suffix := "unknownScript"
		if metadata != nil {
			suffix = metadata.Path().ToModuleName()
		}
		newName := name + missingAliasesPrefix + suffix
		// LINT.ThenChange(//depot/google3/third_party/java_src/clutz/src/main/java/com/google/javascript/clutz/imports/ImportRenameMapBuilder.java)
		initialValue.ReplaceWith(rhino.IR_Name(newName))
		if s.deletedAliasVars == nil {
			s.deletedAliasVars = make(map[*Var]struct{})
		}
		s.deletedAliasVars[aliasVar] = struct{}{}
		s.aliasByID[nameNode.GetStringID()] = aliasVar
		return
	}

	s.aliasByID[nameNode.GetStringID()] = aliasVar

	if rootID, isQualified := aliasedRootID(initialValue, true); isQualified && rootID != 0 {
		if _, ok := s.aliasByID[rootID]; !ok {
			s.forbiddenLocals[strpool.ID(rootID)] = struct{}{}
		}
	}
}

// containsInvalidGoogModuleGet returns whether the rhs contains any
// goog.module.get calls to inexistent namespaces.
//
// Java: ScopedAliases.Traversal.containsInvalidGoogModuleGet(Node)
func (s *ScopedAliases) containsInvalidGoogModuleGet(expression rhino.Node) bool {
	switch expression.GetToken() {
	case rhino.Token_NAME:
		return false
	case rhino.Token_GETPROP:
		return s.containsInvalidGoogModuleGet(ast.MustGetProp(expression).Object().Node())
	case rhino.Token_CALL:
		namespace := ast.MustCall(expression).MustFirstArg().Node().GetString()
		_, ok := s.closureNamespaces[namespace]
		return !ok
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Unrecognized alias rhs " + expression.String()))
	}
}

// findNamespaceShadows finds out if there are any local shadows of
// namespaces.
//
// Java: ScopedAliases.Traversal.findNamespaceShadows(NodeTraversal)
func (s *ScopedAliases) findNamespaceShadows(t *NodeTraversal) {
	if s.hasNamespaceShadows {
		return
	}

	scope := t.GetScope()
	for var_ := range scope.GetVarIterable() {
		if _, ok := s.forbiddenLocals[strpool.ID(var_.nameID)]; ok {
			s.hasNamespaceShadows = true
			return
		}
	}
}

// renameNamespaceShadows renames any local shadows of namespaces. This should
// be a very rare occurrence, so only do this traversal if we know that we
// need it.
//
// Java: ScopedAliases.Traversal.renameNamespaceShadows(NodeTraversal)
func (s *ScopedAliases) renameNamespaceShadows(t *NodeTraversal) {
	jutil.CheckStateMsg1(NodeUtil_IsFunctionBlock(t.GetScopeRoot()), "%s", t.GetScopeRoot())

	if s.hasNamespaceShadows {
		renamer := &MakeDeclaredNamesUnique_TargettedRenamer{
			delegate:   NewMakeDeclaredNamesUnique_ContextualRenamer(),
			targetIDs:  s.forbiddenLocals,
			cachedSize: len(s.forbiddenLocals),
		}
		// order-independent: root ContextualRenamer.addDeclaredNameID sets each independent nameUsage[id] = 1
		for id := range s.forbiddenLocals {
			renamer.addDeclaredNameID(id, false)
		}
		s.uniquifier.reset(renamer, true, false)
		NodeTraversal_TraverseScopeRoots(
			s.compiler, []rhino.Node{t.GetScopeRoot()}, &s.uniquifier, true)
	}
}

var scopedAliasesRenameBleedingTokens = makeVisitTokenSet(rhino.Token_NAME)

// renameBleedingFunctionNameCallback is the callback of
// Traversal.renameBleedingFunctionName.
type renameBleedingFunctionNameCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler
	name     string
	nameID   uint32
	suffix   string
	fnName   rhino.Node
}

var (
	_ nodeTraversal_visitTokens    = (*renameBleedingFunctionNameCallback)(nil)
	_ nodeTraversal_onlyNameLeaves = (*renameBleedingFunctionNameCallback)(nil)
)

func (*renameBleedingFunctionNameCallback) visitTokens() *ast.TokenSet {
	return &scopedAliasesRenameBleedingTokens
}

func (*renameBleedingFunctionNameCallback) onlyNameLeaves() {}

// Java: ScopedAliases.Traversal.renameBleedingFunctionName(NodeTraversal, Node).new AbstractPostOrderCallback() {...}.visit(NodeTraversal, Node, Node)
func (c *renameBleedingFunctionNameCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if name, ok := ast.AsName(n); ok &&
		n.GetStringID() == c.nameID &&
		t.GetScope().GetVarByID(c.nameID).GetNode() == c.fnName {
		name.SetIdent(c.name + "$jscomp$scopedAliases$" + c.suffix)
		c.compiler.ReportChangeToEnclosingScope(n)
	}
}

// Java: ScopedAliases.Traversal.renameBleedingFunctionName(NodeTraversal, Node)
func (s *ScopedAliases) renameBleedingFunctionName(t *NodeTraversal, fnName rhino.Node) {
	compiler := s.compiler
	name := fnName.GetString()
	suffix := compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())
	cb := &renameBleedingFunctionNameCallback{
		compiler: compiler, name: name, nameID: fnName.GetStringID(), suffix: suffix, fnName: fnName}
	NodeTraversal_Builder{
		compiler:     compiler,
		callback:     cb,
		scopeCreator: t.GetScopeCreator(),
	}.TraverseAtScope(t.GetScope())
	fnName.SetString(name + "$jscomp$scopedAliases$" + suffix)
}

// Java: ScopedAliases.Traversal.validateScopeCall(NodeTraversal, Node, Node)
func (s *ScopedAliases) validateScopeCall(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	call := ast.MustCall(n)
	if s.preprocessorSymbolTable != nil {
		s.preprocessorSymbolTable.AddReference(call.Callee().Node())
	}
	if !parent.IsExprResult() {
		s.report(n, ScopedAliases_GOOG_SCOPE_MUST_BE_ALONE)
	}
	if !t.GetEnclosingFunction().IsNil() {
		s.report(n, ScopedAliases_GOOG_SCOPE_MUST_BE_IN_GLOBAL_SCOPE)
	}
	arg, ok := call.OnlyArg()
	if !ok {
		// The goog.scope call should have exactly 1 parameter.
		s.report(n, ScopedAliases_GOOG_SCOPE_HAS_BAD_PARAMETERS)
	} else {
		anonymousFnNode := arg.Node()
		if fn, isFn := ast.AsFunction(anonymousFnNode); !isFn || fn.Params().HasParams() {
			s.report(anonymousFnNode, ScopedAliases_GOOG_SCOPE_HAS_BAD_PARAMETERS)
		} else if _, hasName := NodeUtil_GetNameString(anonymousFnNode); hasName {
			s.report(anonymousFnNode, ScopedAliases_GOOG_SCOPE_HAS_BAD_PARAMETERS)
		} else {
			s.scopeCalls = append(s.scopeCalls, n)
		}
	}
}

// Java: ScopedAliases.Traversal.visit(NodeTraversal, Node, Node)
func (s *ScopedAliases) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if !s.inGoogScopeBody() && !n.IsCall() {
		return
	}
	if n.IsCall() && isCallToScopeMethod(n) {
		s.validateScopeCall(t, n, n.GetParent())
	}
	if !s.inGoogScopeBody() {
		return
	}
	compiler := s.compiler

	type_ := n.GetToken()
	var aliasVar *Var
	if type_ == rhino.Token_NAME {
		id := n.GetStringID()
		if candidate := s.aliasByID[id]; candidate != nil {
			if t.GetScope().GetVarByID(id) == candidate {
				aliasVar = candidate
				// For nodes that are referencing the aliased type, set the original
				// name so it can be accessed later in tools such as the CodePrinter
				// or refactoring tools.
				if compiler.GetOptions().PreservesDetailedSourceInfo() {
					n.SetOriginalNameFromName(n)
				}
			}
		}
	}

	if enclosingFn, ok := ast.AsFunction(t.GetEnclosingFunction()); ok &&
		s.isGoogScopeFunctionBody(enclosingFn.Body().Node()) {
		if aliasVar != nil && NodeUtil_IsLValue(n) {
			if aliasVar.GetNode() == n {
				// Return early, to ensure that we don't record this as an alias
				// usage.
				return
			}
			s.report(n, ScopedAliases_GOOG_SCOPE_ALIAS_REDEFINED, n.GetString())
		}

		switch type_ {
		case rhino.Token_RETURN:
			s.report(n, ScopedAliases_GOOG_SCOPE_USES_RETURN)
		case rhino.Token_THIS:
			s.report(n, ScopedAliases_GOOG_SCOPE_REFERENCES_THIS)
		case rhino.Token_THROW:
			s.report(n, ScopedAliases_GOOG_SCOPE_USES_THROW)
		default:
		}
	}

	// If this is a bleeding function expression, like
	// var x = function y() { ... }
	// then old versions of IE declare "y" in the current scope. We don't
	// want the scope unboxing to add "y" to the global scope, so we
	// need to rename it.
	if NodeUtil_IsBleedingFunctionName(n) {
		s.renameBleedingFunctionName(t, n)
	}

	// Check if this name points to an alias.
	if aliasVar != nil {
		// Note, to support the transitive case, it's important we don't
		// clone aliasedNode here.
		s.aliasUsages = append(s.aliasUsages, aliasUsage{aliasVar: aliasVar, aliasReference: n})
	}

	// When we inject declarations, we duplicate jsdoc. Make sure
	// we only process that jsdoc once.
	if info := n.GetJSDocInfo(); info != nil {
		if _, injected := s.injectedDecls[n]; !injected {
			for it := info.TypeExprIter(); it.Next(); {
				s.fixTypeNode(it.Node())
			}
		}
	}
}

// Java: ScopedAliases.Traversal.fixTypeNode(Node)
func (s *ScopedAliases) fixTypeNode(typeNode rhino.Node) {
	if strLit, ok := ast.AsStringLit(typeNode); ok {
		name := strLit.Value()
		var baseID uint32
		if endIndex := strings.IndexByte(name, '.'); endIndex != -1 {
			baseID = uint32(strpool.Lookup(name[:endIndex]))
		} else {
			baseID = typeNode.GetStringID()
		}
		if aliasVar := s.aliasByID[baseID]; aliasVar != nil {
			s.aliasUsages = append(s.aliasUsages, aliasUsage{
				aliasVar:       aliasVar,
				aliasReference: typeNode,
				isTypeNode:     true,
			})
		}
		// For nodes that are referencing the aliased type, set the original
		// name so it can be accessed later in tools such as the CodePrinter or
		// refactoring tools.
		if s.compiler.GetOptions().PreservesDetailedSourceInfo() {
			typeNode.SetOriginalNameString(name)
		}
	}
	for child := typeNode.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		s.fixTypeNode(child)
	}
}

// Java: ScopedAliases.isScopedAliases(String)
func ScopedAliases_IsScopedAliases(name string) bool {
	return strings.HasPrefix(name, scopedAliasesPrefix)
}
