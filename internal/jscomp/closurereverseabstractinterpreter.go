// Port of com/google/javascript/jscomp/type/ClosureReverseAbstractInterpreter.java
// @ closure-compiler v20260415.
//
// A reverse abstract interpreter (RAI) for specific Closure patterns such
// as goog.isObject. Extends ChainableReverseAbstractInterpreter.
//
// Java's restricters map uses a Function<TypeRestriction,JSType> lambda; Go
// uses a typeRestrictor func type. The two anonymous Visitor fields in Java
// are ported as package-private types in this file.

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

// typeRestriction holds a (type, outcome) pair passed to a restricter.
//
// Port of ClosureReverseAbstractInterpreter.TypeRestriction (private static inner class).
type typeRestriction struct {
	typ     jstype.JSType
	outcome bool
}

// typeRestrictor is the Go equivalent of Java's Function<TypeRestriction,JSType>.
type typeRestrictor func(r typeRestriction) jstype.JSType

// ClosureReverseAbstractInterpreter is a chainable RAI for Closure-specific
// guard functions such as goog.isObject. Currently recognises "isObject"
// (i.e. the callee is a goog.X GETPROP whose member is "isObject").
//
// Port of com/google/javascript/jscomp/type/ClosureReverseAbstractInterpreter.java.
type ClosureReverseAbstractInterpreter struct {
	ChainableReverseAbstractInterpreter

	// restricters maps Closure function names (e.g. "isObject") to
	// the corresponding type-restriction function.
	//
	// Port of ClosureReverseAbstractInterpreter.restricters field.
	restricters map[string]typeRestrictor
}

// Go-only: constructor (Java uses new ClosureReverseAbstractInterpreter(typeRegistry))
// NewClosureReverseAbstractInterpreter constructs a ClosureRAI and wires up
// the restricters map.
//
// Port of ClosureReverseAbstractInterpreter(JSTypeRegistry).
func NewClosureReverseAbstractInterpreter(reg *jstype.Registry) *ClosureReverseAbstractInterpreter {
	c := &ClosureReverseAbstractInterpreter{}
	c.ChainableReverseAbstractInterpreter = newChainableBase(reg, c)

	c.restricters = map[string]typeRestrictor{
		"isObject": func(r typeRestriction) jstype.JSType {
			if r.typ == nil {
				if r.outcome {
					return c.getNativeType(jstype.OBJECT_TYPE)
				}
				return nil
			}
			if r.outcome {
				return asJSType(jstype.VisitJSType(c.newRestrictToObjectVisitor(), r.typ))
			}
			return asJSType(jstype.VisitJSType(c.newRestrictToNotObjectVisitor(), r.typ))
		},
	}
	return c
}

// GetPreciserScopeKnowingConditionOutcome is the main entry point.
// It recognises call expressions of the form goog.X(param) where X is a
// registered Closure guard name, restricts the parameter's type, and returns
// an updated FlowScope. Unrecognised conditions are delegated to the next
// link in the chain.
//
// Port of ClosureReverseAbstractInterpreter.getPreciserScopeKnowingConditionOutcome.
func (c *ClosureReverseAbstractInterpreter) GetPreciserScopeKnowingConditionOutcome(
	condition *rhino.Node, blindScope *FlowScope, outcome bool,
) *FlowScope {
	if condition.IsCall() && len(condition.Children()) == 2 {
		callee := condition.FirstChild()
		param := condition.LastChild()
		if callee.IsGetProp() && param.IsQualifiedName() {
			paramType := c.getTypeIfRefinable(param, blindScope)
			receiver := callee.FirstChild()
			if receiver.IsName() && receiver.String_() == "goog" {
				// callee.String_() on a GETPROP node returns the property name.
				memberName := callee.LastChild().String_()
				if restricter, ok := c.restricters[memberName]; ok {
					return c.restrictParameter(param, paramType, blindScope, restricter, outcome)
				}
			}
		}
	}
	return c.nextPreciserScopeKnowingConditionOutcome(condition, blindScope, outcome)
}

// restrictParameter applies restriction to paramType and, if non-nil,
// declares the refined type for the parameter node in a new FlowScope.
//
// Port of ClosureReverseAbstractInterpreter.restrictParameter (private).
func (c *ClosureReverseAbstractInterpreter) restrictParameter(
	parameter *rhino.Node,
	typ jstype.JSType,
	blindScope *FlowScope,
	restriction typeRestrictor,
	outcome bool,
) *FlowScope {
	restricted := restriction(typeRestriction{typ: typ, outcome: outcome})
	if restricted != nil {
		return c.declareNameInScope(blindScope, parameter, restricted)
	}
	return blindScope
}

// ----------------------------------------------------------------------------
// restrictToObjectVisitor — goog.isObject returns true
// ----------------------------------------------------------------------------

