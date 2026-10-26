package jscomp

// Differential oracle for the reverse abstract interpreters
// (type/ChainableReverseAbstractInterpreter, SemanticReverseAbstractInterpreter,
// ClosureReverseAbstractInterpreter).
//
// This file mirrors tools/oracle/java/ReverseInterpreterDump.java step by step:
// the same type universe (every JSTypeNative, a few user classes, then the
// type expressions of testdata/reverse_interpreter_types.txt), the same seeded
// conditions, the same four interpreter chains and the same copy-on-write
// FlowScope. It must reproduce testdata/reverse_interpreter.tsv.gz (the Java
// output) byte for byte, and its "deps" section (JsIterables, Promises,
// DestructuredTarget, JSTypeColorIdHasher, TypeMatchingStrategy,
// TypeMismatch.Accumulator) testdata/jstype_deps.tsv.gz, and its
// "serialization" section (InvalidatingTypes, JSTypeReconserializer)
// testdata/jstype_serialization.tsv.gz. Regenerate the golden files with:
//
//	java -cp .cache/$CLOSURE_JAR:.cache/oracle-build \
//	  com.google.javascript.jscomp.ReverseInterpreterDump \
//	  internal/jscomp/testdata/reverse_interpreter_types.txt rai.tsv deps.tsv ser.tsv
//	gzip -9nc rai.tsv > internal/jscomp/testdata/reverse_interpreter.tsv.gz
//	gzip -9nc deps.tsv > internal/jscomp/testdata/jstype_deps.tsv.gz
//	gzip -9nc ser.tsv > internal/jscomp/testdata/jstype_serialization.tsv.gz
//
// Set RAI_ORACLE_OUT / DEPS_ORACLE_OUT / SER_ORACLE_OUT=<path> to also write
// the Go outputs there.

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type raiDump struct {
	out        *strings.Builder
	labels     []string
	types      []rhino.JSType
	ctors      []rhino.JSType
	registry   *rhino.JSTypeRegistry
	scope      *raiTestScope
	chainNames []string
	chains     []ReverseAbstractInterpreter
	seed       int64
}

// ---------------------------------------------------------------------
// Output helpers

func raiEsc(s string) string {
	var sb strings.Builder
	for _, c := range jutil.StringToUTF16(s) {
		switch {
		case c == '\\':
			sb.WriteString("\\\\")
		case c == '\n':
			sb.WriteString("\\n")
		case c == '\t':
			sb.WriteString("\\t")
		case c < 0x20 || c >= 0x7f:
			fmt.Fprintf(&sb, "\\u%04x", c)
		default:
			sb.WriteByte(byte(c))
		}
	}
	return sb.String()
}

func raiStr(t rhino.JSType) string {
	if t == nil {
		return "null"
	}
	return raiEsc(t.String())
}

func (d *raiDump) line(s string) {
	d.out.WriteString(s)
	d.out.WriteByte('\n')
}

// ---------------------------------------------------------------------
// Deterministic sequences (the same LCG as TypeAlgebraDump)

func (d *raiDump) nextInt(n int) int {
	d.seed = d.seed*6364136223846793005 + 1442695040888963407
	return int(int64(uint64(d.seed)>>33) % int64(n))
}

// ---------------------------------------------------------------------
// A reporter, a declaration scope and a copy-on-write flow scope

type raiTestReporter struct{ d *raiDump }

func (r *raiTestReporter) Warning(message string, sourceName string, line int, lineOffset int) {
	r.d.line("WARNING\t" + raiEsc(message))
}

func (r *raiTestReporter) Error(message string, sourceName string, line int, lineOffset int) {
	r.d.line("ERROR\t" + raiEsc(message))
}

type raiTestScope struct{ root rhino.Node }

func (s *raiTestScope) GetRootNode() rhino.Node              { return s.root }
func (s *raiTestScope) GetParentScope() rhino.StaticScope    { return nil }
func (s *raiTestScope) GetSlot(name string) rhino.StaticSlot { return nil }
func (s *raiTestScope) GetOwnSlot(name string) rhino.StaticSlot {
	return nil
}
func (s *raiTestScope) GetTypeOfThis() rhino.JSType { return nil }
func (s *raiTestScope) GetTopmostScopeOfEventualDeclaration(name string) rhino.StaticScope {
	return rhino.StaticScope_GetTopmostScopeOfEventualDeclaration(s, name)
}
func (s *raiTestScope) LookupQualifiedName(qname rhino.QualifiedName) rhino.JSType {
	return rhino.StaticTypedScope_LookupQualifiedName(s, qname)
}

// raiFlowScope is a flow scope whose updates copy the slot map
// (insertion-ordered).
type raiFlowScope struct {
	slots *jutil.LinkedHashMap[string, rhino.StaticTypedSlot]
}

func raiCopySlots(m *jutil.LinkedHashMap[string, rhino.StaticTypedSlot]) *jutil.LinkedHashMap[string, rhino.StaticTypedSlot] {
	c := jutil.NewLinkedHashMap[string, rhino.StaticTypedSlot]()
	for k, v := range m.All() {
		c.Put(k, v)
	}
	return c
}

func (s *raiFlowScope) WithSyntacticScope(scope rhino.StaticTypedScope) FlowScope {
	panic(jutil.NewUnsupportedOperationException())
}

func (s *raiFlowScope) InferSlotType(symbol string, type_ rhino.JSType) FlowScope {
	c := raiCopySlots(s.slots)
	c.Put(symbol, rhino.NewSimpleSlot(symbol, type_, true))
	return &raiFlowScope{slots: c}
}

