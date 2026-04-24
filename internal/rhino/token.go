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
	WITH // `with (obj) body`

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

	// Destructuring pattern nodes — mirror Token.java ARRAY_PATTERN / OBJECT_PATTERN.
	ARRAY_PATTERN
	OBJECT_PATTERN
	DEFAULT_VALUE // `x = dflt` inside param/pattern
	ITER_SPREAD   // `...expr` spread in call args / array/obj literals
	OBJECT_REST   // `...target` at tail of object pattern

	// Optional-chaining nodes — mirror Token.java OPTCHAIN_* entries.
	OPTCHAIN_GETPROP // `a?.b`
	OPTCHAIN_GETELEM // `a?.[b]`
	OPTCHAIN_CALL    // `a?.(args)`

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
	// Regular-expression literal  `/pattern/flags`
	REGEXP
	// Template literals.
	TEMPLATELIT
	TEMPLATELIT_STRING
	TEMPLATELIT_SUB
	TAGGED_TEMPLATELIT
	// Prefix/postfix increment and decrement. Not yet emitted by the
	// scanner (TODO: port Scanner.readOperator's ++/-- path); listed
	// here so CFG + AST helpers can reference the token by name.
	INC
	DEC
	// Generator / async expressions. Not yet emitted by the scanner;
	// listed for CFG reference.
	YIELD
	AWAIT
	// Dynamic import expression `import(specifier)`.
	// Mirrors Token.DYNAMIC_IMPORT in Java. The single child is the
	// specifier expression (usually a STRINGLIT).
	DYNAMIC_IMPORT
	// CAST represents `/** @type {T} */(expr)`. The parser synthesises
	// this node when a JSDoc block comment immediately precedes a
	// parenthesised expression; it carries the @type annotation on its
	// JSDocInfo and the inner expression as its sole child.
	CAST

	// AST-only tokens — not emitted by the scanner; used in IR/CFG passes
	// and module-graph plumbing. Mirror Token.java entries of the same name.
	PLACEHOLDER1 // reserved IR placeholder (slot 1)
	PLACEHOLDER2 // reserved IR placeholder (slot 2)
	ROOT         // synthetic root wrapping externsRoot + jsRoot
	INPUT_ID     // synthetic node carrying a CompilerInput id
	IMPORT_STAR  // `* as name` namespace import
	MODULE_BODY  // body of an ES module
	IMPORT_META  // `import.meta` meta-property
	ES6_SUPER    // `super` in ES6 class context (distinct from SUPER in some passes)

	// TS-only tokens — TypeScript constructs that have no JS equivalent.
	// Declared here so passes that inspect TS ASTs can reference them by
	// name; they are never emitted by the JS scanner.
	INTERFACE          // `interface Foo { … }`
	IMPLEMENTS_MEMBERS // implements-clause member list on a class
	DECLARE            // `declare` modifier (ambient declaration)
	TYPE_ALIAS         // `type Foo = Bar`
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
	case WITH:
		return "WITH"
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
	case ARRAY_PATTERN:
		return "ARRAY_PATTERN"
	case OBJECT_PATTERN:
		return "OBJECT_PATTERN"
	case DEFAULT_VALUE:
		return "DEFAULT_VALUE"
	case ITER_SPREAD:
		return "ITER_SPREAD"
	case OBJECT_REST:
		return "OBJECT_REST"
	case OPTCHAIN_GETPROP:
		return "OPTCHAIN_GETPROP"
	case OPTCHAIN_GETELEM:
		return "OPTCHAIN_GETELEM"
	case OPTCHAIN_CALL:
		return "OPTCHAIN_CALL"
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
	case INC:
		return "INC"
	case DEC:
		return "DEC"
	case YIELD:
		return "YIELD"
	case AWAIT:
		return "AWAIT"
	case DYNAMIC_IMPORT:
		return "DYNAMIC_IMPORT"
	case CAST:
		return "CAST"
	case PLACEHOLDER1:
		return "PLACEHOLDER1"
	case PLACEHOLDER2:
		return "PLACEHOLDER2"
	case ROOT:
		return "ROOT"
	case INPUT_ID:
		return "INPUT_ID"
	case IMPORT_STAR:
		return "IMPORT_STAR"
	case MODULE_BODY:
		return "MODULE_BODY"
	case IMPORT_META:
		return "IMPORT_META"
	case ES6_SUPER:
		return "ES6_SUPER"
	case INTERFACE:
		return "INTERFACE"
	case IMPLEMENTS_MEMBERS:
		return "IMPLEMENTS_MEMBERS"
	case DECLARE:
		return "DECLARE"
	case TYPE_ALIAS:
		return "TYPE_ALIAS"
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
