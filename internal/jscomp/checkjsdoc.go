// Port of com/google/javascript/jscomp/CheckJSDoc.java @ closure-compiler v20260415
// (subset: checks implementable against the Go AST and JSDocInfo surface;
// Java-only tokens MEMBER_FUNCTION_DEF / GETTER_DEF / SETTER_DEF /
// COMPUTED_PROP / CAST / EXPORT / WITH / REST / DEFAULT_VALUE are absent from
// the Go parser, so those branches are omitted and noted with TODO comments).

package jscomp

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CheckJSDoc mirrors CheckJSDoc.java's CompilerPass.
// Java walks the AST using NodeTraversal (post-order); here we do a
// plain recursive walk so we avoid coupling to the NodeTraversal
// interface for a pass that doesn't need scope.
type CheckJSDoc struct {
	path      string
	src       string // raw source text (for checkJsDocInBlockComments)
	inExterns bool   // TODO: wire real externs tracking when externs are loaded
	diags     []diag.Diagnostic
}

// NewCheckJSDoc mirrors the Java constructor CheckJSDoc(AbstractCompiler).
// Go-only: language-adapter — takes path and src; Java uses AbstractCompiler which provides both.
func NewCheckJSDoc(path, src string) *CheckJSDoc {
	return &CheckJSDoc{path: path, src: src}
}

// Check is the entry point; mirrors CompilerPass.process(externs, root).
// In the Java pass, process traverses the externs tree with inExterns=true
// and the source tree with inExterns=false. The Go pass receives a single
// SCRIPT node (the source); a future extension can pass an externs node.
func (c *CheckJSDoc) Check(script *rhino.Node) []diag.Diagnostic {
	c.diags = nil
	c.inExterns = false
	if c.src != "" {
		c.checkJsDocInBlockComments()
	}
	c.cjdWalk(script)
	return c.diags
}

// checkJsDocInBlockComments mirrors CheckJSDoc.checkJsDocInBlockComments.
// Scans the raw source for block comments (/* ... */, not JSDoc /** ... */)
// that contain JSDoc-style @annotation patterns. Fires JSC_JSDOC_IN_BLOCK_COMMENT
// when found, suggesting the user use /** instead of /*.
// Java uses a regex: (/|(\n[ \t]*))\*[ \t]*@[a-zA-Z]+[ \t\n{]
func (c *CheckJSDoc) checkJsDocInBlockComments() {
	s := c.src
	i := 0
	for i < len(s) {
		// Skip JSDoc comments /** ... */
		if i+2 < len(s) && s[i] == '/' && s[i+1] == '*' && s[i+2] == '*' {
			if n := strings.Index(s[i+3:], "*/"); n >= 0 {
				i += n + 5
			} else {
				return
			}
			continue
		}
		// Find block comments /* ... */
		if i+1 < len(s) && s[i] == '/' && s[i+1] == '*' {
			commentStart := i
			end := strings.Index(s[i+2:], "*/")
			if end < 0 {
				return
			}
			commentEnd := i + 2 + end
			body := s[i+2 : commentEnd]
			// Check for JSDoc pattern: (*[ \t]*@[a-zA-Z]+[ \t\n{]) or at line start
			if cjdHasJSDocPattern(body) {
				line, col := cjdLineCol(s, commentStart)
				c.diags = append(c.diags, diag.Diagnostic{
					Path:    c.path,
					Line:    line,
					Col:     col,
					Code:    diag.JSC_JSDOC_IN_BLOCK_COMMENT,
					Level:   diag.Warning,
					Message: "Non-JSDoc comment has annotations. Did you mean to start it with '/**'?",
				})
			}
			i = commentEnd + 2
			continue
		}
		// Skip single-line comments
		if i+1 < len(s) && s[i] == '/' && s[i+1] == '/' {
			if n := strings.Index(s[i:], "\n"); n >= 0 {
				i += n + 1
			} else {
				return
			}
			continue
		}
		i++
	}
}