func (s *raiFlowScope) InferQualifiedSlot(
	node rhino.Node, symbol string, bottomType rhino.JSType, inferredType rhino.JSType, declare bool) FlowScope {
	c := raiCopySlots(s.slots)
	c.Put(symbol, rhino.NewSimpleSlot(symbol, inferredType, true))
	return &raiFlowScope{slots: c}
}

func (s *raiFlowScope) GetDeclarationScope() rhino.StaticTypedScope { return nil }
func (s *raiFlowScope) GetRootNode() rhino.Node                     { return rhino.NilNode }
func (s *raiFlowScope) GetParentScope() rhino.StaticScope           { return nil }
func (s *raiFlowScope) GetSlot(name string) rhino.StaticSlot {
	if slot, ok := s.slots.Lookup(name); ok {
		return slot
	}
	return nil
}
func (s *raiFlowScope) GetOwnSlot(name string) rhino.StaticSlot { return s.GetSlot(name) }
func (s *raiFlowScope) GetTypeOfThis() rhino.JSType             { return nil }
func (s *raiFlowScope) GetTopmostScopeOfEventualDeclaration(name string) rhino.StaticScope {
	return rhino.StaticScope_GetTopmostScopeOfEventualDeclaration(s, name)
}
func (s *raiFlowScope) LookupQualifiedName(qname rhino.QualifiedName) rhino.JSType {
	return rhino.StaticTypedScope_LookupQualifiedName(s, qname)
}

func newRaiDump() *raiDump {
	d := &raiDump{out: &strings.Builder{}, seed: 0x5DEECE66D}
	d.scope = &raiTestScope{root: rhino.NewNode(rhino.Token_ROOT)}
	d.registry = rhino.NewJSTypeRegistry(&raiTestReporter{d: d})
	return d
}

// ---------------------------------------------------------------------
// The type universe

func (d *raiDump) nat(n rhino.JSTypeNative) rhino.JSType {
	return d.registry.GetNativeType(n)
}

func (d *raiDump) add(label string, t rhino.JSType) {
	d.labels = append(d.labels, label)
	d.types = append(d.types, t)
	d.line("type\t" + strconv.Itoa(len(d.types)) + "\t" + raiEsc(label) + "\t" + raiStr(t))
}

func raiFn(name string) rhino.Node {
	return rhino.IR_Function(rhino.IR_Name(name), rhino.IR_ParamList(), rhino.IR_Block())
}

func (d *raiDump) build(typeStrings []string) {
	closer := d.registry.GetResolver().OpenForDefinition()
	defer closer.Close()
	d.buildTypes(typeStrings)
}

