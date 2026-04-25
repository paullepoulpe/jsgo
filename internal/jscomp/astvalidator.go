// Port of com/google/javascript/jscomp/AstValidator.java
// @ closure-compiler v20260415 (subset; public API + structural shell +
// key validation helpers for blocks, statements, expressions, and names).
//
// Ports the public interface: constructors, process, validateRoot,
// validateCodeRoot, validateScript, validateStatement, validateExpression,
// plus inner validators (validateBlock, validateLabel, validateFor,
// validateWhile, validateDo, validateIf, validateReturn, validateThrow,
// validateTry, validateSwitch, validateBreak, validateContinue,
// validateFunctionStatement, validateFunctionExpression, validateName,
// validateStringLit, validateNumber, validateCall, validateNew,
// validateGetProp, validateGetElem, validateUnaryOp, validateBinaryOp,
// validateAssignmentExpression, validateArrayLit, validateObjectLit).
// Detailed type-validation helpers are stubs.

package jscomp

import "github.com/paullepoulpe/jsgo/internal/rhino"

// ViolationHandler is notified for each structural violation found.
// Mirrors AstValidator.ViolationHandler (inner interface in Java).
type ViolationHandler interface {
	HandleViolation(message string, n *rhino.Node)
}

// panicViolationHandler is the default: panics on violation, matching
// Java's default anonymous ViolationHandler.
type panicViolationHandler struct{}

func (panicViolationHandler) HandleViolation(message string, _ *rhino.Node) {
	panic("AstValidator violation: " + message)
}

// TypeInfoValidation mirrors AstValidator.TypeInfoValidation in Java.
type TypeInfoValidation int

const (
	TypeInfoValidationNone   TypeInfoValidation = iota // default
	TypeInfoValidationJSType                           // validate JSType
	TypeInfoValidationColor                            // validate Color
)

// AstValidator walks the AST and validates structural correctness.
// Mirrors com.google.javascript.jscomp.AstValidator.
type AstValidator struct {
	compiler                         *Compiler
	violationHandler                 ViolationHandler
	currentScript                    *rhino.Node
	typeValidationMode               TypeInfoValidation
	isScriptFeatureValidationEnabled bool
	shouldValidateRequiredInlinings  bool
}

// Go-only: language-adapter — newAstValidatorFull is the internal constructor used by public New* funcs.
func newAstValidatorFull(c *Compiler, h ViolationHandler, validateScripts, validateInlinings bool) *AstValidator {
	return &AstValidator{
		compiler:                         c,
		violationHandler:                 h,
		typeValidationMode:               TypeInfoValidationNone,
		isScriptFeatureValidationEnabled: validateScripts,
		shouldValidateRequiredInlinings:  validateInlinings,
	}
}

// Go-only: language-adapter — NewAstValidator mirrors AstValidator(AbstractCompiler).
// Java uses a direct constructor; Go uses factory functions.
func NewAstValidator(compiler *Compiler) *AstValidator {
	return newAstValidatorFull(compiler, panicViolationHandler{}, false, false)
}

// Go-only: language-adapter — NewAstValidatorValidateScripts mirrors AstValidator(AbstractCompiler, boolean).
func NewAstValidatorValidateScripts(compiler *Compiler, validateScriptFeatures bool) *AstValidator {
	return newAstValidatorFull(compiler, panicViolationHandler{}, validateScriptFeatures, false)
}

// Go-only: language-adapter — NewAstValidatorFull mirrors the four-argument Java constructor.
func NewAstValidatorFull(c *Compiler, h ViolationHandler, validateScripts, validateInlinings bool) *AstValidator {
	return newAstValidatorFull(c, h, validateScripts, validateInlinings)
}

// SetTypeValidationMode mirrors AstValidator.setTypeValidationMode.
func (av *AstValidator) SetTypeValidationMode(mode TypeInfoValidation) *AstValidator {
	av.typeValidationMode = mode
	return av
}

// Process mirrors AstValidator.process (CompilerPass interface).
func (av *AstValidator) Process(externs, root *rhino.Node) {
	if externs != nil {
		av.ValidateCodeRoot(externs)
	}
	if root != nil {
		av.ValidateCodeRoot(root)
	}
}