// cjdHasJSDocPattern mirrors Java's COMMENT_PATTERN:
// (/|(\n[ \t]*))\*[ \t]*@[a-zA-Z]+[ \t\n{]
// Checks if a block comment body contains JSDoc-style annotations.
func cjdHasJSDocPattern(body string) bool {
	for i := 0; i < len(body); i++ {
		ch := body[i]
		// Look for `*[ \t]*@[a-zA-Z]+` after a `/` or a newline+optional-space.
		matchStar := false
		if ch == '/' || ch == '\n' {
			j := i + 1
			// Skip optional whitespace after newline.
			if ch == '\n' {
				for j < len(body) && (body[j] == ' ' || body[j] == '\t') {
					j++
				}
			}
			// Must be `*`.
			if j < len(body) && body[j] == '*' {
				matchStar = true
				j++
			}
			if !matchStar {
				continue
			}
			// Skip optional whitespace.
			for j < len(body) && (body[j] == ' ' || body[j] == '\t') {
				j++
			}
			// Must be `@` followed by at least one letter.
			if j >= len(body) || body[j] != '@' {
				continue
			}
			j++ // consume '@'
			if j >= len(body) || !((body[j] >= 'a' && body[j] <= 'z') || (body[j] >= 'A' && body[j] <= 'Z')) {
				continue
			}
			// Read the annotation name.
			for j < len(body) && ((body[j] >= 'a' && body[j] <= 'z') || (body[j] >= 'A' && body[j] <= 'Z')) {
				j++
			}
			// Must be followed by space, tab, newline, or `{`.
			if j < len(body) && (body[j] == ' ' || body[j] == '\t' || body[j] == '\n' || body[j] == '{') {
				return true
			}
		}
	}
	return false
}

