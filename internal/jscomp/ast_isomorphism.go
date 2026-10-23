package jscomp

import (
	"fmt"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type noopErrorReporter struct{}

func (noopErrorReporter) Error(string, string, int, int)   {}
func (noopErrorReporter) Warning(string, string, int, int) {}

// AstIsomorphismStats records structural statistics verified during a lockstep
// AST alpha-isomorphism check.
type AstIsomorphismStats struct {
	TotalNodes      int `json:"total_nodes"`
	Functions       int `json:"functions"`
	Calls           int `json:"calls"`
	Returns         int `json:"returns"`
	Throws          int `json:"throws"`
	ControlFlow     int `json:"control_flow"`
	StringLiterals  int `json:"string_literals"`
	NumberLiterals  int `json:"number_literals"`
	ExternNames     int `json:"extern_names"`
	BoundVarRefs    int `json:"bound_var_refs"`
	UniqueBoundVars int `json:"unique_bound_vars"`
	RenamedVarsDiff int `json:"renamed_vars_diff"`
}

// VerifyAstAlphaIsomorphism parses two compiled JS strings using jsgo's parser,
// constructs exact lexical scopes via SyntacticScopeCreator, and walks both ASTs
// in strict lockstep.
//
// It verifies that:
//  1. Every AST node has the exact same Token, child count, boolean flags, and
//     literal/property values (NUMBER, BIGINT, STRINGLIT, REGEXP, GETPROP,
//     OPTCHAIN_GETPROP, STRING_KEY, MEMBER_FUNCTION_DEF, GETTER_DEF, SETTER_DEF).
//  2. Every undeclared global/extern identifier (scope.GetVar(name) == nil) has
//     the exact same name (nameA == nameB).
//  3. Every declared variable (scope.GetVar(name) != nil) resolves to the exact
//     corresponding lexical Scope in both ASTs and forms a strict 1-to-1
//     bijection (varA <-> varB) across its declaration and all references.
//  4. Every LABEL_NAME resolves to the exact same enclosing LABEL depth.
func VerifyAstAlphaIsomorphism(codeA, codeB string) (*AstIsomorphismStats, error) {
	rootA, compA, err := parseForIsomorphism("compat.js", codeA)
	if err != nil {
		return nil, fmt.Errorf("failed to parse codeA: %w", err)
	}
	if codeA == codeB {
		stats := &AstIsomorphismStats{}
		countAstStatsOnly(rootA, stats)
		return stats, nil
	}
	rootB, compB, err := parseForIsomorphism("fast.js", codeB)
	if err != nil {
		return nil, fmt.Errorf("failed to parse codeB: %w", err)
	}

	v := &astIsomorphismVerifier{
		scA:       NewSyntacticScopeCreator(compA),
		scB:       NewSyntacticScopeCreator(compB),
		scopeAtoB: make(map[*Scope]*Scope, 64),
		varAtoB:   make(map[*Var]*Var, 256),
		varBtoA:   make(map[*Var]*Var, 256),
	}
	if err := v.walk(rootA, rootB); err != nil {
		return nil, err
	}
	v.stats.UniqueBoundVars = len(v.varAtoB)
	return &v.stats, nil
}

func parseForIsomorphism(name, src string) (*rhino.Node, *Compiler, error) {
	comp := NewCompiler()
	opts := comp.NewCompilerOptions()
	opts.SetLanguage(CompilerOptions_LanguageMode_UNSUPPORTED)
	comp.InitOptions(opts)
	cfg := parsing.ParserRunner_CreateConfigWithModes(
		parsing.Config_LanguageMode_ES_NEXT,
		parsing.Config_JsDocParsing_TYPES_ONLY,
		parsing.Config_RunMode_STOP_AFTER_ERROR,
		nil,
		true,
		parsing.Config_StrictMode_SLOPPY,
	)
	file := SourceFile_FromCode(name, src)
	var res *parsing.ParserRunner_ParseResult
	var caught jutil.Throwable
	jutil.Try(func() {
		res = parsing.ParserRunner_Parse(file, src, cfg, noopErrorReporter{})
	}, func(e jutil.Throwable) {
		caught = e
	})
	if caught != nil {
		return nil, nil, fmt.Errorf("parser exception: %s", jutil.ToJavaString(caught))
	}
	if res == nil || res.Ast == nil {
		return nil, nil, fmt.Errorf("parser returned nil AST")
	}
	script := res.Ast
	script.SetInputId(rhino.NewInputId(name))
	root := rhino.IR_Root(script)
	comp.putCompilerInput(NewCompilerInput(file))
	return root, comp, nil
}

func countAstStatsOnly(n *rhino.Node, s *AstIsomorphismStats) {
	s.TotalNodes++
	switch n.GetToken() {
	case rhino.Token_FUNCTION:
		s.Functions++
	case rhino.Token_CALL, rhino.Token_OPTCHAIN_CALL, rhino.Token_NEW:
		s.Calls++
	case rhino.Token_RETURN:
		s.Returns++
	case rhino.Token_THROW:
		s.Throws++
	case rhino.Token_IF, rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF,
		rhino.Token_WHILE, rhino.Token_DO, rhino.Token_SWITCH, rhino.Token_CASE, rhino.Token_TRY, rhino.Token_CATCH:
		s.ControlFlow++
	case rhino.Token_STRINGLIT, rhino.Token_TEMPLATELIT_STRING:
		s.StringLiterals++
	case rhino.Token_NUMBER, rhino.Token_BIGINT:
		s.NumberLiterals++
	}
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		countAstStatsOnly(c, s)
	}
}

