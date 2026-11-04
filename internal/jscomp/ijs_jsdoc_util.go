// Port of com/google/javascript/jscomp/ijs/JsdocUtil.java @ closure-compiler v20261028

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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

const ijsSyntheticFileName = "<synthetic>"

var ijsSyntheticFile rhino.StaticSourceFile = rhino.NewSimpleSourceFile(
	ijsSyntheticFileName, rhino.StaticSourceFile_SourceKind_EXTERN)

var ijsSyntheticSrcInfoNode = rhino.NewNode(rhino.Token_EMPTY).SetStaticSourceFile(ijsSyntheticFile)

func isPrivateJSDoc(jsdoc *rhino.JSDocInfo) bool {
	return jsdoc != nil && jsdoc.GetVisibility() == rhino.JSDocInfo_Visibility_PRIVATE
}

func unusableTypeJSDoc(oldJSDoc *rhino.JSDocInfo) *rhino.JSDocInfo {
	return constJSDocWithTypeName(oldJSDoc, "UnusableType")
}

func qmarkTypeJSDoc() *rhino.JSDocInfo {
	return makeJSDocBuilderWithType(nil, rhino.NewNode(rhino.Token_QMARK)).Build()
}

func makeJSDocBuilderWithType(oldJSDoc *rhino.JSDocInfo, typeAst rhino.Node) *rhino.JSDocInfo_Builder {
	builder := rhino.JSDocInfo_Builder_MaybeCopyFrom(oldJSDoc)
	builder.RecordType(newSyntheticTypeExpr(typeAst))
	return builder
}

func replaceJSDocWithBottomType(original *rhino.JSTypeExpression) *rhino.JSTypeExpression {
	if isUnknownTypeExpr(original) {
		// Emitting the UnusableType can cause errors in code if overriding a property
		// or method with a more specific type. There's no harm in preserving the `?`, since it's
		// already loose in the original code.
		return original
	}
	return newSyntheticTypeExpr(rhino.Node_NewString("IjsNoneType"))
}

func replaceJSDocWithUnusableType(original *rhino.JSTypeExpression) *rhino.JSTypeExpression {
	if isUnknownTypeExpr(original) {
		// Emitting the UnusableType can cause errors in code if overriding a property
		// or method with a more specific type. There's no harm in preserving the `?`, since it's
		// already loose in the original code.
		return original
	}
	return newSyntheticTypeExpr(rhino.Node_NewString("UnusableType"))
}

func newSyntheticTypeExpr(typeAst rhino.Node) *rhino.JSTypeExpression {
	return rhino.NewJSTypeExpression(typeAst.SrcrefTree(ijsSyntheticSrcInfoNode), ijsSyntheticFileName)
}

func constJSDocWithTypeName(oldJSDoc *rhino.JSDocInfo, contents string) *rhino.JSDocInfo {
	return constJSDocWithTypeNode(oldJSDoc, rhino.Node_NewString(contents))
}

func constJSDocWithTypeNode(oldJSDoc *rhino.JSDocInfo, typeAst rhino.Node) *rhino.JSDocInfo {
	builder := makeJSDocBuilderWithType(oldJSDoc, typeAst)
	builder.RecordConstancy()
	return builder.Build()
}

func markJSDocConstant(oldJSDoc *rhino.JSDocInfo) *rhino.JSDocInfo {
	builder := rhino.JSDocInfo_Builder_MaybeCopyFrom(oldJSDoc)
	builder.RecordConstancy()
	return builder.Build()
}

func mergeIjsJSDocs(classicJsdoc, inlineJsdoc *rhino.JSDocInfo) *rhino.JSDocInfo {
	if inlineJsdoc == nil || !inlineJsdoc.HasType() {
		return classicJsdoc
	}
	builder := rhino.JSDocInfo_Builder_MaybeCopyFrom(classicJsdoc)
	builder.RecordType(inlineJsdoc.GetType())
	return builder.Build()
}

