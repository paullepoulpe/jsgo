// Port of com/google/javascript/rhino/jstype/FunctionTypeTest.java @ closure-compiler v20260415

package rhino_test

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	. "github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
)

type functionTypeTest struct {
	*rhinotesting.BaseJSTypeTestCase
}

func newFunctionTypeTest(t *testing.T) *functionTypeTest {
	s := &functionTypeTest{BaseJSTypeTestCase: rhinotesting.NewBaseJSTypeTestCase(t)}
	return s
}

func TestFunctionType_DefaultReturnType(t *testing.T) {
	s := newFunctionTypeTest(t)
	f := FunctionType_Builder_(s.Registry).Build()
	s.AssertThat(f.GetReturnType()).IsEqualTo(s.UNKNOWN_TYPE)
}

func TestFunctionType_SupAndInfOfReturnTypes(t *testing.T) {
	s := newFunctionTypeTest(t)
	retString :=
		FunctionType_Builder_(s.Registry).
			WithParametersOfList(s.Registry.CreateParameters()).
			WithInferredReturnType(s.STRING_TYPE).
			Build()
	retNumber :=
		FunctionType_Builder_(s.Registry).
			WithParametersOfList(s.Registry.CreateParameters()).
			WithReturnType(s.NUMBER_TYPE).
			Build()

	s.AssertLeastSupertype(
		"function(): (number|string)", retString, retNumber)
	s.AssertGreatestSubtype(
		"function(): None", retString, retNumber)

	s.AssertThat(retString.IsReturnTypeInferred()).IsTrue()
	s.AssertThat(retNumber.IsReturnTypeInferred()).IsFalse()
	s.AssertThat((retString.GetLeastSupertype(retNumber).(FunctionType)).IsReturnTypeInferred()).
		IsTrue()
	s.AssertThat((retString.GetGreatestSubtype(retString).(FunctionType)).IsReturnTypeInferred()).
		IsTrue()
}

func TestFunctionType_SupAndInfOfReturnTypesWithDifferentParams(t *testing.T) {
	s := newFunctionTypeTest(t)
	retString :=
		FunctionType_Builder_(s.Registry).
			WithParametersOfList(s.Registry.CreateParameters(s.NUMBER_TYPE)).
			WithInferredReturnType(s.STRING_TYPE).
			Build()
	retNumber :=
		FunctionType_Builder_(s.Registry).
			WithParametersOfList(s.Registry.CreateParameters()).
			WithReturnType(s.NUMBER_TYPE).
			Build()

	s.AssertLeastSupertype(
		"Function", retString, retNumber)
	s.AssertGreatestSubtype(
		"function(...*): None", retString, retNumber)
}

func TestFunctionType_SupAndInfWithDifferentParams(t *testing.T) {
	s := newFunctionTypeTest(t)
	retString :=
		FunctionType_Builder_(s.Registry).
			WithParametersOfList(s.Registry.CreateParameters(s.NUMBER_TYPE)).
			WithReturnType(s.STRING_TYPE).
			Build()
	retNumber :=
		FunctionType_Builder_(s.Registry).
			WithParametersOfList(s.Registry.CreateParameters(s.STRING_TYPE)).
			WithReturnType(s.NUMBER_TYPE).
			Build()

	s.AssertLeastSupertype(
		"Function", retString, retNumber)
	s.AssertGreatestSubtype(
		"function(...*): None", retString, retNumber)
}

func TestFunctionType_SupAndInfWithDifferentThisTypes(t *testing.T) {
	s := newFunctionTypeTest(t)
	retString :=
		FunctionType_Builder_(s.Registry).
			WithParametersOfList(s.Registry.CreateParameters()).
			WithTypeOfThis(s.OBJECT_TYPE).
			WithReturnType(s.STRING_TYPE).
			Build()
	retNumber :=
		FunctionType_Builder_(s.Registry).
			WithParametersOfList(s.Registry.CreateParameters()).
			WithTypeOfThis(s.DATE_TYPE).
			WithReturnType(s.NUMBER_TYPE).
			Build()

	s.AssertLeastSupertype(
		"function(this:Object): (number|string)", retString, retNumber)
	s.AssertGreatestSubtype(
		"function(this:Date): None", retString, retNumber)
}