func (d *raiDump) buildTypes(typeStrings []string) {
	r := d.registry
	number := d.nat(rhino.JSTypeNative_NUMBER_TYPE)
	str := d.nat(rhino.JSTypeNative_STRING_TYPE)
	boolean := d.nat(rhino.JSTypeNative_BOOLEAN_TYPE)

	fooCtor := r.CreateConstructorType(
		jutil.StringPtr("Foo"), raiFn("Foo"), r.CreateParameters(number), nil, nil, false)
	foo := fooCtor.GetInstanceType()
	fooCtor.GetPrototype().DefineDeclaredProperty("bar", number, rhino.NilNode)
	foo.DefineDeclaredProperty("p", str, rhino.NilNode)
	r.DeclareType(d.scope, "Foo", foo)

	subFooCtor := r.CreateConstructorType(
		jutil.StringPtr("SubFoo"), raiFn("SubFoo"), r.CreateParameters(), nil, nil, false)
	subFooCtor.SetPrototypeBasedOn(foo)
	subFoo := subFooCtor.GetInstanceType()
	subFoo.DefineDeclaredProperty("qux", boolean, rhino.NilNode)
	r.DeclareType(d.scope, "SubFoo", subFoo)

	iCtor := rhino.FunctionType_Builder_(r).ForInterface().WithName(jutil.StringPtr("I")).WithSourceNode(raiFn("I")).Build()
	iCtor.GetPrototype().DefineDeclaredProperty("bar", number, rhino.NilNode)
	iface := iCtor.GetInstanceType()
	r.DeclareType(d.scope, "I", iface)

	implCtor := r.CreateConstructorType(
		jutil.StringPtr("Impl"), raiFn("Impl"), r.CreateParameters(), nil, nil, false)
	implCtor.SetImplementedInterfaces([]rhino.ObjectType{iface})
	implCtor.GetPrototype().DefineDeclaredProperty("bar", number, rhino.NilNode)
	impl := implCtor.GetInstanceType()
	r.DeclareType(d.scope, "Impl", impl)

	recCtor := rhino.FunctionType_Builder_(r).ForInterface().WithName(jutil.StringPtr("Rec")).WithSourceNode(raiFn("Rec")).Build()
	recCtor.SetImplicitMatch(true)
	recCtor.GetPrototype().DefineDeclaredProperty("bar", number, rhino.NilNode)
	rec := recCtor.GetInstanceType()
	r.DeclareType(d.scope, "Rec", rec)

	tKey := r.CreateTemplateType("T")
	boxCtor := rhino.FunctionType_Builder_(r).
		ForConstructor().
		WithName(jutil.StringPtr("Box")).
		WithSourceNode(raiFn("Box")).
		WithParametersOfList(r.CreateParameters(tKey)).
		WithTemplateKeys(tKey).
		Build()
	boxCtor.GetPrototype().DefineDeclaredProperty("value", tKey, rhino.NilNode)
	box := boxCtor.GetInstanceType()
	r.DeclareType(d.scope, "Box", box)

	numEnum := rhino.EnumType_Builder_(r).SetName(jutil.StringPtr("E")).SetElementType(number).Build()
	numEnum.DefineElement("A", rhino.NilNode)
	r.DeclareType(d.scope, "E", numEnum.GetElementsType())
	strEnum := rhino.EnumType_Builder_(r).SetName(jutil.StringPtr("S")).SetElementType(str).Build()
	strEnum.DefineElement("X", rhino.NilNode)
	r.DeclareType(d.scope, "S", strEnum.GetElementsType())
	unionEnum := rhino.EnumType_Builder_(r).
		SetName(jutil.StringPtr("U")).
		SetElementType(r.CreateUnionType(str, number)).
		Build()
	unionEnum.DefineElement("Y", rhino.NilNode)
	r.DeclareType(d.scope, "U", unionEnum.GetElementsType())
	objEnum := rhino.EnumType_Builder_(r).SetName(jutil.StringPtr("O")).SetElementType(foo).Build()
	objEnum.DefineElement("Z", rhino.NilNode)
	r.DeclareType(d.scope, "O", objEnum.GetElementsType())

	for _, n := range rhino.JSTypeNative_values() {
		d.add("native:"+n.String(), d.nat(n))
	}
	d.add("Foo", foo)
	d.add("ctor:Foo", fooCtor)
	d.add("Foo.prototype", fooCtor.GetPrototype())
	d.add("SubFoo", subFoo)
	d.add("ctor:SubFoo", subFooCtor)
	d.add("I", iface)
	d.add("ctor:I", iCtor)
	d.add("Impl", impl)
	d.add("Rec", rec)
	d.add("Box", box)
	d.add("ctor:Box", boxCtor)
	d.add("T", tKey)
	d.add("enum:E", numEnum)
	d.add("E", numEnum.GetElementsType())
	d.add("enum:S", strEnum)
	d.add("S", strEnum.GetElementsType())
	d.add("U", unionEnum.GetElementsType())
	d.add("O", objEnum.GetElementsType())

	// The instanceof targets.
	d.ctors = append(d.ctors,
		fooCtor,
		subFooCtor,
		iCtor,
		implCtor,
		boxCtor,
		d.nat(rhino.JSTypeNative_ARRAY_FUNCTION_TYPE),
		d.nat(rhino.JSTypeNative_OBJECT_FUNCTION_TYPE),
		d.nat(rhino.JSTypeNative_STRING_OBJECT_FUNCTION_TYPE),
		d.nat(rhino.JSTypeNative_FUNCTION_FUNCTION_TYPE),
		d.nat(rhino.JSTypeNative_DATE_FUNCTION_TYPE),
		d.nat(rhino.JSTypeNative_REGEXP_FUNCTION_TYPE),
		d.nat(rhino.JSTypeNative_PROMISE_FUNCTION_TYPE),
		d.nat(rhino.JSTypeNative_I_THENABLE_FUNCTION_TYPE),
		d.nat(rhino.JSTypeNative_LEAST_FUNCTION_TYPE),
		d.nat(rhino.JSTypeNative_FUNCTION_TYPE),
		nil,
		number,
		d.nat(rhino.JSTypeNative_UNKNOWN_TYPE))

	templateScope := r.CreateScopeWithTemplates(d.scope, []*rhino.TemplateType{tKey})
	for _, s := range typeStrings {
		n := parsing.JsDocInfoParser_ParseTypeString(s)
		if n.IsNil() {
			d.line("unparsed\t" + raiEsc(s))
			continue
		}
		var t rhino.JSType
		var exc string
		func() {
			defer func() {
				if rec := recover(); rec != nil {
					exc = raiExceptionName(rec)
				}
			}()
			t = r.CreateTypeFromCommentNodeWithScope(n, "src.js", templateScope)
		}()
		if exc != "" {
			d.line("typeexception\t" + raiEsc(s) + "\t" + exc)
			continue
		}
		d.add("expr:"+s, t)
	}
}

// raiExceptionName is Throwable.getClass().getSimpleName() of the Java
// exception a recovered panic stands for.
func raiExceptionName(r any) string {
	name := jutil.AsThrowable(r).Error()
	if i := strings.Index(name, ":"); i >= 0 {
		name = name[:i]
	}
	if i := strings.LastIndex(name, "."); i >= 0 {
		name = name[i+1:]
	}
	return name
}

// ---------------------------------------------------------------------
// Conditions

var (
	raiTYPEOFS = []string{
		"string", "number", "boolean", "object", "function", "undefined", "symbol", "bigint", "unknown",
	}
	raiEQ_OPS  = []rhino.Token{rhino.Token_EQ, rhino.Token_NE, rhino.Token_SHEQ, rhino.Token_SHNE}
	raiCMP_OPS = []rhino.Token{rhino.Token_LT, rhino.Token_LE, rhino.Token_GT, rhino.Token_GE}
	raiPROPS   = []string{"p", "bar", "length", "then", "prototype", "qux"}
)

// raiCase is the state of one case: the blind scope's slots and their
// declaration order.
type raiCase struct {
	slots *jutil.LinkedHashMap[string, rhino.StaticTypedSlot]
	decls strings.Builder
}

func newRaiCase() *raiCase {
	return &raiCase{slots: jutil.NewLinkedHashMap[string, rhino.StaticTypedSlot]()}
}

func (c *raiCase) declare(name string, t rhino.JSType) {
	c.slots.Put(name, rhino.NewSimpleSlot(name, t, true))
	if c.decls.Len() > 0 {
		c.decls.WriteString("; ")
	}
	c.decls.WriteString(name + ":" + raiStr(t))
}

// var_ is a declared variable whose NAME node is typed like its slot.
func (c *raiCase) var_(name string, t rhino.JSType) rhino.Node {
	c.declare(name, t)
	return raiTyped(rhino.Node_NewStringWithToken(rhino.Token_NAME, name), t)
}

func raiTyped(n rhino.Node, t rhino.JSType) rhino.Node {
	if t != nil {
		n.SetJSType(t)
	}
	return n
}

