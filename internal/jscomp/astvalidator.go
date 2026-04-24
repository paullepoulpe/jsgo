// Port of com/google/javascript/jscomp/AstValidator.java
// @ closure-compiler v20260415 (subset; public API + structural shell).
//
// Ports the public interface: constructors, Process, ValidateRoot,
// ValidateCodeRoot, ValidateScript, ValidateStatement, ValidateExpression.
// Detailed per-token helpers are stubs pending transform-pass ports.

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

// Compiler is the minimal compiler context needed by AstValidator.
// TODO: replace with the concrete Compiler once that port lands.
type Compiler struct{}

func newAstValidatorFull(c *Compiler, h ViolationHandler, validateScripts, validateInlinings bool) *AstValidator {
	return &AstValidator{
		compiler:                         c,
		violationHandler:                 h,
		typeValidationMode:               TypeInfoValidationNone,
		isScriptFeatureValidationEnabled: validateScripts,
		shouldValidateRequiredInlinings:  validateInlinings,
	}
}

// NewAstValidator mirrors AstValidator(AbstractCompiler).
func NewAstValidator(compiler *Compiler) *AstValidator {
	return newAstValidatorFull(compiler, panicViolationHandler{}, false, false)
}

// NewAstValidatorValidateScripts mirrors AstValidator(AbstractCompiler, boolean).
func NewAstValidatorValidateScripts(compiler *Compiler, validateScriptFeatures bool) *AstValidator {
	return newAstValidatorFull(compiler, panicViolationHandler{}, validateScriptFeatures, false)
}

// NewAstValidatorFull mirrors the four-argument Java constructor.
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

// validateStatementHelper dispatches by token. Per-token bodies are
// stubs pending transform-pass ports.
func (av *AstValidator) validateStatementHelper(n *rhino.Node, _ bool) {
	if n != nil {
		_ = n.Token() // TODO: detailed per-token validation
	}
}

// ValidateExpression mirrors AstValidator.validateExpression.
func (av *AstValidator) ValidateExpression(n *rhino.Node) {
	if n != nil {
		_ = n.Token() // TODO: detailed per-token validation
	}
}

// violation dispatches to the registered ViolationHandler.
// Mirrors the private AstValidator.violation helper.
func (av *AstValidator) violation(message string, n *rhino.Node) {
	if av.violationHandler != nil {
		av.violationHandler.HandleViolation(message, n)
	}
}