func TestFunctionType_SupAndInfWithDifferentThisTypes2(t *testing.T) {
	s := newFunctionTypeTest(t)
	retString :=
		FunctionType_Builder_(s.Registry).
			WithParametersOfList(s.Registry.CreateParameters()).
			WithTypeOfThis(s.ARRAY_TYPE).
			WithReturnType(s.STRING_TYPE).
			Build()
	retNumber :=
		FunctionType_Builder_(s.Registry).
			WithParametersOfList(s.Registry.CreateParameters()).
			WithTypeOfThis(s.DATE_TYPE).
			WithReturnType(s.NUMBER_TYPE).
			Build()

	s.AssertLeastSupertype(
		"function(this:(Array|Date)): (number|string)", retString, retNumber)
	s.AssertGreatestSubtype(
		"function(this:NoObject): None", retString, retNumber)
}

func TestFunctionType_SupAndInfOfReturnTypesWithNumOfParams(t *testing.T) {
	s := newFunctionTypeTest(t)
	twoNumbers :=
		FunctionType_Builder_(s.Registry).
			WithParametersOfList(s.Registry.CreateParameters(s.NUMBER_TYPE, s.NUMBER_TYPE)).
			WithReturnType(s.BOOLEAN_TYPE).
			Build()
	oneNumber :=
		FunctionType_Builder_(s.Registry).
			WithParametersOfList(s.Registry.CreateParameters(s.NUMBER_TYPE)).
			WithReturnType(s.BOOLEAN_TYPE).
			Build()

	s.AssertLeastSupertype(
		"function(number, number): boolean", twoNumbers, oneNumber)
	s.AssertGreatestSubtype(
		"function(number): boolean", twoNumbers, oneNumber)
}

func TestFunctionType_SubtypeWithInterfaceThisType(t *testing.T) {
	s := newFunctionTypeTest(t)
	iface := s.CreateInterfaceType("I", []*TemplateType{}, false)
	ifaceReturnBoolean :=
		FunctionType_Builder_(s.Registry).
			WithParametersOfList(s.Registry.CreateParameters()).
			WithTypeOfThis(iface.GetInstanceType()).
			WithReturnType(s.BOOLEAN_TYPE).
			Build()
	objReturnBoolean :=
		FunctionType_Builder_(s.Registry).
			WithParametersOfList(s.Registry.CreateParameters()).
			WithTypeOfThis(s.OBJECT_TYPE).
			WithReturnType(s.BOOLEAN_TYPE).
			Build()
	s.AssertThat(objReturnBoolean.IsSubtype(ifaceReturnBoolean)).IsTrue()
}

func TestFunctionType_OrdinaryFunctionPrototype(t *testing.T) {
	s := newFunctionTypeTest(t)
	oneNumber :=
		FunctionType_Builder_(s.Registry).
			WithParametersOfList(s.Registry.CreateParameters(s.NUMBER_TYPE)).
			WithReturnType(s.BOOLEAN_TYPE).
			Build()
	s.AssertThat(oneNumber.GetOwnPropertyNames()).IsEmpty()
}

func TestFunctionType_CtorWithPrototypeSet(t *testing.T) {
	s := newFunctionTypeTest(t)
	closer := s.Registry.GetResolver().OpenForDefinition()
	defer closer.Close()
	{
		ctor := s.Registry.CreateConstructorType(ptr("Foo"), nil, nil, nil, nil, false)
		s.AssertThat(ctor.GetInstanceType().IsUnknownType()).IsFalse()

		node := NewNode(Token_OBJECTLIT)
		ctor.DefineDeclaredProperty("prototype", s.UNKNOWN_TYPE, node)
		s.AssertThat(ctor.GetInstanceType().IsUnknownType()).IsTrue()

		s.AssertThat(ctor.GetOwnPropertyNames()).IsEqualTo([]string{"prototype"})
		s.AssertThat(ctor.IsPropertyTypeInferred("prototype")).IsTrue()
		s.AssertThat(ctor.GetPropertyType("prototype").IsUnknownType()).IsTrue()

		s.AssertThat(ctor.GetPropertyNode("prototype")).IsEqualTo(node)
	}
}

