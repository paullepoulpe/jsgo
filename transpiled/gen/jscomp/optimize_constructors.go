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
// Port of com/google/javascript/jscomp/OptimizeConstructors.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.OptimizeConstructors
type OptimizeConstructors struct {
	compiler              AbstractCompiler `gson:"compiler"`
	astAnalyzer           *AstAnalyzer     `gson:"astAnalyzer"`
	removableConstructors []pkg_rhino.Node `gson:"removableConstructors"`
	idHash_               atomic.Int32
}

func (this_ *OptimizeConstructors) BindSelf_(s any) {
	_ = s
}
func (this_ *OptimizeConstructors) IsJavaNil_() bool { return this_ == nil }
func (*OptimizeConstructors) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.OptimizeConstructors", "OptimizeConstructors"
}
func (this_ *OptimizeConstructors) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*OptimizeConstructors) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.OptimizeConstructors", "com.google.javascript.jscomp.OptimizeConstructors", func(v_ any) bool {
		return j2gort.InstanceOf[*OptimizeConstructors](v_)
	})
}
func (*OptimizeConstructors) Implements_jscomp_CompilerPass()                        {}
func (*OptimizeConstructors) Implements_jscomp_OptimizeCalls_CallGraphCompilerPass() {}

func (this_ *OptimizeConstructors) initFields_OptimizeConstructors() {
	_ = this_
	this_.removableConstructors = []pkg_rhino.Node{}
}

// Java constructor: OptimizeConstructors(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewOptimizeConstructors(compiler AbstractCompiler) *OptimizeConstructors {
	this_ := &OptimizeConstructors{}
	InitOptimizeConstructors(this_, compiler)
	return this_
}

func InitOptimizeConstructors(this_ *OptimizeConstructors, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_OptimizeConstructors()
	this_.compiler = j2gort.CheckNotNullIface(compiler)
	this_.astAnalyzer = compiler.GetAstAnalyzer()
}

// Java method: OptimizeConstructors.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *OptimizeConstructors) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	if !(this_.compiler.GetLifeCycleStage() == StaticGet_AbstractCompiler_LifeCycleStage_NORMALIZED()) {
		j2gort.CheckState(false)
	}
	OptimizeCalls_S_Builder().SetCompiler(this_.compiler).SetConsiderExterns(false).AddPass(this_).Build().Process_2(externs, root)
}

// Java method: OptimizeConstructors.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root, com.google.javascript.jscomp.OptimizeCalls.ReferenceMap refMap)
func (this_ *OptimizeConstructors) Process_3(externs pkg_rhino.Node, root pkg_rhino.Node, refMap *OptimizeCalls_ReferenceMap) {
	_ = externs
	_ = root
	_ = refMap
	{
		cur_1 := refMap.names.IterEntries()
		for cur_1.NextEntry() {
			this_.priv_OptimizeConstructors_addConstructorsToBeRemoved_4082(cur_1.Key(), cur_1.Value())
		}
	}
	for _, ref := range this_.removableConstructors {
		_ = ref
		this_.priv_OptimizeConstructors_removeConstructorMethod_4087(ref)
	}

}

