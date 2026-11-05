// Port of com/google/javascript/jscomp/ProcessClosureProvidesAndRequires.java @ closure-compiler v20261028

/*
 * Copyright 2019 The Closure Compiler Authors.
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
)

// ProcessClosureProvidesAndRequires replaces `goog.provide` calls and
// removes goog.{require,requireType,forwardDeclare} calls.
//
// We expect all `goog.modules` and `goog.require`s in modules to have been
// rewritten. This is why all remaining require/requireType calls must refer
// to a `goog.provide`, although the original JS code may contain
// `goog.require`s of a `goog.module`.
//
// This also annotates all provided namespace definitions `a.b = {};` with
// Node#IS_NAMESPACE so that later passes know they refer to a
// `goog.provide`d namespace.
//
// If modules have not been rewritten, this pass also includes legacy Closure
// module namespaces in the list of ProvidedNames.
type ProcessClosureProvidesAndRequires struct {
	compiler   AbstractCompiler
	chunkGraph *JSChunkGraph

	// Track goog.provides in deterministic insertion order.
	providedNames OrderedStringMap[*ProcessClosureProvidesAndRequires_ProvidedName]
	providedLeafs map[string]struct{}

	exportedVars OrderedSet[string]

	// If this is true, rewriting will not remove any goog.provide or
	// goog.require calls
	preserveGoogProvidesAndRequires bool
	requiresToBeRemoved             []rhino.Node
	// Whether this instance has already rewritten goog.provides, which can
	// only happen once
	hasRewritingOccurred    bool
	forwardDeclaresToRemove OrderedSet[rhino.Node]
	astFactory              *AstFactory
}

var ProcessClosureProvidesAndRequires_TYPEDEF_CHILD_OF_PROVIDE = DiagnosticType_Error(
	"JSC_TYPEDEF_CHILD_OF_PROVIDE",
	"invalid @typedef goog.provide {0}\n"+
		"Parent namespace {1} is goog.provided and initialized in the same file")

const processClosureProvidesAndRequires_GOOG = "goog"

// Java: ProcessClosureProvidesAndRequires(AbstractCompiler, boolean)
func NewProcessClosureProvidesAndRequires(
	compiler AbstractCompiler,
	preserveGoogProvidesAndRequires bool) *ProcessClosureProvidesAndRequires {
	return &ProcessClosureProvidesAndRequires{
		compiler:                        compiler,
		chunkGraph:                      compiler.GetChunkGraph(),
		providedNames:                   *NewOrderedStringMapWithCapacity[*ProcessClosureProvidesAndRequires_ProvidedName](64),
		providedLeafs:                   make(map[string]struct{}, 64),
		preserveGoogProvidesAndRequires: preserveGoogProvidesAndRequires,
		astFactory:                      compiler.CreateAstFactory(),
	}
}

func (p *ProcessClosureProvidesAndRequires) putProvidedName(ns string, pn *ProcessClosureProvidesAndRequires_ProvidedName) {
	p.providedNames.Put(ns, pn)
	leaf := ns
	if dot := strings.LastIndexByte(ns, '.'); dot >= 0 {
		leaf = ns[dot+1:]
	}
	p.providedLeafs[leaf] = struct{}{}
}

func (p *ProcessClosureProvidesAndRequires) getProvidedNameByNode(node rhino.Node) *ProcessClosureProvidesAndRequires_ProvidedName {
	switch node.GetToken() {
	case rhino.Token_NAME:
		name := ast.MustName(node).Ident()
		if name == "" {
			return nil
		}
		return p.providedNames.Get(name)
	case rhino.Token_GETPROP:
		if _, ok := p.providedLeafs[ast.MustGetProp(node).Property()]; !ok {
			return nil
		}
	}
	name, ok := node.GetQualifiedNameString()
	if !ok {
		return nil
	}
	return p.providedNames.Get(name)
}

// Process: when invoked as compiler pass, we rewrite all provides and requires.
//
// Java: ProcessClosureProvidesAndRequires.process(Node, Node)
func (p *ProcessClosureProvidesAndRequires) Process(externs rhino.Node, root rhino.Node) {
	p.RewriteProvidesAndRequires(externs, root)
}

// Java: ProcessClosureProvidesAndRequires.getExportedVariableNames()
func (p *ProcessClosureProvidesAndRequires) GetExportedVariableNames() *OrderedSet[string] {
	cloned := p.exportedVars.Clone()
	return &cloned
}

func (p *ProcessClosureProvidesAndRequires) collectProvidedNamesInternal(externs rhino.Node, root rhino.Node) {
	if p.providedNames.IsEmpty() {
		// goog is special-cased because it is provided in Closure's base library.
		p.putProvidedName(
			processClosureProvidesAndRequires_GOOG,
			p.newProvidedName(processClosureProvidesAndRequires_GOOG, rhino.NilNode, nil, false, false, false))
		collector := closureProvideCollector{outer: p}
		NodeTraversal_TraverseRoots(p.compiler, &collector, externs, root)
	}
}

// CollectProvidedNamesSlice collects all `goog.provide`s in deterministic
// insertion order as a flat slice.
func (p *ProcessClosureProvidesAndRequires) CollectProvidedNamesSlice(externs rhino.Node, root rhino.Node) []*ProcessClosureProvidesAndRequires_ProvidedName {
	p.collectProvidedNamesInternal(externs, root)
	return p.providedNames.ValueSlice()
}

// CollectProvidedNames collects all `goog.provide`s in the given namespace
// and warns on invalid code.
//
// Java: ProcessClosureProvidesAndRequires.collectProvidedNames(Node, Node)
func (p *ProcessClosureProvidesAndRequires) CollectProvidedNames(externs rhino.Node, root rhino.Node) *OrderedStringMap[*ProcessClosureProvidesAndRequires_ProvidedName] {
	p.collectProvidedNamesInternal(externs, root)
	return &p.providedNames
}

// RewriteProvidesAndRequires rewrites all provides and requires in the given namespace.
//
// Java: ProcessClosureProvidesAndRequires.rewriteProvidesAndRequires(Node, Node)
func (p *ProcessClosureProvidesAndRequires) RewriteProvidesAndRequires(externs rhino.Node, root rhino.Node) {
	jutil.CheckStateMsg(!p.hasRewritingOccurred, "Cannot call rewriteProvidesAndRequires twice per instance")
	p.hasRewritingOccurred = true

	p.collectProvidedNamesInternal(externs, root)

	for pn := range p.providedNames.Values() {
		pn.replace(p.preserveGoogProvidesAndRequires, &p.providedNames)
	}

	for _, closureRequire := range p.requiresToBeRemoved {
		p.compiler.ReportChangeToEnclosingScope(closureRequire)
		closureRequire.Detach()
	}
	for forwardDeclare := range p.forwardDeclaresToRemove.All() {
		NodeUtil_DeleteNode(forwardDeclare, p.compiler)
	}
}

var closureProvideCollectorVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_CALL)
	s.Add(rhino.Token_ASSIGN)
	s.Add(rhino.Token_NAME)
	s.Add(rhino.Token_EXPR_RESULT)
	return s
}()

type closureProvideCollector struct {
	outer *ProcessClosureProvidesAndRequires
}

var (
	_ NodeTraversal_Callback       = (*closureProvideCollector)(nil)
	_ nodeTraversal_visitTokens    = (*closureProvideCollector)(nil)
	_ nodeTraversal_onlyNameLeaves = (*closureProvideCollector)(nil)
)

func (*closureProvideCollector) visitTokens() *ast.TokenSet {
	return &closureProvideCollectorVisitTokens
}

func (*closureProvideCollector) onlyNameLeaves() {}

func (c *closureProvideCollector) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if n.IsScript() && n.IsFromExterns() && c.outer.providedNames.Len() == 1 && !scriptHasGoogOrModuleBody(n) {
		return false
	}
	// Don't recurse into modules, which cannot have goog.provides. We do need to handle legacy
	// goog.modules, but we do that quickly here rather than descending all the way into them.
	// We completely ignore ES modules and CommonJS modules.
	if (n.IsModuleBody() && parent.GetBooleanProp(rhino.Node_GOOG_MODULE)) ||
		NodeUtil_IsBundledGoogModuleScopeRoot(n) {
		if ns, googModuleStmt, ok := legacyGoogModuleCall(n); ok {
			c.outer.processLegacyModuleCall(ns, googModuleStmt, t.GetChunk())
		}
		return false
	}
	return !n.IsModuleBody()
}

func legacyGoogModuleCall(scopeRoot rhino.Node) (namespace string, googModuleStmt rhino.Node, ok bool) {
	var first ast.Statement
	if modBody, isMod := ast.AsModuleBody(scopeRoot); isMod {
		first, ok = modBody.FirstStatement()
	} else if block, isBlock := ast.AsBlock(scopeRoot); isBlock {
		first, ok = block.FirstStatement()
	}
	if !ok {
		return "", rhino.NilNode, false
	}
	exprStmt, ok := first.AsExprStmt()
	if !ok {
		return "", rhino.NilNode, false
	}
	call, ok := exprStmt.Expr().AsCall()
	if !ok {
		return "", rhino.NilNode, false
	}
	arg0, ok := call.FirstArg()
	if !ok {
		return "", rhino.NilNode, false
	}
	nsLit, ok := arg0.AsStringLit()
	if !ok {
		return "", rhino.NilNode, false
	}
	second, ok := first.Next()
	if !ok || !NodeUtil_IsGoogModuleDeclareLegacyNamespaceCall(second.Node()) {
		return "", rhino.NilNode, false
	}
	return nsLit.Value(), first.Node(), true
}

func (c *closureProvideCollector) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_CALL:
		call := ast.MustCall(n)
		if prop, ok := call.Callee().AsGetProp(); ok && prop.Object().Node().MatchesName(processClosureProvidesAndRequires_GOOG) {
			c.outer.visitGoogMethodCall(t, call, parent, prop.Property())
		}
	case rhino.Token_ASSIGN, rhino.Token_NAME:
		// If this is an assignment to a provided name, remove the provided object.
		c.outer.handleCandidateProvideDefinition(t, n, parent)
	case rhino.Token_EXPR_RESULT:
		c.outer.handleStubDefinition(t, ast.MustExprStmt(n))
	}
}

func (p *ProcessClosureProvidesAndRequires) visitGoogMethodCall(t *NodeTraversal, call ast.Call, parent rhino.Node, methodName string) {
	switch methodName {
	case "exportSymbol":
		// Note: exportSymbol is allowed in local scope
		if firstArg, ok := call.FirstArg(); ok {
			if lit, ok := firstArg.AsStringLit(); ok {
				rootVar, _, _ := strings.Cut(lit.Value(), ".")
				p.exportedVars.Add(rootVar)
			}
		}
	case "require", "requireType":
		if p.isValidPrimitiveCall(t, parent) {
			p.processRequireCall(call, parent)
		}
	case "provide":
		if p.isValidPrimitiveCall(t, parent) {
			p.processProvideCall(t, call, parent)
		}
	case "forwardDeclare":
		if p.isValidPrimitiveCall(t, parent) {
			p.processForwardDeclare(call.Node(), parent)
		}
	}
}

func (p *ProcessClosureProvidesAndRequires) isValidPrimitiveCall(t *NodeTraversal, parent rhino.Node) bool {
	// Ignore invalid primitives if we didn't strip module sugar.
	if p.compiler.GetOptions().ShouldPreserveGoogModule() {
		return true
	}
	return t.InGlobalHoistScope() && parent.IsExprResult()
}

func (p *ProcessClosureProvidesAndRequires) processRequireCall(call ast.Call, parent rhino.Node) {
	if _, ok := closureOnlyStringArg(call); !ok {
		return
	}
	if !p.preserveGoogProvidesAndRequires {
		p.requiresToBeRemoved = append(p.requiresToBeRemoved, parent)
	}
}

func (p *ProcessClosureProvidesAndRequires) processLegacyModuleCall(namespace string, googModuleCall rhino.Node, chunk *JSChunk) {
	p.registerAnyProvidedPrefixes(namespace, googModuleCall, chunk)
	p.putProvidedName(
		namespace,
		p.newProvidedName(namespace, googModuleCall, chunk, true, true, false))
}

func (p *ProcessClosureProvidesAndRequires) processProvideCall(t *NodeTraversal, call ast.Call, parent rhino.Node) {
	ns, ok := closureOnlyStringArg(call)
	if !ok {
		return
	}

	info := NodeUtil_GetBestJSDocInfo(call.Node())
	isImplicitlyInitialized := info != nil && info.IsProvideAlreadyProvided()

	if previouslyProvided := p.providedNames.Get(ns); previouslyProvided != nil {
		if !previouslyProvided.IsExplicitlyProvided() {
			previouslyProvided.addProvide(parent, t.GetChunk(), true /* explicit */, p.chunkGraph)
		}
	} else {
		p.registerAnyProvidedPrefixes(ns, parent, t.GetChunk())
		p.putProvidedName(
			ns,
			p.newProvidedName(ns, parent, t.GetChunk(), true, false, isImplicitlyInitialized))
	}
}

