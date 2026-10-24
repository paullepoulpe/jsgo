// Port of com/google/javascript/jscomp/ProcessClosureProvidesAndRequires.java @ closure-compiler v20261021

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

	// Use a LinkedHashMap because the goog.provides must be processed in a
	// deterministic order.
	providedNames *jutil.LinkedHashMap[string, *ProcessClosureProvidesAndRequires_ProvidedName]
	providedLeafs map[string]struct{}

	exportedVariables *jutil.LinkedHashSet[string]

	// If this is true, rewriting will not remove any goog.provide or
	// goog.require calls
	preserveGoogProvidesAndRequires bool
	requiresToBeRemoved             *jutil.ArrayList[*rhino.Node]
	// Whether this instance has already rewritten goog.provides, which can
	// only happen once
	hasRewritingOccurred    bool
	forwardDeclaresToRemove *jutil.LinkedHashSet[*rhino.Node]
	astFactory              *AstFactory
}

var ProcessClosureProvidesAndRequires_TYPEDEF_CHILD_OF_PROVIDE = DiagnosticType_Error(
	"JSC_TYPEDEF_CHILD_OF_PROVIDE",
	"invalid @typedef goog.provide {0}\n"+
		"Parent namespace {1} is goog.provided and initialized in the same file")

// The root Closure namespace
const processClosureProvidesAndRequires_GOOG = "goog"

// Java: ProcessClosureProvidesAndRequires(AbstractCompiler, boolean)
func NewProcessClosureProvidesAndRequires(
	compiler AbstractCompiler,
	preserveGoogProvidesAndRequires bool) *ProcessClosureProvidesAndRequires {
	p := &ProcessClosureProvidesAndRequires{}
	p.providedNames = jutil.NewLinkedHashMap[string, *ProcessClosureProvidesAndRequires_ProvidedName]()
	p.providedLeafs = make(map[string]struct{}, 64)
	p.exportedVariables = jutil.NewLinkedHashSet[string]()
	p.requiresToBeRemoved = jutil.NewArrayList[*rhino.Node]()
	p.hasRewritingOccurred = false
	p.forwardDeclaresToRemove = jutil.NewLinkedHashSet[*rhino.Node]()
	p.compiler = compiler
	p.chunkGraph = compiler.GetChunkGraph()
	p.preserveGoogProvidesAndRequires = preserveGoogProvidesAndRequires
	p.astFactory = compiler.CreateAstFactory()
	return p
}

func (p *ProcessClosureProvidesAndRequires) putProvidedName(ns string, pn *ProcessClosureProvidesAndRequires_ProvidedName) {
	p.providedNames.Put(ns, pn)
	leaf := ns
	if dot := strings.LastIndexByte(ns, '.'); dot >= 0 {
		leaf = ns[dot+1:]
	}
	p.providedLeafs[leaf] = struct{}{}
}

func (p *ProcessClosureProvidesAndRequires) getProvidedNameByNode(node *rhino.Node) *ProcessClosureProvidesAndRequires_ProvidedName {
	switch node.GetToken() {
	case rhino.Token_NAME:
		name := node.GetString()
		if name == "" {
			return nil
		}
		return p.providedNames.Get(name)
	case rhino.Token_GETPROP:
		if _, ok := p.providedLeafs[node.GetString()]; !ok {
			return nil
		}
		name, ok := node.GetQualifiedNameString()
		if !ok {
			return nil
		}
		return p.providedNames.Get(name)
	default:
		name, ok := node.GetQualifiedNameString()
		if !ok {
			return nil
		}
		return p.providedNames.Get(name)
	}
}

// Process: when invoked as compiler pass, we rewrite all provides and
// requires.
//
// Java: ProcessClosureProvidesAndRequires.process(Node, Node)
func (p *ProcessClosureProvidesAndRequires) Process(externs *rhino.Node, root *rhino.Node) {
	p.RewriteProvidesAndRequires(externs, root)
}

// Java: ProcessClosureProvidesAndRequires.getExportedVariableNames()
func (p *ProcessClosureProvidesAndRequires) GetExportedVariableNames() *jutil.LinkedHashSet[string] {
	return p.exportedVariables
}

// CollectProvidedNames collects all `goog.provide`s in the given namespace
// and warns on invalid code
//
// Java: ProcessClosureProvidesAndRequires.collectProvidedNames(Node, Node)
func (p *ProcessClosureProvidesAndRequires) CollectProvidedNames(externs *rhino.Node, root *rhino.Node) *jutil.LinkedHashMap[string, *ProcessClosureProvidesAndRequires_ProvidedName] {
	if p.providedNames.IsEmpty() {
		// goog is special-cased because it is provided in Closure's base library.
		p.putProvidedName(
			processClosureProvidesAndRequires_GOOG, newProcessClosureProvidesAndRequires_ProvidedNameBuilder(p).SetNamespace(processClosureProvidesAndRequires_GOOG).SetNode(nil).SetChunk(nil).Build())
		NodeTraversal_TraverseRoots(p.compiler, newProcessClosureProvidesAndRequires_CollectDefinitions(p), externs, root)
	}
	return p.providedNames
}

