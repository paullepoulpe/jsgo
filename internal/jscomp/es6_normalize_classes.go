// Port of com/google/javascript/jscomp/Es6NormalizeClasses.java @ closure-compiler v20260415

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
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

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
	compiler             AbstractCompiler
	astFactory           *AstFactory
	expressionDecomposer *ExpressionDecomposer
	ctorCreator          *SynthesizeExplicitConstructors
	transpileClassFields bool

	doStaticInheritanceRewrites bool
	doStaticSuperRewrites       bool

	// Go-only: language-adapter — the Java ArrayDeque used as a stack
	// (addFirst/peek/removeFirst, iterated from the first element): the
	// first element is the last one here.
	classStack []*es6NormalizeClasses_ClassRecord
}

// Java: Es6NormalizeClasses.CLASS_DECL_VAR
const es6NormalizeClasses_CLASS_DECL_VAR = "$jscomp$classDecl$"

// Java: Es6NormalizeClasses.CLASS_EXTENDS_VAR
const es6NormalizeClasses_CLASS_EXTENDS_VAR = "$jscomp$classExtends$"

// Java: Es6NormalizeClasses.COMP_FIELD_VAR
const es6NormalizeClasses_COMP_FIELD_VAR = "$jscomp$compField$"

// Java: Es6NormalizeClasses.STATIC_INIT_METHOD_NAME
const es6NormalizeClasses_STATIC_INIT_METHOD_NAME = "$jscomp$staticInit$"

// Java: Es6NormalizeClasses.GENERIC_NAME_REPLACEMENTS
var Es6NormalizeClasses_GENERIC_NAME_REPLACEMENTS = jutil.NewImmutableMapBuilder[string, string]().
	Put("CLASS_DECL", es6NormalizeClasses_CLASS_DECL_VAR).
	Put("CLASS_EXTENDS", es6NormalizeClasses_CLASS_EXTENDS_VAR).
	Put("COMP_FIELD", es6NormalizeClasses_COMP_FIELD_VAR).
	Put("STATIC_INIT", es6NormalizeClasses_STATIC_INIT_METHOD_NAME).
	BuildOrThrow()

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

	e.doStaticInheritanceRewrites =
		options.GetAssumeStaticInheritanceIsNotUsed() && !options.GetSkipNonTranspilationPasses()

	// Es6ConvertSuper handles static super rewrites otherwise.
	e.doStaticSuperRewrites = !options.NeedsTranspilationOf(parser.FeatureSet_Feature_SUPER)
	return e
}

// Java: Es6NormalizeClasses.process(Node, Node)
func (e *Es6NormalizeClasses) Process(externs *rhino.Node, root *rhino.Node) {
	NodeTraversal_TraverseRoots(e.compiler, e, externs, root)

	// We transpiled away class static blocks. We always do this, regardless of language out, to
	// normalize the AST.
	transpiledFeatures := parser.FeatureSet_BARE_MINIMUM.With(parser.FeatureSet_Feature_CLASS_STATIC_BLOCK)
	if e.transpileClassFields {
		// We removed public and static field declarations if they need transpilation.
		transpiledFeatures = transpiledFeatures.With(parser.FeatureSet_Feature_PUBLIC_CLASS_FIELDS)
	}
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(e.compiler, root, transpiledFeatures)
}

// Go-only: language-adapter — Deque.peek(): the first element, or nil.
func (e *Es6NormalizeClasses) peekClassRecord() *es6NormalizeClasses_ClassRecord {
	if len(e.classStack) == 0 {
		return nil
	}
	return e.classStack[len(e.classStack)-1]
}

// Java: Es6NormalizeClasses.shouldTraverse(NodeTraversal, Node, Node)
func (e *Es6NormalizeClasses) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		scriptFeatures := NodeUtil_GetFeatureSetOfScript(n)
		return scriptFeatures == nil ||
			scriptFeatures.ContainsFeature(parser.FeatureSet_Feature_CLASSES) ||
			scriptFeatures.ContainsFeature(parser.FeatureSet_Feature_CLASS_STATIC_BLOCK) ||
			scriptFeatures.ContainsFeature(parser.FeatureSet_Feature_PUBLIC_CLASS_FIELDS)
	case rhino.Token_CLASS:
		// The class name for declarations and class expression inner names.
		var bindingIdentifier *rhino.Node
		if n.GetFirstChild().IsName() {
			bindingIdentifier = n.GetFirstChild()
		}

		// Generate any needed normalized class names in pre-order so that we can rewrite the inner
		// class name to the new outer name if needed.
		// (We don't actually do the normalization in pre-order because that can cause
		// NodeTraversal to skip traversing other parts of the AST.)
		normalized := e.GenerateNormalizedClassNames(t, n)

		// Create the ClassRecord which will be referenced when we visit the class later.
		e.classStack = append(e.classStack, newEs6NormalizeClasses_ClassRecord(normalized, bindingIdentifier /* classNode= */, n))
	case rhino.Token_COMPUTED_FIELD_DEF:
		jutil.CheckState(len(e.classStack) != 0)
		classRecord := e.peekClassRecord()
		classRecord.MaybeRecordComputedPropWithSideEffects(n)
		classRecord.RecordField(n)
	case rhino.Token_MEMBER_FIELD_DEF:
		jutil.CheckState(len(e.classStack) != 0)
		e.peekClassRecord().RecordField(n)
	case rhino.Token_BLOCK:
		if !NodeUtil_IsClassStaticBlock(n) {
			break
		}

		jutil.CheckState(len(e.classStack) != 0)
		e.peekClassRecord().RecordStaticBlock(n)
	// For example, a class method with a computed property name.
	case rhino.Token_COMPUTED_PROP:
		if !n.GetParent().IsClassMembers() {
			break
		}

		jutil.CheckState(len(e.classStack) != 0)
		e.peekClassRecord().MaybeRecordComputedPropWithSideEffects(n)
	default:
	}
	return true
}

