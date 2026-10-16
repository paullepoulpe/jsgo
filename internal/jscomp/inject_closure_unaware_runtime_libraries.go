// Port of com/google/javascript/jscomp/InjectClosureUnawareRuntimeLibraries.java @ closure-compiler v20260415

/*
 * Copyright 2025 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// InjectClosureUnawareRuntimeLibraries takes the set of runtime libraries
// and fields used in a @closureUnaware nested compilation, and then 1)
// ensures the library definitions are injected in the main AST and 2) passes
// all specific $jscomp.* fields used into the @closureUnaware shadow AST as
// an argument.
//
// Before:
//
//	/** @closureUnaware * /
//	(function() {
//	  use($jscomp_someField$$);
//	}).call(undefined);
//
// After:
//
//	/** @closureUnaware * /
//	(function($jscomp_someField$$) {
//	  use($jscomp_someField$$);
//	}).call(undefined, $jscomp.someField);
type InjectClosureUnawareRuntimeLibraries struct {
	mainCompiler        AbstractCompiler
	shadowCompiler      AbstractCompiler
	closureUnawareCalls []*InjectClosureUnawareRuntimeLibraries_ClosureUnawareCallSite
}

// InjectClosureUnawareRuntimeLibraries_ClosureUnawareCallSite is the Java
// record ClosureUnawareCallSite(Node callNode, Node closureUnawareFunction).
type InjectClosureUnawareRuntimeLibraries_ClosureUnawareCallSite struct {
	callNode               *rhino.Node
	closureUnawareFunction *rhino.Node
}

// Java: ClosureUnawareCallSite(Node, Node)
func NewInjectClosureUnawareRuntimeLibraries_ClosureUnawareCallSite(
	callNode *rhino.Node, closureUnawareFunction *rhino.Node) *InjectClosureUnawareRuntimeLibraries_ClosureUnawareCallSite {
	jutil.CheckArgumentMsg1(callNode.IsCall(), "%s", callNode)
	jutil.CheckArgumentMsg1(closureUnawareFunction.IsFunction(), "%s", closureUnawareFunction)
	return &InjectClosureUnawareRuntimeLibraries_ClosureUnawareCallSite{
		callNode:               callNode,
		closureUnawareFunction: closureUnawareFunction,
	}
}

// Java: ClosureUnawareCallSite.callNode()
func (c *InjectClosureUnawareRuntimeLibraries_ClosureUnawareCallSite) CallNode() *rhino.Node {
	return c.callNode
}

// Java: ClosureUnawareCallSite.closureUnawareFunction()
func (c *InjectClosureUnawareRuntimeLibraries_ClosureUnawareCallSite) ClosureUnawareFunction() *rhino.Node {
	return c.closureUnawareFunction
}

// Java: InjectClosureUnawareRuntimeLibraries(AbstractCompiler, AbstractCompiler, List<ClosureUnawareCallSite>)
func NewInjectClosureUnawareRuntimeLibraries(
	mainCompiler AbstractCompiler,
	shadowCompiler AbstractCompiler,
	closureUnawareCalls []*InjectClosureUnawareRuntimeLibraries_ClosureUnawareCallSite) *InjectClosureUnawareRuntimeLibraries {
	i := &InjectClosureUnawareRuntimeLibraries{}
	i.mainCompiler = mainCompiler
	i.shadowCompiler = shadowCompiler
	i.closureUnawareCalls = closureUnawareCalls
	return i
}

// Java: InjectClosureUnawareRuntimeLibraries.process(Node, Node)
func (i *InjectClosureUnawareRuntimeLibraries) Process(externs *rhino.Node, root *rhino.Node) {
	i.InjectLibrariesIntoMainAst()
	i.InjectFieldsIntoCallSites()
}

// Java: InjectClosureUnawareRuntimeLibraries.injectLibrariesIntoMainAst()
func (i *InjectClosureUnawareRuntimeLibraries) InjectLibrariesIntoMainAst() {
	mainManager := i.mainCompiler.GetRuntimeJsLibManager()
	shadowManager := i.shadowCompiler.GetRuntimeJsLibManager()
	for _, runtimeLib := range shadowManager.GetInjectedLibraries() {
		mainManager.EnsureLibraryInjected(runtimeLib, false /* force= */)
	}
}

// Java: InjectClosureUnawareRuntimeLibraries.injectFieldsIntoCallSites()
func (i *InjectClosureUnawareRuntimeLibraries) InjectFieldsIntoCallSites() {
	// stream().collect(toImmutableMap(ExternedField::qualifiedName, identity()))
	fieldsByNameBuilder := jutil.NewImmutableMapBuilder[string, RuntimeJsLibManager_ExternedField]()
	for _, field := range i.shadowCompiler.GetRuntimeJsLibManager().GetExternedFields() {
		fieldsByNameBuilder.Put(field.QualifiedName(), field)
	}
	fieldsByName := fieldsByNameBuilder.BuildOrThrow()

	for _, callSite := range i.closureUnawareCalls {
		toInject := InjectClosureUnawareRuntimeLibraries_GatherRuntimeFields(callSite.ClosureUnawareFunction(), fieldsByName)
		i.InjectAll(callSite, toInject)
	}
}

// Java: InjectClosureUnawareRuntimeLibraries.injectAll(ClosureUnawareCallSite, Set<ExternedField>)
func (i *InjectClosureUnawareRuntimeLibraries) InjectAll(
	callSite *InjectClosureUnawareRuntimeLibraries_ClosureUnawareCallSite,
	toInject *jutil.LinkedHashSet[RuntimeJsLibManager_ExternedField]) {
	if toInject.IsEmpty() {
		return
	}
	parameterParent := NodeUtil_GetFunctionParameters(callSite.ClosureUnawareFunction())
	argumentParent := callSite.CallNode()

	if parameterParent.HasChildren() && parameterParent.GetLastChild().IsRest() {
		panic(jutil.NewIllegalStateExceptionMsg(
			"rest parameters not allowed in closureUnaware function parameters, found " +
				parameterParent.GetLastChild().String()))
	}

	for field := range toInject.All() {
		parameter := rhino.IR_Name(field.QualifiedName()).Srcref(parameterParent)
		argument := NodeUtil_NewQName(i.mainCompiler, field.UncompiledName()).SrcrefTree(argumentParent)

		parameterParent.AddChildToBack(parameter)
		i.shadowCompiler.ReportChangeToEnclosingScope(parameterParent)

		argumentParent.AddChildToBack(argument)
		i.mainCompiler.ReportChangeToEnclosingScope(argumentParent)
	}
}

// InjectClosureUnawareRuntimeLibraries_GatherRuntimeFields returns the
// subset of allInjectedFields that are actually referenced in the given
// shadowAst.
//
// Java: InjectClosureUnawareRuntimeLibraries.gatherRuntimeFields(Node, ImmutableMap<String, ExternedField>)
func InjectClosureUnawareRuntimeLibraries_GatherRuntimeFields(
	root *rhino.Node, allKnownFields *jutil.LinkedHashMap[string, RuntimeJsLibManager_ExternedField]) *jutil.LinkedHashSet[RuntimeJsLibManager_ExternedField] {
	seen := jutil.NewLinkedHashSet[RuntimeJsLibManager_ExternedField]()
	visitor := NodeUtil_VisitorFunc(func(n *rhino.Node) {
		if !n.IsName() {
			return
		}
		field := allKnownFields.Get(n.GetString())
		if field != nil {
			seen.Add(field)
		}
	})
	NodeUtil_VisitPreOrder(root, visitor)
	return seen
}
