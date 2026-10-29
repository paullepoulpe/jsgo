// Port of com/google/javascript/jscomp/parsing/IRFactory.java @ closure-compiler v20261021

/*
 * Copyright 2013 The Closure Compiler Authors.
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

package parsing

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var IRFactory_GOOG_MODULE = rhino.QualifiedName_Of("goog.module")

// Java: IRFactory.TransformDispatcher.processAstRoot(ProgramTree)
func (d *IRFactory_TransformDispatcher) ProcessAstRoot(rootNode *parser.ProgramTree) rhino.Node {
	scriptNode := d.NewNode(rhino.Token_SCRIPT)
	for _, child := range rootNode.SourceElements {
		scriptNode.AddChildToBack(d.Transform(child))
	}
	d.ParseDirectives(scriptNode)
	isGoogModule := d.IsGoogModuleFile(scriptNode)
	if isGoogModule || d.features.Has(parser.FeatureSet_Feature_MODULES) {
		moduleNode := d.NewNode(rhino.Token_MODULE_BODY)
		d.SetSourceInfo(moduleNode, rootNode)
		moduleNode.AddChildrenToBack(scriptNode.RemoveChildren())
		scriptNode.AddChildToBack(moduleNode)
		if isGoogModule {
			scriptNode.PutBooleanProp(rhino.Node_GOOG_MODULE, true)
		} else {
			scriptNode.PutBooleanProp(rhino.Node_ES6_MODULE, true)
		}
	}
	return scriptNode
}

// Java: IRFactory.TransformDispatcher.isGoogModuleFile(Node)
func (d *IRFactory_TransformDispatcher) IsGoogModuleFile(scriptNode rhino.Node) bool {
	jutil.CheckArgument(scriptNode.IsScript())
	if !scriptNode.HasChildren() {
		return false
	}
	exprResult := scriptNode.GetFirstChild()
	if !exprResult.IsExprResult() {
		return false
	}
	call := exprResult.GetFirstChild()
	if !call.IsCall() {
		return false
	}
	return IRFactory_GOOG_MODULE.Matches(call.GetFirstChild()) &&
		call.HasTwoChildren() &&
		call.GetSecondChild().IsStringLit()
}

// ParseDirectives parses the directives, encodes them in the AST, and
// removes their nodes.
//
// The only suported directive is "use strict".
//
// For information on ES5 directives, see section 14.1 of ECMA-262, Edition 5.
//
// Java: IRFactory.TransformDispatcher.parseDirectives(Node)
func (d *IRFactory_TransformDispatcher) ParseDirectives(node rhino.Node) {
	useStrict := false
	for {
		statement := node.GetFirstChild()
		if statement.IsNil() || !statement.IsExprResult() {
			break
		}

		directive := statement.GetFirstChild()
		if !directive.IsStringLit() || directive.GetString() != "use strict" {
			break
		}

		useStrict = true
		// TODO(nickreid): Should we also remove other directives that aren't
		// "use strict"?
		statement.Detach()
	}

	if useStrict {
		node.SetUseStrict(true)
	}
}

// ProcessField processes class fields: general case.
//
// Java: IRFactory.TransformDispatcher.processField(FieldDeclarationTree)
func (d *IRFactory_TransformDispatcher) ProcessField(tree *parser.FieldDeclarationTree) rhino.Node {
	if tree.Name.IsPrivateIdentifier() {
		d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_PRIVATE_ELEMENTS)
	} else {
		d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_PUBLIC_CLASS_FIELDS)
	}

	node := d.MaybeAddNonJsDocComment(
		d.NewStringNodeFromIdentifier(
			rhino.Token_MEMBER_FIELD_DEF, IRFactory_TransformDispatcher_IdentifierType_CAN_BE_PRIVATE, tree.Name),
		tree.GetStart())
	if tree.Initializer != nil {
		initializer := d.Transform(tree.Initializer)
		node.AddChildToBack(initializer)
		d.SetLength(node, tree.Location.Start, tree.Location.End)
	}

	node.PutBooleanProp(rhino.Node_STATIC_MEMBER, tree.IsStatic)

	return node
}

// ProcessComputedPropertyField processes class fields: computed property
// fields.
//
// Java: IRFactory.TransformDispatcher.processComputedPropertyField(ComputedPropertyFieldTree)
func (d *IRFactory_TransformDispatcher) ProcessComputedPropertyField(tree *parser.ComputedPropertyFieldTree) rhino.Node {
	d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_PUBLIC_CLASS_FIELDS)

	key := d.Transform(tree.Property)

	var node rhino.Node
	if tree.Initializer != nil {
		node = d.NewNodeWith2Children(rhino.Token_COMPUTED_FIELD_DEF, key, d.Transform(tree.Initializer))
	} else {
		node = d.NewNodeWithChild(rhino.Token_COMPUTED_FIELD_DEF, key)
	}

	node.PutBooleanProp(rhino.Node_STATIC_MEMBER, tree.IsStatic)

	return node
}

// Java: IRFactory.TransformDispatcher.processFormalParameterList(FormalParameterListTree)
func (d *IRFactory_TransformDispatcher) ProcessFormalParameterList(tree *parser.FormalParameterListTree) rhino.Node {
	params := d.NewNode(rhino.Token_PARAM_LIST)
	params.SetTrailingComma(tree.HasTrailingComma)
	if !d.CheckParameters(tree) {
		return params
	}

	parameters := tree.Parameters
	argCount := 0

	for _, param := range parameters {
		var paramNode rhino.Node
		switch param.GetParseTreeBase().Type {
		case parser.ParseTreeType_DEFAULT_PARAMETER:
			// processDefaultParameter() knows how to find and apply inline
			// JSDoc to the right node
			paramNode = d.ProcessDefaultParameter(param.(*parser.DefaultParameterTree))
		case parser.ParseTreeType_ITER_REST:
			d.MaybeWarnForFeature(param, parser.FeatureSet_Feature_REST_PARAMETERS)
			paramNode = d.TransformNodeWithInlineComments(param)
		default:
			paramNode = d.TransformNodeWithInlineComments(param)
			// Reusing the logic to attach trailing comments used from
			// call-site argsList
			d.attachPossibleTrailingCommentsForArgAt(
				paramNode, parameters, argCount, tree.CommaPositions, tree.Location.End)
		}

		// Children must be simple names, default parameters, rest
		// parameters, or destructuring patterns.
		jutil.CheckState(
			paramNode.IsName() ||
				paramNode.IsRest() ||
				paramNode.IsArrayPattern() ||
				paramNode.IsObjectPattern() ||
				paramNode.IsDefaultValue())
		params.AddChildToBack(paramNode)
		argCount++
	}

	return params
}

// Java: IRFactory.TransformDispatcher.processDefaultParameter(DefaultParameterTree)
func (d *IRFactory_TransformDispatcher) ProcessDefaultParameter(tree *parser.DefaultParameterTree) rhino.Node {
	d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_DEFAULT_PARAMETERS)
	targetTree := tree.Lhs
	var targetNode rhino.Node
	if targetTree.GetParseTreeBase().Type == parser.ParseTreeType_IDENTIFIER_EXPRESSION {
		// allow inline JSDoc on an identifier
		// let { /** inlineType */ x = defaultValue } = someObject;
		// TODO(bradfordcsmith): Do we need to allow inline JSDoc for
		// qualified names, too?
		targetNode = d.ProcessNameWithInlineComments(
			targetTree.(*parser.IdentifierExpressionTree), IRFactory_TransformDispatcher_IdentifierType_STANDARD)
	} else {
		// ({prop: /** string */ ns.a.b = 'foo'} = someObject);
		// NOTE: CheckJSDoc will report an error for this case, since we want
		// qualified names to be declared with individual statements, like
		// `/** @type {string} */ ns.a.b;`
		targetNode = d.TransformNodeWithInlineComments(targetTree)
	}
	defaultValueExpression := d.Transform(tree.DefaultValue)
	defaultValueNode := d.NewNodeWith2Children(rhino.Token_DEFAULT_VALUE, targetNode, defaultValueExpression)
	d.ReportErrorIfYieldOrAwaitInDefaultValue(defaultValueNode)
	d.SetSourceInfo(defaultValueNode, tree)
	return defaultValueNode
}

