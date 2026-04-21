// Port of com/google/javascript/rhino/Token.java @ closure-compiler v20260415
//
// Only the subset of Token enum values jsgo currently needs is defined
// here — the full Java enum has ~130 entries covering every concrete
// and IR node kind. New entries are added as corpus tiers grow. Names
// and dispatch order mirror Java where possible (ADD / SUB / LT / EQ /
// HOOK / ARRAYLIT / etc.).

package rhino

type Token int

const (
	ERROR Token = iota
	SCRIPT

	// Declarations
	VAR
	LET
	CONST

	// Leaves
	NAME
	STRINGLIT
	NUMBER
	TRUE
	FALSE
	NULL
	UNDEFINED

	// Assignment family — ASSIGN is the bare `=`; the _OP variants are
	// the compound assigns (`+=`, `-=`, …).
	ASSIGN
	ASSIGN_ADD
	ASSIGN_SUB
	ASSIGN_MUL
	ASSIGN_DIV
	ASSIGN_MOD
	ASSIGN_BITAND
	ASSIGN_BITOR
	ASSIGN_BITXOR
	ASSIGN_LSH
	ASSIGN_RSH
	ASSIGN_URSH
	ASSIGN_AND      // &&=
	ASSIGN_OR       // ||=
	ASSIGN_COALESCE // ??=

	// Block / script shells
	EXPR_RESULT
	FUNCTION
	PARAM_LIST
	BLOCK
	RETURN
	IF
	HOOK // ternary `cond ? a : b`

	// Loops / switch / try / jump
	WHILE
	DO_WHILE
	FOR
	FOR_IN
	FOR_OF
	SWITCH
	CASE
	DEFAULT_CASE
	TRY
	CATCH
	FINALLY
	THROW
	BREAK
	CONTINUE
	LABEL

	// Calls / new / member access
	CALL
	NEW
	GETPROP
	GETELEM // computed member access: `x[y]`
	THIS

	// Classes
	CLASS
	CLASS_FIELD_DEF

	// Binary arithmetic
	ADD
	SUB
	MUL
	DIV
	MOD

	// Comparison
	LT
	LE
	GT
	GE
	EQ
	NE
	SHEQ       // ===
	SHNE       // !==
	INSTANCEOF // instanceof
	IN         // in

	// Logical
	AND      // &&
	OR       // ||
	COALESCE // ??

	// Bitwise
	BITAND
	BITOR
	BITXOR
	LSH  // <<
	RSH  // >>
	URSH // >>>

	// Unary
	NOT    // !
	BITNOT // ~
	POS    // unary +
	NEG    // unary -
	TYPEOF
	VOID
	DELPROP // delete

	// Array / object literals
	ARRAYLIT
	OBJECTLIT
	STRING_KEY

	// ---- P1-B parser-breadth additions ---------------------------- //
	// Ported from com/google/javascript/rhino/Token.java tokens used by
	// the Parser.java grammar productions covered in this unit. Only the
	// subset the checker / tree-walking code actually observes is added
	// today; additions strictly extend the enum, never renumber earlier
	// members so existing consumers keep matching the same ordinals.

	// Numeric / string literal variants.
	BIGINT // 1n, 0x1fn — distinct from NUMBER per the ES2020 grammar.
	REGEXP // /foo/g — carries the pattern + flags as two STRINGLIT children.

	// Template literals. Matches Java's Parser.parseTemplateLiteral output
	// shape: TEMPLATELIT(TEMPLATELIT_STRING, [TEMPLATELIT_SUB(expr),
	// TEMPLATELIT_STRING]* ). A tag prefix (x`...`) wraps the whole
	// thing in a TAGGED_TEMPLATELIT.
	TEMPLATELIT
	TEMPLATELIT_STRING
	TEMPLATELIT_SUB
	TAGGED_TEMPLATELIT

	// Destructuring. DESTRUCTURING_LHS wraps the pattern + the RHS
	// initializer so the scope builder can peel the binding sites
	// uniformly (Java uses the same wrapper). ARRAY_PATTERN /
	// OBJECT_PATTERN hold the pattern elements themselves; DEFAULT_VALUE
	// wraps an `x = init` default inside either pattern kind.
	ARRAY_PATTERN
	OBJECT_PATTERN
	DESTRUCTURING_LHS
	DEFAULT_VALUE
	COMPUTED_PROP

	// Spread + rest. ITER_SPREAD appears inside array literals / calls
	// (...iter); OBJECT_SPREAD inside object literals ({...obj}).
	// ITER_REST / OBJECT_REST are used as the final element of a
	// pattern (const [a, ...rest] = x, const {a, ...rest} = y) and the
	// final param of a function (function f(...args)).
	ITER_SPREAD
	OBJECT_SPREAD
	ITER_REST
	OBJECT_REST

	// Optional-chaining variants. Each mirrors its non-optional
	// counterpart but is emitted when the `?.` / `?.(` / `?.[` form
	// was used.
	OPTCHAIN_GETPROP
	OPTCHAIN_GETELEM
	OPTCHAIN_CALL

	// Async / generator. YIELD and AWAIT appear as prefix unary
	// expressions; FOR_AWAIT_OF is the async variant of FOR_OF (for
	// await (x of y) { … }).
	YIELD
	AWAIT
	FOR_AWAIT_OF

	// @type-cast expression. `/** @type {T} */ (expr)` becomes
	// CAST(expr) with the @type JSDoc attached to the CAST node.
	CAST

	// ES6 module-ish bits we recognise but mostly pass through. `super`
	// and `new.target` / `import.meta` get their own tokens so the
	// checker doesn't misclassify them as bare NAMEs.
	SUPER
	NEW_TARGET
	IMPORT_META

	// COMMA is the comma-expression operator (a, b, c) — distinct from
	// the argument-list separator. Appears as a binary node with the
	// two operands as children.
	COMMA

	// Increment / decrement. INC / DEC carry a PROP_INCDEC_POSTFIX bit
	// on the node to distinguish `x++` from `++x` per Java's Node API.
	INC
	DEC

	// Exponent (**) — right-assoc binary operator, precedence above
	// multiplicative. Mirrors Java's Token.EXPONENT.
	EXPONENT
	ASSIGN_EXPONENT

	// Empty node — used as a pattern/param placeholder and for elisions
	// in array literals. Java uses `Token.EMPTY`.
	EMPTY
)