// RewriteProvidesAndRequires rewrites all provides and requires in the
// given namespace.
//
// Call this instead of collectProvidedNames(Node, Node) if you want
// rewriting.
//
// Java: ProcessClosureProvidesAndRequires.rewriteProvidesAndRequires(Node, Node)
func (p *ProcessClosureProvidesAndRequires) RewriteProvidesAndRequires(externs *rhino.Node, root *rhino.Node) {
	jutil.CheckStateMsg(!p.hasRewritingOccurred, "Cannot call rewriteProvidesAndRequires twice per instance")
	p.hasRewritingOccurred = true

	p.CollectProvidedNames(externs, root)

	for pn := range p.providedNames.Values() {
		pn.Replace(p.preserveGoogProvidesAndRequires, p.providedNames)
	}

	for closureRequire := range p.requiresToBeRemoved.All() {
		p.compiler.ReportChangeToEnclosingScope(closureRequire)
		closureRequire.Detach()
	}
	for forwardDeclare := range p.forwardDeclaresToRemove.All() {
		NodeUtil_DeleteNode(forwardDeclare, p.compiler)
	}
}

// ProcessClosureProvidesAndRequires_CollectDefinitions is the
// CollectDefinitions callback.
type ProcessClosureProvidesAndRequires_CollectDefinitions struct {
	outer *ProcessClosureProvidesAndRequires
}

// Java: ProcessClosureProvidesAndRequires.CollectDefinitions()
func newProcessClosureProvidesAndRequires_CollectDefinitions(outer *ProcessClosureProvidesAndRequires) *ProcessClosureProvidesAndRequires_CollectDefinitions {
	return &ProcessClosureProvidesAndRequires_CollectDefinitions{outer: outer}
}

// Java: ProcessClosureProvidesAndRequires.CollectDefinitions.shouldTraverse(NodeTraversal, Node, Node)
func (c *ProcessClosureProvidesAndRequires_CollectDefinitions) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	// Don't recurse into modules, which cannot have goog.provides.  We do need to handle legacy
	// goog.modules, but we do that quickly here rather than descending all the way into them.
	// We completely ignore ES modules and CommonJS modules.
	if (n.IsModuleBody() && n.GetParent().GetBooleanProp(rhino.Node_GOOG_MODULE)) ||
		NodeUtil_IsBundledGoogModuleScopeRoot(n) {
		googModuleCall := n.GetFirstChild()
		closureNamespace := googModuleCall.GetFirstChild().GetSecondChild().GetString()
		maybeLegacyNamespaceCall := googModuleCall.GetNext()
		if maybeLegacyNamespaceCall != nil &&
			NodeUtil_IsGoogModuleDeclareLegacyNamespaceCall(maybeLegacyNamespaceCall) {
			c.outer.ProcessLegacyModuleCall(closureNamespace, googModuleCall, t.GetChunk())
		}
		return false
	}
	return !n.IsModuleBody()
}

// Java: ProcessClosureProvidesAndRequires.CollectDefinitions.visit(NodeTraversal, Node, Node)
func (c *ProcessClosureProvidesAndRequires_CollectDefinitions) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_CALL:
		left := n.GetFirstChild()
		if left.IsGetProp() {
			name := left.GetFirstChild()
			if name.MatchesName(processClosureProvidesAndRequires_GOOG) {
				c.outer.VisitGoogMethodCall(t, parent, n, left.GetString())
			}
		}
	case rhino.Token_ASSIGN, rhino.Token_NAME:
		// If this is an assignment to a provided name, remove the provided object.
		c.outer.HandleCandidateProvideDefinition(t, n, parent)
	case rhino.Token_EXPR_RESULT:
		c.outer.HandleStubDefinition(t, n)
	default:
	}
}

// Java: ProcessClosureProvidesAndRequires.visitGoogMethodCall(NodeTraversal, Node, Node, String)
func (p *ProcessClosureProvidesAndRequires) VisitGoogMethodCall(t *NodeTraversal, parent *rhino.Node, n *rhino.Node, methodName string) {
	// For the sake of simplicity, we report code changes
	// when we see a provides/requires, and don't worry about
	// reporting the change when we actually do the replacement.
	switch methodName {
	case "exportSymbol":
		// Note: exportSymbol is allowed in local scope
		arg := n.GetSecondChild()
		if arg.IsStringLit() {
			argString := arg.GetString()
			dot := jutil.StringIndexOfChar(argString, '.')
			if dot == -1 {
				p.exportedVariables.Add(argString)
			} else {
				p.exportedVariables.Add(jutil.StringSubstring(argString, 0, dot))
			}
		}
	case "require", "requireType":
		if p.IsValidPrimitiveCall(t, n) {
			p.ProcessRequireCall(n, parent)
		}
	case "provide":
		if p.IsValidPrimitiveCall(t, n) {
			p.ProcessProvideCall(t, n, parent)
		}
	case "forwardDeclare":
		if p.IsValidPrimitiveCall(t, n) {
			p.ProcessForwardDeclare(n, parent)
		}
	default:
	}
}

// Java: ProcessClosureProvidesAndRequires.isValidPrimitiveCall(NodeTraversal, Node)
func (p *ProcessClosureProvidesAndRequires) IsValidPrimitiveCall(t *NodeTraversal, n *rhino.Node) bool {
	// Ignore invalid primitives if we didn't strip module sugar.
	if p.compiler.GetOptions().ShouldPreserveGoogModule() {
		return true
	}

	return t.InGlobalHoistScope() && n.GetParent().IsExprResult()
}