// Java: IRFactory.TransformDispatcher.processClassDeclaration(ClassDeclarationTree)
func (d *IRFactory_TransformDispatcher) ProcessClassDeclaration(tree *parser.ClassDeclarationTree) rhino.Node {
	d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_CLASSES)

	name := d.TransformOrEmptyIdentifier(tree.Name, tree)
	superClass := d.TransformOrEmpty(tree.SuperClass, tree)

	classMembers := d.NewNode(rhino.Token_CLASS_MEMBERS)
	d.SetSourceInfo(classMembers, tree)

	func() {
		ignored := d.classScope.Increment()
		defer ignored.Close()
		hasConstructor := false
		for _, child := range tree.Elements {
			switch child.GetParseTreeBase().Type {
			case parser.ParseTreeType_COMPUTED_PROPERTY_GETTER, parser.ParseTreeType_COMPUTED_PROPERTY_SETTER,
				parser.ParseTreeType_GET_ACCESSOR, parser.ParseTreeType_SET_ACCESSOR:
				d.features = d.features.With(parser.FeatureSet_Feature_CLASS_GETTER_SETTER)
			case parser.ParseTreeType_BLOCK:
				d.features = d.features.With(parser.FeatureSet_Feature_CLASS_STATIC_BLOCK)
			default:
			}

			childIsCtor := d.ValidateClassConstructorMember(child) // Has side-effects.
			if childIsCtor {
				if hasConstructor {
					d.errorReporter.Error(
						"Class may have only one constructor.", //
						d.sourceName,
						d.Lineno(child),
						d.Charno(child))
				}
				hasConstructor = true
			}

			classMembers.AddChildToBack(d.Transform(child))
		}
	}()

	// When nested inside of multiple classes, only run the validation once
	// the top-most class exits scope.
	if !d.classScope.InScope() {
		d.ValidatePrivatePropertyUsage(d.GetPrivatePropsAndValidateUniqueness(classMembers), classMembers)
	}

	classNode := d.NewNodeWith3Children(rhino.Token_CLASS, name, superClass, classMembers)
	d.AttachPossibleTrailingComment(classNode, tree.GetEnd())

	return classNode
}