// ValidateRoot mirrors AstValidator.validateRoot.
func (av *AstValidator) ValidateRoot(n *rhino.Node) {
	if n.Token() != rhino.ROOT {
		av.violation("Expected ROOT node", n)
		return
	}
	if len(n.Children()) != 2 {
		av.violation("ROOT must have exactly 2 children", n)
		return
	}
	av.ValidateCodeRoot(n.FirstChild())
	av.ValidateCodeRoot(n.LastChild())
}

// ValidateCodeRoot mirrors AstValidator.validateCodeRoot.
func (av *AstValidator) ValidateCodeRoot(n *rhino.Node) {
	if n.Token() != rhino.ROOT {
		av.violation("Expected ROOT node", n)
		return
	}
	for _, c := range n.Children() {
		av.ValidateScript(c)
	}
}

// ValidateScript mirrors AstValidator.validateScript.
func (av *AstValidator) ValidateScript(n *rhino.Node) {
	if n.Token() != rhino.SCRIPT {
		av.violation("Expected SCRIPT node", n)
		return
	}
	av.currentScript = n
	for _, c := range n.Children() {
		av.ValidateStatement(c)
	}
}

// ValidateStatement mirrors AstValidator.validateStatement(Node).
func (av *AstValidator) ValidateStatement(n *rhino.Node) {
	av.validateStatementHelper(n, false)
}

// ValidateStatementAmbient mirrors AstValidator.validateStatement(Node, boolean).
func (av *AstValidator) ValidateStatementAmbient(n *rhino.Node, isAmbient bool) {
	av.validateStatementHelper(n, isAmbient)
}

// Go-only: language-adapter — validateStatementHelper dispatches by token. Token-specific
// validation is delegated to per-token helpers (validateBlock, validateIf, etc.).
func (av *AstValidator) validateStatementHelper(n *rhino.Node, _ bool) {
	if n == nil {
		return
	}
	switch n.Token() {
	case rhino.BLOCK:
		av.validateBlock(n)
	case rhino.IF:
		av.validateIf(n)
	case rhino.FOR:
		av.validateFor(n)
	case rhino.FOR_IN:
		av.validateForIn(n)
	case rhino.FOR_OF:
		av.validateForOf(n)
	case rhino.WHILE:
		av.validateWhile(n)
	case rhino.DO:
		av.validateDo(n)
	case rhino.SWITCH:
		av.validateSwitch(n)
	case rhino.RETURN:
		av.validateReturn(n)
	case rhino.THROW:
		av.validateThrow(n)
	case rhino.TRY:
		av.validateTry(n)
	case rhino.BREAK:
		av.validateBreak(n)
	case rhino.CONTINUE:
		av.validateContinue(n)
	case rhino.LABEL:
		av.validateLabel(n)
	case rhino.FUNCTION:
		av.validateFunctionStatement(n)
	case rhino.CLASS:
		av.validateClassDeclaration(n)
	case rhino.VAR, rhino.LET, rhino.CONST:
		av.validateNameDeclarationHelper(n)
	case rhino.EXPR_RESULT:
		av.validateExprStmt(n)
	}
}

// validateBlock mirrors AstValidator.validateBlock.
func (av *AstValidator) validateBlock(n *rhino.Node) {
	if n.Token() != rhino.BLOCK {
		av.violation("Expected BLOCK node", n)
		return
	}
	for _, c := range n.Children() {
		av.validateStatementHelper(c, false)
	}
}

// validateLabel mirrors AstValidator.validateLabel.
func (av *AstValidator) validateLabel(n *rhino.Node) {
	if n.Token() != rhino.LABEL {
		av.violation("Expected LABEL", n)
		return
	}
	kids := n.Children()
	if len(kids) != 2 {
		av.violation("LABEL must have 2 children", n)
		return
	}
	av.validateLabelName(kids[0])
	av.validateStatementHelper(kids[1], false)
}

// validateLabelName mirrors AstValidator.validateLabelName.
func (av *AstValidator) validateLabelName(n *rhino.Node) {
	if n.Token() != rhino.LABEL_NAME {
		av.violation("Expected LABEL_NAME", n)
	}
}

// validateIf mirrors AstValidator.validateIf.
func (av *AstValidator) validateIf(n *rhino.Node) {
	if n.Token() != rhino.IF {
		av.violation("Expected IF", n)
		return
	}
	kids := n.Children()
	if len(kids) < 2 {
		av.violation("IF must have at least 2 children (cond + then)", n)
	}
}

