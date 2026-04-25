# Body-Level Parser + IRFactory Fidelity Audit

Read-only audit comparing
- `internal/parser/parser.go` (legacy direct-to-rhino parser; the **default**
  pipeline used by `parser.Parse` → `internal/jscomp/compiler.go:203`) against
  `closure-compiler/src/com/google/javascript/jscomp/parsing/parser/Parser.java`.
- `internal/jscomp/parsing/irfactory.go` (the new-pipeline IRFactory used when
  `JSGO_NEW_PIPELINE=1`) against
  `closure-compiler/src/com/google/javascript/jscomp/parsing/IRFactory.java`.
- `internal/parser/parser_new.go` (1:1 port of `Parser.java`, fed into IRFactory).

The legacy `parser.go` is the path triggered by `testinferredreturn2.js`. The
audit therefore weights position-propagation discrepancies in `parser.go`
(the actually-running code) higher than the equivalent gaps in `irfactory.go`
(off-by-default).

Severity legend
- **high** — currently produces wrong diagnostic line/col or different AST shape
- **medium** — currently silent, but will diverge when the gated surface is wired
- **low** — cosmetic / Go-language-adapter difference with no behavioural impact

Glossary: "Java line N" refers to `IRFactory.java` unless the row's File column
points to `Parser.java`.

---

## Part A — `internal/parser/parser.go` (legacy direct-to-rhino path)

The legacy parser fuses Java's `Parser.java` and `IRFactory.java` into one
pass: it builds `rhino.Node` directly without going through `ParseTree`.
`setSourceInfo` does not exist as a single function; positions are set ad-hoc
via `Node.SetPos(line, col)`. This is the source of the
`(new Foo()).bar()` col-3-vs-col-2 bug in `testinferredreturn2.js`.

