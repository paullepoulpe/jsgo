# A typed view of the AST

Status: prototype. Package `internal/ast`, generated from
`internal/ast/grammar.txt`; twelve passes use it (`MinimizeExitPoints`,
`ConvertToDottedProperties` with a typed visitor, the parts of the hot
`InlineVariables` that read tree shapes, `PeepholeMinimizeConditions`,
`PeepholeRemoveDeadCode`, `Normalize`, `PeepholeFoldConstants`,
`PeepholeSubstituteAlternateSyntax`, `Denormalize`,
`PeepholeCollectPropertyAssignments`, `CollapseVariableDeclarations` and
`ExploitAssigns`).

## Why

The AST is a tree of `rhino.Node`s. Every node has a token, children and
generic accessors, so what a child means is a convention the reader has to
know: the second child of an `IF` is its then-block, the second child of a
`TRY` is a `BLOCK` that holds its catch clause (or nothing), an anonymous
function has a `NAME` child with an empty string while an anonymous class has
an `EMPTY` child. Code reads like `n.GetSecondChild().GetNext()`, and nothing
stops a pass from treating a statement as an expression.

The goal is code that a newcomer can read and that the compiler checks: one
Go type per node kind, with methods that name the children.

## Design

**Views, not a new tree.** `ast.If` is `struct{ node rhino.Node }`: 4 bytes,
created and passed for free, and its methods inline to the same rhino calls
as before. The storage does not change, so passes can move to the typed API
one at a time, and `x.Node()` goes back to the rhino API for anything the
typed view does not cover.

**Kinds and categories.** A kind is one shape of node (`If`, `Call`,
`Binary` for all binary operators). A category is the set of kinds allowed
in one place (`Expr`, `Statement`, `Target`, `Arg`...). Go has no sum types,
and an interface would allocate when a 4-byte value is stored in it, so
categories are views too, with checked conversions:

```go
ifs, ok := ast.AsIf(n)    // ok is false if n is not an if (or is the nil node)
s := ifs.AsStatement()    // a kind to its category: always works
if call, ok := e.AsCall(); ok { ... }  // a category to a kind: checked
```

Conversions exist only where they make sense: `Expr` has no `AsBlock`.

**Optional children use comma-ok**, not a nil check:

```go
if els, ok := ifs.Else(); ok { ... }   // els is an ast.Block
name, ok := fn.Name()                  // false for an anonymous function
```

Children that the tree stores as an `EMPTY` placeholder (`for (;;)`, a class
without `extends`) read the same way, so the reader never needs to know
which encoding a kind uses.

**Lists** are iterators plus first/last/append:

```go
for s := range block.Statements() { ... }
last, ok := block.LastStatement()
```

The iterator reads the next element after the loop body runs, like the
classic `for c := n.GetFirstChild(); c != nil; c = c.GetNext()` loop, so
ported loops keep their meaning when the body edits the tree. A loop whose
body detaches or moves the current element uses `ast.ReadAhead`, which
reads the next element first:

```go
for member := range ast.ReadAhead(obj.FirstMember()) {
	props = append(props, member.Detach())
}
```

**Generated from a grammar.** `grammar.txt` lists every kind in one line:

```
If                 IF                    cond:Expr then:Block else:Block?
Try                TRY                   body:Block _catches:TryCatches finally:Block?
Call               CALL                  callee:Expr args:Arg*
GetProp            GETPROP               object:Expr $Property
```

`tools/astgen` turns it into `nodes.gen.go` (`go generate ./internal/ast`):
the types, `AsX`, accessors and setters per role, constructors, category
conversions, accessors shared by all kinds of a category (`SwitchCase.Body()`,
`Jump.Label()`), and the grammar as data. The few shapes the grammar cannot
express are hand-written in `custom.go` (`Try.Catch()`, `Function.Name()`,
`Loop.Body()`). A test fails if `nodes.gen.go` is out of date.

**The grammar is checked.** `ast.Check` walks a tree and reports every node
that does not follow the grammar. In tests, the `AstValidator` (which every
`CompilerTestCase` runs after each pass) also calls it, so a grammar that
drifts from the trees the compiler builds fails the tests. A normal compile
does not (it would add one more walk, about 20 ms on bundle, to the final
validity check); `JSGO_AST_GRAMMAR_CHECK=1` turns it on there.

The types also catch grammar mistakes at compile time: changing `If`'s
`then` to `Expr` in the grammar breaks the build of `MinimizeExitPoints`.

**Typed traversal.** `ast.Visitor` has an optional handler per kind;
`NodeTraversal_TraverseVisitor` runs it as a post-order callback:

```go
NodeTraversal_TraverseVisitor(c.compiler, root, &ast.Visitor[*NodeTraversal]{
	ComputedProp: c.visitComputedProp,   // func(t *NodeTraversal, x ast.ComputedProp)
	GetElem: func(_ *NodeTraversal, x ast.GetElem) {
		c.convertGetElem(x.Node(), x.Object(), x.Index(), rhino.IR_Getprop)
	},
	...
})
```

A pass lists the kinds it handles instead of switching on tokens in one
`Visit`, and each handler gets its node already typed.

