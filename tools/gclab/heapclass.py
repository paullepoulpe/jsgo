#!/usr/bin/env python3
"""Rolls heapana rows (heap profile flattened by site) up into components.

  python3 -I heapclass.py ROWS.tsv [--index inuse_space|alloc_space|alloc_objects|inuse_objects] [--top N] [--by-pass]
"""
import argparse, collections, csv, re, sys

STRING_KINDS = ("runtime.concatstring", "runtime.rawstring", "runtime.slicebytetostring", "runtime.slicerunetostring",
                "runtime.intstring", "runtime.stringtoslice", "runtime.rawbyteslice", "runtime.rawruneslice")

# (component, regex on "site|caller1|caller2"): first match wins. Sites are short function names
# (package.Func or package.(*Type).Method).
RULES = [
    ("AST: property lists", r"^rhino\.(NewNode_\w*PropListItem|\(\*Node\)\.(PutProp|PutBooleanProp|PutIntProp|SetProp|CreateProp|PutObjectProp|AddProp|SetJSDocInfo|SetIsAddedBlock|SetSideEffectFlags|SetDeclaredConstantVar|SetInferredConstantVar|SetLength|Set\w+)\b)"),
    ("AST: node_ext", r"^rhino\.\(\*Node\)\.node_getExt"),
    ("AST: nodes", r"^rhino\.(New\w*Node\w*|new[Nn]ode\w*|NewNode|IR_\w+|\(\*Node\)\.\w+|Node_\w+|\(\*RhinoStringPool\w*\)\.\w+)"),
    ("JSDoc & source info", r"^(rhino\.(\(\*)?(JSDocInfo|JSTypeExpression|NonJSDocComment|SimpleSourceFile|InputId|SourcePosition)\w*|rhino\.(New)?(JSDocInfo|JSTypeExpression|NonJSDocComment|SimpleSourceFile|InputId)\w*|parsing\.\(\*JsDocInfoParser\)|parsing\.(New)?JsDoc\w*|parsing\.\(\*JsDocTokenStream\)|jscomp\.(New)?SourceFile\w*|jscomp\.\(\*SourceFile\)|jscomp\.\(\*CompilerInput\)|jscomp\.NewCompilerInput\w*|jscomp\.\(\*JsAst\)|sourcemap\.)"),
    ("parser (transient trees, tokens)", r"^(parser\.|parsing\.)"),
    ("types & colors", r"^(rhino\.\(\*?\w*(Type|Types|Registry|PropertyMap|Property|Slot|TemplateTypeMap|TemplateTypeReplacer|TypeVisitor|Builder|Checker|Visitor)\w*\)|rhino\.(New|new)\w*(Type|Types|Registry|PropertyMap|Property|TemplateTypeMap|TemplateTypeReplacer|Builder|Checker)\w*|rhino\.\w*JSType\w*|rhino\.Nominal\w*|rhino\.\w*TypeI?\w*_|colors\.|jscomp\.\(\*(TypeInference|TypedScopeCreator|FunctionTypeBuilder|TypeValidator|TypeCheck|JSTypeResolver|ColorPool|ColorSerializer|ColorGraph\w*|DisambiguateProperties|AmbiguateProperties|TypeMismatch|InferJSDocInfo|ClosureReverseAbstractInterpreter|SemanticReverseAbstractInterpreter|ChainableReverseAbstractInterpreter|TypedCodeGenerator|TypeTransformation|TypeInferencePass|TypedScopeCreator_\w+|ConvertTypesToColors\w*|ColorPool_\w+|NativeColorId\w*|PropertyClustering\w*|FixedPointGraphTraversal\w*)\)|jscomp\.(New)?(TypeInference|TypedScopeCreator|FunctionTypeBuilder|TypeValidator|TypeCheck|ColorPool|ConvertTypesToColors|DisambiguateProperties|AmbiguateProperties|TypeMismatch|ClosureReverseAbstractInterpreter|SemanticReverseAbstractInterpreter|PropertyClustering|ColorGraph|InferJSDocInfo)\w*)"),
    ("CFG & data flow", r"^jscomp\.(\(\*)?(New)?(ControlFlowAnalysis|ControlFlowGraph|LinkedDirectedGraph|DiGraph|Graph|DataFlowAnalysis|LinkedFlowScope|LiveVariablesAnalysis|LiveVariableLattice|MaybeReachingVariableUse|MustBeReachingVariableDef|FlowSensitiveInlineVariables|DeadAssignmentsElimination|CoalesceVariableNames|GraphColoring|LinkedUndirectedGraph|UndirectedGraph|BasicBlock|ReachingUses|MustDef|FixedPointGraph|LatticeElement|JoinOp|BinaryJoinOp|Annotat|CheckPathsBetweenNodes)\w*|^rhino\.(\(\*)?(Hamt|PMap|hamt)\w*|^jutil\.(\(\*)?(New)?BitSet\w*"),
    ("scopes, vars, references", r"^jscomp\.(\(\*)?(New)?(Abstract\w*Scope\w*|AbstractScope\w*|AbstractVar\w*|Scope|TypedScope|Var|TypedVar|SyntacticScopeCreator|ScopeCreator|Es6SyntacticScopeCreator|ReferenceCollector|Reference|ReferenceMap|ReferenceCollection|ReferenceCollectingCallback|GlobalNamespace|GlobalNamespace_\w+|NameBasedDefinitionProvider|RemoveUnusedCode_\w*VarInfo|RemoveUnusedCode_\w*Scope\w*|InlineVariables_\w*|Normalize\w*Scope|MakeDeclaredNamesUnique\w*|VarCheck|ScopedCallback|NodeTraversal\w*|ChangeVerifier\w*|AbstractScopeBase\w*)[\w\[\],.* ]*\)?\.?\w*"),
    ("jutil collections", r"^jutil\.(\(\*)?(New|new)?(ArrayList|LinkedHashMap|HashMap|LinkedHashSet|HashSet|ObjectLinkedHashMap|ObjectHashMap|ObjectHashSet|ArrayDeque|ArrayListMultimap|LinkedHashMultimap|HashMultimap|Multimap|Multiset|ImmutableList|ImmutableMap|ImmutableSet|ImmutableSortedSet|TreeMap|TreeSet|IdentityHashMap|PriorityQueue|Lists|Maps|Sets|Collections|Iterables|BiMap|HashBiMap|EnumSet|EnumMap|lhm|olhm|SetMultimap|ListMultimap|Table|LinkedList|Stack|Deque|Queue|Iterator|SortedSet|WeakHashMap|ObjectArrayList|Arrays|CopyOf)\w*"),
    ("strings", r"^(jutil\.(\(\*)?(String\w*|StringBuilder\w*|Char\w*)|strings\.|strconv\.|unicode/utf16\.|unicode/utf8\.|bytes\.|fmt\.)"),
    ("module/chunk/compiler state", r"^jscomp\.(\(\*)?(New)?(ModuleMetadataMap|ModuleMap|ModuleLoader|JSChunk\w*|Compiler|CompilerOptions|CommandLineRunner|AbstractCommandLineRunner|DependencyOptions|Es6RewriteModules|ClosureRewriteModule|ModuleRenaming|GatherModuleMetadata\w*|ModuleMetadata\w*)\w*"),
    ("code printing & source maps", r"^jscomp\.(\(\*)?(New)?(CodeGenerator|CodePrinter\w*|CodeConsumer|CodePrinter_\w+|SourceMap\w*|LineNumberCheck)\w*"),
]
RULES = [(c, re.compile(r)) for c, r in RULES]


