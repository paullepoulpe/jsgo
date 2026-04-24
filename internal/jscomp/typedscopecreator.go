// Port of com/google/javascript/jscomp/TypedScopeCreator.java @ closure-compiler v20260415

package jscomp

import (
	"fmt"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

// TODO(port): TypedScopeCreator.createScopeInternal from TypedScopeCreator.java:157
// Reason: goog.module / ES module scoping not yet ported; createScopeIn handles the common path.

// TODO(port): TypedScopeCreator.addWellKnownSymbols from TypedScopeCreator.java:298
// Reason: well-known symbols (Symbol.iterator etc.) not yet needed for corpus tests.

// TODO(port): TypedScopeCreator.FinishAndFreeze from TypedScopeCreator.java:346
// Reason: registry freeze / two-phase commit not yet wired in Go compiler pipeline.

// TODO(port): TypedScopeCreator.markGoogModuleExportsAsConst from TypedScopeCreator.java:215
// Reason: goog.module const-exports not yet tracked; no corpus failures tied to it.

type TypedScopeCreator struct {
	registry *jstype.Registry
	// Functions associates each AST FUNCTION node with its local
	// (parameter) scope. TypeCheck reads this to type-check function
	// bodies with the right lexical scope active.
	Functions map[*rhino.Node]*TypedScope
	// sourcePath and sourceSrc are set by the compiler before
	// ResolveTypesInScope is called so the JSDoc scanner can find
	// per-name positions in the source text.
	sourcePath string
	sourceSrc  string
	// diags accumulates diagnostics emitted during scope creation
	// (e.g. JSC_CYCLIC_INHERITANCE_ERROR, JSC_CONFLICTING_GETTER_SETTER_TYPE).
	// Drained by the compiler after createScopeInWithParent returns and by
	// ResolveTypesInScope.
	diags []diag.Diagnostic
}

// Go-only: language-adapter — NewTypedScopeCreator is the Go constructor for
// TypedScopeCreator; Java instantiates it via TypedScopeCreator(AbstractCompiler,
// CodingConvention) which wires into the full compiler infrastructure.
func NewTypedScopeCreator(r *jstype.Registry) *TypedScopeCreator {
	return &TypedScopeCreator{registry: r, Functions: map[*rhino.Node]*TypedScope{}}
}

// CreateScope mirrors TypedScopeCreator.createScope for the SCRIPT root.
// Walks top-level declarations and records their types; recurses into
// function bodies to build their parameter scopes.
func (c *TypedScopeCreator) CreateScope(script *rhino.Node) *TypedScope {
	global := c.createScopeIn(script, c.builtins())
	// Walk the full tree after the top-level pass so function
	// expressions / arrow functions nested inside initializers,
	// returns, and argument lists also get their local scopes. Matches
	// how TypedScopeCreator's AbstractScopeBuilder descends into every
	// syntactic child rather than the first-level only.
	c.registerNestedFunctions(script, global)
	return global
}

// Go-only: language-adapter — createScopeInWithParent is the externs-aware scope
// entry point; Java uses createScope(Node, AbstractScope) with the memoized
// parent scope obtained from a previous call, not a dedicated method.
func (c *TypedScopeCreator) createScopeInWithParent(script *rhino.Node, parent *TypedScope) *TypedScope {
	scope := c.createScopeIn(script, parent)
	c.registerNestedFunctions(script, scope)
	return scope
}

// Go-only: language-adapter — createExternScope builds the externs scope from a
// vendored source string; Java's TypedScopeCreator receives a pre-parsed extern
// AST from the compiler and does not need to parse externs separately.
func (c *TypedScopeCreator) createExternScope() *TypedScope {
	externScript, err := parseExternsSource(defaultExternsSrc)
	if err != nil || externScript == nil {
		return c.builtins()
	}
	scope := c.createScopeIn(externScript, c.builtins())
	c.registerNestedFunctions(externScript, scope)
	return scope
}

// Go-only: language-adapter — registerNestedFunctions replaces Java's
// NodeTraversal-driven inner-class traversal (AbstractScopeBuilder.visit)
// that creates function scopes on-the-fly during tree traversal.
func (c *TypedScopeCreator) registerNestedFunctions(n *rhino.Node, enclosing *TypedScope) {
	if n == nil {
		return
	}
	if n.Token() == rhino.FUNCTION {
		if _, seen := c.Functions[n]; !seen {
			c.buildExpressionFunctionScope(n, enclosing)
		}
	}
	next := enclosing
	if local, ok := c.Functions[n]; ok && local != nil {
		next = local
	}
	for _, child := range n.Children() {
		c.registerNestedFunctions(child, next)
	}
}

// Go-only: language-adapter — buildExpressionFunctionScope handles function
// expressions encountered during registerNestedFunctions; Java's
// AbstractScopeBuilder.visit creates these scopes inline during NodeTraversal.
func (c *TypedScopeCreator) buildExpressionFunctionScope(fn *rhino.Node, parent *TypedScope) {
	kids := fn.Children()
	if len(kids) < 2 {
		return
	}
	info := fn.JSDocInfo()
	// In the new pipeline, class methods have JSDoc on the MEMBER_FUNCTION_DEF
	// wrapper, not the inner FUNCTION node. Fall back to parent's JSDoc.
	if info == nil {
		if p := fn.Parent(); p != nil && (p.Token() == rhino.MEMBER_FUNCTION_DEF || p.Token() == rhino.GETTER_DEF || p.Token() == rhino.SETTER_DEF) {
			info = p.JSDocInfo()
		}
	}
	nameNode := kids[0]
	params := kids[1]
	// Named function expression with a VAR parent — promote the
	// VAR's @param / @return JSDoc so the named binding carries
	// the declared signature inside its own body. Anonymous FEEs
	// stay JSDoc-less to avoid flow-insensitive return-type FPs
	// on corpus code that relies on narrowing. Matches Java's
	// NodeUtil.getBestJSDocInfo for the `var foo = function bar()`
	// pattern specifically.
	if info == nil && nameNode != nil && nameNode.String_() != "" {
		if p := fn.Parent(); p != nil && p.Token() == rhino.NAME {
			if gp := p.Parent(); gp != nil {
				switch gp.Token() {
				case rhino.VAR, rhino.LET, rhino.CONST, rhino.ASSIGN:
					info = gp.JSDocInfo()
				}
			}
		}
	}
	local := &TypedScope{
		parent:      parent,
		declared:    make(map[string]jstype.JSType),
		typeAliases: make(map[string]jstype.JSType),
		registry:    c.registry,
	}
	// Named function expression: `var foo = function bar(){ bar() }`
	// binds `bar` inside the body so recursion + self-reference
	// works. Matches Java's AbstractScope.declareFunction. Build a
	// FunctionType carrying the declared @param / @return so bar's
	// call-sites type-check.
	var paramTypes []jstype.JSType
	for _, p := range params.Children() {
		t := resolveParamType(info, p, parent)
		paramTypes = append(paramTypes, t)
		if p.Token() == rhino.NAME {
			local.declared[p.String_()] = t
		} else {
			// Destructuring param — collect all bound names as Unknown.
			// resolveParamType returns Unknown for non-NAME params already.
			collectPatternNames(p, local)
		}
	}
	funcName := ""
	if nameNode != nil && nameNode.String_() != "" {
		funcName = nameNode.String_()
		var retType jstype.JSType = jstype.Void
		if info != nil {
			if r := info.ReturnType(); r != nil {
				retType = parent.ResolveExpr(r.TypeName)
			}
		}
		local.declared[funcName] = jstype.NewFunctionType(paramTypes, retType)
	}
	c.Functions[fn] = local
	if len(kids) >= 3 {
		body := kids[2]
		nested := c.createScopeIn(body, local)
		for k, v := range nested.declared {
			_, alreadyDefined := local.declared[k]
			// `var x` inside `var f = function x(){}` must shadow the
			// function name `x` — Java's scope creator honours the var
			// binding over the bleeding function-name binding.
			// For all other nested declarations keep the existing
			// (param / outer) binding when a name is already defined.
			if !alreadyDefined || k == funcName {
				local.declared[k] = v
			}
		}
		for k, v := range nested.typeAliases {
			if _, exists := local.typeAliases[k]; !exists {
				local.typeAliases[k] = v
			}
		}
	}
}

// Go-only: language-adapter — builtins() builds a minimal synthetic externs scope;
// Java uses DefaultExterns.java populated from type-annotated .js externs files
// through the full compiler pipeline (addWellKnownSymbols, declareNativeFunctionType).
func (c *TypedScopeCreator) builtins() *TypedScope {
	s := &TypedScope{
		declared: make(map[string]jstype.JSType),
		registry: c.registry,
	}
	// Array: minimal surface so tests that reference `@type {Array}`
	// or touch `.length` / `.push` on an array variable don't
	// collapse to Unknown. We intentionally don't register Object /
	// Function / Number / Boolean / String stubs: Java's handling of
	// their boxed vs primitive vs object forms is subtle and any
	// narrow stub we ship would false-positive more than it fixes.
	arr := jstype.NewConstructorType("Array", nil)
	arr.Instance.DefineProperty("length", jstype.Number)
	arr.Instance.DefineProperty("push", jstype.NewFunctionType([]jstype.JSType{jstype.Unknown}, jstype.Number))
	arr.Instance.DefineProperty("pop", jstype.NewFunctionType(nil, jstype.Unknown))
	arr.Instance.DefineProperty("shift", jstype.NewFunctionType(nil, jstype.Unknown))
	arr.Instance.DefineProperty("indexOf", jstype.NewFunctionType([]jstype.JSType{jstype.Unknown}, jstype.Number))
	arr.Instance.DefineProperty("slice", jstype.NewFunctionType(nil, arr.Instance))
	arr.Instance.DefineProperty("concat", jstype.NewFunctionType([]jstype.JSType{jstype.Unknown}, arr.Instance))
	s.declared["Array"] = arr
	return s
}

func (c *TypedScopeCreator) createScopeIn(block *rhino.Node, parent *TypedScope) *TypedScope {
	scope := &TypedScope{
		parent:      parent,
		declared:    make(map[string]jstype.JSType),
		typeAliases: make(map[string]jstype.JSType),
		registry:    c.registry,
	}
	hoisted := hoistDeclarations(block)
	// Pass 1: register every constructor-like binding (`@constructor`
	// functions + `class` declarations + `@record`/`@interface`
	// declarations) under its name so later `@type {Foo}` references
	// resolve nominally even when they appear textually before the
	// declaration.
	for _, child := range hoisted {
		switch child.Token() {
		case rhino.FUNCTION:
			c.preregisterConstructor(child, scope)
		case rhino.CLASS:
			c.preregisterClass(child, scope)
		}
		// Pre-register @enum bindings too so forward `@type {MyEnum}`
		// references can resolve against the EnumType even when the
		// `/** @enum */ var MyEnum = …` decl appears later in the
		// source. Matches how Java's TypedScopeCreator hoists enum
		// slots during its first scope-building pass.
		if NU.IsNameDeclaration(child) {
			c.preregisterEnum(child, scope)
			// Pre-register var-assigned constructors (`var Foo = function(){}`)
			// so within-scope @type forward references resolve correctly.
			c.preregisterVarConstructor(child, scope)
		}
	}
	// Pass 1b: record @typedef aliases. Typedefs resolve against the
	// scope that now has nominal types pre-registered, so
	// `@typedef {Foo|null}` can reference a Foo declared later.
	for _, child := range hoisted {
		if NU.IsNameDeclaration(child) {
			c.recordTypedef(child, scope)
		}
	}
	// Pass 2: record declarations + functions properly.
	for _, child := range hoisted {
		switch {
		case NU.IsNameDeclaration(child):
			c.recordDecl(child, scope)
		case child.Token() == rhino.FUNCTION:
			c.recordFunction(child, scope)
		case child.Token() == rhino.CLASS:
			c.recordClass(child, scope)
		case child.Token() == rhino.EXPR_RESULT:
			// Closure-library-style prototype method attachment:
			//   Foo.prototype.bar = function(…) {…};
			// Register `bar` on Foo's InstanceType so later property
			// access + call-site checks can see it.
			inner := child.FirstChild()
			if inner != nil && inner.Token() == rhino.ASSIGN {
				c.recordPrototypeAssign(inner, scope)
				c.recordPrototypeLiteralAssign(inner, scope)
				c.recordNamespaceAssign(inner, scope)
				c.recordPrototypeObjectLit(inner, scope)
			}
			// `Foo.prototype.bar;` stub-declaration — Java treats
			// this as declaring `bar` on Foo's instance (type
			// inferred from @type JSDoc when present, else
			// Unknown). Used so `@implements` checks pick up
			// stub-declared interface methods.
			if inner != nil && inner.Token() == rhino.GETPROP {
				c.recordPrototypeStub(child, inner, scope)
				c.recordNamespaceStub(child, inner, scope)
				c.recordNamespaceTypedefStub(child, inner, scope)
			}
		}
	}
	// Pass 3: catch parameters always have type * (Unknown) in Java's
	// TypedScopeCreator regardless of any same-named outer declaration.
	// Walk the block to find CATCH nodes and unconditionally bind their
	// parameter names to Unknown, overriding any shadowed outer binding.
	// Mirrors Java's TypedScopeCreator.AbstractScopeBuilder.visitCatch
	// which always assigns JSType.ALL to the catch parameter slot.
	registerCatchParams(block, scope)
	return scope
}

// Go-only: language-adapter — registerCatchParams binds catch params inline;
// Java's NormalScopeBuilder.visit handles CATCH nodes during NodeTraversal.
func registerCatchParams(block *rhino.Node, scope *TypedScope) {
	if block == nil {
		return
	}
	for _, child := range block.Children() {
		if child.Token() == rhino.CATCH {
			if first := child.FirstChild(); first != nil && first.Token() == rhino.NAME && first.String_() != "" {
				scope.declared[first.String_()] = jstype.Unknown
			}
			registerCatchParams(child, scope)
		} else if child.Token() != rhino.FUNCTION && child.Token() != rhino.CLASS {
			registerCatchParams(child, scope)
		}
	}
}

// Go-only: language-adapter — hoistDeclarations collects declarations eagerly;
// Java's NormalScopeBuilder visits them lazily during NodeTraversal.
func hoistDeclarations(block *rhino.Node) []*rhino.Node {
	var out []*rhino.Node
	var walk func(n *rhino.Node)
	walk = func(n *rhino.Node) {
		if n == nil {
			return
		}
		for _, child := range n.Children() {
			switch child.Token() {
			case rhino.VAR, rhino.LET, rhino.CONST, rhino.FUNCTION, rhino.CLASS, rhino.EXPR_RESULT:
				out = append(out, child)
				if child.Token() == rhino.FUNCTION || child.Token() == rhino.CLASS {
					continue
				}
				walk(child)
			case rhino.FOR, rhino.FOR_IN, rhino.FOR_OF:
				// The init / lhs child may be a VAR/LET/CONST — surface it.
				if first := child.FirstChild(); first != nil {
					t := first.Token()
					if t == rhino.VAR || t == rhino.LET || t == rhino.CONST {
						out = append(out, first)
					}
				}
				walk(child)
			case rhino.CATCH:
				// First child is the NAME of the catch parameter.
				if first := child.FirstChild(); first != nil && first.Token() == rhino.NAME && first.String_() != "" {
					synth := rhino.New(rhino.LET)
					synth.SetPos(first.Lineno(), first.Charno())
					nm := rhino.NewString(rhino.NAME, first.String_())
					nm.SetPos(first.Lineno(), first.Charno())
					synth.AddChildToBack(nm)
					out = append(out, synth)
				}
				walk(child)
			case rhino.BLOCK, rhino.IF, rhino.HOOK,
				rhino.WHILE, rhino.DO_WHILE,
				rhino.SWITCH, rhino.CASE, rhino.DEFAULT_CASE,
				rhino.TRY, rhino.FINALLY,
				rhino.LABEL:
				walk(child)
			}
		}
	}
	walk(block)
	return out
}

// Go-only: language-adapter — recordPrototypeAssign handles `Foo.prototype.x = v`
// as a standalone method; Java's NormalScopeBuilder.visit() dispatches this
// inline when it encounters an ASSIGN node during NodeTraversal.
func (c *TypedScopeCreator) recordPrototypeAssign(a *rhino.Node, scope *TypedScope) {
	kids := a.Children()
	if len(kids) < 2 {
		return
	}
	lhs, rhs := kids[0], kids[1]
	if lhs.Token() != rhino.GETPROP {
		return
	}
	lhsKids := lhs.Children()
	if len(lhsKids) < 2 {
		return
	}
	inner := lhsKids[0]
	propName := lhsKids[1].String_()
	if inner.Token() != rhino.GETPROP {
		return
	}
	innerKids := inner.Children()
	if len(innerKids) < 2 {
		return
	}
	if innerKids[1].String_() != "prototype" {
		return
	}
	// Resolve the owner: either a simple NAME (`Foo`) or a two-segment
	// namespace GETPROP (`goog.Foo`). Mirrors how Java's TypedScopeCreator
	// handles qualified prototype assignments like `goog.Foo.prototype.bar`.
	var fn *jstype.FunctionType
	switch innerKids[0].Token() {
	case rhino.NAME:
		fn, _ = scope.declared[innerKids[0].String_()].(*jstype.FunctionType)
	case rhino.GETPROP:
		gpK := innerKids[0].Children()
		if len(gpK) >= 2 && gpK[0].Token() == rhino.NAME {
			fn, _ = scope.LookupNamespaceProp(gpK[0].String_(), gpK[1].String_()).(*jstype.FunctionType)
		}
	}
	if fn == nil || fn.Instance == nil {
		return
	}
	// Use a template-aware scope when the owner has @template params so
	// `@return {T}` on a prototype method resolves to a TemplateType
	// rather than Unknown. Mirrors how Java's TypedScopeCreator uses the
	// constructor's TemplateTypeMap when resolving prototype slot types.
	resolveScope := scope
	if len(fn.TemplateTypes) > 0 {
		aliases := make(map[string]jstype.JSType, len(fn.TemplateTypes))
		for _, tt := range fn.TemplateTypes {
			aliases[tt.Name] = tt
		}
		resolveScope = &TypedScope{parent: scope, declared: map[string]jstype.JSType{}, typeAliases: aliases, registry: c.registry}
	}
	memberType := c.inferPrototypeValueType(a.JSDocInfo(), rhs, resolveScope, fn.Instance)
	fn.Instance.DefineProperty(propName, memberType)
	// Literal-inferred prototype slots (no @type) don't participate
	// in the subclass-override type-match check — Java treats them
	// as open surface in default mode.
	info := a.JSDocInfo()
	if (info == nil || !info.HasType()) && rhs.Token() != rhino.FUNCTION {
		if t := literalType(rhs); !isJSUnknown(t) {
			fn.Instance.MarkInferred(propName)
		}
	}
	// Bind `this` = Instance inside the method body so `this.x`
	// resolves to a slot lookup on Foo. Java's TypedScopeCreator
	// hooks the same thisType onto the method's FunctionType;
	// we stamp it on the expression-function's scope directly.
	if rhs.Token() == rhino.FUNCTION {
		info := a.JSDocInfo()
		var ptypes []jstype.JSType
		params := rhs.Children()
		if len(params) >= 2 {
			for _, p := range params[1].Children() {
				ptypes = append(ptypes, resolveParamType(info, p, scope))
			}
		}
		c.Functions[rhs] = c.buildMethodScope(rhs, fn.Instance, scope, ptypes)
		// Sweep the body for additional `this.<slot> = …` writes.
		// Matches how Java's TypedScopeCreator walks every prototype
		// method, not just the constructor — any `this.X = V` inside
		// a proto method also defines `X` on the Instance.
		if len(params) >= 3 {
			c.collectThisPropertiesIn(params[2], fn.Instance, scope)
		}
	}
}

// Go-only: language-adapter — recordPrototypeLiteralAssign handles `Foo.prototype={...}`
// as a standalone method; Java's TypedScopeCreator.NormalScopeBuilder.visit dispatches
// this inline when seeing an ASSIGN with an object-literal RHS.
func (c *TypedScopeCreator) recordPrototypeLiteralAssign(a *rhino.Node, scope *TypedScope) {
	kids := a.Children()
	if len(kids) < 2 {
		return
	}
	lhs, rhs := kids[0], kids[1]
	if rhs.Token() != rhino.OBJECTLIT {
		return
	}
	if lhs.Token() != rhino.GETPROP {
		return
	}
	lhsKids := lhs.Children()
	if len(lhsKids) < 2 || lhsKids[1].String_() != "prototype" {
		return
	}
	// Resolve the constructor FunctionType from the receiver.
	// Accepts `Foo.prototype` (NAME receiver) and `ns.Foo.prototype`
	// (GETPROP receiver with NAME base).
	recv := lhsKids[0]
	var fn *jstype.FunctionType
	switch recv.Token() {
	case rhino.NAME:
		fn, _ = scope.GetDeclaredType(recv.String_()).(*jstype.FunctionType)
	case rhino.GETPROP:
		recvKids := recv.Children()
		if len(recvKids) >= 2 && recvKids[0].Token() == rhino.NAME {
			ns := recvKids[0].String_()
			prop := recvKids[1].String_()
			fn, _ = scope.LookupNamespaceProp(ns, prop).(*jstype.FunctionType)
		}
	}
	if fn == nil || fn.Instance == nil {
		return
	}
	// Register each key of the object literal on the constructor's
	// InstanceType. Function-expression values get a proper FunctionType;
	// other values fall back to their literal type.
	for _, kv := range rhs.Children() {
		if kv.Token() != rhino.STRING_KEY {
			continue
		}
		keyName := kv.String_()
		val := kv.FirstChild()
		if val == nil {
			continue
		}
		memberType := c.inferPrototypeValueType(a.JSDocInfo(), val, scope, fn.Instance)
		fn.Instance.DefineProperty(keyName, memberType)
		// Bind `this` inside method bodies so `this.x` type-checks.
		if val.Token() == rhino.FUNCTION {
			c.Functions[val] = c.buildMethodScope(val, fn.Instance, scope, nil)
		}
	}
}

// recordNamespaceAssign handles `/** @type {T} */ owner.prop = …` declarations.
// Go-only: language-adapter — Java's NormalScopeBuilder dispatches this inline.
func (c *TypedScopeCreator) recordNamespaceAssign(a *rhino.Node, scope *TypedScope) {
	info := a.JSDocInfo()
	// No @type JSDoc — but if the RHS is a FUNCTION, infer a
	// FunctionType from the JSDoc's @param / @return so arity
	// checks work at call sites. Matches how Java's
	// TypedScopeCreator handles `ns.foo = function(x) {…}`.
	if info == nil || (!info.HasType() && !info.IsEnum()) {
		c.recordNamespaceAssignFromFunction(a, scope)
		return
	}
	// `/** @enum {T} */ ns.MyEnum = {…}` — register the enum as a
	// namespace prop so `@param {ns.MyEnum}` type references resolve.
	if info.IsEnum() {
		kids2 := a.Children()
		if len(kids2) >= 2 {
			lhs2 := kids2[0]
			if lhs2.Token() == rhino.GETPROP {
				lhsKids2 := lhs2.Children()
				if len(lhsKids2) >= 2 && lhsKids2[0].Token() == rhino.NAME {
					ownerName2 := lhsKids2[0].String_()
					propName2 := lhsKids2[1].String_()
					if _, ok := scope.GetDeclaredType(ownerName2).(*jstype.FunctionType); !ok {
						var underlying jstype.JSType = jstype.Number
						if expr := info.EnumParameterType(); expr != nil {
							underlying = scope.ResolveExpr(expr.TypeName)
						}
						enum := jstype.NewEnumType(ownerName2+"."+propName2, underlying)
						if kids2[1].Token() == rhino.OBJECTLIT {
							for _, kv := range kids2[1].Children() {
								if kv.Token() == rhino.STRING_KEY {
									enum.Members[kv.String_()] = true
								}
							}
						}
						scope.recordNamespaceProp(ownerName2, propName2, enum, lhsKids2[1])
					}
				}
			}
		}
		return
	}
	kids := a.Children()
	if len(kids) < 2 {
		return
	}
	lhs := kids[0]
	if lhs.Token() != rhino.GETPROP {
		return
	}
	lhsKids := lhs.Children()
	if len(lhsKids) < 2 {
		return
	}
	recv := lhsKids[0]
	if recv.Token() != rhino.NAME {
		return
	}
	// Skip the `X.prototype.Y` form — recordPrototypeAssign owns it.
	if recv.Token() == rhino.GETPROP {
		return
	}
	propName := lhsKids[1].String_()
	ownerName := recv.String_()
	// Skip when the owner is a known constructor (FunctionType) —
	// that's a static-property assignment and Java routes those
	// through the constructor's Instance.defineStaticProperty.
	if _, ok := scope.GetDeclaredType(ownerName).(*jstype.FunctionType); ok {
		return
	}
	declared := scope.ResolveExpr(info.Type().TypeName)
	scope.recordNamespaceProp(ownerName, propName, declared, lhsKids[1])
}

// Go-only: language-adapter — recordNamespaceAssignFromFunction infers
// FunctionType from `ns.foo = function(){}`; Java's NormalScopeBuilder.defineName
// handles this inline via createFunctionTypeFromNodes.
func (c *TypedScopeCreator) recordNamespaceAssignFromFunction(a *rhino.Node, scope *TypedScope) {
	kids := a.Children()
	if len(kids) < 2 {
		return
	}
	lhs, rhs := kids[0], kids[1]
	// Namespace-to-namespace ctor alias: \`ns.A = ns.B\` or \`ns.A =
	// existingName\`. Share the resolved constructor so both refer
	// to the same Instance.
	if lhs.Token() == rhino.GETPROP && (rhs.Token() == rhino.NAME || rhs.Token() == rhino.GETPROP) {
		lhsKids := lhs.Children()
		if len(lhsKids) >= 2 && lhsKids[0].Token() == rhino.NAME {
			src := resolveQualifiedType(scope, rhs)
			if fn, ok := src.(*jstype.FunctionType); ok && fn.Instance != nil {
				scope.recordNamespaceProp(lhsKids[0].String_(), lhsKids[1].String_(), fn, lhsKids[1])
				return
			}
			// Typedef alias: `ns.alias = original` where original is a typedef NAME.
			// Propagate the typedef into the namespace slot so `@type {ns.alias}`
			// resolves to the same underlying type.
			// Mirrors Java's TypedScopeCreator const-aliased typedef path.
			if rhs.Token() == rhino.NAME {
				if aliasedType := scope.LookupTypeAlias(rhs.String_()); aliasedType != nil {
					qname := lhsKids[0].String_() + "." + lhsKids[1].String_()
					scope.typeAliases[qname] = aliasedType
					return
				}
			}
		}
	}
	// Scalar literal assignments: `ns.x = 3` / `ns.x = "str"` / etc.
	// Record the primitive type so `new ns.x` / `@type {ns.x}` uses
	// can detect the non-constructor value. Mirrors how Java's scope
	// creator assigns the literal type to simple namespace properties.
	// Skip if the receiver is already a known class instance (ObjectType)
	// or constructor (FunctionType) — those assignments go through the
	// instance/prototype property mechanism, not namespace tracking.
	// These slots are marked inferred so the assignment checker treats
	// them as open surfaces (matching Java's default-mode behaviour for
	// non-@type namespace properties).
	if lhs.Token() == rhino.GETPROP {
		if t := literalType(rhs); !isJSUnknown(t) {
			lhsKids := lhs.Children()
			if len(lhsKids) >= 2 && lhsKids[0].Token() == rhino.NAME {
				recv := lhsKids[0].String_()
				recvType := scope.GetDeclaredType(recv)
				if _, isFn := recvType.(*jstype.FunctionType); !isFn {
					if _, isObj := recvType.(*jstype.ObjectType); !isObj {
						propName := lhsKids[1].String_()
						// Only record literal-inferred type when the prop
						// has not already been declared with @type JSDoc —
						// a re-assignment to a declared slot must not
						// overwrite the declared type or mark it inferred.
						if scope.LookupNamespaceProp(recv, propName) == nil {
							scope.recordNamespaceProp(recv, propName, t, lhsKids[1])
							scope.markNamespacePropInferred(recv, propName)
						}
						return
					}
				}
			}
		}
	}
	// Object-literal assignments: `ns.foo = {}` / `ns.foo = {…}`.
	// Record Unknown as a sentinel namespace prop so nameResolves can
	// detect the prop exists but is not a constructor-type. Mirrors how
	// Java's TypedScopeCreator marks namespace object-literal slots as
	// plain values distinct from @constructor declarations.
	if rhs.Token() == rhino.OBJECTLIT && lhs.Token() == rhino.GETPROP {
		lhsKids := lhs.Children()
		if len(lhsKids) >= 2 && lhsKids[0].Token() == rhino.NAME {
			recv := lhsKids[0].String_()
			if _, ok := scope.GetDeclaredType(recv).(*jstype.FunctionType); !ok {
				scope.recordNamespaceProp(recv, lhsKids[1].String_(), jstype.Unknown, lhsKids[1])
			}
		}
		return
	}
	if rhs.Token() != rhino.FUNCTION || lhs.Token() != rhino.GETPROP {
		return
	}
	lhsKids := lhs.Children()
	if len(lhsKids) < 2 {
		return
	}
	recv := lhsKids[0]
	if recv.Token() != rhino.NAME {
		return
	}
	if _, ok := scope.GetDeclaredType(recv.String_()).(*jstype.FunctionType); ok {
		return
	}
	info := a.JSDocInfo()
	params := rhs.Children()
	if len(params) < 2 {
		return
	}
	// @constructor on a \`ns.Foo = function(){…}\` assignment builds
	// a ConstructorType whose Instance is reachable under the
	// qualified name. scope.ResolveExpr consults namespaceProps for
	// dotted references so \`@type {!ns.Foo}\` resolves nominally.
	if info != nil && (info.IsConstructor() || info.IsInterface() || info.IsRecord()) {
		qname := recv.String_() + "." + lhsKids[1].String_()
		var paramTypes []jstype.JSType
		for _, p := range params[1].Children() {
			paramTypes = append(paramTypes, resolveParamType(info, p, scope))
		}
		ctor := jstype.NewConstructorType(qname, paramTypes)
		if info.IsInterface() || info.IsRecord() {
			ctor.Instance.IsStructural = true
			ctor.Instance.IsRecord = info.IsRecord()
		}
		if info.IsInterface() {
			ctor.Instance.IsInterface = true
		}
		if info.IsStruct() {
			ctor.Instance.IsStruct = true
		}
		if info.IsDict() {
			ctor.Instance.IsDict = true
		}
		if base := info.BaseType(); base != nil {
			if super := scope.LookupNominalType(stripTypeArgs(base.TypeName)); super != nil {
				ctor.Instance.SetSuper(super)
			}
		}
		setArityFromParams(ctor, params[1], info)
		scope.recordNamespaceProp(recv.String_(), lhsKids[1].String_(), ctor, lhsKids[1])
		// Build a method scope for the constructor body so `this`
		// resolves to ctor.Instance inside the body. Mirrors how
		// Java's TypedScopeCreator wires thisType for named
		// constructors; without this, `this.method()` calls inside
		// the body cannot resolve the instance type.
		ctorScope := c.buildMethodScope(rhs, ctor.Instance, scope, paramTypes)
		c.Functions[rhs] = ctorScope
		// Also collect `this.X = …` writes inside the ctor body.
		if len(params) >= 3 {
			c.collectThisPropertiesIn(params[2], ctor.Instance, scope)
		}
		return
	}
	var paramTypes []jstype.JSType
	for _, p := range params[1].Children() {
		paramTypes = append(paramTypes, resolveParamType(info, p, scope))
	}
	var retType jstype.JSType = jstype.Void
	if info != nil {
		if r := info.ReturnType(); r != nil {
			retType = scope.ResolveExpr(r.TypeName)
		}
	}
	fnType := jstype.NewFunctionType(paramTypes, retType)
	setArityFromParams(fnType, params[1], info)
	scope.recordNamespaceProp(recv.String_(), lhsKids[1].String_(), fnType, lhsKids[1])
}


// recordNamespaceStub handles the `owner.prop;` stub form — an
// EXPR_RESULT carrying a bare two-segment GETPROP with a @type
// JSDoc. Matches the assignment form's behaviour, used for Closure
// library-style type declarations that don't assign a value.
// Go-only: language-adapter — recordNamespaceStub is a Go-only helper that handles `owner.prop;` stubs as standalone method; Java inlines this in NormalScopeBuilder.visit for EXPR_RESULT nodes.
func (c *TypedScopeCreator) recordNamespaceStub(stmt *rhino.Node, getprop *rhino.Node, scope *TypedScope) {
	info := stmt.JSDocInfo()
	if info == nil {
		info = getprop.JSDocInfo()
	}
	if info == nil || !info.HasType() {
		return
	}
	kids := getprop.Children()
	if len(kids) < 2 {
		return
	}
	recv := kids[0]
	if recv.Token() != rhino.NAME {
		return
	}
	ownerName := recv.String_()
	if _, ok := scope.GetDeclaredType(ownerName).(*jstype.FunctionType); ok {
		return
	}
	propName := kids[1].String_()
	declared := scope.ResolveExpr(info.Type().TypeName)
	scope.recordNamespaceProp(ownerName, propName, declared, kids[1])
}

// recordNamespaceTypedefStub handles `/** @typedef {T} */ ns.Name;`
// stub-declarations: a namespace property whose JSDoc carries @typedef
// rather than @type. Registers the resolved type under the namespace's
// typeAliases map so `@type {ns.Name}` references resolve correctly.
// Go-only: language-adapter — recordNamespaceTypedefStub is a Go-only helper that handles `@typedef` stubs on GETPROP as standalone method; Java inlines this in NormalScopeBuilder.identifyEnumOrTypedefDeclaration.
func (c *TypedScopeCreator) recordNamespaceTypedefStub(stmt *rhino.Node, getprop *rhino.Node, scope *TypedScope) {
	info := stmt.JSDocInfo()
	if info == nil {
		info = getprop.JSDocInfo()
	}
	if info == nil || info.TypedefType() == nil {
		return
	}
	kids := getprop.Children()
	if len(kids) < 2 {
		return
	}
	recv := kids[0]
	if recv.Token() != rhino.NAME {
		return
	}
	ownerName := recv.String_()
	propName := kids[1].String_()
	// Store as a dotted typeAlias so LookupTypeAlias("owner.prop") finds it.
	qname := ownerName + "." + propName
	resolved := scope.ResolveExpr(info.TypedefType().TypeName)
	scope.typeAliases[qname] = resolved
}

// recordPrototypeStub handles the `X.prototype.name;` stub-
// declaration form: an EXPR_RESULT whose child is a GETPROP
// walking `X.prototype.name`. Registers `name` on X's InstanceType
// at the @type-annotated slot (or Unknown when no JSDoc), matching
// how Java's TypedScopeCreator lifts these stubs into the scope so
// `@implements` + property-existence checks see them.
// Go-only: language-adapter — recordPrototypeStub is a Go-only helper that handles `X.prototype.name;` stubs as standalone method; Java inlines this in NormalScopeBuilder.visit for EXPR_RESULT.
func (c *TypedScopeCreator) recordPrototypeStub(stmt *rhino.Node, getprop *rhino.Node, scope *TypedScope) {
	kids := getprop.Children()
	if len(kids) < 2 {
		return
	}
	inner := kids[0]
	propName := kids[1].String_()
	if inner.Token() != rhino.GETPROP {
		return
	}
	innerKids := inner.Children()
	if len(innerKids) < 2 {
		return
	}
	if innerKids[0].Token() != rhino.NAME || innerKids[1].String_() != "prototype" {
		return
	}
	ownerName := innerKids[0].String_()
	fn, ok := scope.declared[ownerName].(*jstype.FunctionType)
	if !ok || fn.Instance == nil {
		return
	}
	if fn.Instance.GetPropertyType(propName) != nil {
		return // already declared by a prior assign
	}
	var t jstype.JSType = jstype.Unknown
	info := stmt.JSDocInfo()
	if info == nil {
		info = getprop.JSDocInfo()
	}
	if info != nil {
		if info.HasType() {
			t = scope.ResolveExpr(info.Type().TypeName)
		} else if info.ReturnType() != nil || len(info.ParameterNames()) > 0 {
			// `/** @return {T} */ X.prototype.m;` — stubbed method.
			var paramTypes []jstype.JSType
			for _, n := range info.ParameterNames() {
				if expr := info.ParameterType(n); expr != nil {
					paramTypes = append(paramTypes, scope.ResolveExpr(expr.TypeName))
				}
			}
			var ret jstype.JSType = jstype.Unknown
			if info.ReturnType() != nil {
				ret = scope.ResolveExpr(info.ReturnType().TypeName)
			}
			t = jstype.NewFunctionTypeWithThis(paramTypes, ret, fn.Instance)
		}
	}
	fn.Instance.DefineProperty(propName, t)
}

// recordPrototypeObjectLit handles `X.prototype = {key: fn, …}` —
// an ASSIGN whose LHS is `X.prototype` and RHS is an OBJECTLIT.
// Registers each STRING_KEY member on X's InstanceType so that
// JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH checks can compare the
// override type against the inherited superclass/interface slot.
// Go-only: language-adapter — recordPrototypeObjectLit is a Go-only helper that handles `X.prototype={key:fn}` as standalone method; Java inlines this in NormalScopeBuilder.defineObjectLiteral.
func (c *TypedScopeCreator) recordPrototypeObjectLit(a *rhino.Node, scope *TypedScope) {
	kids := a.Children()
	if len(kids) < 2 {
		return
	}
	lhs, rhs := kids[0], kids[1]
	if lhs.Token() != rhino.GETPROP || rhs.Token() != rhino.OBJECTLIT {
		return
	}
	lhsKids := lhs.Children()
	if len(lhsKids) < 2 || lhsKids[0].Token() != rhino.NAME || lhsKids[1].String_() != "prototype" {
		return
	}
	ownerName := lhsKids[0].String_()
	fn, ok := scope.declared[ownerName].(*jstype.FunctionType)
	if !ok || fn.Instance == nil {
		return
	}
	for _, kv := range rhs.Children() {
		if kv.Token() != rhino.STRING_KEY {
			continue
		}
		child := kv.FirstChild()
		if child == nil {
			continue
		}
		memberType := c.inferPrototypeValueType(kv.JSDocInfo(), child, scope, fn.Instance)
		fn.Instance.DefineProperty(kv.String_(), memberType)
	}
}

// inferPrototypeValueType computes the member type for a
// `X.prototype.m = value` assignment. Function expressions yield a
// FunctionType derived from the ASSIGN's @param/@return JSDoc; other
// values fall back to @type if present, otherwise Unknown.
// inst is the constructor's InstanceType — when non-nil it is set as
// typeOfThis on the resulting FunctionType so TypeCheck can emit
// JSC_EXPECTED_THIS_TYPE when the method is called without a receiver.
// Go-only: language-adapter — inferPrototypeValueType is a Go-only helper that infers prototype member type as standalone method; Java inlines this in NormalScopeBuilder.defineObjectLiteral.
func (c *TypedScopeCreator) inferPrototypeValueType(info *rhino.JSDocInfo, rhs *rhino.Node, scope *TypedScope, inst *jstype.ObjectType) jstype.JSType {
	// When an explicit @type annotation is present it takes precedence
	// over the FUNCTION-body parameter / return inference. Mirrors
	// Java's TypedScopeCreator using JSDocInfo.getType() to override
	// the inferred FunctionType for prototype-method stubs.
	if info != nil && info.HasType() {
		return scope.ResolveExpr(info.Type().TypeName)
	}
	if rhs.Token() == rhino.FUNCTION {
		kids := rhs.Children()
		if len(kids) < 2 {
			return jstype.Unknown
		}
		params := kids[1]
		var paramTypes []jstype.JSType
		for _, p := range params.Children() {
			paramTypes = append(paramTypes, resolveParamType(info, p, scope))
		}
		// Default return type: Void (undefined) when no @return JSDoc is
		// present, matching recordFunction's Void-default policy and Java's
		// FunctionTypeBuilder.inferReturnType. Functions with an explicit
		// `return expr;` statement get Unknown (inferred) rather than Void.
		var ret jstype.JSType = jstype.Void
		if info != nil && info.ReturnType() != nil {
			ret = scope.ResolveExpr(info.ReturnType().TypeName)
		} else if functionReturnsExplicitly(rhs) {
			ret = jstype.Unknown
		}
		// Set typeOfThis so EXPECTED_THIS_TYPE fires when the method is
		// detached and called without a receiver. Mirrors Java's
		// FunctionType.getTypeOfThis() returning the instance type for
		// prototype methods (FunctionTypeBuilder.inferThisType).
		var typeOfThis jstype.JSType
		if inst != nil {
			typeOfThis = inst
		}
		fnT := jstype.NewFunctionTypeWithThis(paramTypes, ret, typeOfThis)
		setArityFromParams(fnT, params, info)
		return fnT
	}
	if info != nil && info.HasType() {
		return scope.ResolveExpr(info.Type().TypeName)
	}
	// Data slot without @type — infer from the RHS literal if we
	// can. Java's AssignmentInferrer flows the primitive type onto
	// the prototype's slot for later @return-style checks.
	if t := literalType(rhs); !isJSUnknown(t) {
		return t
	}
	return jstype.Unknown
}

// preregisterVarConstructor pre-registers `/** @constructor */ var Foo = fn`
// in Pass 1 so forward `@type {?Foo}` refs in the same scope resolve to the
// correct inner ConstructorType instead of falling back to an outer scope.
// Mirrors Java's two-pass constructor hoisting for VAR-assigned forms.
// Go-only: language-adapter — preregisterVarConstructor is a Go-only helper that pre-registers @constructor var declarations; Java inlines this in first-pass GatherScopeCreator traversal.
func (c *TypedScopeCreator) preregisterVarConstructor(v *rhino.Node, scope *TypedScope) {
	info := v.JSDocInfo()
	if !info.IsConstructor() && !info.IsInterface() && !info.IsRecord() {
		return
	}
	name := v.FirstChild()
	if name == nil || name.Token() != rhino.NAME || name.String_() == "" {
		return
	}
	init := name.FirstChild()
	if init == nil || init.Token() != rhino.FUNCTION {
		return
	}
	ct := jstype.NewConstructorType(name.String_(), nil)
	if info.IsInterface() || info.IsRecord() {
		ct.Instance.IsStructural = true
		ct.Instance.IsRecord = info.IsRecord()
	}
	if info.IsInterface() {
		ct.Instance.IsInterface = true
	}
	scope.declared[name.String_()] = ct
}

// preregisterConstructor reserves a slot in `scope` for any function
// carrying `@constructor`, `@interface`, or `@record`, so later
// `@type {Foo}` references resolve to the real InstanceType rather
// than collapsing to Unknown. The param/return types stay Unknown
// until pass 2 fills them in.
// Go-only: language-adapter — preregisterConstructor is a Go-only helper that pre-registers @constructor FUNCTION nodes; Java uses NormalScopeBuilder in a two-pass traversal via NodeTraversal.
func (c *TypedScopeCreator) preregisterConstructor(fn *rhino.Node, scope *TypedScope) {
	info := fn.JSDocInfo()
	if !info.IsConstructor() && !info.IsInterface() && !info.IsRecord() {
		return
	}
	kids := fn.Children()
	if len(kids) < 3 {
		return
	}
	name := kids[0]
	if name.String_() == "" {
		return
	}
	// Stub: real param/return types are computed in recordFunction.
	scope.declared[name.String_()] = jstype.NewConstructorType(name.String_(), nil)
}

// preregisterEnum reserves a slot for an `@enum`-annotated var /
// let / const so forward `@type {MyEnum}` references resolve to
// the shared EnumType. The Underlying primitive is computed from
// the `@enum {T}` argument (Number when omitted); aliased enums
// (`/** @enum */ var Alias = Source`) are handled in recordDecl's
// second pass, where we can see the RHS type.
// Go-only: language-adapter — preregisterEnum is a Go-only helper that pre-registers @enum declarations; Java uses NormalScopeBuilder.identifyEnumOrTypedefDeclaration.
func (c *TypedScopeCreator) preregisterEnum(v *rhino.Node, scope *TypedScope) {
	info := v.JSDocInfo()
	if !info.IsEnum() {
		return
	}
	name := v.FirstChild()
	if name == nil || name.Token() != rhino.NAME || name.String_() == "" {
		return
	}
	var underlying jstype.JSType = jstype.Number
	if expr := info.EnumParameterType(); expr != nil {
		underlying = scope.ResolveExpr(expr.TypeName)
		// Detect self-referential enum: `/** @enum {MyEnum} */ var MyEnum`.
		// Java's TypedScopeCreator emits JSC_CYCLIC_INHERITANCE_ERROR when
		// the enum's type argument names the enum itself. The position is
		// at the type argument inside the @enum braces in the JSDoc.
		if expr.TypeName == name.String_() && c.sourcePath != "" {
			line, col := enumTypeArgPos(c.sourceSrc, name.Lineno(), expr.TypeName)
			c.diags = append(c.diags, diag.Diagnostic{
				Path:    c.sourcePath,
				Line:    line,
				Col:     col,
				Code:    diag.JSC_CYCLIC_INHERITANCE_ERROR,
				Level:   diag.Error,
				Message: fmt.Sprintf("Cycle detected in inheritance chain of type %s", name.String_()),
			})
		}
	}
	scope.declared[name.String_()] = jstype.NewEnumType(name.String_(), underlying)
}

// enumTypeArgPos finds the line/col of the type-argument inside
// `@enum {TypeName}` on the given source line. Returns the name node's
// own position as fallback. lineno is 1-indexed; returned col is 0-indexed.
// Go-only: language-adapter — enumTypeArgPos is a Go-only helper that computes position within source for enum @type arg; Java uses SourceFile position directly from the compiler.
func enumTypeArgPos(src string, lineno int, typeName string) (line, col int) {
	// Split source into lines (0-indexed slice, 1-indexed lineno).
	lines := strings.Split(src, "\n")
	if lineno < 1 || lineno > len(lines) {
		return lineno, 0
	}
	l := lines[lineno-1]
	// Find @enum on this line, then the opening brace after it.
	idx := strings.Index(l, "@enum")
	if idx < 0 {
		return lineno, 0
	}
	brace := strings.Index(l[idx:], "{")
	if brace < 0 {
		return lineno, 0
	}
	return lineno, idx + brace + 1 // col of first char after '{'
}

// classMemberFunctionNode returns the FUNCTION node inside a class body member
// node. In the new pipeline, class methods are MEMBER_FUNCTION_DEF(FUNCTION);
// in the legacy pipeline, they are plain FUNCTION nodes.
// Returns nil if m is not a function-bearing member.
// Go-only: language-adapter — classMemberFunctionNode is a Go-only helper that extracts the FUNCTION node from a class member; Java accesses class member AST via Node.getLastChild() directly.
func classMemberFunctionNode(m *rhino.Node) *rhino.Node {
	if m == nil {
		return nil
	}
	if m.Token() == rhino.FUNCTION {
		return m
	}
	if m.Token() == rhino.MEMBER_FUNCTION_DEF || m.Token() == rhino.GETTER_DEF || m.Token() == rhino.SETTER_DEF {
		fc := m.FirstChild()
		if fc != nil && fc.Token() == rhino.FUNCTION {
			return fc
		}
	}
	return nil
}

// classMemberFunctionOrSelf is like classMemberFunctionNode but returns m itself
// when it is already a FUNCTION, and returns nil for non-method nodes.
// Used to normalize class body iteration across old and new pipeline ASTs.
// Go-only: language-adapter — classMemberFunctionOrSelf is a Go-only helper that returns FUNCTION node or the member itself; Java accesses the class member structure directly in NormalScopeBuilder.
func classMemberFunctionOrSelf(m *rhino.Node) *rhino.Node {
	if fn := classMemberFunctionNode(m); fn != nil {
		return fn
	}
	return m
}

// preregisterClass reserves a slot for a `class` declaration, mirroring
// the @constructor pre-registration path so forward `@type {C}`
// references resolve to C's InstanceType.
// Go-only: language-adapter — preregisterClass is a Go-only helper that pre-registers a CLASS node; Java uses a first-pass traversal in GatherScopeCreator to pre-register class names.
func (c *TypedScopeCreator) preregisterClass(cls *rhino.Node, scope *TypedScope) {
	kids := cls.Children()
	if len(kids) < 3 {
		return
	}
	name := kids[0]
	if name.String_() == "" {
		return
	}
	ct := jstype.NewConstructorType(name.String_(), nil)
	ct.IsES6Class = true
	scope.declared[name.String_()] = ct
}

// recordClass fills in a class's constructor params, methods, fields,
// and super link. Mirrors how TypedScopeCreator processes a CLASS node
// in Java (TypedScopeCreator.java, visitClass).
// Go-only: language-adapter — recordClass is a Go-only helper that processes a CLASS node to build its FunctionType; Java uses NormalScopeBuilder.visit for CLASS tokens in TypedScopeCreator.
func (c *TypedScopeCreator) recordClass(cls *rhino.Node, scope *TypedScope) {
	kids := cls.Children()
	if len(kids) < 3 {
		return
	}
	nameNode := kids[0]
	superNode := kids[1]
	body := kids[2]
	if nameNode.String_() == "" {
		return // anonymous classes — skip for Tier-2
	}
	preStub, ok := scope.declared[nameNode.String_()].(*jstype.FunctionType)
	if !ok || preStub.Instance == nil {
		return
	}
	inst := preStub.Instance
	// Mark as structural when the class carries @record or @interface.
	if info := cls.JSDocInfo(); info != nil && (info.IsRecord() || info.IsInterface()) {
		inst.IsStructural = true
		inst.IsRecord = info.IsRecord()
		if info.IsInterface() {
			inst.IsInterface = true
		}
	}
	if info := cls.JSDocInfo(); info != nil {
		if info.IsStruct() {
			inst.IsStruct = true
		}
		if info.IsDict() {
			inst.IsDict = true
		}
	}

	// Locate the `constructor` method; otherwise the class has an
	// implicit zero-arg constructor. In the new pipeline, class methods
	// are wrapped in MEMBER_FUNCTION_DEF; unwrap to find the FUNCTION node.
	var ctor *rhino.Node
	for _, mOuter := range body.Children() {
		fn := classMemberFunctionNode(mOuter)
		if fn == nil {
			continue
		}
		// Legacy: name is the FUNCTION's first child.
		if len(fn.Children()) >= 1 && fn.FirstChild().String_() == "constructor" {
			ctor = fn
			break
		}
		// New pipeline: name is the MEMBER_FUNCTION_DEF's String_().
		if mOuter.Token() == rhino.MEMBER_FUNCTION_DEF && mOuter.String_() == "constructor" {
			ctor = fn
			break
		}
	}
	var ctorParams []jstype.JSType
	if ctor != nil {
		ctorInfo := ctor.JSDocInfo()
		// In the new pipeline, JSDoc may be on the MEMBER_FUNCTION_DEF wrapper.
		if ctorInfo == nil && ctor.Parent() != nil {
			ctorInfo = ctor.Parent().JSDocInfo()
		}
		for _, p := range ctor.Children()[1].Children() {
			ctorParams = append(ctorParams, resolveParamType(ctorInfo, p, scope))
		}
	}
	preStub.Params = ctorParams
	preStub.MinArity = len(ctorParams)
	preStub.MaxArity = len(ctorParams)
	if ctor != nil {
		setArityFromParamNames(preStub, ctor.Children()[1])
	}

	// Wire @extends if present. Struct/dict flags propagate down the
	// inheritance chain — Java's FunctionType carries them as sticky
	// traits on the whole ancestry.
	if superNode.String_() != "" {
		if super := scope.LookupNominalType(superNode.String_()); super != nil {
			inst.SetSuper(super)
			if super.IsStruct {
				inst.IsStruct = true
			}
			if super.IsDict {
				inst.IsDict = true
			}
		}
	}

	// Walk body: methods attach to the InstanceType (or the static
	// slot when tagged `static`); fields attach with their declared
	// @type to the same slot.
	// Track getter/setter property types for JSC_CONFLICTING_GETTER_SETTER_TYPE.
	type accessorEntry struct {
		typ  jstype.JSType
		node *rhino.Node
	}
	getterTypes := map[string]accessorEntry{}
	setterTypes := map[string]accessorEntry{}
	for _, mOrig := range body.Children() {
		// In the new pipeline, class methods are MEMBER_FUNCTION_DEF nodes
		// wrapping a FUNCTION; fields are MEMBER_FIELD_DEF nodes.
		// In the legacy pipeline, methods are direct FUNCTION nodes and
		// fields are CLASS_FIELD_DEF nodes.
		// We normalize by checking the original token first.
		origTok := mOrig.Token()

		// Handle field declarations.
		if origTok == rhino.CLASS_FIELD_DEF || origTok == rhino.MEMBER_FIELD_DEF {
			info := mOrig.JSDocInfo()
			if info == nil {
				info = mOrig.FirstChild().JSDocInfo() // sometimes on inner node
			}
			fieldName := mOrig.String_()
			if fieldName == "" && mOrig.FirstChild() != nil {
				fieldName = mOrig.FirstChild().String_()
			}
			var t jstype.JSType = jstype.Unknown
			if info != nil && info.HasType() {
				t = scope.ResolveExpr(info.Type().TypeName)
			}
			if mOrig.IsStaticMember() || (info != nil && info.IsStatic()) {
				inst.DefineStaticProperty(fieldName, t)
			} else {
				inst.DefineProperty(fieldName, t)
			}
			continue
		}

		// Handle method declarations.
		mInfo := mOrig.JSDocInfo()
		isStatic := mOrig.IsStaticMember()
		isGetter := mOrig.Token() == rhino.GETTER_DEF
		isSetter := mOrig.Token() == rhino.SETTER_DEF
		m := classMemberFunctionOrSelf(mOrig)
		switch m.Token() {
		case rhino.FUNCTION:
			if len(m.Children()) < 2 {
				continue
			}
			methodName := m.FirstChild().String_()
			if methodName == "" {
				methodName = mOrig.String_() // from MEMBER_FUNCTION_DEF string value
			}
			if methodName == "constructor" {
				continue // handled above
			}
			info := m.JSDocInfo()
			if info == nil {
				info = mInfo // fall back to MEMBER_FUNCTION_DEF's JSDoc
			}
			if !isStatic && info != nil && info.IsStatic() {
				isStatic = true
			}
			params := m.Children()[1]
			var paramTypes []jstype.JSType
			for _, p := range params.Children() {
				paramTypes = append(paramTypes, resolveParamType(info, p, scope))
			}
			var ret jstype.JSType = jstype.Unknown
			if info != nil && info.ReturnType() != nil {
				ret = scope.ResolveExpr(info.ReturnType().TypeName)
			}
			fnT := jstype.NewFunctionType(paramTypes, ret)
			setArityFromParams(fnT, params, info)
			if isStatic {
				inst.DefineStaticProperty(methodName, fnT)
			} else {
				inst.DefineProperty(methodName, fnT)
			}
			// Record getter/setter property types for conflict detection.
			if isGetter || m.IsGetterDef() {
				var propType jstype.JSType = jstype.Unknown
				if info != nil && info.ReturnType() != nil {
					propType = ret
				}
				getterTypes[methodName] = accessorEntry{propType, mOrig}
			} else if isSetter || m.IsSetterDef() {
				var propType jstype.JSType = jstype.Unknown
				if len(paramTypes) > 0 {
					propType = paramTypes[0]
				}
				setterTypes[methodName] = accessorEntry{propType, mOrig}
			}
		case rhino.CLASS_FIELD_DEF, rhino.MEMBER_FIELD_DEF:
			// Should have been handled above, but catch stragglers.
			info := m.JSDocInfo()
			if info == nil {
				info = mInfo
			}
			fieldName := m.String_()
			var t jstype.JSType = jstype.Unknown
			if info != nil && info.HasType() {
				t = scope.ResolveExpr(info.Type().TypeName)
			}
			if isStatic || (info != nil && info.IsStatic()) {
				inst.DefineStaticProperty(fieldName, t)
			} else {
				inst.DefineProperty(fieldName, t)
			}
		}
	}
	// Emit JSC_CONFLICTING_GETTER_SETTER_TYPE when getter/setter types differ.
	// Mirrors Java TypeCheck.checkGetterOrSetterType. Skip when the class
	// carries @suppress {checkTypes}.
	clsInfo := cls.JSDocInfo()
	checkTypesSuppressed := clsInfo != nil && clsInfo.HasSuppression("checkTypes")
	for propName, ge := range getterTypes {
		se, hasSetter := setterTypes[propName]
		if !hasSetter {
			continue
		}
		getType := ge.typ
		setType := se.typ
		// Compare by TypeName() since primitive singletons print as "{}" via %v.
		// Report at the SETTER position (matching Java's checkGetterOrSetterType
		// which reports on the child that doesn't match the "official" type; the
		// getter's return type is registered first and the setter's param type
		// is then compared against it → report fires on the setter).
		if !checkTypesSuppressed && getType.TypeName() != setType.TypeName() {
			c.diags = append(c.diags, diag.Diagnostic{
				Path:  c.sourcePath,
				Line:  se.node.Lineno(),
				Col:   se.node.Charno(),
				Code:  diag.JSC_CONFLICTING_GETTER_SETTER_TYPE,
				Level: diag.Error,
				Message: fmt.Sprintf(
					"The types of the getter and setter for property '%s' do not match.\ngetter type is: %s\nsetter type is: %s",
					propName, getType.TypeName(), setType.TypeName()),
			})
		}
	}

	// Build local scopes for the constructor and each method so
	// TypeCheck can walk their bodies with `this` bound.
	c.Functions[cls] = &TypedScope{
		parent:   scope,
		declared: make(map[string]jstype.JSType),
		registry: c.registry,
		thisType: inst,
	}
	if ctor != nil {
		ctorScope := c.buildMethodScope(ctor, inst, scope, ctorParams)
		// Bind `super` inside the derived-class ctor to the base's
		// FunctionType so `super(args)` + `super.method(args)` flow
		// through the normal checkCall / checkGetProp paths.
		c.bindSuper(ctorScope, superNode.String_(), scope)
		c.Functions[ctor] = ctorScope
		// Also collect `this.x = …` property writes inside the
		// constructor body.
		if ctorBody := ctor.Children()[2]; ctorBody != nil {
			c.collectThisPropertiesIn(ctorBody, inst, scope)
		}
	}
	for _, m := range body.Children() {
		if m.Token() != rhino.FUNCTION || m.FirstChild().String_() == "constructor" {
			continue
		}
		params := m.Children()[1]
		var ptypes []jstype.JSType
		info := m.JSDocInfo()
		for _, p := range params.Children() {
			ptypes = append(ptypes, resolveParamType(info, p, scope))
		}
		c.Functions[m] = c.buildMethodScope(m, inst, scope, ptypes)
	}
}

// buildMethodScope wraps the parameter + this binding needed for a
// method body to type-check under the same local-scope model as
// top-level @constructor functions.
// Go-only: language-adapter — buildMethodScope is a Go-only helper that builds a parameter scope for a class method; Java uses NormalScopeBuilder and createScope(Node, AbstractScope) inline.
func (c *TypedScopeCreator) buildMethodScope(fn *rhino.Node, inst *jstype.ObjectType, parent *TypedScope, ptypes []jstype.JSType) *TypedScope {
	local := &TypedScope{
		parent:      parent,
		declared:    make(map[string]jstype.JSType),
		typeAliases: make(map[string]jstype.JSType),
		registry:    c.registry,
		thisType:    inst,
	}
	kids := fn.Children()
	if len(kids) < 2 {
		return local
	}
	params := kids[1]
	for i, p := range params.Children() {
		if i < len(ptypes) {
			local.declared[p.String_()] = ptypes[i]
		}
	}
	// Also hoist var declarations from the method body into the local scope
	// so VarCheck sees them as declared when traversing the body. Without
	// this, `var` declarations inside for-loops and nested blocks in a
	// prototype method body would be reported as JSC_UNDEFINED_VARIABLE.
	// We use a narrow hoist (only var/let/const names) rather than full
	// createScopeIn to avoid side-effects on the type system from
	// processing EXPR_RESULT nodes inside the method body.
	if len(kids) >= 3 {
		body := kids[2]
		c.hoistVarNamesInto(body, local)
	}
	return local
}

// hoistVarNamesInto walks a function body and adds every VAR/LET/CONST
// binding name to `scope` with Unknown type. This is a narrow variant of
// createScopeIn that only declares names (not resolves types), avoiding
// side-effects from processing EXPR_RESULT nodes inside method bodies.
// Used by buildMethodScope to ensure VarCheck sees `var i` inside a
// for-loop as declared when checking the loop body. Mirrors
// AbstractScopeBuilder.hoistVars's narrow first pass (names only).
// Go-only: language-adapter — hoistVarNamesInto is a Go-only helper that hoists `var` declarations into a scope; Java does this via NormalScopeBuilder walking the function body on demand.
func (c *TypedScopeCreator) hoistVarNamesInto(body *rhino.Node, scope *TypedScope) {
	if body == nil {
		return
	}
	for _, child := range body.Children() {
		switch child.Token() {
		case rhino.VAR, rhino.LET, rhino.CONST:
			if name := child.FirstChild(); name != nil && name.Token() == rhino.NAME && name.String_() != "" {
				if _, exists := scope.declared[name.String_()]; !exists {
					scope.declared[name.String_()] = jstype.Unknown
				}
			}
			c.hoistVarNamesInto(child, scope)
		case rhino.FOR, rhino.FOR_IN, rhino.FOR_OF:
			if first := child.FirstChild(); first != nil && NU.IsNameDeclaration(first) {
				if name := first.FirstChild(); name != nil && name.Token() == rhino.NAME && name.String_() != "" {
					if _, exists := scope.declared[name.String_()]; !exists {
						// Java's TypeInference: FOR_IN key is always string type;
						// FOR_OF key type is the element type of the Iterable.
						keyType := jstype.JSType(jstype.Unknown)
						if child.Token() == rhino.FOR_IN {
							keyType = jstype.String
						} else if child.Token() == rhino.FOR_OF {
							// Extract element type from the iterable (second child).
							kids := child.Children()
							if len(kids) >= 2 {
								iterType := scope.inferExprType(kids[1])
								if elem := extractIterableElemType(iterType); elem != nil {
									keyType = elem
								}
							}
						}
						scope.declared[name.String_()] = keyType
					}
				}
			}
			c.hoistVarNamesInto(child, scope)
		case rhino.BLOCK, rhino.IF, rhino.WHILE, rhino.DO_WHILE, rhino.SWITCH, rhino.CASE, rhino.TRY, rhino.FINALLY:
			c.hoistVarNamesInto(child, scope)
		case rhino.FUNCTION, rhino.CLASS:
			// Don't descend into nested function/class bodies.
		}
	}
}

// templateResolver returns (resolver, templateMap) — the resolver is
// a closure that resolves a type expression against a shim scope
// where each of the function's @template names is bound to a fresh
// TemplateType. The templateMap lets the caller stash those
// TemplateTypes on the resulting FunctionType so call-site
// unification can refer to them later. Port of the per-function
// TemplateTypeMap construction in TypedScopeCreator.java.
// Go-only: language-adapter — templateResolver is a Go-only helper that builds a @template name resolver for a scope; Java creates a JSTypeRegistry scope with templates via typeRegistry.createScopeWithTemplates.
func (c *TypedScopeCreator) templateResolver(parent *TypedScope, info *rhino.JSDocInfo) (func(string) jstype.JSType, map[string]*jstype.TemplateType) {
	if info == nil || len(info.TemplateNames()) == 0 {
		return parent.ResolveExpr, nil
	}
	tmpls := make(map[string]*jstype.TemplateType, len(info.TemplateNames()))
	aliases := make(map[string]jstype.JSType, len(info.TemplateNames()))
	for _, tn := range info.TemplateNames() {
		tt := jstype.NewTemplateType(tn)
		tmpls[tn] = tt
		aliases[tn] = tt
	}
	shim := &TypedScope{
		parent:      parent,
		declared:    map[string]jstype.JSType{},
		typeAliases: aliases,
		registry:    c.registry,
	}
	return shim.ResolveExpr, tmpls
}

// bindSuper resolves the named superclass into its constructor's
// FunctionType and stores it under `super` in the given (ctor-local)
// scope. Harmless no-op when the class has no base.
// Go-only: language-adapter — bindSuper is a Go-only helper that binds `super` in a class method scope; Java handles `super` via the implicit variable system in TypedScope.makeImplicitVar.
func (c *TypedScopeCreator) bindSuper(local *TypedScope, superName string, outer *TypedScope) {
	if superName == "" {
		return
	}
	baseCtor, ok := outer.declared[superName].(*jstype.FunctionType)
	if !ok {
		return
	}
	local.declared["super"] = baseCtor
}

// recordTypedef stores a `@typedef {T}`-annotated `var NAME;` under
// the scope's typeAliases map. Mirrors how Java's TypedScopeCreator
// registers typedefs ahead of the normal scope-building passes so
// `@type {NAME}` lookups resolve to the aliased type.
// Go-only: language-adapter — recordTypedef is a Go-only helper that records a @typedef declaration; Java uses NormalScopeBuilder.identifyEnumOrTypedefDeclaration via NodeTraversal.
func (c *TypedScopeCreator) recordTypedef(v *rhino.Node, scope *TypedScope) {
	info := v.JSDocInfo()
	if info == nil || info.TypedefType() == nil {
		return
	}
	name := v.FirstChild()
	if name == nil || name.Token() != rhino.NAME {
		return
	}
	scope.typeAliases[name.String_()] = scope.ResolveExpr(info.TypedefType().TypeName)
}

// Go-only: language-adapter — recordDecl is a Go-only helper that dispatches a VAR/LET/CONST declaration; Java uses NormalScopeBuilder.defineName inline during NodeTraversal.
func (c *TypedScopeCreator) recordDecl(v *rhino.Node, scope *TypedScope) {
	// Multi-binding forms (`var a, b = 1, c`) each declare their own
	// slot. Run the full recordOneDecl per NAME child so the scope
	// sees every name; inner helpers still use the first child as
	// the primary (matches Java's visitDeclarationsOrList loop).
	for _, child := range v.Children() {
		if child.Token() != rhino.NAME {
			continue
		}
		c.recordOneDecl(v, child, scope)
	}
}

// recordOneDecl processes a single declaration binding. Split out
// from recordDecl so the multi-binding loop can reuse the same
// @type / @enum / @constructor handling per name.
// Go-only: language-adapter — recordOneDecl is a Go-only helper that processes a single declaration binding; Java uses NormalScopeBuilder.defineName which handles both simple and destructured names.
func (c *TypedScopeCreator) recordOneDecl(v *rhino.Node, name *rhino.Node, scope *TypedScope) {
	// VAR-level JSDoc is the primary source; NAME-level JSDoc
	// (set by the parser for inline `var /** @type {T} */ x = …`
	// annotations) takes precedence for @type when present on the
	// specific NAME node. This mirrors Java's per-binding type
	// annotation: `var /** @type {!Date} */ foo = 'abc'` attaches
	// the @type to foo's NAME node, not the shared VAR.
	info := v.JSDocInfo()
	if nameInfo := name.JSDocInfo(); nameInfo != nil && nameInfo.HasType() && (info == nil || !info.HasType()) {
		info = nameInfo
	}
	if name == nil || name.Token() != rhino.NAME {
		return
	}
	// Inline JSDoc (`var /** @type {T} */ n`) may be attached to the NAME
	// node by the new pipeline's IRFactory. Fall back to the NAME's JSDoc
	// when the declaration container (VAR/LET/CONST) carries none.
	if info == nil {
		info = name.JSDocInfo()
	}
	// Typedef declarations primarily populate the typeAliases
	// namespace, but we still record the NAME at Unknown in the
	// value namespace so later `var alias = original;` references
	// don't fire UNDEFINED_VARIABLE. Java's VarCheck treats
	// `var original` (with @typedef) as a normal declaration slot.
	if info != nil && info.TypedefType() != nil {
		scope.declared[name.String_()] = jstype.Unknown
		return
	}
	// `/** @enum [{T}] */ var MyEnum = {FOO: 1, BAR: 2}` — the
	// preregister pass already allocated an EnumType for MyEnum;
	// the second-pass job is (a) scanning the initializer
	// OBJECTLIT to populate the Members set, (b) emitting a
	// TYPE_MISMATCH at the initializer when the RHS isn't an
	// object literal / matching enum. Aliased form
	// (`/** @enum */ var A = B`) reuses B's EnumType pointer so
	// `A.FOO` and `B.FOO` share the same static type.
	if info.IsEnum() && name.String_() != "" {
		c.finishEnumDecl(name, info, scope)
		return
	}
	// `/** @constructor */ var Foo = function() {}` — Closure-library
	// form of a constructor declaration. Builds a ConstructorType so
	// `new Foo()` and `@type {Foo}` both resolve to the same instance.
	// Mirrors TypedScopeCreator's function-expression-as-constructor
	// path in AbstractScopeBuilder.
	if (info.IsConstructor() || info.IsInterface() || info.IsRecord()) && name.String_() != "" {
		init := name.FirstChild()
		// Constructor alias: `@constructor var Alias = ExistingCtor`.
		// Java shares the existing ConstructorType so Alias and the
		// source are interchangeable.
		if init != nil && init.Token() == rhino.NAME {
			if src, ok := scope.GetDeclaredType(init.String_()).(*jstype.FunctionType); ok && src.Instance != nil {
				scope.declared[name.String_()] = src
				return
			}
		}
		// Non-function init (OBJECTLIT / NUMBER / …) — create a
		// ConstructorType slot anyway so visitVar's subtype check
		// fires JSC_TYPE_MISMATCH, and the name stays declared.
		if init != nil && init.Token() != rhino.FUNCTION {
			ct := jstype.NewConstructorType(name.String_(), nil)
			if info.IsInterface() || info.IsRecord() {
				ct.Instance.IsStructural = true
				ct.Instance.IsRecord = info.IsRecord()
			}
			if info.IsInterface() {
				ct.Instance.IsInterface = true
			}
			scope.declared[name.String_()] = ct
			return
		}
		if init != nil && init.Token() == rhino.FUNCTION {
			fnType := c.buildDeclaredFunctionConstructor(init, info, name.String_(), scope)
			// Reuse the Pass-1 stub's Instance so forward @type references and
			// `new Foo()` share the same ObjectType pointer. Mirrors recordFunction.
			if preStub, ok := scope.declared[name.String_()].(*jstype.FunctionType); ok && preStub.Instance != nil && preStub.Instance != fnType.Instance {
				fnType.Instance = preStub.Instance
				fnType.Return = preStub.Instance
			}
			fnType.SetJSDocInfo(info)
			if info.IsInterface() || info.IsRecord() {
				fnType.Instance.IsStructural = true
				fnType.Instance.IsRecord = info.IsRecord()
			}
			if info.IsInterface() {
				fnType.Instance.IsInterface = true
			}
			if info.IsStruct() {
				fnType.Instance.IsStruct = true
			}
			if info.IsDict() {
				fnType.Instance.IsDict = true
			}
			if base := info.BaseType(); base != nil {
				if super := scope.LookupNominalType(stripTypeArgs(base.TypeName)); super != nil {
					fnType.Instance.SetSuper(super)
					if super.IsStruct {
						fnType.Instance.IsStruct = true
					}
					if super.IsDict {
						fnType.Instance.IsDict = true
					}
				}
			}
			scope.declared[name.String_()] = fnType
			c.collectThisPropertiesIn(init.Children()[2], fnType.Instance, scope)
			return
		}
	}
	// `@type {T}` pins the binding's declared type precisely. Java's
	// VarCheck treats variables whose initializer references an
	// undeclared NAME as Unknown — the @type annotation is silently
	// discarded so subsequent uses don't double-report. We follow
	// the same rule by downgrading to Unknown when the initializer
	// is a bare reference we can't resolve.
	if info.HasType() {
		declared := scope.ResolveExpr(info.Type().TypeName)
		// When the declared type is Unknown but the raw annotation mentions
		// ReadonlyArray (e.g. cast `/** @type {!ReadonlyArray<T>} */ (expr)`),
		// create a synthetic ObjectType named "ReadonlyArray" so that
		// JSC_PROPERTY_ASSIGNMENT_TO_READONLY_VALUE can fire. Java registers
		// ReadonlyArray in its full type registry; Go does it on-demand only
		// for cast-hoisted declarations, to avoid false-positives from
		// adding ReadonlyArray as a global builtin (which would cause
		// JSC_TYPE_MISMATCH for `@const {!ReadonlyArray<T>} const a = []`).
		if info.IsCastHoisted() && isJSUnknown(declared) {
			rawTypeName := info.Type().TypeName
			if synthetic := syntheticReadonlyArrayType(rawTypeName); synthetic != nil {
				declared = synthetic
			}
			// else: keep declared as jstype.Unknown so VarCheck doesn't
			// fire for variables declared via cast-hoisted @type {?} or
			// other non-ReadonlyArray annotations.
		}
		if init := name.FirstChild(); init != nil && initIsUndeclaredRef(init, scope) {
			declared = jstype.Unknown
		}
		scope.declared[name.String_()] = declared
		return
	}
	// No @type: try to recover a sharper type from a simple initializer
	// (`var x = new Foo()` → Foo instance). For everything else, we
	// still record the name at Unknown so VarCheck doesn't flag later
	// references to it as undeclared. Mirrors how TypedScopeCreator
	// always adds a slot for every `var` even when its type stays
	// inferred.
	init := name.FirstChild()
	if init != nil && init.Token() == rhino.NEW {
		if callee := init.FirstChild(); callee != nil && callee.Token() == rhino.NAME {
			if fn, ok := scope.declared[callee.String_()].(*jstype.FunctionType); ok && fn.Instance != nil {
				scope.declared[name.String_()] = fn.Instance
				return
			}
		}
	}
	// Bare constructor alias: \`var FooAlias = Foo\` / \`const FooAlias = Foo\`
	// — when the RHS is a known constructor NAME, share the same
	// FunctionType pointer so FooAlias.prototype / \`new FooAlias()\`
	// reach Foo's Instance. Java tracks this through @const-aliased
	// nominal types.
	if init != nil && init.Token() == rhino.NAME {
		if src, ok := scope.GetDeclaredType(init.String_()).(*jstype.FunctionType); ok && src.Instance != nil {
			scope.declared[name.String_()] = src
			return
		}
		// Typedef alias: `const alias = original` where original has a
		// @typedef. Propagate the typedef into the alias slot so
		// `@type {alias}` resolves to the same underlying type.
		// Mirrors Java's TypedScopeCreator const-aliased typedef path.
		if alias := scope.LookupTypeAlias(init.String_()); alias != nil {
			scope.typeAliases[name.String_()] = alias
			scope.declared[name.String_()] = jstype.Unknown
			return
		}
		// Enum alias: `const MyEnum = YourEnum` where YourEnum has
		// @enum. Share the EnumType so `@type {MyEnum}` resolves.
		if enum, ok := scope.GetDeclaredType(init.String_()).(*jstype.EnumType); ok {
			scope.declared[name.String_()] = enum
			return
		}
		// Instance alias: `var y = x` where x is typed as a structural
		// ObjectType (@record / @interface). Propagate the instance type
		// so that subsequent `/** @type {T} */ y.prop = val` assignments
		// in recordNamespaceAssign can register the property on the
		// ObjectType. Mirrors Java's TypedScopeCreator which propagates
		// variable types through assignments via getDeclaredType.
		if obj, ok := scope.GetDeclaredType(init.String_()).(*jstype.ObjectType); ok && obj.IsStructural {
			scope.declared[name.String_()] = obj
			return
		}
	}
	// `var foo = function bar(…){…}` — build a FunctionType from the
	// function's shape even when JSDoc is absent. Java's TypedScopeCreator
	// always registers a callable type for function-valued variables so
	// `new f` on a plain (non-constructor) function expression correctly
	// fires JSC_NOT_A_CONSTRUCTOR. Matches NodeUtil.getBestJSDocInfo
	// inheriting the VAR's JSDoc onto the function expression.
	if init != nil && init.Token() == rhino.FUNCTION {
		kids := init.Children()
		if len(kids) >= 2 {
			params := kids[1]
			var paramTypes []jstype.JSType
			for _, p := range params.Children() {
				paramTypes = append(paramTypes, resolveParamType(info, p, scope))
			}
			var retType jstype.JSType = jstype.Void
			if info != nil {
				if r := info.ReturnType(); r != nil {
					retType = scope.ResolveExpr(r.TypeName)
				} else if functionReturnsExplicitly(init) {
					retType = jstype.Unknown
				}
			} else if functionReturnsExplicitly(init) {
				retType = jstype.Unknown
			}
			fnType := jstype.NewFunctionType(paramTypes, retType)
			setArityFromParams(fnType, params, info)
			scope.declared[name.String_()] = fnType
			return
		}
	}
	// GETPROP init: `const E = ns.Prop` — try to resolve the qualified
	// name. If it resolves to something non-Unknown (e.g. an enum or
	// constructor), carry that type over so downstream checks work.
	// When the GETPROP can't be resolved (e.g. `const E = Foo.E` where
	// Foo.E is a static enum we don't track), fall through and let the
	// caller see the slot as unknown-but-non-empty (use All to signal
	// "could be a type, don't flag as unknown value binding").
	// Mirrors Java's TypedScopeCreator const-GETPROP aliasing.
	if init != nil && init.Token() == rhino.GETPROP {
		t := resolveQualifiedType(scope, init)
		if t != nil && !isJSUnknown(t) {
			scope.declared[name.String_()] = t
		} else {
			// Unknown GETPROP init — use All to distinguish from a
			// plain-value Unknown (avoids false-positive in the
			// ResolveTypesInScope emission pass).
			scope.declared[name.String_()] = jstype.All
		}
		return
	}
	if _, exists := scope.declared[name.String_()]; !exists {
		// `undefined` is a special well-known value with type void.
		if name.String_() == "undefined" {
			scope.declared[name.String_()] = jstype.Void
		} else if p := v.Parent(); p != nil && p.Token() == rhino.FOR_IN {
			// FOR_IN key type is string (property names of the object).
			// Mirrors TypeInference.java case FOR_IN: iterKeyType = STRING_TYPE.
			scope.declared[name.String_()] = jstype.String
		} else if p := v.Parent(); p != nil && p.Token() == rhino.FOR_OF {
			// Extract element type from the Iterable (second child of FOR_OF).
			kids := p.Children()
			if len(kids) >= 2 {
				iterType := scope.inferExprType(kids[1])
				if elem := extractIterableElemType(iterType); elem != nil {
					scope.declared[name.String_()] = elem
					return
				}
			}
			scope.declared[name.String_()] = jstype.Unknown
		} else {
			scope.declared[name.String_()] = jstype.Unknown
		}
	}
}


// extractIterableElemType returns the element type T from an Iterable<T>,
// Array<T>, or similar iteratable type. For union types, returns the union
// of element types from each member. Returns nil if not determinable.
// Mirrors Java's TypeInference case FOR_OF element-type extraction.
// Go-only: language-adapter — extractIterableElemType is a Go-only helper that extracts element type from an iterable type; Java uses JSTypeRegistry.getIteratedTypeOfIterable(JSType).
func extractIterableElemType(t jstype.JSType) jstype.JSType {
	if t == nil {
		return nil
	}
	// Handle TemplatizedType: Array<T>, Iterable<T>, etc.
	if tt, ok := t.(*jstype.TemplatizedType); ok && len(tt.Args) > 0 {
		return tt.Args[0]
	}
	// Handle union type: collect element types from each member and union them.
	if u, ok := t.(*jstype.UnionType); ok {
		var elems []jstype.JSType
		for _, m := range u.Members() {
			if elem := extractIterableElemType(m); elem != nil {
				elems = append(elems, elem)
			}
		}
		if len(elems) == 1 {
			return elems[0]
		}
		if len(elems) > 1 {
			return jstype.NewUnion(elems...)
		}
	}
	return nil
}

// finishEnumDecl resolves `/** @enum */ var MyEnum = …` by
// completing the pre-registered EnumType: aliasing replaces the
// slot's EnumType with the RHS enum's pointer, object-literal
// initialisers populate the Members map. Matches how Java's
// TypedScopeCreator.EnumBuilder hoists the enum slot in pass 1
// and fills its members in pass 2.
// Go-only: language-adapter — finishEnumDecl is a Go-only helper that finalizes an @enum type declaration; Java uses NormalScopeBuilder.createEnumTypeFromNodes inline.
func (c *TypedScopeCreator) finishEnumDecl(nameNode *rhino.Node, info *rhino.JSDocInfo, scope *TypedScope) {
	name := nameNode.String_()
	enum, ok := scope.declared[name].(*jstype.EnumType)
	if !ok {
		// Pre-register missed this binding (e.g. nested inside a
		// block and the hoist pass didn't pick it up). Create one.
		var underlying jstype.JSType = jstype.Number
		if expr := info.EnumParameterType(); expr != nil {
			underlying = scope.ResolveExpr(expr.TypeName)
		}
		enum = jstype.NewEnumType(name, underlying)
		scope.declared[name] = enum
	}
	init := nameNode.FirstChild()
	if init == nil {
		return
	}
	// Alias form: the RHS is another enum NAME reference. Replace
	// this enum's slot with the source enum so `Alias.FOO` and
	// `Source.FOO` share their EnumElementType pointer (matches
	// how Java treats `/** @enum */ var A = B`).
	if init.Token() == rhino.NAME {
		if other, ok := scope.declared[init.String_()].(*jstype.EnumType); ok {
			scope.declared[name] = other
			return
		}
	}
	// Literal form: scan the OBJECTLIT entries and record each
	// STRING_KEY child on the enum's Members set.
	if init.Token() == rhino.OBJECTLIT {
		for _, kv := range init.Children() {
			if kv.Token() == rhino.STRING_KEY {
				enum.Members[kv.String_()] = true
			}
		}
	}
}

// setArityFromParamNames computes MinArity / MaxArity on `fn` from
// the param-name prefixes of its PARAM_LIST and the JSDoc on
// `info` (optional). Closure-convention:
//
//   - `opt_<name>` — optional; MinArity stops counting at the
//     first opt_ param.
//   - `var_args`    — rest; MaxArity becomes MaxArgsUnlimited.
//
// JSDoc takes precedence when present:
//
//   - `@param {...T} x`  — rest; same effect as var_args.
//
// Matches FunctionParamBuilder's prefix detection path in Java.
// Called for every FunctionType we build from a PARAM_LIST node.
// Go-only: language-adapter — setArityFromParamNames is a Go-only helper that sets arity from param name nodes; Java uses FunctionType.Builder.withParamsNode in createFunctionTypeFromNodes.
func setArityFromParamNames(fn *jstype.FunctionType, params *rhino.Node) {
	setArityFromParams(fn, params, nil)
}

// setArityFromParams is the doc-aware variant of
// setArityFromParamNames; callers pass the owning @param-bearing
// JSDocInfo so the `{...T}` rest marker can widen MaxArity.
// resolveParamType returns the declared type for a FUNCTION's
// parameter node, honoring (in order):
//  1. Inline JSDoc on the param itself: `function f(/** T */ x)`.
//  2. Outer function @param {T} x JSDoc on `info`.
//  3. Unknown (unannotated).
// Go-only: language-adapter — resolveParamType is a Go-only helper that resolves a @param type for a function parameter; Java resolves this via JSDocInfo.getParameterType inline in createFunctionTypeFromNodes.
func resolveParamType(info *rhino.JSDocInfo, p *rhino.Node, scope *TypedScope) jstype.JSType {
	if p != nil {
		if pj := p.JSDocInfo(); pj != nil && pj.HasType() {
			return scope.ResolveExpr(pj.Type().TypeName)
		}
	}
	if info != nil {
		if expr := info.ParameterType(p.String_()); expr != nil {
			return scope.ResolveExpr(expr.TypeName)
		}
	}
	return jstype.Unknown
}

// Go-only: language-adapter — setArityFromParams is a Go-only helper that sets arity from both param nodes and JSDoc; Java handles this inline in NormalScopeBuilder.createFunctionTypeFromNodes.
func setArityFromParams(fn *jstype.FunctionType, params *rhino.Node, info *rhino.JSDocInfo) {
	if params == nil {
		return
	}
	kids := params.Children()
	fn.MinArity = 0
	fn.MaxArity = len(kids)
	sawOptional := false
	for i, p := range kids {
		name := p.String_()
		// JSDoc explicit rest marker ({...T}) takes highest priority.
		if info != nil && info.IsRestParam(name) {
			fn.MaxArity = jstype.MaxArgsUnlimited
			if !sawOptional {
				fn.MinArity = i
			}
			return
		}
		// JSDoc explicit optional marker ({T=}) is next.
		if info != nil && info.IsOptionalParam(name) {
			sawOptional = true
			continue
		}
		// Naming conventions (opt_, var_args) apply only when there is NO
		// explicit @param JSDoc for this name — mirrors Java's behaviour
		// where an explicit @param {T} opt_y (without =) keeps the param
		// required. When no JSDoc exists at all, the naming convention
		// determines optionality.
		hasExplicitParam := info != nil && info.ParameterType(name) != nil
		if !hasExplicitParam {
			if name == "var_args" || (len(name) > 8 && name[:8] == "var_args") {
				fn.MaxArity = jstype.MaxArgsUnlimited
				if !sawOptional {
					fn.MinArity = i
				}
				return
			}
			if len(name) >= 4 && name[:4] == "opt_" {
				sawOptional = true
				continue
			}
		}
		// Required param.
		if !sawOptional {
			fn.MinArity = i + 1
		}
	}
	if !sawOptional {
		fn.MinArity = len(kids)
	}
}

// buildDeclaredFunctionConstructor turns a `var Foo = function(a,b){…}`
// node into a FunctionType whose params are resolved from the
// @param JSDoc on the surrounding VAR declaration and which carries
// a fresh InstanceType named after the variable. Mirrors how Java's
// TypedScopeCreator promotes a plain `function(){}` on the RHS of
// a constructor-annotated `var`.
// Go-only: language-adapter — buildDeclaredFunctionConstructor is a Go-only helper that builds a FunctionType for a declared @constructor; Java uses NormalScopeBuilder.createFunctionTypeFromNodes.
func (c *TypedScopeCreator) buildDeclaredFunctionConstructor(fn *rhino.Node, info *rhino.JSDocInfo, name string, scope *TypedScope) *jstype.FunctionType {
	kids := fn.Children()
	var paramTypes []jstype.JSType
	if len(kids) >= 2 {
		params := kids[1]
		for _, p := range params.Children() {
			paramTypes = append(paramTypes, resolveParamType(info, p, scope))
		}
	}
	return jstype.NewConstructorType(name, paramTypes)
}

// Go-only: language-adapter — recordFunction is a Go-only helper that records a FUNCTION node and builds its type; Java uses NormalScopeBuilder.visit for FUNCTION tokens and createFunctionTypeFromNodes.
func (c *TypedScopeCreator) recordFunction(fn *rhino.Node, scope *TypedScope) {
	info := fn.JSDocInfo()
	kids := fn.Children()
	if len(kids) < 3 {
		return
	}
	name := kids[0]
	params := kids[1]
	body := kids[2]

	// @template names bind to fresh TemplateType values inside a shim
	// scope; the same TemplateType pointers are recorded on the
	// resulting FunctionType so call-site unification can refer to
	// them. Mirrors how Java's JSTypeRegistry builds a
	// TemplateTypeMap before resolving per-function annotations.
	resolve, templateMap := c.templateResolver(scope, info)
	var paramTypes []jstype.JSType
	for _, p := range params.Children() {
		var t jstype.JSType = jstype.Unknown
		// Inline param JSDoc (`function f(/** T */ x)`) wins over
		// absence of an outer @param on the function.
		if pj := p.JSDocInfo(); pj != nil && pj.HasType() {
			t = scope.ResolveExpr(pj.Type().TypeName)
		} else if pexpr := info.ParameterType(p.String_()); pexpr != nil {
			t = resolve(pexpr.TypeName)
		}
		paramTypes = append(paramTypes, t)
	}
	// JS semantics for a function without `return expr` is to return
	// `undefined`. When @return is absent we default to Void instead
	// of Unknown so a `+` / assignment chain that consumes the return
	// value (`a + f()`) still type-checks against `string` /
	// `number`. Constructors keep the existing Instance return.
	var retType jstype.JSType = jstype.Void
	if r := info.ReturnType(); r != nil {
		retType = resolve(r.TypeName)
	} else if functionReturnsExplicitly(fn) {
		// Body has at least one `return expr;` — without dataflow we
		// can't infer the return type cheaply, so stay Unknown.
		retType = jstype.Unknown
	}

	var fnType *jstype.FunctionType
	// @record / @interface functions are also InstanceType carriers.
	if (info.IsConstructor() || info.IsInterface() || info.IsRecord()) && name.String_() != "" {
		// Pass 1 stub already registered the InstanceType under this
		// name; reuse the same Instance pointer so forward
		// `@type {Foo}` refs resolve to the same object as `new Foo()`.
		preStub, _ := scope.declared[name.String_()].(*jstype.FunctionType)
		fnType = &jstype.FunctionType{
			Params:        paramTypes,
			Return:        preStub.Instance,
			IsConstructor: true,
			Instance:      preStub.Instance,
		}
		// Wire @extends: look up the base ctor's Instance and set as
		// super on this ctor's Instance, mirroring
		// FunctionType.setExtendedClass in Java. @struct / @dict flags
		// inherit down the super chain. Strip template args from the base
		// type name (e.g. `Foo<string>` → `Foo`) before lookup so that
		// generic-instantiation annotations don't prevent super resolution.
		// Mirrors TypedScopeCreator's stripTemplateArgs call.
		if base := info.BaseType(); base != nil {
			if super := scope.LookupNominalType(stripTypeArgs(base.TypeName)); super != nil {
				fnType.Instance.SetSuper(super)
				if super.IsStruct {
					fnType.Instance.IsStruct = true
				}
				if super.IsDict {
					fnType.Instance.IsDict = true
				}
			}
		}
		// Mark structural when annotated @record / @interface.
		if info.IsInterface() || info.IsRecord() {
			fnType.Instance.IsStructural = true
			fnType.Instance.IsRecord = info.IsRecord()
		}
		if info.IsInterface() {
			fnType.Instance.IsInterface = true
		}
		// Propagate @struct / @dict onto the instance so TypeCheck can
		// gate JSC_ILLEGAL_PROPERTY_ACCESS at GETPROP / GETELEM sites.
		if info.IsStruct() {
			fnType.Instance.IsStruct = true
		}
		if info.IsDict() {
			fnType.Instance.IsDict = true
		}
	} else if info.HasType() {
		// `@type {function(this:T, …)}` on a plain function — resolve the
		// declared type so typeOfThis is preserved. When the resolved type
		// is a FunctionType, use it directly; otherwise fall back to the
		// param/return-inferred type. Mirrors Java's TypedScopeCreator
		// using the JSDoc type annotation over the inferred shape.
		if declFn, ok := scope.ResolveExpr(info.Type().TypeName).(*jstype.FunctionType); ok {
			fnType = declFn
		} else {
			fnType = jstype.NewFunctionType(paramTypes, retType)
		}
	} else {
		fnType = jstype.NewFunctionType(paramTypes, retType)
	}
	setArityFromParams(fnType, params, info)
	if templateMap != nil {
		fnType.TemplateTypes = make([]*jstype.TemplateType, 0, len(templateMap))
		for _, tn := range info.TemplateNames() {
			if tt := templateMap[tn]; tt != nil {
				fnType.TemplateTypes = append(fnType.TemplateTypes, tt)
			}
		}
	}
	if name.String_() != "" {
		scope.declared[name.String_()] = fnType
	}

	// Per-function local scope contains the parameters and a binding
	// for `this` (only meaningful inside constructors).
	localScope := &TypedScope{
		parent:   scope,
		declared: make(map[string]jstype.JSType),
		registry: c.registry,
		thisType: fnType.Instance, // nil for non-ctors
	}
	for i, p := range params.Children() {
		if p.Token() == rhino.NAME {
			if i < len(paramTypes) {
				localScope.declared[p.String_()] = paramTypes[i]
			}
		} else {
			// Destructuring parameter: `({a, b}) =>` or `([x, y]) =>`.
			// Collect all NAME bindings from the pattern and declare them
			// as Unknown (we don't have per-binding type info without full
			// destructuring type inference). Mirrors Java's
			// TypedScopeCreator.AbstractScopeBuilder.visitParameters()
			// handling of non-simple param patterns.
			collectPatternNames(p, localScope)
		}
	}
	// Declare @template names in the function body scope (Java's
	// declareTemplatesInFunctionBody) so `/** @type {T} */` inside
	// the body doesn't fire JSC_UNRECOGNIZED_TYPE_ERROR.
	c.declareTemplatesInFunctionBody(fn, info, localScope, scope)
	// Walk the body for `this.x = …` patterns so we can populate the
	// instance type's member map. Mirrors TypedScopeCreator's
	// AbstractScopeBuilder.attachLiteralTypes path.
	if fnType.Instance != nil {
		c.collectThisPropertiesIn(body, fnType.Instance, scope)
	}
	nestedScope := c.createScopeIn(body, localScope)
	for k, v := range nestedScope.declared {
		if _, exists := localScope.declared[k]; !exists {
			localScope.declared[k] = v
		}
	}
	c.Functions[fn] = localScope
}

// declareTemplatesInFunctionBody ports Java's TypedScopeCreator
// FunctionScopeBuilder.declareTemplatesInFunctionBody. Declares each
// @template name from the function's JSDoc as Void in the local scope so
// `/** @type {T} */ var x` inside the body doesn't fire
// JSC_UNRECOGNIZED_TYPE_ERROR. Void defers concrete unification to the
// call-site. Mirrors Java's `currentScope.declare(name, fn, voidType, ...)`.
// Go-only: language-adapter — declareTemplatesInFunctionBody is a Go-only helper that declares @template types in a function body scope; Java creates a template scope via typeRegistry.createScopeWithTemplates.
func (c *TypedScopeCreator) declareTemplatesInFunctionBody(_ *rhino.Node, info *rhino.JSDocInfo, local, _ *TypedScope) {
	if info == nil {
		return
	}
	for _, tname := range info.TemplateNames() {
		if tname != "" {
			if _, exists := local.declared[tname]; !exists {
				local.declared[tname] = jstype.Void
			}
		}
	}
}

// collectPatternNames recursively visits a destructuring pattern node
// (OBJECT_PATTERN, ARRAY_PATTERN, or DEFAULT_VALUE) and declares every
// NAME leaf as Unknown in the supplied scope. This mirrors Java's
// TypedScopeCreator.AbstractScopeBuilder.visitParameters() for
// non-simple (destructuring) function parameters so that VarCheck
// and VariableReferenceCheck do not flag the bound names as undeclared
// when they appear as references inside the function body.
// Go-only: language-adapter — collectPatternNames is a Go-only helper that collects names from destructuring patterns; Java handles destructuring in NormalScopeBuilder.defineName for LET/CONST/VAR.
func collectPatternNames(n *rhino.Node, scope *TypedScope) {
	if n == nil {
		return
	}
	switch n.Token() {
	case rhino.NAME:
		// Leaf binding — register in scope unless already declared.
		name := n.String_()
		if name != "" {
			if _, exists := scope.declared[name]; !exists {
				scope.declared[name] = jstype.Unknown
			}
		}
	case rhino.OBJECT_PATTERN:
		for _, ch := range n.Children() {
			collectPatternNames(ch, scope)
		}
	case rhino.ARRAY_PATTERN:
		for _, ch := range n.Children() {
			collectPatternNames(ch, scope)
		}
	case rhino.STRING_KEY:
		// `{key: pattern}` or `{key}` shorthand — descend into the
		// value child (the binding target).
		for _, ch := range n.Children() {
			collectPatternNames(ch, scope)
		}
	case rhino.DEFAULT_VALUE:
		// `{a = dflt}` — first child is the binding, second is the
		// default expression (a reference, not a binding).
		kids := n.Children()
		if len(kids) >= 1 {
			collectPatternNames(kids[0], scope)
		}
		// Do NOT recurse into kids[1] — it's a reference expression.
	case rhino.OBJECT_REST:
		// `{...rest}` — the single child is the binding name.
		for _, ch := range n.Children() {
			collectPatternNames(ch, scope)
		}
	}
}

// collectThisProperties walks a constructor body looking for
// `this.<name> = <expr>` statements (with optional `@type` JSDoc) and
// adds each property to the instance's member map. Mirrors how
// TypedScopeCreator's AssignmentInferrer attaches inferred property
// types to `this`.
// Go-only: language-adapter — collectThisProperties is a Go-only helper that collects `this.prop` assignments in a constructor; Java uses AssignmentInferrer which runs separately from scope creation.
func (c *TypedScopeCreator) collectThisProperties(node *rhino.Node, inst *jstype.ObjectType) {
	c.collectThisPropertiesIn(node, inst, nil)
}

// resolveTypeExpr resolves `expr` against `scope` when non-nil,
// falling back to the primitive-only registry. Used by the this-
// property collector so inline \`@type {Name}\` JSDocs see the
// scope's nominal bindings (Number, user classes, etc.).
// Go-only: language-adapter — resolveTypeExpr is a Go-only helper that resolves a type expression string in a scope; Java uses JSTypeRegistry.getType(StaticTypedScope, String).
func resolveTypeExpr(scope *TypedScope, reg *jstype.Registry, expr string) jstype.JSType {
	if scope != nil {
		return scope.ResolveExpr(expr)
	}
	return reg.ResolveExpr(expr)
}

// collectThisPropertiesIn is the scope-aware variant. When a scope
// is passed, \`@type {Name}\` references inside \`this.X = …\` JSDocs
// resolve through the scope chain (which knows about \`Number\`,
// user-defined classes, etc.) — not through a bare primitive-only
// registry. Callers without a scope reference use the plain
// collectThisProperties overload above, whose behaviour is the
// primitive-registry fallback.
// Go-only: language-adapter — collectThisPropertiesIn is a Go-only helper that recursively collects `this.prop` assignments; Java uses AssignmentInferrer which walks the AST separately.
func (c *TypedScopeCreator) collectThisPropertiesIn(node *rhino.Node, inst *jstype.ObjectType, scope *TypedScope) {
	if node == nil {
		return
	}
	if node.Token() == rhino.EXPR_RESULT {
		if a := node.FirstChild(); a != nil && a.Token() == rhino.ASSIGN {
			kids := a.Children()
			if len(kids) >= 2 && kids[0].Token() == rhino.GETPROP {
				obj := kids[0].FirstChild()
				if obj != nil && obj.Token() == rhino.THIS {
					prop := kids[0].Children()[1]
					existing := inst.GetPropertyType(prop.String_())
					alreadyDeclared := existing != nil && !isJSUnknown(existing) && !inst.IsPropertyInferred(prop.String_())
					info := a.JSDocInfo()
					if info == nil {
						info = node.JSDocInfo()
					}
					hasInfoType := info != nil && info.HasType()
					// A class-field \`/** @type {T} */ x = …\` (declared,
					// not inferred) wins over any \`this.x = …\`
					// inside the ctor body. Leaves the existing slot
					// untouched so visitAssign fires on incompatible
					// writes. Matches Java's class-field authority rule.
					if alreadyDeclared {
						// fall through: don't overwrite
						_ = hasInfoType
					} else {
						var t jstype.JSType = jstype.Unknown
						if hasInfoType {
							t = existing
							if t == nil {
								t = jstype.Unknown
							}
							declared := resolveTypeExpr(scope, c.registry, info.Type().TypeName)
							if !isJSUnknown(declared) {
								t = declared
							}
						} else {
							if existing == nil || isJSUnknown(existing) {
								t = literalType(kids[1])
								inst.MarkInferred(prop.String_())
							} else {
								t = existing
							}
						}
						inst.DefineProperty(prop.String_(), t)
					}
				}
			}
		}
		// Stub form: `/** @type {T} */ this.prop;` (no RHS). Java's
		// TypedScopeCreator records this as a property declaration on
		// the enclosing instance. Matches how @record / @interface
		// classes define their surface via constructor-body stubs.
		if gp := node.FirstChild(); gp != nil && gp.Token() == rhino.GETPROP {
			gpKids := gp.Children()
			if len(gpKids) >= 2 && gpKids[0].Token() == rhino.THIS {
				propName := gpKids[1].String_()
				var t jstype.JSType = jstype.Unknown
				info := node.JSDocInfo()
				if info == nil {
					info = gp.JSDocInfo()
				}
				if info != nil && info.HasType() {
					declared := resolveTypeExpr(scope, c.registry, info.Type().TypeName)
					if !isJSUnknown(declared) {
						t = declared
					}
				}
				inst.DefineProperty(propName, t)
			}
		}
	}
	for _, k := range node.Children() {
		c.collectThisPropertiesIn(k, inst, scope)
	}
}

// literalType returns the JSType of a literal or simple RHS node.
// Used by collectThisProperties to sharpen `this.x = <lit>` slot
// types when no @type JSDoc is present.
// Go-only: language-adapter — literalType is a Go-only helper that maps literal node tokens to primitive types; Java uses TypedScopeCreator.getDeclaredTypeInAnnotation inline for literal nodes.
func literalType(n *rhino.Node) jstype.JSType {
	if n == nil {
		return jstype.Unknown
	}
	switch n.Token() {
	case rhino.NUMBER:
		return jstype.Number
	case rhino.STRINGLIT:
		return jstype.String
	case rhino.TRUE, rhino.FALSE:
		return jstype.Boolean
	case rhino.NULL:
		return jstype.Null
	case rhino.UNDEFINED:
		return jstype.Void
	}
	return jstype.Unknown
}

// initIsUndeclaredRef reports whether `init` is a bare NAME node that
// resolves to no value-namespace binding in `scope` — used so
// recordDecl can downgrade `var x = unknown;` to Unknown when
// `unknown` itself isn't reachable.
// Go-only: language-adapter — initIsUndeclaredRef is a Go-only helper that checks if an initializer is an undeclared reference; Java uses similar logic inline in NormalScopeBuilder.defineName.
func initIsUndeclaredRef(init *rhino.Node, scope *TypedScope) bool {
	if init == nil || init.Token() != rhino.NAME {
		return false
	}
	if scope.GetDeclaredType(init.String_()) != nil {
		return false
	}
	if defaultExternNames()[init.String_()] {
		return false
	}
	return true
}

// functionReturnsExplicitly reports whether a FUNCTION body contains
// any `return EXPR` statement. Used so `function f(){}` is treated as
// returning Void while `function g(){return x;}` stays Unknown when
// no @return JSDoc is present. Mirrors the body-walk at the start of
// TypedScopeCreator's recordFunction.
// Go-only: language-adapter — functionReturnsExplicitly is a Go-only helper that checks if a function has explicit returns; Java uses NodeUtil.functionHasExplicitReturn as a static utility.
func functionReturnsExplicitly(fn *rhino.Node) bool {
	kids := fn.Children()
	if len(kids) < 3 {
		return false
	}
	body := kids[2]
	return hasReturnExpr(body)
}

// Go-only: language-adapter — hasReturnExpr is a Go-only helper that recursively checks for return-with-expression nodes; Java uses NodeUtil.functionHasExplicitReturn which calls into CheckMissingReturn logic.
func hasReturnExpr(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	if n.Token() == rhino.FUNCTION {
		return false // don't descend into nested functions
	}
	if n.Token() == rhino.RETURN && n.FirstChild() != nil {
		return true
	}
	for _, c := range n.Children() {
		if hasReturnExpr(c) {
			return true
		}
	}
	return false
}


// typeNameResolvesInScope reports whether a type name is valid when
// resolved strictly as a TYPE — not just as any declared identifier.
// This is stricter than CheckJsDocTypes.nameResolves which also
// accepts any identifier that appears anywhere in source text (the
// allNames gate). We check only: primitive, built-in extern, typedef,
// enum, constructor, or a dotted namespace ctor reference.
// Go-only: language-adapter — typeNameResolvesInScope is a Go-only helper that checks if a type name is resolvable in scope; Java resolves type names eagerly through JSTypeRegistry during scope building.
func typeNameResolvesInScope(scope *TypedScope, name string) bool {
	// Structural keywords that can appear inside a JSDoc expression.
	if name == "function" || name == "this" || name == "new" || name == "typeof" {
		return true
	}
	// Primitive / sentinel names (canonical + capitalised aliases).
	switch name {
	case "number", "string", "boolean", "null", "undefined", "void",
		"symbol", "bigint", "any",
		"Null", "Undefined", "Void", "Any":
		return true
	}
	// Well-known built-in externs that are always available.
	switch name {
	case "Object", "Array", "Function", "Date", "RegExp", "Error",
		"Map", "Set", "WeakMap", "WeakSet", "WeakRef",
		"Promise", "Symbol", "BigInt",
		"Math", "JSON", "Reflect", "Proxy",
		"Int8Array", "Uint8Array", "Uint8ClampedArray",
		"Int16Array", "Uint16Array",
		"Int32Array", "Uint32Array",
		"Float32Array", "Float64Array",
		"BigInt64Array", "BigUint64Array",
		"ArrayBuffer", "SharedArrayBuffer", "DataView",
		"Atomics",
		"IObject", "IArrayLike", "IThenable", "IIterable",
		"IterableIterator", "Iterator", "Iterable",
		"Generator", "AsyncIterator", "AsyncIterable", "AsyncGenerator",
		"Disposable", "IDisposable",
		"ReadonlyArray", "ReadonlyMap", "ReadonlySet",
		"Arguments", "Arguments_",
		"Window", "Document", "Element", "Event",
		"EventTarget", "EventListener",
		"Node", "NodeList", "HTMLElement",
		"XMLHttpRequest":
		return true
	}
	if defaultExternNames()[name] {
		return true
	}
	if scope == nil {
		return false
	}
	// Dotted name: check namespace constructor props, namespace typedefs,
	// and namespace enums.
	if dot := strings.IndexByte(name, '.'); dot >= 0 {
		// Check dotted typedef alias (e.g. "goog.MyTypedef").
		if scope.LookupTypeAlias(name) != nil {
			return true
		}
		// Walk to the last segment for constructor / enum prop lookup.
		// For multi-segment paths like `ns.ns2.Foo`, split at last dot.
		lastDot := strings.LastIndexByte(name, '.')
		owner := name[:lastDot]
		prop := name[lastDot+1:]
		// Only count it as a type if the prop is a constructor.
		if t := scope.LookupNamespaceProp(owner, prop); t != nil {
			if fn, ok := t.(*jstype.FunctionType); ok && fn.Instance != nil {
				return true
			}
			// Also accept enum types (e.g. `ns.MyEnum = {@enum}…`).
			if _, ok := t.(*jstype.EnumType); ok {
				return true
			}
		}
		return false
	}
	// Typedef alias?
	if scope.LookupTypeAlias(name) != nil {
		return true
	}
	// Enum type?
	if scope.LookupEnum(name) != nil {
		return true
	}
	// Constructor / interface?
	if scope.LookupNominalType(name) != nil {
		return true
	}
	return false
}

// positionInSrc returns (line, col) for absolute offset off in src.
// Line is 1-indexed, col 0-indexed, matching Java Node positions.
// Go-only: language-adapter — positionInSrc is a Go-only helper that computes line/column position in source; Java uses SourceFile.getLineOfOffset / getColumnOfOffset for diagnostics.
func positionInSrc(src string, off int) (int, int) {
	if off > len(src) {
		off = len(src)
	}
	line := 1
	col := 0
	for i := 0; i < off; i++ {
		if src[i] == '\n' {
			line++
			col = 0
			continue
		}
		col++
	}
	return line, col
}

// ResolveTypesInScope is the post-scope-creation named-type validation
// pass. It re-scans JSDoc annotations in the source and, for any type
// name that:
//   (a) appears somewhere in source as a declared identifier (allNames)
//       — so CheckJsDocTypes already suppressed it via its allNames gate
//   (b) does NOT resolve as a genuine type (constructor, typedef, enum,
//       or primitive)
//
// emits JSC_UNRECOGNIZED_TYPE_ERROR, matching Java's TypeValidator
// behaviour where value bindings used as type names are flagged even
// when they appear as declared identifiers.
//
// Names NOT in allNames are already caught by CheckJsDocTypes and are
// not duplicated here.
// Go-only: language-adapter — ResolveTypesInScope is a Go-only helper that post-processes a scope to resolve forward type references; Java resolves types eagerly during scope construction via TypedScopeCreator.createScope.
func (c *TypedScopeCreator) ResolveTypesInScope(scope *TypedScope) []diag.Diagnostic {
	if c.sourcePath == "" || c.sourceSrc == "" {
		return nil
	}
	path := c.sourcePath
	src := c.sourceSrc

	allTemplates := rtsCollectAllTemplateNames(src)
	allNames := rtsCollectDeclaredNames(src)

	// Scan all JSDoc blocks and register pending types for names that
	// are in allNames but do not resolve as genuine types in scope.
	s := src
	i := 0
	for i < len(s) {
		start := strings.Index(s[i:], "/**")
		if start < 0 {
			break
		}
		start += i
		end := strings.Index(s[start:], "*/")
		if end < 0 {
			break
		}
		end += start
		rtsCheckComment(path, src, s, start, end, allTemplates, allNames, scope, c.registry)
		i = end + 2
	}

	// Retry any pending types now that the full scope is available.
	c.registry.ResolveWorklist(scope.ResolveName)

	// Drain remaining unresolved entries and emit diagnostics.
	var out []diag.Diagnostic
	// Emit scope-creation diagnostics (e.g. JSC_CONFLICTING_GETTER_SETTER_TYPE).
	out = append(out, c.diags...)
	c.diags = nil
	for _, nt := range c.registry.DrainPendingNamedTypes() {
		out = append(out, diag.Diagnostic{
			Path:    nt.SourceName(),
			Line:    nt.Lineno(),
			Col:     nt.Charno(),
			Code:    diag.JSC_UNRECOGNIZED_TYPE_ERROR,
			Level:   diag.Error,
			Message: fmt.Sprintf("Bad type annotation. Unknown type %s", nt.Reference()),
		})
	}
	return out
}

// rtsCollectDeclaredNames mirrors CheckJsDocTypes.collectDeclaredNames:
// scans source for every identifier (including JSDoc-comment text is
// skipped). Used to determine which names CheckJsDocTypes would have
// suppressed via its allNames gate.
// Go-only: language-adapter — rtsCollectDeclaredNames is a Go-only helper that collects declared names from JSDoc comments; Go-specific post-processing step not needed when types resolve eagerly.
func rtsCollectDeclaredNames(s string) map[string]bool {
	out := map[string]bool{}
	i := 0
	for i < len(s) {
		// Skip JSDoc / block comments.
		if i+2 < len(s) && s[i] == '/' && s[i+1] == '*' && s[i+2] == '*' {
			if n := strings.Index(s[i+3:], "*/"); n >= 0 {
				i += n + 5
				continue
			}
			return out
		}
		if i+1 < len(s) && s[i] == '/' && s[i+1] == '*' {
			if n := strings.Index(s[i+2:], "*/"); n >= 0 {
				i += n + 4
			} else {
				return out
			}
			continue
		}
		if i+1 < len(s) && s[i] == '/' && s[i+1] == '/' {
			if n := strings.Index(s[i:], "\n"); n >= 0 {
				i += n + 1
			} else {
				return out
			}
			continue
		}
		if s[i] == '"' || s[i] == '\'' {
			i = rtsSkipQuote(s, i)
			continue
		}
		if isIdentStart(s[i]) {
			start := i
			for i < len(s) && isIdentChar(s[i]) {
				i++
			}
			out[s[start:i]] = true
			continue
		}
		i++
	}
	return out
}

// Go-only: language-adapter — rtsSkipQuote is a Go-only helper that skips quoted strings in JSDoc text; Go-specific helper for rtsCheckComment string scanning.
func rtsSkipQuote(s string, i int) int {
	quote := s[i]
	i++
	for i < len(s) {
		if s[i] == '\\' {
			i += 2
			continue
		}
		if s[i] == quote {
			return i + 1
		}
		i++
	}
	return len(s)
}

// rtsCollectAllTemplateNames collects every @template name across all
// JSDoc blocks in src. Mirrors CheckJsDocTypes.collectAllTemplateNames.
// Go-only: language-adapter — rtsCollectAllTemplateNames is a Go-only helper that collects all @template names in a source file; Go-specific analysis for forward-reference resolution.
func rtsCollectAllTemplateNames(src string) map[string]bool {
	out := map[string]bool{}
	i := 0
	for i < len(src) {
		start := strings.Index(src[i:], "/**")
		if start < 0 {
			return out
		}
		start += i
		end := strings.Index(src[start:], "*/")
		if end < 0 {
			return out
		}
		end += start
		for k, v := range rtsCollectTemplateNames(src[start:end]) {
			out[k] = v
		}
		i = end + 2
	}
	return out
}

// Go-only: language-adapter — rtsCollectTemplateNames is a Go-only helper that collects @template names from a JSDoc comment body; Go-specific helper for ResolveTypesInScope.
func rtsCollectTemplateNames(body string) map[string]bool {
	out := map[string]bool{}
	i := 0
	for i < len(body) {
		idx := strings.Index(body[i:], "@template")
		if idx < 0 {
			return out
		}
		idx += i
		nl := strings.Index(body[idx:], "\n")
		if nl < 0 {
			nl = len(body) - idx
		}
		rest := body[idx+len("@template") : idx+nl]
		// Strip bounded-generic `{B}` clause.
		if lb := strings.Index(rest, "{"); lb >= 0 {
			if rb := strings.Index(rest[lb:], "}"); rb > 0 {
				rest = rest[:lb] + rest[lb+rb+1:]
			}
		}
		for _, part := range strings.Split(rest, ",") {
			part = strings.TrimSpace(part)
			if sp := strings.IndexAny(part, " \t*"); sp >= 0 {
				part = part[:sp]
			}
			if part != "" && isIdentStart(part[0]) {
				out[part] = true
			}
		}
		i = idx + nl
	}
	return out
}

// rtsCheckComment scans one JSDoc comment block and registers pending
// NamedType entries for names that are in allNames but do not resolve
// strictly as types in scope.
// Go-only: language-adapter — rtsCheckComment is a Go-only helper that checks a JSDoc comment block for unresolved types; Go-specific post-pass not present in Java TypedScopeCreator.
func rtsCheckComment(path, src, s string, start, end int, allTemplates, allNames map[string]bool, scope *TypedScope, reg *jstype.Registry) {
	body := s[start:end]
	tmpls := rtsCollectTemplateNames(body)
	i := 0
	for i < len(body) {
		at := strings.Index(body[i:], "@")
		if at < 0 {
			return
		}
		at += i
		j := at + 1
		for j < len(body) && isTagChar(body[j]) {
			j++
		}
		tag := body[at+1 : j]
		if !isTypeTag(tag) {
			i = j
			continue
		}
		k := j
		for k < len(body) && (body[k] == ' ' || body[k] == '\t' || body[k] == '\n' || body[k] == '\r' || body[k] == '*') {
			k++
		}
		if k >= len(body) || body[k] != '{' {
			i = j
			continue
		}
		braceStart := k
		braceEnd := rtsFindMatchingBrace(body, braceStart)
		if braceEnd < 0 {
			return
		}
		inner := body[braceStart+1 : braceEnd]
		rtsCheckTypeExpr(path, src, start+braceStart+1, inner, tmpls, allTemplates, allNames, scope, reg)
		i = braceEnd + 1
	}
}

// rtsCheckTypeExpr walks a type expression string and registers any
// name that is in allNames (declared identifier in source) but fails
// the strict type-only check. Names NOT in allNames are handled by
// CheckJsDocTypes already.
// Go-only: language-adapter — rtsCheckTypeExpr is a Go-only helper that checks a single type expression for resolution; Go-specific helper invoked by rtsCheckComment.
func rtsCheckTypeExpr(path, src string, base int, inner string, tmpls, allTemplates, allNames map[string]bool, scope *TypedScope, reg *jstype.Registry) {
	type frame struct {
		kind       byte
		expectType bool
	}
	stack := []frame{{kind: 0, expectType: true}}
	top := func() *frame { return &stack[len(stack)-1] }

	i := 0
	for i < len(inner) {
		ch := inner[i]
		switch {
		case ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' || ch == '*':
			i++
		case ch == '!' || ch == '?':
			i++
		case ch == '|':
			top().expectType = true
			i++
		case ch == ',':
			switch top().kind {
			case '{', '(':
				top().expectType = false
			default:
				top().expectType = true
			}
			i++
		case ch == ':':
			top().expectType = true
			i++
		case ch == '=':
			i++
		case ch == '{':
			stack = append(stack, frame{kind: '{', expectType: false})
			i++
		case ch == '(':
			stack = append(stack, frame{kind: '(', expectType: false})
			i++
		case ch == '<':
			stack = append(stack, frame{kind: '<', expectType: true})
			i++
		case ch == '}' || ch == ')' || ch == '>':
			if len(stack) > 1 {
				stack = stack[:len(stack)-1]
			}
			i++
		case ch == '.' && i+2 < len(inner) && inner[i+1] == '.' && inner[i+2] == '.':
			i += 3
		case isIdentStart(ch):
			start := i
			for i < len(inner) && (isIdentChar(inner[i]) || inner[i] == '.') {
				i++
			}
			name := inner[start:i]
			if !top().expectType {
				continue
			}
			if name == "typeof" {
				// Skip the operand — not a type name.
				k := i
				for k < len(inner) && (inner[k] == ' ' || inner[k] == '\t') {
					k++
				}
				for k < len(inner) && (isIdentChar(inner[k]) || inner[k] == '.') {
					k++
				}
				i = k
				continue
			}
			// Skip template names declared in this comment or globally.
			if tmpls[name] || allTemplates[name] {
				continue
			}
			// Only register as pending for FLAT (non-dotted) names.
			// Dotted names like `ns.Foo` have too many edge cases
			// (namespace constructors inside IIFEs, aliased namespaces)
			// that the global scope cannot reliably distinguish from
			// genuine types. Only flat names that are top-level plain
			// value bindings (scope declares them as Unknown) are safe.
			if strings.IndexByte(name, '.') >= 0 {
				continue
			}
			// Only register as pending when the name:
			//   1. Appears somewhere in source as a declared identifier
			//      (allNames gate) — otherwise CheckJsDocTypes already
			//      handles it (no allNames suppression happened).
			//   2. Does NOT resolve strictly as a type.
			//   3. IS explicitly declared in global scope as Unknown
			//      (not nil — nil means locally scoped or undeclared,
			//      which CheckJsDocTypes would NOT suppress).
			if allNames[name] && !typeNameResolvesInScope(scope, name) &&
				scope.GetRawDeclaredType(name) == jstype.Unknown {
				line, col := positionInSrc(src, base+start)
				reg.RegisterPending(name, path, line, col)
			}
		default:
			i++
		}
	}
}

// getLValueRootScope mirrors TypedScopeCreator.getLValueRootScope in Java.
// Returns the scope in which the lvalue declaration belongs (VAR → hoist,
// LET/CONST → block). In our flat-scope model this always returns scope.
// Go-only: language-adapter — getLValueRootScope is a Go-only helper that returns the root scope for an LValue; Java uses TypedScopeCreator.getDeclaredType inline without a separate helper.
func getLValueRootScope(_ *rhino.Node, scope *TypedScope) *TypedScope { return scope }

// inferTypeForDestructuredTarget mirrors the narrow path in Java's
// TypedScopeCreator: for `const {id} = obj`, returns the declared type of
// obj.id from patternType when available, nil otherwise (defers to inference).
// Go-only: language-adapter — inferTypeForDestructuredTarget is a Go-only helper that infers type for a destructuring target property; Java uses NormalScopeBuilder.RValueInfo inline.
func inferTypeForDestructuredTarget(patternType jstype.JSType, propertyName string) jstype.JSType {
	if patternType == nil || isJSUnknown(patternType) || propertyName == "" {
		return nil
	}
	var obj *jstype.ObjectType
	switch t := patternType.(type) {
	case *jstype.ObjectType:
		obj = t
	case *jstype.TemplatizedType:
		obj = t.Base
	}
	if obj != nil {
		return obj.GetPropertyType(propertyName)
	}
	return nil
}

// rtsFindMatchingBrace returns the offset of the closing `}` matching
// the `{` at open within s, or -1.
// Go-only: language-adapter — rtsFindMatchingBrace is a Go-only helper that finds the matching closing brace in source; Go-specific helper for rtsCheckComment brace scanning.
func rtsFindMatchingBrace(s string, open int) int {
	depth := 0
	for i := open; i < len(s); i++ {
		switch s[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}
