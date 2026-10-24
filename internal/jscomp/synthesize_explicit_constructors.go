// Port of com/google/javascript/jscomp/SynthesizeExplicitConstructors.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// SynthesizeExplicitConstructors adds explicit constructors to classes that
// lack them.
type SynthesizeExplicitConstructors struct {
	compiler   AbstractCompiler
	astFactory *AstFactory
}

// Java: SynthesizeExplicitConstructors(AbstractCompiler)
func NewSynthesizeExplicitConstructors(compiler AbstractCompiler) *SynthesizeExplicitConstructors {
	s := &SynthesizeExplicitConstructors{}
	s.compiler = compiler
	s.astFactory = compiler.CreateAstFactory()
	return s
}

// SynthesizeClassConstructorIfMissing adds a synthetic constructor to the
// given class if no constructor is present.
//
// Java: SynthesizeExplicitConstructors.synthesizeClassConstructorIfMissing(NodeTraversal, Node)
func (s *SynthesizeExplicitConstructors) SynthesizeClassConstructorIfMissing(t *NodeTraversal, classNode rhino.Node) {
	jutil.CheckArgument(classNode.IsClass())
	if NodeUtil_GetEs6ClassConstructorMemberFunctionDef(classNode).IsNil() {
		s.AddSyntheticConstructor(t, classNode)
	}
}

// Java: SynthesizeExplicitConstructors.addSyntheticConstructor(NodeTraversal, Node)
func (s *SynthesizeExplicitConstructors) AddSyntheticConstructor(t *NodeTraversal, classNode rhino.Node) {
	superClass := classNode.GetSecondChild()
	classMembers := classNode.GetLastChild()
	var memberDef rhino.Node
	if superClass.IsEmpty() {
		function := s.astFactory.CreateEmptyFunction(AstFactory_Type_(classNode))
		s.compiler.ReportChangeToChangeScope(function)
		memberDef = s.astFactory.CreateMemberFunctionDef("constructor", function)
	} else {
		jutil.CheckStateMsg1(
			superClass.IsQualifiedName(),
			"Expected Es6NormalizeClasses to make all extends clauses into qualified names, found %s",
			superClass)
		body := rhino.IR_Block()

		// If a class is defined in an externs file or as an interface, it's only a stub, not an
		// implementation that should be instantiated.
		// A call to super() shouldn't actually exist for these cases and is problematic to
		// transpile, so don't generate it.
		if !classNode.IsFromExterns() && !s.IsInterface(classNode) {
			// Generate required call to super()
			// `super(...arguments);`
			// Note that transpilation of spread args must occur after this pass for this to work.
			exprResult := rhino.IR_ExprResult(
				s.astFactory.CreateConstructorCall(
					AstFactory_Type_(classNode), // returned type is the subclass
					s.astFactory.CreateSuper(AstFactory_Type_(superClass)),
					rhino.IR_IterSpread(s.astFactory.CreateArgumentsReference())))
			body.AddChildToFront(exprResult)
			NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_SUPER, s.compiler)
			NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_SPREAD_EXPRESSIONS, s.compiler)
		}
		constructor := s.astFactory.CreateFunction("", rhino.IR_ParamList(), body, AstFactory_Type_(classNode))
		memberDef = s.astFactory.CreateMemberFunctionDef("constructor", constructor)
	}
	memberDef.SrcrefTreeIfMissing(classNode)
	memberDef.MakeNonIndexableRecursive()
	classMembers.AddChildToFront(memberDef)
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_MEMBER_DECLARATIONS, s.compiler)
	// report newly created constructor
	s.compiler.ReportChangeToChangeScope(memberDef.GetOnlyChild())
	// report change to scope containing the class
	s.compiler.ReportChangeToEnclosingScope(memberDef)
}

// Java: SynthesizeExplicitConstructors.isInterface(Node)
func (s *SynthesizeExplicitConstructors) IsInterface(classNode rhino.Node) bool {
	classJsDocInfo := NodeUtil_GetBestJSDocInfo(classNode)
	return classJsDocInfo != nil && classJsDocInfo.IsInterface()
}
