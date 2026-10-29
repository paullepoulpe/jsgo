// Port of com/google/javascript/jscomp/AstFactory.java @ closure-compiler v20261021

/*
 * Copyright 2018 The Closure Compiler Authors.
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
	"sync"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// AstFactory creates AST nodes and subtrees.
//
// This class supports creating nodes either with or without type
// information.
//
// The idea is that client code can create the trees of nodes it needs
// without having to contain logic for deciding whether type information
// should be added or not, and only minimal logic for determining which types
// to add when they are necessary. Most methods in this class are able to
// determine the correct type information from already existing AST nodes and
// the current scope.
//
// AstFactory supports both Closure types (see JSType) and optimization-only
// types (see Color)s. Colors contain less information than JSTypes which
// puts some restrictions on the amount of inference this class can do. For
// example, there's no way to ask "What is the color of property 'x' on
// receiver color 'obj'". This is why many methods accept a StaticScope
// instead of a Scope: you may pass in a GlobalNamespace or similar object
// which contains fully qualified names, to look up colors for an entire
// property chain.
//
// IMPORTANT: The methods in this class should never set source reference
// information. It is the responsibility of the client code to set the
// correct source reference information. If we were to make guesses here,
// that would lead to client code that sometimes does and sometimes doesn't
// set the source reference and force the reader of that code to determine
// whether the default guess we have here is really correct or not. It's
// better to have the decision made explicitly in the client code.
//
// TODO(b/193800507): delete the methods in this class that only work for
// JSTypes but not colors.
type AstFactory struct {
	colorRegistry *colors.ColorRegistry
	registry      *rhino.JSTypeRegistry
	// We need the unknown type so frequently, it's worth caching it.
	unknownType         rhino.JSType
	typeMode            AstFactory_TypeMode
	lifeCycleStage      AbstractCompiler_LifeCycleStage
	runtimeJsLibManager *RuntimeJsLibManager
}

// Java: AstFactory.bigintNumberStringColor (Suppliers.memoize)
var AstFactory_bigintNumberStringColor = sync.OnceValue(func() *colors.Color {
	return colors.Color_CreateUnionSlice(
		[]*colors.Color{
			colors.StandardColors_BIGINT, colors.StandardColors_STRING, colors.StandardColors_NUMBER})
})

// AstFactory_TypeMode is AstFactory.TypeMode.
type AstFactory_TypeMode int32

const (
	AstFactory_TypeMode_JSTYPE AstFactory_TypeMode = iota
	AstFactory_TypeMode_COLOR
	AstFactory_TypeMode_NONE
)

var astFactory_TypeMode_names = [...]string{"JSTYPE", "COLOR", "NONE"}

// Java: AstFactory.TypeMode.name()
func (m AstFactory_TypeMode) String() string { return astFactory_TypeMode_names[m] }

// Java: AstFactory.TypeMode.ordinal()
func (m AstFactory_TypeMode) Ordinal() int { return int(m) }

// Java: AstFactory(LifeCycleStage, JSTypeRegistry, RuntimeJsLibManager)
func NewAstFactoryWithRegistry(
	lifeCycleStage AbstractCompiler_LifeCycleStage,
	registry *rhino.JSTypeRegistry,
	runtimeJsLibManager *RuntimeJsLibManager) *AstFactory {
	a := &AstFactory{}
	a.lifeCycleStage = lifeCycleStage
	a.registry = registry
	a.runtimeJsLibManager = runtimeJsLibManager
	a.colorRegistry = nil
	a.unknownType = a.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
	a.typeMode = AstFactory_TypeMode_JSTYPE
	return a
}

// Java: AstFactory(LifeCycleStage, RuntimeJsLibManager)
func NewAstFactory(
	lifeCycleStage AbstractCompiler_LifeCycleStage,
	runtimeJsLibManager *RuntimeJsLibManager) *AstFactory {
	a := &AstFactory{}
	a.lifeCycleStage = lifeCycleStage
	a.registry = nil
	a.runtimeJsLibManager = runtimeJsLibManager
	a.colorRegistry = nil
	a.unknownType = nil
	a.typeMode = AstFactory_TypeMode_NONE
	return a
}

// Java: AstFactory(LifeCycleStage, ColorRegistry, RuntimeJsLibManager)
func NewAstFactoryWithColorRegistry(
	lifeCycleStage AbstractCompiler_LifeCycleStage,
	colorRegistry *colors.ColorRegistry,
	runtimeJsLibManager *RuntimeJsLibManager) *AstFactory {
	a := &AstFactory{}
	a.lifeCycleStage = lifeCycleStage
	a.registry = nil
	a.runtimeJsLibManager = runtimeJsLibManager
	a.colorRegistry = colorRegistry
	a.unknownType = nil
	a.typeMode = AstFactory_TypeMode_COLOR
	return a
}

// Java: AstFactory.createFactoryWithoutTypes(LifeCycleStage, RuntimeJsLibManager)
func AstFactory_CreateFactoryWithoutTypes(
	lifeCycleStage AbstractCompiler_LifeCycleStage,
	runtimeJsLibManager *RuntimeJsLibManager) *AstFactory {
	return NewAstFactory(lifeCycleStage, runtimeJsLibManager)
}

// Java: AstFactory.createFactoryWithTypes(LifeCycleStage, JSTypeRegistry, RuntimeJsLibManager)
func AstFactory_CreateFactoryWithTypes(
	lifeCycleStage AbstractCompiler_LifeCycleStage,
	registry *rhino.JSTypeRegistry,
	runtimeJsLibManager *RuntimeJsLibManager) *AstFactory {
	return NewAstFactoryWithRegistry(lifeCycleStage, registry, runtimeJsLibManager)
}

// Java: AstFactory.createFactoryWithColors(LifeCycleStage, ColorRegistry, RuntimeJsLibManager)
func AstFactory_CreateFactoryWithColors(
	lifeCycleStage AbstractCompiler_LifeCycleStage,
	colorRegistry *colors.ColorRegistry,
	runtimeJsLibManager *RuntimeJsLibManager) *AstFactory {
	return NewAstFactoryWithColorRegistry(lifeCycleStage, colorRegistry, runtimeJsLibManager)
}

// IsAddingTypes answers: does this class instance add types to the nodes it
// creates?
//
// Java: AstFactory.isAddingTypes()
func (a *AstFactory) IsAddingTypes() bool {
	return AstFactory_TypeMode_JSTYPE == a.typeMode
}

// IsAddingColors answers: does this class instance add optimization colors
// to the nodes it creates?
//
// Java: AstFactory.isAddingColors()
func (a *AstFactory) IsAddingColors() bool {
	return AstFactory_TypeMode_COLOR == a.typeMode
}

// Java: AstFactory.assertNotAddingColors()
func (a *AstFactory) AssertNotAddingColors() {
	jutil.CheckStateMsg(!a.IsAddingColors(), "method not supported for colors")
}

// Java: AstFactory.createString(String)
func (a *AstFactory) CreateString(value string) rhino.Node {
	result := rhino.IR_String(value)
	a.SetJSTypeOrColor(
		AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_STRING_TYPE, colors.StandardColors_STRING), result)
	return result
}

// Java: AstFactory.createNumber(double)
func (a *AstFactory) CreateNumber(value float64) rhino.Node {
	result := rhino.IR_Number(value)
	a.SetJSTypeOrColor(
		AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_NUMBER_TYPE, colors.StandardColors_NUMBER), result)
	return result
}

// Java: AstFactory.createBoolean(boolean)
func (a *AstFactory) CreateBoolean(value bool) rhino.Node {
	var result rhino.Node
	if value {
		result = rhino.IR_TrueNode()
	} else {
		result = rhino.IR_FalseNode()
	}
	a.SetJSTypeOrColor(
		AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_BOOLEAN_TYPE, colors.StandardColors_BOOLEAN), result)
	return result
}

// Java: AstFactory.createNull()
func (a *AstFactory) CreateNull() rhino.Node {
	result := rhino.IR_NullNode()
	a.SetJSTypeOrColor(
		AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_NULL_TYPE, colors.StandardColors_NULL_OR_VOID), result)
	return result
}

// Java: AstFactory.createVoid(Node)
func (a *AstFactory) CreateVoid(child rhino.Node) rhino.Node {
	result := rhino.IR_VoidNode(child)
	a.SetJSTypeOrColor(
		AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_VOID_TYPE, colors.StandardColors_NULL_OR_VOID), result)
	return result
}

// CreateUndefinedValue returns a new Node representing the undefined value.
//
// Java: AstFactory.createUndefinedValue()
func (a *AstFactory) CreateUndefinedValue() rhino.Node {
	// We prefer `void 0` as being shorter than `undefined`.
	// Also, it's technically possible for malicious code to assign a value to `undefined`.
	return a.CreateVoid(a.CreateNumber(0))
}

// Java: AstFactory.createThis(AstFactory.Type)
func (a *AstFactory) CreateThis(thisType AstFactory_Type) rhino.Node {
	result := rhino.IR_ThisNode()
	a.SetJSTypeOrColor(thisType, result)
	return result
}

// Java: AstFactory.createSuper(AstFactory.Type)
func (a *AstFactory) CreateSuper(superType AstFactory_Type) rhino.Node {
	result := rhino.IR_SuperNode()
	a.SetJSTypeOrColor(superType, result)
	return result
}

// Java: AstFactory.createNewTarget(AstFactory.Type)
func (a *AstFactory) CreateNewTarget(newTargetType AstFactory_Type) rhino.Node {
	result := rhino.IR_NewTarget()
	a.SetJSTypeOrColor(newTargetType, result)
	return result
}

// CreateThisForEs6Class creates a THIS node with the correct type for the
// given ES6 class node.
//
// With the optimization colors type system, we can support inferring the
// type of this for constructors but not generic functions annotated @this
//
// Java: AstFactory.createThisForEs6Class(Node)
func (a *AstFactory) CreateThisForEs6Class(functionNode rhino.Node) rhino.Node {
	jutil.CheckStateMsg1(functionNode.IsClass(), "%s", functionNode)
	result := rhino.IR_ThisNode()
	a.SetJSTypeOrColor(a.GetTypeOfThisForEs6Class(functionNode), result)
	return result
}

// CreateThisForEs6ClassMember creates a THIS node with the correct type for
// the given ES6 class node.
//
// With the optimization colors type system, we can support inferring the
// type of this for constructors but not generic functions annotated @this
//
// Java: AstFactory.createThisForEs6ClassMember(Node)
func (a *AstFactory) CreateThisForEs6ClassMember(memberNode rhino.Node) rhino.Node {
	jutil.CheckArgument(memberNode.GetParent().IsClassMembers())
	jutil.CheckArgument(
		memberNode.IsMemberFunctionDef() ||
			memberNode.IsMemberFieldDef() ||
			memberNode.IsComputedFieldDef())
	classNode := memberNode.GetGrandparent()
	if memberNode.IsStaticMember() {
		result := rhino.IR_ThisNode()
		a.SetJSTypeOrColor(AstFactory_Type_(classNode), result)
		return result
	} else {
		return a.CreateThisForEs6Class(classNode)
	}
}

// Java: AstFactory.getTypeOfThisForFunctionNode(Node)
func (a *AstFactory) GetTypeOfThisForFunctionNode(functionNode rhino.Node) rhino.JSType {
	a.AssertNotAddingColors()
	if a.IsAddingTypes() {
		functionType := a.GetFunctionType(functionNode)
		return jutil.CheckNotNullMsg1(functionType.GetTypeOfThis(), "%s", functionType)
	} else {
		return nil // not adding type information
	}
}

// Java: AstFactory.getTypeOfThisForEs6Class(Node)
func (a *AstFactory) GetTypeOfThisForEs6Class(functionNode rhino.Node) AstFactory_Type {
	// Go-only: language-adapter — the JVM's "helpful" NullPointerException
	// message (jutil.NewHelpfulNullPointerException). An async method of an
	// object literal that calls super has no enclosing class, so
	// RewriteAsyncFunctions passes null here and Java dereferences it; the
	// INTERNAL COMPILER ERROR report prints the message the JVM derives from
	// this call. A Go nil dereference has none.
	if functionNode.IsNil() {
		panic(jutil.NewHelpfulNullPointerException(`Cannot invoke "com.google.javascript.rhino.Node.isClass()" because "functionNode" is null`))
	}
	jutil.CheckArgumentMsg1(functionNode.IsClass(), "%s", functionNode)
	switch a.typeMode {
	case AstFactory_TypeMode_JSTYPE:
		return AstFactory_Type_OfJSType(a.GetTypeOfThisForFunctionNode(functionNode))
	case AstFactory_TypeMode_COLOR:
		return AstFactory_Type_OfColor(a.GetInstanceOfColor(functionNode.GetColor()))
	case AstFactory_TypeMode_NONE:
		return AstFactory_NoTypeInformation()
	}
	// Go-only: language-adapter — unreachable for a valid enum value: the
	// Java switch expression is exhaustive.
	panic(jutil.NewIllegalStateException())
}

// Java: AstFactory.getFunctionType(Node)
func (a *AstFactory) GetFunctionType(functionNode rhino.Node) rhino.FunctionType {
	jutil.CheckStateMsg1(
		functionNode.IsFunction() || functionNode.IsClass(),
		"not a function or class: %s",
		functionNode)
	a.AssertNotAddingColors()
	// If the function declaration was cast to a different type, we want the original type
	// from before the cast.
	typeBeforeCast := functionNode.GetJSTypeBeforeCast()
	var functionType rhino.FunctionType
	if typeBeforeCast != nil {
		functionType = typeBeforeCast.AssertFunctionType()
	} else {
		functionType = functionNode.GetJSTypeRequired().AssertFunctionType()
	}
	return functionType
}

// CreateThisAliasReferenceForEs6Class creates a NAME node having the type of
// "this" appropriate for the given ES6 class node
//
// With the optimization colors type system, we can support inferring the
// type of this for classes but not generic functions annotated @this.
//
// Java: AstFactory.createThisAliasReferenceForEs6Class(String, Node)
func (a *AstFactory) CreateThisAliasReferenceForEs6Class(aliasName string, functionNode rhino.Node) rhino.Node {
	return a.CreateName(aliasName, a.GetTypeOfThisForEs6Class(functionNode))
}

// CreateArgumentsReference creates a reference to "arguments" with the type
// specified in externs, or unknown if the externs for it weren't included.
//
// Java: AstFactory.createArgumentsReference()
func (a *AstFactory) CreateArgumentsReference() rhino.Node {
	result := rhino.IR_Name("arguments")
	switch a.typeMode {
	case AstFactory_TypeMode_JSTYPE:
		result.SetJSType(a.registry.GetNativeType(rhino.JSTypeNative_ARGUMENTS_TYPE))
	case AstFactory_TypeMode_COLOR:
		result.SetColor(a.colorRegistry.Get(colors.StandardColors_ARGUMENTS_ID))
	case AstFactory_TypeMode_NONE:
	}
	return result
}

// CreateName creates a name node with the provided type information.
//
// NOTE: You should use CreateNameInScope when creating a new name node that
// references an existing variable. That version will look up the
// declaration of the variable and ensure this reference to it is correct by
// adding type information and / or any other necessary data.
//
// If you are creating a reference to a variable that won't be visible in
// the current scope (e.g. because you are creating a new variable), you can
// use this method. However, if you've just created the variable
// declaration, you could also just clone the NAME node from it to create the
// new reference.
//
// This method assumes that if the AST is normalized and the name starts
// with "$jscomp", then it must be a constant name.
//
// Java: AstFactory.createName(String, AstFactory.Type)
func (a *AstFactory) CreateName(name string, type_ AstFactory_Type) rhino.Node {
	jutil.CheckArgumentMsg1(
		name != "$jscomp",
		"Use createQName(RuntimeJsLibManager.Field) to reference $jscomp.* methods. Found: %s",
		name)
	result := rhino.IR_Name(name)
	a.SetJSTypeOrColor(type_, result)
	if a.lifeCycleStage.IsNormalized() && strings.HasPrefix(name, "$jscomp") {
		// $jscomp will always be a constant and needs to be marked that way to satisfy
		// the normalization invariants.
		// TODO: b/322009741 - Stop depending on lifeCycleStage.isNormalized() and "$jscomp" prefix to
		// decide constness. The callers must explicitly use `createConstantName` is they need a NAME
		// node that's set with IS_CONSTANT_NAME prop.
		result.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
	}
	return result
}

// CreateConstantName: use this when you know you need to create a constant
// name for const declarations
//
// Java: AstFactory.createConstantName(String, AstFactory.Type)
func (a *AstFactory) CreateConstantName(name string, type_ AstFactory_Type) rhino.Node {
	jutil.CheckArgumentMsg1(
		name != "$jscomp",
		"Use createQName(RuntimeJsLibManager.Field) to reference $jscomp.* methods. Found: %s",
		name)
	result := rhino.IR_Name(name)
	a.SetJSTypeOrColor(type_, result)
	if a.lifeCycleStage.IsNormalized() {
		// TODO: b/322009741 - Stop depending on lifeCycleStage.isNormalized() to decide constness
		result.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
	}
	return result
}

// Java: AstFactory.createName(StaticScope, String)
func (a *AstFactory) CreateNameInScope(scope rhino.StaticScope, name string) rhino.Node {
	jutil.CheckArgumentMsg1(
		name != "$jscomp",
		"Use createQName(RuntimeJsLibManager.Field) to reference $jscomp.* methods. Found: %s",
		name)
	result := rhino.IR_Name(name)

	if a.lifeCycleStage.IsNormalized() || a.typeMode != AstFactory_TypeMode_NONE {
		// We need a scope to maintain normalization and / or propagate type information
		jutil.CheckNotNullMsg2(
			scope,
			"A scope is required [lifeCycleStage: %s, typeMode: %s]",
			a.lifeCycleStage,
			a.typeMode)
		var_ := scope.GetSlot(name)
		if var_ == nil {
			// TODO(bradfordcsmith): Why is this special exception needed?
			// There are a few cases where `$jscomp` isn't found in the code (implying that runtime
			// library injection somehow didn't happen), but we do perform transpilations which require
			// it to exist. Can we fix that?
			// This only happens when type checking is not being done.
			jutil.CheckStateMsg2(a.typeMode == AstFactory_TypeMode_NONE, "Missing var %s in scope %s", name, scope)
		} else {
			declaration := jutil.CheckNotNullMsg1(
				var_.GetDeclaration(), "Cannot find type for var with missing declaration %s", var_)
			varDefinitionNode := rhino.CheckNotNullMsg1(
				declaration.GetNode(), "Missing node for declaration %s", declaration)

			// Normalization requires that all references to a constant variable have this property.
			if varDefinitionNode.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME) {
				result.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
			}
			switch a.typeMode {
			case AstFactory_TypeMode_JSTYPE:
				definitionType := varDefinitionNode.GetJSType()
				// TODO(b/149843534): crash instead of defaulting to unknown
				if definitionType != nil {
					result.SetJSType(definitionType)
				} else {
					result.SetJSType(a.unknownType)
				}
			case AstFactory_TypeMode_COLOR:
				definitionColor := varDefinitionNode.GetColor()
				// TODO(b/149843534): crash instead of defaulting to unknown
				if definitionColor != nil {
					result.SetColor(definitionColor)
				} else {
					result.SetColor(colors.StandardColors_UNKNOWN)
				}
			case AstFactory_TypeMode_NONE:
			}
		}
	}

	return result
}

// Java: AstFactory.createNameWithUnknownType(String)
func (a *AstFactory) CreateNameWithUnknownType(name string) rhino.Node {
	return a.CreateName(name, AstFactory_Type_OfJSTypeAndColor(a.unknownType, colors.StandardColors_UNKNOWN))
}

// CreateQNameUsingJSTypeInfo looks up the type of a name from a TypedScope
// created from typechecking, using the JSType API. Will crash if is called
// on an AstFactory that is created after JSType -> color conversion.
//
// Prefer CreateQName if running after JSType -> color conversion.
//
// globalTypedScope must be the top, global scope.
//
// Java: AstFactory.createQNameUsingJSTypeInfo(TypedScope, String)
func (a *AstFactory) CreateQNameUsingJSTypeInfo(globalTypedScope *TypedScope, qname string) rhino.Node {
	jutil.CheckArgumentMsg1(globalTypedScope == nil || globalTypedScope.IsGlobal(), "%s", globalTypedScope)
	a.AssertNotAddingColors()

	dot := strings.IndexByte(qname, '.')
	receiverPart := qname
	if dot >= 0 {
		receiverPart = qname[:dot]
	}

	receiver := rhino.IR_Name(receiverPart)
	if a.IsAddingTypes() {
		var_ := jutil.CheckNotNullMsg1(globalTypedScope.GetVar(receiverPart), "%s", receiverPart)
		receiver.SetJSType(jutil.CheckNotNullMsg1(var_.GetType(), "%s", var_))
	}
	if dot < 0 {
		return receiver
	}

	result := receiver
	rest := qname[dot+1:]
	for {
		nextDot := strings.IndexByte(rest, '.')
		if nextDot < 0 {
			return a.CreateGetPropWithoutColor(result, rest)
		}
		result = a.CreateGetPropWithoutColor(result, rest[:nextDot])
		rest = rest[nextDot+1:]
	}
}

// CreateQName creates a qualfied name in the given scope.
//
// Only works if StaticScope#getSlot(String) returns a name. In practice,
// that means this throws an exception for instance methods and properties,
// for example.
//
// Java: AstFactory.createQName(StaticScope, String)
func (a *AstFactory) CreateQName(scope rhino.StaticScope, qname string) rhino.Node {
	dot := strings.IndexByte(qname, '.')
	if dot < 0 {
		return a.CreateNameInScope(scope, qname)
	}
	baseNameNode := a.CreateNameInScope(scope, qname[:dot])
	return a.createQNamePropsFromString(scope, baseNameNode, qname, dot+1)
}

func (a *AstFactory) createQNamePropsFromString(
	scope rhino.StaticScope, baseNameNode rhino.Node, qname string, propStart int) rhino.Node {
	result := baseNameNode
	addingTypeInfo := a.IsAddingTypes() || a.IsAddingColors()
	for {
		relDot := strings.IndexByte(qname[propStart:], '.')
		var propertyName, prefix string
		if relDot < 0 {
			propertyName = qname[propStart:]
			prefix = qname
		} else {
			nextDot := propStart + relDot
			propertyName = qname[propStart:nextDot]
			prefix = qname[:nextDot]
			propStart = nextDot + 1
		}
		var type_ AstFactory_Type
		if addingTypeInfo {
			def := jutil.CheckNotNullMsg1(scope.GetSlot(prefix), "Cannot find name %s in StaticScope.", prefix).
				GetDeclaration().
				GetNode()
			type_ = AstFactory_Type_(def)
		}
		result = a.CreateGetProp(result, propertyName, type_)
		if relDot < 0 {
			return result
		}
	}
}

// CreateQNameOfNames creates a qualfied name in the given scope.
//
// Only works if StaticScope#getSlot(String) returns a name. In practice,
// that means this does not work for instance methods or properties, for
// example.
//
// Java: AstFactory.createQName(StaticScope, Iterable<String>)
func (a *AstFactory) CreateQNameOfNames(scope rhino.StaticScope, names []string) rhino.Node {
	if len(names) == 0 {
		panic(jutil.NewNullPointerException())
	}
	return a.CreateQNameWithPropertyList(scope, names[0], names[1:])
}

// CreateQNameWithProperties creates a qualfied name in the given scope.
//
// Only works if StaticScope#getSlot(String) returns a name. In practice,
// that means this does not work for instance methods or properties, for
// example.
//
// Java: AstFactory.createQName(StaticScope, String, String...)
func (a *AstFactory) CreateQNameWithProperties(
	scope rhino.StaticScope, baseName string, propertyNames ...string) rhino.Node {
	// Go-only: language-adapter — checkNotNull(baseName): a Go string is never null.
	return a.CreateQNameWithPropertyList(scope, baseName, propertyNames)
}

// CreateQNameWithPropertyList creates a qualfied name in the given scope.
//
// Only works if StaticScope#getSlot(String) returns a name. In practice,
// that means this does not work for instance methods or properties, for
// example.
//
// Java: AstFactory.createQName(StaticScope, String, Iterable<String>)
func (a *AstFactory) CreateQNameWithPropertyList(
	scope rhino.StaticScope, baseName string, propertyNames []string) rhino.Node {
	baseNameNode := a.CreateNameInScope(scope, baseName)
	return a.CreateQNameWithBaseNameNode(scope, baseName, baseNameNode, propertyNames)
}

// CreateQNameWithBaseNameNode creates a qualfied name in the given scope.
//
// Only works if StaticScope#getSlot(String) returns a name. In practice,
// that means this does not work for instance methods or properties, for
// example.
//
// Java: AstFactory.createQName(StaticScope, String, Node, Iterable<String>)
func (a *AstFactory) CreateQNameWithBaseNameNode(
	scope rhino.StaticScope, baseName string, baseNameNode rhino.Node, propertyNames []string) rhino.Node {
	qname := baseNameNode
	if !a.IsAddingTypes() && !a.IsAddingColors() {
		for _, propertyName := range propertyNames {
			qname = a.CreateGetProp(qname, propertyName, AstFactory_Type{})
		}
		return qname
	}
	name := baseName
	for _, propertyName := range propertyNames {
		name += "." + propertyName
		def := jutil.CheckNotNullMsg1(scope.GetSlot(name), "Cannot find name %s in StaticScope.", name).
			GetDeclaration().
			GetNode()
		qname = a.CreateGetProp(qname, propertyName, AstFactory_Type_(def))
	}
	return qname
}

func (a *AstFactory) createJscompFieldBaseNode(baseName string) rhino.Node {
	jutil.CheckStateMsg1(strings.HasPrefix(baseName, "$jscomp"), "Unexpected Field name %s", baseName)
	var baseNameNode rhino.Node
	if baseName == "$jscomp" {
		baseNameNode = a.CreateJscomp()
	} else {
		baseNameNode = a.CreateName(baseName, AstFactory_Type_OfJSTypeAndColor(a.unknownType, colors.StandardColors_UNKNOWN))
	}
	if a.lifeCycleStage.IsNormalized() {
		baseNameNode.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
	}
	return baseNameNode
}

// CreateQNameOfField creates a qualfied name in the given scope.
//
// All $jscomp runtime methods must be created using this method.
//
// Java: AstFactory.createQName(StaticScope, RuntimeJsLibManager.JsLibField)
func (a *AstFactory) CreateQNameOfField(scope rhino.StaticScope, field RuntimeJsLibManager_JsLibField) rhino.Node {
	qname := field.AssertInjected().QualifiedName()
	dot := strings.IndexByte(qname, '.')
	if dot < 0 {
		return a.createJscompFieldBaseNode(qname)
	}
	baseNameNode := a.createJscompFieldBaseNode(qname[:dot])
	return a.createQNamePropsFromString(scope, baseNameNode, qname, dot+1)
}

// Java: AstFactory.createQNameWithUnknownType(String)
func (a *AstFactory) CreateQNameWithUnknownType(qname string) rhino.Node {
	dot := strings.IndexByte(qname, '.')
	if dot < 0 {
		return a.CreateNameWithUnknownType(qname)
	}
	return a.createGetPropsWithUnknownTypeFromString(a.CreateNameWithUnknownType(qname[:dot]), qname[dot+1:])
}

// Java: AstFactory.createQNameWithUnknownType(RuntimeJsLibManager.JsLibField)
func (a *AstFactory) CreateQNameWithUnknownTypeOfField(field RuntimeJsLibManager_JsLibField) rhino.Node {
	qname := field.AssertInjected().QualifiedName()
	dot := strings.IndexByte(qname, '.')
	if dot < 0 {
		return a.createJscompFieldBaseNode(qname)
	}
	return a.createGetPropsWithUnknownTypeFromString(a.createJscompFieldBaseNode(qname[:dot]), qname[dot+1:])
}

func (a *AstFactory) createGetPropsWithUnknownTypeFromString(receiver rhino.Node, rest string) rhino.Node {
	result := receiver
	for {
		dot := strings.IndexByte(rest, '.')
		if dot < 0 {
			return a.CreateGetPropWithUnknownType(result, rest)
		}
		result = a.CreateGetPropWithUnknownType(result, rest[:dot])
		rest = rest[dot+1:]
	}
}

// Java: AstFactory.createQNameWithUnknownType(Iterable<String>)
func (a *AstFactory) CreateQNameWithUnknownTypeOfNames(names []string) rhino.Node {
	if len(names) == 0 {
		panic(jutil.NewNullPointerException())
	}
	return a.CreateQNameWithUnknownTypeWithPropertyList(names[0], names[1:])
}

// Java: AstFactory.createQNameWithUnknownType(String, Iterable<String>)
func (a *AstFactory) CreateQNameWithUnknownTypeWithPropertyList(baseName string, propertyNames []string) rhino.Node {
	baseNameNode := a.CreateNameWithUnknownType(baseName)
	return a.CreateGetPropsWithUnknownType(baseNameNode, propertyNames)
}

// Java: AstFactory.createJscomp()
func (a *AstFactory) CreateJscomp() rhino.Node {
	jscomp := rhino.IR_Name("$jscomp")
	a.SetJSTypeOrColor(AstFactory_Type_OfJSTypeAndColor(a.unknownType, colors.StandardColors_UNKNOWN), jscomp)
	if a.lifeCycleStage.IsNormalized() {
		jscomp.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
	}
	return jscomp
}

// CreatePrototypeAccess creates a Node representing <receiver>.prototype
//
// For example, given the AST for `Foo`, returns `Foo.prototype`
//
// Java: AstFactory.createPrototypeAccess(Node)
func (a *AstFactory) CreatePrototypeAccess(receiver rhino.Node) rhino.Node {
	result := rhino.IR_Getprop(receiver, "prototype")
	switch a.typeMode {
	case AstFactory_TypeMode_JSTYPE:
		result.SetJSType(a.GetJsTypeForProperty(receiver, "prototype"))
	case AstFactory_TypeMode_COLOR:
		jutil.CheckNotNullMsg1(receiver.GetColor(), "Missing color on %s", receiver)
		possiblePrototypes := receiver.GetColor().GetPrototypes()
		if len(possiblePrototypes) == 0 {
			result.SetColor(colors.StandardColors_UNKNOWN)
		} else {
			result.SetColor(colors.Color_CreateUnionSlice(possiblePrototypes))
		}
	case AstFactory_TypeMode_NONE:
	}
	return result
}

// CreateJSCompDotGlobalAccess creates an access of the given qname on
// $jscomp.global
//
// For example, given "Object.defineProperties", returns an AST
// representation of "$jscomp.global.Object.defineProperties".
//
// This may be useful if adding synthetic code to a local scope, which can
// shadow the global like Object you're trying to access. The $jscomp global
// should not be shadowed.
//
// Java: AstFactory.createJSCompDotGlobalAccess(StaticScope, String)
func (a *AstFactory) CreateJSCompDotGlobalAccess(scope rhino.StaticScope, qname string) rhino.Node {
	global := a.runtimeJsLibManager.GetJsLibField("$jscomp.global")
	jscompDotGlobal := a.CreateQNameOfField(scope, global)
	result := a.CreateQName(scope, qname)
	// Move the fully qualified qname onto the $jscomp.global getprop
	qnameRoot := NodeUtil_GetRootOfQualifiedName(result)
	qnameRoot.ReplaceWith(a.CreateGetProp(jscompDotGlobal, qnameRoot.GetString(), AstFactory_Type_(qnameRoot)))
	return result
}

// Deprecated: Use CreateGetProp instead.
//
// Java: AstFactory.createGetPropWithoutColor(Node, String)
func (a *AstFactory) CreateGetPropWithoutColor(receiver rhino.Node, propertyName string) rhino.Node {
	a.AssertNotAddingColors()
	result := rhino.IR_Getprop(receiver, propertyName)
	if a.IsAddingTypes() {
		result.SetJSType(a.GetJsTypeForProperty(receiver, propertyName))
	}
	return result
}

// Java: AstFactory.createGetProp(Node, String, AstFactory.Type)
func (a *AstFactory) CreateGetProp(receiver rhino.Node, propertyName string, type_ AstFactory_Type) rhino.Node {
	result := rhino.IR_Getprop(receiver, propertyName)
	a.SetJSTypeOrColor(type_, result)
	return result
}

// CreateGetPropsWithoutColors creates a tree of nodes representing
// `receiver.name1.name2.etc`.
//
// Deprecated: use individual CreateGetProp calls or CreateQName instead.
//
// Java: AstFactory.createGetPropsWithoutColors(Node, Iterable<String>)
func (a *AstFactory) CreateGetPropsWithoutColors(receiver rhino.Node, propertyNames []string) rhino.Node {
	a.AssertNotAddingColors()
	result := receiver
	for _, propertyName := range propertyNames {
		result = a.CreateGetPropWithoutColor(result, propertyName)
	}
	return result
}

// CreateGetPropsWithUnknownType creates a tree of nodes representing
// `receiver.name1.name2.etc`.
//
// Java: AstFactory.createGetPropsWithUnknownType(Node, Iterable<String>)
func (a *AstFactory) CreateGetPropsWithUnknownType(receiver rhino.Node, propertyNames []string) rhino.Node {
	result := receiver
	for _, propertyName := range propertyNames {
		result = a.CreateGetPropWithUnknownType(result, propertyName)
	}
	return result
}

// Java: AstFactory.createGetPropWithUnknownType(Node, String)
func (a *AstFactory) CreateGetPropWithUnknownType(receiver rhino.Node, propertyName string) rhino.Node {
	result := rhino.IR_Getprop(receiver, propertyName)
	a.SetJSTypeOrColor(AstFactory_Type_OfJSTypeAndColor(a.unknownType, colors.StandardColors_UNKNOWN), result)
	return result
}

// Java: AstFactory.createStartOptChainGetprop(Node, String, AstFactory.Type)
func (a *AstFactory) CreateStartOptChainGetprop(receiver rhino.Node, propertyName string, type_ AstFactory_Type) rhino.Node {
	result := rhino.IR_StartOptChainGetprop(receiver, propertyName)
	a.SetJSTypeOrColor(type_, result)
	return result
}

// Java: AstFactory.createContinueOptChainGetprop(Node, String, AstFactory.Type)
func (a *AstFactory) CreateContinueOptChainGetprop(receiver rhino.Node, propertyName string, type_ AstFactory_Type) rhino.Node {
	result := rhino.IR_ContinueOptChainGetprop(receiver, propertyName)
	a.SetJSTypeOrColor(type_, result)
	return result
}

// Java: AstFactory.createStartOptChainGetelem(Node, Node, AstFactory.Type)
func (a *AstFactory) CreateStartOptChainGetelem(receiver rhino.Node, elem rhino.Node, type_ AstFactory_Type) rhino.Node {
	result := rhino.IR_StartOptChainGetelem(receiver, elem)
	a.SetJSTypeOrColor(type_, result)
	return result
}

// Java: AstFactory.createContinueOptChainGetelem(Node, Node, AstFactory.Type)
func (a *AstFactory) CreateContinueOptChainGetelem(receiver rhino.Node, elem rhino.Node, type_ AstFactory_Type) rhino.Node {
	result := rhino.IR_ContinueOptChainGetelem(receiver, elem)
	a.SetJSTypeOrColor(type_, result)
	return result
}

// Java: AstFactory.createStartOptChainCall(Node, AstFactory.Type, Node...)
func (a *AstFactory) CreateStartOptChainCall(receiver rhino.Node, type_ AstFactory_Type, args ...rhino.Node) rhino.Node {
	result := rhino.IR_StartOptChainCall(receiver, args...)
	a.SetJSTypeOrColor(type_, result)
	return result
}

// Java: AstFactory.createContinueOptChainCall(Node, AstFactory.Type, Node...)
func (a *AstFactory) CreateContinueOptChainCall(receiver rhino.Node, type_ AstFactory_Type, args ...rhino.Node) rhino.Node {
	result := rhino.IR_ContinueOptChainCall(receiver, args...)
	a.SetJSTypeOrColor(type_, result)
	return result
}

// Java: AstFactory.createGetElem(Node, Node)
func (a *AstFactory) CreateGetElem(receiver rhino.Node, key rhino.Node) rhino.Node {
	result := rhino.IR_Getelem(receiver, key)
	// TODO(bradfordcsmith): When receiver is an Array<T> or an Object<K, V>, use the template
	// type here.
	a.SetJSTypeOrColor(AstFactory_Type_OfJSTypeAndColor(a.unknownType, colors.StandardColors_UNKNOWN), result)
	return result
}

// Java: AstFactory.createCall(Node, AstFactory.Type, Node...)
func (a *AstFactory) CreateCall(callee rhino.Node, resultType AstFactory_Type, args ...rhino.Node) rhino.Node {
	result := NodeUtil_NewCallNode(callee, args...)
	a.SetJSTypeOrColor(resultType, result)
	return result
}

// Java: AstFactory.createCallWithUnknownType(Node, Node...)
func (a *AstFactory) CreateCallWithUnknownType(callee rhino.Node, args ...rhino.Node) rhino.Node {
	return a.CreateCall(callee, AstFactory_Type_OfJSTypeAndColor(a.unknownType, colors.StandardColors_UNKNOWN), args...)
}

// CreateObjectDotAssignCall creates a call to Object.assign that returns the
// specified type.
//
// Object.assign returns !Object in the externs, which can lose type
// information if the actual type is known.
//
// Java: AstFactory.createObjectDotAssignCall(StaticScope, AstFactory.Type, Node...)
func (a *AstFactory) CreateObjectDotAssignCall(
	scope rhino.StaticScope, returnType AstFactory_Type, args ...rhino.Node) rhino.Node {
	objAssign := a.CreateQName(scope, "Object.assign")
	result := a.CreateCall(objAssign, returnType, args...)

	switch a.typeMode {
	case AstFactory_TypeMode_JSTYPE:
		// Make a unique function type that returns the exact type we've inferred it to be.
		// Object.assign in the externs just returns !Object, which loses type information.
		returnJSType := returnType.GetJSType(a.registry)
		objAssignType :=
			a.registry.CreateFunctionTypeWithVarArgs(
				returnJSType,
				a.registry.GetNativeType(rhino.JSTypeNative_OBJECT_TYPE),
				a.registry.CreateUnionTypeOfNatives(rhino.JSTypeNative_OBJECT_TYPE, rhino.JSTypeNative_NULL_TYPE))
		objAssign.SetJSType(objAssignType)
	case AstFactory_TypeMode_COLOR, AstFactory_TypeMode_NONE:
	}

	return result
}

// Java: AstFactory.createNewNode(Node, Node...)
func (a *AstFactory) CreateNewNode(target rhino.Node, args ...rhino.Node) rhino.Node {
	result := rhino.IR_NewNode(target, args...)
	switch a.typeMode {
	case AstFactory_TypeMode_JSTYPE:
		instanceType := target.GetJSType()
		if instanceType.IsFunctionType() {
			instanceType = astFactoryObjectTypeAsJSType(instanceType.ToMaybeFunctionType().GetInstanceType())
		} else {
			instanceType = a.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
		}
		result.SetJSType(instanceType)
	case AstFactory_TypeMode_COLOR:
		result.SetColor(a.GetInstanceOfColor(target.GetColor()))
	case AstFactory_TypeMode_NONE:
	}
	return result
}

// astFactoryObjectTypeAsJSType widens an ObjectType to a JSType.
//
// Go-only: language-adapter — keeps a nil ObjectType a nil JSType.
func astFactoryObjectTypeAsJSType(t rhino.ObjectType) rhino.JSType {
	if t == nil {
		return nil
	}
	return t
}

// Java: AstFactory.getInstanceOfColor(Color)
func (a *AstFactory) GetInstanceOfColor(color *colors.Color) *colors.Color {
	possibleInstanceColors := color.GetInstanceColors()
	if len(possibleInstanceColors) == 0 {
		return colors.StandardColors_UNKNOWN
	}
	return colors.Color_CreateUnionSlice(possibleInstanceColors)
}

// CreateConstructorCall creates a call that returns an instance of the
// given class type.
//
// This method is intended for use in special cases, such as calling
// `super()` in a constructor.
//
// Java: AstFactory.createConstructorCall(AstFactory.Type, Node, Node...)
func (a *AstFactory) CreateConstructorCall(classType AstFactory_Type, callee rhino.Node, args ...rhino.Node) rhino.Node {
	result := NodeUtil_NewCallNode(callee, args...)
	switch a.typeMode {
	case AstFactory_TypeMode_JSTYPE:
		classJSType := classType.GetJSType(a.registry)
		constructorType := jutil.CheckNotNull(classJSType.ToMaybeFunctionType())
		instanceType := jutil.CheckNotNull(constructorType.GetInstanceType())
		result.SetJSType(instanceType)
	case AstFactory_TypeMode_COLOR:
		result.SetColor(a.GetInstanceOfColor(classType.GetColor(a.colorRegistry)))
	case AstFactory_TypeMode_NONE:
	}
	return result
}

// Java: AstFactory.createJSCompMakeIteratorCall(Node, StaticScope)
func (a *AstFactory) CreateJSCompMakeIteratorCall(iterable rhino.Node, scope rhino.StaticScope) rhino.Node {
	makeIterator := a.runtimeJsLibManager.GetJsLibField("$jscomp.makeIterator")
	makeIteratorName := a.CreateQNameOfField(scope, makeIterator)
	var type_ AstFactory_Type
	switch a.typeMode {
	case AstFactory_TypeMode_JSTYPE:
		// Since createCall (currently) doesn't handle templated functions, fill in the template
		// types
		// of makeIteratorName manually.
		// e.g get `number` from `Iterable<number>`
		iterableType :=
			iterable.
				GetJSType().
				GetTemplateTypeMap().
				GetResolvedTemplateType(a.registry.GetIterableValueTemplate())
		makeIteratorType := makeIteratorName.GetJSType()
		// e.g. replace
		//   function(Iterable<T>): Iterator<T>
		// with
		//   function(Iterable<number>): Iterator<number>
		makeIteratorName.SetJSType(
			a.ReplaceTemplate(makeIteratorType, []rhino.JSType{iterableType}))
		type_ = AstFactory_Type_OfJSType(makeIteratorName.GetJSType().AssertFunctionType().GetReturnType())
	case AstFactory_TypeMode_COLOR:
		type_ = AstFactory_Type_OfColor(a.colorRegistry.Get(colors.StandardColors_ITERATOR_ID))
	case AstFactory_TypeMode_NONE:
		type_ = AstFactory_NoTypeInformation()
	default:
		panic(jutil.NewAssertionErrorNoMessage())
	}
	call := a.CreateCall(makeIteratorName, type_, iterable)
	call.PutBooleanProp(rhino.Node_FREE_CALL, true)
	return call
}

// Java: AstFactory.createJscompIteratorCloseCall(Node, Node, StaticScope)
func (a *AstFactory) CreateJscompIteratorCloseCall(iterator rhino.Node, iterResult rhino.Node, scope rhino.StaticScope) rhino.Node {
	iteratorClose := a.runtimeJsLibManager.GetJsLibField("$jscomp.iteratorClose")
	iteratorCloseName := a.CreateQNameOfField(scope, iteratorClose)
	return a.CreateCall(
		iteratorCloseName,
		AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_VOID_TYPE, colors.StandardColors_NULL_OR_VOID),
		iterator,
		iterResult)
}

// Java: AstFactory.createJscompArrayFromIteratorCall(Node, StaticScope)
func (a *AstFactory) CreateJscompArrayFromIteratorCall(iterator rhino.Node, scope rhino.StaticScope) rhino.Node {
	arrayFromIterator := a.runtimeJsLibManager.GetJsLibField("$jscomp.arrayFromIterator")
	makeIteratorName := a.CreateQNameOfField(scope, arrayFromIterator)

	var resultType AstFactory_Type
	switch a.typeMode {
	case AstFactory_TypeMode_JSTYPE:
		// Since createCall (currently) doesn't handle templated functions, fill in the template
		// types of makeIteratorName manually.
		iterableType :=
			iterator.
				GetJSType().
				GetTemplateTypeMap().
				GetResolvedTemplateType(a.registry.GetIteratorValueTemplate())
		makeIteratorType := makeIteratorName.GetJSType()
		// e.g. replace
		//   function(Iterator<T>): Array<T>
		// with
		//   function(Iterator<number>): Array<number>
		makeIteratorName.SetJSType(
			a.ReplaceTemplate(makeIteratorType, []rhino.JSType{iterableType}))
		resultType = AstFactory_Type_OfJSType(makeIteratorName.GetJSType().AssertFunctionType().GetReturnType())
	// colors don't include generics, so just set the return type to Array.
	case AstFactory_TypeMode_COLOR:
		resultType = AstFactory_Type_OfColor(a.colorRegistry.Get(colors.StandardColors_ARRAY_ID))
	case AstFactory_TypeMode_NONE:
		resultType = AstFactory_NoTypeInformation()
	default:
		// Go-only: language-adapter — unreachable for a valid enum value: the
		// Java switch expression is exhaustive.
		panic(jutil.NewIllegalStateException())
	}
	call := a.CreateCall(makeIteratorName, resultType, iterator)
	call.PutBooleanProp(rhino.Node_FREE_CALL, true)
	return call
}

// Java: AstFactory.createJscompArrayFromIterableCall(Node, StaticScope)
func (a *AstFactory) CreateJscompArrayFromIterableCall(iterable rhino.Node, scope rhino.StaticScope) rhino.Node {
	arrayFromIterable := a.runtimeJsLibManager.GetJsLibField("$jscomp.arrayFromIterable")
	makeIterableName := a.CreateQNameOfField(scope, arrayFromIterable)

	var resultType AstFactory_Type
	switch a.typeMode {
	case AstFactory_TypeMode_JSTYPE:
		// Since createCall (currently) doesn't handle templated functions, fill in the template
		// types of makeIteratorName manually.
		iterableType :=
			iterable.
				GetJSType().
				GetTemplateTypeMap().
				GetResolvedTemplateType(a.registry.GetIterableValueTemplate())
		makeIterableType := makeIterableName.GetJSType()
		// e.g. replace
		//   function(Iterable<T>): Array<T>
		// with
		//   function(Iterable<number>): Array<number>
		makeIterableName.SetJSType(
			a.ReplaceTemplate(makeIterableType, []rhino.JSType{iterableType}))
		resultType = AstFactory_Type_OfJSType(makeIterableName.GetJSType().AssertFunctionType().GetReturnType())
	// colors don't include generics, so just set the return type to Array.
	case AstFactory_TypeMode_COLOR:
		resultType = AstFactory_Type_OfColor(a.colorRegistry.Get(colors.StandardColors_ARRAY_ID))
	case AstFactory_TypeMode_NONE:
		resultType = AstFactory_NoTypeInformation()
	default:
		// Go-only: language-adapter — unreachable for a valid enum value: the
		// Java switch expression is exhaustive.
		panic(jutil.NewIllegalStateException())
	}
	call := a.CreateCall(makeIterableName, resultType, iterable)
	call.PutBooleanProp(rhino.Node_FREE_CALL, true)
	return call
}

// CreateJSCompMakeAsyncIteratorCall: given an iterable like rhs in
//
//	for await (lhs of rhs) { block(); }
//
// returns a call node for the rhs wrapped in a $jscomp.makeAsyncIterator
// call.
//
//	$jscomp.makeAsyncIterator(rhs)
//
// Java: AstFactory.createJSCompMakeAsyncIteratorCall(Node, StaticScope)
func (a *AstFactory) CreateJSCompMakeAsyncIteratorCall(iterable rhino.Node, scope rhino.StaticScope) rhino.Node {
	makeAsyncIterator := a.runtimeJsLibManager.GetJsLibField("$jscomp.makeAsyncIterator")
	makeIteratorAsyncName := a.CreateQNameOfField(scope, makeAsyncIterator)
	var resultType AstFactory_Type
	switch a.typeMode {
	case AstFactory_TypeMode_JSTYPE:
		// Since createCall (currently) doesn't handle templated functions, fill in the template
		// types of makeIteratorName manually.
		// e.g get `number` from `AsyncIterable<number>`
		asyncIterableType :=
			JsIterables_MaybeBoxIterableOrAsyncIterable(iterable.GetJSType(), a.registry).
				OrElse(a.unknownType)
		makeAsyncIteratorType := makeIteratorAsyncName.GetJSType()
		// e.g. replace
		//   function(AsyncIterable<T>): AsyncIterator<T>
		// with
		//   function(AsyncIterable<number>): AsyncIterator<number>
		makeIteratorAsyncName.SetJSType(
			a.ReplaceTemplate(makeAsyncIteratorType, []rhino.JSType{asyncIterableType}))
		resultType = AstFactory_Type_OfJSType(makeIteratorAsyncName.GetJSType().AssertFunctionType().GetReturnType())
	case AstFactory_TypeMode_COLOR:
		resultType = AstFactory_Type_OfColor(a.colorRegistry.Get(colors.StandardColors_ASYNC_ITERATOR_ITERABLE_ID))
	case AstFactory_TypeMode_NONE:
		resultType = AstFactory_NoTypeInformation()
	default:
		panic(jutil.NewAssertionErrorNoMessage())
	}
	call := a.CreateCall(makeIteratorAsyncName, resultType, iterable)
	call.PutBooleanProp(rhino.Node_FREE_CALL, true)
	return call
}

// Java: AstFactory.replaceTemplate(JSType, ImmutableList<JSType>)
func (a *AstFactory) ReplaceTemplate(templatedType rhino.JSType, templateTypes []rhino.JSType) rhino.JSType {
	typeMap :=
		a.registry.
			GetEmptyTemplateTypeMap().
			CopyWithExtensionOfKeys(templatedType.GetTemplateTypeMap().GetTemplateKeys(), templateTypes)
	replacer := rhino.TemplateTypeReplacer_MakeForPartialReplacement(a.registry, typeMap)
	return replacer.VisitJSType(templatedType)
}

// Java: AstFactory.createJscompAsyncExecutePromiseGeneratorFunctionCall(StaticScope, Node)
func (a *AstFactory) CreateJscompAsyncExecutePromiseGeneratorFunctionCall(
	scope rhino.StaticScope, generatorFunction rhino.Node) rhino.Node {
	method := a.runtimeJsLibManager.GetJsLibField("$jscomp.asyncExecutePromiseGeneratorFunction")
	jscompDotAsyncExecutePromiseGeneratorFunction := a.CreateQNameOfField(scope, method)
	var resultType AstFactory_Type
	switch a.typeMode {
	case AstFactory_TypeMode_JSTYPE:
		// TODO(bradfordcsmith): Maybe update the type to be more specific
		// Currently this method expects `function(): !Generator<?>` and returns `Promise<?>`.
		// Since we propagate type information only if type checking has already run,
		// these unknowns probably don't matter, but we should be able to be more specific
		// with the return type at least.
		resultType =
			AstFactory_Type_OfJSType(
				jscompDotAsyncExecutePromiseGeneratorFunction.
					GetJSType().
					AssertFunctionType().
					GetReturnType())
	case AstFactory_TypeMode_COLOR:
		resultType = AstFactory_Type_OfColor(a.colorRegistry.Get(colors.StandardColors_PROMISE_ID))
	case AstFactory_TypeMode_NONE:
		resultType = AstFactory_Type_OfJSType(a.unknownType)
	default:
		// Go-only: language-adapter — unreachable for a valid enum value: the
		// Java switch expression is exhaustive.
		panic(jutil.NewIllegalStateException())
	}
	call :=
		a.CreateCall(jscompDotAsyncExecutePromiseGeneratorFunction, resultType, generatorFunction)
	call.PutBooleanProp(rhino.Node_FREE_CALL, true)
	return call
}

// Java: AstFactory.getNativeType(JSTypeNative)
func (a *AstFactory) GetNativeType(nativeType rhino.JSTypeNative) rhino.JSType {
	jutil.CheckNotNullMsg(a.registry, "registry is null")
	return jutil.CheckNotNullMsg1(
		a.registry.GetNativeType(nativeType), "native type not found: %s", nativeType)
}

// Java: AstFactory.getJsTypeForProperty(Node, String)
func (a *AstFactory) GetJsTypeForProperty(receiver rhino.Node, propertyName string) rhino.JSType {
	// NOTE: we use both findPropertyType and getPropertyType because they are subtly
	// different: findPropertyType works on JSType, autoboxing scalars and joining unions,
	// but it returns null if the type is not found and does not handle dynamic types of
	// Function.prototype.call and .apply; whereas getPropertyType does not autobox nor
	// iterate over unions, but it does synthesize the function properties correctly, and
	// it returns unknown instead of null if the property is missing.
	var getpropType rhino.JSType = nil
	receiverJSType := receiver.GetJSType()
	if receiverJSType != nil {
		getpropType = receiverJSType.FindPropertyType(propertyName)
		if getpropType == nil {
			receiverObjectType := rhino.ObjectType_Cast(receiverJSType.Autobox())
			if receiverObjectType == nil {
				getpropType = a.unknownType
			} else {
				getpropType = receiverObjectType.GetPropertyType(propertyName)
			}
		}
	}
	if getpropType == nil {
		getpropType = a.unknownType
	}
	// TODO(bradfordcsmith): Special case $jscomp.global until we annotate its type correctly.
	if getpropType.IsUnknownType() &&
		propertyName == "global" &&
		receiver.MatchesName("$jscomp") {
		getpropType = a.GetNativeType(rhino.JSTypeNative_GLOBAL_THIS)
	}
	return getpropType
}

// Java: AstFactory.setJSTypeOrColor(AstFactory.Type, Node)
func (a *AstFactory) SetJSTypeOrColor(type_ AstFactory_Type, result rhino.Node) {
	switch a.typeMode {
	case AstFactory_TypeMode_JSTYPE:
		result.SetJSType(type_.GetJSType(a.registry))
	case AstFactory_TypeMode_COLOR:
		result.SetColor(type_.GetColor(a.colorRegistry))
	case AstFactory_TypeMode_NONE:
	}
}

type astFactoryTypeKind uint8

const (
	astFactoryTypeZero     astFactoryTypeKind = 0
	astFactoryTypeFromNode astFactoryTypeKind = 1 << iota
	astFactoryTypeExplicit
	astFactoryTypeHasJSTypeNative
	astFactoryTypeHasColorId
)

// AstFactory_Type is AstFactory.Type, represented as a zero-allocation value.
type AstFactory_Type struct {
	jstype       rhino.JSType
	color        *colors.Color
	colorId      colors.ColorId
	node         rhino.Node
	jstypeNative rhino.JSTypeNative
	kind         astFactoryTypeKind
}

// IsZero reports whether t is the zero AstFactory_Type value.
func (t AstFactory_Type) IsZero() bool {
	return t.kind == astFactoryTypeZero
}

// Java: AstFactory.Type.getJSType(JSTypeRegistry)
func (t AstFactory_Type) GetJSType(registry *rhino.JSTypeRegistry) rhino.JSType {
	if t.kind == astFactoryTypeFromNode {
		jstype := t.node.GetJSType()
		// TODO(b/149843534): crash instead of defaulting to unknown
		if jstype != nil {
			return jstype
		}
		return registry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
	}
	if t.jstype != nil {
		return t.jstype
	}
	if t.kind&astFactoryTypeHasJSTypeNative == 0 {
		// Go-only: language-adapter — registry.getNativeType(null) throws a
		// NullPointerException (typeId.ordinal()).
		panic(jutil.NewNullPointerException())
	}
	return registry.GetNativeType(t.jstypeNative)
}

// Java: AstFactory.Type.getColor(ColorRegistry)
func (t AstFactory_Type) GetColor(registry *colors.ColorRegistry) *colors.Color {
	if t.kind == astFactoryTypeFromNode {
		color := t.node.GetColor()
		// TODO(b/149843534): crash instead of defaulting to unknown
		if color != nil {
			return color
		}
		return colors.StandardColors_UNKNOWN
	}
	if t.color != nil {
		return t.color
	}
	if t.kind&astFactoryTypeHasColorId == 0 {
		panic(jutil.NewNullPointerException())
	}
	return registry.Get(t.colorId)
}

// AstFactory_Type_ uses the JSType or Color of the given node as a template
// for adding type information
//
// Java: AstFactory.type(Node)
func AstFactory_Type_(node rhino.Node) AstFactory_Type {
	return AstFactory_Type{kind: astFactoryTypeFromNode, node: node}
}

// Java: AstFactory.type(JSType)
func AstFactory_Type_OfJSType(type_ rhino.JSType) AstFactory_Type {
	return AstFactory_Type{kind: astFactoryTypeExplicit, jstype: type_}
}

// Java: AstFactory.type(JSTypeNative)
func AstFactory_Type_OfJSTypeNative(type_ rhino.JSTypeNative) AstFactory_Type {
	return AstFactory_Type{
		kind:         astFactoryTypeExplicit | astFactoryTypeHasJSTypeNative,
		jstypeNative: type_,
	}
}

// Java: AstFactory.type(Color)
func AstFactory_Type_OfColor(type_ *colors.Color) AstFactory_Type {
	return AstFactory_Type{kind: astFactoryTypeExplicit, color: type_}
}

// Java: AstFactory.type(ColorId)
func AstFactory_Type_OfColorId(type_ colors.ColorId) AstFactory_Type {
	return AstFactory_Type{
		kind:    astFactoryTypeExplicit | astFactoryTypeHasColorId,
		colorId: type_,
	}
}

// Java: AstFactory.type(JSType, Color)
func AstFactory_Type_OfJSTypeAndColor(type_ rhino.JSType, color *colors.Color) AstFactory_Type {
	return AstFactory_Type{kind: astFactoryTypeExplicit, jstype: type_, color: color}
}

// Java: AstFactory.type(JSTypeNative, Color)
func AstFactory_Type_OfJSTypeNativeAndColor(type_ rhino.JSTypeNative, color *colors.Color) AstFactory_Type {
	return AstFactory_Type{
		kind:         astFactoryTypeExplicit | astFactoryTypeHasJSTypeNative,
		jstypeNative: type_,
		color:        color,
	}
}

// Java: AstFactory.noTypeInformation()
func AstFactory_NoTypeInformation() AstFactory_Type {
	return AstFactory_Type{kind: astFactoryTypeExplicit}
}
