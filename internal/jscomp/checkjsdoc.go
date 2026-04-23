// Port of com/google/javascript/jscomp/CheckJSDoc.java @ closure-compiler v20260415
// (subset: checks implementable against the Go AST and JSDocInfo surface;
// Java-only tokens MEMBER_FUNCTION_DEF / GETTER_DEF / SETTER_DEF /
// COMPUTED_PROP / CAST / EXPORT / WITH / REST / DEFAULT_VALUE are absent from
// the Go parser, so those branches are omitted and noted with TODO comments).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CheckJSDoc mirrors CheckJSDoc.java's CompilerPass.
// Java walks the AST using NodeTraversal (post-order); here we do a
// plain recursive walk so we avoid coupling to the NodeTraversal
// interface for a pass that doesn't need scope.
type CheckJSDoc struct {
	path      string
	inExterns bool // TODO: wire real externs tracking when externs are loaded
	diags     []diag.Diagnostic
}

// NewCheckJSDoc mirrors the Java constructor CheckJSDoc(AbstractCompiler).
func NewCheckJSDoc(path string) *CheckJSDoc {
	return &CheckJSDoc{path: path}
}

// Check is the entry point; mirrors CompilerPass.process(externs, root).
// In the Java pass, process traverses the externs tree with inExterns=true
// and the source tree with inExterns=false. The Go pass receives a single
// SCRIPT node (the source); a future extension can pass an externs node.
func (c *CheckJSDoc) Check(script *rhino.Node) []diag.Diagnostic {
	c.diags = nil
	c.inExterns = false
	c.cjdWalk(script)
	return c.diags
}

// cjdWalk is the post-order visit loop, mirroring
// AbstractPostOrderCallback.visit called bottom-up by NodeTraversal.
func (c *CheckJSDoc) cjdWalk(n *rhino.Node) {
	if n == nil {
		return
	}
	// Recurse into children first (post-order).
	for _, ch := range n.Children() {
		c.cjdWalk(ch)
	}
	c.visit(n)
}

// visit is the main dispatch, mirroring CheckJSDoc.visit(NodeTraversal, Node, Node).
func (c *CheckJSDoc) visit(n *rhino.Node) {
	info := n.JSDocInfo()
	c.validateTypeAnnotations(n, info)
	c.validateFunctionJsDoc(n, info)
	c.validateMsgJsDoc(n, info)
	c.validateClassLevelJsDoc(n, info)
	c.validateArrowFunction(n)
	// validateRestParameter: TODO — REST token absent from Go AST.
	// validateDefaultValue: TODO — DEFAULT_VALUE token absent from Go AST.
	c.validateTemplates(n, info)
	c.validateTypedefs(n, info)
	c.validateNoSideEffects(n, info)
	c.validateAbstractJsDoc(n, info)
	c.validateDefinesDeclaration(n, info)
	c.validateSuppress(n, info)
	// validateImplicitCast: TODO — @implicitCast not in Go JSDocInfo.
	// validateClosurePrimitive: TODO — @closurePrimitive not in Go JSDocInfo.
	c.validateReturnJsDoc(n, info)
	// validateTsType: TODO — @tsType not in Go JSDocInfo.
	// validateJsDocTypeNames: handled by CheckJsDocTypes.
	// validateIsUsedViaDotConstructor: TODO — @usedViaDotConstructor not in Go JSDocInfo.
}

