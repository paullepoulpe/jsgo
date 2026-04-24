// java2go output: CRASH - tool panicked before producing output.
// Error: panic: Unknown node type: line_comment
// Root cause: java2go@v0.0.0-20251207015311-79857fbb7b4c handles "comment" AST node
// type from an older tree-sitter-java grammar version, but the current grammar
// emits "line_comment" and "block_comment" as separate node types. Any Java file
// containing a // single-line comment causes an immediate panic.