// Java: Es6NormalizeClasses.shouldExtractExtends(Node)
func (e *Es6NormalizeClasses) ShouldExtractExtends(classNode *rhino.Node) bool {
	jutil.CheckArgument(classNode.IsClass())
	superClassNode := classNode.GetSecondChild()
	return !superClassNode.IsEmpty() && !superClassNode.IsQualifiedName()
}

// CanExtractExtendsSimply finds common cases where we can safely decompose
// class extends expressions which are not qualified names. Enables
// transpilation of complex extends expressions.
//
// We can only decompose the expression in a limited set of cases to avoid
// changing evaluation order of side-effect causing statements.
//
// Java: Es6NormalizeClasses.canExtractExtendsSimply(Node)
func (e *Es6NormalizeClasses) CanExtractExtendsSimply(classNode *rhino.Node) bool {
	enclosingStatement := jutil.CheckNotNullMsg1(NodeUtil_GetEnclosingStatement(classNode), "%s", classNode)
	if enclosingStatement == classNode {
		// `class Foo extends some_expression {}`
		return true
	}

	classNodeParent := classNode.GetParent()
	if NodeUtil_IsNameDeclaration(enclosingStatement) &&
		classNodeParent.IsName() &&
		classNodeParent.IsFirstChildOf(enclosingStatement) {
		// `const Foo = class extends some_expression {}, maybe_other_var;`
		return true
	}

	if enclosingStatement.IsExprResult() &&
		classNodeParent.IsOnlyChildOf(enclosingStatement) &&
		classNodeParent.IsAssign() &&
		classNode.IsSecondChildOf(classNodeParent) {
		// `lhs = class extends some_expression {};`
		lhsNode := classNodeParent.GetFirstChild()
		// We can extract a temporary variable for some_expression as long as lhs expression
		// has no side effects.
		return !e.compiler.GetAstAnalyzer().MayHaveSideEffects(lhsNode)
	}

	return false
}

// Java: Es6NormalizeClasses.extractExtends(NodeTraversal, ClassRecord)
func (e *Es6NormalizeClasses) ExtractExtends(t *NodeTraversal, classRecord *es6NormalizeClasses_ClassRecord) {
	classNode := classRecord.classNode

	statement := NodeUtil_GetEnclosingStatement(classNode)
	originalExtends := classNode.GetSecondChild()

	// If the extracted superclass is itself a named class expression, remove its inner name to
	// avoid creating a new, useless variable in the scope.
	if originalExtends.IsClass() && originalExtends.GetFirstChild().IsName() {
		originalExtends.
			GetFirstChild().
			ReplaceWith(rhino.IR_Empty().Srcref(originalExtends.GetFirstChild()))
	}

	nameNode := es6NormalizeClasses_optionalGet(classRecord.normalizedNames.superClassNode)
	originalExtends.ReplaceWith(nameNode)
	extendsAlias := rhino.IR_ConstNode(nameNode.CloneNode(), originalExtends).Srcref(originalExtends)
	extendsAlias.InsertBefore(statement)
	NodeUtil_AddFeatureToScript(
		NodeUtil_GetEnclosingScript(classNode), parser.FeatureSet_Feature_CONST_DECLARATIONS, e.compiler)
	t.ReportCodeChangeWithNode(classNode)
}

// Go-only: language-adapter — Optional<Node>.get(): an absent Optional is
// a nil node.
func es6NormalizeClasses_optionalGet(n *rhino.Node) *rhino.Node {
	if n == nil {
		panic(jutil.NewNoSuchElementExceptionMsg("No value present"))
	}
	return n
}