const (
	privStatic uint8 = 1 << iota
	privGetter
	privSetter
	privOther
)

type privatePropScope struct {
	props  map[string]uint8
	parent *privatePropScope
}

func (s *privatePropScope) contains(name string) bool {
	for cur := s; cur != nil; cur = cur.parent {
		if cur.props[name] != 0 {
			return true
		}
	}
	return false
}

// ValidatePrivatePropertyUsage validates that all referenced private
// properties within the scope of classMembers are found in privatePropNames
// and private fields are not deleted (both of these cases are syntax
// errors).
//
// Recurses into inner classes and successively chains privatePropNames scopes
// such that outer class private properties are available to inner classes.
//
// Java: IRFactory.TransformDispatcher.validatePrivatePropertyUsage(Set<String>, Node)
func (d *IRFactory_TransformDispatcher) ValidatePrivatePropertyUsage(privatePropNames map[string]uint8, classMembers rhino.Node) {
	d.validatePrivatePropertyUsage(&privatePropScope{props: privatePropNames}, classMembers)
}

func (d *IRFactory_TransformDispatcher) validatePrivatePropertyUsage(scope *privatePropScope, classMembers rhino.Node) {
	jutil.CheckState(classMembers.IsClassMembers())
	if !classMembers.HasChildren() {
		return
	}

	// Uses a stack object instead of recursion to avoid a stack overflow.
	var visitStack []rhino.Node
	visitStack = append(visitStack, classMembers.GetFirstChild())
	for len(visitStack) != 0 {
		node := visitStack[len(visitStack)-1]
		visitStack = visitStack[:len(visitStack)-1]

		// Check that these usages reference a private property declared in
		// the class:
		//   this.#privateProp
		//   x.#privateProp
		//   x?.#privateProp
		//   #privateProp in x
		if node.IsPrivateIdentifier() &&
			(node.IsGetProp() ||
				node.IsOptChainGetProp() ||
				(node.IsName() && node.HasParent() && node.GetParent().IsIn())) &&
			!scope.contains(node.GetString()) {
			propType := "fields"
			if node.GetParent().IsCall() {
				propType = "methods"
			}
			d.errorReporter.Error(
				"Private "+propType+" must be declared in an enclosing class",
				d.sourceName,
				node.GetLineno(),
				node.GetCharno())
		}

		// Checks that private fields are not deleted:
		//   delete this.#privateProp
		//   delete x.#privateProp
		//   delete x?.#privateProp
		if node.IsDelProp() {
			firstChild := rhino.CheckNotNull(node.GetFirstChild())
			if firstChild.IsPrivateIdentifier() &&
				(firstChild.IsGetProp() || firstChild.IsOptChainGetProp()) {
				d.errorReporter.Error(
					"Private fields cannot be deleted", d.sourceName, node.GetLineno(), node.GetCharno())
			}
		}

		if node.IsClassMembers() {
			// Recurse into the inner class with the outer and inner class
			// private properties.
			innerProps := d.GetPrivatePropsAndValidateUniqueness(node)
			d.validatePrivatePropertyUsage(&privatePropScope{props: innerProps, parent: scope}, node)
		} else {
			if !node.GetNext().IsNil() {
				visitStack = append(visitStack, node.GetNext())
			}
			if node.HasChildren() {
				visitStack = append(visitStack, node.GetFirstChild())
			}
		}
	}
}

