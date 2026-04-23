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
	FOR_AWAIT_OF
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
	SUPER // `super` keyword (super(), super.prop)

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

	// Debugger statement
	DEBUGGER
	// Class / object members
	MEMBER_FUNCTION_DEF
	COMPUTED_PROP
	// Type-expression AST tokens — used in parsed JSDoc type trees, not in
	// the JS source AST.
	ITER_REST // `{...T}` rest parameter type
	QMARK     // nullable / unknown type marker `?`
	EQUALS    // optional parameter type marker `=`
	// Cast — /** @type {T} */ (expr). Wraps the inner expression as
	// its first child; the @type JSDocInfo is attached to this node.
	CAST
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
	case FOR_AWAIT_OF:
		return "FOR_AWAIT_OF"
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
	case SUPER:
		return "SUPER"
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
	case DEBUGGER:
		return "DEBUGGER"
	case MEMBER_FUNCTION_DEF:
		return "MEMBER_FUNCTION_DEF"
	case COMPUTED_PROP:
		return "COMPUTED_PROP"
	case ITER_REST:
		return "ITER_REST"
	case QMARK:
		return "QMARK"
	case EQUALS:
		return "EQUALS"
	case CAST:
		return "CAST"
	}
	return "ERROR"
}

// IsAssign reports whether t is one of the assignment family tokens
// (bare or compound). Mirrors NodeUtil.isAssignmentOp in Java.
func (t Token) IsAssign() bool {
	switch t {
	case ASSIGN, ASSIGN_ADD, ASSIGN_SUB, ASSIGN_MUL, ASSIGN_DIV, ASSIGN_MOD,
		ASSIGN_BITAND, ASSIGN_BITOR, ASSIGN_BITXOR, ASSIGN_LSH, ASSIGN_RSH, ASSIGN_URSH:
		return true
	}
	return false
}