func (p *ProcessClosureProvidesAndRequires) handleStubDefinition(t *NodeTraversal, exprStmt ast.ExprStmt) {
	if !t.InGlobalHoistScope() {
		return
	}
	exprResult := exprStmt.Node()
	var isExternStub bool
	if !t.currentScript.IsNil() {
		isExternStub = t.currentScript.IsFromExterns()
	} else {
		isExternStub = exprResult.IsFromExterns()
	}
	isTypedefStub := isTypedefStubDeclaration(exprResult)
	// recognize @typedefs only if using this pass for typechecking via
	// collectProvidedNames. We don't want rewriting to depend on @typedef annotations.
	isValidTypedefStubDefinition := isTypedefStub && !p.hasRewritingOccurred

	if isValidTypedefStubDefinition || isExternStub {
		if pn := p.getProvidedNameByNode(exprStmt.Expr().Node()); pn != nil {
			pn.addDefinition(exprResult, t.GetChunk(), p.chunkGraph)
		}
	}
	if isTypedefStub {
		p.checkNestedTypedefProvide(exprStmt)
	}
}

func (p *ProcessClosureProvidesAndRequires) checkNestedTypedefProvide(exprStmt ast.ExprStmt) {
	// forbid this pattern:
	//   goog.provide('my.parent');
	//   goog.provide('my.parent.ChildTypedef');
	//   my.parent = [...]; // some initialization, doesn't matter what
	//   /** @typedef {...} */
	//   my.parent.ChildType;
	// at one point, this pattern was supported. now it would produce code that crashes at
	// runtime because the compiler would initialize `my.parent.ChildType = {}` before
	// `my.parent = [...]`, so report an error.

	prop, ok := exprStmt.Expr().AsGetProp()
	if !ok {
		return
	}
	if _, ok := p.providedLeafs[prop.Property()]; !ok {
		return
	}
	name, ok := prop.Node().GetQualifiedNameString()
	if !ok {
		return
	}
	lastDot := strings.LastIndexByte(name, '.')
	if lastDot < 0 || !p.providedNames.ContainsKey(name) {
		return
	}
	parentName := name[:lastDot]
	parent := p.providedNames.Get(parentName)
	if parent == nil {
		return
	}
	parentDefinition := parent.GetCandidateDefinition()
	exprResult := exprStmt.Node()
	if parentDefinition.IsNil() ||
		parentDefinition.GetStaticSourceFile() != exprResult.GetStaticSourceFile() ||
		isTypedefStubDeclaration(parentDefinition) {
		return
	}
	p.compiler.Report(JSError_MakeWithNode(exprResult, ProcessClosureProvidesAndRequires_TYPEDEF_CHILD_OF_PROVIDE, name, parentName))
}

