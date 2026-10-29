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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var (
	inlineAliasesCollectorVisitTokens = func() ast.TokenSet {
		var s ast.TokenSet
		s.Add(rhino.Token_VAR)
		s.Add(rhino.Token_CONST)
		s.Add(rhino.Token_LET)
		s.Add(rhino.Token_ASSIGN)
		return s
	}()
	inlineAliasesInlinerVisitTokens = func() ast.TokenSet {
		var s ast.TokenSet
		s.Add(rhino.Token_NAME)
		s.Add(rhino.Token_GETPROP)
		return s
	}()
)

// InlineAndCollapseProperties_InlineAliases inlines constant aliases
//
// This pass was originally necessary because typechecking did not handle
// type aliases well. Now typechecking understands type aliases. In theory,
// this pass can be deleted, but in practice this pass affects some check
// passes that run post-typechecking.
//
// This alias inliner is not very aggressive. It will only inline explicitly
// const aliases but not effectively const ones (for example ones that are
// only ever assigned a value once). This is done to be conservative since
// it's not a good idea to be making dramatic AST changes during checks (or
// really, any AST changes at all). There is a more aggressive alias inliner
// that runs at the start of optimization.
//
// TODO(b/124915436): Delete this pass.
type InlineAndCollapseProperties_InlineAliases struct {
	outer *InlineAndCollapseProperties

	aliases    *jutil.LinkedHashMap[string, string]
	namespace  *GlobalNamespace
	astFactory *AstFactory
}

// Java: InlineAndCollapseProperties.InlineAliases()
func NewInlineAndCollapseProperties_InlineAliases(outer *InlineAndCollapseProperties) *InlineAndCollapseProperties_InlineAliases {
	a := &InlineAndCollapseProperties_InlineAliases{
		outer:   outer,
		aliases: jutil.NewLinkedHashMap[string, string](),
	}
	a.astFactory = outer.compiler.CreateAstFactory()
	return a
}

// Java: InlineAndCollapseProperties.InlineAliases.process(Node, Node)
func (a *InlineAndCollapseProperties_InlineAliases) Process(externs rhino.Node, root rhino.Node) {
	a.namespace = NewGlobalNamespaceWithRoots(a.outer.compiler, externs, root)
	defer func() {
		a.namespace = nil
	}()
	NodeTraversal_TraverseRoots(a.outer.compiler, NewInlineAndCollapseProperties_InlineAliases_AliasesCollector(a), externs, root)
	if !a.aliases.IsEmpty() {
		NodeTraversal_TraverseRoots(a.outer.compiler, NewInlineAndCollapseProperties_InlineAliases_AliasesInliner(a), externs, root)
	}
}

// InlineAndCollapseProperties_InlineAliases_AliasesCollector is
// InlineAliases.AliasesCollector.
type InlineAndCollapseProperties_InlineAliases_AliasesCollector struct {
	NodeTraversal_ExternsSkippingCallbackBase
	outer *InlineAndCollapseProperties_InlineAliases
}

// Java: InlineAndCollapseProperties.InlineAliases.AliasesCollector()
func NewInlineAndCollapseProperties_InlineAliases_AliasesCollector(outer *InlineAndCollapseProperties_InlineAliases) *InlineAndCollapseProperties_InlineAliases_AliasesCollector {
	return &InlineAndCollapseProperties_InlineAliases_AliasesCollector{outer: outer}
}

var (
	_ nodeTraversal_noVisitLeaves  = (*InlineAndCollapseProperties_InlineAliases_AliasesCollector)(nil)
	_ nodeTraversal_visitTokens    = (*InlineAndCollapseProperties_InlineAliases_AliasesCollector)(nil)
	_ nodeTraversal_onlyNameLeaves = (*InlineAndCollapseProperties_InlineAliases_AliasesInliner)(nil)
	_ nodeTraversal_visitTokens    = (*InlineAndCollapseProperties_InlineAliases_AliasesInliner)(nil)
)

func (*InlineAndCollapseProperties_InlineAliases_AliasesCollector) noVisitLeaves() {}

func (*InlineAndCollapseProperties_InlineAliases_AliasesCollector) visitTokens() *ast.TokenSet {
	return &inlineAliasesCollectorVisitTokens
}

// Java: InlineAndCollapseProperties.InlineAliases.AliasesCollector.visit(NodeTraversal, Node, Node)
func (c *InlineAndCollapseProperties_InlineAliases_AliasesCollector) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_VAR, rhino.Token_CONST, rhino.Token_LET:
		if n.HasOneChild() && t.InGlobalScope() {
			c.VisitAliasDefinition(n.GetFirstChild(), NodeUtil_GetBestJSDocInfo(n.GetFirstChild()))
		}
	case rhino.Token_ASSIGN:
		if !parent.IsNil() && parent.IsExprResult() && t.InGlobalScope() {
			c.VisitAliasDefinition(n.GetFirstChild(), n.GetJSDocInfo())
		}
	default:
	}
}

