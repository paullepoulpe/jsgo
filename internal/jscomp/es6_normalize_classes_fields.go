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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// RewriteSideEffectedComputedProps rewrites and moves all side effected
// computed field keys to the top.
//
// Java: Es6NormalizeClasses.rewriteSideEffectedComputedProps(NodeTraversal, ClassRecord, Node)
func (e *Es6NormalizeClasses) RewriteSideEffectedComputedProps(
	t *NodeTraversal, record *es6NormalizeClasses_ClassRecord, insertionPoint rhino.Node) {
	computedPropsWithSideEffects := &record.computedPropsWithSideEffects

	if len(*computedPropsWithSideEffects) == 0 {
		return
	}

	for len(*computedPropsWithSideEffects) != 0 {
		computedPropMember := es6NormalizeClasses_dequeRemove(computedPropsWithSideEffects)
		e.ExtractExpressionFromCompField(t, record, computedPropMember, insertionPoint)
	}
	t.ReportCodeChange()
}

// Go-only: language-adapter — ArrayDeque.remove(): removes and returns the
// first element, throwing NoSuchElementException when empty.
func es6NormalizeClasses_dequeRemove(deque *[]ast.ClassMember) ast.ClassMember {
	if len(*deque) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	first := (*deque)[0]
	(*deque)[0] = ast.ClassMember{}
	*deque = (*deque)[1:]
	return first
}

// ExtractExpressionFromCompField extracts the expression in the LHS of a
// computed field to not disturb possible side effects and allow for this
// and super to be used in the LHS of a computed field in certain cases.
//
// Java: Es6NormalizeClasses.extractExpressionFromCompField(NodeTraversal, ClassRecord, Node, Node)
func (e *Es6NormalizeClasses) ExtractExpressionFromCompField(
	t *NodeTraversal, record *es6NormalizeClasses_ClassRecord, memberField ast.ClassMember, insertionPoint rhino.Node) {
	compExpression := es6NormalizeClasses_computedKey(memberField).Detach()
	compFieldVar := e.astFactory.
		CreateSingleVarNameDeclarationWithValue(e.GenerateUniqueCompFieldVarName(t), compExpression.Node()).
		SrcrefTreeIfMissing(record.classNode.Node())
	compFieldName := compFieldVar.GetFirstChild()
	memberField.Node().AddChildToFront(compFieldName.CloneNode())
	compFieldVar.InsertBefore(insertionPoint)
	compFieldVar.SrcrefTreeIfMissing(record.classNode.Node())
}

func es6NormalizeClasses_computedKey(member ast.ClassMember) ast.Expr {
	if field, ok := member.AsComputedFieldDef(); ok {
		return field.Key()
	}
	if prop, ok := member.AsComputedProp(); ok {
		return prop.Key()
	}
	panic(jutil.NewIllegalArgumentExceptionMsg(jutil.ToJavaString(member.Node())))
}