func raiName(name string) rhino.Node { return rhino.Node_NewStringWithToken(rhino.Token_NAME, name) }

func (d *raiDump) num(v float64) rhino.Node {
	return raiTyped(rhino.Node_NewNumber(v), d.nat(rhino.JSTypeNative_NUMBER_TYPE))
}

func (d *raiDump) strLit(s string) rhino.Node {
	return raiTyped(rhino.Node_NewString(s), d.nat(rhino.JSTypeNative_STRING_TYPE))
}

func (d *raiDump) nul() rhino.Node {
	return raiTyped(rhino.NewNode(rhino.Token_NULL), d.nat(rhino.JSTypeNative_NULL_TYPE))
}

func (d *raiDump) undef() rhino.Node {
	return raiTyped(raiName("undefined"), d.nat(rhino.JSTypeNative_VOID_TYPE))
}

func (d *raiDump) voidZero() rhino.Node {
	return raiTyped(rhino.NewNodeWithChild(rhino.Token_VOID, d.num(0)), d.nat(rhino.JSTypeNative_VOID_TYPE))
}

func (d *raiDump) tru() rhino.Node {
	return raiTyped(rhino.NewNode(rhino.Token_TRUE), d.nat(rhino.JSTypeNative_BOOLEAN_TYPE))
}

func (d *raiDump) literal(k int) rhino.Node {
	switch k {
	case 0:
		return d.strLit("x")
	case 1:
		return d.num(0)
	case 2:
		return d.nul()
	case 3:
		return d.undef()
	case 4:
		return d.voidZero()
	default:
		return d.tru()
	}
}

func raiNode2(op rhino.Token, a, b rhino.Node) rhino.Node {
	return rhino.NewNodeWith2Children(op, a, b)
}

func raiBinop(op rhino.Token, a, b rhino.Node, swap bool) rhino.Node {
	if swap {
		return raiNode2(op, b, a)
	}
	return raiNode2(op, a, b)
}

func (d *raiDump) typeofStr(operand rhino.Node, op rhino.Token, s string, swap bool) rhino.Node {
	return raiBinop(op, rhino.NewNodeWithChild(rhino.Token_TYPEOF, operand), d.strLit(s), swap)
}

func raiCall(callee, arg rhino.Node) rhino.Node { return raiNode2(rhino.Token_CALL, callee, arg) }

func raiGetprop(obj rhino.Node, prop string) rhino.Node { return rhino.IR_Getprop(obj, prop) }

// raiGoog is goog.<fn>(arg).
func raiGoog(fn string, arg rhino.Node) rhino.Node {
	return raiCall(raiGetprop(raiName("goog"), fn), arg)
}

// raiSwitchCase is `switch (cond) { case expr: }`; returns the CASE node.
func raiSwitchCase(cond, expr rhino.Node) rhino.Node {
	c := raiNode2(rhino.Token_CASE, expr, rhino.NewNode(rhino.Token_BLOCK))
	raiNode2(rhino.Token_SWITCH, cond, rhino.NewNodeWithChild(rhino.Token_SWITCH_BODY, c))
	return c
}

// qualified is o.p with its slot and node type chosen by variant.
func (d *raiDump) qualified(c *raiCase, obj, prop string, t rhino.JSType) rhino.Node {
	variant := d.nextInt(3)
	n := raiGetprop(raiName(obj), prop)
	if variant != 1 {
		c.declare(obj+"."+prop, t)
	}
	if variant != 2 {
		raiTyped(n, t)
	}
	return n
}

func (d *raiDump) randomType() rhino.JSType {
	return d.types[d.nextInt(len(d.types))]
}

const raiUNARY = 37

