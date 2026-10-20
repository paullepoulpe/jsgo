# Genuine BigQuery v2 generated-code fixture

This fixture contains the complete `google/cloud/bigquery/v2` schema family
from public Googleapis revision `ba591f1d7fdedc610639425448756cd8a64adc0f`:
48 upstream family schemas plus their 13 actual transitive imports. Files remain
byte-identical. Googleapis uses Apache-2.0; the matching protobuf sources/runtime
retain their upstream license. No schemas, fields, runtime classes or externs
are invented, deleted or patched.

`generated/googleapis_bigquery_pin.json` freezes every family file's published
git blob, length and source revision. Fetch records SHA-256 for the entire
transitive source closure. Well-known protobuf imports come from the already
pinned genuine protobuf 3.20.3 source archive, not a replacement implementation.

Preparation uses these separate commands:

1. `googleapis_family_fixture.py fetch --out <new-source> --protobuf-archive <pinned-source-tar>`
2. Run `generate --source <source> --protoc <pinned-native-protoc> --out <new-generation>` through the existing capped setup launcher.
3. Run `wire-oracle` with the same three arguments in a new output directory through that launcher.
4. `finalize --source <source> --generation <generation> --wire <oracle> --runtime <validated-protobuf-fixture> --out <new-fixture>`
5. `specs --fixture <fixture> --out <new-specs> --cache <shared-cache> --jar <pinned-stock-jar> --threads 1`

The native generator and wire-reference stages require a real 4 GiB cgroup,
zero swap and a two-CPU quota/affinity, and record the exact executable SHA.
Downloads and generation never enter compiler timing. Current preparation
successfully generates 227 genuine Closure-style JS files, about 3.75 MB, from
61 schemas without version suppression or source edits.

The small owned harness exports a catalog of 349 actual BigQuery API types and
exercises real Table/Job payloads: table/schema references, repeated fields,
query parameters, present false-valued wrappers and Unicode text. Catalog
entries reference actual generated types rather than duplicate their code.
This is a coherent generated API model library, not a production application
or a BigQuery network request.

Independent native protobuf encodes benchmark-owned text payloads using these
unchanged schemas. The generated JS must decode those native payloads, preserve
all data through serialization and return the expected semantic result.
`bigquery_native_crosscheck.py` then independently decodes the actual JS-produced
bytes with native protobuf and compares the complete canonical text payloads.
Field order in protobuf wire output is not a semantic invariant: the retained
first attempt compared bytes too strictly, while the corrected independent
native decode confirms equivalent complete messages.

Before any performance claim, all of these gates remain required:

- Fresh unchanged stock Java and Go compile the same complete arguments and inputs under the same real cap; every stream, exit, JavaScript/map and manifest byte agrees.
- The custom `bigquery_compiled_check.js` runs independently on both compiled outputs.
- Native protobuf crosschecks both of those actual compiled-output wire records.
- Actual compiler manifests report the selected source and generated-family work. Supplied candidates are not presented as already measured compiler work.

The spec deliberately uses ordinary ADVANCED/PRUNE/type-check settings. It
currently supplies 1,209 JS candidates (about 21.6 MB including the real Closure
library) plus 61 schema provenance inputs. Actual manifest counts and compiler
compatibility remain pending the stock Java/Go gates. Normal pruning must be
reported, not hidden. The general enterprise gate queue rejects this spec's
additional native-crosscheck requirement so a partial workflow cannot silently
claim complete semantic validation.