// RewriteInstanceMembers rewrites and moves all instance fields into the
// constructor.
//
// Java: Es6NormalizeClasses.rewriteInstanceMembers(NodeTraversal, ClassRecord)
func (e *Es6NormalizeClasses) RewriteInstanceMembers(t *NodeTraversal, record *es6NormalizeClasses_ClassRecord) {
	instanceMembers := &record.instanceMembers
	hasPriv := e.transpilePrivateClassFields && record.hasInstancePrivateMembers

	if len(*instanceMembers) == 0 && !hasPriv {
		return
	}
	classNode := record.classNode.Node()
	e.ctorCreator.SynthesizeClassConstructorIfMissing(t, classNode)
	ctorFn := ast.MustMemberFunctionDef(NodeUtil_GetEs6ClassConstructorMemberFunctionDef(classNode)).Function()
	ctorBlock, ok := ctorFn.Body().AsBlock()
	jutil.CheckState(ok) // a constructor is a method: its body is a block
	_, isDerived := record.classNode.Extends()

	insertionPoint := e.AddTemporaryInsertionPoint(ctorBlock, record)

	if hasPriv {
		mapName := record.GetPrivateMapVarName(e.compiler, t)
		privName := record.GetPrivateVarName(e.compiler, t)

		var privInit rhino.Node
		if record.hasInstancePrivateMethodsOrAccessors {
			// Initializer: Object.create(PRIVATE_PROTO$uid)
			protoName :=
				e.astFactory.CreateNameWithUnknownType(record.GetPrivateProtoVarName(e.compiler, t))
			createProp :=
				e.astFactory.CreateGetPropWithUnknownType(
					e.astFactory.CreateNameWithUnknownType("Object"), "create")
			privInit = e.astFactory.CreateCallWithUnknownType(createProp, protoName)
		} else {
			// Initializer: Object.create(null)
			createProp :=
				e.astFactory.CreateGetPropWithUnknownType(
					e.astFactory.CreateNameWithUnknownType("Object"), "create")
			privInit = e.astFactory.CreateCallWithUnknownType(createProp, e.astFactory.CreateNull())
		}
		// Add `const PRIVATE$uid = privInit;` to the original code.
		privDecl :=
			e.astFactory.
				CreateSingleConstNameDeclaration(privName, privInit).
				SrcrefTreeIfMissing(record.classNode.Node())
		privDecl.InsertBefore(insertionPoint)

		if record.needsInstanceSelf {
			// Add `PRIVATE$uid.$self = this;` to the original code.
			privSelf :=
				e.astFactory.CreateGetPropWithUnknownType(
					e.astFactory.CreateConstantName(privName, AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN)), "$self")
			assignSelf :=
				e.astFactory.
					CreateAssignStatement(privSelf, e.astFactory.CreateThisForEs6Class(record.classNode.Node())).
					SrcrefTreeIfMissing(record.classNode.Node())
			assignSelf.InsertBefore(insertionPoint)
		}

		// Add `PRIVATE_MAP$uid.set(this, PRIVATE$uid);` to the original code.
		mapNameNode := e.astFactory.CreateNameWithUnknownType(mapName)
		mapSetCallee := e.astFactory.CreateGetPropWithUnknownType(mapNameNode, "set")
		mapSetCall :=
			e.astFactory.CreateCallWithUnknownType(
				mapSetCallee,
				e.astFactory.CreateThisForEs6Class(record.classNode.Node()),
				e.astFactory.CreateConstantName(privName, AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN)))
		mapSetStmt := e.astFactory.ExprResult(mapSetCall).SrcrefTreeIfMissing(record.classNode.Node())
		mapSetStmt.InsertBefore(insertionPoint)
	}

	for len(*instanceMembers) != 0 {
		instanceMember := es6NormalizeClasses_dequeRemove(instanceMembers)
		field, isField := instanceMember.AsMemberFieldDef()
		computedField, isComputedField := instanceMember.AsComputedFieldDef()
		jutil.CheckStateMsg1(isField || isComputedField, "%s", instanceMember.Node())

		if isField && field.Node().IsPrivateIdentifier() { // a computed field is never private
			if e.transpilePrivateClassFields {
				// Add `PRIVATE$uid.field = fieldValue;` to the original code.
				assign := e.CreatePrivateFieldAssignment(record.GetPrivateVarName(e.compiler, t), field)
				assign.InsertBefore(insertionPoint)
				field.Detach()
			}
		} else {
			// Add `this.field = fieldValue;` or `this[field] = fieldValue;` to the original code.
			thisNode := e.astFactory.CreateThisForEs6ClassMember(instanceMember.Node())
			var transpiledNode rhino.Node
			if isField {
				transpiledNode = e.ConvNonCompFieldToGetProp(thisNode, field /* isStatic= */, false)
			} else {
				transpiledNode = e.ConvCompFieldToGetElem(thisNode, computedField /* isStatic= */, false)
			}

			if !transpiledNode.IsNil() {
				transpiledNode.InsertBefore(insertionPoint)
			}
		}
	}

	if e.transpilePrivateClassFields {
		// Remove original non-static private member declaration nodes (fields and methods) from the
		// class body AST now that their transpiled equivalents have been generated.
		for _, pm := range record.privateMembers {
			if !pm.isStatic {
				for _, declNode := range pm.declarationNodes {
					if !declNode.GetParent().IsNil() {
						NodeUtil_DeleteNode(declNode, e.compiler)
					}
				}
			}
		}
	}

	// ES2022 § 10.2.2 ([[Construct]]) runs InitializeInstanceElements before
	// OrdinaryCallEvaluateBody, which binds the parameters.
	if !isDerived && NodeUtil_ReferencesEnclosingReceiver(ctorFn.Params().Node()) {
		e.MoveConstructorParamDefaultsAfterFields(t, ctorFn, insertionPoint)
	}

	insertionPoint.Detach()

	// We moved the field from the class body to the constructor, so we need both.
	t.ReportCodeChange()
	t.ReportCodeChangeWithNode(ctorBlock.Node())
}

