# A typed view of the AST

Status: prototype. Package `internal/ast`, generated from
`internal/ast/grammar.txt`; three passes use it (`MinimizeExitPoints`,
`ConvertToDottedProperties` with a typed visitor, and the parts of the hot
`InlineVariables` that read tree shapes).

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
ported loops keep their meaning when the body edits the tree.

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
- **Iteration semantics need a decision.** The iterators copy the classic
  loop (the next element is read after the body). A loop that removes the
  current element still has to save the next one first, as today.
- **API size.** 7,500 generated lines, about 860 conversion methods: a
  `Name` converts to the 13 categories it belongs to. Generated code is cheap
  to keep, but the method lists are long in an editor; categories that only
  exist to type one child could stop generating conversions.
- **Same token, different shapes.** A `BLOCK` is a block statement, the
  catch holder of a `TRY`, or a class static block. The grammar marks the
  catch holder as a kind that only exists inside `Try`, so `AsBlock` stays
  unambiguous.

## Next steps

1. **Pre-order handlers and declared kinds** in the typed visitor
   (`ShouldTraverse` per kind; skip the callbacks of other kinds), replacing
   the traversal's marker interfaces.
2. **Port more passes**, starting with those that are mostly about shapes
   (peephole passes, `Normalize`, `ControlFlowAnalysis`), and measure each.
3. **Change the storage behind the API.** Callers ask for `ifs.Else()`, not
   for the third child, so the layout becomes an implementation detail:
   fixed fields per kind, child arrays or no `previous` link can be tried
   without touching the passes that use the typed view.
4. **Retire structural checks from AstValidator**: the grammar covers child
   counts and kinds; the validator keeps the semantic rules (features, scopes,
   types).
5. **Inlining.** Make `LastX` and the hand-written accessors fit the
   inlining budget (a cheaper last-child primitive in rhino).
