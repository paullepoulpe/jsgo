package ast

import (
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Accessors that the grammar cannot describe on its own.

// Catch returns the catch clause of x; ok is false if x has none.
//
// (The tree keeps the catch clause in a BLOCK, the second child of the TRY,
// which is empty when there is no catch.)
func (x Try) Catch() (Catch, bool) {
	return x.catches().Catch()
}

// SetCatch sets the catch clause of x to v (which must not be in a tree).
func (x Try) SetCatch(v Catch) {
	x.catches().SetCatch(v)
}

// Catches returns the BLOCK that holds the catch clause of x, empty when x
// has none. Most code wants Catch(); this is for code that treats the holder
// as a node of its own, such as the control-flow graph (an exception goes to
// the holder, which goes to the catch clause).
func (x Try) Catches() TryCatches { return x.catches() }

// NewTry returns try body, without catch or finally yet (use SetCatch and
// SetFinally: a try needs at least one of them).
func NewTry(body Block) Try {
	n := rhino.NewNode(rhino.Token_TRY)
	n.AddChildToBack(body.node)
	n.AddChildToBack(rhino.NewNode(rhino.Token_BLOCK))
	return Try{n}
}

// NewFunction returns function name(params) body; an empty name makes an
// anonymous function.
func NewFunction(name string, params ParamList, body FunctionBody) Function {
	return Function{rhino.IR_Function(rhino.IR_Name(name), params.node, body.node)}
}

// Name returns the name of x; ok is false if x is anonymous.
//
// (An anonymous function has a NAME child with an empty string.)
func (x Function) Name() (Name, bool) {
	n := x.name()
	return n, n.Ident() != ""
}

// NameNode returns the NAME child of x, including when x is anonymous (its
// Ident is ""). Prefer Name unless traversing the raw child node itself.
func (x Function) NameNode() Name { return x.name() }

// RemoveName makes x anonymous.
//
// (The NAME child stays, with an empty string.)
func (x Function) RemoveName() { x.name().SetIdent("") }

// IsArrow reports whether x is an arrow function.
func (x Function) IsArrow() bool { return x.node.IsArrowFunction() }

// SetArrow makes x an arrow function, or a normal one.
func (x Function) SetArrow(arrow bool) { x.node.SetIsArrowFunction(arrow) }

// IsAsync reports whether x is an async function.
func (x Function) IsAsync() bool { return x.node.IsAsyncFunction() }

// SetAsync makes x an async function, or a synchronous one.
func (x Function) SetAsync(async bool) { x.node.SetIsAsyncFunction(async) }

// IsGenerator reports whether x is a generator function.
func (x Function) IsGenerator() bool { return x.node.IsGeneratorFunction() }

// SetGenerator makes x a generator function, or a regular one.
func (x Function) SetGenerator(gen bool) { x.node.SetIsGeneratorFunction(gen) }

// NameNode returns the first child of x: its NAME node, or its placeholder
// EMPTY node when x is anonymous. Prefer Name unless traversing the raw child
// node itself.
func (x Class) NameNode() rhino.Node { return x.node.GetFirstChild() }

// ExtendsNode returns the second child of x: its superclass expression, or its
// placeholder EMPTY node when x has no extends clause. Prefer Extends unless
// traversing the raw child node itself.
func (x Class) ExtendsNode() rhino.Node { return x.node.GetSecondChild() }

// RemoveName makes x anonymous by replacing its NAME child with an EMPTY
// placeholder.
func (x Class) RemoveName() {
	if c := x.node.GetFirstChild(); !c.IsEmpty() {
		c.ReplaceWith(rhino.IR_Empty().Srcref(c))
	}
}

// RemoveExtends removes the extends clause of x by replacing it with an EMPTY
// placeholder.
func (x Class) RemoveExtends() {
	if c := x.node.GetSecondChild(); !c.IsEmpty() {
		c.ReplaceWith(rhino.IR_Empty().Srcref(c))
	}
}

// IsSynthetic reports whether x is a synthetic block that does not correspond
// to a lexical block in source.
func (x Block) IsSynthetic() bool { return x.node.IsSyntheticBlock() }

// SetSynthetic marks x as a synthetic block, or clears the mark.
func (x Block) SetSynthetic(synthetic bool) { x.node.SetIsSyntheticBlock(synthetic) }

// IsAdded reports whether x was added around a single-statement branch body.
func (x Block) IsAdded() bool { return x.node.IsAddedBlock() }

// SetAdded marks x as an added block, or clears the mark.
func (x Block) SetAdded(added bool) { x.node.SetIsAddedBlock(added) }

// Value returns the value or initializer of x; ok is false if x has none.
func (x ComputedMember) Value() (v PropValue, ok bool) {
	c := x.node.GetSecondChild()
	return PropValue{c}, c != rhino.NilNode
}

// IsStatic reports whether x is a static class member.
func (x MethodDef) IsStatic() bool { return x.node.IsStaticMember() }

// SetStatic marks x as a static class member, or an instance one.
func (x MethodDef) SetStatic(static bool) { x.node.SetStaticMember(static) }

// IsStatic reports whether x is a static class method.
func (x MemberFunctionDef) IsStatic() bool { return x.node.IsStaticMember() }

// SetStatic marks x as a static class method, or an instance one.
func (x MemberFunctionDef) SetStatic(static bool) { x.node.SetStaticMember(static) }

// IsStatic reports whether x is a static class getter.
func (x GetterDef) IsStatic() bool { return x.node.IsStaticMember() }

// SetStatic marks x as a static class getter, or an instance one.
func (x GetterDef) SetStatic(static bool) { x.node.SetStaticMember(static) }

// IsStatic reports whether x is a static class setter.
func (x SetterDef) IsStatic() bool { return x.node.IsStaticMember() }

// SetStatic marks x as a static class setter, or an instance one.
func (x SetterDef) SetStatic(static bool) { x.node.SetStaticMember(static) }

// IsStatic reports whether x is a static class field.
func (x MemberFieldDef) IsStatic() bool { return x.node.IsStaticMember() }

// SetStatic marks x as a static class field, or an instance one.
func (x MemberFieldDef) SetStatic(static bool) { x.node.SetStaticMember(static) }

// IsStatic reports whether x is a static computed class field.
func (x ComputedFieldDef) IsStatic() bool { return x.node.IsStaticMember() }

// SetStatic marks x as a static computed class field, or an instance one.
func (x ComputedFieldDef) SetStatic(static bool) { x.node.SetStaticMember(static) }

// IsStatic reports whether x is a static computed class member.
func (x ComputedProp) IsStatic() bool { return x.node.IsStaticMember() }

// SetStatic marks x as a static computed class member, or an instance one.
func (x ComputedProp) SetStatic(static bool) { x.node.SetStaticMember(static) }

// IsStatic reports whether x is a static computed class member.
func (x ComputedMember) IsStatic() bool { return x.node.IsStaticMember() }

// SetStatic marks x as a static computed class member, or an instance one.
func (x ComputedMember) SetStatic(static bool) { x.node.SetStaticMember(static) }

// IsPostfix reports whether x is a postfix operator (a++ or a--) rather than
// prefix (++a or --a).
func (x IncDec) IsPostfix() bool { return x.node.GetBooleanProp(rhino.Node_INCRDECR_PROP) }

// SetPostfix marks x as postfix (true) or prefix (false).
func (x IncDec) SetPostfix(postfix bool) { x.node.PutBooleanProp(rhino.Node_INCRDECR_PROP, postfix) }

// IsFreeCall reports whether x is a free call (called without a receiver this,
// such as (0, o.f)()).
func (x Call) IsFreeCall() bool { return x.node.GetBooleanProp(rhino.Node_FREE_CALL) }

// SetFreeCall marks x as a free call (true) or bound method call (false).
func (x Call) SetFreeCall(free bool) { x.node.PutBooleanProp(rhino.Node_FREE_CALL, free) }

// IsFreeCall reports whether x is a free optional-chain call.
func (x OptChainCall) IsFreeCall() bool { return x.node.GetBooleanProp(rhino.Node_FREE_CALL) }

// SetFreeCall marks x as a free call (true) or bound method call (false).
func (x OptChainCall) SetFreeCall(free bool) { x.node.PutBooleanProp(rhino.Node_FREE_CALL, free) }

// IsFreeCall reports whether x is a free call.
func (x NormalOrOptChainCall) IsFreeCall() bool { return x.node.GetBooleanProp(rhino.Node_FREE_CALL) }

// SetFreeCall marks x as a free call (true) or bound method call (false).
func (x NormalOrOptChainCall) SetFreeCall(free bool) {
	x.node.PutBooleanProp(rhino.Node_FREE_CALL, free)
}

// IsFreeCall reports whether x is a free tagged-template call.
func (x TaggedTemplateLit) IsFreeCall() bool { return x.node.GetBooleanProp(rhino.Node_FREE_CALL) }

// SetFreeCall marks x as a free call (true) or bound method call (false).
func (x TaggedTemplateLit) SetFreeCall(free bool) { x.node.PutBooleanProp(rhino.Node_FREE_CALL, free) }

// IsYieldAll reports whether x is yield* rather than plain yield.
func (x Yield) IsYieldAll() bool { return x.node.IsYieldAll() }

// SetYieldAll marks x as yield* (true) or plain yield (false).
func (x Yield) SetYieldAll(all bool) { x.node.SetYieldAll(all) }

// IsDefault reports whether x is an export default declaration.
func (x Export) IsDefault() bool { return x.node.GetBooleanProp(rhino.Node_EXPORT_DEFAULT) }

// SetDefault marks x as export default (true) or not (false).
func (x Export) SetDefault(def bool) { x.node.PutBooleanProp(rhino.Node_EXPORT_DEFAULT, def) }

// IsAllFrom reports whether x is export * from "source".
func (x Export) IsAllFrom() bool { return x.node.GetBooleanProp(rhino.Node_EXPORT_ALL_FROM) }

// SetAllFrom marks x as export * from (true) or not (false).
func (x Export) SetAllFrom(all bool) { x.node.PutBooleanProp(rhino.Node_EXPORT_ALL_FROM, all) }

// Value returns the number.
func (x Number) Value() float64 { return x.node.GetDouble() }

// IsQuoted reports whether the key of x was written as a string literal:
// {'a': b} rather than {a: b}.
func (x StringKey) IsQuoted() bool { return x.node.IsQuotedStringKey() }

// SetQuoted marks x as having a quoted string-literal key (true) or plain
// identifier key (false).
func (x StringKey) SetQuoted(quoted bool) { x.node.PutBooleanProp(rhino.Node_QUOTED_PROP, quoted) }

// IsQuoted reports whether the key of x was written as a string literal.
func (x GetterDef) IsQuoted() bool { return x.node.IsQuotedStringKey() }

// SetQuoted marks x as having a quoted string-literal key (true) or plain
// identifier key (false).
func (x GetterDef) SetQuoted(quoted bool) { x.node.PutBooleanProp(rhino.Node_QUOTED_PROP, quoted) }

// IsQuoted reports whether the key of x was written as a string literal.
func (x SetterDef) IsQuoted() bool { return x.node.IsQuotedStringKey() }

// SetQuoted marks x as having a quoted string-literal key (true) or plain
// identifier key (false).
func (x SetterDef) SetQuoted(quoted bool) { x.node.PutBooleanProp(rhino.Node_QUOTED_PROP, quoted) }

// IsQuoted reports whether the key of x was written as a string literal.
func (x MemberFunctionDef) IsQuoted() bool { return x.node.IsQuotedStringKey() }

// SetQuoted marks x as having a quoted string-literal key (true) or plain
// identifier key (false).
func (x MemberFunctionDef) SetQuoted(quoted bool) {
	x.node.PutBooleanProp(rhino.Node_QUOTED_PROP, quoted)
}

// IsQuoted reports whether the key of x was written as a string literal.
func (x MemberFieldDef) IsQuoted() bool { return x.node.IsQuotedStringKey() }

// SetQuoted marks x as having a quoted string-literal key (true) or plain
// identifier key (false).
func (x MemberFieldDef) SetQuoted(quoted bool) { x.node.PutBooleanProp(rhino.Node_QUOTED_PROP, quoted) }

// IsQuoted reports whether the key of x was written as a string literal.
func (x MethodDef) IsQuoted() bool { return x.node.IsQuotedStringKey() }

// SetQuoted marks x as having a quoted string-literal key (true) or plain
// identifier key (false).
func (x MethodDef) SetQuoted(quoted bool) { x.node.PutBooleanProp(rhino.Node_QUOTED_PROP, quoted) }

// Value returns the property value or method function of x.
func (x PlainKeyMember) Value() PropValue { return PropValue{x.node.GetFirstChild()} }

// IsQuoted reports whether the key of x was written as a string literal.
func (x PlainKeyMember) IsQuoted() bool { return x.node.IsQuotedStringKey() }

// SetQuoted marks x as having a quoted string-literal key (true) or plain
// identifier key (false).
func (x PlainKeyMember) SetQuoted(quoted bool) { x.node.PutBooleanProp(rhino.Node_QUOTED_PROP, quoted) }

// IsQuoted reports whether the key of x was written as a string literal.
func (x NamedMember) IsQuoted() bool { return x.node.IsQuotedStringKey() }

// SetQuoted marks x as having a quoted string-literal key (true) or plain
// identifier key (false).
func (x NamedMember) SetQuoted(quoted bool) { x.node.PutBooleanProp(rhino.Node_QUOTED_PROP, quoted) }

// Value returns the initializer expression of x; ok is false if x has none.
func (x FieldDef) Value() (v Expr, ok bool) {
	if x.node.IsMemberFieldDef() {
		c := x.node.GetFirstChild()
		return Expr{c}, c != rhino.NilNode
	}
	c := x.node.GetSecondChild()
	return Expr{c}, c != rhino.NilNode
}

// IsStatic reports whether x is a static class field.
func (x FieldDef) IsStatic() bool { return x.node.IsStaticMember() }

// SetStatic marks x as a static class field, or an instance one.
func (x FieldDef) SetStatic(static bool) { x.node.SetStaticMember(static) }

// IsMethod reports whether x is a computed method ([k]() {}).
func (x ComputedProp) IsMethod() bool { return x.node.GetBooleanProp(rhino.Node_COMPUTED_PROP_METHOD) }

// SetMethod marks x as a computed method (true) or not (false).
func (x ComputedProp) SetMethod(v bool) { x.node.PutBooleanProp(rhino.Node_COMPUTED_PROP_METHOD, v) }

// IsGetter reports whether x is a computed getter (get [k]() {}).
func (x ComputedProp) IsGetter() bool { return x.node.GetBooleanProp(rhino.Node_COMPUTED_PROP_GETTER) }

// SetGetter marks x as a computed getter (true) or not (false).
func (x ComputedProp) SetGetter(v bool) { x.node.PutBooleanProp(rhino.Node_COMPUTED_PROP_GETTER, v) }

// IsSetter reports whether x is a computed setter (set [k](v) {}).
func (x ComputedProp) IsSetter() bool { return x.node.GetBooleanProp(rhino.Node_COMPUTED_PROP_SETTER) }

// SetSetter marks x as a computed setter (true) or not (false).
func (x ComputedProp) SetSetter(v bool) { x.node.PutBooleanProp(rhino.Node_COMPUTED_PROP_SETTER, v) }

// IsStart reports whether x starts an optional chain.
func (x OptChainGetProp) IsStart() bool { return x.node.IsOptionalChainStart() }

// SetStart marks x as starting an optional chain (true) or continuing one (false).
func (x OptChainGetProp) SetStart(start bool) { x.node.SetIsOptionalChainStart(start) }

// IsStart reports whether x starts an optional chain.
func (x OptChainGetElem) IsStart() bool { return x.node.IsOptionalChainStart() }

// SetStart marks x as starting an optional chain (true) or continuing one (false).
func (x OptChainGetElem) SetStart(start bool) { x.node.SetIsOptionalChainStart(start) }

// IsStart reports whether x starts an optional chain.
func (x OptChainCall) IsStart() bool { return x.node.IsOptionalChainStart() }

// SetStart marks x as starting an optional chain (true) or continuing one (false).
func (x OptChainCall) SetStart(start bool) { x.node.SetIsOptionalChainStart(start) }

// IsStart reports whether x starts an optional chain.
func (x OptChainGet) IsStart() bool { return x.node.IsOptionalChainStart() }

// SetStart marks x as starting an optional chain (true) or continuing one (false).
func (x OptChainGet) SetStart(start bool) { x.node.SetIsOptionalChainStart(start) }

// IsStart reports whether x starts an optional chain.
func (x OptChain) IsStart() bool { return x.node.IsOptionalChainStart() }

// SetStart marks x as starting an optional chain (true) or continuing one (false).
func (x OptChain) SetStart(start bool) { x.node.SetIsOptionalChainStart(start) }

// IsUnused reports whether x is marked as an unused function parameter.
func (x Param) IsUnused() bool { return x.node.IsUnusedParameter() }

// SetUnused marks x as an unused function parameter (true) or used (false).
func (x Param) SetUnused(unused bool) { x.node.SetUnusedParameter(unused) }

// IsUnusedParam reports whether x is marked as an unused function parameter.
func (x Name) IsUnusedParam() bool { return x.node.IsUnusedParameter() }

// SetUnusedParam marks x as an unused function parameter (true) or used (false).
func (x Name) SetUnusedParam(unused bool) { x.node.SetUnusedParameter(unused) }

// IsConstantName reports whether x is annotated with IS_CONSTANT_NAME.
func (x Name) IsConstantName() bool { return x.node.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME) }

