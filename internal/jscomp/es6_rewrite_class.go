// Port of com/google/javascript/jscomp/Es6RewriteClass.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/AutoBuilder_Es6RewriteClass_ClassProperty_Builder.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Es6RewriteClass converts ES6 classes to valid ES5 or ES3 code.
type Es6RewriteClass struct {
	compiler                     AbstractCompiler
	astFactory                   *AstFactory
	convertSuperConstructorCalls *Es6ConvertSuperConstructorCalls
	transpilationNamespace       rhino.StaticScope
	jscompInherits               RuntimeJsLibManager_JsLibField
	isES3Output                  bool
	metadataBuf                  Es6RewriteClass_ClassDeclarationMetadata
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
func (e *Es6RewriteClass) Process(externs *rhino.Node, root *rhino.Node) {
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
	defer func() {
		ns.Release()
		e.convertSuperConstructorCalls.SetGlobalNamespace(nil)
	}()
	e.convertSuperConstructorCalls.SetGlobalNamespace(ns)
	NodeTraversal_Traverse(e.compiler, root, e.convertSuperConstructorCalls)
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(e.compiler, root, es6RewriteClass_features)
}

// Java: Es6RewriteClass.shouldTraverse(NodeTraversal, Node, Node)
func (e *Es6RewriteClass) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
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
func (e *Es6RewriteClass) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if n.IsClass() {
		e.VisitClass(t, n, parent)
	}
}