// ProcessRequireCall handles a goog.require or goog.requireType call.
//
// Java: ProcessClosureProvidesAndRequires.processRequireCall(Node, Node)
func (p *ProcessClosureProvidesAndRequires) ProcessRequireCall(call *rhino.Node, parent *rhino.Node) {
	if !p.VerifyOnlyArgumentIsString(call) {
		return
	}

	if !p.preserveGoogProvidesAndRequires {
		p.requiresToBeRemoved.Add(parent)
	}
}

// ProcessLegacyModuleCall handles a goog.module that is a legacy namespace.
//
// Java: ProcessClosureProvidesAndRequires.processLegacyModuleCall(String, Node, JSChunk)
func (p *ProcessClosureProvidesAndRequires) ProcessLegacyModuleCall(namespace string, googModuleCall *rhino.Node, chunk *JSChunk) {
	p.RegisterAnyProvidedPrefixes(namespace, googModuleCall, chunk)
	p.putProvidedName(
		namespace,
		newProcessClosureProvidesAndRequires_ProvidedNameBuilder(p).
			SetNamespace(namespace).
			SetNode(googModuleCall).
			SetChunk(chunk).
			SetExplicit(true).
			SetFromLegacyModule(true).
			Build())
}

// ProcessProvideCall handles a goog.provide call.
//
// Java: ProcessClosureProvidesAndRequires.processProvideCall(NodeTraversal, Node, Node)
func (p *ProcessClosureProvidesAndRequires) ProcessProvideCall(t *NodeTraversal, call *rhino.Node, parent *rhino.Node) {
	jutil.CheckState(call.IsCall())
	if !p.VerifyOnlyArgumentIsString(call) {
		return
	}
	left := call.GetFirstChild()
	arg := left.GetNext()
	ns := arg.GetString()

	info := NodeUtil_GetBestJSDocInfo(call)
	isImplicitlyInitialized := info != nil && info.IsProvideAlreadyProvided()

	if p.providedNames.ContainsKey(ns) {
		previouslyProvided := p.providedNames.Get(ns)
		if !previouslyProvided.IsExplicitlyProvided() {
			previouslyProvided.AddProvide(parent, t.GetChunk(), true /* explicit */, p.chunkGraph)
		}
	} else {
		p.RegisterAnyProvidedPrefixes(ns, parent, t.GetChunk())
		p.putProvidedName(
			ns,
			newProcessClosureProvidesAndRequires_ProvidedNameBuilder(p).
				SetNamespace(ns).
				SetNode(parent).
				SetChunk(t.GetChunk()).
				SetExplicit(true).
				SetHasImplicitInitialization(isImplicitlyInitialized).
				Build())
	}
}

// HandleStubDefinition handles a stub definition for a goog.provided name
// (e.g. a @typedef or a definition from externs)
//
// exprResult is an EXPR_RESULT node.
//
// Java: ProcessClosureProvidesAndRequires.handleStubDefinition(NodeTraversal, Node)
func (p *ProcessClosureProvidesAndRequires) HandleStubDefinition(t *NodeTraversal, exprResult *rhino.Node) {
	if !t.InGlobalHoistScope() {
		return
	}
	isExternStub := exprResult.IsFromExterns()
	isTypedefStub := p.IsTypedefStubDeclaration(exprResult)
	// recognize @typedefs only if using this pass for typechecking via
	// collectProvidedNames. We don't want rewriting to depend on @typedef annotations.
	isValidTypedefStubDefinition := isTypedefStub && !p.hasRewritingOccurred

	if isValidTypedefStubDefinition || isExternStub {
		if pn := p.getProvidedNameByNode(exprResult.GetFirstChild()); pn != nil {
			pn.AddDefinition(exprResult, t.GetChunk(), p.chunkGraph)
		}
	}
	if isTypedefStub {
		p.CheckNestedTypedefProvide(exprResult)
	}
}

// CheckNestedTypedefProvide checks that code doesn't use goog.provides in a
// way that will cause broken output code.
//
// exprResult is an EXPR_RESULT node with an @typedef type
//
// Java: ProcessClosureProvidesAndRequires.checkNestedTypedefProvide(Node)
func (p *ProcessClosureProvidesAndRequires) CheckNestedTypedefProvide(exprResult *rhino.Node) {
	// forbid this pattern:
	//   goog.provide('my.parent');
	//   goog.provide('my.parent.ChildTypedef');
	//   my.parent = [...]; // some initialization, doesn't matter what
	//   /** @typedef {...} */
	//   my.parent.ChildType;
	// at one point, this pattern was supported. now it would produce code that crashes at
	// runtime because the compiler would initializer `my.parent.ChildType = {}` before
	// `my.parent = [...]`, so report an error.

	child := exprResult.GetFirstChild()
	if !child.IsGetProp() {
		// @typedefs on simple names are okay.
		return
	}
	if _, ok := p.providedLeafs[child.GetString()]; !ok {
		return
	}
	name, ok := child.GetQualifiedNameString()
	if !ok || !(jutil.StringIndexOf(name, ".") >= 0) {
		return
	}
	if !p.providedNames.ContainsKey(name) {
		// non-provided names don't matter.
		return
	}
	parentName := jutil.StringSubstring(name, 0, jutil.StringLastIndexOf(name, "."))
	parent := p.providedNames.Get(parentName)
	parentDefinition := parent.GetCandidateDefinition()
	if parentDefinition == nil ||
		!(parentDefinition.GetStaticSourceFile() == exprResult.GetStaticSourceFile()) ||
		p.IsTypedefStubDeclaration(parentDefinition) {
		return
	}
	p.compiler.Report(JSError_MakeWithNode(exprResult, ProcessClosureProvidesAndRequires_TYPEDEF_CHILD_OF_PROVIDE, name, parentName))
}