// SetConstantName marks x with IS_CONSTANT_NAME (true) or clears it (false).
func (x Name) SetConstantName(constant bool) {
	x.node.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, constant)
}

// IsConstantName reports whether x is annotated with IS_CONSTANT_NAME.
func (x GetProp) IsConstantName() bool { return x.node.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME) }

// SetConstantName marks x with IS_CONSTANT_NAME (true) or clears it (false).
func (x GetProp) SetConstantName(constant bool) {
	x.node.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, constant)
}

// IsConstantName reports whether x is annotated with IS_CONSTANT_NAME.
func (x NormalOrOptChainGetProp) IsConstantName() bool {
	return x.node.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME)
}

// SetConstantName marks x with IS_CONSTANT_NAME (true) or clears it (false).
func (x NormalOrOptChainGetProp) SetConstantName(constant bool) {
	x.node.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, constant)
}

// IsEs6Constructor reports whether x is an ES6 class constructor method
// definition (non-static constructor() inside a class body).
func (x MemberFunctionDef) IsEs6Constructor() bool {
	return x.node.GetParent().IsClassMembers() && !x.IsStatic() && x.Key() == "constructor"
}

// IsEs6Constructor reports whether x is the function of an ES6 class
// constructor method definition.
func (x Function) IsEs6Constructor() bool {
	m, ok := AsMemberFunctionDef(x.node.GetParent())
	return ok && m.IsEs6Constructor()
}