func TestFunctionType_CtorWithInstanceInheritance(t *testing.T) {
	s := newFunctionTypeTest(t)
	closer := s.Registry.GetResolver().OpenForDefinition()
	defer closer.Close()
	{
		fooCtor :=
			FunctionType_Builder_(s.Registry).ForConstructor().WithName(ptr("Foo")).Build()
		barCtor :=
			FunctionType_Builder_(s.Registry).ForConstructor().WithName(ptr("Bar")).Build()
		barCtor.SetPrototypeBasedOn(fooCtor.GetInstanceType())
		fooCtor.GetPrototype().DefineDeclaredProperty("bar", s.STRING_TYPE, nil)

		s.AssertThat(barCtor.GetPrototype().GetImplicitPrototype()).
			IsEqualTo(fooCtor.GetInstanceType())
		s.AssertThat(fooCtor.GetInstanceType().GetSlot("bar").GetType()).IsEqualTo(s.STRING_TYPE)
	}
}

func TestFunctionType_CtorWithClassSideInheritance(t *testing.T) {
	s := newFunctionTypeTest(t)
	closer := s.Registry.GetResolver().OpenForDefinition()
	defer closer.Close()
	{
		fooCtor :=
			FunctionType_Builder_(s.Registry).ForConstructor().WithName(ptr("Foo")).Build()
		// NOTE: FunctionType does not look into the node, only at its token.
		source := NewNode(Token_CLASS)
		barCtor :=
			FunctionType_Builder_(s.Registry).
				WithSourceNode(source).
				ForConstructor().
				WithName(ptr("Bar")).
				Build()
		barCtor.SetPrototypeBasedOn(fooCtor.GetInstanceType())
		fooCtor.DefineDeclaredProperty("foo", s.NUMBER_TYPE, nil)

		s.AssertThat(barCtor.GetImplicitPrototype()).IsEqualTo(fooCtor)
		s.AssertThat(barCtor.GetSlot("foo").GetType()).IsEqualTo(s.NUMBER_TYPE)
	}
}

func TestFunctionType_EqualityOfProxyForCtor(t *testing.T) {
	s := newFunctionTypeTest(t)
	var fooCtor FunctionType
	// Bar points to `typeof Foo`.
	var barType *NamedType

	func() {
		closer := s.Registry.GetResolver().OpenForDefinition()
		defer closer.Close()
		emptyScope := rhinotesting.NewMapBasedScope(mapOf())
		key := s.Registry.CreateTemplateType("KEY")
		fooCtor =
			FunctionType_Builder_(s.Registry).
				ForConstructor().
				WithName(ptr("Foo")).
				WithTemplateKeys(key).
				Build()
		barType = s.Registry.CreateNamedType(emptyScope, "Bar", ptr(""), -1, -1)
		s.Registry.DeclareType(emptyScope, "Bar", fooCtor)
	}()

	s.AssertType(fooCtor).IsEqualTo(barType)
	s.AssertType(fooCtor).IsSubtypeOf(barType)
	s.AssertThat(barType.IsFunctionType()).IsTrue()
}

func TestFunctionType_CtorsSpecializedOnTemplateTypes(t *testing.T) {
	s := newFunctionTypeTest(t)
	closer := s.Registry.GetResolver().OpenForDefinition()
	defer closer.Close()
	{
		key := s.Registry.CreateTemplateType("KEY")
		fooCtor :=
			FunctionType_Builder_(s.Registry).
				ForConstructor().
				WithName(ptr("Foo")).
				WithTemplateKeys(key).
				WithReturnType(key).
				Build()

		stringReplacer :=
			TemplateTypeReplacer_ForInference(s.Registry, templateMapOf(key, s.STRING_TYPE))
		numberReplacer :=
			TemplateTypeReplacer_ForInference(s.Registry, templateMapOf(key, s.NUMBER_TYPE))
		fooOfString := fooCtor.Visit(stringReplacer).(JSType).ToMaybeFunctionType()
		fooOfNumber := fooCtor.Visit(numberReplacer).(JSType).ToMaybeFunctionType()

		rhinotesting.NewEqualsTester().
			AddEqualityGroup(fooCtor).
			AddEqualityGroup(fooOfString).
			AddEqualityGroup(fooOfNumber).
			TestEquals(s.T)
	}
}

