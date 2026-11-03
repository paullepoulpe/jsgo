package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type j2clChecksPassTest struct {
	CompilerTestCaseBase
}

func newJ2clChecksPassTest() *j2clChecksPassTest {
	c := &j2clChecksPassTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func (c *j2clChecksPassTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableTypeCheck()
}

func (c *j2clChecksPassTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewJ2clChecksPass(compiler)
}

func (c *j2clChecksPassTest) TestWarnsOnBoxedReferenceEquality() {
	intSrc := SourceFile_FromCode(
		"java/lang/Integer.impl.java.js",
		"/** @constructor */ var Integer = function() {};")
	c.TestWarning(
		c.Srcs(
			intSrc,
			SourceFile_FromCode("test.js", "var a = new Integer(); var b = new Integer(); a == b;")),
		J2clChecksPass_J2CL_REFERENCE_EQUALITY)
	c.TestWarning(
		c.Srcs(
			intSrc,
			SourceFile_FromCode("test.js", "var a = new Integer(); a !== new Integer();")),
		J2clChecksPass_J2CL_REFERENCE_EQUALITY)
}

func (c *j2clChecksPassTest) TestAllowsNullAndUndefinedComparison() {
	intSrc := SourceFile_FromCode(
		"java/lang/Integer.impl.java.js",
		"/** @constructor */ var Integer = function() {};")
	c.TestNoWarning(
		c.Srcs(
			intSrc,
			SourceFile_FromCode("test.js", "var a = new Integer(); a == null; a === undefined;")))
}

func (c *j2clChecksPassTest) TestNilRootNoop() {
	comp := NewCompiler()
	comp.InitOptions(NewCompilerOptions())
	NewJ2clChecksPass(comp).Process(rhino.NilNode, rhino.NilNode)
}

func TestJ2clChecksPass(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return newJ2clChecksPassTest() })
}
