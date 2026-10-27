// Port of com/google/javascript/jscomp/Es6RewriteClass.java @ closure-compiler v20261021

/*
 * Copyright 2014 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: this pass uses the typed AST view (internal/ast) where it inspects
// class structure, members, constructor returns, and enclosing declarations,
// and replaces Java's AutoBuilder metadata/property records with compact value
// structs.

// Es6RewriteClass converts ES6 classes to valid ES5 or ES3 code.
type Es6RewriteClass struct {
	compiler                     AbstractCompiler
	astFactory                   *AstFactory
	convertSuperConstructorCalls *Es6ConvertSuperConstructorCalls
	transpilationNamespace       rhino.StaticScope
	jscompInherits               RuntimeJsLibManager_JsLibField
	isES3Output                  bool
	metadataBuf                  es6RewriteClassMetadata
}

var es6RewriteClass_noCollapseJSDoc = func() *rhino.JSDocInfo {
	b := rhino.JSDocInfo_Builder_()
	b.RecordNoCollapse()
	return b.Build()
}()

// Java: Es6RewriteClass.features
var es6RewriteClass_features = parser.FeatureSet_BARE_MINIMUM.WithFeatures(
	parser.FeatureSet_Feature_CLASSES, parser.FeatureSet_Feature_CLASS_GETTER_SETTER, parser.FeatureSet_Feature_NEW_TARGET, parser.FeatureSet_Feature_SUPER)

// Java: Es6RewriteClass.CONSTRUCTOR_PRIMITIVE_RETURN
var Es6RewriteClass_CONSTRUCTOR_PRIMITIVE_RETURN = DiagnosticType_Error(
	"JSC_CONSTRUCTOR_PRIMITIVE_RETURN",
	"Class constructors may only return an object or undefined")

// Java: Es6RewriteClass(AbstractCompiler, Es6SubclassTranspilation)
func NewEs6RewriteClass(
	compiler AbstractCompiler, es6SubclassTranspilation CompilerOptions_Es6SubclassTranspilation) *Es6RewriteClass {
	e := &Es6RewriteClass{}
	e.compiler = compiler
	e.astFactory = compiler.CreateAstFactory()
	e.transpilationNamespace = compiler.GetTranspilationNamespace()
	e.convertSuperConstructorCalls =
		NewEs6ConvertSuperConstructorCalls(compiler, es6SubclassTranspilation)
	e.jscompInherits = compiler.GetRuntimeJsLibManager().GetJsLibField("$jscomp.inherits")
	return e
}

// Java: Es6RewriteClass.process(Node, Node)
func (e *Es6RewriteClass) Process(externs rhino.Node, root rhino.Node) {
	if opts := e.compiler.GetOptions(); opts != nil && opts.GetOutputFeatureSet() != nil {
		e.isES3Output = parser.FeatureSet_ES3.Contains(opts.GetOutputFeatureSet())
	} else {
		e.isES3Output = false
	}
	// TODO(b/171853310): This transpilation should be turned off in externs
	TranspilationPasses_ProcessTranspile(e.compiler, externs, es6RewriteClass_features, e)
	TranspilationPasses_ProcessTranspile(e.compiler, root, es6RewriteClass_features, e)
	// Super constructor calls are done all at once as a separate step largely for historical
	// reasons. It used to be an entirely separate pass, but that has been fixed so we no longer
	// have an invalid AST state between passes.
	// TODO(bradfordcsmith): It would probably be more readable and efficient to merge the super
	//     constructor rewriting logic into this class.
	// The code here only creates the GlobalNamespace object which is very cheap. The expensive
	// building of global namespace happens inside es6ConvertSuperConstructorCalls pass.
	ns := NewGlobalNamespaceWithRoots(e.compiler, externs, root)
	defer e.convertSuperConstructorCalls.SetGlobalNamespace(nil)
	e.convertSuperConstructorCalls.SetGlobalNamespace(ns)
	NodeTraversal_Traverse(e.compiler, root, e.convertSuperConstructorCalls)
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(e.compiler, root, es6RewriteClass_features)
}

// Java: Es6RewriteClass.shouldTraverse(NodeTraversal, Node, Node)
func (e *Es6RewriteClass) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF:
		if e.isES3Output {
			TranspilationUtil_CannotConvert(e.compiler, n, "ES5 getters/setters (consider using --language_out=ES5)")
			return false
		}
	default:
	}
	return true
}

// Java: Es6RewriteClass.visit(NodeTraversal, Node, Node)
func (e *Es6RewriteClass) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if cls, ok := ast.AsClass(n); ok {
		e.VisitClass(t, cls, parent)
	}
}

// VisitClass processes classes in 3 phases:
//
//  1. The class name is extracted.
//  2. Class members are processed and rewritten.
//  3. The constructor is built.
//
// Java: Es6RewriteClass.visitClass(NodeTraversal, Node, Node)
func (e *Es6RewriteClass) VisitClass(t *NodeTraversal, cls ast.Class, parent rhino.Node) {
	classNode := cls.Node()
	// Collect Metadata
	metadata := es6RewriteClassCreateMetadata(&e.metadataBuf, cls, parent, e.astFactory)

	if metadata == nil {
		panic(jutil.NewIllegalStateExceptionMsg(
			"Can only convert classes that are declarations or the right hand" +
				" side of a simple assignment: " +
				jutil.ToJavaString(classNode)))
	}
	if metadata.hasSuperClass() {
		jutil.CheckStateMsg1(
			metadata.superClassNameNode.IsQualifiedName(),
			"Expected Es6NormalizeClasses to make all extends clauses into qualified  names, found"+
				" %s",
			metadata.superClassNameNode)
	}

	jutil.CheckStateMsg1(
		NodeUtil_IsStatement(metadata.insertionPoint.node),
		"insertion point must be a statement: %s",
		metadata.insertionPoint.node)

	var constructor rhino.Node
	// Process all members of the class
	for typedMember := range ast.ReadAhead(cls.Body().FirstMember()) {
		member := typedMember.Node()
		if prop, ok := typedMember.AsComputedProp(); ok && (prop.IsGetter() || prop.IsSetter()) {
			e.extractComputedProperty(prop, metadata)
		} else if typedMember.Token() == rhino.Token_GETTER_DEF || typedMember.Token() == rhino.Token_SETTER_DEF {
			e.visitNonMethodMember(member, metadata)
		} else if mfd, ok := typedMember.AsMemberFunctionDef(); ok && !mfd.IsStatic() && mfd.Key() == "constructor" {
			ctorFn := mfd.Function()
			if !classNode.IsFromExterns() {
				e.ReportPrimitiveReturns(ctorFn)
			}
			constructor = ctorFn.Detach().Node().SetColor(classNode.GetColor())
			if !metadata.anonymous {
				// Turns class Foo { constructor: function() {} } into function Foo() {},
				// i.e. attaches the name to the ctor function.
				ast.MustFunction(constructor).NameNode().Node().ReplaceWith(metadata.classNameNode.CloneNode())
			}
		} else if typedMember.Token() == rhino.Token_EMPTY {
			// Do nothing.
		} else {
			jutil.CheckStateMsg1(
				member.IsMemberFunctionDef() || member.IsComputedProp(),
				"Unexpected class member: (%s)",
				member)
			jutil.CheckStateMsg1(
				!member.GetBooleanProp(rhino.Node_COMPUTED_PROP_VARIABLE),
				"Member variables should have been transpiled earlier: (%s)",
				member)
			e.visitMethod(member, metadata)
		}
	}
	rhino.CheckNotNullMsg(
		constructor,
		"Es6RewriteClasses expects all classes to have (possibly synthetic) constructors")

	e.flushDefinePropertiesForPrototype(metadata, classNode)
	e.flushDefinePropertiesForClass(metadata, classNode)

	classJSDoc := NodeUtil_GetBestJSDocInfo(classNode)
	newInfo := rhino.JSDocInfo_Builder_MaybeCopyFrom(classJSDoc)
	newInfo.RecordConstructor()

	enclosingStatement := NodeUtil_GetEnclosingStatement(classNode)
	if metadata.hasSuperClass() && !classNode.IsFromExterns() {
		inheritsCall :=
			rhino.IR_ExprResult(
				e.astFactory.CreateCall(
					e.astFactory.CreateQNameOfField(e.transpilationNamespace, e.jscompInherits),
					AstFactory_Type_OfColor(colors.StandardColors_NULL_OR_VOID),
					metadata.fullClassNameNode.CloneTree(),
					metadata.superClassNameNode.CloneTree())).
				SrcrefTreeIfMissing(metadata.superClassNameNode)
		inheritsCall.InsertAfter(enclosingStatement)
	}

	e.addTypeDeclarations(metadata, enclosingStatement)

	if NodeUtil_IsStatement(classNode) {
		ast.MustFunction(constructor).RemoveName()
		ctorVar := rhino.IR_LetWithValue(metadata.classNameNode.CloneNode(), constructor)
		ctorVar.SrcrefTreeIfMissing(classNode)
		classNode.ReplaceWith(ctorVar)
		NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_LET_DECLARATIONS, e.compiler)
	} else {
		classNode.ReplaceWith(constructor)
	}
	NodeUtil_MarkFunctionsDeleted(classNode, e.compiler)

	if NodeUtil_IsStatement(constructor) {
		constructor.SetJSDocInfo(newInfo.Build())
	} else if parent.IsName() {
		// The constructor function is the RHS of a var statement.
		// Add the JSDoc to the VAR node.
		parent.GetParent().SetJSDocInfo(newInfo.Build())
	} else if constructor.GetParent().IsName() {
		// Is a newly created VAR node.
		constructor.GetGrandparent().SetJSDocInfo(newInfo.Build())
	} else if parent.IsAssign() {
		// The constructor function is the RHS of an assignment.
		// Add the JSDoc to the ASSIGN node.
		parent.SetJSDocInfo(newInfo.Build())
	} else {
		panic(jutil.NewIllegalStateExceptionMsg("Unexpected parent node " + jutil.ToJavaString(parent)))
	}

	t.ReportCodeChange()
}

// ReportPrimitiveReturns reports an error for each primitive value returned
// from a class constructor. A native class constructor ignores such a value
// in a base class and throws a TypeError in a derived class (ES2022 §
// 10.2.2). After transpilation the constructor is a plain function, so a
// subclass's `Base.call(this) || this` can evaluate to the value, and a
// derived class constructor returns the new object instead of throwing.
//
// Only side-effect-free primitive values (NodeUtil.isImmutableValue) are
// recognized, except `undefined`, which both kinds of constructor may
// return. A computed return value would need a runtime check in every
// constructor.
//
// Java: Es6RewriteClass.reportPrimitiveReturns(Node)
func (e *Es6RewriteClass) ReportPrimitiveReturns(constructor ast.Function) {
	// TODO(b/573571176): A derived constructor that returns `undefined`, explicitly or by reaching
	// the end of its body, before calling `super()` must throw a ReferenceError.
	NodeUtil_VisitPostOrderWithPred(
		constructor.Body().Node(),
		NodeUtil_VisitorFunc(func(n rhino.Node) {
			if ret, ok := ast.AsReturn(n); ok {
				if val, hasVal := ret.Value(); hasVal &&
					NodeUtil_IsImmutableValue(val.Node()) &&
					NodeUtil_GetKnownValueType(val.Node()) != NodeUtil_ValueType_VOID {
					e.compiler.Report(JSError_MakeWithNode(n, Es6RewriteClass_CONSTRUCTOR_PRIMITIVE_RETURN))
				}
			}
		}),
		// A `return` in a nested function belongs to that function.
		func(n rhino.Node) bool { return !n.IsFunction() })
}

// Java: Es6RewriteClass.createObjectDotDefineProperties()
func (e *Es6RewriteClass) createObjectDotDefineProperties() rhino.Node {
	return e.astFactory.CreateJSCompDotGlobalAccess(
		e.transpilationNamespace, "Object.defineProperties")
}

// Java: Es6RewriteClass.createObjectDotDefineProperty()
func (e *Es6RewriteClass) createObjectDotDefineProperty() rhino.Node {
	return e.astFactory.CreateJSCompDotGlobalAccess(
		e.transpilationNamespace, "Object.defineProperty")
}

// Java: Es6RewriteClass.flushDefinePropertiesForPrototype(ClassDeclarationMetadata, Node)
func (e *Es6RewriteClass) flushDefinePropertiesForPrototype(
	metadata *es6RewriteClassMetadata, srcrefNode rhino.Node) {
	if !metadata.definePropertiesObjForPrototype.IsNil() && metadata.definePropertiesObjForPrototype.HasChildren() {
		definePropsCall :=
			rhino.IR_ExprResult(
				e.astFactory.CreateCall(
					e.createObjectDotDefineProperties(),
					AstFactory_Type_(metadata.classPrototypeNode),
					metadata.classPrototypeNode.CloneTree(),
					metadata.definePropertiesObjForPrototype))
		definePropsCall.SrcrefTreeIfMissing(srcrefNode)
		metadata.insertNodeAndAdvance(definePropsCall)
		metadata.definePropertiesObjForPrototype = rhino.NilNode
	}
}

// Java: Es6RewriteClass.flushDefinePropertiesForClass(ClassDeclarationMetadata, Node)
func (e *Es6RewriteClass) flushDefinePropertiesForClass(metadata *es6RewriteClassMetadata, srcrefNode rhino.Node) {
	if !metadata.definePropertiesObjForClass.IsNil() && metadata.definePropertiesObjForClass.HasChildren() {
		definePropsCall :=
			rhino.IR_ExprResult(
				e.astFactory.CreateCall(
					e.createObjectDotDefineProperties(),
					AstFactory_Type_(metadata.fullClassNameNode),
					metadata.fullClassNameNode.CloneTree(),
					metadata.definePropertiesObjForClass))
		definePropsCall.SrcrefTreeIfMissing(srcrefNode)
		metadata.insertNodeAndAdvance(definePropsCall)
		metadata.definePropertiesObjForClass = rhino.NilNode
	}
}

// Java: Es6RewriteClass.flushPendingDefineProperties(ClassDeclarationMetadata, boolean, Node)
func (e *Es6RewriteClass) flushPendingDefineProperties(
	metadata *es6RewriteClassMetadata, isStatic bool, srcrefNode rhino.Node) {
	if isStatic {
		e.flushDefinePropertiesForClass(metadata, srcrefNode)
	} else {
		e.flushDefinePropertiesForPrototype(metadata, srcrefNode)
	}
}

// addToDefinePropertiesObject takes member, a non-computed getter or setter.
//
// Java: Es6RewriteClass.addToDefinePropertiesObject(ClassDeclarationMetadata, Node)
func (e *Es6RewriteClass) addToDefinePropertiesObject(metadata *es6RewriteClassMetadata, member rhino.Node) {
	jutil.CheckArgument(!member.IsComputedProp())
	var obj rhino.Node
	if member.IsStaticMember() {
		obj = metadata.definePropertiesForClass()
	} else {
		obj = metadata.definePropertiesForPrototype()
	}
	prop := NodeUtil_GetFirstPropMatchingKey(obj, member.GetString())
	accessorKind := "set"
	var function ast.Function
	if getter, ok := ast.AsGetterDef(member); ok {
		accessorKind = "get"
		function = getter.Function()
	} else {
		function = ast.MustSetterDef(member).Function()
	}
	if !prop.IsNil() && !NodeUtil_GetFirstPropMatchingKey(prop, accessorKind).IsNil() {
		e.flushPendingDefineProperties(metadata, member.IsStaticMember(), member)
		if member.IsStaticMember() {
			obj = metadata.definePropertiesForClass()
		} else {
			obj = metadata.definePropertiesForPrototype()
		}
		prop = rhino.NilNode
	}
	if prop.IsNil() {
		prop = e.createPropertyDescriptor()
		if member.GetString() == "__proto__" {
			stringKey :=
				e.astFactory.CreateComputedProperty(e.astFactory.CreateString("__proto__"), prop)
			obj.AddChildToBack(stringKey)
			scriptNode := NodeUtil_GetEnclosingScript(member)
			if !scriptNode.IsNil() {
				NodeUtil_AddFeatureToScript(scriptNode, parser.FeatureSet_Feature_COMPUTED_PROPERTIES, e.compiler)
			}
		} else {
			stringKey := e.astFactory.CreateStringKey(member.GetString(), prop)
			if member.IsQuotedStringKey() {
				stringKey.PutBooleanProp(rhino.Node_QUOTED_PROP, true)
			}
			obj.AddChildToBack(stringKey)
		}
	}

	info := NodeUtil_GetBestJSDocInfo(function.Node())
	stringKey := e.astFactory.CreateStringKey(accessorKind, function.Detach().Node())
	stringKey.SetJSDocInfo(info)
	prop.AddChildToBack(stringKey)
	prop.SrcrefTreeIfMissing(member)
}

// extractComputedProperty appends an Object.defineProperty call defining
// the given computed getter or setter.
//
// Java: Es6RewriteClass.extractComputedProperty(Node, ClassDeclarationMetadata)
func (e *Es6RewriteClass) extractComputedProperty(computedMember ast.ComputedProp, metadata *es6RewriteClassMetadata) {
	memberNode := computedMember.Detach().Node()
	e.flushPendingDefineProperties(metadata, memberNode.IsStaticMember(), memberNode)

	var owner rhino.Node
	if memberNode.IsStaticMember() {
		owner = metadata.fullClassNameNode
	} else {
		owner = metadata.classPrototypeNode
	}
	keyNode := computedMember.Key()
	valNode, ok := computedMember.Value()
	jutil.CheckState(ok)
	property := keyNode.Detach().Node()
	propertyValue := valNode.Detach().Node()

	propertyDescriptor := e.createPropertyDescriptor()
	key := "set"
	if computedMember.IsGetter() {
		key = "get"
	}
	stringKey := e.astFactory.CreateStringKey(key, propertyValue)
	propertyDescriptor.AddChildToBack(stringKey)

	objectDefinePropertyCall :=
		e.astFactory.CreateCall(
			e.createObjectDotDefineProperty(),
			AstFactory_Type_(owner),
			owner.CloneTree(),
			property,
			propertyDescriptor)

	metadata.insertNodeAndAdvance(
		rhino.IR_ExprResult(objectDefinePropertyCall).SrcrefTreeIfMissing(memberNode))
}

// visitNonMethodMember visits non-computed getters and setters, including both
// static and instance properties.
//
// Non-computed getters and setters are aggregated into two
// Object.defineProperties calls. One defines static members and the other
// instance members.
//
// Java: Es6RewriteClass.visitNonMethodMember(Node, ClassDeclarationMetadata)
func (e *Es6RewriteClass) visitNonMethodMember(member rhino.Node, metadata *es6RewriteClassMetadata) {
	e.addToDefinePropertiesObject(metadata, member)

	if !member.IsStaticMember() {
		return
	}
	// Add stub declarations of static properties so that they are not broken by property collapsing
	membersToDeclare := metadata.membersToDeclare()
	memberName := member.GetString()
	membersToDeclare.Put(memberName, es6RewriteClassProperty{
		propertyKey:  memberName,
		quoted:       member.IsQuotedStringKey(),
		jsDocInfo:    es6RewriteClass_noCollapseJSDoc,
		propertyType: member.GetColor(),
	})
}

// visitMethod handles transpilation of a standard class member function.
// Getters, setters, and the constructor are not handled here.
//
// Java: Es6RewriteClass.visitMethod(Node, ClassDeclarationMetadata)
func (e *Es6RewriteClass) visitMethod(member rhino.Node, metadata *es6RewriteClassMetadata) {
	e.flushPendingDefineProperties(metadata, member.IsStaticMember(), member)

	qualifiedMemberAccess := e.getQualifiedMemberAccess(member, metadata)
	method := member.GetLastChild().Detach()

	// Use the source info from the method (a FUNCTION) not the MEMBER_FUNCTION_DEF
	// because the MEMBER_FUNCTION_DEF source info only corresponds to the identifier
	assign := e.astFactory.CreateAssign(qualifiedMemberAccess, method).SrcrefIfMissing(method)

	info := member.GetJSDocInfo()
	if member.IsStaticMember() && NodeUtil_ReferencesOwnReceiver(method) {
		memberDoc := rhino.JSDocInfo_Builder_MaybeCopyFrom(info)
		// Go-only: language-adapter — rhino.JSTypeExpression holds a non-null
		// source name; a null getSourceFileName() becomes "".
		sourceName := ""
		if s := member.GetSourceFileName(); s != nil {
			sourceName = *s
		}
		// adding an @this type prevents a JSC_UNSAFE_THIS error later on.
		memberDoc.RecordThisType(
			rhino.NewJSTypeExpression(
				rhino.NewNodeWithChild(rhino.Token_BANG, rhino.NewNode(rhino.Token_QMARK)).SrcrefTree(member),
				sourceName))
		info = memberDoc.Build()
	}
	if info != nil {
		assign.SetJSDocInfo(info)
	}

	metadata.insertNodeAndAdvance(NodeUtil_NewExpr(assign))
}

// addTypeDeclarations adds declarations for static properties defined with
// a getter or setter, so that optimizations like property collapsing
// recognize the properties' existence.
//
// Java: Es6RewriteClass.addTypeDeclarations(ClassDeclarationMetadata, Node)
func (e *Es6RewriteClass) addTypeDeclarations(metadata *es6RewriteClassMetadata, insertionPoint rhino.Node) {
	if metadata.classMembersToDeclare == nil {
		return
	}
	for property := range metadata.classMembersToDeclare.Values() {
		declaration := property.declaration(e.astFactory, metadata.fullClassNameNode.CloneTree())
		declaration.SrcrefTreeIfMissing(metadata.classNameNode)
		declaration.InsertAfter(insertionPoint)
		insertionPoint = declaration
	}
}

// getQualifiedMemberAccess constructs a Node that represents an access to
// the given class member, qualified by either the static or the instance
// access context, depending on whether the member is static.
//
// WARNING: member may be modified/destroyed by this method, do not use it
// afterwards.
//
// Java: Es6RewriteClass.getQualifiedMemberAccess(Node, ClassDeclarationMetadata)
func (e *Es6RewriteClass) getQualifiedMemberAccess(member rhino.Node, metadata *es6RewriteClassMetadata) rhino.Node {
	var context rhino.Node
	if member.IsStaticMember() {
		context = metadata.fullClassNameNode.CloneTree()
	} else {
		context = metadata.classPrototypeNode.CloneTree()
	}

	context.MakeNonIndexableRecursive()
	if prop, ok := ast.AsComputedProp(member); ok {
		return e.astFactory.
			CreateGetElem(context, prop.Key().Detach().Node()).
			SrcrefTreeIfMissing(member)
	}
	mfd := ast.MustMemberFunctionDef(member)
	methodName := mfd.Function().NameNode().Node()
	return e.astFactory.
		CreateGetProp(context, mfd.Key(), AstFactory_Type_(member)).
		SrcrefTree(methodName)
}

// Java: Es6RewriteClass.createPropertyDescriptor()
func (e *Es6RewriteClass) createPropertyDescriptor() rhino.Node {
	return e.astFactory.CreateObjectLit(
		e.astFactory.CreateStringKey("configurable", e.astFactory.CreateBoolean(true)),
		e.astFactory.CreateStringKey("enumerable", e.astFactory.CreateBoolean(true)))
}

// es6RewriteClassProperty describes a static non-computed getter/setter stub
// property to declare after the class so property collapsing preserves it.
type es6RewriteClassProperty struct {
	propertyKey  string
	quoted       bool
	jsDocInfo    *rhino.JSDocInfo
	propertyType *colors.Color
}

// Java: Es6RewriteClass.ClassProperty.getDeclaration(AstFactory, Node)
func (c es6RewriteClassProperty) declaration(astFactory *AstFactory, toDeclareOn rhino.Node) rhino.Node {
	var decl rhino.Node
	if c.quoted {
		decl = astFactory.CreateGetElem(toDeclareOn, astFactory.CreateString(c.propertyKey))
	} else {
		decl = astFactory.CreateGetProp(toDeclareOn, c.propertyKey, AstFactory_Type_OfColor(c.propertyType))
	}
	decl.SetJSDocInfo(c.jsDocInfo)
	return astFactory.ExprResult(decl)
}

// es6RewriteClassInsertionPoint tracks the statement node after which the next
// transpiled method of a class should be inserted.
type es6RewriteClassInsertionPoint struct {
	node rhino.Node
}

func (p *es6RewriteClassInsertionPoint) insertNodeAndAdvance(newNode rhino.Node) {
	newNode.InsertAfter(p.node)
	p.node = newNode
}

// es6RewriteClassMetadata holds the per-class transpilation state for
// Es6RewriteClass.VisitClass.
type es6RewriteClassMetadata struct {
	insertionPoint                  es6RewriteClassInsertionPoint
	astFactory                      *AstFactory
	definePropertiesObjForPrototype rhino.Node
	definePropertiesObjForClass     rhino.Node
	classMembersToDeclare           *jutil.LinkedHashMap[string, es6RewriteClassProperty]
	fullClassNameNode               rhino.Node
	classPrototypeNode              rhino.Node
	anonymous                       bool
	classNameNode                   rhino.Node
	superClassNameNode              rhino.Node
}

func (m *es6RewriteClassMetadata) definePropertiesForPrototype() rhino.Node {
	if m.definePropertiesObjForPrototype.IsNil() && m.astFactory != nil {
		m.definePropertiesObjForPrototype =
			m.astFactory.CreateObjectLitWithType(AstFactory_Type_(m.fullClassNameNode))
	}
	return m.definePropertiesObjForPrototype
}

func (m *es6RewriteClassMetadata) definePropertiesForClass() rhino.Node {
	if m.definePropertiesObjForClass.IsNil() && m.astFactory != nil {
		m.definePropertiesObjForClass =
			m.astFactory.CreateObjectLitWithType(AstFactory_Type_(m.fullClassNameNode))
	}
	return m.definePropertiesObjForClass
}

func (m *es6RewriteClassMetadata) membersToDeclare() *jutil.LinkedHashMap[string, es6RewriteClassProperty] {
	if m.classMembersToDeclare == nil {
		m.classMembersToDeclare = jutil.NewLinkedHashMap[string, es6RewriteClassProperty]()
	}
	return m.classMembersToDeclare
}

func (m *es6RewriteClassMetadata) insertNodeAndAdvance(newNode rhino.Node) {
	m.insertionPoint.insertNodeAndAdvance(newNode)
}

func (m *es6RewriteClassMetadata) hasSuperClass() bool {
	return !m.superClassNameNode.IsEmpty()
}

// Java: Es6RewriteClass.ClassDeclarationMetadata.create(Node, Node, AstFactory)
func es6RewriteClassCreateMetadata(
	m *es6RewriteClassMetadata,
	cls ast.Class,
	parent rhino.Node,
	astFactory *AstFactory) *es6RewriteClassMetadata {
	classNode := cls.Node()
	classNameNode := cls.NameNode()
	superClassNameNode := cls.ExtendsNode()

	var insertionStart rhino.Node
	var fullClassNameNode rhino.Node
	var anonymous bool

	// If this is a class statement, or a class expression in a simple
	// assignment or var statement, convert it. In any other case, the
	// code is too dynamic, so return nil.
	if NodeUtil_IsClassDeclaration(classNode) {
		insertionStart = classNode
		fullClassNameNode = classNameNode
		anonymous = false
	} else if assign, ok := ast.AsAssign(parent); ok && parent.GetParent().IsExprResult() {
		// Add members after the EXPR_RESULT node:
		// example.C = class {}; example.C.prototype.foo = function() {};
		fullClassNameNode = assign.Target().Node()
		if !fullClassNameNode.IsQualifiedName() {
			return nil
		}
		insertionStart = parent.GetParent()
		anonymous = true
	} else if parent.IsExport() {
		insertionStart = classNode
		fullClassNameNode = classNameNode
		anonymous = false
	} else if parent.IsName() {
		// Add members after the 'var' statement.
		// var C = class {}; C.prototype.foo = function() {};
		insertionStart = parent.GetParent()
		fullClassNameNode = parent.CloneNode() // specifically don't want children
		anonymous = true
	} else {
		// Cannot handle this class declaration.
		return nil
	}

	m.insertionPoint.node = insertionStart
	m.astFactory = astFactory
	m.definePropertiesObjForPrototype = rhino.NilNode
	m.definePropertiesObjForClass = rhino.NilNode
	m.classMembersToDeclare = nil
	m.fullClassNameNode = fullClassNameNode
	m.classPrototypeNode = astFactory.CreatePrototypeAccess(fullClassNameNode.CloneTree())
	m.anonymous = anonymous
	m.classNameNode = classNameNode
	m.superClassNameNode = superClassNameNode
	return m
}