func TestFunctionType_EmptyFunctionTypes(t *testing.T) {
	s := newFunctionTypeTest(t)
	s.AssertThat(s.LEAST_FUNCTION_TYPE.IsEmptyType()).IsTrue()
	s.AssertThat(s.GREATEST_FUNCTION_TYPE.IsEmptyType()).IsFalse()
}

func TestFunctionType_InterfacePrototypeChain1(t *testing.T) {
	s := newFunctionTypeTest(t)
	iface := s.CreateInterfaceType("I", []*TemplateType{}, false)
	s.AssertTypeEquals(
		iface.GetPrototype(),
		iface.GetInstanceType().GetImplicitPrototype())
	s.AssertTypeEquals(
		s.OBJECT_TYPE,
		iface.GetPrototype().GetImplicitPrototype())
}

func TestFunctionType_InterfacePrototypeChain2(t *testing.T) {
	s := newFunctionTypeTest(t)
	iface := s.CreateInterfaceType("I", []*TemplateType{}, false)
	iface.GetPrototype().DefineDeclaredProperty(
		"numberProp", s.NUMBER_TYPE, nil)

	subIface := s.CreateInterfaceType("SubI", []*TemplateType{}, false)
	subIface.SetExtendedInterfaces(
		[]ObjectType{iface.GetInstanceType()})
	s.AssertTypeEquals(
		subIface.GetPrototype(),
		subIface.GetInstanceType().GetImplicitPrototype())
	s.AssertTypeEquals(
		s.OBJECT_TYPE,
		subIface.GetPrototype().GetImplicitPrototype())

	subIfaceInst := subIface.GetInstanceType()
	s.AssertThat(subIfaceInst.HasProperty("numberProp")).IsTrue()
	s.AssertThat(subIfaceInst.IsPropertyTypeDeclared("numberProp")).IsTrue()
	s.AssertThat(subIfaceInst.IsPropertyTypeInferred("numberProp")).IsFalse()
}

func TestFunctionType_InterfacePrototypeChain3(t *testing.T) {
	s := newFunctionTypeTest(t)
	templateT := s.Registry.CreateTemplateType("T")
	iface := s.CreateInterfaceType("I", []*TemplateType{templateT}, false)
	iface.GetPrototype().DefineDeclaredProperty(
		"genericProp", templateT, nil)

	subIface := s.CreateInterfaceType("SubI", []*TemplateType{}, false)
	subIface.SetExtendedInterfaces(
		[]ObjectType{iface.GetInstanceType()})
	s.AssertTypeEquals(
		subIface.GetPrototype(),
		subIface.GetInstanceType().GetImplicitPrototype())
	s.AssertTypeEquals(
		s.OBJECT_TYPE,
		subIface.GetPrototype().GetImplicitPrototype())

	subIfaceInst := subIface.GetInstanceType()
	s.AssertThat(subIfaceInst.HasProperty("genericProp")).IsTrue()
	s.AssertThat(subIfaceInst.IsPropertyTypeDeclared("genericProp")).IsTrue()
	s.AssertThat(subIfaceInst.IsPropertyTypeInferred("genericProp")).IsFalse()
	s.AssertThat(subIfaceInst.GetPropertyType("genericProp")).IsEqualTo(templateT)
}

func (s *functionTypeTest) AssertLeastSupertype(s_ string, t1 JSType, t2 JSType) {
	s.AssertThat(t1.GetLeastSupertype(t2).String()).IsEqualTo(s_)
	s.AssertThat(t2.GetLeastSupertype(t1).String()).IsEqualTo(s_)
}

func (s *functionTypeTest) AssertGreatestSubtype(s_ string, t1 JSType, t2 JSType) {
	s.AssertThat(t1.GetGreatestSubtype(t2).String()).IsEqualTo(s_)
	s.AssertThat(t2.GetGreatestSubtype(t1).String()).IsEqualTo(s_)
}

func TestFunctionType_Equals(t *testing.T) {
	s := newFunctionTypeTest(t)
	type_ := FunctionType_Builder_(s.Registry).Build()
	s.AssertThat(type_.Equals(nil)).IsFalse()
	s.AssertThat(type_.Equals(type_)).IsTrue()
}

