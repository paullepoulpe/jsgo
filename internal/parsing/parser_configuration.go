// Port of com/google/javascript/jscomp/parsing/ParserConfiguration.java @ closure-compiler v20260415

package parsing

import (
	_ "embed"
	"sync"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Go-only: language-adapter — classpath resource via go:embed. This is
// com/google/javascript/jscomp/parsing/ParserConfig.properties, extracted
// byte for byte from the pinned closure-compiler.jar.
//
//go:embed ParserConfig.properties
var parserConfiguration_properties string

// Go-only: language-adapter — ResourceBundle.getBundle loads (and caches) the
// PropertyResourceBundle once; the mutex is the JDK's own synchronization of
// that cache.
var (
	parserConfiguration_bundleMu sync.Mutex
	parserConfiguration_bundle   *jutil.LinkedHashMap[string, string]
)

// ParserConfiguration reads the parser's ParserConfig.properties resource.
//
// Java: final class ParserConfiguration
type ParserConfiguration struct{}

// ParserConfiguration_GetString is
// ResourceBundle.getBundle("com.google.javascript.jscomp.parsing.ParserConfig").getString(name).
//
// Java: ParserConfiguration.getString(String)
func ParserConfiguration_GetString(name string) string {
	parserConfiguration_bundleMu.Lock()
	if parserConfiguration_bundle == nil {
		parserConfiguration_bundle = jutil.PropertiesLoad(parserConfiguration_properties)
	}
	bundle := parserConfiguration_bundle
	parserConfiguration_bundleMu.Unlock()
	value, ok := bundle.Lookup(name)
	if !ok {
		panic(jutil.NewMissingResourceException(
			"Can't find resource for bundle java.util.PropertyResourceBundle, key "+name,
			"java.util.PropertyResourceBundle", name))
	}
	return value
}

// Java: ParserConfiguration()
func NewParserConfiguration() *ParserConfiguration { return &ParserConfiguration{} }