func (d *raiDump) unary(c *raiCase, tmpl int, t rhino.JSType) rhino.Node {
	switch tmpl {
	case 0:
		return c.var_("a", t)
	case 1:
		return rhino.NewNodeWithChild(rhino.Token_NOT, c.var_("a", t))
	case 2:
		return rhino.NewNodeWithChild(rhino.Token_NOT, rhino.NewNodeWithChild(rhino.Token_NOT, c.var_("a", t)))
	case 3:
		op := raiEQ_OPS[d.nextInt(4)]
		s := raiTYPEOFS[d.nextInt(len(raiTYPEOFS))]
		a := c.var_("a", t)
		swap := d.nextInt(2) == 1
		return d.typeofStr(a, op, s, swap)
	case 4, 5, 6, 7, 8, 9:
		op := raiEQ_OPS[d.nextInt(4)]
		swap := d.nextInt(2) == 1
		var lit rhino.Node
		switch tmpl {
		case 4:
			lit = d.nul()
		case 5:
			lit = d.undef()
		case 6:
			lit = d.voidZero()
		case 7:
			lit = d.num(0)
		case 8:
			lit = d.strLit("")
		default:
			lit = d.tru()
		}
		return raiBinop(op, c.var_("a", t), lit, swap)
	case 10:
		op := raiCMP_OPS[d.nextInt(4)]
		a := c.var_("a", t)
		three := d.num(3)
		return raiBinop(op, a, three, d.nextInt(2) == 1)
	case 11:
		ctor := d.ctors[d.nextInt(len(d.ctors))]
		return raiNode2(rhino.Token_INSTANCEOF, c.var_("a", t), raiTyped(raiName("C"), ctor))
	case 12:
		prop := d.strLit(raiPROPS[d.nextInt(len(raiPROPS))])
		return raiNode2(rhino.Token_IN, prop, c.var_("a", t))
	case 13:
		return raiCall(raiGetprop(raiName("Array"), "isArray"), c.var_("a", t))
	case 14:
		return raiGoog("isObject", c.var_("a", t))
	case 15:
		return raiGoog("isString", c.var_("a", t))
	case 16:
		return d.qualified(c, "o", "p", t)
	case 17:
		op := raiEQ_OPS[d.nextInt(4)]
		s := raiTYPEOFS[d.nextInt(len(raiTYPEOFS))]
		q := d.qualified(c, "o", "p", t)
		swap := d.nextInt(2) == 1
		return d.typeofStr(q, op, s, swap)
	case 18:
		op := raiEQ_OPS[d.nextInt(4)]
		q := d.qualified(c, "o", "p", t)
		nul := d.nul()
		return raiBinop(op, q, nul, d.nextInt(2) == 1)
	case 19:
		return raiGoog("isObject", d.qualified(c, "o", "p", t))
	case 20:
		return rhino.NewNodeWithChild(rhino.Token_NOT, d.qualified(c, "o", "p", t))
	case 21:
		a := c.var_("a", t)
		return raiSwitchCase(a, d.literal(d.nextInt(6)))
	case 22:
		a := rhino.NewNodeWithChild(rhino.Token_TYPEOF, c.var_("a", t))
		return raiSwitchCase(a, d.strLit(raiTYPEOFS[d.nextInt(len(raiTYPEOFS))]))
	case 23:
		a := c.var_("a", t)
		return raiNode2(rhino.Token_ASSIGN, a, c.var_("b", d.randomType()))
	case 24:
		object := raiGetprop(raiName("o"), "q")
		raiTyped(object, t)
		prop := raiPROPS[d.nextInt(len(raiPROPS))]
		if d.nextInt(2) == 1 {
			c.declare("o.q."+prop, d.nat(rhino.JSTypeNative_NUMBER_TYPE))
		}
		return raiNode2(rhino.Token_IN, d.strLit(prop), object)
	case 25:
		a := c.var_("a", t)
		return raiGetprop(a, "b")
	case 26:
		return raiTyped(raiGetprop(raiCall(raiName("f"), d.num(1)), "p"), t)
	case 27:
		if d.nextInt(2) == 1 {
			return rhino.NewNode(rhino.Token_THIS)
		}
		return raiGoog("isObject", rhino.NewNode(rhino.Token_THIS))
	case 28:
		op := raiEQ_OPS[d.nextInt(4)]
		a := c.var_("a", t)
		return raiNode2(op, a, raiTyped(raiName("a"), t))
	case 29:
		op := raiEQ_OPS[d.nextInt(4)]
		s := raiTYPEOFS[d.nextInt(len(raiTYPEOFS))]
		return d.typeofStr(raiTyped(raiName("x"), t), op, s, false)
	case 30:
		return raiTyped(raiName("x"), t)
	case 31:
		c.declare("a", nil)
		a := raiTyped(raiName("a"), t)
		if d.nextInt(2) == 1 {
			return a
		}
		return rhino.NewNodeWithChild(rhino.Token_NOT, a)
	case 32:
		op := raiEQ_OPS[d.nextInt(4)]
		a := c.var_("a", t)
		return raiBinop(op, a, raiTyped(raiName("y"), d.randomType()), false)
	case 33:
		return raiCall(raiGetprop(raiName("x"), "isObject"), c.var_("a", t))
	case 34:
		call := raiGoog("isObject", c.var_("a", t))
		call.AddChildToBack(d.num(1))
		return call
	case 35:
		abc := raiGetprop(raiGetprop(raiName("a"), "b"), "c")
		if d.nextInt(2) == 1 {
			c.declare("a.b.c", t)
		}
		return raiGoog("isObject", raiTyped(abc, t))
	default:
		return raiGoog("isObject", raiTyped(raiCall(raiName("f"), d.num(2)), t))
	}
}

const raiBINARY = 20

func (d *raiDump) binary(c *raiCase, tmpl int, t1, t2 rhino.JSType) rhino.Node {
	switch tmpl {
	case 0:
		op := raiEQ_OPS[d.nextInt(4)]
		a := c.var_("a", t1)
		return raiNode2(op, a, c.var_("b", t2))
	case 1:
		op := raiCMP_OPS[d.nextInt(4)]
		a := c.var_("a", t1)
		return raiNode2(op, a, c.var_("b", t2))
	case 2:
		a := c.var_("a", t1)
		return raiNode2(rhino.Token_AND, a, c.var_("b", t2))
	case 3:
		a := c.var_("a", t1)
		return raiNode2(rhino.Token_OR, a, c.var_("b", t2))
	case 4:
		a := c.var_("a", t1)
		return rhino.NewNodeWithChild(rhino.Token_NOT, raiNode2(rhino.Token_AND, a, c.var_("b", t2)))
	case 5:
		a := c.var_("a", t1)
		return rhino.NewNodeWithChild(rhino.Token_NOT, raiNode2(rhino.Token_OR, a, c.var_("b", t2)))
	case 6:
		a := c.var_("a", t1)
		b := c.var_("b", t2)
		return raiNode2(rhino.Token_AND, a, raiNode2(rhino.Token_OR, b, c.var_("c", d.randomType())))
	case 7:
		a := c.var_("a", t1)
		b := c.var_("b", t2)
		return raiNode2(rhino.Token_AND, raiNode2(rhino.Token_OR, a, b), c.var_("c", d.randomType()))
	case 8:
		a := c.var_("a", t1)
		return raiNode2(rhino.Token_INSTANCEOF, a, c.var_("b", t2))
	case 9:
		a := c.var_("a", t1)
		return raiNode2(rhino.Token_AND, a, raiGoog("isObject", c.var_("b", t2)))
	case 10:
		a := raiGoog("isObject", c.var_("a", t1))
		return raiNode2(rhino.Token_OR, a, c.var_("b", t2))
	case 11:
		a := c.var_("a", t1)
		b := c.var_("b", t2)
		return raiNode2(
			rhino.Token_AND, raiNode2(rhino.Token_ASSIGN, a, b), raiTyped(raiName("b"), t2))
	case 12:
		av := c.var_("a", t1)
		a := d.typeofStr(av, rhino.Token_EQ, raiTYPEOFS[d.nextInt(len(raiTYPEOFS))], false)
		return raiNode2(rhino.Token_AND, a, c.var_("b", t2))
	case 13:
		p := d.qualified(c, "a", "p", t2)
		c.declare("a", t1)
		return raiNode2(rhino.Token_AND, p, raiTyped(raiName("a"), t1))
	case 14:
		op := raiEQ_OPS[d.nextInt(4)]
		a := c.var_("a", t1)
		return rhino.NewNodeWithChild(rhino.Token_NOT, raiNode2(op, a, c.var_("b", t2)))
	case 15:
		op := raiEQ_OPS[d.nextInt(4)]
		a := c.var_("a", t1)
		return raiNode2(op, a, raiTyped(rhino.Node_NewNumber(1), t2))
	case 16:
		a := c.var_("a", t1)
		return raiSwitchCase(a, c.var_("b", t2))
	case 17:
		a := c.var_("a", t1)
		return raiNode2(rhino.Token_AND, a, rhino.NewNodeWithChild(rhino.Token_NOT, c.var_("b", t2)))
	case 18:
		a := raiCall(raiGetprop(raiName("Array"), "isArray"), c.var_("a", t1))
		return raiNode2(rhino.Token_OR, a, c.var_("b", t2))
	default:
		a := c.var_("a", t1)
		b := c.var_("b", t2)
		c3 := c.var_("c", d.randomType())
		return raiNode2(
			rhino.Token_OR,
			raiNode2(rhino.Token_AND, a, b),
			raiNode2(rhino.Token_AND, c3, raiTyped(raiName("a"), t1)))
	}
}