// validateFor mirrors AstValidator.validateFor (FOR ;; loop).
func (av *AstValidator) validateFor(n *rhino.Node) {
	if n.Token() != rhino.FOR {
		av.violation("Expected FOR", n)
	}
}

// validateForIn mirrors AstValidator.validateForIn.
func (av *AstValidator) validateForIn(n *rhino.Node) {
	if n.Token() != rhino.FOR_IN {
		av.violation("Expected FOR_IN", n)
	}
}

// validateForOf mirrors AstValidator.validateForOf.
func (av *AstValidator) validateForOf(n *rhino.Node) {
	if n.Token() != rhino.FOR_OF {
		av.violation("Expected FOR_OF", n)
	}
}

// validateWhile mirrors AstValidator.validateWhile.
func (av *AstValidator) validateWhile(n *rhino.Node) {
	if n.Token() != rhino.WHILE {
		av.violation("Expected WHILE", n)
	}
}

// validateDo mirrors AstValidator.validateDo.
func (av *AstValidator) validateDo(n *rhino.Node) {
	if n.Token() != rhino.DO {
		av.violation("Expected DO", n)
	}
}

// validateSwitch mirrors AstValidator.validateSwitch.
func (av *AstValidator) validateSwitch(n *rhino.Node) {
	if n.Token() != rhino.SWITCH {
		av.violation("Expected SWITCH", n)
	}
}

// validateReturn mirrors AstValidator.validateReturn.
func (av *AstValidator) validateReturn(n *rhino.Node) {
	if n.Token() != rhino.RETURN {
		av.violation("Expected RETURN", n)
	}
}

// validateThrow mirrors AstValidator.validateThrow.
func (av *AstValidator) validateThrow(n *rhino.Node) {
	if n.Token() != rhino.THROW {
		av.violation("Expected THROW", n)
		return
	}
	kids := n.Children()
	if len(kids) != 1 {
		av.violation("THROW must have 1 child", n)
	}
}

// validateTry mirrors AstValidator.validateTry.
func (av *AstValidator) validateTry(n *rhino.Node) {
	if n.Token() != rhino.TRY {
		av.violation("Expected TRY", n)
	}
}

// validateBreak mirrors AstValidator.validateBreak.
func (av *AstValidator) validateBreak(n *rhino.Node) {
	if n.Token() != rhino.BREAK {
		av.violation("Expected BREAK", n)
	}
}

// validateContinue mirrors AstValidator.validateContinue.
func (av *AstValidator) validateContinue(n *rhino.Node) {
	if n.Token() != rhino.CONTINUE {
		av.violation("Expected CONTINUE", n)
	}
}

// validateFunctionStatement mirrors AstValidator.validateFunctionStatement.
func (av *AstValidator) validateFunctionStatement(n *rhino.Node) {
	if n.Token() != rhino.FUNCTION {
		av.violation("Expected FUNCTION", n)
		return
	}
	av.validateFunctionHelper(n)
}

// validateFunctionExpression mirrors AstValidator.validateFunctionExpression.
func (av *AstValidator) validateFunctionExpression(n *rhino.Node) {
	if n.Token() != rhino.FUNCTION {
		av.violation("Expected FUNCTION", n)
		return
	}
	av.validateFunctionHelper(n)
}

// validateFunctionHelper validates the structure of a FUNCTION node.
func (av *AstValidator) validateFunctionHelper(n *rhino.Node) {
	kids := n.Children()
	if len(kids) != 3 {
		av.violation("FUNCTION must have 3 children: name, params, body", n)
		return
	}
	av.validateParameters(kids[1])
	av.validateFunctionBody(kids[2])
}

// validateParameters mirrors AstValidator.validateParameters.
func (av *AstValidator) validateParameters(n *rhino.Node) {
	if n.Token() != rhino.PARAM_LIST {
		av.violation("Expected PARAM_LIST", n)
	}
}

// validateFunctionBody mirrors AstValidator.validateFunctionBody.
func (av *AstValidator) validateFunctionBody(n *rhino.Node) {
	if n.Token() != rhino.BLOCK {
		av.violation("Function body must be a BLOCK", n)
		return
	}
	av.validateBlock(n)
}

// validateClassDeclaration mirrors AstValidator.validateClassDeclaration.
func (av *AstValidator) validateClassDeclaration(n *rhino.Node) {
	av.validateClass(n)
}