// VisitAliasDefinition maybe records that given lvalue is an alias of the
// qualified name on its rhs. Note that since we are doing a post-order
// traversal, any previous aliases contained in the rhs will have already
// been substituted by the time we record the new alias.
//
// Java: InlineAndCollapseProperties.InlineAliases.AliasesCollector.visitAliasDefinition(Node, JSDocInfo)
func (c *InlineAndCollapseProperties_InlineAliases_AliasesCollector) VisitAliasDefinition(lhs rhino.Node, info *rhino.JSDocInfo) {
	if c.IsDeclaredConst(lhs, info) &&
		(info == nil || !info.HasTypeInformation()) {
		if lhsQName, ok := lhs.GetQualifiedNameString(); ok {
			rhs := NodeUtil_GetRValueOfLValue(lhs)
			if rhsQName, ok := rhs.GetQualifiedNameString(); ok {
				lhsName, _ := c.outer.namespace.GetOwnSlot(lhsQName).(*GlobalNamespace_Name)
				rhsName, _ := c.outer.namespace.GetOwnSlot(rhsQName).(*GlobalNamespace_Name)
				if lhsName != nil &&
					lhsName.CalculateInlinability().ShouldInlineUsages() &&
					rhsName != nil &&
					rhsName.CalculateInlinability().ShouldInlineUsages() {
					c.outer.aliases.Put(lhsQName, rhsQName)
				}
			}
		}
	}
}

// Java: InlineAndCollapseProperties.InlineAliases.AliasesCollector.isDeclaredConst(Node, JSDocInfo)
func (c *InlineAndCollapseProperties_InlineAliases_AliasesCollector) IsDeclaredConst(lhs rhino.Node, info *rhino.JSDocInfo) bool {
	if info != nil && info.HasConstAnnotation() {
		return true
	}
	return lhs.GetParent().IsConst()
}

// InlineAndCollapseProperties_InlineAliases_AliasesInliner is
// InlineAliases.AliasesInliner.
type InlineAndCollapseProperties_InlineAliases_AliasesInliner struct {
	NodeTraversal_ExternsSkippingCallbackBase
	outer *InlineAndCollapseProperties_InlineAliases
}

// Java: InlineAndCollapseProperties.InlineAliases.AliasesInliner()
func NewInlineAndCollapseProperties_InlineAliases_AliasesInliner(outer *InlineAndCollapseProperties_InlineAliases) *InlineAndCollapseProperties_InlineAliases_AliasesInliner {
	return &InlineAndCollapseProperties_InlineAliases_AliasesInliner{outer: outer}
}

func (*InlineAndCollapseProperties_InlineAliases_AliasesInliner) onlyNameLeaves() {}

func (*InlineAndCollapseProperties_InlineAliases_AliasesInliner) visitTokens() *ast.TokenSet {
	return &inlineAliasesInlinerVisitTokens
}

// Java: InlineAndCollapseProperties.InlineAliases.AliasesInliner.visit(NodeTraversal, Node, Node)
func (i *InlineAndCollapseProperties_InlineAliases_AliasesInliner) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_NAME, rhino.Token_GETPROP:
		if i.outer.aliases.IsEmpty() {
			return
		}
		qname, ok := n.GetQualifiedNameString()
		if ok && i.outer.aliases.ContainsKey(qname) {
			if i.IsLeftmostNameLocal(t, n) {
				// The alias is shadowed by a local variable. Don't rewrite.
				return
			}
			if NodeUtil_IsNameDeclOrSimpleAssignLhs(n, parent) {
				// The node defines an alias. Don't rewrite.
				return
			}

			newNode := i.outer.astFactory.CreateQName(i.outer.namespace, i.ResolveAlias(qname, n))
			if i.IsLeftmostNameLocal(t, newNode) {
				// The aliased name is shadowed by a local variable. Don't rewrite.
				return
			}

			// If n is get_prop like "obj.foo" then newNode should use only location of foo, not
			// obj.foo.
			newNode.SrcrefTree(n)
			// Similarly if n is get_prop like "obj.foo" we should index only foo. obj should not
			// be indexed as it's invisible to users.
			if newNode.IsGetProp() {
				newNode.GetFirstChild().MakeNonIndexableRecursive()
			}
			n.ReplaceWith(newNode)
			t.ReportCodeChange()
		}
	default:
	}
}

// Java: InlineAndCollapseProperties.InlineAliases.AliasesInliner.isLeftmostNameLocal(NodeTraversal, Node)
func (i *InlineAndCollapseProperties_InlineAliases_AliasesInliner) IsLeftmostNameLocal(t *NodeTraversal, n rhino.Node) bool {
	jutil.CheckState(n.IsQualifiedName())
	leftmostName := NodeUtil_GetRootOfQualifiedName(n).GetString()
	v := t.LookupVar(leftmostName)
	return v != nil && v.IsLocal()
}

// ResolveAlias uses the alias table to look up the resolved name of the
// given alias. If the result is also an alias repeat until the real name is
// resolved.
//
// Java: InlineAndCollapseProperties.InlineAliases.AliasesInliner.resolveAlias(String, Node)
func (i *InlineAndCollapseProperties_InlineAliases_AliasesInliner) ResolveAlias(name string, n rhino.Node) string {
	aliasPath := jutil.NewLinkedHashSet[string]()
	for i.outer.aliases.ContainsKey(name) {
		if !aliasPath.Add(name) {
			i.outer.outer.compiler.Report(JSError_MakeWithNode(n, InlineAndCollapseProperties_ALIAS_CYCLE, aliasPath.String(), name))

			// Cut the cycle so that it doesn't get reported more than once.
			i.outer.aliases.Remove(name)
			break
		}

		name = i.outer.aliases.Get(name)
	}
	return name
}
