#!/usr/bin/env bash
# Bulk-generate skeletons for all checker-relevant closure-compiler
# Java sources. Output lands under internal/jscomp/skel/ with a
# `//go:build skeletons` tag so the files stay out of the default
# build — the main checker pipeline isn't affected. Re-run is
# idempotent: the tool rewrites each skeleton in place.
#
# Usage:
#   bash scripts/gen-skeletons.sh
#
# Requires:
#   .cache/closure-src/  (extracted closure-compiler sources)

set -eu
repo_root="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$repo_root"

SRC=".cache/closure-src"
OUT="internal/jscomp/skel"
TAG="skeletons"

if [[ ! -d "$SRC" ]]; then
    echo "error: $SRC not found — run the oracle regen once to populate it"
    exit 2
fi

# Curated list of type-checker-adjacent files. Keeping this explicit
# rather than a wildcard so the skel tree doesn't drift into
# transpilation / optimization passes that aren't part of the
# checker's oracle.
declare -a targets=(
    # rhino AST + type lattice
    com/google/javascript/rhino/ClosurePrimitive.java
    com/google/javascript/rhino/ErrorReporter.java
    com/google/javascript/rhino/IR.java
    com/google/javascript/rhino/InputId.java
    com/google/javascript/rhino/JSDocInfo.java
    com/google/javascript/rhino/JSTypeExpression.java
    com/google/javascript/rhino/Node.java
    com/google/javascript/rhino/NominalTypeBuilder.java
    com/google/javascript/rhino/NonJSDocComment.java
    com/google/javascript/rhino/Outcome.java
    com/google/javascript/rhino/QualifiedName.java
    com/google/javascript/rhino/SimpleSourceFile.java
    com/google/javascript/rhino/SourcePosition.java
    com/google/javascript/rhino/StaticRef.java
    com/google/javascript/rhino/StaticScope.java
    com/google/javascript/rhino/StaticSlot.java
    com/google/javascript/rhino/StaticSourceFile.java
    com/google/javascript/rhino/StaticSymbolTable.java
    com/google/javascript/rhino/Token.java
    com/google/javascript/rhino/TokenStream.java
    com/google/javascript/rhino/TokenUtil.java
    com/google/javascript/rhino/TypeDeclarationsIR.java

    # rhino/jstype — the full lattice
    com/google/javascript/rhino/jstype/AllType.java
    com/google/javascript/rhino/jstype/ArrowType.java
    com/google/javascript/rhino/jstype/BigIntType.java
    com/google/javascript/rhino/jstype/BooleanLiteralSet.java
    com/google/javascript/rhino/jstype/BooleanType.java
    com/google/javascript/rhino/jstype/CanCastToVisitor.java
    com/google/javascript/rhino/jstype/ContainsUpperBoundSuperTypeVisitor.java
    com/google/javascript/rhino/jstype/EnumElementType.java
    com/google/javascript/rhino/jstype/EnumType.java
    com/google/javascript/rhino/jstype/EqualityChecker.java
    com/google/javascript/rhino/jstype/EquivalenceMethod.java
    com/google/javascript/rhino/jstype/FunctionParamBuilder.java
    com/google/javascript/rhino/jstype/FunctionType.java
    com/google/javascript/rhino/jstype/InstanceObjectType.java
    com/google/javascript/rhino/jstype/JSType.java
    com/google/javascript/rhino/jstype/JSTypeClass.java
    com/google/javascript/rhino/jstype/JSTypeIterations.java
    com/google/javascript/rhino/jstype/JSTypeNative.java
    com/google/javascript/rhino/jstype/JSTypeRegistry.java
    com/google/javascript/rhino/jstype/JSTypeResolver.java
    com/google/javascript/rhino/jstype/KnownSymbolType.java
    com/google/javascript/rhino/jstype/NamedType.java
    com/google/javascript/rhino/jstype/NoObjectType.java
    com/google/javascript/rhino/jstype/NoResolvedType.java
    com/google/javascript/rhino/jstype/NoType.java
    com/google/javascript/rhino/jstype/NullType.java
    com/google/javascript/rhino/jstype/NumberType.java
    com/google/javascript/rhino/jstype/ObjectType.java
    com/google/javascript/rhino/jstype/Property.java
    com/google/javascript/rhino/jstype/PropertyMap.java
    com/google/javascript/rhino/jstype/PrototypeObjectType.java
    com/google/javascript/rhino/jstype/ProxyObjectType.java
    com/google/javascript/rhino/jstype/RecordType.java
    com/google/javascript/rhino/jstype/RecordTypeBuilder.java
    com/google/javascript/rhino/jstype/RelationshipVisitor.java
    com/google/javascript/rhino/jstype/StringType.java
    com/google/javascript/rhino/jstype/SubtypeChecker.java
    com/google/javascript/rhino/jstype/SymbolType.java
    com/google/javascript/rhino/jstype/TemplateType.java
    com/google/javascript/rhino/jstype/TemplateTypeMap.java
    com/google/javascript/rhino/jstype/TemplateTypeReplacer.java
    com/google/javascript/rhino/jstype/TemplatizedType.java
    com/google/javascript/rhino/jstype/TypeStringBuilder.java
    com/google/javascript/rhino/jstype/UnionType.java
    com/google/javascript/rhino/jstype/UnknownType.java
    com/google/javascript/rhino/jstype/ValueType.java
    com/google/javascript/rhino/jstype/Visitor.java
    com/google/javascript/rhino/jstype/VoidType.java

    # jscomp/type — FlowScope layer
    com/google/javascript/jscomp/type/FlowScope.java
    com/google/javascript/jscomp/type/ReverseAbstractInterpreter.java
    com/google/javascript/jscomp/type/ChainableReverseAbstractInterpreter.java
    com/google/javascript/jscomp/type/ClosureReverseAbstractInterpreter.java
    com/google/javascript/jscomp/type/SemanticReverseAbstractInterpreter.java

    # jscomp/ top-level: Check-* passes
    com/google/javascript/jscomp/CheckAccessControls.java
    com/google/javascript/jscomp/CheckClosureImports.java
    com/google/javascript/jscomp/CheckConformance.java
    com/google/javascript/jscomp/CheckDebuggerStatement.java
    com/google/javascript/jscomp/CheckGlobalThis.java
    com/google/javascript/jscomp/CheckJSDoc.java
    com/google/javascript/jscomp/CheckLevel.java
    com/google/javascript/jscomp/CheckMissingRequires.java
    com/google/javascript/jscomp/CheckMissingReturn.java
    com/google/javascript/jscomp/CheckRegExp.java
    com/google/javascript/jscomp/CheckSideEffects.java
    com/google/javascript/jscomp/CheckSuper.java
    com/google/javascript/jscomp/CheckSuspiciousCode.java
    com/google/javascript/jscomp/CheckTypeImportCodeReferences.java
    com/google/javascript/jscomp/CheckUnreachableCode.java

    # jscomp/ Type / Scope core
    com/google/javascript/jscomp/TypeCheck.java
    com/google/javascript/jscomp/TypeInference.java
    com/google/javascript/jscomp/TypeInferencePass.java
    com/google/javascript/jscomp/TypeMatchingStrategy.java
    com/google/javascript/jscomp/TypeMismatch.java
    com/google/javascript/jscomp/TypeTransformation.java
    com/google/javascript/jscomp/TypeValidator.java
    com/google/javascript/jscomp/TypedScope.java
    com/google/javascript/jscomp/TypedScopeCreator.java
    com/google/javascript/jscomp/TypedVar.java
    com/google/javascript/jscomp/InferJSDocInfo.java
    com/google/javascript/jscomp/InferConsts.java

    # jscomp/ Scope / Var machinery
    com/google/javascript/jscomp/AbstractScope.java
    com/google/javascript/jscomp/AbstractVar.java
    com/google/javascript/jscomp/Scope.java
    com/google/javascript/jscomp/ScopeCreator.java
    com/google/javascript/jscomp/ScopedName.java
    com/google/javascript/jscomp/SymbolTable.java
    com/google/javascript/jscomp/Var.java
    com/google/javascript/jscomp/VarCheck.java
    com/google/javascript/jscomp/VariableReferenceCheck.java

    # jscomp/ AST + traversal helpers
    com/google/javascript/jscomp/AstAnalyzer.java
    com/google/javascript/jscomp/AstFactory.java
    com/google/javascript/jscomp/AstManipulations.java
    com/google/javascript/jscomp/AstValidator.java
    com/google/javascript/jscomp/NodeIterators.java
    com/google/javascript/jscomp/NodeNameExtractor.java
    com/google/javascript/jscomp/NodeTraversal.java
    com/google/javascript/jscomp/NodeUtil.java

    # jscomp/ Diagnostic machinery
    com/google/javascript/jscomp/DiagnosticGroup.java
    com/google/javascript/jscomp/DiagnosticGroupPathSuppressingWarningsGuard.java
    com/google/javascript/jscomp/DiagnosticGroupWarningsGuard.java
    com/google/javascript/jscomp/DiagnosticGroups.java
    com/google/javascript/jscomp/DiagnosticToSuppressionMapper.java
    com/google/javascript/jscomp/DiagnosticType.java
    com/google/javascript/jscomp/JSError.java
    com/google/javascript/jscomp/AbstractMessageFormatter.java
    com/google/javascript/jscomp/AllowlistWarningsGuard.java

    # jscomp/ Control-flow / data-flow
    com/google/javascript/jscomp/ControlFlowAnalysis.java
    com/google/javascript/jscomp/ControlFlowGraph.java
    com/google/javascript/jscomp/DataFlowAnalysis.java
    com/google/javascript/jscomp/BasicBlock.java

    # jscomp/ Compiler pipeline
    com/google/javascript/jscomp/AbstractCompiler.java
    com/google/javascript/jscomp/Compiler.java
    com/google/javascript/jscomp/CompilerInput.java
    com/google/javascript/jscomp/CompilerInputProvider.java
    com/google/javascript/jscomp/CompilerOptions.java
    com/google/javascript/jscomp/CompilerOptionsPreprocessor.java
    com/google/javascript/jscomp/CompilerPass.java
    com/google/javascript/jscomp/PassConfig.java
    com/google/javascript/jscomp/DefaultPassConfig.java
    com/google/javascript/jscomp/PassFactory.java
    com/google/javascript/jscomp/PassListBuilder.java
    com/google/javascript/jscomp/PassNames.java
    com/google/javascript/jscomp/DefaultExterns.java

    # jscomp/parsing — parser / JSDoc
    com/google/javascript/jscomp/parsing/JsDocInfoParser.java
    com/google/javascript/jscomp/parsing/JsDocToken.java
    com/google/javascript/jscomp/parsing/JsDocTokenStream.java
    com/google/javascript/jscomp/parsing/JsDocInfoWriter.java
    com/google/javascript/jscomp/parsing/TypeTransformationParser.java
    com/google/javascript/jscomp/parsing/ParserRunner.java
    com/google/javascript/jscomp/parsing/IRFactory.java

    # jscomp/parsing/parser — parser grammar
    com/google/javascript/jscomp/parsing/parser/Parser.java
    com/google/javascript/jscomp/parsing/parser/Scanner.java
    com/google/javascript/jscomp/parsing/parser/IdentifierToken.java
    com/google/javascript/jscomp/parsing/parser/LiteralToken.java
    com/google/javascript/jscomp/parsing/parser/TemplateLiteralToken.java
    com/google/javascript/jscomp/parsing/parser/Token.java
    com/google/javascript/jscomp/parsing/parser/TokenType.java
)