## How the grammar was validated

The rules were transcribed from `AstValidator`, then checked against real
trees:

- the whole jscomp test suite: one violation, in
  `NormalizeTest.testLogicalAssignInLabeledLoop`, a known upstream bug the
  test expects (Normalize builds a `LABEL` with three children there; the
  validator throws first);
- the typed-app, bundle and typed-full workloads, checked before and after
  compilation (`--jscomp_dev_mode=START_AND_END`): no violation;
- a deliberately wrong rule (`then:Expr`) was reported 469 times on
  typed-app, so the check does run.

## The pilot: MinimizeExitPoints

The pass moves the statements that follow an `if` whose branch ends in an
exit into the other branch. It is all about statement shapes, so it shows the
difference well.

Before:

```go
if n.IsTry() {
	tryBlock := n.GetFirstChild()
	m.TryMinimizeExits(tryBlock, exitType, labelName)
	allCatchNodes := NodeUtil_GetCatchBlock(n)
	if NodeUtil_HasCatchHandler(allCatchNodes) {
		jutil.CheckState(allCatchNodes.HasOneChild())
		catchNode := allCatchNodes.GetFirstChild()
		catchCodeBlock := catchNode.GetLastChild()
		m.TryMinimizeExits(catchCodeBlock, exitType, labelName)
	}
}
```

After:

```go
if try, ok := n.AsTry(); ok {
	m.TryMinimizeExits(try.Body().AsStatement(), exitType, labelName)
	if catch, ok := try.Catch(); ok {
		m.TryMinimizeExits(catch.Body().AsStatement(), exitType, labelName)
	}
	return
}
```

`TryMinimizeIfBlockExits` took a source branch, a destination branch that
could be missing, an `EMPTY` node, a block or a single statement, and handled
all four. With typed branches it takes the `if` and which branch to look at:
both branches are `ast.Block`s (the grammar check confirms it for every tree
above), so three of the four cases were dead code, and they are gone.

Results:

- tests: `MinimizeExitPointsTest` (85 cases) and the whole suite pass (the
  same 5 golden tests fail as on main);
- output: byte-identical on typed-app, checks-only, bundle and typed-full;
- speed: no measurable difference (5 interleaved runs, medians): bundle
  wall +0.4%, typed-full +0.7%, CPU +0.3% and 0.0%, within the run-to-run
  noise; this includes the extra grammar walk of the final validity check;
- the typed accessors inline to the rhino calls they replace (`go build
  -gcflags=-m`), except `Block.LastStatement` and the hand-written
  `Try.Catch` and `Loop.Body`, which are just over Go's inlining budget.

## Second pilot: ConvertToDottedProperties with a typed visitor

The pass handled seven kinds in one `Visit` with a token switch; it now
registers one handler per kind (or a shared helper for kinds with the same
shape). Results: the 21 test methods pass, output is
byte-identical on all four workloads, and the pass takes the same time
(about 60 ms on bundle).

One trap showed up while porting: accessors are positional underneath, so
`x.Value()` (the second child) is wrong once `x.Key()` (the first) has been
detached. Read the children you need before editing the node: the same rule
as with the rhino API, but less visible when the positions are hidden.

## Third port: InlineVariables

A hot pass (about 1.1 s on bundle, 1.8 s on typed-full) of 1,700 lines,
most of it reference bookkeeping that never touches the tree. The eight
places that read shapes around a reference now use the typed view:

```go
// before
callNode := n.GetParent()
if NodeUtil_IsNormalOrOptChainCall(callNode) {
	calleeNode := callNode.GetFirstChild()
	if NodeUtil_IsNormalOrOptChainGetProp(calleeNode) {
		if calleeNode.GetString() == "apply" {
			thisArgNode := calleeNode.GetNext()
			...
// after
call, ok := ast.AsNormalOrOptChainCall(n.GetParent())
if !ok { return false }
callee, ok := call.Callee().AsNormalOrOptChainGetProp()
if !ok || callee.Property() != "apply" { return false }
thisArg, ok := call.FirstArg()
if !ok { return false }
secondArg, ok := thisArg.Next()
return ok && secondArg.Is(n)
```

Output and source maps byte-identical, same pass time (medians of 11-19
interleaved runs within 1.2%). It added to the API: categories for
NodeUtil's `isNormalOrOptChain*` predicates (`NormalOrOptChainCall`,
`...GetProp`, `...Get`), list accessors shared by a category's kinds
(`call.Args()` on both calls), `Function.RemoveName()`, and `x.Is(n)` on
every view.

Its feedback on the API:

- Checked conversions do not chain: "is the declared value a name?" is
  three comma-ok steps; small helpers per pass are the answer today.
- Code that looks up from a node ("is n the callee of its parent?") reads
  `AsCall(n.GetParent())` then compares; parent-role queries would help.
- A switch that groups kinds by "the first child is the declared name"
  (VAR, LET, CONST, FUNCTION, CLASS) is not clearer typed: each kind names
  that child differently.
- `fn, _ := ast.AsFunction(x)` after a NodeUtil check is safe but a bad
  pattern to copy; typed versions of the common NodeUtil predicates, which
  return the view, would avoid it.