func (t Token) String() string {
	switch t {
	case SCRIPT:
		return "SCRIPT"
	case VAR:
		return "VAR"
	case LET:
		return "LET"
	case CONST:
		return "CONST"
	case NAME:
		return "NAME"
	case STRINGLIT:
		return "STRINGLIT"
	case NUMBER:
		return "NUMBER"
	case TRUE:
		return "TRUE"
	case FALSE:
		return "FALSE"
	case NULL:
		return "NULL"
	case UNDEFINED:
		return "UNDEFINED"
	case ASSIGN:
		return "ASSIGN"
	case ASSIGN_ADD:
		return "ASSIGN_ADD"
	case ASSIGN_SUB:
		return "ASSIGN_SUB"
	case ASSIGN_MUL:
		return "ASSIGN_MUL"
	case ASSIGN_DIV:
		return "ASSIGN_DIV"
	case ASSIGN_MOD:
		return "ASSIGN_MOD"
	case ASSIGN_BITAND:
		return "ASSIGN_BITAND"
	case ASSIGN_BITOR:
		return "ASSIGN_BITOR"
	case ASSIGN_BITXOR:
		return "ASSIGN_BITXOR"
	case ASSIGN_LSH:
		return "ASSIGN_LSH"
	case ASSIGN_RSH:
		return "ASSIGN_RSH"
	case ASSIGN_URSH:
		return "ASSIGN_URSH"
	case ASSIGN_AND:
		return "ASSIGN_AND"
	case ASSIGN_OR:
		return "ASSIGN_OR"
	case ASSIGN_COALESCE:
		return "ASSIGN_COALESCE"
	case EXPR_RESULT:
		return "EXPR_RESULT"
	case FUNCTION:
		return "FUNCTION"
	case PARAM_LIST:
		return "PARAM_LIST"
	case BLOCK:
		return "BLOCK"
	case RETURN:
		return "RETURN"
	case IF:
		return "IF"
	case HOOK:
		return "HOOK"
	case WHILE:
		return "WHILE"
	case DO_WHILE:
		return "DO_WHILE"
	case FOR:
		return "FOR"
	case FOR_IN:
		return "FOR_IN"
	case FOR_OF:
		return "FOR_OF"
	case SWITCH:
		return "SWITCH"
	case CASE:
		return "CASE"
	case DEFAULT_CASE:
		return "DEFAULT_CASE"
	case TRY:
		return "TRY"
	case CATCH:
		return "CATCH"
	case FINALLY:
		return "FINALLY"
	case THROW:
		return "THROW"
	case BREAK:
		return "BREAK"
	case CONTINUE:
		return "CONTINUE"
	case LABEL:
		return "LABEL"
	case CALL:
		return "CALL"
	case NEW:
		return "NEW"
	case GETPROP:
		return "GETPROP"
	case GETELEM:
		return "GETELEM"
	case THIS:
		return "THIS"
	case CLASS:
		return "CLASS"
	case CLASS_FIELD_DEF:
		return "CLASS_FIELD_DEF"
	case ADD:
		return "ADD"
	case SUB:
		return "SUB"
	case MUL:
		return "MUL"
	case DIV:
		return "DIV"
	case MOD:
		return "MOD"
	case LT:
		return "LT"
	case LE:
		return "LE"
	case GT:
		return "GT"
	case GE:
		return "GE"
	case EQ:
		return "EQ"
	case NE:
		return "NE"
	case SHEQ:
		return "SHEQ"
	case SHNE:
		return "SHNE"
	case INSTANCEOF:
		return "INSTANCEOF"
	case IN:
		return "IN"
	case AND:
		return "AND"
	case OR:
		return "OR"
	case COALESCE:
		return "COALESCE"
	case BITAND:
		return "BITAND"
	case BITOR:
		return "BITOR"
	case BITXOR:
		return "BITXOR"
	case LSH:
		return "LSH"
	case RSH:
		return "RSH"
	case URSH:
		return "URSH"
	case NOT:
		return "NOT"
	case BITNOT:
		return "BITNOT"
	case POS:
		return "POS"
	case NEG:
		return "NEG"
	case TYPEOF:
		return "TYPEOF"
	case VOID:
		return "VOID"
	case DELPROP:
		return "DELPROP"
	case ARRAYLIT:
		return "ARRAYLIT"
	case OBJECTLIT:
		return "OBJECTLIT"
	case STRING_KEY:
		return "STRING_KEY"
	case BIGINT:
		return "BIGINT"
	case REGEXP:
		return "REGEXP"
	case TEMPLATELIT:
		return "TEMPLATELIT"
	case TEMPLATELIT_STRING:
		return "TEMPLATELIT_STRING"
	case TEMPLATELIT_SUB:
		return "TEMPLATELIT_SUB"
	case TAGGED_TEMPLATELIT:
		return "TAGGED_TEMPLATELIT"
	case ARRAY_PATTERN:
		return "ARRAY_PATTERN"
	case OBJECT_PATTERN:
		return "OBJECT_PATTERN"
	case DESTRUCTURING_LHS:
		return "DESTRUCTURING_LHS"
	case DEFAULT_VALUE:
		return "DEFAULT_VALUE"
	case COMPUTED_PROP:
		return "COMPUTED_PROP"
	case ITER_SPREAD:
		return "ITER_SPREAD"
	case OBJECT_SPREAD:
		return "OBJECT_SPREAD"
	case ITER_REST:
		return "ITER_REST"
	case OBJECT_REST:
		return "OBJECT_REST"
	case OPTCHAIN_GETPROP:
		return "OPTCHAIN_GETPROP"
	case OPTCHAIN_GETELEM:
		return "OPTCHAIN_GETELEM"
	case OPTCHAIN_CALL:
		return "OPTCHAIN_CALL"
	case YIELD:
		return "YIELD"
	case AWAIT:
		return "AWAIT"
	case FOR_AWAIT_OF:
		return "FOR_AWAIT_OF"
	case CAST:
		return "CAST"
	case SUPER:
		return "SUPER"
	case NEW_TARGET:
		return "NEW_TARGET"
	case IMPORT_META:
		return "IMPORT_META"
	case COMMA:
		return "COMMA"
	case INC:
		return "INC"
	case DEC:
		return "DEC"
	case EXPONENT:
		return "EXPONENT"
	case ASSIGN_EXPONENT:
		return "ASSIGN_EXPONENT"
	case EMPTY:
		return "EMPTY"
	}
	return "ERROR"
}