// Java: ProcessClosureProvidesAndRequires.isTypedefStubDeclaration(Node)
func (p *ProcessClosureProvidesAndRequires) IsTypedefStubDeclaration(statement *rhino.Node) bool {
	if !statement.IsExprResult() {
		return false
	}
	info := NodeUtil_GetBestJSDocInfo(statement)
	return info != nil && info.HasTypedefType()
}

// HandleCandidateProvideDefinition handles a candidate definition for a
// goog.provided name.
//
// Java: ProcessClosureProvidesAndRequires.handleCandidateProvideDefinition(NodeTraversal, Node, Node)
func (p *ProcessClosureProvidesAndRequires) HandleCandidateProvideDefinition(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if !t.InGlobalHoistScope() {
		return
	}
	var pn *ProcessClosureProvidesAndRequires_ProvidedName
	switch n.GetParent().GetToken() {
	case rhino.Token_LET, rhino.Token_CONST:
		if t.InGlobalScope() {
			pn = p.providedNames.Get(n.GetString())
		}
	case rhino.Token_VAR:
		pn = p.providedNames.Get(n.GetString())
	case rhino.Token_EXPR_RESULT:
		if n.IsAssign() {
			pn = p.getProvidedNameByNode(n.GetFirstChild())
		}
	case rhino.Token_CLASS, rhino.Token_FUNCTION:
		// Class and function provides are forbidden; see ProcessClosurePrimitives's
		// CLASS_NAMESPACE_ERROR and FUNCTION_NAMESPACE_ERROR.
	default:
	}

	if pn != nil {
		pn.AddDefinition(parent, t.GetChunk(), p.chunkGraph)
	}
}

// ProcessForwardDeclare marks a goog.forwardDeclare call for removal.
//
// Java: ProcessClosureProvidesAndRequires.processForwardDeclare(Node, Node)
func (p *ProcessClosureProvidesAndRequires) ProcessForwardDeclare(n *rhino.Node, parent *rhino.Node) {
	convention := p.compiler.GetCodingConvention()

	typeDeclarations := convention.IdentifyTypeDeclarationCall(n)

	if !p.preserveGoogProvidesAndRequires &&
		typeDeclarations != nil &&
		len(typeDeclarations) == 1 {
		// Forward declaration was recorded and we can remove the call.
		toRemove := parent.GetParent()
		if parent.IsExprResult() {
			toRemove = parent
		}
		p.forwardDeclaresToRemove.Add(toRemove)
	}
}

// VerifyOnlyArgumentIsString verifies that a method call has exactly one
// argument, and that it's a string literal.
//
// Returns whether the argument checked out okay.
//
// Java: ProcessClosureProvidesAndRequires.verifyOnlyArgumentIsString(Node)
func (p *ProcessClosureProvidesAndRequires) VerifyOnlyArgumentIsString(call *rhino.Node) bool {
	arg := call.GetSecondChild()
	return arg != nil && arg.IsStringLit() && arg.GetNext() == nil
}

// RegisterAnyProvidedPrefixes registers ProvidedNames for prefix namespaces
// if they haven't already been defined. The prefix namespaces must be
// registered in order from shortest to longest.
//
// ns is the namespace whose prefixes may need to be provided; node is the
// EXPR of the provide call; chunk is the current chunk.
//
// Java: ProcessClosureProvidesAndRequires.registerAnyProvidedPrefixes(String, Node, JSChunk)
func (p *ProcessClosureProvidesAndRequires) RegisterAnyProvidedPrefixes(ns string, node *rhino.Node, chunk *JSChunk) {
	pos := jutil.StringIndexOfChar(ns, '.')
	for pos != -1 {
		prefixNs := jutil.StringSubstring(ns, 0, pos)
		pos = jutil.StringIndexOfCharFrom(ns, '.', pos+1)
		if p.providedNames.ContainsKey(prefixNs) {
			p.providedNames.Get(prefixNs).AddProvide(node, chunk, false /* explicit */, p.chunkGraph)
		} else {
			p.putProvidedName(
				prefixNs,
				newProcessClosureProvidesAndRequires_ProvidedNameBuilder(p).
					SetNamespace(prefixNs).
					SetNode(node).
					SetChunk(chunk).
					SetExplicit(false).
					Build())
		}
	}
}

// ProcessClosureProvidesAndRequires_ProvidedNameBuilder is
// ProvidedNameBuilder.
type ProcessClosureProvidesAndRequires_ProvidedNameBuilder struct {
	outer *ProcessClosureProvidesAndRequires

	namespace                 string
	node                      *rhino.Node
	chunk                     *JSChunk
	explicit                  bool
	fromLegacyModule          bool
	hasImplicitInitialization bool
}

// Java: ProcessClosureProvidesAndRequires.ProvidedNameBuilder()
func newProcessClosureProvidesAndRequires_ProvidedNameBuilder(outer *ProcessClosureProvidesAndRequires) *ProcessClosureProvidesAndRequires_ProvidedNameBuilder {
	return &ProcessClosureProvidesAndRequires_ProvidedNameBuilder{outer: outer}
}

// Java: ProcessClosureProvidesAndRequires.ProvidedNameBuilder.setNamespace(String)
func (b *ProcessClosureProvidesAndRequires_ProvidedNameBuilder) SetNamespace(namespace string) *ProcessClosureProvidesAndRequires_ProvidedNameBuilder {
	b.namespace = namespace
	return b
}