// raiSexp serializes a condition: (TOKEN[ "string"][:type] child...).
func raiSexp(n rhino.Node, sb *strings.Builder) {
	sb.WriteByte('(')
	sb.WriteString(n.GetToken().Name())
	if n.IsName() || n.IsStringLit() || n.IsGetProp() {
		sb.WriteString(" \"" + raiEsc(n.GetString()) + "\"")
	} else if n.IsNumber() {
		sb.WriteString(" " + jutil.DoubleToString(n.GetDouble()))
	}
	if n.GetJSType() != nil {
		sb.WriteString(" :" + raiStr(n.GetJSType()))
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		sb.WriteByte(' ')
		raiSexp(c, sb)
	}
	sb.WriteByte(')')
}

// raiDiff prints how result differs from blind.
func raiDiff(blind *raiFlowScope, result FlowScope) string {
	if r, ok := result.(*raiFlowScope); ok && r == blind {
		return "="
	}
	r, ok := result.(*raiFlowScope)
	if !ok {
		return "W"
	}
	var sb strings.Builder
	sb.WriteString("~")
	for name, slot := range r.slots.All() {
		old, present := blind.slots.Lookup(name)
		t := slot.GetType()
		if !present {
			sb.WriteString(" +" + raiEsc(name) + ":" + raiStr(t))
		} else if old.GetType() != t {
			sb.WriteString(" " + raiEsc(name) + ":" + raiStr(t))
		}
	}
	return sb.String()
}

func (d *raiDump) runCase(kind string, c *raiCase, cond rhino.Node) {
	var sb strings.Builder
	// A CASE is shown with its switch.
	shown := cond
	if cond.IsCase() {
		shown = cond.GetGrandparent()
	}
	raiSexp(shown, &sb)
	d.line("case\t" + kind + "\t" + sb.String() + "\t" + c.decls.String())
	blind := &raiFlowScope{slots: c.slots}
	for i, chain := range d.chains {
		var r strings.Builder
		r.WriteString(d.chainNames[i])
		for _, o := range rhino.Outcome_values() {
			var v string
			func() {
				defer func() {
					if rec := recover(); rec != nil {
						v = "!" + raiExceptionName(rec)
					}
				}()
				v = raiDiff(blind, chain.GetPreciserScopeKnowingConditionOutcome(cond, blind, o))
			}()
			r.WriteString("\t" + v)
		}
		d.line(r.String())
	}
}

func (d *raiDump) run(typeStrings []string) {
	d.build(typeStrings)

	d.chainNames = append(d.chainNames, "S")
	d.chains = append(d.chains, NewSemanticReverseAbstractInterpreter(d.registry))
	d.chainNames = append(d.chainNames, "C")
	d.chains = append(d.chains, NewClosureReverseAbstractInterpreter(d.registry))
	d.chainNames = append(d.chainNames, "CS")
	d.chains = append(d.chains,
		NewClosureReverseAbstractInterpreter(d.registry).
			Append(NewSemanticReverseAbstractInterpreter(d.registry)).
			GetFirst())
	d.chainNames = append(d.chainNames, "SC")
	d.chains = append(d.chains,
		NewSemanticReverseAbstractInterpreter(d.registry).
			Append(NewClosureReverseAbstractInterpreter(d.registry)).
			GetFirst())

	n := len(d.types)
	// Every unary template on every type, then every typeof string on every type.
	for i := 0; i < n; i++ {
		for tmpl := 0; tmpl < raiUNARY; tmpl++ {
			c := newRaiCase()
			cond := d.unary(c, tmpl, d.types[i])
			d.runCase("u"+strconv.Itoa(tmpl)+":"+strconv.Itoa(i+1), c, cond)
		}
		for _, s := range raiTYPEOFS {
			c := newRaiCase()
			op := raiEQ_OPS[d.nextInt(4)]
			a := c.var_("a", d.types[i])
			swap := d.nextInt(2) == 1
			cond := d.typeofStr(a, op, s, swap)
			d.runCase("t:"+strconv.Itoa(i+1), c, cond)
		}
	}
	for k := 0; k < 6000; k++ {
		i := d.nextInt(n)
		j := d.nextInt(n)
		tmpl := d.nextInt(raiBINARY)
		c := newRaiCase()
		cond := d.binary(c, tmpl, d.types[i], d.types[j])
		d.runCase("b"+strconv.Itoa(tmpl)+":"+strconv.Itoa(i+1)+","+strconv.Itoa(j+1), c, cond)
	}
}