// Java method: OptimizeConstructors.addConstructorsToBeRemoved(java.lang.String name, com.google.javascript.jscomp.OptimizeCalls.NodeList refs)
func (this_ *OptimizeConstructors) priv_OptimizeConstructors_addConstructorsToBeRemoved_4082(name string, refs *OptimizeCalls_NodeList) {
	_ = name
	_ = refs
	if !(OptimizeCalls_MayBeOptimizableName(this_.compiler, name)) {
		return
	}
	var candidateClassDefinition pkg_rhino.Node = pkg_rhino.Node(0)
	_ = candidateClassDefinition
	var candidateConstructor pkg_rhino.Node = pkg_rhino.Node(0)
	_ = candidateConstructor
	var subclassConstructors []pkg_rhino.Node = nil
	_ = subclassConstructors
	{
		var i int32 = int32(0)
		_ = i
		var size int32 = refs.Size()
		_ = size
		for ; i < size; i = i + 1 {
			var n pkg_rhino.Node = refs.Get(i)
			_ = n
			var definition pkg_rhino.Node = this_.priv_OptimizeConstructors_getClassDefinitionOrFunction_4084(n)
			_ = definition
			if definition != 0 {
				if candidateClassDefinition != 0 {
					return
				}
				candidateClassDefinition = definition
				if candidateClassDefinition.IsClass() {
					var extendsExpr pkg_rhino.Node = candidateClassDefinition.GetSecondChild()
					_ = extendsExpr
					if extendsExpr.IsEmpty() {
						var constructor pkg_rhino.Node = NodeUtil_GetEs6ClassConstructorMemberFunctionDef(candidateClassDefinition)
						_ = constructor
						if constructor != 0 {
							var fn pkg_rhino.Node = constructor.GetLastChild()
							_ = fn
							var body pkg_rhino.Node = fn.GetLastChild()
							_ = body
							if !(body.HasChildren()) && this_.priv_OptimizeConstructors_hasRemovableParameterList_4086(constructor) {
								candidateConstructor = constructor
							}
						}
					} else {
					}
				}
			} else {
				if this_.priv_OptimizeConstructors_isClassExtendsExpression_4083(n) {
					var subclassLiteral pkg_rhino.Node = n.GetParent()
					_ = subclassLiteral
					var subclassConstructor pkg_rhino.Node = NodeUtil_GetEs6ClassConstructorMemberFunctionDef(subclassLiteral)
					_ = subclassConstructor
					if (subclassConstructor != 0) && OptimizeConstructors_ConstructorHasRemovableDefinition(subclassConstructor) {
						if subclassConstructors == nil {
							subclassConstructors = []pkg_rhino.Node{}
						}
						subclassConstructors = append(subclassConstructors, subclassConstructor)
					}
				} else {
					if OptimizeConstructors_IsAssigningReference(n) {
						return
					}
				}
			}
		}
	}
	if candidateClassDefinition == 0 {
		return
	}
	if candidateConstructor != 0 {
		this_.removableConstructors = append(this_.removableConstructors, candidateConstructor)
	}
	if subclassConstructors != nil {
		var summary *OptimizeConstructors_ClassConstructorSummary = OptimizeConstructors_ClassConstructorSummary_Build(candidateClassDefinition)
		_ = summary
		for _, n_2 := range subclassConstructors {
			_ = n_2
			if summary.IsEquivalentConstructorDefinition(n_2) {
				this_.removableConstructors = append(this_.removableConstructors, n_2)
			}
		}

	}
}