// GetPrivatePropsAndValidateUniqueness returns the map of private properties
// declared in the class corresponding to classMembers.
//
// While gathering private properties, validates that there are no duplicate
// property names. Note that, for private getters and setters, one getter and
// one setter can have the same name provided they have the same static or
// non-static modifier.
//
// Java: IRFactory.TransformDispatcher.getPrivatePropsAndValidateUniqueness(Node)
func (d *IRFactory_TransformDispatcher) GetPrivatePropsAndValidateUniqueness(classMembers rhino.Node) map[string]uint8 {
	var props map[string]uint8
	for curNode := classMembers.GetFirstChild(); !curNode.IsNil(); curNode = curNode.GetNext() {
		if !curNode.IsPrivateIdentifier() {
			continue
		}

		propName := curNode.GetString()
		jutil.CheckStateMsg2(
			curNode.IsGetterDef() ||
				curNode.IsSetterDef() ||
				curNode.IsMemberFieldDef() ||
				curNode.IsMemberFunctionDef(),
			"Private property '%s' has an unsupported token type: %s",
			propName,
			curNode.GetToken())

		if props == nil {
			props = make(map[string]uint8)
		}
		prev := props[propName]
		var cur uint8
		switch curNode.GetToken() {
		case rhino.Token_GETTER_DEF:
			cur = privGetter
		case rhino.Token_SETTER_DEF:
			cur = privSetter
		case rhino.Token_MEMBER_FIELD_DEF, rhino.Token_MEMBER_FUNCTION_DEF:
			cur = privOther
		default:
		}
		if curNode.IsStaticMember() {
			cur |= privStatic
		}
		alreadyDeclared := false
		if prev != 0 {
			staticMismatch := (prev & privStatic) != (cur & privStatic)
			roleConflict := prev&privOther != 0 || cur&privOther != 0 || prev&cur&(privGetter|privSetter) != 0
			if staticMismatch || roleConflict {
				alreadyDeclared = true
			}
			if staticMismatch {
				cur |= privOther
			}
		}
		if alreadyDeclared {
			d.errorReporter.Error(
				"Identifier '"+propName+"' has already been declared",
				d.sourceName,
				curNode.GetLineno(),
				curNode.GetCharno())
		}
		props[propName] = prev | cur
	}
	return props
}