func TestFunctionType_EqualsParams(t *testing.T) {
	s := newFunctionTypeTest(t)
	oneNum :=
		FunctionType_Builder_(s.Registry).
			WithParametersOfList(s.Registry.CreateParameters(s.NUMBER_TYPE)).
			Build()
	optNum :=
		FunctionType_Builder_(s.Registry).
			WithParametersOfList(s.Registry.CreateOptionalParameters(s.NUMBER_TYPE)).
			Build()
	varNum :=
		FunctionType_Builder_(s.Registry).
			WithParametersOfList(s.Registry.CreateParametersWithVarArgs(s.NUMBER_TYPE)).
			Build()
	rhinotesting.Asserts_AssertEquivalenceOperations(s.T, oneNum, oneNum)
	rhinotesting.Asserts_AssertEquivalenceOperations(s.T, optNum, optNum)
	rhinotesting.Asserts_AssertEquivalenceOperations(s.T, varNum, varNum)
	s.AssertThat(oneNum.Equals(optNum)).IsFalse()
	s.AssertThat(oneNum.Equals(varNum)).IsFalse()
	s.AssertThat(optNum.Equals(varNum)).IsFalse()
}

func TestFunctionType_IsEquivalentOptAndVarArgs(t *testing.T) {
	s := newFunctionTypeTest(t)
	varNum :=
		FunctionType_Builder_(s.Registry).
			WithParametersOfList(s.Registry.CreateParametersWithVarArgs(s.NUMBER_TYPE)).
			Build()

	builder := NewFunctionParamBuilder(s.Registry)
	builder.AddOptionalParams(s.NUMBER_TYPE)
	builder.AddVarArgs(s.NUMBER_TYPE)
	optAndVarNum :=
		FunctionType_Builder_(s.Registry).WithParametersOfList(builder.Build()).Build()

	// We currently do not consider function(T=, ...T) and function(...T)
	// equivalent. This may change.
	s.AssertThat(varNum.Equals(optAndVarNum)).IsFalse()
	s.AssertThat(optAndVarNum.Equals(varNum)).IsFalse()
}

func TestFunctionType_RecursiveFunction(t *testing.T) {
	s := newFunctionTypeTest(t)
	loop := NewProxyObjectType(s.Registry, s.NUMBER_TYPE)
	fn :=
		FunctionType_Builder_(s.Registry).
			WithParametersOfList(s.Registry.CreateParameters(loop)).
			WithReturnType(loop).
			Build()

	loop.SetReferencedType(fn)
	s.AssertThat(fn.String()).IsEqualTo("function(Function): Function")

	rhinotesting.Asserts_AssertEquivalenceOperations(s.T, fn, loop)
}

func TestFunctionType_BindSignature(t *testing.T) {
	s := newFunctionTypeTest(t)
	fn :=
		FunctionType_Builder_(s.Registry).
			WithTypeOfThis(s.DATE_TYPE).
			WithParametersOfList(s.Registry.CreateParameters(s.STRING_TYPE, s.NUMBER_TYPE)).
			WithReturnType(s.BOOLEAN_TYPE).
			Build()

	s.AssertThat(fn.GetPropertyType("bind").String()).
		IsEqualTo(
			"function((Date|null|undefined), string=, number=):" + " function(...?): boolean")
}

func TestFunctionType_CallSignature1(t *testing.T) {
	s := newFunctionTypeTest(t)
	fn :=
		FunctionType_Builder_(s.Registry).
			WithTypeOfThis(s.DATE_TYPE).
			WithParametersOfList(s.Registry.CreateParameters(s.STRING_TYPE, s.NUMBER_TYPE)).
			WithReturnType(s.BOOLEAN_TYPE).
			Build()

	s.AssertThat(fn.GetPropertyType("call").String()).
		IsEqualTo("function((Date|null|undefined), string, number): boolean")
}

func TestFunctionType_CallSignature2(t *testing.T) {
	s := newFunctionTypeTest(t)
	fn :=
		FunctionType_Builder_(s.Registry).
			WithTypeOfThis(s.DATE_TYPE).
			WithParametersOfList(s.Registry.CreateParameters()).
			WithReturnType(s.BOOLEAN_TYPE).
			Build()

	s.AssertThat(fn.GetPropertyType("call").String()).IsEqualTo("function((Date|null)=): boolean")
}