// SetNode: node can be null (for GOOG or an implicit name), an EXPR_RESULT
// for a goog.provide, or an EXPR_RESULT or name declaration for a
// previously provided name.
//
// Java: ProcessClosureProvidesAndRequires.ProvidedNameBuilder.setNode(Node)
func (b *ProcessClosureProvidesAndRequires_ProvidedNameBuilder) SetNode(node *rhino.Node) *ProcessClosureProvidesAndRequires_ProvidedNameBuilder {
	b.node = node
	return b
}

// Java: ProcessClosureProvidesAndRequires.ProvidedNameBuilder.setChunk(JSChunk)
func (b *ProcessClosureProvidesAndRequires_ProvidedNameBuilder) SetChunk(chunk *JSChunk) *ProcessClosureProvidesAndRequires_ProvidedNameBuilder {
	b.chunk = chunk
	return b
}

// SetExplicit: explicit is whether this came from an actual
// goog.provide('a.b.c'); call
//
// Java: ProcessClosureProvidesAndRequires.ProvidedNameBuilder.setExplicit(boolean)
func (b *ProcessClosureProvidesAndRequires_ProvidedNameBuilder) SetExplicit(explicit bool) *ProcessClosureProvidesAndRequires_ProvidedNameBuilder {
	b.explicit = explicit
	return b
}

// SetHasImplicitInitialization: alreadyInitialized is whether this came
// from an actual goog.provide('a.b.c'); call
//
// Java: ProcessClosureProvidesAndRequires.ProvidedNameBuilder.setHasImplicitInitialization(boolean)
func (b *ProcessClosureProvidesAndRequires_ProvidedNameBuilder) SetHasImplicitInitialization(alreadyInitialized bool) *ProcessClosureProvidesAndRequires_ProvidedNameBuilder {
	b.hasImplicitInitialization = alreadyInitialized
	return b
}

// SetFromLegacyModule is whether this comes from a legacy goog.module
//
// Java: ProcessClosureProvidesAndRequires.ProvidedNameBuilder.setFromLegacyModule(boolean)
func (b *ProcessClosureProvidesAndRequires_ProvidedNameBuilder) SetFromLegacyModule(fromLegacyModule bool) *ProcessClosureProvidesAndRequires_ProvidedNameBuilder {
	b.fromLegacyModule = fromLegacyModule
	return b
}

// Java: ProcessClosureProvidesAndRequires.ProvidedNameBuilder.build()
func (b *ProcessClosureProvidesAndRequires_ProvidedNameBuilder) Build() *ProcessClosureProvidesAndRequires_ProvidedName {
	return newProcessClosureProvidesAndRequires_ProvidedName(b, b.outer.chunkGraph, b.outer.compiler, b.outer.astFactory)
}

// -------------------------------------------------------------------------

// ProcessClosureProvidesAndRequires_ProvidedName stores information about a
// Closure namespace created by a goog.provide
//
// There are three ways that we find these namespaces in the AST:
//
//   - An explicit goog.provide. `goog.provide('a.b.c');` creates a
//     ProvidedName for 'a.b.c'.
//   - An implicit parent namespace. `goog.provide('a.b.c');` creates a
//     ProvidedName for 'a.b'.
//   - A provide definition processed by an earlier run. `a.b.c = {};` when
//     annotated IS_NAMESPACE
type ProcessClosureProvidesAndRequires_ProvidedName struct {
	// The Closure namespace this name represents, e.g. `a.b` for
	// `goog.provide('a.b');`
	namespace string

	// The first node in the AST that creates this ProvidedName.
	// This is always a goog.provide('a.b'), null (for implicit namespaces and
	// 'goog'), or an assignment or declaration for a 'previously provided'
	// name or parent namespace of such. This should only be used for source
	// info and a place to hang namespace definitions.
	firstNode *rhino.Node
	// The chunk where this namespace was first goog.provided, if chunks
	// exist.
	firstChunk *JSChunk

	// When set the namespace will not have an initialization added, even if
	// the namespace initialization is not otherwise visibile to this pass.
	// This allows for trivial aliasing of a tree of namespaces.
	hasImplicitInitialization bool

	// The node where the call was explicitly goog.provided. Null if the
	// namespace is implicit. If this is previously provided, this will
	// instead be the expression or declaration marked as IS_NAMESPACE.
	explicitNode *rhino.Node

	// The candidate definition for this namespace. For example, given
	//      goog.provide('a.b');
	//      /** @constructor * /
	//      a.b = function() {};
	// the 'candidate definition' of 'a.b' is the GETPROP 'a.b' from the
	// constructor declaration.
	candidateDefinition *rhino.Node

	// The minimum chunk where the provide namespace definition must appear.
	// If child namespaces of this provide appear in multiple chunks, this
	// chunk must be earlier than all child namespace's chunks.
	minimumChunk *JSChunk

	// The replacement declaration. Null until replace() has been called.
	replacementNode *rhino.Node

	// Whether this comes from a goog.module with declareLegacyNamespace.
	fromLegacyModule bool

	compiler   AbstractCompiler
	astFactory *AstFactory
}

