# internal/jscomp/skel — generated Java→Go skeletons

Staging area for 158 closure-compiler Java files that haven't been
fully ported yet. Each file was produced by `tools/javaport` from
the upstream source at `v20260415` and carries a
`//go:build skeletons` tag so it's excluded from the default build.

## Status

- **Structs + method signatures**: generated mechanically by
  `javaport`. Types are mapped to existing jsgo packages where
  possible (`Node` → `*rhino.Node`, `JSType` → `jstype.JSType`, …);
  unmapped types leak through verbatim with a `// TODO: not
  mapped` comment.
- **Bodies**: empty `// TODO: port body` stubs. Every method
  returns the zero value for its declared Go return type.
- **Javadoc**: preserved as Go doc comments above each declaration.
- **Diagnostic codes**: `DiagnosticType.{error,warning,disabled}(…)`
  calls extracted into a commented catalogue at the top of each
  file so you can paste them into `internal/diag/diag.go` as you
  port the pass.

## Layout

```
internal/jscomp/skel/
├── jscomp/                    # com/google/javascript/jscomp/*
│   ├── checkaccesscontrols.go
│   ├── typecheck.go           # (full, un-narrowed TypeCheck port target)
│   ├── parsing/               # com/google/javascript/jscomp/parsing/*
│   │   └── parser/            # com/google/javascript/jscomp/parsing/parser/*
│   └── typeinfer/             # com/google/javascript/jscomp/type/*
│                              # (renamed — `type` is a Go keyword)
└── rhino/                     # com/google/javascript/rhino/*
    └── jstype/                # com/google/javascript/rhino/jstype/*
```

## Workflow

1. Pick a skeleton, e.g. `skel/jscomp/checkaccesscontrols.go`.
2. Move it out of the skel tree (e.g. into
   `internal/jscomp/checkaccesscontrols.go`), drop the
   `//go:build skeletons` tag.
3. Port the body — consult the upstream Java source at
   `.cache/closure-src/…` (already extracted if you've run
   `make oracle`) and translate the logic. Keep
   method/field names recognisable to Java readers per SPEC.md.
4. Add the discovered `JSC_*` codes to `internal/diag/diag.go`.
5. Tighten the `// Port of` qualifier from `(skeleton)` to
   `(subset)` / `(narrow; …)` once the method bodies land.
6. Wire the pass into `compiler.go`'s `CheckSourceWithOptions`.

## Regenerating

The skeletons are generated from the curated target list in
`scripts/gen-skeletons.sh`. Re-run to refresh:

```sh
bash scripts/gen-skeletons.sh
```

The tool rewrites each file in place — any hand-edits are lost, so
move a skeleton out of `skel/` before editing it seriously.

## Build tag

Skeletons are gated behind `//go:build skeletons` so the main build
stays clean. They won't compile as-is (unmapped types cascade across
the set) even with the tag set — they're starting points for manual
translation, not a parallel implementation.
