// Port of com/google/javascript/jscomp/TypeMatchingStrategyTest.java @ closure-compiler v20260415

package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
)

// Java: TypeMatchingStrategyTest.EXTERNS
const typeMatchingStrategyTest_EXTERNS = `/** @constructor */
var OtherType = function() {};
/** @constructor */
var SuperType = function() {};
/** @constructor @extends {SuperType} */
var SubType = function() {};
`

func TestTypeMatchingStrategy_Match_default(t *testing.T) {
	assertMatch(t, TypeMatchingStrategy_LOOSE, "!SuperType", "!SuperType", true, false)
	assertMatch(t, TypeMatchingStrategy_LOOSE, "!SuperType", "?SuperType", true, false)
	assertMatch(t, TypeMatchingStrategy_LOOSE, "!SuperType", "SuperType|undefined", true, false)
	assertMatch(t, TypeMatchingStrategy_LOOSE, "!SuperType", "SuperType|void", true, false)
	assertMatch(t, TypeMatchingStrategy_LOOSE, "!SuperType", "!SubType", true, false)
	assertMatch(t, TypeMatchingStrategy_LOOSE, "!SuperType", "?SubType", true, false)
	assertMatch(t, TypeMatchingStrategy_LOOSE, "!SuperType", "SubType|undefined", true, false)
	assertMatch(t, TypeMatchingStrategy_LOOSE, "!SuperType", "SubType|void", true, false)
	assertMatch(t, TypeMatchingStrategy_LOOSE, "!SuperType", "number", false, false)
	assertMatch(t, TypeMatchingStrategy_LOOSE, "!SuperType", "*", true, true)
	assertMatch(t, TypeMatchingStrategy_LOOSE, "!SuperType", "?", true, true)
	assertMatch(t, TypeMatchingStrategy_LOOSE, "?", "string", true, false)
}

func TestTypeMatchingStrategy_Match_respectNullability(t *testing.T) {
	assertMatch(t, TypeMatchingStrategy_STRICT_NULLABILITY, "!SuperType", "!SuperType", true, false)
	assertMatch(t, TypeMatchingStrategy_STRICT_NULLABILITY, "!SuperType", "?SuperType", false, false)
	assertMatch(t, TypeMatchingStrategy_STRICT_NULLABILITY, "!SuperType", "SuperType|undefined", false, false)
	assertMatch(t, TypeMatchingStrategy_STRICT_NULLABILITY, "!SuperType", "SuperType|void", false, false)
	assertMatch(t, TypeMatchingStrategy_STRICT_NULLABILITY, "!SuperType", "!SubType", true, false)
	assertMatch(t, TypeMatchingStrategy_STRICT_NULLABILITY, "!SuperType", "?SubType", false, false)
	assertMatch(t, TypeMatchingStrategy_STRICT_NULLABILITY, "!SuperType", "SubType|undefined", false, false)
	assertMatch(t, TypeMatchingStrategy_STRICT_NULLABILITY, "!SuperType", "SubType|void", false, false)
	assertMatch(t, TypeMatchingStrategy_STRICT_NULLABILITY, "!SuperType", "number", false, false)
	assertMatch(t, TypeMatchingStrategy_STRICT_NULLABILITY, "!SuperType", "*", true, true)
	assertMatch(t, TypeMatchingStrategy_STRICT_NULLABILITY, "!SuperType", "?", true, true)
	assertMatch(t, TypeMatchingStrategy_STRICT_NULLABILITY, "?", "string", true, false)
}

func TestTypeMatchingStrategy_Match_subtypes(t *testing.T) {
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "!SuperType", "!SuperType", true, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "!SuperType", "?SuperType", false, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "!SuperType", "!OtherType", false, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "!SuperType", "SuperType|undefined", false, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "!SuperType", "SuperType|void", false, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "!SuperType", "!SubType", true, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "!SuperType", "?SubType", false, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "!SuperType", "SubType|undefined", false, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "!SuperType", "SubType|void", false, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "!SuperType", "number", false, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "!SuperType", "*", false, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "!SuperType", "?", false, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "?SuperType", "!SuperType", true, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "?SuperType", "?SuperType", true, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "?SuperType", "?OtherType", false, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "?SuperType", "SuperType|undefined", false, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "?SuperType", "SuperType|void", false, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "?SuperType", "SuperType|OtherType", false, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "?SuperType", "!SubType", true, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "?SuperType", "?SubType", true, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "?SuperType", "SubType|undefined", false, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "?SuperType", "SubType|void", false, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "?SuperType", "number", false, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "?SuperType", "*", false, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "?SuperType", "?", false, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "?", "string", true, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "?", "?", true, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "?", "*", true, false)
	assertMatch(t, TypeMatchingStrategy_SUBTYPES, "*", "?", false, false)
}

