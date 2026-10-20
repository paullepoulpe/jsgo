# Targeted original fragments for builtin extern bindings

This prototype keeps native special-name binding rejection intact. It routes
only two proven top-level EXTERN statement forms to the existing original
statement parser, before VAR/FUNCTION or leading JSDoc state is consumed:

* `var arguments;`, with the exact plain identifier, no initializer or other
  declarators, and an explicit semicolon;
* `function eval(code) {}`, with those exact plain identifiers and empty body.

Whitespace and comments are allowed between tokens. Source scripts, nested
statements/functions, function expressions, escaped identifiers, other eval
signatures/bodies, eval/arguments parameters, let/const declarations,
initializers, multiple declarators and implicit-semicolon forms retain their
previous native rejection. No new environment flag is added; this is a narrow
coverage extension of the existing native parser experiment.

The original parser accepts these names as ordinary identifier bindings.
`StrictModeCheck.NonExternChecks` reports eval/arguments declarations later and
visits only the source root, so the standard extern declarations are permitted.
Implicit function arguments are a later scope-creation concern. Reference uses
of arguments and eval calls already use the native expression path; shared
`AnnotateCalls` retains DIRECT_EVAL/FREE_CALL metadata. The original fragment
builds the same NAME/function nodes, JSDoc attachments, source positions and
feature metadata, and shared final validation still runs on the finished AST.

The fresh stock inventory at commit051fa548 found only es3.js falling back:
76,751 bytes, initially rejected at `var arguments;` (line275, UTF-16 offset6807).
It also contains `function eval(code) {}` at line409. Supporting only the first
fragment would expose the second rejection and leave the whole file on fallback.
This prototype includes both narrow forms; full-file acceptance must still be
verified rather than inferred from source presence.

Focused tests compare exact original AST/JSDoc/source/properties/features,
ordered diagnostics and exceptions in strict/sloppy modes. They require explicit
fragment counts and preserve rejection boundaries. The pinned full builtin
inventory is the next whole-file oracle. No build, runtime performance or
coverage gain is claimed for this source-only prototype yet.

The potential remains bounded: this file was3.88% of BROWSER extern source bytes
and26.03% of CUSTOM extern source bytes. Those fractions are not CPU fractions;
no AST is cached or preparsed for a later compilation, and all fragments are
parsed freshly in the current request.
