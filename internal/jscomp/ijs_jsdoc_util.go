// Port of com/google/javascript/jscomp/ijs/JsdocUtil.java @ closure-compiler v20261021

/*
 * Copyright 2017 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// JsdocUtil holds static utility methods for dealing with inspecting and
// constructing JSDoc objects.
//
// Go-only: language-adapter — the Java class has only a private
// constructor and static members; it has no Go type.

const jsdocUtil_SYNTHETIC_FILE_NAME = "<synthetic>"

var jsdocUtil_SYNTHETIC_FILE rhino.StaticSourceFile = rhino.NewSimpleSourceFile(jsdocUtil_SYNTHETIC_FILE_NAME, rhino.StaticSourceFile_SourceKind_EXTERN)

var jsdocUtil_SYNTETIC_SRCINFO_NODE = rhino.NewNode(rhino.Token_EMPTY).SetStaticSourceFile(jsdocUtil_SYNTHETIC_FILE)

// Java: JsdocUtil.isPrivate(JSDocInfo)
func JsdocUtil_IsPrivate(jsdoc *rhino.JSDocInfo) bool {
	return jsdoc != nil && jsdoc.GetVisibility() == rhino.JSDocInfo_Visibility_PRIVATE
}

// Java: JsdocUtil.getUnusableTypeJSDoc(JSDocInfo)
func JsdocUtil_GetUnusableTypeJSDoc(oldJSDoc *rhino.JSDocInfo) *rhino.JSDocInfo {
	return jsdocUtil_getConstJSDoc(oldJSDoc, "UnusableType")
}

// Java: JsdocUtil.getQmarkTypeJSDoc()
func JsdocUtil_GetQmarkTypeJSDoc() *rhino.JSDocInfo {
	return jsdocUtil_makeBuilderWithType(nil, rhino.NewNode(rhino.Token_QMARK)).Build()
}

// Java: JsdocUtil.makeBuilderWithType(JSDocInfo, Node)
func jsdocUtil_makeBuilderWithType(oldJSDoc *rhino.JSDocInfo, typeAst rhino.Node) *rhino.JSDocInfo_Builder {
	builder := rhino.JSDocInfo_Builder_MaybeCopyFrom(oldJSDoc)
	builder.RecordType(
		rhino.NewJSTypeExpression(typeAst.SrcrefTree(jsdocUtil_SYNTETIC_SRCINFO_NODE), jsdocUtil_SYNTHETIC_FILE_NAME))
	return builder
}

// Java: JsdocUtil.replaceWithBottomType(JSTypeExpression)
func jsdocUtil_replaceWithBottomType(original *rhino.JSTypeExpression) *rhino.JSTypeExpression {
	if jsdocUtil_isUnknownType(original) {
		// Emitting the UnusableType can cause errors in code if overriding a property
		// or method with a more specific type. There's no harm in preserving the `?`, since it's
		// already loose in the original code.
		return original
	}
	return jsdocUtil_newTypeExpr(rhino.Node_NewString("IjsNoneType"))
}

// Java: JsdocUtil.replaceWithUnusableType(JSTypeExpression)
func jsdocUtil_replaceWithUnusableType(original *rhino.JSTypeExpression) *rhino.JSTypeExpression {
	if jsdocUtil_isUnknownType(original) {
		// Emitting the UnusableType can cause errors in code if overriding a property
		// or method with a more specific type. There's no harm in preserving the `?`, since it's
		// already loose in the original code.
		return original
	}
	return jsdocUtil_newTypeExpr(rhino.Node_NewString("UnusableType"))
}

// Java: JsdocUtil.newTypeExpr(Node)
func jsdocUtil_newTypeExpr(typeAst rhino.Node) *rhino.JSTypeExpression {
	return rhino.NewJSTypeExpression(typeAst.SrcrefTree(jsdocUtil_SYNTETIC_SRCINFO_NODE), jsdocUtil_SYNTHETIC_FILE_NAME)
}

// Java: JsdocUtil.getConstJSDoc(JSDocInfo, String)
func jsdocUtil_getConstJSDoc(oldJSDoc *rhino.JSDocInfo, contents string) *rhino.JSDocInfo {
	return jsdocUtil_getConstJSDocOfNode(oldJSDoc, rhino.Node_NewString(contents))
}

// Java: JsdocUtil.getConstJSDoc(JSDocInfo, Node)
func jsdocUtil_getConstJSDocOfNode(oldJSDoc *rhino.JSDocInfo, typeAst rhino.Node) *rhino.JSDocInfo {
	builder := jsdocUtil_makeBuilderWithType(oldJSDoc, typeAst)
	builder.RecordConstancy()
	return builder.Build()
}

// Java: JsdocUtil.markConstant(JSDocInfo)
func JsdocUtil_MarkConstant(oldJSDoc *rhino.JSDocInfo) *rhino.JSDocInfo {
	builder := rhino.JSDocInfo_Builder_MaybeCopyFrom(oldJSDoc)
	builder.RecordConstancy()
	return builder.Build()
}

// Java: JsdocUtil.mergeJsdocs(JSDocInfo, JSDocInfo)
func JsdocUtil_MergeJsdocs(classicJsdoc *rhino.JSDocInfo, inlineJsdoc *rhino.JSDocInfo) *rhino.JSDocInfo {
	if inlineJsdoc == nil || !inlineJsdoc.HasType() {
		return classicJsdoc
	}
	builder := rhino.JSDocInfo_Builder_MaybeCopyFrom(classicJsdoc)
	builder.RecordType(inlineJsdoc.GetType())
	return builder.Build()
}

// Java: JsdocUtil.replaceAnnotatedTypesWithUnusableType(JSDocInfo)
func JsdocUtil_ReplaceAnnotatedTypesWithUnusableType(jsdoc *rhino.JSDocInfo) *rhino.JSDocInfo {
	builder := rhino.JSDocInfo_Builder_()
	if jsdoc.HasType() {
		builder.RecordType(jsdocUtil_replaceWithUnusableType(jsdoc.GetType()))
	}
	if jsdoc.HasReturnType() {
		builder.RecordReturnType(jsdocUtil_replaceWithUnusableType(jsdoc.GetReturnType()))
	}
	for i := 0; i < jsdoc.GetParameterCount(); i++ {
		name, _ := jsdoc.GetParameterNameValueAt(i)
		original := jsdoc.GetParameterType(name)
		if original == nil {
			continue
		}
		newType := jsdocUtil_replaceWithBottomType(original)
		if original.IsOptionalArg() {
			newType = rhino.JSTypeExpression_MakeOptionalArg(newType)
		} else if original.IsVarArgs() {
			newType = rhino.JSTypeExpression_MakeVarArgs(newType)
		}
		builder.RecordParameter(name, newType)
	}
	if jsdoc.HasThisType() {
		builder.RecordThisType(jsdocUtil_replaceWithUnusableType(jsdoc.GetThisType()))
	}
	// NOTE: intentionally *not* replacing @enum, @extends, @implements & @typedef. We have decided
	// that anything *declaring* a new type name should be preserved. For example - given:
	//   class MyClass {}
	//   /** @private @typedef {!Foo|!Bar} */
	//   MyClass.NestedType; */
	//   /** @return {!MyClass.NestedType} */
	//   MyClass.SomeFn = function() { // ...
	// we choose to preserve the definition of `MyClass.NestedType` in order to make
	// `MyClass.SomeFn` well-typed.

	// Preserve existing @const or @private/@package/@protected annotations.
	if jsdoc.IsConstant() {
		builder.RecordConstancy()
	}
	// Go-only: language-adapter — getVisibility() never returns null (it
	// returns INHERITED when unset), so Java's null check always passes.
	builder.RecordVisibility(jsdoc.GetVisibility())
	return builder.Build()
}

