// Port of com/google/javascript/jscomp/Es6NormalizeClasses.java @ closure-compiler v20261021

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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Java: Es6NormalizeClasses.ClassRecord.PrivateMemberKind
type es6NormalizeClasses_ClassRecord_PrivateMemberKind int32

const (
	es6NormalizeClasses_ClassRecord_PrivateMemberKind_FIELD es6NormalizeClasses_ClassRecord_PrivateMemberKind = iota
	es6NormalizeClasses_ClassRecord_PrivateMemberKind_METHOD
	es6NormalizeClasses_ClassRecord_PrivateMemberKind_ACCESSOR
)

// Java: Es6NormalizeClasses.ClassRecord.PrivateMemberKind.name()
func (k es6NormalizeClasses_ClassRecord_PrivateMemberKind) String() string {
	switch k {
	case es6NormalizeClasses_ClassRecord_PrivateMemberKind_FIELD:
		return "FIELD"
	case es6NormalizeClasses_ClassRecord_PrivateMemberKind_METHOD:
		return "METHOD"
	case es6NormalizeClasses_ClassRecord_PrivateMemberKind_ACCESSOR:
		return "ACCESSOR"
	}
	panic(jutil.NewIllegalStateException())
}

// es6NormalizeClasses_ClassRecord_PrivateMember is information about a
// private member declaration in the class.
//
// Java: Es6NormalizeClasses.ClassRecord.PrivateMember
type es6NormalizeClasses_ClassRecord_PrivateMember struct {
	name     string
	kind     es6NormalizeClasses_ClassRecord_PrivateMemberKind
	isStatic bool

	methodNode ast.MemberFunctionDef
	getterNode ast.GetterDef
	setterNode ast.SetterDef

	declarationNodes []rhino.Node
}

func es6NormalizeClasses_unmovedFunction(member rhino.Node) (ast.Function, bool) {
	if member.IsNil() || !member.HasChildren() {
		return ast.Function{}, false
	}
	return ast.MustFunction(member.GetFirstChild()), true
}

// Java: Es6NormalizeClasses.ClassRecord.PrivateMember(String, PrivateMemberKind, boolean, Node)
func newEs6NormalizeClasses_ClassRecord_PrivateMember(
	name string, kind es6NormalizeClasses_ClassRecord_PrivateMemberKind, isStatic bool, initialNode ast.ClassMember) *es6NormalizeClasses_ClassRecord_PrivateMember {
	m := &es6NormalizeClasses_ClassRecord_PrivateMember{
		name:             name,
		kind:             kind,
		isStatic:         isStatic,
		declarationNodes: []rhino.Node{initialNode.Node()},
	}
	if kind == es6NormalizeClasses_ClassRecord_PrivateMemberKind_METHOD {
		m.methodNode = ast.MustMemberFunctionDef(initialNode.Node())
	} else {
		m.recordAccessor(initialNode)
	}
	return m
}

func (m *es6NormalizeClasses_ClassRecord_PrivateMember) recordAccessor(member ast.ClassMember) {
	if getter, ok := member.AsGetterDef(); ok {
		m.getterNode = getter
	} else if setter, ok := member.AsSetterDef(); ok {
		m.setterNode = setter
	}
}

// Java: Es6NormalizeClasses.ClassRecord.PrivateMember.propertyName()
func (m *es6NormalizeClasses_ClassRecord_PrivateMember) PropertyName() string {
	return strings.TrimPrefix(m.name, "#")
}

// Java: Es6NormalizeClasses.ClassRecord.PrivateMember.getPrimaryDeclarationNode()
func (m *es6NormalizeClasses_ClassRecord_PrivateMember) GetPrimaryDeclarationNode() rhino.Node {
	if getter := m.getterNode.Node(); !getter.IsNil() {
		return getter
	}
	if setter := m.setterNode.Node(); !setter.IsNil() {
		return setter
	}
	if len(m.declarationNodes) == 0 {
		return rhino.NilNode
	}
	return m.declarationNodes[0]
}