// CreatePrivateFieldAssignment returns `PRIVATE$uid.field = fieldValue;`
// for the private field `#field = fieldValue`, whose value it moves (or
// `void 0` if it has none), with the field's source position.
func (e *Es6NormalizeClasses) CreatePrivateFieldAssignment(privName string, field ast.MemberFieldDef) rhino.Node {
	propName := jutil.StringSubstringFrom(field.Key(), 1)
	lhs :=
		e.astFactory.CreateGetPropWithUnknownType(
			e.astFactory.CreateConstantName(privName, AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN)), propName)
	var fieldValue rhino.Node
	if value, ok := field.Value(); ok {
		fieldValue = value.Detach().Node()
	} else {
		fieldValue = e.astFactory.CreateUndefinedValue()
	}
	assign := e.astFactory.CreateAssignStatement(lhs, fieldValue)
	assign.SrcrefTreeIfMissing(field.Node())
	return assign
}

// MoveConstructorParamDefaultsAfterFields moves a constructor's parameter
// default initializers and destructuring patterns into statements before
// insertionPoint, so that the instance elements installed before
// insertionPoint are initialized before the defaults run.
//
// Java: Es6NormalizeClasses.moveConstructorParamDefaultsAfterFields(NodeTraversal, Node, Node)
func (e *Es6NormalizeClasses) MoveConstructorParamDefaultsAfterFields(
	t *NodeTraversal, ctorFn ast.Function, insertionPoint rhino.Node) {
	// The loop replaces the current parameter: read the next one first.
	for typedParam := range ast.ReadAhead(ctorFn.Params().FirstParam()) {
		param := typedParam.Node()
		if defaultParam, ok := typedParam.AsDefaultValue(); ok {
			target, value := defaultParam.Target(), defaultParam.Value() // read both before detaching either
			lhs := target.Detach().Node()
			jsDoc := lhs.GetJSDocInfo()
			lhs.SetJSDocInfo(nil)
			defaultValue := value.Detach().Node()
			if _, isName := target.AsName(); isName {
				rhs := e.CreateDefaultValueHook(lhs, defaultValue)
				assignStmt :=
					e.astFactory.CreateAssignStatement(lhs.CloneNode(), rhs).SrcrefTreeIfMissing(param)
				assignStmt.InsertBefore(insertionPoint)
				param.ReplaceWith(lhs)
				lhs.SetJSDocInfo(jsDoc)
			} else {
				_, isPattern := target.AsPattern()
				jutil.CheckStateMsg1(isPattern, "%s", lhs)
				tempParam := e.astFactory.CreateName(e.GenerateUniqueTmpVarName(t), AstFactory_Type_(lhs))
				rhs := e.CreateDefaultValueHook(tempParam, defaultValue)
				letStmt := rhino.IR_LetWithValue(lhs, rhs).SrcrefTreeIfMissing(param)
				letStmt.InsertBefore(insertionPoint)
				NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_LET_DECLARATIONS, e.compiler)
				param.ReplaceWith(tempParam)
				tempParam.SetJSDocInfo(jsDoc)
			}
		} else if paramPattern, ok := es6NormalizeClasses_paramPattern(typedParam); ok {
			pattern := paramPattern.Node()
			jsDoc := pattern.GetJSDocInfo()
			pattern.SetJSDocInfo(nil)
			tempParam := e.astFactory.CreateName(e.GenerateUniqueTmpVarName(t), AstFactory_Type_(pattern))
			pattern.ReplaceWith(tempParam)
			tempParam.SetJSDocInfo(jsDoc)
			letStmt := rhino.IR_LetWithValue(pattern, tempParam.CloneNode()).SrcrefTreeIfMissing(param)
			letStmt.InsertBefore(insertionPoint)
			NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_LET_DECLARATIONS, e.compiler)
		}
	}
}

func es6NormalizeClasses_paramPattern(param ast.Param) (ast.Pattern, bool) {
	if rest, ok := param.AsIterRest(); ok {
		return rest.Target().AsPattern()
	}
	return param.AsPattern()
}