// Constructor returns the non-static constructor method definition of x; ok is
// false if x has no explicit constructor.
func (x Class) Constructor() (MemberFunctionDef, bool) {
	for c := x.Body().Node().GetFirstChild(); c != rhino.NilNode; c = c.GetNext() {
		if m, ok := AsMemberFunctionDef(c); ok && !m.IsStatic() && m.Key() == "constructor" {
			return m, true
		}
	}
	return MemberFunctionDef{}, false
}

// StaticPropertyName returns the static property name accessed by x: the
// property of o.p / o?.p, or the string literal index of o['p'] / o?.['p'].
func (x NormalOrOptChainGet) StaticPropertyName() (string, bool) {
	if x.node.IsGetProp() || x.node.IsOptChainGetProp() {
		return x.node.GetString(), true
	}
	if s, ok := AsStringLit(x.node.GetSecondChild()); ok {
		return s.Value(), true
	}
	return "", false
}

// StaticPropertyName returns the static property name accessed by x: the
// property of o.p, or the string literal index of o['p'].
func (x NormalGet) StaticPropertyName() (string, bool) {
	return x.AsNormalOrOptChainGet().StaticPropertyName()
}

// ParamAt returns the i-th parameter of x (0-based); ok is false if i is out
// of range.
func (x ParamList) ParamAt(i int) (Param, bool) {
	if i < 0 {
		return Param{}, false
	}
	c := x.node.GetFirstChild()
	for ; i > 0 && c != rhino.NilNode; i-- {
		c = c.GetNext()
	}
	return Param{c}, c != rhino.NilNode
}