// Java: private static class Es6NormalizeClasses.ClassRecord
type es6NormalizeClasses_ClassRecord struct {
	privateMembers       []*es6NormalizeClasses_ClassRecord_PrivateMember
	privateMembersByName map[string]*es6NormalizeClasses_ClassRecord_PrivateMember

	privateMapVarName       string
	staticPrivateMapVarName string

	hasInstancePrivateMembers            bool
	hasInstancePrivateMethodsOrAccessors bool
	hasStaticPrivateMembers              bool
	needsInstanceSelf                    bool
	needsStaticSelf                      bool

	privateVarName       string
	staticPrivateVarName string
	privateProtoVarName  string

	bindingIdentifier    ast.Name
	hasBindingIdentifier bool
	classNode            ast.Class
	normalizedNames      es6NormalizeClasses_NormalizedNames

	instanceMembers              []ast.ClassMember
	staticMembers                []ast.ClassMember
	computedPropsWithSideEffects []ast.ClassMember
	superCalls                   []rhino.Node
}

// Java: Es6NormalizeClasses.ClassRecord.recordPrivateMember(Node)
func (r *es6NormalizeClasses_ClassRecord) RecordPrivateMember(member ast.ClassMember) {
	n := member.Node()
	jutil.CheckArgument(n.IsPrivateIdentifier())
	name := n.GetString()
	jutil.CheckArgumentMsg1(strings.HasPrefix(name, "#"), "%s", name)
	isStatic := n.IsStaticMember()

	var kind es6NormalizeClasses_ClassRecord_PrivateMemberKind
	switch member.Token() {
	case rhino.Token_MEMBER_FIELD_DEF:
		kind = es6NormalizeClasses_ClassRecord_PrivateMemberKind_FIELD
	case rhino.Token_MEMBER_FUNCTION_DEF:
		kind = es6NormalizeClasses_ClassRecord_PrivateMemberKind_METHOD
	case rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF:
		kind = es6NormalizeClasses_ClassRecord_PrivateMemberKind_ACCESSOR
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Unexpected node for private member: " + jutil.ToJavaString(n)))
	}

	if pm := r.privateMembersByName[name]; pm != nil {
		pm.declarationNodes = append(pm.declarationNodes, n)
		pm.recordAccessor(member)
	} else {
		pm = newEs6NormalizeClasses_ClassRecord_PrivateMember(name, kind, isStatic, member)
		r.privateMembers = append(r.privateMembers, pm)
		if r.privateMembersByName == nil {
			r.privateMembersByName = make(map[string]*es6NormalizeClasses_ClassRecord_PrivateMember, 4)
		}
		r.privateMembersByName[name] = pm
	}

	if (kind == es6NormalizeClasses_ClassRecord_PrivateMemberKind_METHOD || kind == es6NormalizeClasses_ClassRecord_PrivateMemberKind_ACCESSOR) && !isStatic {
		r.hasInstancePrivateMethodsOrAccessors = true
	}

	if isStatic {
		r.hasStaticPrivateMembers = true
	} else {
		r.hasInstancePrivateMembers = true
	}
}

// Java: Es6NormalizeClasses.ClassRecord.hasPrivateMember(String)
func (r *es6NormalizeClasses_ClassRecord) HasPrivateMember(name string) bool {
	_, ok := r.privateMembersByName[name]
	return ok
}

// Java: Es6NormalizeClasses.ClassRecord.getPrivateMember(String)
func (r *es6NormalizeClasses_ClassRecord) GetPrivateMember(name string) *es6NormalizeClasses_ClassRecord_PrivateMember {
	return r.privateMembersByName[name]
}

// Java: Es6NormalizeClasses.ClassRecord.getPrivateMapVarName(AbstractCompiler, NodeTraversal)
func (r *es6NormalizeClasses_ClassRecord) GetPrivateMapVarName(compiler AbstractCompiler, t *NodeTraversal) string {
	if r.privateMapVarName == "" {
		r.privateMapVarName = Es6NormalizeClasses_GenerateUniquePrivateMapVarName(compiler, t)
	}
	return r.privateMapVarName
}