// Java: ProcessClosureProvidesAndRequires.ProvidedName(ProvidedNameBuilder, JSChunkGraph, AbstractCompiler, AstFactory)
func newProcessClosureProvidesAndRequires_ProvidedName(
	builder *ProcessClosureProvidesAndRequires_ProvidedNameBuilder,
	chunkGraph *JSChunkGraph,
	compiler AbstractCompiler,
	astFactory *AstFactory) *ProcessClosureProvidesAndRequires_ProvidedName {
	pn := &ProcessClosureProvidesAndRequires_ProvidedName{}
	pn.explicitNode = nil
	pn.candidateDefinition = nil
	pn.minimumChunk = nil
	pn.replacementNode = nil
	pn.compiler = compiler
	pn.astFactory = astFactory
	node := builder.node
	jutil.CheckArgumentMsg1(
		node == nil ||
			NodeUtil_IsExprCall(node) ||
			(!builder.explicit &&
				(NodeUtil_IsExprAssign(node) ||
					NodeUtil_IsNameDeclaration(node) ||
					(node.IsExprResult() && node.GetFirstChild().IsQualifiedName()))),
		"%s",
		node)
	pn.namespace = builder.namespace
	pn.firstNode = builder.node
	pn.firstChunk = builder.chunk
	pn.fromLegacyModule = builder.fromLegacyModule
	pn.hasImplicitInitialization = builder.hasImplicitInitialization

	pn.AddProvide(node, builder.chunk, builder.explicit, chunkGraph)
	return pn
}

// AddProvide adds an implicit or explicit provide.
//
// Every provided name can have multiple implicit provides but a maximum of
// one explicit provide.
//
// node is the EXPR_RESULT representing this provide or possible a VAR for a
// previously provided name. null if implicit.
//
// Java: ProcessClosureProvidesAndRequires.ProvidedName.addProvide(Node, JSChunk, boolean, JSChunkGraph)
func (pn *ProcessClosureProvidesAndRequires_ProvidedName) AddProvide(node *rhino.Node, chunk *JSChunk, explicit bool, chunkGraph *JSChunkGraph) {
	if explicit {
		// goog.provide('name.space');
		jutil.CheckState(pn.explicitNode == nil)
		jutil.CheckArgumentMsg1(node.IsExprResult(), "%s", node)
		pn.explicitNode = node
	}
	pn.UpdateMinimumChunk(chunk, chunkGraph)
}

// IsExplicitlyProvided is whether there existed a `goog.provide('a.b');`
// for this name 'a.b'
//
// Java: ProcessClosureProvidesAndRequires.ProvidedName.isExplicitlyProvided()
func (pn *ProcessClosureProvidesAndRequires_ProvidedName) IsExplicitlyProvided() bool {
	return pn.explicitNode != nil
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
	return pn.candidateDefinition != nil
}

// GetFirstProvideCall returns the `goog.provide` or legacy namespace
// `goog.module` call that created this name, if any, or otherwise the first
// 'previous provide' assignment that created this name.
//
// Java: ProcessClosureProvidesAndRequires.ProvidedName.getFirstProvideCall()
func (pn *ProcessClosureProvidesAndRequires_ProvidedName) GetFirstProvideCall() *rhino.Node {
	return pn.firstNode
}

// GetCandidateDefinition returns the definition of this provided namespace
// in the input code, if any, or null.
//
// For example, this returns `a.b = class {};` given 'a.b' in
//
//	goog.provide('a.b');
//	a.b = class {};
//
// Note: this method will only return candidate definitions that count
// towards provide rewriting. If a name is defined, then provided, the
// candidate definition will not be the early definition. This doesn't
// completely mimic uncompiled behavior, but supports some legacy code.
// Externs definitions never count.
//
// Java: ProcessClosureProvidesAndRequires.ProvidedName.getCandidateDefinition()
func (pn *ProcessClosureProvidesAndRequires_ProvidedName) GetCandidateDefinition() *rhino.Node {
	return pn.candidateDefinition
}

// GetNamespace returns the Closure namespace of this provide, e.g. "a.b"
// for `goog.provide('a.b');`
//
// Java: ProcessClosureProvidesAndRequires.ProvidedName.getNamespace()
func (pn *ProcessClosureProvidesAndRequires_ProvidedName) GetNamespace() string {
	return pn.namespace
}

// AddDefinition records function declaration, variable declarations, and
// assignments that refer to this provided namespace.
//
// This pass gives preference to declarations. If no declaration exists,
// records a reference to an assignment so it can be repurposed later into a
// declaration.
//
// Java: ProcessClosureProvidesAndRequires.ProvidedName.addDefinition(Node, JSChunk, JSChunkGraph)
func (pn *ProcessClosureProvidesAndRequires_ProvidedName) AddDefinition(node *rhino.Node, chunk *JSChunk, chunkGraph *JSChunkGraph) {
	jutil.CheckArgument(
		node.IsExprResult() || // assign
			node.IsFunction() ||
			NodeUtil_IsNameDeclaration(node))
	jutil.CheckArgument(pn.explicitNode != node)
	if (pn.candidateDefinition == nil) || !node.IsExprResult() {
		pn.candidateDefinition = node
		pn.UpdateMinimumChunk(chunk, chunkGraph)
	}
}

// Java: ProcessClosureProvidesAndRequires.ProvidedName.updateMinimumChunk(JSChunk, JSChunkGraph)
func (pn *ProcessClosureProvidesAndRequires_ProvidedName) UpdateMinimumChunk(newChunk *JSChunk, chunkGraph *JSChunkGraph) {
	if pn.minimumChunk == nil {
		pn.minimumChunk = newChunk
	} else if chunkGraph.GetChunkCount() > 1 {
		pn.minimumChunk = chunkGraph.GetDeepestCommonDependencyInclusive(pn.minimumChunk, newChunk)
	} else {
		// If there is no chunk graph, then there must be exactly one chunk in the program.
		jutil.CheckStateMsg(newChunk == pn.minimumChunk, "Missing chunk graph")
	}
}