// closureRestrictToObjectVisitor restricts a type for when goog.isObject
// returns true. This includes functions, but not null.
//
// Port of ClosureReverseAbstractInterpreter.restrictToObjectVisitor field
// (anonymous RestrictByTrueTypeOfResultVisitor subclass).
type closureRestrictToObjectVisitor struct {
	jstype.WithDefaultCase
	c *ClosureReverseAbstractInterpreter
}

// Go-only: factory method for the anonymous restrictToObjectVisitor field (Java uses anonymous class init in constructor)
func (c *ClosureReverseAbstractInterpreter) newRestrictToObjectVisitor() *closureRestrictToObjectVisitor {
	v := &closureRestrictToObjectVisitor{c: c}
	// RestrictByTrueTypeOfResultVisitor: all unhandled leaf cases → nil.
	v.DefaultCase = func(typ jstype.JSType) any { return nil }
	return v
}

// Go-only: Visitor interface implementation (Java uses anonymous RestrictByTrueTypeOfResultVisitor subclass)
func (v *closureRestrictToObjectVisitor) CaseAllType() any {
	return v.c.getNativeType(jstype.NO_OBJECT_TYPE)
}

// Go-only: Visitor interface implementation (Java uses anonymous RestrictByTrueTypeOfResultVisitor subclass)
func (v *closureRestrictToObjectVisitor) CaseUnknownType() any {
	return v.c.getNativeType(jstype.NO_OBJECT_TYPE)
}

func (v *closureRestrictToObjectVisitor) CaseObjectType(typ *jstype.ObjectType) any { return typ }

func (v *closureRestrictToObjectVisitor) CaseFunctionType(typ *jstype.FunctionType) any { return typ }

// Go-only: Visitor interface implementation for union types (Java anonymous class inherits from RestrictByTrueTypeOfResultVisitor)
// Union, enum, named, proxy, templatized, template: shared helpers.
func (v *closureRestrictToObjectVisitor) CaseUnionType(typ *jstype.UnionType) any {
	return visitUnionRestrictingType(v, typ)
}

// Go-only: Visitor interface implementation (Java anonymous class inherits default)
func (v *closureRestrictToObjectVisitor) CaseEnumElementType(typ *jstype.EnumElementType) any {
	return visitEnumElement(v, typ)
}

// Go-only: Visitor interface implementation (Java anonymous class inherits default)
func (v *closureRestrictToObjectVisitor) CaseTemplatizedType(typ *jstype.TemplatizedType) any {
	return v.CaseObjectType(typ.Base)
}

// Go-only: Visitor interface implementation (Java anonymous class inherits default)
func (v *closureRestrictToObjectVisitor) CaseTemplateType(typ *jstype.TemplateType) any {
	_ = typ // template types are not object types; mirrors Java RestrictByTrueTypeOfResultVisitor default
	return nil
}

// Go-only: Visitor interface implementation (Java anonymous class inherits default)
func (v *closureRestrictToObjectVisitor) CaseNamedType(typ jstype.JSType) any {
	if nt, ok := typ.(*jstype.NamedType); ok {
		return visitNamedType(v, nt)
	}
	return nil
}

// Go-only: Visitor interface implementation (Java anonymous class inherits default)
func (v *closureRestrictToObjectVisitor) CaseProxyObjectType(typ jstype.JSType) any {
	if pt, ok := typ.(*jstype.ProxyObjectType); ok {
		return visitProxyObjectType(v, pt)
	}
	return nil
}

// Go-only: Visitor interface implementation (Java anonymous class inherits default)
func (v *closureRestrictToObjectVisitor) CaseNoType(typ jstype.JSType) any { return typ }

// ----------------------------------------------------------------------------
// restrictToNotObjectVisitor — goog.isObject returns false
// ----------------------------------------------------------------------------

// closureRestrictToNotObjectVisitor restricts a type for when goog.isObject
// returns false. Objects and functions are filtered out; primitive types remain.
//
// Port of ClosureReverseAbstractInterpreter.restrictToNotObjectVisitor field
// (anonymous RestrictByFalseTypeOfResultVisitor subclass).
type closureRestrictToNotObjectVisitor struct {
	jstype.WithDefaultCase
	c *ClosureReverseAbstractInterpreter
}

// Go-only: factory method for the anonymous restrictToNotObjectVisitor field (Java uses anonymous class init in constructor)
func (c *ClosureReverseAbstractInterpreter) newRestrictToNotObjectVisitor() *closureRestrictToNotObjectVisitor {
	v := &closureRestrictToNotObjectVisitor{c: c}
	// RestrictByFalseTypeOfResultVisitor: default keeps the type as-is.
	v.DefaultCase = func(typ jstype.JSType) any { return typ }
	return v
}

// Go-only: Visitor interface implementation (Java uses anonymous RestrictByFalseTypeOfResultVisitor subclass)
func (v *closureRestrictToNotObjectVisitor) CaseAllType() any {
	return v.c.getNativeType(jstype.ALL_TYPE)
}

