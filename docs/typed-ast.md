# A typed view of the AST

Status: prototype. Package `internal/ast`, generated from
`internal/ast/grammar.txt`; one pass (`MinimizeExitPoints`) uses it.

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
that does not follow the grammar. The `AstValidator` (which every
`CompilerTestCase` runs after each pass, and every optimized compile runs once
at the end) now calls it, so a grammar that drifts from the trees the
compiler builds fails the tests.

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

1. **Typed traversal.** Generate a dispatcher so a pass can write
   `VisitIf(t, ifs ast.If)` instead of switching on tokens.
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