func isTypedefStubDeclaration(statement rhino.Node) bool {
	if !statement.IsExprResult() {
		return false
	}
	info := NodeUtil_GetBestJSDocInfo(statement)
	return info != nil && info.HasTypedefType()
}

func (p *ProcessClosureProvidesAndRequires) handleCandidateProvideDefinition(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if !t.InGlobalHoistScope() {
		return
	}
	var pn *ProcessClosureProvidesAndRequires_ProvidedName
	switch parent.GetToken() {
	case rhino.Token_LET, rhino.Token_CONST:
		if t.InGlobalScope() {
			pn = p.providedNames.Get(ast.MustName(n).Ident())
		}
	case rhino.Token_VAR:
		pn = p.providedNames.Get(ast.MustName(n).Ident())
	case rhino.Token_EXPR_RESULT:
		if assign, ok := ast.AsAssign(n); ok {
			pn = p.getProvidedNameByNode(assign.Target().Node())
		}
	case rhino.Token_CLASS, rhino.Token_FUNCTION:
		// Class and function provides are forbidden; see ProcessClosurePrimitives's
		// CLASS_NAMESPACE_ERROR and FUNCTION_NAMESPACE_ERROR.
	}

	if pn != nil {
		pn.addDefinition(parent, t.GetChunk(), p.chunkGraph)
	}
}