// VisitClass processes classes in 3 phases:
//
//  1. The class name is extracted.
//  2. Class members are processed and rewritten.
//  3. The constructor is built.
//
// Java: Es6RewriteClass.visitClass(NodeTraversal, Node, Node)
func (e *Es6RewriteClass) VisitClass(t *NodeTraversal, classNode *rhino.Node, parent *rhino.Node) {
	// Collect Metadata
	metadata := es6RewriteClass_ClassDeclarationMetadata_CreateInto(&e.metadataBuf, classNode, parent, e.astFactory)

	if metadata == nil {
		panic(jutil.NewIllegalStateExceptionMsg(
			"Can only convert classes that are declarations or the right hand" +
				" side of a simple assignment: " +
				jutil.ToJavaString(classNode)))
	}
	if metadata.HasSuperClass() {
		jutil.CheckStateMsg1(
			metadata.SuperClassNameNode().IsQualifiedName(),
			"Expected Es6NormalizeClasses to make all extends clauses into qualified  names, found"+
				" %s",
			metadata.SuperClassNameNode())
	}

	jutil.CheckStateMsg1(
		NodeUtil_IsStatement(metadata.InsertionPoint().GetNode()),
		"insertion point must be a statement: %s",
		metadata.InsertionPoint().GetNode())

	var constructor *rhino.Node
	// Process all members of the class
	classMembers := classNode.GetLastChild()
	for member := classMembers.GetFirstChild(); member != nil; {
		next := member.GetNext()
		if (member.IsComputedProp() &&
			(member.GetBooleanProp(rhino.Node_COMPUTED_PROP_GETTER) ||
				member.GetBooleanProp(rhino.Node_COMPUTED_PROP_SETTER))) ||
			(member.IsGetterDef() || member.IsSetterDef()) {
			e.VisitNonMethodMember(member, metadata)
		} else if NodeUtil_IsEs6ConstructorMemberFunctionDef(member) {
			if !classNode.IsFromExterns() {
				e.ReportPrimitiveReturns(member.GetFirstChild())
			}
			constructor = member.RemoveFirstChild().SetColor(classNode.GetColor())
			if !metadata.Anonymous() {
				// Turns class Foo { constructor: function() {} } into function Foo() {},
				// i.e. attaches the name to the ctor function.
				constructor.GetFirstChild().ReplaceWith(metadata.ClassNameNode().CloneNode())
			}
		} else if member.IsEmpty() {
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
			e.VisitMethod(member, metadata)
		}
		member = next
	}
	jutil.CheckNotNullMsg(
		constructor,
		"Es6RewriteClasses expects all classes to have (possibly synthetic) constructors")

	e.FlushDefinePropertiesForPrototype(metadata, classNode)
	e.FlushDefinePropertiesForClass(metadata, classNode)

	classJSDoc := NodeUtil_GetBestJSDocInfo(classNode)
	newInfo := rhino.JSDocInfo_Builder_MaybeCopyFrom(classJSDoc)
	newInfo.RecordConstructor()

	enclosingStatement := NodeUtil_GetEnclosingStatement(classNode)
	if metadata.HasSuperClass() {
		if !classNode.IsFromExterns() {
			inheritsCall :=
				rhino.IR_ExprResult(
					e.astFactory.CreateCall(
						e.astFactory.CreateQNameOfField(e.transpilationNamespace, e.jscompInherits),
						AstFactory_Type_OfColor(colors.StandardColors_NULL_OR_VOID),
						metadata.FullClassNameNode().CloneTree(),
						metadata.SuperClassNameNode().CloneTree())).
					SrcrefTreeIfMissing(metadata.SuperClassNameNode())
			inheritsCall.InsertAfter(enclosingStatement)
		}
	}

	e.AddTypeDeclarations(metadata, enclosingStatement)

	if NodeUtil_IsStatement(classNode) {
		constructor.GetFirstChild().SetString("")
		ctorVar := rhino.IR_LetWithValue(metadata.ClassNameNode().CloneNode(), constructor)
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
		var_ := parent.GetParent()
		var_.SetJSDocInfo(newInfo.Build())
	} else if constructor.GetParent().IsName() {
		// Is a newly created VAR node.
		var_ := constructor.GetGrandparent()
		var_.SetJSDocInfo(newInfo.Build())
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
func (e *Es6RewriteClass) ReportPrimitiveReturns(constructor *rhino.Node) {
	// TODO(b/573571176): A derived constructor that returns `undefined`, explicitly or by reaching
	// the end of its body, before calling `super()` must throw a ReferenceError.
	NodeUtil_VisitPostOrderWithPred(
		NodeUtil_GetFunctionBody(constructor),
		NodeUtil_VisitorFunc(func(n *rhino.Node) {
			if n.IsReturn() &&
				n.HasChildren() &&
				NodeUtil_IsImmutableValue(n.GetFirstChild()) &&
				NodeUtil_GetKnownValueType(n.GetFirstChild()) != NodeUtil_ValueType_VOID {
				e.compiler.Report(JSError_MakeWithNode(n, Es6RewriteClass_CONSTRUCTOR_PRIMITIVE_RETURN))
			}
		}),
		// A `return` in a nested function belongs to that function.
		func(n *rhino.Node) bool { return !n.IsFunction() })
}

// Java: Es6RewriteClass.createObjectDotDefineProperties()
func (e *Es6RewriteClass) CreateObjectDotDefineProperties() *rhino.Node {
	return e.astFactory.CreateJSCompDotGlobalAccess(
		e.transpilationNamespace, "Object.defineProperties")
}

// Java: Es6RewriteClass.createObjectDotDefineProperty()
func (e *Es6RewriteClass) CreateObjectDotDefineProperty() *rhino.Node {
	return e.astFactory.CreateJSCompDotGlobalAccess(
		e.transpilationNamespace, "Object.defineProperty")
}

// Java: Es6RewriteClass.flushDefinePropertiesForPrototype(ClassDeclarationMetadata, Node)
func (e *Es6RewriteClass) FlushDefinePropertiesForPrototype(
	metadata *Es6RewriteClass_ClassDeclarationMetadata, srcrefNode *rhino.Node) {
	if metadata.definePropertiesObjForPrototype != nil && metadata.definePropertiesObjForPrototype.HasChildren() {
		definePropsCall :=
			rhino.IR_ExprResult(
				e.astFactory.CreateCall(
					e.CreateObjectDotDefineProperties(),
					AstFactory_Type_(metadata.ClassPrototypeNode()),
					metadata.ClassPrototypeNode().CloneTree(),
					metadata.definePropertiesObjForPrototype))
		definePropsCall.SrcrefTreeIfMissing(srcrefNode)
		metadata.InsertNodeAndAdvance(definePropsCall)
		metadata.definePropertiesObjForPrototype = nil
	}
}

// Java: Es6RewriteClass.flushDefinePropertiesForClass(ClassDeclarationMetadata, Node)
func (e *Es6RewriteClass) FlushDefinePropertiesForClass(metadata *Es6RewriteClass_ClassDeclarationMetadata, srcrefNode *rhino.Node) {
	if metadata.definePropertiesObjForClass != nil && metadata.definePropertiesObjForClass.HasChildren() {
		definePropsCall :=
			rhino.IR_ExprResult(
				e.astFactory.CreateCall(
					e.CreateObjectDotDefineProperties(),
					AstFactory_Type_(metadata.FullClassNameNode()),
					metadata.FullClassNameNode().CloneTree(),
					metadata.definePropertiesObjForClass))
		definePropsCall.SrcrefTreeIfMissing(srcrefNode)
		metadata.InsertNodeAndAdvance(definePropsCall)
		metadata.definePropertiesObjForClass = nil
	}
}

// Java: Es6RewriteClass.flushPendingDefineProperties(ClassDeclarationMetadata, boolean, Node)
func (e *Es6RewriteClass) FlushPendingDefineProperties(
	metadata *Es6RewriteClass_ClassDeclarationMetadata, isStatic bool, srcrefNode *rhino.Node) {
	if isStatic {
		e.FlushDefinePropertiesForClass(metadata, srcrefNode)
	} else {
		e.FlushDefinePropertiesForPrototype(metadata, srcrefNode)
	}
}

// AddToDefinePropertiesObject takes member, a getter or setter.
//
// Java: Es6RewriteClass.addToDefinePropertiesObject(ClassDeclarationMetadata, Node)
func (e *Es6RewriteClass) AddToDefinePropertiesObject(metadata *Es6RewriteClass_ClassDeclarationMetadata, member *rhino.Node) {
	jutil.CheckArgument(!member.IsComputedProp())
	var obj *rhino.Node
	if member.IsStaticMember() {
		obj = metadata.DefinePropertiesObjForClass()
	} else {
		obj = metadata.DefinePropertiesObjForPrototype()
	}
	prop := NodeUtil_GetFirstPropMatchingKey(obj, member.GetString())
	accessorKind := "set"
	if member.IsGetterDef() {
		accessorKind = "get"
	}
	if prop != nil && NodeUtil_GetFirstPropMatchingKey(prop, accessorKind) != nil {
		e.FlushPendingDefineProperties(metadata, member.IsStaticMember(), member)
		if member.IsStaticMember() {
			obj = metadata.DefinePropertiesObjForClass()
		} else {
			obj = metadata.DefinePropertiesObjForPrototype()
		}
		prop = nil
	}
	if prop == nil {
		prop = e.CreatePropertyDescriptor()
		if member.GetString() == "__proto__" {
			stringKey :=
				e.astFactory.CreateComputedProperty(e.astFactory.CreateString("__proto__"), prop)
			obj.AddChildToBack(stringKey)
			scriptNode := NodeUtil_GetEnclosingScript(member)
			if scriptNode != nil {
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

	function := member.GetLastChild()
	info := NodeUtil_GetBestJSDocInfo(function)

	stringKey := e.astFactory.CreateStringKey(accessorKind, function.Detach())
	stringKey.SetJSDocInfo(info)
	prop.AddChildToBack(stringKey)
	prop.SrcrefTreeIfMissing(member)
}

// ExtractComputedProperty appends an Object.defineProperty call defining
// the given computed getter or setter.
//
// Java: Es6RewriteClass.extractComputedProperty(Node, ClassDeclarationMetadata)
func (e *Es6RewriteClass) ExtractComputedProperty(computedMember *rhino.Node, metadata *Es6RewriteClass_ClassDeclarationMetadata) {
	e.FlushPendingDefineProperties(metadata, computedMember.IsStaticMember(), computedMember)

	var owner *rhino.Node
	if computedMember.IsStaticMember() {
		owner = metadata.FullClassNameNode()
	} else {
		owner = metadata.ClassPrototypeNode()
	}
	property := computedMember.RemoveFirstChild()
	propertyValue := computedMember.RemoveFirstChild()

	propertyDescriptor := e.CreatePropertyDescriptor()
	key := "set"
	if computedMember.GetBooleanProp(rhino.Node_COMPUTED_PROP_GETTER) {
		key = "get"
	}
	stringKey := e.astFactory.CreateStringKey(key, propertyValue)
	propertyDescriptor.AddChildToBack(stringKey)

	objectDefinePropertyCall :=
		e.astFactory.CreateCall(
			e.CreateObjectDotDefineProperty(),
			AstFactory_Type_(owner),
			owner.CloneTree(),
			property,
			propertyDescriptor)

	metadata.InsertNodeAndAdvance(
		rhino.IR_ExprResult(objectDefinePropertyCall).SrcrefTreeIfMissing(computedMember))
}

// VisitNonMethodMember visits getters and setters, including both static
// and instance properties, and computed and non-computed properties.
//
// Non-computed getters and setters are aggregated into two
// Object.defineProperties calls. One defines static members and the other
// instance members. This is just an optimization; it's just as sound to
// append a single Object.defineProperty definition for each here.
//
// Computed getters and setters are defined in individual
// Object.defineProperty calls
//
// Java: Es6RewriteClass.visitNonMethodMember(Node, ClassDeclarationMetadata)
func (e *Es6RewriteClass) VisitNonMethodMember(member *rhino.Node, metadata *Es6RewriteClass_ClassDeclarationMetadata) {
	if member.IsComputedProp() {
		e.ExtractComputedProperty(member.Detach(), metadata)
		return
	}

	e.AddToDefinePropertiesObject(metadata, member)

	if !member.IsStaticMember() {
		return
	}
	jutil.CheckStateMsg1(!member.IsComputedProp(), "%s", member)
	// Add stub declarations of static properties so that they are not broken by property collapsing

	membersToDeclare := metadata.ClassMembersToDeclare()
	memberName := member.GetString()
	kind := Es6RewriteClass_ClassProperty_PropertyKind_NORMAL_PROPERTY
	if member.IsQuotedStringKey() {
		kind = Es6RewriteClass_ClassProperty_PropertyKind_QUOTED_PROPERTY
	}
	membersToDeclare.Put(memberName, &Es6RewriteClass_ClassProperty{
		propertyKey:  memberName,
		kind:         kind,
		jsDocInfo:    es6RewriteClass_noCollapseJSDoc,
		propertyType: member.GetColor(),
	})
}

// VisitMethod handles transpilation of a standard class member function.
// Getters, setters, and the constructor are not handled here.
//
// Java: Es6RewriteClass.visitMethod(Node, ClassDeclarationMetadata)
func (e *Es6RewriteClass) VisitMethod(member *rhino.Node, metadata *Es6RewriteClass_ClassDeclarationMetadata) {
	e.FlushPendingDefineProperties(metadata, member.IsStaticMember(), member)

	qualifiedMemberAccess := e.GetQualifiedMemberAccess(member, metadata)
	method := member.GetLastChild().Detach()

	// Use the source info from the method (a FUNCTION) not the MEMBER_FUNCTION_DEF
	// because the MEMBER_FUNCTION_DEf source info only corresponds to the identifier
	assign := e.astFactory.CreateAssign(qualifiedMemberAccess, method).SrcrefIfMissing(method)

	info := member.GetJSDocInfo()
	if member.IsStaticMember() && NodeUtil_ReferencesOwnReceiver(assign.GetLastChild()) {
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

	newNode := NodeUtil_NewExpr(assign)
	metadata.InsertNodeAndAdvance(newNode)
}

// AddTypeDeclarations adds declarations for static properties defined with
// a getter or setter, so that optimizations like property collapsing
// recognize the properties' existence.
//
// Java: Es6RewriteClass.addTypeDeclarations(ClassDeclarationMetadata, Node)
func (e *Es6RewriteClass) AddTypeDeclarations(metadata *Es6RewriteClass_ClassDeclarationMetadata, insertionPoint *rhino.Node) {
	if metadata.classMembersToDeclare == nil {
		return
	}
	for property := range metadata.classMembersToDeclare.Values() {
		declaration :=
			property.GetDeclaration(e.astFactory, metadata.FullClassNameNode().CloneTree())
		declaration.SrcrefTreeIfMissing(metadata.ClassNameNode())
		declaration.InsertAfter(insertionPoint)
		insertionPoint = declaration
	}
}

// GetQualifiedMemberAccess constructs a Node that represents an access to
// the given class member, qualified by either the static or the instance
// access context, depending on whether the member is static.
//
// WARNING: member may be modified/destroyed by this method, do not use it
// afterwards.
//
// Java: Es6RewriteClass.getQualifiedMemberAccess(Node, ClassDeclarationMetadata)
func (e *Es6RewriteClass) GetQualifiedMemberAccess(member *rhino.Node, metadata *Es6RewriteClass_ClassDeclarationMetadata) *rhino.Node {
	var context *rhino.Node
	if member.IsStaticMember() {
		context = metadata.FullClassNameNode().CloneTree()
	} else {
		context = metadata.ClassPrototypeNode().CloneTree()
	}

	context.MakeNonIndexableRecursive()
	if member.IsComputedProp() {
		return e.astFactory.
			CreateGetElem(context, member.RemoveFirstChild()).
			SrcrefTreeIfMissing(member)
	} else {
		methodName := member.GetFirstFirstChild()
		return e.astFactory.
			CreateGetProp(context, member.GetString(), AstFactory_Type_(member)).
			SrcrefTree(methodName)
	}
}

// Java: Es6RewriteClass.createPropertyDescriptor()
func (e *Es6RewriteClass) CreatePropertyDescriptor() *rhino.Node {
	return e.astFactory.CreateObjectLit(
		e.astFactory.CreateStringKey("configurable", e.astFactory.CreateBoolean(true)),
		e.astFactory.CreateStringKey("enumerable", e.astFactory.CreateBoolean(true)))
}

// Es6RewriteClass_ClassProperty is the record ClassProperty: propertyKey is
// the name of this ClassProperty for NORMAL_PROPERTY, the string value of
// this property if QUOTED_PROPERTY, or the qualified name of the computed
// property.
type Es6RewriteClass_ClassProperty struct {
	propertyKey  string
	kind         Es6RewriteClass_ClassProperty_PropertyKind
	jsDocInfo    *rhino.JSDocInfo
	propertyType *colors.Color
}

// Java: Es6RewriteClass.ClassProperty(String, PropertyKind, JSDocInfo, Color)
func NewEs6RewriteClass_ClassProperty(
	propertyKey string,
	kind Es6RewriteClass_ClassProperty_PropertyKind,
	jsDocInfo *rhino.JSDocInfo,
	propertyType *colors.Color) *Es6RewriteClass_ClassProperty {
	// Go-only: language-adapter — requireNonNull(propertyKey) and
	// requireNonNull(kind): a Go string and enum value cannot be null.
	jutil.CheckNotNullMsg(jsDocInfo, "jsDocInfo")
	return &Es6RewriteClass_ClassProperty{
		propertyKey:  propertyKey,
		kind:         kind,
		jsDocInfo:    jsDocInfo,
		propertyType: propertyType,
	}
}

// Java: Es6RewriteClass.ClassProperty.propertyKey()
func (c *Es6RewriteClass_ClassProperty) PropertyKey() string { return c.propertyKey }

// Java: Es6RewriteClass.ClassProperty.kind()
func (c *Es6RewriteClass_ClassProperty) Kind() Es6RewriteClass_ClassProperty_PropertyKind {
	return c.kind
}

// Java: Es6RewriteClass.ClassProperty.jsDocInfo()
func (c *Es6RewriteClass_ClassProperty) JsDocInfo() *rhino.JSDocInfo { return c.jsDocInfo }

// Java: Es6RewriteClass.ClassProperty.propertyType()
func (c *Es6RewriteClass_ClassProperty) PropertyType() *colors.Color { return c.propertyType }

// String is the record's toString().
//
// Java: Es6RewriteClass.ClassProperty.toString()
func (c *Es6RewriteClass_ClassProperty) String() string {
	return "ClassProperty[propertyKey=" + c.propertyKey +
		", kind=" + c.kind.String() +
		", jsDocInfo=" + jutil.ToJavaString(c.jsDocInfo) +
		", propertyType=" + jutil.ToJavaString(c.propertyType) + "]"
}

// Es6RewriteClass_ClassProperty_PropertyKind is ClassProperty.PropertyKind.
type Es6RewriteClass_ClassProperty_PropertyKind int32

const (
	// Any kind of quoted property, which can include numeric properties
	// that we treated as quoted.
	//
	//	class Example {
	//	  'quoted'() {}
	//	  42() { return 'the answer'; }
	//	}
	Es6RewriteClass_ClassProperty_PropertyKind_QUOTED_PROPERTY Es6RewriteClass_ClassProperty_PropertyKind = iota

	// A computed property, e.g. using bracket [] access.
	//
	// Computed properties *must* currently be qualified names, and not
	// literals, function calls, etc.
	//
	//	class Example {
	//	  [variable]() {}
	//	  [another.example]() {}
	//	}
	Es6RewriteClass_ClassProperty_PropertyKind_COMPUTED_PROPERTY

	// A normal property definition.
	//
	//	class Example {
	//	  normal() {}
	//	}
	Es6RewriteClass_ClassProperty_PropertyKind_NORMAL_PROPERTY
)

var es6RewriteClass_ClassProperty_PropertyKind_names = [...]string{"QUOTED_PROPERTY", "COMPUTED_PROPERTY", "NORMAL_PROPERTY"}

// Es6RewriteClass_ClassProperty_PropertyKind_values is PropertyKind.values().
func Es6RewriteClass_ClassProperty_PropertyKind_values() []Es6RewriteClass_ClassProperty_PropertyKind {
	return []Es6RewriteClass_ClassProperty_PropertyKind{
		Es6RewriteClass_ClassProperty_PropertyKind_QUOTED_PROPERTY,
		Es6RewriteClass_ClassProperty_PropertyKind_COMPUTED_PROPERTY,
		Es6RewriteClass_ClassProperty_PropertyKind_NORMAL_PROPERTY,
	}
}

// String is Enum.name().
func (k Es6RewriteClass_ClassProperty_PropertyKind) String() string {
	return es6RewriteClass_ClassProperty_PropertyKind_names[k]
}

// Ordinal is Enum.ordinal().
func (k Es6RewriteClass_ClassProperty_PropertyKind) Ordinal() int { return int(k) }

// GetDeclaration returns an EXPR_RESULT node that declares this property on
// the given node.
//
// Examples:
//
//	/** @type {string} */
//	Class.prototype.property;
//
//	/** @type {string} */
//	Class.staticProperty;
//
// toDeclareOn is the node to declare the property on. This should either be
// a reference to the class (if a static property) or a class' prototype (if
// non-static). This should always be a new node, as this method will insert
// it into the returned EXPR_RESULT.
//
// Java: Es6RewriteClass.ClassProperty.getDeclaration(AstFactory, Node)
func (c *Es6RewriteClass_ClassProperty) GetDeclaration(astFactory *AstFactory, toDeclareOn *rhino.Node) *rhino.Node {
	var decl *rhino.Node
	switch c.Kind() {
	case Es6RewriteClass_ClassProperty_PropertyKind_QUOTED_PROPERTY:
		decl = astFactory.CreateGetElem(toDeclareOn, astFactory.CreateString(c.PropertyKey()))
	case Es6RewriteClass_ClassProperty_PropertyKind_COMPUTED_PROPERTY:
		// No need to declare computed properties as they're unaffected by property
		// collapsing
		panic(jutil.NewUnsupportedOperationExceptionMsg(c.String()))
	case Es6RewriteClass_ClassProperty_PropertyKind_NORMAL_PROPERTY:
		decl = astFactory.CreateGetProp(toDeclareOn, c.PropertyKey(), AstFactory_Type_OfColor(c.PropertyType()))
	}

	decl.SetJSDocInfo(c.JsDocInfo())
	decl = astFactory.ExprResult(decl)
	return decl
}

// Java: Es6RewriteClass.ClassProperty.builder()
func Es6RewriteClass_ClassProperty_Builder_() *Es6RewriteClass_ClassProperty_Builder {
	return newAutoBuilder_Es6RewriteClass_ClassProperty_Builder()
}

// Es6RewriteClass_ClassProperty_Builder is ClassProperty.Builder.
//
// Go-only: language-adapter — the AutoBuilder-generated
// AutoBuilder_Es6RewriteClass_ClassProperty_Builder is folded into it.
type Es6RewriteClass_ClassProperty_Builder struct {
	propertyKey  *string
	kind         *Es6RewriteClass_ClassProperty_PropertyKind
	jsDocInfo    *rhino.JSDocInfo
	propertyType *colors.Color
}

// Java: AutoBuilder_Es6RewriteClass_ClassProperty_Builder()
func newAutoBuilder_Es6RewriteClass_ClassProperty_Builder() *Es6RewriteClass_ClassProperty_Builder {
	return &Es6RewriteClass_ClassProperty_Builder{}
}

// Java: AutoBuilder_Es6RewriteClass_ClassProperty_Builder(Es6RewriteClass.ClassProperty)
func newAutoBuilder_Es6RewriteClass_ClassProperty_BuilderFrom(source *Es6RewriteClass_ClassProperty) *Es6RewriteClass_ClassProperty_Builder {
	b := &Es6RewriteClass_ClassProperty_Builder{}
	propertyKey := source.PropertyKey()
	b.propertyKey = &propertyKey
	kind := source.Kind()
	b.kind = &kind
	b.jsDocInfo = source.JsDocInfo()
	b.propertyType = source.PropertyType()
	return b
}

// Java: AutoBuilder_Es6RewriteClass_ClassProperty_Builder.propertyKey(String)
func (b *Es6RewriteClass_ClassProperty_Builder) PropertyKey(propertyKey string) *Es6RewriteClass_ClassProperty_Builder {
	b.propertyKey = &propertyKey
	return b
}

// Java: AutoBuilder_Es6RewriteClass_ClassProperty_Builder.kind(Es6RewriteClass.ClassProperty.PropertyKind)
func (b *Es6RewriteClass_ClassProperty_Builder) Kind(kind Es6RewriteClass_ClassProperty_PropertyKind) *Es6RewriteClass_ClassProperty_Builder {
	b.kind = &kind
	return b
}

// Java: AutoBuilder_Es6RewriteClass_ClassProperty_Builder.jsDocInfo(JSDocInfo)
func (b *Es6RewriteClass_ClassProperty_Builder) JsDocInfo(jsDocInfo *rhino.JSDocInfo) *Es6RewriteClass_ClassProperty_Builder {
	if jsDocInfo == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null jsDocInfo"))
	}
	b.jsDocInfo = jsDocInfo
	return b
}

// Java: AutoBuilder_Es6RewriteClass_ClassProperty_Builder.propertyType(Color)
func (b *Es6RewriteClass_ClassProperty_Builder) PropertyType(propertyType *colors.Color) *Es6RewriteClass_ClassProperty_Builder {
	b.propertyType = propertyType
	return b
}

// Java: AutoBuilder_Es6RewriteClass_ClassProperty_Builder.build()
func (b *Es6RewriteClass_ClassProperty_Builder) Build() *Es6RewriteClass_ClassProperty {
	if b.propertyKey == nil ||
		b.kind == nil ||
		b.jsDocInfo == nil {
		missing := jutil.NewStringBuilder()
		if b.propertyKey == nil {
			missing.Append(" propertyKey")
		}
		if b.kind == nil {
			missing.Append(" kind")
		}
		if b.jsDocInfo == nil {
			missing.Append(" jsDocInfo")
		}
		panic(jutil.NewIllegalStateExceptionMsg("Missing required properties:" + missing.String()))
	}
	return NewEs6RewriteClass_ClassProperty(
		*b.propertyKey,
		*b.kind,
		b.jsDocInfo,
		b.propertyType)
}

// Es6RewriteClass_ClassDeclarationMetadata represents static metadata on a
// class declaration expression - i.e. the qualified name that a class
// declares (directly or by assignment), whether it's anonymous, and where
// transpiled code should be inserted (i.e. which object will hold the
// prototype after transpilation).
//
// Note that this class is NOT deeply immutable! Don't use it in a Map. The
// AutoValue(.Builder) is just used to simplify creating instances.
//
// insertionPoint: a statement node. Transpiled methods etc of the class are
// inserted after this node. definePropertiesObjForPrototype: an object
// literal node that will be used in a call to Object.defineProperties, to
// add getters and setters to the prototype. definePropertiesObjForClass: an
// object literal node that will be used in a call to
// Object.defineProperties, to add getters and setters to the class.
// classMembersToDeclare: property declarations to be added to the class.
// fullClassNameNode: the fully qualified name of the class, as a cloneable
// node. May come from the class itself or the LHS of an assignment.
// classPrototypeNode: the fully qualified name of this class, plus
// ".prototype", as a cloneable node with type information as needed.
// anonymous: whether the constructor function in the output should be
// anonymous.
type Es6RewriteClass_ClassDeclarationMetadata struct {
	insertionPoint                  *Es6RewriteClass_InsertionPoint
	insertionPointVal               Es6RewriteClass_InsertionPoint
	astFactory                      *AstFactory
	definePropertiesObjForPrototype *rhino.Node
	definePropertiesObjForClass     *rhino.Node
	classMembersToDeclare           *jutil.LinkedHashMap[string, *Es6RewriteClass_ClassProperty]
	fullClassNameNode               *rhino.Node
	classPrototypeNode              *rhino.Node
	anonymous                       bool
	classNameNode                   *rhino.Node
	superClassNameNode              *rhino.Node
}

// Java: Es6RewriteClass.ClassDeclarationMetadata(InsertionPoint, Node, Node, Map<String, ClassProperty>, Node, Node, boolean, Node, Node)
func NewEs6RewriteClass_ClassDeclarationMetadata(
	insertionPoint *Es6RewriteClass_InsertionPoint,
	definePropertiesObjForPrototype *rhino.Node,
	definePropertiesObjForClass *rhino.Node,
	classMembersToDeclare *jutil.LinkedHashMap[string, *Es6RewriteClass_ClassProperty],
	fullClassNameNode *rhino.Node,
	classPrototypeNode *rhino.Node,
	anonymous bool,
	classNameNode *rhino.Node,
	superClassNameNode *rhino.Node) *Es6RewriteClass_ClassDeclarationMetadata {
	jutil.CheckNotNullMsg(insertionPoint, "insertionPoint")
	jutil.CheckNotNullMsg(definePropertiesObjForPrototype, "definePropertiesObjForPrototype")
	jutil.CheckNotNullMsg(definePropertiesObjForClass, "definePropertiesObjForClass")
	jutil.CheckNotNullMsg(classMembersToDeclare, "classMembersToDeclare")
	jutil.CheckNotNullMsg(fullClassNameNode, "fullClassNameNode")
	jutil.CheckNotNullMsg(classPrototypeNode, "classPrototypeNode")
	jutil.CheckNotNullMsg(classNameNode, "classNameNode")
	jutil.CheckNotNullMsg(superClassNameNode, "superClassNameNode")
	return &Es6RewriteClass_ClassDeclarationMetadata{
		insertionPoint:                  insertionPoint,
		definePropertiesObjForPrototype: definePropertiesObjForPrototype,
		definePropertiesObjForClass:     definePropertiesObjForClass,
		classMembersToDeclare:           classMembersToDeclare,
		fullClassNameNode:               fullClassNameNode,
		classPrototypeNode:              classPrototypeNode,
		anonymous:                       anonymous,
		classNameNode:                   classNameNode,
		superClassNameNode:              superClassNameNode,
	}
}

// Java: Es6RewriteClass.ClassDeclarationMetadata.insertionPoint()
func (m *Es6RewriteClass_ClassDeclarationMetadata) InsertionPoint() *Es6RewriteClass_InsertionPoint {
	return m.insertionPoint
}

// Java: Es6RewriteClass.ClassDeclarationMetadata.definePropertiesObjForPrototype()
func (m *Es6RewriteClass_ClassDeclarationMetadata) DefinePropertiesObjForPrototype() *rhino.Node {
	if m.definePropertiesObjForPrototype == nil && m.astFactory != nil {
		m.definePropertiesObjForPrototype =
			m.astFactory.CreateObjectLitWithType(AstFactory_Type_(m.fullClassNameNode))
	}
	return m.definePropertiesObjForPrototype
}

// Java: Es6RewriteClass.ClassDeclarationMetadata.resetDefinePropertiesObjForPrototype(AstFactory)
func (m *Es6RewriteClass_ClassDeclarationMetadata) ResetDefinePropertiesObjForPrototype(astFactory *AstFactory) {
	m.definePropertiesObjForPrototype =
		astFactory.CreateObjectLitWithType(AstFactory_Type_(m.fullClassNameNode))
}

// Java: Es6RewriteClass.ClassDeclarationMetadata.definePropertiesObjForClass()
func (m *Es6RewriteClass_ClassDeclarationMetadata) DefinePropertiesObjForClass() *rhino.Node {
	if m.definePropertiesObjForClass == nil && m.astFactory != nil {
		m.definePropertiesObjForClass =
			m.astFactory.CreateObjectLitWithType(AstFactory_Type_(m.fullClassNameNode))
	}
	return m.definePropertiesObjForClass
}

// Java: Es6RewriteClass.ClassDeclarationMetadata.resetDefinePropertiesObjForClass(AstFactory)
func (m *Es6RewriteClass_ClassDeclarationMetadata) ResetDefinePropertiesObjForClass(astFactory *AstFactory) {
	m.definePropertiesObjForClass = astFactory.CreateObjectLitWithType(AstFactory_Type_(m.fullClassNameNode))
}

// Java: Es6RewriteClass.ClassDeclarationMetadata.classMembersToDeclare()
func (m *Es6RewriteClass_ClassDeclarationMetadata) ClassMembersToDeclare() *jutil.LinkedHashMap[string, *Es6RewriteClass_ClassProperty] {
	if m.classMembersToDeclare == nil {
		m.classMembersToDeclare = jutil.NewLinkedHashMap[string, *Es6RewriteClass_ClassProperty]()
	}
	return m.classMembersToDeclare
}

// Java: Es6RewriteClass.ClassDeclarationMetadata.fullClassNameNode()
func (m *Es6RewriteClass_ClassDeclarationMetadata) FullClassNameNode() *rhino.Node {
	return m.fullClassNameNode
}

// Java: Es6RewriteClass.ClassDeclarationMetadata.classPrototypeNode()
func (m *Es6RewriteClass_ClassDeclarationMetadata) ClassPrototypeNode() *rhino.Node {
	return m.classPrototypeNode
}

// Java: Es6RewriteClass.ClassDeclarationMetadata.anonymous()
func (m *Es6RewriteClass_ClassDeclarationMetadata) Anonymous() bool { return m.anonymous }

// Java: Es6RewriteClass.ClassDeclarationMetadata.classNameNode()
func (m *Es6RewriteClass_ClassDeclarationMetadata) ClassNameNode() *rhino.Node {
	return m.classNameNode
}

// Java: Es6RewriteClass.ClassDeclarationMetadata.superClassNameNode()
func (m *Es6RewriteClass_ClassDeclarationMetadata) SuperClassNameNode() *rhino.Node {
	return m.superClassNameNode
}

// Es6RewriteClass_ClassDeclarationMetadata_Builder is
// ClassDeclarationMetadata.Builder.
//
// Go-only: language-adapter — the AutoBuilder-generated
// AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder is folded
// into it.
type Es6RewriteClass_ClassDeclarationMetadata_Builder struct {
	insertionPoint                  *Es6RewriteClass_InsertionPoint
	definePropertiesObjForPrototype *rhino.Node
	definePropertiesObjForClass     *rhino.Node
	classMembersToDeclare           *jutil.LinkedHashMap[string, *Es6RewriteClass_ClassProperty]
	fullClassNameNode               *rhino.Node
	classPrototypeNode              *rhino.Node
	anonymous                       bool
	classNameNode                   *rhino.Node
	superClassNameNode              *rhino.Node
	set_0                           byte
}

// Java: AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder()
func newAutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder() *Es6RewriteClass_ClassDeclarationMetadata_Builder {
	return &Es6RewriteClass_ClassDeclarationMetadata_Builder{}
}

// Java: AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder(Es6RewriteClass.ClassDeclarationMetadata)
func newAutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_BuilderFrom(source *Es6RewriteClass_ClassDeclarationMetadata) *Es6RewriteClass_ClassDeclarationMetadata_Builder {
	b := &Es6RewriteClass_ClassDeclarationMetadata_Builder{}
	b.insertionPoint = source.InsertionPoint()
	b.definePropertiesObjForPrototype = source.DefinePropertiesObjForPrototype()
	b.definePropertiesObjForClass = source.DefinePropertiesObjForClass()
	b.classMembersToDeclare = source.ClassMembersToDeclare()
	b.fullClassNameNode = source.FullClassNameNode()
	b.classPrototypeNode = source.ClassPrototypeNode()
	b.anonymous = source.Anonymous()
	b.classNameNode = source.ClassNameNode()
	b.superClassNameNode = source.SuperClassNameNode()
	b.set_0 = 1
	return b
}

// Java: AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder.setInsertionPoint(Es6RewriteClass.InsertionPoint)
func (b *Es6RewriteClass_ClassDeclarationMetadata_Builder) SetInsertionPoint(insertionPoint *Es6RewriteClass_InsertionPoint) *Es6RewriteClass_ClassDeclarationMetadata_Builder {
	if insertionPoint == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null insertionPoint"))
	}
	b.insertionPoint = insertionPoint
	return b
}

// Java: AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder.setDefinePropertiesObjForPrototype(Node)
func (b *Es6RewriteClass_ClassDeclarationMetadata_Builder) SetDefinePropertiesObjForPrototype(definePropertiesObjForPrototype *rhino.Node) *Es6RewriteClass_ClassDeclarationMetadata_Builder {
	if definePropertiesObjForPrototype == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null definePropertiesObjForPrototype"))
	}
	b.definePropertiesObjForPrototype = definePropertiesObjForPrototype
	return b
}

// Java: AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder.setDefinePropertiesObjForClass(Node)
func (b *Es6RewriteClass_ClassDeclarationMetadata_Builder) SetDefinePropertiesObjForClass(definePropertiesObjForClass *rhino.Node) *Es6RewriteClass_ClassDeclarationMetadata_Builder {
	if definePropertiesObjForClass == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null definePropertiesObjForClass"))
	}
	b.definePropertiesObjForClass = definePropertiesObjForClass
	return b
}

// Java: AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder.setClassMembersToDeclare(Map<String, Es6RewriteClass.ClassProperty>)
func (b *Es6RewriteClass_ClassDeclarationMetadata_Builder) SetClassMembersToDeclare(classMembersToDeclare *jutil.LinkedHashMap[string, *Es6RewriteClass_ClassProperty]) *Es6RewriteClass_ClassDeclarationMetadata_Builder {
	if classMembersToDeclare == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null classMembersToDeclare"))
	}
	b.classMembersToDeclare = classMembersToDeclare
	return b
}

// Java: AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder.setFullClassNameNode(Node)
func (b *Es6RewriteClass_ClassDeclarationMetadata_Builder) SetFullClassNameNode(fullClassNameNode *rhino.Node) *Es6RewriteClass_ClassDeclarationMetadata_Builder {
	if fullClassNameNode == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null fullClassNameNode"))
	}
	b.fullClassNameNode = fullClassNameNode
	return b
}

