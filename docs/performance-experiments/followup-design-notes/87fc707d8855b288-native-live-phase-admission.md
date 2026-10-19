# Native liveness admission after color conversion

Status: isolated source change from c7998d45, unbuilt/untested. No runtime GC,
BitSet representation, scratch reservation, publication or solver changes.

The old nativeLiveVariablesCompilerSafe rejected any compiler with nonnil
runtimeLibraryTypedAsts. ConvertTypesToColors.Process always calls
InitRuntimeLibraryTypedAsts in both ordinary paths: after removing JSTypes when
type checking did not run, or while building colors when it did. This precedes
dead-assignment elimination and coalescing. Consequently an ordinary full stock
CLI pipeline could set the experiment flag yet never execute the native solver;
earlier parse/normalize-only fixtures did not reproduce this admission state.
Previous flag-on performance results alone cannot establish native liveness
performance. Real full-pipeline dispatch counts are required before retiming.

The fix removes this one native-liveness-specific nonnil-table rejection. The
runtime table is consulted by Compiler.loadResourceContents for explicit runtime
library injection. Neither native/generic liveness, dead-assignment elimination
nor coalescing accesses it or invokes its suppliers. Their scope/CFA/transfer/
mutation calls do not inject runtime code. Table presence, table size and
supplier identity therefore cannot affect this phase's semantic equivalence.
The original fresh private CLI source ownership, exact compiler/default factory,
normalized input, known convention/reporter/warnings/change-observer, and
non-typed-AST-input proof remain mandatory. Public consumer constructors and
custom callbacks still use the original analysis. No shared/global admission
fence or unrelated experiment guard changes.

The complete-lifecycle test uses actual stock runner compiler/options/extern
factories, installs the genuine ownCLISourceFiles lease before Compiler.Compile,
and releases it afterward. It attaches existing numeric probes to the original
default DAE/coalescing factories; it does not use customPasses, fake ownership,
test mode or a global observer. Both consumers must see the real populated
post-conversion runtime table and perform native analysis at N1/N2 for SIMPLE
and ADVANCED. JS, full AST, source maps and diagnostics must match the generic
solver. Separate throwing-supplier fixtures prove that the inert table is never
consulted and public/unowned/custom/typed-input fallbacks remain original.
Untouched CommandLineRunner.DoRun diagnostics on pinned real tasks should also
count solver dispatch before any timing claim.