type labelPair struct {
	nameA string
	nameB string
}

type astIsomorphismVerifier struct {
	scA       *SyntacticScopeCreator
	scB       *SyntacticScopeCreator
	scopeA    *Scope
	scopeB    *Scope
	scopeAtoB map[*Scope]*Scope
	varAtoB   map[*Var]*Var
	varBtoA   map[*Var]*Var
	labels    []labelPair
	stats     AstIsomorphismStats
}

func (v *astIsomorphismVerifier) pushScope(nA, nB *rhino.Node) {
	v.scopeA = v.scA.CreateScope(nA, v.scopeA).(*Scope)
	v.scopeB = v.scB.CreateScope(nB, v.scopeB).(*Scope)
	v.scopeAtoB[v.scopeA] = v.scopeB
}

func (v *astIsomorphismVerifier) popScope() {
	v.scopeA = v.scopeA.parent
	v.scopeB = v.scopeB.parent
}

func (v *astIsomorphismVerifier) walk(nA, nB *rhino.Node) error {
	v.stats.TotalNodes++
	tokA := nA.GetToken()
	tokB := nB.GetToken()
	if tokA != tokB {
		return fmt.Errorf(
			"token mismatch at compat %d:%d (%v) vs fast %d:%d (%v)",
			nA.GetLineno(), nA.GetCharno(), tokA,
			nB.GetLineno(), nB.GetCharno(), tokB,
		)
	}

	cntA := nA.GetChildCount()
	cntB := nB.GetChildCount()
	if cntA != cntB {
		return fmt.Errorf(
			"child count mismatch on %v at compat %d:%d (%d children) vs fast %d:%d (%d children)",
			tokA, nA.GetLineno(), nA.GetCharno(), cntA,
			nB.GetLineno(), nB.GetCharno(), cntB,
		)
	}

	// Verify node-level boolean flags that affect semantics.
	if nA.IsArrowFunction() != nB.IsArrowFunction() ||
		nA.IsAsyncFunction() != nB.IsAsyncFunction() ||
		nA.IsGeneratorFunction() != nB.IsGeneratorFunction() ||
		nA.IsStaticMember() != nB.IsStaticMember() ||
		nA.IsQuotedStringKey() != nB.IsQuotedStringKey() {
		return fmt.Errorf(
			"boolean property flag mismatch on %v at compat %d:%d vs fast %d:%d",
			tokA, nA.GetLineno(), nA.GetCharno(), nB.GetLineno(), nB.GetCharno(),
		)
	}

	switch tokA {
	case rhino.Token_FUNCTION:
		v.stats.Functions++
		// Matches NodeTraversal.TraverseFunction exactly:
		// fnName is resolved in the enclosing scope (for named function declarations)
		// or the function scope (for function expressions).
		fnNameA := nA.GetFirstChild()
		argsA := fnNameA.GetNext()
		bodyA := argsA.GetNext()

		fnNameB := nB.GetFirstChild()
		argsB := fnNameB.GetNext()
		bodyB := argsB.GetNext()

		isFnExprA := NodeUtil_IsFunctionExpression(nA)
		isFnExprB := NodeUtil_IsFunctionExpression(nB)
		if isFnExprA != isFnExprB {
			return fmt.Errorf("function expression vs declaration mismatch at %d:%d", nA.GetLineno(), nA.GetCharno())
		}

		if !isFnExprA {
			if err := v.walk(fnNameA, fnNameB); err != nil {
				return err
			}
		}

		v.pushScope(nA, nB)
		if isFnExprA {
			if err := v.walk(fnNameA, fnNameB); err != nil {
				return err
			}
		}
		if err := v.walk(argsA, argsB); err != nil {
			return err
		}
		if err := v.walk(bodyA, bodyB); err != nil {
			return err
		}
		v.popScope()
		return nil

	case rhino.Token_CALL, rhino.Token_OPTCHAIN_CALL, rhino.Token_NEW:
		v.stats.Calls++
	case rhino.Token_RETURN:
		v.stats.Returns++
	case rhino.Token_THROW:
		v.stats.Throws++
	case rhino.Token_IF, rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF,
		rhino.Token_WHILE, rhino.Token_DO, rhino.Token_SWITCH, rhino.Token_CASE, rhino.Token_TRY, rhino.Token_CATCH:
		v.stats.ControlFlow++

	case rhino.Token_NUMBER:
		v.stats.NumberLiterals++
		dA, dB := nA.GetDouble(), nB.GetDouble()
		if dA != dB && !(dA != dA && dB != dB) {
			return fmt.Errorf(
				"NUMBER literal mismatch at compat %d:%d (%v) vs fast %d:%d (%v)",
				nA.GetLineno(), nA.GetCharno(), dA, nB.GetLineno(), nB.GetCharno(), dB,
			)
		}

	case rhino.Token_BIGINT:
		v.stats.NumberLiterals++
		if nA.GetBigInt().Cmp(nB.GetBigInt()) != 0 {
			return fmt.Errorf(
				"BIGINT literal mismatch at compat %d:%d (%v) vs fast %d:%d (%v)",
				nA.GetLineno(), nA.GetCharno(), nA.GetBigInt(), nB.GetLineno(), nB.GetCharno(), nB.GetBigInt(),
			)
		}

	case rhino.Token_STRINGLIT, rhino.Token_TEMPLATELIT_STRING:
		v.stats.StringLiterals++
		if nA.GetString() != nB.GetString() {
			return fmt.Errorf(
				"%v literal mismatch at compat %d:%d (%q) vs fast %d:%d (%q)",
				tokA, nA.GetLineno(), nA.GetCharno(), nA.GetString(),
				nB.GetLineno(), nB.GetCharno(), nB.GetString(),
			)
		}

	case rhino.Token_GETPROP, rhino.Token_OPTCHAIN_GETPROP, rhino.Token_STRING_KEY,
		rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF,
		rhino.Token_MEMBER_FIELD_DEF:
		if nA.GetString() != nB.GetString() {
			return fmt.Errorf(
				"%v property name mismatch at compat %d:%d (%q) vs fast %d:%d (%q)",
				tokA, nA.GetLineno(), nA.GetCharno(), nA.GetString(),
				nB.GetLineno(), nB.GetCharno(), nB.GetString(),
			)
		}

	case rhino.Token_LABEL:
		lblA := nA.GetFirstChild()
		lblB := nB.GetFirstChild()
		v.labels = append(v.labels, labelPair{nameA: lblA.GetString(), nameB: lblB.GetString()})
		err := v.walkChildren(nA, nB)
		v.labels = v.labels[:len(v.labels)-1]
		return err

	case rhino.Token_LABEL_NAME:
		sA, sB := nA.GetString(), nB.GetString()
		idxA, idxB := -1, -1
		for i := len(v.labels) - 1; i >= 0; i-- {
			if idxA == -1 && v.labels[i].nameA == sA {
				idxA = i
			}
			if idxB == -1 && v.labels[i].nameB == sB {
				idxB = i
			}
		}
		if idxA != idxB || idxA == -1 {
			return fmt.Errorf(
				"LABEL_NAME binding mismatch at compat %d:%d (%q -> depth %d) vs fast %d:%d (%q -> depth %d)",
				nA.GetLineno(), nA.GetCharno(), sA, idxA,
				nB.GetLineno(), nB.GetCharno(), sB, idxB,
			)
		}

	case rhino.Token_NAME:
		if err := v.verifyNameNode(nA, nB); err != nil {
			return err
		}
	}

	createsScope := nA.IsRoot() || NodeUtil_CreatesBlockScope(nA)
	if createsScope {
		v.pushScope(nA, nB)
	}
	err := v.walkChildren(nA, nB)
	if createsScope {
		v.popScope()
	}
	return err
}