// ---------------------------------------------------------------------------
// validateSuppress
// Mirrors CheckJSDoc.validateSuppress. In Java, suppressions are valid on
// FUNCTION, CLASS, VAR, LET, CONST, SCRIPT, MEMBER_FUNCTION_DEF,
// GETTER_DEF, SETTER_DEF, MEMBER_FIELD_DEF, COMPUTED_FIELD_DEF,
// COMPUTED_PROP (only when last child is a function), STRING_KEY (in
// object literals), and WITH (only for "with"). Anywhere else is
// MISPLACED_SUPPRESS.
// Go AST subset: FUNCTION, CLASS, VAR, LET, CONST, SCRIPT, CLASS_FIELD_DEF,
// STRING_KEY (in OBJECTLIT), ASSIGN (expr-result), all qualify; everything
// else fires MISPLACED_SUPPRESS.
// ---------------------------------------------------------------------------
func (c *CheckJSDoc) validateSuppress(n *rhino.Node, info *rhino.JSDocInfo) {
	if info == nil || len(info.Suppressions()) == 0 {
		return
	}
	tok := n.Token()
	switch tok {
	case rhino.FUNCTION, rhino.CLASS, rhino.VAR, rhino.LET, rhino.CONST,
		rhino.SCRIPT, rhino.CLASS_FIELD_DEF:
		// Always valid.
		return
	case rhino.STRING_KEY:
		if p := n.Parent(); p != nil && p.Token() == rhino.OBJECTLIT {
			return
		}
	default:
	}
	// A suppression placed on an expression-result is valid (Java allows
	// @suppress on EXPR_RESULT via its parent check).
	if p := n.Parent(); p != nil && p.Token() == rhino.EXPR_RESULT {
		return
	}
	c.report(n, diag.JSC_MISPLACED_SUPPRESS, diag.Warning,
		"@suppress annotation not allowed here. See"+
			" https://github.com/google/closure-compiler/wiki/@suppress-annotations")
}

// cjdContainsOnlySuppressionFor mirrors Java's containsOnlySuppressionFor.
func cjdContainsOnlySuppressionFor(info *rhino.JSDocInfo, allowed string) bool {
	s := info.Suppressions()
	if len(s) != 1 {
		return false
	}
	for k := range s {
		return k == allowed
	}
	return false
}

// ---------------------------------------------------------------------------
// validateTypedefs
// Mirrors CheckJSDoc.validateTypedefs.
// ---------------------------------------------------------------------------
func (c *CheckJSDoc) validateTypedefs(n *rhino.Node, info *rhino.JSDocInfo) {
	if info == nil || info.TypedefType() == nil {
		return
	}
	if cjdIsClassDecl(n) {
		c.reportMisplaced(n, "typedef", "@typedef is not allowed on a class declaration.")
		return
	}
	// lvalue: for a name-declaration or assign, descend to the LHS.
	lvalue := n
	if cjdIsNameDeclaration(n) || n.Token() == rhino.ASSIGN {
		lvalue = n.FirstChild()
	}
	if lvalue == nil {
		return
	}
	if !cjdIsQualifiedName(lvalue) {
		c.reportMisplaced(n, "typedef",
			"@typedef is only allowed on qualified name declarations. Did you mean @type?")
	} else if cjdIsPrototypeOrInstanceDecl(lvalue) {
		c.reportMisplaced(n, "typedef",
			"@typedef is not allowed on instance or prototype properties. Did you mean @type?")
	}
}

// ---------------------------------------------------------------------------
// validateTemplates
// Mirrors CheckJSDoc.validateTemplates.
// ---------------------------------------------------------------------------
func (c *CheckJSDoc) validateTemplates(n *rhino.Node, info *rhino.JSDocInfo) {
	if info == nil || len(info.TemplateNames()) == 0 {
		return
	}
	if info.IsConstructor() || info.IsInterface() {
		return
	}
	if cjdIsClassDecl(n) {
		return
	}
	// containsFunctionDeclaration / getFunctionDecl approximation:
	if n.Token() == rhino.FUNCTION {
		return
	}
	if cjdGetFunctionDecl(n) != nil {
		return
	}
	c.reportMisplaced(n, "template",
		"@template is only allowed in class, constructor, interface, function "+
			"or method declarations")
}

// ---------------------------------------------------------------------------
// validateClassLevelJsDoc
// Mirrors CheckJSDoc.validateClassLevelJsDoc.
// Java checks for MEMBER_FUNCTION_DEF nodes; the Go AST represents class
// methods as plain FUNCTION children inside a CLASS BLOCK. We fire on a
// FUNCTION whose parent is a CLASS body BLOCK.
// ---------------------------------------------------------------------------
func (c *CheckJSDoc) validateClassLevelJsDoc(n *rhino.Node, info *rhino.JSDocInfo) {
	if info == nil || !cjdHasClassLevelJsDoc(info) {
		return
	}
	// In the Go AST a class member function sits inside the CLASS body block.
	// Pattern: n.Token()==FUNCTION && n.Parent().Token()==BLOCK &&
	//          n.Parent().Parent().Token()==CLASS.
	if n.Token() != rhino.FUNCTION {
		return
	}
	p := n.Parent()
	if p == nil {
		return
	}
	gp := p.Parent()
	if gp == nil {
		return
	}
	if p.Token() == rhino.BLOCK && gp.Token() == rhino.CLASS {
		c.report(n, diag.JSC_DISALLOWED_MEMBER_JSDOC, diag.Warning,
			"Class level JSDocs (@interface, @extends, etc.) are not allowed on class members")
	}
}