func (p *ProcessClosureProvidesAndRequires) processForwardDeclare(n rhino.Node, parent rhino.Node) {
	if p.preserveGoogProvidesAndRequires {
		return
	}
	if _, ok := p.compiler.GetCodingConvention().IdentifyTypeDeclarationCall(n); ok {
		toRemove := parent.GetParent()
		if parent.IsExprResult() {
			toRemove = parent
		}
		p.forwardDeclaresToRemove.Add(toRemove)
	}
}

func closureOnlyStringArg(call ast.Call) (string, bool) {
	arg, ok := call.OnlyArg()
	if !ok {
		return "", false
	}
	lit, ok := arg.AsStringLit()
	if !ok {
		return "", false
	}
	return lit.Value(), true
}

// registerAnyProvidedPrefixes registers ProvidedNames for prefix namespaces
// if they haven't already been defined, from shortest to longest.
func (p *ProcessClosureProvidesAndRequires) registerAnyProvidedPrefixes(ns string, node rhino.Node, chunk *JSChunk) {
	for start := 0; ; {
		rel := strings.IndexByte(ns[start:], '.')
		if rel < 0 {
			break
		}
		pos := start + rel
		prefixNs := ns[:pos]
		start = pos + 1
		if existing := p.providedNames.Get(prefixNs); existing != nil {
			existing.addProvide(node, chunk, false /* explicit */, p.chunkGraph)
		} else {
			p.putProvidedName(
				prefixNs,
				p.newProvidedName(prefixNs, node, chunk, false, false, false))
		}
	}
}

