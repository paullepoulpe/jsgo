// Port of com/google/javascript/jscomp/deps/BrowserWithTransformedPrefixesModuleResolver.java @ closure-compiler v20260415

package jscomp

import (
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// BrowserWithTransformedPrefixesModuleResolver is a limited superset of the
// BrowserModuleResolver that allows for replacing some path prefixes before
// resolving.
type BrowserWithTransformedPrefixesModuleResolver struct {
	ModuleResolverBase
	prefixReplacements []*BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement
	expectedPrefixes   string
}

var _ ModuleResolver = (*BrowserWithTransformedPrefixesModuleResolver)(nil)

var BrowserWithTransformedPrefixesModuleResolver_INVALID_AMBIGUOUS_PATH = DiagnosticType_Error(
	"JSC_INVALID_AMBIGUOUS_PATH",
	"The path \"{0}\" is invalid. Expected any of the following prefixes for non-relative "+
		"paths: {1}.")

// BrowserWithTransformedPrefixesModuleResolver_Factory is the factory for
// BrowserWithTransformedPrefixesModuleResolver.
type BrowserWithTransformedPrefixesModuleResolver_Factory struct {
	prefixReplacements *jutil.LinkedHashMap[string, string]
}

var _ ModuleLoader_ModuleResolverFactory = (*BrowserWithTransformedPrefixesModuleResolver_Factory)(nil)

// Java: BrowserWithTransformedPrefixesModuleResolver.Factory(ImmutableMap<String, String>)
func NewBrowserWithTransformedPrefixesModuleResolver_Factory(prefixReplacements *jutil.LinkedHashMap[string, string]) *BrowserWithTransformedPrefixesModuleResolver_Factory {
	return &BrowserWithTransformedPrefixesModuleResolver_Factory{prefixReplacements: prefixReplacements}
}

// Java: BrowserWithTransformedPrefixesModuleResolver.Factory.create(ImmutableSet<String>, ImmutableList<String>, ErrorHandler, PathEscaper)
func (f *BrowserWithTransformedPrefixesModuleResolver_Factory) Create(
	modulePaths *jutil.LinkedHashSet[string],
	moduleRootPaths []string,
	errorHandler ErrorHandler,
	pathEscaper ModuleLoader_PathEscaper) ModuleResolver {
	return NewBrowserWithTransformedPrefixesModuleResolver(
		modulePaths, moduleRootPaths, errorHandler, pathEscaper, f.prefixReplacements)
}

// BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement is the
// struct of prefix and replacement. Has a natural ordering from longest
// prefix to shortest prefix.
type BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement struct {
	prefix      string
	replacement string
}

// Java: BrowserWithTransformedPrefixesModuleResolver.PrefixReplacement.of(String, String)
func BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement_Of(prefix string, replacement string) *BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement {
	return &BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement{prefix: prefix, replacement: replacement}
}

// Java: BrowserWithTransformedPrefixesModuleResolver.PrefixReplacement.prefix()
func (r *BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement) Prefix() string {
	return r.prefix
}

// Java: BrowserWithTransformedPrefixesModuleResolver.PrefixReplacement.replacement()
func (r *BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement) Replacement() string {
	return r.replacement
}

// Equals is the record's componentwise equals.
func (r *BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement) Equals(o any) bool {
	other, ok := o.(*BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement)
	return ok && other != nil && r.prefix == other.prefix && r.replacement == other.replacement
}

// HashCode is the record's hashCode.
func (r *BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement) HashCode() int32 {
	return 31*jutil.StringHashCode(r.prefix) + jutil.StringHashCode(r.replacement)
}

// String is the record's toString.
func (r *BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement) String() string {
	return "PrefixReplacement[prefix=" + r.prefix + ", replacement=" + r.replacement + "]"
}

// Java: BrowserWithTransformedPrefixesModuleResolver(ImmutableSet<String>, ImmutableList<String>, ErrorHandler, PathEscaper, ImmutableMap<String, String>)
func NewBrowserWithTransformedPrefixesModuleResolver(
	modulePaths *jutil.LinkedHashSet[string],
	moduleRootPaths []string,
	errorHandler ErrorHandler,
	pathEscaper ModuleLoader_PathEscaper,
	prefixReplacements *jutil.LinkedHashMap[string, string]) *BrowserWithTransformedPrefixesModuleResolver {
	r := &BrowserWithTransformedPrefixesModuleResolver{}
	r.InitModuleResolverBase(modulePaths, moduleRootPaths, errorHandler, pathEscaper)
	var replacements []*BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement
	for k, v := range prefixReplacements.All() {
		replacements = append(replacements, BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement_Of(k, v))
	}
	// Sort by length in descending order to prefixes are applied most specific to
	// least specific.
	slices.SortStableFunc(replacements, func(a, b *BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement) int {
		la, lb := jutil.StringLength(a.prefix), jutil.StringLength(b.prefix)
		if la != lb {
			if lb < la {
				return -1
			}
			return 1
		}
		return jutil.StringCompareTo(a.prefix, b.prefix)
	})
	// toImmutableSet(): the map keys are distinct, so are the records.
	r.prefixReplacements = replacements
	var prefixes []string
	for _, p := range r.prefixReplacements {
		prefixes = append(prefixes, p.prefix)
	}
	slices.SortStableFunc(prefixes, jutil.StringCompareTo)
	r.expectedPrefixes = strings.Join(prefixes, ", ")
	return r
}

// replacePathPrefixes may return nil when moduleAddress is nil in Java;
// Go strings are never null.
//
// Java: BrowserWithTransformedPrefixesModuleResolver.replacePathPrefixes(String)
func (r *BrowserWithTransformedPrefixesModuleResolver) replacePathPrefixes(moduleAddress string) string {
	transformedAddress := moduleAddress
	for _, prefixReplacement := range r.prefixReplacements {
		if strings.HasPrefix(moduleAddress, prefixReplacement.Prefix()) {
			transformedAddress =
				prefixReplacement.Replacement() +
					moduleAddress[len(prefixReplacement.Prefix()):]
			break
		}
	}
	return transformedAddress
}

// ResolveJsModule may return nil.
//
// Java: BrowserWithTransformedPrefixesModuleResolver.resolveJsModule(String, String, String, int, int)
func (r *BrowserWithTransformedPrefixesModuleResolver) ResolveJsModule(
	scriptAddress string, moduleAddress string, sourcename string, lineno int, colno int) *string {
	transformedAddress := r.replacePathPrefixes(moduleAddress)
	replaced := transformedAddress != moduleAddress
	// If ambiguous after the loop it was not transformed and the original moduleAddress is
	// ambiguous with an unrecognized prefix. Allow transformed paths to be ambiguous after
	// transformation to maybe match how files are passed to the compiler.
	if !replaced && ModuleLoader_IsAmbiguousIdentifier(transformedAddress) {
		r.errorHandler.Report(
			CheckLevel_WARNING.Ptr(),
			JSError_MakeWithSourceLocation(
				&sourcename,
				lineno,
				colno,
				BrowserWithTransformedPrefixesModuleResolver_INVALID_AMBIGUOUS_PATH,
				transformedAddress,
				r.expectedPrefixes))
		return nil
	}
	// A null load address is not an error for this resolver, because it is expected that the
	// module may be supplied at runtime.
	return r.locate(scriptAddress, transformedAddress)
}

// ResolveJsModuleSilently may return nil.
//
// Java: BrowserWithTransformedPrefixesModuleResolver.resolveJsModuleSilently(String, String)
func (r *BrowserWithTransformedPrefixesModuleResolver) ResolveJsModuleSilently(scriptAddress string, moduleAddress string) *string {
	transformedAddress := r.replacePathPrefixes(moduleAddress)
	replaced := transformedAddress != moduleAddress
	// If ambiguous after the loop it was not transformed and the original moduleAddress is
	// ambiguous with an unrecognized prefix. Allow transformed paths to be ambiguous after
	// transformation to maybe match how files are passed to the compiler.
	if !replaced && ModuleLoader_IsAmbiguousIdentifier(transformedAddress) {
		return nil
	}
	return r.locate(scriptAddress, transformedAddress)
}

// Java: BrowserWithTransformedPrefixesModuleResolver.resolveModuleAsPath(String, String)
func (r *BrowserWithTransformedPrefixesModuleResolver) ResolveModuleAsPath(scriptAddress string, moduleAddress string) string {
	if ModuleLoader_IsRelativeIdentifier(moduleAddress) {
		return r.ModuleResolverBase.ResolveModuleAsPath(scriptAddress, moduleAddress)
	}
	transformedAddress := moduleAddress
	for _, prefixReplacement := range r.prefixReplacements {
		if strings.HasPrefix(moduleAddress, prefixReplacement.Prefix()) {
			transformedAddress =
				prefixReplacement.Replacement() +
					moduleAddress[len(prefixReplacement.Prefix()):]
			break
		}
	}
	return ModuleLoader_Normalize(transformedAddress, r.moduleRootPaths)
}