// Java: Es6NormalizeClasses.createDefaultValueHook(Node, Node)
func (e *Es6NormalizeClasses) CreateDefaultValueHook(paramName rhino.Node, defaultValue rhino.Node) rhino.Node {
	undefined := e.astFactory.CreateUndefinedValue()
	undefined.MakeNonIndexable()
	return e.astFactory.CreateHook(
		e.astFactory.CreateSheq(paramName.CloneNode(), undefined),
		defaultValue,
		paramName.CloneNode())
}

// AddTemporaryInsertionPoint finds the location of super() call in the
// constructor and adds a temporary empty node after it (or after all
// hoisted inner function declarations).
//
// Java: Es6NormalizeClasses.addTemporaryInsertionPoint(Node, ClassRecord)
func (e *Es6NormalizeClasses) AddTemporaryInsertionPoint(ctorBlock ast.Block, classRecord *es6NormalizeClasses_ClassRecord) rhino.Node {
	jutil.CheckStateMsg1(
		len(classRecord.superCalls) <= 1,
		"classes with public fields must have only one super() call at the constructor root: %s",
		classRecord.classNode.Node())

	tempNode := rhino.IR_Empty()
	if len(classRecord.superCalls) == 1 {
		stmt := e.EnsureSuperCallIsOwnStatement(classRecord.superCalls[0])
		tempNode.InsertAfter(stmt)
		return tempNode
	}
	// Search for super() calls added by SynthesizeExplicitConstructors.
	for stmt := range ctorBlock.Statements() {
		if es6NormalizeClasses_isSuperCallStatement(stmt) {
			tempNode.InsertAfter(stmt.Node())
			return tempNode
		}
	}
	insertionPoint := NodeUtil_GetInsertionPointAfterAllInnerFunctionDeclarations(ctorBlock.Node())
	if !insertionPoint.IsNil() {
		tempNode.InsertBefore(insertionPoint)
	} else {
		ctorBlock.Node().AddChildToBack(tempNode)
	}
	return tempNode
}

func es6NormalizeClasses_isSuperCallStatement(stmt ast.Statement) bool {
	exprStmt, ok := stmt.AsExprStmt()
	if !ok {
		return false
	}
	call, ok := exprStmt.Expr().AsCall()
	if !ok {
		return false
	}
	_, ok = call.Callee().AsSuper()
	return ok
}

// EnsureSuperCallIsOwnStatement returns a statement containing either
// `super(...);` or `var tmp = super(...)`, extracting the given super call
// into a `tmp` variable if necessary.
//
// Java: Es6NormalizeClasses.ensureSuperCallIsOwnStatement(Node)
func (e *Es6NormalizeClasses) EnsureSuperCallIsOwnStatement(superCall rhino.Node) rhino.Node {
	parent := superCall.GetParent()
	if parent.IsExprResult() {
		return parent // No extraction needed.
	}
	jutil.CheckStateMsg1(
		e.expressionDecomposer.CanExposeExpression(superCall) != ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE,
		"Cannot decompose super() call in a class with class fields. Move super() call to the"+
			" root of the constructor. %s",
		superCall)
	e.expressionDecomposer.MaybeExposeExpression(superCall)
	e.expressionDecomposer.MoveExpression(superCall)
	return NodeUtil_GetEnclosingStatement(superCall)
}

// Java: Es6NormalizeClasses.createObjectDotDefineProperty()
func (e *Es6NormalizeClasses) CreateObjectDotDefineProperty() rhino.Node {
	e.compiler.GetRuntimeJsLibManager().InjectLibForField("$jscomp.global")
	global := e.compiler.GetRuntimeJsLibManager().GetJsLibField("$jscomp.global")
	jscompDotGlobal := e.astFactory.CreateQNameWithUnknownTypeOfField(global)
	objectNode := e.astFactory.CreateGetPropWithUnknownType(jscompDotGlobal, "Object")
	return e.astFactory.CreateGetPropWithUnknownType(objectNode, "defineProperty")
}