// cjdHasClassLevelJsDoc mirrors the private Java helper hasClassLevelJsDoc.
func cjdHasClassLevelJsDoc(info *rhino.JSDocInfo) bool {
	return info.IsConstructor() || info.IsInterface() ||
		info.BaseType() != nil ||
		len(info.Implementations()) != 0 ||
		len(info.BaseTypes()) > 1 // multiple @extends = interface extends
}

// ---------------------------------------------------------------------------
// validateAbstractJsDoc
// Mirrors CheckJSDoc.validateAbstractJsDoc.
// @abstract is not yet in the Go JSDocInfo, so this check is a TODO stub.
// ---------------------------------------------------------------------------
func (c *CheckJSDoc) validateAbstractJsDoc(_ *rhino.Node, _ *rhino.JSDocInfo) {
	// TODO: port when @abstract is added to rhino.JSDocInfo.
}

// ---------------------------------------------------------------------------
// validateFunctionJsDoc
// Mirrors CheckJSDoc.validateFunctionJsDoc.
// Fires MISPLACED_ANNOTATION when a function-bearing JSDoc
// (containsFunctionDeclaration) is not on a node that can host a function.
// In Go JSDocInfo, "containsFunctionDeclaration" means the comment has
// @constructor, @interface, @record, @param, @return, or @extends (but NOT
// @type — the Java check requires !info.hasType()).
// ---------------------------------------------------------------------------
func (c *CheckJSDoc) validateFunctionJsDoc(n *rhino.Node, info *rhino.JSDocInfo) {
	if info == nil {
		return
	}
	if !cjdContainsFunctionDeclaration(info) {
		return
	}
	if info.HasType() {
		return
	}
	if cjdIsJSDocOnFunctionNode(n, info) {
		return
	}
	c.reportMisplaced(n, "function",
		"This JSDoc is not attached to a function node. Are you missing parentheses?")
}

// cjdContainsFunctionDeclaration approximates Java's
// JSDocInfo.containsFunctionDeclaration(). Java returns true when any of:
// @constructor, @interface, @record, @param, @return, @template (on a
// function), or @abstract is set.
// Go approximation: true when @constructor/@interface/@record is set, or
// when @param or @return annotations are present.
func cjdContainsFunctionDeclaration(info *rhino.JSDocInfo) bool {
	if info.IsConstructor() || info.IsInterface() || info.IsRecord() {
		return true
	}
	if len(info.ParameterNames()) > 0 || info.ReturnType() != nil {
		return true
	}
	return false
}

// cjdIsJSDocOnFunctionNode mirrors Java's isJSDocOnFunctionNode.
// Returns true when the node can be the host for function-bearing JSDoc.
// Java's switch includes FUNCTION, GETTER_DEF, SETTER_DEF, MEMBER_FUNCTION_DEF,
// STRING_KEY, COMPUTED_PROP, EXPORT, MEMBER_FIELD_DEF, COMPUTED_FIELD_DEF,
// GETELEM/GETPROP (when first child is qualified), VAR/LET/CONST/ASSIGN.
// Go subset: FUNCTION, CLASS (for @constructor/@interface on class decls),
// STRING_KEY, CLASS_FIELD_DEF, VAR/LET/CONST, ASSIGN,
// GETPROP (when qualified and parent is EXPR_RESULT).
func cjdIsJSDocOnFunctionNode(n *rhino.Node, info *rhino.JSDocInfo) bool {
	switch n.Token() {
	case rhino.FUNCTION, rhino.CLASS, rhino.STRING_KEY, rhino.CLASS_FIELD_DEF:
		return true
	case rhino.GETPROP, rhino.GETELEM:
		if fc := n.FirstChild(); fc != nil && cjdIsQualifiedName(fc) {
			return true
		}
		return false
	case rhino.VAR, rhino.LET, rhino.CONST, rhino.ASSIGN:
		lhs := n.FirstChild()
		rhs := cjdGetRValueOfLValue(lhs)
		if rhs != nil && cjdIsClassNode(rhs) && !info.IsConstructor() {
			return false
		}
		return true
	default:
		return false
	}
}