// Java static method: OptimizeConstructors.isAssigningReference(com.google.javascript.rhino.Node n)
func OptimizeConstructors_IsAssigningReference(n pkg_rhino.Node) bool {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	var gparent pkg_rhino.Node = parent.GetParent()
	_ = gparent
	switch parent.TokenSwitchOrd_() {
	case int32(111) /* Token.LET */, int32(97) /* Token.CONST */, int32(84) /* Token.VAR */ :
		return n.HasChildren()
	case int32(100) /* Token.STRING_KEY */ :
		return gparent.IsObjectPattern()
	case int32(129) /* Token.COMPUTED_PROP */ :
		return ((parent.GetLastChild() == n) && gparent.IsObjectPattern())
	case int32(102) /* Token.ARRAY_PATTERN */, int32(134) /* Token.DEFAULT_VALUE */, int32(49) /* Token.PARAM_LIST */, int32(126) /* Token.OBJECT_REST */, int32(125) /* Token.ITER_REST */, int32(71) /* Token.INC */, int32(72) /* Token.DEC */ :
		return true
	case int32(73) /* Token.FUNCTION */, int32(105) /* Token.CLASS */, int32(86) /* Token.CATCH */ :
		return (parent.GetFirstChild() == n)
	case int32(51) /* Token.ASSIGN */, int32(52) /* Token.ASSIGN_BITOR */, int32(53) /* Token.ASSIGN_BITXOR */, int32(54) /* Token.ASSIGN_BITAND */, int32(55) /* Token.ASSIGN_LSH */, int32(56) /* Token.ASSIGN_RSH */, int32(57) /* Token.ASSIGN_URSH */, int32(58) /* Token.ASSIGN_ADD */, int32(59) /* Token.ASSIGN_SUB */, int32(60) /* Token.ASSIGN_MUL */, int32(61) /* Token.ASSIGN_DIV */, int32(62) /* Token.ASSIGN_MOD */, int32(63) /* Token.ASSIGN_EXPONENT */ :
		return (parent.GetFirstChild() == n)
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: OptimizeConstructors.isClassExtendsExpression(com.google.javascript.rhino.Node n)
func (this_ *OptimizeConstructors) priv_OptimizeConstructors_isClassExtendsExpression_4083(n pkg_rhino.Node) bool {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	return (parent.IsClass() && (parent.GetSecondChild() == n))
}

// Java method: OptimizeConstructors.getClassDefinitionOrFunction(com.google.javascript.rhino.Node n)
func (this_ *OptimizeConstructors) priv_OptimizeConstructors_getClassDefinitionOrFunction_4084(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	var expr pkg_rhino.Node = pkg_rhino.Node(0)
	_ = expr
	if OptimizeCalls_ReferenceMap_IsSimpleAssignmentTarget(n) {
		expr = parent.GetLastChild()
	} else if n.IsName() && n.HasChildren() {
		expr = n.GetFirstChild()
	} else if parent.IsFunction() && n.IsFirstChildOf(parent) {
		expr = parent
	} else if parent.IsClass() && n.IsFirstChildOf(parent) {
		expr = parent
	} else {
		return pkg_rhino.Node(0)
	}
	expr = this_.priv_OptimizeConstructors_unwrap_4085(expr)
	if OptimizeConstructors_IsDefinitionClassLiteralOrFunction(expr) {
		return expr
	} else {
		return pkg_rhino.Node(0)
	}
}

// Java method: OptimizeConstructors.unwrap(com.google.javascript.rhino.Node expr)
func (this_ *OptimizeConstructors) priv_OptimizeConstructors_unwrap_4085(expr pkg_rhino.Node) pkg_rhino.Node {
	_ = expr
	for expr.IsCast() || expr.IsComma() {
		expr = expr.GetLastChild()
	}
	return expr
}

// Java static method: OptimizeConstructors.isDefinitionClassLiteralOrFunction(com.google.javascript.rhino.Node n)
func OptimizeConstructors_IsDefinitionClassLiteralOrFunction(n pkg_rhino.Node) bool {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(73) /* Token.FUNCTION */ :
		return true
	case int32(105) /* Token.CLASS */ :
		var constructorMemberFunctionDef pkg_rhino.Node = NodeUtil_GetEs6ClassConstructorMemberFunctionDef(n)
		_ = constructorMemberFunctionDef
		return (constructorMemberFunctionDef != 0)
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: OptimizeConstructors.functionHasRest(com.google.javascript.rhino.Node fn)
func OptimizeConstructors_FunctionHasRest(fn pkg_rhino.Node) bool {
	_ = fn
	if !(fn.IsFunction()) {
		j2gort.CheckState(false)
	}
	var params pkg_rhino.Node = NodeUtil_GetFunctionParameters(fn)
	_ = params
	var lastParam pkg_rhino.Node = params.GetLastChild()
	_ = lastParam
	return ((lastParam != 0) && lastParam.IsRest())
}

// Java static method: OptimizeConstructors.constructorHasRemovableDefinition(com.google.javascript.rhino.Node member)
func OptimizeConstructors_ConstructorHasRemovableDefinition(member pkg_rhino.Node) bool {
	_ = member
	var fn pkg_rhino.Node = member.GetFirstChild()
	_ = fn
	var superCall pkg_rhino.Node = OptimizeConstructors_GetOnlySuperCall(fn)
	_ = superCall
	if superCall == 0 {
		return false
	}
	var paramList pkg_rhino.Node = NodeUtil_GetFunctionParameters(fn)
	_ = paramList
	if paramList.GetChildCount() != (superCall.GetChildCount() - int32(1)) {
		return false
	}
	var param pkg_rhino.Node = paramList.GetFirstChild()
	_ = param
	var arg pkg_rhino.Node = superCall.GetSecondChild()
	_ = arg
	for param != 0 {
		if param.IsRest() {
			if !(arg.IsSpread()) || !(param.GetFirstChild().MatchesName_1_Node(arg.GetFirstChild())) {
				return false
			}
		} else if !(param.MatchesName_1_Node(arg)) {
			return false
		}
		param = param.GetNext()
		arg = arg.GetNext()
	}
	return true
}

// Java method: OptimizeConstructors.hasRemovableParameterList(com.google.javascript.rhino.Node member)
func (this_ *OptimizeConstructors) priv_OptimizeConstructors_hasRemovableParameterList_4086(member pkg_rhino.Node) bool {
	_ = member
	var fn pkg_rhino.Node = member.GetFirstChild()
	_ = fn
	var paramList pkg_rhino.Node = NodeUtil_GetFunctionParameters(fn)
	_ = paramList
	{
		var param pkg_rhino.Node = paramList.GetFirstChild()
		_ = param
		for ; param != 0; param = param.GetNext() {
			if param.IsName() {
				continue
			}
			if param.IsRest() && param.GetFirstChild().IsName() {
				continue
			}
			if (param.IsDefaultValue() && param.GetFirstChild().IsName()) && !(this_.astAnalyzer.MayHaveSideEffects(param.GetLastChild())) {
				continue
			}
			return false
		}
	}
	return true
}

// Java static method: OptimizeConstructors.getOnlySuperCall(com.google.javascript.rhino.Node fn)
func OptimizeConstructors_GetOnlySuperCall(fn pkg_rhino.Node) pkg_rhino.Node {
	_ = fn
	var body pkg_rhino.Node = fn.GetLastChild()
	_ = body
	if body.IsBlock() && body.HasOneChild() {
		var stmt pkg_rhino.Node = body.GetFirstChild()
		_ = stmt
		if stmt.IsExprResult() {
			var call pkg_rhino.Node = stmt.GetFirstChild()
			_ = call
			if call.IsCall() && call.GetFirstChild().IsSuper() {
				return call
			}
		}
	}
	return pkg_rhino.Node(0)
}

// Java method: OptimizeConstructors.removeConstructorMethod(com.google.javascript.rhino.Node member)
func (this_ *OptimizeConstructors) priv_OptimizeConstructors_removeConstructorMethod_4087(member pkg_rhino.Node) {
	_ = member
	if !(member.IsMemberFunctionDef()) {
		j2gort.CheckState(false)
	}
	this_.compiler.ReportFunctionDeleted(member.GetFirstChild())
	this_.compiler.ReportChangeToEnclosingScope(member)
	member.Detach()
}

// Java class: com.google.javascript.jscomp.OptimizeConstructors.ClassConstructorSummary
type OptimizeConstructors_ClassConstructorSummary struct {
	formalParameterCount int32 `gson:"formalParameterCount"`
	idHash_              atomic.Int32
	isVarArgs            bool `gson:"isVarArgs"`
}

func (this_ *OptimizeConstructors_ClassConstructorSummary) BindSelf_(s any) {
	_ = s
}
func (this_ *OptimizeConstructors_ClassConstructorSummary) IsJavaNil_() bool { return this_ == nil }
func (*OptimizeConstructors_ClassConstructorSummary) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.OptimizeConstructors$ClassConstructorSummary", "ClassConstructorSummary"
}
func (this_ *OptimizeConstructors_ClassConstructorSummary) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*OptimizeConstructors_ClassConstructorSummary) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.OptimizeConstructors.ClassConstructorSummary", "com.google.javascript.jscomp.OptimizeConstructors$ClassConstructorSummary", func(v_ any) bool {
		return j2gort.InstanceOf[*OptimizeConstructors_ClassConstructorSummary](v_)
	})
}

