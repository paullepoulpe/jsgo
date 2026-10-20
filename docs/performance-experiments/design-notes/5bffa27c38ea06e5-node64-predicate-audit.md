# Compact Node predicate audit

Source baseline: `3594959c` (compact Node plus common cold payload header).
No production change or performance claim results from this audit.

The suspected simple-predicate overhead is already eliminated in this baseline.
There are 119 `Node.Is*` methods that compare packed token bits directly:

```go
return n.meta & nodeMetaTokenMask == uint32(Token_NAME)
```

This includes the two-token OR predicates IsRest and IsSpread. Composed helpers
such as IsDestructuringPattern and IsNormalBlock call these already-specialized
predicates. They do not call GetToken or read the cold overflow payload.

All current enum tokens are below 255. For any comparison constant `c` with
`0 <= c < 255`, direct packed equality is exactly equivalent to GetToken()==c.
A narrow stored token has its exact value in the low byte; an arbitrary negative
or large token stores sentinel 255 there and its actual int32 value in the cold
header. Sentinel 255 can never equal c. Token 255 itself must also use the cold
path. Reassigning an extended node to a narrow token updates the low byte and
retains this proof even if the cold record still exists. No payload kind,
original name, length, GC bitmap, or type/color state participates in equality.
GetToken/SetToken must continue preserving all arbitrary int32 token values.
The existing TestNodePayloadAllCurrentTokensFitFastPredicates guards the enum
range, and TestNodePayloadKindsAndMetadataTransitions covers wide/narrow token
reassignment across all payload kinds.

The remaining Is* methods that inspect GetToken are IsQualifiedName,
IsUnscopedQualifiedName and IsValidAssignmentTarget. Their fixed-case switches
could inspect a private packed discriminator instead: every accepted case is
below 255 and every overflow token follows the original default branch. Similar
constant-only switches occur in getQualifiedNameString,
MatchesQualifiedNameWithEnd and MatchesQualifiedNameNode. Six NAME inequalities
also retain GetToken (closure-unaware shadow lookup, two printing helpers,
MatchesName and both operands of MatchesNameNode). These are possible bounded
followups only if profiles attribute material CPU to these exact helpers.

Dynamic token identity comparisons, cloning, stringification and the public
GetToken API require the full token and cannot be replaced with byte equality:
distinct wide values share sentinel 255. In particular the dynamic equality
check in MatchesQualifiedNameNode should remain original unless its entire
control flow is proven, rather than applying a mechanical equality rewrite.

No helper abstraction, predicate rewrites, build, test, profile or benchmark was
added or run. The common-header optimization's smaller GetToken machine code
therefore cannot be multiplied by simple IsName/IsCall invocation counts: those
predicates already bypassed GetToken. This audit provides no evidence that a
further predicate change would recover the compact layout's macro regression.