// ValidateClassConstructorMember returns true iff this member is a legal
// class constructor.
//
// Java: IRFactory.TransformDispatcher.validateClassConstructorMember(ParseTree)
func (d *IRFactory_TransformDispatcher) ValidateClassConstructorMember(member parser.ParseTree) bool {
	var memberName parser.Token
	var isStatic bool
	var hasIllegalModifier bool
	switch member.GetParseTreeBase().Type {
	case parser.ParseTreeType_GET_ACCESSOR:
		getter := member.(*parser.GetAccessorTree)
		memberName = getter.PropertyName
		isStatic = getter.IsStatic
		hasIllegalModifier = true
	case parser.ParseTreeType_SET_ACCESSOR:
		setter := member.(*parser.SetAccessorTree)
		memberName = setter.PropertyName
		isStatic = setter.IsStatic
		hasIllegalModifier = true
	case parser.ParseTreeType_FUNCTION_DECLARATION:
		method := member.(*parser.FunctionDeclarationTree)
		memberName = method.Name
		isStatic = method.IsStatic
		hasIllegalModifier = method.IsGenerator || method.IsAsync
	default:
		// Computed properties aren't an issue here because they aren't used
		// as the class constructor, regardless of their name.
		return false
	}

	if isStatic {
		// Statics are fine because they're never the class constructor.
		return false
	}

	checkTokenNotNull(memberName)
	if memberName.Type() != parser.TokenType_IDENTIFIER ||
		!memberName.ValueEquals("constructor") {
		// There's only a potential issue if the member is named
		// "constructor".
		// TODO(b/123769080): Also check for quoted string literals with the
		// value "constructor".
		return false
	}

	if hasIllegalModifier {
		d.errorReporter.Error(
			"Class constructor may not be getter, setter, async, or generator.",
			d.sourceName,
			d.Lineno(member),
			d.Charno(member))
		return false
	}

	return true
}

// Java: IRFactory.TransformDispatcher.processExportDecl(ExportDeclarationTree)
func (d *IRFactory_TransformDispatcher) ProcessExportDecl(tree *parser.ExportDeclarationTree) rhino.Node {
	d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_MODULES)
	var decls rhino.Node
	if tree.IsExportAll {
		jutil.CheckState(tree.Declaration == nil && tree.ExportSpecifierList == nil)
	} else if tree.Declaration != nil {
		jutil.CheckState(tree.ExportSpecifierList == nil)
		decls = d.Transform(tree.Declaration)
	} else {
		decls = d.TransformList(rhino.Token_EXPORT_SPECS, tree.ExportSpecifierList)
	}
	if decls.IsNil() {
		decls = d.NewNode(rhino.Token_EMPTY)
	}
	d.SetSourceInfo(decls, tree)
	export := d.NewNodeWithChild(rhino.Token_EXPORT, decls)
	if tree.From.IsPresent() {
		from := d.ProcessString(tree.From)
		export.AddChildToBack(from)
	}

	export.PutBooleanProp(rhino.Node_EXPORT_ALL_FROM, tree.IsExportAll)
	export.PutBooleanProp(rhino.Node_EXPORT_DEFAULT, tree.IsDefault)
	return export
}

// Java: IRFactory.TransformDispatcher.processExportSpec(ExportSpecifierTree)
func (d *IRFactory_TransformDispatcher) ProcessExportSpec(tree *parser.ExportSpecifierTree) rhino.Node {
	importedName := d.ProcessNameOfToken(tree.ImportedName, rhino.Token_NAME, IRFactory_TransformDispatcher_IdentifierType_STANDARD)
	exportSpec := d.NewNodeWithChild(rhino.Token_EXPORT_SPEC, importedName)
	if tree.DestinationName.IsNull() {
		exportSpec.SetShorthandProperty(true)
		exportSpec.AddChildToBack(importedName.CloneTree())
	} else {
		destinationName := d.ProcessNameOfToken(tree.DestinationName, rhino.Token_NAME, IRFactory_TransformDispatcher_IdentifierType_STANDARD)
		exportSpec.AddChildToBack(destinationName)
	}
	return exportSpec
}