// ---------------------------------------------------------------------------
// validateMsgJsDoc
// Mirrors CheckJSDoc.validateMsgJsDoc.
// @desc/@meaning/@alternateMessageId are not in Go JSDocInfo yet; this is
// a TODO stub.
// ---------------------------------------------------------------------------
func (c *CheckJSDoc) validateMsgJsDoc(_ *rhino.Node, _ *rhino.JSDocInfo) {
	// TODO: port when @desc / @meaning / @alternateMessageId are added to
	// rhino.JSDocInfo. Fire JSC_MISPLACED_MSG_ANNOTATION on name nodes
	// that don't start with MSG_.
}

// ---------------------------------------------------------------------------
// validateTypeAnnotations
// Mirrors CheckJSDoc.validateTypeAnnotations.
// Fires MISPLACED_ANNOTATION "type" when @type is in an invalid position.
// ---------------------------------------------------------------------------
func (c *CheckJSDoc) validateTypeAnnotations(n *rhino.Node, info *rhino.JSDocInfo) {
	if info == nil || !info.HasType() {
		return
	}
	valid := false
	switch n.Token() {
	case rhino.FUNCTION:
		// Function declarations are valid.
		valid = cjdIsFunctionDeclaration(n)
	case rhino.NAME:
		valid = cjdIsTypeAnnotationAllowedForName(n)
	case rhino.STRING_KEY, rhino.CLASS_FIELD_DEF:
		// Object literal properties and class fields are valid.
		valid = true
	case rhino.VAR, rhino.LET, rhino.CONST:
		// Declarations valid iff non-destructuring (Go parser doesn't
		// produce destructuring declarations, so always true here).
		valid = true
	case rhino.ASSIGN:
		lvalue := n.FirstChild()
		if lvalue != nil {
			valid = n.Parent() != nil &&
				n.Parent().Token() == rhino.EXPR_RESULT &&
				(lvalue.Token() == rhino.GETPROP ||
					lvalue.Token() == rhino.GETELEM ||
					(lvalue.Token() == rhino.NAME && lvalue.String_() == "exports"))
		}
	case rhino.GETPROP:
		valid = n.Parent() != nil &&
			n.Parent().Token() == rhino.EXPR_RESULT &&
			cjdIsQualifiedName(n)
	case rhino.CALL:
		// valid when @define (not in Go JSDocInfo yet); keep false.
		valid = false
	case rhino.CAST:
		// Cast expressions `/** @type {T} */(expr)` are always valid.
		valid = true
	// EXPORT, GETTER_DEF, SETTER_DEF, MEMBER_FIELD_DEF,
	// COMPUTED_FIELD_DEF, ARRAY_PATTERN, OBJECT_PATTERN: TODO — absent.
	default:
		valid = false
	}
	if !valid {
		c.reportMisplaced(n, "type",
			"Type annotations are not allowed here. Are you missing parentheses?")
	}
}

// cjdIsTypeAnnotationAllowedForName mirrors Java's isTypeAnnotationAllowedForName.
func cjdIsTypeAnnotationAllowedForName(n *rhino.Node) bool {
	if n.Token() != rhino.NAME {
		return false
	}
	if !cjdIsLValue(n) {
		return false
	}
	// Don't allow JSDoc on a name in an assignment; simple names should
	// only have JSDoc on them when originally declared.
	root := cjdGetRootTarget(n)
	return !cjdIsLhsOfAssign(root)
}

// ---------------------------------------------------------------------------
// validateArrowFunction
// Mirrors CheckJSDoc.validateArrowFunction.
// Fires JSC_ARROW_FUNCTION_AS_CONSTRUCTOR when an arrow function carries
// @constructor or @interface.
// ---------------------------------------------------------------------------
func (c *CheckJSDoc) validateArrowFunction(n *rhino.Node) {
	if n.Token() != rhino.FUNCTION {
		return
	}
	if !n.GetBooleanProp(rhino.PropArrowFn) {
		return
	}
	// Use node's own JSDoc or the best available from parent.
	info := cjdGetBestJSDocInfo(n)
	if info != nil && (info.IsConstructor() || info.IsInterface()) {
		c.report(n, diag.JSC_ARROW_FUNCTION_AS_CONSTRUCTOR, diag.Error,
			"Arrow functions cannot be used as constructors")
	}
}