- Some typed helpers exceed Go's inlining budget (each accessor brings in
  the record lookup); fine per variable, to watch in per-node loops.

## Fourth to sixth ports: PeepholeMinimizeConditions, PeepholeRemoveDeadCode, Normalize

Three passes that are mostly about shapes, ported in parallel (2,400 lines
together). Output, source maps and stderr byte-identical; pass times equal
within noise (the machine was loaded, so medians of interleaved runs, and
for Normalize an in-process harness: 238.8 vs 236.1 ms median over 48 runs).

What the typed view changed in them:

- **Reading by role.** `assign.Target()` and `assign.Value()` instead of
  two `RemoveFirstChild()` calls in a row; `try.Finally()`, `try.Catch()`
  instead of `HasXChildren(3)` and `GetSecondChild().GetOnlyChild()`;
  `ifs.Else()` instead of `GetChildCount() < 3`.
- **Java check-then-get pairs merged.** `isFoldableExpressBlock` +
  `getBlockExpression` (and three more pairs) are one comma-ok function that
  returns the typed node: no shape test, re-test in a `checkState`, then
  positional reads.
- **Token switches as dispatch tables.** Normalize had two copies of a token
  switch (fast renamer and `MakeDeclaredNamesUnique` paths); both now run one
  pre-order and one post-order `ast.Visitor`, called by hand from the
  callback that keeps `EnterScope`/`ExitScope` (a pattern worth reusing when
  a pass needs pre-order or scope callbacks).
- **Dead checks gone, missing checks visible.** `IsBlock` tests of if
  branches, `checkNotNull` of the then branch, `HasTwoChildren` on typed
  arguments. In the other direction, facts the types cannot know (the body
  is a block after normalization) became explicit `CheckState`s.
- **A crash found.** `tryFoldLabel` read the label of a `break` without
  checking it had one (Java does the same): `switch(x){case 1: a: break;}`
  crashed the compiler. `brk.Label()` is comma-ok, so the case is now handled
  (Go-only fix, with a test).
- **Java quirks made visible.** The switch folding passes each case's first
  child as its tag; for a default case that is its body (kept, for identical
  output, behind a named helper). A Java "one case left" test reads the
  wrong sibling and is always true (now a real check, same result).

Their feedback, and what changed in the API since:

| Friction | Now |
|---|---|
| no "only element" accessor (`HasOneChild()` + `GetFirstChild()`, six uses in one pass) | `OnlyX()` on every list (`block.OnlyStatement()`) |
| categories have no `Token()`: switching on a statement's kind needs `.Node()` | `Token()` on every category |
| every edit goes through `.Node()` | `Detach()` on every view, `ReplaceWith(v)` on categories; the rest (`InsertBefore`, `Srcref`, the `IR_*` builders) is still rhino |
| `AsExpr` (a switch over 70 tokens) and `Statement.Next` (31) over the inlining budget | categories test a 256-bit token set: `AsExpr`, `AsStatement` inline; `Next()` calls an inlined check |
| setters are positional, so they hit the wrong child mid-edit | documented in the package doc: read every child before the first edit |
| peephole passes get one node at a time and dispatch with `x, _ := ast.AsHook(n) // the case proves the kind` | `x := ast.MustHook(n)` in the case of a token switch (panics, naming both kinds, if wrong; inlines) |
| a visitor handler does not get the parent; no handler per category | open |
| a typed visitor costs more than a token switch in cheap passes | the traversal skips the kinds a visitor has no handler for: now cheaper than a token switch |
| optional children stored as an EMPTY placeholder have no `RemoveX` (putting the EMPTY back, with its srcref, stays rhino) | open |
| RETURN and THROW share no category; neither do the optional chains | left on rhino where typing would mean two conversions for one `GetFirstChild` |

## Seventh to twelfth ports: constant folding, alternate syntax, denormalization, and three small passes

`PeepholeFoldConstants` (2,200 lines), `PeepholeSubstituteAlternateSyntax`,
`Denormalize`, `PeepholeCollectPropertyAssignments`,
`CollapseVariableDeclarations` and `ExploitAssigns`, ported in parallel with
the API above. Output, source maps and stderr byte-identical; pass times
equal within noise, checked with in-process harnesses (the tracer was too
noisy on a loaded machine): for example `PeepholeSubstituteAlternateSyntax`
45.96 vs 46.28 ms and `Denormalize` 26.59 vs 26.79 ms (medians of 30 runs).

- **Operands by role.** Constant folding reads `op.Left()`/`op.Right()`,
  `unary.Operand()`, `assign.Target()`/`Value()`; "the right side has two
  children" became the `AsBinary` conversion. A new category,
  `NormalOrOptChainGetElem`, replaced positional `(left, right)` passing.
- **Typed parameters** (`TryFoldArrayAccess(ast.ArrayLit, ...)`,
  `tryMinimizeIf`-style methods taking `ast.Call`, `ast.Return`...) removed
  the Java re-checks (`checkState(isString)`, `isObjectLit`).