func (this_ *OptimizeConstructors_ClassConstructorSummary) initFields_OptimizeConstructors_ClassConstructorSummary() {
	_ = this_
}

// Java constructor: ClassConstructorSummary(boolean isVarArgs, int formalParameterCount)
func NewOptimizeConstructors_ClassConstructorSummary(isVarArgs bool, formalParameterCount int32) *OptimizeConstructors_ClassConstructorSummary {
	this_ := &OptimizeConstructors_ClassConstructorSummary{}
	InitOptimizeConstructors_ClassConstructorSummary(this_, isVarArgs, formalParameterCount)
	return this_
}

func InitOptimizeConstructors_ClassConstructorSummary(this_ *OptimizeConstructors_ClassConstructorSummary, isVarArgs bool, formalParameterCount int32) {
	_ = isVarArgs
	_ = formalParameterCount
	this_.initFields_OptimizeConstructors_ClassConstructorSummary()
	this_.isVarArgs = isVarArgs
	this_.formalParameterCount = formalParameterCount
}

// Java static method: ClassConstructorSummary.build(com.google.javascript.rhino.Node classDefinition)
func OptimizeConstructors_ClassConstructorSummary_Build(classDefinition pkg_rhino.Node) *OptimizeConstructors_ClassConstructorSummary {
	_ = classDefinition
	if !(classDefinition.IsClass() || classDefinition.IsFunction()) {
		j2gort.CheckState(false)
	}
	var fn pkg_rhino.Node = pkg_rhino.Node(0)
	_ = fn
	if classDefinition.IsClass() {
		var member pkg_rhino.Node = NodeUtil_GetEs6ClassConstructorMemberFunctionDef(classDefinition)
		_ = member
		fn = member.GetFirstChild()
	} else {
		fn = classDefinition
	}
	var argumentsReference bool = NodeUtil_DoesFunctionReferenceOwnArgumentsObject(fn)
	_ = argumentsReference
	var hasVarArgs bool = (argumentsReference || OptimizeConstructors_FunctionHasRest(fn))
	_ = hasVarArgs
	return NewOptimizeConstructors_ClassConstructorSummary(hasVarArgs, NodeUtil_GetFunctionParameters(fn).GetChildCount())
}

// Java method: ClassConstructorSummary.isEquivalentConstructorDefinition(com.google.javascript.rhino.Node constructorMember)
func (this_ *OptimizeConstructors_ClassConstructorSummary) IsEquivalentConstructorDefinition(constructorMember pkg_rhino.Node) bool {
	_ = constructorMember
	var fn pkg_rhino.Node = constructorMember.GetFirstChild()
	_ = fn
	var paramList pkg_rhino.Node = NodeUtil_GetFunctionParameters(fn)
	_ = paramList
	var hasRest bool = OptimizeConstructors_FunctionHasRest(fn)
	_ = hasRest
	if hasRest {
		return true
	}
	if this_.isVarArgs {
		return false
	}
	if this_.formalParameterCount != paramList.GetChildCount() {
		return false
	}
	return true
}