// Java: AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder.fullClassNameNode()
func (b *Es6RewriteClass_ClassDeclarationMetadata_Builder) FullClassNameNode() *rhino.Node {
	if b.fullClassNameNode == nil {
		panic(jutil.NewIllegalStateExceptionMsg("Property \"fullClassNameNode\" has not been set"))
	}
	return b.fullClassNameNode
}

// Java: AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder.setClassPrototypeNode(Node)
func (b *Es6RewriteClass_ClassDeclarationMetadata_Builder) SetClassPrototypeNode(classPrototypeNode *rhino.Node) *Es6RewriteClass_ClassDeclarationMetadata_Builder {
	if classPrototypeNode == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null classPrototypeNode"))
	}
	b.classPrototypeNode = classPrototypeNode
	return b
}

// Java: AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder.setAnonymous(boolean)
func (b *Es6RewriteClass_ClassDeclarationMetadata_Builder) SetAnonymous(anonymous bool) *Es6RewriteClass_ClassDeclarationMetadata_Builder {
	b.anonymous = anonymous
	b.set_0 |= 1
	return b
}

// Java: AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder.setClassNameNode(Node)
func (b *Es6RewriteClass_ClassDeclarationMetadata_Builder) SetClassNameNode(classNameNode *rhino.Node) *Es6RewriteClass_ClassDeclarationMetadata_Builder {
	if classNameNode == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null classNameNode"))
	}
	b.classNameNode = classNameNode
	return b
}