| File | Method | Java line | Discrepancy | Severity | Suggested Fix |
|------|--------|-----------|-------------|----------|---------------|
| parser.go | `tryParen` (~2391) | Parser.java:2132 / IRFactory.java:2699 | Only sets pos on inner if `Lineno()==0 && Charno()==0`; Java's `ParenExpressionTree.location.start` records the LP and `processParenthesizedExpression` returns the inner node WITHOUT moving its pos — but the **outer** call/member chain re-reads `getTreeStartLocation()` which IS the LP. Go's outer `parseUpdateExpression` instead reads `expr.Lineno()/Charno()` after `tryParen` returns, losing the LP. | **high** | Make `tryParen` unconditionally `inner.SetPos(openLine, openCol)` (record the LP) so the chained `.bar` / `(args)` propagation in `parseUpdateExpression` inherits the correct start. This is the direct cause of the `testinferredreturn2.js` arg-col-3 failure. |
| parser.go | `parseUpdateExpression` `.` chain (~2056-2069) | Parser.java:2898 (`MemberExpressionTree(getTreeLocation(start), …)`) | GETPROP node positioned at `expr.Lineno(), expr.Charno()` — i.e. the receiver's pos, mutated by every previous chain step. Java records `start = getTreeStartLocation()` ONCE (the start of the leftmost operand or its LP) and reuses it for every chained tree. Once `tryParen` is fixed the chain inherits LP, but the chain logic itself should record `start` once at the top of `parseUpdateExpression` rather than threading through `expr`. | **high** | Capture `chainStartLine, chainStartCol := p.scanner.Line(), p.scanner.Col()` BEFORE `parsePrimaryExpression` and use that for every chain node's `SetPos`. |
| parser.go | `parseUpdateExpression` CALL (~2094-2103) | Parser.java:2887 | Same as `.` chain: `call.SetPos(expr.Lineno(), expr.Charno())`. CallExpressionTree in Java is `getTreeLocation(start)` where `start` is the leftmost operand. | **high** | Use the captured chain start, see above. |
| parser.go | `parseUpdateExpression` GETELEM (~2087-2092) | Parser.java:2893 | Same: positioned at `expr.Lineno(), expr.Charno()`. | **high** | Use captured chain start. |
| parser.go | `parseUpdateExpression` OPTCHAIN_* (~1979-2050) | Parser.java:2924-2966 | Optional-chain segments also use `expr.Lineno(), expr.Charno()`. Java uses `start = operand.location.start` for the *first* `?.` and `getTreeLocation(start)` thereafter — semantically the same as the chain start. | **medium** (correct as long as chain start is correct) | Once chain-start fix lands, audit the OPTCHAIN positions for parity. |
| parser.go | `parseArguments` (~2141-2166) | Parser.java:parseArguments / IRFactory.java:1854 | Captures `argLine, argCol` BEFORE `parseExpression()` and unconditionally `arg.SetPos(argLine, argCol)`. This is correct but masked by the receiver chain pulling its CALL node's pos from the inside-out — see test failure reproduction. The fix to `tryParen` makes this code path's overwrite redundant but not incorrect. | low | None (this method is OK; the bug is upstream). |
| parser.go | `parsePrimaryExpression` `(` branch (~2235-2251) | Parser.java:1657 | `startLine, startCol` captured but NOT used to position the returned paren-expression's pos (only used for arrow). Java threads `getTreeStartLocation()` down to `ParenExpressionTree`. | **high** | Pass `startLine, startCol` into `tryParen` so the LP can be authoritatively assigned. |
| parser.go | `parseNewExpression` (~3128-3163) | Parser.java:parseNewExpression / IRFactory.java:2516 | Positions NEW node at the constructor IDENT (`Foo`) on line 3162, then caller in `parseUpdateExpression:1948` re-positions to the `new` keyword. Two writes is fragile. | low | Set NEW pos at the `new` keyword in `parseNewExpression` itself; caller no-op. |
| parser.go | `parseArrayLiteral` (~2424-2458) | Parser.java:parseArrayLiteral / IRFactory.java:1481 | No `setTrailingComma` propagation; each hole emits `UNDEFINED` instead of `EMPTY`. Java emits `EMPTY` for elisions; Go's `UNDEFINED` is semantically identical for type checking but diverges from Rhino IR shape. | medium | Track `hasTrailingComma` and emit `rhino.EMPTY` for holes once any pass reads either. |
| parser.go | `parseObjectLiteral` (~2496) | IRFactory.java:2542 | Java calls `maybeWarnKeywordProperty` per non-computed/non-quoted/non-spread key; Go does not. ECMASCRIPT3 mode produces `INVALID_ES3_PROP_NAME` warnings in Java. | medium | Add keyword-property warning when ES3 mode lands. |
| parser.go | `parseUnaryExpression` (~1842-1928) | Parser.java:parseUnaryExpression / IRFactory.java:2911 | No `maybeWarnForFeature` calls. Java warns for `EXPONENT_OP` (`**`), `await`-without-async, `delete` on private fields, etc. | medium | Wire feature warnings when ES-mode gating is implemented. |
| parser.go | `tryBinary` (~1772-1840) | IRFactory.java:2288 (`markBinaryExpressionFeatures`) | No feature marking for `**`, `**=`, `??`, `??=`, `||=`, `&&=`. | medium | Mirror `markBinaryExpressionFeatures` once `features` set is exposed. |
| parser.go | `parseExpressionStatement` (~281) | IRFactory.java:processExpressionStatement | EXPR_RESULT pos is set from the inner expression rather than the statement's own start. Java `transform(exprStmtTree)` would set EXPR_RESULT pos from the ExpressionStatementTree.location.start. | low | Verify diagnostics that ride the EXPR_RESULT pos (pretty-printer, JSC_USELESS_CODE) match Java; align if any drift. |
| parser.go | `parseReturnStatement` (~1502) | IRFactory.java:processReturnStatement | RETURN node correctly positioned at `return` keyword; child expression positioned via parseExpression result. Looks parity-clean. | low | None. |
| parser.go | `finishFunction` (~714) | IRFactory.java:processFunction | Missing `node.setTrailingComma(...)` from formal-param list, and missing `attachPossibleTrailingComment(node, functionTree.getEnd())`. Trailing-comment attachment is silent today. | medium | Wire trailing-comment attachment when a non-JSDoc-comment pass lands. |
| parser.go | `finishClass` (~820) | IRFactory.java:processClassDeclaration is implemented elsewhere (see parserrunner) | `EMPTY` superclass and `CLASS_MEMBERS` pos derived from current scanner position; Java derives from class tree itself (start of `class` keyword). | low | Cosmetic — diagnostics rarely cite super-EMPTY pos. |
| parser.go | `parseClassMember` (~869) | IRFactory.java:processField / processGetAccessor / processSetAccessor | `MEMBER_FIELD_DEF` does NOT call `setLength` (Java IRFactory.java:2168 sets length from tree.location). | low | No length field in `rhino.Node`; document divergence and skip. |
| parser.go | `parseTemplateLiteral` (~3222) | IRFactory.java:processTemplateLiteral (1369 in irfactory.go) | No `TAGGED_TEMPLATELIT` handling at this level (legacy parser.go); only bare `TEMPLATELIT`. Tag handling is in the `\`` arm of the chain at line 2109. | low | Accept divergence; the arm is functional for current corpus. |
| parser.go | `tryArrow` (~3029) | Parser.java:parseArrowFunction (post-cover) | Speculative parse; on success returns FUNCTION at `(` start. Java uses the `(` LP from `getTreeStartLocation()`. Currently positions at `(` which matches. | low | None. |
| parser.go | `consumeParenExpression` (~1475) | Parser.java:eat(OPEN_PAREN)+parseExpression+eat(CLOSE_PAREN) | Returns inner node without recording LP (used by `if/while/for/switch` headers). Diagnostic positions fall on the inner expression, matching Java's behaviour for those statements (the IF/WHILE/etc. node carries the keyword's pos). | low | None. |
| parser.go | `parseLeftHandSideExpression` (~1654) | Parser.java:2854 | Go's helper is a thin wrapper around `parseUpdateExpression`. Missing the speculative NEW vs CALL branching from Java; Go's `parseUpdateExpression` already handles both inline. | low | Naming convention difference only. |