// Java: JsdocUtil.hasAnnotatedType(JSDocInfo)
func JsdocUtil_HasAnnotatedType(jsdoc *rhino.JSDocInfo) bool {
	if jsdoc == nil {
		return false
	}
	return jsdoc.HasType() ||
		jsdoc.HasReturnType() ||
		jsdoc.GetParameterCount() > 0 ||
		jsdoc.IsConstructorOrInterface() ||
		jsdoc.HasTypedefType() ||
		jsdoc.HasThisType() ||
		jsdoc.HasEnumParameterType()
}

// Java: JsdocUtil.getJSDocForRhs(Node, JSDocInfo)
func JsdocUtil_GetJSDocForRhs(rhs rhino.Node, oldJSDoc *rhino.JSDocInfo) *rhino.JSDocInfo {
	switch NodeUtil_GetKnownValueType(rhs) {
	case NodeUtil_ValueType_BOOLEAN:
		return jsdocUtil_getConstJSDoc(oldJSDoc, "boolean")
	case NodeUtil_ValueType_NUMBER:
		return jsdocUtil_getConstJSDoc(oldJSDoc, "number")
	case NodeUtil_ValueType_BIGINT:
		return jsdocUtil_getConstJSDoc(oldJSDoc, "bigint")
	case NodeUtil_ValueType_STRING:
		return jsdocUtil_getConstJSDoc(oldJSDoc, "string")
	case NodeUtil_ValueType_NULL:
		return jsdocUtil_getConstJSDoc(oldJSDoc, "null")
	case NodeUtil_ValueType_VOID:
		return jsdocUtil_getConstJSDoc(oldJSDoc, "void")
	case NodeUtil_ValueType_OBJECT:
		if rhs.IsRegExp() {
			return jsdocUtil_getConstJSDocOfNode(oldJSDoc, rhino.NewNodeWithChild(rhino.Token_BANG, rhino.IR_String("RegExp")))
		}
	case NodeUtil_ValueType_UNDETERMINED:
		if oldJSDoc != nil && oldJSDoc.GetDescription() != nil {
			return jsdocUtil_getConstJSDoc(oldJSDoc, "string")
		}
	}
	if rhs.IsCast() {
		return jsdocUtil_getConstJSDocOfNode(oldJSDoc, rhs.GetJSDocInfo().GetType().GetRoot())
	}
	return nil
}