// ---------------------------------------------------------------------------
// validateNoSideEffects
// Mirrors CheckJSDoc.validateNoSideEffects.
// Fires JSC_INVALID_MODIFIES_ANNOTATION when @modifies appears outside externs.
// Go JSDocInfo doesn't yet have hasSideEffectsArgumentsAnnotation /
// modifiesThis — this is a TODO stub.
// ---------------------------------------------------------------------------
func (c *CheckJSDoc) validateNoSideEffects(_ *rhino.Node, _ *rhino.JSDocInfo) {
	// TODO: port when @modifies / @nosideeffects are added to rhino.JSDocInfo.
	// Fire JSC_INVALID_MODIFIES_ANNOTATION when not inExterns.
}

// ---------------------------------------------------------------------------
// validateDefinesDeclaration
// Mirrors CheckJSDoc.validateDefinesDeclaration.
// Fires JSC_INVALID_DEFINE_ON_LET when @define appears on a let declaration.
// Go JSDocInfo doesn't yet have isDefine() — this is a TODO stub.
// ---------------------------------------------------------------------------
func (c *CheckJSDoc) validateDefinesDeclaration(_ *rhino.Node, _ *rhino.JSDocInfo) {
	// TODO: port when @define is added to rhino.JSDocInfo.
	// if info != nil && info.IsDefine() && n.Token() == rhino.LET {
	//     c.report(n, diag.JSC_INVALID_DEFINE_ON_LET, diag.Error, ...)
	// }
}

// ---------------------------------------------------------------------------
// validateReturnJsDoc
// Mirrors CheckJSDoc.validateReturnJsDoc.
// Fires JSC_JSDOC_ON_RETURN when annotations other than @type / @typedef
// appear directly on a return statement.
// Java checks info.containsDeclaration() && !info.hasType() && !info.hasTypedefType().
// Go approximation: containsDeclaration ≈ @constructor/@interface/@record/@param/@return.
// ---------------------------------------------------------------------------
func (c *CheckJSDoc) validateReturnJsDoc(n *rhino.Node, info *rhino.JSDocInfo) {
	if n.Token() != rhino.RETURN || info == nil {
		return
	}
	if !cjdContainsDeclaration(info) {
		return
	}
	if info.HasType() || info.TypedefType() != nil {
		return
	}
	c.report(n, diag.JSC_JSDOC_ON_RETURN, diag.Warning,
		"JSDoc annotations are not supported on return.")
}

// cjdContainsDeclaration approximates Java's JSDocInfo.containsDeclaration().
// Java returns true for many annotations; Go subset: @constructor, @interface,
// @record, @param, @return, @extends, @implements, @typedef, @enum.
func cjdContainsDeclaration(info *rhino.JSDocInfo) bool {
	if info.IsConstructor() || info.IsInterface() || info.IsRecord() || info.IsEnum() {
		return true
	}
	if len(info.ParameterNames()) > 0 || info.ReturnType() != nil {
		return true
	}
	if info.BaseType() != nil || len(info.Implementations()) > 0 {
		return true
	}
	if info.TypedefType() != nil {
		return true
	}
	return false
}

// ---------------------------------------------------------------------------
// report helpers (mirror Java's private report / reportMisplaced)
// ---------------------------------------------------------------------------

func (c *CheckJSDoc) reportMisplaced(n *rhino.Node, annotationName, note string) {
	c.report(n, diag.JSC_MISPLACED_ANNOTATION, diag.Warning,
		"Misplaced "+annotationName+" annotation. "+note)
}

func (c *CheckJSDoc) report(n *rhino.Node, code string, level diag.Level, msg string) {
	c.diags = append(c.diags, diag.Diagnostic{
		Path:    c.path,
		Line:    n.Lineno(),
		Col:     n.Charno(),
		Code:    code,
		Level:   level,
		Message: msg,
	})
}

// ---------------------------------------------------------------------------
// cjd* local AST helpers (mirror Java's NodeUtil / Node predicate methods)
// Prefixed cjd to avoid colliding with any package-level helpers.
// ---------------------------------------------------------------------------