// Java: AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder.setSuperClassNameNode(Node)
func (b *Es6RewriteClass_ClassDeclarationMetadata_Builder) SetSuperClassNameNode(superClassNameNode *rhino.Node) *Es6RewriteClass_ClassDeclarationMetadata_Builder {
	if superClassNameNode == nil {
		panic(jutil.NewNullPointerExceptionMsg("Null superClassNameNode"))
	}
	b.superClassNameNode = superClassNameNode
	return b
}

// Java: AutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder.build()
func (b *Es6RewriteClass_ClassDeclarationMetadata_Builder) Build() *Es6RewriteClass_ClassDeclarationMetadata {
	if b.set_0 != 1 ||
		b.insertionPoint == nil ||
		b.definePropertiesObjForPrototype == nil ||
		b.definePropertiesObjForClass == nil ||
		b.classMembersToDeclare == nil ||
		b.fullClassNameNode == nil ||
		b.classPrototypeNode == nil ||
		b.classNameNode == nil ||
		b.superClassNameNode == nil {
		missing := jutil.NewStringBuilder()
		if b.insertionPoint == nil {
			missing.Append(" insertionPoint")
		}
		if b.definePropertiesObjForPrototype == nil {
			missing.Append(" definePropertiesObjForPrototype")
		}
		if b.definePropertiesObjForClass == nil {
			missing.Append(" definePropertiesObjForClass")
		}
		if b.classMembersToDeclare == nil {
			missing.Append(" classMembersToDeclare")
		}
		if b.fullClassNameNode == nil {
			missing.Append(" fullClassNameNode")
		}
		if b.classPrototypeNode == nil {
			missing.Append(" classPrototypeNode")
		}
		if (b.set_0 & 1) == 0 {
			missing.Append(" anonymous")
		}
		if b.classNameNode == nil {
			missing.Append(" classNameNode")
		}
		if b.superClassNameNode == nil {
			missing.Append(" superClassNameNode")
		}
		panic(jutil.NewIllegalStateExceptionMsg("Missing required properties:" + missing.String()))
	}
	return NewEs6RewriteClass_ClassDeclarationMetadata(
		b.insertionPoint,
		b.definePropertiesObjForPrototype,
		b.definePropertiesObjForClass,
		b.classMembersToDeclare,
		b.fullClassNameNode,
		b.classPrototypeNode,
		b.anonymous,
		b.classNameNode,
		b.superClassNameNode)
}

