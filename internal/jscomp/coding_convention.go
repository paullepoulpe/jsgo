// Port of com/google/javascript/jscomp/CodingConvention.java @ closure-compiler v20261104

/*
 * Copyright 2007 The Closure Compiler Authors.
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
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// CodingConvention defines a set of hooks to customize the behavior of the
// Compiler for a specific team/company.
//
// Overloads (PORTING.md §3): isExported(String) keeps the plain name and
// isExported(String, boolean) is IsExportedWithLocal.
type CodingConvention interface {
	// IsConstant checks whether a given variable name, such as a name in
	// all-caps should be treated as if it had the @const annotation.
	//
	// variableName is a potentially constant variable name; returns true if
	// the name should be treated as a constant.
	//
	// Java: CodingConvention.isConstant(String)
	IsConstant(variableName string) bool

	// IsConstantKey checks whether a given key of an object literal, such as
	// a name in all-caps should be treated as if it had the @const
	// annotation.
	//
	// Java: CodingConvention.isConstantKey(String)
	IsConstantKey(keyName string) bool

	// IsValidEnumKey checks that a given key may be used as a key for an
	// enum.
	//
	// key is the potential key to an enum; returns true if the key may be
	// used as an enum key, false otherwise.
	//
	// Java: CodingConvention.isValidEnumKey(String)
	IsValidEnumKey(key string) bool

	// IsOptionalParameter checks whether a given parameter name should be
	// treated as an optional parameter as far as type checking or function
	// call arg count checking is concerned. Note that an optional function
	// parameter may be declared as a simple type and is automatically
	// converted to a union of the declared type and Undefined.
	//
	// parameter is the parameter's node; returns true if the parameter
	// should be treated as an optional parameter.
	//
	// Java: CodingConvention.isOptionalParameter(Node)
	IsOptionalParameter(parameter rhino.Node) bool

	// IsVarArgsParameter checks whether a given parameter should be treated
	// as a marker for a variable argument list function. A VarArgs parameter
	// must be the last parameter in a function declaration.
	//
	// parameter is the parameter's node; returns true if the parameter
	// should be treated as a variable length parameter.
	//
	// Java: CodingConvention.isVarArgsParameter(Node)
	IsVarArgsParameter(parameter rhino.Node) bool

	// IsFunctionCallThatAlwaysThrows is used by CheckMissingReturn. When a
	// function call always throws an error, it can be the last stm of a
	// block and we don't warn about missing return.
	//
	// Java: CodingConvention.isFunctionCallThatAlwaysThrows(Node)
	IsFunctionCallThatAlwaysThrows(n rhino.Node) bool

	// IsExportedWithLocal checks whether a global variable or function name
	// should be treated as exported, or externally referenceable.
	//
	// name is a global variable or function name; local is true if the name
	// is a local variable. Returns true if the name should be considered
	// exported.
	//
	// Java: CodingConvention.isExported(String, boolean)
	IsExportedWithLocal(name string, local bool) bool

	// IsExported is equivalent to `isExported(name, true) ||
	// isExported(name, false);`
	//
	// Should only be used to check if a property is exported. Variables
	// should always use IsExportedWithLocal, as in most cases local
	// variables should not be treated as exported.
	//
	// Do not override! Unfortunately, that cannot be enforced without making
	// this an abstract class.
	//
	// A Java default method: implementations that do not override it call
	// CodingConvention_IsExported.
	//
	// Java: CodingConvention.isExported(String)
	IsExported(name string) bool

	// GetPackageName returns the package name for the given source file and
	// true, or ("", false) if no package name is known.
	//
	// Java: CodingConvention.getPackageName(StaticSourceFile)
	GetPackageName(source rhino.StaticSourceFile) (string, bool)

	// GetClassesDefinedByCall checks if the given method defines a subclass
	// relationship, and if it does, returns information on that
	// relationship. By default, always returns null. Meant to be overridden
	// by subclasses.
	//
	// callNode is a CALL node.
	//
	// Java: CodingConvention.getClassesDefinedByCall(Node)
	GetClassesDefinedByCall(callNode rhino.Node) *CodingConvention_SubclassRelationship

	// IsClassFactoryCall checks if the given method is a call to a class
	// factory, such a factory returns a unique class.
	//
	// callNode is a CALL node.
	//
	// Java: CodingConvention.isClassFactoryCall(Node)
	IsClassFactoryCall(callNode rhino.Node) bool

	// IsSuperClassReference returns true if passed a string referring to the
	// superclass. The string will usually be from the string node at the
	// right of a GETPROP, e.g. this.superClass_.
	//
	// Java: CodingConvention.isSuperClassReference(String)
	IsSuperClassReference(propertyName string) bool

	// GetExportPropertyFunction returns the function name used when exporting
	// properties, or ("", false) if none. Signature: fn(object, publicName, symbol).
	//
	// Java: CodingConvention.getExportPropertyFunction()
	GetExportPropertyFunction() (string, bool)

	// GetExportSymbolFunction returns the function name used when exporting
	// symbols, or ("", false) if none. Signature: fn(publicPath, object).
	//
	// Java: CodingConvention.getExportSymbolFunction()
	GetExportSymbolFunction() (string, bool)

	// IdentifyTypeDeclarationCall checks if the given CALL node is
	// forward-declaring a single type, and returns its name and true if so.
	//
	// Java: CodingConvention.identifyTypeDeclarationCall(Node)
	IdentifyTypeDeclarationCall(n rhino.Node) (string, bool)

	// ApplySubclassRelationship: in many JS libraries, the function that
	// produces inheritance also adds properties to the superclass and/or
	// subclass.
	//
	// Java: CodingConvention.applySubclassRelationship(NominalTypeBuilder, NominalTypeBuilder, SubclassType)
	ApplySubclassRelationship(
		parent *rhino.NominalTypeBuilder, child *rhino.NominalTypeBuilder, type_ CodingConvention_SubclassType)

	// GetAbstractMethodName is the function name for abstract methods, or "" if
	// none. An abstract method can be assigned to an interface method instead
	// of a function expression in order to avoid linter warnings produced by
	// assigning a function without a return value where a return value is
	// expected.
	//
	// Java: CodingConvention.getAbstractMethodName()
	GetAbstractMethodName() string

	// GetSingletonGetterClassName checks if the given method defines a
	// singleton getter, and if it does, returns the name of the class with
	// the singleton getter and true. By default, always returns ("", false).
	// Meant to be overridden by subclasses.
	//
	// addSingletonGetter needs a coding convention because in the general
	// case, it can't be inlined. The function inliner sees that it creates
	// an alias to the given class in an inner closure, and bails out.
	//
	// callNode is a CALL node.
	//
	// Java: CodingConvention.getSingletonGetterClassName(Node)
	GetSingletonGetterClassName(callNode rhino.Node) (string, bool)

	// ApplySingletonGetter: in many JS libraries, the function that adds a
	// singleton getter to a class adds properties to the class.
	//
	// Java: CodingConvention.applySingletonGetter(NominalTypeBuilder, FunctionType)
	ApplySingletonGetter(classType *rhino.NominalTypeBuilder, getterType rhino.FunctionType)

	// DescribeFunctionBind returns a Bind description and true if n is a
	// recognized function bind call.
	//
	// When seeing an expression exp1.bind(recv, arg1, ...); we only know
	// that it's a function bind if exp1 has type function. W/out type info,
	// exp1 has certainly a function type only if it's a function literal.
	//
	// If type checking has already happened, exp1's type is attached to the
	// AST node. When checkTypes is true, describeFunctionBind looks for that
	// type.
	//
	// checkTypes: check that the bound node has a function type.
	//
	// Java: CodingConvention.describeFunctionBind(Node, boolean)
	DescribeFunctionBind(n rhino.Node, checkTypes bool) (CodingConvention_Bind, bool)

	// DescribeCachingCall builds a Cache description from the given call node
	// and returns it with true, or false if the Node does not resemble a
	// cache utility call.
	//
	// This should match calls to a cache utility method. This type of node
	// is specially considered for side-effects since conventionally storing
	// something on a cache object would be seen as a side-effect.
	//
	// Java: CodingConvention.describeCachingCall(Node)
	DescribeCachingCall(node rhino.Node) (CodingConvention_Cache, bool)

	// IsPropertyTestFunction: whether this CALL function is testing for the
	// existence of a property.
	//
	// Java: CodingConvention.isPropertyTestFunction(Node)
	IsPropertyTestFunction(call rhino.Node) bool

	// IsPropertyRenameFunction: whether this GETPROP or NAME node is the
	// function is returning the string name for a property, but allows
	// renaming.
	//
	// Java: CodingConvention.isPropertyRenameFunction(Node)
	IsPropertyRenameFunction(nameNode rhino.Node) bool

	// GetObjectLiteralCast checks if the given method performs an object
	// literal cast, and if it does, returns information on the cast and true.
	// By default, always returns false. Meant to be overridden by subclasses.
	//
	// callNode is a CALL node.
	//
	// Java: CodingConvention.getObjectLiteralCast(Node)
	GetObjectLiteralCast(callNode rhino.Node) (CodingConvention_ObjectLiteralCast, bool)

	// GetAssertionFunctions returns the set of AssertionFunction.
	//
	// Java: CodingConvention.getAssertionFunctions()
	GetAssertionFunctions() []*CodingConvention_AssertionFunctionSpec
}

// CodingConvention_IsExported is the default body of
// CodingConvention.isExported(String).
//
// Java: CodingConvention.isExported(String) (default method)
func CodingConvention_IsExported(c CodingConvention, name string) bool {
	return c.IsExportedWithLocal(name, true) || c.IsExportedWithLocal(name, false)
}

// CodingConvention_Bind is the Bind class.
type CodingConvention_Bind struct {
	// The target of the bind action
	target rhino.Node
	// The node representing the "this" value, maybe NilNode
	thisValue rhino.Node
	// The head of a Node list representing the parameters
	parameters rhino.Node
}

// Java: CodingConvention.Bind(Node, Node, Node)
func NewCodingConvention_Bind(target rhino.Node, thisValue rhino.Node, parameters rhino.Node) CodingConvention_Bind {
	return CodingConvention_Bind{
		target:     target,
		thisValue:  thisValue,
		parameters: parameters,
	}
}

// GetBoundParameterCount is the number of parameters bound (not including
// the 'this' value).
//
// Java: CodingConvention.Bind.getBoundParameterCount()
func (b CodingConvention_Bind) GetBoundParameterCount() int {
	if b.parameters.IsNil() {
		return 0
	}
	paramParent := b.parameters.GetParent()
	return paramParent.GetChildCount() - paramParent.GetIndexOfChild(b.parameters)
}

// CodingConvention_Cache is the Cache class.
type CodingConvention_Cache struct {
	cacheObj rhino.Node
	key      rhino.Node
	valueFn  rhino.Node
	keyFn    rhino.Node
}

// Java: CodingConvention.Cache(Node, Node, Node, Node)
func NewCodingConvention_Cache(cacheObj rhino.Node, key rhino.Node, valueFn rhino.Node, keyFn rhino.Node) CodingConvention_Cache {
	return CodingConvention_Cache{
		cacheObj: cacheObj,
		key:      key,
		valueFn:  valueFn,
		keyFn:    keyFn,
	}
}

// CodingConvention_SubclassType specifies the kind of inheritance.
type CodingConvention_SubclassType int32

const (
	CodingConvention_SubclassType_INHERITS CodingConvention_SubclassType = iota
	CodingConvention_SubclassType_MIXIN
)

var codingConvention_SubclassType_names = [...]string{"INHERITS", "MIXIN"}

// Name is Enum.name().
func (t CodingConvention_SubclassType) Name() string { return codingConvention_SubclassType_names[t] }

// String is Enum.toString().
func (t CodingConvention_SubclassType) String() string { return t.Name() }

// CodingConvention_SubclassRelationship records subclass relations.
type CodingConvention_SubclassRelationship struct {
	type_          CodingConvention_SubclassType
	subclassName   string
	superclassName string
}

// Java: CodingConvention.SubclassRelationship(SubclassType, Node, Node)
func NewCodingConvention_SubclassRelationship(
	type_ CodingConvention_SubclassType, subclassNode rhino.Node, superclassNode rhino.Node) *CodingConvention_SubclassRelationship {
	subName, subOk := subclassNode.GetQualifiedNameString()
	if !subOk {
		jutil.CheckArgumentMsg1(false, "Expected qualified name, found: %s", subclassNode)
	}
	supName, supOk := superclassNode.GetQualifiedNameString()
	if !supOk {
		jutil.CheckArgumentMsg1(false, "Expected qualified name, found: %s", superclassNode)
	}
	return &CodingConvention_SubclassRelationship{
		type_:          type_,
		subclassName:   subName,
		superclassName: supName,
	}
}

// CodingConvention_ObjectLiteralCast: an object literal cast provides a
// mechanism to cast object literals to other types without a warning.
type CodingConvention_ObjectLiteralCast struct {
	// Type to cast to.
	typeName string

	// Object to cast.
	objectNode rhino.Node

	// Error message
	diagnosticType *DiagnosticType
}

// Java: CodingConvention.ObjectLiteralCast(String, Node, DiagnosticType)
func NewCodingConvention_ObjectLiteralCast(
	typeName string, objectNode rhino.Node, diagnosticType *DiagnosticType) CodingConvention_ObjectLiteralCast {
	return CodingConvention_ObjectLiteralCast{
		typeName:       typeName,
		objectNode:     objectNode,
		diagnosticType: diagnosticType,
	}
}

// CodingConvention_AssertionFunctionSpec describes a JavaScript function that
// throws when either one of its parameters does not match its return type or
// one of its parameters is falsy.
type CodingConvention_AssertionFunctionSpec struct {
	functionName     string
	hasFunctionName  bool
	qname            rhino.QualifiedName
	closurePrimitive rhino.ClosurePrimitive
	hasPrimitive     bool
	assertionKind    CodingConvention_AssertionFunctionSpec_AssertionKind
	paramIndex       int
}

func newAssertionSpecPrimitive(
	primitive rhino.ClosurePrimitive,
	kind CodingConvention_AssertionFunctionSpec_AssertionKind,
) *CodingConvention_AssertionFunctionSpec {
	return &CodingConvention_AssertionFunctionSpec{
		closurePrimitive: primitive,
		hasPrimitive:     true,
		assertionKind:    kind,
	}
}

func newAssertionSpecFunc(
	name string,
	kind CodingConvention_AssertionFunctionSpec_AssertionKind,
) *CodingConvention_AssertionFunctionSpec {
	var qname rhino.QualifiedName
	if strings.IndexByte(name, '.') >= 0 {
		qname = rhino.QualifiedName_Of(name)
	}
	return &CodingConvention_AssertionFunctionSpec{
		functionName:    name,
		hasFunctionName: true,
		qname:           qname,
		assertionKind:   kind,
	}
}

func (s *CodingConvention_AssertionFunctionSpec) GetAssertionKind() CodingConvention_AssertionFunctionSpec_AssertionKind {
	return s.assertionKind
}

// GetFunctionName returns the assertion function's qualified name if present.
//
// Java: CodingConvention.AssertionFunctionSpec.getFunctionName()
func (s *CodingConvention_AssertionFunctionSpec) GetFunctionName() base.OptString {
	if s.hasFunctionName {
		return base.SomeString(s.functionName)
	}
	return base.NoString()
}

// GetParamIndex returns the zero-based index of the asserted parameter.
//
// Java: CodingConvention.AssertionFunctionSpec.getParamIndex()
func (s *CodingConvention_AssertionFunctionSpec) GetParamIndex() int {
	return s.paramIndex
}

// Equals reports whether o is an equivalent CodingConvention_AssertionFunctionSpec.
//
// Java: CodingConvention.AssertionFunctionSpec.equals(Object)
func (s *CodingConvention_AssertionFunctionSpec) Equals(o any) bool {
	that, ok := o.(*CodingConvention_AssertionFunctionSpec)
	return ok && s.equals(that)
}

// CodingConvention_AssertionFunctionSpec_Builder builds a
// CodingConvention_AssertionFunctionSpec.
//
// Java: CodingConvention.AssertionFunctionSpec.Builder
type CodingConvention_AssertionFunctionSpec_Builder struct {
	functionName     string
	hasFunctionName  bool
	closurePrimitive rhino.ClosurePrimitive
	hasPrimitive     bool
	assertionKind    CodingConvention_AssertionFunctionSpec_AssertionKind
	hasAssertionKind bool
	paramIndex       int
	hasParamIndex    bool
}

// CodingConvention_AssertionFunctionSpec_ForTruthy returns a builder for an
// assertion function that checks truthiness of its parameter (default index 0).
//
// Java: CodingConvention.AssertionFunctionSpec.forTruthy()
func CodingConvention_AssertionFunctionSpec_ForTruthy() *CodingConvention_AssertionFunctionSpec_Builder {
	return &CodingConvention_AssertionFunctionSpec_Builder{
		assertionKind:    CodingConvention_AssertionFunctionSpec_AssertionKind_TRUTHY,
		hasAssertionKind: true,
		paramIndex:       0,
		hasParamIndex:    true,
	}
}

// CodingConvention_AssertionFunctionSpec_ForMatchesReturn returns a builder for
// an assertion function whose asserted parameter matches its return type.
//
// Java: CodingConvention.AssertionFunctionSpec.forMatchesReturn()
func CodingConvention_AssertionFunctionSpec_ForMatchesReturn() *CodingConvention_AssertionFunctionSpec_Builder {
	return &CodingConvention_AssertionFunctionSpec_Builder{
		assertionKind:    CodingConvention_AssertionFunctionSpec_AssertionKind_MATCHES_RETURN_TYPE,
		hasAssertionKind: true,
		paramIndex:       0,
		hasParamIndex:    true,
	}
}

// SetFunctionName sets the function name on the builder.
//
// Java: CodingConvention.AssertionFunctionSpec.Builder.setFunctionName(String)
func (b *CodingConvention_AssertionFunctionSpec_Builder) SetFunctionName(functionName string) *CodingConvention_AssertionFunctionSpec_Builder {
	b.functionName = functionName
	b.hasFunctionName = true
	return b
}

// SetClosurePrimitive sets the ClosurePrimitive on the builder.
//
// Java: CodingConvention.AssertionFunctionSpec.Builder.setClosurePrimitive(ClosurePrimitive)
func (b *CodingConvention_AssertionFunctionSpec_Builder) SetClosurePrimitive(primitive rhino.ClosurePrimitive) *CodingConvention_AssertionFunctionSpec_Builder {
	b.closurePrimitive = primitive
	b.hasPrimitive = true
	return b
}

// SetParamIndex sets the zero-based asserted parameter index on the builder.
//
// Java: CodingConvention.AssertionFunctionSpec.Builder.setParamIndex(int)
func (b *CodingConvention_AssertionFunctionSpec_Builder) SetParamIndex(paramIndex int) *CodingConvention_AssertionFunctionSpec_Builder {
	b.paramIndex = paramIndex
	b.hasParamIndex = true
	return b
}

// Build validates and constructs the CodingConvention_AssertionFunctionSpec.
//
// Java: CodingConvention.AssertionFunctionSpec.Builder.build()
func (b *CodingConvention_AssertionFunctionSpec_Builder) Build() *CodingConvention_AssertionFunctionSpec {
	if !b.hasAssertionKind || !b.hasParamIndex {
		var missing strings.Builder
		if !b.hasAssertionKind {
			missing.WriteString(" assertionKind")
		}
		if !b.hasParamIndex {
			missing.WriteString(" paramIndex")
		}
		panic(jutil.NewIllegalStateExceptionMsg("Missing required properties:" + missing.String()))
	}
	jutil.CheckStateMsg(
		b.hasFunctionName || b.hasPrimitive,
		"Must provide a function name or ClosurePrimitive for each spec")
	var qname rhino.QualifiedName
	if b.hasFunctionName && strings.IndexByte(b.functionName, '.') >= 0 {
		qname = rhino.QualifiedName_Of(b.functionName)
	}
	return &CodingConvention_AssertionFunctionSpec{
		functionName:     b.functionName,
		hasFunctionName:  b.hasFunctionName,
		qname:            qname,
		closurePrimitive: b.closurePrimitive,
		hasPrimitive:     b.hasPrimitive,
		assertionKind:    b.assertionKind,
		paramIndex:       b.paramIndex,
	}
}

// GetAssertedArg returns which argument is actually being asserted, or NilNode
// if fewer args than expected.
func (s *CodingConvention_AssertionFunctionSpec) GetAssertedArg(firstArg rhino.Node) rhino.Node {
	for i := 0; i < s.paramIndex; i++ {
		if firstArg.IsNil() {
			return rhino.NilNode
		}
		firstArg = firstArg.GetNext()
	}
	return firstArg
}

func (s *CodingConvention_AssertionFunctionSpec) equals(that *CodingConvention_AssertionFunctionSpec) bool {
	return s == that || (that != nil &&
		s.hasFunctionName == that.hasFunctionName &&
		s.functionName == that.functionName &&
		s.hasPrimitive == that.hasPrimitive &&
		s.closurePrimitive == that.closurePrimitive &&
		s.assertionKind == that.assertionKind &&
		s.paramIndex == that.paramIndex)
}

func (s *CodingConvention_AssertionFunctionSpec) String() string {
	functionName := "null"
	if s.hasFunctionName {
		functionName = s.functionName
	}
	closurePrimitive := "null"
	if s.hasPrimitive {
		closurePrimitive = s.closurePrimitive.String()
	}
	return "AssertionFunctionSpec{functionName=" + functionName +
		", closurePrimitive=" + closurePrimitive +
		", assertionKind=" + s.assertionKind.String() +
		", paramIndex=" + strconv.Itoa(s.paramIndex) + "}"
}

// CodingConvention_AssertionFunctionSpec_AssertionKind is
// AssertionFunctionSpec.AssertionKind.
type CodingConvention_AssertionFunctionSpec_AssertionKind int32

const (
	// an assertion that the parameter is 'truthy' and returns the parameter
	CodingConvention_AssertionFunctionSpec_AssertionKind_TRUTHY CodingConvention_AssertionFunctionSpec_AssertionKind = iota
	// an assertion that the parameter matches the inferred return kind
	CodingConvention_AssertionFunctionSpec_AssertionKind_MATCHES_RETURN_TYPE
)

var codingConvention_AssertionFunctionSpec_AssertionKind_names = [...]string{"TRUTHY", "MATCHES_RETURN_TYPE"}

func (k CodingConvention_AssertionFunctionSpec_AssertionKind) String() string {
	return codingConvention_AssertionFunctionSpec_AssertionKind_names[k]
}

// CodingConvention_AssertionFunctionLookup stores a relation from either
// name or Closure Primitive to assertion function.
type CodingConvention_AssertionFunctionLookup struct {
	byPrimitive  map[rhino.ClosurePrimitive]*CodingConvention_AssertionFunctionSpec
	byNameID     map[strpool.ID]*CodingConvention_AssertionFunctionSpec
	byLastPropID map[strpool.ID][]*CodingConvention_AssertionFunctionSpec
}

// CodingConvention_AssertionFunctionLookup_Of returns a new lookup containing
// all the given AssertionFunctionSpecs.
func CodingConvention_AssertionFunctionLookup_Of(
	specs []*CodingConvention_AssertionFunctionSpec) *CodingConvention_AssertionFunctionLookup {
	l := &CodingConvention_AssertionFunctionLookup{}
	for _, spec := range specs {
		jutil.CheckNotNullPtr(spec)
		switch {
		case spec.hasPrimitive:
			key := spec.closurePrimitive
			if l.byPrimitive == nil {
				l.byPrimitive = make(map[rhino.ClosurePrimitive]*CodingConvention_AssertionFunctionSpec)
			}
			if prev, exists := l.byPrimitive[key]; exists {
				panic(jutil.NewIllegalArgumentExceptionMsg(
					"Multiple entries with same key: " + key.String() + "=" + spec.String() + " and " + key.String() + "=" + prev.String()))
			}
			l.byPrimitive[key] = spec
		case spec.hasFunctionName:
			key := spec.functionName
			if lastDot := strings.LastIndexByte(key, '.'); lastDot >= 0 {
				if l.byLastPropID == nil {
					l.byLastPropID = make(map[strpool.ID][]*CodingConvention_AssertionFunctionSpec)
				}
				lastPropID := strpool.Intern(key[lastDot+1:])
				for _, prev := range l.byLastPropID[lastPropID] {
					if prev.hasFunctionName && prev.functionName == key {
						panic(jutil.NewIllegalArgumentExceptionMsg(
							"Multiple entries with same key: " + key + "=" + spec.String() + " and " + key + "=" + prev.String()))
					}
				}
				l.byLastPropID[lastPropID] = append(l.byLastPropID[lastPropID], spec)
			} else {
				if l.byNameID == nil {
					l.byNameID = make(map[strpool.ID]*CodingConvention_AssertionFunctionSpec)
				}
				nameID := strpool.Intern(key)
				if prev, exists := l.byNameID[nameID]; exists {
					panic(jutil.NewIllegalArgumentExceptionMsg(
						"Multiple entries with same key: " + key + "=" + spec.String() + " and " + key + "=" + prev.String()))
				}
				l.byNameID[nameID] = spec
			}
		default:
			panic(jutil.NewNullPointerException())
		}
	}
	return l
}

// LookupByCallee returns the AssertionFunctionSpec matching the given
// function reference, first by ClosurePrimitive, then by qualified name.
func (l *CodingConvention_AssertionFunctionLookup) LookupByCallee(callee rhino.Node) *CodingConvention_AssertionFunctionSpec {
	fnType := rhino.JSType_ToMaybeFunctionType(callee.GetJSType())
	if fnType != nil {
		if prim, ok := fnType.GetClosurePrimitive(); ok {
			if spec := l.byPrimitive[prim]; spec != nil {
				return spec
			}
		}
	}
	if callee.IsName() {
		return l.byNameID[nodeStringID(callee)]
	}
	if callee.IsGetProp() && len(l.byLastPropID) > 0 {
		for _, spec := range l.byLastPropID[nodeStringID(callee)] {
			if spec.qname != nil && spec.qname.Matches(callee) {
				return spec
			}
		}
	}
	return nil
}
