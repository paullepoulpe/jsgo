// Port of com/google/javascript/jscomp/bundle/CachedTransformer.java @ closure-compiler v20260415

package jscomp

import "github.com/paullepoulpe/jsgo/internal/jutil"

// CachedTransformer is a Transformer that caches output from a delegate
// transformer.
type CachedTransformer struct {
	cache *jutil.LoadingCache[*Source, *Source]
}

var _ Source_Transformer = (*CachedTransformer)(nil)

const cachedTransformer_DEFAULT_CACHE_SPEC = "maximumSize=10000"

// Java: CachedTransformer(Source.Transformer, CacheBuilder<? super Source, ? super Source>)
func NewCachedTransformerWithBuilder(delegate Source_Transformer, builder *jutil.CacheBuilder) *CachedTransformer {
	return &CachedTransformer{cache: jutil.CacheBuilder_Build(builder, delegate.Transform)}
}

// Java: CachedTransformer(Source.Transformer, String)
func NewCachedTransformerWithSpec(delegate Source_Transformer, spec string) *CachedTransformer {
	return NewCachedTransformerWithBuilder(delegate, jutil.CacheBuilder_From(spec))
}

// Java: CachedTransformer(Source.Transformer)
func NewCachedTransformer(delegate Source_Transformer) *CachedTransformer {
	return NewCachedTransformerWithSpec(delegate, cachedTransformer_DEFAULT_CACHE_SPEC)
}

// Java: CachedTransformer.transform(Source)
func (t *CachedTransformer) Transform(input *Source) *Source {
	return t.cache.GetUnchecked(input)
}