func (v *astIsomorphismVerifier) walkChildren(nA, nB *rhino.Node) error {
	cA := nA.GetFirstChild()
	cB := nB.GetFirstChild()
	for cA != nil && cB != nil {
		if err := v.walk(cA, cB); err != nil {
			return err
		}
		cA = cA.GetNext()
		cB = cB.GetNext()
	}
	return nil
}

func (v *astIsomorphismVerifier) verifyNameNode(nA, nB *rhino.Node) error {
	nameA := nA.GetString()
	nameB := nB.GetString()
	if nameA == "" || nameB == "" {
		if nameA != nameB {
			return fmt.Errorf(
				"empty vs non-empty NAME mismatch at compat %d:%d (%q) vs fast %d:%d (%q)",
				nA.GetLineno(), nA.GetCharno(), nameA, nB.GetLineno(), nB.GetCharno(), nameB,
			)
		}
		return nil
	}

	varA := v.scopeA.GetVar(nameA)
	varB := v.scopeB.GetVar(nameB)
	if varA == nil || varB == nil {
		if varA != nil || varB != nil {
			return fmt.Errorf(
				"declared vs undeclared NAME mismatch at compat %d:%d (%q, declared=%v) vs fast %d:%d (%q, declared=%v)",
				nA.GetLineno(), nA.GetCharno(), nameA, varA != nil,
				nB.GetLineno(), nB.GetCharno(), nameB, varB != nil,
			)
		}
		if nameA != nameB {
			return fmt.Errorf(
				"undeclared global/extern NAME mismatch at compat %d:%d (%q) vs fast %d:%d (%q)",
				nA.GetLineno(), nA.GetCharno(), nameA, nB.GetLineno(), nB.GetCharno(), nameB,
			)
		}
		v.stats.ExternNames++
		return nil
	}

	v.stats.BoundVarRefs++
	if v.scopeAtoB[varA.scope] != varB.scope {
		return fmt.Errorf(
			"lexical scope depth mismatch for NAME at compat %d:%d (%q in scope depth %d) vs fast %d:%d (%q in scope depth %d)",
			nA.GetLineno(), nA.GetCharno(), nameA, varA.scope.depth,
			nB.GetLineno(), nB.GetCharno(), nameB, varB.scope.depth,
		)
	}

	mappedB, okA := v.varAtoB[varA]
	mappedA, okB := v.varBtoA[varB]
	if !okA && !okB {
		v.varAtoB[varA] = varB
		v.varBtoA[varB] = varA
		if nameA != nameB {
			v.stats.RenamedVarsDiff++
		}
		return nil
	}
	if mappedB != varB || mappedA != varA {
		expectedB := "<unmapped>"
		if mappedB != nil {
			expectedB = mappedB.GetName()
		}
		return fmt.Errorf(
			"variable bijection violation at compat %d:%d (%q) vs fast %d:%d (%q): compat %q was already bound to fast %q",
			nA.GetLineno(), nA.GetCharno(), nameA,
			nB.GetLineno(), nB.GetCharno(), nameB,
			nameA, expectedB,
		)
	}
	return nil
}