// Java: Es6NormalizeClasses.ClassRecord.getPrivateProtoVarName(AbstractCompiler, NodeTraversal)
func (r *es6NormalizeClasses_ClassRecord) GetPrivateProtoVarName(compiler AbstractCompiler, t *NodeTraversal) string {
	if r.privateProtoVarName == "" {
		r.privateProtoVarName = Es6NormalizeClasses_GenerateUniquePrivateProtoVarName(compiler, t)
	}
	return r.privateProtoVarName
}

// Java: Es6NormalizeClasses.ClassRecord.getStaticPrivateMapVarName(AbstractCompiler, NodeTraversal)
func (r *es6NormalizeClasses_ClassRecord) GetStaticPrivateMapVarName(compiler AbstractCompiler, t *NodeTraversal) string {
	if r.staticPrivateMapVarName == "" {
		r.staticPrivateMapVarName = Es6NormalizeClasses_GenerateUniqueStaticPrivateMapVarName(compiler, t)
	}
	return r.staticPrivateMapVarName
}

// Java: Es6NormalizeClasses.ClassRecord.getPrivateVarName(AbstractCompiler, NodeTraversal)
func (r *es6NormalizeClasses_ClassRecord) GetPrivateVarName(compiler AbstractCompiler, t *NodeTraversal) string {
	if r.privateVarName == "" {
		r.privateVarName = Es6NormalizeClasses_GenerateUniquePrivateVarName(compiler, t)
	}
	return r.privateVarName
}

// Java: Es6NormalizeClasses.ClassRecord.getStaticPrivateVarName(AbstractCompiler, NodeTraversal)
func (r *es6NormalizeClasses_ClassRecord) GetStaticPrivateVarName(compiler AbstractCompiler, t *NodeTraversal) string {
	if r.staticPrivateVarName == "" {
		r.staticPrivateVarName = Es6NormalizeClasses_GenerateUniquePrivateVarName(compiler, t)
	}
	return r.staticPrivateVarName
}

// Java: Es6NormalizeClasses.ClassRecord(NormalizedNames, Optional<Node>, Node)
func newEs6NormalizeClasses_ClassRecord(normalizedNames es6NormalizeClasses_NormalizedNames, bindingIdentifier ast.Name, hasBindingIdentifier bool, classNode ast.Class) *es6NormalizeClasses_ClassRecord {
	return &es6NormalizeClasses_ClassRecord{
		normalizedNames:      normalizedNames,
		bindingIdentifier:    bindingIdentifier,
		hasBindingIdentifier: hasBindingIdentifier,
		classNode:            classNode,
	}
}

// Java: Es6NormalizeClasses.ClassRecord.recordField(Node)
func (r *es6NormalizeClasses_ClassRecord) RecordField(field ast.ClassMember) {
	jutil.CheckArgument(field.Token() == rhino.Token_COMPUTED_FIELD_DEF || field.Token() == rhino.Token_MEMBER_FIELD_DEF)
	if field.Node().IsStaticMember() {
		r.staticMembers = append(r.staticMembers, field)
	} else {
		r.instanceMembers = append(r.instanceMembers, field)
	}
}

// Java: Es6NormalizeClasses.ClassRecord.maybeRecordComputedPropWithSideEffects(Node)
func (r *es6NormalizeClasses_ClassRecord) MaybeRecordComputedPropWithSideEffects(computedProp ast.ClassMember) {
	if NodeUtil_CanBeSideEffected(es6NormalizeClasses_computedKey(computedProp).Node()) {
		r.computedPropsWithSideEffects = append(r.computedPropsWithSideEffects, computedProp)
	}
}

// Java: Es6NormalizeClasses.ClassRecord.recordStaticBlock(Node)
func (r *es6NormalizeClasses_ClassRecord) RecordStaticBlock(block ast.Block) {
	jutil.CheckArgument(NodeUtil_IsClassStaticBlock(block.Node()))
	r.staticMembers = append(r.staticMembers, block.AsClassMember())
}