// Replace replaces the provide statement.
//
// If we're providing a name with no definition, then create one. If we're
// providing a name with a duplicate definition, then make sure that
// definition becomes a declaration.
//
// Java: ProcessClosureProvidesAndRequires.ProvidedName.replace(boolean, Map<String, ProvidedName>)
func (pn *ProcessClosureProvidesAndRequires_ProvidedName) Replace(
	preserveGoogProvidesAndRequires bool, providedNames *jutil.LinkedHashMap[string, *ProcessClosureProvidesAndRequires_ProvidedName]) {
	jutil.CheckStateMsg1(
		!pn.IsFromLegacyModule(),
		"Cannot rewrite provides without having rewritten goog.modules, found %s",
		pn.firstNode)
	if pn.firstNode == nil {
		// Don't touch the base case ('goog').
		pn.replacementNode = pn.candidateDefinition
		return
	}

	// Handle the case where there is a duplicate definition for an explicitly
	// provided symbol.
	if pn.HasCandidateDefinition() && pn.explicitNode != nil {
		// Does this need a VAR keyword?
		pn.replacementNode = pn.candidateDefinition
		if pn.candidateDefinition.IsExprResult() {
			exprNode := pn.candidateDefinition.GetOnlyChild()
			if exprNode.IsAssign() {
				nameNode := exprNode.GetFirstChild()
				if nameNode.IsName() {
					// In the case of a simple name, `name = value;`, we need to ensure the name is
					// actually declared with `var`.
					pn.ConvertProvideAssignmentToVarDeclaration(exprNode, nameNode)
				} else {
					// `some.provided.namespace = value;`
					// We don't need to change the definition, but mark it as 'IS_NAMESPACE' so that
					// future passes know this was originally provided.
					pn.candidateDefinition.PutBooleanProp(rhino.Node_IS_NAMESPACE, true)
				}
			}
		}
	} else {
		// Handle the case where there's not an existing definition.
		if !pn.hasImplicitInitialization {
			pn.CreateNamespaceInitialization(
				pn.CreateDeclarationNode(pn.astFactory.CreateObjectLit()), providedNames)
		}
	}

	// Remove the `goog.provide('a.b.c');` call.
	if pn.explicitNode != nil {
		if preserveGoogProvidesAndRequires {
			return
		}
		pn.compiler.ReportChangeToEnclosingScope(pn.explicitNode)
		pn.explicitNode.Detach()
	}
}

// Java: ProcessClosureProvidesAndRequires.ProvidedName.convertProvideAssignmentToVarDeclaration(Node, Node)
func (pn *ProcessClosureProvidesAndRequires_ProvidedName) ConvertProvideAssignmentToVarDeclaration(assignNode *rhino.Node, nameNode *rhino.Node) {
	// Convert `providedName = value;` into `var providedName = value;`.
	jutil.CheckArgumentMsg1(assignNode.IsAssign(), "%s", assignNode)
	jutil.CheckArgumentMsg1(nameNode.IsName(), "%s", nameNode)
	valueNode := nameNode.GetNext()
	nameNode.Detach()
	valueNode.Detach()

	varNode := rhino.IR_VarWithValue(nameNode, valueNode).Srcref(pn.candidateDefinition)
	varNode.SetJSDocInfo(assignNode.GetJSDocInfo())
	varNode.PutBooleanProp(rhino.Node_IS_NAMESPACE, true)

	pn.candidateDefinition.ReplaceWith(varNode)
	pn.replacementNode = varNode
	pn.compiler.ReportChangeToEnclosingScope(varNode)
}

// CreateNamespaceInitialization adds an assignment or declaration to this
// namespace to the AST, using the provided value
//
// Java: ProcessClosureProvidesAndRequires.ProvidedName.createNamespaceInitialization(Node, Map<String, ProvidedName>)
func (pn *ProcessClosureProvidesAndRequires_ProvidedName) CreateNamespaceInitialization(
	replacement *rhino.Node, providedNames *jutil.LinkedHashMap[string, *ProcessClosureProvidesAndRequires_ProvidedName]) {
	pn.replacementNode = replacement
	if pn.firstChunk == pn.minimumChunk {
		pn.replacementNode.InsertBefore(pn.firstNode)
	} else {
		// In this case, the name was implicitly provided by two independent
		// modules. We need to move this code up to a common module.
		indexOfDot := jutil.StringLastIndexOf(pn.namespace, ".")
		if indexOfDot == -1 {
			// Any old place is fine.
			pn.compiler.GetNodeForCodeInsertion(pn.minimumChunk).AddChildToBack(pn.replacementNode)
		} else {
			// Add it after the parent namespace.
			parentName := providedNames.Get(jutil.StringSubstring(pn.namespace, 0, indexOfDot))
			jutil.CheckNotNull(parentName)
			jutil.CheckNotNull(parentName.replacementNode)
			pn.replacementNode.InsertAfter(parentName.replacementNode)
		}
	}
	pn.compiler.ReportChangeToEnclosingScope(pn.replacementNode)
}