processed=0
skipped=0
for rel in "${targets[@]}"; do
    in="$SRC/$rel"
    if [[ ! -f "$in" ]]; then
        echo "skip (missing): $rel"
        skipped=$((skipped + 1))
        continue
    fi
    # Output path mirrors the upstream layout under internal/jscomp/skel.
    # Rewrite Go-reserved directory names (closure has `jscomp/type/`,
    # but `type` is a Go keyword) so the generated package is legal.
    rel_path="${rel#com/google/javascript/}"
    rel_path="${rel_path/\/type\//\/typeinfer\/}"
    out="$OUT/${rel_path}"
    out="${out%.java}.go"
    # Go lower-case filenames for the Go idiom.
    out_dir="$(dirname "$out")"
    out_base="$(basename "$out")"
    out_base_lc="$(echo "$out_base" | tr '[:upper:]' '[:lower:]')"
    out="$out_dir/$out_base_lc"
    # Package name = leaf directory. e.g. rhino/jstype → jstype.
    pkg="$(basename "$out_dir")"
    mkdir -p "$out_dir"
    go run ./tools/javaport \
        -in "$in" \
        -out "$out" \
        -pkg "$pkg" \
        -build-tag "$TAG"
    processed=$((processed + 1))
done

# ---------------------------------------------------------------------------
# Section A: parsing/parser/trees/*.java → internal/parser/trees/
# 86 AST-node tree classes; package "trees".
# ---------------------------------------------------------------------------
TREES_SRC="$SRC/com/google/javascript/jscomp/parsing/parser/trees"
TREES_OUT="internal/parser/trees"
mkdir -p "$TREES_OUT"
echo ""
echo "--- parser/trees ---"
for java_file in "$TREES_SRC"/*.java; do
    if [[ ! -f "$java_file" ]]; then
        continue
    fi
    base="$(basename "$java_file" .java)"
    # CamelCase → snake_case for Go file names.
    snake="$(echo "$base" | sed 's/\([A-Z]\)/_\1/g' | sed 's/^_//' | tr '[:upper:]' '[:lower:]')"
    out="$TREES_OUT/${snake}.go"
    go run ./tools/javaport \
        -in "$java_file" \
        -out "$out" \
        -pkg "trees" \
        -build-tag "$TAG"
    processed=$((processed + 1))
done

# ---------------------------------------------------------------------------
# Section B: jscomp/parsing/{IRFactory,ParserRunner,TypeTransformationParser}
# → internal/jscomp/parsing/skel/
# ---------------------------------------------------------------------------
PARSING_SRC="$SRC/com/google/javascript/jscomp/parsing"
PARSING_OUT="internal/jscomp/parsing/skel"
mkdir -p "$PARSING_OUT"
echo ""
echo "--- jscomp/parsing/skel ---"
declare -a parsing_targets=(
    IRFactory.java
    ParserRunner.java
    TypeTransformationParser.java
)
for fname in "${parsing_targets[@]}"; do
    in="$PARSING_SRC/$fname"
    if [[ ! -f "$in" ]]; then
        echo "skip (missing): $fname"
        skipped=$((skipped + 1))
        continue
    fi
    base="$(basename "$fname" .java)"
    out_base_lc="$(echo "$base" | tr '[:upper:]' '[:lower:]').go"
    out="$PARSING_OUT/$out_base_lc"
    go run ./tools/javaport \
        -in "$in" \
        -out "$out" \
        -pkg "parsing" \
        -build-tag "$TAG"
    processed=$((processed + 1))
done

# ---------------------------------------------------------------------------
# Section C: parsing/parser/{ClassOrObjectElementInfo,FeatureSet,Identifiers,
#   Keywords,LineNumberScanner,PredefinedName,SourceFile,StringLiteralToken}
#   + parsing/parser/util/*.java → internal/parser/skel/
# ---------------------------------------------------------------------------
PARSER_SRC="$SRC/com/google/javascript/jscomp/parsing/parser"
PARSER_SKEL_OUT="internal/parser/skel"
mkdir -p "$PARSER_SKEL_OUT"
echo ""
echo "--- parser/skel ---"
declare -a parser_skel_targets=(
    ClassOrObjectElementInfo.java
    FeatureSet.java
    Identifiers.java
    Keywords.java
    LineNumberScanner.java
    PredefinedName.java
    SourceFile.java
    StringLiteralToken.java
)
for fname in "${parser_skel_targets[@]}"; do
    in="$PARSER_SRC/$fname"
    if [[ ! -f "$in" ]]; then
        echo "skip (missing): $fname"
        skipped=$((skipped + 1))
        continue
    fi
    base="$(basename "$fname" .java)"
    out_base_lc="$(echo "$base" | tr '[:upper:]' '[:lower:]').go"
    out="$PARSER_SKEL_OUT/$out_base_lc"
    go run ./tools/javaport \
        -in "$in" \
        -out "$out" \
        -pkg "parser" \
        -build-tag "$TAG"
    processed=$((processed + 1))
done
# util/ subdirectory
PARSER_UTIL_SRC="$PARSER_SRC/util"
PARSER_UTIL_OUT="$PARSER_SKEL_OUT/util"
mkdir -p "$PARSER_UTIL_OUT"
for java_file in "$PARSER_UTIL_SRC"/*.java; do
    if [[ ! -f "$java_file" ]]; then
        continue
    fi
    base="$(basename "$java_file" .java)"
    out_base_lc="$(echo "$base" | tr '[:upper:]' '[:lower:]').go"
    out="$PARSER_UTIL_OUT/$out_base_lc"
    go run ./tools/javaport \
        -in "$java_file" \
        -out "$out" \
        -pkg "util" \
        -build-tag "$TAG"
    processed=$((processed + 1))
done

echo ""
echo "generated $processed skeletons, skipped $skipped"