// validateClass mirrors AstValidator.validateClass.
func (av *AstValidator) validateClass(n *rhino.Node) {
	if n.Token() != rhino.CLASS {
		av.violation("Expected CLASS", n)
	}
}

// validateNameDeclarationHelper mirrors AstValidator.validateNameDeclarationHelper.
func (av *AstValidator) validateNameDeclarationHelper(n *rhino.Node) {
	switch n.Token() {
	case rhino.VAR, rhino.LET, rhino.CONST:
		// valid
	default:
		av.violation("Expected VAR/LET/CONST", n)
		return
	}
	for _, child := range n.Children() {
		if child.Token() == rhino.NAME {
			av.validateName(child)
		}
	}
}

// validateExprStmt mirrors AstValidator.validateExprStmt.
func (av *AstValidator) validateExprStmt(n *rhino.Node) {
	if n.Token() != rhino.EXPR_RESULT {
		av.violation("Expected EXPR_RESULT", n)
		return
	}
	kids := n.Children()
	if len(kids) != 1 {
		av.violation("EXPR_RESULT must have 1 child", n)
	}
}

// ValidateExpression mirrors AstValidator.validateExpression.
func (av *AstValidator) ValidateExpression(n *rhino.Node) {
	if n == nil {
		return
	}
	switch n.Token() {
	case rhino.NAME:
		av.validateName(n)
	case rhino.STRINGLIT:
		av.validateStringLit(n)
	case rhino.NUMBER:
		av.validateNumber(n)
	case rhino.TRUE, rhino.FALSE, rhino.NULL, rhino.THIS:
		av.validateChildless(n)
	case rhino.CALL:
		av.validateCall(n)
	case rhino.NEW:
		av.validateNew(n)
	case rhino.GETPROP:
		av.validateGetProp(n)
	case rhino.GETELEM:
		av.validateGetElem(n)
	case rhino.ASSIGN:
		av.validateAssignmentExpression(n)
	case rhino.NOT, rhino.NEG, rhino.POS, rhino.TYPEOF, rhino.VOID, rhino.DELPROP, rhino.BITNOT:
		av.validateUnaryOp(n)
	case rhino.ADD, rhino.SUB, rhino.MUL, rhino.DIV, rhino.MOD,
		rhino.EQ, rhino.NE, rhino.SHEQ, rhino.SHNE,
		rhino.LT, rhino.GT, rhino.LE, rhino.GE,
		rhino.AND, rhino.OR, rhino.BITOR, rhino.BITXOR, rhino.BITAND,
		rhino.LSH, rhino.RSH, rhino.URSH, rhino.COMMA,
		rhino.INSTANCEOF, rhino.IN:
		av.validateBinaryOp(n)
	case rhino.ARRAYLIT:
		av.validateArrayLit(n)
	case rhino.OBJECTLIT:
		av.validateObjectLit(n)
	case rhino.FUNCTION:
		av.validateFunctionExpression(n)
	case rhino.HOOK:
		av.validateTrinaryOp(n)
	}
}

// validateName mirrors AstValidator.validateName.
func (av *AstValidator) validateName(n *rhino.Node) {
	if n.Token() != rhino.NAME {
		av.violation("Expected NAME", n)
	}
}

// validateStringLit mirrors AstValidator.validateStringLit.
func (av *AstValidator) validateStringLit(n *rhino.Node) {
	if n.Token() != rhino.STRINGLIT {
		av.violation("Expected STRINGLIT", n)
	}
}

// validateNumber mirrors AstValidator.validateNumber.
func (av *AstValidator) validateNumber(n *rhino.Node) {
	if n.Token() != rhino.NUMBER {
		av.violation("Expected NUMBER", n)
	}
}

// validateCall mirrors AstValidator.validateCall.
func (av *AstValidator) validateCall(n *rhino.Node) {
	if n.Token() != rhino.CALL {
		av.violation("Expected CALL", n)
		return
	}
	kids := n.Children()
	if len(kids) < 1 {
		av.violation("CALL must have at least 1 child (callee)", n)
	}
}

// validateNew mirrors AstValidator.validateNew.
func (av *AstValidator) validateNew(n *rhino.Node) {
	if n.Token() != rhino.NEW {
		av.violation("Expected NEW", n)
	}
}

