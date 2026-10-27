package refactoring

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/testjunit"
)

type JUnitTestCase = testjunit.JUnitTestCase

func RunJUnit4(t *testing.T, newTest func() JUnitTestCase) {
	testjunit.RunJUnit4(t, newTest)
}