// Go-only: Visitor interface implementation (Java uses anonymous RestrictByFalseTypeOfResultVisitor subclass)
func (v *closureRestrictToNotObjectVisitor) CaseUnknownType() any {
	return v.c.getNativeType(jstype.CHECKED_UNKNOWN_TYPE)
}

func (v *closureRestrictToNotObjectVisitor) CaseObjectType(typ *jstype.ObjectType) any {
	_ = typ // objects are eliminated when goog.isObject is false; mirrors Java anonymous class override
	return nil
}
func (v *closureRestrictToNotObjectVisitor) CaseFunctionType(typ *jstype.FunctionType) any {
	_ = typ // functions are eliminated when goog.isObject is false; mirrors Java anonymous class override
	return nil
}

// Go-only: explicit case (Java RestrictByFalseTypeOfResultVisitor base class handles via default)
// The following are RestrictByFalseTypeOfResultVisitor defaults — return native type.
func (v *closureRestrictToNotObjectVisitor) CaseNoObjectType() any {
	return v.c.getNativeType(jstype.NO_OBJECT_TYPE)
}

// Go-only: explicit case (Java RestrictByFalseTypeOfResultVisitor base class handles via default)
func (v *closureRestrictToNotObjectVisitor) CaseBooleanType() any {
	return v.c.getNativeType(jstype.BOOLEAN_TYPE)
}

// Go-only: explicit case (Java RestrictByFalseTypeOfResultVisitor base class handles via default)
func (v *closureRestrictToNotObjectVisitor) CaseNullType() any {
	return v.c.getNativeType(jstype.NULL_TYPE)
}

// Go-only: explicit case (Java RestrictByFalseTypeOfResultVisitor base class handles via default)
func (v *closureRestrictToNotObjectVisitor) CaseNumberType() any {
	return v.c.getNativeType(jstype.NUMBER_TYPE)
}
// Go-only: explicit case (Java RestrictByFalseTypeOfResultVisitor base class handles via default)
func (v *closureRestrictToNotObjectVisitor) CaseBigIntType() any {
	return v.c.getNativeType(jstype.BIGINT_TYPE)
}

// Go-only: explicit case (Java RestrictByFalseTypeOfResultVisitor base class handles via default)
func (v *closureRestrictToNotObjectVisitor) CaseStringType() any {
	return v.c.getNativeType(jstype.STRING_TYPE)
}

// Go-only: explicit case (Java RestrictByFalseTypeOfResultVisitor base class handles via default)
func (v *closureRestrictToNotObjectVisitor) CaseSymbolType(typ jstype.JSType) any { return typ }

// Go-only: explicit case (Java RestrictByFalseTypeOfResultVisitor base class handles via default)
func (v *closureRestrictToNotObjectVisitor) CaseVoidType() any {
	return v.c.getNativeType(jstype.VOID_TYPE)
}

// Go-only: Visitor interface implementation for union types (Java anonymous class inherits default)
// Union, enum, named, proxy, templatized, template: shared helpers.
func (v *closureRestrictToNotObjectVisitor) CaseUnionType(typ *jstype.UnionType) any {
	return visitUnionRestrictingType(v, typ)
}

// Go-only: Visitor interface implementation (Java anonymous class inherits default)
func (v *closureRestrictToNotObjectVisitor) CaseEnumElementType(typ *jstype.EnumElementType) any {
	return visitEnumElement(v, typ)
}

// Go-only: Visitor interface implementation (Java anonymous class inherits default)
func (v *closureRestrictToNotObjectVisitor) CaseTemplatizedType(typ *jstype.TemplatizedType) any {
	return v.CaseObjectType(typ.Base)
}

// Go-only: Visitor interface implementation (Java anonymous class inherits default)
func (v *closureRestrictToNotObjectVisitor) CaseTemplateType(typ *jstype.TemplateType) any {
	_ = typ // template types are not constrained by goog.isObject; mirrors RestrictByFalseTypeOfResultVisitor default
	return nil
}

// Go-only: Visitor interface implementation (Java anonymous class inherits default)
func (v *closureRestrictToNotObjectVisitor) CaseNamedType(typ jstype.JSType) any {
	if nt, ok := typ.(*jstype.NamedType); ok {
		return visitNamedType(v, nt)
	}
	return nil
}

// Go-only: Visitor interface implementation (Java anonymous class inherits default)
func (v *closureRestrictToNotObjectVisitor) CaseProxyObjectType(typ jstype.JSType) any {
	if pt, ok := typ.(*jstype.ProxyObjectType); ok {
		return visitProxyObjectType(v, pt)
	}
	return nil
}

// Go-only: Visitor interface implementation (Java anonymous class inherits default)
func (v *closureRestrictToNotObjectVisitor) CaseNoType(typ jstype.JSType) any { return typ }