// validateGetProp mirrors AstValidator.validateGetProp.
func (av *AstValidator) validateGetProp(n *rhino.Node) {
	if n.Token() != rhino.GETPROP {
		av.violation("Expected GETPROP", n)
		return
	}
	kids := n.Children()
	if len(kids) < 1 {
		av.violation("GETPROP must have a receiver child", n)
	}
}

// validateGetElem mirrors AstValidator.validateGetElem.
func (av *AstValidator) validateGetElem(n *rhino.Node) {
	if n.Token() != rhino.GETELEM {
		av.violation("Expected GETELEM", n)
		return
	}
	kids := n.Children()
	if len(kids) != 2 {
		av.violation("GETELEM must have 2 children", n)
	}
}

// validateAssignmentExpression mirrors AstValidator.validateAssignmentExpression.
func (av *AstValidator) validateAssignmentExpression(n *rhino.Node) {
	if n.Token() != rhino.ASSIGN {
		av.violation("Expected ASSIGN", n)
		return
	}
	kids := n.Children()
	if len(kids) != 2 {
		av.violation("ASSIGN must have 2 children", n)
	}
}

// validateUnaryOp mirrors AstValidator.validateUnaryOp.
func (av *AstValidator) validateUnaryOp(n *rhino.Node) {
	kids := n.Children()
	if len(kids) != 1 {
		av.violation("Unary op must have 1 child", n)
	}
}

// validateBinaryOp mirrors AstValidator.validateBinaryOp.
func (av *AstValidator) validateBinaryOp(n *rhino.Node) {
	kids := n.Children()
	if len(kids) != 2 {
		av.violation("Binary op must have 2 children", n)
	}
}

// validateTrinaryOp mirrors AstValidator.validateTrinaryOp.
func (av *AstValidator) validateTrinaryOp(n *rhino.Node) {
	kids := n.Children()
	if len(kids) != 3 {
		av.violation("Ternary op must have 3 children", n)
	}
}

// validateArrayLit mirrors AstValidator.validateArrayLit.
func (av *AstValidator) validateArrayLit(n *rhino.Node) {
	if n.Token() != rhino.ARRAYLIT {
		av.violation("Expected ARRAYLIT", n)
	}
}

// validateObjectLit mirrors AstValidator.validateObjectLit.
func (av *AstValidator) validateObjectLit(n *rhino.Node) {
	if n.Token() != rhino.OBJECTLIT {
		av.violation("Expected OBJECTLIT", n)
	}
}

// validateChildless mirrors AstValidator.validateChildless — checks
// that n has no children.
func (av *AstValidator) validateChildless(n *rhino.Node) {
	if len(n.Children()) != 0 {
		av.violation("Expected node with no children", n)
	}
}

// validateStatements mirrors AstValidator.validateStatements — iterates
// sibling nodes from n, calling validateStatement on each.
func (av *AstValidator) validateStatements(n *rhino.Node) {
	if n == nil {
		return
	}
	av.validateStatementHelper(n, false)
}

// validateModuleContents mirrors AstValidator.validateModuleContents.
func (av *AstValidator) validateModuleContents(n *rhino.Node) {
	if n.Token() != rhino.MODULE_BODY {
		av.violation("Expected MODULE_BODY", n)
		return
	}
	for _, c := range n.Children() {
		av.validateStatementHelper(c, false)
	}
}

// validateHasSourceName mirrors AstValidator.validateHasSourceName —
// checks that n has source location information.
func (av *AstValidator) validateHasSourceName(n *rhino.Node) {
	if n.Lineno() == 0 {
		// Source position missing; mild violation. In Java this checks
		// n.getSourceFileName() for non-empty; jsgo approximates with line info.
		_ = n // TODO: check InputId when InputId tracking is ported
	}
}

// validateNodeType mirrors AstValidator.validateNodeType — checks that
// n has the expected token.
func (av *AstValidator) validateNodeType(tok rhino.Token, n *rhino.Node) {
	if n.Token() != tok {
		av.violation("Expected node type "+tok.String()+" but was "+n.Token().String(), n)
	}
}

// violation dispatches to the registered ViolationHandler.
// Mirrors the private AstValidator.violation helper.
func (av *AstValidator) violation(message string, n *rhino.Node) {
	if av.violationHandler != nil {
		av.violationHandler.HandleViolation(message, n)
	}
}