- **Another crash found:** `Boolean(...a)` was folded to `!!...a`, which
  `IR.not` rejects. `call.OnlyArg()` returns an `Arg`, and its conversion to
  `Expr` (comma-ok) made the spread case visible. Now left alone (Go-only,
  with a test).
- **A miscompilation found, not fixed** (same output as Java):
  `tryMinimizeWindowRefs` marks the enclosing call as a free call whatever
  the role of the `window.x` reference, so `obj.f(window.Array)` becomes
  `(0,obj.f)(Array)`. Documented in a Go-only note.
- **Two more categories**: `NormalGet` (`NodeUtil.isNormalGet`) and, from
  Normalize, `EnhancedFor`.

What they reported:

- **Dispatch cost matters for cheap passes.** `ast.Visitor` made
  Denormalize's traversal ~10% slower and a one-handler visitor in
  CollapseVariableDeclarations ~20% slower than a token switch, because every
  node goes through the dispatch. Both kept a token switch that hands typed
  nodes to the rewrites. A visitor that declares its kinds (next steps, 1)
  would fix it.
- **The per-node dispatch of peephole passes is still the main gap**:
  OptimizeSubtree returns the replacement node, so `ast.Visitor` does not
  fit, and the `x, _ := ast.AsX(n) // the case proves the kind` lines remain.
- **One token within a kind takes two steps** (`void`, `!`, `var`, `||`):
  `u, ok := e.AsUnary(); ok && u.Token() == rhino.Token_VOID`.
  Helpers such as `AsUnaryOp(n, tok)` were asked for.
- **Missing pieces**: constructors for kinds with a string payload
  (`NewStringKey`), raw/cooked accessors on template strings, an index
  accessor and element count on lists, `Target.AsExpr` (patterns are
  targets, not expressions), a common view of SCRIPT items and BLOCK
  statements.

## The API after this feedback

Added before the next ports, each for a pattern the ports kept writing:

| ask | added |
|---|---|
| `x, _ := ast.AsX(n) // the case proves the kind` in token switches | `ast.MustX(n)`: panics if n is not an X; inlines to the same check |
| one token within a kind (`void`, `!`, `var`) | `ast.AsUnaryOf(n, rhino.Token_VOID)` and the same on categories (`e.AsBinaryOf(tok)`) |
| constructors for kinds with a string payload | generated: `NewStringKey(key, value)`, `NewMemberFunctionDef(key, fn)`, `NewOptChainGetProp(obj, prop)`, ... |
| template strings | `TemplateLitString.Raw()`, `Cooked()` (ok is false for an invalid escape) |
| "is this the last argument / statement" | `IsLast()` on list elements |
| "how many arguments" (Java `getChildCount() - 1`) | `Num<Role>()` on every list: `call.NumArgs()`, `block.NumStatements()` |
| loops that move the current element | `for x := range ast.ReadAhead(list.FirstX())` reads the next element first |
| a visitor's dispatch cost | the traversal skips the kinds a visitor has no handler for (next steps, 1) |