// IsAssign reports whether t is one of the assignment family tokens
// (bare or compound). Mirrors NodeUtil.isAssignmentOp in Java.
func (t Token) IsAssign() bool {
	switch t {
	case ASSIGN, ASSIGN_ADD, ASSIGN_SUB, ASSIGN_MUL, ASSIGN_DIV, ASSIGN_MOD,
		ASSIGN_BITAND, ASSIGN_BITOR, ASSIGN_BITXOR, ASSIGN_LSH, ASSIGN_RSH, ASSIGN_URSH,
		ASSIGN_AND, ASSIGN_OR, ASSIGN_COALESCE, ASSIGN_EXPONENT:
		return true
	}
	return false
}

// IsOptChain reports whether t is one of the optional-chaining node
// kinds (OPTCHAIN_GETPROP / OPTCHAIN_GETELEM / OPTCHAIN_CALL). Mirrors
// NodeUtil.isOptChainNode in Java.
func (t Token) IsOptChain() bool {
	switch t {
	case OPTCHAIN_GETPROP, OPTCHAIN_GETELEM, OPTCHAIN_CALL:
		return true
	}
	return false
}

// IsPattern reports whether t is ARRAY_PATTERN or OBJECT_PATTERN —
// the two destructuring binding shapes. Matches NodeUtil.isPattern.
func (t Token) IsPattern() bool {
	return t == ARRAY_PATTERN || t == OBJECT_PATTERN
}