func TestTypeMatchingStrategy_Match_exact(t *testing.T) {
	assertMatch(t, TypeMatchingStrategy_EXACT, "!SuperType", "!SuperType", true, false)
	assertMatch(t, TypeMatchingStrategy_EXACT, "!SuperType", "?SuperType", false, false)
	assertMatch(t, TypeMatchingStrategy_EXACT, "!SuperType", "SuperType|undefined", false, false)
	assertMatch(t, TypeMatchingStrategy_EXACT, "!SuperType", "SuperType|void", false, false)
	assertMatch(t, TypeMatchingStrategy_EXACT, "!SuperType", "!SubType", false, false)
	assertMatch(t, TypeMatchingStrategy_EXACT, "!SuperType", "?SubType", false, false)
	assertMatch(t, TypeMatchingStrategy_EXACT, "!SuperType", "SubType|undefined", false, false)
	assertMatch(t, TypeMatchingStrategy_EXACT, "!SuperType", "SubType|void", false, false)
	assertMatch(t, TypeMatchingStrategy_EXACT, "!SuperType", "number", false, false)
	assertMatch(t, TypeMatchingStrategy_EXACT, "!SuperType", "*", false, false)
	assertMatch(t, TypeMatchingStrategy_EXACT, "!SuperType", "?", false, false)
	assertMatch(t, TypeMatchingStrategy_EXACT, "?", "string", true, false)
}

// typeMatchingStrategyTest_scope is the global scope the externs declare
// their constructors in.
type typeMatchingStrategyTest_scope struct{ root *rhino.Node }

func (s *typeMatchingStrategyTest_scope) GetRootNode() *rhino.Node                { return s.root }
func (s *typeMatchingStrategyTest_scope) GetParentScope() rhino.StaticScope       { return nil }
func (s *typeMatchingStrategyTest_scope) GetSlot(name string) rhino.StaticSlot    { return nil }
func (s *typeMatchingStrategyTest_scope) GetOwnSlot(name string) rhino.StaticSlot { return nil }
func (s *typeMatchingStrategyTest_scope) GetTypeOfThis() rhino.JSType             { return nil }
func (s *typeMatchingStrategyTest_scope) GetTopmostScopeOfEventualDeclaration(name string) rhino.StaticScope {
	return rhino.StaticScope_GetTopmostScopeOfEventualDeclaration(s, name)
}
func (s *typeMatchingStrategyTest_scope) LookupQualifiedName(qname rhino.QualifiedName) rhino.JSType {
	return rhino.StaticTypedScope_LookupQualifiedName(s, qname)
}

// compileTypes stands in for compiling EXTERNS and
// `/** @type {templateType} */ var x; /** @type {type} */ var y;` with
// type checking and reading the JSTypes of x and y.
//
// Go-only: test adapter — Compiler is not ported yet. A fresh registry
// declares the three extern constructors as TypedScopeCreator would
// (SubType's prototype based on SuperType), and the two @type expressions
// are evaluated in that scope (JSTypeExpression.evaluate →
// JSTypeRegistry.createTypeFromCommentNode, as for a declared var).
func compileTypes(templateType string, type_ string) (rhino.JSType, rhino.JSType) {
	registry := rhino.NewJSTypeRegistry(rhinotesting.NewTestErrorReporter())
	scope := &typeMatchingStrategyTest_scope{root: rhino.NewNode(rhino.Token_ROOT)}
	closer := registry.GetResolver().OpenForDefinition()
	ctor := func(name string) rhino.FunctionType {
		fn := rhino.IR_Function(rhino.IR_Name(""), rhino.IR_ParamList(), rhino.IR_Block())
		c := registry.CreateConstructorType(jutil.StringPtr(name), fn, registry.CreateParameters(), nil, nil, false)
		registry.DeclareType(scope, name, c.GetInstanceType())
		return c
	}
	ctor("OtherType")
	superType := ctor("SuperType")
	subType := ctor("SubType")
	subType.SetPrototypeBasedOn(superType.GetInstanceType())
	closer.Close()

	eval := func(s string) rhino.JSType {
		return registry.CreateTypeFromCommentNodeWithScope(parsing.JsDocInfoParser_ParseTypeString(s), "test", scope)
	}
	return eval(templateType), eval(type_)
}

func assertMatch(
	t *testing.T,
	typeMatchingStrategy TypeMatchingStrategy,
	templateType string,
	type_ string,
	isMatch bool,
	isLooseMatch bool) {
	t.Helper()
	truth := rhinotesting.NewTruth(t)

	// It's important that the test uses the same compiler to compile the template type and the
	// type to be matched. Otherwise, equal types won't be considered equal.
	templateJsType, jsType := compileTypes(templateType, type_)

	matchResult := typeMatchingStrategy.Match(templateJsType, jsType)
	if isMatch {
		truth.AssertWithMessage("%s", "'"+templateJsType.String()+"' should match '"+jsType.String()+"'").
			That(matchResult.IsMatch()).IsEqualTo(isMatch)
	} else {
		truth.AssertWithMessage("%s", "'"+templateType+"' should not match '"+type_+"'").
			That(matchResult.IsMatch()).IsEqualTo(isMatch)
	}
	if isLooseMatch {
		truth.AssertWithMessage("%s", "'"+templateType+"' should loosely match '"+type_+"'").
			That(matchResult.IsLooseMatch()).IsEqualTo(isLooseMatch)
	} else {
		truth.AssertWithMessage("%s", "'"+templateType+"' should not loosely match '"+type_+"'").
			That(matchResult.IsLooseMatch()).IsEqualTo(isLooseMatch)
	}
}