func (p *ProcessClosureProvidesAndRequires) newProvidedName(
	namespace string,
	node rhino.Node,
	chunk *JSChunk,
	explicit bool,
	fromLegacyModule bool,
	hasImplicitInitialization bool) *ProcessClosureProvidesAndRequires_ProvidedName {
	jutil.CheckArgumentMsg1(
		node.IsNil() ||
			NodeUtil_IsExprCall(node) ||
			(!explicit &&
				(NodeUtil_IsExprAssign(node) ||
					NodeUtil_IsNameDeclaration(node) ||
					isClosureExprQualifiedName(node))),
		"%s",
		node)
	pn := &ProcessClosureProvidesAndRequires_ProvidedName{
		namespace:                 namespace,
		firstNode:                 node,
		firstChunk:                chunk,
		hasImplicitInitialization: hasImplicitInitialization,
		explicitNode:              rhino.NilNode,
		candidateDefinition:       rhino.NilNode,
		replacementNode:           rhino.NilNode,
		fromLegacyModule:          fromLegacyModule,
		compiler:                  p.compiler,
		astFactory:                p.astFactory,
	}
	pn.addProvide(node, chunk, explicit, p.chunkGraph)
	return pn
}

func isClosureExprQualifiedName(node rhino.Node) bool {
	exprStmt, ok := ast.AsExprStmt(node)
	return ok && exprStmt.Expr().Node().IsQualifiedName()
}

// -------------------------------------------------------------------------