// cjdLineCol returns (line, col) for the absolute offset in src.
// Go-only: language-adapter — Java uses Node.getLineno/getCharno; Go uses raw offset.
func cjdLineCol(src string, off int) (int, int) {
	if off > len(src) {
		off = len(src)
	}
	line := 1
	col := 0
	for i := 0; i < off; i++ {
		if src[i] == '\n' {
			line++
			col = 0
		} else {
			col++
		}
	}
	return line, col
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
	c.validateNoCollapse(n, info)
	c.validateClassLevelJsDoc(n, info)
	c.validateArrowFunction(n)
	c.validateRestParameter(n)
	// validateDefaultValue: TODO — DEFAULT_VALUE token absent from Go parser (default params produce NAME, not DEFAULT_VALUE).
	c.validateTemplates(n, info)
	c.validateTypedefs(n, info)
	c.validateNoSideEffects(n, info)
	c.validateAbstractJsDoc(n, info)
	c.validateDefinesDeclaration(n, info)
	c.validateSuppress(n, info)
	c.validateImplicitCast(n, info)
	c.validateClosurePrimitive(n, info)
	c.validateReturnJsDoc(n, info)
	c.validateTsType(n, info)
	// validateJsDocTypeNames: handled by CheckJsDocTypes.
	c.validateIsUsedViaDotConstructor(n, info)
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

// containsOnlySuppressionFor mirrors Java's containsOnlySuppressionFor.
func containsOnlySuppressionFor(info *rhino.JSDocInfo, allowed string) bool {
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
	if isClassDecl(n) {
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
	if isClassDecl(n) {
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
	if info == nil || !hasClassLevelJsDoc(info) {
		return
	}
	// Legacy AST: FUNCTION → BLOCK → CLASS.
	// New pipeline AST: FUNCTION → MEMBER_FUNCTION_DEF → CLASS_MEMBERS → CLASS.
	if n.Token() != rhino.FUNCTION && n.Token() != rhino.MEMBER_FUNCTION_DEF {
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
	// Legacy: FUNCTION inside CLASS body BLOCK.
	if n.Token() == rhino.FUNCTION && p.Token() == rhino.BLOCK && gp.Token() == rhino.CLASS {
		c.report(n, diag.JSC_DISALLOWED_MEMBER_JSDOC, diag.Warning,
			"Class level JSDocs (@interface, @extends, etc.) are not allowed on class members")
		return
	}
	// New pipeline: MEMBER_FUNCTION_DEF inside CLASS_MEMBERS inside CLASS.
	if n.Token() == rhino.MEMBER_FUNCTION_DEF && p.Token() == rhino.CLASS_MEMBERS && gp.Token() == rhino.CLASS {
		c.report(n, diag.JSC_DISALLOWED_MEMBER_JSDOC, diag.Warning,
			"Class level JSDocs (@interface, @extends, etc.) are not allowed on class members")
	}
}

// hasClassLevelJsDoc mirrors the private Java helper hasClassLevelJsDoc.
func hasClassLevelJsDoc(info *rhino.JSDocInfo) bool {
	return info.IsConstructor() || info.IsInterface() ||
		info.BaseType() != nil ||
		len(info.Implementations()) != 0 ||
		len(info.BaseTypes()) > 1 // multiple @extends = interface extends
}

// ---------------------------------------------------------------------------
// validateAbstractJsDoc
// Mirrors CheckJSDoc.validateAbstractJsDoc.
// Fires MISPLACED_ANNOTATION when @abstract appears in an invalid position:
//   - On a non-class, non-function node.
//   - On an ES6 constructor member.
//   - On a function with a non-empty body.
//   - On a static method.
//   - On a function that is not a method.
// ---------------------------------------------------------------------------
func (c *CheckJSDoc) validateAbstractJsDoc(n *rhino.Node, info *rhino.JSDocInfo) {
	if info == nil || !info.IsAbstract() {
		return
	}
	// Class declarations are always valid for @abstract.
	if isClassDecl(n) {
		return
	}
	// @abstract on a non-function node.
	functionNode := cjdGetFunctionDecl(n)
	if functionNode == nil {
		c.reportMisplaced(n, "@abstract",
			"only functions or non-static methods can be abstract")
		return
	}
	// @abstract on a function with a non-empty body (unless it's a constructor).
	if !info.IsConstructor() {
		body := functionNode.LastChild()
		if body != nil && body.Token() == rhino.BLOCK && body.HasChildren() {
			c.reportMisplaced(n, "@abstract",
				"function with a non-empty body cannot be abstract")
			return
		}
	}
	// @abstract on a static member is invalid.
	// Check both the AST property (for MEMBER_FUNCTION_DEF) and the JSDoc
	// flag (for FUNCTION with `static` parsed into JSDocInfo).
	if n.GetBooleanProp(rhino.PropStaticMember) || info.IsStatic() {
		c.reportMisplaced(n, "@abstract", "static methods cannot be abstract")
		return
	}
	// @abstract is only valid on methods (MEMBER_FUNCTION_DEF, STRING_KEY in
	// object literal, GETTER_DEF, SETTER_DEF, FUNCTION inside a CLASS body,
	// or a constructor) or prototype functions.
	// Java's check: !info.isConstructor() && !n.isMemberFunctionDef() && ...
	// Skip the method check when @constructor is present.
	if !info.IsConstructor() {
		switch n.Token() {
		case rhino.MEMBER_FUNCTION_DEF, rhino.GETTER_DEF, rhino.SETTER_DEF, rhino.STRING_KEY:
			// Methods in class body or object literal — valid.
		case rhino.FUNCTION:
			// FUNCTION directly inside a CLASS body BLOCK (legacy AST for class
			// methods — no MEMBER_FUNCTION_DEF wrapper) is a class method — valid.
			if p := n.Parent(); p != nil && p.Token() == rhino.BLOCK {
				if gp := p.Parent(); gp != nil && gp.Token() == rhino.CLASS {
					return // class method — valid
				}
			}
			// Check if it's a prototype method: parent is ASSIGN whose LHS is
			// X.prototype.method (GETPROP chain ending in "prototype").
			if !cjdIsPrototypeMethod(functionNode) {
				c.reportMisplaced(n, "@abstract",
					"only functions or non-static methods can be abstract")
			}
		}
	}
}

// cjdIsPrototypeMethod reports whether fn is assigned as a prototype method
// (i.e. the enclosing ASSIGN's LHS is X.prototype.y). Mirrors
// NodeUtil.isPrototypeMethod in Java.
func cjdIsPrototypeMethod(fn *rhino.Node) bool {
	if fn == nil {
		return false
	}
	p := fn.Parent()
	if p == nil {
		return false
	}
	if p.Token() != rhino.ASSIGN {
		return false
	}
	lhs := p.FirstChild()
	return lhs != nil && cjdIsPrototypeProperty(lhs)
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
	if isJSDocOnFunctionNode(n, info) {
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

// isJSDocOnFunctionNode mirrors Java's isJSDocOnFunctionNode.
// Returns true when the node can be the host for function-bearing JSDoc.
// Java's switch includes FUNCTION, GETTER_DEF, SETTER_DEF, MEMBER_FUNCTION_DEF,
// STRING_KEY, COMPUTED_PROP, EXPORT, MEMBER_FIELD_DEF, COMPUTED_FIELD_DEF,
// GETELEM/GETPROP (when first child is qualified), VAR/LET/CONST/ASSIGN.
// Go subset: FUNCTION, CLASS (for @constructor/@interface on class decls),
// STRING_KEY, CLASS_FIELD_DEF, VAR/LET/CONST, ASSIGN,
// GETPROP (when qualified and parent is EXPR_RESULT).
func isJSDocOnFunctionNode(n *rhino.Node, info *rhino.JSDocInfo) bool {
	switch n.Token() {
	case rhino.FUNCTION, rhino.CLASS, rhino.STRING_KEY, rhino.CLASS_FIELD_DEF,
		rhino.MEMBER_FUNCTION_DEF, rhino.GETTER_DEF, rhino.SETTER_DEF:
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
// Fires JSC_MISPLACED_MSG_ANNOTATION when @desc/@meaning/@alternateMessageId
// appear on a node that is not a valid message declaration (name starts with
// MSG_). Java also checks info.getDescription() != null, but only for non-TS
// source; we mirror that by always checking description.
// ---------------------------------------------------------------------------
func (c *CheckJSDoc) validateMsgJsDoc(n *rhino.Node, info *rhino.JSDocInfo) {
	if info == nil {
		return
	}
	hasNonDescMsgTag := info.GetMeaning() != "" || info.GetAlternateMessageId() != ""
	hasDesc := info.GetDescription() != ""
	if !hasNonDescMsgTag && !hasDesc {
		return
	}
	descOkay := false
	switch n.Token() {
	case rhino.ASSIGN, rhino.VAR, rhino.LET, rhino.CONST:
		if fc := n.FirstChild(); fc != nil {
			descOkay = isValidMsgName(fc)
		}
	case rhino.STRING_KEY:
		descOkay = isValidMsgName(n)
	case rhino.GETPROP:
		// Java: only when n.isFromExterns() && n.isQualifiedName().
		// Go: only for qualified name in externs context (inExterns).
		if c.inExterns && cjdIsQualifiedName(n) {
			descOkay = isValidMsgName(n)
		}
	}
	if !descOkay {
		c.report(n, diag.JSC_MISPLACED_MSG_ANNOTATION, diag.Warning,
			"Misplaced message annotation. @desc, @meaning, and @alternateMessageId"+
				" annotations should be only on message nodes."+
				"\nMessage constants must be prefixed with 'MSG_'.")
	}
}

// isValidMsgName mirrors CheckJSDoc.isValidMsgName.
// Returns true when the node's name starts with "MSG_".
func isValidMsgName(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	switch n.Token() {
	case rhino.NAME, rhino.STRING_KEY:
		return strings.HasPrefix(n.String_(), "MSG_")
	case rhino.GETPROP:
		// Qualified name: check the full dotted string.
		return strings.HasPrefix(n.GetQualifiedName(), "MSG_")
	}
	return false
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
	case rhino.STRING_KEY, rhino.CLASS_FIELD_DEF, rhino.MEMBER_FIELD_DEF:
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
// validateRestParameter
// Mirrors CheckJSDoc.validateRestParameter.
// Fires JSC_BAD_REST_PARAMETER_ANNOTATION when a rest parameter has a JSDoc
// type annotation that is not wrapped in `...` (i.e. the annotation is
// `{number}` instead of `{...number}`). Java uses Token.ITER_REST on the
// type-expression root to detect this; Go checks IsRestParam.
// ---------------------------------------------------------------------------
func (c *CheckJSDoc) validateRestParameter(n *rhino.Node) {
	// Java: if (!restParam.isRest() || !restParam.getParent().isParamList())
	if n.Token() != rhino.ITER_REST {
		return
	}
	p := n.Parent()
	if p == nil || p.Token() != rhino.PARAM_LIST {
		return
	}
	// The first child of ITER_REST is the bound NAME (or destructuring pattern).
	firstChild := n.FirstChild()
	if firstChild == nil {
		return
	}
	// Prefer inline JSDoc on the rest-param child.
	var paramTypeAnnotation *rhino.JSTypeExpression
	if inlineInfo := firstChild.JSDocInfo(); inlineInfo != nil {
		paramTypeAnnotation = inlineInfo.Type()
	} else {
		// Fall back to the enclosing function's JSDoc @param annotation.
		fn := p.Parent() // PARAM_LIST → FUNCTION
		if fn != nil {
			fnInfo := cjdGetBestJSDocInfo(fn)
			if fnInfo != nil && firstChild.Token() == rhino.NAME {
				paramTypeAnnotation = fnInfo.ParameterType(firstChild.String_())
				// If no named param found, check IsRestParam.
				if paramTypeAnnotation != nil && fnInfo.IsRestParam(firstChild.String_()) {
					// Already a rest param annotation — valid.
					return
				}
			}
		}
	}
	// Java fires BAD_REST_PARAMETER_ANNOTATION when the annotation root token
	// is NOT Token.ITER_REST. In Go we approximate: if the annotation is
	// present and does NOT carry the rest-param marker, it's bad.
	if paramTypeAnnotation == nil {
		return
	}
	// If the type name starts with `...` it was parsed as a rest annotation.
	typeName := paramTypeAnnotation.TypeName
	if len(typeName) >= 3 && typeName[:3] == "..." {
		return
	}
	c.report(n, diag.JSC_BAD_REST_PARAMETER_ANNOTATION, diag.Warning,
		"Missing \"...\" in type annotation for rest parameter.")
}

// ---------------------------------------------------------------------------
// validateNoSideEffects
// Mirrors CheckJSDoc.validateNoSideEffects.
// Fires JSC_INVALID_MODIFIES_ANNOTATION when @modifies appears outside externs.
// ---------------------------------------------------------------------------
func (c *CheckJSDoc) validateNoSideEffects(n *rhino.Node, info *rhino.JSDocInfo) {
	if info == nil {
		return
	}
	if c.inExterns {
		return
	}
	if info.HasSideEffectsArgumentsAnnotation() || info.ModifiesThis() {
		c.report(n, diag.JSC_INVALID_MODIFIES_ANNOTATION, diag.Error,
			"@modifies may only appear in externs files.")
	}
}

// ---------------------------------------------------------------------------
// validateDefinesDeclaration
// Mirrors CheckJSDoc.validateDefinesDeclaration.
// Fires JSC_INVALID_DEFINE_ON_LET when @define appears on a let declaration.
// ---------------------------------------------------------------------------
func (c *CheckJSDoc) validateDefinesDeclaration(n *rhino.Node, info *rhino.JSDocInfo) {
	if info != nil && info.IsDefine() && n.Token() == rhino.LET {
		c.report(n, diag.JSC_INVALID_DEFINE_ON_LET, diag.Error,
			"variables annotated with @define may only be declared with VARs, ASSIGNs, or CONSTs")
	}
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
// validateImplicitCast
// Mirrors CheckJSDoc.validateImplicitCast.
// Fires JSC_ILLEGAL_IMPLICIT_CAST when @implicitCast appears outside externs.
// ---------------------------------------------------------------------------
func (c *CheckJSDoc) validateImplicitCast(n *rhino.Node, info *rhino.JSDocInfo) {
	if !c.inExterns && info != nil && info.IsImplicitCast() {
		// Java reports at the node `n`. For GETPROP nodes (e.g. the JSDoc
		// is on `Element.prototype.innerHTML`), the property name STRINGLIT
		// holds the position Java uses — not the GETPROP itself.
		// Mirror that by using the last child (property name) of GETPROP.
		reportNode := n
		if n.Token() == rhino.GETPROP {
			if last := n.LastChild(); last != nil {
				reportNode = last
			}
		}
		c.report(reportNode, diag.JSC_ILLEGAL_IMPLICIT_CAST, diag.Error,
			"Illegal annotation on {0}. @implicitCast may only be used in externs.")
	}
}

// ---------------------------------------------------------------------------
// validateClosurePrimitive
// Mirrors CheckJSDoc.validateClosurePrimitive.
// Fires MISPLACED_ANNOTATION when @closurePrimitive is not on a function node.
// ---------------------------------------------------------------------------
func (c *CheckJSDoc) validateClosurePrimitive(n *rhino.Node, info *rhino.JSDocInfo) {
	if info == nil || !info.HasClosurePrimitiveId() {
		return
	}
	if !isJSDocOnFunctionNode(n, info) {
		c.reportMisplaced(n, "closurePrimitive", "must be on a function node")
	}
}

// ---------------------------------------------------------------------------
// validateNoCollapse
// Mirrors CheckJSDoc.validateNoCollapse.
// Fires MISPLACED_ANNOTATION when @nocollapse is in an invalid position:
//   - In externs (no effect there).
//   - On prototype properties or non-static class fields.
//   - On an assignment whose LHS root is not a NAME.
//
// ---------------------------------------------------------------------------
func (c *CheckJSDoc) validateNoCollapse(n *rhino.Node, info *rhino.JSDocInfo) {
	if info == nil || !info.IsNoCollapse() {
		return
	}
	if c.inExterns {
		c.reportMisplaced(n, "nocollapse", "This JSDoc has no effect in externs.")
		return
	}
	// isPrototypePropertyDeclaration: parent of n is a prototype property.
	p := n.Parent()
	if p != nil && cjdIsPrototypeProperty(p) {
		c.reportMisplaced(n, "nocollapse",
			"This JSDoc has no effect on prototype properties and non-static fields.")
		return
	}
	// Non-static class field: parent is CLASS_MEMBERS block and n is non-static.
	if p != nil && p.Token() == rhino.BLOCK {
		gp := p.Parent()
		if gp != nil && gp.Token() == rhino.CLASS {
			if !n.GetBooleanProp(rhino.PropStaticMember) {
				c.reportMisplaced(n, "nocollapse",
					"This JSDoc has no effect on prototype properties and non-static fields.")
				return
			}
		}
	}
	if n.Token() == rhino.ASSIGN {
		assignee := n.FirstChild()
		if assignee != nil && cjdIsQualifiedName(assignee) {
			root := cjdGetRootOfQualifiedName(assignee)
			if root != nil && root.Token() != rhino.NAME {
				c.reportMisplaced(n, "nocollapse", "This JSDoc has no effect.")
			}
		}
	}
}

// ---------------------------------------------------------------------------
// validateIsUsedViaDotConstructor
// Mirrors CheckJSDoc.validateIsUsedViaDotConstructor.
// Fires MISPLACED_ANNOTATION when @usedViaDotConstructor is not on a
// constructor FUNCTION or ES6 constructor member.
// ---------------------------------------------------------------------------
func (c *CheckJSDoc) validateIsUsedViaDotConstructor(n *rhino.Node, info *rhino.JSDocInfo) {
	if info == nil || !info.IsUsedViaDotConstructor() {
		return
	}
	if n.Token() == rhino.FUNCTION && info.IsConstructor() {
		return
	}
	c.reportMisplaced(n, "usedViaDotConstructor", "must be on a constructor")
}

// ---------------------------------------------------------------------------
// validateTsType
// Mirrors CheckJSDoc.validateTsType.
// Fires MISPLACED_ANNOTATION when @tsType is not on a function node.
// ---------------------------------------------------------------------------
func (c *CheckJSDoc) validateTsType(n *rhino.Node, info *rhino.JSDocInfo) {
	if info == nil || len(info.GetTsTypes()) == 0 {
		return
	}
	if !isJSDocOnFunctionNode(n, info) {
		c.reportMisplaced(n, "tsType", "must be on a function node")
	}
}

// isFromTs mirrors CheckJSDoc.isFromTs — returns whether the node's source is
// a TypeScript file. Go does not distinguish TS from JS sources, so this
// always returns false. Java calls n.getStaticSourceFile().isTypeScriptSource().
// Go-only: language-adapter — Java checks StaticSourceFile.isTypeScriptSource(); Go has no TS source-kind.
func isFromTs(n *rhino.Node) bool {
	// Go jsgo does not load TypeScript sources — all files are treated as JS.
	// Return false to match the default (non-TS) path in Java.
	_ = n // n is not used; kept for API symmetry with Java.
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

// isClassDecl mirrors CheckJSDoc.isClassDecl.
func isClassDecl(n *rhino.Node) bool {
	return cjdIsClassNode(n) ||
		(n.Token() == rhino.ASSIGN && n.LastChild() != nil && cjdIsClassNode(n.LastChild())) ||
		(cjdIsNameDeclaration(n) && n.FirstChild() != nil &&
			isNameInitializeWithClass(n.FirstChild())) ||
		isNameInitializeWithClass(n)
}

// isNameInitializeWithClass mirrors CheckJSDoc.isNameInitializeWithClass.
func isNameInitializeWithClass(n *rhino.Node) bool {
	return n != nil && n.Token() == rhino.NAME &&
		n.FirstChild() != nil && isClassDecl(n.FirstChild())
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
	case rhino.MEMBER_FUNCTION_DEF, rhino.GETTER_DEF, rhino.SETTER_DEF:
		// Class method definitions: MEMBER_FUNCTION_DEF wraps a FUNCTION node.
		fc := n.FirstChild()
		if fc != nil && fc.Token() == rhino.FUNCTION {
			return fc
		}
	}
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