// Java: JsdocUtil.getJSDocForName(Var, JSDocInfo)
func JsdocUtil_GetJSDocForName(decl *Var, oldJSDoc *rhino.JSDocInfo) *rhino.JSDocInfo {
	if decl == nil {
		return nil
	}
	expr := NodeUtil_GetDeclaredTypeExpression(decl.GetNameNode())
	if expr == nil {
		return nil
	}
	typeAst := expr.GetRoot()
	switch typeAst.GetToken() {
	case rhino.Token_EQUALS:
		typeRoot := typeAst.GetFirstChild().CloneTree()
		if !decl.IsDefaultParam() {
			typeRoot = rhino.NewNodeWith2Children(rhino.Token_PIPE, typeRoot, rhino.IR_String("undefined"))
		}
		typeAst = typeRoot
	case rhino.Token_ITER_REST:
		newType := rhino.NewNode(rhino.Token_BANG)
		array := rhino.IR_String("Array")
		newType.AddChildToBack(array)
		block := rhino.NewNodeWithChild(rhino.Token_BLOCK, typeAst.GetFirstChild().CloneTree())
		array.AddChildToBack(block)
		typeAst = newType
	default:
	}
	return jsdocUtil_getConstJSDocOfNode(oldJSDoc, typeAst)
}

// Java: JsdocUtil.isUnknownType(JSTypeExpression)
func jsdocUtil_isUnknownType(type_ *rhino.JSTypeExpression) bool {
	typeAst := type_.GetRoot()
	return typeAst.GetToken() == rhino.Token_QMARK && !typeAst.HasChildren()
}