// Java: Es6NormalizeClasses.ClassRecord.createNewNameReferenceNode()
func (r *es6NormalizeClasses_ClassRecord) CreateNewNameReferenceNode() rhino.Node {
	if r.normalizedNames.ClassNameNode().IsName() {
		return r.normalizedNames.ClassNameNode().CloneNode()
	}
	return r.normalizedNames.ClassNameNode().CloneTree()
}

// Java: Es6NormalizeClasses.ClassRecord.createNewSuperclassReferenceNode()
func (r *es6NormalizeClasses_ClassRecord) CreateNewSuperclassReferenceNode() rhino.Node {
	if r.normalizedNames.SuperClassNode().IsNil() {
		panic(jutil.NewIllegalArgumentExceptionMsg("no superclass"))
	}
	return es6NormalizeClasses_optionalGet(r.normalizedNames.SuperClassNode()).CloneTree()
}

// es6NormalizeClasses_NormalizedNames holds the normalized class and superclass
// reference templates computed before traversing a class body.
type es6NormalizeClasses_NormalizedNames struct {
	classNameNode        rhino.Node
	superClassNode       rhino.Node
	shouldExtractClass   bool
	shouldExtractExtends bool
}

// Java: Es6NormalizeClasses.NormalizedNames(Node, Optional<Node>, boolean, boolean)
func newEs6NormalizeClasses_NormalizedNames(
	classNameNode rhino.Node,
	superClassNode rhino.Node,
	shouldExtractClass bool,
	shouldExtractExtends bool) es6NormalizeClasses_NormalizedNames {
	jutil.CheckArgumentMsg1(classNameNode.IsQualifiedName(), "%s", classNameNode)
	if !(superClassNode.IsNil() || superClassNode.IsQualifiedName()) {
		jutil.CheckArgumentMsg1(false, "Optional[%s]", superClassNode)
	}
	return es6NormalizeClasses_NormalizedNames{
		classNameNode:        classNameNode,
		superClassNode:       superClassNode,
		shouldExtractClass:   shouldExtractClass,
		shouldExtractExtends: shouldExtractExtends,
	}
}

// Java: Es6NormalizeClasses.NormalizedNames.classNameNode()
func (n es6NormalizeClasses_NormalizedNames) ClassNameNode() rhino.Node { return n.classNameNode }

// Java: Es6NormalizeClasses.NormalizedNames.superClassNode()
func (n es6NormalizeClasses_NormalizedNames) SuperClassNode() rhino.Node {
	return n.superClassNode
}

// Java: Es6NormalizeClasses.generateNormalizedClassNames(NodeTraversal, Node)
func (e *Es6NormalizeClasses) GenerateNormalizedClassNames(t *NodeTraversal, classNode ast.Class) es6NormalizeClasses_NormalizedNames {
	extendsClause, hasExtends := classNode.Extends()

	shouldExtractExtends := e.ShouldExtractExtends(classNode)
	var superClass rhino.Node
	if shouldExtractExtends {
		name := e.GenerateUniqueClassExtendsVarName(t)
		superClass = e.astFactory.CreateName(name, AstFactory_Type_(extendsClause.Node()))
	} else if hasExtends {
		superClass = extendsClause.Node()
	}

	shouldExtractClass :=
		e.ShouldExtractClass(classNode) ||
			(shouldExtractExtends && !e.CanExtractExtendsSimply(classNode))
	var classNameNodeTemplate rhino.Node
	if shouldExtractClass {
		name := e.GenerateUniqueClassDeclVarName(t)
		classNameNodeTemplate = e.astFactory.CreateConstantName(name, AstFactory_Type_(classNode.Node()))
	} else {
		classNameNodeTemplate = NodeUtil_GetNameNode(classNode.Node())
	}

	return newEs6NormalizeClasses_NormalizedNames(
		classNameNodeTemplate, superClass, shouldExtractClass, shouldExtractExtends)
}