// ProcessClosureProvidesAndRequires_ProvidedName stores information about a
// Closure namespace created by a goog.provide.
type ProcessClosureProvidesAndRequires_ProvidedName struct {
	namespace                 string
	firstNode                 rhino.Node
	firstChunk                *JSChunk
	hasImplicitInitialization bool
	explicitNode              rhino.Node
	candidateDefinition       rhino.Node
	minimumChunk              *JSChunk
	replacementNode           rhino.Node
	fromLegacyModule          bool

	compiler   AbstractCompiler
	astFactory *AstFactory
}

func (pn *ProcessClosureProvidesAndRequires_ProvidedName) addProvide(node rhino.Node, chunk *JSChunk, explicit bool, chunkGraph *JSChunkGraph) {
	if explicit {
		jutil.CheckState(pn.explicitNode.IsNil())
		jutil.CheckArgumentMsg1(node.IsExprResult(), "%s", node)
		pn.explicitNode = node
	}
	pn.updateMinimumChunk(chunk, chunkGraph)
}

// Java: ProcessClosureProvidesAndRequires.ProvidedName.isExplicitlyProvided()
func (pn *ProcessClosureProvidesAndRequires_ProvidedName) IsExplicitlyProvided() bool {
	return !pn.explicitNode.IsNil()
}

// Java: ProcessClosureProvidesAndRequires.ProvidedName.hasImplicitInitialization()
func (pn *ProcessClosureProvidesAndRequires_ProvidedName) HasImplicitInitialization() bool {
	return pn.hasImplicitInitialization
}

// Java: ProcessClosureProvidesAndRequires.ProvidedName.isFromLegacyModule()
func (pn *ProcessClosureProvidesAndRequires_ProvidedName) IsFromLegacyModule() bool {
	return pn.fromLegacyModule
}

// Java: ProcessClosureProvidesAndRequires.ProvidedName.hasCandidateDefinition()
func (pn *ProcessClosureProvidesAndRequires_ProvidedName) HasCandidateDefinition() bool {
	return !pn.candidateDefinition.IsNil()
}

// Java: ProcessClosureProvidesAndRequires.ProvidedName.getFirstProvideCall()
func (pn *ProcessClosureProvidesAndRequires_ProvidedName) GetFirstProvideCall() rhino.Node {
	return pn.firstNode
}

// Java: ProcessClosureProvidesAndRequires.ProvidedName.getCandidateDefinition()
func (pn *ProcessClosureProvidesAndRequires_ProvidedName) GetCandidateDefinition() rhino.Node {
	return pn.candidateDefinition
}

// Java: ProcessClosureProvidesAndRequires.ProvidedName.getNamespace()
func (pn *ProcessClosureProvidesAndRequires_ProvidedName) GetNamespace() string {
	return pn.namespace
}

func (pn *ProcessClosureProvidesAndRequires_ProvidedName) addDefinition(node rhino.Node, chunk *JSChunk, chunkGraph *JSChunkGraph) {
	jutil.CheckArgument(
		node.IsExprResult() ||
			node.IsFunction() ||
			NodeUtil_IsNameDeclaration(node))
	jutil.CheckArgument(pn.explicitNode != node)
	if pn.candidateDefinition.IsNil() || !node.IsExprResult() {
		pn.candidateDefinition = node
		pn.updateMinimumChunk(chunk, chunkGraph)
	}
}

func (pn *ProcessClosureProvidesAndRequires_ProvidedName) updateMinimumChunk(newChunk *JSChunk, chunkGraph *JSChunkGraph) {
	if pn.minimumChunk == nil {
		pn.minimumChunk = newChunk
	} else if chunkGraph.GetChunkCount() > 1 {
		pn.minimumChunk = chunkGraph.GetDeepestCommonDependencyInclusive(pn.minimumChunk, newChunk)
	} else {
		jutil.CheckStateMsg(newChunk == pn.minimumChunk, "Missing chunk graph")
	}
}

