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
	"strconv"
	"strings"

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

	dynamicallyImportedModules *jutil.LinkedHashSet[string]
}

// Java: InlineAndCollapseProperties.CollapseProperties()
func NewInlineAndCollapseProperties_CollapseProperties(outer *InlineAndCollapseProperties) *InlineAndCollapseProperties_CollapseProperties {
	return &InlineAndCollapseProperties_CollapseProperties{
		outer:                      outer,
		dynamicallyImportedModules: jutil.NewLinkedHashSet[string](),
	}
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
		} else if !c.dynamicallyImportedModules.IsEmpty() && c.dynamicallyImportedModules.Contains(name.GetBaseName()) {
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

	if rvalue.IsObjectLit() {
		for child := rvalue.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
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
	if !c.dynamicallyImportedModules.IsEmpty() {
		// When the output chunk type is ES_MODULES, properties of the module namespace
		// must not be collapsed as they are referenced off the namespace object. The namespace
		// objects escape via the dynamic import expression.
		dynNames := make([]*GlobalNamespace_Name, 0, c.dynamicallyImportedModules.Size())
		for modName := range c.dynamicallyImportedModules.All() {
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
	if comp, ok := c.outer.compiler.(*Compiler); ok && comp.currentPassIndex <= 1 && !comp.options.isFastMode() {
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

// CollapseDeclarationOfNameAndDescendants collapses definitions of the
// collapsible properties of a global name. Recurs on subnames that also
// represent JavaScript objects with collapsible properties.
//
// n is a node representing a global name; alias the flattened name for n.
//
// Java: InlineAndCollapseProperties.CollapseProperties.collapseDeclarationOfNameAndDescendants(Name, String, Set<Name>)
func (c *InlineAndCollapseProperties_CollapseProperties) CollapseDeclarationOfNameAndDescendants(
	n *GlobalNamespace_Name, alias string) {
	childNameInlinability := n.CanCollapseOrInlineChildNames()
	var canCollapseChildNames bool
	if !childNameInlinability.CanCollapse() {
		if c.outer.decisionsLog.IsLogging() {
			c.LogDecisionForNameSupplier(
				n,
				func() string {
					return "child name inlinability: " + childNameInlinability.String() + ": will not collapse child names"
				})
		}
		canCollapseChildNames = false
	} else if n.isEscaped() {
		c.LogDecisionForName(n, "escapes: will not collapse child names")
		canCollapseChildNames = false
	} else {
		canCollapseChildNames = true
	}

	// Handle this name first so that nested object literals get unrolled.
	if c.CanCollapse(n) {
		c.LogDecisionForName(n, "collapsing")
		c.UpdateGlobalNameDeclaration(n, alias, canCollapseChildNames)
	}

	if !n.HasProps() || n.isEscaped() {
		return
	}
	c.LogDecisionForName(n, "collapsing descendants")
	isLogging := c.outer.decisionsLog.IsLogging()
	for p := range n.Props() {
		if !isLogging && !p.HasProps() && !c.CanCollapse(p) {
			continue
		}
		var propAlias string
		if c.CanCollapse(p) || (p.HasProps() && !p.isEscaped()) {
			propAlias = c.AppendPropForAlias(alias, p.GetBaseName())
		}
		c.CollapseDeclarationOfNameAndDescendants(p, propAlias)
	}
}

// UpdateTwinnedDeclaration updates the initial assignment to a collapsible
// property at global scope by adding a VAR stub and collapsing the
// property. e.g. c = a.b = 1; => var a$b; c = a$b = 1; This specifically
// handles "twinned" assignments, which are those where the assignment is
// also used as a reference and which need special handling.
//
// alias is the flattened property name (e.g. "a$b"); refName the name for
// the reference being updated; ref an object containing information about
// the assignment getting updated.
//
// Java: InlineAndCollapseProperties.CollapseProperties.updateTwinnedDeclaration(String, Name, Ref)
func (c *InlineAndCollapseProperties_CollapseProperties) UpdateTwinnedDeclaration(alias string, refName *GlobalNamespace_Name, ref *GlobalNamespace_Ref) {
	jutil.CheckStateMsg1(ref.IsTwin(), "%s", ref)
	// Don't handle declarations of an already flat name, just qualified names.
	if !ref.GetNode().IsGetProp() {
		return
	}
	rvalue := ref.GetNode().GetNext()
	parent := ref.GetNode().GetParent()
	grandparent := parent.GetParent()

	if !rvalue.IsNil() && rvalue.IsFunction() {
		c.CheckForReceiverAffectedByCollapse(rvalue, refName.GetJSDocInfo(), refName)
	}

	// Create the new alias node.
	nameNode := NodeUtil_NewNameWithOriginalName(c.outer.compiler, alias, grandparent.GetFirstChild(), refName.GetFullName())
	NodeUtil_CopyNameAnnotations(ref.GetNode(), nameNode)

	// BEFORE:
	// ... (x.y = 3);
	//
	// AFTER:
	// var x$y;
	// ... (x$y = 3);

	current := grandparent
	currentParent := grandparent.GetParent()
	for ; !currentParent.IsScript() && !currentParent.IsBlock(); current, currentParent = currentParent, currentParent.GetParent() {
	}

	// Create a stub variable declaration right
	// before the current statement.
	stubVar := rhino.IR_Var(nameNode.CloneTree()).SrcrefIfMissing(nameNode)
	stubVar.InsertBefore(current)

	ref.GetNode().ReplaceWith(nameNode)
	c.outer.compiler.ReportChangeToEnclosingScope(nameNode)
}

// UpdateGlobalNameDeclaration updates the first initialization (a.k.a
// "declaration") of a global name. This involves flattening the global name
// (if it's not just a global variable name already), collapsing object
// literal keys into global variables, declaring stub global variables for
// properties added later in a local scope.
//
// It may seem odd that this function also takes care of declaring stubs for
// direct children. The ultimate goal of this function is to eliminate the
// global name entirely (when possible), so that "middlemen" namespaces
// disappear, and to do that we need to make sure that all the direct
// children will be collapsed as well.
//
// n is an object representing a global name (e.g. "a", "a.b.c"); alias the
// flattened name for n (e.g. "a", "a$b$c"); canCollapseChildNames whether
// it's possible to collapse children of this name. (This is mostly passed
// for convenience; it's equivalent to n.canCollapseChildNames()).
//
// Java: InlineAndCollapseProperties.CollapseProperties.updateGlobalNameDeclaration(Name, String, boolean)
func (c *InlineAndCollapseProperties_CollapseProperties) UpdateGlobalNameDeclaration(n *GlobalNamespace_Name, alias string, canCollapseChildNames bool) {
	decl := n.getDeclarationRef()
	if decl == nil {
		// Some names do not have declarations, because they
		// are only defined in local scopes.
		c.LogDecisionForName(n, "no global declaration found")
		return
	}

	isSimpleName := n.IsSimpleName()
	declNode := decl.GetNode()
	switch declNode.GetParent().GetToken() {
	case rhino.Token_ASSIGN:
		c.LogDeclarationAction(n, declNode, "updating assignment")
		c.UpdateGlobalNameDeclarationAtAssignNode(n, alias, canCollapseChildNames)
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		c.LogDeclarationAction(n, declNode, "updating variable declaration")
		c.UpdateGlobalNameDeclarationAtVariableNode(n, canCollapseChildNames)
	case rhino.Token_FUNCTION:
		c.LogDeclarationAction(n, declNode, "updating function declaration")
		c.UpdateGlobalNameDeclarationAtFunctionNode(n, canCollapseChildNames)
	case rhino.Token_CLASS:
		c.LogDeclarationAction(n, declNode, "updating class declaration")
		c.UpdateGlobalNameDeclarationAtClassNode(n, canCollapseChildNames)
	case rhino.Token_CLASS_MEMBERS:
		c.LogDeclarationAction(n, declNode, "updating static member declaration")
		c.UpdateGlobalNameDeclarationAtStaticMemberNode(n, alias, canCollapseChildNames)
	default:
		c.LogDeclarationAction(
			n, declNode, "not updating an unsupported type of declaration node")
	}

	if !isSimpleName {
		commonAncestor := n.GetDeepestCommonAncestorChunk(c.outer.compiler.GetChunkGraph())
		if commonAncestor != decl.GetChunk() {
			c.MoveDeclarationIntoHigherChunk(decl, commonAncestor)
		}
	}
}

// Java: InlineAndCollapseProperties.CollapseProperties.logDeclarationAction(Name, Node, String)
func (c *InlineAndCollapseProperties_CollapseProperties) LogDeclarationAction(name *GlobalNamespace_Name, declarationNode rhino.Node, message string) {
	if c.outer.decisionsLog.IsLogging() {
		c.LogDecisionForNameSupplier(name, func() string { return declarationNode.String() + ": " + message })
	}
}

// Java: InlineAndCollapseProperties.CollapseProperties.moveDeclarationIntoHigherChunk(Ref, JSChunk)
func (c *InlineAndCollapseProperties_CollapseProperties) MoveDeclarationIntoHigherChunk(decl *GlobalNamespace_Ref, commonAncestor *JSChunk) {
	// Var hoisting means that early references are generally ok. E.g. `a$b; var a$b;` is valid.
	// However, if `var a$b;` were in a different, conditionally loaded chunk, then it might not
	// be parsed before `a$b;`. So we need to avoid a ReferenceError when referencing `a$b`.
	newNode := decl.GetNode()
	varNode := newNode.GetParent()
	initializer := newNode.GetFirstChild()
	jutil.CheckStateMsg1(varNode.IsVar(), "Unexpected declaration %s", varNode)

	newRoot := commonAncestor.GetFirst().GetAstRoot(c.outer.compiler)
	if initializer.IsNil() {
		// Already a stub declaration - "var a$b;". Move it to the ancestor chunk.
		newRoot.AddChildToFront(varNode.Detach())
	} else {
		// Define `var a$b;` in the ancestor chunk.
		newVarNode := rhino.IR_Var(newNode.CloneNode()).Srcref(newNode)
		newRoot.AddChildToFront(newVarNode)
		// Replace `var a$b = () => {};` with `a$b = () = {};`. Having multiple `var` declarations
		// of the same name is valid JavaScript but violates Normalize constraints.
		newAssign := rhino.IR_ExprResult(rhino.IR_Assign(newNode.Detach(), initializer.Detach())).
			SrcrefTreeIfMissing(newNode)
		varNode.ReplaceWith(newAssign)
	}
	c.outer.compiler.ReportChangeToChangeScope(newRoot)
}

// UpdateGlobalNameDeclarationAtAssignNode updates the first initialization
// (a.k.a "declaration") of a global name that occurs at an ASSIGN node. See
// comment for updateGlobalNameDeclaration.
//
// n is an object representing a global name (e.g. "a", "a.b.c"); alias the
// flattened name for n (e.g. "a", "a$b$c").
//
// Java: InlineAndCollapseProperties.CollapseProperties.updateGlobalNameDeclarationAtAssignNode(Name, String, boolean)
func (c *InlineAndCollapseProperties_CollapseProperties) UpdateGlobalNameDeclarationAtAssignNode(
	n *GlobalNamespace_Name, alias string, canCollapseChildNames bool) {
	// NOTE: It's important that we don't add additional nodes
	// (e.g. a var node before the exprstmt) because the exprstmt might be
	// the child of an if statement that's not inside a block).

	// All qualified names - even for variables that are initially declared as LETS and CONSTS -
	// are being declared as VAR statements, but this is not incorrect because
	// we are only collapsing for global names.
	ref := n.getDeclarationRef()
	refNode := ref.GetNode()
	rvalue := refNode.GetNext()
	if ref.IsTwin() {
		c.UpdateTwinnedDeclaration(alias, n, ref)
		return
	}
	assignNode := refNode.GetParent()
	grandparent := assignNode.GetParent()
	varParent := grandparent.GetParent()
	isObjLit := rvalue.IsObjectLit()
	insertedVarNode := false
	var varNode rhino.Node

	if isObjLit && c.CanEliminate(n, rvalue) {
		// Eliminate the object literal altogether.
		if grandparent.CanResetInPlace() {
			grandparent.MutateAttachedToFreshVarInPlace(rhino.NilNode, nil)
			varNode = grandparent
		} else {
			varNode = rhino.NewNode(rhino.Token_VAR)
			grandparent.ReplaceWith(varNode)
		}
		n.UpdateRefNode(ref, rhino.NilNode)
		insertedVarNode = true
		c.outer.compiler.ReportChangeToEnclosingScope(varNode)
	} else if !n.IsSimpleName() {
		// Create a VAR node to declare the name.
		if rvalue.IsFunction() {
			c.CheckForReceiverAffectedByCollapse(rvalue, n.GetJSDocInfo(), n)
		}

		c.outer.compiler.ReportChangeToEnclosingScope(rvalue)
		rvalue.Detach()

		isConstName := c.ShouldAddConstantName(refNode, alias)
		info := NodeUtil_GetBestJSDocInfo(assignNode)
		var nameNode rhino.Node
		if assignNode.CanResetInPlace() && grandparent.CanResetInPlace() {
			assignNode.Detach()
			nameNode = assignNode.ResetDetachedToNameWithOriginalNameInPlace(alias, n.GetFullName(), grandparent, isConstName)
			nameNode.AddChildToFront(rvalue)
			grandparent.MutateAttachedToFreshVarInPlace(nameNode, info)
			varNode = grandparent
		} else {
			nameNode = NodeUtil_NewNameWithOriginalName(c.outer.compiler, alias, grandparent, n.GetFullName())
			nameNode.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, isConstName)
			varNode = rhino.NewNode(rhino.Token_VAR)
			if info != nil {
				varNode.SetJSDocInfo(info)
			}
			varNode.AddChildToBack(nameNode)
			nameNode.AddChildToFront(rvalue)
			grandparent.ReplaceWith(varNode)
		}

		// Update the node ancestry stored in the reference.
		n.UpdateRefNode(ref, nameNode)
		insertedVarNode = true
		c.outer.compiler.ReportChangeToEnclosingScope(varNode)
	} else if canCollapseChildNames {
		varNode = rhino.NewNode(rhino.Token_VAR)
	}

	if canCollapseChildNames {
		if isObjLit {
			c.DeclareVariablesForObjLitValues(n, alias, rvalue, varNode, varNode.GetPrevious())
		}

		c.AddStubsForUndeclaredProperties(n, alias, varParent, varNode)
	}

	if insertedVarNode {
		if !varNode.HasChildren() {
			varNode.Detach()
		}
	}
}

// CheckForReceiverAffectedByCollapse warns about any references to "this" in
// the given FUNCTION. The function is getting collapsed, so the references
// will change.
//
// Java: InlineAndCollapseProperties.CollapseProperties.checkForReceiverAffectedByCollapse(Node, JSDocInfo, Name)
func (c *InlineAndCollapseProperties_CollapseProperties) CheckForReceiverAffectedByCollapse(function rhino.Node, docInfo *rhino.JSDocInfo, name *GlobalNamespace_Name) {
	jutil.CheckState(function.IsFunction())

	if docInfo != nil {
		// Don't rely on type inference for this check.

		if docInfo.IsConstructorOrInterface() {
			return // Ctors and interfaces need to be able to reference `this`
		}
		if docInfo.HasThisType() {
			/*
			 * Use `@this` as a signal that the reference is intentional.
			 *
			 * <p>TODO(b/156823102): This signal also silences the check on all transpiled static
			 * methods.
			 */
			return
		}
	}

	// Use the NodeUtil method so we don't forget to update this logic.
	if NodeUtil_ReferencesOwnReceiver(function) {
		c.outer.compiler.Report(JSError_MakeWithNode(function, InlineAndCollapseProperties_RECEIVER_AFFECTED_BY_COLLAPSE, name.GetFullName()))
	}
}

// UpdateGlobalNameDeclarationAtVariableNode updates the first
// initialization (a.k.a "declaration") of a global name that occurs at a VAR
// node. See comment for updateGlobalNameDeclaration.
//
// n is an object representing a global name (e.g. "a").
//
// Java: InlineAndCollapseProperties.CollapseProperties.updateGlobalNameDeclarationAtVariableNode(Name, boolean)
func (c *InlineAndCollapseProperties_CollapseProperties) UpdateGlobalNameDeclarationAtVariableNode(n *GlobalNamespace_Name, canCollapseChildNames bool) {
	if !canCollapseChildNames {
		c.LogDecisionForName(n, "cannot collapse child names: skipping")
		return
	}

	ref := n.getDeclarationRef()
	name := ref.GetNode().GetString()
	rvalue := ref.GetNode().GetFirstChild()
	variableNode := ref.GetNode().GetParent()
	grandparent := variableNode.GetParent()

	isObjLit := rvalue.IsObjectLit()

	if isObjLit {
		c.DeclareVariablesForObjLitValues(n, name, rvalue, variableNode, variableNode.GetPrevious())
	}

	c.AddStubsForUndeclaredProperties(n, name, grandparent, variableNode)

	if isObjLit && c.CanEliminate(n, rvalue) {
		ref.GetNode().Detach()
		c.outer.compiler.ReportChangeToEnclosingScope(variableNode)
		if !variableNode.HasChildren() {
			variableNode.Detach()
		}

		// Clear out the object reference, since we've eliminated it from the
		// parse tree.
		n.UpdateRefNode(ref, rhino.NilNode)
	}
}

// UpdateGlobalNameDeclarationAtFunctionNode updates the first
// initialization (a.k.a "declaration") of a global name that occurs at a
// FUNCTION node. See comment for updateGlobalNameDeclaration.
//
// n is an object representing a global name (e.g. "a").
//
// Java: InlineAndCollapseProperties.CollapseProperties.updateGlobalNameDeclarationAtFunctionNode(Name, boolean)
func (c *InlineAndCollapseProperties_CollapseProperties) UpdateGlobalNameDeclarationAtFunctionNode(n *GlobalNamespace_Name, canCollapseChildNames bool) {
	if !canCollapseChildNames || !c.CanCollapse(n) {
		return
	}

	ref := n.getDeclarationRef()
	fnName := ref.GetNode().GetString()
	c.AddStubsForUndeclaredProperties(
		n, fnName, ref.GetNode().GetAncestor(2), ref.GetNode().GetParent())
}

// UpdateGlobalNameDeclarationAtClassNode updates the first initialization
// (a.k.a "declaration") of a global name that occurs at a CLASS node. See
// comment for updateGlobalNameDeclaration.
//
// n is an object representing a global name (e.g. "a").
//
// Java: InlineAndCollapseProperties.CollapseProperties.updateGlobalNameDeclarationAtClassNode(Name, boolean)
func (c *InlineAndCollapseProperties_CollapseProperties) UpdateGlobalNameDeclarationAtClassNode(n *GlobalNamespace_Name, canCollapseChildNames bool) {
	if !canCollapseChildNames || !c.CanCollapse(n) {
		return
	}

	ref := n.getDeclarationRef()
	className := ref.GetNode().GetString()
	c.AddStubsForUndeclaredProperties(
		n, className, ref.GetNode().GetAncestor(2), ref.GetNode().GetParent())
}

// UpdateGlobalNameDeclarationAtStaticMemberNode updates the first
// initialization (a.k.a "declaration") of a global name that occurs in a
// static MEMBER_FUNCTION_DEF or MEMBER_FIELD_DEF in a class. See comment for
// updateGlobalNameDeclaration.
//
// n is a static MEMBER_FUNCTION_DEF or MEMBER_FIELD_DEF in a class assigned
// to a global name (e.g. `a.b`); alias the new flattened name for `n` (e.g.
// "a$b"); canCollapseChildNames whether properties of `n` are also
// collapsible, meaning that any properties only assigned locally need stub
// declarations.
//
// Java: InlineAndCollapseProperties.CollapseProperties.updateGlobalNameDeclarationAtStaticMemberNode(Name, String, boolean)
func (c *InlineAndCollapseProperties_CollapseProperties) UpdateGlobalNameDeclarationAtStaticMemberNode(
	n *GlobalNamespace_Name, alias string, canCollapseChildNames bool) {

	declaration := n.getDeclarationRef()
	classNode := declaration.GetNode().GetGrandparent()
	jutil.CheckStateMsg1(classNode.IsClass(), "%s", classNode)
	enclosingStatement := NodeUtil_GetEnclosingStatement(classNode)

	if canCollapseChildNames {
		c.AddStubsForUndeclaredProperties(n, alias, enclosingStatement.GetParent(), classNode)
	}

	// detach `static m() {}` from `class Foo { static m() {} }`
	member := declaration.GetNode().Detach()
	var varDecl rhino.Node
	if member.HasChildren() {
		memberVal := member.GetOnlyChild()
		c.CheckForReceiverAffectedByCollapse(memberVal, member.GetJSDocInfo(), n)
		memberVal = memberVal.Detach()

		// add a var declaration, creating `var Foo$m = function() {}; class Foo {}`
		varDecl = rhino.IR_VarWithValue(NodeUtil_NewNameWithSrcref(c.outer.compiler, alias, member), memberVal).Srcref(member)
	} else {
		// Member field defs may not have initializers.
		varDecl = rhino.IR_Var(NodeUtil_NewNameWithSrcref(c.outer.compiler, alias, member)).Srcref(member)
	}
	varDecl.InsertBefore(enclosingStatement)
	// We would lose optimization-relevant jsdoc tags here because they are stored on the class
	// member node, not the function node. Copy them over to the new declaration statement so
	// later passes can make use of them.
	varDecl.SetJSDocInfo(member.GetJSDocInfo())
	c.outer.compiler.ReportChangeToEnclosingScope(varDecl)

	// collapsing this name's properties requires updating this Ref
	n.UpdateRefNode(declaration, varDecl.GetFirstChild())
}

// DeclareVariablesForObjLitValues declares global variables to serve as
// aliases for the values in an object literal, optionally removing all of
// the object literal's keys and values.
//
// alias is the object literal's flattened name (e.g. "a$b$c"); objlit the
// OBJLIT node; varNode the VAR node to which new global variables should be
// added as children; nameToAddAfter the child of varNode after which new
// variables should be added (may be null).
//
// Java: InlineAndCollapseProperties.CollapseProperties.declareVariablesForObjLitValues(Name, String, Node, Node, Node)
func (c *InlineAndCollapseProperties_CollapseProperties) DeclareVariablesForObjLitValues(
	objlitName *GlobalNamespace_Name, alias string, objlit rhino.Node, varNode rhino.Node, nameToAddAfter rhino.Node) {
	arbitraryNameCounter := 0
	discardKeys := !objlitName.ShouldKeepKeys()

	var nextKey rhino.Node
	for key := objlit.GetFirstChild(); !key.IsNil(); key = nextKey {
		value := key.GetFirstChild()
		nextKey = key.GetNext()

		// A computed property, or a get or a set can not be rewritten as a VAR. We don't know what
		// properties will be generated by a spread.
		switch key.GetToken() {
		case rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_COMPUTED_PROP, rhino.Token_OBJECT_SPREAD:
			continue
		case rhino.Token_STRING_KEY, rhino.Token_MEMBER_FUNCTION_DEF:
		default:
			panic(jutil.NewIllegalStateExceptionMsg(
				"Unexpected child of OBJECTLIT: " + key.ToStringTree()))
		}

		// We generate arbitrary names for keys that aren't valid JavaScript
		// identifiers, since those keys are never referenced. (If they were,
		// this object literal's child names wouldn't be collapsible.) The only
		// reason that we don't eliminate them entirely is the off chance that
		// their values are expressions that have side effects.
		isJsIdentifier := !key.IsNumber() && rhino.TokenStream_IsJSIdentifier(key.GetString())
		var propName string
		if isJsIdentifier {
			propName = key.GetString()
		} else {
			arbitraryNameCounter++
			propName = strconv.Itoa(arbitraryNameCounter)
		}

		// If the name cannot be collapsed, skip it.
		p := objlitName.LookupChild(propName)
		if p != nil && !c.CanCollapse(p) {
			continue
		}

		propAlias := c.AppendPropForAlias(alias, propName)
		var refNode rhino.Node
		isConstantName := c.ShouldAddConstantName(key, propAlias)
		keyDoc := key.GetJSDocInfo()
		if discardKeys {
			key.Detach()
			value.Detach()
			// Don't report a change here because the objlit has already been removed from the tree.
		} else {
			// Substitute a reference for the value.
			refNode = rhino.IR_Name(propAlias)
			if isConstantName {
				refNode.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
			}

			value.ReplaceWith(refNode)
			c.outer.compiler.ReportChangeToEnclosingScope(refNode)
		}

		// Declare the collapsed name as a variable with the original value.
		nameNode := rhino.IR_Name(propAlias)
		nameNode.AddChildToFront(value)
		if isConstantName {
			nameNode.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
		}
		var newVar rhino.Node
		if discardKeys && key.CanResetInPlace() && (p == nil || p.getDeclarationRef() == nil || p.getDeclarationRef().GetNode() == key) {
			nameNode.SrcrefTreeIfMissing(key)
			key.MutateAttachedStmtToVarFromSrcref(nameNode, key, nil)
			newVar = key
		} else {
			newVar = rhino.IR_Var(nameNode).SrcrefTreeIfMissing(key)
		}
		if !nameToAddAfter.IsNil() {
			newVar.InsertAfter(nameToAddAfter)
		} else {
			newVar.InsertBefore(varNode)
		}
		c.outer.compiler.ReportChangeToEnclosingScope(newVar)
		nameToAddAfter = newVar

		// Update the global name's node ancestry if it hasn't already been
		// done. (Duplicate keys in an object literal can bring us here twice
		// for the same global name.)
		if isJsIdentifier && p != nil {
			if !discardKeys {
				p.AddAliasingGetClonedFromDeclaration(refNode)
			}

			p.UpdateRefNode(p.getDeclarationRef(), nameNode)

			if value.IsFunction() {
				c.CheckForReceiverAffectedByCollapse(value, keyDoc, p)
			}
		}
	}
}

// AddStubsForUndeclaredProperties adds global variable "stubs" for any
// properties of a global name that are only set in a local scope or read
// but never set.
//
// n is an object representing a global name (e.g. "a", "a.b.c"); alias the
// flattened name of the object whose properties we are adding stubs for
// (e.g. "a$b$c"); parent the node to which new global variables should be
// added as children; addAfter the child of after which new variables should
// be added.
//
// Java: InlineAndCollapseProperties.CollapseProperties.addStubsForUndeclaredProperties(Name, String, Node, Node)
func (c *InlineAndCollapseProperties_CollapseProperties) AddStubsForUndeclaredProperties(
	n *GlobalNamespace_Name, alias string, parent rhino.Node, addAfter rhino.Node) {
	jutil.CheckStateMsg1(n.CanCollapseUnannotatedChildNames(), "%s", n)
	jutil.CheckArgumentMsg1(NodeUtil_IsStatementBlock(parent), "%s", parent)
	rhino.CheckNotNull(addAfter)
	if !n.HasProps() {
		return
	}
	for p := range n.Props() {
		if !p.NeedsToBeStubbed() {
			continue
		}

		propAlias := c.AppendPropForAlias(alias, p.GetBaseName())
		nameNode := rhino.IR_Name(propAlias)
		newVar := rhino.IR_Var(nameNode).SrcrefTreeIfMissing(addAfter)
		newVar.InsertAfter(addAfter)

		// Determine if this is a constant var by checking the first
		// reference to it. Don't check the declaration, as it might be null.
		constPropNode := p.GetFirstRef().GetNode()
		isConstantName :=
			// Don't call shouldAddConstantName because it expects to see a declaration node for
			// a name, and this name is, by definition, undeclared.
			constPropNode.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME) ||
				c.outer.compiler.GetCodingConvention().IsConstant(propAlias)
		if isConstantName {
			nameNode.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
		}

		c.outer.compiler.ReportChangeToEnclosingScope(newVar)
		addAfter = newVar
	}
}

// Java: InlineAndCollapseProperties.CollapseProperties.appendPropForAlias(String, String)
func (c *InlineAndCollapseProperties_CollapseProperties) AppendPropForAlias(root string, prop string) string {
	if strings.IndexByte(prop, '$') != -1 {
		// Encode '$' in a property as '$0'. Because '0' cannot be the
		// start of an identifier, this will never conflict with our
		// encoding from '.' -> '$'.
		prop = strings.ReplaceAll(prop, "$", "$0")
	}
	result := root + "$" + prop
	id := 1
	for c.namespace.ContainsName(result) {
		result = root + "$" + prop + "$" + strconv.Itoa(id)
		id++
	}
	return result
}

// collapseProperties_findDynamicallyImportedModules finds all the module
// namespace objects which are referenced by a dynamic import.
type collapseProperties_findDynamicallyImportedModules struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *InlineAndCollapseProperties_CollapseProperties

	processCommonJSModules bool
	moduleResolutionMode   ModuleLoader_ResolutionMode
}

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
			f.outer.dynamicallyImportedModules.Add(NodeUtil_GetRootOfQualifiedNameString(*n.GetQualifiedName()))
		}
	} else if ConvertChunksToESModules_IsDynamicImportCallback(n) {
		moduleNamespace := ConvertChunksToESModules_GetDynamicImportCallbackModuleNamespace(f.outer.outer.compiler, n)
		if !moduleNamespace.IsNil() {
			f.outer.dynamicallyImportedModules.Add(*moduleNamespace.GetQualifiedName())
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