Not done: an index accessor on lists (no port needed one after `IsLast`),
`Target.AsExpr`, a common view of SCRIPT and BLOCK statements, and a shared
view of `Assign` and `DestructuringLhs` (their roles differ: `target:Target`
against `pattern:Pattern`, and a declaration's value is optional).

## Thirteenth to eighteenth ports: known methods and five ES6 passes

`PeepholeReplaceKnownMethods`, then `Es6RewriteArrowFunction`,
`RewriteObjectSpread`, `LateEs6ToEs3Converter`, `CreateSyntheticBlocks` and
`Es6RewriteDestructuring`,
ported on the API above. Output, source maps and stderr byte-identical on
the four gate workloads; pass times equal within noise or better, measured
with in-process A/B harnesses (old and new pass on clones of the same tree,
alternating).

- **Calls read by role.** ReplaceKnownMethods reads the callee, receiver,
  method name and arguments as `call.Callee()`, `getProp.Object()`,
  `getProp.Property()`, `call.FirstArg()`/`OnlyArg()`/`arg.Next()` instead
  of `GetNext` chains and nil checks. Each `tryFoldX` takes the kinds it
  folds (`ast.Call`, `ast.StringLit`, `ast.Arg`), so the Java re-checks
  (`checkArgument(isCall)`, `!arg1.IsNil()`) are gone, and so are their
  per-call `Check*Msg` costs: 0 to 17% faster on the 85K calls of the
  bundle inputs.
- **Two output quirks found, kept for identical output** (Go-only notes):
  `"a,b".split(",", 1, f())` folds to `["a"]` and drops the call to `f`, and
  `Array.join` counts the cost of a separator it already removed. Same as
  Java.
- **Shapes stated instead of assumed.** A computed property in an object
  literal has a value, and the members left after computed properties,
  getters and setters are string keys: two explicit checks in
  LateEs6ToEs3Converter where Java crashed later on a null.
- **Patterns by kind.** The destructuring rewrites take the kind they work
  on (`ast.ObjectPattern`, `ast.ArrayPattern`, `ast.Catch`,
  `ast.EnhancedFor`), so the Java re-checks of the pattern and its parent are
  types now.
- **A crash fixed (Go-only, with a test):** `({a, ...arr[0]} = foo())` read
  the string of the `GETELEM` rest target. Java throws a
  `ClassCastException` on the same input.

What they reported:

- **The AstFactory returns `rhino.Node`.** The transpilation passes build
  most of their output with it (it sets colors), so the typed view stops at
  every node they build. Typed factory methods would help these passes most.
- **Node flags have no typed accessors** (JSDoc, quoted keys, the computed
  getter flag); `Function.SetArrow` was added, the others still need
  `.Node()`.
- **"Read the next element first" loops**, three times in four passes: now
  `ast.ReadAhead`.
- **No equality between two views**: `x.Is(y.Node())`.
- **Read both children before detaching either.** `DefaultValue.Target()`
  and `Value()`, `ComputedProp.Key()` and `Value()` are positional: after
  detaching the first, the second accessor reads the wrong child. Five places
  in the destructuring pass; Java's `removeFirstChild()` twice hides it.
- **No shared view for two kinds with the same layout**: an `Assign` and a
  `DestructuringLhs` (target first, value last) are passed as `rhino.Node`.
  A category with `Target()` and `Value()` would type them.

## Nineteenth to twenty-second ports: three checks and the module metadata

`CheckSuspiciousCode`, `StrictModeCheck`, `ProcessClosurePrimitives` and
`GatherModuleMetadata`, ported on the API above. They report warnings, so
the stderr comparison on checks-only and typed-app is their main check:
byte-identical, and output and source maps too.

- **One dispatch instead of seven token tests.** CheckSuspiciousCode ran
  seven checks on every node, each testing the token itself. Its Visit now
  switches once and hands each check the kind it reads: the `if` (then and
  else blocks), the body of a loop (`ast.Loop.Body()`), the operands of a
  comparison, of `in` and `instanceof`, of `&&` and `||`. A node has one
  token, so the reports come out in the same order. The pass got faster:
  median 83 vs 100 ms on checks-only (5 interleaved runs, loaded machine).
- **Report targets passed, not read back.** `reportIfNaN` and
  `reportIfNonObject` reported on `operand.getParent()`; they now take the
  operator they report on.
- **Closure calls read by role.** ProcessClosurePrimitives reads
  `goog.x(...)` as a GetProp callee whose object is the name `goog`, and
  `goog.base(this, 'method', ...)` as the first two arguments. Five
  "qualified name or `null`" fallbacks in messages became one helper.
- **One idiom named.** GatherModuleMetadata tested "exactly one argument, a
  string literal" ten times as `n.hasTwoChildren() &&
  n.getLastChild().isStringLit()` (the callee is the first of the two
  children). It is now `onlyStringArg(call)`, built on `call.OnlyArg()`.
- **Kept on rhino**: code that reads a node's parent or siblings of any kind
  (is this name a declaration, the enclosing class, the statements after a
  constructor), and the duplicate-member scan of StrictModeCheck, which
  walks object literal members and class members, two different lists.

## Twenty-third to twenty-fifth ports: ProcessDefines, ControlFlowAnalysis, SyntacticScopeCreator

Two of the compiler's core pieces: the control-flow graph that type
inference, the data-flow analyses and several optimizations use, and the
scope creator that almost every pass uses. Both were checked beyond the
gate, by comparing old and new code in process on the same trees:

- **ControlFlowAnalysis**: identical graphs, priorities and finally maps
  for all 81,110 function and script roots of the four workloads, and the
  exported helpers (follow node, exception handler, fall-through,
  break/continue targets) agree on all 5.59 M nodes. The golden test passes
  all 131 files against the Java dump when pointed at the extracted corpus
  (`CFA_CORPUS`; it fails in this environment only because
  `.cache/corpus` is empty). CFG time: new/old median 0.98-1.02, inside the
  ±1 % band of an A/A control.
- **SyntacticScopeCreator**: the same vars in the same order (name, node,
  index, input, flags), the same redeclaration events and block-declaration
  marks for every scope root, on the four workloads plus 797 ES-module
  files (2,750 imports, 393 object patterns). Scope creation: new/old
  median 0.95-1.02.

What got clearer:

- **Handlers name what they connect**: `ifs.Else()`, `while.Cond()`,
  `do.Body()`, `sw.Discriminant()`, `ret.Value()`, `brk.Label()`; a `for`
  without a condition is the `ok` of `For.Cond()`, not an `IsEmpty` test.
  The follow node of a TRY child was three positional branches, each
  repeating the "has finally" test; it is one switch on `try.Body()`,
  `try.Catches()`, `try.Finally()`.
- **One method per declaration shape** in the scope creator
  (`declareParam`, `declareCatchParam`, `declareImport`, `declarePattern`,
  ...) instead of Java's generic switch over first, second and last
  children.
- **Dead code the types exposed**: CFA searched a switch body for an EMPTY
  node the grammar never allows, and called `computeFallThrough` on blocks
  (a block falls through to itself); the scope creator scanned for-loop
  bodies that are not blocks.
- **A bug caught by the port, not the tests.** Normalize's redeclaration
  handler rewrites a duplicate `var a = b` into `a = b` while the scope is
  being built, so a declarator loop that reads the next element after the
  body would declare `b`. The loop uses `ast.ReadAhead`, with a new test.
- **ProcessDefines** reads the default value of `goog.define('name', v)` as
  the second argument (Java: child 2), and its CLOSURE_DEFINES handling,
  duplicated between a Go-only walk and Visit, is two shared helpers.

What they reported:

- **Loops that call out must read ahead.** A callback (here a redeclaration
  handler) can move the current node, which the loop cannot see: use
  `ast.ReadAhead` whenever the body calls code outside the pass.
- **`ReadAhead` calls `Next()` through the generics dictionary** (an
  indirect call), and `Declarator.Next`, `Statement.Next` and
  `Function.Name()` are just over the inlining budget. Unmeasurable in these
  passes; worth knowing for hot loops.
- **Parts that are missing to the reader are nodes to the graph.** The
  EMPTY init and update of a `for` are CFG nodes, and so is the BLOCK that
  holds a catch; `Try.Catches()` was exported for graph code, and two
  positional helpers remain for the `for` parts. An accessor returning the
  placeholder node of a `~` role would remove them.
- **No common view of statement containers** (block, script, module body,
  root, catch holder), so statement-list walks stay on rhino in both
  passes; a "block scope root" view is the same request.
- **Same child, different role names**: `object` in a for-in, `iterable` in
  a for-of; `EnhancedFor.Collection()` was added by hand. Grammar aliases
  would cover this in general. Likewise Param, ArrayPatternElem and
  PropValue all mean "a target, a target with a default, or a rest", and a
  shared binding-element view would merge three if-chains.
- **"Is n the body of its parent?"** takes one `MustX(parent).Body().Is(n)`
  per kind (nine in CFA's traversal pruning); a query for a child's role in
  its parent would make it one line.

## Twenty-sixth and twenty-seventh ports: TypeCheck and ClosureCheckModule

`TypeCheck` (3,700 lines, the pass that reports type errors) in five
commits, one per group of visit methods, and `ClosureCheckModule`. The
warnings, their order and their text are byte-identical on checks-only and
typed-app, which are full of type warnings; TypeCheck's speed is unchanged
(in-process A/B after inferTypes, pooled medians: typed-app +0.4 %,
checks-only +0.1 %, typed-full -0.6 %; same allocations).

- **Visit stays a token switch** (every node also goes through the checks
  before and after it, so a visitor could not skip any kind) and hands each
  visit method its kind: `op.Left()`/`Right()`, `call.Callee()`,
  `getProp.Object()`/`Property()`, `assign.Target()`/`Value()`; optional
  children are comma-ok (`ret.Value()`, `class.Extends()`, `name.Init()`).
- **Two copies of the argument check became one.** A Java-shaped
  `CheckArgumentsMatchParameters` nothing called any more sat next to a
  fast copy that walked siblings and skipped template strings by hand. The
  arguments of a call or `new` are `ast.MustCallOrNew(call).Args()`, those of
  a tagged template the substitutions of `tagged.Template().Parts()`.
- **Five Java argument checks became parameter types.**
- **A Java quirk kept for identical output:** `VisitOptChainGetProp`
  reports a dict access on the object, because Java's `getLastChild()` of a
  property access is its only child.
- **New categories**, each from a token list jscomp repeats:
  `CallOrNew = Call | OptChainCall | New` (`NodeUtil.isCallOrNew`, 31
  uses), `ForOfLoop = ForOf | ForAwaitOf` (42 case clauses list the two),
  `GetterOrSetterDef = GetterDef | SetterDef`.
- **ClosureCheckModule** reads a destructuring import as an object pattern
  of string keys whose values are names, and the namespace of
  `goog.require` as the call's argument.

What they reported:

- **`.Node()` is everywhere in TypeCheck (120 calls)**: the type checker
  works on JSTypes and reports on nodes, so the usual line is
  `left := op.Left().Node()`. Reading by role is the gain.
- **Do not pass a list iterator as a parameter in hot code.** Ranging over
  an `iter.Seq` parameter moves the loop state to the heap (5 allocations
  per call in a test program); the shared loop body became a local closure
  called from two typed loops.
- **Parent-role questions** (about ten here: is this getprop an assignment
  target, is this the arrow's body, the switch of a case) take the two-line
  `x, ok := ast.AsAssign(parent); ok && x.Target().Is(n)`.
- **No Must conversions on categories**: "a declarator that is not a
  destructuring is a name" is `ast.MustName(declarator.Node())`;
  `declarator.MustName()` would read better.
- **Same shape, different role names**: Binary is left/right, AssignOp is
  target/value, so the shared operator check takes rhino nodes.

## API round three: roles, aliases, Must on categories

Added before the next ports (branch `typed/a`), from the requests listed
in next steps, item 6:

| ask | added |
|---|---|
| "is n the target of its assignment?" (`x, ok := ast.AsAssign(parent); ok && x.Target().Is(n)`) | `ast.RoleOf(n) == ast.RoleAssignTarget`: one `Role` constant per child of each kind, generated with a table per parent token; EMPTY placeholders have the role they stand for (`RoleForCond`), the catch of a try is `RoleTryCatchesCatch`, and `Role.Name()` asks across kinds (`"body"`) |
| `ast.MustName(declarator.Node())` | `declarator.MustName()`: every checked conversion of a category has a Must form (they inline) |
| same child, different role names (for-in `object`, for-of `iterable`) | grammar aliases, `object\|collection:Expr`: each name gets a getter, a category gets the getter of a name all its kinds have; `EnhancedFor.Collection()` is generated now, and `RoleForInCollection` is a constant |

Not done: a common view of statement containers (none of the three ports
below walks statement lists of several container kinds; their statement
walks dispatch per kind anyway), Binary/AssignOp aliases (their `left` and
`target` have different types, Expr and Target, so no shared getter
would type both), the `~` placeholder accessor and typed AstFactory
methods.

`RoleOf` reads the parent's token and compares n with the parent's fixed
children; a parent whose children are one list (a block, a var, a
parameter list) answers from the table without reading them. It is not
free: one more record load than a token test of the parent.

## Twenty-eighth to thirtieth ports: FlowSensitiveInlineVariables, InlineFunctions, RemoveUnusedCode

`FlowSensitiveInlineVariables`, `InlineFunctions` with `FunctionInjector`
and `FunctionArgumentInjector` (5,600 lines), and `RemoveUnusedCode`
(3,700 lines), on the API above. Output, source maps and stderr
byte-identical on the four gate workloads at 1 and 4 threads, clie2e
fixtures 2,974/2,976. Pass times measured in process: the base commit's
pass, a copy with every top-level name renamed, and the new one run on
clones of the same tree, alternating, 30 pairs, at the point where the
compile first runs the pass (an A/A control of the same harness moves
1-3 %):

| pass | bundle | typed-full | typed-app |
|---|---|---|---|
| FlowSensitiveInlineVariables | 0.95-0.99 | 0.98 | 0.98 |
| InlineFunctions | 0.97-1.02 | 1.00 | 1.00 |
| RemoveUnusedCode | 0.98-0.99 | 1.00 | 0.99 |

(new/old, median of pairs; same allocation counts.)

What got clearer:

- **Roles instead of parent tests.** "Is this name purely a read" was
  seven token tests of the parent, written twice in
  FlowSensitiveInlineVariables; it is one switch on `ast.RoleOf(n)`
  (assignment target, declarator, ++/-- target, param, catch param).
  InlineFunctions' `isCandidateUsage` (declarator, function name, callee,
  object of `.call`) and the enclosing-scope skips (class extends,
  computed field key) read the same way.
- **One helper per shape the passes kept spelling out**: "a var whose
  value is a function" (three places in InlineFunctions), "the function a
  call calls directly", "the name x = f() assigns" and "the name var x =
  f() declares" (FunctionInjector's call-site classification, reused for
  the result name), "the owner of x.prototype" and "the name of
  name.prototype" (eight places in RemoveUnusedCode), "the parameter name
  or REST_MARKER" (three copies).
- **Missing parts read as missing.** A class without `extends`, a `for`
  without condition, a catch without binding: the old code traversed the
  EMPTY node (a no-op); the typed code skips it.
- **Dead walks gone**: FlowSensitiveInlineVariables visited the names of
  functions and classes (NAMEs without children) on two identical
  branches, and the gather walk had a root-function case it never takes.
- **A Java quirk kept**: InlineFunctions counts any NAME child of a
  function as a declaration, so the expression body of an arrow
  (`() => fn`) is a candidate usage.

What they reported:

- **`RoleOf` in a per-name loop costs.** RemoveUnusedCode asks about the
  parent of every name it traverses; with `RoleOf` the pass was 2-3 %
  slower, with a token switch on the parent first (and `RoleOf` only for a
  function or class parent) it is back to equal. A cheaper form when the
  caller has the parent (`RoleIn(n, parent)`) or a test for one role
  (`ast.IsRole(n, RoleAssignTarget)` reading only that position) would
  help hot code.
- **Ranging over a list whose loop body returns** compiles to a closure
  call per element (`appendCallParamBindings-range1` in the profile);
  the hot parameter-binding loop uses `FirstParam()`/`Next()`.
- **`Function.Name()` reads the string** to answer ok (an anonymous
  function has a NAME with ""), so `name, _ := fn.Name(); name.Ident()`
  resolves the string twice; InlineFunctions' NamedFunction keeps the
  name node. A `NameNode()` without the check would help.
- **Lists have no Must**: "the declaration's only declarator, which
  normalization guarantees" is `d, ok := decl.OnlyDeclarator();
  jutil.CheckState(ok)` (three times), and `a, _ := call.FirstArg()` where
  a NodeUtil test proved the argument.
- **A placeholder's node** was wanted again: RemoveUnusedCode adds a
  continuation for a class's `extends` even when it is EMPTY; the typed
  code skips it (equivalent, as traversing EMPTY does nothing).

## Performance outlook

The typed view compiles to the same code as before, so on its own it neither
speeds up nor slows down the compiler. What it can enable was measured on the
bundle and typed-full workloads:

- **Where the time goes.** Reading node records (`Node.rec` and the block
  directory load) is 16% of the CPU on bundle and 11% on typed-full;
  resolving string ids to strings (`GetString`) another 8% and 6.5%; the
  traversal functions themselves about 6%.
- **Callback dispatch is a small lever.** A compile makes 84 M (bundle) and
  225 M (typed-full) traversal callbacks. Letting a visitor declare its kinds
  (`Visitor.Tokens`) so that the traversal skips the others saves a few
  nanoseconds per skipped callback: about 1-3% at best. The traversal already
  has ad-hoc versions of this (marker interfaces such as `onlyNameLeaves`,
  used by about 28 passes); declared kinds would replace them with one
  mechanism.
- **The layout is the big lever.** Walking the whole tree costs 13-15 ns per
  node through today's linked records (two dependent loads per hop: the block
  directory, then the 40-byte record). The same recursive walk over a compact
  array in visit order (token and subtree size, first child next to its
  parent, next sibling after the subtree) costs 8-10 ns, and a linear scan
  1-1.6 ns. Walking is about 1.1 s of the 10 s bundle compile (84 M visits),
  so a compact layout could save on the order of 5%, more for passes that can
  scan instead of recursing, plus less memory and GC work.
- **Why the typed view matters for the layout.** Changes inside the linked
  model (smaller records, hot/cold field splits) are possible today, behind
  the rhino accessors. Changing the model itself (fixed child slots per kind,
  child arrays, a flat array for read-only phases) is not, because code walks
  siblings with `GetNext` and expects it to be O(1). Code written against
  `ifs.Else()` and `block.Statements()` does not care how children are
  stored, so the more passes use the typed view, the more of the layout
  becomes free to change.
- **Fewer checks.** Shapes the grammar guarantees need no defensive code
  (three of the four cases in `TryMinimizeIfBlockExits`); small, but free.

## What we learned

- **The types expose dead code.** Defensive branches for shapes the tree
  never has (an `if` branch that is not a block) become visible once the
  accessor says `Block`.
- **The port drifts from the Java.** The typed pass keeps the Java method
  names and `Java:` comments, but it no longer follows the Java line by line.
  That is the price of readable Go; a pass ported this way should say so
  (`Go-only:` comment at the top).
- **Iteration semantics need both kinds.** The iterators copy the classic
  loop (the next element is read after the body), which is right when the
  body edits what follows (merging the next declaration into the current
  one). A loop that moves the current element, or calls code that might
  (a callback, a handler), needs the next one read first: `ast.ReadAhead`.
  Neither choice fits every loop, so both exist.
- **API size.** 7,500 generated lines, about 860 conversion methods: a
  `Name` converts to the 13 categories it belongs to. Generated code is cheap
  to keep, but the method lists are long in an editor; categories that only
  exist to type one child could stop generating conversions.
- **Same token, different shapes.** A `BLOCK` is a block statement, the
  catch holder of a `TRY`, or a class static block. The grammar marks the
  catch holder as a kind that only exists inside `Try`, so `AsBlock` stays
  unambiguous.

## Next steps

1. **Pre-order handlers** in the typed visitor (`ShouldTraverse` per kind).
   Declared kinds are done: `NodeTraversal_TraverseVisitor` hands the
   traversal the visitor's token set, and the traversal does not call the
   callback for other nodes. A one-handler visitor over three.js's AST now
   takes 2.21 ms, against 2.43 ms for a hand-written token switch and
   3.07 ms before (minimums of 15 interleaved runs).
2. **Port more passes.** 27 are on the typed view, among them the CFG,
   the scope creator and the type checker. Left with the most positional
   code: `RemoveUnusedCode`, `ProcessCommonJSModules`, `TypeInference`,
   `FlowSensitiveInlineVariables`, `CodeGenerator`, `TypedScopeCreator`,
   `InlineFunctions`, `ClosureRewriteModule`, `Es6NormalizeClasses`.
6. **The API requests that came back in several ports**, by count
   (done: `RoleOf`, Must conversions on categories, role aliases; see
   "API round three"):
   - a cheaper role test for hot loops (`RoleIn(n, parent)`, one role at
     one position), `Function.NameNode()`, Must on list accessors;
   - a common view of statement containers (block, script, module body,
     root, catch holder), so statement-list walks can be typed;
   - aliases for Binary left/right and AssignOp target/value (different
     types: needs a Target-or-Expr view);
   - the placeholder node of a `~` role, for graph code that needs the
     EMPTY init of a `for`;
   - typed AstFactory methods for the transpilation passes.
3. **Change the storage behind the API.** Callers ask for `ifs.Else()`, not
   for the third child, so the layout becomes an implementation detail:
   fixed fields per kind, child arrays or no `previous` link can be tried
   without touching the passes that use the typed view.
4. **Retire structural checks from AstValidator**: the grammar covers child
   counts and kinds; the validator keeps the semantic rules (features, scopes,
   types).
5. **Inlining.** Done for the category conversions and `LastX` (token sets,
   a cheaper `GetLastChild`); `If.Else`, `For.Cond` and `Next()` are still
   just over the budget.