// cjdIsNameDeclaration mirrors NodeUtil.isNameDeclaration.
func cjdIsNameDeclaration(n *rhino.Node) bool {
	tok := n.Token()
	return tok == rhino.VAR || tok == rhino.LET || tok == rhino.CONST
}

// cjdIsQualifiedName mirrors Node.isQualifiedName.
// A qualified name is a NAME, THIS, or a GETPROP chain whose leaves are
// all NAMEs or THIS. Java's Node.isQualifiedName includes THIS as a root.
func cjdIsQualifiedName(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	switch n.Token() {
	case rhino.NAME, rhino.THIS:
		return true
	case rhino.GETPROP:
		kids := n.Children()
		if len(kids) < 2 {
			return false
		}
		return cjdIsQualifiedName(kids[0])
	default:
		return false
	}
}

// cjdIsPrototypeOrInstanceDecl mirrors NodeUtil.isPrototypeOrInstanceDecl.
func cjdIsPrototypeOrInstanceDecl(n *rhino.Node) bool {
	if n.Token() == rhino.STRING_KEY {
		return false
	}
	if cjdIsPrototypeProperty(n) {
		return true
	}
	root := cjdGetRootOfQualifiedName(n)
	if root == nil {
		return false
	}
	return root.Token() == rhino.THIS
}

// cjdIsPrototypeProperty detects `X.prototype.y` chains.
func cjdIsPrototypeProperty(n *rhino.Node) bool {
	// We're looking for GETPROP nodes of the form (GETPROP (GETPROP ? "prototype") name)
	if n.Token() != rhino.GETPROP {
		return false
	}
	kids := n.Children()
	if len(kids) < 2 {
		return false
	}
	inner := kids[0]
	if inner.Token() != rhino.GETPROP {
		return false
	}
	innerKids := inner.Children()
	if len(innerKids) < 2 {
		return false
	}
	return innerKids[1].String_() == "prototype"
}

// cjdGetRootOfQualifiedName mirrors NodeUtil.getRootOfQualifiedName.
func cjdGetRootOfQualifiedName(n *rhino.Node) *rhino.Node {
	for n != nil && n.Token() == rhino.GETPROP {
		kids := n.Children()
		if len(kids) == 0 {
			return nil
		}
		n = kids[0]
	}
	return n
}

// cjdIsClassDecl mirrors CheckJSDoc.isClassDecl.
func cjdIsClassDecl(n *rhino.Node) bool {
	return cjdIsClassNode(n) ||
		(n.Token() == rhino.ASSIGN && n.LastChild() != nil && cjdIsClassNode(n.LastChild())) ||
		(cjdIsNameDeclaration(n) && n.FirstChild() != nil &&
			cjdIsNameInitializeWithClass(n.FirstChild())) ||
		cjdIsNameInitializeWithClass(n)
}

// cjdIsNameInitializeWithClass mirrors CheckJSDoc.isNameInitializeWithClass.
func cjdIsNameInitializeWithClass(n *rhino.Node) bool {
	return n != nil && n.Token() == rhino.NAME &&
		n.FirstChild() != nil && cjdIsClassDecl(n.FirstChild())
}

// cjdIsClassNode mirrors CheckJSDoc.isClass.
// Java also checks compiler.getCodingConvention().isClassFactoryCall(n);
// we omit that since we don't have a coding convention.
func cjdIsClassNode(n *rhino.Node) bool {
	return n != nil && n.Token() == rhino.CLASS
}

// cjdGetFunctionDecl mirrors CheckJSDoc.getFunctionDecl.
// Returns the FUNCTION node associated with the declaration node n,
// or nil if none.
func cjdGetFunctionDecl(n *rhino.Node) *rhino.Node {
	if n == nil {
		return nil
	}
	switch n.Token() {
	case rhino.FUNCTION:
		return n
	case rhino.VAR, rhino.LET, rhino.CONST:
		// `var f = function() {}` — first child is NAME, its first child is FUNCTION.
		fc := n.FirstChild()
		if fc == nil {
			return nil
		}
		ffc := fc.FirstChild()
		if ffc != nil && ffc.Token() == rhino.FUNCTION {
			return ffc
		}
	case rhino.ASSIGN:
		// `x.y = function() {}` — LHS is qualified name, RHS is FUNCTION.
		lhs := n.FirstChild()
		rhs := n.LastChild()
		if lhs != nil && cjdIsQualifiedName(lhs) && rhs != nil && rhs.Token() == rhino.FUNCTION {
			return rhs
		}
	case rhino.CLASS_FIELD_DEF:
		// Go's class methods are FUNCTION nodes inside CLASS_FIELD_DEF.
		fc := n.FirstChild()
		if fc != nil && fc.Token() == rhino.FUNCTION {
			return fc
		}
	}
	// GETTER_DEF, SETTER_DEF, MEMBER_FUNCTION_DEF: TODO — not in Go AST.
	return nil
}