func replaceAnnotatedTypesWithUnusableType(jsdoc *rhino.JSDocInfo) *rhino.JSDocInfo {
	if jsdoc == nil {
		return nil
	}
	builder := rhino.JSDocInfo_Builder_()
	if jsdoc.HasType() {
		builder.RecordType(replaceJSDocWithUnusableType(jsdoc.GetType()))
	}
	if jsdoc.HasReturnType() {
		builder.RecordReturnType(replaceJSDocWithUnusableType(jsdoc.GetReturnType()))
	}
	for i := 0; i < jsdoc.GetParameterCount(); i++ {
		name, _ := jsdoc.GetParameterNameValueAt(i)
		original := jsdoc.GetParameterType(name)
		if original == nil {
			continue
		}
		newType := replaceJSDocWithBottomType(original)
		if original.IsOptionalArg() {
			newType = rhino.JSTypeExpression_MakeOptionalArg(newType)
		} else if original.IsVarArgs() {
			newType = rhino.JSTypeExpression_MakeVarArgs(newType)
		}
		builder.RecordParameter(name, newType)
	}
	if jsdoc.HasThisType() {
		builder.RecordThisType(replaceJSDocWithUnusableType(jsdoc.GetThisType()))
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
	builder.RecordVisibility(jsdoc.GetVisibility())
	return builder.Build()
}

func hasAnnotatedType(jsdoc *rhino.JSDocInfo) bool {
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

func jsdocForRhs(rhs rhino.Node, oldJSDoc *rhino.JSDocInfo) *rhino.JSDocInfo {
	if rhs.IsNil() {
		return nil
	}
	switch NodeUtil_GetKnownValueType(rhs) {
	case NodeUtil_ValueType_BOOLEAN:
		return constJSDocWithTypeName(oldJSDoc, "boolean")
	case NodeUtil_ValueType_NUMBER:
		return constJSDocWithTypeName(oldJSDoc, "number")
	case NodeUtil_ValueType_BIGINT:
		return constJSDocWithTypeName(oldJSDoc, "bigint")
	case NodeUtil_ValueType_STRING:
		return constJSDocWithTypeName(oldJSDoc, "string")
	case NodeUtil_ValueType_NULL:
		return constJSDocWithTypeName(oldJSDoc, "null")
	case NodeUtil_ValueType_VOID:
		return constJSDocWithTypeName(oldJSDoc, "void")
	case NodeUtil_ValueType_OBJECT:
		if rhs.IsRegExp() {
			return constJSDocWithTypeNode(oldJSDoc, rhino.NewNodeWithChild(rhino.Token_BANG, rhino.IR_String("RegExp")))
		}
	case NodeUtil_ValueType_UNDETERMINED:
		if oldJSDoc != nil {
			if _, ok := oldJSDoc.Description(); ok {
				return constJSDocWithTypeName(oldJSDoc, "string")
			}
		}
	}
	if _, ok := ast.AsCast(rhs); ok {
		if castDoc := rhs.GetJSDocInfo(); castDoc != nil && castDoc.HasType() {
			if root := castDoc.GetType().GetRoot(); !root.IsNil() {
				return constJSDocWithTypeNode(oldJSDoc, root.CloneTree())
			}
		}
	}
	return nil
}

func jsdocForName(decl *Var, oldJSDoc *rhino.JSDocInfo) *rhino.JSDocInfo {
	if decl == nil || decl.GetNameNode().IsNil() {
		return nil
	}
	expr := NodeUtil_GetDeclaredTypeExpression(decl.GetNameNode())
	if expr == nil {
		return nil
	}
	typeAst := expr.GetRoot()
	if typeAst.IsNil() {
		return nil
	}
	switch typeAst.GetToken() {
	case rhino.Token_EQUALS:
		typeRoot := cloneJSDocOperandOrUnknown(typeAst)
		if !decl.IsDefaultParam() {
			typeRoot = rhino.NewNodeWith2Children(rhino.Token_PIPE, typeRoot, rhino.IR_String("undefined"))
		}
		typeAst = typeRoot
	case rhino.Token_ITER_REST:
		newType := rhino.NewNode(rhino.Token_BANG)
		array := rhino.IR_String("Array")
		newType.AddChildToBack(array)
		block := rhino.NewNodeWithChild(rhino.Token_BLOCK, cloneJSDocOperandOrUnknown(typeAst))
		array.AddChildToBack(block)
		typeAst = newType
	default:
		typeAst = typeAst.CloneTree()
	}
	return constJSDocWithTypeNode(oldJSDoc, typeAst)
}

// cloneJSDocOperandOrUnknown clones the inner JSDoc type node of an optional (`=`)
// or rest (`...`) parameter annotation, falling back to `?` for untyped `@param {...}`.
func cloneJSDocOperandOrUnknown(typeAst rhino.Node) rhino.Node {
	for child := range typeAst.Children() {
		if !child.IsEmpty() {
			return child.CloneTree()
		}
		break
	}
	return rhino.NewNode(rhino.Token_QMARK)
}

func isUnknownTypeExpr(typeExpr *rhino.JSTypeExpression) bool {
	if typeExpr == nil {
		return false
	}
	typeAst := typeExpr.GetRoot()
	return !typeAst.IsNil() && typeAst.GetToken() == rhino.Token_QMARK && !typeAst.HasChildren()
}