// PushDownNodeIntoNewIife replaces a node with an IIFE that returns that
// node.
//
// Java: Es6NormalizeClasses.pushDownNodeIntoNewIife(Node)
func (e *Es6NormalizeClasses) PushDownNodeIntoNewIife(n *rhino.Node) {
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
// Es6ToEs3Converter.ClassDeclarationMetadata#create.
//
// Java: Es6NormalizeClasses.shouldExtractClass(Node)
func (e *Es6NormalizeClasses) ShouldExtractClass(classNode *rhino.Node) bool {
	parent := classNode.GetParent()
	if NodeUtil_IsClassDeclaration(classNode) ||
		parent.IsName() ||
		(parent.IsAssign() &&
			parent.GetFirstChild().IsQualifiedName() &&
			parent.GetParent().IsExprResult()) {
		// No need to extract. Handled directly by Es6ToEs3Converter.ClassDeclarationMetadata#create.
		return false
	}

	return true
}

// Java: Es6NormalizeClasses.extractClass(NodeTraversal, ClassRecord)
func (e *Es6NormalizeClasses) ExtractClass(t *NodeTraversal, classRecord *es6NormalizeClasses_ClassRecord) {
	classNode := classRecord.classNode
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

	// If the original statement was a variable declaration or qualified name assignment like
	// like these:
	// var ClassName = class {...
	// OR
	// some.qname.ClassName = class {...
	//
	// We will have changed the original statement to
	//
	// var ClassName = generatedName;
	// OR
	// some.qname.ClassName = generatedName;
	//
	// This is creating a type alias for a class, but since there's no literal class on the RHS,
	// it doesn't look like one. Add at-constructor JSDoc to make it clear that this is happening.
	//
	// This was added to fix a specific problem where the original definition was for an abstract
	// class, so its JSDoc included at-abstract.
	// This caused ClosureCodeRemoval to think this rewritten assignment was a removable abstract
	// method definition instead of the definition of an abstract class.
	//
	// TODO(b/117292942): Make ClosureCodeRemoval smarter so this hack isn't necessary to
	// prevent incorrect removal of assignments.
	if NodeUtil_IsNameDeclaration(statement) &&
		statement.HasOneChild() &&
		statement.GetOnlyChild() == parent {
		// var ClassName = generatedName;
		e.AddAtConstructor(statement)
	} else if statement.IsExprResult() {
		expr := statement.GetOnlyChild()
		if expr.IsAssign() &&
			expr.GetFirstChild().IsQualifiedName() &&
			expr.GetSecondChild() == classNameRhs {
			// some.qname.ClassName = generatedName;
			e.AddAtConstructor(expr)
		}
	}
	e.compiler.ReportChangeToEnclosingScope(classDeclaration)
}

// AddAtConstructor adds at-constructor to the JSDoc of the given node.
//
// Java: Es6NormalizeClasses.addAtConstructor(Node)
func (e *Es6NormalizeClasses) AddAtConstructor(node *rhino.Node) {
	builder := rhino.JSDocInfo_Builder_MaybeCopyFrom(node.GetJSDocInfo())
	builder.RecordConstructor()
	node.SetJSDocInfo(builder.Build())
}

// GetStatementDeclaringClass gets the location of the statement declaring
// the class; nil if the class cannot be extracted.
//
// Java: Es6NormalizeClasses.getStatementDeclaringClass(Node, Node)
func (e *Es6NormalizeClasses) GetStatementDeclaringClass(classNode *rhino.Node, classNameNode *rhino.Node) *rhino.Node {
	if NodeUtil_IsClassDeclaration(classNode) {
		// `class C {}` -> can use `C.staticMember` to extract static fields
		jutil.CheckState(NodeUtil_IsStatement(classNode))
		return classNode
	}
	parent := classNode.GetParent()
	if parent.IsName() {
		// `let C = class {};`
		// We can use `C.staticMemberName = ...` to extract static fields
		jutil.CheckState(parent == classNameNode)
		jutil.CheckState(NodeUtil_IsStatement(classNameNode.GetParent()))
		return classNameNode.GetParent()
	}
	if parent.IsAssign() &&
		parent.GetFirstChild() == classNameNode &&
		parent.GetParent().IsExprResult() {
		// `something.C = class {}`
		// we can use `something.C.staticMemberName = ...` to extract static fields
		jutil.CheckState(NodeUtil_IsStatement(classNameNode.GetGrandparent()))
		return classNameNode.GetGrandparent()
	}
	return nil
}

// Java: Es6NormalizeClasses.enterScope(NodeTraversal)
func (e *Es6NormalizeClasses) EnterScope(t *NodeTraversal) {
	scopeRoot := t.GetScopeRoot()
	if NodeUtil_IsFunctionBlock(scopeRoot) && NodeUtil_IsEs6Constructor(scopeRoot.GetParent()) {
		e.peekClassRecord().RecordConstructorScope(t.GetScope())
	}
}

// Java: Es6NormalizeClasses.exitScope(NodeTraversal)
func (e *Es6NormalizeClasses) ExitScope(t *NodeTraversal) {}

// Java: Es6NormalizeClasses.visit(NodeTraversal, Node, Node)
func (e *Es6NormalizeClasses) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_CLASS:
		if len(e.classStack) == 0 {
			panic(jutil.NewNoSuchElementException())
		}
		currClassRecord := e.classStack[len(e.classStack)-1]
		e.classStack[len(e.classStack)-1] = nil
		e.classStack = e.classStack[:len(e.classStack)-1]
		jutil.CheckStateMsg1(currClassRecord.classNode == n, "unexpected node: %s", n)

		// First, rewrite complex extends expressions. This must happen before class extraction
		// to ensure the class's structure is simplified first.
		if currClassRecord.normalizedNames.shouldExtractExtends {
			if !e.CanExtractExtendsSimply(n) {
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
		jutil.CheckStateMsg1(classNameNode != nil, "Class missing a name: %s", n)
		classInsertionPoint := e.GetStatementDeclaringClass(n, classNameNode)
		jutil.CheckStateMsg1(classInsertionPoint != nil, "Class was not extracted: %s", n)

		// For class expressions, remove the inner name if it's present.
		// Note: Safe because we removed all static initialization logic, ensured the class has a
		// name, and we replaced all references to the inner class name.
		if currClassRecord.bindingIdentifier != nil &&
			(parent.IsName() || parent.IsAssign()) {
			n.GetFirstChild().ReplaceWith(rhino.IR_Empty().Srcref(n.GetFirstChild()))
			e.compiler.ReportChangeToEnclosingScope(n)
		}

		// Rewrite parts of the class based on the record.
		e.RewriteSideEffectedComputedProps(t, currClassRecord, classInsertionPoint)
		if e.transpileClassFields {
			e.RewriteInstanceMembers(t, currClassRecord)
		}
		e.RewriteStaticMembers(t, currClassRecord, classInsertionPoint)
	case rhino.Token_NAME:
		e.MaybeUpdateClassSelfRef(t, n)
	case rhino.Token_THIS, rhino.Token_SUPER:
		e.VisitThisAndSuper(t, n)
	default:
	}
}

// RewriteSideEffectedComputedProps rewrites and moves all side effected
// computed field keys to the top
//
// Java: Es6NormalizeClasses.rewriteSideEffectedComputedProps(NodeTraversal, ClassRecord, Node)
func (e *Es6NormalizeClasses) RewriteSideEffectedComputedProps(
	t *NodeTraversal, record *es6NormalizeClasses_ClassRecord, insertionPoint *rhino.Node) {
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
func es6NormalizeClasses_dequeRemove(deque *[]*rhino.Node) *rhino.Node {
	if len(*deque) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	first := (*deque)[0]
	(*deque)[0] = nil
	*deque = (*deque)[1:]
	return first
}

// ExtractExpressionFromCompField extracts the expression in the LHS of a
// computed field to not disturb possible side effects and allow for this
// and super to be used in the LHS of a computed field in certain cases. Does
// not extract a computed field that was already moved into a computed
// function.
//
// E.g.
//
//	class Foo {
//	  [bar('str')] = 4;
//	}
//
// becomes
//
//	var $jscomp$computedfield$0 = bar('str');
//	class Foo {
//	  [$jscomp$computedfield$0] = 4;
//	}
//
// Java: Es6NormalizeClasses.extractExpressionFromCompField(NodeTraversal, ClassRecord, Node, Node)
func (e *Es6NormalizeClasses) ExtractExpressionFromCompField(
	t *NodeTraversal, record *es6NormalizeClasses_ClassRecord, memberField *rhino.Node, insertionPoint *rhino.Node) {
	jutil.CheckArgumentMsg1(memberField.IsComputedFieldDef() || memberField.IsComputedProp(), "%s", memberField)

	compExpression := memberField.RemoveFirstChild()
	compFieldVar := e.astFactory.
		CreateSingleVarNameDeclarationWithValue(e.GenerateUniqueCompFieldVarName(t), compExpression).
		SrcrefTreeIfMissing(record.classNode)
	compFieldName := compFieldVar.GetFirstChild()
	memberField.AddChildToFront(compFieldName.CloneNode())
	compFieldVar.InsertBefore(insertionPoint)
	compFieldVar.SrcrefTreeIfMissing(record.classNode)
}

// RewriteInstanceMembers rewrites and moves all instance fields
//
// Java: Es6NormalizeClasses.rewriteInstanceMembers(NodeTraversal, ClassRecord)
func (e *Es6NormalizeClasses) RewriteInstanceMembers(t *NodeTraversal, record *es6NormalizeClasses_ClassRecord) {
	instanceMembers := &record.instanceMembers

	if len(*instanceMembers) == 0 {
		return
	}
	e.ctorCreator.SynthesizeClassConstructorIfMissing(t, record.classNode)
	ctor := NodeUtil_GetEs6ClassConstructorMemberFunctionDef(record.classNode)
	ctorBlock := ctor.GetFirstChild().GetLastChild()
	insertionPoint := e.AddTemporaryInsertionPoint(ctorBlock, record)

	for len(*instanceMembers) != 0 {
		instanceMember := es6NormalizeClasses_dequeRemove(instanceMembers)
		jutil.CheckStateMsg1(
			instanceMember.IsMemberFieldDef() || instanceMember.IsComputedFieldDef(), "%s", instanceMember)

		thisNode := e.astFactory.CreateThisForEs6ClassMember(instanceMember)
		var transpiledNode *rhino.Node
		if instanceMember.IsMemberFieldDef() {
			transpiledNode = e.ConvNonCompFieldToGetProp(thisNode, instanceMember /* isStatic= */, false)
		} else {
			transpiledNode = e.ConvCompFieldToGetElem(thisNode, instanceMember /* isStatic= */, false)
		}

		transpiledNode.InsertBefore(jutil.CheckNotNull(insertionPoint))
	}

	insertionPoint.Detach()

	// We moved the field from the class body to the constructor, so we need both.
	t.ReportCodeChange()
	t.ReportCodeChangeWithNode(ctorBlock)
}

// AddTemporaryInsertionPoint finds the location of super() call in the
// constructor and add a temporary empty node after the super() call. If
// there is no super() call, add a temporary empty node at the beginning of
// the constructor body after all function declarations.
//
// Returns the added temporary empty node.
//
// Java: Es6NormalizeClasses.addTemporaryInsertionPoint(Node, ClassRecord)
func (e *Es6NormalizeClasses) AddTemporaryInsertionPoint(ctorBlock *rhino.Node, classRecord *es6NormalizeClasses_ClassRecord) *rhino.Node {
	jutil.CheckStateMsg1(
		len(classRecord.superCalls) <= 1,
		"classes with public fields must have only one super() call at the constructor root: %s",
		classRecord.classNode)

	tempNode := rhino.IR_Empty()
	if len(classRecord.superCalls) == 1 {
		stmt := e.EnsureSuperCallIsOwnStatement(classRecord.superCalls[0])
		tempNode.InsertAfter(stmt)
		return tempNode
	}
	// Search for super() calls added by SynthesizeExplicitConstructors.
	for stmt := ctorBlock.GetFirstChild(); stmt != nil; stmt = stmt.GetNext() {
		if NodeUtil_IsExprCall(stmt) && stmt.GetFirstFirstChild().IsSuper() {
			tempNode.InsertAfter(stmt)
			return tempNode
		}
	}
	insertionPoint := NodeUtil_GetInsertionPointAfterAllInnerFunctionDeclarations(ctorBlock)
	if insertionPoint != nil {
		tempNode.InsertBefore(insertionPoint)
	} else {
		// functionBody only contains hoisted function declarations
		ctorBlock.AddChildToBack(tempNode)
	}
	return tempNode
}

// EnsureSuperCallIsOwnStatement returns a statement containing either
// `super(...);` or `var tmp = super(...)`, extracting the given super call
// into a `tmp` variable if necessary.
//
// Java: Es6NormalizeClasses.ensureSuperCallIsOwnStatement(Node)
func (e *Es6NormalizeClasses) EnsureSuperCallIsOwnStatement(superCall *rhino.Node) *rhino.Node {
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

// ConvNonCompFieldToGetProp creates a node that represents receiver.key =
// value; where the key and value comes from the non-computed field
//
// Java: Es6NormalizeClasses.convNonCompFieldToGetProp(Node, Node, boolean)
func (e *Es6NormalizeClasses) ConvNonCompFieldToGetProp(
	receiver *rhino.Node, noncomputedField *rhino.Node, isStatic bool) *rhino.Node {
	jutil.CheckArgument(noncomputedField.IsMemberFieldDef())
	jutil.CheckArgumentMsg1(receiver.GetParent() == nil, "%s", receiver)

	var fieldValue *rhino.Node
	// We always move out the static field initialization.
	if isStatic || e.transpileClassFields {
		fieldValue = noncomputedField.GetFirstChild()
	}

	return e.CreateFieldAssignment(
		noncomputedField,
		fieldValue,
		func() *rhino.Node {
			return e.astFactory.CreateGetProp(
				receiver, noncomputedField.GetString(), AstFactory_Type_(noncomputedField))
		})
}

// ConvCompFieldToGetElem creates a node that represents receiver[key] =
// value; where the key and value comes from the computed field
//
// Java: Es6NormalizeClasses.convCompFieldToGetElem(Node, Node, boolean)
func (e *Es6NormalizeClasses) ConvCompFieldToGetElem(
	receiver *rhino.Node, computedField *rhino.Node, isStatic bool) *rhino.Node {
	jutil.CheckArgumentMsg1(computedField.IsComputedFieldDef(), "%s", computedField)
	jutil.CheckArgumentMsg1(receiver.GetParent() == nil, "%s", receiver)

	var fieldValue *rhino.Node
	if computedField.GetChildCount() == 2 && (isStatic || e.transpileClassFields) {
		fieldValue = computedField.GetLastChild()
	}

	return e.CreateFieldAssignment(
		computedField,
		fieldValue,
		func() *rhino.Node {
			var compFieldNameExpr *rhino.Node
			if e.transpileClassFields {
				compFieldNameExpr = computedField.RemoveFirstChild()
			} else {
				compFieldNameExpr = computedField.GetFirstChild().CloneTree()
			}
			return e.astFactory.CreateGetElem(receiver, compFieldNameExpr)
		})
}

// CreateFieldAssignment is a shared helper to create an assignment
// statement for a class field. Handles logic for creating undefined values
// for uninitialized fields during transpilation, as well as JSDoc transfer
// and source reference mapping.
//
// Java: Es6NormalizeClasses.createFieldAssignment(Node, Node, Supplier<Node>)
func (e *Es6NormalizeClasses) CreateFieldAssignment(
	field *rhino.Node, fieldValue *rhino.Node, createAssignLhs func() *rhino.Node) *rhino.Node {
	if e.transpileClassFields {
		// We don't need to keep around the field declaration.
		field.Detach()
	}

	// We don't make a copy of the static field declaration if it doesn't have an initializer and
	// we're not transpiling class fields.
	if fieldValue == nil {
		if !e.transpileClassFields {
			return nil
		}

		// If we are transpiling class fields, we want an explicit undefined value so we properly
		// initialize the field.
		fieldValue = e.astFactory.CreateUndefinedValue()
	} else {
		// We are going to move fieldValue to a new parent, so detach it now.
		fieldValue.Detach()
	}

	result := e.astFactory.CreateAssignStatement(createAssignLhs(), fieldValue)

	if e.transpileClassFields {
		// Move any JSDoc from the field declaration to the child of the EXPR_RESULT, which represents
		// the new declaration.
		// field is already detached, so there's no benefit to calling
		// NodeUtil.getBestJSDocInfo(field). For now at least, the JSDocInfo we want will always be
		// directly on field in all cases.
		result.GetFirstChild().SetJSDocInfo(field.GetJSDocInfo())
	}
	result.SrcrefTreeIfMissing(field)
	return result
}

// RewriteStaticMembers rewrites and moves all static blocks and fields
//
// Java: Es6NormalizeClasses.rewriteStaticMembers(NodeTraversal, ClassRecord, Node)
func (e *Es6NormalizeClasses) RewriteStaticMembers(t *NodeTraversal, record *es6NormalizeClasses_ClassRecord, insertionPoint *rhino.Node) {
	staticMembers := &record.staticMembers
	if len(*staticMembers) == 0 {
		return
	}

	staticInitTempParent := rhino.IR_Block()
	for len(*staticMembers) != 0 {
		staticMember := es6NormalizeClasses_dequeRemove(staticMembers)
		// if the name is a property access, we want the whole chain of accesses, while for other
		// cases we only want the name node
		nameToUse := record.CreateNewNameReferenceNode().SrcrefTree(staticMember)

		var transpiledNode *rhino.Node
		switch staticMember.GetToken() {
		case rhino.Token_BLOCK:
			block := staticMember.Detach()
			t.ReportCodeChange()
			// We ignore empty static initialization blocks.
			if block.HasChildren() {
				transpiledNode = block
			}
		case rhino.Token_MEMBER_FIELD_DEF:
			transpiledNode = e.ConvNonCompFieldToGetProp(nameToUse, staticMember /* isStatic= */, true)
		case rhino.Token_COMPUTED_FIELD_DEF:
			transpiledNode = e.ConvCompFieldToGetElem(nameToUse, staticMember /* isStatic= */, true)
		default:
			panic(jutil.NewIllegalStateExceptionMsg(jutil.ToJavaString(staticMember)))
		}
		if transpiledNode != nil {
			staticInitTempParent.AddChildToBack(transpiledNode)
		}
	}

	// If there is no static initialization logic, we are done.
	if !staticInitTempParent.HasChildren() {
		return
	}

	// If there are no super references, place all the static initialization code after the class.
	// TODO: b/236744850 - Also skip if there are private field refs.
	if !NodeUtil_Has(
		staticInitTempParent,
		/* pred= */ (*rhino.Node).IsSuper,
		/* traverseChildrenPred= */ func(n *rhino.Node) bool {
			return !n.IsClass() && (!n.IsFunction() || n.IsArrowFunction())
		}) {
		for staticInitTempParent.HasChildren() {
			staticInitTempParent.GetLastChild().Detach().InsertAfter(insertionPoint)
		}
		t.ReportCodeChangeWithNode(insertionPoint.GetParent())
		return
	}

	// We will introduce a new static method and a call to it to contain static initialization code.

	// Add new static method "$jscomp$staticInit$[FILE_ID]$[number]".
	// All the static initialization logic goes into this method. We need a static method as opposed
	// to doing this logic alongside the class because #private members can only be accessed within
	// the class.
	staticInitMethodName := e.GenerateUniqueStaticInitMethodName(t)
	staticInitMethodFunc :=
		e.astFactory.CreateEmptyFunction(AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_FUNCTION_TYPE, colors.StandardColors_TOP_OBJECT))
	staticInitMethod :=
		e.astFactory.CreateMemberFunctionDef(staticInitMethodName, staticInitMethodFunc)
	staticInitMethod.SetStaticMember(true)

	staticInitBlock := staticInitMethod.GetFirstChild().GetLastChild()
	staticInitBlock.AddChildrenToFront(staticInitTempParent.RemoveChildren())

	classMembers := NodeUtil_GetClassMembers(record.classNode)
	staticInitMethod.SrcrefTree(classMembers)
	classMembers.AddChildToBack(staticInitMethod)
	e.compiler.ReportChangeToChangeScope(staticInitMethodFunc)
	t.ReportCodeChangeWithNode(staticInitMethod)
	// Record the feature corresponding to a static method.
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_MEMBER_DECLARATIONS, e.compiler)

	// Add a call to the static initialization method after the class definition.
	staticInitCall := e.astFactory.
		CreateCall(
			e.astFactory.CreateGetProp(
				record.CreateNewNameReferenceNode(),
				staticInitMethodName,
				AstFactory_Type_(staticInitMethodFunc)),
			AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_VOID_TYPE, colors.StandardColors_NULL_OR_VOID)).
		SrcrefTree(classMembers)
	staticInitCallStmt := e.astFactory.ExprResult(staticInitCall)
	staticInitCallStmt.InsertAfter(insertionPoint)
}