## Part B — `internal/jscomp/parsing/irfactory.go` (new pipeline)

This file is exercised only when `JSGO_NEW_PIPELINE=1`. Its `setSourceInfo`
helper exists and is called in many places, but the guard differs from Java
and several Java-side hooks are missing.

| File | Method | Java line | Discrepancy | Severity | Suggested Fix |
|------|--------|-----------|-------------|----------|---------------|
| irfactory.go | `setSourceInfo(node, tree)` (~406) | IRFactory.java:1226 | Guard is `node.Lineno() == 0`; Java guard is `irNode.getLineno() == -1`. A node legitimately on line 1 col 0 looks "unset" to Go and gets re-positioned. | **medium** (currently no rhino.Node defaults to lineno=0+charno=0 in valid input, but boundary cases like SCRIPT root drift) | Initialise `rhino.Node.lineno` to `-1` (or add an explicit `posSet` flag) and use Java's sentinel semantics. |
| irfactory.go | `setSourceInfoFromToken` (~1466) | IRFactory.java:1232 | Same `Lineno() == 0` guard issue. | medium | Same fix. |
| irfactory.go | `Transform` (~155) | IRFactory.java:1045-1059 (`transformDispatcher.transform`) | Missing the `setSourceInfo(node, tree)` call at the **end** of Transform when JSDoc isn't attached: Java unconditionally calls `setSourceInfo(node, tree)` on the way out. Go calls it inside the `if node != nil` branch but only once and the guard inside `setSourceInfo` (Lineno==0) often skips the assignment. | **high** when paired with the guard bug | Match Java's order: process → optional JSDoc/comment attachment → unconditional `setSourceInfo`. Drop the `Lineno==0` guard or use `-1` sentinel. |
| irfactory.go | `processName` (~1063) | IRFactory.java:2406 | Java calls `setSourceInfo(node, identifierToken)` **and** `maybeWarnReservedKeyword`, plus parses inline JSDoc on the token. Go does none of these inside processName; it relies on outer `Transform` to set pos. Ident-level positions are correct because outer Transform runs setSourceInfo, but the warning surface and JSDoc attachment are missing. | medium | Add `f.setSourceInfoFromToken(node, tree.IdentifierToken)` and a TODO for `maybeWarnReservedKeyword`. |
| irfactory.go | `processLiteralExpression` (~1072) | IRFactory.java:2527 (number) / 2534 (bigint) / 2440 (template) / 2427 (string) | Java's processString/processNumber calls `setSourceInfo(node, literal)` and parses non-JSDoc inline comments on the token. Go does neither inside the function. Outer Transform compensates for pos but not for comment attachment. | medium | Add `f.setSourceInfoFromToken(node, tok)` and TODOs for comment attachment. |
| irfactory.go | `processParenthesizedExpression` (~1190) | IRFactory.java:2699 | Java calls `checkParenthesizedExpression(exprNode)` first (errors on `(a, ...rest)`); Go does not. Outer Transform's `setSourceInfo` won't move the inner pos because Java's guard is `lineno==-1` and Go's is `lineno==0` — both effectively "do nothing" since the inner was already positioned. So Go's processParenthesizedExpression matches Java's behaviour for the *position* but skips the rest-in-comma-paren error. | medium | Port `checkParenthesizedExpression`. Note: Java's processParenthesizedExpression does NOT itself set LP — the LP comes from the outer Transform setting `parenTree`'s start onto the returned node IFF lineno was unset, which is then a no-op. So neither path records LP for the inner node. The LP propagation in Java actually rides on `getTreeLocation(start)` in `parseLeftHandSideExpression` (Parser.java:2887). Mirror that in `parser_new.go`. |
| irfactory.go | `processFunctionCall` (~1223) | IRFactory.java:1854 | Missing `node.setTrailingComma(callNode.arguments.hasTrailingComma)`, `getEndOfArgCommentZones` + `attachPossibleTrailingCommentsForArg` per arg. Free-call annotation is via `setFreeCallIfNeeded` (similar to Java's `annotateCalls`) but missing the special-case where `(0, real.callee)(...)` is collapsed (IRFactory.java:1903). | medium | Port `annotateCalls` collapse; wire trailing-comma. |
| irfactory.go | `processNewExpression` (~1241) | IRFactory.java:2516 | Missing `node.setTrailingComma(exprNode.hasTrailingComma)`. | low | Add when trailing-comma serializer pass needs it. |
| irfactory.go | `processArrayLiteral` (~1281) | IRFactory.java:1481 | Missing `node.setTrailingComma(tree.hasTrailingComma)`. | low | Wire when needed. |
| irfactory.go | `processObjectLiteral` (~1291) | IRFactory.java:2542 | Missing `setTrailingComma`, `maybeWarnKeywordProperty` per key, `maybeWarnForFeature(SHORTHAND_OBJECT_PROPERTIES)`, error on `DEFAULT_PARAMETER` at top level (`"Default value cannot appear at top level of an object literal."`), `maybeReportGetter`/`maybeReportSetter`. | medium | Stage: error first (high-fidelity diagnostic), then warnings, then trailing-comma. |
| irfactory.go | `processArrayPattern` (~not present at array-pattern path) | IRFactory.java:1491 | `maybeWarnForFeature(ARRAY_DESTRUCTURING)` and `transformNodeWithInlineComments` per element. Go's array-pattern path goes through `Transform` directly without the inline-JSDoc handling. | medium | Route array/object pattern elements through an `transformNodeWithInlineComments` helper. |
| irfactory.go | `processObjectPattern` | IRFactory.java:1513 | `maybeWarnForFeature(OBJECT_DESTRUCTURING)`, `processObjectPatternElement` dispatch (DEFAULT_PARAMETER, PROPERTY_NAME_ASSIGNMENT, COMPUTED_PROPERTY_DEFINITION, OBJECT_REST). | medium | Port the dispatcher; `OBJECT_REST` is currently best-effort. |
| irfactory.go | `processFunction` (~740) | IRFactory.java:processFunction (around 2000-2154) | Missing: error reporting on unnamed declarations (`"unnamed function statement"` IRFactory.java:2100), `attachPossibleTrailingComment`, `parseDirectives(bodyNode)`, multiple `maybeWarnForFeature` (ARROW_FUNCTIONS, ASYNC_FUNCTIONS, ASYNC_GENERATORS), `IR.block()` repair under KEEP_GOING. | medium | Port the error path first; warnings later. |
| irfactory.go | `processClassDeclaration` (~841) | (no direct match in IRFactory.java for `process` w/ this name; ClassDeclaration is dispatched at ~3534) | Missing `maybeWarnForFeature(CLASSES)`. | low | Wire when feature gating lands. |
| irfactory.go | `processVariableStatement` (~890) | IRFactory.java has no such method; statements route through transform | Java handles `var/let/const` via `transformList` of `processVariableDeclaration` — does not delegate from a wrapper. Go's wrapper-then-delegate pattern produces correct AST shape but the delegate's setSourceInfo runs against the inner declaration list, not the wrapping statement. | low | Inline the body of `processVariableDeclarationList` here so the outer `Transform`'s setSourceInfo fires against the VariableStatementTree. |
| irfactory.go | `processBinaryExpression` (~1144) | IRFactory.java:2309 | Missing the comment-tracking branch (`jsdocTracker.hasPendingCommentBefore(exprNode.right.location.start)`) which determines in-order vs out-of-order traversal, AND missing `processBinaryExpressionHelper` iterative deep-binary-op handling that prevents stack overflow on long string concatenations. Currently recursive Go can blow the stack for 2000-deep `a + b + c + ... + z`. | medium | Port `processBinaryExpressionHelper`. |
| irfactory.go | `processBinaryExpression` (~1144) | IRFactory.java:2288 | Missing `markBinaryExpressionFeatures` (warns on `**`, `??`, `??=`, `||=`, `&&=`). | medium | Port the feature switch. |
| irfactory.go | `processForLoop` (~613) | IRFactory.java:1825 | Java uses `transformOrEmpty(loopNode.initializer, loopNode)` for each of init/cond/incr; Go's path does the same via `transformOrEmpty` — looks fine. Verify EMPTY pos copies parent. | low | Verify in test. |
| irfactory.go | `processGetAccessor` (~1319) | IRFactory.java:2629 | Java passes `IdentifierType.CAN_BE_PRIVATE` for class-member case; Go ignores private-id distinction. | medium | Track private-name flag once private-class-fields land. |
| irfactory.go | `processSetAccessor` (~1340) | IRFactory.java:2647 | Same as above. | medium | Same. |
| irfactory.go | `processTemplateLiteralPortion` (~1396) | IRFactory.java:2440 | Returns `NewTemplateLitString("", "")` — empty cooked + raw because Token has no Value field. Java's `processTemplateLiteralToken` calls `normalizeString(token, true)`. | **high** (template literals stringify to empty) | Plumb cooked/raw values through TemplateLiteralPortionTree's token. |
| irfactory.go | `transformBlock` semantic | IRFactory.java:718 | Java's `transformBlock` calls `setSourceInfo(newBlock, irNode)` (the **node** overload that copies lineno/charno from `irNode`). Go's `processBlockTree` calls `setSourceInfo(block, tree)` (the **tree** overload). Same start line/col in normal cases; diverges if the inner node was relocated. | low | Mirror the node-overload to be paranoid; add `setSourceInfoFromNode(node, ref)` helper. |
| irfactory.go | `processAstRoot` (~434) | IRFactory.java:1653 | Missing `parseDirectives(scriptNode)` and the `isGoogModule` rewrite. | medium | Port directives parsing when ES module gating lands. |

## Part C — `internal/parser/parser_new.go` (1:1 port of Parser.java, fed into IRFactory)

This file IS structurally faithful (it uses `getTreeStartLocation` /
`getTreeLocation` and constructs `*Tree` records), so the fidelity story
is much shorter.

| File | Method | Java line | Discrepancy | Severity | Suggested Fix |
|------|--------|-----------|-------------|----------|---------------|
| parser_new.go | `parseCoverParenthesizedExpressionAndArrowParameterList` (~2025) | Parser.java:2132 | Faithful port — `getTreeStartLocation` records LP, `ParenExpressionTree(getTreeLocation(start), result)` carries LP→RP range. This is the layer that **does** propagate the LP correctly. | low | None. |
| parser_new.go | `parseLeftHandSideExpression` (~1266) | Parser.java:2854 | Faithful port — `start := p.getTreeStartLocation()` BEFORE `parseNewExpression` and reused for every chain step. `CallExpressionTree(getTreeLocation(start), …)` correctly carries the LP. | low | None. |
| parser_new.go | `getTreeStartLocation` / `getTreeLocation` (~275-290) | Parser.java:3542-3556 | Faithful 1:1 port. | low | None. |
| parser_new.go | `parseStatementStandard` etc. | Parser.java:1140 | Structural shape matches; statement bodies in `parser_statements.go` are 1:1. | low | None. |

## Cross-cutting items

| File | Method | Java ref | Discrepancy | Severity | Suggested Fix |
|------|--------|----------|-------------|----------|---------------|
| rhino/node.go (referenced) | `SetPos` | Node.setLinenoCharno + sentinel | `lineno` defaults to **0**; Java's `Node` uses **-1** as "unset". Every `setSourceInfo` guard in irfactory.go therefore has wrong semantics on the SCRIPT root and on any genuine line-1-col-0 token. | **medium** (couples to all irfactory guards) | Initialise rhino.Node fields to -1 in `rhino.New`/`rhino.NewString`/`rhino.NewNumber`. |
| irfactory.go (multiple) | trailing-comment attachment | IRFactory.java:`attachPossibleTrailingComment{,sForArg}` | Entirely absent from Go IRFactory. | medium | Out of scope for current corpus; flag and revisit when comment-preserving printer lands. |
| irfactory.go (multiple) | `maybeWarnForFeature(...)` | IRFactory.java:1175-1219 | Missing across all process* methods. The Go IRFactory has no `features` set at all. | medium | Add a feature-set field and wire through the most user-visible warnings (ARROW, ASYNC, SPREAD, REST, EXPONENT, NULL_COALESCE) first. |
| irfactory.go | inline-JSDoc on `processNameWithInlineComments` | IRFactory.java:2461 | Go has only `findPrecedingJSDoc`, which skips by offset and walks the JSDoc list. Java parses JSDoc directly off the token's start position. The two strategies converge for typical input but diverge when comments are stacked. | low | Document; revisit if a corpus diff reproduces. |

---

## Reproduction note for `testinferredreturn2.js`

`g((new Foo()).bar())` at line 4. Java reports col 2 (the inner LP). Go
reports col 3 (`new`). The decisive call chain in the **legacy**
`parser.go` is:

1. `parseArguments` captures `argLine, argCol = (1, 2)` at the LP.
2. `parseExpression` → `parseAssignment` → ... → `parseUpdateExpression`.
3. `parseUpdateExpression` does NOT capture a chain-start before
   `parsePrimaryExpression`. After the inner `parsePrimaryExpression →
   tryParen` returns the bare `NEW` node positioned at col 3, every
   subsequent `.bar` and `(args)` step copies `expr.Lineno()/Charno()`
   = col 3.
4. Back in `parseArguments`, the unconditional `arg.SetPos(argLine, argCol)`
   *should* overwrite to col 2, but downstream the type checker reads the
   CALL's children positions (the receiver chain), not the CALL's own pos,
   for the actual-parameter diagnostic — so the col-3 leak surfaces.

Two fix paths, both inside `parser.go`:
- **A (minimal):** in `tryParen`, **unconditionally** set `inner.SetPos(openLine, openCol)` so the receiver chain inherits LP from the inside.
- **B (Java-faithful):** in `parseUpdateExpression`, capture `chainStartLine, chainStartCol := p.scanner.Line(), p.scanner.Col()` before `parsePrimaryExpression` and use that instead of `expr.Lineno()/Charno()` on every chain node — mirrors `getTreeStartLocation()` + `getTreeLocation(start)` in `Parser.java:2854-2898`.

Path B matches Java more precisely; Path A is one line.