func (pn *ProcessClosureProvidesAndRequires_ProvidedName) replace(
	preserveGoogProvidesAndRequires bool, providedNames *OrderedStringMap[*ProcessClosureProvidesAndRequires_ProvidedName]) {
	jutil.CheckStateMsg1(
		!pn.IsFromLegacyModule(),
		"Cannot rewrite provides without having rewritten goog.modules, found %s",
		pn.firstNode)
	if pn.firstNode.IsNil() {
		// Don't touch the base case ('goog').
		pn.replacementNode = pn.candidateDefinition
		return
	}

	// Handle the case where there is a duplicate definition for an explicitly
	// provided symbol.
	if pn.HasCandidateDefinition() && !pn.explicitNode.IsNil() {
		pn.replacementNode = pn.candidateDefinition
		if exprStmt, ok := ast.AsExprStmt(pn.candidateDefinition); ok {
			if assign, ok := exprStmt.Expr().AsAssign(); ok {
				if name, ok := assign.Target().AsName(); ok {
					// In the case of a simple name, `name = value;`, ensure the name is declared with `var`.
					pn.convertProvideAssignmentToVarDeclaration(assign, name)
				} else {
					// `some.provided.namespace = value;`
					pn.candidateDefinition.PutBooleanProp(rhino.Node_IS_NAMESPACE, true)
				}
			}
		}
	} else if !pn.hasImplicitInitialization {
		pn.createNamespaceInitialization(
			pn.createDeclarationNode(pn.astFactory.CreateObjectLit()), providedNames)
	}

	// Remove the `goog.provide('a.b.c');` call.
	if !pn.explicitNode.IsNil() {
		if preserveGoogProvidesAndRequires {
			return
		}
		pn.compiler.ReportChangeToEnclosingScope(pn.explicitNode)
		pn.explicitNode.Detach()
	}
}

func (pn *ProcessClosureProvidesAndRequires_ProvidedName) convertProvideAssignmentToVarDeclaration(assign ast.Assign, name ast.Name) {
	nameNode := name.Node()
	valueNode := assign.Value().Node()
	nameNode.Detach()
	valueNode.Detach()

	varNode := rhino.IR_VarWithValue(nameNode, valueNode).Srcref(pn.candidateDefinition)
	varNode.SetJSDocInfo(assign.Node().GetJSDocInfo())
	varNode.PutBooleanProp(rhino.Node_IS_NAMESPACE, true)

	pn.candidateDefinition.ReplaceWith(varNode)
	pn.replacementNode = varNode
	pn.compiler.ReportChangeToEnclosingScope(varNode)
}

func (pn *ProcessClosureProvidesAndRequires_ProvidedName) createNamespaceInitialization(
	replacement rhino.Node, providedNames *OrderedStringMap[*ProcessClosureProvidesAndRequires_ProvidedName]) {
	pn.replacementNode = replacement
	if pn.firstChunk == pn.minimumChunk {
		pn.replacementNode.InsertBefore(pn.firstNode)
	} else {
		// In this case, the name was implicitly provided by two independent
		// modules. We need to move this code up to a common module.
		indexOfDot := strings.LastIndexByte(pn.namespace, '.')
		if indexOfDot == -1 {
			pn.compiler.GetNodeForCodeInsertion(pn.minimumChunk).AddChildToBack(pn.replacementNode)
		} else {
			parentName := providedNames.Get(pn.namespace[:indexOfDot])
			jutil.CheckNotNull(parentName)
			rhino.CheckNotNull(parentName.replacementNode)
			pn.replacementNode.InsertAfter(parentName.replacementNode)
		}
	}
	pn.compiler.ReportChangeToEnclosingScope(pn.replacementNode)
}

func (pn *ProcessClosureProvidesAndRequires_ProvidedName) createDeclarationNode(value rhino.Node) rhino.Node {
	jutil.CheckArgumentMsg1(value.IsObjectLit() || value.IsCast(), "%s", value)
	if strings.IndexByte(pn.namespace, '.') == -1 {
		return pn.makeVarDeclNode(value)
	}
	return pn.makeAssignmentExprNode(value)
}