// ArgAt returns the i-th argument of x (0-based, after the callee); ok is
// false if i is out of range.
func (x CallOrNew) ArgAt(i int) (Arg, bool) {
	if i < 0 {
		return Arg{}, false
	}
	c := x.node.GetSecondChild()
	for ; i > 0 && c != rhino.NilNode; i-- {
		c = c.GetNext()
	}
	return Arg{c}, c != rhino.NilNode
}

// Body returns the body of the loop x.
func (x Loop) Body() Block {
	if x.node.IsDo() {
		return Block{x.node.GetFirstChild()}
	}
	return Block{x.node.GetLastChild()}
}

// NewName returns a new Name.
func NewName(ident string) Name { return Name{rhino.IR_Name(ident)} }

// NewStringLit returns a new StringLit.
func NewStringLit(value string) StringLit { return StringLit{rhino.IR_String(value)} }

// NewGetProp returns object.property.
func NewGetProp(object Expr, property string) GetProp {
	return GetProp{rhino.IR_Getprop(object.node, property)}
}

// NewLabelName returns a new LabelName.
func NewLabelName(ident string) LabelName { return LabelName{rhino.IR_LabelName(ident)} }

// Raw returns the raw text of x, as written in the source.
func (x TemplateLitString) Raw() string { return x.node.GetRawString() }

// Cooked returns the cooked text of x (escapes processed); ok is false if
// it has none (an invalid escape in a tagged template).
func (x TemplateLitString) Cooked() (s string, ok bool) {
	if p := x.node.GetCookedString(); p != nil {
		return *p, true
	}
	return "", false
}