// ConvNonCompFieldToGetProp creates a node that represents receiver.key =
// value; where the key and value come from the non-computed field, or an
// Object.defineProperty call if the property is "__proto__".
//
// Java: Es6NormalizeClasses.convNonCompFieldToGetProp(Node, Node, boolean)
func (e *Es6NormalizeClasses) ConvNonCompFieldToGetProp(
	receiver rhino.Node, noncomputedField ast.MemberFieldDef, isStatic bool) rhino.Node {
	jutil.CheckArgumentMsg1(receiver.GetParent().IsNil(), "%s", receiver)

	var fieldValue rhino.Node
	if value, ok := noncomputedField.Value(); ok && (isStatic || e.transpileClassFields) {
		fieldValue = value.Node()
	}

	field := noncomputedField.Node()
	if field.GetStringID() == es6ClassProtoID {
		return e.CreateFieldDefineProperty(
			field, fieldValue, receiver, e.astFactory.CreateString("__proto__"))
	}

	val, ok := e.prepareFieldInitValue(field, fieldValue)
	if !ok {
		return rhino.NilNode
	}
	lhs := e.astFactory.CreateGetProp(receiver, noncomputedField.Key(), AstFactory_Type_(field))
	return e.finishFieldInitStmt(field, e.astFactory.CreateAssignStatement(lhs, val))
}

// ConvCompFieldToGetElem creates a node that represents receiver[key] =
// value; where the key and value come from the computed field, or an
// Object.defineProperty call if the property evaluates to "__proto__".
//
// Java: Es6NormalizeClasses.convCompFieldToGetElem(Node, Node, boolean)
func (e *Es6NormalizeClasses) ConvCompFieldToGetElem(
	receiver rhino.Node, computedField ast.ComputedFieldDef, isStatic bool) rhino.Node {
	jutil.CheckArgumentMsg1(receiver.GetParent().IsNil(), "%s", receiver)

	// Read the key and the value before detaching the key.
	key := computedField.Key()
	var fieldValue rhino.Node
	if value, ok := computedField.Value(); ok && (isStatic || e.transpileClassFields) {
		fieldValue = value.Node()
	}

	var compFieldNameExpr rhino.Node
	if e.transpileClassFields {
		compFieldNameExpr = key.Detach().Node()
	} else {
		compFieldNameExpr = key.Node().CloneTree()
	}

	field := computedField.Node()
	if es6IsProtoKeyExpr(compFieldNameExpr) {
		return e.CreateFieldDefineProperty(field, fieldValue, receiver, compFieldNameExpr)
	}

	val, ok := e.prepareFieldInitValue(field, fieldValue)
	if !ok {
		return rhino.NilNode
	}
	lhs := e.astFactory.CreateGetElem(receiver, compFieldNameExpr)
	return e.finishFieldInitStmt(field, e.astFactory.CreateAssignStatement(lhs, val))
}

func es6IsProtoKeyExpr(expr rhino.Node) bool {
	if expr.IsStringLit() {
		return expr.GetStringID() == es6ClassProtoID
	}
	v := NodeUtil_GetStringValue(expr)
	return v != nil && *v == "__proto__"
}

// Java: Es6NormalizeClasses.createFieldDefineProperty(Node, Node, Node, Node)
func (e *Es6NormalizeClasses) CreateFieldDefineProperty(
	field rhino.Node, fieldValue rhino.Node, receiver rhino.Node, keyNode rhino.Node) rhino.Node {
	val, ok := e.prepareFieldInitValue(field, fieldValue)
	if !ok {
		return rhino.NilNode
	}
	descriptor :=
		e.astFactory.CreateObjectLit(
			e.astFactory.CreateStringKey("configurable", e.astFactory.CreateBoolean(true)),
			e.astFactory.CreateStringKey("enumerable", e.astFactory.CreateBoolean(true)),
			e.astFactory.CreateStringKey("writable", e.astFactory.CreateBoolean(true)),
			e.astFactory.CreateStringKey("value", val))
	definePropertyCall :=
		e.astFactory.CreateCallWithUnknownType(
			e.CreateObjectDotDefineProperty(), receiver, keyNode, descriptor)
	return e.finishFieldInitStmt(field, e.astFactory.ExprResult(definePropertyCall))
}

// Java: Es6NormalizeClasses.createFieldAssignment(Node, Node, Supplier<Node>)
func (e *Es6NormalizeClasses) CreateFieldAssignment(
	field rhino.Node, fieldValue rhino.Node, createAssignLhs func() rhino.Node) rhino.Node {
	val, ok := e.prepareFieldInitValue(field, fieldValue)
	if !ok {
		return rhino.NilNode
	}
	return e.finishFieldInitStmt(field, e.astFactory.CreateAssignStatement(createAssignLhs(), val))
}

