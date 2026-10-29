// Port of com/google/javascript/jscomp/Es6NormalizeClasses.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/Es6NormalizeShorthandProperties.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/Es6RenameReferences.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/SynthesizeExplicitConstructors.java @ closure-compiler v20261021

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

/*
 * Copyright 2015 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: this pass uses the typed AST view (internal/ast) where it reads
// class shapes: the class (cls.Name(), cls.Extends(), cls.Body()), its
// members (ast.MemberFieldDef, ast.ComputedFieldDef, ast.MemberFunctionDef,
// ast.GetterDef, ast.SetterDef, static blocks), the constructor's function,
// parameters and body, the `in`, property access, call and assignment
// around a private name or `super`, and the declaration or assignment that
// holds a class. Private-member lowering, field/static-block lowering, and
// per-class record state live in es6_normalize_classes_private.go,
// es6_normalize_classes_fields.go, and es6_normalize_classes_records.go.

// Es6NormalizeClasses normalizes ES6 classes to a simpler form for easier
// optimization.
//
// This pass handles four main transformations which are performed when
// needed:
//
//  1. Extracts class extends expressions into variables.
//  2. Extracts classes into variables.
//  3. Rewrites static class field initializers and static blocks into a
//     static initialization method that is called after class definition.
//  4. If transpilation of Feature#PUBLIC_CLASS_FIELDS is needed, rewrites
//     public class instance fields into assignments in the constructor and
//     removes static field declarations.
//
// These transformations must be done before passes like Es6RewriteClass and
// Es6ConvertSuper, which only handle classes in simpler forms.
type Es6NormalizeClasses struct {
	compiler                    AbstractCompiler
	astFactory                  *AstFactory
	expressionDecomposer        *ExpressionDecomposer
	ctorCreator                 *SynthesizeExplicitConstructors
	transpileClassFields        bool
	transpilePrivateClassFields bool

	doStaticInheritanceRewrites bool
	doStaticSuperRewrites       bool
	jscompPrivateMap            RuntimeJsLibManager_JsLibField

	// Go-only: language-adapter — the Java ArrayDeque used as a stack
	// (addFirst/peek/removeFirst, iterated from the first element): the
	// first element is the last one here.
	classStack    []*es6NormalizeClasses_ClassRecord
	classStackBuf [4]*es6NormalizeClasses_ClassRecord
}

// Java: Es6NormalizeClasses.CLASS_DECL_VAR
const es6NormalizeClasses_CLASS_DECL_VAR = "$jscomp$classDecl$"

// Java: Es6NormalizeClasses.CLASS_EXTENDS_VAR
const es6NormalizeClasses_CLASS_EXTENDS_VAR = "$jscomp$classExtends$"

// Java: Es6NormalizeClasses.COMP_FIELD_VAR
const es6NormalizeClasses_COMP_FIELD_VAR = "$jscomp$compField$"

// Java: Es6NormalizeClasses.STATIC_INIT_METHOD_NAME
const es6NormalizeClasses_STATIC_INIT_METHOD_NAME = "$jscomp$staticInit$"

// Java: Es6NormalizeClasses.PRIVATE_MAP_VAR
const es6NormalizeClasses_PRIVATE_MAP_VAR = "$jscomp$privateMap$"

// Java: Es6NormalizeClasses.STATIC_PRIVATE_MAP_VAR
const es6NormalizeClasses_STATIC_PRIVATE_MAP_VAR = "$jscomp$staticPrivateMap$"

// Java: Es6NormalizeClasses.PRIVATE_PROTO_VAR
const es6NormalizeClasses_PRIVATE_PROTO_VAR = "$jscomp$priv$proto$"

// Java: Es6NormalizeClasses.PRIVATE_VAR
const es6NormalizeClasses_PRIVATE_VAR = "$jscomp$priv$"

// Java: Es6NormalizeClasses.TMP_VAR
const es6NormalizeClasses_TMP_VAR = "$jscomp$tmp$"

// Java: Es6NormalizeClasses.ILLEGAL_PRIVATE_MEMBER_ASSIGNMENT
var Es6NormalizeClasses_ILLEGAL_PRIVATE_MEMBER_ASSIGNMENT = DiagnosticType_Error(
	"JSC_ILLEGAL_PRIVATE_MEMBER_ASSIGNMENT",
	"Invalid assignment to private method ''#{0}''. Private methods cannot be reassigned.")

// Es6NormalizeClasses_GENERIC_NAME_REPLACEMENTS is the mapping from generic
// placeholder names used in unit tests to the variable name prefixes
// generated by class normalization and rewriting passes.
//
// Used in unit tests via setGenericNameReplacements so that test expected
// output can use readable placeholder names (such as CLASS_DECL or PRIVATE)
// rather than hardcoding compiler-generated variable names.
//
// Java: Es6NormalizeClasses.GENERIC_NAME_REPLACEMENTS
var Es6NormalizeClasses_GENERIC_NAME_REPLACEMENTS = jutil.NewImmutableMapBuilder[string, string]().
	Put("CLASS_DECL", es6NormalizeClasses_CLASS_DECL_VAR).
	Put("CLASS_EXTENDS", es6NormalizeClasses_CLASS_EXTENDS_VAR).
	Put("COMP_FIELD", es6NormalizeClasses_COMP_FIELD_VAR).
	Put("STATIC_INIT", es6NormalizeClasses_STATIC_INIT_METHOD_NAME).
	Put("STATIC_PRIVATE_MAP", es6NormalizeClasses_STATIC_PRIVATE_MAP_VAR).
	Put("PRIVATE_MAP", es6NormalizeClasses_PRIVATE_MAP_VAR).
	Put("PRIVATE_PROTO", es6NormalizeClasses_PRIVATE_PROTO_VAR).
	Put("PRIVATE", es6NormalizeClasses_PRIVATE_VAR).
	Put("TMP", es6NormalizeClasses_TMP_VAR).
	BuildOrThrow()

// GenerateUniqueTmpVarName returns $jscomp$tmp$[FILE_ID]$[number]
//
// Java: Es6NormalizeClasses.generateUniqueTmpVarName(NodeTraversal)
func (e *Es6NormalizeClasses) GenerateUniqueTmpVarName(t *NodeTraversal) string {
	return es6NormalizeClasses_TMP_VAR + e.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())
}

// GenerateUniqueClassDeclVarName returns $jscomp$classDecl$[FILE_ID]$[number]
//
// Java: Es6NormalizeClasses.generateUniqueClassDeclVarName(NodeTraversal)
func (e *Es6NormalizeClasses) GenerateUniqueClassDeclVarName(t *NodeTraversal) string {
	return es6NormalizeClasses_CLASS_DECL_VAR + e.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())
}

// GenerateUniqueClassExtendsVarName returns
// $jscomp$classExtends$[FILE_ID]$[number]
//
// Java: Es6NormalizeClasses.generateUniqueClassExtendsVarName(NodeTraversal)
func (e *Es6NormalizeClasses) GenerateUniqueClassExtendsVarName(t *NodeTraversal) string {
	return es6NormalizeClasses_CLASS_EXTENDS_VAR + e.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())
}

// GenerateUniqueCompFieldVarName returns $jscomp$compfield$[FILE_ID]$[number]
//
// Java: Es6NormalizeClasses.generateUniqueCompFieldVarName(NodeTraversal)
func (e *Es6NormalizeClasses) GenerateUniqueCompFieldVarName(t *NodeTraversal) string {
	return es6NormalizeClasses_COMP_FIELD_VAR + e.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())
}

// GenerateUniqueStaticInitMethodName returns
// $jscomp$staticInit$[FILE_ID]$[number]
//
// Java: Es6NormalizeClasses.generateUniqueStaticInitMethodName(NodeTraversal)
func (e *Es6NormalizeClasses) GenerateUniqueStaticInitMethodName(t *NodeTraversal) string {
	return es6NormalizeClasses_STATIC_INIT_METHOD_NAME + e.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())
}

// Es6NormalizeClasses_GenerateUniquePrivateMapVarName returns
// $jscomp$privateMap$[FILE_ID]$[number]
//
// Java: Es6NormalizeClasses.generateUniquePrivateMapVarName(AbstractCompiler, NodeTraversal)
func Es6NormalizeClasses_GenerateUniquePrivateMapVarName(
	compiler AbstractCompiler, t *NodeTraversal) string {
	return es6NormalizeClasses_PRIVATE_MAP_VAR + compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())
}

// Es6NormalizeClasses_GenerateUniqueStaticPrivateMapVarName returns
// $jscomp$staticPrivateMap$[FILE_ID]$[number]
//
// Java: Es6NormalizeClasses.generateUniqueStaticPrivateMapVarName(AbstractCompiler, NodeTraversal)
func Es6NormalizeClasses_GenerateUniqueStaticPrivateMapVarName(
	compiler AbstractCompiler, t *NodeTraversal) string {
	return es6NormalizeClasses_STATIC_PRIVATE_MAP_VAR + compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())
}

// Es6NormalizeClasses_GenerateUniquePrivateProtoVarName returns
// $jscomp$priv$proto$[FILE_ID]$[number]
//
// Java: Es6NormalizeClasses.generateUniquePrivateProtoVarName(AbstractCompiler, NodeTraversal)
func Es6NormalizeClasses_GenerateUniquePrivateProtoVarName(
	compiler AbstractCompiler, t *NodeTraversal) string {
	return es6NormalizeClasses_PRIVATE_PROTO_VAR + compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())
}

// Es6NormalizeClasses_GenerateUniquePrivateVarName returns
// $jscomp$priv$[FILE_ID]$[number]
//
// Java: Es6NormalizeClasses.generateUniquePrivateVarName(AbstractCompiler, NodeTraversal)
func Es6NormalizeClasses_GenerateUniquePrivateVarName(compiler AbstractCompiler, t *NodeTraversal) string {
	return es6NormalizeClasses_PRIVATE_VAR + compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())
}

// Java: Es6NormalizeClasses(AbstractCompiler)
func NewEs6NormalizeClasses(compiler AbstractCompiler) *Es6NormalizeClasses {
	e := &Es6NormalizeClasses{}
	e.compiler = compiler
	e.astFactory = compiler.CreateAstFactory()
	e.expressionDecomposer = compiler.CreateDefaultExpressionDecomposer()
	e.ctorCreator = NewSynthesizeExplicitConstructors(compiler)
	// Note: This covers both public fields as well as static fields (they launched together).
	options := compiler.GetOptions()
	e.transpileClassFields = options.NeedsTranspilationOf(parser.FeatureSet_Feature_PUBLIC_CLASS_FIELDS)
	e.transpilePrivateClassFields = options.NeedsTranspilationOf(parser.FeatureSet_Feature_PRIVATE_ELEMENTS)

	e.doStaticInheritanceRewrites =
		options.GetAssumeStaticInheritanceIsNotUsed() && !options.GetSkipNonTranspilationPasses()

	// Es6ConvertSuper handles static super rewrites otherwise.
	e.doStaticSuperRewrites = !options.NeedsTranspilationOf(parser.FeatureSet_Feature_SUPER)
	e.jscompPrivateMap = compiler.GetRuntimeJsLibManager().GetJsLibField("$jscomp.PrivateMap")
	e.classStack = e.classStackBuf[:0]
	return e
}

// Java: Es6NormalizeClasses.process(Node, Node)
func (e *Es6NormalizeClasses) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_TraverseRoots(e.compiler, e, externs, root)

	// We transpiled away class static blocks. We always do this, regardless of language out, to
	// normalize the AST.
	transpiledFeatures := parser.FeatureSet_BARE_MINIMUM.With(parser.FeatureSet_Feature_CLASS_STATIC_BLOCK)
	if e.transpileClassFields {
		// We removed public and static field declarations if they need transpilation.
		transpiledFeatures = transpiledFeatures.With(parser.FeatureSet_Feature_PUBLIC_CLASS_FIELDS)
	}
	if e.transpilePrivateClassFields {
		transpiledFeatures = transpiledFeatures.With(parser.FeatureSet_Feature_PRIVATE_ELEMENTS)
	}
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(e.compiler, root, transpiledFeatures)
	// Update the compiler's summary of getter and setter property names in the AST since class
	// getter/setter definitions may have been detached or rewritten into normal
	// methods/descriptors.
	GatherGetterAndSetterProperties_Update(e.compiler, externs, root)
}

// Go-only: language-adapter — Deque.peek(): the first element, or nil.
func (e *Es6NormalizeClasses) peekClassRecord() *es6NormalizeClasses_ClassRecord {
	if len(e.classStack) == 0 {
		return nil
	}
	return e.classStack[len(e.classStack)-1]
}

// Java: Es6NormalizeClasses.shouldTraverse(NodeTraversal, Node, Node)
func (e *Es6NormalizeClasses) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		scriptFeatures := NodeUtil_GetFeatureSetOfScript(n)
		return scriptFeatures == nil ||
			scriptFeatures.ContainsFeature(parser.FeatureSet_Feature_CLASSES) ||
			scriptFeatures.ContainsFeature(parser.FeatureSet_Feature_CLASS_STATIC_BLOCK) ||
			scriptFeatures.ContainsFeature(parser.FeatureSet_Feature_PUBLIC_CLASS_FIELDS) ||
			scriptFeatures.ContainsFeature(parser.FeatureSet_Feature_PRIVATE_ELEMENTS)
	case rhino.Token_CLASS:
		cls := ast.MustClass(n)
		// The class name for declarations and class expression inner names.
		bindingIdentifier, hasBindingIdentifier := cls.Name()

		// Generate any needed normalized class names in pre-order so that we can rewrite the inner
		// class name to the new outer name if needed.
		// (We don't actually do the normalization in pre-order because that can cause
		// NodeTraversal to skip traversing other parts of the AST.)
		normalized := e.GenerateNormalizedClassNames(t, cls)

		record := newEs6NormalizeClasses_ClassRecord(normalized, bindingIdentifier, hasBindingIdentifier, cls)
		// Pre-scan all class members (fields, methods, accessors) upon entering the CLASS node
		// to record all private identifiers (#member) declared within the class body.
		// Pre-registering private declarations up front ensures that forward brand checks
		// (e.g., `#field in this` inside a method defined above `#field`) and cross-member
		// references can resolve the private member record regardless of AST traversal order.
		//
		// Go-only: Java first checks that the last child of the class is its body; the grammar
		// says it always is.
		for member := range cls.Body().Members() {
			if Es6NormalizeClasses_IsPrivateMember(member.Node()) {
				record.RecordPrivateMember(member)
			}
		}
		e.classStack = append(e.classStack, record)
	case rhino.Token_COMPUTED_FIELD_DEF:
		jutil.CheckState(len(e.classStack) != 0)
		field := ast.MustComputedFieldDef(n).AsClassMember()
		classRecord := e.peekClassRecord()
		classRecord.MaybeRecordComputedPropWithSideEffects(field)
		classRecord.RecordField(field)
	case rhino.Token_MEMBER_FIELD_DEF:
		jutil.CheckState(len(e.classStack) != 0)
		field := ast.MustMemberFieldDef(n).AsClassMember()
		classRecord := e.peekClassRecord()
		if n.IsPrivateIdentifier() {
			classRecord.RecordPrivateMember(field)
		}
		classRecord.RecordField(field)
	case rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF:
		if n.IsPrivateIdentifier() {
			jutil.CheckState(len(e.classStack) != 0)
			// A private method or accessor is a class member: object literals have no private names.
			e.peekClassRecord().RecordPrivateMember(ast.MustClassMember(n))
		}
	case rhino.Token_BLOCK:
		if !NodeUtil_IsClassStaticBlock(n) {
			break
		}

		jutil.CheckState(len(e.classStack) != 0)
		e.peekClassRecord().RecordStaticBlock(ast.MustBlock(n))
	// For example, a class method with a computed property name.
	case rhino.Token_COMPUTED_PROP:
		if !parent.IsClassMembers() {
			break
		}

		jutil.CheckState(len(e.classStack) != 0)
		e.peekClassRecord().MaybeRecordComputedPropWithSideEffects(ast.MustComputedProp(n).AsClassMember())
	default:
	}
	return true
}

// Java: Es6NormalizeClasses.shouldExtractExtends(Node)
func (e *Es6NormalizeClasses) ShouldExtractExtends(classNode ast.Class) bool {
	superClassNode, ok := classNode.Extends()
	return ok && !superClassNode.Node().IsQualifiedName()
}

// CanExtractExtendsSimply finds common cases where we can safely decompose
// class extends expressions which are not qualified names. Enables
// transpilation of complex extends expressions.
//
// We can only decompose the expression in a limited set of cases to avoid
// changing evaluation order of side-effect causing statements.
//
// Java: Es6NormalizeClasses.canExtractExtendsSimply(Node)
func (e *Es6NormalizeClasses) CanExtractExtendsSimply(classNode ast.Class) bool {
	enclosingStatement := rhino.CheckNotNullMsg1(NodeUtil_GetEnclosingStatement(classNode.Node()), "%s", classNode.Node())
	if classNode.Is(enclosingStatement) {
		// `class Foo extends some_expression {}`
		return true
	}

	classNodeParent := classNode.Node().GetParent()
	if declaration, ok := ast.AsNameDeclaration(enclosingStatement); ok {
		if first, ok := declaration.FirstDeclarator(); ok && first.Is(classNodeParent) && classNodeParent.IsName() {
			// `const Foo = class extends some_expression {}, maybe_other_var;`
			return true
		}
	}

	if exprStmt, ok := ast.AsExprStmt(enclosingStatement); ok {
		if assign, ok := exprStmt.Expr().AsAssign(); ok && assign.Is(classNodeParent) && assign.Value().Is(classNode.Node()) {
			// `lhs = class extends some_expression {};`
			// We can extract a temporary variable for some_expression as long as lhs expression
			// has no side effects.
			return !e.compiler.GetAstAnalyzer().MayHaveSideEffects(assign.Target().Node())
		}
	}

	return false
}

// Java: Es6NormalizeClasses.extractExtends(NodeTraversal, ClassRecord)
func (e *Es6NormalizeClasses) ExtractExtends(t *NodeTraversal, classRecord *es6NormalizeClasses_ClassRecord) {
	classNode := classRecord.classNode

	statement := NodeUtil_GetEnclosingStatement(classNode.Node())
	originalExtends, ok := classNode.Extends()
	jutil.CheckState(ok) // shouldExtractExtends: the class extends an expression

	// If the extracted superclass is itself a named class expression, remove its inner name to
	// avoid creating a new, useless variable in the scope.
	if superClass, ok := originalExtends.AsClass(); ok {
		superClass.RemoveName()
	}

	nameNode := es6NormalizeClasses_optionalGet(classRecord.normalizedNames.superClassNode)
	originalExtends.Node().ReplaceWith(nameNode)
	extendsAlias := rhino.IR_ConstNode(nameNode.CloneNode(), originalExtends.Node()).Srcref(originalExtends.Node())
	extendsAlias.InsertBefore(statement)
	NodeUtil_AddFeatureToScript(
		NodeUtil_GetEnclosingScript(classNode.Node()), parser.FeatureSet_Feature_CONST_DECLARATIONS, e.compiler)
	t.ReportCodeChangeWithNode(classNode.Node())
}

// Go-only: language-adapter — Optional<Node>.get(): an absent Optional is
// a nil node.
func es6NormalizeClasses_optionalGet(n rhino.Node) rhino.Node {
	if n.IsNil() {
		panic(jutil.NewNoSuchElementExceptionMsg("No value present"))
	}
	return n
}

// PushDownNodeIntoNewIife replaces a node with an IIFE that returns that
// node.
//
// Java: Es6NormalizeClasses.pushDownNodeIntoNewIife(Node)
func (e *Es6NormalizeClasses) PushDownNodeIntoNewIife(n rhino.Node) {
	// Create the shell of the IIFE.
	functionBody := rhino.IR_Block()
	arrowFn := rhino.IR_ArrowFunction(rhino.IR_Name(""), rhino.IR_ParamList(), functionBody).Srcref(n)
	arrowFn.SetColor(colors.StandardColors_UNKNOWN)
	iife := e.astFactory.CreateCallWithUnknownType(arrowFn).SrcrefTreeIfMissing(n)

	// This replaces n with the IIFE and detaches n from the AST.
	n.ReplaceWith(iife)

	// Now that n is detached, it's safe to make it the child of the RETURN node.
	functionBody.AddChildToBack(e.astFactory.CreateReturn(n).Srcref(n))

	NodeUtil_AddFeatureToScript(
		NodeUtil_GetEnclosingScript(iife), parser.FeatureSet_Feature_ARROW_FUNCTIONS, e.compiler)

	e.compiler.ReportChangeToEnclosingScope(iife)
}

// ShouldExtractClass determines if a class node needs to be extracted into
// a separate variable declaration.
//
// Extraction is necessary for class expressions in contexts where their
// side effects or structure might interfere with subsequent transpilation
// passes. However, some forms of class expressions are already handled
// correctly by other parts of the transpilation, specifically
// Es6RewriteClass.ClassDeclarationMetadata#create.
//
// Java: Es6NormalizeClasses.shouldExtractClass(Node)
func (e *Es6NormalizeClasses) ShouldExtractClass(classNode ast.Class) bool {
	parent := classNode.Node().GetParent()
	if NodeUtil_IsClassDeclaration(classNode.Node()) ||
		parent.IsName() ||
		es6NormalizeClasses_isQualifiedNameAssignStatement(parent) {
		// No need to extract. Handled directly by Es6RewriteClass.
		return false
	}

	return true
}

// es6NormalizeClasses_isQualifiedNameAssignStatement reports whether n is
// the assignment of an expression statement and assigns to a qualified
// name: `some.qname = value;`.
func es6NormalizeClasses_isQualifiedNameAssignStatement(n rhino.Node) bool {
	assign, ok := ast.AsAssign(n)
	return ok && assign.Target().Node().IsQualifiedName() && n.GetParent().IsExprResult()
}

// Java: Es6NormalizeClasses.extractClass(NodeTraversal, ClassRecord)
func (e *Es6NormalizeClasses) ExtractClass(t *NodeTraversal, classRecord *es6NormalizeClasses_ClassRecord) {
	classNode := classRecord.classNode.Node()
	parent := classNode.GetParent()

	info := NodeUtil_GetBestJSDocInfo(classNode)

	statement := NodeUtil_GetEnclosingStatement(parent)
	// class name node used as LHS in newly created assignment
	classNameLhs := classRecord.CreateNewNameReferenceNode()
	// class name node that replaces the class literal in the original statement
	classNameRhs := classNameLhs.CloneTree()
	classNode.ReplaceWith(classNameRhs)
	classDeclaration := rhino.IR_ConstNode(classNameLhs, classNode).SrcrefTreeIfMissing(classNode)
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_CONST_DECLARATIONS, e.compiler)
	classDeclaration.SetJSDocInfo(rhino.JSDocInfo_Builder_MaybeCopyFrom(info).Build())
	classDeclaration.InsertBefore(statement)

	// If the original statement was a variable declaration or qualified name assignment like:
	// var ClassName = class {...
	// OR
	// some.qname.ClassName = class {...
	//
	// Add @constructor JSDoc to the rewritten alias statement so ClosureCodeRemoval does not
	// mistake an @abstract class alias for a removable abstract method definition.
	if declaration, ok := ast.AsNameDeclaration(statement); ok {
		if only, ok := declaration.OnlyDeclarator(); ok && only.Is(parent) {
			e.AddAtConstructor(statement)
		}
	} else if exprStmt, ok := ast.AsExprStmt(statement); ok {
		if assign, ok := exprStmt.Expr().AsAssign(); ok &&
			assign.Target().Node().IsQualifiedName() &&
			assign.Value().Is(classNameRhs) {
			e.AddAtConstructor(assign.Node())
		}
	}
	e.compiler.ReportChangeToEnclosingScope(classDeclaration)
}

// AddAtConstructor adds at-constructor to the JSDoc of the given node.
//
// Java: Es6NormalizeClasses.addAtConstructor(Node)
func (e *Es6NormalizeClasses) AddAtConstructor(node rhino.Node) {
	builder := rhino.JSDocInfo_Builder_MaybeCopyFrom(node.GetJSDocInfo())
	builder.RecordConstructor()
	node.SetJSDocInfo(builder.Build())
}

// GetStatementDeclaringClass gets the location of the statement declaring
// the class; nil if the class cannot be extracted.
//
// Java: Es6NormalizeClasses.getStatementDeclaringClass(Node, Node)
func (e *Es6NormalizeClasses) GetStatementDeclaringClass(classNode ast.Class, classNameNode rhino.Node) rhino.Node {
	if NodeUtil_IsClassDeclaration(classNode.Node()) {
		// `class C {}` -> can use `C.staticMember` to extract static fields
		jutil.CheckState(NodeUtil_IsStatement(classNode.Node()))
		return classNode.Node()
	}
	parent := classNode.Node().GetParent()
	if parent.IsName() {
		// `let C = class {};`
		// We can use `C.staticMemberName = ...` to extract static fields
		jutil.CheckState(parent == classNameNode)
		jutil.CheckState(NodeUtil_IsStatement(classNameNode.GetParent()))
		return classNameNode.GetParent()
	}
	if assign, ok := ast.AsAssign(parent); ok &&
		assign.Target().Is(classNameNode) &&
		parent.GetParent().IsExprResult() {
		// `something.C = class {}`
		// we can use `something.C.staticMemberName = ...` to extract static fields
		jutil.CheckState(NodeUtil_IsStatement(classNameNode.GetGrandparent()))
		return classNameNode.GetGrandparent()
	}
	return rhino.NilNode
}

// Java: Es6NormalizeClasses.visit(NodeTraversal, Node, Node)
func (e *Es6NormalizeClasses) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_IN:
		if in := ast.MustBinary(n); e.transpilePrivateClassFields && in.Left().Node().IsPrivateIdentifier() {
			e.RewritePrivateIn(t, in)
		}
	case rhino.Token_GETPROP, rhino.Token_OPTCHAIN_GETPROP:
		if e.transpilePrivateClassFields && n.IsPrivateIdentifier() {
			e.RewritePrivateGetProp(t, ast.MustNormalOrOptChainGetProp(n))
		}
	case rhino.Token_CLASS:
		if len(e.classStack) == 0 {
			panic(jutil.NewNoSuchElementException())
		}
		currClassRecord := e.classStack[len(e.classStack)-1]
		e.classStack[len(e.classStack)-1] = nil
		e.classStack = e.classStack[:len(e.classStack)-1]
		jutil.CheckStateMsg1(currClassRecord.classNode.Is(n), "unexpected node: %s", n)
		cls := currClassRecord.classNode

		// First, rewrite complex extends expressions. This must happen before class extraction
		// to ensure the class's structure is simplified first.
		if currClassRecord.normalizedNames.shouldExtractExtends {
			if !e.CanExtractExtendsSimply(cls) {
				// When a class is used in an expressions where adding an alias as the previous
				// statement might change execution order of a side-effect causing statement, wrap the
				// class in an IIFE so that decomposition can happen safely.
				e.PushDownNodeIntoNewIife(n)

				// The parent of the class has been updated so reflect that.
				parent = n.GetParent()
			}
			e.ExtractExtends(t, currClassRecord)
		}

		// Second, extract the class itself if it's in an expression context or it needs a name.
		if currClassRecord.normalizedNames.shouldExtractClass {
			if e.expressionDecomposer.CanExposeExpression(n) != ExpressionDecomposer_DecompositionType_MOVABLE {
				// When class is not movable, we wrap it inside an IIFE. We have observed unsafe
				// circumstances where decomposing causes issues so this is safer. See b/417772606.
				e.PushDownNodeIntoNewIife(n)
			}
			e.ExtractClass(t, currClassRecord)

			// The parent of the class has been updated so reflect that.
			parent = n.GetParent()
		}

		classNameNode := NodeUtil_GetNameNode(n)
		jutil.CheckStateMsg1(!classNameNode.IsNil(), "Class missing a name: %s", n)
		classInsertionPoint := e.GetStatementDeclaringClass(cls, classNameNode)
		jutil.CheckStateMsg1(!classInsertionPoint.IsNil(), "Class was not extracted: %s", n)

		// For class expressions, remove the inner name if it's present.
		// Note: Safe because we removed all static initialization logic, ensured the class has a
		// name, and we replaced all references to the inner class name.
		if currClassRecord.hasBindingIdentifier && (parent.IsName() || parent.IsAssign()) {
			cls.RemoveName()
			e.compiler.ReportChangeToEnclosingScope(n)
		}

		if e.transpilePrivateClassFields &&
			(currClassRecord.hasInstancePrivateMembers ||
				currClassRecord.hasStaticPrivateMembers) {
			e.InsertPrivateDeclarations(t, currClassRecord, classInsertionPoint)
		}

		// Rewrite parts of the class based on the record.
		e.RewriteSideEffectedComputedProps(t, currClassRecord, classInsertionPoint)
		if e.transpileClassFields ||
			(e.transpilePrivateClassFields && currClassRecord.hasInstancePrivateMembers) {
			e.RewriteInstanceMembers(t, currClassRecord)
		}
		e.RewriteStaticMembers(t, currClassRecord, classInsertionPoint)
	case rhino.Token_NAME:
		e.MaybeUpdateClassSelfRef(t, ast.MustName(n))
	case rhino.Token_THIS, rhino.Token_SUPER:
		e.VisitThisAndSuper(t, n)
	default:
	}
}

// Java: Es6NormalizeClasses.maybeUpdateClassSelfRef(NodeTraversal, Node)
func (e *Es6NormalizeClasses) MaybeUpdateClassSelfRef(t *NodeTraversal, name ast.Name) {
	nameNode := name.Node()
	for i := len(e.classStack) - 1; i >= 0; i-- {
		klass := e.classStack[i]
		if !klass.hasBindingIdentifier {
			continue
		}

		bindingIdentifier := klass.bindingIdentifier.Node()
		if
		// Skip if the inner class name is the same as the class name (either this is a class
		// declaration or the outer and inner class names are the same).
		bindingIdentifier.MatchesNameNode(klass.normalizedNames.ClassNameNode()) ||
			// Skip if this IS the inner class name node (removed later when we visit the class).
			nameNode == bindingIdentifier ||
			// Skip if this nameNode is not equal to the inner class name.
			!nameNode.MatchesQualifiedNameNode(bindingIdentifier) {
			continue
		}

		// Ensure that this nameNode is not being shadowed by a different variable with the same name
		// as the inner class name.
		var_ := t.LookupVar(name.Ident())
		if var_ == nil || var_.GetNameNode() != bindingIdentifier {
			continue
		}

		newNameNode := klass.CreateNewNameReferenceNode().Srcref(nameNode)
		nameNode.ReplaceWith(newNameNode)
		e.compiler.ReportChangeToEnclosingScope(newNameNode)

		return // We made the replacement so bail.
	}
}

// VisitThisAndSuper visits `this` and `super` nodes in the class. Only
// makes changes in static contexts and records super calls in the
// constructor for later.
//
// For static `this`, we can only replace with the class name if in a static
// initialization context (i.e. in a static field initializer or a static
// block).
//
// For static `super`, we can always replace with the extended class name.
//
// Note: This runs before #rewriteStaticMembers so the static members are
// still there.
//
// Java: Es6NormalizeClasses.visitThisAndSuper(NodeTraversal, Node)
func (e *Es6NormalizeClasses) VisitThisAndSuper(t *NodeTraversal, n rhino.Node) {
	if Es6NormalizeClasses_IsSuperCall(n.GetParent()) {
		// super(...); calls are only legal in a constructor so we don't need to check for static
		// context.
		classRecord := e.peekClassRecord()
		classRecord.superCalls = append(classRecord.superCalls, n.GetParent())
		return
	}
	rootNode := t.GetClosestScopeRootNodeBindingThisOrSuper()
	rootParent := rootNode.GetParent()
	if rootNode.IsFunction() &&
		(rootParent.IsMemberFunctionDef() ||
			rootParent.IsComputedProp() ||
			rootParent.IsGetterDef() ||
			rootParent.IsSetterDef()) {
		rootNode = rootParent
	}
	inStaticBlock := rootNode.IsBlock() && rootNode.GetParent().IsClassMembers()
	if !rootNode.IsStaticMember() && !inStaticBlock {
		return
	}

	// We always rewrite static this to the class name in a static initialization context so that we
	// can normalize-away static initialization from field initializers and static blocks.
	forceRewriteStaticThis := n.IsThis() && Es6NormalizeClasses_IsInStaticInitializationContext(rootNode)

	if !e.doStaticInheritanceRewrites && !forceRewriteStaticThis {
		return
	}

	classNode := rootNode.GetGrandparent()
	classRecord := e.peekClassRecord()
	jutil.CheckStateMsg2(
		classRecord.classNode.Is(classNode),
		"wrong class node: %s != %s",
		classRecord.classNode.Node(),
		classNode)

	var newNameNode rhino.Node
	if forceRewriteStaticThis {
		// For static this, we can only universally replace with the class name in a static
		// initialization context.
		newNameNode = classRecord.CreateNewNameReferenceNode().SrcrefTree(n)
	} else if e.doStaticSuperRewrites && n.IsSuper() {
		// For a static super, we can always rewrite the super as something else.
		newNameNode = classRecord.CreateNewSuperclassReferenceNode()
	}

	if !newNameNode.IsNil() {
		originalName := jutil.CheckNotNull(n.GetQualifiedName())
		newNameNode.SrcrefTree(n).SetOriginalName(originalName)
		n.ReplaceWith(newNameNode)
		t.ReportCodeChangeWithNode(newNameNode)
	}
}

// Java: Es6NormalizeClasses.isSuperCall(Node)
func Es6NormalizeClasses_IsSuperCall(n rhino.Node) bool {
	return n.IsCall() && n.GetFirstChild().IsSuper()
}

// Es6NormalizeClasses_IsInStaticInitializationContext assumes that we have
// already filtered out non-static members and non-static blocks.
//
// Java: Es6NormalizeClasses.isInStaticInitializationContext(Node)
func Es6NormalizeClasses_IsInStaticInitializationContext(rootNode rhino.Node) bool {
	switch rootNode.GetToken() {
	case rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_COMPUTED_PROP, rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF:
		jutil.CheckState(rootNode.IsStaticMember())
		return false
	case rhino.Token_MEMBER_FIELD_DEF, rhino.Token_COMPUTED_FIELD_DEF:
		jutil.CheckState(rootNode.IsStaticMember())
		return true
	case rhino.Token_BLOCK:
		jutil.CheckState(rootNode.GetParent().IsClassMembers())
		return true
	default:
		panic(jutil.NewIllegalArgumentExceptionMsg("invalid root " + jutil.ToJavaString(rootNode)))
	}
}

// Es6NormalizeShorthandProperties normalizes shorthand object properties.
// This should be one of the first things done when transpiling from ES6 down
// to ES5, as it allows all the following checks and transpilations to not
// care about shorthand and destructured assignments.
type Es6NormalizeShorthandProperties struct{}

var es6NormalizeShorthandTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_STRING_KEY)
	return s
}()

// Java: Es6NormalizeShorthandProperties(AbstractCompiler)
func NewEs6NormalizeShorthandProperties(compiler AbstractCompiler) *Es6NormalizeShorthandProperties {
	return &Es6NormalizeShorthandProperties{}
}

func (*Es6NormalizeShorthandProperties) HandledTokens() ast.TokenSet {
	return es6NormalizeShorthandTokens
}

func (*Es6NormalizeShorthandProperties) GetAdditionalFeaturesToRunOn() *parser.FeatureSet {
	return parser.FeatureSet_BARE_MINIMUM
}

// Java: Es6NormalizeShorthandProperties.transpileSubtree(Node)
func (e *Es6NormalizeShorthandProperties) TranspileSubtree(subtree rhino.Node) rhino.Node {
	if subtree.IsStringKey() {
		subtree.SetShorthandProperty(false)
	}
	return subtree
}

// Java: Es6NormalizeShorthandProperties.getTranspiledAwayFeatures()
func (e *Es6NormalizeShorthandProperties) GetTranspiledAwayFeatures() *parser.FeatureSet {
	return parser.FeatureSet_BARE_MINIMUM.With(parser.FeatureSet_Feature_SHORTHAND_OBJECT_PROPERTIES)
}

type Es6ScopeRenameKey struct {
	Root rhino.Node
	Name string
}

// Es6RenameReferences renames references in code.
type Es6RenameReferences struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	// Map from the root node of a scope + the original variable name to a
	// new name for the variable.
	//
	// An entry in this map means we should rename a variable that was
	// originally declared in the scope with the given root and had the
	// original variable name to the new name.
	renameTable map[Es6ScopeRenameKey]string
}

// Java: Es6RenameReferences(Table<Node, String, String>)
func NewEs6RenameReferences(renameTable map[Es6ScopeRenameKey]string) *Es6RenameReferences {
	return &Es6RenameReferences{renameTable: renameTable}
}

// Java: Es6RenameReferences.visit(NodeTraversal, Node, Node)
func (e *Es6RenameReferences) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if NodeUtil_IsReferenceName(n) {
		e.RenameReference(t, n, n.GetString())
	}
}

// RenameReference: t is the NodeTraversal; oldName is the root name of a
// qualified name.
//
// Java: Es6RenameReferences.renameReference(NodeTraversal, Node, String)
func (e *Es6RenameReferences) RenameReference(t *NodeTraversal, n rhino.Node, oldName string) {
	current := t.GetScope()
	// You should be wondering:
	//
	// Why are we searching up the stack of scopes here instead of just looking up the Var for
	// oldName and getting the scope from that?
	//
	// The answer is:
	//
	// When Es6RewriteBlockScopedDeclaration uses this class, it has already modified the scopes so
	// that `oldName` is no longer declared in its original scope.
	//
	// I don't know if the other users of this class are also modifying the scopes in advance.
	// - bradfordcsmith@google.com
	for current != nil {
		newName, ok := e.renameTable[Es6ScopeRenameKey{Root: current.GetRootNode(), Name: oldName}]
		if ok {
			n.SetString(newName)
			t.ReportCodeChange()
			return
		} else if current.HasOwnSlot(oldName) {
			return
		} else {
			current = current.GetParent()
		}
	}
}

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