// Java: Es6NormalizeClasses.maybeUpdateClassSelfRef(NodeTraversal, Node)
func (e *Es6NormalizeClasses) MaybeUpdateClassSelfRef(t *NodeTraversal, nameNode *rhino.Node) {
	for i := len(e.classStack) - 1; i >= 0; i-- {
		klass := e.classStack[i]
		if klass.bindingIdentifier == nil {
			continue
		}

		bindingIdentifier := klass.bindingIdentifier
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
		var_ := t.GetScope().GetVar(nameNode.GetString())
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
func (e *Es6NormalizeClasses) VisitThisAndSuper(t *NodeTraversal, n *rhino.Node) {
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
		classRecord.classNode == classNode,
		"wrong class node: %s != %s",
		classRecord.classNode,
		classNode)

	var newNameNode *rhino.Node
	if forceRewriteStaticThis {
		// For static this, we can only universally replace with the class name in a static
		// initialization context.
		newNameNode = classRecord.CreateNewNameReferenceNode().SrcrefTree(n)
	} else if e.doStaticSuperRewrites && n.IsSuper() {
		// For a static super, we can always rewrite the super as something else.
		newNameNode = classRecord.CreateNewSuperclassReferenceNode()
	}

	if newNameNode != nil {
		originalName := jutil.CheckNotNull(n.GetQualifiedName())
		newNameNode.SrcrefTree(n).SetOriginalName(originalName)
		n.ReplaceWith(newNameNode)
		t.ReportCodeChangeWithNode(newNameNode)
	}
}

// Java: Es6NormalizeClasses.isSuperCall(Node)
func Es6NormalizeClasses_IsSuperCall(n *rhino.Node) bool {
	return n.IsCall() && n.GetFirstChild().IsSuper()
}

// Es6NormalizeClasses_IsInStaticInitializationContext assumes that we have
// already filtered out non-static members and non-static blocks.
//
// Java: Es6NormalizeClasses.isInStaticInitializationContext(Node)
func Es6NormalizeClasses_IsInStaticInitializationContext(rootNode *rhino.Node) bool {
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

// Java: private static class Es6NormalizeClasses.ClassRecord
type es6NormalizeClasses_ClassRecord struct {
	// The binding identifier (a.k.a. name node) of the class itself.
	//
	// The node for `C` as in:
	//
	//   - `class C {}`
	//   - `const D = class C {};` // Inner name node in this case.
	//
	// Absent as in:
	//
	//   - `const C = class {};`
	//
	// Go-only: language-adapter — Optional<Node>: nil is absent.
	bindingIdentifier *rhino.Node

	// The node for `class` in both class declarations and class expressions.
	classNode *rhino.Node

	// The final post-normalized names for the class and extends clause.
	// These may or may not match the original names (if there were any
	// original names).
	normalizedNames *es6NormalizeClasses_NormalizedNames

	// Instance fields in the class.
	//
	// Go-only: language-adapter — the ArrayDeque queues below are slices
	// (add appends, remove takes the first element).
	instanceMembers []*rhino.Node

	// Static fields + static blocks
	staticMembers []*rhino.Node

	// computed props with side effects
	computedPropsWithSideEffects []*rhino.Node

	// Set of all the Vars defined in the constructor arguments scope and
	// constructor body scope
	constructorVars *jutil.LinkedHashSet[*Var]

	// List of all the `super()` calls within the constructor, if any.
	//
	// It's legal for classes to have more than one super() call, hence this
	// is a list rather than an `Optional<Node>`. However, public class field
	// transpilation will throw an error if this list has more than one entry
	// in a class with a public field; see
	// Es6NormalizeClasses#ensureSuperCallIsOwnStatement. We wait to emit this
	// error until after traversing the entire class body to find public
	// fields.
	superCalls []*rhino.Node
}

// Java: Es6NormalizeClasses.ClassRecord(NormalizedNames, Optional<Node>, Node)
func newEs6NormalizeClasses_ClassRecord(normalizedNames *es6NormalizeClasses_NormalizedNames, bindingIdentifier *rhino.Node, classNode *rhino.Node) *es6NormalizeClasses_ClassRecord {
	r := &es6NormalizeClasses_ClassRecord{}
	r.instanceMembers = []*rhino.Node{}
	r.staticMembers = []*rhino.Node{}
	r.computedPropsWithSideEffects = []*rhino.Node{}
	r.constructorVars = jutil.NewLinkedHashSet[*Var]()
	r.superCalls = []*rhino.Node{}
	r.normalizedNames = normalizedNames
	r.bindingIdentifier = bindingIdentifier
	r.classNode = classNode
	return r
}

// Java: Es6NormalizeClasses.ClassRecord.recordField(Node)
func (r *es6NormalizeClasses_ClassRecord) RecordField(field *rhino.Node) {
	jutil.CheckArgument(field.IsComputedFieldDef() || field.IsMemberFieldDef())
	if field.IsStaticMember() {
		r.staticMembers = append(r.staticMembers, field)
	} else {
		r.instanceMembers = append(r.instanceMembers, field)
	}
}

// Java: Es6NormalizeClasses.ClassRecord.maybeRecordComputedPropWithSideEffects(Node)
func (r *es6NormalizeClasses_ClassRecord) MaybeRecordComputedPropWithSideEffects(computedProp *rhino.Node) {
	jutil.CheckArgument(computedProp.IsComputedProp() || computedProp.IsComputedFieldDef())

	if NodeUtil_CanBeSideEffected(computedProp.GetFirstChild()) {
		r.computedPropsWithSideEffects = append(r.computedPropsWithSideEffects, computedProp)
	}
}

// Java: Es6NormalizeClasses.ClassRecord.recordStaticBlock(Node)
func (r *es6NormalizeClasses_ClassRecord) RecordStaticBlock(block *rhino.Node) {
	jutil.CheckArgument(NodeUtil_IsClassStaticBlock(block))
	r.staticMembers = append(r.staticMembers, block)
}

// Java: Es6NormalizeClasses.ClassRecord.recordConstructorScope(Scope)
func (r *es6NormalizeClasses_ClassRecord) RecordConstructorScope(s *Scope) {
	jutil.CheckArgumentMsg1(s.IsFunctionBlockScope(), "%s", s)
	jutil.CheckStateMsg1(r.constructorVars.IsEmpty(), "%s", r.constructorVars)
	argsScope := s.GetParent()
	constructorVars := jutil.NewLinkedHashSet[*Var]()
	for v := range s.GetAllSymbols() {
		constructorVars.Add(v)
	}
	for v := range argsScope.GetAllSymbols() {
		constructorVars.Add(v)
	}
	r.constructorVars = constructorVars
}

// Java: Es6NormalizeClasses.ClassRecord.createNewNameReferenceNode()
func (r *es6NormalizeClasses_ClassRecord) CreateNewNameReferenceNode() *rhino.Node {
	if r.normalizedNames.ClassNameNode().IsName() {
		// Don't cloneTree() here, because the name may have a child node, the class itself.
		return r.normalizedNames.ClassNameNode().CloneNode()
	}

	// Must cloneTree() for a qualified name.
	return r.normalizedNames.ClassNameNode().CloneTree()
}

// Java: Es6NormalizeClasses.ClassRecord.createNewSuperclassReferenceNode()
func (r *es6NormalizeClasses_ClassRecord) CreateNewSuperclassReferenceNode() *rhino.Node {
	if r.normalizedNames.SuperClassNode() == nil {
		panic(jutil.NewIllegalArgumentExceptionMsg("no superclass"))
	}
	return es6NormalizeClasses_optionalGet(r.normalizedNames.SuperClassNode()).CloneTree()
}

// es6NormalizeClasses_NormalizedNames is information about a class *after*
// it has been normalized.
//
// This is used to store information about the class that is potentially
// computed during the normalization process. The nodes are not guaranteed
// to be present in the original AST.
//
// classNameNode: the node for the qualified name of the class. "C" as in
// `class C {}`, `const C = class D {};`, `const C = class {};`; "ns.C" as in
// `const ns = {}; ns.C = class {};`.
//
// superClassNode: the name node of the super class, if present. The node
// for `E` as in `class C extends E {}`, `const D = class C extends E {};`,
// `const D = class extends E {};`.
//
// shouldExtractClass: whether the class needs to be extracted into a new
// variable. shouldExtractExtends: whether the extends clause needs to be
// extracted into a new variable.
type es6NormalizeClasses_NormalizedNames struct {
	classNameNode *rhino.Node
	// Go-only: language-adapter — Optional<Node>: nil is absent.
	superClassNode       *rhino.Node
	shouldExtractClass   bool
	shouldExtractExtends bool
}

// Java: Es6NormalizeClasses.NormalizedNames(Node, Optional<Node>, boolean, boolean)
func newEs6NormalizeClasses_NormalizedNames(
	classNameNode *rhino.Node,
	superClassNode *rhino.Node,
	shouldExtractClass bool,
	shouldExtractExtends bool) *es6NormalizeClasses_NormalizedNames {
	jutil.CheckArgumentMsg1(classNameNode.IsQualifiedName(), "%s", classNameNode)
	if !(superClassNode == nil || superClassNode.IsQualifiedName()) {
		// Go-only: language-adapter — Optional.toString() of a present value.
		jutil.CheckArgumentMsg1(false, "Optional[%s]", superClassNode)
	}
	return &es6NormalizeClasses_NormalizedNames{
		classNameNode:        classNameNode,
		superClassNode:       superClassNode,
		shouldExtractClass:   shouldExtractClass,
		shouldExtractExtends: shouldExtractExtends,
	}
}

// Java: Es6NormalizeClasses.NormalizedNames.classNameNode()
func (n *es6NormalizeClasses_NormalizedNames) ClassNameNode() *rhino.Node { return n.classNameNode }

// Java: Es6NormalizeClasses.NormalizedNames.superClassNode()
func (n *es6NormalizeClasses_NormalizedNames) SuperClassNode() *rhino.Node {
	return n.superClassNode
}

// Java: Es6NormalizeClasses.NormalizedNames.shouldExtractClass()
func (n *es6NormalizeClasses_NormalizedNames) ShouldExtractClass() bool {
	return n.shouldExtractClass
}

// Java: Es6NormalizeClasses.NormalizedNames.shouldExtractExtends()
func (n *es6NormalizeClasses_NormalizedNames) ShouldExtractExtends() bool {
	return n.shouldExtractExtends
}

// GenerateNormalizedClassNames checks whether the class and/or extends
// clause need to be normalized into a separate variable declaration. If so,
// generates a new unique name/names. Returns the NormalizedNames with the
// results. (If the class does not need normalization, then the
// NormalizedNames will match the original class name and extends clause.)
//
// Java: Es6NormalizeClasses.generateNormalizedClassNames(NodeTraversal, Node)
func (e *Es6NormalizeClasses) GenerateNormalizedClassNames(t *NodeTraversal, classNode *rhino.Node) *es6NormalizeClasses_NormalizedNames {
	extendsClause := classNode.GetSecondChild()

	shouldExtractExtends := e.ShouldExtractExtends(classNode)
	var superClass *rhino.Node
	if shouldExtractExtends {
		name := e.GenerateUniqueClassExtendsVarName(t)
		superClass = e.astFactory.CreateName(name, AstFactory_Type_(extendsClause))
	} else {
		if !extendsClause.IsEmpty() {
			superClass = extendsClause
		}
	}

	shouldExtractClass :=
		e.ShouldExtractClass(classNode) ||
			(shouldExtractExtends && !e.CanExtractExtendsSimply(classNode))
	var classNameNodeTemplate *rhino.Node
	if shouldExtractClass {
		name := e.GenerateUniqueClassDeclVarName(t)
		classNameNodeTemplate = e.astFactory.CreateConstantName(name, AstFactory_Type_(classNode))
	} else {
		classNameNodeTemplate = NodeUtil_GetNameNode(classNode)
	}

	return newEs6NormalizeClasses_NormalizedNames(
		classNameNodeTemplate, superClass, shouldExtractClass, shouldExtractExtends)
}
