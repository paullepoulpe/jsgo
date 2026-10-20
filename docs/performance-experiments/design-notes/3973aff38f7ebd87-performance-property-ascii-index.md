# PropertyMap ASCII point lookup experiment

`JSGO_EXPERIMENT_PROPERTY_ASCII_INDEX=1` enables a private auxiliary lookup
index only on maps created by `NewPropertyMap` with their own stock
`String.compareTo` backing. The original TreeMap remains the source of sorted
iteration and mutation behavior. Maps provided to the backing-map constructor,
including custom comparators and maps with aliases, remain on the original path.

After four ASCII lookups while the map contains at least16 properties, a native
Go string index is built for its ASCII keys. Non-ASCII and malformed/WTF-8 keys
always use the TreeMap so Java UTF16 equality and ordering remain exact. All
production string writes occur through PutProperty; successful writes update the
shared index after the original TreeMap write. Symbol writes are independent.
There is no public mutable raw-map or entry view; Values is a live value iterator
and callbacks that overwrite via PutProperty update both representations.

An independent shared sidecar preserves Go shallow-copy aliases. It contains a
backing identity guard and atomic lazy publication so previously safe concurrent
readers remain safe. The original map's unsupported concurrent read/write case
is not broadened. The sidecar never points back at PropertyMap: changing a copied
map's parent must not retain an obsolete original parent/cache graph.

Memory is a deliberate tradeoff. The extra PropertyMap pointer changes its
source-derived amd64 size88->96, within the same96-byte allocation class. There
is no TreeMap layout change. Opted-in maps allocate a separate24-byte sidecar;
only hot maps allocate the native index plus an8-byte wrapper. The experiment
must be measured for wall time, allocations, live heap and peak RSS before it is
retained. Class sizes are source-derived until the candidate is compiled.

Focused differential tests compare the original TreeMap lookup against the
index under random overwrites, Unicode/WTF-8 equivalent keys, missing/nil values,
ordered/live views, custom comparator calls, aliased removal/clear, shallow copies,
JSDoc preservation and immutable exceptions. Concurrent-reader race tests cover
lazy publication; a weak-reference GC test covers obsolete-parent retention.

Fresh labeled eecbcf44 profiles retain0.34s of checks CPU under property point
lookups (4.3% whole), and0.25s typed (2.1%). Those are sampled attribution, not a
prediction that this whole cost disappears. No performance result has yet been
measured for this prototype.