// cjdIsFunctionDeclaration mirrors NodeUtil.isFunctionDeclaration.
// A function declaration is a FUNCTION whose parent is SCRIPT or BLOCK (not
// an expression context).
func cjdIsFunctionDeclaration(n *rhino.Node) bool {
	if n == nil || n.Token() != rhino.FUNCTION {
		return false
	}
	p := n.Parent()
	if p == nil {
		return true
	}
	return p.Token() == rhino.SCRIPT || p.Token() == rhino.BLOCK
}

// cjdIsLValue mirrors NodeUtil.isLValue — returns true when n appears in
// a left-value position (as the target of an assignment or declaration).
func cjdIsLValue(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	p := n.Parent()
	if p == nil {
		return false
	}
	switch p.Token() {
	case rhino.VAR, rhino.LET, rhino.CONST:
		return true
	case rhino.ASSIGN, rhino.ASSIGN_ADD, rhino.ASSIGN_SUB, rhino.ASSIGN_MUL,
		rhino.ASSIGN_DIV, rhino.ASSIGN_MOD, rhino.ASSIGN_BITAND, rhino.ASSIGN_BITOR,
		rhino.ASSIGN_BITXOR, rhino.ASSIGN_LSH, rhino.ASSIGN_RSH, rhino.ASSIGN_URSH:
		// n is the LHS when it's the first child.
		kids := p.Children()
		return len(kids) > 0 && kids[0] == n
	case rhino.PARAM_LIST:
		return true
	case rhino.FOR_IN, rhino.FOR_OF:
		kids := p.Children()
		return len(kids) > 0 && kids[0] == n
	}
	return false
}

// cjdGetRootTarget mirrors NodeUtil.getRootTarget.
// Climbs destructuring patterns; for plain names returns n itself.
func cjdGetRootTarget(n *rhino.Node) *rhino.Node {
	// The Go parser doesn't produce destructuring, so n is always the root.
	return n
}

// cjdIsLhsOfAssign mirrors NodeUtil.isLhsOfAssign.
func cjdIsLhsOfAssign(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	p := n.Parent()
	if p == nil {
		return false
	}
	if p.Token() != rhino.ASSIGN {
		return false
	}
	kids := p.Children()
	return len(kids) > 0 && kids[0] == n
}

// cjdGetRValueOfLValue mirrors NodeUtil.getRValueOfLValue.
// For a declaration LHS name, returns its initialiser; for an assign LHS,
// returns the RHS.
func cjdGetRValueOfLValue(n *rhino.Node) *rhino.Node {
	if n == nil {
		return nil
	}
	p := n.Parent()
	if p == nil {
		return nil
	}
	switch p.Token() {
	case rhino.VAR, rhino.LET, rhino.CONST:
		// The NAME node's first child is the initialiser.
		return n.FirstChild()
	case rhino.ASSIGN:
		kids := p.Children()
		if len(kids) >= 2 && kids[0] == n {
			return kids[1]
		}
	}
	return nil
}

// cjdGetBestJSDocInfo mirrors NodeUtil.getBestJSDocInfo — returns the
// JSDoc on n, or climbs to the parent EXPR_RESULT / ASSIGN to find it.
func cjdGetBestJSDocInfo(n *rhino.Node) *rhino.JSDocInfo {
	if n == nil {
		return nil
	}
	if info := n.JSDocInfo(); info != nil {
		return info
	}
	p := n.Parent()
	if p == nil {
		return nil
	}
	switch p.Token() {
	case rhino.EXPR_RESULT, rhino.ASSIGN, rhino.VAR, rhino.LET, rhino.CONST:
		if info := p.JSDocInfo(); info != nil {
			return info
		}
	}
	return nil
}

