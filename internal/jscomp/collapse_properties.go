// Port of com/google/javascript/jscomp/InlineAndCollapseProperties.java @ closure-compiler v20261021

/*
 * Copyright 2021 The Closure Compiler Authors.
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
	"cmp"
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// InlineAndCollapseProperties_CollapseProperties flattens global
// objects/namespaces by replacing each '.' with '$' in their names.
//
// This reduces the number of property lookups the browser has to do and
// allows the RenameVars pass to shorten namespaced names. For example,
// goog.events.handleEvent() -> goog$events$handleEvent() -> Za().
//
// If a global object's name is assigned to more than once, or if a property
// is added to the global object in a complex expression, then none of its
// properties will be collapsed (for safety/correctness).
//
// If, after a global object is declared, it is never referenced except when
// its properties are read or set, then the object will be removed after its
// properties have been collapsed.
//
// Uninitialized variable stubs are created at a global object's declaration
// site for any of its properties that are added late in a local scope.
//
// Static properties of constructors are always collapsed, unsafely! For
// other objects: if, after an object is declared, it is referenced directly
// in a way that might create an alias for it, then none of its properties
// will be collapsed. This behavior is a safeguard to prevent the values
// associated with the flattened names from getting out of sync with the
// object's actual property values. For example, in the following case, an
// alias a$b, if created, could easily keep the value 0 even after a.b
// became 5: `a = {b: 0}; c = a; c.b = 5;`.
//
// This pass may break code, but relies on AggressiveInlineAliases running
// before this pass to make some common patterns safer.
//
// This pass doesn't flatten property accesses of the form: a[b].
//
// For lots of examples, see the unit test.
type InlineAndCollapseProperties_CollapseProperties struct {
	outer *InlineAndCollapseProperties

	namespace *GlobalNamespace

	dynamicallyImportedModules map[string]bool
}

// Java: InlineAndCollapseProperties.CollapseProperties()
func NewInlineAndCollapseProperties_CollapseProperties(outer *InlineAndCollapseProperties) *InlineAndCollapseProperties_CollapseProperties {
	return &InlineAndCollapseProperties_CollapseProperties{
		outer: outer,
	}
}

func (c *InlineAndCollapseProperties_CollapseProperties) addDynamicallyImportedModule(modName string) {
	if c.dynamicallyImportedModules == nil {
		c.dynamicallyImportedModules = make(map[string]bool)
	}
	c.dynamicallyImportedModules[modName] = true
}

// Java: InlineAndCollapseProperties.CollapseProperties.process(Node, Node)
func (c *InlineAndCollapseProperties_CollapseProperties) Process(externs rhino.Node, root rhino.Node) {
	if c.outer.propertyCollapseLevel == CompilerOptions_PropertyCollapseLevel_MODULE_EXPORT ||
		c.outer.chunkOutputType == CompilerOptions_ChunkOutputType_ES_MODULES {
		NodeTraversal_Traverse(
			c.outer.compiler,
			root,
			newCollapseProperties_findDynamicallyImportedModules(c, c.outer.haveModulesBeenRewritten, c.outer.moduleResolutionMode))
	}

	c.namespace = jutil.CheckNotNullMsg(c.outer.namespace, "namespace was not initialized")
	c.namespace.EnsureGenerated()
	c.CheckNamespaces()
	for name := range c.namespace.NameForest() {
		c.FlattenReferencesToCollapsibleDescendantNames(name, name.GetBaseName())
		// We collapse property definitions after collapsing property references
		// because this step can alter the parse tree above property references,
		// invalidating the node ancestry stored with each reference.
		c.CollapseDeclarationOfNameAndDescendants(name, name.GetBaseName())
	}
}

// Java: InlineAndCollapseProperties.CollapseProperties.canCollapse(Name)
func (c *InlineAndCollapseProperties_CollapseProperties) CanCollapse(name *GlobalNamespace_Name) bool {
	inlinability := name.CanCollapseOrInline()
	if !inlinability.CanCollapse() {
		c.LogDecisionForNameWithInlinability(name, inlinability, "canCollapse() returns false")
		return false
	}

	if c.outer.propertyCollapseLevel == CompilerOptions_PropertyCollapseLevel_MODULE_EXPORT {
		if !name.IsModuleExport() {
			c.LogDecisionForNameWithInlinability(name, inlinability, "module export: canCollapse() returns false")
			return false
		} else if c.dynamicallyImportedModules[name.GetBaseName()] {
			c.LogDecisionForNameWithInlinability(
				name, inlinability, "dynamic module export: canCollapse() returns false")
			return false
		}
	}

	c.LogDecisionForNameWithInlinability(name, inlinability, "canCollapse() returns true")
	return true
}

// CanEliminate returns true if it is safe to completely remove the rvalue
// assigned to the name from the AST.
//
// Java: InlineAndCollapseProperties.CollapseProperties.canEliminate(Name, Node)
func (c *InlineAndCollapseProperties_CollapseProperties) CanEliminate(name *GlobalNamespace_Name, rvalue rhino.Node) bool {
	if !name.CanEliminate() {
		return false
	}

	if objLit, ok := ast.AsObjectLit(rvalue); ok {
		for member := range objLit.Members() {
			child := member.Node()
			switch child.GetToken() {
			case rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_COMPUTED_PROP, rhino.Token_OBJECT_SPREAD:
				// We cannot completely eliminate the object literal, because it contains
				// at least one feature that we cannot extract as a collapsed variable name.
				return false
			case rhino.Token_STRING_KEY, rhino.Token_MEMBER_FUNCTION_DEF:
			default:
				panic(jutil.NewIllegalStateExceptionMsg(
					"Unexpected child of OBJECTLIT: " + child.ToStringTree()))
			}
		}
	}
	if !name.HasProps() ||
		c.outer.propertyCollapseLevel != CompilerOptions_PropertyCollapseLevel_MODULE_EXPORT {
		return true
	}

	return false
}

// CheckNamespaces runs through all namespaces (prefixes of classes and
// enums), and checks if any of them have been used in an unsafe way.
//
// Java: InlineAndCollapseProperties.CollapseProperties.checkNamespaces()
func (c *InlineAndCollapseProperties_CollapseProperties) CheckNamespaces() {
	if len(c.dynamicallyImportedModules) > 0 {
		// When the output chunk type is ES_MODULES, properties of the module namespace
		// must not be collapsed as they are referenced off the namespace object. The namespace
		// objects escape via the dynamic import expression.
		dynNames := make([]*GlobalNamespace_Name, 0, len(c.dynamicallyImportedModules))
		for modName := range c.dynamicallyImportedModules { // order-independent: dynNames is immediately sorted by unique GlobalNamespace_Name.selfID below
			if name := c.namespace.LookupName(modName); name != nil {
				dynNames = append(dynNames, name)
			}
		}
		slices.SortFunc(dynNames, func(a, b *GlobalNamespace_Name) int {
			return cmp.Compare(a.selfID, b.selfID)
		})
		for _, name := range dynNames {
			// Mark the rewritten module namespace variable as escaped to prevent any property
			// collapsing (e.g. `/** @const */ var module$foo = {};`).
			c.LogDecisionForName(name, "escapes - dynamically imported module namespace")
			name.markEscaped()
			if !name.HasProps() {
				continue
			}
			for prop := range name.Props() {
				propDeclaration := prop.getDeclarationRef()
				if propDeclaration == nil || propDeclaration.GetNode().IsNil() {
					continue
				}
				// ES Module rewriting creates aliases on the module namespace object. These aliased
				// names also escape and their properties may not be collapsed.
				//
				// example:
				// class Foo$$module$foo {
				//   static bar() { return 'bar'; }
				// }
				// /** @const */ var module$foo = {};
				// /** @const */ module$foo.Foo = Foo$$module$foo;
				//
				// While module$foo.Foo cannot be collapsed because we marked the module namespace
				// as escaped, we also need to prevent any property collapsing on the
				// Foo$$module$foo
				// class itself.
				rValue := NodeUtil_GetRValueOfLValue(propDeclaration.GetNode())
				if rValue.IsName() {
					c.LogDecisionForName(
						name, "escapes - dynamically imported module namespace property alias")
					if target := c.namespace.LookupName(rValue.GetString()); target != nil {
						target.markEscaped()
					}
				}
			}
		}
	}

	for name := range c.namespace.AllNames() {
		if !name.IsNamespaceObjectLit() {
			continue
		}
		if name.GetAliasingGets() == 0 &&
			name.GetLocalSets()+name.GetGlobalSets() <= 1 &&
			name.GetDeleteProps() == 0 {
			continue
		}
		initialized := name.GetDeclaration() != nil
		for ref := range name.Refs() {
			if ref.IsDeleteProp() {
				if initialized {
					c.WarnAboutNamespaceRedefinition(name, ref)
				}
			} else if ref.IsSet() && ref != name.getDeclarationRef() {
				if initialized && !InlineAndCollapseProperties_IsSafeNamespaceReinit(ref) {
					c.WarnAboutNamespaceRedefinition(name, ref)
				}

				initialized = true
			} else if ref.IsAliasingGet() {
				c.WarnAboutNamespaceAliasing(name, ref)
				c.LogDecisionForName(name, "escapes")
				name.markEscaped()
				break
			}
		}
	}
}

// WarnAboutNamespaceAliasing reports a warning because a namespace was
// aliased.
//
// nameObj is a namespace that is being aliased; ref the reference that
// forced the alias.
//
// Java: InlineAndCollapseProperties.CollapseProperties.warnAboutNamespaceAliasing(Name, Ref)
func (c *InlineAndCollapseProperties_CollapseProperties) WarnAboutNamespaceAliasing(nameObj *GlobalNamespace_Name, ref *GlobalNamespace_Ref) {
	c.outer.compiler.Report(
		JSError_MakeWithNode(ref.GetNode(), InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING, nameObj.GetFullName()))
}

// WarnAboutNamespaceRedefinition reports a warning because a namespace was
// redefined.
//
// nameObj is a namespace that is being redefined; ref the reference that set
// the namespace.
//
// Java: InlineAndCollapseProperties.CollapseProperties.warnAboutNamespaceRedefinition(Name, Ref)
func (c *InlineAndCollapseProperties_CollapseProperties) WarnAboutNamespaceRedefinition(nameObj *GlobalNamespace_Name, ref *GlobalNamespace_Ref) {
	c.outer.compiler.Report(
		JSError_MakeWithNode(ref.GetNode(), InlineAndCollapseProperties_NAMESPACE_REDEFINED_WARNING, nameObj.GetFullName()))
}

// FlattenReferencesToCollapsibleDescendantNames flattens all references to
// collapsible properties of a global name except their initial definitions.
// Recurs on subnames.
//
// n is an object representing a global name; alias the flattened name for
// n.
//
// Java: InlineAndCollapseProperties.CollapseProperties.flattenReferencesToCollapsibleDescendantNames(Name, String, Set<Name>)
func (c *InlineAndCollapseProperties_CollapseProperties) FlattenReferencesToCollapsibleDescendantNames(
	n *GlobalNamespace_Name, alias string) {
	if !n.HasProps() {
		return
	}
	if n.IsCollapsingExplicitlyDenied() {
		c.LogDecisionForName(n, "@nocollapse: will not flatten descendant name references")
		return
	}
	if n.isEscaped() {
		c.LogDecisionForName(n, "escapes: will not flatten descendant name references")
		return
	}

	for p := range n.Props() {
		inlinability := p.CanCollapseOrInline()

		isAllowedToCollapse :=
			c.outer.propertyCollapseLevel != CompilerOptions_PropertyCollapseLevel_MODULE_EXPORT || p.IsModuleExport()
		willCollapse := isAllowedToCollapse && inlinability.CanCollapse()
		willStub := isAllowedToCollapse && !inlinability.CanCollapse() && !p.IsCollapsingExplicitlyDenied() && p.IsSimpleStubDeclaration()
		willRecurse := p.HasProps() && (!p.IsCollapsingExplicitlyDenied() && !p.isEscaped() || c.outer.decisionsLog.IsLogging())

		var propAlias string
		if willCollapse || willStub || (p.HasProps() && !p.IsCollapsingExplicitlyDenied() && !p.isEscaped()) {
			propAlias = c.AppendPropForAlias(alias, p.GetBaseName())
		}

		if isAllowedToCollapse {
			if willCollapse {
				c.LogDecisionForNameWithInlinability(p, inlinability, "will flatten references")
				c.flattenReferencesTo(p, propAlias, true)
			} else if p.IsCollapsingExplicitlyDenied() {
				c.LogDecisionForName(p, "@nocollapse: will not flatten references")
			} else if willStub {
				c.LogDecisionForName(p, "simple stub declaration: will flatten references")
				c.FlattenSimpleStubDeclaration(p, propAlias)
			} else {
				c.LogDecisionForNameWithInlinability(p, inlinability, "will not flatten references")
			}
		}

		if willRecurse {
			c.FlattenReferencesToCollapsibleDescendantNames(p, propAlias)
		}
	}
}

// Java: InlineAndCollapseProperties.CollapseProperties.logDecisionForName(Name, Inlinability, String)
func (c *InlineAndCollapseProperties_CollapseProperties) LogDecisionForNameWithInlinability(
	name *GlobalNamespace_Name, inlinability GlobalNamespace_Inlinability, message string) {
	if c.outer.decisionsLog.IsLogging() {
		c.LogDecisionForNameSupplier(name, func() string { return "inlinability " + inlinability.String() + ": " + message })
	}
}

// Java: InlineAndCollapseProperties.CollapseProperties.logDecisionForName(Name, String)
func (c *InlineAndCollapseProperties_CollapseProperties) LogDecisionForName(name *GlobalNamespace_Name, message string) {
	if c.outer.decisionsLog.IsLogging() {
		c.outer.decisionsLog.LogSupplier(func() string { return name.GetFullName() + ": " + message })
	}
}

// Java: InlineAndCollapseProperties.CollapseProperties.logDecisionForName(Name, Supplier<String>)
func (c *InlineAndCollapseProperties_CollapseProperties) LogDecisionForNameSupplier(name *GlobalNamespace_Name, messageSupplier func() string) {
	if c.outer.decisionsLog.IsLogging() {
		c.outer.decisionsLog.LogSupplier(func() string { return name.GetFullName() + ": " + messageSupplier() })
	}
}

// FlattenSimpleStubDeclaration flattens a stub declaration. This is mostly a
// hack to support legacy users.
//
// Java: InlineAndCollapseProperties.CollapseProperties.flattenSimpleStubDeclaration(Name, String)
func (c *InlineAndCollapseProperties_CollapseProperties) FlattenSimpleStubDeclaration(name *GlobalNamespace_Name, alias string) {
	// Iterables.getOnlyElement(name.getRefs())
	nRefs := name.refCount()
	if nRefs == 0 {
		panic(jutil.NewNoSuchElementException())
	} else if nRefs > 1 {
		panic(jutil.NewIllegalArgumentExceptionMsg(inlineAndCollapseProperties_onlyElementMessage(name.GetRefs())))
	}
	ref := name.GetFirstRef()
	nameNode := NodeUtil_NewNameWithOriginalName(c.outer.compiler, alias, ref.GetNode(), name.GetFullName())
	varNode := rhino.IR_Var(nameNode).SrcrefIfMissing(nameNode)

	jutil.CheckState(ref.GetNode().GetParent().IsExprResult())
	parent := ref.GetNode().GetParent()
	parent.ReplaceWith(varNode)
	c.outer.compiler.ReportChangeToEnclosingScope(varNode)
}

// inlineAndCollapseProperties_onlyElementMessage is the message of Guava's
// Iterables.getOnlyElement for more than one element: "expected one element
// but was: <e1, e2, ...>" (at most five elements, then ", ...").
//
// Go-only: language-adapter — Guava's message, computed only on failure.
func inlineAndCollapseProperties_onlyElementMessage(refs []*GlobalNamespace_Ref) string {
	var sb strings.Builder
	sb.WriteString("expected one element but was: <")
	for i, ref := range refs {
		if i == 5 {
			sb.WriteString(", ...")
			break
		}
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(ref.String())
	}
	sb.WriteByte('>')
	return sb.String()
}

// ShouldAddConstantName returns whether a NAME node with the given
// collapsedName, for a name originally declared at declarationNode, should
// be annotated with Node#IS_CONSTANT_NAME
//
// Validity checks in unit tests will otherwise throw an exception.
//
// Java: InlineAndCollapseProperties.CollapseProperties.shouldAddConstantName(Node, String)
func (c *InlineAndCollapseProperties_CollapseProperties) ShouldAddConstantName(declarationNode rhino.Node, collapsedName string) bool {
	if declarationNode.IsNil() {
		return false
	}
	info := NodeUtil_GetBestJSDocInfo(declarationNode.GetParent())
	return (info != nil && info.IsConstant()) ||
		declarationNode.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME) ||
		c.outer.compiler.GetCodingConvention().IsConstant(collapsedName)
}

func (c *InlineAndCollapseProperties_CollapseProperties) willFlattenReferences(p *GlobalNamespace_Name) bool {
	if comp := c.outer.compiler; comp != nil && !comp.isParallelWorker() && comp.currentPassIndex <= 1 && !comp.options.isFastMode() {
		return false
	}
	isAllowedToCollapse :=
		c.outer.propertyCollapseLevel != CompilerOptions_PropertyCollapseLevel_MODULE_EXPORT || p.IsModuleExport()
	return isAllowedToCollapse && p.CanCollapseOrInline().CanCollapse()
}

// FlattenReferencesTo flattens all references to a collapsible property of a
// global name except its initial definition.
//
// n is a global property name (e.g. "a.b" or "a.b.c.d"); alias the
// flattened name (e.g. "a$b" or "a$b$c$d").
//
// Java: InlineAndCollapseProperties.CollapseProperties.flattenReferencesTo(Name, String)
func (c *InlineAndCollapseProperties_CollapseProperties) FlattenReferencesTo(n *GlobalNamespace_Name, alias string) {
	c.flattenReferencesTo(n, alias, false)
}

func (c *InlineAndCollapseProperties_CollapseProperties) flattenReferencesTo(
	n *GlobalNamespace_Name, alias string, checkEscaped bool) {
	var originalName string
	var hasOriginalName bool
	getOriginalName := func() string {
		if !hasOriginalName {
			originalName = n.GetFullName()
			hasOriginalName = true
		}
		return originalName
	}
	var declarationNode rhino.Node
	if n.GetDeclaration() != nil {
		declarationNode = n.getDeclarationRef().GetNode()
	}
	isConstantName := c.ShouldAddConstantName(declarationNode, alias)
	for r := range n.Refs() {
		if r == n.getDeclarationRef() {
			// Declarations are handled separately.
			continue
		}
		rParent := r.GetNode().GetParent()
		// We shouldn't flatten a reference that's an object literal key, because duplicate keys
		// show up as refs.
		if !NodeUtil_MayBeObjectLitKey(r.GetNode()) {
			c.FlattenNameRef(alias, r.GetNode(), rParent, getOriginalName(), isConstantName)
		} else if r.GetNode().IsStringKey() && r.GetNode().GetParent().IsObjectPattern() {
			newNode := rhino.IR_Name(alias).Srcref(r.GetNode())
			NodeUtil_CopyNameAnnotations(r.GetNode(), newNode)
			DestructuringGlobalNameExtractor_ReassignDestructringLvalue(
				r.GetNode(), newNode, nil, r, c.outer.compiler)
		}
	}

	// Flatten all occurrences of a name as a prefix of its subnames. For
	// example, if `n` corresponds to the name "a.b", then "a.b" will be
	// replaced with "a$b" in all occurrences of "a.b.c", "a.b.c.d", etc.
	if n.HasProps() {
		descendantsVisited := checkEscaped && !n.IsCollapsingExplicitlyDenied() && !n.isEscaped()
		for p := range n.Props() {
			if flattenPrefixesHasWork(p) {
				c.flattenPrefixesLazy(alias, getOriginalName(), p, 1, 1, isConstantName, descendantsVisited)
			}
		}
	}
}

func flattenPrefixesHasWork(n *GlobalNamespace_Name) bool {
	if n.HasProps() {
		return true
	}
	decl := n.getDeclarationRef()
	if decl != nil && !decl.GetNode().IsNil() && decl.GetNode().IsGetProp() {
		return true
	}
	nRefs := n.refCount()
	return nRefs > 1 || (nRefs == 1 && n.GetFirstRef() != decl)
}

// FlattenPrefixes flattens all occurrences of a name as a prefix of subnames
// beginning with a particular subname.
//
// n is a global property name (e.g. "a.b.c.d"); alias a flattened prefix
// name (e.g. "a$b"); originalName the full original name of the global
// property (e.g. "a.b.c.d") equivalent to n.getFullName(), but pre-computed
// to save on intermediate string allocation; depth the difference in depth
// between the property name and the prefix name (e.g. 2).
//
// Java: InlineAndCollapseProperties.CollapseProperties.flattenPrefixes(String, String, Name, int, boolean)
func (c *InlineAndCollapseProperties_CollapseProperties) FlattenPrefixes(
	alias string, originalName string, n *GlobalNamespace_Name, depth int, isConstantName bool) {
	c.flattenPrefixesLazy(alias, originalName, n, depth, 0, isConstantName, false)
}

func (c *InlineAndCollapseProperties_CollapseProperties) buildLazyPrefixOriginalName(baseOriginalName string, n *GlobalNamespace_Name, extraDepth int) string {
	if extraDepth <= 0 {
		return baseOriginalName
	}
	if extraDepth == 1 {
		return baseOriginalName + n.GetBaseName()
	}
	var stackBuf [16]string
	var parts []string
	if extraDepth <= len(stackBuf) {
		parts = stackBuf[:extraDepth]
	} else {
		parts = make([]string, extraDepth)
	}
	totalLen := len(baseOriginalName)
	cur := n
	for i := extraDepth - 1; i >= 0; i-- {
		parts[i] = cur.GetBaseName()
		totalLen += len(parts[i])
		cur = cur.GetParent()
	}
	var sb strings.Builder
	sb.Grow(totalLen)
	sb.WriteString(baseOriginalName)
	for _, s := range parts {
		sb.WriteString(s)
	}
	return sb.String()
}

func (c *InlineAndCollapseProperties_CollapseProperties) flattenPrefixesLazy(
	alias string,
	baseOriginalName string,
	n *GlobalNamespace_Name,
	depth int,
	extraDepth int,
	isConstantName bool,
	descendantsVisited bool) {
	var cachedOrig string
	hasOrig := false

	decl := n.getDeclarationRef()
	if decl != nil && !decl.GetNode().IsNil() && decl.GetNode().IsGetProp() {
		c.flattenNameRefAtDepthLazy(alias, decl.GetNode(), depth, baseOriginalName, n, extraDepth, &cachedOrig, &hasOrig, isConstantName)
	}

	if descendantsVisited && c.willFlattenReferences(n) {
		return
	}

	for r := range n.Refs() {
		if r == decl {
			continue
		}
		c.flattenNameRefAtDepthLazy(alias, r.GetNode(), depth, baseOriginalName, n, extraDepth, &cachedOrig, &hasOrig, isConstantName)
	}

	if n.HasProps() {
		nextDescendantsVisited := descendantsVisited && !n.IsCollapsingExplicitlyDenied() && !n.isEscaped()
		for p := range n.Props() {
			if flattenPrefixesHasWork(p) {
				c.flattenPrefixesLazy(alias, baseOriginalName, p, depth+1, extraDepth+1, isConstantName, nextDescendantsVisited)
			}
		}
	}
}

func (c *InlineAndCollapseProperties_CollapseProperties) flattenNameRefAtDepthLazy(
	alias string,
	n rhino.Node,
	depth int,
	baseOriginalName string,
	nameObj *GlobalNamespace_Name,
	extraDepth int,
	cachedOrig *string,
	hasOrig *bool,
	isConstantName bool) {
	nType := n.GetToken()
	isQName := nType == rhino.Token_NAME || nType == rhino.Token_GETPROP
	if !isQName {
		jutil.CheckState(NodeUtil_MayBeObjectLitKey(n))
		return
	}
	for i := 1; i < depth && n.HasChildren(); i++ {
		n = n.GetFirstChild()
	}
	if n.IsGetProp() && n.GetFirstChild().IsGetProp() {
		if !*hasOrig {
			*cachedOrig = c.buildLazyPrefixOriginalName(baseOriginalName, nameObj, extraDepth)
			*hasOrig = true
		}
		c.FlattenNameRef(alias, n.GetFirstChild(), n, *cachedOrig, isConstantName)
	}
}

// FlattenNameRefAtDepth flattens a particular prefix of a single name
// reference.
//
// alias is a flattened prefix name (e.g. "a$b"); n the node corresponding to
// a subproperty name (e.g. "a.b.c.d"); depth the difference in depth between
// the property name and the prefix name (e.g. 2); originalName the string
// version of the property name.
//
// Java: InlineAndCollapseProperties.CollapseProperties.flattenNameRefAtDepth(String, Node, int, String, boolean)
func (c *InlineAndCollapseProperties_CollapseProperties) FlattenNameRefAtDepth(
	alias string, n rhino.Node, depth int, originalName string, isConstantName bool) {
	// This method has to work for both GETPROP chains and, in rare cases,
	// OBJLIT keys, possibly nested. That's why we check for children before
	// proceeding. In the OBJLIT case, we don't need to do anything.
	nType := n.GetToken()
	isQName := nType == rhino.Token_NAME || nType == rhino.Token_GETPROP
	isObjKey := NodeUtil_MayBeObjectLitKey(n)
	jutil.CheckState(isObjKey || isQName)
	if isQName {
		for i := 1; i < depth && n.HasChildren(); i++ {
			n = n.GetFirstChild()
		}
		if n.IsGetProp() && n.GetFirstChild().IsGetProp() {
			c.FlattenNameRef(alias, n.GetFirstChild(), n, originalName, isConstantName)
		}
	}
}

// FlattenNameRef replaces a GETPROP a.b.c with a NAME a$b$c.
//
// alias is a flattened prefix name (e.g. "a$b"); n the GETPROP node
// corresponding to the original name (e.g. "a.b"); parent n's parent;
// originalName the string version of the property name; isConstantName
// whether to annotate with Node#IS_CONSTANT_NAME.
//
// Java: InlineAndCollapseProperties.CollapseProperties.flattenNameRef(String, Node, Node, String, boolean)
func (c *InlineAndCollapseProperties_CollapseProperties) FlattenNameRef(
	alias string, n rhino.Node, parent rhino.Node, originalName string, isConstantName bool) {
	jutil.CheckArgumentMsg2(
		n.IsGetProp(), "Expected GETPROP, found %s. Node: %s", n.GetToken(), n)

	// BEFORE:
	//   getprop
	//     getprop
	//       name a
	//       string b
	//     string c
	// AFTER:
	//   name a$b$c
	isConst := isConstantName || n.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME) || c.outer.compiler.GetCodingConvention().IsConstant(alias)
	if NodeUtil_IsNormalOrOptChainCall(parent) && n.IsFirstChildOf(parent) {
		// The node was a call target. We are deliberately flattening these as
		// the "this" isn't provided by the namespace. Mark it as such:
		parent.PutBooleanProp(rhino.Node_FREE_CALL, true)
	}
	if n.MutateGetPropToNameInPlace(alias, originalName, isConst) {
		c.outer.compiler.ReportChangeToEnclosingScope(n)
		return
	}

	ref := NodeUtil_NewNameWithOriginalName(c.outer.compiler, alias, n, originalName).CopyTypeFrom(n)
	NodeUtil_CopyNameAnnotations(n, ref)
	if isConstantName {
		ref.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, isConstantName)
	}
	n.ReplaceWith(ref)
	c.outer.compiler.ReportChangeToEnclosingScope(ref)
}

var (
	collapseFindDynamicImportCallTokens    = makeVisitTokenSet(rhino.Token_CALL)
	collapseFindDynamicImportCjsCallTokens = makeVisitTokenSet(rhino.Token_GETPROP, rhino.Token_CALL)
)

// collapseProperties_findDynamicallyImportedModules finds all the module
// namespace objects which are referenced by a dynamic import.
type collapseProperties_findDynamicallyImportedModules struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *InlineAndCollapseProperties_CollapseProperties

	processCommonJSModules bool
	moduleResolutionMode   ModuleLoader_ResolutionMode
}

var (
	_ nodeTraversal_visitTokens   = (*collapseProperties_findDynamicallyImportedModules)(nil)
	_ nodeTraversal_noVisitLeaves = (*collapseProperties_findDynamicallyImportedModules)(nil)
)

func (f *collapseProperties_findDynamicallyImportedModules) visitTokens() *ast.TokenSet {
	if f.processCommonJSModules {
		return &collapseFindDynamicImportCjsCallTokens
	}
	return &collapseFindDynamicImportCallTokens
}

func (*collapseProperties_findDynamicallyImportedModules) noVisitLeaves() {}

// Java: InlineAndCollapseProperties.CollapseProperties.FindDynamicallyImportedModules(boolean, ResolutionMode)
func newCollapseProperties_findDynamicallyImportedModules(
	outer *InlineAndCollapseProperties_CollapseProperties,
	processCommonJSModules bool,
	resolutionMode ModuleLoader_ResolutionMode) *collapseProperties_findDynamicallyImportedModules {
	return &collapseProperties_findDynamicallyImportedModules{
		outer:                  outer,
		processCommonJSModules: processCommonJSModules,
		moduleResolutionMode:   resolutionMode,
	}
}

// Java: InlineAndCollapseProperties.CollapseProperties.FindDynamicallyImportedModules.visit(NodeTraversal, Node, Node)
func (f *collapseProperties_findDynamicallyImportedModules) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	// After rewriting, CommonJS dynamic imports are of the form
	// __webpack_require__.e(2).then(function() { return module$mod1.default; })
	//
	// Mark the base module name as being dynamically imported.
	if f.processCommonJSModules &&
		n.IsGetProp() &&
		n.IsQualifiedName() &&
		!n.GetParent().IsNil() &&
		n.GetParent().IsReturn() &&
		n.GetGrandparent().IsBlock() &&
		n.GetGrandparent().HasOneChild() &&
		n.GetGrandparent().GetParent().IsFunction() {
		potentialCallback := NodeUtil_GetEnclosingFunction(n)
		if !potentialCallback.IsNil() &&
			ProcessCommonJSModules_IsCommonJsDynamicImportCallback(
				NodeUtil_GetEnclosingFunction(potentialCallback), f.moduleResolutionMode) {
			qname, ok := n.GetQualifiedNameString()
			jutil.CheckState(ok)
			f.outer.addDynamicallyImportedModule(NodeUtil_GetRootOfQualifiedNameString(qname))
		}
	} else if ConvertChunksToESModules_IsDynamicImportCallback(n) {
		moduleNamespace := ConvertChunksToESModules_GetDynamicImportCallbackModuleNamespace(f.outer.outer.compiler, n)
		if !moduleNamespace.IsNil() {
			qname, ok := moduleNamespace.GetQualifiedNameString()
			jutil.CheckState(ok)
			f.outer.addDynamicallyImportedModule(qname)
		}
	}
}

// Java: InlineAndCollapseProperties.isSafeNamespaceReinit(Ref)
func InlineAndCollapseProperties_IsSafeNamespaceReinit(ref *GlobalNamespace_Ref) bool {
	// allow "a = a || {}" or "var a = a || {}" or "var a;"
	valParent := InlineAndCollapseProperties_GetValueParent(ref)
	val := valParent.GetLastChild()
	if !val.IsNil() && val.IsOr() {
		maybeName := val.GetFirstChild()
		if ref.GetNode().MatchesQualifiedNameNode(maybeName) {
			return true
		}
	}

	return false
}

// InlineAndCollapseProperties_GetValueParent gets the parent node of the
// value for any assignment to a Name. For example, in the assignment `var x
// = 3;` the parent would be the NAME node.
//
// Java: InlineAndCollapseProperties.getValueParent(Ref)
func InlineAndCollapseProperties_GetValueParent(ref *GlobalNamespace_Ref) rhino.Node {
	// there are four types of declarations: VARs, LETs, CONSTs, and ASSIGNs
	n := ref.GetNode().GetParent()
	if !n.IsNil() && NodeUtil_IsNameDeclaration(n) {
		return ref.GetNode()
	}
	return ref.GetNode().GetParent()
}