// Java: IRFactory.TransformDispatcher.processImportDecl(ImportDeclarationTree)
func (d *IRFactory_TransformDispatcher) ProcessImportDecl(tree *parser.ImportDeclarationTree) rhino.Node {
	d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_MODULES)

	firstChild := d.TransformOrEmptyIdentifier(tree.DefaultBindingIdentifier, tree)
	var secondChild rhino.Node
	if tree.NameSpaceImportIdentifier.IsNull() {
		secondChild = d.TransformListOrEmpty(rhino.Token_IMPORT_SPECS, tree.ImportSpecifierList)
		// Currently source info is "import {foo} from '...';" expression. If
		// needed this should be changed to use only "{foo}" part.
		d.SetSourceInfo(secondChild, tree)
	} else {
		secondChild = d.NewStringNodeFromIdentifier(
			rhino.Token_IMPORT_STAR, IRFactory_TransformDispatcher_IdentifierType_STANDARD, tree.NameSpaceImportIdentifier)
		d.SetSourceInfoFromToken(secondChild, tree.NameSpaceImportIdentifier)
	}
	thirdChild := d.ProcessString(tree.ModuleSpecifier)

	return d.NewNodeWith3Children(rhino.Token_IMPORT, firstChild, secondChild, thirdChild)
}

// Java: IRFactory.TransformDispatcher.processImportSpec(ImportSpecifierTree)
func (d *IRFactory_TransformDispatcher) ProcessImportSpec(tree *parser.ImportSpecifierTree) rhino.Node {
	importedName := d.ProcessNameOfToken(tree.ImportedName, rhino.Token_NAME, IRFactory_TransformDispatcher_IdentifierType_STANDARD)
	importSpec := d.NewNodeWithChild(rhino.Token_IMPORT_SPEC, importedName)
	if tree.DestinationName.IsNull() {
		importSpec.SetShorthandProperty(true)
		importSpec.AddChildToBack(importedName.CloneTree())
	} else {
		importSpec.AddChildToBack(
			d.ProcessNameOfToken(tree.DestinationName, rhino.Token_NAME, IRFactory_TransformDispatcher_IdentifierType_STANDARD))
	}
	return importSpec
}

// Java: IRFactory.TransformDispatcher.checkParameters(FormalParameterListTree)
func (d *IRFactory_TransformDispatcher) CheckParameters(tree *parser.FormalParameterListTree) bool {
	params := tree.Parameters
	good := true
	for i := 0; i < len(params); i++ {
		param := params[i]
		if param.GetParseTreeBase().Type == parser.ParseTreeType_ITER_REST {
			if i != len(params)-1 {
				d.errorReporter.Error(
					"A rest parameter must be last in a parameter list.",
					d.sourceName,
					d.Lineno(param),
					d.Charno(param))
				good = false
			} else if tree.HasTrailingComma {
				d.errorReporter.Error(
					"A trailing comma must not follow a rest parameter.",
					d.sourceName,
					d.Lineno(param),
					d.Charno(param))
				good = false
			}
		}
	}
	return good
}

// Java: IRFactory.TransformDispatcher.transformList(Token, ImmutableList<ParseTree>)
func (d *IRFactory_TransformDispatcher) TransformList(type_ rhino.Token, list []parser.ParseTree) rhino.Node {
	n := d.NewNode(type_)
	// Go-only: language-adapter — a nil slice is Java null (the parser keeps
	// empty lists non-nil), and an enhanced for statement over a null list
	// throws NullPointerException where a Go range does nothing. Of the lists
	// this file iterates, only ExportDeclarationTree.exportSpecifierList
	// reaches here as null (an export the parser recovered from in
	// KEEP_GOING mode).
	if list == nil {
		panic(jutil.NewNullPointerException())
	}
	for _, tree := range list {
		n.AddChildToBack(d.Transform(tree))
	}
	return n
}

// Java: IRFactory.TransformDispatcher.transformListOrEmpty(Token, ImmutableList<ParseTree>)
func (d *IRFactory_TransformDispatcher) TransformListOrEmpty(type_ rhino.Token, list []parser.ParseTree) rhino.Node {
	if list == nil || len(list) == 0 {
		return d.NewNode(rhino.Token_EMPTY)
	} else {
		return d.TransformList(type_, list)
	}
}