func (pn *ProcessClosureProvidesAndRequires_ProvidedName) makeVarDeclNode(value rhino.Node) rhino.Node {
	name := rhino.IR_Name(pn.namespace)
	name.AddChildToFront(value)

	decl := rhino.IR_Var(name)
	decl.PutBooleanProp(rhino.Node_IS_NAMESPACE, true)

	if pn.compiler.GetCodingConvention().IsConstant(pn.namespace) {
		name.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
	}
	if !pn.HasCandidateDefinition() {
		decl.SetJSDocInfo(NodeUtil_CreateConstantJsDoc())
	}

	jutil.CheckState(isNamespacePlaceholder(decl))
	pn.setSourceInfo(decl)
	return decl
}

func (pn *ProcessClosureProvidesAndRequires_ProvidedName) makeAssignmentExprNode(value rhino.Node) rhino.Node {
	lhs := pn.astFactory.CreateQNameWithUnknownType(pn.namespace).SrcrefTree(pn.firstNode)
	assign := pn.astFactory.CreateAssign(lhs, value)
	if !pn.HasCandidateDefinition() {
		assign.SetJSDocInfo(NodeUtil_CreateConstantJsDoc())
	}
	decl := rhino.IR_ExprResult(assign)
	decl.PutBooleanProp(rhino.Node_IS_NAMESPACE, true)
	jutil.CheckState(isNamespacePlaceholder(decl))
	pn.setSourceInfo(decl)
	ast.MustGetProp(lhs).Object().Node().MakeNonIndexableRecursive()
	return decl
}

func (pn *ProcessClosureProvidesAndRequires_ProvidedName) setSourceInfo(newNode rhino.Node) {
	sourceInfoNode := pn.firstNode

	provideStringNode := pn.getProvideStringNode()
	if !provideStringNode.IsNil() {
		// Given namespace "foo.bar.baz" we create node for "baz" here and need to calculate
		// UTF-16 length and start of the last component which is "baz".
		firstCharIndex := jutil.StringLastIndexOf(pn.namespace, ".") + 1 // If no dots, then 0.

		sourceInfoNode = provideStringNode.CloneNode()
		sourceInfoNode.SetLinenoCharno(
			sourceInfoNode.GetLineno(),
			sourceInfoNode.GetCharno()+firstCharIndex+1) // +1 for quote
		sourceInfoNode.SetLength(jutil.StringLength(pn.namespace) - firstCharIndex)
	}

	newNode.SrcrefTree(sourceInfoNode)
}

func (pn *ProcessClosureProvidesAndRequires_ProvidedName) getProvideStringNode() rhino.Node {
	if exprStmt, ok := ast.AsExprStmt(pn.firstNode); ok {
		if call, ok := exprStmt.Expr().AsCall(); ok {
			if arg, ok := call.LastArg(); ok {
				return arg.Node()
			}
			return call.Callee().Node()
		}
	}
	return rhino.NilNode
}

// Java: ProcessClosureProvidesAndRequires.ProvidedName.toString()
func (pn *ProcessClosureProvidesAndRequires_ProvidedName) String() string {
	explicitOrImplicit := "implicit"
	if pn.IsExplicitlyProvided() {
		explicitOrImplicit = "explicit"
	}
	return "ProvidedName: " + pn.namespace + ", " + explicitOrImplicit
}

func isNamespacePlaceholder(n rhino.Node) bool {
	if !n.GetBooleanProp(rhino.Node_IS_NAMESPACE) {
		return false
	}

	var value ast.Expr
	var ok bool
	if exprStmt, isExpr := ast.AsExprStmt(n); isExpr {
		if assign, isAssign := exprStmt.Expr().AsAssign(); isAssign {
			value, ok = assign.Value(), true
		}
	} else if decl, isVar := ast.AsNameDeclarationOf(n, rhino.Token_VAR); isVar {
		if first, hasFirst := decl.FirstDeclarator(); hasFirst {
			if name, isName := first.AsName(); isName {
				value, ok = name.Init()
			}
		}
	}
	if !ok {
		return false
	}
	if cast, isCast := value.AsCast(); isCast {
		value = cast.Expr()
	}
	obj, isObj := value.AsObjectLit()
	return isObj && !obj.HasMembers()
}
