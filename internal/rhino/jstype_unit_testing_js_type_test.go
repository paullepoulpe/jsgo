// Port of com/google/javascript/rhino/jstype/UnitTestingJSType.java @ closure-compiler v20260415

package rhino_test

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	. "github.com/paullepoulpe/jsgo/internal/rhino"
)

// unitTestingJSType is a minimal implementation of JSType for unit tests
// and nothing else.
//
// This class has no innate behaviour. It is intended as a stand-in for
// testing behaviours on JSType that require a concrete instance. Test cases
// are responsible for any configuration.
//
// Go-only: language-adapter — the tests' anonymous subclasses override
// isNoResolvedType() and getTypeClass(); those overrides are the optional
// fields.
type unitTestingJSType struct {
	JSTypeBase
	noResolvedType *bool
	typeClass      *JSTypeClass
}

// Java: UnitTestingJSType(JSTypeRegistry)
func newUnitTestingJSType(registry *JSTypeRegistry) *unitTestingJSType {
	u := &unitTestingJSType{}
	u.InitJSTypeBase(u, registry)
	return u
}

func (u *unitTestingJSType) IsNoResolvedType() bool {
	if u.noResolvedType != nil {
		return *u.noResolvedType
	}
	return u.JSTypeBase.IsNoResolvedType()
}

// Java: UnitTestingJSType.getTypeClass()
func (u *unitTestingJSType) GetTypeClass() JSTypeClass {
	if u.typeClass != nil {
		return *u.typeClass
	}
	panic(jutil.NewUnsupportedOperationException())
}

// Java: UnitTestingJSType.recursionUnsafeHashCode()
func (u *unitTestingJSType) RecursionUnsafeHashCode() int32 {
	panic(jutil.NewUnsupportedOperationException())
}

// Java: UnitTestingJSType.getPossibleToBooleanOutcomes()
func (u *unitTestingJSType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	panic(jutil.NewUnsupportedOperationException())
}

// Java: UnitTestingJSType.visit(Visitor<T>)
func (u *unitTestingJSType) Visit(visitor Visitor) any {
	panic(jutil.NewUnsupportedOperationException())
}

// Java: UnitTestingJSType.visit(RelationshipVisitor<T>, JSType)
func (u *unitTestingJSType) VisitRelationship(visitor RelationshipVisitor, that JSType) any {
	panic(jutil.NewUnsupportedOperationException())
}

// Java: UnitTestingJSType.resolveInternal(ErrorReporter)
func (u *unitTestingJSType) ResolveInternal(reporter ErrorReporter) JSType {
	return u
}

// Java: UnitTestingJSType.appendTo(TypeStringBuilder)
func (u *unitTestingJSType) AppendTo(sb *TypeStringBuilder) {
	panic(jutil.NewUnsupportedOperationException())
}