def component(kind, site, c1, c2):
    if kind.startswith(STRING_KINDS) or (kind.startswith("runtime.") and "string" in kind.lower()):
        return "strings"
    for c, r in RULES:
        if r.match(site):
            return c
    if site.startswith("jscomp."):
        return "other (pass-local)"
    return "other"


def load(path):
    with open(path) as f:
        rd = csv.reader(f, delimiter="\t")
        hdr = next(rd)
        for row in rd:
            d = dict(zip(hdr, row))
            for k in hdr[6:]:
                d[k] = int(d[k])
            yield d


IDX = {"inuse_space": "inuse_space_bytes", "alloc_space": "alloc_space_bytes",
       "inuse_objects": "inuse_objects_count", "alloc_objects": "alloc_objects_count"}


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("rows")
    ap.add_argument("--index", default="inuse_space")
    ap.add_argument("--top", type=int, default=40)
    ap.add_argument("--by-pass", action="store_true")
    ap.add_argument("--component")
    a = ap.parse_args()
    k = IDX[a.index]
    rows = list(load(a.rows))
    tot = sum(r[k] for r in rows)
    comp = collections.Counter()
    site = collections.Counter()
    for r in rows:
        c = component(r["kind"], r["site"], r["caller1"], r["caller2"])
        r["component"] = c
        comp[c] += r[k]
        site[(c, r["kind"], r["site"], r["caller1"])] += r[k]
    unit = 2**20 if "space" in a.index else 1e6
    uname = "MiB" if "space" in a.index else "M"
    print(f"# {a.index}: total {tot / unit:.1f} {uname}")
    print("component\t" + uname + "\tshare")
    for c, v in comp.most_common():
        print(f"{c}\t{v / unit:.1f}\t{100 * v / tot:.1f}%")
    print(f"\n# top sites by {a.index}\n{uname}\tshare\tcomponent\tkind\tsite\tcaller1")
    for (c, kd, s, c1), v in site.most_common(a.top):
        if a.component and c != a.component:
            continue
        print(f"{v / unit:.2f}\t{100 * v / tot:.1f}%\t{c}\t{kd}\t{s}\t{c1}")
    if a.by_pass:
        bp = collections.Counter()
        bpc = collections.defaultdict(collections.Counter)
        for r in rows:
            bp[r["pass"]] += r[k]
            bpc[r["pass"]][r["component"]] += r[k]
        print(f"\n# by pass ({a.index})\npass\t{uname}\tshare\ttop components")
        for p, v in bp.most_common(30):
            tops = ", ".join(f"{c} {100 * x / v:.0f}%" for c, x in bpc[p].most_common(3))
            print(f"{p}\t{v / unit:.1f}\t{100 * v / tot:.1f}%\t{tops}")


if __name__ == "__main__":
    main()