// Java: Es6NormalizeClasses.createFieldInitializationStatement(Node, Node, Function<Node, Node>)
func (e *Es6NormalizeClasses) CreateFieldInitializationStatement(
	field rhino.Node, fieldValue rhino.Node, createStatement func(rhino.Node) rhino.Node) rhino.Node {
	val, ok := e.prepareFieldInitValue(field, fieldValue)
	if !ok {
		return rhino.NilNode
	}
	return e.finishFieldInitStmt(field, createStatement(val))
}

func (e *Es6NormalizeClasses) prepareFieldInitValue(field rhino.Node, fieldValue rhino.Node) (rhino.Node, bool) {
	if e.transpileClassFields && !field.IsPrivateIdentifier() {
		field.Detach()
	}
	if fieldValue.IsNil() {
		if !e.transpileClassFields {
			return rhino.NilNode, false
		}
		return e.astFactory.CreateUndefinedValue(), true
	}
	fieldValue.Detach()
	return fieldValue, true
}

func (e *Es6NormalizeClasses) finishFieldInitStmt(field rhino.Node, result rhino.Node) rhino.Node {
	if e.transpileClassFields {
		result.GetFirstChild().SetJSDocInfo(field.GetJSDocInfo())
	}
	result.SrcrefTreeIfMissing(field)
	return result
}

func es6StaticInitNeedsMethodPred(n rhino.Node) bool {
	return n.IsSuper() || n.IsPrivateIdentifier()
}