// Java: Es6RewriteClass.ClassDeclarationMetadata.builder()
func Es6RewriteClass_ClassDeclarationMetadata_Builder_() *Es6RewriteClass_ClassDeclarationMetadata_Builder {
	return newAutoBuilder_Es6RewriteClass_ClassDeclarationMetadata_Builder().
		SetClassMembersToDeclare(jutil.NewLinkedHashMap[string, *Es6RewriteClass_ClassProperty]())
}

// Java: Es6RewriteClass.ClassDeclarationMetadata.create(Node, Node, AstFactory)
func Es6RewriteClass_ClassDeclarationMetadata_Create(
	classNode *rhino.Node, parent *rhino.Node, astFactory *AstFactory) *Es6RewriteClass_ClassDeclarationMetadata {
	return es6RewriteClass_ClassDeclarationMetadata_CreateInto(&Es6RewriteClass_ClassDeclarationMetadata{}, classNode, parent, astFactory)
}

func es6RewriteClass_ClassDeclarationMetadata_CreateInto(
	m *Es6RewriteClass_ClassDeclarationMetadata,
	classNode *rhino.Node,
	parent *rhino.Node,
	astFactory *AstFactory) *Es6RewriteClass_ClassDeclarationMetadata {
	classNameNode := classNode.GetFirstChild()
	superClassNameNode := classNameNode.GetNext()

	var insertionStart *rhino.Node
	var fullClassNameNode *rhino.Node
	var anonymous bool

	// If this is a class statement, or a class expression in a simple
	// assignment or var statement, convert it. In any other case, the
	// code is too dynamic, so return null.
	if NodeUtil_IsClassDeclaration(classNode) {
		insertionStart = classNode
		fullClassNameNode = classNameNode
		anonymous = false
	} else if parent.IsAssign() && parent.GetParent().IsExprResult() {
		// Add members after the EXPR_RESULT node:
		// example.C = class {}; example.C.prototype.foo = function() {};
		fullClassNameNode = parent.GetFirstChild()
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

	m.insertionPointVal.insertionPoint = insertionStart
	m.insertionPoint = &m.insertionPointVal
	m.astFactory = astFactory
	m.definePropertiesObjForPrototype = nil
	m.definePropertiesObjForClass = nil
	m.classMembersToDeclare = nil
	m.fullClassNameNode = fullClassNameNode
	m.classPrototypeNode = astFactory.CreatePrototypeAccess(fullClassNameNode.CloneTree())
	m.anonymous = anonymous
	m.classNameNode = classNameNode
	m.superClassNameNode = superClassNameNode
	return m
}

// Java: Es6RewriteClass.ClassDeclarationMetadata.insertNodeAndAdvance(Node)
func (m *Es6RewriteClass_ClassDeclarationMetadata) InsertNodeAndAdvance(newNode *rhino.Node) {
	m.InsertionPoint().InsertNodeAndAdvance(newNode)
}

// Java: Es6RewriteClass.ClassDeclarationMetadata.hasSuperClass()
func (m *Es6RewriteClass_ClassDeclarationMetadata) HasSuperClass() bool {
	return !m.SuperClassNameNode().IsEmpty()
}

// Es6RewriteClass_InsertionPoint is used by ClassDeclarationMetadata to
// represent the point where we insert the next transpiled method of a class
type Es6RewriteClass_InsertionPoint struct {
	insertionPoint *rhino.Node
}

// Java: Es6RewriteClass.InsertionPoint(Node)
func newEs6RewriteClass_InsertionPoint(insertionPoint *rhino.Node) *Es6RewriteClass_InsertionPoint {
	return &Es6RewriteClass_InsertionPoint{insertionPoint: insertionPoint}
}

// Java: Es6RewriteClass.InsertionPoint.insertNodeAndAdvance(Node)
func (p *Es6RewriteClass_InsertionPoint) InsertNodeAndAdvance(newNode *rhino.Node) {
	newNode.InsertAfter(p.insertionPoint)
	p.insertionPoint = newNode
}

// Java: Es6RewriteClass.InsertionPoint.from(Node)
func Es6RewriteClass_InsertionPoint_From(start *rhino.Node) *Es6RewriteClass_InsertionPoint {
	return newEs6RewriteClass_InsertionPoint(start)
}

// Java: Es6RewriteClass.InsertionPoint.getNode()
func (p *Es6RewriteClass_InsertionPoint) GetNode() *rhino.Node {
	return p.insertionPoint
}