// CreateDeclarationNode creates the declaration node for this name, without
// inserting it into the AST.
//
// value is the object literal namespace, possibly in a CAST
//
// Java: ProcessClosureProvidesAndRequires.ProvidedName.createDeclarationNode(Node)
func (pn *ProcessClosureProvidesAndRequires_ProvidedName) CreateDeclarationNode(value *rhino.Node) *rhino.Node {
	jutil.CheckArgumentMsg1(value.IsObjectLit() || value.IsCast(), "%s", value)
	if jutil.StringIndexOfChar(pn.namespace, '.') == -1 {
		return pn.MakeVarDeclNode(value)
	} else {
		return pn.MakeAssignmentExprNode(value)
	}
}

// MakeVarDeclNode creates a simple namespace variable declaration (e.g.
// var foo = {};).
//
// Java: ProcessClosureProvidesAndRequires.ProvidedName.makeVarDeclNode(Node)
func (pn *ProcessClosureProvidesAndRequires_ProvidedName) MakeVarDeclNode(value *rhino.Node) *rhino.Node {
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

	jutil.CheckState(ProcessClosureProvidesAndRequires_IsNamespacePlaceholder(decl))
	pn.SetSourceInfo(decl)
	return decl
}

// MakeAssignmentExprNode creates a dotted namespace assignment expression
// (e.g. foo.bar = {};).
//
// Java: ProcessClosureProvidesAndRequires.ProvidedName.makeAssignmentExprNode(Node)
func (pn *ProcessClosureProvidesAndRequires_ProvidedName) MakeAssignmentExprNode(value *rhino.Node) *rhino.Node {
	// Note: as of May 2021, using the unknown type vs. the actual inferred type both produced the
	// same optimized JS after type-based optimizations. So the lack of type info is intentional.
	lhs := pn.astFactory.CreateQNameWithUnknownType(pn.namespace).SrcrefTree(pn.firstNode)
	decl := rhino.IR_ExprResult(pn.astFactory.CreateAssign(lhs, value))
	decl.PutBooleanProp(rhino.Node_IS_NAMESPACE, true)
	if !pn.HasCandidateDefinition() {
		decl.GetFirstChild().SetJSDocInfo(NodeUtil_CreateConstantJsDoc())
	}
	jutil.CheckState(ProcessClosureProvidesAndRequires_IsNamespacePlaceholder(decl))
	pn.SetSourceInfo(decl)
	// This function introduces artifical nodes and we don't need them for indexing.
	// Marking all but the last one as non-indexable. So if this function adds:
	// foo.bar.baz = {};
	// then we mark foo and bar as non-indexable.
	lhs.GetFirstChild().MakeNonIndexableRecursive()
	return decl
}

// SetSourceInfo copies source info to the new node.
//
// Java: ProcessClosureProvidesAndRequires.ProvidedName.setSourceInfo(Node)
func (pn *ProcessClosureProvidesAndRequires_ProvidedName) SetSourceInfo(newNode *rhino.Node) {
	sourceInfoNode := pn.firstNode

	provideStringNode := pn.GetProvideStringNode()
	if provideStringNode != nil {
		// Given namespace "foo.bar.baz" we create node for "baz" here and need to calculate
		// length and start of the last component which is "baz".
		firstCharIndex := jutil.StringLastIndexOf(pn.namespace, ".") + 1 // If no dots, then 0.

		sourceInfoNode = provideStringNode.CloneNode()
		sourceInfoNode.SetLinenoCharno(
			sourceInfoNode.GetLineno(),
			sourceInfoNode.GetCharno()+firstCharIndex+1) // +1 for quote
		sourceInfoNode.SetLength(jutil.StringLength(pn.namespace) - firstCharIndex)
	}

	newNode.SrcrefTree(sourceInfoNode)
}

// Java: ProcessClosureProvidesAndRequires.ProvidedName.getProvideStringNode()
func (pn *ProcessClosureProvidesAndRequires_ProvidedName) GetProvideStringNode() *rhino.Node {
	if pn.firstNode.HasChildren() && NodeUtil_IsExprCall(pn.firstNode) {
		return pn.firstNode.GetFirstChild().GetLastChild()
	}
	return nil
}

// Java: ProcessClosureProvidesAndRequires.ProvidedName.toString()
func (pn *ProcessClosureProvidesAndRequires_ProvidedName) String() string {
	explicitOrImplicit := "implicit"
	if pn.IsExplicitlyProvided() {
		explicitOrImplicit = "explicit"
	}
	return "ProvidedName: " + pn.namespace + ", " + explicitOrImplicit
}

// ProcessClosureProvidesAndRequires_IsNamespacePlaceholder returns whether
// the node initializes a goog.provide'd namespace (e.g. `a.b = {};`) with a
// simple namespace object literal (e.g. not `a.b = class {}`;)
//
// Java: ProcessClosureProvidesAndRequires.isNamespacePlaceholder(Node)
func ProcessClosureProvidesAndRequires_IsNamespacePlaceholder(n *rhino.Node) bool {
	if !n.GetBooleanProp(rhino.Node_IS_NAMESPACE) {
		return false
	}

	var value *rhino.Node = nil
	if n.IsExprResult() {
		assign := n.GetFirstChild()
		value = assign.GetLastChild()
	} else if n.IsVar() {
		name := n.GetFirstChild()
		value = name.GetFirstChild()
	}

	if value == nil {
		return false
	}
	if value.IsCast() {
		// There may be a cast to unknown type wrapped around the value.
		value = value.GetOnlyChild()
	}
	return value.IsObjectLit() && !value.HasChildren()
}
