// Port of com/google/javascript/jscomp/TypeValidator.java @ closure-compiler v20261021

/*
 * Copyright 2009 The Closure Compiler Authors.
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
	"slices"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ExpectUndeclaredVariable expects that the given variable has not been
// declared with a type.
//
// sourceName is the name of the source file we're in; n is the node where
// warnings should point to; parent is the parent of n; var_ is the variable
// that we're checking; variableName is the name of the variable; newType is
// the type being applied to the variable. Mostly just here for the benefit
// of the warning. Returns the variable we end up with. Most of the time,
// this will just be var, but in some rare cases we will need to declare a
// new var with new source info.
//
// Java: TypeValidator.expectUndeclaredVariable(String, CompilerInput, Node, Node, TypedVar, String, JSType)
func (v *TypeValidator) ExpectUndeclaredVariable(
	sourceName string,
	input *CompilerInput,
	n rhino.Node,
	parent rhino.Node,
	var_ *TypedVar,
	variableName string,
	newType rhino.JSType) *TypedVar {
	newVar := var_
	varType := var_.GetType()

	// Only report duplicate declarations that have types. Other duplicates
	// will be reported by the syntactic scope creator later in the
	// compilation process.
	if varType != nil &&
		varType != v.typeRegistry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE) &&
		newType != nil &&
		newType != v.typeRegistry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE) {
		// If there are two typed declarations of the same variable, that
		// is an error and the second declaration is ignored, except in the
		// case of native types. A null input type means that the declaration
		// was made in TypedScopeCreator#createInitialScope and is a
		// native type. We should redeclare it at the new input site.
		if var_.GetInput() == nil {
			// Go-only: language-adapter — getScope() returns S (TypedScope)
			// in Java; the Go method keeps rhino.StaticScope (see AbstractVar).
			s := var_.GetScope().(*TypedScope)
			s.Undeclare(var_)
			newVar = s.Declare(variableName, n, varType, input, false)

			n.SetJSType(varType)
			if parent.IsVar() {
				if n.HasChildren() {
					n.GetFirstChild().SetJSType(varType)
				}
			} else if parent.IsExprResult() {
				n.SetJSType(varType)
			} else {
				jutil.CheckState(parent.IsFunction() || parent.IsClass())
				parent.SetJSType(varType)
			}
		} else {
			// Check for @suppress duplicate or similar warnings guard on the previous variable
			// declaration location.
			allowDupe := TypeValidator_HasDuplicateDeclarationSuppression(v.compiler, var_.GetNameNode())
			// If the previous definition doesn't suppress the warning, emit it here (i.e. always emit
			// on the second of the duplicate definitions). The warning might still be suppressed by an
			// @suppress tag on this declaration.
			if !allowDupe {
				// Report specifically if it is not just a duplicate, but types also don't mismatch.
				if !newType.Equals(varType) {
					v.compiler.Report(
						JSError_MakeWithNode(
							n,
							TypeValidator_DUP_VAR_DECLARATION_TYPE_MISMATCH,
							variableName,
							newType.String(),
							var_.GetInputName(),
							strconv.Itoa(var_.GetNameNode().GetLineno()),
							varType.String()))
				} else if !var_.GetParentNode().IsExprResult() {
					// If the type matches and the previous declaration was a stub declaration
					// (isExprResult), then ignore the duplicate, otherwise emit an error.
					v.compiler.Report(
						JSError_MakeWithNode(
							n,
							TypeValidator_DUP_VAR_DECLARATION,
							variableName,
							var_.GetInputName(),
							strconv.Itoa(var_.GetNameNode().GetLineno())))
				}
			}
		}
	}

	return newVar
}

// ExpectAllInterfaceProperties expects that all properties on interfaces
// that this type implements are implemented and correctly typed.
//
// Java: TypeValidator.expectAllInterfaceProperties(Node, FunctionType)
func (v *TypeValidator) ExpectAllInterfaceProperties(n rhino.Node, type_ rhino.FunctionType) {
	instance := type_.GetInstanceType()
	for _, implemented := range type_.GetAllImplementedInterfaces() {
		v.ExpectInterfaceProperties(n, instance, implemented)
	}
	for _, extended := range type_.GetExtendedInterfaces() {
		v.ExpectInterfaceProperties(n, instance, extended)
	}
}

// Java: TypeValidator.expectInterfaceProperties(Node, ObjectType, ObjectType)
func (v *TypeValidator) ExpectInterfaceProperties(
	n rhino.Node, instance rhino.ObjectType, ancestorInterface rhino.ObjectType) {
	// Case: `/** @interface */ class Foo { constructor() { this.prop; } }`
	for _, prop := range ancestorInterface.GetOwnPropertyKeys() {
		v.ExpectInterfaceProperty(n, instance, ancestorInterface, prop)
	}
	if ancestorInterface.GetImplicitPrototype() != nil {
		// Case: `/** @interface */ class Foo { prop() { } }`
		for _, prop := range ancestorInterface.GetImplicitPrototype().GetOwnPropertyKeys() {
			v.ExpectInterfaceProperty(n, instance, ancestorInterface, prop)
		}
	}
}

// ExpectInterfaceProperty expects that the property in an interface that
// this type implements is implemented and correctly typed.
//
// Java: TypeValidator.expectInterfaceProperty(Node, ObjectType, ObjectType, Property.Key)
func (v *TypeValidator) ExpectInterfaceProperty(
	n rhino.Node, instance rhino.ObjectType, implementedInterface rhino.ObjectType, propName rhino.Property_Key) {
	propSlot := instance.FindClosestDefinitionByKey(propName)
	if propSlot == nil ||
		(!instance.GetConstructor().IsInterface() && propSlot.IsOwnedByInterface()) {

		if instance.GetConstructor().IsAbstract() || instance.GetConstructor().IsInterface() {
			// Abstract classes and interfaces are not required to implement interface properties.
			return
		}
		if implementedInterface.GetPropertyTypeByKey(propName).IsVoidable() {
			// Voidable properties don't require explicit initialization in type constructors.
			return
		}

		v.RegisterMismatchAndReport(
			instance,
			implementedInterface,
			JSError_MakeWithNode(
				n,
				TypeValidator_INTERFACE_METHOD_NOT_IMPLEMENTED,
				propName.HumanReadableName(),
				jutil.ToJavaString(implementedInterface.GetReferenceName()),
				instance.String()))
	} else {
		local := propSlot.GetOwnerInstanceType().Equals(instance)
		if !local && instance.GetConstructor().IsInterface() {
			// non-local interface mismatches already handled via interface property conflict checks.
			return
		}

		prop := propSlot.GetValue()
		var propNode rhino.Node
		if prop.GetDeclaration() != nil {
			propNode = prop.GetDeclaration().GetNode()
		}
		// Fall back on the constructor node if we can't find a node for the property.
		if propNode.IsNil() {
			propNode = n
		}

		v.CheckPropertyType(propNode, instance, implementedInterface, propName, prop.GetType())
	}
}

// CheckPropertyType checks the property is correctly typed (i.e. subtype of
// the parent property's type declaration).
//
// Java: TypeValidator.checkPropertyType(Node, JSType, ObjectType, Property.Key, JSType)
func (v *TypeValidator) CheckPropertyType(
	n rhino.Node, instance rhino.JSType, parent rhino.ObjectType, propertyName rhino.Property_Key, found rhino.JSType) {
	required := parent.GetPropertyTypeByKey(propertyName)
	typeMap := instance.GetTemplateTypeMap()
	if !typeMap.IsEmpty() && required.HasAnyTemplateTypes() {
		replacer := rhino.TemplateTypeReplacer_MakeForPartialReplacement(v.typeRegistry, typeMap)
		required = replacer.VisitJSType(required)
	}

	if found.IsSubtypeOfWithMode(required, v.subtypingMode) {
		return
	}

	// Implemented, but not correctly typed
	var diagnostic *DiagnosticType
	if parent.GetConstructor().IsInterface() {
		diagnostic = TypeValidator_HIDDEN_INTERFACE_PROPERTY_MISMATCH
	} else {
		diagnostic = TypeValidator_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH
	}
	err := JSError_MakeWithNode(
		n,
		diagnostic,
		propertyName.HumanReadableName(),
		jutil.ToJavaString(parent.GetReferenceName()),
		required.String(),
		found.String(),
		instance.String())
	v.RegisterMismatchAndReport(found, required, err)
}

// ExpectAbstractMethodsImplemented: for a concrete class, expect that all
// abstract methods that haven't been implemented by any of the super classes
// on the inheritance chain are implemented.
//
// Java: TypeValidator.expectAbstractMethodsImplemented(Node, FunctionType)
func (v *TypeValidator) ExpectAbstractMethodsImplemented(n rhino.Node, ctorType rhino.FunctionType) {
	jutil.CheckArgument(ctorType.IsConstructor())

	currSuperCtor := ctorType.GetSuperClassConstructor()
	if currSuperCtor == nil || !currSuperCtor.IsAbstract() {
		return
	}

	var abstractMethodSuperTypeMap propertyKeyMap

	for currSuperCtor != nil && currSuperCtor.IsAbstract() {
		superType := currSuperCtor.GetInstanceType()
		for _, prop := range currSuperCtor.GetPrototype().GetOwnPropertyKeys() {
			maybeAbstractMethod := superType.FindPropertyTypeByKey(prop).ToMaybeFunctionType()
			if maybeAbstractMethod != nil &&
				maybeAbstractMethod.IsAbstract() &&
				!abstractMethodSuperTypeMap.containsKey(prop) {
				abstractMethodSuperTypeMap.put(prop, superType)
			}
		}
		currSuperCtor = currSuperCtor.GetSuperClassConstructor()
	}

	instance := ctorType.GetInstanceType()
	for _, entry := range abstractMethodSuperTypeMap.entries {
		method, superType := entry.key, entry.val
		abstractMethod := instance.FindPropertyTypeByKey(method).ToMaybeFunctionType()
		if abstractMethod == nil || abstractMethod.IsAbstract() {
			v.RegisterMismatchAndReport(
				instance,
				superType,
				JSError_MakeWithNode(
					n,
					TypeValidator_ABSTRACT_METHOD_NOT_IMPLEMENTED,
					method.HumanReadableName(),
					superType.String(),
					instance.String()))
		}
	}
}

// Java: TypeValidator.mismatch(Node, String, JSType, JSTypeNative)
func (v *TypeValidator) MismatchNative(n rhino.Node, msg string, found rhino.JSType, required rhino.JSTypeNative) {
	v.Mismatch(n, msg, found, v.GetNativeType(required))
}

// Java: TypeValidator.mismatch(Node, String, JSType, JSType)
func (v *TypeValidator) Mismatch(n rhino.Node, msg string, found rhino.JSType, required rhino.JSType) {
	if !found.IsSubtypeOfWithMode(required, v.subtypingMode) {
		var missing []string
		var mismatch []string
		var structuralType rhino.JSType
		if required.IsStructuralType() {
			structuralType = required
		} else if required.IsUnionType() {
			restrictedType := required.RestrictByNotNullOrUndefined()
			if restrictedType.IsStructuralType() {
				// restrict to unions created by undefined `{!Type=}` type or by nullable `{?Type}` type.
				structuralType = restrictedType
			}
		}

		if structuralType != nil {
			requiredObject := structuralType.ToMaybeObjectType()
			foundObject := found.ToMaybeObjectType()
			if requiredObject != nil && foundObject != nil {
				for _, property := range requiredObject.GetPropertyNames() {
					propRequired := requiredObject.GetPropertyType(property)
					hasProperty := foundObject.HasProperty(property)
					if !propRequired.IsExplicitlyVoidable() || hasProperty {
						if hasProperty {
							if !foundObject.
								GetPropertyType(property).
								IsSubtypeOfWithMode(propRequired, v.subtypingMode) {
								mismatch = append(mismatch, property)
							}
						} else {
							missing = append(missing, property)
						}
					}
				}
				if len(missing) > 1 {
					slices.SortFunc(missing, jutil.StringCompareTo)
					missing = slices.Compact(missing)
				}
				if len(mismatch) > 1 {
					slices.SortFunc(mismatch, jutil.StringCompareTo)
					mismatch = slices.Compact(mismatch)
				}
			}
		}
		v.RegisterMismatchAndReportWithDiagnostic(n, TypeValidator_TYPE_MISMATCH_WARNING, msg, found, required, missing, mismatch)
	}
}

// RegisterMismatchAndReportWithDiagnostic is used both for
// TYPE_MISMATCH_WARNING and INVALID_OPERAND_TYPE.
//
// Java: TypeValidator.registerMismatchAndReport(Node, DiagnosticType, String, JSType, JSType, Set<String>, Set<String>)
func (v *TypeValidator) RegisterMismatchAndReportWithDiagnostic(
	n rhino.Node,
	diagnostic *DiagnosticType,
	msg string,
	found rhino.JSType,
	required rhino.JSType,
	missing []string,
	mismatch []string) {
	foundRequiredFormatted := TypeValidator_FormatFoundRequired(msg, found, required, missing, mismatch)
	err := JSError_MakeWithNode(n, diagnostic, foundRequiredFormatted)
	v.RegisterMismatchAndReport(found, required, err)
}

// RegisterMismatchAndReport registers a type mismatch into the universe of
// mismatches owned by this pass.
//
// Java: TypeValidator.registerMismatchAndReport(JSType, JSType, JSError)
func (v *TypeValidator) RegisterMismatchAndReport(found rhino.JSType, required rhino.JSType, error *JSError) {
	v.compiler.Report(error)
	v.mismatches.RegisterMismatch(error.Node(), found, required)
}

// TypeValidator_FormatNodeLocation is
// String.format("%s:%s:%s", node.getSourceFileName(), node.getLineno(), node.getCharno()).
//
// Java: TypeValidator.formatNodeLocation(Node)
func TypeValidator_FormatNodeLocation(node rhino.Node) string {
	return jutil.ToJavaString(node.GetSourceFileName()) + ":" +
		strconv.Itoa(node.GetLineno()) + ":" + strconv.Itoa(node.GetCharno())
}

// TypeValidator_FormatFoundRequired formats a found/required error message.
//
// Java: TypeValidator.formatFoundRequired(String, JSType, JSType, Set<String>, Set<String>)
func TypeValidator_FormatFoundRequired(
	description string,
	found rhino.JSType,
	required rhino.JSType,
	missing []string,
	mismatch []string) string {
	foundStr := found.String()
	requiredStr := required.String()
	if foundStr == requiredStr {
		foundStr = found.ToAnnotationString(rhino.JSType_Nullability_IMPLICIT)
		requiredStr = required.ToAnnotationString(rhino.JSType_Nullability_IMPLICIT)
	}
	if foundStr == requiredStr {
		// For some types, their annotation strings will be identical, but they might be identically
		// named structures in different scopes (and have theoretically different "fully-qualified"
		// names).
		// The general case of giving a correct fully qualified name is quite difficult, so instead we
		// address some common cases where we know that types have a concrete name to use.
		if found.IsEnumElementType() {
			foundStr = foundStr +
				" (enum definition: " +
				TypeValidator_FormatNodeLocation(found.ToMaybeEnumElementType().GetEnumType().GetSource()) +
				")"
		}
		if required.IsEnumElementType() {
			requiredStr = requiredStr +
				" (enum definition: " +
				TypeValidator_FormatNodeLocation(required.ToMaybeEnumElementType().GetEnumType().GetSource()) +
				")"
		}
	}
	missingStr := strings.Join(missing, ",")
	mismatchStr := strings.Join(mismatch, ",")
	if len(missingStr) > 0 || len(mismatchStr) > 0 {
		return description + "\n" +
			"found   : " + foundStr + "\n" +
			"required: " + requiredStr + "\n" +
			"missing : [" + missingStr + "]\n" +
			"mismatch: [" + mismatchStr + "]"
	}
	return description + "\n" +
		"found   : " + foundStr + "\n" +
		"required: " + requiredStr
}

// GetJSType gets the JSType from the Node argument and verifies that it is
// present.
//
// Java: TypeValidator.getJSType(Node)
func (v *TypeValidator) GetJSType(n rhino.Node) rhino.JSType {
	return jutil.CheckNotNullMsg1(n.GetJSType(), "%s has no JSType attached", n)
}

// Java: TypeValidator.getNativeType(JSTypeNative)
func (v *TypeValidator) GetNativeType(typeId rhino.JSTypeNative) rhino.JSType {
	return v.typeRegistry.GetNativeType(typeId)
}

// TypeValidator_HasDuplicateDeclarationSuppression reports whether
// duplicated declarations warnings should be suppressed for the given node.
//
// decl is the declaration to check.
//
// Java: TypeValidator.hasDuplicateDeclarationSuppression(AbstractCompiler, Node)
func TypeValidator_HasDuplicateDeclarationSuppression(compiler AbstractCompiler, decl rhino.Node) bool {
	// NB: DUP_VAR_DECLARATION is somewhat arbitrary here, but it must be one of the errors
	// suppressed by the "duplicate" group.
	jutil.CheckNotNull(compiler.options)
	var lazyErr *JSError
	level := dupVarDeclGuardLevel(compiler.warningsGuard, decl, &lazyErr)
	return level != nil && *level == CheckLevel_OFF
}

func makeDupVarDeclProbeError(decl rhino.Node) *JSError {
	return &JSError{
		type_:        TypeValidator_DUP_VAR_DECLARATION,
		description:  TypeValidator_DUP_VAR_DECLARATION.Format("dummy", "dummy"),
		sourceName:   decl.GetSourceFileName(),
		lineno:       decl.GetLineno(),
		charno:       decl.GetCharno(),
		length:       decl.GetLength(),
		node:         decl,
		defaultLevel: TypeValidator_DUP_VAR_DECLARATION.level,
	}
}

func dupVarDeclGuardLevel(guard WarningsGuard, decl rhino.Node, lazyErr **JSError) *CheckLevel {
	switch g := guard.(type) {
	case *ComposeWarningsGuard:
		for _, entry := range g.entries {
			newLevel := dupVarDeclGuardLevel(entry.guard, decl, lazyErr)
			if newLevel != nil {
				if g.demoteErrors && *newLevel == CheckLevel_ERROR {
					return CheckLevel_WARNING.Ptr()
				}
				return newLevel
			}
		}
		return nil
	case *ClosureUnawareCodeWarningsGuard:
		return nil
	case *J2clSuppressWarningsGuard:
		if !J2clSourceUtils_IsJ2clSource(decl.GetSourceFileName()) {
			return nil
		}
		if J2clSuppressWarningsGuard_DEFAULT_J2CL_SUPRRESIONS.MatchesType(TypeValidator_DUP_VAR_DECLARATION) {
			return CheckLevel_OFF.Ptr()
		}
		return nil
	case *DiagnosticGroupWarningsGuard:
		if g.group.MatchesType(TypeValidator_DUP_VAR_DECLARATION) {
			return g.level.Ptr()
		}
		return nil
	case *DiagnosticGroupPathSuppressingWarningsGuard:
		src := decl.GetSourceFileName()
		if src == nil || !strings.Contains(*src, g.part) {
			return nil
		}
		if g.group.MatchesType(TypeValidator_DUP_VAR_DECLARATION) {
			return CheckLevel_OFF.Ptr()
		}
		return nil
	case *ByPathWarningsGuard:
		if TypeValidator_DUP_VAR_DECLARATION.level != CheckLevel_ERROR {
			if errorPath := decl.GetSourceFileName(); errorPath != nil {
				src := *errorPath
				inPath := false
				for _, path := range g.paths {
					if strings.Contains(src, path) {
						inPath = true
						break
					}
				}
				if inPath == g.include {
					return g.level.Ptr()
				}
			}
		}
		return nil
	case *ShowByPathWarningsGuard:
		return dupVarDeclGuardLevel(g.warningsGuard, decl, lazyErr)
	case *CheckTypeSummaryWarningsGuard:
		if src := decl.GetSourceFileName(); src != nil && strings.HasSuffix(*src, ".i.js") {
			return g.level.Ptr()
		}
		return nil
	case *StrictWarningsGuard:
		if TypeValidator_DUP_VAR_DECLARATION.level.IsOn() {
			return CheckLevel_ERROR.Ptr()
		}
		return nil
	case *SuppressDocWarningsGuard:
		return dupVarDeclSuppressDocLevel(g, decl)
	default:
		if *lazyErr == nil {
			*lazyErr = makeDupVarDeclProbeError(decl)
		}
		return guard.Level(*lazyErr)
	}
}

func dupVarDeclSuppressDocLevel(g *SuppressDocWarningsGuard, decl rhino.Node) *CheckLevel {
	node := decl
	src := decl.GetSourceFileName()
	if node.IsNil() && src != nil {
		if sn := g.compiler.GetScriptNode(*src); !sn.IsNil() {
			jutil.CheckState(sn.IsScript())
			node = sn
		}
	}
	if node.IsNil() {
		return nil
	}
	for current := node; !current.IsNil(); current = current.GetParent() {
		var info *rhino.JSDocInfo
		if current.IsFunction() || current.IsClass() {
			info = NodeUtil_GetBestJSDocInfo(current)
		} else if current.IsScript() {
			info = current.GetJSDocInfo()
		} else if NodeUtil_IsNameDeclaration(current) ||
			NodeUtil_MayBeObjectLitKey(current) ||
			current.IsComputedProp() ||
			current.IsMemberFieldDef() ||
			current.IsComputedFieldDef() ||
			(current.HasParent() && current.GetParent().IsExprResult()) {
			info = NodeUtil_GetBestJSDocInfo(current)
		}
		if info != nil {
			if level := dupVarDeclLevelFromInfo(g, info); level != nil {
				return level
			}
		}
	}
	if src != nil {
		if scriptNode := g.compiler.GetScriptNode(*src); !scriptNode.IsNil() {
			jutil.CheckState(scriptNode.IsScript())
			if info := scriptNode.GetJSDocInfo(); info != nil {
				return dupVarDeclLevelFromInfo(g, info)
			}
		}
	}
	return nil
}

func dupVarDeclLevelFromInfo(g *SuppressDocWarningsGuard, info *rhino.JSDocInfo) *CheckLevel {
	var matched *CheckLevel
	info.ForEachSuppression(func(suppressor string) bool {
		group := g.suppressors[suppressor]
		if group == nil {
			return true
		}
		if group.MatchesType(TypeValidator_DUP_VAR_DECLARATION) {
			matched = CheckLevel_OFF.Ptr()
			return false
		}
		return true
	})
	return matched
}

// Java: TypeValidator.expectWellFormedTemplatizedType(Node)
func (v *TypeValidator) ExpectWellFormedTemplatizedType(n rhino.Node) {
	if nodeType := n.GetJSType(); nodeType != nil {
		v.checkWellFormedTemplatizedType(n, nodeType)
	}
}

// checkWellFormedTemplatizedType verifies bounded generic type arguments on t
// without allocating a visitor, caching positive results on resolved JSTypes.
//
// Java: TypeValidator.WellFormedTemplatizedTypeVerifier
func (v *TypeValidator) checkWellFormedTemplatizedType(n rhino.Node, t rhino.JSType) bool {
	tb := t.GetJSTypeBase()
	if tb.IsWellFormedTemplateValid() {
		return true
	}
	var valid bool
	switch x := t.(type) {
	case *rhino.EnumElementType:
		prim := x.GetPrimitiveType()
		valid = prim == nil || v.checkWellFormedTemplatizedType(n, prim)
	case *rhino.FunctionTypeBase:
		if x.IsNoObjectType() || x.IsNoType() {
			valid = true
		} else {
			valid = v.checkWellFormedFunctionType(n, x)
		}
	case *rhino.NamedType:
		valid = v.checkWellFormedTemplatizedType(n, x.GetReferencedType())
	case *rhino.UnionType:
		valid = true
		for _, alt := range x.GetAlternates() {
			if !v.checkWellFormedTemplatizedType(n, alt) {
				valid = false
				break
			}
		}
	case *rhino.TemplatizedType:
		valid = v.checkWellFormedTemplatizedInstance(n, x)
	case *rhino.TemplateType:
		valid = !x.ContainsCycle() && v.checkWellFormedTemplatizedType(n, x.GetBound())
	case *rhino.AllType, *rhino.BooleanType, *rhino.UnknownType, *rhino.NullType,
		*rhino.NumberType, *rhino.BigIntType, *rhino.StringType, *rhino.VoidType,
		*rhino.InstanceObjectType, *rhino.PrototypeObjectTypeBase, *rhino.RecordType,
		*rhino.EnumType, *rhino.NoResolvedType, *rhino.ProxyObjectTypeBase,
		*rhino.SymbolTypeBase, *rhino.KnownSymbolType:
		valid = true
	default:
		switch t.(type) {
		case rhino.ProxyObjectType, rhino.ObjectType:
			valid = true
		default:
			panic(jutil.NewUnsupportedOperationException())
		}
	}
	if valid {
		tb.MarkWellFormedTemplateValidIfResolved()
	}
	return valid
}

func (v *TypeValidator) checkWellFormedFunctionType(n rhino.Node, type_ rhino.FunctionType) bool {
	for _, param := range type_.GetParameters() {
		if !v.checkWellFormedTemplatizedType(n, param.GetJSType()) {
			return false
		}
	}
	return v.checkWellFormedTemplatizedType(n, type_.GetReturnType())
}

func (v *TypeValidator) checkWellFormedTemplatizedInstance(n rhino.Node, type_ *rhino.TemplatizedType) bool {
	referencedTemplates := type_.GetReferencedType().GetTemplateTypeMap().GetTemplateKeys()
	assignedTypes := type_.GetTemplateTypes()
	for i, assignedType := range assignedTypes {
		if i < len(referencedTemplates) {
			templateType := referencedTemplates[i]
			if !assignedType.IsSubtype(templateType.GetBound()) {
				v.RegisterMismatchAndReport(
					assignedType,
					templateType.GetBound(),
					JSError_MakeWithNode(
						n,
						RhinoErrorReporter_BOUNDED_GENERIC_TYPE_ERROR,
						assignedType.String(),
						jutil.ToJavaString(templateType.GetReferenceName()),
						templateType.GetBound().String()))
				return false
			}
		}
	}
	return true
}