// ---------------------------------------------------------------------
// The "deps" section

// fact runs s, printing its result or the class of the exception it threw.
func (d *raiDump) fact(name string, s func() string) {
	var v string
	func() {
		defer func() {
			if rec := recover(); rec != nil {
				v = "!" + raiExceptionName(rec)
			}
		}()
		v = s()
	}()
	d.line(name + "\t" + v)
}

// raiTry is a Java try/catch (RuntimeException e) around s.
func raiTry(s func() string) (v string) {
	defer func() {
		if rec := recover(); rec != nil {
			v = "!" + raiExceptionName(rec)
		}
	}()
	return s()
}

func raiBoxed(m *JsIterables_MaybeBoxedType) string {
	templated := raiTry(func() string { return raiStr(m.GetTemplatedType()) })
	return strconv.FormatBool(m.IsMatch()) + " " + templated + " " + raiStr(m.GetMismatchType()) + " " + raiStr(m.OrElse(nil))
}

// destructure is const {p: x, ['k']: y = 1, ...r} or [a, b = 1, , ...c],
// typed by t.
func (d *raiDump) destructure(t rhino.JSType, object bool) string {
	var pattern rhino.Node
	if object {
		k := raiNode2(rhino.Token_COMPUTED_PROP, d.strLit("k"), raiNode2(rhino.Token_DEFAULT_VALUE, raiName("y"), d.num(1)))
		pattern = raiNode2(
			rhino.Token_OBJECT_PATTERN,
			rhino.IR_StringKeyWithValue("p", raiName("x")),
			k)
		pattern.AddChildToBack(rhino.IR_StringKeyWithValue("q", raiNode2(rhino.Token_DEFAULT_VALUE, raiName("z"), d.strLit(""))))
		pattern.AddChildToBack(rhino.NewNodeWithChild(rhino.Token_OBJECT_REST, raiName("r")))
	} else {
		pattern = rhino.NewNodeWith3Children(
			rhino.Token_ARRAY_PATTERN,
			raiName("a"),
			raiNode2(rhino.Token_DEFAULT_VALUE, raiName("b"), d.num(1)),
			rhino.NewNode(rhino.Token_EMPTY))
		pattern.AddChildToBack(rhino.NewNodeWithChild(rhino.Token_ITER_REST, raiName("c")))
	}
	var sb strings.Builder
	for _, target := range DestructuredTarget_CreateAllNonEmptyTargetsInPattern(d.registry, t, pattern) {
		sb.WriteString(" [" + target.GetNode().GetString() + " " +
			strconv.FormatBool(target.HasComputedProperty()) +
			strconv.FormatBool(target.HasStringKey()) +
			strconv.FormatBool(target.HasDefaultValue()))
		a := raiTry(func() string { return raiStr(target.InferTypeWithoutUsingDefaultValue()) })
		b := raiTry(func() string { return raiStr(target.InferType()) })
		sb.WriteString(" " + a + " | " + b + "]")
	}
	return sb.String()
}

func (d *raiDump) runDeps(typeStrings []string) {
	d.build(typeStrings)
	hasher := NewJSTypeColorIdHasher(d.registry)
	r := d.registry

	n := len(d.types)
	for i := 0; i < n; i++ {
		t := d.types[i]
		p := "t" + strconv.Itoa(i+1) + "\t"
		d.fact(p+"elem", func() string { return raiStr(JsIterables_GetElementType(t, r)) })
		d.fact(p+"ret", func() string { return raiStr(JsIterables_GetReturnElementType(t, r)) })
		d.fact(p+"box", func() string { return raiBoxed(JsIterables_MaybeBoxIterable(t, r)) })
		d.fact(p+"boxAsync", func() string { return raiBoxed(JsIterables_MaybeBoxIterableOrAsyncIterable(t, r)) })
		d.fact(p+"iterOf", func() string { return raiStr(JsIterables_CreateIterableTypeOf(t, r)) })
		d.fact(p+"thenable", func() string { return raiStr(Promises_GetTemplateTypeOfThenable(r, t)) })
		d.fact(p+"resolved", func() string { return raiStr(Promises_GetResolvedType(r, t)) })
		d.fact(p+"wrap", func() string { return raiStr(Promises_WrapInIThenable(r, t)) })
		d.fact(p+"asyncRet", func() string { return raiStr(Promises_CreateAsyncReturnableType(r, t)) })
		d.fact(p+"objPattern", func() string { return d.destructure(t, true) })
		d.fact(p+"arrPattern", func() string { return d.destructure(t, false) })
		d.fact(p+"colorId", func() string {
			o := t.ToMaybeObjectType()
			if o == nil {
				return "-"
			}
			return hasher.HashObjectType(o).String()
		})
	}
	location := rhino.IR_Name("loc")
	for k := 0; k < 8000; k++ {
		i := d.nextInt(n)
		j := d.nextInt(n)
		t1 := d.types[i]
		t2 := d.types[j]
		p := "p" + strconv.Itoa(i+1) + "," + strconv.Itoa(j+1) + "\t"
		d.fact(p+"match", func() string {
			var sb strings.Builder
			for _, s := range TypeMatchingStrategy_values() {
				sb.WriteString(" " + s.Match(t1, t2).String())
			}
			return sb.String()
		})
		d.fact(p+"mismatches", func() string {
			acc := NewTypeMismatch_Accumulator()
			acc.RegisterMismatch(location, t1, t2)
			var sb strings.Builder
			for _, m := range acc.GetMismatches() {
				sb.WriteString(" <" + raiStr(m.Found()) + " ; " + raiStr(m.Required()) + ">")
			}
			return sb.String()
		})
	}
}

