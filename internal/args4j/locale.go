package args4j

import (
	"sync"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Go-only: language-adapter — java.util.Locale and ResourceBundle, as
// args4j's message bundles use them. The port runs as Java does in an
// English default locale: ResourceBundle.getBundle(name, locale) then finds
// the bundle name_en, which is empty in args4j, and falls back to the base
// bundle name.properties for every key. The base bundles are embedded by
// messages.go and spi_messages.go (messages.properties and
// spi_messages.properties: the entries as the pinned closure-compiler.jar
// shades them, under a comment with args4j's notice); the German and Russian
// bundles are not ported, so every locale resolves to the base bundle.

// Locale is java.util.Locale.
type Locale struct {
	tag string
}

// Locale_ENGLISH is Locale.ENGLISH.
var Locale_ENGLISH = &Locale{tag: "en"}

// Locale_GetDefault is Locale.getDefault().
func Locale_GetDefault() *Locale { return Locale_ENGLISH }

// ResourceBundle is the java.util.ResourceBundle that CmdLineParser may be
// given to localize usage texts and meta variables: GetString is
// ResourceBundle.getString(String) (it throws MissingResourceException for
// an unknown key).
type ResourceBundle interface {
	GetString(key string) string
}

// propertyBundle is a PropertyResourceBundle loaded once, like
// ResourceBundle.getBundle's cache.
type propertyBundle struct {
	text    string
	once    sync.Once
	entries *jutil.LinkedHashMap[string, string]
}

// GetString is ResourceBundle.getString(String).
func (b *propertyBundle) GetString(key string) string {
	b.once.Do(func() { b.entries = jutil.PropertiesLoad(b.text) })
	value, ok := b.entries.Lookup(key)
	if !ok {
		panic(jutil.NewMissingResourceException(
			"Can't find resource for bundle java.util.PropertyResourceBundle, key "+key,
			"java.util.PropertyResourceBundle", key))
	}
	return value
}