// RewriteStaticMembers rewrites and moves all static blocks and fields.
//
// Java: Es6NormalizeClasses.rewriteStaticMembers(NodeTraversal, ClassRecord, Node)
func (e *Es6NormalizeClasses) RewriteStaticMembers(t *NodeTraversal, record *es6NormalizeClasses_ClassRecord, insertionPoint rhino.Node) {
	staticMembers := &record.staticMembers
	hasPriv := e.transpilePrivateClassFields && record.hasStaticPrivateMembers
	if len(*staticMembers) == 0 && !hasPriv {
		return
	}

	staticInitTempParent := rhino.IR_Block()

	if hasPriv {
		staticMapName := record.GetStaticPrivateMapVarName(e.compiler, t)
		privName := record.GetStaticPrivateVarName(e.compiler, t)

		for _, pm := range record.privateMembers {
			if pm.isStatic && pm.kind == es6NormalizeClasses_ClassRecord_PrivateMemberKind_ACCESSOR {
				if getFn, ok := es6NormalizeClasses_unmovedFunction(pm.getterNode.Node()); ok && e.HasThisReference(getFn) {
					record.needsStaticSelf = true
				}
				if setFn, ok := es6NormalizeClasses_unmovedFunction(pm.setterNode.Node()); ok && e.HasThisReference(setFn) {
					record.needsStaticSelf = true
				}
			}
		}

		createProp :=
			e.astFactory.CreateGetPropWithUnknownType(
				e.astFactory.CreateNameWithUnknownType("Object"), "create")
		staticPrivInit :=
			e.astFactory.CreateCallWithUnknownType(createProp, e.astFactory.CreateNull())
		privDecl :=
			e.astFactory.
				CreateSingleConstNameDeclaration(privName, staticPrivInit).
				SrcrefTreeIfMissing(record.classNode.Node())
		staticInitTempParent.AddChildToBack(privDecl)

		if record.needsStaticSelf {
			privSelf :=
				e.astFactory.CreateGetPropWithUnknownType(
					e.astFactory.CreateConstantName(privName, AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN)), "$self")
			assignSelf :=
				e.astFactory.
					CreateAssignStatement(
						privSelf,
						record.CreateNewNameReferenceNode().SrcrefTreeIfMissing(record.classNode.Node())).
					SrcrefTreeIfMissing(record.classNode.Node())
			staticInitTempParent.AddChildToBack(assignSelf)
		}

		mapNameNode := e.astFactory.CreateNameWithUnknownType(staticMapName)
		mapSetCallee := e.astFactory.CreateGetPropWithUnknownType(mapNameNode, "set")
		mapSetCall :=
			e.astFactory.CreateCallWithUnknownType(
				mapSetCallee,
				record.CreateNewNameReferenceNode().SrcrefTreeIfMissing(record.classNode.Node()),
				e.astFactory.CreateConstantName(privName, AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN)))
		mapSetStmt := e.astFactory.ExprResult(mapSetCall).SrcrefTreeIfMissing(record.classNode.Node())
		staticInitTempParent.AddChildToBack(mapSetStmt)

		e.AddStaticPrivateMembers(record, privName, staticInitTempParent)
	}

	for len(*staticMembers) != 0 {
		staticMember := es6NormalizeClasses_dequeRemove(staticMembers)
		if staticMember.Node().IsPrivateIdentifier() {
			if e.transpilePrivateClassFields {
				field := ast.MustMemberFieldDef(staticMember.Node())
				assign := e.CreatePrivateFieldAssignment(record.GetStaticPrivateVarName(e.compiler, t), field)
				staticInitTempParent.AddChildToBack(assign)
				field.Detach()
			}
		} else {
			nameToUse := record.CreateNewNameReferenceNode().SrcrefTree(staticMember.Node())

			var transpiledNode rhino.Node
			switch staticMember.Token() {
			case rhino.Token_BLOCK:
				block := ast.MustBlock(staticMember.Node()).Detach()
				t.ReportCodeChange()
				if block.HasStatements() {
					transpiledNode = block.Node()
				}
			case rhino.Token_MEMBER_FIELD_DEF:
				transpiledNode = e.ConvNonCompFieldToGetProp(nameToUse, ast.MustMemberFieldDef(staticMember.Node()) /* isStatic= */, true)
			case rhino.Token_COMPUTED_FIELD_DEF:
				transpiledNode = e.ConvCompFieldToGetElem(nameToUse, ast.MustComputedFieldDef(staticMember.Node()) /* isStatic= */, true)
			default:
				panic(jutil.NewIllegalStateExceptionMsg(jutil.ToJavaString(staticMember.Node())))
			}
			if !transpiledNode.IsNil() {
				staticInitTempParent.AddChildToBack(transpiledNode)
			}
		}
	}

	if e.transpilePrivateClassFields {
		for _, pm := range record.privateMembers {
			if pm.isStatic {
				for _, declNode := range pm.declarationNodes {
					if !declNode.GetParent().IsNil() {
						NodeUtil_DeleteNode(declNode, e.compiler)
					}
				}
			}
		}
	}

	if !staticInitTempParent.HasChildren() {
		return
	}

	if !hasPriv &&
		!NodeUtil_Has(
			staticInitTempParent,
			es6StaticInitNeedsMethodPred,
			es6NonClassOrNonArrowFnPred) {
		for staticInitTempParent.HasChildren() {
			staticInitTempParent.GetLastChild().Detach().InsertAfter(insertionPoint)
		}
		t.ReportCodeChangeWithNode(insertionPoint.GetParent())
		return
	}

	staticInitMethodName := e.GenerateUniqueStaticInitMethodName(t)
	staticInitMethodFunc :=
		e.astFactory.CreateEmptyFunction(AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_FUNCTION_TYPE, colors.StandardColors_TOP_OBJECT))
	staticInitMethod :=
		e.astFactory.CreateMemberFunctionDef(staticInitMethodName, staticInitMethodFunc)
	staticInitMethod.SetStaticMember(true)

	staticInitBlock := ast.MustFunction(staticInitMethodFunc).Body().Node()
	staticInitBlock.AddChildrenToFront(staticInitTempParent.RemoveChildren())

	classMembers := record.classNode.Body().Node()
	staticInitMethod.SrcrefTreeIfMissing(classMembers)
	classMembers.AddChildToBack(staticInitMethod)
	e.compiler.ReportChangeToChangeScope(staticInitMethodFunc)
	t.ReportCodeChangeWithNode(staticInitMethod)
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_MEMBER_DECLARATIONS, e.compiler)

	staticInitCall := e.astFactory.
		CreateCall(
			e.astFactory.CreateGetProp(
				record.CreateNewNameReferenceNode(),
				staticInitMethodName,
				AstFactory_Type_(staticInitMethodFunc)),
			AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_VOID_TYPE, colors.StandardColors_NULL_OR_VOID)).
		SrcrefTreeIfMissing(classMembers)
	staticInitCallStmt := e.astFactory.ExprResult(staticInitCall)
	staticInitCallStmt.InsertAfter(insertionPoint)
}