// ---------------------------------------------------------------------
// The "serialization" section: InvalidatingTypes and JSTypeReconserializer

func (d *raiDump) runSerialization(typeStrings []string) {
	d.build(typeStrings)
	n := len(d.types)
	r := d.registry

	// Seeded mismatches, each with its own location.
	mismatches := []*TypeMismatch{}
	for k := 0; k < 300; k++ {
		found := d.types[d.nextInt(n)]
		required := d.types[d.nextInt(n)]
		mismatches = append(mismatches, NewTypeMismatch(found, required, rhino.IR_Name("loc"+strconv.Itoa(k))))
	}
	invalidating := NewInvalidatingTypes_Builder(r).AddAllTypeMismatches(mismatches).Build()
	for i := 0; i < n; i++ {
		t := d.types[i]
		d.fact("t"+strconv.Itoa(i+1)+"\tinvalidating", func() string {
			return strconv.FormatBool(invalidating.IsInvalidating(t))
		})
	}
	for key := range invalidating.GetMismatchLocations().KeySet() {
		var sb strings.Builder
		for loc := range invalidating.GetMismatchLocations().Get(key).All() {
			sb.WriteString(" " + loc.GetString())
		}
		d.line("location\t" + raiStr(key) + "\t" + sb.String())
	}

	stringPool := StringPool_Builder_()
	propagate := func(name string) bool { return !strings.HasPrefix(name, "q") }
	serializer := JSTypeReconserializer_Create(
		r,
		invalidating,
		stringPool,
		propagate,
		SerializationOptions_Builder_().
			SetIncludeDebugInfo(true).
			SetRunValidation(true).
			Build())
	for i := 0; i < n; i++ {
		t := d.types[i]
		d.fact("t"+strconv.Itoa(i+1)+"\tpointer", func() string {
			return strconv.Itoa(serializer.SerializeType(t))
		})
	}
	d.fact("pool", func() string { return raiEsc(serializer.GenerateTypePool().String()) })
	d.line("strings\t" + raiEsc(stringPool.Build().ToProto().String()))
	for key, values := range serializer.GetColorIdToJSTypeMapForDebugging().All() {
		var sb strings.Builder
		for _, t := range values {
			sb.WriteString(" " + raiStr(t))
		}
		d.line("debug\t" + key + "\t" + sb.String())
	}
}

func TestJstypeSerializationOracle(t *testing.T) {
	want := raiReadGolden(t, "jstype_serialization.tsv.gz")
	d := newRaiDump()
	d.runSerialization(raiReadTypeStrings(t))
	got := d.out.String()
	raiWriteOut(t, "SER_ORACLE_OUT", got)
	raiCompare(t, want, got, "\x00")
}

func raiReadGolden(t *testing.T, name string) string {
	f, err := os.Open("testdata/" + name)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	zr, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	b, err := io.ReadAll(zr)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// raiCompare reports the lines where got differs from want (the Java
// output); ctxPrefix selects the header line shown for context.
func raiCompare(t *testing.T, want, got, ctxPrefix string) {
	t.Helper()
	if got == want {
		return
	}
	gotLines := strings.Split(got, "\n")
	wantLines := strings.Split(want, "\n")
	diffs := 0
	for i := 0; i < len(gotLines) || i < len(wantLines); i++ {
		var g, w string
		if i < len(gotLines) {
			g = gotLines[i]
		}
		if i < len(wantLines) {
			w = wantLines[i]
		}
		if g != w {
			if diffs < 20 {
				ctx := ""
				for j := i - 1; j >= 0 && j >= i-5 && j < len(wantLines); j-- {
					if strings.HasPrefix(wantLines[j], ctxPrefix) {
						ctx = wantLines[j]
						break
					}
				}
				t.Errorf("line %d:\n  java: %s\n  go:   %s\n  in:   %s", i+1, w, g, ctx)
			}
			diffs++
		}
	}
	t.Errorf("%d differing lines (java %d lines, go %d lines)", diffs, len(wantLines), len(gotLines))
}

func raiWriteOut(t *testing.T, env, got string) {
	if p := os.Getenv(env); p != "" {
		if err := os.WriteFile(p, []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func TestJstypeDepsOracle(t *testing.T) {
	want := raiReadGolden(t, "jstype_deps.tsv.gz")
	d := newRaiDump()
	d.runDeps(raiReadTypeStrings(t))
	got := d.out.String()
	raiWriteOut(t, "DEPS_ORACLE_OUT", got)
	raiCompare(t, want, got, "\x00")
}

func raiReadTypeStrings(t *testing.T) []string {
	f, err := os.Open("testdata/reverse_interpreter_types.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var out []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		l := sc.Text()
		if l != "" && !strings.HasPrefix(l, "#") {
			out = append(out, l)
		}
	}
	return out
}

func TestReverseInterpreterOracle(t *testing.T) {
	want := raiReadGolden(t, "reverse_interpreter.tsv.gz")
	d := newRaiDump()
	d.run(raiReadTypeStrings(t))
	got := d.out.String()
	raiWriteOut(t, "RAI_ORACLE_OUT", got)
	raiCompare(t, want, got, "case\t")
}