func TestFunctionType_TemplatedFunctionDerivedFunctions(t *testing.T) {
	s := newFunctionTypeTest(t)
	template := s.Registry.CreateTemplateType("T")
	fn :=
		FunctionType_Builder_(s.Registry).
			WithTypeOfThis(template).
			WithTemplateKeys(template).
			WithReturnType(s.BOOLEAN_TYPE).
			Build()

	s.AssertThat(rhinotesting.Show(fn.GetPropertyType("call").GetTemplateTypeMap().GetTemplateKeys())).
		IsEqualTo("[T]")
	s.AssertThat(rhinotesting.Show(fn.GetPropertyType("apply").GetTemplateTypeMap().GetTemplateKeys())).
		IsEqualTo("[T]")
	s.AssertThat(rhinotesting.Show(fn.GetPropertyType("bind").GetTemplateTypeMap().GetTemplateKeys())).
		IsEqualTo("[T]")
	s.AssertThat(rhinotesting.Show(fn.GetBindReturnType(0).GetTemplateTypeMap().GetTemplateKeys())).
		IsEqualTo("[T]")
}

func TestFunctionType_Print_ordinaryFunction_withTemplatedThis(t *testing.T) {
	s := newFunctionTypeTest(t)
	fn :=
		FunctionType_Builder_(s.Registry).
			WithTypeOfThis(NewTemplateType(s.Registry, "T")).
			WithReturnType(s.BOOLEAN_TYPE).
			Build()
	s.AssertThat(fn.String()).IsEqualTo("function(this:T, ...?): boolean")
}

func TestFunctionType_Print_constructorFunction_withoutSource(t *testing.T) {
	s := newFunctionTypeTest(t)
	closer := s.Registry.GetResolver().OpenForDefinition()
	defer closer.Close()
	{
		fn := FunctionType_Builder_(s.Registry).ForConstructor().WithName(ptr("Foo")).Build()
		s.AssertType(fn).ToStringIsEqualTo("function(new:Foo, ...?): ?")
	}
}

func TestFunctionType_Print_interfaceFunction_withoutSource(t *testing.T) {
	s := newFunctionTypeTest(t)
	fn := FunctionType_Builder_(s.Registry).ForInterface().WithName(ptr("Foo")).Build()
	s.AssertType(fn).ToStringIsEqualTo("function(this:Foo): ?")
}

func TestFunctionType_Print_constructorFunction_withSource(t *testing.T) {
	s := newFunctionTypeTest(t)
	closer := s.Registry.GetResolver().OpenForDefinition()
	defer closer.Close()
	{
		fn :=
			FunctionType_Builder_(s.Registry).
				ForConstructor().
				WithName(ptr("Foo")).
				WithSourceNode(IR_Function(IR_Name(""), IR_ParamList(), IR_Block())).
				Build()
		s.AssertType(fn).ToStringIsEqualTo("(typeof Foo)")
	}
}

func TestFunctionType_Print_interfaceFunction_withSource(t *testing.T) {
	s := newFunctionTypeTest(t)
	closer := s.Registry.GetResolver().OpenForDefinition()
	defer closer.Close()
	{
		fn :=
			FunctionType_Builder_(s.Registry).
				ForInterface().
				WithName(ptr("Foo")).
				WithSourceNode(IR_Function(IR_Name(""), IR_ParamList(), IR_Block())).
				Build()
		s.AssertType(fn).ToStringIsEqualTo("(typeof Foo)")
	}
}

func TestFunctionType_SetImplementsOnInterface(t *testing.T) {
	s := newFunctionTypeTest(t)
	iface := s.CreateInterfaceType("I", []*TemplateType{}, false)
	subIface := s.CreateInterfaceType("SubI", []*TemplateType{}, false)
	rhinotesting.AssertThrows[jutil.Throwable](s.T, func() { subIface.SetImplementedInterfaces([]ObjectType{iface.GetInstanceType()}) })
}

func (s *functionTypeTest) CreateInterfaceType(name string, templateKeys []*TemplateType, struct_ bool) FunctionType {
	fn :=
		FunctionType_Builder_(s.Registry).
			ForInterface().
			WithName(ptr(name)).
			WithParameters().
			WithTemplateKeysOfList(templateKeys).
			Build()
	if struct_ {
		fn.SetStruct()
	}
	return fn
}
